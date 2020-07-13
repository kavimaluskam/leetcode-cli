package model

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kavimaluskam/leetcode-cli/pkg/utils"
)

// ProblemDetail is the response from leetcode GraphQL API
// concerning individual problem detail
type ProblemDetail struct {
	QuestionID            string                `json:"questionId"`
	QuestionFrontendID    string                `json:"questionFrontendId"`
	BoundTopicID          string                `json:"boundTopicId"`
	Title                 string                `json:"title"`
	TitleSlug             string                `json:"titleSlug"`
	Content               string                `json:"content"`
	TranslatedTitle       string                `json:"translatedTitle"`
	TranslatedContent     string                `json:"translatedContent"`
	IsPaidOnly            bool                  `json:"isPaidOnly"`
	Difficulty            string                `json:"difficulty"`
	Likes                 int                   `json:"likes"`
	Dislikes              int                   `json:"dislikes"`
	IsLiked               bool                  `json:"isLiked"`
	SimilarQuestions      string                `json:"similarQuestions"`
	Contributors          []ProblemContributor  `json:"contributors"`
	LangToValidPlayground string                `json:"langToValidPlayground"`
	TopicTags             []ProblemTag          `json:"topicTags"`
	CompanyTagStats       string                `json:"companyTagStats"`
	CodeSnippets          []ProblemCodeSnippets `json:"codeSnippets"`
	Stats                 string                `json:"stats"`
	Hints                 []string              `json:"hints"`
	Solution              ProblemSolution       `json:"solution"`
	Status                string                `json:"status"`
	SampleTestCase        string                `json:"sampleTestCase"`
	MetaData              string                `json:"metaData"`
	JudgerAvailable       bool                  `json:"judgerAvailable"`
	JudgeType             string                `json:"judgeType"`
	MySQLSchemas          []string              `json:"mysqlSchemas"`
	EnableRuneCode        bool                  `json:"enableRunCode"`
	EnableTestMode        bool                  `json:"enableTestMode"`
	EnableDebugger        bool                  `json:"enableDebugger"`
	EnvInfo               string                `json:"envInfo"`
	LibraryURL            string                `json:"libraryUrl"`
	AdminURL              string                `json:"adminUrl"`
	TypeName              string                `json:"__typename"`
}

// ProblemContributor is the response from leetcode GraphQL API
// concering problem contributor
type ProblemContributor struct {
	Username   string `json:"username"`
	ProfileURL string `json:"profileUrl"`
	AvatarURL  string `json:"avatarUrl"`
	TypeName   string `json:"__typename"`
}

// ProblemTag is the response from leetcode GraphQL API
// concerning problem tagging
type ProblemTag struct {
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	TranslatedName string `json:"translatedName"`
	TypeName       string `json:"__typename"`
}

// ProblemCodeSnippets is the response from leetcode GraphQL API
// concerning problem codes
type ProblemCodeSnippets struct {
	Lang     string `json:"lang"`
	LangSlug string `json:"langSlug"`
	Code     string `json:"code"`
	TypeName string `json:"__typename"`
}

// GetLanguageExt is a mapper function mapping question lang slug to file ext
func (pcs ProblemCodeSnippets) GetLanguageExt() string {
	switch pcs.LangSlug {
	case "cpp":
		return "cpp"
	case "java":
		return "java"
	case "python":
		return "py"
	case "python3":
		return "py3"
	case "c":
		return "c"
	case "csharp":
		return "cs"
	case "javascript":
		return "js"
	case "ruby":
		return "rb"
	case "swift":
		return "swift"
	case "golang":
		return "go"
	case "scala":
		return "scala"
	case "kotlin":
		return "kt"
	case "rust":
		return "rs"
	case "php":
		return "php"
	case "typescript":
		return "ts"
	default:
		return "txt"
	}
}

// ProblemSolution is the response from leetcode GraphQL API
// concerning problem solutions
type ProblemSolution struct {
	ID           string `json:"id"`
	CanSeeDetail bool   `json:"canSeeDetail"`
	PaidOnly     bool   `json:"paidOnly"`
	TypeName     string `json:"__typename"`
}

