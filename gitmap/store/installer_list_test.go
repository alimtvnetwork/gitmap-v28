// Package store — installer_list_test.go tests ListInstallers.
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

func setupInstallerListTestDB(t *testing.T) *DB {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "installer_list_test.sqlite")
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

func TestListInstallersSuccess(t *testing.T) {
	db := setupInstallerListTestDB(t)

	script1 := &model.InstallerScript{
		Name:         "Composer",
		Slug:         "composer",
		Description:  "PHP dependency manager",
		TargetOS:     "all",
		Version:      "2.7.0",
		Instructions: "echo install composer",
	}
	script2 := &model.InstallerScript{
		Name:         "NodeJS",
		Slug:         "nodejs",
		Description:  "JavaScript runtime",
		TargetOS:     "win",
		Version:      "20.11.0",
		Instructions: "echo install nodejs",
	}

	if err := db.CreateInstaller(script1); err != nil {
		t.Fatalf("CreateInstaller failed for script1: %v", err)
	}
	if err := db.CreateInstaller(script2); err != nil {
		t.Fatalf("CreateInstaller failed for script2: %v", err)
	}

	scripts, err := db.ListInstallers()
	if err != nil {
		t.Fatalf("ListInstallers failed unexpectedly: %v", err)
	}

	if len(scripts) != 2 {
		t.Fatalf("expected 2 scripts, got %d", len(scripts))
	}

	if scripts[0].ID != script1.ID || scripts[0].Slug != "composer" || scripts[0].Name != "Composer" {
		t.Errorf("script[0] mismatch: %+v", scripts[0])
	}
	if scripts[0].Description != "PHP dependency manager" || scripts[0].TargetOS != "all" || scripts[0].Version != "2.7.0" {
		t.Errorf("script[0] field mismatch: %+v", scripts[0])
	}
	if scripts[0].Instructions != "echo install composer" || scripts[0].CreatedAt == "" || scripts[0].UpdatedAt == "" {
		t.Errorf("script[0] metadata mismatch: %+v", scripts[0])
	}

	if scripts[1].ID != script2.ID || scripts[1].Slug != "nodejs" || scripts[1].Name != "NodeJS" {
		t.Errorf("script[1] mismatch: %+v", scripts[1])
	}
	if scripts[1].Description != "JavaScript runtime" || scripts[1].TargetOS != "win" || scripts[1].Version != "20.11.0" {
		t.Errorf("script[1] field mismatch: %+v", scripts[1])
	}
	if scripts[1].Instructions != "echo install nodejs" || scripts[1].CreatedAt == "" || scripts[1].UpdatedAt == "" {
		t.Errorf("script[1] metadata mismatch: %+v", scripts[1])
	}
}

func TestListInstallersEmpty(t *testing.T) {
	db := setupInstallerListTestDB(t)

	scripts, err := db.ListInstallers()
	if err != nil {
		t.Fatalf("expected no error on empty list, got: %v", err)
	}

	if len(scripts) != 0 {
		t.Errorf("expected 0 scripts, got %d", len(scripts))
	}
}

func TestListInstallersClosedDB(t *testing.T) {
	dbConn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("sql.Open failed: %v", err)
	}
	db := &DB{conn: dbConn}
	_ = dbConn.Close()

	scripts, err := db.ListInstallers()
	if err == nil {
		t.Fatalf("expected error on closed db, got nil")
	}
	if scripts != nil {
		t.Errorf("expected nil scripts on error, got %+v", scripts)
	}

	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != "E_INSTALLER_LIST_FAILED" {
		t.Errorf("expected code E_INSTALLER_LIST_FAILED, got %s", appErr.Code)
	}
}

func TestListInstallersNilDB(t *testing.T) {
	var db *DB
	scripts, err := db.ListInstallers()
	if err == nil {
		t.Fatalf("expected error on nil db, got nil")
	}
	if scripts != nil {
		t.Errorf("expected nil scripts on error, got %+v", scripts)
	}

	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != "E_INSTALLER_NIL_DB" {
		t.Errorf("expected code E_INSTALLER_NIL_DB, got %s", appErr.Code)
	}
}
