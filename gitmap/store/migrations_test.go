package store

import (
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

func setupMigrationTestMemoryDB(testingT *testing.T) *sql.DB {
	dbConn, errOpen := sql.Open("sqlite", ":memory:")
	if errOpen != nil {
		testingT.Fatalf("failed to open in-memory sqlite db: %v", errOpen)
	}

	return dbConn
}

func TestMigrationsInstallerScriptsTable(testingT *testing.T) {
	dbConn := setupMigrationTestMemoryDB(testingT)
	defer dbConn.Close()

	errMigration := RegisterInstallerScriptsMigration(dbConn, 1, false)
	if errMigration != nil {
		testingT.Fatalf("expected no error from RegisterInstallerScriptsMigration, got: %v", errMigration)
	}

	row := dbConn.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='installer_scripts'")
	var tableName string
	if errScan := row.Scan(&tableName); errScan != nil {
		testingT.Fatalf("expected table installer_scripts to exist: %v", errScan)
	}

	// Test idempotence
	errIdempotent := RegisterInstallerScriptsMigration(dbConn, 1, false)
	if errIdempotent != nil {
		testingT.Fatalf("expected no error on idempotent migration, got: %v", errIdempotent)
	}
}

func TestMigrationsInstallerVersionsTable(testingT *testing.T) {
	dbConn := setupMigrationTestMemoryDB(testingT)
	defer dbConn.Close()

	errMigration := RegisterInstallerVersionsMigration(dbConn, 1, false)
	if errMigration != nil {
		testingT.Fatalf("expected no error from RegisterInstallerVersionsMigration, got: %v", errMigration)
	}

	row := dbConn.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='installer_versions'")
	var tableName string
	if errScan := row.Scan(&tableName); errScan != nil {
		testingT.Fatalf("expected table installer_versions to exist: %v", errScan)
	}

	slugIdxRow := dbConn.QueryRow("SELECT name FROM sqlite_master WHERE type='index' AND name='idx_installer_versions_slug'")
	var slugIdxName string
	if errScan := slugIdxRow.Scan(&slugIdxName); errScan != nil {
		testingT.Fatalf("expected index idx_installer_versions_slug to exist: %v", errScan)
	}

	scriptIDIdxRow := dbConn.QueryRow("SELECT name FROM sqlite_master WHERE type='index' AND name='idx_installer_versions_script_id'")
	var scriptIDIdxName string
	if errScan := scriptIDIdxRow.Scan(&scriptIDIdxName); errScan != nil {
		testingT.Fatalf("expected index idx_installer_versions_script_id to exist: %v", errScan)
	}

	// Test idempotence
	errIdempotent := RegisterInstallerVersionsMigration(dbConn, 1, false)
	if errIdempotent != nil {
		testingT.Fatalf("expected no error on idempotent migration, got: %v", errIdempotent)
	}
}

func TestMigrationsAllInstallerTables(testingT *testing.T) {
	dbConn := setupMigrationTestMemoryDB(testingT)
	defer dbConn.Close()

	errMigration := RegisterInstallerMigration(dbConn, 1, false)
	if errMigration != nil {
		testingT.Fatalf("expected no error from RegisterInstallerMigration, got: %v", errMigration)
	}

	// Verify inserting and reading from both tables with model-compatible columns
	const insertScriptSQL = `INSERT INTO installer_scripts (name, slug, description, target_os, version, instructions)
		VALUES ('Go Tools', 'go-tools', 'Go tooling suite', 'win', 'v1.0.0', '{"steps":[]}')`
	res, errExec := dbConn.Exec(insertScriptSQL)
	if errExec != nil {
		testingT.Fatalf("failed to insert into installer_scripts: %v", errExec)
	}

	scriptID, errLastID := res.LastInsertId()
	if errLastID != nil {
		testingT.Fatalf("failed to get last insert id: %v", errLastID)
	}

	const insertVersionSQL = `INSERT INTO installer_versions (script_id, slug, version, target_os, instructions)
		VALUES (?, 'go-tools', 'v1.0.0', 'win', '{"steps":[]}')`
	if _, errExecVersion := dbConn.Exec(insertVersionSQL, scriptID); errExecVersion != nil {
		testingT.Fatalf("failed to insert into installer_versions: %v", errExecVersion)
	}

	var scriptName string
	errQueryScript := dbConn.QueryRow("SELECT name FROM installer_scripts WHERE slug = 'go-tools'").Scan(&scriptName)
	if errQueryScript != nil || scriptName != "Go Tools" {
		testingT.Fatalf("expected Go Tools, got %q (err: %v)", scriptName, errQueryScript)
	}

	var versionString string
	errQueryVer := dbConn.QueryRow("SELECT version FROM installer_versions WHERE script_id = ?", scriptID).Scan(&versionString)
	if errQueryVer != nil || versionString != "v1.0.0" {
		testingT.Fatalf("expected v1.0.0, got %q (err: %v)", versionString, errQueryVer)
	}

	// Test RegisterInstallerMigrations alias
	errAliasMigration := RegisterInstallerMigrations(dbConn, 2, true)
	if errAliasMigration != nil {
		testingT.Fatalf("expected no error from RegisterInstallerMigrations, got: %v", errAliasMigration)
	}
}

func TestMigrationsDBMigrateInstallers(testingT *testing.T) {
	tempDir := testingT.TempDir()
	dbPath := filepath.Join(tempDir, "test_migrations.db")

	dbInstance, errOpen := OpenAt(dbPath)
	if errOpen != nil {
		testingT.Fatalf("failed to open db: %v", errOpen)
	}
	defer dbInstance.Close()

	errMigrate := dbInstance.MigrateInstallers()
	if errMigrate != nil {
		testingT.Fatalf("expected MigrateInstallers to succeed, got: %v", errMigrate)
	}

	if isScriptsPresent := dbInstance.tableExists("installer_scripts"); !isScriptsPresent {
		testingT.Fatalf("expected installer_scripts table to exist")
	}

	if isVersionsPresent := dbInstance.tableExists("installer_versions"); !isVersionsPresent {
		testingT.Fatalf("expected installer_versions table to exist")
	}
}

func TestMigrationsErrorHandling(testingT *testing.T) {
	dbConn := setupMigrationTestMemoryDB(testingT)
	// Close connection to force an execution failure
	dbConn.Close()

	errScripts := RegisterInstallerScriptsMigration(dbConn, 1, false)
	if errScripts == nil {
		testingT.Fatalf("expected error on closed db for scripts migration")
	}

	var appErrScripts *apperror.AppError
	if !errors.As(errScripts, &appErrScripts) {
		testingT.Fatalf("expected AppError, got %T", errScripts)
	}
	if appErrScripts.Code != "E_INSTALLER_SCRIPTS_MIGRATION_FAILED" {
		testingT.Fatalf("expected E_INSTALLER_SCRIPTS_MIGRATION_FAILED, got %s", appErrScripts.Code)
	}

	errVersions := RegisterInstallerVersionsMigration(dbConn, 1, false)
	if errVersions == nil {
		testingT.Fatalf("expected error on closed db for versions migration")
	}

	var appErrVersions *apperror.AppError
	if !errors.As(errVersions, &appErrVersions) {
		testingT.Fatalf("expected AppError, got %T", errVersions)
	}
	if appErrVersions.Code != "E_INSTALLER_VERSIONS_MIGRATION_FAILED" {
		testingT.Fatalf("expected E_INSTALLER_VERSIONS_MIGRATION_FAILED, got %s", appErrVersions.Code)
	}

	errAll := RegisterInstallerMigration(dbConn, 1, false)
	if errAll == nil {
		testingT.Fatalf("expected error on closed db for all installer migrations")
	}
}
