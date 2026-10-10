// Package cmdsummary implements SQLite Split-DB caching for git release summaries and repository states.
package cmdsummary

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

var (
	summaryDBInstance *sql.DB
	summaryDBLock     sync.Mutex
)

// resolveSummaryDBPath discovers the target Split-DB SQLite storage file path.
func resolveSummaryDBPath() string {
	cwd, err := os.Getwd()
	if err == nil && cwd != "" {
		localGitmapDir := filepath.Join(cwd, ".gitmap")
		if info, errStat := os.Stat(localGitmapDir); errStat == nil && info.IsDir() {
			return filepath.Join(localGitmapDir, "summary.db")
		}
	}

	userHome, errHome := os.UserHomeDir()
	if errHome == nil && userHome != "" {
		homeGitmapDir := filepath.Join(userHome, ".gitmap")
		_ = os.MkdirAll(homeGitmapDir, 0755)
		return filepath.Join(homeGitmapDir, "summary.db")
	}

	return "summary.db"
}

// OpenSummaryDB initializes or retrieves the shared Split SQLite database handle.
func OpenSummaryDB() (*sql.DB, error) {
	summaryDBLock.Lock()
	defer summaryDBLock.Unlock()

	if summaryDBInstance != nil {
		if errPing := summaryDBInstance.Ping(); errPing == nil {
			return summaryDBInstance, nil
		}
	}

	dbPath := resolveSummaryDBPath()
	dir := filepath.Dir(dbPath)
	if dir != "" {
		_ = os.MkdirAll(dir, 0755)
	}

	dsn := fmt.Sprintf("%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)", dbPath)
	db, errOpen := sql.Open("sqlite", dsn)
	if errOpen != nil {
		return nil, fmt.Errorf("failed to open summary.db at %s: %w", dbPath, errOpen)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(10 * time.Minute)

	if errMigrate := initSummaryDBSchema(db); errMigrate != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to migrate summary.db schema: %w", errMigrate)
	}

	summaryDBInstance = db
	return summaryDBInstance, nil
}

func initSummaryDBSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS Repositories (
		repo_id INTEGER PRIMARY KEY AUTOINCREMENT,
		repo_url TEXT UNIQUE NOT NULL,
		canonical_slug TEXT NOT NULL,
		local_path TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS ReleaseSummaries (
		summary_id INTEGER PRIMARY KEY AUTOINCREMENT,
		repo_id INTEGER NOT NULL REFERENCES Repositories(repo_id) ON DELETE CASCADE,
		tag_name TEXT NOT NULL,
		tag_commit_hash TEXT NOT NULL,
		release_date DATETIME NOT NULL,
		summary_gist TEXT NOT NULL,
		heated_files_json TEXT NOT NULL,
		word_count INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(repo_id, tag_name, tag_commit_hash)
	);

	CREATE TABLE IF NOT EXISTS RepoHeadStates (
		repo_id INTEGER PRIMARY KEY REFERENCES Repositories(repo_id) ON DELETE CASCADE,
		head_commit_hash TEXT NOT NULL,
		is_dirty BOOLEAN NOT NULL DEFAULT 0,
		dirty_files_count INTEGER NOT NULL DEFAULT 0,
		last_activity_at DATETIME NOT NULL,
		last_scanned_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE VIEW IF NOT EXISTS ViewRepoReleaseSummaries AS
	SELECT 
		r.repo_url,
		r.canonical_slug,
		r.local_path,
		h.head_commit_hash,
		h.is_dirty,
		h.dirty_files_count,
		h.last_activity_at,
		s.tag_name,
		s.tag_commit_hash,
		s.release_date,
		s.summary_gist,
		s.heated_files_json,
		s.word_count
	FROM Repositories r
	JOIN RepoHeadStates h ON r.repo_id = h.repo_id
	LEFT JOIN ReleaseSummaries s ON r.repo_id = s.repo_id
	ORDER BY r.canonical_slug ASC, s.release_date DESC;

	CREATE INDEX IF NOT EXISTS idx_release_summaries_lookup 
	ON ReleaseSummaries (repo_id, tag_name, tag_commit_hash);

	CREATE INDEX IF NOT EXISTS idx_repo_head_states_activity 
	ON RepoHeadStates (last_activity_at);
	`
	_, err := db.Exec(schema)
	return err
}

func getOrCreateRepoID(db *sql.DB, repoURL, slug, localPath string) (int64, error) {
	if repoURL == "" {
		repoURL = slug
	}

	var repoID int64
	query := `SELECT repo_id FROM Repositories WHERE repo_url = ? LIMIT 1;`
	if err := db.QueryRow(query, repoURL).Scan(&repoID); err == nil {
		return repoID, nil
	}

	insertSQL := `INSERT INTO Repositories (repo_url, canonical_slug, local_path) VALUES (?, ?, ?)
	              ON CONFLICT(repo_url) DO UPDATE SET canonical_slug = excluded.canonical_slug, local_path = excluded.local_path;`
	res, err := db.Exec(insertSQL, repoURL, slug, localPath)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetCachedRelease checks SQLite summary.db for a matching release gist and heated files.
func GetCachedRelease(repoURL, tag, tagCommitHash string) (*ReleaseSummaryRecord, bool) {
	db, errDB := OpenSummaryDB()
	if errDB != nil {
		return nil, false
	}

	var row *sql.Row
	if tagCommitHash != "" {
		query := `
		SELECT s.tag_name, s.tag_commit_hash, s.release_date, s.summary_gist, s.heated_files_json, s.word_count
		FROM ReleaseSummaries s
		JOIN Repositories r ON s.repo_id = r.repo_id
		WHERE r.repo_url = ? AND s.tag_name = ? AND s.tag_commit_hash = ?
		LIMIT 1;`
		row = db.QueryRow(query, repoURL, tag, tagCommitHash)
	} else {
		query := `
		SELECT s.tag_name, s.tag_commit_hash, s.release_date, s.summary_gist, s.heated_files_json, s.word_count
		FROM ReleaseSummaries s
		JOIN Repositories r ON s.repo_id = r.repo_id
		WHERE r.repo_url = ? AND s.tag_name = ?
		ORDER BY s.release_date DESC
		LIMIT 1;`
		row = db.QueryRow(query, repoURL, tag)
	}

	var rec ReleaseSummaryRecord
	var heatedJSON string
	var relDate time.Time

	if err := row.Scan(&rec.TagName, &rec.TagCommitHash, &relDate, &rec.SummaryGist, &heatedJSON, &rec.WordCount); err != nil {
		return nil, false
	}

	rec.ReleaseDate = relDate.Format("2006-01-02")
	rec.IsCacheHit = true
	if heatedJSON != "" {
		_ = json.Unmarshal([]byte(heatedJSON), &rec.HeatedFiles)
	}

	return &rec, true
}

// RepoHeadStateRecord holds cached repository head state from SQLite summary.db.
type RepoHeadStateRecord struct {
	HeadCommitHash  string
	IsDirty         bool
	DirtyFilesCount int
	LastActivityAt  time.Time
}

// GetCachedHeadState retrieves the cached head state for a repository URL.
func GetCachedHeadState(repoURL string) (*RepoHeadStateRecord, bool) {
	db, errDB := OpenSummaryDB()
	if errDB != nil {
		return nil, false
	}

	query := `
	SELECT h.head_commit_hash, h.is_dirty, h.dirty_files_count, h.last_activity_at
	FROM RepoHeadStates h
	JOIN Repositories r ON h.repo_id = r.repo_id
	WHERE r.repo_url = ?
	LIMIT 1;`

	var rec RepoHeadStateRecord
	row := db.QueryRow(query, repoURL)
	if err := row.Scan(&rec.HeadCommitHash, &rec.IsDirty, &rec.DirtyFilesCount, &rec.LastActivityAt); err != nil {
		return nil, false
	}

	return &rec, true
}

// SaveCachedRelease persists a release gist and its heated files into Split SQLite summary.db.
func SaveCachedRelease(repoURL, slug, localPath string, rec ReleaseSummaryRecord) error {
	db, errDB := OpenSummaryDB()
	if errDB != nil {
		return errDB
	}

	repoID, errRepo := getOrCreateRepoID(db, repoURL, slug, localPath)
	if errRepo != nil {
		return errRepo
	}

	heatedJSON, _ := json.Marshal(rec.HeatedFiles)
	parsedDate, errDate := time.Parse("2006-01-02", rec.ReleaseDate)
	if errDate != nil {
		parsedDate = time.Now()
	}

	upsert := `
	INSERT INTO ReleaseSummaries (repo_id, tag_name, tag_commit_hash, release_date, summary_gist, heated_files_json, word_count)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(repo_id, tag_name, tag_commit_hash) DO UPDATE SET
		summary_gist = excluded.summary_gist,
		heated_files_json = excluded.heated_files_json,
		word_count = excluded.word_count;`

	_, err := db.Exec(upsert, repoID, rec.TagName, rec.TagCommitHash, parsedDate, rec.SummaryGist, string(heatedJSON), rec.WordCount)
	return err
}

// UpdateRepoHeadState saves or updates the head state of a repository in SQLite.
func UpdateRepoHeadState(repoURL, slug, localPath, headHash string, isDirty bool, dirtyCount int, lastActivity time.Time) error {
	db, errDB := OpenSummaryDB()
	if errDB != nil {
		return errDB
	}

	repoID, errRepo := getOrCreateRepoID(db, repoURL, slug, localPath)
	if errRepo != nil {
		return errRepo
	}

	upsert := `
	INSERT INTO RepoHeadStates (repo_id, head_commit_hash, is_dirty, dirty_files_count, last_activity_at, last_scanned_at)
	VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	ON CONFLICT(repo_id) DO UPDATE SET
		head_commit_hash = excluded.head_commit_hash,
		is_dirty = excluded.is_dirty,
		dirty_files_count = excluded.dirty_files_count,
		last_activity_at = excluded.last_activity_at,
		last_scanned_at = CURRENT_TIMESTAMP;`

	_, err := db.Exec(upsert, repoID, headHash, isDirty, dirtyCount, lastActivity)
	return err
}
