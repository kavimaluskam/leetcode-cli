package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.PersistentFlags().Bool("help", false, "Show help for command")
}

var RootCmd = &cobra.Command{
	Use:   "lc <command> <subcommand> [flags]",
	Short: "LeetCode CLI",
	Long:  `Work seamlessly with LeetCode from the command line.`,
}
