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
	title, err := cmd.Flags().GetString("title")
	if err != nil {
		return err
	}
	random, err := cmd.Flags().GetBool("random")
	if err != nil {
		return err
	}

	if id == 0 && title == "" && !random {
		return fmt.Errorf("invalid arguments: either 'id', 'title', 'random' should be applied")
	}

	generate, err := cmd.Flags().GetBool("generate")
	if err != nil {
		return err
	}
	language, err := cmd.Flags().GetString("language")
	if err != nil {
		return err
	}
	if language != "" && generate == false {
		return fmt.Errorf("invalid arguments: 'language' should only be applied with 'generate'")
	}

	return nil
}
