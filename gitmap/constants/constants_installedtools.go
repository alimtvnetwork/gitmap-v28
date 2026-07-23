package constants

// SQL: create InstalledTool table (v15: singular + InstalledToolId PK).
const SQLCreateInstalledTool = `CREATE TABLE IF NOT EXISTS InstalledTool (
	InstalledToolId INTEGER PRIMARY KEY AUTOINCREMENT,
	Tool            TEXT NOT NULL UNIQUE,
	VersionMajor    INTEGER NOT NULL DEFAULT 0,
	VersionMinor    INTEGER NOT NULL DEFAULT 0,
	VersionPatch    INTEGER NOT NULL DEFAULT 0,
	VersionBuild    INTEGER NOT NULL DEFAULT 0,
	VersionString   TEXT NOT NULL DEFAULT '',
	PackageManager  TEXT NOT NULL DEFAULT '',
	InstallPath     TEXT NOT NULL DEFAULT '',
	InstalledAt     TEXT NOT NULL DEFAULT '',
	UpdatedAt       TEXT NOT NULL DEFAULT ''
)`

// Legacy plural retained for migration detection.
const LegacyTableInstalledTools = "InstalledTools"

// SQL: InstalledTool queries (v15).
const (
	SQLInsertInstalledTool = `INSERT OR REPLACE INTO InstalledTool
		(Tool, VersionMajor, VersionMinor, VersionPatch, VersionBuild, VersionString, PackageManager, InstallPath, InstalledAt, UpdatedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))`

	SQLSelectInstalledTool = `SELECT InstalledToolId, Tool, VersionMajor, VersionMinor, VersionPatch, VersionBuild, VersionString, PackageManager, InstallPath, InstalledAt, UpdatedAt FROM InstalledTool WHERE Tool = ?`
	SQLSelectAllInstalled  = `SELECT InstalledToolId, Tool, VersionMajor, VersionMinor, VersionPatch, VersionBuild, VersionString, PackageManager, InstallPath, InstalledAt, UpdatedAt FROM InstalledTool ORDER BY Tool`
	SQLDeleteInstalledTool = `DELETE FROM InstalledTool WHERE Tool = ?`
	SQLExistsInstalledTool = `SELECT COUNT(*) FROM InstalledTool WHERE Tool = ?`
	SQLDropInstalledTool   = `DROP TABLE IF EXISTS InstalledTool`
	SQLDropInstalledTools  = `DROP TABLE IF EXISTS InstalledTools` // legacy
)
