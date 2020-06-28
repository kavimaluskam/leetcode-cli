package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/kavimaluskam/leetcode-cli/pkg/cmd/util"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var authConfigPath = fmt.Sprintf("%s/.lc/leetcode/user.json", os.Getenv("HOME"))

func init() {
	RootCmd.AddCommand(userCmd)
	userCmd.AddCommand(userSignInCmd)
	userCmd.AddCommand(userSignOutCmd)
}

var userCmd = &cobra.Command{
	Use:   `user <commands>`,
	Short: `Sign In, Sing Out on cli`,
	Long:  `Work with user authentication`,
}

var userSignInCmd = &cobra.Command{
	Use:   `signin`,
	Short: `Sign in to leetcode on cli`,
	Args:  util.NoArgsQuoteReminder,
	RunE:  userSignIn,
}

var userSignOutCmd = &cobra.Command{
	Use:   `signout`,
	Short: `Sign out from leetcode on cli`,
	Args:  util.NoArgsQuoteReminder,
}

func userSignIn(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Please enter your username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSuffix(username, "\n")

	fmt.Printf("Please enter your password: \n")
	password, _ := terminal.ReadPassword(0)
	passwordStr := fmt.Sprintf("%s", password)

	fmt.Printf("Loggin into leetcode as %s...\n", username)
	csrfToken, sessionID := login(username, passwordStr)

	data := auth{
		Login:       username,
		LoginCSRF:   "",
		SessionCSRF: csrfToken,
		SessionID:   sessionID,
	}

	file, err := json.Marshal(data)
	if err != nil {
		fmt.Errorf("Error on processing authencation json: %s", err.Error())
		return err
	}

	err = ioutil.WriteFile(authConfigPath, file, os.ModePerm)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully signed in as %s\n", username)
	return nil
}

func userSignOut(cmd *cobra.Command, args []string) error {
	// TODO: add login implementation
	fmt.Println("User sign out!")

	return nil
}

type auth struct {
	Login       string `json:"login"`
	LoginCSRF   string `json:"loginCSRF"`
	SessionCSRF string `json:"sessionCSRF"`
	SessionID   string `json:"sessionId"`
}

func login(username string, password string) (csrfToken string, sessionID string) {
	// Launch a new browser with default options, and connect to it.
	browser := rod.New().Connect()

	// Even you forget to close, rod will close it after main process ends.
	defer browser.Close()

	// Timeout will be passed to all chained function calls.
	// The code will panic out if any chained call is used after the timeout.
	page := browser.Timeout(time.Minute).Page("https://leetcode.com/accounts/login/")

	// Resize the window make sure window size is always consistent.
	page.Window(0, 0, 1200, 600)

	page.Element("#signin_btn").WaitVisible()

	page.Element("#id_login").Input(username)
	page.Element("#id_password").Input(password).Press(input.Enter)

	page.Element("#nav-user-app").WaitVisible()

	for _, cookie := range page.Cookies() {
		if cookie.Name == "csrftoken" {
			csrfToken = cookie.Value
		} else if cookie.Name == "LEETCODE_SESSION" {
			sessionID = cookie.Value
		}
	}

	return csrfToken, sessionID
}
