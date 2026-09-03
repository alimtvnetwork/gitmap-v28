// Package installer — export_all_test.go tests bulk installer export logic.
package installer

import (
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestExportAllToZip(t *testing.T) {
	db := setupInstallerTestDB(t)
	mgr, _ := NewManager(db)

	script1 := &model.InstallerScript{
		Name:     "App 1",
		Slug:     "app-1",
		TargetOS: "all",
		Version:  "v1.0.0",
	}

	script2 := &model.InstallerScript{
		Name:     "App 2",
		Slug:     "app-2",
		TargetOS: "win",
		Version:  "v1.0.0",
	}

	db.CreateInstaller(script1)
	db.CreateInstaller(script2)

	tempDir := t.TempDir()
	outZip := filepath.Join(tempDir, "all.zip")

	if errExport := mgr.ExportAllToZip(outZip); errExport != nil {
		t.Fatalf("ExportAllToZip failed: %v", errExport)
	}
}
