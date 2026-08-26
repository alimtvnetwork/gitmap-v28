// Package cmd — installer_update_win_test.go tests the installer update-win CLI subcommand.
package cmd

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func setupInstallerUpdateWinTestDB(t *testing.T) *store.DB {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_installer_update_win.db")
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

func TestInstallerUpdateWinCmd(t *testing.T) {
	if installerUpdateWinCmd == nil {
		t.Fatal("expected installerUpdateWinCmd to be initialized")
	}

	flags, err := parseUpdateWinFlags([]string{"my-app", "-v", "v2.0.0"})
	if err != nil {
		t.Fatalf("unexpected error parsing flags: %v", err)
	}
	if flags.Slug != "my-app" || flags.Version != "v2.0.0" {
		t.Errorf("unexpected parsed values: %+v", flags)
	}

	// Test missing slug
	_, errMissing := parseUpdateWinFlags([]string{})
	if errMissing == nil {
		t.Fatal("expected error on missing slug")
	}

	// Test executeUpdateWin
	db := setupInstallerUpdateWinTestDB(t)
	script := &model.InstallerScript{
		Name:     "My App",
		Slug:     "my-app",
		TargetOS: "ubuntu",
		Version:  "v1.0.0",
	}
	if errCreate := db.CreateInstaller(script); errCreate != nil {
		t.Fatalf("failed to seed installer: %v", errCreate)
	}

	ctx := context.Background()
	if errUpdate := executeUpdateWin(ctx, db, flags); errUpdate != nil {
		t.Fatalf("executeUpdateWin failed: %v", errUpdate)
	}

	// Verify not found
	flagsNotFound := &UpdateWinInstallerFlags{Slug: "non-existent"}
	errNotFound := executeUpdateWin(ctx, db, flagsNotFound)
	if errNotFound == nil {
		t.Fatal("expected not found error")
	}
	appErr, ok := errNotFound.(*apperror.AppError)
	if !ok || appErr.Code != "E_INSTALLER_NOT_FOUND" {
		t.Errorf("expected E_INSTALLER_NOT_FOUND code, got: %v", errNotFound)
	}
}
