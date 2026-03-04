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

// resetListFlags resets listCmd's boolean flags to their defaults between tests.
// pflag does not automatically reset flag values between cobra Execute calls.
func resetListFlags() {
	refresh = false
	_ = listCmd.Flags().Set("all", "false")
}

func TestListCommandEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	config.SetTestConfig(&config.TestConfig{PVMPath: tmpDir})
	defer config.ResetConfig()
	resetListFlags()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"list"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "No versions installed") {
		t.Errorf("expected no-versions message, got: %s", buf.String())
	}
}

func TestListCommandWithInstalled(t *testing.T) {
	tmpDir := t.TempDir()
	for _, v := range []string{"3.78.1", "3.77.0"} {
		if err := os.MkdirAll(filepath.Join(tmpDir, "versions", v), 0755); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}
	config.SetTestConfig(&config.TestConfig{PVMPath: tmpDir})
	defer config.ResetConfig()
	resetListFlags()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"list"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "3.78.1") {
		t.Errorf("expected 3.78.1 in output, got: %s", out)
	}
	if !strings.Contains(out, "3.77.0") {
		t.Errorf("expected 3.77.0 in output, got: %s", out)
	}
}

func TestListCommandAll(t *testing.T) {
	tmpDir := t.TempDir()
	config.SetTestConfig(&config.TestConfig{PVMPath: tmpDir})
	defer config.ResetConfig()

	cleanup := utils.MockVersionOperations(t)
	defer cleanup()
	resetListFlags()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"list", "--all"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	// MockVersionOperations returns {"3.78.1", "3.78.0", "3.77.0"}
	if !strings.Contains(out, "3.78.1") {
		t.Errorf("expected available version in output, got: %s", out)
	}
	if !strings.Contains(out, "Available versions") {
		t.Errorf("expected 'Available versions' header, got: %s", out)
	}
}

func TestListCommandSortOrder(t *testing.T) {
	tmpDir := t.TempDir()
	// Install versions in non-sorted order
	for _, v := range []string{"3.77.0", "3.100.0", "3.78.1"} {
		if err := os.MkdirAll(filepath.Join(tmpDir, "versions", v), 0755); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}
	config.SetTestConfig(&config.TestConfig{PVMPath: tmpDir})
	defer config.ResetConfig()
	resetListFlags()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"list"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	pos100 := strings.Index(out, "3.100.0")
	pos78 := strings.Index(out, "3.78.1")
	pos77 := strings.Index(out, "3.77.0")

	if pos100 < 0 || pos78 < 0 || pos77 < 0 {
		t.Fatalf("versions missing from output: %s", out)
	}
	// Descending order: 3.100.0 > 3.78.1 > 3.77.0
	if pos100 >= pos78 || pos78 >= pos77 {
		t.Errorf("versions not in descending order; positions: 3.100.0=%d, 3.78.1=%d, 3.77.0=%d\noutput:\n%s",
			pos100, pos78, pos77, out)
	}
}
