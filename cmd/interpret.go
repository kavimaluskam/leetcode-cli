package cmd

import (
	"path/filepath"

	"github.com/ckidckidckid/leetcode-cli/pkg/api"
	"github.com/ckidckidckid/leetcode-cli/pkg/arg"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(interpretCmd)
	interpretCmd.PersistentFlags().IntP("id", "i", 0, "ID of problem to be submitted")
	interpretCmd.PersistentFlags().StringP("file", "f", "", "path of file to be submitted")
	interpretCmd.PersistentFlags().StringP("test_input", "t", "", "test input to be submitted")
}

var interpretCmd = &cobra.Command{
	Use:   `interpret`,
	Short: `Interpret code`,
	Long:  `Interpret local code to leetcode problem with testing input`,
	Args:  arg.Interpret,
	RunE:  interpret,
}

func interpret(cmd *cobra.Command, args []string) error {
	id, _ := cmd.Flags().GetInt("id")
	file, _ := cmd.Flags().GetString("file")
	testInput, _ := cmd.Flags().GetString("test_input")

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

	err = sClient.InterpretCode(problemDetail, fp, testInput)
	if err != nil {
		return err
	}

	return nil
}
