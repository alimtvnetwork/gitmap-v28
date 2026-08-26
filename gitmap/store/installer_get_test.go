// Package store — installer_get_test.go tests GetInstallerBySlug.
package store

import (
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	_ "modernc.org/sqlite"
)

func setupInstallerGetTestDB(t *testing.T) *DB {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "installer_get_test.sqlite")
	db, err := OpenAt(dbPath)
	if err != nil {
		t.Fatalf("OpenAt failed: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if err := db.MigrateInstallers(); err != nil {
		t.Fatalf("MigrateInstallers failed: %v", err)
	}

	return db
}

func TestGetInstallerSuccess(t *testing.T) {
	db := setupInstallerGetTestDB(t)

	insertQuery := `
	INSERT INTO installer_scripts (name, slug, description, target_os, version, instructions)
	VALUES (?, ?, ?, ?, ?, ?)
	`
	res, err := db.conn.Exec(insertQuery, "Composer", "composer", "PHP dependency manager", "all", "2.7.0", "echo install composer")
	if err != nil {
		t.Fatalf("failed to insert test installer script: %v", err)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		t.Fatalf("failed to get last insert id: %v", err)
	}

	script, err := db.GetInstallerBySlug("composer")
	if err != nil {
		t.Fatalf("GetInstallerBySlug failed unexpectedly: %v", err)
	}

	if script == nil {
		t.Fatalf("expected script to be non-nil")
	}

	if script.ID != lastID {
		t.Errorf("expected ID %d, got %d", lastID, script.ID)
	}
	if script.Name != "Composer" {
		t.Errorf("expected Name 'Composer', got %q", script.Name)
	}
	if script.Slug != "composer" {
		t.Errorf("expected Slug 'composer', got %q", script.Slug)
	}
	if script.Description != "PHP dependency manager" {
		t.Errorf("expected Description 'PHP dependency manager', got %q", script.Description)
	}
	if script.TargetOS != "all" {
		t.Errorf("expected TargetOS 'all', got %q", script.TargetOS)
	}
	if script.Version != "2.7.0" {
		t.Errorf("expected Version '2.7.0', got %q", script.Version)
	}
	if script.Instructions != "echo install composer" {
		t.Errorf("expected Instructions 'echo install composer', got %q", script.Instructions)
	}
	if script.CreatedAt == "" {
		t.Errorf("expected CreatedAt to be non-empty")
	}
	if script.UpdatedAt == "" {
		t.Errorf("expected UpdatedAt to be non-empty")
	}
}

func TestGetInstallerNotFound(t *testing.T) {
	db := setupInstallerGetTestDB(t)

	script, err := db.GetInstallerBySlug("non-existent-slug")
	if err == nil {
		t.Fatalf("expected error for non-existent slug, got nil")
	}
	if script != nil {
		t.Errorf("expected nil script on error, got %+v", script)
	}

	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != "E_INSTALLER_NOT_FOUND" {
		t.Errorf("expected code E_INSTALLER_NOT_FOUND, got %s", appErr.Code)
	}
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected errors.Is(err, apperror.ErrNotFound) to be true")
	}
}

func TestGetInstallerClosedDB(t *testing.T) {
	dbConn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("sql.Open failed: %v", err)
	}
	db := &DB{conn: dbConn}
	_ = dbConn.Close()

	script, err := db.GetInstallerBySlug("composer")
	if err == nil {
		t.Fatalf("expected error on closed db, got nil")
	}
	if script != nil {
		t.Errorf("expected nil script on error, got %+v", script)
	}

	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != "E_INSTALLER_GET_FAILED" {
		t.Errorf("expected code E_INSTALLER_GET_FAILED, got %s", appErr.Code)
	}
}

func TestGetInstallerNilDB(t *testing.T) {
	var db *DB
	script, err := db.GetInstallerBySlug("composer")
	if err == nil {
		t.Fatalf("expected error on nil db, got nil")
	}
	if script != nil {
		t.Errorf("expected nil script on error, got %+v", script)
	}

	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != "E_INSTALLER_NIL_DB" {
		t.Errorf("expected code E_INSTALLER_NIL_DB, got %s", appErr.Code)
	}
}
