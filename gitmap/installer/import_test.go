// Package installer — import_test.go tests zip archive import business logic.
package installer

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestImportFromZip(t *testing.T) {
	db := setupInstallerTestDB(t)
	mgr, _ := NewManager(db)

	tempDir := t.TempDir()
	zipPath := filepath.Join(tempDir, "test_import.zip")

	f, _ := os.Create(zipPath)
	zw := zip.NewWriter(f)
	w, _ := zw.Create("zip-app.json")
	w.Write([]byte(`{"name": "Zip App", "slug": "zip-app", "target_os": "all", "version": "v1.0.0"}`))
	zw.Close()
	f.Close()

	if errImport := mgr.ImportFromZip(zipPath); errImport != nil {
		t.Fatalf("ImportFromZip failed: %v", errImport)
	}

	imported, errGet := db.GetInstallerBySlug("zip-app")
	if errGet != nil || imported == nil {
		t.Fatalf("failed to retrieve imported script: %v", errGet)
	}
}
