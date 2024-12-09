package utils

import (
	"testing"
)

// Mock function types
type (
	mockInstallVersionFunc func(version string) error
	mockUseVersionFunc    func(version string) error
	mockGetLatestVersionFunc func() (string, error)
)

// Mock function variables
var (
	mockInstallVersionFn mockInstallVersionFunc
	mockUseVersionFn    mockUseVersionFunc  
	mockGetLatestVersionFn mockGetLatestVersionFunc
)

// MockVersionOperations sets up mock functions for version operations and returns cleanup function
func MockVersionOperations(t testing.TB) func() {
	originalInstallVersion := InstallVersion
	originalUseVersion := UseVersion
	originalGetLatestVersion := GetLatestVersion

	InstallVersion = func(version string) error {
		if mockInstallVersionFn != nil {
			return mockInstallVersionFn(version)
		}
		return nil
	}

	UseVersion = func(version string) error {
		if mockUseVersionFn != nil {
			return mockUseVersionFn(version)
		}
		return nil
	}

	GetLatestVersion = func() (string, error) {
		if mockGetLatestVersionFn != nil {
			return mockGetLatestVersionFn()
		}
		return "3.78.1", nil
	}

	return func() {
		InstallVersion = originalInstallVersion
		UseVersion = originalUseVersion
		GetLatestVersion = originalGetLatestVersion
	}
}