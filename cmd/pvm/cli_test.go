// Package main contains integration tests that build the pvm binary and
// exercise it as a subprocess. These tests verify the end-to-end CLI
// behaviour (argument parsing, exit codes, output formatting) without
// mocking any internal packages.
package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// binaryPath is the compiled pvm binary used by all integration tests.
var binaryPath string

// TestMain builds the binary once for the whole test suite.
func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp("", "pvm-integ-*")
	if err != nil {
		panic("failed to create temp dir for binary: " + err.Error())
	}
	defer os.RemoveAll(tmpDir)

	binaryPath = filepath.Join(tmpDir, "pvm")
	build := exec.Command("go", "build", "-o", binaryPath, ".")
	if out, err := build.CombinedOutput(); err != nil {
		panic("failed to build pvm binary: " + err.Error() + "\n" + string(out))
	}

	os.Exit(m.Run())
}

// runPVMInDir executes the binary inside pvmHome (PVM_HOME is set to pvmHome).
func runPVMInDir(pvmHome string, args ...string) (string, int) {
	cmd := exec.Command(binaryPath, args...)
	cmd.Env = append(os.Environ(), "PVM_HOME="+pvmHome, "NO_COLOR=1")
	out, err := cmd.CombinedOutput()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}
	return string(out), exitCode
}

// runPVM executes the binary with a fresh isolated PVM_HOME for each test.
func runPVM(t *testing.T, args ...string) (string, int) {
	t.Helper()
	return runPVMInDir(t.TempDir(), args...)
}

// releaseCacheEntry mirrors config.ReleaseCache for JSON serialisation.
type releaseCacheEntry struct {
	Versions  []string  `json:"versions"`
	Timestamp time.Time `json:"timestamp"`
}

// primeCache writes a releases.cache file into pvmHome so the binary never
// needs to contact GitHub to resolve versions during integration tests.
func primeCache(t *testing.T, pvmHome string, versions []string) {
	t.Helper()
	cache := releaseCacheEntry{
		Versions:  versions,
		Timestamp: time.Now(),
	}
	data, err := json.Marshal(cache)
	if err != nil {
		t.Fatalf("marshal cache: %v", err)
	}
	if err := os.WriteFile(filepath.Join(pvmHome, "releases.cache"), data, 0644); err != nil {
		t.Fatalf("write cache: %v", err)
	}
}

// ── Basic invocations ────────────────────────────────────────────────────────

func TestCLIHelp(t *testing.T) {
	out, code := runPVM(t, "--help")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d\noutput: %s", code, out)
	}
	for _, phrase := range []string{"Usage:", "Available Commands:", "install", "use", "list", "current"} {
		if !strings.Contains(out, phrase) {
			t.Errorf("expected %q in help output, got:\n%s", phrase, out)
		}
	}
}

func TestCLINoArgs(t *testing.T) {
	out, code := runPVM(t)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d\noutput: %s", code, out)
	}
	if !strings.Contains(out, "Usage:") {
		t.Errorf("expected usage in output, got:\n%s", out)
	}
}

func TestCLIVersion(t *testing.T) {
	out, code := runPVM(t, "version")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d\noutput: %s", code, out)
	}
	// Binary is built without -ldflags, so version is the default "dev"
	if strings.TrimSpace(out) == "" {
		t.Errorf("expected non-empty version output")
	}
}

func TestCLICurrentNoVersion(t *testing.T) {
	out, code := runPVM(t, "current")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d\noutput: %s", code, out)
	}
	if !strings.Contains(out, "No Pulumi version currently selected") {
		t.Errorf("expected no-version message, got:\n%s", out)
	}
}

func TestCLIListEmpty(t *testing.T) {
	out, code := runPVM(t, "list")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d\noutput: %s", code, out)
	}
	if !strings.Contains(out, "No versions installed") {
		t.Errorf("expected empty-list message, got:\n%s", out)
	}
}

// ── Error paths ──────────────────────────────────────────────────────────────

func TestCLIInstallMissingArg(t *testing.T) {
	_, code := runPVM(t, "install")
	if code == 0 {
		t.Error("expected non-zero exit for missing argument, got 0")
	}
}

func TestCLIRemoveMissingArg(t *testing.T) {
	_, code := runPVM(t, "remove")
	if code == 0 {
		t.Error("expected non-zero exit for missing argument, got 0")
	}
}

