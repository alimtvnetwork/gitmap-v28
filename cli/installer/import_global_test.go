// Package installer — import_global_test.go tests global state archive restoration.
package installer

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestImportGlobalState(t *testing.T) {
	db := setupInstallerTestDB(t)
	mgr, _ := NewManager(db)

	tempDir := t.TempDir()
	zipPath := filepath.Join(tempDir, "test_global.zip")

	f, _ := os.Create(zipPath)
	zw := zip.NewWriter(f)
	w, _ := zw.Create("installers/global-app.json")
	w.Write([]byte(`{"name": "Global Restored App", "slug": "global-restored", "target_os": "all", "version": "v1.0.0"}`))
	zw.Close()
	f.Close()

	if errImport := mgr.ImportGlobalState(zipPath); errImport != nil {
		t.Fatalf("ImportGlobalState failed: %v", errImport)
	}

	imported, errGet := db.GetInstallerBySlug("global-restored")
	if errGet != nil || imported == nil {
		t.Fatalf("failed to retrieve restored script: %v", errGet)
	}
}
