// Package store — installer_reset_test.go tests ResetInstallers.
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

func setupInstallerResetTestDB(t *testing.T) *DB {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "installer_reset_test.sqlite")
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

func TestResetInstallersAllSuccess(t *testing.T) {
	db := setupInstallerResetTestDB(t)

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

	v1 := &model.InstallerVersion{
		ScriptID:     script1.ID,
		Slug:         "composer",
		Version:      "2.7.0",
		TargetOS:     "all",
		Instructions: "echo install composer",
	}
	v2 := &model.InstallerVersion{
		ScriptID:     script2.ID,
		Slug:         "nodejs",
		Version:      "20.11.0",
		TargetOS:     "win",
		Instructions: "echo install nodejs",
	}

	if err := db.SaveVersion(v1); err != nil {
		t.Fatalf("SaveVersion failed for v1: %v", err)
	}
	if err := db.SaveVersion(v2); err != nil {
		t.Fatalf("SaveVersion failed for v2: %v", err)
	}

	if err := db.ResetInstallers("", true); err != nil {
		t.Fatalf("ResetInstallers(all=true) failed: %v", err)
	}

	scripts, err := db.ListInstallers()
	if err != nil {
		t.Fatalf("ListInstallers failed after reset: %v", err)
	}
	if len(scripts) != 0 {
		t.Errorf("expected 0 installer scripts after reset all, got %d", len(scripts))
	}

	var versionCount int
	row := db.conn.QueryRow("SELECT COUNT(*) FROM installer_versions")
	if err := row.Scan(&versionCount); err != nil {
		t.Fatalf("failed to query installer_versions count: %v", err)
	}
	if versionCount != 0 {
		t.Errorf("expected 0 installer versions after reset all, got %d", versionCount)
	}
}

func TestResetInstallersSpecificSlugSuccess(t *testing.T) {
	db := setupInstallerResetTestDB(t)

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

	v1 := &model.InstallerVersion{
		ScriptID:     script1.ID,
		Slug:         "composer",
		Version:      "2.7.0",
		TargetOS:     "all",
		Instructions: "echo install composer",
	}
	v2 := &model.InstallerVersion{
		ScriptID:     script2.ID,
		Slug:         "nodejs",
		Version:      "20.11.0",
		TargetOS:     "win",
		Instructions: "echo install nodejs",
	}

	if err := db.SaveVersion(v1); err != nil {
		t.Fatalf("SaveVersion failed for v1: %v", err)
	}
	if err := db.SaveVersion(v2); err != nil {
		t.Fatalf("SaveVersion failed for v2: %v", err)
	}

	if err := db.ResetInstallers("composer", false); err != nil {
		t.Fatalf("ResetInstallers(composer, false) failed: %v", err)
	}

	// Composer should be gone
	_, err := db.GetInstallerBySlug("composer")
	if err == nil {
		t.Errorf("expected composer to not be found after reset")
	}

	// Node should still exist
	nodeScript, err := db.GetInstallerBySlug("nodejs")
	if err != nil {
		t.Fatalf("expected nodejs to exist after resetting composer: %v", err)
	}
	if nodeScript.Slug != "nodejs" {
		t.Errorf("expected slug 'nodejs', got %q", nodeScript.Slug)
	}

	// Composer version should be gone
	var composerVersionCount int
	row := db.conn.QueryRow("SELECT COUNT(*) FROM installer_versions WHERE slug = ?", "composer")
	if err := row.Scan(&composerVersionCount); err != nil {
		t.Fatalf("failed to query composer version count: %v", err)
	}
	if composerVersionCount != 0 {
		t.Errorf("expected 0 versions for composer, got %d", composerVersionCount)
	}

	// Node version should still exist
	var nodeVersionCount int
	row = db.conn.QueryRow("SELECT COUNT(*) FROM installer_versions WHERE slug = ?", "nodejs")
	if err := row.Scan(&nodeVersionCount); err != nil {
		t.Fatalf("failed to query nodejs version count: %v", err)
	}
	if nodeVersionCount != 1 {
		t.Errorf("expected 1 version for nodejs, got %d", nodeVersionCount)
	}
}

