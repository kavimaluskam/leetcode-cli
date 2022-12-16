package api

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ckidckidckid/leetcode-cli/pkg/model"
	"github.com/ckidckidckid/leetcode-cli/pkg/utils"
)

// ProblemDetailCollection is the response from leetcode GraphQL API
// concerning problem deatil
type ProblemDetailCollection struct {
	Question model.ProblemDetail `json:"question"`
}

// GetProblemDetail is the graphql query function fetching leetcode Individual Problem
func (client *Client) GetProblemDetail(id int, random bool) (*model.ProblemDetail, error) {
	var titleSlug string
	var problemDetailCollection ProblemDetailCollection

	if random { // randomly pick problem title slug
		problemCollection, err := client.GetProblemCollection("all", "", "", "free", "new")
		if err != nil {
			return nil, err
		}

		rand.Seed(time.Now().Unix())
		i := rand.Int() % len(problemCollection.Problems)

		titleSlug = problemCollection.Problems[i].Stat.QuestionTitleSlug

	} else { // pick title slug by id
		problemCollection, err := client.GetProblemCollection("all", "", "", "all", "all")
		if err != nil {
			return nil, err
		}

		for _, problem := range problemCollection.Problems {
			if problem.Stat.FrontendQuestionID == id {
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

	return &problemDetailCollection.Question, nil
}
