package config

import (
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	PVMDir           = ".pvm"
	VersionsDir      = "versions"
	BinDir           = "bin"
	PulumiBinary     = "pulumi"
	GithubReleaseURL = "https://github.com/pulumi/pulumi/releases/download/v%s/pulumi-v%s-%s-%s.tar.gz"
	GithubZipURL     = "https://github.com/pulumi/pulumi/releases/download/v%s/pulumi-v%s-%s-%s.zip"
	CacheFile        = "releases.cache"
	CacheTTL         = 24 * time.Hour
)

type ReleaseCache struct {
	Versions  []string  `json:"versions"`
	Timestamp time.Time `json:"timestamp"`
}

// TestConfig holds configuration for testing
type TestConfig struct {
	PVMPath string
}

var testConfig *TestConfig

// SetTestConfig sets the test configuration
func SetTestConfig(cfg *TestConfig) {
	testConfig = cfg
}

// ResetConfig resets the test configuration
func ResetConfig() {
	testConfig = nil
}

// GetHomeDir returns the user's home directory
func GetHomeDir() string {
	if testConfig != nil {
		return testConfig.PVMPath
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return os.Getenv("HOME")
	}
	return home
}

// GetPVMPath returns the PVM root directory path
func GetPVMPath() string {
	if testConfig != nil {
		return testConfig.PVMPath
	}
	return filepath.Join(GetHomeDir(), PVMDir)
}

// GetVersionsPath returns the versions directory path
func GetVersionsPath() string {
	return filepath.Join(GetPVMPath(), VersionsDir)
}

// GetBinPath returns the bin directory path
func GetBinPath() string {
	return filepath.Join(GetPVMPath(), BinDir)
}

// GetPlatformInfo returns the current OS and architecture
func GetPlatformInfo() (string, string) {
	return runtime.GOOS, runtime.GOARCH
}
