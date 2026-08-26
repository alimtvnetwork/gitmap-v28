// Package installer — update_test.go tests installer update business logic.
package installer

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestManagerUpdate(t *testing.T) {
	db := setupInstallerTestDB(t)
	mgr, _ := NewManager(db)

	script := &model.InstallerScript{
		Name:     "Update Test",
		Slug:     "update-test",
		TargetOS: "win",
		Version:  "v1.0.0",
	}
	db.CreateInstaller(script)

	updated, errUpdate := mgr.Update("update-test", "ubuntu")
	if errUpdate != nil {
		t.Fatalf("Update failed: %v", errUpdate)
	}
	if updated.Version != "v1.0.1" || updated.TargetOS != "ubuntu" {
		t.Errorf("unexpected updated version: %+v", updated)
	}

	// Empty slug
	_, errEmpty := mgr.Update("", "")
	if errEmpty == nil {
		t.Fatal("expected error on empty slug")
	}
}
