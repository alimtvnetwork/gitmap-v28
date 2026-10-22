// Package store manages the SQLite database for gitmap.
package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"

	_ "modernc.org/sqlite"
)

// DB wraps the SQLite database connection.
type DB struct {
	conn  *sql.DB
	path  string
	dbDir string
	isMem bool
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

// OpenInMemory opens an in-memory SQLite database instance for testing.
func OpenInMemory() (*DB, error) {
	connPath := fmt.Sprintf("file:mem_%d?mode=memory&cache=shared", time.Now().UnixNano())
	retainMemAnchor(connPath)

	return openDBConnection(connPath, ":memory:", "", true)
}

var (
	memAnchorMu sync.Mutex
	memAnchors  = make(map[string]*sql.DB)
)

func retainMemAnchor(connPath string) {
	memAnchorMu.Lock()
	defer memAnchorMu.Unlock()

	if _, hasAnchor := memAnchors[connPath]; hasAnchor {
		return
	}

	anchor, err := sql.Open("sqlite", connPath)
	if err != nil {
		return
	}

	_ = ConfigureSQLiteConn(anchor)
	memAnchors[connPath] = anchor
}

func resolveConnPath(dbPath string, isMem bool) string {
	if isMem {
		norm := filepath.ToSlash(filepath.Clean(dbPath))
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(norm)))[:16]
		connPath := "file:" + hash + "?mode=memory&cache=shared"
		retainMemAnchor(connPath)

		return connPath
	}

	return dbPath
}

// openDBAt opens a database at an exact path.
func openDBAt(dbPath string) (*DB, error) {
	dbDir := filepath.Dir(dbPath)
	isMem := os.Getenv("GITMAP_IN_MEMORY_DB") == "1"

	if err := ensureDir(dbDir); err != nil {
		return nil, fmt.Errorf(constants.ErrDBCreateDir, dbDir, err)
	}

	if err := lockDBIfNotMem(dbDir, isMem); err != nil {
		return nil, err
	}

	connPath := resolveConnPath(dbPath, isMem)

	return openDBConnection(connPath, dbPath, dbDir, isMem)
}

func closeOnConfigError(conn *sql.DB, dbDir string, isMem bool, err error) (*DB, error) {
	conn.Close()
	releaseLockIfNotMem(dbDir, isMem)

	return nil, err
}

func openDBConnection(connPath, dbPath, dbDir string, isMem bool) (*DB, error) {
	conn, err := sql.Open("sqlite", connPath)
	if err != nil {
		releaseLockIfNotMem(dbDir, isMem)

		return nil, fmt.Errorf(constants.ErrDBOpen, dbPath, err)
	}

	if err := ConfigureSQLiteConn(conn); err != nil {
		return closeOnConfigError(conn, dbDir, isMem, err)
	}

	return &DB{conn: conn, path: dbPath, dbDir: dbDir, isMem: isMem}, nil
}

// GetDatabasePath returns the filesystem path of the active SQLite database file.
func (db *DB) GetDatabasePath() string {
	if db == nil {
		return ""
	}

	return db.path
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
		getRepoAbsPathIndexDDL(),
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
		getScanFolderPathIndexQuery(),
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
		constants.SQLCreateClusterNode,
		constants.SQLCreateClusterRun,
		constants.SQLCreateClusterExecResult,
		SQLCreateWorkDirsTable,
		SQLCreateSSHHostsTable,
		SQLCreateSSHHistoryTable,
		constants.SQLCreateSplitDatabaseRegistry,
		constants.SQLCreateSplitDatabaseRegistryTypeIndex,
		constants.SQLCreateSplitDatabaseRegistryStatusIndex,
		constants.SQLCreateSplitDatabaseRegistryUpdatedAtIndex,
		SQLCreateReleaseCache,
	}

	for _, stmt := range statements {
		if _, err := ExecWrapper(db.conn, stmt).Destruct(); err != nil {
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
	db.migrateAliasColumns()
	db.purgeOrphanRepoReferences()

	if err := db.EnsurePurgeHistoryTable(); err != nil {
		return fmt.Errorf("ensure purge history table: %w", err)
	}

	if err := db.SeedProjectTypes(); err != nil {
		return err
	}

	if err := db.SeedTaskTypes(); err != nil {
		return err
	}

	// commit-in (02-spec/03-commit-in/) — adds 18 tables + seeds. Keep
	// AFTER the v15 / standard CREATE pass so commit-in DDL never
	// races with the legacy rebuilds; keep BEFORE the schema-version
	// stamp so a failure forces the next run to retry.
	if err := db.migrateCommitIn(); err != nil {
		return err
	}

	if err := db.SyncKnownSplitDatabases(); err != nil {
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
		"DROP TABLE IF EXISTS PurgeHistoryLog",
	}

	for _, stmt := range drops {
		if _, err := ExecWrapper(db.conn, stmt).Destruct(); err != nil {
			return fmt.Errorf(constants.ErrDBMigrate, err)
		}
	}

	return db.Migrate()
}

// Close shuts down the database connection and releases the OS directory lock.
func (db *DB) Close() error {
	if !db.isMem {
		releaseLock(db.dbDir)
	}

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

func (db *DB) SQL() *sql.DB {
	return db.conn
}

func (db *DB) Context() context.Context {
	return context.Background()
}

func releaseLockIfNotMem(dbDir string, isMem bool) {
	if !isMem {
		releaseLock(dbDir)
	}
}

func lockDBIfNotMem(dbDir string, isMem bool) error {
	if isMem {
		return nil
	}

	return acquireLock(dbDir)
}

func getRepoAbsPathIndexDDL() string {
	if runtime.GOOS == "windows" {
		return constants.SQLCreateAbsPathIndexWindows
	}

	return constants.SQLCreateAbsPathIndexUnix
}

func getScanFolderPathIndexQuery() string {
	if runtime.GOOS == "windows" {
		return constants.SQLCreateScanFolderPathIndexWindows
	}

	return constants.SQLCreateScanFolderPathIndexUnix
}

// GetCachedRelease retrieves a valid unexpired cached release from the database.
func (db *DB) GetCachedRelease(tag, platform, arch string) (*ReleaseCacheRecord, error) {
	return GetCachedRelease(db.conn, tag, platform, arch)
}

// UpsertReleaseCache inserts or updates a release cache entry with 24-hour expiration.
func (db *DB) UpsertReleaseCache(record ReleaseCacheRecord) error {
	return UpsertReleaseCache(db.conn, record)
}

// ListCachedReleases lists recent release cache records.
func (db *DB) ListCachedReleases(limit int) ([]ReleaseCacheRecord, error) {
	return ListCachedReleases(db.conn, limit)
}
