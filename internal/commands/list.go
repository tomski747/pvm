package commands

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/tomski747/pvm/internal/utils"
)

var (
	refresh bool
)

func init() {
	listCmd.Flags().BoolVar(&refresh, "refresh", false, "Force refresh the version cache")
	listCmd.Flags().Bool("all", false, "Show all available versions")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Pulumi versions",
	Long:  "List installed Pulumi versions. Use --all to show all available versions.",
	RunE: func(cmd *cobra.Command, args []string) error {
		showAll, _ := cmd.Flags().GetBool("all")

		installed := utils.GetInstalledVersions()
		current, err := utils.GetCurrentVersion()
		if err != nil {
			return fmt.Errorf("failed to get current version: %v", err)
		}

		if showAll {
			versions, err := utils.GetAvailableVersions(refresh)
			if err != nil {
				return fmt.Errorf("failed to fetch available versions: %v", err)
			}

			sort.Slice(versions, func(i, j int) bool {
				return utils.SemverGreater(versions[i], versions[j])
			})

			fmt.Fprintln(cmd.OutOrStdout(), utils.Info("Available versions:"))
			for _, version := range versions {
				prefix := "  "
				if installed[version] {
					if version == current {
						prefix = utils.Current("→ ")
					} else {
						prefix = utils.Success("* ")
					}
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s%s\n", prefix, version)
			}
		} else {
			if len(installed) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), utils.Warning("No versions installed. Use 'pvm install <version>' to install one."))
				fmt.Fprintln(cmd.OutOrStdout(), utils.Info("Run 'pvm list --all' to see all available versions."))
				return nil
			}

			installedVersions := make([]string, 0, len(installed))
			for version := range installed {
				installedVersions = append(installedVersions, version)
			}

			sort.Slice(installedVersions, func(i, j int) bool {
				return utils.SemverGreater(installedVersions[i], installedVersions[j])
			})

			fmt.Fprintln(cmd.OutOrStdout(), utils.Info("Installed versions:"))
			for _, version := range installedVersions {
				prefix := "  "
				if version == current {
					prefix = utils.Current("→ ")
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s%s\n", prefix, version)
			}
		}

		fmt.Fprintln(cmd.OutOrStdout(), utils.Info("\nLegend:"))
		fmt.Fprintln(cmd.OutOrStdout(), utils.Current("  →  current"))
		if showAll {
			fmt.Fprintln(cmd.OutOrStdout(), utils.Success("  *  installed"))
		}

		return nil
	},
}
