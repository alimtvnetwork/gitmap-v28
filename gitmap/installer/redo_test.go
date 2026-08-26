// Package installer — redo_test.go tests installer redo version logic.
package installer

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestManagerRedo(t *testing.T) {
	db := setupInstallerTestDB(t)
	mgr, _ := NewManager(db)

	script := &model.InstallerScript{
		Name:     "Redo Test",
		Slug:     "redo-test",
		TargetOS: "win",
		Version:  "v1.0.0",
	}
	db.CreateInstaller(script)

	if errRedo := mgr.Redo("redo-test"); errRedo != nil {
		t.Fatalf("Redo failed: %v", errRedo)
	}

	if errEmpty := mgr.Redo(""); errEmpty == nil {
		t.Fatal("expected error on empty slug")
	}
}
