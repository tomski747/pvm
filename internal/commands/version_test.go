package commands

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"version"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// config.Version defaults to "dev" in tests
	if !strings.Contains(buf.String(), "dev") {
		t.Errorf("expected version string in output, got: %s", buf.String())
	}
}
