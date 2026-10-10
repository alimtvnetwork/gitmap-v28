package cmdclone

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

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

// OpenRepoCacheDBAt opens or creates the SQLite database at the target path.
func OpenRepoCacheDBAt(dbPath string) (*sql.DB, error) {
	dir := filepath.Dir(dbPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create repocache db directory %s: %w", dir, err)
		}
	}

	dsn := fmt.Sprintf("%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)", dbPath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database at %s: %w", dbPath, err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(10 * time.Minute)

	if err := initRepoCacheDBSchema(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize repocache schema in %s: %w", dbPath, err)
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
	upsertManifestQuery := `
	INSERT INTO RepoCacheManifest (ManifestPath, FileHash, FileMtime, TotalRepos, HasValidData, CreatedAt, UpdatedAt)
	VALUES (?, ?, ?, ?, 1, ?, ?)
	ON CONFLICT(ManifestPath) DO UPDATE SET
		FileHash = excluded.FileHash,
		FileMtime = excluded.FileMtime,
		TotalRepos = excluded.TotalRepos,
		HasValidData = 1,
		UpdatedAt = excluded.UpdatedAt
	RETURNING RepoCacheManifestId;
	`
	var manifestId int64
	err = tx.QueryRow(upsertManifestQuery, manifestPath, fileHash, fileMtime, len(entries), now, now).Scan(&manifestId)
	if err != nil {
		fallbackUpsert := `
		INSERT INTO RepoCacheManifest (ManifestPath, FileHash, FileMtime, TotalRepos, HasValidData, CreatedAt, UpdatedAt)
		VALUES (?, ?, ?, ?, 1, ?, ?)
		ON CONFLICT(ManifestPath) DO UPDATE SET
			FileHash = excluded.FileHash,
			FileMtime = excluded.FileMtime,
			TotalRepos = excluded.TotalRepos,
			HasValidData = 1,
			UpdatedAt = excluded.UpdatedAt;
		`
		if _, execErr := tx.Exec(fallbackUpsert, manifestPath, fileHash, fileMtime, len(entries), now, now); execErr != nil {
			return execErr
		}
		if scanErr := tx.QueryRow(`SELECT RepoCacheManifestId FROM RepoCacheManifest WHERE ManifestPath = ?`, manifestPath).Scan(&manifestId); scanErr != nil {
			return scanErr
		}
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

// ParseManifestFile parses a JSON manifest from disk into normalized RepoCacheEntry objects.
func ParseManifestFile(manifestPath string) ([]RepoCacheEntry, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}

	var rawItems []genericManifestItem
	if err := json.Unmarshal(data, &rawItems); err != nil {
		var wrapper genericManifestWrapper
		if wrapErr := json.Unmarshal(data, &wrapper); wrapErr != nil {
			return nil, fmt.Errorf("unable to parse manifest JSON %s: %w", manifestPath, err)
		}
		if len(wrapper.Records) > 0 {
			rawItems = wrapper.Records
		} else if len(wrapper.Repositories) > 0 {
			rawItems = wrapper.Repositories
		} else if len(wrapper.Items) > 0 {
			rawItems = wrapper.Items
		} else if len(wrapper.Repos) > 0 {
			rawItems = wrapper.Repos
		}
	}

	entries := make([]RepoCacheEntry, 0, len(rawItems))
	for _, item := range rawItems {
		cloneURL := resolveItemCloneURL(item)
		if len(cloneURL) == 0 {
			continue
		}
		owner, repo := FormatOwnerRepo(cloneURL)
		if len(item.RepoName) > 0 {
			repo = item.RepoName
		} else if len(item.Name) > 0 {
			repo = item.Name
		}
		if len(item.Owner) > 0 {
			owner = item.Owner
		}
		isSSH := isSSHProtocol(cloneURL)
		isHTTPS := isHTTPSProtocol(cloneURL)
		entries = append(entries, RepoCacheEntry{
			RepoName:     repo,
			Owner:        owner,
			CloneUrl:     cloneURL,
			IsSSH:        isSSH,
			IsHTTPS:      isHTTPS,
			Branch:       item.Branch,
			RelativePath: item.RelativePath,
		})
	}

	return entries, nil
}

func resolveItemCloneURL(item genericManifestItem) string {
	if len(item.URL) > 0 {
		return item.URL
	}
	if len(item.CloneUrl) > 0 {
		return item.CloneUrl
	}
	if len(item.HTTPSUrl) > 0 {
		return item.HTTPSUrl
	}
	if len(item.SSHUrl) > 0 {
		return item.SSHUrl
	}
	if len(item.DiscoveredURL) > 0 {
		return item.DiscoveredURL
	}
	return item.RemoteURL
}

func isSSHProtocol(url string) bool {
	low := strings.ToLower(url)
	return strings.HasPrefix(low, "git@") || strings.HasPrefix(low, "ssh://")
}

func isHTTPSProtocol(url string) bool {
	low := strings.ToLower(url)
	return strings.HasPrefix(low, "https://") || strings.HasPrefix(low, "http://")
}

// GetOrSyncManifestEntries gets cached entries if valid, or reads disk, syncs to DB, and returns them.
func GetOrSyncManifestEntries(db *sql.DB, manifestPath string) ([]RepoCacheEntry, error) {
	fileHash, fileMtime, err := ComputeFileHashAndMtime(manifestPath)
	if err != nil {
		return nil, err
	}

	entries, isHit, err := GetCachedManifestEntries(db, manifestPath, fileHash, fileMtime)
	if err == nil && isHit {
		return entries, nil
	}

	parsedEntries, parseErr := ParseManifestFile(manifestPath)
	if parseErr != nil {
		return nil, parseErr
	}

	_ = StoreManifestEntries(db, manifestPath, fileHash, fileMtime, parsedEntries)
	return parsedEntries, nil
}
