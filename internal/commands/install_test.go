package commands

import (
	"bytes"
	"strings"
	"testing"

	"github.com/tomski747/pvm/internal/utils"
)

func TestInstallCommand(t *testing.T) {
	cleanup := utils.MockVersionOperations(t)
	defer cleanup()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"install", "3.78.1"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Successfully installed Pulumi") {
		t.Errorf("expected success message in output, got: %s", out)
	}
	if !strings.Contains(out, "pvm use 3.78.1") {
		t.Errorf("expected use hint in output, got: %s", out)
	}
}

func TestInstallCommandWithUseFlag(t *testing.T) {
	cleanup := utils.MockVersionOperations(t)
	defer cleanup()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"install", "3.78.1", "--use"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Successfully installed Pulumi") {
		t.Errorf("expected install message, got: %s", out)
	}
	if !strings.Contains(out, "Switched to Pulumi") {
		t.Errorf("expected switch message, got: %s", out)
	}
}

func TestInstallCommandLatest(t *testing.T) {
	cleanup := utils.MockVersionOperations(t)
	defer cleanup()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"install", "latest"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	// GetLatestVersion mock returns "3.78.1"
	if !strings.Contains(out, "3.78.1") {
		t.Errorf("expected resolved version in output, got: %s", out)
	}
}

func TestInstallCommandMissingArg(t *testing.T) {
	cleanup := utils.MockVersionOperations(t)
	defer cleanup()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"install"})

	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing version argument, got nil")
	}
}
