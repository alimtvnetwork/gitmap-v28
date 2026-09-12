package constants

// gitmap:cmd top-level
// Zip group command names.
const (
	CmdZipGroup      = "zip-group"
	CmdZipGroupShort = "z"
	SubCmdZGCreate   = "create"
	SubCmdZGAdd      = "add"
	SubCmdZGRemove   = "remove"
	SubCmdZGList     = "list"
	SubCmdZGShow     = "show"
	SubCmdZGDelete   = "delete"
	SubCmdZGRename   = "rename"
)

// Zip group table names (v15: PascalCase singular + {Table}Id PK).
const (
	TableZipGroup     = "ZipGroup"
	TableZipGroupItem = "ZipGroupItem"
)

// Legacy plural names retained for migration detection.
const (
	LegacyTableZipGroups     = "ZipGroups"
	LegacyTableZipGroupItems = "ZipGroupItems"
)

// SQL: create ZipGroup table (v15: singular + ZipGroupId PK).
const SQLCreateZipGroup = `CREATE TABLE IF NOT EXISTS ZipGroup (
	ZipGroupId  INTEGER PRIMARY KEY AUTOINCREMENT,
	Name        TEXT NOT NULL UNIQUE,
	ArchiveName TEXT DEFAULT '',
	CreatedAt   TEXT DEFAULT CURRENT_TIMESTAMP
)`

// SQL: create ZipGroupItem table (v15 singular). Composite PK retained.
const SQLCreateZipGroupItem = `CREATE TABLE IF NOT EXISTS ZipGroupItem (
	ZipGroupId   INTEGER NOT NULL REFERENCES ZipGroup(ZipGroupId) ON DELETE CASCADE,
	RepoPath     TEXT NOT NULL DEFAULT '',
	RelativePath TEXT NOT NULL DEFAULT '',
	FullPath     TEXT NOT NULL DEFAULT '',
	IsFolder     INTEGER DEFAULT 0,
	PRIMARY KEY (ZipGroupId, FullPath)
)`

// SQL: legacy ALTERs for pre-v15 ZipGroupItems (still target legacy plural —
// run BEFORE v15 rebuild copies the table). Idempotent.
const (
	SQLMigrateZGIRepoPath     = `ALTER TABLE ZipGroupItems ADD COLUMN RepoPath TEXT NOT NULL DEFAULT ''`
	SQLMigrateZGIRelativePath = `ALTER TABLE ZipGroupItems ADD COLUMN RelativePath TEXT NOT NULL DEFAULT ''`
	SQLMigrateZGIFullPath     = `ALTER TABLE ZipGroupItems ADD COLUMN FullPath TEXT NOT NULL DEFAULT ''`
	SQLMigrateZGICopyPath     = `UPDATE ZipGroupItems SET FullPath = Path WHERE FullPath = '' AND Path IS NOT NULL AND Path != ''`
	SQLMigrateZGIDropPath     = `ALTER TABLE ZipGroupItems DROP COLUMN Path`
)

// SQL: zip group operations (v15 singular tables + ZipGroupId PK).
const (
	SQLInsertZipGroup = `INSERT INTO ZipGroup (Name, ArchiveName) VALUES (?, ?)`

	SQLSelectAllZipGroups = `SELECT ZipGroupId, Name, ArchiveName, CreatedAt FROM ZipGroup ORDER BY Name`

	SQLSelectZipGroupByName = `SELECT ZipGroupId, Name, ArchiveName, CreatedAt FROM ZipGroup WHERE Name = ?`

	SQLDeleteZipGroup = `DELETE FROM ZipGroup WHERE Name = ?`

	SQLUpdateZipGroupArchive = `UPDATE ZipGroup SET ArchiveName = ? WHERE Name = ?`
)

// SQL: zip group item operations.
const (
	SQLInsertZipGroupItem = `INSERT OR IGNORE INTO ZipGroupItem (ZipGroupId, RepoPath, RelativePath, FullPath, IsFolder) VALUES (?, ?, ?, ?, ?)`

	SQLDeleteZipGroupItem = `DELETE FROM ZipGroupItem WHERE ZipGroupId = ? AND FullPath = ?`

	SQLSelectZipGroupItems = `SELECT ZipGroupId, RepoPath, RelativePath, FullPath, IsFolder FROM ZipGroupItem WHERE ZipGroupId = ? ORDER BY FullPath`

	SQLCountZipGroupItems = `SELECT COUNT(*) FROM ZipGroupItem WHERE ZipGroupId = ?`

	SQLSelectAllZipGroupsWithCount = `SELECT g.ZipGroupId, g.Name, g.ArchiveName, g.CreatedAt,
		(SELECT COUNT(*) FROM ZipGroupItem i WHERE i.ZipGroupId = g.ZipGroupId) AS ItemCount
		FROM ZipGroup g ORDER BY g.Name`
)

