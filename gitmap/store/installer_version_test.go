// Package store — installer_version_test.go tests SaveVersion.
package store

import (
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	_ "modernc.org/sqlite"
)

func setupInstallerVersionTestDB(testingT *testing.T) *DB {
	testingT.Helper()
	tempDir := testingT.TempDir()
	dbPath := filepath.Join(tempDir, "test_installer_version.db")

	dbInstance, errOpen := OpenAt(dbPath)
	if errOpen != nil {
		testingT.Fatalf("failed to open test db: %v", errOpen)
	}
	testingT.Cleanup(func() { _ = dbInstance.Close() })

	if errMigrate := dbInstance.MigrateInstallers(); errMigrate != nil {
		testingT.Fatalf("failed to migrate installers: %v", errMigrate)
	}

	return dbInstance
}

func TestSaveVersionSuccess(testingT *testing.T) {
	dbInstance := setupInstallerVersionTestDB(testingT)

	script := &model.InstallerScript{
		Name:         "Golang Suite",
		Slug:         "go-suite",
		Description:  "Go tooling installer",
		TargetOS:     "win",
		Version:      "v1.0.0",
		Instructions: `{"steps":[{"action":"install","pkg":"go"}]}`,
	}
	if err := dbInstance.CreateInstaller(script); err != nil {
		testingT.Fatalf("failed to create parent installer script: %v", err)
	}

	versionRecord := &model.InstallerVersion{
		ScriptID:     script.ID,
		Slug:         "go-suite",
		Version:      "v1.0.0",
		TargetOS:     "win",
		Instructions: `{"steps":[{"action":"install","pkg":"go"}]}`,
	}

	errSave := dbInstance.SaveVersion(versionRecord)
	if errSave != nil {
		testingT.Fatalf("expected SaveVersion to succeed, got: %v", errSave)
	}

	if versionRecord.ID <= 0 {
		testingT.Fatalf("expected versionRecord.ID to be populated > 0, got %d", versionRecord.ID)
	}

	row := dbInstance.conn.QueryRow(
		"SELECT id, script_id, slug, version, target_os, instructions, created_at FROM installer_versions WHERE id = ?",
		versionRecord.ID,
	)
	var (
		fetchedID           int64
		fetchedScriptID     int64
		fetchedSlug         string
		fetchedVersion      string
		fetchedTargetOS     string
		fetchedInstructions string
		fetchedCreatedAt    string
	)

	errScan := row.Scan(
		&fetchedID,
		&fetchedScriptID,
		&fetchedSlug,
		&fetchedVersion,
		&fetchedTargetOS,
		&fetchedInstructions,
		&fetchedCreatedAt,
	)
	if errScan != nil {
		testingT.Fatalf("failed to scan inserted installer version: %v", errScan)
	}

	if fetchedID != versionRecord.ID {
		testingT.Errorf("expected ID %d, got %d", versionRecord.ID, fetchedID)
	}
	if fetchedScriptID != script.ID {
		testingT.Errorf("expected script_id %d, got %d", script.ID, fetchedScriptID)
	}
	if fetchedSlug != "go-suite" {
		testingT.Errorf("expected slug 'go-suite', got %q", fetchedSlug)
	}
	if fetchedVersion != "v1.0.0" {
		testingT.Errorf("expected version 'v1.0.0', got %q", fetchedVersion)
	}
	if fetchedTargetOS != "win" {
		testingT.Errorf("expected target_os 'win', got %q", fetchedTargetOS)
	}
	if fetchedInstructions != versionRecord.Instructions {
		testingT.Errorf("expected instructions %q, got %q", versionRecord.Instructions, fetchedInstructions)
	}
	if fetchedCreatedAt == "" {
		testingT.Errorf("expected created_at to be populated")
	}
}

func TestSaveVersionMultipleVersions(testingT *testing.T) {
	dbInstance := setupInstallerVersionTestDB(testingT)

	v1 := &model.InstallerVersion{
		ScriptID:     10,
		Slug:         "node",
		Version:      "v18.0.0",
		TargetOS:     "ubuntu",
		Instructions: "apt install nodejs",
	}
	v2 := &model.InstallerVersion{
		ScriptID:     10,
		Slug:         "node",
		Version:      "v20.0.0",
		TargetOS:     "ubuntu",
		Instructions: "apt install nodejs-20",
	}

	if err := dbInstance.SaveVersion(v1); err != nil {
		testingT.Fatalf("expected v1 to save successfully, got: %v", err)
	}
	if err := dbInstance.SaveVersion(v2); err != nil {
		testingT.Fatalf("expected v2 to save successfully, got: %v", err)
	}

	if v1.ID <= 0 || v2.ID <= 0 || v1.ID == v2.ID {
		testingT.Fatalf("expected distinct valid IDs for v1 and v2, got %d and %d", v1.ID, v2.ID)
	}
}

func TestSaveVersionNilInput(testingT *testing.T) {
	dbInstance := setupInstallerVersionTestDB(testingT)

	errSave := dbInstance.SaveVersion(nil)
	if errSave == nil {
		testingT.Fatalf("expected error on nil version, got nil")
	}

	var appErr *apperror.AppError
	if !errors.As(errSave, &appErr) {
		testingT.Fatalf("expected AppError, got %T", errSave)
	}
	if appErr.Code != "E_INSTALLER_INVALID_INPUT" {
		testingT.Errorf("expected error code E_INSTALLER_INVALID_INPUT, got %s", appErr.Code)
	}
}

func TestSaveVersionNilDB(testingT *testing.T) {
	var dbInstance *DB

	versionRecord := &model.InstallerVersion{
		Slug:    "python",
		Version: "3.11",
	}
	errSave := dbInstance.SaveVersion(versionRecord)
	if errSave == nil {
		testingT.Fatalf("expected error on nil db, got nil")
	}

	var appErr *apperror.AppError
	if !errors.As(errSave, &appErr) {
		testingT.Fatalf("expected AppError, got %T", errSave)
	}
	if appErr.Code != "E_INSTALLER_NIL_DB" {
		testingT.Errorf("expected error code E_INSTALLER_NIL_DB, got %s", appErr.Code)
	}
}

func TestSaveVersionClosedDB(testingT *testing.T) {
	dbConn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		testingT.Fatalf("sql.Open failed: %v", err)
	}
	dbInstance := &DB{conn: dbConn}
	_ = dbConn.Close()

	versionRecord := &model.InstallerVersion{
		Slug:    "python",
		Version: "3.11",
	}
	errSave := dbInstance.SaveVersion(versionRecord)
	if errSave == nil {
		testingT.Fatalf("expected error on closed db, got nil")
	}

	var appErr *apperror.AppError
	if !errors.As(errSave, &appErr) {
		testingT.Fatalf("expected AppError, got %T", errSave)
	}
	if appErr.Code != "E_INSTALLER_SAVE_VERSION_FAILED" {
		testingT.Errorf("expected error code E_INSTALLER_SAVE_VERSION_FAILED, got %s", appErr.Code)
	}
}
