package api

import (
	"strings"

	"github.com/kavimaluskam/leetcode-cli/pkg/utils"
	"github.com/kyokomi/emoji"
)

// ProblemCollection is the response from leetcode API concerning problem set
type ProblemCollection struct {
	Username  string    `json:"user_name"`
	NumSolved int       `json:"num_solved"`
	NumTotal  int       `json:"num_total"`
	AcEasy    int       `json:"ac_easy"`
	AcMedium  int       `json:"ac_medium"`
	AcHard    int       `json:"ac_hard"`
	Problems  []Problem `json:"stat_status_pairs"`
}

// Problem is the response from leetcode API concerning individual problems
type Problem struct {
	Stat struct {
		QuestionID          int    `json:"question_id"`
		QuestionArticleLive bool   `json:"question__article__live"`
		QuestionArticleSlug string `json:"question__article__slug"`
		QuestionTitle       string `json:"question__title"`
		QuestionTitleSlug   string `json:"question__title_slug"`
		TotalAcs            int    `json:"total_acs"`
		TotalSubmitted      int    `json:"total_submitted"`
		FrontendQuestionID  int    `json:"frontend_question_id"`
		IsNewQuestion       bool   `json:"is_new_question"`
	} `json:"stat"`
	Status    string `json:"status"`
	Diffculty struct {
		Level int `json:"level"`
	} `json:"difficulty"`
	PaidOnly  bool `json:"paid_only"`
	IsFavor   bool `json:"is_favor"`
	Frequency int  `json:"frequency"`
	Progress  int  `json:"progress"`
}

// GetDiffculty is a mapper function from problem diffculty level to string
func (p Problem) GetDiffculty(format string) string {
	switch p.Diffculty.Level {
	case 1:
		return utils.GreenFormatted("Easy", format)
	case 2:
		return utils.YellowFormatted("Medium", format)
	default:
		return utils.RedFormatted("Hard", format)
	}
}

// GetStatus is a mapper function from problem status to emoji
func (p Problem) GetStatus() string {
	switch p.Status {
	case "ac": // Problem approved
		return emoji.Sprint(":heavy_check_mark: ")
	case "notac": // Problem WIP
		return emoji.Sprint(":question:")
	default:
		return "   "
	}
}

// GetIsFavor is a mapper function from `is_favor` status to emoji
func (p Problem) GetIsFavor() string {
	if p.IsFavor {
		return emoji.Sprint(":heart: ")
	}
	return "   "
}

// GetLockedStatus is a mapper function from `paid_only` status to emoji
func (p Problem) GetLockedStatus() string {
	if p.PaidOnly {
		return emoji.Sprint(":locked:")
	}
	return "   "
}

// GetProblemCollection is the query function fetching leetcode Problem List
func (client *Client) GetProblemCollection(category string, query string, name string) (*ProblemCollection, error) {
	var problemCollection ProblemCollection
	var problemIDList []int

	url := strings.Replace(utils.ProblemListingURL, "$category", category, 1)
	err := client.REST("GET", url, nil, &problemCollection)

	if err != nil {
		return nil, err
	}

	// filter problems by name
	if name != "" {
		name = strings.ToLower(name)
		var queriedProblems []Problem

		for _, problem := range problemCollection.Problems {
			if strings.Contains(strings.ToLower(problem.Stat.QuestionTitle), name) {
				queriedProblems = append(queriedProblems, problem)
			}
		}
		problemCollection.Problems = queriedProblems
	}

	// filter problems by queried IDs
	if query != "" {
		queryURL := strings.Replace(utils.ProblemQueryURL, "$query", query, 1)

		err = client.REST("GET", queryURL, nil, &problemIDList)
		if err != nil {
			return nil, err
		}

		var queriedProblems []Problem

		for _, problem := range problemCollection.Problems {
			for _, queryQuestionID := range problemIDList {
				if problem.Stat.QuestionID == queryQuestionID {
					queriedProblems = append(queriedProblems, problem)
				}
			}
		}
		problemCollection.Problems = queriedProblems
	}

	return &problemCollection, nil
}
