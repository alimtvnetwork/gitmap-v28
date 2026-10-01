package store

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	_ "modernc.org/sqlite"
)

// GitIgnoreRecord models a single repository's cached ignore check metadata.
type GitIgnoreRecord struct {
	ID              int64  `json:"id"`
	RepoPath        string `json:"repoPath"`
	RepoSlug        string `json:"repoSlug"`
	LastCheckedAt   int64  `json:"lastCheckedAt"`
	Status          string `json:"status"`
	RemediatedCount int    `json:"remediatedCount"`
	DurationMs      int64  `json:"durationMs"`
	IsActive        bool   `json:"isActive"`
}

// GitIgnoreSplitDB provides managed operations against .gitmap/data/gitignore/cache/sql.db.
type GitIgnoreSplitDB struct {
	conn *sql.DB
	Path string
}

const sqlCreateGitIgnoreTable = `
CREATE TABLE IF NOT EXISTS gitignore_repo_cache (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    repo_path TEXT NOT NULL UNIQUE COLLATE NOCASE,
    repo_slug TEXT NOT NULL,
    last_checked_at INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'clean',
    remediated_count INTEGER NOT NULL DEFAULT 0,
    duration_ms INTEGER NOT NULL DEFAULT 0,
    is_active INTEGER NOT NULL DEFAULT 1
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_gitignore_cache_repo_path ON gitignore_repo_cache(repo_path COLLATE NOCASE);
CREATE INDEX IF NOT EXISTS idx_gitignore_cache_last_checked ON gitignore_repo_cache(last_checked_at);
CREATE INDEX IF NOT EXISTS idx_gitignore_cache_is_active ON gitignore_repo_cache(is_active);
`

const sqlSelectLastIgnoreCheck = `
SELECT id, repo_path, repo_slug, last_checked_at, status, remediated_count, duration_ms, is_active
FROM gitignore_repo_cache
WHERE repo_path = ? COLLATE NOCASE AND is_active = 1
LIMIT 1;
`

const sqlUpsertIgnoreCheck = `
INSERT INTO gitignore_repo_cache (
    repo_path, repo_slug, last_checked_at, status, remediated_count, duration_ms, is_active
) VALUES (?, ?, ?, ?, ?, ?, 1)
ON CONFLICT(repo_path) DO UPDATE SET
    repo_slug = excluded.repo_slug,
    last_checked_at = excluded.last_checked_at,
    status = excluded.status,
    remediated_count = excluded.remediated_count,
    duration_ms = excluded.duration_ms,
    is_active = 1;
`

const sqlSelectListCached = `
SELECT id, repo_path, repo_slug, last_checked_at, status, remediated_count, duration_ms, is_active
FROM gitignore_repo_cache
ORDER BY last_checked_at DESC;
`

// ResolveGitIgnoreSplitDbPath returns the canonical path to the gitignore cache SQLite database.
func ResolveGitIgnoreSplitDbPath() string {
	return filepath.ToSlash(filepath.Join(BinaryDataDir(), "gitignore", "cache", DbFileName))
}

// OpenGitIgnoreSplitDB opens or creates the isolated SQLite DB for gitignore caching.
func OpenGitIgnoreSplitDB() (*GitIgnoreSplitDB, error) {
	dbPath := ResolveGitIgnoreSplitDbPath()
	_ = os.MkdirAll(filepath.Dir(dbPath), 0755)
	conn, appErr := OpenSQLiteDB(dbPath)
	if appErr != nil {
		return nil, appErr
	}
	db := &GitIgnoreSplitDB{conn: conn, Path: dbPath}
	if err := db.EnsureGitIgnoreTable(); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return db, nil
}

// Close closes the underlying SQLite database connection.
func (s *GitIgnoreSplitDB) Close() error {
	if s == nil || s.conn == nil {
		return nil
	}
	return s.conn.Close()
}

// EnsureGitIgnoreTable idempotently creates the cache table and indices.
func (s *GitIgnoreSplitDB) EnsureGitIgnoreTable() error {
	if s == nil || s.conn == nil {
		return apperror.NewSimple("split db connection is nil", "E1020")
	}
	if _, err := s.conn.Exec(sqlCreateGitIgnoreTable); err != nil {
		return apperror.WrapSimple(err, "ensure gitignore_repo_cache table")
	}
	return nil
}

// GetLastIgnoreCheck retrieves the cached check record for the specified repository path.
func (s *GitIgnoreSplitDB) GetLastIgnoreCheck(repoPath string) (*GitIgnoreRecord, error) {
	if s == nil || s.conn == nil {
		return nil, apperror.NewSimple("split db connection is nil", "E1020")
	}
	cleanPath := filepath.ToSlash(filepath.Clean(repoPath))
	row := s.conn.QueryRow(sqlSelectLastIgnoreCheck, cleanPath)
	return scanIgnoreRecord(row)
}

func scanIgnoreRecord(row *sql.Row) (*GitIgnoreRecord, error) {
	var r GitIgnoreRecord
	var activeInt int
	err := row.Scan(&r.ID, &r.RepoPath, &r.RepoSlug, &r.LastCheckedAt, &r.Status, &r.RemediatedCount, &r.DurationMs, &activeInt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, apperror.WrapSimple(err, "scan gitignore record")
	}
	r.IsActive = (activeInt == 1)
	return &r, nil
}

