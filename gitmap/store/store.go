// Package store manages the SQLite database for gitmap.
package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"

	_ "modernc.org/sqlite"
)

// DB wraps the SQLite database connection.
type DB struct {
	conn  *sql.DB
	dbDir string
}

// Open creates or opens the SQLite database for the active profile.
func Open(outputDir string) (*DB, error) {
	dbFile := ActiveProfileDBFile(outputDir)
	dbPath := filepath.Join(outputDir, constants.DBDir, dbFile)

	return openDBAt(dbPath)
}

// OpenProfile opens the database for a specific named profile.
func OpenProfile(outputDir, profileName string) (*DB, error) {
	dbFile := ProfileDBFile(profileName)
	dbPath := filepath.Join(outputDir, constants.DBDir, dbFile)

	return openDBAt(dbPath)
}

// OpenAt opens (or creates) a database at an exact filesystem path.
// Used by features (e.g. commit-in) that anchor SQLite under a path
// other than the default <outputDir>/data layout.
func OpenAt(dbPath string) (*DB, error) {
	return openDBAt(dbPath)
}

// openDBAt opens a database at an exact path.
func openDBAt(dbPath string) (*DB, error) {
	dbDir := filepath.Dir(dbPath)
	if err := ensureDir(dbDir); err != nil {
		return nil, fmt.Errorf(constants.ErrDBCreateDir, dbDir, err)
	}

	if err := acquireLock(dbDir); err != nil {
		return nil, err
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		releaseLock(dbDir)

		return nil, fmt.Errorf(constants.ErrDBOpen, dbPath, err)
	}

	// SQLite does not support concurrent writes; pin to one connection
	// so PRAGMAs (foreign_keys, etc.) persist across all operations.
	conn.SetMaxOpenConns(1)

	err = enableFK(conn)
	if err != nil {
		conn.Close()
		releaseLock(dbDir)

		return nil, err
	}

	return &DB{conn: conn, dbDir: dbDir}, nil
}

