package cmd

import (
	"github.com/kavimaluskam/leetcode-cli/pkg/api"
	"github.com/kavimaluskam/leetcode-cli/pkg/arg"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(showCmd)
	showCmd.Flags().IntP("id", "i", 0, "ID of problem to be shown")
	showCmd.Flags().BoolP("random", "r", false, "Random choice of problem to be shown")
	showCmd.Flags().StringP("language", "l", "", "Open source code in editor")
	showCmd.Flags().BoolP("summary", "s", false, "Print out generation summary")
}

var showCmd = &cobra.Command{
	Use:     `show`,
	Aliases: []string{`dl`, `pick`, `show`},
	Short:   `Show individual problem`,
	Long:    `Show or download individual problem description and code template`,
	Args:    arg.Show,
	RunE:    show,
}

func show(cmd *cobra.Command, args []string) error {
	id, _ := cmd.Flags().GetInt("id")
	random, _ := cmd.Flags().GetBool("random")
	language, _ := cmd.Flags().GetString("language")
	summary, _ := cmd.Flags().GetBool("summary")

	client, err := api.GetAuthClient()
	if err != nil {
		return err
	}

	problemDetail, err := client.GetProblemDetail(id, random)
	if err != nil {
		return err
	}

	err = problemDetail.ExportDetail(language, summary)
	if err != nil {
		return err
	}

	return nil
}