// ProblemStats is the string response from leetcode GraphQL API
// concerning problem stats
type ProblemStats struct {
	TotalAccepted      string `json:"totalAccepted"`
	TotalSubmission    string `json:"totalSubmission"`
	TotalAcceptedRaw   int    `json:"totalAcceptedRaw"`
	TotalSubmissionRaw int    `json:"totalSubmissionRaw"`
	AcceptRate         string `json:"acRate"`
}

// GetDifficulty is a mapper function from problem Difficulty level to string
func (pd ProblemDetail) GetDifficulty() string {
	switch pd.Difficulty {
	case "Easy":
		return utils.Green("Easy")
	case "Medium":
		return utils.Yellow("Medium")
	default:
		return utils.Red("Hard")
	}
}

// GetStats is a property function unmarshal json string field `stats`
func (pd ProblemDetail) GetStats() (*ProblemStats, error) {
	ps := &ProblemStats{}
	err := json.Unmarshal([]byte(pd.Stats), ps)
	if err != nil {
		return nil, err
	}
	return ps, nil
}

// ExportDetail generate source code in local directory
func (pd ProblemDetail) ExportDetail(generate bool, language string, summary bool) error {
	sourceCodePath := ""

	if generate == false {
		return pd.exportStdoutDetail()
	}

	t, err := GetFileTemplate(pd)
	if err != nil {
		return err
	}

	if language != "" {
		sourceCodePath, err = pd.generateSourceCode(t, language)
		if err != nil {
			return err
		}
	}

	err = pd.generateMarkdown(t, sourceCodePath)
	if err != nil {
		return err
	}

	if summary {
		pd.exportGenerateSummary(t)
	}

	return nil
}

func (pd ProblemDetail) exportStdoutDetail() error {
	pds, err := pd.GetStats()
	if err != nil {
		return err
	}

	var tags []string
	for _, t := range pd.TopicTags {
		tags = append(tags, utils.Yellow(t.Name))
	}

	p := strings.NewReader(pd.Content)
	parsedContent, err := goquery.NewDocumentFromReader(p)
	if err != nil {
		return err
	}

	fmt.Printf("[%s] %s\n\n", pd.QuestionID, pd.Title)
	fmt.Printf("%s\n\n", utils.Gray(strings.Replace(utils.ProblemURL, "$slug", pd.TitleSlug, 1)))
	fmt.Printf("Tags: %s \n\n", strings.Join(tags, ", "))
	fmt.Printf("* %s (%s)\n", pd.GetDifficulty(), pds.AcceptRate)
	fmt.Printf("* Total Accepted:    %d\n", pds.TotalAcceptedRaw)
	fmt.Printf("* Total Submissions: %d\n", pds.TotalSubmissionRaw)
	fmt.Printf("* Testcase Example: %s\n\n", strings.ReplaceAll(pd.SampleTestCase, "\n", "\\n"))
	fmt.Println(parsedContent.Text())

	return nil
}

