package installer

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestManagerDelete(t *testing.T) {
	db := setupInstallerTestDB(t)
	mgr, _ := NewManager(db)

	script := &model.InstallerScript{
		Name:     "Delete App",
		Slug:     "delete-app",
		TargetOS: "all",
		Version:  "v1.0.0",
	}
	db.CreateInstaller(script)

	if err := mgr.Delete("delete-app"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if errEmpty := mgr.Delete(""); errEmpty == nil {
		t.Fatal("expected error on empty slug")
	}

	if errVer := mgr.DeleteVersion("delete-app", "v1.0.0"); errVer != nil {
		t.Fatalf("DeleteVersion failed: %v", errVer)
	}
}
