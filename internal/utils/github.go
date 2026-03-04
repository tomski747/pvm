package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/tomski747/pvm/internal/config"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
}

// githubAPIURL can be overridden in tests.
var githubAPIURL = "https://api.github.com/repos/pulumi/pulumi/releases"

// FetchGitHubReleases fetches all available Pulumi versions from GitHub.
func FetchGitHubReleases(refresh bool) ([]string, error) {
	// If refresh is true, skip cache and fetch directly from GitHub
	if !refresh {
		if versions, err := readCache(); err == nil {
			return versions, nil
		}
	}

	versions, err := fetchFromGitHub()
	if err != nil {
		return nil, err
	}

	if err := saveCache(versions); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to save cache: %v\n", err)
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

	sort.Slice(cache.Versions, func(i, j int) bool {
		return SemverGreater(cache.Versions[i], cache.Versions[j])
	})

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
	var versions []string
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

		if !strings.Contains(resp.Header.Get("Link"), `rel="next"`) {
			break
		}

		page++
	}

	return versions, nil
}

// FindLatestMatchingVersion finds the latest version that matches the given prefix.
func FindLatestMatchingVersion(prefix string, versions []string) (string, error) {
	if prefix == "" {
		return "", fmt.Errorf("version prefix cannot be empty")
	}

	parts := strings.Split(prefix, ".")

	var matchingVersions []string
	for _, version := range versions {
		if len(parts) == 1 {
			// For a single segment like "3", match any version starting with "3."
			if strings.HasPrefix(version+".", parts[0]+".") {
				matchingVersions = append(matchingVersions, version)
			}
		} else {
			// For partial versions like "3.1", match exact prefix
			if strings.HasPrefix(version+".", prefix+".") {
				matchingVersions = append(matchingVersions, version)
			}
		}
	}

	if len(matchingVersions) == 0 {
		return "", fmt.Errorf("no versions found matching prefix %s", prefix)
	}

	sort.Slice(matchingVersions, func(i, j int) bool {
		return SemverGreater(matchingVersions[i], matchingVersions[j])
	})

	return matchingVersions[0], nil
}
