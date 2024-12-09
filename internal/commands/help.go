package commands

import (
	"github.com/spf13/cobra"
)

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Help about any command",
	Long: `Help provides help for any command in the application.
Simply type pvm help [path to command] for full details.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var helpText = `PVM - Pulumi Version Manager

Usage:
  pvm [command]

Available Commands:
  install     Install a specific version of Pulumi (use 'latest' for most recent)
  use        Switch to a specific version of Pulumi
  list       List available Pulumi versions
  current    Show current active version
  help       Show help information

Examples:
  # Install the latest version
  pvm install latest

  # Install and switch to a specific version
  pvm install 3.91.1 --use

  # Switch to an installed version
  pvm use 3.91.1

  # Show current version
  pvm current

For more information, visit: https://github.com/tomski747/pvm
`
