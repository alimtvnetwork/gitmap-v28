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
	SettingActiveGroup      = "active_group"
	SettingActiveMultiGroup = "active_multi_group"
	SettingSourceRepoPath   = "source_repo_path"
	SettingSchemaVersion    = "schema_version"
)

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
const SchemaVersionCurrent = 28

// Schema-version log strings.
const (
	MsgSchemaVersionUpToDateFmt = "  ✓ Schema version %d is current — skipping migration pipeline.\n"
	MsgSchemaVersionAdvanceFmt  = "  ▸ Schema version %d → %d — running migration pipeline...\n"
	WarnSchemaVersionWriteFmt   = "  ⚠ Could not record schema version %d: %v\n"
)

// Settings error messages.
const (
	ErrDBSettingUpsert = "failed to save setting: %v"
	ErrDBSettingQuery  = "failed to read setting: %v"
)
