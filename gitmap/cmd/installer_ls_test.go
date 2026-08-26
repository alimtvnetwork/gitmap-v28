// Package cmd — installer_ls_test.go tests the installer list CLI subcommand.
package cmd

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func setupInstallerLsTestDB(t *testing.T) *store.DB {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_installer_ls.db")
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

func TestInstallerLsCmd(t *testing.T) {
	if installerLsCmd == nil {
		t.Fatal("expected installerLsCmd to be initialized")
	}

	db := setupInstallerLsTestDB(t)
	script1 := &model.InstallerScript{
		Name:     "Win App",
		Slug:     "win-app",
		TargetOS: "win",
		Version:  "v1.0.0",
	}
	script2 := &model.InstallerScript{
		Name:     "Ubuntu App",
		Slug:     "ubuntu-app",
		TargetOS: "ubuntu",
		Version:  "v1.0.0",
	}
	db.CreateInstaller(script1)
	db.CreateInstaller(script2)

	ctx := context.Background()
	if errLs := executeInstallerLs(ctx, db, ""); errLs != nil {
		t.Fatalf("executeInstallerLs failed: %v", errLs)
	}
	if errLsWin := executeInstallerLs(ctx, db, "win"); errLsWin != nil {
		t.Fatalf("executeInstallerLs win failed: %v", errLsWin)
	}
}
