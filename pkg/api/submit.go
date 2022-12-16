package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ckidckidckid/leetcode-cli/pkg/model"
	"github.com/ckidckidckid/leetcode-cli/pkg/utils"
	"github.com/kyokomi/emoji"
)

type submitInitResp struct {
	SubmissionID int `json:"submission_id"`
}

type submitResp struct {
	State             string  `json:"state"`
	CodeOutput        string  `json:"code_output"`
	CompareResult     string  `json:"compare_result"`
	ElapsedTime       int     `json:"elapsed_time"`
	ExpectedOutput    string  `json:"expected_output"`
	FullRuntimeError  string  `json:"full_runtime_error"`
	Lang              string  `json:"lang"`
	LastTestcase      string  `json:"last_testcase"`
	Memory            int     `json:"memory"`
	MemoryPercentile  float32 `json:"memory_percentile"`
	PrettyLang        string  `json:"pretty_lang"`
	QuestionID        string  `json:"question_id"`
	RunSuccess        bool    `json:"run_success"`
	RuntimeError      string  `json:"runtime_error"`
	RuntimePercentile float32 `json:"runtime_percentile"`
	StatusCode        int     `json:"status_code"`
	StatusMemory      string  `json:"status_memory"`
	StatusMsg         string  `json:"status_msg"`
	StatusRuntime     string  `json:"status_runtime"`
	StdOutput         string  `json:"std_output"`
	SubmissionID      string  `json:"submission_id"`
	TaskFinishTime    int     `json:"task_finish_time"`
	TotalCorrect      int     `json:"total_correct"`
	TotalTestcases    int     `json:"total_testcases"`
}

// SubmitCode to leetcode judge
func (c *Client) SubmitCode(pd *model.ProblemDetail, fp string) error {
	ext := filepath.Ext(fp)
	lang, err := pd.GetLanguageSlug(ext)
	if err != nil {
		return err
	}

	file, err := os.ReadFile(fp)
	if err != nil {
		return err
	}

	url := strings.Replace(utils.SubmitURL, "$slug", pd.TitleSlug, 1)

	reqBody, err := json.Marshal(
		map[string]interface{}{
			"lang":        lang,
			"question_id": pd.QuestionID,
			"typed_code":  string(file),
		},
	)
	if err != nil {
		return err
	}

	sr := &submitInitResp{}
	err = c.REST("POST", url, bytes.NewBuffer(reqBody), sr)
	if err != nil {
		return err
	}

	for {
		vr, err := c.verifySubmission(sr.SubmissionID)
		if err != nil {
			return err
		}
		switch vr.State {
		case "PENDING", "STARTED":
		case "SUCCESS":
			vr.exportSdtoutSubmission()
			return nil
		default:
			return fmt.Errorf("failure code submission. unexpected submission state: %s", vr.State)
		}
		time.Sleep(2 * time.Second)
	}
}

func (c *Client) verifySubmission(id int) (*submitResp, error) {
	idstr := fmt.Sprintf("%d", id)
	url := strings.Replace(utils.VerifyURL, "$id", idstr, 1)
	vr := &submitResp{}
	err := c.REST("GET", url, nil, vr)
	if err != nil {
		return nil, err
	}
	return vr, nil
}

func (vr *submitResp) exportSdtoutSubmission() {
	if vr.StatusMsg == "Accepted" {
		emoji.Printf("%s :heavy_check_mark:\n", utils.Green("Accepted"))
		fmt.Printf("%d/%d test cases passed\n\n", vr.TotalCorrect, vr.TotalTestcases)
		fmt.Printf("%s\n", utils.Blue("Runtime"))
		fmt.Printf("%s, faster than %.2f%% submissions\n\n", vr.StatusRuntime, vr.RuntimePercentile)
		fmt.Printf("%s\n", utils.Blue("Memory"))
		fmt.Printf("%s, less than %.2f%% submissions\n", vr.StatusMemory, vr.MemoryPercentile)
	} else {
		emoji.Printf("%s :x:\n", utils.Red("Rejected"))
		fmt.Printf("%d/%d test cases passed\n\n", vr.TotalCorrect, vr.TotalTestcases)
		fmt.Printf(
			"%s\n%s",
			utils.Cyan("Last Test Case"),
			fmt.Sprintf("%s\n\n", strings.ReplaceAll(vr.LastTestcase, "\n", "\\n")),
		)

		if vr.FullRuntimeError != "" {
			fmt.Printf("%s\n%s\n", utils.Red("Runtime Error"), utils.Magenta(vr.FullRuntimeError))
		} else {
			fmt.Printf("%s\n", utils.Red("Wrong Answer"))
			fmt.Printf("Expected   %s\n", vr.ExpectedOutput)
			fmt.Printf("Actual     %s\n", vr.CodeOutput)
		}
	}
}
