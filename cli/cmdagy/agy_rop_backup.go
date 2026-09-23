package cmdagy

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const sqlCreateAGYBackupTable = `CREATE TABLE IF NOT EXISTS AGYConversationBackup (
	ConversationID TEXT PRIMARY KEY,
	ProjectSlug TEXT NOT NULL,
	ProjectPath TEXT NOT NULL,
	Title TEXT,
	StepCount INTEGER,
	TranscriptJSON TEXT,
	CreatedAt TIMESTAMP NOT NULL
);`

const sqlInsertAGYBackup = `INSERT OR REPLACE INTO AGYConversationBackup (
	ConversationID, ProjectSlug, ProjectPath, Title, StepCount, TranscriptJSON, CreatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?);`

// ResolveRepoSlug generates a safe folder name from project name or path.
func ResolveRepoSlug(projectName, projectPath string) string {
	raw := projectName
	if raw == "" {
		raw = filepath.Base(projectPath)
	}
	slug := strings.ToLower(strings.TrimSpace(raw))
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "\\", "-")
	slug = strings.ReplaceAll(slug, "/", "-")
	slug = strings.ReplaceAll(slug, ":", "-")
	if slug == "" || slug == "." {
		return "default"
	}
	return slug
}

// GetAGYBackupDBPath returns the Split-DB path at data/agy/<repo-slug>/agy.db.
func GetAGYBackupDBPath(repoSlug string) string {
	dataDir := store.BinaryDataDir()
	return filepath.Join(dataDir, "agy", repoSlug, "agy.db")
}

// OpenAGYBackupDB opens or initializes the SQLite database for a repo slug.
func OpenAGYBackupDB(repoSlug string) (*sql.DB, error) {
	dbPath := GetAGYBackupDBPath(repoSlug)
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, apperror.WrapSimple(err, "mkdir agy backup dir")
	}
	db, err := store.OpenSQLiteDB(dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open agy backup db")
	}
	if _, err := db.Exec(sqlCreateAGYBackupTable); err != nil {
		_ = db.Close()
		return nil, apperror.WrapSimple(err, "create agy backup table")
	}
	return db, nil
}

// ReadTranscriptForConv attempts to load the conversation transcript text.
func ReadTranscriptForConv(convID string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	candidate1 := filepath.Join(home, ".gemini", "antigravity", "brain", convID, ".system_generated", "logs", "transcript.jsonl")
	if data, readErr := os.ReadFile(candidate1); readErr == nil && len(data) > 0 {
		return string(data)
	}
	candidate2 := filepath.Join(home, ".gemini", "antigravity", "brain", convID, ".system_generated", "logs", "transcript_full.jsonl")
	if data, readErr := os.ReadFile(candidate2); readErr == nil && len(data) > 0 {
		return string(data)
	}
	return ""
}

// BackupConversation saves conversation state into the dedicated Split-DB.
func BackupConversation(db *sql.DB, backup AGYConversationBackup) error {
	now := backup.CreatedAt
	if now.IsZero() {
		now = time.Now().UTC()
	}
	_, err := db.Exec(sqlInsertAGYBackup,
		backup.ConversationID,
		backup.ProjectSlug,
		backup.ProjectPath,
		backup.Title,
		backup.StepCount,
		backup.TranscriptJSON,
		now,
	)
	if err != nil {
		return apperror.WrapSimple(err, "insert agy conversation backup")
	}
	return nil
}

// ClearAGYBackups removes all backup databases under data/agy/ or for a specific slug.
func ClearAGYBackups(repoSlug string) (int, int64, error) {
	baseDir := filepath.Join(store.BinaryDataDir(), "agy")
	if repoSlug != "" {
		targetDir := filepath.Join(baseDir, repoSlug)
		return purgeSingleAGYDir(targetDir)
	}
	return purgeAllAGYDirs(baseDir)
}

func purgeSingleAGYDir(targetDir string) (int, int64, error) {
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return 0, 0, nil
	}
	var totalBytes int64
	_ = filepath.Walk(targetDir, func(p string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			totalBytes += info.Size()
		}
		return nil
	})
	removeErr := os.RemoveAll(targetDir)
	if removeErr != nil {
		return 0, 0, apperror.WrapSimple(removeErr, "remove agy backup dir")
	}
	return 1, totalBytes, nil
}

func purgeAllAGYDirs(baseDir string) (int, int64, error) {
	entries, err := os.ReadDir(baseDir)
	if os.IsNotExist(err) {
		return 0, 0, nil
	}
	if err != nil {
		return 0, 0, apperror.WrapSimple(err, "read agy base dir")
	}
	count := 0
	var totalBytes int64
	for _, e := range entries {
		if e.IsDir() {
			c, bytesFreed, _ := purgeSingleAGYDir(filepath.Join(baseDir, e.Name()))
			count += c
			totalBytes += bytesFreed
		}
	}
	return count, totalBytes, nil
}
