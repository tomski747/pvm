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
		installIfMissing, _ := cmd.Flags().GetBool("install")

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

		// Check if version is installed
		installed := utils.GetInstalledVersions()
		if !installed[resolvedVersion] {
			if !installIfMissing {
				return fmt.Errorf("version %s is not installed. Use 'pvm install %s' first or retry with --install flag", resolvedVersion, resolvedVersion)
			}

			// Install the version if --install flag is used
			if err := utils.InstallVersion(resolvedVersion); err != nil {
				return fmt.Errorf("failed to install version %s: %w", resolvedVersion, err)
			}
			fmt.Printf("%s %s\n", utils.Success("Successfully installed Pulumi"), resolvedVersion)
		}

		if err := utils.UseVersion(resolvedVersion); err != nil {
			return fmt.Errorf("failed to switch to version %s: %w", resolvedVersion, err)
		}

		fmt.Printf("%s %s\n", utils.Success("Switched to Pulumi"), resolvedVersion)
		return nil
	},
}

func init() {
	useCmd.Flags().Bool("install", false, "Install the version if not already installed")
}
