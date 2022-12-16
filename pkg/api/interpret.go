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

type interpretInitResp struct {
	InterpretID string `json:"interpret_id"`
	TestCase    string `json:"test_case"`
}

type interpretResp struct {
	State                  string   `json:"state"`
	CodeAnswer             []string `json:"code_answer"`
	CodeOutput             []string `json:"code_output"`
	CorrectAnswer          bool     `json:"correct_answer"`
	ElapsedTime            int      `json:"elapsed_time"`
	ExpectedCodeAnswer     []string `json:"expected_code_answer"`
	ExpectedCodeOutput     []string `json:"expected_code_output"`
	ExpectedElapsedTime    int      `json:"expected_elapsed_time"`
	ExpectedLang           string   `json:"expected_lang"`
	ExpectedMemory         int      `json:"expected_memory"`
	ExpectedRunSuccess     bool     `json:"expected_run_success"`
	ExpectedStatusCode     int      `json:"expected_status_code"`
	ExpectedStatusRuntime  string   `json:"expected_status_runtime"`
	ExpectedTaskFinishTime int      `json:"expected_task_finish_time"`
	FullRuntimeError       string   `json:"full_runtime_error"`
	Lang                   string   `json:"lang"`
	Memory                 int      `json:"memory"`
	MemoryPercentile       string   `json:"memory_percentile"`
	PrettyLang             string   `json:"pretty_lang"`
	RuntimeError           string   `json:"runtime_error"`
	RunSuccess             bool     `json:"run_success"`
	RuntimePercentile      string   `json:"runtime_percentile"`
	StatusCode             int      `json:"status_code"`
	StatusMemory           string   `json:"status_memory"`
	StatusMsg              string   `json:"status_msg"`
	StatusRuntime          string   `json:"status_runtime"`
	SubmissionID           string   `json:"submission_id"`
	TaskFinishTime         int      `json:"task_finish_time"`
	TotalCorrect           string   `json:"total_correct"`
	TotalTestcases         string   `json:"total_testcases"`
}

// InterpretCode with leetcode judge and input testcase
func (c *Client) InterpretCode(pd *model.ProblemDetail, fp string, dataInput string) error {
	ext := filepath.Ext(fp)
	lang, err := pd.GetLanguageSlug(ext)
	if err != nil {
		return err
	}

	file, err := os.ReadFile(fp)
	if err != nil {
		return err
	}

	dataInput = strings.ReplaceAll(dataInput, "\\n", "\n")

	url := strings.Replace(utils.InterpretURL, "$slug", pd.TitleSlug, 1)

	reqBody, err := json.Marshal(
		map[string]interface{}{
			"data_input":  dataInput,
			"judge_type":  "large",
			"lang":        lang,
			"question_id": pd.QuestionID,
			"typed_code":  string(file),
		},
	)
	if err != nil {
		return err
	}

	iir := &interpretInitResp{}
	err = c.REST("POST", url, bytes.NewBuffer(reqBody), iir)
	if err != nil {
		return err
	}

	for {
		ir, err := c.verifyInterpretation(iir.InterpretID)
		if err != nil {
			return err
		}
		switch ir.State {
		case "PENDING", "STARTED":
		case "SUCCESS":
			ir.exportSdtoutInterpretation(dataInput)
			return nil
		default:
			return fmt.Errorf("failure code submission. unexpected submission state: %s", ir.State)
		}
		time.Sleep(2 * time.Second)
	}
}

func (c *Client) verifyInterpretation(id string) (*interpretResp, error) {
	url := strings.Replace(utils.VerifyURL, "$id", id, 1)
	ir := &interpretResp{}
	err := c.REST("GET", url, nil, ir)
	if err != nil {
		return nil, err
	}
	return ir, nil
}

func (ir *interpretResp) exportSdtoutInterpretation(t string) {
	if ir.CorrectAnswer {
		emoji.Printf("%s :heavy_check_mark:\n\n", utils.Green("Accepted"))
	} else {
		emoji.Printf("%s :x:\n\n", utils.Red("Rejected"))
	}

	fmt.Printf(
		"%s\n%s",
		utils.Cyan("Test Case"),
		fmt.Sprintf("%s\n\n", strings.ReplaceAll(t, "\n", "\\n")),
	)

	if ir.FullRuntimeError != "" {
		fmt.Printf("%s\n%s\n", utils.Red("Runtime Error"), utils.Magenta(ir.FullRuntimeError))
	} else {
		fmt.Printf("%s\n", utils.Blue("Answer"))
		fmt.Printf("Expected: %s\n", strings.Join(ir.ExpectedCodeAnswer, ", "))
		fmt.Printf("Actual:   %s\n\n", strings.Join(ir.CodeAnswer, ", "))
		fmt.Printf("%s\n", utils.Blue("Runtime"))
		fmt.Printf("Expected: %s ms\n", ir.ExpectedStatusRuntime)
		fmt.Printf("Actual:   %s\n\n", ir.StatusRuntime)
		fmt.Printf("%s\n", utils.Blue("Memory"))
		fmt.Printf("Expected: %.2f MB\n", float32(ir.ExpectedMemory)/float32(1024)/float32(1024))
		fmt.Printf("Actual:   %s\n", ir.StatusMemory)
	}
}
