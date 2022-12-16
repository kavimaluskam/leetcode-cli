package cmd

import (
	"path/filepath"

	"github.com/ckidckidckid/leetcode-cli/pkg/api"
	"github.com/ckidckidckid/leetcode-cli/pkg/arg"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(submitCmd)
	submitCmd.PersistentFlags().IntP("id", "i", 0, "ID of problem to be submitted")
	submitCmd.PersistentFlags().StringP("file", "f", "", "path of file to be submitted")
}

var submitCmd = &cobra.Command{
	Use:   `submit`,
	Short: `Submit code`,
	Long:  `Submit local code to leetcode problem`,
	Args:  arg.Submit,
	RunE:  submit,
}

func submit(cmd *cobra.Command, args []string) error {
	id, _ := cmd.Flags().GetInt("id")
	file, _ := cmd.Flags().GetString("file")

	fp, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	client, err := api.GetAuthClient()
	if err != nil {
		return err
	}

	problemDetail, err := client.GetProblemDetail(id, false)
	if err != nil {
		return err
	}

	sClient, err := api.GetSubmitClient(problemDetail)
	if err != nil {
		return err
	}

	err = sClient.SubmitCode(problemDetail, fp)
	if err != nil {
		return err
	}

	return nil
}
