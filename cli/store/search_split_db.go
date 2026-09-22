// Package store — search_split_db.go: isolated SQLite database for search queries and audit logging.
package store

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

const (
	sqlCreateSearchCategory = `CREATE TABLE IF NOT EXISTS SearchCategory (
    SearchCategoryId INTEGER PRIMARY KEY AUTOINCREMENT,
    CategoryCode     TEXT NOT NULL UNIQUE,
    Name             TEXT NOT NULL,
    Description      TEXT NOT NULL DEFAULT '',
    IsActive         INTEGER NOT NULL DEFAULT 1,
    CreatedAt        TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS IdxSearchCategory_Code ON SearchCategory(CategoryCode);`

	sqlSeedSearchCategory = `INSERT OR IGNORE INTO SearchCategory (CategoryCode, Name, Description) VALUES
('user', 'User Search', 'Interactive user-initiated search'),
('ai', 'AI Search', 'AI-initiated search and retrieval'),
('automated_audit', 'Automated Audit', 'Automated scanner and audit search');`

	sqlCreateSearchRecord = `CREATE TABLE IF NOT EXISTS SearchRecord (
    SearchRecordId INTEGER PRIMARY KEY AUTOINCREMENT,
    CategoryCode   TEXT NOT NULL DEFAULT 'user',
    QueryText      TEXT NOT NULL DEFAULT '',
    RegexPattern   TEXT NOT NULL DEFAULT '',
    SearchType     TEXT NOT NULL DEFAULT 'keyword',
    CallerIp       TEXT NOT NULL DEFAULT '',
    IsAiCaller     INTEGER NOT NULL DEFAULT 0,
    DurationMs     INTEGER NOT NULL DEFAULT 0,
    ResultCount    INTEGER NOT NULL DEFAULT 0,
    CreatedAt      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS IdxSearchRecord_Category ON SearchRecord(CategoryCode);
CREATE INDEX IF NOT EXISTS IdxSearchRecord_IsAiCaller ON SearchRecord(IsAiCaller);
CREATE INDEX IF NOT EXISTS IdxSearchRecord_CreatedAt ON SearchRecord(CreatedAt);`

	sqlCreateSearchLogView = `CREATE VIEW IF NOT EXISTS SearchLogView AS
SELECT
    r.SearchRecordId,
    r.CategoryCode,
    COALESCE(c.Name, r.CategoryCode) AS CategoryName,
    r.QueryText,
    r.RegexPattern,
    r.SearchType,
    r.CallerIp,
    r.IsAiCaller,
    r.DurationMs,
    r.ResultCount,
    r.CreatedAt
FROM SearchRecord r
LEFT JOIN SearchCategory c ON r.CategoryCode = c.CategoryCode;`
)

// SearchSplitDB wraps an isolated SQLite database connection for search logging.
type SearchSplitDB struct {
	conn *sql.DB
	Path string
}

// SearchDbPath returns the full path to the search SQLite DB.
func SearchDbPath() string {
	return ResolveSearchDbPath("")
}

// OpenSearchSplitDB opens or initializes the search split database.
func OpenSearchSplitDB() (*SearchSplitDB, error) {
	return OpenSearchSplitDBAt(SearchDbPath())
}

// OpenSearchSplitDBAt opens or creates a search split database at a specific path.
func OpenSearchSplitDBAt(dbPath string) (*SearchSplitDB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, apperror.WrapSimple(err, "search_split.mkdir")
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "search_split.open")
	}

	return initSearchSplitConn(conn, dbPath)
}

func initSearchSplitConn(conn *sql.DB, dbPath string) (*SearchSplitDB, error) {
	if err := ConfigureSQLiteConn(conn); err != nil {
		_ = conn.Close()
		return nil, apperror.WrapSimple(err, "search_split.config")
	}

	if err := initSearchSchema(conn); err != nil {
		_ = conn.Close()
		return nil, err
	}

	db := &SearchSplitDB{conn: conn, Path: dbPath}
	db.registerWithRegistry()
	return db, nil
}

func initSearchSchema(conn *sql.DB) error {
	statements := []string{
		sqlCreateSearchCategory,
		sqlSeedSearchCategory,
		sqlCreateSearchRecord,
		sqlCreateSearchLogView,
	}
	for _, stmt := range statements {
		res := ExecWrapper(conn, stmt)
		if res.IsFailure {
			return apperror.WrapSimple(res.Error, "search_split.initSchema")
		}
	}
	return nil
}

func (db *SearchSplitDB) registerWithRegistry() {
	master, err := OpenDefault()
	if err != nil {
		return
	}
	defer master.Close()
	_ = master.RegisterSplitDB(db.buildRegistryEntry())
}

func (db *SearchSplitDB) buildRegistryEntry() SplitDatabaseEntry {
	return SplitDatabaseEntry{
		DatabaseType:  "search",
		DatabaseKey:   "search_master",
		DatabasePath:  db.Path,
		Status:        "active",
		IsActive:      true,
		Description:   "Search queries, caller tracking, and audit split database",
		SchemaVersion: 1,
	}
}

// Close terminates the database connection.
func (db *SearchSplitDB) Close() error {
	if db == nil || db.conn == nil {
		return nil
	}
	return db.conn.Close()
}

// Conn returns the raw database connection.
func (db *SearchSplitDB) Conn() *sql.DB {
	return db.conn
}

// RecordSearchQuery logs a search query into the database instance.
func (db *SearchSplitDB) RecordSearchQuery(
	categoryCode, queryText, regexPattern, searchType, callerIp string,
	isAiCaller bool,
	durationMs, resultCount int,
) error {
	cat := resolveSearchCategory(categoryCode, isAiCaller)
	aiVal := boolToInt(isAiCaller)
	q := `INSERT INTO SearchRecord (
		CategoryCode, QueryText, RegexPattern, SearchType,
		CallerIp, IsAiCaller, DurationMs, ResultCount
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	res := ExecWrapper(db.conn, q, cat, queryText, regexPattern, searchType, callerIp, aiVal, durationMs, resultCount)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "search_split.record_query")
	}
	return nil
}

// RecordSearchQuery opens the default search database and logs a query.
func RecordSearchQuery(
	categoryCode, queryText, regexPattern, searchType, callerIp string,
	isAiCaller bool,
	durationMs, resultCount int,
) error {
	db, err := OpenSearchSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()
	return db.RecordSearchQuery(categoryCode, queryText, regexPattern, searchType, callerIp, isAiCaller, durationMs, resultCount)
}

func resolveSearchCategory(categoryCode string, isAiCaller bool) string {
	if categoryCode != "" {
		return categoryCode
	}
	if isAiCaller {
		return "ai"
	}
	return "user"
}
