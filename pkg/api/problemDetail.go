package api

import (
	"fmt"

	"github.com/kavimaluskam/leetcode-cli/pkg/utils"
)

// ProblemDetailCollection is the response from leetcode GraphQL API
// concerning problem deatil
type ProblemDetailCollection struct {
	Question ProblemDetail `json:"question"`
}

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
	TrnslatedContent      string                `json:"translatedContent"`
	IsPaidOnly            bool                  `json:"isPaidOnly"`
	Diffculty             string                `json:"difficulty"`
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

// ProblemSolution is the response from leetcode GraphQL API
// concerning problem solutions
type ProblemSolution struct {
	ID           string `json:"id"`
	CanSeeDetail bool   `json:"canSeeDetail"`
	PaidOnly     bool   `json:"paidOnly"`
	TypeName     string `json:"__typename"`
}

// GetProblemDetail is the graphql query function fetching leetcode Individual Problem
func (client *Client) GetProblemDetail(id int, titleSlug string) (*ProblemDetail, error) {
	var problemDetailCollection ProblemDetailCollection

	if titleSlug == "" && id != 0 {
		problemCollection, err := client.GetProblemCollection("all", "", "", "all", "all")
		if err != nil {
			return nil, err
		}

		for _, problem := range problemCollection.Problems {
			if problem.Stat.QuestionID == id {
				titleSlug = problem.Stat.QuestionTitleSlug
			}
		}

		if titleSlug == "" {
			// TODO: enhance error type handling
			return nil, fmt.Errorf(
				"Failed to find problem with ID %d",
				id,
			)
		}
	}

	variables := make(map[string]interface{})
	variables["titleSlug"] = titleSlug

	err := client.GraphQL(
		utils.QuestionDataOperation,
		utils.QuestionDataQuery,
		variables,
		&problemDetailCollection,
	)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", problemDetailCollection)

	return &problemDetailCollection.Question, nil
}
