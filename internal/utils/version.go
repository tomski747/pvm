package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tomski747/pvm/internal/config"
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
func UseVersion(version string) error {
	// Verify version is installed
	versionsPath := config.GetVersionsPath()
	versionDir := filepath.Join(versionsPath, version)
	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		return fmt.Errorf("version %s is not installed", version)
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
func InstallVersion(version string) error {
	versionsPath := config.GetVersionsPath()
	if err := os.MkdirAll(versionsPath, 0755); err != nil {
		return fmt.Errorf("failed to create versions directory: %v", err)
	}

	goos, arch := config.GetPlatformInfo()
	versionDir := filepath.Join(versionsPath, version)
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		return fmt.Errorf("failed to create version directory: %v", err)
	}

	// Determine download URL based on platform
	var downloadURL string
	if goos == "windows" {
		downloadURL = fmt.Sprintf(config.GithubZipURL, version, version, goos, arch)
	} else {
		downloadURL = fmt.Sprintf(config.GithubReleaseURL, version, version, goos, arch)
	}

	// Download and extract
	if err := downloadAndExtract(downloadURL, versionDir, goos == "windows"); err != nil {
		os.RemoveAll(versionDir) // Clean up on failure
		return fmt.Errorf("failed to download and extract: %v", err)
	}

	return nil
} 