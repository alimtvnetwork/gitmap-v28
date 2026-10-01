package store

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

// CacheFileRecord represents metadata of a cached file in the root split-db.
type CacheFileRecord struct {
	FileId       int64  `json:"fileId"`
	RepoUrl      string `json:"repoUrl"`
	RelativePath string `json:"relativePath"`
	AbsolutePath string `json:"absolutePath"`
	FileSize     int64  `json:"fileSize"`
	ModifiedTime int64  `json:"modifiedTime"`
	FolderSlug   string `json:"folderSlug"`
	IsKeep       bool   `json:"isKeep"`
}

// CachedRepoSummary summarizes a cached repository.
type CachedRepoSummary struct {
	RepoSlug    string `json:"repoSlug"`
	RepoPath    string `json:"repoPath"`
	FileCount   int    `json:"fileCount"`
	TotalBytes  int64  `json:"totalBytes"`
	LastIndexed string `json:"lastIndexed"`
}

// CacheContextLine holds context around a matching line.
type CacheContextLine struct {
	LineNumber int    `json:"lineNumber"`
	Content    string `json:"content"`
	IsMatch    bool   `json:"isMatch"`
}

// ResolveRepoSlug derives a sanitized slug from the repository root directory.
func ResolveRepoSlug(repoRoot string) string {
	abs, err := filepath.Abs(repoRoot)
	if err != nil {
		return "default"
	}
	base := filepath.Base(abs)
	slug := SanitizeSlug(base)
	if slug == "" {
		return "default"
	}
	return slug
}

// ResolveCacheReposBaseDir returns .gitmap/cache/repos directory.
func ResolveCacheReposBaseDir(repoRoot string) string {
	absRoot, _ := filepath.Abs(repoRoot)
	dir := filepath.Join(absRoot, ".gitmap", "cache", "repos")
	_ = os.MkdirAll(dir, 0755)
	return filepath.ToSlash(dir)
}

// ResolveCacheRepoDir returns .gitmap/cache/repos/<slug> directory.
func ResolveCacheRepoDir(repoRoot string) string {
	baseDir := ResolveCacheReposBaseDir(repoRoot)
	slug := ResolveRepoSlug(repoRoot)
	dir := filepath.Join(baseDir, slug)
	_ = os.MkdirAll(dir, 0755)
	return filepath.ToSlash(dir)
}

func openSqliteConn(dbPath string) (*sql.DB, *appfault.AppError) {
	connStr := dbPath + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)"
	db, err := sql.Open("sqlite", connStr)
	if err != nil {
		return nil, appfault.WrapSimple(err, "open sqlite conn")
	}
	return db, nil
}

// OpenRootCacheDB opens or creates the root split database for cache metadata.
func OpenRootCacheDB(repoRoot string) (*sql.DB, *appfault.AppError) {
	repoDir := ResolveCacheRepoDir(repoRoot)
	dbPath := filepath.ToSlash(filepath.Join(repoDir, DbFileName))

	db, err := openSqliteConn(dbPath)
	if err != nil {
		return nil, err
	}
	if initErr := InitRootCacheSchema(db); initErr != nil {
		_ = db.Close()
		return nil, initErr
	}
	return db, nil
}

// OpenSlugCacheDB opens or creates a partition split database for content lines.
func OpenSlugCacheDB(slug, repoRoot string) (*sql.DB, *appfault.AppError) {
	cleanSlug := SanitizeSlug(slug)
	repoDir := ResolveCacheRepoDir(repoRoot)
	dbPath := filepath.ToSlash(filepath.Join(repoDir, cleanSlug+".db"))

	db, err := openSqliteConn(dbPath)
	if err != nil {
		return nil, err
	}
	if initErr := InitSlugCacheSchema(db); initErr != nil {
		_ = db.Close()
		return nil, initErr
	}
	return db, nil
}

