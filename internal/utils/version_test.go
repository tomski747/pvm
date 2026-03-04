package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/tomski747/pvm/internal/config"
)

// setupVersionsDir creates a temp dir with the given versions installed and
// points testConfig at it. Returns a cleanup function.
func setupVersionsDir(t *testing.T, versions []string) string {
	t.Helper()
	tmpDir := t.TempDir()
	for _, v := range versions {
		if err := os.MkdirAll(filepath.Join(tmpDir, "versions", v), 0755); err != nil {
			t.Fatalf("create version dir: %v", err)
		}
	}
	config.SetTestConfig(&config.TestConfig{PVMPath: tmpDir})
	t.Cleanup(config.ResetConfig)
	return tmpDir
}

func TestGetInstalledVersions(t *testing.T) {
	testVersions := []string{"3.78.1", "3.78.0"}
	setupVersionsDir(t, testVersions)

	installed := GetInstalledVersions()
	for _, v := range testVersions {
		if !installed[v] {
			t.Errorf("expected version %s to be installed", v)
		}
	}
}

func TestGetCurrentVersion(t *testing.T) {
	setupVersionsDir(t, nil)

	version, err := GetCurrentVersion()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version != "" {
		t.Errorf("expected empty version, got: %s", version)
	}
}

func TestGetCurrentVersionWithSymlink(t *testing.T) {
	tmpDir := setupVersionsDir(t, []string{"3.78.1"})

	// Create a fake pulumi binary and a symlink in bin/
	binDir := filepath.Join(tmpDir, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		t.Fatalf("create bin dir: %v", err)
	}
	src := filepath.Join(tmpDir, "versions", "3.78.1", "pulumi")
	if err := os.WriteFile(src, []byte("#!/bin/sh"), 0755); err != nil {
		t.Fatalf("create fake binary: %v", err)
	}
	if err := os.Symlink(src, filepath.Join(binDir, "pulumi")); err != nil {
		t.Fatalf("create symlink: %v", err)
	}

	version, err := GetCurrentVersion()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version != "3.78.1" {
		t.Errorf("expected 3.78.1, got: %s", version)
	}
}

func TestGetLatestVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		release := struct {
			TagName string `json:"tag_name"`
		}{TagName: "v3.78.1"}
		_ = json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	orig := githubLatestReleaseURL
	githubLatestReleaseURL = server.URL + "/latest"
	defer func() { githubLatestReleaseURL = orig }()

	version, err := GetLatestVersion()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version != "3.78.1" {
		t.Errorf("expected 3.78.1, got %s", version)
	}
}

func TestRemoveVersion(t *testing.T) {
	tmpDir := setupVersionsDir(t, []string{"3.78.1", "3.78.0"})

	// Remove 3.78.0 (not active)
	if err := RemoveVersion("3.78.0"); err != nil {
		t.Fatalf("RemoveVersion: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "versions", "3.78.0")); !os.IsNotExist(err) {
		t.Error("expected 3.78.0 directory to be gone")
	}
	// 3.78.1 should still be there
	if _, err := os.Stat(filepath.Join(tmpDir, "versions", "3.78.1")); err != nil {
		t.Error("expected 3.78.1 directory to still exist")
	}
}

func TestRemoveVersionNotInstalled(t *testing.T) {
	setupVersionsDir(t, nil)

	if err := RemoveVersion("9.9.9"); err == nil {
		t.Error("expected error for non-installed version, got nil")
	}
}

func TestRemoveVersionCurrentlyActive(t *testing.T) {
	tmpDir := setupVersionsDir(t, []string{"3.78.1"})

	// Make 3.78.1 the active version via a symlink
	binDir := filepath.Join(tmpDir, "bin")
	_ = os.MkdirAll(binDir, 0755)
	src := filepath.Join(tmpDir, "versions", "3.78.1", "pulumi")
	_ = os.WriteFile(src, []byte("#!/bin/sh"), 0755)
	_ = os.Symlink(src, filepath.Join(binDir, "pulumi"))

	if err := RemoveVersion("3.78.1"); err == nil {
		t.Error("expected error when removing active version, got nil")
	}
}

func TestUseVersion(t *testing.T) {
	tmpDir := setupVersionsDir(t, []string{"3.78.1"})

	// Create a fake binary in the version dir
	src := filepath.Join(tmpDir, "versions", "3.78.1", "pulumi")
	if err := os.WriteFile(src, []byte("#!/bin/sh"), 0755); err != nil {
		t.Fatalf("create fake binary: %v", err)
	}

	// Override ResolveVersion so it doesn't hit the network
	orig := ResolveVersion
	ResolveVersion = func(v string) (string, error) { return v, nil }
	defer func() { ResolveVersion = orig }()

	if err := UseVersion("3.78.1"); err != nil {
		t.Fatalf("UseVersion: %v", err)
	}

	// Verify symlink was created
	link := filepath.Join(tmpDir, "bin", "pulumi")
	target, err := os.Readlink(link)
	if err != nil {
		t.Fatalf("readlink: %v", err)
	}
	if target != src {
		t.Errorf("symlink target = %s, want %s", target, src)
	}

	// GetCurrentVersion should now return 3.78.1
	version, err := GetCurrentVersion()
	if err != nil {
		t.Fatalf("GetCurrentVersion: %v", err)
	}
	if version != "3.78.1" {
		t.Errorf("expected current version 3.78.1, got %s", version)
	}
}

func TestResolveVersionExact(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		releases := []githubRelease{{TagName: "v3.78.1"}, {TagName: "v3.78.0"}}
		_ = json.NewEncoder(w).Encode(releases)
	}))
	defer server.Close()

	origURL := githubAPIURL
	githubAPIURL = server.URL
	defer func() { githubAPIURL = origURL }()

	// Use a fresh temp dir so cache misses and hits the server
	tmpDir := t.TempDir()
	config.SetTestConfig(&config.TestConfig{PVMPath: tmpDir})
	defer config.ResetConfig()

	version, err := resolveVersion("3.78.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version != "3.78.1" {
		t.Errorf("expected 3.78.1, got %s", version)
	}
}

func TestResolveVersionPrefix(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		releases := []githubRelease{
			{TagName: "v3.78.1"},
			{TagName: "v3.78.0"},
			{TagName: "v3.77.5"},
		}
		_ = json.NewEncoder(w).Encode(releases)
	}))
	defer server.Close()

	origURL := githubAPIURL
	githubAPIURL = server.URL
	defer func() { githubAPIURL = origURL }()

	tmpDir := t.TempDir()
	config.SetTestConfig(&config.TestConfig{PVMPath: tmpDir})
	defer config.ResetConfig()

	version, err := resolveVersion("3.78")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should resolve to the latest 3.78.x
	if version != "3.78.1" {
		t.Errorf("expected 3.78.1, got %s", version)
	}
}
