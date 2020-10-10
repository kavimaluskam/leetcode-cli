package arg

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Show cmd argument checking
func Show(cmd *cobra.Command, args []string) error {
	id, err := cmd.Flags().GetInt("id")
	if err != nil {
		return err
	}
	random, err := cmd.Flags().GetBool("random")
	if err != nil {
		return err
	}

	if id == 0 && !random {
		return fmt.Errorf("invalid arguments: either 'id', 'random' should be applied")
	}

	_, err = cmd.Flags().GetString("language")
	if err != nil {
		return err
	}

	return nil
}
