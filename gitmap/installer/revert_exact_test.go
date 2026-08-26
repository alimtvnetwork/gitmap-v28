// Package installer — revert_exact_test.go tests exact version reversion logic.
package installer

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestManagerRevertExact(t *testing.T) {
	db := setupInstallerTestDB(t)
	mgr, _ := NewManager(db)

	script := &model.InstallerScript{
		Name:     "Revert Exact Test",
		Slug:     "revert-exact-test",
		TargetOS: "win",
		Version:  "v2.0.0",
	}
	db.CreateInstaller(script)

	if errRevert := mgr.RevertTo("revert-exact-test", "v1.0.0"); errRevert != nil {
		t.Fatalf("RevertTo failed: %v", errRevert)
	}

	if errEmpty := mgr.RevertTo("", ""); errEmpty == nil {
		t.Fatal("expected error on empty slug/version")
	}
}
