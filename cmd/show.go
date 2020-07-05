package cmd

import (
	"fmt"

	"github.com/kavimaluskam/leetcode-cli/pkg/api"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:     `show`,
	Aliases: []string{`dl`, `pick`, `show`},
	Short:   `Show individual problem`,
	RunE:    show,
}

func show(cmd *cobra.Command, args []string) error {
	client, err := api.GetAuthClient()
	if err != nil {
		return err
	}

	problemDetail, err := client.GetProblemDetail(1, "")
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", problemDetail)
	return nil
}
