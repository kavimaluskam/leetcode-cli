package cmd

import (
	"fmt"

	"github.com/kavimaluskam/leetcode-cli/pkg/api"
	"github.com/kavimaluskam/leetcode-cli/pkg/arg"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("category", "c", "all", "problem categories: {all|algorithms|database|shell}")
	listCmd.Flags().StringP("query", "q", "", "problem query string")
	listCmd.Flags().StringP("name", "n", "", "problem name query string")
	listCmd.Flags().StringP("lock", "l", "all", "problem lock status: {all|free|locked}")
	listCmd.Flags().StringP("status", "s", "all", "problem status: {all|approved|rejected|new}")
}

var listCmd = &cobra.Command{
	Use:     `list`,
	Aliases: []string{`li`},
	Short:   `Listing questions`,
	Args:    arg.List,
	RunE:    list,
}

func list(cmd *cobra.Command, args []string) error {
	category, _ := cmd.Flags().GetString("category")
	name, _ := cmd.Flags().GetString("name")
	query, _ := cmd.Flags().GetString("query")
	lock, _ := cmd.Flags().GetString("lock")
	status, _ := cmd.Flags().GetString("status")

	client, err := api.GetAuthClient()
	if err != nil {
		return err
	}

	problemCollection, err := client.GetProblemCollection(category, query, name, lock, status)
	if err != nil {
		return err
	}

	for _, problem := range problemCollection.Problems {
		fmt.Printf(
			"%2s%2s%2s [%4d] %-60s %s (%.2f %%)\n",
			problem.GetLockStatus(),
			problem.GetIsFavor(),
			problem.GetStatus(),
			problem.Stat.QuestionID,
			problem.Stat.QuestionTitle,
			problem.GetDiffculty("%-6s"),
			(float64(problem.Stat.TotalAcs) / float64(problem.Stat.TotalSubmitted)),
		)
	}

	return nil
}
