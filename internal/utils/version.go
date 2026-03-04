package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/tomski747/pvm/internal/config"
)

// Function variables allow tests to inject mocks without build tags.
var (
	InstallVersion         = installVersion
	UseVersion             = useVersion
	GetLatestVersion       = getLatestVersion
	ResolveVersion         = resolveVersion
	GetAvailableVersions   = getAvailableVersions
	githubLatestReleaseURL = "https://api.github.com/repos/pulumi/pulumi/releases/latest"
)

// GetInstalledVersions returns a map of installed versions.
func GetInstalledVersions() map[string]bool {
	installed := make(map[string]bool)
	versionsPath := config.GetVersionsPath()

	files, err := os.ReadDir(versionsPath)
	if err != nil && !os.IsNotExist(err) {
		return installed
	}

	for _, file := range files {
		if file.IsDir() {
			installed[file.Name()] = true
		}
	}

	return installed
}

// GetCurrentVersion returns the currently active version.
func GetCurrentVersion() (string, error) {
	binPath := filepath.Join(config.GetBinPath(), config.PulumiBinary)
	linkTarget, err := os.Readlink(binPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	// Extract version from the symlink path: ~/.pvm/versions/<version>/pulumi
	parts := strings.Split(linkTarget, string(filepath.Separator))
	for i, part := range parts {
		if part == config.VersionsDir && i+1 < len(parts) {
			return parts[i+1], nil
		}
	}

	return "", nil
}

func useVersion(version string) error {
	resolvedVersion, err := ResolveVersion(version)
	if err != nil {
		return err
	}

	versionsPath := config.GetVersionsPath()
	versionDir := filepath.Join(versionsPath, resolvedVersion)
	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		return fmt.Errorf("version %s is not installed", resolvedVersion)
	}

	binPath := config.GetBinPath()
	if err := os.MkdirAll(binPath, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %v", err)
	}

	// Remove existing symlinks
	files, err := os.ReadDir(binPath)
	if err != nil {
		return fmt.Errorf("failed to read bin directory: %v", err)
	}

	for _, file := range files {
		filePath := filepath.Join(binPath, file.Name())
		if fileInfo, err := os.Lstat(filePath); err == nil {
			if fileInfo.Mode()&os.ModeSymlink != 0 {
				if err := os.Remove(filePath); err != nil {
					return fmt.Errorf("failed to remove existing symlink %s: %v", file.Name(), err)
				}
			}
		}
	}

	// Create new symlinks for every file in the version directory
	binFiles, err := os.ReadDir(versionDir)
	if err != nil {
		return fmt.Errorf("failed to read version directory: %v", err)
	}

	for _, file := range binFiles {
		if !file.IsDir() {
			sourcePath := filepath.Join(versionDir, file.Name())
			symlinkPath := filepath.Join(binPath, file.Name())

			if err := os.Symlink(sourcePath, symlinkPath); err != nil {
				return fmt.Errorf("failed to create symlink for %s: %v", file.Name(), err)
			}
		}
	}

	return nil
}

func installVersion(version string) error {
	resolvedVersion, err := ResolveVersion(version)
	if err != nil {
		return err
	}

	versionsPath := config.GetVersionsPath()
	if err := os.MkdirAll(versionsPath, 0755); err != nil {
		return fmt.Errorf("failed to create versions directory: %v", err)
	}

	goos, arch := config.GetPlatformInfo()
	versionDir := filepath.Join(versionsPath, resolvedVersion)
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		return fmt.Errorf("failed to create version directory: %v", err)
	}

	// Pulumi's release naming uses "x64" instead of "amd64"
	if arch == "amd64" {
		arch = "x64"
	}

	var downloadURL string
	if goos == "windows" {
		downloadURL = fmt.Sprintf(config.GithubZipURL, resolvedVersion, resolvedVersion, goos, arch)
	} else {
		downloadURL = fmt.Sprintf(config.GithubReleaseURL, resolvedVersion, resolvedVersion, goos, arch)
	}

	if err := downloadAndExtract(downloadURL, versionDir, goos == "windows"); err != nil {
		os.RemoveAll(versionDir) // clean up partial download
		return fmt.Errorf("failed to download and extract: %v", err)
	}

	return nil
}

func getLatestVersion() (string, error) {
	resp, err := http.Get(githubLatestReleaseURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return strings.TrimPrefix(release.TagName, "v"), nil
}

// RemoveVersion removes a specific version of Pulumi.
func RemoveVersion(version string) error {
	current, err := GetCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to check current version: %w", err)
	}
	if current == version {
		return fmt.Errorf("cannot remove version %s: currently in use", version)
	}

	versionsPath := config.GetVersionsPath()
	versionDir := filepath.Join(versionsPath, version)

	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		return fmt.Errorf("version %s is not installed", version)
	}

	if err := os.RemoveAll(versionDir); err != nil {
		return fmt.Errorf("failed to remove version %s: %w", version, err)
	}

	return nil
}

func getAvailableVersions(refresh bool) ([]string, error) {
	return FetchGitHubReleases(refresh)
}

func resolveVersion(versionOrPrefix string) (string, error) {
	versions, err := FetchGitHubReleases(false)
	if err != nil {
		return "", fmt.Errorf("failed to fetch versions: %v", err)
	}

	// Exact match first
	for _, v := range versions {
		if v == versionOrPrefix {
			return versionOrPrefix, nil
		}
	}

	// Fall back to prefix match
	return FindLatestMatchingVersion(versionOrPrefix, versions)
}
