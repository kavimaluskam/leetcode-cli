package utils

import (
	"fmt"
	"os"
)

// URLs supported by leetcode api
const (
	BaseURL           = "https://leetcode.com"
	GraphQLURL        = "https://leetcode.com/graphql"
	LoginURL          = "https://leetcode.com/accounts/login/"
	ProblemListingURL = "https://leetcode.com/api/problems/$category/"
	ProblemQueryURL   = "https://leetcode.com/problems/api/filter-questions/$query"
	ProblemURL        = "https://leetcode.com/problems/$slug/description/"
	TestURL           = "https://leetcode.com/problems/$slug/interpret_solution/"
	SessionURL        = "https://leetcode.com/session/"
	SubmitURL         = "https://leetcode.com/problems/$slug/submit/"
	SubmissionsURL    = "https://leetcode.com/api/submissions/$slug"
	SubmissionURL     = "https://leetcode.com/submissions/detail/$id/"
	VerifyURL         = "https://leetcode.com/submissions/detail/$id/check/"
	FavoritesURL      = "https://leetcode.com/list/api/questions"
	FavoriteDeleteURL = "https://leetcode.com/list/api/questions/$hash/$id"
)

var (
	AuthConfigPath = fmt.Sprintf("%s/.lc/leetcode/user.json", os.Getenv("HOME"))
)
