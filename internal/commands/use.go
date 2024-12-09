package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomski747/pvm/internal/utils"
)

var useCmd = &cobra.Command{
	Use:   "use <version>",
	Short: "Switch to a specific version of Pulumi",
	Long:  `Switch to a specific version of Pulumi that has been installed.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]
		fmt.Printf("Switching to Pulumi version %s...\n", version)
		
		if err := utils.UseVersion(version); err != nil {
			return fmt.Errorf("failed to switch to version %s: %v", version, err)
		}

		fmt.Printf("Successfully switched to Pulumi version %s\n", version)
		return nil
	},
} 