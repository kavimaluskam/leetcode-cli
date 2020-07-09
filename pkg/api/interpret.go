package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/kavimaluskam/leetcode-cli/pkg/model"
	"github.com/kavimaluskam/leetcode-cli/pkg/utils"
)

type interpretInitResp struct {
	InterpretID string `json:"interpret_id"`
	TestCase    string `json:"test_case"`
}

type interpretResp struct {
	State                  string   `json:"state"`
	CodeAnswer             []string `json:"code_answer"`               //: ["[0,1]"]
	CodeOutput             []string `json:"code_output"`               //: []
	CorrectAnswer          bool     `json:"correct_answer"`            //: true
	ElapsedTime            int      `json:"elapsed_time"`              //: 70
	ExpectedCodeAnswer     []string `json:"expected_code_answer"`      //: ["[0,1]"]
	ExpectedCodeOutput     []string `json:"expected_code_output"`      //: []
	ExpectedElapsedTime    int      `json:"expected_elapsed_time"`     //: 14
	ExpectedLang           string   `json:"expected_lang"`             //: "cpp"
	ExpectedMemory         int      `json:"expected_memory"`           //: 6636000
	ExpectedRunSuccess     bool     `json:"expected_run_success"`      //: true
	ExpectedStatusCode     int      `json:"expected_status_code"`      //: 10
	ExpectedStatusRuntime  string   `json:"expected_status_runtime"`   //: "0"
	ExpectedTaskFinishTime int      `json:"expected_task_finish_time"` //: 1594210446118
	FullRuntimeError       string   `json:"full_runtime_error"`        //: "IndentationError: unexpected indent↵    ^↵    diff_key = nums_dict.get(diff)↵Line 8  (Solution.py)"
	Lang                   string   `json:"lang"`                      //: "python3"
	Memory                 int      `json:"memory"`                    //: 13684000
	MemoryPercentile       string   `json:"memory_percentile"`         //: null
	PrettyLang             string   `json:"pretty_lang"`               //: "Python3"
	RuntimeError           string   `json:"runtime_error"`             //: "Line 8: IndentationError: unexpected indent"
	RunSuccess             bool     `json:"run_success"`               //: true
	RuntimePercentile      string   `json:"runtime_percentile"`        //: null
	StatusCode             int      `json:"status_code"`               //: 10
	StatusMemory           string   `json:"status_memory"`             //: "13.7 MB"
	StatusMsg              string   `json:"status_msg"`                //: "Accepted"
	StatusRuntime          string   `json:"status_runtime"`            //: "40 ms"
	SubmissionID           string   `json:"submission_id"`             //: "runcode_1594212035.1193292_VvybRjW7Lk"
	TaskFinishTime         int      `json:"task_finish_time"`          //: 1594212037297
	TotalCorrect           string   `json:"total_correct"`             //: null
	TotalTestcases         string   `json:"total_testcases"`           //: null
}

func (c *Client) InterpretCode(pd *model.ProblemDetail, fp string, dataInput string) error {
	ext := filepath.Ext(fp)
	lang, err := pd.GetLanguageSlug(ext)
	if err != nil {
		return err
	}

	file, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}

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

	ir := &interpretInitResp{}
	err = c.REST("POST", url, bytes.NewBuffer(reqBody), ir)
	if err != nil {
		return err
	}

	for {
		vr, err := c.verifyInterpretation(ir.InterpretID)
		if err != nil {
			return err
		}
		switch vr.State {
		case "PENDING", "STARTED":
		case "SUCCESS":
			// TODO: handle resp properly
			fmt.Printf("%+v\n", vr)
			return nil
		default:
			return fmt.Errorf("failure code submission. unexpected submission state: %s", vr.State)
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
