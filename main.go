package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/kavimaluskam/leetcode-cli/cmd"
	"github.com/kavimaluskam/leetcode-cli/pkg/cmd/util"
	"github.com/spf13/cobra"
)

var updaterEnabled = ""

func main() {
	hasDebug := os.Getenv("DEBUG") != ""

	if cmd, err := cmd.RootCmd.ExecuteC(); err != nil {
		printError(os.Stderr, err, cmd, hasDebug)
		os.Exit(1)
	}
}

func printError(out io.Writer, err error, cmd *cobra.Command, debug bool) {
	if err == util.SilentError {
		return
	}

	var dnsError *net.DNSError
	if errors.As(err, &dnsError) {
		fmt.Fprintf(out, "error connecting to %s\n", dnsError.Name)
		if debug {
			fmt.Fprintln(out, dnsError)
		}
		fmt.Fprintln(out, "check your internet connection or leetcode.com")
		return
	}

	fmt.Fprintln(out, err)

	var flagError *util.FlagError
	if errors.As(err, &flagError) || strings.HasPrefix(err.Error(), "unknown command ") {
		if !strings.HasSuffix(err.Error(), "\n") {
			fmt.Fprintln(out)
		}
		fmt.Fprintln(out, cmd.UsageString())
	}
}
