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

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
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

