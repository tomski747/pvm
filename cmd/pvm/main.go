package main

import (
	"os"

	"github.com/tomski747/pvm/internal/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