// InitRootCacheSchema initializes tables and indices for file metadata.
func InitRootCacheSchema(db *sql.DB) *appfault.AppError {
	const ddl = `
	CREATE TABLE IF NOT EXISTS RepoMetadata (
		Key TEXT PRIMARY KEY,
		Value TEXT
	);
	CREATE TABLE IF NOT EXISTS FolderTree (
		FolderPath TEXT PRIMARY KEY,
		FolderSlug TEXT,
		FileCount INTEGER,
		TotalBytes INTEGER,
		UpdatedAt INTEGER
	);
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
	CREATE INDEX IF NOT EXISTS idx_files_mtime ON Files(ModifiedTime);
	`
	_, err := db.Exec(ddl)
	if err != nil {
		return appfault.WrapSimple(err, "init root cache schema")
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
	CREATE INDEX IF NOT EXISTS idx_lines_path_num ON Lines(RelativePath, LineNumber);
	`
	_, err := db.Exec(ddl)
	if err != nil {
		return appfault.WrapSimple(err, "init slug cache schema")
	}
	return nil
}

// InsertCacheFile upserts a file record in the root cache database.
func InsertCacheFile(db *sql.DB, rec CacheFileRecord) *appfault.AppError {
	const query = `
	INSERT INTO Files (RepoUrl, RelativePath, AbsolutePath, FileSize, ModifiedTime, FolderSlug, IsKeep)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(RelativePath) DO UPDATE SET
		AbsolutePath=excluded.AbsolutePath,
		FileSize=excluded.FileSize,
		ModifiedTime=excluded.ModifiedTime,
		FolderSlug=excluded.FolderSlug,
		IsKeep=MAX(Files.IsKeep, excluded.IsKeep);
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

// DeleteCacheFile removes a file record from the root cache database.
func DeleteCacheFile(db *sql.DB, relPath string) *appfault.AppError {
	_, err := db.Exec("DELETE FROM Files WHERE RelativePath = ?", relPath)
	if err != nil {
		return appfault.WrapSimple(err, "delete cache file")
	}
	return nil
}

// DeleteCachedLines deletes lines for a given file from a slug database.
func DeleteCachedLines(db *sql.DB, relPath string) *appfault.AppError {
	_, err := db.Exec("DELETE FROM Lines WHERE RelativePath = ?", relPath)
	if err != nil {
		return appfault.WrapSimple(err, "delete cached lines")
	}
	return nil
}

// InsertCachedLines writes file lines in a batch into the slug cache database.
func InsertCachedLines(db *sql.DB, relPath string, lines []string) *appfault.AppError {
	_ = DeleteCachedLines(db, relPath)
	tx, err := db.Begin()
	if err != nil {
		return appfault.WrapSimple(err, "begin lines tx")
	}
	defer tx.Rollback()

	if insertErr := insertLinesInTx(tx, relPath, lines); insertErr != nil {
		return insertErr
	}
	if commitErr := tx.Commit(); commitErr != nil {
		return appfault.WrapSimple(commitErr, "commit lines tx")
	}
	return nil
}

func insertLinesInTx(tx *sql.Tx, relPath string, lines []string) *appfault.AppError {
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

	return scanCacheFiles(rows)
}

func scanCacheFiles(rows *sql.Rows) ([]CacheFileRecord, *appfault.AppError) {
	var results []CacheFileRecord
	for rows.Next() {
		var r CacheFileRecord
		var isKeepInt int
		if scanErr := rows.Scan(&r.FileId, &r.RepoUrl, &r.RelativePath, &r.AbsolutePath, &r.FileSize, &r.ModifiedTime, &r.FolderSlug, &isKeepInt); scanErr == nil {
			r.IsKeep = isKeepInt > 0
			results = append(results, r)
		}
	}
	return results, nil
}

// GetCacheFile retrieves a single file record by RelativePath.
func GetCacheFile(db *sql.DB, relPath string) (*CacheFileRecord, *appfault.AppError) {
	const query = "SELECT FileId, RepoUrl, RelativePath, AbsolutePath, FileSize, ModifiedTime, FolderSlug, IsKeep FROM Files WHERE RelativePath = ?"
	row := db.QueryRow(query, relPath)
	var r CacheFileRecord
	var isKeepInt int
	if err := row.Scan(&r.FileId, &r.RepoUrl, &r.RelativePath, &r.AbsolutePath, &r.FileSize, &r.ModifiedTime, &r.FolderSlug, &isKeepInt); err != nil {
		return nil, appfault.WrapSimple(err, "get cache file")
	}
	r.IsKeep = isKeepInt > 0
	return &r, nil
}

// UpdateFolderTree upserts folder statistics into the FolderTree table.
func UpdateFolderTree(db *sql.DB, folderPath, folderSlug string, fileCount int, totalBytes int64) *appfault.AppError {
	const query = `
	INSERT INTO FolderTree (FolderPath, FolderSlug, FileCount, TotalBytes, UpdatedAt)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(FolderPath) DO UPDATE SET
		FolderSlug=excluded.FolderSlug,
		FileCount=excluded.FileCount,
		TotalBytes=excluded.TotalBytes,
		UpdatedAt=excluded.UpdatedAt;
	`
	now := time.Now().Unix()
	_, err := db.Exec(query, folderPath, folderSlug, fileCount, totalBytes, now)
	if err != nil {
		return appfault.WrapSimple(err, "update folder tree")
	}
	return nil
}

// SetRepoMetadata sets a key-value pair in RepoMetadata.
func SetRepoMetadata(db *sql.DB, key, value string) *appfault.AppError {
	const query = `INSERT INTO RepoMetadata (Key, Value) VALUES (?, ?) ON CONFLICT(Key) DO UPDATE SET Value=excluded.Value;`
	_, err := db.Exec(query, key, value)
	if err != nil {
		return appfault.WrapSimple(err, "set repo metadata")
	}
	return nil
}

// GetRepoMetadata retrieves a value from RepoMetadata.
func GetRepoMetadata(db *sql.DB, key string) (string, *appfault.AppError) {
	var val string
	err := db.QueryRow("SELECT Value FROM RepoMetadata WHERE Key = ?", key).Scan(&val)
	if err != nil {
		return "", appfault.WrapSimple(err, "get repo metadata")
	}
	return val, nil
}

// ListCachedRepos inspects the repos folder and lists all cached repository summaries.
func ListCachedRepos(repoRoot string) ([]CachedRepoSummary, *appfault.AppError) {
	baseDir := ResolveCacheReposBaseDir(repoRoot)
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, appfault.WrapSimple(err, "read cached repos dir")
	}

	var summaries []CachedRepoSummary
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		summary, hasSummary := summarizeRepoDir(baseDir, entry.Name())
		if hasSummary {
			summaries = append(summaries, summary)
		}
	}
	return summaries, nil
}

