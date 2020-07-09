package arg

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Interpret cmd argument checking
func Interpret(cmd *cobra.Command, args []string) error {
	id, err := cmd.Flags().GetInt("id")
	if err != nil {
		return err
	}
	if id == 0 {
		return fmt.Errorf("missing required parameter: 'id'")
	}

	file, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}
	if file == "" {
		return fmt.Errorf("missing required parameter: 'file'")
	}

	testInput, err := cmd.Flags().GetString("test_input")
	if err != nil {
		return err
	}
	if testInput == "" {
		return fmt.Errorf("missing required parameter: 'test_input'")
	}

	return nil
}
