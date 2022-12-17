package model

import (
	"fmt"

	"github.com/ckidckidckid/leetcode-cli/pkg/utils"
	"github.com/kyokomi/emoji"
)

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
	Status     string `json:"status"`
	Difficulty struct {
		Level int `json:"level"`
	} `json:"difficulty"`
	PaidOnly  bool    `json:"paid_only"`
	IsFavor   bool    `json:"is_favor"`
	Frequency float64 `json:"frequency"`
	Progress  float64 `json:"progress"`
}

// GetDifficulty is a mapper function from problem Difficulty level to string
func (p Problem) GetDifficulty(format string) string {
	switch p.Difficulty.Level {
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
		return emoji.Sprint(":x:")
	default:
		return "   "
	}
}

// CheckStatus is a switcher function checking problem status with `status` checker
func (p Problem) CheckStatus(checker string) bool {
	switch checker {
	case "approved":
		return p.Status == "ac"
	case "rejected":
		return p.Status == "notac"
	case "new":
		return p.Status == ""
	default:
		return true
	}
}

// GetIsFavor is a mapper function from `is_favor` status to emoji
func (p Problem) GetIsFavor() string {
	if p.IsFavor {
		return emoji.Sprint(":heart: ")
	}
	return "   "
}

// GetLockStatus is a mapper function from `paid_only` status to emoji
func (p Problem) GetLockStatus() string {
	if p.PaidOnly {
		return emoji.Sprint(":locked:")
	}
	return "   "
}

// CheckLockStatus is a switcher function checking `paid_only` with `lock` checker
func (p Problem) CheckLockStatus(checker string) bool {
	switch checker {
	case "locked":
		return p.PaidOnly
	case "free":
		return !p.PaidOnly
	default:
		return true
	}
}

// ExportStdoutListing prints problem as row in stdout table
func (p Problem) ExportStdoutListing() {
	fmt.Printf(
		"%2s%2s%2s [%4d] %-60s %s (%.2f %%)\n",
		p.GetLockStatus(),
		p.GetIsFavor(),
		p.GetStatus(),
		p.Stat.QuestionID,
		p.Stat.QuestionTitle,
		p.GetDifficulty("%-6s"),
		(float64(p.Stat.TotalAcs) / float64(p.Stat.TotalSubmitted)),
	)

}