func TestCLIUseMissingArg(t *testing.T) {
	_, code := runPVM(t, "use")
	if code == 0 {
		t.Error("expected non-zero exit for missing argument, got 0")
	}
}

func TestCLIUseNotInstalled(t *testing.T) {
	pvmHome := t.TempDir()
	// Prime the cache so the binary resolves "3.78.1" without hitting GitHub.
	primeCache(t, pvmHome, []string{"3.78.1", "3.78.0"})

	out, code := runPVMInDir(pvmHome, "use", "3.78.1")
	if code == 0 {
		t.Errorf("expected non-zero exit for uninstalled version, got 0\noutput: %s", out)
	}
	if !strings.Contains(out, "not installed") {
		t.Errorf("expected 'not installed' message, got:\n%s", out)
	}
}

func TestCLIRemoveNotInstalled(t *testing.T) {
	pvmHome := t.TempDir()
	primeCache(t, pvmHome, []string{"3.78.1"})

	out, code := runPVMInDir(pvmHome, "remove", "3.78.1")
	if code == 0 {
		t.Errorf("expected non-zero exit for uninstalled version, got 0\noutput: %s", out)
	}
}

// ── Feature tests ─────────────────────────────────────────────────────────────

func TestCLIListInstalledVersions(t *testing.T) {
	pvmHome := t.TempDir()
	// Create fake installed versions
	for _, v := range []string{"3.78.1", "3.77.0"} {
		dir := filepath.Join(pvmHome, "versions", v)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	out, code := runPVMInDir(pvmHome, "list")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d\noutput: %s", code, out)
	}
	for _, v := range []string{"3.78.1", "3.77.0"} {
		if !strings.Contains(out, v) {
			t.Errorf("expected %s in list output, got:\n%s", v, out)
		}
	}
}

func TestCLICurrentAfterUse(t *testing.T) {
	pvmHome := t.TempDir()
	// Create a fake version with a pulumi binary and prime cache
	versionDir := filepath.Join(pvmHome, "versions", "3.78.1")
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	fakeBin := filepath.Join(versionDir, "pulumi")
	if err := os.WriteFile(fakeBin, []byte("#!/bin/sh\necho 3.78.1"), 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	primeCache(t, pvmHome, []string{"3.78.1"})

	// Switch to the version
	if out, code := runPVMInDir(pvmHome, "use", "3.78.1"); code != 0 {
		t.Fatalf("pvm use 3.78.1 failed (exit %d): %s", code, out)
	}

	// Now current should report 3.78.1
	out, code := runPVMInDir(pvmHome, "current")
	if code != 0 {
		t.Fatalf("pvm current failed (exit %d): %s", code, out)
	}
	if !strings.Contains(out, "3.78.1") {
		t.Errorf("expected current version 3.78.1, got:\n%s", out)
	}
}

func TestCLIRemoveVersion(t *testing.T) {
	pvmHome := t.TempDir()
	versionDir := filepath.Join(pvmHome, "versions", "3.77.0")
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	out, code := runPVMInDir(pvmHome, "remove", "3.77.0")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d\noutput: %s", code, out)
	}
	if !strings.Contains(out, "Successfully removed") {
		t.Errorf("expected success message, got:\n%s", out)
	}
	if _, err := os.Stat(versionDir); !os.IsNotExist(err) {
		t.Error("expected version directory to be gone after remove")
	}
}

// ── Subcommand help ───────────────────────────────────────────────────────────

func TestCLISubcommandHelp(t *testing.T) {
	// Cobra's built-in help supports `pvm help <subcommand>` properly.
	for _, sub := range []string{"install", "use", "list", "current", "remove"} {
		out, code := runPVM(t, "help", sub)
		if code != 0 {
			t.Errorf("pvm help %s: expected exit 0, got %d\noutput: %s", sub, code, out)
			continue
		}
		if !strings.Contains(out, "Usage:") {
			t.Errorf("pvm help %s: expected Usage: in output, got:\n%s", sub, out)
		}
	}
}

func TestCLINoColorFlag(t *testing.T) {
	out, code := runPVM(t, "--no-color", "current")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d\noutput: %s", code, out)
	}
	if strings.Contains(out, "\x1b[") {
		t.Errorf("expected no ANSI escape codes with --no-color, got:\n%s", out)
	}
}
