package arg

import (
	"fmt"

	"github.com/ckidckidckid/leetcode-cli/pkg/utils"
	"github.com/spf13/cobra"
)

// List cmd argument checking
func List(cmd *cobra.Command, args []string) error {
	category, err := cmd.Flags().GetString("category")
	if err != nil {
		return err
	}
	if !utils.Contains([]interface{}{"all", "algorithms", "database", "shell"}, category) {
		return fmt.Errorf("invalid arguments: %s = %s", "category", category)
	}

	_, err = cmd.Flags().GetString("name")
	if err != nil {
		return err
	}

	_, err = cmd.Flags().GetString("query")
	if err != nil {
		return err
	}

	lock, err := cmd.Flags().GetString("lock")
	if err != nil {
		return err
	}
	if !utils.Contains([]interface{}{"all", "free", "locked"}, lock) {
		return fmt.Errorf("invalid arguments: %s = %s", "lock", lock)
	}

	status, err := cmd.Flags().GetString("status")
	if err != nil {
		return err
	}
	if !utils.Contains([]interface{}{"all", "approved", "rejected", "new"}, status) {
		return fmt.Errorf("invalid arguments: %s = %s", "status", status)
	}

	return nil
}
