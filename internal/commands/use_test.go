package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tomski747/pvm/internal/config"
	"github.com/tomski747/pvm/internal/utils"
)

func TestUseCommand(t *testing.T) {
	tmpDir := t.TempDir()
	// Create the version directory so GetInstalledVersions finds it
	if err := os.MkdirAll(filepath.Join(tmpDir, "versions", "3.78.1"), 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	config.SetTestConfig(&config.TestConfig{PVMPath: tmpDir})
	defer config.ResetConfig()

	cleanup := utils.MockVersionOperations(t)
	defer cleanup()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"use", "3.78.1"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Switched to Pulumi") {
		t.Errorf("expected switch message, got: %s", buf.String())
	}
}

func TestUseCommandNotInstalled(t *testing.T) {
	tmpDir := t.TempDir()
	config.SetTestConfig(&config.TestConfig{PVMPath: tmpDir})
	defer config.ResetConfig()

	cleanup := utils.MockVersionOperations(t)
	defer cleanup()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"use", "9.9.9"})

	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for uninstalled version, got nil")
	}
}

func TestUseCommandWithInstallFlag(t *testing.T) {
	tmpDir := t.TempDir()
	config.SetTestConfig(&config.TestConfig{PVMPath: tmpDir})
	defer config.ResetConfig()

	cleanup := utils.MockVersionOperations(t)
	defer cleanup()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"use", "3.78.1", "--install"})

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

func TestUseCommandLatest(t *testing.T) {
	tmpDir := t.TempDir()
	// Mock GetLatestVersion returns "3.78.1"; pre-install it
	if err := os.MkdirAll(filepath.Join(tmpDir, "versions", "3.78.1"), 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	config.SetTestConfig(&config.TestConfig{PVMPath: tmpDir})
	defer config.ResetConfig()

	cleanup := utils.MockVersionOperations(t)
	defer cleanup()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"use", "latest"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "3.78.1") {
		t.Errorf("expected version in output, got: %s", buf.String())
	}
}
