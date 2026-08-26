package cmd

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestChromeProfileZipExport(t *testing.T) {
	src := t.TempDir()

	// Create mock sqlite database files
	os.MkdirAll(filepath.Join(src, "Extensions"), 0755)
	os.WriteFile(filepath.Join(src, "History"), []byte("history-blob"), 0644)
	os.WriteFile(filepath.Join(src, "Web Data"), []byte("web-data-blob"), 0644)
	os.WriteFile(filepath.Join(src, "Login Data"), []byte("secret"), 0644) // Should NOT be in zip

	outDir := t.TempDir()
	outZip := filepath.Join(outDir, "test.zip")

	// Export ZIP
	n, err := writeChromeExportZIP(src, "testprofile", outZip)
	if err != nil {
		t.Fatalf("writeChromeExportZIP failed: %v", err)
	}
	if n == 0 {
		t.Fatalf("expected non-zero bytes")
	}

	// Verify ZIP contents
	r, err := zip.OpenReader(outZip)
	if err != nil {
		t.Fatalf("open zip failed: %v", err)
	}
	defer r.Close()

	foundJSON := false
	foundHistory := false
	foundWebData := false
	foundLoginData := false

	for _, f := range r.File {
		if f.Name == "testprofile.json" {
			foundJSON = true
		}
		if f.Name == "History" {
			foundHistory = true
		}
		if f.Name == "Web Data" {
			foundWebData = true
		}
		if f.Name == "Login Data" {
			foundLoginData = true
		}
	}

	if !foundJSON {
		t.Errorf("missing JSON snapshot in zip")
	}
	if !foundHistory {
		t.Errorf("missing History in zip")
	}
	if !foundWebData {
		t.Errorf("missing Web Data in zip")
	}
	if foundLoginData {
		t.Errorf("Login Data should be omitted from export")
	}
}
