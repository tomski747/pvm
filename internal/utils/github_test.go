package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchFromGitHub(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/pulumi/pulumi/releases" {
			t.Errorf("Expected path /repos/pulumi/pulumi/releases, got %s", r.URL.Path)
		}
		releases := []githubRelease{
			{TagName: "v3.78.1"},
			{TagName: "v3.78.0"},
		}
		json.NewEncoder(w).Encode(releases)
	}))
	defer server.Close()

	// Override GitHub API URL for testing
	originalURL := githubAPIURL
	githubAPIURL = server.URL + "/repos/pulumi/pulumi/releases"
	defer func() { githubAPIURL = originalURL }()

	// Test fetching versions
	versions, err := fetchFromGitHub()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedVersions := []string{"3.78.1", "3.78.0"}
	if len(versions) != len(expectedVersions) {
		t.Errorf("Expected %d versions, got %d", len(expectedVersions), len(versions))
	}

	for i, version := range versions {
		if version != expectedVersions[i] {
			t.Errorf("Expected version %s, got %s", expectedVersions[i], version)
		}
	}
} 