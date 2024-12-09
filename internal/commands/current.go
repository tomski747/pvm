package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomski747/pvm/internal/utils"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current Pulumi version",
	Long:  `Display the currently active version of Pulumi.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		version, err := utils.GetCurrentVersion()
		if err != nil {
			return fmt.Errorf("failed to get current version: %v", err)
		}

		if version == "" {
			fmt.Println("No Pulumi version currently selected. Use 'pvm use <version>' to select one.")
			return nil
		}

		fmt.Printf("Current Pulumi version: %s\n", version)
		return nil
	},
} 