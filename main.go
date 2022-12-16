package main

import (
	"fmt"
	"io"
	"os"

	"github.com/ckidckidckid/leetcode-cli/cmd"
	"github.com/spf13/cobra"
)

var updaterEnabled = ""

func main() {
	if cmd, err := cmd.RootCmd.ExecuteC(); err != nil {
		printError(os.Stderr, err, cmd)
	}
}

func printError(out io.Writer, err error, cmd *cobra.Command) {
	fmt.Fprintln(out, err)
	fmt.Fprintln(out)
	fmt.Fprintln(out, cmd.UsageString())
}
