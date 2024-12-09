package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomski747/pvm/internal/utils"
)

var installCmd = &cobra.Command{
	Use:   "install <version>",
	Short: "Install a specific version of Pulumi",
	Long:  `Install a specific version of Pulumi from GitHub releases.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]
		fmt.Printf("Installing Pulumi version %s...\n", version)
		
		// Verify version exists
		versions, err := utils.FetchGitHubReleases()
		if err != nil {
			return fmt.Errorf("failed to verify version: %v", err)
		}

		versionExists := false
		for _, v := range versions {
			if v == version {
				versionExists = true
				break
			}
		}

		if !versionExists {
			return fmt.Errorf("version %s does not exist", version)
		}

		// Install the version
		if err := utils.InstallVersion(version); err != nil {
			return fmt.Errorf("failed to install version %s: %v", version, err)
		}

		fmt.Printf("Successfully installed Pulumi version %s\n", version)
		return nil
	},
} 