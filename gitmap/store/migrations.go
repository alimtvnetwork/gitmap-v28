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
	"fmt"
	"os"
	"strings"
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
	rows, err := db.conn.Query(fmt.Sprintf("PRAGMA table_info(%q)", table))
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
	row := db.conn.QueryRow(
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
