// Package installer — export_test.go tests single installer export logic.
package installer

import (
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestExportToZip(t *testing.T) {
	db := setupInstallerTestDB(t)
	mgr, _ := NewManager(db)

	script := &model.InstallerScript{
		Name:     "Export Test",
		Slug:     "export-test",
		TargetOS: "all",
		Version:  "v1.0.0",
	}

	db.CreateInstaller(script)

	tempDir := t.TempDir()
	outZip := filepath.Join(tempDir, "out.zip")

	if errExport := mgr.ExportToZip("export-test", outZip); errExport != nil {
		t.Fatalf("ExportToZip failed: %v", errExport)
	}

	if errEmpty := mgr.ExportToZip("", outZip); errEmpty == nil {
		t.Fatal("expected error on empty slug")
	}
}