func summarizeRepoDir(baseDir, slug string) (CachedRepoSummary, bool) {
	dbPath := filepath.Join(baseDir, slug, DbFileName)
	summary := CachedRepoSummary{RepoSlug: slug}

	db, err := sql.Open("sqlite", dbPath+"?_pragma=busy_timeout(2000)")
	if err != nil {
		return summary, false
	}
	defer db.Close()

	_ = db.QueryRow("SELECT Value FROM RepoMetadata WHERE Key = 'repo_path'").Scan(&summary.RepoPath)
	_ = db.QueryRow("SELECT Value FROM RepoMetadata WHERE Key = 'last_indexed_at'").Scan(&summary.LastIndexed)
	_ = db.QueryRow("SELECT COUNT(*), COALESCE(SUM(FileSize), 0) FROM Files").Scan(&summary.FileCount, &summary.TotalBytes)
	return summary, true
}

// FetchContextLines retrieves context lines around a specific line number in a slug DB.
func FetchContextLines(slugDB *sql.DB, relPath string, matchLine, radius int) []CacheContextLine {
	minLine := matchLine - radius
	if minLine < 1 {
		minLine = 1
	}
	maxLine := matchLine + radius
	query := "SELECT LineNumber, Content FROM Lines WHERE RelativePath = ? AND LineNumber BETWEEN ? AND ? ORDER BY LineNumber ASC"
	rows, err := slugDB.Query(query, relPath, minLine, maxLine)
	if err != nil {
		return nil
	}
	defer rows.Close()

	return scanContextLines(rows, matchLine)
}

func scanContextLines(rows *sql.Rows, matchLine int) []CacheContextLine {
	var lines []CacheContextLine
	for rows.Next() {
		var num int
		var content string
		if rows.Scan(&num, &content) == nil {
			lines = append(lines, CacheContextLine{
				LineNumber: num,
				Content:    content,
				IsMatch:    num == matchLine,
			})
		}
	}
	return lines
}