// Migrate creates all required tables if they don't exist.
//
// Order: legacy UUID migration → v15 Repo rename → v15 Phase 1.2 (Group/
// Release/Alias/Bookmark) → v15 Phase 1.3 (Amendment/CommitTemplate/Setting/
// SshKey/InstalledTool/TempRelease) → standard CREATE TABLE pass → ALTER-based
// column additions → seed data. Every v15 step is idempotent and a no-op on
// fresh installs.
func (db *DB) Migrate() error {
	// Fast path: skip the entire pipeline when the persisted schema
	// version marker matches the current target. This makes every
	// gitmap subcommand that calls openDB() pay only one Setting SELECT
	// instead of the full v15 phase walk. See migrate_schemaversion.go
	// for the marker semantics and bump policy.
	if db.isSchemaUpToDate() {
		// Schema is current, but the downloader seed file may have
		// changed since the last run. The seeder is cheap (one Setting
		// SELECT for the hash) and is the only way config tweaks ship
		// to existing installs without a schema bump.
		db.SeedDownloaderConfig(constants.DefaultDownloaderConfigSeedPath)

		return nil
	}

	db.migrateLegacyIDs()

	if err := db.migrateV15Repo(); err != nil {
		return fmt.Errorf(constants.ErrV15RepoMigration, err)
	}

	// Pre-Phase-1.2: legacy `Releases` may be missing `Source` and/or `Notes`
	// columns on very old installs. Phase 1.2 SELECTs every column by name
	// when copying into the new `Release`, so add them here first. Idempotent.
	db.preV15Phase2EnsureReleaseColumns()

	if err := db.migrateV15Phase2(); err != nil {
		return fmt.Errorf(constants.ErrV15Phase2Migration, err)
	}

	// Phase 1.3 reads the legacy Commit column on TempReleases if it exists,
	// so rename that column BEFORE the v15 rebuild copies the table.
	db.migrateTRCommitSha()

	if err := db.migrateV15Phase3(); err != nil {
		return fmt.Errorf(constants.ErrV15Phase3Migration, err)
	}

	// Phase 1.4 reads legacy ZipGroupItems columns + pre-Csharp table names,
	// so run the column-shape migrations BEFORE the rebuild copies tables.
	db.migrateZipGroupItemPaths()
	db.migratePendingTaskColumns()

	if err := db.migrateV15Phase4(); err != nil {
		return fmt.Errorf(constants.ErrV15Phase4Migration, err)
	}

	if err := db.migrateV15Phase5(); err != nil {
		return fmt.Errorf(constants.ErrV15Phase5Migration, err)
	}

	if err := db.migrateV15Phase6(); err != nil {
		return fmt.Errorf(constants.ErrV15Phase6Migration, err)
	}

	statements := []string{
		constants.SQLCreateRepo,
		constants.SQLCreateAbsPathIndex,
		constants.SQLCreateGroup,
		constants.SQLCreateGroupRepo,
		constants.SQLCreateRelease,
		constants.SQLCreateReleaseRepoIdIndex,
		constants.SQLCreateCommitTemplate,
		constants.SQLCreateAmendment,
		constants.SQLCreateCommandHistory,
		constants.SQLCreateBookmark,
		constants.SQLCreateProjectType,
		constants.SQLCreateDetectedProject,
		constants.SQLCreateGoProjectMetadata,
		constants.SQLCreateGoRunnableFile,
		constants.SQLCreateCsharpProjectMeta,
		constants.SQLCreateCsharpProjectFile,
		constants.SQLCreateCsharpKeyFile,
		constants.SQLCreateSetting,
		constants.SQLCreateAlias,
		constants.SQLCreateZipGroup,
		constants.SQLCreateZipGroupItem,
		constants.SQLCreateSshKey,
		constants.SQLCreateTempRelease,
		constants.SQLCreateInstalledTool,
		constants.SQLCreateTaskType,
		constants.SQLCreatePendingTask,
		constants.SQLCreateCompletedTask,
		constants.SQLCreateRepoVersionHistory,
		constants.SQLCreateScanFolder,
		constants.SQLCreateScanFolderPathIndex,
		constants.SQLCreateVersionProbe,
		constants.SQLCreateVersionProbeRepoIndex,
		constants.SQLCreateVSCodeProject,
		constants.SQLCreateVSCodeProjectRootPathIndex,
		constants.SQLCreateArchiveHistory,
		constants.SQLCreateCloneInteractiveSelection,
		constants.SQLCreateClonePickRepoCanonIndex,
		constants.SQLCreateClonePickNameIndex,
		constants.SQLCreateTransaction,
		constants.SQLCreateTransactionFile,
		constants.SQLCreateTransactionStatusIndex,
		constants.SQLCreateTransactionAction,
		constants.SQLCreateTransactionActionTxnSeqIndex,
		constants.SQLCreateTransactionActionKindIndex,
		constants.SQLCreateMakeAllVisibilityRun,
		constants.SQLCreateMakeAllVisibilityResult,
		constants.SQLCreateMakeAllVisibilityResultRunIndex,
		constants.SQLCreateOwnerRepoListCache,
		constants.SQLCreateOwnerRepoNameIndex,
	}

	for _, stmt := range statements {
		if _, err := db.conn.Exec(stmt); err != nil {
			return fmt.Errorf(constants.ErrDBMigrate, err)
		}
	}

	db.migrateSourceColumn()
	db.migrateNotesColumn()
	db.migrateRepoVersionColumns()
	db.migrateRepoScanFolderID()
	db.migrateRepoInjectTimestamps()
	db.migrateRepoLastClonedAt()
	db.migrateRepoIdentifiedTransport()
	db.migrateVSCodeProjectPaths()

	if err := db.SeedProjectTypes(); err != nil {
		return err
	}

	if err := db.SeedTaskTypes(); err != nil {
		return err
	}

	// commit-in (spec/03-commit-in/) — adds 18 tables + seeds. Keep
	// AFTER the v15 / standard CREATE pass so commit-in DDL never
	// races with the legacy rebuilds; keep BEFORE the schema-version
	// stamp so a failure forces the next run to retry.
	if err := db.migrateCommitIn(); err != nil {
		return err
	}

	// Stamp the marker LAST so any earlier failure leaves the previous
	// (or empty) marker in place, ensuring the next run retries.
	db.writeSchemaVersion(constants.SchemaVersionCurrent)

	// Apply the downloader Seedable-Config AFTER the schema marker is
	// stamped. It writes only Setting rows (no schema change) and is
	// safe to repeat — the seeder hashes the seed file and skips when
	// nothing has changed since the last apply.
	db.SeedDownloaderConfig(constants.DefaultDownloaderConfigSeedPath)

	return nil
}

// addColumnIfNotExists runs an ALTER TABLE ADD COLUMN statement.
// It silently ignores "duplicate column" errors (expected on repeat runs)
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

// migrateZipGroupItemPaths adds RepoPath, RelativePath, FullPath columns
// to existing ZipGroupItems tables and copies Path into FullPath.
func (db *DB) migrateZipGroupItemPaths() {
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

	if _, err := db.conn.Exec(constants.SQLMigrateTRCommitSha); err != nil && !isBenignAlterError(err) {
		logMigrationFailure("TempReleases", "Commit",
			"rename->CommitSha", err, constants.SQLMigrateTRCommitSha)
	}
}

