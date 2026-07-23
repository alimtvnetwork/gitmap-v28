package constants

// Bookmark table name (v15: singular).
const TableBookmark = "Bookmark"

// Legacy plural retained for migration detection.
const LegacyTableBookmarks = "Bookmarks"

// SQL: create Bookmark table (v15: singular + BookmarkId PK).
const SQLCreateBookmark = `CREATE TABLE IF NOT EXISTS Bookmark (
	BookmarkId INTEGER PRIMARY KEY AUTOINCREMENT,
	Name       TEXT NOT NULL UNIQUE,
	Command    TEXT NOT NULL,
	Args       TEXT DEFAULT '',
	Flags      TEXT DEFAULT '',
	CreatedAt  TEXT DEFAULT CURRENT_TIMESTAMP
)`

// SQL: bookmark operations (v15).
const (
	SQLInsertBookmark = `INSERT INTO Bookmark (Name, Command, Args, Flags)
		VALUES (?, ?, ?, ?)`

	SQLSelectAllBookmarks = `SELECT BookmarkId, Name, Command, Args, Flags, CreatedAt
		FROM Bookmark ORDER BY Name`

	SQLSelectBookmarkByName = `SELECT BookmarkId, Name, Command, Args, Flags, CreatedAt
		FROM Bookmark WHERE Name = ?`

	SQLDeleteBookmark = "DELETE FROM Bookmark WHERE Name = ?"

	SQLDropBookmark  = "DROP TABLE IF EXISTS Bookmark"
	SQLDropBookmarks = "DROP TABLE IF EXISTS Bookmarks" // legacy
)

// SQL: import-side bookmark insert.
const SQLImportInsertBookmark = `INSERT OR IGNORE INTO Bookmark (Name, Command, Args, Flags) VALUES (?, ?, ?, ?)`

// gitmap:cmd top-level
// Bookmark CLI commands.
const (
	CmdBookmark      = "bookmark"
	CmdBookmarkAlias = "bk"
)

// gitmap:cmd top-level
// Bookmark subcommands.
const (
	CmdBookmarkSave   = "save"   // gitmap:cmd skip
	CmdBookmarkList   = "list"   // gitmap:cmd skip
	CmdBookmarkRun    = "run"    // gitmap:cmd skip
	CmdBookmarkDelete = "delete" // gitmap:cmd skip
)

// Bookmark help text.
const (
	HelpBookmark = "  bookmark (bk) <sub> Save and replay command+flag combinations (save, list, run, delete)"
)

// Bookmark messages.
const (
	MsgBookmarkSaved     = "Bookmark saved: %s → gitmap %s %s %s\n"
	MsgBookmarkDeleted   = "Bookmark deleted: %s\n"
	MsgBookmarkEmpty     = "No bookmarks saved.\n"
	MsgBookmarkRunning   = "Running bookmark: %s → gitmap %s %s %s\n"
	MsgBookmarkColumns   = "NAME                 COMMAND         ARGS             FLAGS"
	MsgBookmarkRowFmt    = "%-20s %-15s %-16s %s\n"
	ErrBookmarkUsage     = "usage: gitmap bookmark <save|list|run|delete> [args]\n"
	ErrBookmarkSaveUsage = "usage: gitmap bookmark save <name> <command> [args...] [--flags...]\n"
	ErrBookmarkRunUsage  = "usage: gitmap bookmark run <name>\n"
	ErrBookmarkDelUsage  = "usage: gitmap bookmark delete <name>\n"
	ErrBookmarkNotFound  = "bookmark not found: %s\n"
	ErrBookmarkExists    = "bookmark already exists: %s (delete it first)\n"
	ErrBookmarkQuery     = "failed to query bookmarks: %v"
	ErrBookmarkSave      = "failed to save bookmark: %v\n"
	ErrBookmarkDelete    = "failed to delete bookmark: %v\n"
)
