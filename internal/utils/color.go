package utils

import (
	"os"

	"github.com/fatih/color"
)

var (
	useColor = true
	// Define colors
	Success = color.New(color.FgGreen).SprintfFunc()
	Info    = color.New(color.FgCyan).SprintfFunc()
	Warning = color.New(color.FgYellow).SprintfFunc()
	Error   = color.New(color.FgRed).SprintfFunc()
	Current = color.New(color.FgGreen, color.Bold).SprintfFunc()
)

func init() {
	// Disable colors if NO_COLOR is set
	if os.Getenv("NO_COLOR") != "" {
		DisableColors()
	}
}

func DisableColors() {
	useColor = false
	color.NoColor = true
}

func EnableColors() {
	useColor = true
	color.NoColor = false
} 