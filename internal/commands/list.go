package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomski747/pvm/internal/utils"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available Pulumi versions",
	Long:  `List all available Pulumi versions from GitHub releases and show which ones are installed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		versions, err := utils.FetchGitHubReleases()
		if err != nil {
			return fmt.Errorf("failed to fetch versions: %v", err)
		}

		installed := utils.GetInstalledVersions()
		current, _ := utils.GetCurrentVersion()

		fmt.Println("Available Pulumi versions:")
		for _, version := range versions {
			indicators := []string{}
			if installed[version] {
				indicators = append(indicators, "*")
			}
			if version == current {
				indicators = append(indicators, "(current)")
			}

			if len(indicators) > 0 {
				fmt.Printf("  %s %s\n", version, indicators)
			} else {
				fmt.Printf("  %s\n", version)
			}
		}

		fmt.Println("\nLegend:")
		fmt.Println("  * downloaded")
		fmt.Println("  (current) currently active version")
		return nil
	},
} 