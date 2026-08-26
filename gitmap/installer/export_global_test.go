// Package installer — export_global_test.go tests global system export logic.
package installer

import (
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestExportGlobalState(t *testing.T) {
	db := setupInstallerTestDB(t)
	mgr, _ := NewManager(db)

	script := &model.InstallerScript{
		Name:     "Global App",
		Slug:     "global-app",
		TargetOS: "all",
		Version:  "v1.0.0",
	}
	db.CreateInstaller(script)

	tempDir := t.TempDir()
	outZip := filepath.Join(tempDir, "global.zip")

	if errExport := mgr.ExportGlobalState(outZip); errExport != nil {
		t.Fatalf("ExportGlobalState failed: %v", errExport)
	}
}
