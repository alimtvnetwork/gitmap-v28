// Package cmd — installer_reset_test.go tests the installer reset CLI subcommand.
package cmd

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func setupInstallerResetTestDB(t *testing.T) *store.DB {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_installer_reset.db")
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

func TestInstallerResetCmd(t *testing.T) {
	if installerResetCmd == nil {
		t.Fatal("expected installerResetCmd to be initialized")
	}

	flags, err := parseInstallerResetFlags([]string{"my-app"})
	if err != nil {
		t.Fatalf("unexpected error parsing flags: %v", err)
	}
	if flags.Slug != "my-app" || flags.ResetAll {
		t.Errorf("unexpected parsed values: %+v", flags)
	}

	// Test missing slug when not all
	_, errMissing := parseInstallerResetFlags([]string{})
	if errMissing == nil {
		t.Fatal("expected error on missing slug without --all")
	}

	// Test executeInstallerReset
	db := setupInstallerResetTestDB(t)
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
	if errReset := executeInstallerReset(ctx, db, flags); errReset != nil {
		t.Fatalf("executeInstallerReset failed: %v", errReset)
	}

	// Verify not found after reset
	_, errGet := db.GetInstallerBySlug("my-app")
	if errGet == nil {
		t.Fatal("expected script to be removed after reset")
	}
}
