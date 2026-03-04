package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tomski747/pvm/internal/config"
)

func TestCurrentCommandNoVersion(t *testing.T) {
	tmpDir := t.TempDir()
	config.SetTestConfig(&config.TestConfig{PVMPath: tmpDir})
	defer config.ResetConfig()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"current"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "No Pulumi version currently selected") {
		t.Errorf("expected no-version message, got: %s", buf.String())
	}
}

func TestCurrentCommandWithVersion(t *testing.T) {
	tmpDir := t.TempDir()
	// Set up a fake active version symlink
	binDir := filepath.Join(tmpDir, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	versionDir := filepath.Join(tmpDir, "versions", "3.78.1")
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	src := filepath.Join(versionDir, "pulumi")
	if err := os.WriteFile(src, []byte("#!/bin/sh"), 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := os.Symlink(src, filepath.Join(binDir, "pulumi")); err != nil {
		t.Fatalf("setup symlink: %v", err)
	}

	config.SetTestConfig(&config.TestConfig{PVMPath: tmpDir})
	defer config.ResetConfig()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"current"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "3.78.1") {
		t.Errorf("expected version 3.78.1 in output, got: %s", out)
	}
}
