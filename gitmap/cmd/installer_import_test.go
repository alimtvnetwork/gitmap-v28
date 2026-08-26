// Package cmd — installer_import_test.go tests the installer import CLI subcommand.
package cmd

import (
	"archive/zip"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func setupInstallerImportTestDB(t *testing.T) *store.DB {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_installer_import.db")
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

func TestInstallerImportCmd(t *testing.T) {
	if installerImportCmd == nil {
		t.Fatal("expected installerImportCmd to be initialized")
	}

	flags, err := parseInstallerImportFlags([]string{"custom.zip"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flags.InputPath != "custom.zip" {
		t.Errorf("unexpected input path: %s", flags.InputPath)
	}

	// Create test zip with one installer
	tempDir := t.TempDir()
	zipPath := filepath.Join(tempDir, "test_in.zip")
	fZip, errCreate := os.Create(zipPath)
	if errCreate != nil {
		t.Fatalf("failed to create zip: %v", errCreate)
	}
	zw := zip.NewWriter(fZip)
	w, errEntry := zw.Create("imported-app.json")
	if errEntry != nil {
		t.Fatalf("failed to create zip entry: %v", errEntry)
	}
	script := model.InstallerScript{
		Name:     "Imported App",
		Slug:     "imported-app",
		TargetOS: "all",
		Version:  "v1.0.0",
	}
	data, _ := json.Marshal(script)
	w.Write(data)
	zw.Close()
	fZip.Close()

	db := setupInstallerImportTestDB(t)
	flagsImport := &ImportInstallerFlags{InputPath: zipPath}
	ctx := context.Background()

	if errExec := executeInstallerImport(ctx, db, flagsImport); errExec != nil {
		t.Fatalf("executeImport failed: %v", errExec)
	}

	imported, errGet := db.GetInstallerBySlug("imported-app")
	if errGet != nil || imported == nil {
		t.Fatalf("failed to find imported script: %v", errGet)
	}
	if imported.Name != "Imported App" {
		t.Errorf("unexpected imported name: %s", imported.Name)
	}
}