func TestResetInstallersNonExistentSlug(t *testing.T) {
	db := setupInstallerResetTestDB(t)

	script := &model.InstallerScript{
		Name:         "NodeJS",
		Slug:         "nodejs",
		Description:  "JavaScript runtime",
		TargetOS:     "win",
		Version:      "20.11.0",
		Instructions: "echo install nodejs",
	}
	if err := db.CreateInstaller(script); err != nil {
		t.Fatalf("CreateInstaller failed: %v", err)
	}

	if err := db.ResetInstallers("nonexistent", false); err != nil {
		t.Fatalf("expected ResetInstallers for nonexistent slug to return nil, got: %v", err)
	}

	existingScript, err := db.GetInstallerBySlug("nodejs")
	if err != nil {
		t.Fatalf("expected nodejs to still exist: %v", err)
	}
	if existingScript == nil || existingScript.Slug != "nodejs" {
		t.Errorf("nodejs script corrupted")
	}
}

func TestResetInstallersEmptySlugWhenNotAll(t *testing.T) {
	db := setupInstallerResetTestDB(t)

	err := db.ResetInstallers("", false)
	if err == nil {
		t.Fatalf("expected error when slug is empty and all is false, got nil")
	}

	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != "E_INSTALLER_INVALID_INPUT" {
		t.Errorf("expected code E_INSTALLER_INVALID_INPUT, got %s", appErr.Code)
	}
}

func TestResetInstallersNilDB(t *testing.T) {
	var db *DB

	err1 := db.ResetInstallers("", true)
	if err1 == nil {
		t.Fatalf("expected error on nil db with all=true, got nil")
	}
	var appErr1 *apperror.AppError
	if !errors.As(err1, &appErr1) {
		t.Fatalf("expected AppError, got %T", err1)
	}
	if appErr1.Code != "E_INSTALLER_NIL_DB" {
		t.Errorf("expected code E_INSTALLER_NIL_DB, got %s", appErr1.Code)
	}

	err2 := db.ResetInstallers("composer", false)
	if err2 == nil {
		t.Fatalf("expected error on nil db with all=false, got nil")
	}
	var appErr2 *apperror.AppError
	if !errors.As(err2, &appErr2) {
		t.Fatalf("expected AppError, got %T", err2)
	}
	if appErr2.Code != "E_INSTALLER_NIL_DB" {
		t.Errorf("expected code E_INSTALLER_NIL_DB, got %s", appErr2.Code)
	}
}

func TestResetInstallersClosedDB(t *testing.T) {
	dbConn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("sql.Open failed: %v", err)
	}
	db := &DB{conn: dbConn}
	_ = dbConn.Close()

	err1 := db.ResetInstallers("", true)
	if err1 == nil {
		t.Fatalf("expected error on closed db with all=true, got nil")
	}
	var appErr1 *apperror.AppError
	if !errors.As(err1, &appErr1) {
		t.Fatalf("expected AppError, got %T", err1)
	}
	if appErr1.Code != "E_INSTALLER_RESET_FAILED" {
		t.Errorf("expected code E_INSTALLER_RESET_FAILED, got %s", appErr1.Code)
	}

	err2 := db.ResetInstallers("composer", false)
	if err2 == nil {
		t.Fatalf("expected error on closed db with all=false, got nil")
	}
	var appErr2 *apperror.AppError
	if !errors.As(err2, &appErr2) {
		t.Fatalf("expected AppError, got %T", err2)
	}
	if appErr2.Code != "E_INSTALLER_RESET_FAILED" {
		t.Errorf("expected code E_INSTALLER_RESET_FAILED, got %s", appErr2.Code)
	}
}
