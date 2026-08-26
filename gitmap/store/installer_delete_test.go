// Package store — installer_delete_test.go tests DeleteInstaller.
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

func setupInstallerDeleteTestDB(testingT *testing.T) *DB {
	testingT.Helper()
	tempDir := testingT.TempDir()
	dbPath := filepath.Join(tempDir, "test_installer_delete.db")

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

func TestDeleteInstallerSuccess(testingT *testing.T) {
	dbInstance := setupInstallerDeleteTestDB(testingT)

	script := &model.InstallerScript{
		Name:         "Golang Suite",
		Slug:         "go-suite",
		Description:  "Go tooling installer",
		TargetOS:     "win",
		Version:      "v1.0.0",
		Instructions: `{"steps":[{"action":"install","pkg":"go"}]}`,
	}
	if errCreate := dbInstance.CreateInstaller(script); errCreate != nil {
		testingT.Fatalf("CreateInstaller failed: %v", errCreate)
	}

	version := &model.InstallerVersion{
		ScriptID:     script.ID,
		Slug:         "go-suite",
		Version:      "v1.0.0",
		TargetOS:     "win",
		Instructions: `{"steps":[{"action":"install","pkg":"go"}]}`,
	}
	if errVersion := dbInstance.SaveVersion(version); errVersion != nil {
		testingT.Fatalf("SaveVersion failed: %v", errVersion)
	}

	if errDelete := dbInstance.DeleteInstaller("go-suite"); errDelete != nil {
		testingT.Fatalf("expected DeleteInstaller to succeed, got: %v", errDelete)
	}

	_, errGet := dbInstance.GetInstallerBySlug("go-suite")
	if errGet == nil {
		testingT.Fatalf("expected GetInstallerBySlug to fail after deletion, got nil")
	}
	if !errors.Is(errGet, apperror.ErrNotFound) {
		testingT.Errorf("expected ErrNotFound for deleted installer, got %v", errGet)
	}

	var versionCount int
	row := dbInstance.conn.QueryRow("SELECT COUNT(*) FROM installer_versions WHERE slug = ?", "go-suite")
	if errScan := row.Scan(&versionCount); errScan != nil {
		testingT.Fatalf("failed to count versions: %v", errScan)
	}
	if versionCount != 0 {
		testingT.Errorf("expected 0 versions remaining, got %d", versionCount)
	}
}

func TestDeleteInstallerNotFound(testingT *testing.T) {
	dbInstance := setupInstallerDeleteTestDB(testingT)

	errDelete := dbInstance.DeleteInstaller("nonexistent-slug")
	if errDelete == nil {
		testingT.Fatalf("expected error for non-existent slug, got nil")
	}

	var appErr *apperror.AppError
	if !errors.As(errDelete, &appErr) {
		testingT.Fatalf("expected AppError, got %T", errDelete)
	}
	if appErr.Code != "E_INSTALLER_NOT_FOUND" {
		testingT.Errorf("expected code E_INSTALLER_NOT_FOUND, got %s", appErr.Code)
	}
	if !errors.Is(errDelete, apperror.ErrNotFound) {
		testingT.Errorf("expected errors.Is(err, apperror.ErrNotFound) to be true")
	}
}

func TestDeleteInstallerEmptySlug(testingT *testing.T) {
	dbInstance := setupInstallerDeleteTestDB(testingT)

	errDelete := dbInstance.DeleteInstaller("")
	if errDelete == nil {
		testingT.Fatalf("expected error for empty slug, got nil")
	}

	var appErr *apperror.AppError
	if !errors.As(errDelete, &appErr) {
		testingT.Fatalf("expected AppError, got %T", errDelete)
	}
	if appErr.Code != "E_INSTALLER_INVALID_INPUT" {
		testingT.Errorf("expected code E_INSTALLER_INVALID_INPUT, got %s", appErr.Code)
	}
}

func TestDeleteInstallerNilDB(testingT *testing.T) {
	var dbInstance *DB
	errDelete := dbInstance.DeleteInstaller("go-suite")
	if errDelete == nil {
		testingT.Fatalf("expected error on nil db, got nil")
	}

	var appErr *apperror.AppError
	if !errors.As(errDelete, &appErr) {
		testingT.Fatalf("expected AppError, got %T", errDelete)
	}
	if appErr.Code != "E_INSTALLER_NIL_DB" {
		testingT.Errorf("expected code E_INSTALLER_NIL_DB, got %s", appErr.Code)
	}
}

func TestDeleteInstallerClosedDB(testingT *testing.T) {
	dbConn, errOpen := sql.Open("sqlite", ":memory:")
	if errOpen != nil {
		testingT.Fatalf("sql.Open failed: %v", errOpen)
	}
	dbInstance := &DB{conn: dbConn}
	_ = dbConn.Close()

	errDelete := dbInstance.DeleteInstaller("go-suite")
	if errDelete == nil {
		testingT.Fatalf("expected error on closed db, got nil")
	}

	var appErr *apperror.AppError
	if !errors.As(errDelete, &appErr) {
		testingT.Fatalf("expected AppError, got %T", errDelete)
	}
	if appErr.Code != "E_INSTALLER_DELETE_FAILED" {
		testingT.Errorf("expected code E_INSTALLER_DELETE_FAILED, got %s", appErr.Code)
	}
}
