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

type submissionResp struct {
	SubmissionID int `json:"submission_id"`
}

type submitResp struct {
	State             string  `json:"state"`
	CodeOutput        string  `json:"code_output"`
	CompareResult     string  `json:"compare_result"`
	ElapsedTime       int     `json:"elapsed_time"`
	Lang              string  `json:"lang"`
	LastTestcase      string  `json:"last_testcase"`
	Memory            int     `json:"memory"`
	MemoryPercentile  float32 `json:"memory_percentile"`
	PrettyLang        string  `json:"pretty_lang"`
	QuestionID        string  `json:"question_id"`
	RunSuccess        bool    `json:"run_success"`
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

func (c *Client) SubmitCode(pd *model.ProblemDetail, fp string) error {
	ext := filepath.Ext(fp)
	lang, err := pd.GetLanguageSlug(ext)
	if err != nil {
		return err
	}

	file, err := ioutil.ReadFile(fp)
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

	sr := &submissionResp{}
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
		case "PENDING":
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

type interpretResp struct {
	State                  string   `json:"state"`
	CodeAnswer             []string `json:"code_answer"`               //: ["[0,1]"]
	CodeOutput             []string `json:"code_output"`               //: []
	CorrectAnswer          bool     `json:"correct_answer"`            //: true
	ElapsedTime            int      `json:"elapsed_time"`              //: 70
	ExpectedCodeAnswer     string   `json:"expected_code_answer"`      //: ["[0,1]"]
	ExpectedCodeOutput     string   `json:"expected_code_output"`      //: []
	ExpectedElapsedTime    string   `json:"expected_elapsed_time"`     //: 14
	ExpectedLang           string   `json:"expected_lang"`             //: "cpp"
	ExpectedMemory         string   `json:"expected_memory"`           //: 6636000
	ExpectedRunSuccess     string   `json:"expected_run_success"`      //: true
	ExpectedStatusCode     string   `json:"expected_status_code"`      //: 10
	ExpectedStatusRuntime  string   `json:"expected_status_runtime"`   //: "0"
	ExpectedTaskFinishTime string   `json:"expected_task_finish_time"` //: 1594210446118
	Lang                   string   `json:"lang"`                      //: "python3"
	Memory                 string   `json:"memory"`                    //: 13684000
	MemoryPercentile       string   `json:"memory_percentile"`         //: null
	PrettyLang             string   `json:"pretty_lang"`               //: "Python3"
	RunSuccess             string   `json:"run_success"`               //: true
	RuntimePercentile      string   `json:"runtime_percentile"`        //: null
	StatusCode             string   `json:"status_code"`               //: 10
	StatusMemory           string   `json:"status_memory"`             //: "13.7 MB"
	StatusMsg              string   `json:"status_msg"`                //: "Accepted"
	StatusRuntime          string   `json:"status_runtime"`            //: "40 ms"
	SubmissionID           string   `json:"submission_id"`             //: "runcode_1594212035.1193292_VvybRjW7Lk"
	TaskFinishTime         string   `json:"task_finish_time"`          //: 1594212037297
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

	sr := &submissionResp{}
	err = c.REST("POST", url, bytes.NewBuffer(reqBody), sr)
	if err != nil {
		return err
	}

	for {
		vr, err := c.verifyInterpretation(sr.SubmissionID)
		if err != nil {
			return err
		}
		switch vr.State {
		case "PENDING":
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

func (c *Client) verifyInterpretation(id int) (*interpretResp, error) {
	idstr := fmt.Sprintf("%d", id)
	url := strings.Replace(utils.VerifyURL, "$id", idstr, 1)
	ir := &interpretResp{}
	err := c.REST("GET", url, nil, ir)
	if err != nil {
		return nil, err
	}
	return ir, nil
}
