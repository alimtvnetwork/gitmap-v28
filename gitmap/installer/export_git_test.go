package installer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestExportGitFolder(t *testing.T) {
	db := setupInstallerTestDB(t)
	mgr, _ := NewManager(db)

	script := &model.InstallerScript{
		Name:     "Git Export App",
		Slug:     "git-export-app",
		TargetOS: "all",
		Version:  "v1.0.0",
	}
	db.CreateInstaller(script)

	tempDir := t.TempDir()
	if err := mgr.ExportToGitFolder("git-export-app", tempDir, "app.json", "test commit"); err != nil {
		t.Fatalf("ExportToGitFolder failed: %v", err)
	}

	exportedFile := filepath.Join(tempDir, "app.json")
	if _, errStat := os.Stat(exportedFile); errStat != nil {
		t.Fatalf("expected exported file to exist: %v", errStat)
	}
}
