package utils

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func downloadAndExtract(url string, destDir string, isZip bool) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	if isZip {
		return extractZip(resp.Body, destDir)
	}
	return extractTarGz(resp.Body, destDir)
}

// safeJoin joins destDir and relPath and verifies the result stays inside destDir.
// It returns an error if the path escapes the destination (zip-slip prevention).
func safeJoin(destDir, relPath string) (string, error) {
	destDir = filepath.Clean(destDir)
	path := filepath.Join(destDir, relPath)
	if !strings.HasPrefix(filepath.Clean(path)+string(filepath.Separator), destDir+string(filepath.Separator)) {
		return "", fmt.Errorf("invalid path in archive: %q escapes destination", relPath)
	}
	return path, nil
}

func extractTarGz(r io.Reader, destDir string) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar: %v", err)
		}

		// Skip the top-level directory entry
		parts := strings.Split(header.Name, string(filepath.Separator))
		if len(parts) <= 1 {
			continue
		}
		relPath := filepath.Join(parts[1:]...)

		path, err := safeJoin(destDir, relPath)
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
			if err := writeFile(path, tr); err != nil {
				return err
			}
		}
	}
	return nil
}

// writeFile creates path and copies content from r into it, then chmods it executable.
func writeFile(path string, r io.Reader) error {
	outFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	_, copyErr := io.Copy(outFile, r)
	if closeErr := outFile.Close(); closeErr != nil && copyErr == nil {
		return fmt.Errorf("failed to close file: %v", closeErr)
	}
	if copyErr != nil {
		return fmt.Errorf("failed to write file: %v", copyErr)
	}
	if err := os.Chmod(path, 0755); err != nil {
		return fmt.Errorf("failed to set permissions: %v", err)
	}
	return nil
}

func extractZip(r io.Reader, destDir string) error {
	// zip.Reader requires io.ReaderAt, so buffer to a temp file first.
	tmpFile, err := os.CreateTemp("", "pulumi-*.zip")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, r); err != nil {
		return fmt.Errorf("failed to write temp file: %v", err)
	}

	zipReader, err := zip.OpenReader(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("failed to open zip: %v", err)
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		// Skip the top-level directory entry
		parts := strings.Split(file.Name, string(filepath.Separator))
		if len(parts) <= 1 {
			continue
		}
		relPath := filepath.Join(parts[1:]...)

		path, err := safeJoin(destDir, relPath)
		if err != nil {
			return err
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}

		rc, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open zip entry: %v", err)
		}
		if err := writeFile(path, rc); err != nil {
			_ = rc.Close()
			return err
		}
		if err := rc.Close(); err != nil {
			return fmt.Errorf("failed to close zip entry: %v", err)
		}
	}
	return nil
}
