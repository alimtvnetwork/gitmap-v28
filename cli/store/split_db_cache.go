package store

import (
	"database/sql"
	"os"
	"path/filepath"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

// CacheFileRecord represents metadata of a cached file in the root split-db.
type CacheFileRecord struct {
	FileId       int64
	RepoUrl      string
	RelativePath string
	AbsolutePath string
	FileSize     int64
	ModifiedTime int64
	FolderSlug   string
	IsKeep       bool
}

// OpenRootCacheDB opens or creates the root split database for cache metadata.
func OpenRootCacheDB(repoRoot string) (*sql.DB, *appfault.AppError) {
	baseDir := resolveBaseDataDir(repoRoot)
	dir := filepath.Join(baseDir, "cache")
	_ = os.MkdirAll(dir, 0755)
	dbPath := filepath.ToSlash(filepath.Join(dir, DbFileName))

	db, err := sql.Open("sqlite", dbPath+"?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, appfault.WrapSimple(err, "open root cache db")
	}
	if initErr := InitRootCacheSchema(db); initErr != nil {
		db.Close()
		return nil, initErr
	}
	return db, nil
}

// OpenSlugCacheDB opens or creates a partition split database for content lines.
func OpenSlugCacheDB(slug, repoRoot string) (*sql.DB, *appfault.AppError) {
	cleanSlug := SanitizeSlug(slug)
	baseDir := resolveBaseDataDir(repoRoot)
	dir := filepath.Join(baseDir, "cache", "slugs")
	_ = os.MkdirAll(dir, 0755)
	dbPath := filepath.ToSlash(filepath.Join(dir, cleanSlug+".db"))

	db, err := sql.Open("sqlite", dbPath+"?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, appfault.WrapSimple(err, "open slug cache db")
	}
	if initErr := InitSlugCacheSchema(db); initErr != nil {
		db.Close()
		return nil, initErr
	}
	return db, nil
}

// InitRootCacheSchema initializes tables and indices for file metadata.
func InitRootCacheSchema(db *sql.DB) *appfault.AppError {
	const ddl = `
	CREATE TABLE IF NOT EXISTS Files (
		FileId INTEGER PRIMARY KEY AUTOINCREMENT,
		RepoUrl TEXT,
		RelativePath TEXT UNIQUE,
		AbsolutePath TEXT,
		FileSize INTEGER,
		ModifiedTime INTEGER,
		FolderSlug TEXT,
		IsKeep INTEGER DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_files_slug ON Files(FolderSlug);
	`
	_, err := db.Exec(ddl)
	if err != nil {
		return appfault.WrapSimple(err, "init root schema")
	}
	return nil
}

// InitSlugCacheSchema initializes tables for indexed line contents.
func InitSlugCacheSchema(db *sql.DB) *appfault.AppError {
	const ddl = `
	CREATE TABLE IF NOT EXISTS Lines (
		LineId INTEGER PRIMARY KEY AUTOINCREMENT,
		RelativePath TEXT,
		LineNumber INTEGER,
		Content TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_lines_path ON Lines(RelativePath);
	`
	_, err := db.Exec(ddl)
	if err != nil {
		return appfault.WrapSimple(err, "init slug schema")
	}
	return nil
}

// InsertCacheFile upserts a file record in the root cache database.
func InsertCacheFile(db *sql.DB, rec CacheFileRecord) *appfault.AppError {
	const query = `
	INSERT INTO Files (RepoUrl, RelativePath, AbsolutePath, FileSize, ModifiedTime, FolderSlug, IsKeep)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(RelativePath) DO UPDATE SET
		FileSize=excluded.FileSize,
		ModifiedTime=excluded.ModifiedTime,
		FolderSlug=excluded.FolderSlug;
	`
	isKeepInt := 0
	if rec.IsKeep {
		isKeepInt = 1
	}
	_, err := db.Exec(query, rec.RepoUrl, rec.RelativePath, rec.AbsolutePath, rec.FileSize, rec.ModifiedTime, rec.FolderSlug, isKeepInt)
	if err != nil {
		return appfault.WrapSimple(err, "insert cache file")
	}
	return nil
}

// InsertCachedLines writes file lines in a batch into the slug cache database.
func InsertCachedLines(db *sql.DB, relPath string, lines []string) *appfault.AppError {
	if _, delErr := db.Exec("DELETE FROM Lines WHERE RelativePath = ?", relPath); delErr != nil {
		return appfault.WrapSimple(delErr, "delete old lines")
	}
	tx, err := db.Begin()
	if err != nil {
		return appfault.WrapSimple(err, "begin lines tx")
	}
	defer tx.Rollback()

	stmt, prepErr := tx.Prepare("INSERT INTO Lines (RelativePath, LineNumber, Content) VALUES (?, ?, ?)")
	if prepErr != nil {
		return appfault.WrapSimple(prepErr, "prepare lines stmt")
	}
	defer stmt.Close()

	for idx, line := range lines {
		if _, execErr := stmt.Exec(relPath, idx+1, line); execErr != nil {
			return appfault.WrapSimple(execErr, "exec line insert")
		}
	}
	if commitErr := tx.Commit(); commitErr != nil {
		return appfault.WrapSimple(commitErr, "commit lines tx")
	}
	return nil
}

// ListCacheFiles retrieves all indexed files from the root database.
func ListCacheFiles(db *sql.DB) ([]CacheFileRecord, *appfault.AppError) {
	const query = "SELECT FileId, RepoUrl, RelativePath, AbsolutePath, FileSize, ModifiedTime, FolderSlug, IsKeep FROM Files ORDER BY RelativePath ASC"
	rows, err := db.Query(query)
	if err != nil {
		return nil, appfault.WrapSimple(err, "query cache files")
	}
	defer rows.Close()

	var results []CacheFileRecord
	for rows.Next() {
		var r CacheFileRecord
		var isKeepInt int
		if scanErr := rows.Scan(&r.FileId, &r.RepoUrl, &r.RelativePath, &r.AbsolutePath, &r.FileSize, &r.ModifiedTime, &r.FolderSlug, &isKeepInt); scanErr == nil {
			r.IsKeep = isKeepInt == 1
			results = append(results, r)
		}
	}
	return results, nil
}
