package cmd

import (
	"fmt"

	"github.com/kavimaluskam/leetcode-cli/pkg/api"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("category", "c", "all", "problem categories: {''|all|algorithms|database|shell}")
	listCmd.Flags().StringP("query", "q", "", "problem query string")
}

var listCmd = &cobra.Command{
	Use:     `list`,
	Aliases: []string{`li`},
	Short:   `Listing questions`,
	RunE:    list,
}

func list(cmd *cobra.Command, args []string) error {
	category, err := cmd.Flags().GetString("category")
	if err != nil {
		return err
	}

	query, err := cmd.Flags().GetString("query")
	if err != nil {
		return err
	}

	client := api.NewClient()
	var problemCollection *api.ProblemCollection

	switch category {
	case "all":
		problemCollection, err = client.GetProblemCollection("all", query)
	case "algorithms", "database", "shell":
		problemCollection, err = client.GetProblemCollection(category, query)
	default:
		return fmt.Errorf("unsupported category %s", category)
	}
	if err != nil {
		return err
	}

	for _, problem := range problemCollection.Problems {
		fmt.Printf(
			"%s [%4d] %-60s %s (%.2f %%)\n",
			problem.GetLockedStatus(),
			problem.Stat.QuestionID,
			problem.Stat.QuestionTitle,
			problem.GetDiffculty("%-6s"),
			(float64(problem.Stat.TotalAcs) / float64(problem.Stat.TotalSubmitted)),
		)
	}

	return nil
}
