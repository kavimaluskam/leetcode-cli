package arg

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Submit cmd argument checking
func Submit(cmd *cobra.Command, args []string) error {
	id, err := cmd.Flags().GetInt("id")
	if err != nil {
		return err
	}
	if id == 0 {
		return fmt.Errorf("missing required parameter: 'id'")
	}

	_, err = cmd.Flags().GetString("file")
	if err != nil {
		return err
	}

	return nil
}
