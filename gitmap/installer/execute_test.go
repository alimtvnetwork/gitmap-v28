// Package installer — execute_test.go tests installer execution dispatch.
package installer

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestExecuteDispatch(t *testing.T) {
	db := setupInstallerTestDB(t)
	mgr, _ := NewManager(db)

	script := &model.InstallerScript{
		Name:     "Exec App",
		Slug:     "exec-app",
		TargetOS: "win",
		Version:  "v1.0.0",
	}
	db.CreateInstaller(script)

	if errExec := mgr.Execute("exec-app", "win"); errExec != nil {
		t.Fatalf("Execute failed: %v", errExec)
	}

	if errMismatch := mgr.Execute("exec-app", "ubuntu"); errMismatch == nil {
		t.Fatal("expected OS mismatch error")
	}
}
