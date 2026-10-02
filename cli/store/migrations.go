// Package store — migration helpers.
//
// Each migration follows the detect-then-act pattern:
//
//  1. Inspect schema with PRAGMA table_info to learn the current shape.
//  2. Only run ALTER if it is actually required.
//  3. If a write still fails, log a *contextual* warning that names the
//     table and column so users (and downstream tooling) can act on it.
//
// This avoids spurious "no such column" warnings on fresh installs and
// makes the migration log self-explanatory across every OS / SQLite
// driver variant (Windows mingw vs. Linux glibc vs. macOS).
package store

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// MigrationReport summarizes a single Migrate() run for `gitmap db-migrate`.
type MigrationReport struct {
	TablesEnsured int
	StepsRun      []string
	StepsSkipped  []string
	Warnings      []string
}

// columnExists reports whether table.column exists. Returns false on any
// query error (treated as "not present" so callers can skip safely).
func (db *DB) columnExists(table, column string) bool {
	rows, err := QueryWrapper(db.conn, fmt.Sprintf("PRAGMA table_info(%q)", table)).Destruct()
	if err != nil {
		return false
	}

	defer rows.Close()

	for rows.Next() {
		var (
			cid     int
			name    string
			ctype   string
			notnull int
			dflt    any
			pk      int
		)

		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			continue
		}

		if name == column {
			return true
		}
	}

	return false
}

// tableExists reports whether a table is present in the active database.
func (db *DB) tableExists(table string) bool {
	row := QueryRowWrapper(db.conn,
		"SELECT 1 FROM sqlite_master WHERE type='table' AND name=?", table)

	var seen int

	return row.Scan(&seen) == nil && seen == 1
}

// logMigrationFailure prints a uniform, contextual warning for any migration
// statement that fails for an unexpected reason.
func logMigrationFailure(table, column, action string, err error, stmt string) {
	fmt.Fprintf(os.Stderr,
		"  ⚠ Migration failed: table=%s column=%s action=%s: %v\n"+
			"      statement: %s\n"+
			"      hint: run `gitmap db-migrate --verbose` to retry, "+
			"or `gitmap db-reset --confirm` to rebuild the schema.\n",
		table, column, action, err, stmt)
}

// isBenignAlterError reports whether err can be safely ignored for ALTER
// migrations: the column is already missing, already renamed, or duplicate.
func isBenignAlterError(err error) bool {
	if err == nil {
		return true
	}

	msg := strings.ToLower(err.Error())
	for _, needle := range []string{
		"no such column",
		"no such table",
		"duplicate column",
		"already exists",
	} {
		if strings.Contains(msg, needle) {
			return true
		}
	}

	return false
}

// SQLCreateInstallerScripts creates the installer_scripts table.
const SQLCreateInstallerScripts = `CREATE TABLE IF NOT EXISTS installer_scripts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	slug TEXT NOT NULL UNIQUE,
	description TEXT DEFAULT '',
	target_os TEXT DEFAULT '',
	version TEXT DEFAULT '',
	instructions TEXT DEFAULT '',
	created_at TEXT DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);`

// SQLCreateInstallerVersions creates the installer_versions table.
const SQLCreateInstallerVersions = `CREATE TABLE IF NOT EXISTS installer_versions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	script_id INTEGER NOT NULL DEFAULT 0,
	slug TEXT NOT NULL,
	version TEXT NOT NULL,
	target_os TEXT NOT NULL DEFAULT '',
	instructions TEXT NOT NULL DEFAULT '',
	created_at TEXT DEFAULT CURRENT_TIMESTAMP
);`

// SQLCreateInstallerVersionsIndex creates the index on installer_versions(slug).
const SQLCreateInstallerVersionsIndex = `CREATE INDEX IF NOT EXISTS idx_installer_versions_slug ON installer_versions(slug);`

// SQLCreateInstallerVersionsScriptIDIndex creates the index on installer_versions(script_id).
const SQLCreateInstallerVersionsScriptIDIndex = `CREATE INDEX IF NOT EXISTS idx_installer_versions_script_id ON installer_versions(script_id);`

// RegisterInstallerScriptsMigration creates the installer_scripts table if it does not exist.
func RegisterInstallerScriptsMigration(dbConn *sql.DB, migrationVersion int, isForce bool) error {
	if _, errExec := dbConn.Exec(SQLCreateInstallerScripts); errExec != nil {
		migrationErr := apperror.Wrap(errExec, "RegisterInstallerScriptsMigration", map[string]any{
			"table":   "installer_scripts",
			"version": migrationVersion,
			"force":   isForce,
		})
		migrationErr.Code = "E_INSTALLER_SCRIPTS_MIGRATION_FAILED"

		return migrationErr
	}

	return nil
}

