// Package cmd — installer_revert_test.go tests undo, redo, and revert subcommands.
package cmd

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func setupInstallerRevertTestDB(t *testing.T) *store.DB {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_installer_revert.db")
	db, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	if errMigrate := db.MigrateInstallers(); errMigrate != nil {
		t.Fatalf("failed to migrate installers: %v", errMigrate)
	}
	return db
}

func TestInstallerRevertCmd(t *testing.T) {
	if installerUndoCmd == nil || installerRedoCmd == nil || installerRevertCmd == nil {
		t.Fatal("expected revert commands to be initialized")
	}

	db := setupInstallerRevertTestDB(t)
	script := &model.InstallerScript{
		Name:     "Revert App",
		Slug:     "revert-app",
		TargetOS: "win",
		Version:  "v2.0.0",
	}
	if errCreate := db.CreateInstaller(script); errCreate != nil {
		t.Fatalf("failed to seed script: %v", errCreate)
	}

	ctx := context.Background()
	if errUndo := executeRevertAction(ctx, db, "undo", "revert-app", ""); errUndo != nil {
		t.Fatalf("executeRevertAction undo failed: %v", errUndo)
	}
	if errRedo := executeRevertAction(ctx, db, "redo", "revert-app", ""); errRedo != nil {
		t.Fatalf("executeRevertAction redo failed: %v", errRedo)
	}
	if errRevert := executeRevertAction(ctx, db, "revert", "revert-app", "v1.0.0"); errRevert != nil {
		t.Fatalf("executeRevertAction revert failed: %v", errRevert)
	}
}
