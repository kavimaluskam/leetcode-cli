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
	SubmitURL         = "https://leetcode.com/problems/$slug/submit/"
	SubmitRefererURL  = "https://leetcode.com/problems/$slug/submissions/"
	InterpretURL      = "https://leetcode.com/problems/$slug/interpret_solution/"
	VerifyURL         = "https://leetcode.com/submissions/detail/$id/check/"
)

// Local Path for configuration
var (
	AuthConfigPath     = fmt.Sprintf("%s/.lc/leetcode/user.json", os.Getenv("HOME"))
	TemplateConfigPath = fmt.Sprintf("%s/.lc/leetcode/template.json", os.Getenv("HOME"))
)

// GraphQL related query, operation string
const (
	QuestionDataQuery = `
		query questionData($titleSlug: String!) {
		    question(titleSlug: $titleSlug) {
		        questionId
		        questionFrontendId
		        boundTopicId
		        title
		        titleSlug
		        content
		        translatedTitle
		        translatedContent
		        isPaidOnly
		        difficulty
		        likes
		        dislikes
		        isLiked
		        similarQuestions
		        contributors {
		            username
		            profileUrl
		            avatarUrl
		            __typename
		        }
		        langToValidPlayground
		        topicTags {
		            name
		            slug
		            translatedName
		            __typename
		        }
		        companyTagStats
		        codeSnippets {
		            lang
		            langSlug
		            code
		            __typename
		        }
		        stats
		        hints
		        solution {
		            id
		            canSeeDetail
		            paidOnly
		            __typename
		        }
		        status
		        sampleTestCase
		        metaData
		        judgerAvailable
		        judgeType
		        mysqlSchemas
		        enableRunCode
		        enableTestMode
		        enableDebugger
		        envInfo
		        libraryUrl
		        adminUrl
		        __typename
		    }
		}`
	QuestionDataOperation = "questionData"
)
