package cmdcache

import (
	"database/sql"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

// OpenRootCacheDB opens or creates the root split database for cache metadata.
func OpenRootCacheDB(repoRoot string) (*sql.DB, error) {
	dir := store.ResolveSplitDbDir("cache", "root", repoRoot)
	dbPath := filepath.Join(dir, store.DbFileName)
	db, err := sql.Open("sqlite", dbPath+"?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, apperror.WrapSimple(err, "open root cache db")
	}
	if initErr := InitRootCacheSchema(db); initErr != nil {
		db.Close()
		return nil, initErr
	}
	return db, nil
}

// OpenSlugCacheDB opens or creates a partition split database for content lines.
func OpenSlugCacheDB(slug, repoRoot string) (*sql.DB, error) {
	dir := store.ResolveSplitDbDir("cache", slug, repoRoot)
	dbPath := filepath.Join(dir, store.DbFileName)
	db, err := sql.Open("sqlite", dbPath+"?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, apperror.WrapSimple(err, "open slug cache db")
	}
	if initErr := InitSlugCacheSchema(db); initErr != nil {
		db.Close()
		return nil, initErr
	}
	return db, nil
}

// InitRootCacheSchema initializes tables and indices for file metadata.
func InitRootCacheSchema(db *sql.DB) error {
	const ddl = `
	CREATE TABLE IF NOT EXISTS Files (
		FileId INTEGER PRIMARY KEY AUTOINCREMENT,
		RepoUrl TEXT,
		RelativePath TEXT UNIQUE,
		AbsolutePath TEXT,
		FileSize INTEGER,
		ModifiedTime INTEGER,
		FolderSlug TEXT,
		IsKeep BOOLEAN DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_files_slug ON Files(FolderSlug);
	`
	_, err := db.Exec(ddl)
	return err
}

// InitSlugCacheSchema initializes tables for indexed line contents.
func InitSlugCacheSchema(db *sql.DB) error {
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
	return err
}

// InsertCacheFile upserts a file record in the root cache database.
func InsertCacheFile(db *sql.DB, rec CacheFileRecord) error {
	const query = `
	INSERT INTO Files (RepoUrl, RelativePath, AbsolutePath, FileSize, ModifiedTime, FolderSlug, IsKeep)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(RelativePath) DO UPDATE SET
		FileSize=excluded.FileSize,
		ModifiedTime=excluded.ModifiedTime,
		FolderSlug=excluded.FolderSlug;
	`
	_, err := db.Exec(query, rec.RepoUrl, rec.RelativePath, rec.AbsolutePath, rec.FileSize, rec.ModifiedTime, rec.FolderSlug, rec.IsKeep)
	return err
}

// InsertCachedLines writes file lines in a batch into the slug cache database.
func InsertCachedLines(db *sql.DB, relPath string, lines []string) error {
	_, _ = db.Exec("DELETE FROM Lines WHERE RelativePath = ?", relPath)
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, prepErr := tx.Prepare("INSERT INTO Lines (RelativePath, LineNumber, Content) VALUES (?, ?, ?)")
	if prepErr != nil {
		return prepErr
	}
	defer stmt.Close()

	for idx, line := range lines {
		if _, execErr := stmt.Exec(relPath, idx+1, line); execErr != nil {
			return execErr
		}
	}
	return tx.Commit()
}

// ListCacheFiles retrieves all indexed files from the root database.
func ListCacheFiles(db *sql.DB) ([]CacheFileRecord, error) {
	const query = "SELECT FileId, RepoUrl, RelativePath, AbsolutePath, FileSize, ModifiedTime, FolderSlug, IsKeep FROM Files ORDER BY RelativePath ASC"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []CacheFileRecord
	for rows.Next() {
		var r CacheFileRecord
		if scanErr := rows.Scan(&r.FileId, &r.RepoUrl, &r.RelativePath, &r.AbsolutePath, &r.FileSize, &r.ModifiedTime, &r.FolderSlug, &r.IsKeep); scanErr == nil {
			results = append(results, r)
		}
	}
	return results, nil
}
