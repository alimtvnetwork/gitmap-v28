// Package installer — revert_test.go tests installer undo version logic.
package installer

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestManagerUndo(t *testing.T) {
	db := setupInstallerTestDB(t)
	mgr, _ := NewManager(db)

	script := &model.InstallerScript{
		Name:     "Undo Test",
		Slug:     "undo-test",
		TargetOS: "win",
		Version:  "v1.0.0",
	}

	db.CreateInstaller(script)

	if errUndo := mgr.Undo("undo-test"); errUndo != nil {
		t.Fatalf("Undo failed: %v", errUndo)
	}

	if errEmpty := mgr.Undo(""); errEmpty == nil {
		t.Fatal("expected error on empty slug")
	}
}