// SQL: drop zip group tables (v15 + legacy plurals retained for Reset).
const (
	SQLDropZipGroup      = "DROP TABLE IF EXISTS ZipGroup"
	SQLDropZipGroups     = "DROP TABLE IF EXISTS ZipGroups" // legacy
	SQLDropZipGroupItem  = "DROP TABLE IF EXISTS ZipGroupItem"
	SQLDropZipGroupItems = "DROP TABLE IF EXISTS ZipGroupItems" // legacy
)

// Zip group flag descriptions.
const (
	FlagDescZGArchive  = "Custom output archive filename"
	FlagDescZGZipGroup = "Include a persistent zip group as a release asset"
	FlagDescZGZipItem  = "Add ad-hoc file or folder to zip as a release asset"
	FlagDescZGBundle   = "Bundle all -Z items into a single named archive"
)

// Zip group JSON persistence directory/file.
const (
	ZGJSONDir  = ".gitmap"
	ZGJSONFile = "zip-groups.json"
)

// Zip group messages.
const (
	MsgZGCreated      = "  ✓ Created zip group %q\n"
	MsgZGCreatedPath  = "  ✓ Created zip group %q with %s %s\n"
	MsgZGDeleted      = "  ✓ Deleted zip group %q\n"
	MsgZGItemAdded    = "  ✓ Added %s to %q (%s)\n"
	MsgZGItemRemoved  = "  ✓ Removed %s from %q\n"
	MsgZGArchiveSet   = "  ✓ Archive name set to %q for group %q\n"
	MsgZGListHeader   = "\n  Zip Groups (%d):\n\n"
	MsgZGListRow      = "  %-20s %3d item(s)  %s\n"
	MsgZGShowHeader   = "\n  %s (%d item(s)):\n\n"
	MsgZGShowFile     = "    📄 %s\n"
	MsgZGShowFolder   = "    📁 %s\n"
	MsgZGShowArchive  = "  Archive: %s\n"
	MsgZGShowPaths    = "    repo:     %s\n    relative: %s\n    full:     %s\n"
	MsgZGCompressed   = "  ✓ Compressed %s → %s\n"
	MsgZGDryRunHeader = "  [dry-run] Would create %d zip archive(s):\n"
	MsgZGDryRunEntry  = "    → %s (%d items: %s)\n"
	MsgZGSkipEmpty    = "  ⚠ Skipping empty group %q\n"
	MsgZGSkipMissing  = "  ⚠ Skipping missing item: %s\n"
	MsgZGProcessing   = "  Processing %d zip group(s)...\n"
	MsgZGNoArchives   = "  ⚠ No zip archives were produced from %d group(s)\n"
	ErrZGStagingDir   = "  ✗ Cannot create staging dir at %s: %v (operation: mkdir)\n"
	MsgZGTypeFolder   = "folder"
	MsgZGTypeFile     = "file"
	MsgZGJSONWritten  = "  ✓ Saved %s\n"
	MsgZGShowExpanded = "    Contents (%d files):\n"
	MsgZGShowExpFile  = "      %s\n"
)

// Zip group error messages.
const (
	ErrZGNotFound    = "no zip group found: %s"
	ErrZGEmpty       = "zip group name cannot be empty"
	ErrZGCreate      = "failed to create zip group: %v"
	ErrZGQuery       = "failed to query zip groups: %v"
	ErrZGDelete      = "failed to delete zip group: %v"
	ErrZGAddItem     = "failed to add item to zip group: %v"
	ErrZGRemoveItem  = "failed to remove item from zip group: %v"
	ErrZGCompress    = "  ✗ Failed to create archive for %s: %v (operation: write)\n"
	ErrZGGroupNotDB  = "zip group %q not found in database"
	ErrZGPathResolve = "cannot resolve path %q: %v (operation: resolve)"
	ErrZGJSONWrite   = "failed to write zip-groups.json at %s: %v (operation: write)"
)
