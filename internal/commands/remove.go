package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomski747/pvm/internal/utils"
)

var removeCmd = &cobra.Command{
	Use:   "remove <version>",
	Short: "Remove a specific version of Pulumi",
	Long:  "Remove a specific version of Pulumi that has been installed.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]

		if err := utils.RemoveVersion(version); err != nil {
			return fmt.Errorf("failed to remove version %s: %w", version, err)
		}

		fmt.Printf("Successfully removed Pulumi %s\n", version)
		return nil
	},
}
