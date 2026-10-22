package store

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// addColumnIfNotExists executes stmt and swallows benign "duplicate column"
// and the broader benign-ALTER family (table missing, etc.) so that fresh
// installs do not surface migration warnings. Any other failure is logged
// with the offending statement so users can diagnose it.
func (db *DB) addColumnIfNotExists(stmt string) {
	_, err := db.conn.Exec(stmt)
	if err == nil || isBenignAlterError(err) {
		return
	}

	fmt.Fprintf(os.Stderr, "  ⚠ Migration failed: %v (statement: %s)\n", err, stmt)
}

// migrateSourceColumn adds the Source column to existing Releases tables.
func (db *DB) migrateSourceColumn() {
	db.addColumnIfNotExists(constants.SQLAddSourceColumn)
}

// migrateNotesColumn adds the Notes column to existing Release tables.
func (db *DB) migrateNotesColumn() {
	db.addColumnIfNotExists(constants.SQLAddNotesColumn)
}

// preV15Phase2EnsureReleaseColumns ensures the legacy `Releases` table has
// `Source` and `Notes` columns BEFORE the v15 rebuild copies it into
// `Release`. Without this, very old installs (pre-Source/pre-Notes) would
// fail the column-by-name SELECT in migrateV15Phase2. No-op when the legacy
// table is absent (fresh install).
func (db *DB) preV15Phase2EnsureReleaseColumns() {
	if !db.tableExists("Releases") {
		return
	}

	// Use raw ALTERs targeting the legacy plural table directly. The
	// constants.SQLAddSourceColumn / SQLAddNotesColumn now target the new
	// singular `Release` table and would fail here.
	db.addColumnIfNotExists(`ALTER TABLE Releases ADD COLUMN Source TEXT DEFAULT 'release'`)
	db.addColumnIfNotExists(`ALTER TABLE Releases ADD COLUMN Notes TEXT DEFAULT ''`)
}

// migrateRepoVersionColumns adds CurrentVersionTag and CurrentVersionNum to Repos.
func (db *DB) migrateRepoVersionColumns() {
	db.addColumnIfNotExists(constants.SQLAddCurrentVersionTag)
	db.addColumnIfNotExists(constants.SQLAddCurrentVersionNum)
}

// migrateRepoIdentifiedTransport adds the IdentifiedTransport column
// to existing Repo tables (migration 008). Plan 03 step 2: persist
// the SSH-vs-HTTPS verdict so reclone-class commands (cfr, cfrp,
// clone-now) can honor it without re-detecting from disk. Empty
// string = "unknown, backfill from origin on next scan"; the column
// is intentionally non-null with a default so older read paths that
// SELECT * keep working. Idempotent via the detect-then-act pattern
// in addColumnIfNotExists (catches "duplicate column" benignly).
func (db *DB) migrateRepoIdentifiedTransport() {
	db.addColumnIfNotExists(constants.SQLAddRepoIdentifiedTransport)
}

// migrateRepoScanFolderID adds the nullable ScanFolderId FK column to Repo
// (v3.7.0, Phase 2.1). No backfill — existing rows stay NULL until the next
// `gitmap scan` re-discovers them.
func (db *DB) migrateRepoScanFolderID() {
	db.addColumnIfNotExists(constants.SQLAddRepoScanFolderId)
}

// migrateAliasColumns adds IsPrimary and Source columns to Alias.
func (db *DB) migrateAliasColumns() {
	db.addColumnIfNotExists(constants.SQLAddAliasIsPrimary)
	db.addColumnIfNotExists(constants.SQLAddAliasSource)
}

// migrateRepoInjectTimestamps adds LastInjectedDesktopAt + LastInjectedVSCodeAt
// to existing Repo rows (schema v25). Both default to ” so the inject /
// open idempotency check treats legacy rows as "never injected".
func (db *DB) migrateRepoInjectTimestamps() {
	db.addColumnIfNotExists(constants.SQLAddRepoLastInjectedDesktopAt)
	db.addColumnIfNotExists(constants.SQLAddRepoLastInjectedVSCodeAt)
}

