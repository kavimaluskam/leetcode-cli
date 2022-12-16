package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ckidckidckid/leetcode-cli/pkg/api"
	"github.com/ckidckidckid/leetcode-cli/pkg/utils"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

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
	Args:  cobra.NoArgs,
	RunE:  userSignIn,
}

var userSignOutCmd = &cobra.Command{
	Use:   `signout`,
	Short: `Sign out from leetcode on cli`,
	Args:  cobra.NoArgs,
}

func userSignIn(cmd *cobra.Command, args []string) (err error) {
	var username string
	var passwordStr string
	auth, err := api.GetAuthCredentials()

	if auth.Username == "" || auth.Password == "" || err != nil {
		fmt.Println(
			utils.Gray("Cannot read auth credentials in local config, asking for manual input"),
		)
		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("Please enter your username: ")
		username, err = reader.ReadString('\n')
		if err != nil {
			return err
		}
		username = strings.TrimSuffix(username, "\n")

		fmt.Printf("Please enter your password: \n")
		password, err := terminal.ReadPassword(0)
		if err != nil {
			return err
		}
		passwordStr = fmt.Sprintf("%s", password)
	} else {
		username = auth.Username
		passwordStr = auth.Password
	}

	a := api.Auth{
		Username: username,
		Password: passwordStr,
	}
	fmt.Printf("Loggin into leetcode as %s...\n", username)

	err = a.Login()
	if err != nil {
		return err
	}

	err = a.SetAuthCredentials()
	if err != nil {
		return err
	}

	fmt.Printf("Successfully signed in as %s\n", username)
	return nil
}
