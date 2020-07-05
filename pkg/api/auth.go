package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/kavimaluskam/leetcode-cli/pkg/utils"
)

// Auth is the config of leetcode stored in local
type Auth struct {
	Login       string `json:"login"`
	LoginCSRF   string `json:"loginCSRF"`
	SessionCSRF string `json:"sessionCSRF"`
	SessionID   string `json:"sessionId"`
}

// GetAuthClient returns a basic API Client based on local auth config
func GetAuthClient() (*Client, error) {
	file, _ := ioutil.ReadFile(utils.AuthConfigPath)
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

	return NewClient(opts...), nil
}
