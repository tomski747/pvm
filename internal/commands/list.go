package commands

import (
	"fmt"
	"sort"
	"strings"
	"strconv"

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

// semverCompare compares two version strings
func semverCompare(v1, v2 string) bool {
	v1Parts := strings.Split(v1, ".")
	v2Parts := strings.Split(v2, ".")
	
	for k := 0; k < len(v1Parts) && k < len(v2Parts); k++ {
		n1, _ := strconv.Atoi(v1Parts[k])
		n2, _ := strconv.Atoi(v2Parts[k])
		if n1 != n2 {
			return n1 > n2
		}
	}
	return len(v1Parts) > len(v2Parts)
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

			// Sort versions using semver comparison
			sort.Slice(versions, func(i, j int) bool {
				return semverCompare(versions[i], versions[j])
			})

			fmt.Println(utils.Info("Available versions:"))
			for _, version := range versions {
				prefix := "  "
				if installed[version] {
					if version == current {
						prefix = utils.Current("→ ")
					} else {
						prefix = utils.Success("* ")
					}
				}
				fmt.Printf("%s%s\n", prefix, version)
			}
		} else {
			if len(installed) == 0 {
				fmt.Println(utils.Warning("No versions installed. Use 'pvm install <version>' to install one."))
				fmt.Println(utils.Info("Run 'pvm list --all' to see all available versions."))
				return nil
			}

			// Convert installed map to slice for sorting
			installedVersions := make([]string, 0, len(installed))
			for version := range installed {
				installedVersions = append(installedVersions, version)
			}

			// Sort installed versions
			sort.Slice(installedVersions, func(i, j int) bool {
				return semverCompare(installedVersions[i], installedVersions[j])
			})

			fmt.Println(utils.Info("Installed versions:"))
			for _, version := range installedVersions {
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
