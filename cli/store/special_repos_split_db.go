// Package store — special_repos_split_db.go manages SQLite Split-DB persistence for special repositories (repo-secrets & repo-cache).
package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

// SpecialReposDBFileName is the SQLite database filename for special repositories.
const SpecialReposDBFileName = "gitmap-special-repos.db"

const (
	sqlCreateSpecialRepository = `CREATE TABLE IF NOT EXISTS SpecialRepository (
    RepoKey TEXT PRIMARY KEY,
    ShortKey TEXT NOT NULL UNIQUE,
    DefaultName TEXT NOT NULL,
    ConfiguredName TEXT NOT NULL,
    Category TEXT NOT NULL,
    LocalPath TEXT NOT NULL DEFAULT '',
    RemoteURL TEXT NOT NULL DEFAULT '',
    IsPromptAnswered INTEGER NOT NULL DEFAULT 0,
    UserDecision TEXT NOT NULL DEFAULT 'pending',
    PromptedAt TEXT NOT NULL DEFAULT '',
    UpdatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);`

	sqlCreateSpecialRepoFolderSeq = `CREATE TABLE IF NOT EXISTS SpecialRepoFolderSeq (
    SpecialKey TEXT NOT NULL,
    RepoName TEXT NOT NULL,
    FolderPrefix TEXT NOT NULL,
    SeqNumber INTEGER NOT NULL,
    CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (SpecialKey, RepoName)
);`

	sqlSeedSpecialSecret = `INSERT OR IGNORE INTO SpecialRepository (
    RepoKey, ShortKey, DefaultName, ConfiguredName, Category, UserDecision
) VALUES ('repo-secrets', 'rs', 'repo-secrets', 'repo-secrets', 'secrets', 'pending');`

	sqlSeedSpecialCache = `INSERT OR IGNORE INTO SpecialRepository (
    RepoKey, ShortKey, DefaultName, ConfiguredName, Category, UserDecision
) VALUES ('repo-cache', 'rc', 'repo-cache', 'repo-cache', 'cache', 'pending');`
)

// SpecialRepositoryRecord represents a configured special repository in SQLite.
type SpecialRepositoryRecord struct {
	RepoKey          string `json:"repoKey"`
	ShortKey         string `json:"shortKey"`
	DefaultName      string `json:"defaultName"`
	ConfiguredName   string `json:"configuredName"`
	Category         string `json:"category"`
	LocalPath        string `json:"localPath"`
	RemoteURL        string `json:"remoteUrl"`
	IsPromptAnswered int    `json:"isPromptAnswered"`
	UserDecision     string `json:"userDecision"`
	PromptedAt       string `json:"promptedAt"`
	UpdatedAt        string `json:"updatedAt"`
}

// HasAnsweredPrompt returns true when the first-scan prompt has already been answered.
func (r SpecialRepositoryRecord) HasAnsweredPrompt() bool {
	return r.IsPromptAnswered > 0
}

// SpecialReposSplitDB manages dedicated SQLite storage for repo-secrets and repo-cache metadata.
type SpecialReposSplitDB struct {
	conn *sql.DB
	path string
}

// NormalizeSpecialRepoKey maps any alias to its canonical (repoKey, shortKey).
func NormalizeSpecialRepoKey(keyOrShort string) (string, string) {
	cleaned := strings.ToLower(strings.TrimSpace(keyOrShort))
	switch cleaned {
	case "rc", "repo-cache", "repo-storage", "cache", "storage":
		return "repo-cache", "rc"
	default:
		return "repo-secrets", "rs"
	}
}

// OpenSpecialReposSplitDB opens the special repositories database at BinaryDataDir().
func OpenSpecialReposSplitDB() (*SpecialReposSplitDB, error) {
	defaultPath := filepath.Join(BinaryDataDir(), SpecialReposDBFileName)
	return OpenSpecialReposSplitDBAt(defaultPath)
}

// OpenSpecialReposSplitDBAt opens or creates the special repositories database at dbPath.
func OpenSpecialReposSplitDBAt(dbPath string) (*SpecialReposSplitDB, error) {
	targetPath := resolveSpecialReposDBPath(dbPath)
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return nil, apperror.WrapSimple(err, "mkdir special-repos db dir")
	}
	conn, err := sql.Open("sqlite", targetPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open special-repos sqlite")
	}
	return initSpecialReposConn(conn, targetPath)
}

func resolveSpecialReposDBPath(dbPath string) string {
	if len(strings.TrimSpace(dbPath)) > 0 {
		return dbPath
	}
	return filepath.Join(BinaryDataDir(), SpecialReposDBFileName)
}

func initSpecialReposConn(conn *sql.DB, dbPath string) (*SpecialReposSplitDB, error) {
	if err := ConfigureSQLiteConn(conn); err != nil {
		_ = conn.Close()
		return nil, apperror.WrapSimple(err, "configure special-repos sqlite")
	}
	db := &SpecialReposSplitDB{conn: conn, path: dbPath}
	if err := db.InitSchema(); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return db, nil
}

