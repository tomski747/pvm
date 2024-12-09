package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tomski747/pvm/internal/config"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
}

// For testing purposes, this can be overridden
var githubAPIURL = "https://api.github.com/repos/pulumi/pulumi/releases"

// FetchGitHubReleases fetches all available Pulumi versions from GitHub
func FetchGitHubReleases() ([]string, error) {
	// Try cache first
	if versions, err := readCache(); err == nil {
		return versions, nil
	}

	// Fetch from GitHub
	versions, err := fetchFromGitHub()
	if err != nil {
		return nil, err
	}

	// Save to cache
	if err := saveCache(versions); err != nil {
		fmt.Printf("Warning: Failed to save cache: %v\n", err)
	}

	return versions, nil
}

func readCache() ([]string, error) {
	cachePath := filepath.Join(config.GetPVMPath(), config.CacheFile)
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}

	var cache config.ReleaseCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}

	if time.Since(cache.Timestamp) > config.CacheTTL {
		return nil, fmt.Errorf("cache expired")
	}

	return cache.Versions, nil
}

func saveCache(versions []string) error {
	cache := config.ReleaseCache{
		Versions:  versions,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	cachePath := filepath.Join(config.GetPVMPath(), config.CacheFile)
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0644)
}

func fetchFromGitHub() ([]string, error) {
	client := &http.Client{}
	versions := []string{}
	page := 1
	perPage := 100

	for {
		url := fmt.Sprintf("%s?page=%d&per_page=%d", githubAPIURL, page, perPage)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}

		req.Header.Set("Accept", "application/vnd.github.v3+json")
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error fetching releases: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
		}

		var releases []githubRelease
		if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
			return nil, fmt.Errorf("error decoding response: %v", err)
		}

		if len(releases) == 0 {
			break
		}

		for _, release := range releases {
			versions = append(versions, strings.TrimPrefix(release.TagName, "v"))
		}

		linkHeader := resp.Header.Get("Link")
		if !strings.Contains(linkHeader, `rel="next"`) {
			break
		}

		page++
	}

	return versions, nil
} 