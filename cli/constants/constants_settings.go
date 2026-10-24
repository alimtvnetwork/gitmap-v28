package constants

// Settings table (v15: Setting singular). The PK is Key (TEXT) so no
// {Table}Id rename is needed.
const TableSetting = "Setting"

// Legacy plural retained for migration detection.
const LegacyTableSettings = "Settings"

// SQL: create Setting table (v15).
const SQLCreateSetting = `CREATE TABLE IF NOT EXISTS Setting (
	Key   TEXT PRIMARY KEY,
	Value TEXT NOT NULL
)`

// SQL: setting operations (v15).
const (
	SQLUpsertSetting = `INSERT INTO Setting (Key, Value) VALUES (?, ?)
		ON CONFLICT(Key) DO UPDATE SET Value=excluded.Value`

	SQLSelectSetting = "SELECT Value FROM Setting WHERE Key = ?"

	SQLDeleteSetting = "DELETE FROM Setting WHERE Key = ?"
)

// SQL: reset (v15 + legacy).
const (
	SQLDropSetting  = "DROP TABLE IF EXISTS Setting"
	SQLDropSettings = "DROP TABLE IF EXISTS Settings" // legacy
)

// Settings keys.
const (
	SettingActiveGroup             = "active_group"
	SettingActiveMultiGroup        = "active_multi_group"
	SettingSourceRepoPath          = "source_repo_path"
	SettingSchemaVersion           = "schema_version"
	SettingCustomRepoExclusions    = "custom_repo_exclusions"
	SettingPullActivityWindowHours = "pae_recency_window_hours"
	SettingNodeResourceThreshold   = "node_resource_threshold"
	SettingGoEmailReaderBin        = "go_email_reader_bin"
	SettingPullCacheWindowSeconds  = "pull_cache_window_seconds"
	SettingPullRetentionDays       = "pull_retention_days"
)

// Default retention period in days for git pull execution history.
const DefaultPullRetentionDays = 7

// Default activity window for efficient pull operations.
const DefaultPullActivityWindowHours = 48

// Default cache window in seconds for CWD git pull operations.
const DefaultPullCacheWindowSeconds = 3

// DefaultNodeResourceThreshold defines the baseline benchmark utilization limit (30%).
const DefaultNodeResourceThreshold = 30

// SchemaVersionCurrent is the target schema version produced by the current
// build of Migrate(). Bump this integer whenever a NEW migration step is
// added to gitmap/store/store.go:Migrate (a new v15 phase, a new ALTER, a
// new seed table, etc.).
//
// Migrate() short-circuits when Setting[schema_version] == this value, so
// every gitmap subcommand that calls openDB() pays only a single SELECT
// against Setting instead of re-running the full v15 phase pipeline.
//
// Bump policy:
//   - Bump on ANY structural change to Migrate() — new CREATE TABLE,
//     new ALTER TABLE, new v15 phase, new seed call, new ID rename.
//   - Do NOT bump for cosmetic changes (comments, log strings, code moves
//     that produce identical SQL).
//   - The marker is cleared by `gitmap db-reset` and by migrateLegacyIDs()
//     when it detects pre-integer-PK rows, so legacy databases will always
//     re-run the full pipeline regardless of this number.
// SchemaVersionCurrent advances to 34 to incorporate PascalCase WorkDirectory,
// FolderWorkDir, FolderGroup, and SshKeyAssociation tables and indexes.
const SchemaVersionCurrent = 34

// Schema-version log strings.
const (
	MsgSchemaVersionUpToDateFmt = "  ✓ Schema version %d is current — skipping migration pipeline.\n"
	MsgSchemaVersionAdvanceFmt  = "  ▸ Schema version %d → %d — running migration pipeline...\n"
	WarnSchemaVersionWriteFmt   = "  ⚠ Could not record schema version %d: %v\n"
)

// Settings error messages.
const (
	ErrDBSettingUpsert = "failed to save setting: %w"
	ErrDBSettingQuery  = "failed to read setting: %w"
)