// RecordIgnoreCheck upserts a repository check record into the database.
func (s *GitIgnoreSplitDB) RecordIgnoreCheck(record GitIgnoreRecord) error {
	if s == nil || s.conn == nil {
		return apperror.NewSimple("split db connection is nil", "E1020")
	}
	cleanPath := filepath.ToSlash(filepath.Clean(record.RepoPath))
	_, err := s.conn.Exec(sqlUpsertIgnoreCheck,
		cleanPath, record.RepoSlug, record.LastCheckedAt,
		record.Status, record.RemediatedCount, record.DurationMs,
	)
	if err != nil {
		return apperror.WrapSimple(err, "record ignore check")
	}
	return nil
}

// IsCheckRecent reports whether an active check exists within the TTL duration and is clean or skipped.
func (s *GitIgnoreSplitDB) IsCheckRecent(repoPath string, ttl time.Duration) bool {
	if ttl <= 0 {
		return false
	}
	record, err := s.GetLastIgnoreCheck(repoPath)
	if err != nil || record == nil || !record.IsActive {
		return false
	}
	elapsed := time.Since(time.Unix(record.LastCheckedAt, 0))
	if elapsed >= ttl {
		return false
	}
	return record.Status == "clean" || record.Status == "skipped"
}

// InvalidateCache marks a single repository's cache record as inactive.
func (s *GitIgnoreSplitDB) InvalidateCache(repoPath string) error {
	if s == nil || s.conn == nil {
		return apperror.NewSimple("split db connection is nil", "E1020")
	}
	cleanPath := filepath.ToSlash(filepath.Clean(repoPath))
	_, err := s.conn.Exec("UPDATE gitignore_repo_cache SET is_active = 0 WHERE repo_path = ? COLLATE NOCASE", cleanPath)
	if err != nil {
		return apperror.WrapSimple(err, "invalidate cache for "+cleanPath)
	}
	return nil
}

// InvalidateAll marks all repository cache records as inactive.
func (s *GitIgnoreSplitDB) InvalidateAll() error {
	if s == nil || s.conn == nil {
		return apperror.NewSimple("split db connection is nil", "E1020")
	}
	_, err := s.conn.Exec("UPDATE gitignore_repo_cache SET is_active = 0")
	if err != nil {
		return apperror.WrapSimple(err, "invalidate all ignore cache")
	}
	return nil
}

// ListCachedChecks retrieves all cached repository records ordered by last check timestamp descending.
func (s *GitIgnoreSplitDB) ListCachedChecks() ([]GitIgnoreRecord, error) {
	if s == nil || s.conn == nil {
		return nil, apperror.NewSimple("split db connection is nil", "E1020")
	}
	rows, err := s.conn.Query(sqlSelectListCached)
	if err != nil {
		return nil, apperror.WrapSimple(err, "list cached checks")
	}
	defer rows.Close()
	return iterateCachedChecks(rows)
}

func iterateCachedChecks(rows *sql.Rows) ([]GitIgnoreRecord, error) {
	var records []GitIgnoreRecord
	for rows.Next() {
		r, err := scanCachedRow(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	return records, rows.Err()
}

func scanCachedRow(rows *sql.Rows) (GitIgnoreRecord, error) {
	var r GitIgnoreRecord
	var activeInt int
	err := rows.Scan(&r.ID, &r.RepoPath, &r.RepoSlug, &r.LastCheckedAt, &r.Status, &r.RemediatedCount, &r.DurationMs, &activeInt)
	if err != nil {
		return r, apperror.WrapSimple(err, "scan cached row")
	}
	r.IsActive = (activeInt == 1)
	return r, nil
}

// FilterReposNeedingIgnoreCheck filters out repositories whose .gitignore status
// has already been audited within the specified TTL duration.
//
// Fails open (returns all repos) if TTL <= 0 or if split-DB fails to open.
func FilterReposNeedingIgnoreCheck(repos []model.ScanRecord, ttl time.Duration) ([]model.ScanRecord, error) {
	if ttl <= 0 || len(repos) == 0 {
		return repos, nil
	}
	db, err := OpenGitIgnoreSplitDB()
	if err != nil {
		return repos, nil
	}
	defer db.Close()
	return filterReposWithDB(db, repos, ttl), nil
}

func filterReposWithDB(db *GitIgnoreSplitDB, repos []model.ScanRecord, ttl time.Duration) []model.ScanRecord {
	var needingCheck []model.ScanRecord
	for _, repo := range repos {
		isRecent := db.IsCheckRecent(repo.AbsolutePath, ttl)
		if !isRecent {
			needingCheck = append(needingCheck, repo)
		}
	}
	return needingCheck
}

// RecordIgnoreCheckResult persists a repository inspection outcome into the Split-DB.
func RecordIgnoreCheckResult(repoPath, slug, status string, count int, dur time.Duration) error {
	db, err := OpenGitIgnoreSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()
	record := GitIgnoreRecord{
		RepoPath:        repoPath,
		RepoSlug:        slug,
		LastCheckedAt:   time.Now().Unix(),
		Status:          status,
		RemediatedCount: count,
		DurationMs:      dur.Milliseconds(),
		IsActive:        true,
	}
	return db.RecordIgnoreCheck(record)
}
