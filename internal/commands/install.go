package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomski747/pvm/internal/utils"
)

func installCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install <version>",
		Short: "Install a specific version of Pulumi",
		Long:  "Install a specific version of Pulumi. Use 'latest' to install the most recent version.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			version := args[0]
			useAfterInstall, _ := cmd.Flags().GetBool("use")

			if version == "latest" {
				latest, err := utils.GetLatestVersion()
				if err != nil {
					return fmt.Errorf("failed to get latest version: %w", err)
				}
				version = latest
			}

			if err := utils.InstallVersion(version); err != nil {
				return err
			}
			resolvedVersion, err := utils.ResolveVersion(version)
			if err != nil {
				return fmt.Errorf("failed to resolve version: %w", err)
			}

			fmt.Printf("%s %s\n", utils.Success("Successfully installed Pulumi"), resolvedVersion)

			if useAfterInstall {
				if err := utils.UseVersion(version); err != nil {
					return fmt.Errorf("failed to switch to version %s: %w", resolvedVersion, err)
				}
				fmt.Printf("%s %s\n", utils.Success("Switched to Pulumi"), resolvedVersion)
			} else {
				fmt.Printf("\n%s pvm use %s\n", utils.Info("To use this version, run:"), resolvedVersion)
			}

			return nil
		},
	}

	cmd.Flags().Bool("use", false, "Switch to this version after installing")
	return cmd
}
