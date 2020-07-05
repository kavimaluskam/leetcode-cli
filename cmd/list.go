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
	listCmd.Flags().StringP("name", "n", "", "problem name query")
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

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}

	query, err := cmd.Flags().GetString("query")
	if err != nil {
		return err
	}

	acceptedCategory := map[string]bool{
		"all":        true,
		"algorithms": true,
		"database":   true,
		"shell":      true,
	}
	if acceptedCategory[category] == false {
		// TODO: enhance error type handling
		return fmt.Errorf(
			"Failed to list problems with unsupported category %s",
			category,
		)
	}

	client, err := api.GetAuthClient()
	if err != nil {
		return err
	}

	problemCollection, err := client.GetProblemCollection(category, query, name)
	if err != nil {
		return err
	}

	for _, problem := range problemCollection.Problems {
		fmt.Printf(
			"%2s%2s%2s [%4d] %-60s %s (%.2f %%)\n",
			problem.GetLockedStatus(),
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
