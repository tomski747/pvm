package utils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"
)

// buildTarGz creates an in-memory tar.gz archive with the given files.
// The archive has a single top-level directory "pulumi/" that gets stripped.
func buildTarGz(t *testing.T, files map[string]string) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	gzw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gzw)

	for name, content := range files {
		hdr := &tar.Header{
			Name: "pulumi/" + name,
			Mode: 0755,
			Size: int64(len(content)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatalf("tar write header: %v", err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			t.Fatalf("tar write content: %v", err)
		}
	}

	if err := tw.Close(); err != nil {
		t.Fatalf("close tar: %v", err)
	}
	if err := gzw.Close(); err != nil {
		t.Fatalf("close gzip: %v", err)
	}
	return &buf
}

func TestExtractTarGz(t *testing.T) {
	destDir := t.TempDir()

	archive := buildTarGz(t, map[string]string{
		"pulumi":          "#!/bin/sh\necho pulumi",
		"pulumi-language": "#!/bin/sh\necho lang",
	})

	if err := extractTarGz(archive, destDir); err != nil {
		t.Fatalf("extractTarGz: %v", err)
	}

	for _, name := range []string{"pulumi", "pulumi-language"} {
		p := filepath.Join(destDir, name)
		if _, err := os.Stat(p); err != nil {
			t.Errorf("expected file %s to exist: %v", name, err)
		}
	}
}

func TestExtractTarGzPathTraversal(t *testing.T) {
	destDir := t.TempDir()

	// Craft an archive with a path-traversal entry
	var buf bytes.Buffer
	gzw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gzw)

	hdr := &tar.Header{
		Name: "pulumi/../../evil",
		Mode: 0644,
		Size: 4,
	}
	_ = tw.WriteHeader(hdr)
	_, _ = tw.Write([]byte("evil"))
	_ = tw.Close()
	_ = gzw.Close()

	err := extractTarGz(&buf, destDir)
	if err == nil {
		t.Fatal("expected error for path traversal, got nil")
	}
}

func TestSafeJoin(t *testing.T) {
	tests := []struct {
		dest    string
		rel     string
		wantErr bool
	}{
		{"/dest", "file.txt", false},
		{"/dest", "sub/file.txt", false},
		{"/dest", "../escape", true},
		{"/dest", "sub/../../escape", true},
	}

	for _, tc := range tests {
		_, err := safeJoin(tc.dest, tc.rel)
		if (err != nil) != tc.wantErr {
			t.Errorf("safeJoin(%q, %q) error = %v, wantErr %v", tc.dest, tc.rel, err, tc.wantErr)
		}
	}
}
