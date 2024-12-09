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

		// Skip the top-level directory
		parts := strings.Split(header.Name, string(filepath.Separator))
		if len(parts) <= 1 {
			continue
		}
		relPath := filepath.Join(parts[1:]...)
		path := filepath.Join(destDir, relPath)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
		case tar.TypeReg:
			outFile, err := os.Create(path)
			if err != nil {
				return fmt.Errorf("failed to create file: %v", err)
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to write file: %v", err)
			}
			outFile.Close()
			if err := os.Chmod(path, 0755); err != nil {
				return fmt.Errorf("failed to set permissions: %v", err)
			}
		}
	}
	return nil
}

func extractZip(r io.Reader, destDir string) error {
	// Create a temporary file to store the zip content
	tmpFile, err := os.CreateTemp("", "pulumi-*.zip")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Copy the content to the temp file
	if _, err := io.Copy(tmpFile, r); err != nil {
		return fmt.Errorf("failed to write temp file: %v", err)
	}

	// Open the zip file
	zipReader, err := zip.OpenReader(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("failed to open zip: %v", err)
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		// Skip the top-level directory
		parts := strings.Split(file.Name, string(filepath.Separator))
		if len(parts) <= 1 {
			continue
		}
		relPath := filepath.Join(parts[1:]...)
		path := filepath.Join(destDir, relPath)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}

		outFile, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("failed to open zip file: %v", err)
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()
		if err != nil {
			return fmt.Errorf("failed to write file: %v", err)
		}

		if err := os.Chmod(path, 0755); err != nil {
			return fmt.Errorf("failed to set permissions: %v", err)
		}
	}
	return nil
}
