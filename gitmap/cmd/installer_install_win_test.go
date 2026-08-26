// Package cmd — installer_install_win_test.go tests the installer install-win CLI subcommand.
package cmd

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func setupInstallerInstallWinTestDB(t *testing.T) *store.DB {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_installer_install_win.db")
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

func TestInstallerInstallWinCmd(t *testing.T) {
	if installerInstallWinCmd == nil {
		t.Fatal("expected installerInstallWinCmd to be initialized")
	}

	flags, err := parseInstallWinFlags([]string{"my-app", "-d"})
	if err != nil {
		t.Fatalf("unexpected error parsing flags: %v", err)
	}
	if flags.Slug != "my-app" || !flags.DryRun {
		t.Errorf("unexpected parsed values: %+v", flags)
	}

	// Test missing slug
	_, errMissing := parseInstallWinFlags([]string{})
	if errMissing == nil {
		t.Fatal("expected error on missing slug")
	}

	// Test executeInstallWin
	db := setupInstallerInstallWinTestDB(t)
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
	if errInstall := executeInstallWin(ctx, db, flags); errInstall != nil {
		t.Fatalf("executeInstallWin failed: %v", errInstall)
	}

	// Verify OS mismatch
	scriptUbuntu := &model.InstallerScript{
		Name:     "Ubuntu App",
		Slug:     "ubuntu-app",
		TargetOS: "ubuntu",
		Version:  "v1.0.0",
	}
	if errCreate := db.CreateInstaller(scriptUbuntu); errCreate != nil {
		t.Fatalf("failed to seed installer: %v", errCreate)
	}
	flagsMismatch := &InstallWinFlags{Slug: "ubuntu-app"}
	errMismatch := executeInstallWin(ctx, db, flagsMismatch)
	if errMismatch == nil {
		t.Fatal("expected OS mismatch error")
	}
	appErr, ok := errMismatch.(*apperror.AppError)
	if !ok || appErr.Code != "E_INSTALLER_OS_MISMATCH" {
		t.Errorf("expected E_INSTALLER_OS_MISMATCH code, got: %v", errMismatch)
	}
}