// InitSchema creates the tables and seeds default special repository records.
func (db *SpecialReposSplitDB) InitSchema() error {
	stmts := []string{
		sqlCreateSpecialRepository,
		sqlCreateSpecialRepoFolderSeq,
		sqlSeedSpecialSecret,
		sqlSeedSpecialCache,
	}
	for _, stmt := range stmts {
		if _, err := db.conn.Exec(stmt); err != nil {
			return apperror.WrapSimple(err, "init special-repos schema")
		}
	}
	return nil
}

// Close closes the underlying database connection.
func (db *SpecialReposSplitDB) Close() error {
	if db.conn == nil {
		return nil
	}
	return db.conn.Close()
}

// Conn returns the underlying SQL connection.
func (db *SpecialReposSplitDB) Conn() *sql.DB {
	return db.conn
}

// Path returns the SQLite database path.
func (db *SpecialReposSplitDB) Path() string {
	return db.path
}

// GetSpecialRepo retrieves a single special repository record by key or short alias.
func (db *SpecialReposSplitDB) GetSpecialRepo(keyOrShort string) (*SpecialRepositoryRecord, error) {
	repoKey, _ := NormalizeSpecialRepoKey(keyOrShort)
	query := `SELECT RepoKey, ShortKey, DefaultName, ConfiguredName, Category,
		LocalPath, RemoteURL, IsPromptAnswered, UserDecision, PromptedAt, UpdatedAt
		FROM SpecialRepository WHERE RepoKey = ?`
	row := db.conn.QueryRow(query, repoKey)
	return scanSpecialRepoRow(row)
}

func scanSpecialRepoRow(scanner interface{ Scan(dest ...any) error }) (*SpecialRepositoryRecord, error) {
	var rec SpecialRepositoryRecord
	err := scanner.Scan(
		&rec.RepoKey, &rec.ShortKey, &rec.DefaultName, &rec.ConfiguredName, &rec.Category,
		&rec.LocalPath, &rec.RemoteURL, &rec.IsPromptAnswered, &rec.UserDecision, &rec.PromptedAt, &rec.UpdatedAt,
	)
	if err != nil {
		return nil, apperror.WrapSimple(err, "scan special repository")
	}
	return &rec, nil
}

// ListSpecialRepos returns all configured special repository records.
func (db *SpecialReposSplitDB) ListSpecialRepos() ([]SpecialRepositoryRecord, error) {
	query := `SELECT RepoKey, ShortKey, DefaultName, ConfiguredName, Category,
		LocalPath, RemoteURL, IsPromptAnswered, UserDecision, PromptedAt, UpdatedAt
		FROM SpecialRepository ORDER BY RepoKey DESC`
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, apperror.WrapSimple(err, "list special repositories")
	}
	defer rows.Close()
	return collectSpecialRepoRows(rows)
}

func collectSpecialRepoRows(rows *sql.Rows) ([]SpecialRepositoryRecord, error) {
	var list []SpecialRepositoryRecord
	for rows.Next() {
		rec, err := scanSpecialRepoRow(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *rec)
	}
	return list, nil
}

// MarkPromptAnswered records the one-time first-scan user decision in SQLite.
func (db *SpecialReposSplitDB) MarkPromptAnswered(keyOrShort, decision, localPath, remoteURL string) error {
	repoKey, _ := NormalizeSpecialRepoKey(keyOrShort)
	now := time.Now().UTC().Format(time.RFC3339)
	query := `UPDATE SpecialRepository
		SET IsPromptAnswered = 1, UserDecision = ?, LocalPath = ?, RemoteURL = ?, PromptedAt = ?, UpdatedAt = ?
		WHERE RepoKey = ?`
	if _, err := db.conn.Exec(query, decision, localPath, remoteURL, now, now, repoKey); err != nil {
		return apperror.WrapSimple(err, "mark special repo prompt answered")
	}
	return nil
}

// UpdateConfiguredName updates the custom repository folder name and local path.
func (db *SpecialReposSplitDB) UpdateConfiguredName(keyOrShort, newName, localPath string) error {
	repoKey, _ := NormalizeSpecialRepoKey(keyOrShort)
	now := time.Now().UTC().Format(time.RFC3339)
	query := `UPDATE SpecialRepository SET ConfiguredName = ?, LocalPath = ?, UpdatedAt = ? WHERE RepoKey = ?`
	if _, err := db.conn.Exec(query, strings.TrimSpace(newName), localPath, now, repoKey); err != nil {
		return apperror.WrapSimple(err, "update special repo configured name")
	}
	return nil
}

