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