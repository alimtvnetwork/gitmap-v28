// Package cmd — installer_update_test.go tests the installer update CLI subcommand.
package cmd

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func setupInstallerUpdateTestDB(t *testing.T) *store.DB {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_installer_update.db")
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

func TestInstallerUpdateCmd(t *testing.T) {
	if installerUpdateCmd == nil {
		t.Fatal("expected installerUpdateCmd to be initialized")
	}

	flags, err := parseUpdateFlags([]string{"my-app", "-v", "v2.0.0", "-os", "ubuntu"})
	if err != nil {
		t.Fatalf("unexpected error parsing flags: %v", err)
	}
	if flags.Slug != "my-app" || flags.Version != "v2.0.0" || flags.TargetOS != "ubuntu" {
		t.Errorf("unexpected parsed values: %+v", flags)
	}

	// Test missing slug
	_, errMissing := parseUpdateFlags([]string{})
	if errMissing == nil {
		t.Fatal("expected error on missing slug")
	}

	// Test executeUpdate
	db := setupInstallerUpdateTestDB(t)
	script := &model.InstallerScript{
		Name:     "My App",
		Slug:     "my-app",
		TargetOS: "win",
		Version:  "v1.0.0",
	}
	if errCreate := db.CreateInstaller(script); errCreate != nil {
		t.Fatalf("failed to seed installer: %v", errCreate)
	}

	ctx := context.Background()
	if errUpdate := executeInstallerUpdate(ctx, db, flags); errUpdate != nil {
		t.Fatalf("executeUpdate failed: %v", errUpdate)
	}

	// Verify not found
	flagsNotFound := &UpdateInstallerFlags{Slug: "non-existent"}
	errNotFound := executeInstallerUpdate(ctx, db, flagsNotFound)
	if errNotFound == nil {
		t.Fatal("expected not found error")
	}
	appErr, ok := errNotFound.(*apperror.AppError)
	if !ok || appErr.Code != "E_INSTALLER_NOT_FOUND" {
		t.Errorf("expected E_INSTALLER_NOT_FOUND code, got: %v", errNotFound)
	}
}
