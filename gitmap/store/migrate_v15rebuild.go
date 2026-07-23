// Package store — migrate_v15rebuild.go provides the generic table-rebuild
// helper used by Phase 1.2+ v15 migrations. SQLite has no ALTER COLUMN, so
// renaming a primary key (Id → {Table}Id) and the table itself requires:
//
//  1. CREATE the new singular table with v15 schema.
//  2. INSERT INTO new (newCols) SELECT oldCols FROM old.
//  3. Verify row-count parity.
//  4. DROP the legacy plural table.
//
// Foreign keys are temporarily disabled so child tables (which still
// reference the legacy name at this point) survive the rename. Subsequent
// phases rebuild those child tables with the new REFERENCES clause.
//
// All steps are idempotent: detect-then-act via tableExists() means fresh
// installs and re-runs are safe no-ops.
package store

import (
	"fmt"
	"os"
	"strings"
)

// v15RebuildSpec describes one legacy-plural → v15-singular table rebuild.
type v15RebuildSpec struct {
	OldTable      string // e.g. "Groups"
	NewTable      string // e.g. "Group"
	NewCreateSQL  string // full CREATE TABLE for the new singular table
	OldColumnList string // exact column list to SELECT from old (e.g. "Id, Name, ...")
	NewColumnList string // exact column list to INSERT into new (e.g. "GroupId, Name, ...")
	StartMsg      string
	DoneMsg       string
}

// runV15Rebuild executes a single rebuild spec idempotently.
func (db *DB) runV15Rebuild(spec v15RebuildSpec) error {
	if !db.tableExists(spec.OldTable) {
		return nil // fresh install, nothing to migrate
	}

	if db.hasCanonicalV15Shape(spec) {
		return nil
	}

	if db.tableExists(spec.NewTable) {
		// Both exist — the new table was created by a prior partial run
		// (e.g., via the standard CREATE TABLE IF NOT EXISTS pass). Drop
		// the legacy and let the standard pass own the new one going
		// forward. Data preservation is impossible in this edge case
		// because the new table is presumed empty/fresh.
		_, _ = db.conn.Exec("DROP TABLE IF EXISTS " + spec.OldTable)

		return nil
	}

	if spec.StartMsg != "" {
		fmt.Println(spec.StartMsg)
	}

	// Adapt the SELECT list when the legacy `Id` column has already been
	// renamed to `{OldTable}Id` (e.g. on a database first created at v3.5.0+
	// where the standard CREATE TABLE pass produced the v15 schema before
	// this rebuild ran). Without this the INSERT ... SELECT fails with
	// "no such column: Id".
	spec = adaptOldColumnList(db, spec)

	oldCount, err := db.countRows(spec.OldTable)
	if err != nil {
		return fmt.Errorf("count %s: %w", spec.OldTable, err)
	}

	if err := db.execV15Rebuild(spec); err != nil {
		return err
	}

	newCount, err := db.countRows(spec.NewTable)
	if err != nil {
		return fmt.Errorf("count %s: %w", spec.NewTable, err)
	}

	if oldCount != newCount {
		fmt.Fprintf(os.Stderr,
			"  ✗ v15 %s→%s row-count mismatch: old=%d new=%d\n",
			spec.OldTable, spec.NewTable, oldCount, newCount)

		return fmt.Errorf("v15 %s→%s row-count mismatch: old=%d new=%d",
			spec.OldTable, spec.NewTable, oldCount, newCount)
	}

	if spec.DoneMsg != "" {
		fmt.Println(spec.DoneMsg)
	}

	return nil
}

// hasCanonicalV15Shape reports whether a same-name rebuild spec can be skipped
// because the canonical table already exposes the migrated PK shape. This
// prevents staging-table rebuilds (TaskType -> TaskType_v15 -> TaskType, etc.)
// from re-running on every openDB() after the first successful migration.
func (db *DB) hasCanonicalV15Shape(spec v15RebuildSpec) bool {
	if !strings.Contains(spec.NewTable, "_v15") {
		return false
	}
	if db.columnExists(spec.OldTable, "Id") {
		return false
	}

	v15PK := derivePKColumnName(spec)

	return v15PK != "" && db.columnExists(spec.OldTable, v15PK)
}

// execV15Rebuild performs the table-rebuild dance for one spec.
func (db *DB) execV15Rebuild(spec v15RebuildSpec) error {
	if _, err := db.conn.Exec("PRAGMA foreign_keys = OFF"); err != nil {
		return fmt.Errorf("disable foreign keys: %w", err)
	}

	defer func() {
		_, _ = db.conn.Exec("PRAGMA foreign_keys = ON")
	}()

	if _, err := db.conn.Exec(spec.NewCreateSQL); err != nil {
		return fmt.Errorf("create %s: %w", spec.NewTable, err)
	}

	copySQL := fmt.Sprintf( //nolint:gosec // G201: identifiers come from internal RebuildSpec, not user input.
		`INSERT INTO %q (%s) SELECT %s FROM %q`,
		spec.NewTable, spec.NewColumnList, spec.OldColumnList, spec.OldTable,
	)

	if _, err := db.conn.Exec(copySQL); err != nil {
		return fmt.Errorf("copy %s→%s: %w", spec.OldTable, spec.NewTable, err)
	}

	if _, err := db.conn.Exec(`DROP TABLE "` + spec.OldTable + `"`); err != nil {
		return fmt.Errorf("drop %s: %w", spec.OldTable, err)
	}

	return nil
}

// adaptOldColumnList rewrites the leading `Id` token in spec.OldColumnList to
// `{OldTable}Id` when the existing OldTable no longer has an `Id` column but
// already exposes the v15 PK column. This handles two scenarios:
//
//  1. Same-name PK-rename rebuilds (e.g. GoProjectMetadata,  PendingTask)
//     where the standard CREATE TABLE IF NOT EXISTS pass in Migrate() may
//     have already produced the v15-shaped table on a fresh-since-v3.5.0
//     install.
//  2. Idempotent re-runs against a partially-migrated database.
//
// The function is a no-op when the legacy `Id` column is still present, so it
// never breaks the genuine legacy → v15 migration path.
func adaptOldColumnList(db *DB, spec v15RebuildSpec) v15RebuildSpec {
	hasLegacyID := db.columnExists(spec.OldTable, "Id")
	if hasLegacyID {
		return spec
	}

	v15PK := derivePKColumnName(spec)
	if v15PK == "" || !db.columnExists(spec.OldTable, v15PK) {
		return spec
	}

	spec.OldColumnList = replaceLeadingIDToken(spec.OldColumnList, v15PK)

	return spec
}

// derivePKColumnName returns the v15 PK column name for spec, derived from
// the first token of NewColumnList. Returns "" if NewColumnList is empty.
func derivePKColumnName(spec v15RebuildSpec) string {
	first, _, found := strings.Cut(spec.NewColumnList, ",")
	if !found {
		first = spec.NewColumnList
	}

	return strings.TrimSpace(first)
}

// replaceLeadingIDToken swaps the first comma-separated token of list with
// replacement, preserving the rest of the column list verbatim. When the
// leading token is not exactly "Id", list is returned unchanged.
func replaceLeadingIDToken(list, replacement string) string {
	first, rest, found := strings.Cut(list, ",")
	if strings.TrimSpace(first) != "Id" {
		return list
	}
	if !found {
		return replacement
	}

	return replacement + "," + rest
}
