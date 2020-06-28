package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/joho/godotenv"
)

type auth struct {
	Login       string `json:"login"`
	LoginCSRF   string `json:"loginCSRF"`
	SessionCSRF string `json:"sessionCSRF"`
	SessionID   string `json:"sessionId"`
}

func login(username string, password string) (csrftoken string, sessionID string) {
	// Launch a new browser with default options, and connect to it.
	browser := rod.New().Connect()

	// Even you forget to close, rod will close it after main process ends.
	defer browser.Close()

	// Timeout will be passed to all chained function calls.
	// The code will panic out if any chained call is used after the timeout.
	page := browser.Timeout(time.Minute).Page(`https://leetcode.com/accounts/login/`)

	// Resize the window make sure window size is always consistent.
	page.Window(0, 0, 1200, 600)

	page.Element(`#signin_btn`).WaitVisible()

	page.Element(`#id_login`).Input(username)
	page.Element(`#id_password`).Input(password).Press(input.Enter)

	page.Element(`#nav-user-app`).WaitVisible()

	for _, cookie := range page.Cookies() {
		if cookie.Name == `csrftoken` {
			csrftoken = cookie.Value
		} else if cookie.Name == `LEETCODE_SESSION` {
			sessionID = cookie.Value
		}
	}

	return csrftoken, sessionID
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	username := os.Getenv("LEETCODE_USERNAME")
	password := os.Getenv("LEETCODE_PASSWORD")
	authConfigPath := os.Getenv("LEETCODE_CONFIG_PATH")

	log.Printf("Signing in leetcode...")
	csrftoken, sessionID := login(username, password)

	data := auth{
		Login:       username,
		LoginCSRF:   ``,
		SessionCSRF: csrftoken,
		SessionID:   sessionID,
	}

	file, err := json.Marshal(data)
	if err != nil {
		log.Panicf(`Error on processing authencation json: %s`, err.Error())
		return
	}

	err = ioutil.WriteFile(authConfigPath, file, os.ModePerm)
	if err != nil {
		log.Panicf(`Error on writing authencation json to path %s: %s`, authConfigPath, err.Error())
		return
	}

	log.Printf("Successfully signed in as %s", username)
}
