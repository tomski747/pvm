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

// Exported variables for testing
var (
	InstallVersion   = installVersion
	UseVersion       = useVersion
	GetLatestVersion = getLatestVersion
)

// GetInstalledVersions returns a map of installed versions
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

// GetCurrentVersion returns the currently active version
func GetCurrentVersion() (string, error) {
	binPath := filepath.Join(config.GetBinPath(), config.PulumiBinary)
	linkTarget, err := os.Readlink(binPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	// Extract version from the symlink path
	// Path format is ~/.pvm/versions/<version>/pulumi
	parts := strings.Split(linkTarget, string(filepath.Separator))
	for i, part := range parts {
		if part == config.VersionsDir && i+1 < len(parts) {
			return parts[i+1], nil
		}
	}

	return "", nil
}

// UseVersion switches to a specific version of Pulumi
func useVersion(version string) error {
	resolvedVersion, err := ResolveVersion(version)
	if err != nil {
		return err
	}

	// Verify version is installed
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

	// Create new symlinks
	versionBinDir := filepath.Join(versionDir)
	binFiles, err := os.ReadDir(versionBinDir)
	if err != nil {
		return fmt.Errorf("failed to read version directory: %v", err)
	}

	for _, file := range binFiles {
		if !file.IsDir() {
			sourcePath := filepath.Join(versionBinDir, file.Name())
			symlinkPath := filepath.Join(binPath, file.Name())

			if err := os.Symlink(sourcePath, symlinkPath); err != nil {
				return fmt.Errorf("failed to create symlink for %s: %v", file.Name(), err)
			}
		}
	}

	return nil
}

// InstallVersion installs a specific version of Pulumi
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

	// Determine download URL based on platform
	var downloadURL string

	// Convert amd64 to x64
	if arch == "amd64" {
		arch = "x64"
	}

	if goos == "windows" {
		downloadURL = fmt.Sprintf(config.GithubZipURL, resolvedVersion, resolvedVersion, goos, arch)
	} else {
		downloadURL = fmt.Sprintf(config.GithubReleaseURL, resolvedVersion, resolvedVersion, goos, arch)
	}

	// Download and extract
	if err := downloadAndExtract(downloadURL, versionDir, goos == "windows"); err != nil {
		os.RemoveAll(versionDir) // Clean up on failure
		return fmt.Errorf("failed to download and extract: %v", err)
	}

	return nil
}

func getLatestVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/pulumi/pulumi/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var release struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return strings.TrimPrefix(release.TagName, "v"), nil
}

// RemoveVersion removes a specific version of Pulumi
func RemoveVersion(version string) error {
	// Check if version is currently in use
	current, err := GetCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to check current version: %w", err)
	}
	if current == version {
		return fmt.Errorf("cannot remove version %s: currently in use", version)
	}

	versionsPath := config.GetVersionsPath()
	versionDir := filepath.Join(versionsPath, version)

	// Check if version exists
	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		return fmt.Errorf("version %s is not installed", version)
	}

	// Remove the version directory
	if err := os.RemoveAll(versionDir); err != nil {
		return fmt.Errorf("failed to remove version %s: %w", version, err)
	}

	return nil
}

// GetAvailableVersions fetches all available Pulumi versions from GitHub
func GetAvailableVersions(refresh bool) ([]string, error) {
	return FetchGitHubReleases(refresh)
}

// ResolveVersion takes a version or version prefix and returns the full version
func ResolveVersion(versionOrPrefix string) (string, error) {
	versions, err := FetchGitHubReleases(false)
	if err != nil {
		return "", fmt.Errorf("failed to fetch versions: %v", err)
	}

	// If it's an exact version, return it
	for _, v := range versions {
		if v == versionOrPrefix {
			return versionOrPrefix, nil
		}
	}

	// Try to find the latest matching version
	version, err := FindLatestMatchingVersion(versionOrPrefix, versions)
	if err != nil {
		return "", err
	}

	return version, nil
}
