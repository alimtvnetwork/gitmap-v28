package installer

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestExportGitRemote(t *testing.T) {
	db := setupInstallerTestDB(t)
	mgr, _ := NewManager(db)

	script := &model.InstallerScript{
		Name:     "Remote App",
		Slug:     "remote-app",
		TargetOS: "all",
		Version:  "v1.0.0",
	}

	db.CreateInstaller(script)

	// Simulated remote repo export
	if err := mgr.ExportToRemoteGitRepo("remote-app", "https://example.com/repo.git", "main", "app.json", "msg", false); err != nil {
		t.Fatalf("ExportToRemoteGitRepo failed: %v", err)
	}
}
