package api

import (
	"strings"

	"github.com/ckidckidckid/leetcode-cli/pkg/model"
	"github.com/ckidckidckid/leetcode-cli/pkg/utils"
)

// ProblemCollection is the response from leetcode API concerning problem set
type ProblemCollection struct {
	Username  string          `json:"user_name"`
	NumSolved int             `json:"num_solved"`
	NumTotal  int             `json:"num_total"`
	AcEasy    int             `json:"ac_easy"`
	AcMedium  int             `json:"ac_medium"`
	AcHard    int             `json:"ac_hard"`
	Problems  []model.Problem `json:"stat_status_pairs"`
}

// GetProblemCollection is the query function fetching leetcode Problem List
func (client *Client) GetProblemCollection(category string, query string, name string, lock string, status string) (*ProblemCollection, error) {
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
		var queriedProblems []model.Problem

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

		var queriedProblems []model.Problem

		for _, problem := range problemCollection.Problems {
			for _, queryQuestionID := range problemIDList {
				if problem.Stat.FrontendQuestionID == queryQuestionID {
					queriedProblems = append(queriedProblems, problem)
				}
			}
		}
		problemCollection.Problems = queriedProblems
	}

	// filter problems by lock status
	if lock != "all" {
		var queriedProblems []model.Problem
		for _, problem := range problemCollection.Problems {
			if problem.CheckLockStatus(lock) {
				queriedProblems = append(queriedProblems, problem)
			}
		}
		problemCollection.Problems = queriedProblems
	}

	// filter problems by status
	if status != "all" {
		var queriedProblems []model.Problem
		for _, problem := range problemCollection.Problems {
			if problem.CheckStatus(status) {
				queriedProblems = append(queriedProblems, problem)
			}
		}
		problemCollection.Problems = queriedProblems
	}

	return &problemCollection, nil
}
