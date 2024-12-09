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
			
			fmt.Printf("Successfully installed Pulumi %s\n", version)
			
			if useAfterInstall {
				if err := utils.UseVersion(version); err != nil {
					return fmt.Errorf("failed to switch to version %s: %w", version, err)
				}
				fmt.Printf("Switched to Pulumi %s\n", version)
			} else {
				fmt.Printf("\nTo use this version, run: pvm use %s\n", version)
			}
			
			return nil
		},
	}
	
	cmd.Flags().Bool("use", false, "Switch to this version after installing")
	return cmd
} 