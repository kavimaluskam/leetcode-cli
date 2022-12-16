package model

import (
	"encoding/json"
	"os"
	"strings"
	"text/template"

	"github.com/ckidckidckid/leetcode-cli/pkg/utils"
)

// Template is the config of leetcode stored in local
type FileTemplate struct {
	MarkdownPath     string `json:"markDownPath"`
	SourceCodePath   string `json:"sourceCodePath"`
	MarkdownTemplate *template.Template
}

// GetFileTemplate returns a basic API template struct based on local template config
func GetFileTemplate(pd ProblemDetail) (*FileTemplate, error) {
	t := FileTemplate{}

	file, err := os.ReadFile(utils.TemplateConfigPath)
	if err != nil {
		return &t, err
	}

	err = json.Unmarshal([]byte(file), &t)
	if err != nil {
		return &t, err
	}

	if t.MarkdownPath != "" {
		t.MarkdownPath = strings.ReplaceAll(t.MarkdownPath, "$questionID", pd.QuestionFrontendID)
		t.MarkdownPath = strings.ReplaceAll(t.MarkdownPath, "$questionSlug", pd.TitleSlug)
	}

	if t.SourceCodePath != "" {
		t.SourceCodePath = strings.ReplaceAll(t.SourceCodePath, "$questionID", pd.QuestionFrontendID)
		t.SourceCodePath = strings.ReplaceAll(t.SourceCodePath, "$questionSlug", pd.TitleSlug)
		t.SourceCodePath = strings.ReplaceAll(t.SourceCodePath, "$submissionID", "1")
	}

	md, err := template.ParseFiles(utils.MarkdownTemplatePath)
	if err != nil {
		return &t, err
	}
	t.MarkdownTemplate = md

	return &t, nil
}
