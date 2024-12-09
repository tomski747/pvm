package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetHomeDir(t *testing.T) {
	// Test with test config
	testDir := "/test/home"
	testConfig = &TestConfig{
		PVMPath: testDir,
	}
	if got := GetHomeDir(); got != testDir {
		t.Errorf("GetHomeDir() = %v, want %v", got, testDir)
	}

	// Test without test config
	testConfig = nil
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME")
	}
	if got := GetHomeDir(); got != home {
		t.Errorf("GetHomeDir() = %v, want %v", got, home)
	}
}

func TestGetPVMPath(t *testing.T) {
	// Test with test config
	testDir := "/test/pvm"
	testConfig = &TestConfig{
		PVMPath: testDir,
	}
	if got := GetPVMPath(); got != testDir {
		t.Errorf("GetPVMPath() = %v, want %v", got, testDir)
	}

	// Test without test config
	testConfig = nil
	expected := filepath.Join(GetHomeDir(), PVMDir)
	if got := GetPVMPath(); got != expected {
		t.Errorf("GetPVMPath() = %v, want %v", got, expected)
	}
}

func TestGetVersionsPath(t *testing.T) {
	expected := filepath.Join(GetPVMPath(), VersionsDir)
	if got := GetVersionsPath(); got != expected {
		t.Errorf("GetVersionsPath() = %v, want %v", got, expected)
	}
}

func TestGetBinPath(t *testing.T) {
	expected := filepath.Join(GetPVMPath(), BinDir)
	if got := GetBinPath(); got != expected {
		t.Errorf("GetBinPath() = %v, want %v", got, expected)
	}
}

func TestGetPlatformInfo(t *testing.T) {
	goos, arch := GetPlatformInfo()
	if goos == "" {
		t.Error("GetPlatformInfo() returned empty GOOS")
	}
	if arch == "" {
		t.Error("GetPlatformInfo() returned empty architecture")
	}
} 