package commands

import (
	"github.com/spf13/cobra"
	"github.com/tomski747/pvm/internal/utils"
)

var rootCmd = &cobra.Command{
	Use:   "pvm",
	Short: "Pulumi Version Manager",
	Long: `PVM is a version manager for Pulumi CLI.
It allows you to install and switch between different versions of Pulumi.

Examples:
  pvm install 3.78.1    Install Pulumi version 3.78.1
  pvm use 3.78.1        Switch to Pulumi version 3.78.1
  pvm list              List installed versions
  pvm list --all        List all available versions
  pvm current           Show current version`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().Bool("no-color", false, "Disable color output")
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		noColor, _ := cmd.Flags().GetBool("no-color")
		if noColor {
			utils.DisableColors()
		}
	}

	rootCmd.AddCommand(installCmd())
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(removeCmd)
}