// RegisterInstallerVersionsMigration creates the installer_versions table and its index if they do not exist.
func RegisterInstallerVersionsMigration(dbConn *sql.DB, migrationVersion int, isForce bool) error {
	if _, errExec := dbConn.Exec(SQLCreateInstallerVersions); errExec != nil {
		migrationErr := apperror.Wrap(errExec, "RegisterInstallerVersionsMigration", map[string]any{
			"table":   "installer_versions",
			"version": migrationVersion,
			"force":   isForce,
		})
		migrationErr.Code = "E_INSTALLER_VERSIONS_MIGRATION_FAILED"

		return migrationErr
	}

	if _, errIdx := dbConn.Exec(SQLCreateInstallerVersionsIndex); errIdx != nil {
		indexErr := apperror.Wrap(errIdx, "RegisterInstallerVersionsMigration", map[string]any{
			"table":   "installer_versions",
			"index":   "idx_installer_versions_slug",
			"version": migrationVersion,
			"force":   isForce,
		})
		indexErr.Code = "E_INSTALLER_VERSIONS_MIGRATION_FAILED"

		return indexErr
	}

	if _, errIdx := dbConn.Exec(SQLCreateInstallerVersionsScriptIDIndex); errIdx != nil {
		indexErr := apperror.Wrap(errIdx, "RegisterInstallerVersionsMigration", map[string]any{
			"table":   "installer_versions",
			"index":   "idx_installer_versions_script_id",
			"version": migrationVersion,
			"force":   isForce,
		})
		indexErr.Code = "E_INSTALLER_VERSIONS_MIGRATION_FAILED"

		return indexErr
	}

	return nil
}

// RegisterInstallerMigration creates installer_scripts and installer_versions tables if they do not exist.
func RegisterInstallerMigration(dbConn *sql.DB, migrationVersion int, isForce bool) error {
	if errScripts := RegisterInstallerScriptsMigration(dbConn, migrationVersion, isForce); errScripts != nil {
		return errScripts
	}

	if errVersions := RegisterInstallerVersionsMigration(dbConn, migrationVersion, isForce); errVersions != nil {
		return errVersions
	}

	if errPins := RegisterInstallerPinsMigration(dbConn, migrationVersion, isForce); errPins != nil {
		return errPins
	}

	return nil
}

// RegisterInstallerMigrations applies all installer-related migrations.
func RegisterInstallerMigrations(dbConn *sql.DB, migrationVersion int, isForce bool) error {
	return RegisterInstallerMigration(dbConn, migrationVersion, isForce)
}

// MigrateInstallers creates the installer_scripts and installer_versions tables on the DB.
func (dbInstance *DB) MigrateInstallers() error {
	return RegisterInstallerMigration(dbInstance.conn, 1, false)
}

func isTablePresent(dbConn *sql.DB, tableName string) bool {
	var count int
	query := "SELECT 1 FROM sqlite_master WHERE type='table' AND name=?"
	err := dbConn.QueryRow(query, tableName).Scan(&count)
	hasTable := (err == nil && count == 1)

	return hasTable
}

func deleteConflictingReleases(dbConn *sql.DB) error {
	const sqlPrune = `DELETE FROM Release
	WHERE RepoId IN (
		SELECT RepoId FROM Repo WHERE RepoId NOT IN (
			SELECT MIN(RepoId) FROM Repo GROUP BY LOWER(AbsolutePath)
		)
	) AND EXISTS (
		SELECT 1 FROM Release r2
		JOIN Repo rDup ON rDup.RepoId = Release.RepoId
		JOIN Repo rSurv ON LOWER(rSurv.AbsolutePath) = LOWER(rDup.AbsolutePath)
		WHERE r2.RepoId = rSurv.RepoId AND r2.Tag = Release.Tag AND rSurv.RepoId != rDup.RepoId
	)`
	if _, err := dbConn.Exec(sqlPrune); err != nil {
		return apperror.WrapSimple(err, "deleteConflictingReleases")
	}

	return nil
}

