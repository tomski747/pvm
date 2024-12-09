package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tomski747/pvm/internal/config"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print pvm version",
	Long:  `Print the version information of pvm`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", config.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
} 