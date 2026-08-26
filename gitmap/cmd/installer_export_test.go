// Package cmd — installer_export_test.go tests installer export CLI subcommands.
package cmd

import (
	"archive/zip"
	"context"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func setupInstallerExportTestDB(t *testing.T) *store.DB {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_installer_export.db")
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

func TestInstallerExportCmd(t *testing.T) {
	if installerExportCmd == nil || installerExportAllCmd == nil {
		t.Fatal("expected export commands to be initialized")
	}

	flags, err := parseExportFlags([]string{"my-app", "-o", "out.zip"}, false)
	if err != nil {
		t.Fatalf("unexpected error parsing flags: %v", err)
	}
	if flags.Slug != "my-app" || flags.OutputPath != "out.zip" {
		t.Errorf("unexpected parsed values: %+v", flags)
	}

	// Test missing slug for single export
	_, errMissing := parseExportFlags([]string{}, false)
	if errMissing == nil {
		t.Fatal("expected error on missing slug")
	}

	// Test executeExport single and all
	db := setupInstallerExportTestDB(t)
	script := &model.InstallerScript{
		Name:     "Export App",
		Slug:     "export-app",
		TargetOS: "all",
		Version:  "v1.0.0",
	}
	if errCreate := db.CreateInstaller(script); errCreate != nil {
		t.Fatalf("failed to seed installer: %v", errCreate)
	}

	tempDir := t.TempDir()
	zipPath := filepath.Join(tempDir, "test_export.zip")
	flagsExport := &ExportInstallerFlags{
		Slug:       "export-app",
		OutputPath: zipPath,
		ExportAll:  false,
	}

	ctx := context.Background()
	if errExport := executeExport(ctx, db, flagsExport); errExport != nil {
		t.Fatalf("executeExport failed: %v", errExport)
	}

	// Verify zip content
	zr, errZip := zip.OpenReader(zipPath)
	if errZip != nil {
		t.Fatalf("failed to read created zip: %v", errZip)
	}
	defer zr.Close()

	if len(zr.File) != 1 || zr.File[0].Name != "export-app.json" {
		t.Errorf("unexpected zip contents: %+v", zr.File)
	}

	// Test export-all
	zipAllPath := filepath.Join(tempDir, "test_export_all.zip")
	flagsExportAll := &ExportInstallerFlags{
		OutputPath: zipAllPath,
		ExportAll:  true,
	}
	if errExportAll := executeExport(ctx, db, flagsExportAll); errExportAll != nil {
		t.Fatalf("executeExport all failed: %v", errExportAll)
	}
}
