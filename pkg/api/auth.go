package api

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ckidckidckid/leetcode-cli/pkg/model"
	"github.com/ckidckidckid/leetcode-cli/pkg/utils"
	"github.com/go-rod/rod"
)

// Auth is the config of leetcode stored in local
type Auth struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	SessionCSRF string `json:"sessionCSRF"`
	SessionID   string `json:"sessionId"`
}

// GetAuthCredentials retrieve auth information from local config
func GetAuthCredentials() (*Auth, error) {
	file, err := os.ReadFile(utils.AuthConfigPath)
	if err != nil {
		return nil, err
	}

	a := Auth{}
	err = json.Unmarshal([]byte(file), &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// SetAuthCredentials update auth information to local config
func (a *Auth) SetAuthCredentials() error {
	file, err := json.Marshal(a)
	if err != nil {
		// TODO: enhance error type handling
		return fmt.Errorf("Error on processing authentication json: %s", err.Error())
	}

	err = os.WriteFile(utils.AuthConfigPath, file, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// GetAuthClient returns a basic API Client based on local auth config
func GetAuthClient() (*Client, error) {
	a, err := GetAuthCredentials()
	if err != nil {
		return nil, err
	}

	var opts []ClientOption

	opts = append(
		opts,
		AddHeader(
			"Cookie",
			fmt.Sprintf(
				"LEETCODE_SESSION=%s;csrftoken=%s;",
				a.SessionID,
				a.SessionCSRF,
			),
		),
	)

	opts = append(
		opts,
		AddHeader("X-Requested-With", "XMLHttpRequest"),
	)

	opts = append(
		opts,
		AddHeader("X-CSRFToken", a.SessionCSRF),
	)

	return NewClient(opts...), nil
}

// Login to leetcode with Rod headless browser
func (a *Auth) Login() error {
	// Launch a new browser with default options, and connect to it.
	browser := rod.New().MustConnect()

	defer browser.MustClose()

	page := browser.MustPage(utils.LoginURL)
	page.MustWaitLoad()

	// Resize the window make sure window size is always consistent.
	page.MustSetWindow(0, 0, 1200, 600)

	signInBtn := page.MustElement("#signin_btn").MustWaitVisible()

	idLogin := page.MustElement("#id_login").MustWaitVisible()
	idLogin.Input(a.Username)

	password := page.MustElement("#id_password").MustWaitVisible()
	password.Input(a.Password)

	signInBtn.MustClick()

	var err error
	select {
	case err = <-goWaitVisible(
		page,
		".error-message__27FL",
		"failed signing into leetcode. incorrect auth credentials",
	):
	case err = <-goWaitVisible(page, "#navbar-right-container", ""):
	}

	if err != nil {
		return err
	}

	cookies, _ := page.Cookies([]string{utils.BaseURL})
	for _, cookie := range cookies {
		if cookie.Name == "csrftoken" {
			a.SessionCSRF = cookie.Value
		} else if cookie.Name == "LEETCODE_SESSION" {
			a.SessionID = cookie.Value
		}
	}

	return nil
}

func goWaitVisible(page *rod.Page, selector string, errMsg string) <-chan error {
	e := make(chan error)
	go func() {
		defer close(e)
		page.MustElement(selector).MustWaitVisible()
		if errMsg != "" {
			e <- fmt.Errorf(errMsg)
		} else {
			e <- nil
		}
	}()
	return e
}

func GetSubmitClient(pd *model.ProblemDetail) (*Client, error) {
	file, _ := os.ReadFile(utils.AuthConfigPath)
	a := Auth{}
	_ = json.Unmarshal([]byte(file), &a)

	var opts []ClientOption

	opts = append(
		opts,
		AddHeader(
			"Cookie",
			fmt.Sprintf(
				"LEETCODE_SESSION=%s;csrftoken=%s;",
				a.SessionID,
				a.SessionCSRF,
			),
		),
	)

	opts = append(
		opts,
		AddHeader("X-Requested-With", "XMLHttpRequest"),
	)

	opts = append(
		opts,
		AddHeader("X-CSRFToken", a.SessionCSRF),
	)

	opts = append(
		opts,
		AddHeader("Origin", utils.BaseURL),
	)

	opts = append(
		opts,
		AddHeader("Referer", strings.Replace(utils.SubmitRefererURL, "$slug", pd.TitleSlug, 1)),
	)

	return NewClient(opts...), nil
}
