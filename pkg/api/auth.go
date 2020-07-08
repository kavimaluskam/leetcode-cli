package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/kavimaluskam/leetcode-cli/pkg/utils"
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
	file, err := ioutil.ReadFile(utils.AuthConfigPath)
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
		return fmt.Errorf("Error on processing authencation json: %s", err.Error())
	}

	err = ioutil.WriteFile(utils.AuthConfigPath, file, os.ModePerm)
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
	browser := rod.New().Connect()

	// Timeout will be passed to all chained function calls.
	// The code will panic out if any chained call is used after the timeout.
	page := browser.Timeout(30 * time.Second).Page(utils.LoginURL)

	// Resize the window make sure window size is always consistent.
	page.Window(0, 0, 1200, 600)

	page.Element("#signin_btn").WaitVisible()

	page.Element("#id_login").Input(a.Username)
	page.Element("#id_password").Input(a.Password).Press(input.Enter)

	var err error
	select {
	case err = <-goWaitVisible(
		page,
		".error-message__27FL",
		"failed signing into leetcode. incorrect auth credentials",
	):
	case err = <-goWaitVisible(page, "#nav-user-app", ""):
	}

	if err != nil {
		return err
	}

	for _, cookie := range page.Cookies() {
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
		page.Element(selector).WaitVisible()
		if errMsg != "" {
			e <- fmt.Errorf(errMsg)
		} else {
			e <- nil
		}
	}()
	return e
}
