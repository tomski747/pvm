package utils

import (
	"testing"
)

// Mock function types
type (
	mockInstallVersionFunc      func(version string) error
	mockUseVersionFunc          func(version string) error
	mockGetLatestVersionFunc    func() (string, error)
	mockResolveVersionFunc      func(version string) (string, error)
	mockGetAvailableVersionFunc func(refresh bool) ([]string, error)
)

// Mock function variables – set these in tests before calling Execute/RunE.
var (
	mockInstallVersionFn      mockInstallVersionFunc
	mockUseVersionFn          mockUseVersionFunc
	mockGetLatestVersionFn    mockGetLatestVersionFunc
	mockResolveVersionFn      mockResolveVersionFunc
	mockGetAvailableVersionFn mockGetAvailableVersionFunc
)

// MockVersionOperations replaces network-dependent function variables with
// no-op stubs (or delegates to the mock*Fn variables when set) and returns
// a cleanup function that restores the originals.
func MockVersionOperations(t testing.TB) func() {
	t.Helper()

	origInstall := InstallVersion
	origUse := UseVersion
	origLatest := GetLatestVersion
	origResolve := ResolveVersion
	origAvailable := GetAvailableVersions

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

	ResolveVersion = func(version string) (string, error) {
		if mockResolveVersionFn != nil {
			return mockResolveVersionFn(version)
		}
		// Default: treat the version string as already resolved
		return version, nil
	}

	GetAvailableVersions = func(refresh bool) ([]string, error) {
		if mockGetAvailableVersionFn != nil {
			return mockGetAvailableVersionFn(refresh)
		}
		return []string{"3.78.1", "3.78.0", "3.77.0"}, nil
	}

	return func() {
		InstallVersion = origInstall
		UseVersion = origUse
		GetLatestVersion = origLatest
		ResolveVersion = origResolve
		GetAvailableVersions = origAvailable
		// Clear per-test overrides
		mockInstallVersionFn = nil
		mockUseVersionFn = nil
		mockGetLatestVersionFn = nil
		mockResolveVersionFn = nil
		mockGetAvailableVersionFn = nil
	}
}
