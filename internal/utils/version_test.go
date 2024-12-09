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

func TestGetInstalledVersions(t *testing.T) {
	// Create temporary test directory
	tmpDir, err := os.MkdirTemp("", "pvm-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test versions
	testVersions := []string{"3.78.1", "3.78.0"}
	for _, version := range testVersions {
		versionDir := filepath.Join(tmpDir, "versions", version)
		if err := os.MkdirAll(versionDir, 0755); err != nil {
			t.Fatalf("Failed to create version dir: %v", err)
		}
	}

	// Create test environment
	testConfig := &config.TestConfig{
		PVMPath: tmpDir,
	}
	config.SetTestConfig(testConfig)
	defer config.ResetConfig()

	// Test GetInstalledVersions
	installed := GetInstalledVersions()
	for _, version := range testVersions {
		if !installed[version] {
			t.Errorf("Expected version %s to be installed", version)
		}
	}
}

func TestGetCurrentVersion(t *testing.T) {
	// Create temporary test directory
	tmpDir, err := os.MkdirTemp("", "pvm-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test environment
	testConfig := &config.TestConfig{
		PVMPath: tmpDir,
	}
	config.SetTestConfig(testConfig)
	defer config.ResetConfig()

	// Test when no version is set
	version, err := GetCurrentVersion()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if version != "" {
		t.Errorf("Expected empty version, got: %s", version)
	}
}

func TestGetLatestVersion(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/pulumi/pulumi/releases/latest" {
			t.Errorf("Expected path /repos/pulumi/pulumi/releases/latest, got %s", r.URL.Path)
		}
		release := struct {
			TagName string `json:"tag_name"`
		}{
			TagName: "v3.78.1",
		}
		if err := json.NewEncoder(w).Encode(release); err != nil {
			t.Errorf("Failed to encode release: %v", err)
		}
	}))
	defer server.Close()

	// Override GitHub API URL for testing
	originalURL := "https://api.github.com"
	githubAPIURL = server.URL
	defer func() { githubAPIURL = originalURL }()

	// Test getting latest version
	version, err := GetLatestVersion()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedVersion := "3.142.0"
	if version != expectedVersion {
		t.Errorf("Expected version %s, got %s", expectedVersion, version)
	}
}
