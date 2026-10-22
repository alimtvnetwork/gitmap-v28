package cmdclone

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"

	_ "modernc.org/sqlite"
)

const (
	DefaultRepoCacheDBRelativePath = ".gitmap/data/repocache/sql.db"
	FallbackRepoCacheDBPath        = ".gitmap/repocache.db"
)

// CustomRepoCacheDBPath allows tests to point the repo-cache database to an isolated location.
var CustomRepoCacheDBPath string

// RepoCacheManifest represents a cached repository manifest row in SQLite.
type RepoCacheManifest struct {
	RepoCacheManifestId int64     `json:"repoCacheManifestId"`
	ManifestPath        string    `json:"manifestPath"`
	FileHash            string    `json:"fileHash"`
	FileMtime           int64     `json:"fileMtime"`
	TotalRepos          int       `json:"totalRepos"`
	HasValidData        bool      `json:"hasValidData"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

// RepoCacheEntry represents a single cached repository entry in SQLite.
type RepoCacheEntry struct {
	RepoCacheEntryId    int64  `json:"repoCacheEntryId"`
	RepoCacheManifestId int64  `json:"repoCacheManifestId"`
	RepoName            string `json:"repoName"`
	Owner               string `json:"owner"`
	CloneUrl            string `json:"cloneUrl"`
	IsSSH               bool   `json:"isSSH"`
	IsHTTPS             bool   `json:"isHTTPS"`
	Branch              string `json:"branch"`
	RelativePath        string `json:"relativePath"`
}

type genericManifestItem struct {
	RepoName      string `json:"repoName"`
	Name          string `json:"name"`
	Owner         string `json:"owner"`
	URL           string `json:"url"`
	CloneUrl      string `json:"cloneUrl"`
	HTTPSUrl      string `json:"httpsUrl"`
	SSHUrl        string `json:"sshUrl"`
	DiscoveredURL string `json:"discoveredUrl"`
	RemoteURL     string `json:"remoteUrl"`
	Branch        string `json:"branch"`
	RelativePath  string `json:"relativePath"`
}

type genericManifestWrapper struct {
	Records      []genericManifestItem `json:"records"`
	Repositories []genericManifestItem `json:"repositories"`
	Items        []genericManifestItem `json:"items"`
	Repos        []genericManifestItem `json:"repos"`
}

// ResolveRepoCacheDBPath returns the filesystem path for the repo-cache Split-DB.
func ResolveRepoCacheDBPath() string {
	if len(strings.TrimSpace(CustomRepoCacheDBPath)) > 0 {
		return CustomRepoCacheDBPath
	}
	if _, err := os.Stat(DefaultRepoCacheDBRelativePath); err == nil {
		return DefaultRepoCacheDBRelativePath
	}
	if _, err := os.Stat(FallbackRepoCacheDBPath); err == nil {
		return FallbackRepoCacheDBPath
	}
	return DefaultRepoCacheDBRelativePath
}

// OpenRepoCacheDB opens and initializes the Split SQLite database with WAL mode and pragmas.
func OpenRepoCacheDB() (*sql.DB, error) {
	return OpenRepoCacheDBAt(ResolveRepoCacheDBPath())
}

func ensureRepoCacheDir(dbPath string) error {
	dir := filepath.Dir(dbPath)
	if dir == "" || dir == "." {
		return nil
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("failed to create repocache db directory %s", dir))
	}
	return nil
}

// OpenRepoCacheDBAt opens or creates the SQLite database at the target path.
func OpenRepoCacheDBAt(dbPath string) (*sql.DB, error) {
	if err := ensureRepoCacheDir(dbPath); err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)", dbPath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, apperror.WrapSimple(err, fmt.Sprintf("failed to open sqlite database at %s", dbPath))
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(10 * time.Minute)

	if err := initRepoCacheDBSchema(db); err != nil {
		db.Close()
		return nil, apperror.WrapSimple(err, fmt.Sprintf("failed to initialize repocache schema in %s", dbPath))
	}

	return db, nil
}

func initRepoCacheDBSchema(db *sql.DB) error {
	ddl := `
	CREATE TABLE IF NOT EXISTS RepoCacheManifest (
		RepoCacheManifestId INTEGER PRIMARY KEY AUTOINCREMENT,
		ManifestPath        TEXT NOT NULL UNIQUE,
		FileHash            TEXT NOT NULL,
		FileMtime           INTEGER NOT NULL,
		TotalRepos          INTEGER NOT NULL DEFAULT 0,
		HasValidData        INTEGER NOT NULL DEFAULT 1,
		CreatedAt           TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UpdatedAt           TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS RepoCacheEntry (
		RepoCacheEntryId    INTEGER PRIMARY KEY AUTOINCREMENT,
		RepoCacheManifestId INTEGER NOT NULL,
		RepoName            TEXT NOT NULL,
		Owner               TEXT NOT NULL,
		CloneUrl            TEXT NOT NULL,
		IsSSH               INTEGER NOT NULL DEFAULT 0,
		IsHTTPS             INTEGER NOT NULL DEFAULT 1,
		Branch              TEXT NOT NULL DEFAULT '',
		RelativePath        TEXT NOT NULL DEFAULT '',
		CreatedAt           TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(RepoCacheManifestId) REFERENCES RepoCacheManifest(RepoCacheManifestId) ON DELETE CASCADE
	);

	CREATE UNIQUE INDEX IF NOT EXISTS IdxRepoCacheManifest_ManifestPath ON RepoCacheManifest(ManifestPath);
	CREATE INDEX IF NOT EXISTS IdxRepoCacheEntry_ManifestId ON RepoCacheEntry(RepoCacheManifestId);
	CREATE INDEX IF NOT EXISTS IdxRepoCacheEntry_Owner ON RepoCacheEntry(Owner);
	CREATE INDEX IF NOT EXISTS IdxRepoCacheEntry_RepoName ON RepoCacheEntry(RepoName);
	`
	_, err := db.Exec(ddl)
	return err
}

// ComputeFileHashAndMtime calculates the SHA256 hex digest and modification time of a file.
func ComputeFileHashAndMtime(filePath string) (string, int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return "", 0, err
	}
	mtime := info.ModTime().Unix()

	f, err := os.Open(filePath)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", 0, err
	}
	return hex.EncodeToString(h.Sum(nil)), mtime, nil
}

// GetCachedManifestEntries queries cached repository entries if file hash and mtime match.
func GetCachedManifestEntries(db *sql.DB, manifestPath, currentHash string, currentMtime int64) ([]RepoCacheEntry, bool, error) {
	query := `SELECT RepoCacheManifestId, FileHash, FileMtime, TotalRepos, HasValidData
	          FROM RepoCacheManifest WHERE ManifestPath = ?`
	var manifestId int64
	var hash string
	var mtime int64
	var totalRepos int
	var hasValidData int

	err := db.QueryRow(query, manifestPath).Scan(&manifestId, &hash, &mtime, &totalRepos, &hasValidData)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	hasMatchingHash := (hash == currentHash)
	hasMatchingMtime := (mtime == currentMtime)
	isValidData := (hasValidData == 1)

	if !hasMatchingHash || !hasMatchingMtime || !isValidData {
		return nil, false, nil
	}

	entryQuery := `SELECT RepoCacheEntryId, RepoCacheManifestId, RepoName, Owner, CloneUrl, IsSSH, IsHTTPS, Branch, RelativePath
	               FROM RepoCacheEntry WHERE RepoCacheManifestId = ? ORDER BY RepoCacheEntryId ASC`
	rows, err := db.Query(entryQuery, manifestId)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var entries []RepoCacheEntry
	for rows.Next() {
		var e RepoCacheEntry
		var isSSHInt, isHTTPSInt int
		if err := rows.Scan(&e.RepoCacheEntryId, &e.RepoCacheManifestId, &e.RepoName, &e.Owner, &e.CloneUrl, &isSSHInt, &isHTTPSInt, &e.Branch, &e.RelativePath); err != nil {
			return nil, false, err
		}
		e.IsSSH = (isSSHInt == 1)
		e.IsHTTPS = (isHTTPSInt == 1)
		entries = append(entries, e)
	}

	return entries, true, nil
}

// StoreManifestEntries saves manifest metadata and all entries into the SQLite cache transactionally.
func StoreManifestEntries(db *sql.DB, manifestPath, fileHash string, fileMtime int64, entries []RepoCacheEntry) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now().UTC().Format(time.RFC3339)
	manifestId, upsertErr := upsertManifestRecord(tx, manifestPath, fileHash, fileMtime, len(entries), now)
	if upsertErr != nil {
		return upsertErr
	}

	if _, err := tx.Exec(`DELETE FROM RepoCacheEntry WHERE RepoCacheManifestId = ?`, manifestId); err != nil {
		return err
	}

	insertStmt, err := tx.Prepare(`
	INSERT INTO RepoCacheEntry (RepoCacheManifestId, RepoName, Owner, CloneUrl, IsSSH, IsHTTPS, Branch, RelativePath, CreatedAt)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer insertStmt.Close()

	for _, e := range entries {
		isSSHInt := 0
		if e.IsSSH {
			isSSHInt = 1
		}
		isHTTPSInt := 0
		if e.IsHTTPS {
			isHTTPSInt = 1
		}
		if _, err := insertStmt.Exec(manifestId, e.RepoName, e.Owner, e.CloneUrl, isSSHInt, isHTTPSInt, e.Branch, e.RelativePath, now); err != nil {
			return err
		}
	}

	return tx.Commit()
}