// migratePendingTaskColumns adds WorkingDirectory and CommandArgs to existing tables.
func (db *DB) migratePendingTaskColumns() {
	db.addColumnIfNotExists(constants.SQLMigratePendingWorkDir)
	db.addColumnIfNotExists(constants.SQLMigratePendingCmdArgs)
	db.addColumnIfNotExists(constants.SQLMigrateCompletedWorkDir)
	db.addColumnIfNotExists(constants.SQLMigrateCompletedCmdArgs)
}

// Reset drops all tables and recreates them for a fresh start. Lists v15
// singular drops first, followed by legacy plural drops (which are safe
// no-ops when the table does not exist) so installations at any migration
// state can be reset cleanly.
func (db *DB) Reset() error {
	drops := []string{
		// Children first (FK order). v15 names + legacy names retained as no-ops.
		constants.SQLDropCompletedTask,
		constants.SQLDropPendingTask,
		constants.SQLDropTaskType,
		constants.SQLDropSetting,
		constants.SQLDropSettings, // legacy
		constants.SQLDropGoRunnableFile,
		constants.SQLDropGoRunnableFiles, // legacy
		constants.SQLDropGoProjectMetadata,
		constants.SQLDropCsharpKeyFile,
		constants.SQLDropCsharpKeyFiles, // legacy (pre-Csharp + plural)
		constants.SQLDropCsharpProjectFile,
		constants.SQLDropCsharpProjectFiles, // legacy
		constants.SQLDropCsharpProjectMeta,
		constants.SQLDropCsharpProjectMetaLegacy, // legacy (pre-Csharp spelling)
		constants.SQLDropDetectedProject,
		constants.SQLDropDetectedProjects, // legacy
		constants.SQLDropProjectType,
		constants.SQLDropProjectTypes, // legacy
		constants.SQLDropGroupRepo,
		constants.SQLDropGroupRepos, // legacy
		constants.SQLDropGroup,
		constants.SQLDropGroups, // legacy
		constants.SQLDropRelease,
		constants.SQLDropReleases, // legacy
		constants.SQLDropAmendment,
		constants.SQLDropAmendments, // legacy
		constants.SQLDropCommitTemplate,
		constants.SQLDropCommitTemplates, // legacy
		constants.SQLDropCommandHistory,
		constants.SQLDropBookmark,
		constants.SQLDropBookmarks, // legacy
		constants.SQLDropAlias,
		constants.SQLDropAliases, // legacy
		constants.SQLDropZipGroupItem,
		constants.SQLDropZipGroupItems, // legacy
		constants.SQLDropZipGroup,
		constants.SQLDropZipGroups, // legacy
		constants.SQLDropTempRelease,
		constants.SQLDropTempReleases, // legacy
		constants.SQLDropSshKey,
		constants.SQLDropSSHKeys, // legacy
		constants.SQLDropInstalledTool,
		constants.SQLDropInstalledTools, // legacy
		constants.SQLDropRepoVersionHistory,
		constants.SQLDropVersionProbe,
		constants.SQLDropVSCodeProject,
		constants.SQLDropScanFolder,
		constants.SQLDropRepo,
		constants.SQLDropRepos, // legacy
		constants.SQLDropArchiveHistory,
	}

	for _, stmt := range drops {
		if _, err := db.conn.Exec(stmt); err != nil {
			return fmt.Errorf(constants.ErrDBMigrate, err)
		}
	}

	return db.Migrate()
}

// Close closes the database connection and releases the lock.
func (db *DB) Close() error {
	releaseLock(db.dbDir)

	return db.conn.Close()
}

// Conn returns the underlying sql.DB for advanced queries.
func (db *DB) Conn() *sql.DB {
	return db.conn
}

// ensureDir creates the directory tree if it doesn't exist.
func ensureDir(dir string) error {
	return os.MkdirAll(dir, constants.DirPermission)
}

// enableFK turns on SQLite foreign key enforcement and applies the
// WAL + relaxed-sync pragmas. WAL + synchronous=NORMAL eliminates the
// per-commit FlushFileBuffers stall that caused 10-min Windows CI
// timeouts in store/ and cmd/ test packages (see panic stacks in
// _winSync → FlushFileBuffers). busy_timeout gives concurrent test
// connections a 5s grace window instead of immediate SQLITE_BUSY.
func enableFK(conn *sql.DB) error {
	for _, pragma := range []string{
		constants.SQLPragmaJournalWAL,
		constants.SQLPragmaSynchronousNor,
		constants.SQLPragmaBusyTimeout5s,
		constants.SQLEnableFK,
	} {
		if _, err := conn.Exec(pragma); err != nil {
			return fmt.Errorf("apply pragma %q: %w", pragma, err)
		}
	}

	return nil
}
