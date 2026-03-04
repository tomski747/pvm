package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tomski747/pvm/internal/config"
)

func TestRemoveCommand(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmpDir, "versions", "3.78.1"), 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	config.SetTestConfig(&config.TestConfig{PVMPath: tmpDir})
	defer config.ResetConfig()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"remove", "3.78.1"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Successfully removed Pulumi 3.78.1") {
		t.Errorf("expected success message, got: %s", buf.String())
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "versions", "3.78.1")); !os.IsNotExist(err) {
		t.Error("expected version directory to be removed")
	}
}

func TestRemoveCommandNotInstalled(t *testing.T) {
	tmpDir := t.TempDir()
	config.SetTestConfig(&config.TestConfig{PVMPath: tmpDir})
	defer config.ResetConfig()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"remove", "9.9.9"})

	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for non-installed version, got nil")
	}
}

func TestRemoveCommandMissingArg(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"remove"})

	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing argument, got nil")
	}
}