// ResolveOrCreateRepoFolder resolves or allocates the sequenced XX-<repoName> folder inside specialRepoRoot.
func (db *SpecialReposSplitDB) ResolveOrCreateRepoFolder(specialKey, specialRepoRoot, repoName string) (string, error) {
	_, shortKey := NormalizeSpecialRepoKey(specialKey)
	cleanRepo := sanitizeRepoFolderSlug(repoName)
	if err := os.MkdirAll(specialRepoRoot, 0755); err != nil {
		return "", apperror.WrapSimple(err, "mkdir special repo root")
	}
	if diskPrefix, seqNum, isFound := findExistingRepoFolderOnDisk(specialRepoRoot, cleanRepo); isFound {
		_ = db.saveRepoFolderSeq(shortKey, cleanRepo, diskPrefix, seqNum)
		return filepath.Join(specialRepoRoot, diskPrefix), nil
	}
	return db.allocateOrLoadRepoFolder(shortKey, specialRepoRoot, cleanRepo)
}

func (db *SpecialReposSplitDB) allocateOrLoadRepoFolder(shortKey, rootDir, cleanRepo string) (string, error) {
	if existingPrefix, isFound := db.queryStoredFolderPrefix(shortKey, cleanRepo); isFound {
		fullDir := filepath.Join(rootDir, existingPrefix)
		return fullDir, os.MkdirAll(fullDir, 0755)
	}
	nextSeq := db.computeNextRepoSeq(shortKey, rootDir)
	prefix := fmt.Sprintf("%02d-%s", nextSeq, cleanRepo)
	if err := db.saveRepoFolderSeq(shortKey, cleanRepo, prefix, nextSeq); err != nil {
		return "", err
	}
	fullDir := filepath.Join(rootDir, prefix)
	return fullDir, os.MkdirAll(fullDir, 0755)
}

func (db *SpecialReposSplitDB) queryStoredFolderPrefix(shortKey, cleanRepo string) (string, bool) {
	var prefix string
	err := db.conn.QueryRow(
		`SELECT FolderPrefix FROM SpecialRepoFolderSeq WHERE SpecialKey = ? AND RepoName = ?`,
		shortKey, cleanRepo,
	).Scan(&prefix)
	return prefix, err == nil && len(prefix) > 0
}

func (db *SpecialReposSplitDB) computeNextRepoSeq(shortKey, rootDir string) int {
	var maxDB int
	_ = db.conn.QueryRow(`SELECT COALESCE(MAX(SeqNumber), 0) FROM SpecialRepoFolderSeq WHERE SpecialKey = ?`, shortKey).Scan(&maxDB)
	maxDisk := maxSeqOnDisk(rootDir)
	if maxDisk > maxDB {
		return maxDisk + 1
	}
	return maxDB + 1
}

func (db *SpecialReposSplitDB) saveRepoFolderSeq(shortKey, cleanRepo, prefix string, seqNum int) error {
	query := `INSERT OR REPLACE INTO SpecialRepoFolderSeq (SpecialKey, RepoName, FolderPrefix, SeqNumber) VALUES (?, ?, ?, ?)`
	if _, err := db.conn.Exec(query, shortKey, cleanRepo, prefix, seqNum); err != nil {
		return apperror.WrapSimple(err, "save special repo folder sequence")
	}
	return nil
}

func findExistingRepoFolderOnDisk(rootDir, cleanRepo string) (string, int, bool) {
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return "", 0, false
	}
	for _, entry := range entries {
		if prefix, seq, isMatch := matchRepoFolderEntry(entry, cleanRepo); isMatch {
			return prefix, seq, true
		}
	}
	return "", 0, false
}

func matchRepoFolderEntry(entry os.DirEntry, cleanRepo string) (string, int, bool) {
	if !entry.IsDir() {
		return "", 0, false
	}
	seq, slug, isParsed := parseSequencedName(entry.Name())
	if !isParsed {
		return "", 0, false
	}
	return entry.Name(), seq, strings.EqualFold(slug, cleanRepo)
}

func maxSeqOnDisk(rootDir string) int {
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return 0
	}
	maxVal := 0
	for _, entry := range entries {
		seq, _, isParsed := parseSequencedName(entry.Name())
		if isParsed && seq > maxVal {
			maxVal = seq
		}
	}
	return maxVal
}

func parseSequencedName(name string) (int, string, bool) {
	parts := strings.SplitN(name, "-", 2)
	if len(parts) < 2 || len(parts[0]) < 2 {
		return 0, "", false
	}
	seq, err := strconv.Atoi(parts[0])
	if err != nil || seq <= 0 {
		return 0, "", false
	}
	return seq, parts[1], true
}

func sanitizeRepoFolderSlug(repoName string) string {
	cleaned := strings.ToLower(strings.TrimSpace(repoName))
	cleaned = strings.ReplaceAll(cleaned, " ", "-")
	cleaned = strings.ReplaceAll(cleaned, "_", "-")
	if len(cleaned) == 0 {
		return "default-repo"
	}
	return cleaned
}
