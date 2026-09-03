// Package installer — conflict_test.go tests import conflict resolution.
package installer

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestConflictResolver(t *testing.T) {
	db := setupInstallerTestDB(t)
	mgr, _ := NewManager(db)

	script := &model.InstallerScript{
		Name:     "Conflict App",
		Slug:     "conflict-app",
		TargetOS: "all",
		Version:  "v1.0.0",
	}

	db.CreateInstaller(script)

	if errConflict := mgr.resolveImportConflict("conflict-app"); errConflict != nil {
		t.Fatalf("resolveImportConflict failed: %v", errConflict)
	}
}
