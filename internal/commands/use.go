package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomski747/pvm/internal/utils"
)

var useCmd = &cobra.Command{
	Use:   "use <version>",
	Short: "Switch to a specific version of Pulumi",
	Long:  "Switch to a specific version of Pulumi. Use 'latest' to switch to the most recent version.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]

		if version == "latest" {
			latest, err := utils.GetLatestVersion()
			if err != nil {
				return fmt.Errorf("failed to get latest version: %w", err)
			}
			version = latest
		}

		resolvedVersion, err := utils.ResolveVersion(version)
		if err != nil {
			return fmt.Errorf("failed to resolve version: %w", err)
		}

		if err := utils.UseVersion(resolvedVersion); err != nil {
			return fmt.Errorf("failed to switch to version %s: %w", resolvedVersion, err)
		}

		fmt.Printf("%s %s\n", utils.Success("Switched to Pulumi"), resolvedVersion)
		return nil
	},
}