func (pd ProblemDetail) generateMarkdown(t *Template, sourceCodePath string) error {
	markdown := ""

	pds, err := pd.GetStats()
	if err != nil {
		return err
	}

	var tags []string
	for _, t := range pd.TopicTags {
		tags = append(tags, fmt.Sprintf("`%s`", t.Name))
	}

	markdown += fmt.Sprintf("# [%s] %s\n\n", pd.QuestionID, pd.Title)
	markdown += fmt.Sprintf(
		"<%s>\n\n",
		strings.Replace(utils.ProblemURL, "$slug", pd.TitleSlug, 1),
	)
	markdown += fmt.Sprintf("- Tags: %s\n\n", strings.Join(tags, ", "))
	markdown += fmt.Sprintf("- Difficulty: %s\n\n", pd.Difficulty)
	if sourceCodePath != "" {
		markdown += fmt.Sprintf("- Source Code: [./%s](./%s)\n\n", sourceCodePath, sourceCodePath)
	}
	markdown += fmt.Sprintf("- Acceptance: %s\n\n", pds.AcceptRate)
	markdown += fmt.Sprintf("- Total Accepted: %d\n\n", pds.TotalAcceptedRaw)
	markdown += fmt.Sprintf("- Total Submissions: %d\n\n", pds.TotalSubmissionRaw)
	markdown += fmt.Sprintf("- Testcase Example: %s\n\n", strings.ReplaceAll(pd.SampleTestCase, "\n", "\\n"))
	markdown += fmt.Sprintf("## Description\n\n")
	markdown += fmt.Sprintln(pd.Content)

	if _, err := os.Stat(t.DirTemplate); os.IsNotExist(err) {
		os.Mkdir(t.DirTemplate, os.ModePerm)
	}

	f, err := os.Create(
		fmt.Sprintf(
			"%s/%s",
			t.DirTemplate,
			t.MarkdownTemplate,
		),
	)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	defer f.Close()

	_, err = f.WriteString(markdown)
	if err != nil {
		return err
	}

	f.Sync()

	return nil
}

func (pd ProblemDetail) generateSourceCode(t *Template, language string) (string, error) {
	var supportedLanguage []string
	for _, codeSnippet := range pd.CodeSnippets {
		supportedLanguage = append(
			supportedLanguage,
			fmt.Sprintf("%s(%s)", codeSnippet.Lang, codeSnippet.LangSlug),
		)
		if codeSnippet.Lang == language || codeSnippet.LangSlug == language {
			if _, err := os.Stat(t.DirTemplate); os.IsNotExist(err) {
				os.Mkdir(t.DirTemplate, os.ModePerm)
			}

			t.SourceCodeTemplate = strings.ReplaceAll(
				t.SourceCodeTemplate,
				"$ext",
				codeSnippet.GetLanguageExt(),
			)

			f, err := os.Create(
				fmt.Sprintf(
					"%s/%s",
					t.DirTemplate,
					t.SourceCodeTemplate,
				),
			)
			if err != nil {
				return "", fmt.Errorf(err.Error())
			}

			defer f.Close()

			_, err = f.WriteString(codeSnippet.Code)
			if err != nil {
				return "", fmt.Errorf(err.Error())
			}

			f.Sync()

			return t.SourceCodeTemplate, nil
		}
	}

	errMessage := fmt.Sprintf("invalid language '%s' for problem: '%s'", language, pd.Title)
	errMessage += fmt.Sprintf(" with supported language:\n[%s]", strings.Join(supportedLanguage, ", "))
	return "", fmt.Errorf(errMessage)
}

func (pd ProblemDetail) exportGenerateSummary(t *Template) {
	var tags []string
	for _, tag := range pd.TopicTags {
		tags = append(tags, tag.Name)
	}

	fmt.Printf(
		"| [%s](%s) | %s | %s |",
		pd.Title,
		t.DirTemplate,
		strings.Join(tags, ", "),
		pd.Difficulty,
	)
}

// GetLanguageSlug is a mapper function mapping file ext to question slug, with checking
func (pd ProblemDetail) GetLanguageSlug(ext string) (string, error) {
	var slug string
	switch ext {
	case ".cpp":
		slug = "cpp"
	case ".java":
		slug = "java"
	case ".py":
		slug = "python"
	case ".py3":
		slug = "python3"
	case ".c":
		slug = "c"
	case ".cs":
		slug = "csharp"
	case ".js":
		slug = "javascript"
	case ".rb":
		slug = "ruby"
	case ".swift":
		slug = "swift"
	case ".go":
		slug = "golang"
	case ".scala":
		slug = "scala"
	case ".kt":
		slug = "kotlin"
	case ".rs":
		slug = "rust"
	case ".php":
		slug = "php"
	case ".ts":
		slug = "typescript"
	default:
		slug = ""
	}

	for _, pcs := range pd.CodeSnippets {
		if slug == pcs.LangSlug {
			return slug, nil
		}
	}

	return "", fmt.Errorf("question %s does not support file format %s", pd.QuestionID, ext)
}