func executeReleaseRemap(dbConn *sql.DB) error {
	const sqlRemap = `UPDATE Release SET RepoId = (
		SELECT MIN(r2.RepoId) FROM Repo r1
		JOIN Repo r2 ON LOWER(r1.AbsolutePath) = LOWER(r2.AbsolutePath)
		WHERE r1.RepoId = Release.RepoId
	) WHERE RepoId IN (
		SELECT RepoId FROM Repo WHERE RepoId NOT IN (
			SELECT MIN(RepoId) FROM Repo GROUP BY LOWER(AbsolutePath)
		)
	)`
	if _, err := dbConn.Exec(sqlRemap); err != nil {
		return apperror.WrapSimple(err, "remapReleaseRepoRefs")
	}

	return nil
}

func remapReleaseRepoRefs(dbConn *sql.DB) error {
	hasTable := isTablePresent(dbConn, "Release")
	if !hasTable {
		return nil
	}
	if err := deleteConflictingReleases(dbConn); err != nil {
		return err
	}

	return executeReleaseRemap(dbConn)
}

func insertRemappedGroupRepos(dbConn *sql.DB) error {
	const sqlInsert = `INSERT OR IGNORE INTO GroupRepo (GroupId, RepoId)
	SELECT gr.GroupId, (
		SELECT MIN(r2.RepoId) FROM Repo r1
		JOIN Repo r2 ON LOWER(r1.AbsolutePath) = LOWER(r2.AbsolutePath)
		WHERE r1.RepoId = gr.RepoId
	) FROM GroupRepo gr WHERE gr.RepoId IN (
		SELECT RepoId FROM Repo WHERE RepoId NOT IN (
			SELECT MIN(RepoId) FROM Repo GROUP BY LOWER(AbsolutePath)
		)
	)`
	if _, err := dbConn.Exec(sqlInsert); err != nil {
		return apperror.WrapSimple(err, "insertRemappedGroupRepos")
	}

	return nil
}

func deleteStaleGroupRepos(dbConn *sql.DB) error {
	const sqlDelete = `DELETE FROM GroupRepo WHERE RepoId IN (
		SELECT RepoId FROM Repo WHERE RepoId NOT IN (
			SELECT MIN(RepoId) FROM Repo GROUP BY LOWER(AbsolutePath)
		)
	)`
	if _, err := dbConn.Exec(sqlDelete); err != nil {
		return apperror.WrapSimple(err, "deleteStaleGroupRepos")
	}

	return nil
}

func remapGroupRepoRefs(dbConn *sql.DB) error {
	hasTable := isTablePresent(dbConn, "GroupRepo")
	if !hasTable {
		return nil
	}
	if err := insertRemappedGroupRepos(dbConn); err != nil {
		return err
	}

	return deleteStaleGroupRepos(dbConn)
}

func remapVersionProbeRefs(dbConn *sql.DB) error {
	hasTable := isTablePresent(dbConn, "VersionProbe")
	if !hasTable {
		return nil
	}
	const sqlRemap = `UPDATE VersionProbe SET RepoId = (
		SELECT MIN(r2.RepoId) FROM Repo r1
		JOIN Repo r2 ON LOWER(r1.AbsolutePath) = LOWER(r2.AbsolutePath)
		WHERE r1.RepoId = VersionProbe.RepoId
	) WHERE RepoId IN (
		SELECT RepoId FROM Repo WHERE RepoId NOT IN (
			SELECT MIN(RepoId) FROM Repo GROUP BY LOWER(AbsolutePath)
		)
	)`
	if _, err := dbConn.Exec(sqlRemap); err != nil {
		return apperror.WrapSimple(err, "remapVersionProbeRefs")
	}

	return nil
}

func remapAllRepoChildRefs(dbConn *sql.DB) error {
	if err := remapReleaseRepoRefs(dbConn); err != nil {
		return err
	}
	if err := remapGroupRepoRefs(dbConn); err != nil {
		return err
	}

	return remapVersionProbeRefs(dbConn)
}

func deduplicateRepoRows(dbConn *sql.DB) error {
	hasTable := isTablePresent(dbConn, "Repo")
	if !hasTable {
		return nil
	}
	if _, err := dbConn.Exec(constants.SQLDeduplicateRepos); err != nil {
		return apperror.WrapSimple(err, "deduplicateRepoRows")
	}

	return nil
}

