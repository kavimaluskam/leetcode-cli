package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.PersistentFlags().Bool("help", false, "Show help for command")
}

var RootCmd = &cobra.Command{
	Use:   "lc <command> <subcommand> [flags]",
	Short: "leetcode CLI",
	Long:  `Work seamlessly with leetcode from the command line.`,
}