// migrateRepoLastClonedAt adds LastClonedAt (schema v26). Powers the
// `gitmap release` auto-cd-into-most-recent-clone fallback. Defaults
// to ” so legacy rows are treated as "never cloned by gitmap".
func (db *DB) migrateRepoLastClonedAt() {
	db.addColumnIfNotExists(constants.SQLAddRepoLastClonedAt)
}

// migrateVSCodeProjectPaths adds the JSON-encoded Paths TEXT column to
// existing VSCodeProject tables (schema v20+, v3.39.0). No-op on fresh
// installs and on already-migrated databases (handled by addColumnIfNotExists).
func (db *DB) migrateVSCodeProjectPaths() {
	db.addColumnIfNotExists(constants.SQLAddVSCodeProjectPathsColumn)
}

// purgeOrphanRepoReferences removes rows in child tables referencing deleted or missing RepoIds.
func (db *DB) purgeOrphanRepoReferences() {
	orphanTables := []string{
		"DetectedProject",
		"Release",
		"VersionProbe",
		"Alias",
		"GroupRepo",
		"RepoVersionHistory",
	}

	for _, tbl := range orphanTables {
		if db.tableExists(tbl) {
			query := fmt.Sprintf("DELETE FROM %s WHERE RepoId NOT IN (SELECT RepoId FROM Repo)", tbl)
			_, _ = ExecWrapper(db.conn, query).Destruct()
		}
	}
}

// migrateZipGroupItemPaths adds RepoPath, RelativePath, FullPath columns
// to existing ZipGroupItems tables and copies Path into FullPath.
func (db *DB) migrateZipGroupItemPaths() {
	if !db.tableExists("ZipGroupItems") {
		return
	}

	db.addColumnIfNotExists(constants.SQLMigrateZGIRepoPath)
	db.addColumnIfNotExists(constants.SQLMigrateZGIRelativePath)
	db.addColumnIfNotExists(constants.SQLMigrateZGIFullPath)

	// Copy existing Path values into FullPath — data migration. Skip silently
	// when the legacy Path column does not exist (fresh installs).
	if _, err := db.conn.Exec(constants.SQLMigrateZGICopyPath); err != nil && !isBenignAlterError(err) {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not copy ZipGroupItem paths: %v\n", err)
	}
}

// migrateTRCommitSha renames the legacy Commit column to CommitSha in
// TempReleases. Uses detect-then-act via PRAGMA table_info so fresh installs
// (where only CommitSha exists) never trigger a SQLite warning, regardless
// of OS or driver-specific error message wording.
func (db *DB) migrateTRCommitSha() {
	if !db.tableExists("TempReleases") {
		return
	}

	if !db.columnExists("TempReleases", "Commit") {
		return // already migrated, or freshly created with CommitSha.
	}

	if db.columnExists("TempReleases", "CommitSha") {
		return // both columns exist — refuse to clobber.
	}

	if _, err := ExecWrapper(db.conn, constants.SQLMigrateTRCommitSha).Destruct(); err != nil && !isBenignAlterError(err) {
		logMigrationFailure("TempReleases", "Commit",
			"rename->CommitSha", err, constants.SQLMigrateTRCommitSha)
	}
}

// migratePendingTaskColumns adds WorkingDirectory and CommandArgs to existing tables.
func (db *DB) migratePendingTaskColumns() {
	if db.tableExists("PendingTask") {
		db.addColumnIfNotExists(constants.SQLMigratePendingWorkDir)
		db.addColumnIfNotExists(constants.SQLMigratePendingCmdArgs)
	}

	if db.tableExists("CompletedTask") {
		db.addColumnIfNotExists(constants.SQLMigrateCompletedWorkDir)
		db.addColumnIfNotExists(constants.SQLMigrateCompletedCmdArgs)
	}
}
