package commands

import (
	"bytes"
	"testing"
)

func TestRootCommand(t *testing.T) {
	// Save original output and restore after test
	oldOut := rootCmd.OutOrStdout()
	oldErr := rootCmd.ErrOrStderr()
	defer func() {
		rootCmd.SetOut(oldOut)
		rootCmd.SetErr(oldErr)
	}()

	// Create buffer for test output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Test help output
	rootCmd.SetArgs([]string{"--help"})
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	output := buf.String()
	expectedPhrases := []string{
		"Usage:",
		"Available Commands:",
		"Flags:",
	}

	for _, phrase := range expectedPhrases {
		if !bytes.Contains(buf.Bytes(), []byte(phrase)) {
			t.Errorf("Expected output to contain '%s', got:\n%s", phrase, output)
		}
	}
}

func TestRootCommandNoArgs(t *testing.T) {
	// Save original output and restore after test
	oldOut := rootCmd.OutOrStdout()
	oldErr := rootCmd.ErrOrStderr()
	defer func() {
		rootCmd.SetOut(oldOut)
		rootCmd.SetErr(oldErr)
	}()

	// Create buffer for test output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Test with no arguments
	rootCmd.SetArgs([]string{})
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	output := buf.String()
	expectedPhrases := []string{
		"Usage:",
		"Available Commands:",
		"install",
		"use",
		"list",
		"current",
	}

	for _, phrase := range expectedPhrases {
		if !bytes.Contains(buf.Bytes(), []byte(phrase)) {
			t.Errorf("Expected output to contain '%s', got:\n%s", phrase, output)
		}
	}
}