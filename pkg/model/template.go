package model

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/kavimaluskam/leetcode-cli/pkg/utils"
)

// Template is the config of leetcode stored in local
type Template struct {
	DirTemplate        string `json:"dirTemplate"`
	MarkdownTemplate   string `json:"markdownTemplate"`
	SourceCodeTemplate string `json:"sourceCodeTemplate"`
}

// GetFileTemplate returns a basic API template struct based on local template config
func GetFileTemplate(pd ProblemDetail) (*Template, error) {
	t := Template{}

	file, err := ioutil.ReadFile(utils.TemplateConfigPath)
	if err != nil {
		return &t, err
	}

	err = json.Unmarshal([]byte(file), &t)
	if err != nil {
		return &t, err
	}

	if t.DirTemplate != "" {
		t.DirTemplate = strings.ReplaceAll(t.DirTemplate, "$questionID", pd.QuestionFrontendID)
		t.DirTemplate = strings.ReplaceAll(t.DirTemplate, "$questionSlug", pd.TitleSlug)
	}

	if t.MarkdownTemplate != "" {
		t.MarkdownTemplate = strings.ReplaceAll(t.MarkdownTemplate, "$questionID", pd.QuestionFrontendID)
		t.MarkdownTemplate = strings.ReplaceAll(t.MarkdownTemplate, "$questionSlug", pd.TitleSlug)
	}

	if t.SourceCodeTemplate != "" {
		t.SourceCodeTemplate = strings.ReplaceAll(t.SourceCodeTemplate, "$questionID", pd.QuestionFrontendID)
		t.SourceCodeTemplate = strings.ReplaceAll(t.SourceCodeTemplate, "$questionSlug", pd.TitleSlug)
		t.SourceCodeTemplate = strings.ReplaceAll(t.SourceCodeTemplate, "$submissionID", "1")
	}

	return &t, nil
}
