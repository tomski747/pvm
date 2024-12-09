package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomski747/pvm/internal/utils"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Pulumi versions",
	Long:  "List installed Pulumi versions. Use --all to show all available versions.",
	RunE: func(cmd *cobra.Command, args []string) error {
		showAll, _ := cmd.Flags().GetBool("all")

		installed := utils.GetInstalledVersions()
		current, err := utils.GetCurrentVersion()
		if err != nil {
			return fmt.Errorf("failed to get current version: %w", err)
		}

		if showAll {
			available, err := utils.GetAvailableVersions()
			if err != nil {
				return fmt.Errorf("failed to fetch available versions: %w", err)
			}

			fmt.Println(utils.Info("Available versions:"))
			for _, version := range available {
				prefix := "  "
				if installed[version] {
					prefix = utils.Success("* ")
				}
				if version == current {
					prefix = utils.Current("→ ")
				}
				fmt.Printf("%s%s\n", prefix, version)
			}
		} else {
			if len(installed) == 0 {
				fmt.Println(utils.Warning("No versions installed. Use 'pvm install <version>' to install one."))
				fmt.Println(utils.Info("Run 'pvm list --all' to see all available versions."))
				return nil
			}

			fmt.Println(utils.Info("Installed versions:"))
			for version := range installed {
				prefix := "  "
				if version == current {
					prefix = utils.Current("→ ")
				}
				fmt.Printf("%s%s\n", prefix, version)
			}
		}

		fmt.Println(utils.Info("\nLegend:"))
		fmt.Println(utils.Current("  →  current"))
		if showAll {
			fmt.Println(utils.Success("  *  installed"))
		}

		return nil
	},
}

func init() {
	listCmd.Flags().Bool("all", false, "Show all available versions")
}