func recreateRepoPathIndex(dbConn *sql.DB) error {
	hasTable := isTablePresent(dbConn, "Repo")
	if !hasTable {
		return nil
	}
	if _, err := dbConn.Exec(constants.SQLDropRepoAbsPathIndex); err != nil {
		return apperror.WrapSimple(err, "dropRepoAbsPathIndex")
	}
	if _, err := dbConn.Exec(constants.SQLCreateAbsPathIndex); err != nil {
		return apperror.WrapSimple(err, "createRepoAbsPathIndex")
	}

	return nil
}

func executeScanFolderRemap(dbConn *sql.DB) error {
	const sqlRemap = `UPDATE Repo SET ScanFolderId = (
		SELECT MIN(sf2.ScanFolderId) FROM ScanFolder sf1
		JOIN ScanFolder sf2 ON LOWER(sf1.AbsolutePath) = LOWER(sf2.AbsolutePath)
		WHERE sf1.ScanFolderId = Repo.ScanFolderId
	) WHERE ScanFolderId IS NOT NULL AND ScanFolderId IN (
		SELECT ScanFolderId FROM ScanFolder WHERE ScanFolderId NOT IN (
			SELECT MIN(ScanFolderId) FROM ScanFolder GROUP BY LOWER(AbsolutePath)
		)
	)`
	if _, err := dbConn.Exec(sqlRemap); err != nil {
		return apperror.WrapSimple(err, "remapScanFolderRefs")
	}

	return nil
}

func remapScanFolderRefs(dbConn *sql.DB) error {
	hasRepo := isTablePresent(dbConn, "Repo")
	hasScanFolder := isTablePresent(dbConn, "ScanFolder")
	canRemap := hasRepo && hasScanFolder
	if !canRemap {
		return nil
	}

	return executeScanFolderRemap(dbConn)
}

func deduplicateScanFolderRows(dbConn *sql.DB) error {
	hasTable := isTablePresent(dbConn, "ScanFolder")
	if !hasTable {
		return nil
	}
	if _, err := dbConn.Exec(constants.SQLDeduplicateScanFolders); err != nil {
		return apperror.WrapSimple(err, "deduplicateScanFolderRows")
	}

	return nil
}

func recreateScanFolderPathIndex(dbConn *sql.DB) error {
	hasTable := isTablePresent(dbConn, "ScanFolder")
	if !hasTable {
		return nil
	}
	if _, err := dbConn.Exec(constants.SQLDropScanFolderPathIndex); err != nil {
		return apperror.WrapSimple(err, "dropScanFolderPathIndex")
	}
	if _, err := dbConn.Exec(constants.SQLCreateScanFolderPathIndex); err != nil {
		return apperror.WrapSimple(err, "createScanFolderPathIndex")
	}

	return nil
}

func bumpSchemaVersionTo33(dbConn *sql.DB) error {
	hasTable := isTablePresent(dbConn, "Setting")
	if !hasTable {
		return nil
	}
	const sqlBump = `INSERT INTO Setting (Key, Value) VALUES ('schema_version', '33')
		ON CONFLICT(Key) DO UPDATE SET Value='33'`
	if _, err := dbConn.Exec(sqlBump); err != nil {
		return apperror.WrapSimple(err, "bumpSchemaVersionTo33")
	}

	return nil
}

func deduplicateRepoAndScanFolder(dbConn *sql.DB) error {
	if err := deduplicateRepoRows(dbConn); err != nil {
		return err
	}
	if err := recreateRepoPathIndex(dbConn); err != nil {
		return err
	}
	if err := remapScanFolderRefs(dbConn); err != nil {
		return err
	}
	if err := deduplicateScanFolderRows(dbConn); err != nil {
		return err
	}

	return recreateScanFolderPathIndex(dbConn)
}

// Migration_AddRepoAbsolutePathCollateNoCase performs schema version 33 migration:
// deduplicating Repo & ScanFolder records, remapping foreign keys, and applying COLLATE NOCASE.
func Migration_AddRepoAbsolutePathCollateNoCase(dbConn *sql.DB) error {
	if err := remapAllRepoChildRefs(dbConn); err != nil {
		return err
	}
	if err := deduplicateRepoAndScanFolder(dbConn); err != nil {
		return err
	}

	return bumpSchemaVersionTo33(dbConn)
}

// Migration_AddRepoAbsolutePathCollateNoCase runs the collation and deduplication migration on the DB instance.
func (dbInstance *DB) Migration_AddRepoAbsolutePathCollateNoCase() error {
	return Migration_AddRepoAbsolutePathCollateNoCase(dbInstance.conn)
}
