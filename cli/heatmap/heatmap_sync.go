package heatmap

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// HeatmapLocalTrackingRecord maps to table HeatmapLocalTracking.
type HeatmapLocalTrackingRecord struct {
	HeatmapLocalTrackingId int64   `json:"heatmapLocalTrackingId"`
	RepoSlug               string  `json:"repoSlug"`
	AiName                 string  `json:"aiName"`
	SequencePrefix         string  `json:"sequencePrefix"`
	RepoCachePath          string  `json:"repoCachePath"`
	LastCommittedSha       string  `json:"lastCommittedSha"`
	LastCommitTimestamp    string  `json:"lastCommitTimestamp"`
	UncommittedChangeCount int     `json:"uncommittedChangeCount"`
	SyncStatus             string  `json:"syncStatus"`
	IsDirty                bool    `json:"isDirty"`
	LastSyncAttemptAt      *string `json:"lastSyncAttemptAt,omitempty"`
	LastSuccessfulSyncAt   *string `json:"lastSuccessfulSyncAt,omitempty"`
	ErrorMessage           *string `json:"errorMessage,omitempty"`
	CreatedAt              string  `json:"createdAt"`
	UpdatedAt              string  `json:"updatedAt"`
}

const (
	sqlCreateHeatmapLocalTracking = `CREATE TABLE IF NOT EXISTS HeatmapLocalTracking (
    HeatmapLocalTrackingId INTEGER PRIMARY KEY AUTOINCREMENT, RepoSlug TEXT NOT NULL,
    AiName TEXT NOT NULL, SequencePrefix TEXT NOT NULL DEFAULT '01', RepoCachePath TEXT NOT NULL DEFAULT '',
    LastCommittedSha TEXT NOT NULL DEFAULT '', LastCommitTimestamp TEXT NOT NULL DEFAULT '1970-01-01T00:00:00Z',
    UncommittedChangeCount INTEGER NOT NULL DEFAULT 0, SyncStatus TEXT NOT NULL DEFAULT 'IDLE',
    IsDirty INTEGER NOT NULL DEFAULT 0, LastSyncAttemptAt TEXT NULL, LastSuccessfulSyncAt TEXT NULL,
    ErrorMessage TEXT NULL, CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, UpdatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT UqHeatmapLocalTracking_Repo_AI UNIQUE (RepoSlug, AiName)
);
CREATE UNIQUE INDEX IF NOT EXISTS IdxHeatmapLocalTracking_Repo_AI ON HeatmapLocalTracking(RepoSlug, AiName);
CREATE INDEX IF NOT EXISTS IdxHeatmapLocalTracking_SyncStatus ON HeatmapLocalTracking(SyncStatus);
CREATE INDEX IF NOT EXISTS IdxHeatmapLocalTracking_IsDirty ON HeatmapLocalTracking(IsDirty);`

	sqlCreateHeatmapCommitLog = `CREATE TABLE IF NOT EXISTS HeatmapCommitLog (
    HeatmapCommitLogId INTEGER PRIMARY KEY AUTOINCREMENT, TaskSlug TEXT NOT NULL, RepoSlug TEXT NOT NULL,
    CommitSha TEXT NOT NULL, CommitMessage TEXT NOT NULL, FilesCommittedCount INTEGER NOT NULL DEFAULT 0,
    TransactionsCommittedCount INTEGER NOT NULL DEFAULT 0, CommittedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PushedSuccessfully INTEGER NOT NULL DEFAULT 1, RemoteOrigin TEXT NOT NULL DEFAULT '',
    Notes TEXT NULL, Comments TEXT NULL, CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS IdxHeatmapCommitLog_TaskSlug ON HeatmapCommitLog(TaskSlug);
CREATE INDEX IF NOT EXISTS IdxHeatmapCommitLog_CommitSha ON HeatmapCommitLog(CommitSha);`
)

// EnsureLocalTrackingTable creates the tracking table and indices if absent.
func EnsureLocalTrackingTable(conn *sql.DB) error {
	if conn == nil {
		return apperror.NewSimple("database connection is nil", "E_NIL_CONN")
	}
	if _, err := conn.Exec(sqlCreateHeatmapLocalTracking); err != nil {
		return apperror.WrapSimple(err, "heatmap_sync.ensureLocalTracking")
	}
	return nil
}

// OpenLocalTrackingDB opens or creates the local machine tracking SQLite database.
func OpenLocalTrackingDB() (*sql.DB, error) {
	dataDir := store.EffectiveDataDir()
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, apperror.WrapSimple(err, "heatmap_sync.mkdir")
	}
	openRes := store.OpenSQLiteDB(filepath.Join(dataDir, "heatmap.db"))
	if openRes.IsFailure() {
		return nil, openRes.AppError()
	}
	conn := openRes.Value
	if err := EnsureLocalTrackingTable(conn); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return conn, nil
}

var aiSequencePrefixes = map[string]string{
	"antigravity": "01",
	"muse":        "02",
	"cloud-code":  "03",
	"cloudcode":   "03",
	"codex":       "04",
	"grok":        "05",
	"cursor":      "06",
}

// ResolveAiSequencePrefix returns the canonical two-digit prefix for an AI agent.
func ResolveAiSequencePrefix(aiName string) string {
	if prefix, ok := aiSequencePrefixes[strings.ToLower(strings.TrimSpace(aiName))]; ok {
		return prefix
	}

	return "01"
}

// ResolveAiFolderName formats the AI folder name with its sequence prefix.
func ResolveAiFolderName(aiName string) string {
	clean := strings.ToLower(strings.TrimSpace(aiName))
	if clean == "" {
		clean = "antigravity"
	}
	return ResolveAiSequencePrefix(clean) + "-" + clean
}

// RecordLocalHeatmapAction updates local tracking when changes occur.
func RecordLocalHeatmapAction(repoSlug, aiName string, count int) error {
	conn, err := OpenLocalTrackingDB()
	if err != nil {
		return err
	}
	defer conn.Close()
	prefix := ResolveAiSequencePrefix(aiName)
	query := `INSERT INTO HeatmapLocalTracking (RepoSlug, AiName, SequencePrefix, UncommittedChangeCount, SyncStatus, IsDirty, UpdatedAt)
		VALUES (?, ?, ?, ?, 'DIRTY', 1, CURRENT_TIMESTAMP)
		ON CONFLICT(RepoSlug, AiName) DO UPDATE SET UncommittedChangeCount = UncommittedChangeCount + excluded.UncommittedChangeCount,
		SyncStatus = 'DIRTY', IsDirty = 1, UpdatedAt = CURRENT_TIMESTAMP`
	if _, execErr := conn.Exec(query, repoSlug, aiName, prefix, count); execErr != nil {
		return apperror.WrapSimple(execErr, "heatmap_sync.recordLocalAction")
	}
	return nil
}

// IncrementUncommittedChanges is an alias for RecordLocalHeatmapAction.
func IncrementUncommittedChanges(repoSlug, aiName string, count int) error {
	return RecordLocalHeatmapAction(repoSlug, aiName, count)
}

// GetLocalTrackingStatus retrieves current tracking status or seeds default.
func GetLocalTrackingStatus(repoSlug, aiName string) (*HeatmapLocalTrackingRecord, error) {
	conn, err := OpenLocalTrackingDB()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	rec, err := scanTrackingRow(conn, repoSlug, aiName)
	if err == nil {
		return rec, nil
	}
	if err != sql.ErrNoRows {
		return nil, apperror.WrapSimple(err, "heatmap_sync.getStatus")
	}
	return seedDefaultTracking(conn, repoSlug, aiName)
}

// GetOrCreateLocalTracking is an alias for GetLocalTrackingStatus.
func GetOrCreateLocalTracking(repoSlug, aiName string) (*HeatmapLocalTrackingRecord, error) {
	return GetLocalTrackingStatus(repoSlug, aiName)
}

func scanTrackingRow(conn *sql.DB, repoSlug, aiName string) (*HeatmapLocalTrackingRecord, error) {
	query := `SELECT HeatmapLocalTrackingId, RepoSlug, AiName, SequencePrefix, RepoCachePath,
		LastCommittedSha, LastCommitTimestamp, UncommittedChangeCount, SyncStatus, IsDirty,
		LastSyncAttemptAt, LastSuccessfulSyncAt, ErrorMessage, CreatedAt, UpdatedAt
		FROM HeatmapLocalTracking WHERE RepoSlug = ? AND AiName = ?`
	var r HeatmapLocalTrackingRecord
	var isDirtyInt int
	err := conn.QueryRow(query, repoSlug, aiName).Scan(
		&r.HeatmapLocalTrackingId, &r.RepoSlug, &r.AiName, &r.SequencePrefix, &r.RepoCachePath,
		&r.LastCommittedSha, &r.LastCommitTimestamp, &r.UncommittedChangeCount, &r.SyncStatus,
		&isDirtyInt, &r.LastSyncAttemptAt, &r.LastSuccessfulSyncAt, &r.ErrorMessage, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	r.IsDirty = isDirtyInt == 1
	return &r, nil
}

func seedDefaultTracking(conn *sql.DB, repoSlug, aiName string) (*HeatmapLocalTrackingRecord, error) {
	prefix := ResolveAiSequencePrefix(aiName)
	query := `INSERT OR IGNORE INTO HeatmapLocalTracking (RepoSlug, AiName, SequencePrefix, LastCommittedSha, LastCommitTimestamp,
		UncommittedChangeCount, SyncStatus, IsDirty) VALUES (?, ?, ?, '', '1970-01-01T00:00:00Z', 0, 'IDLE', 0)`
	if _, err := conn.Exec(query, repoSlug, aiName, prefix); err != nil {
		return nil, apperror.WrapSimple(err, "heatmap_sync.seedDefault")
	}
	return scanTrackingRow(conn, repoSlug, aiName)
}

// UpdateLocalTrackingCommit records a completed commit SHA into local tracking.
func UpdateLocalTrackingCommit(repoSlug, aiName, commitSha string) error {
	conn, err := OpenLocalTrackingDB()
	if err != nil {
		return err
	}
	defer conn.Close()
	now := time.Now().UTC().Format(time.RFC3339)
	query := `UPDATE HeatmapLocalTracking SET LastCommittedSha = ?, LastCommitTimestamp = ?,
		UncommittedChangeCount = 0, IsDirty = 0, SyncStatus = 'COMMITTED', LastSuccessfulSyncAt = ?, UpdatedAt = ?
		WHERE RepoSlug = ? AND AiName = ?`
	if _, execErr := conn.Exec(query, commitSha, now, now, now, repoSlug, aiName); execErr != nil {
		return apperror.WrapSimple(execErr, "heatmap_sync.updateCommit")
	}
	return nil
}

// SetSyncStatus updates sync status and dirty flag for a repository.
func SetSyncStatus(repoSlug, aiName, status string, isDirty bool) error {
	conn, err := OpenLocalTrackingDB()
	if err != nil {
		return err
	}
	defer conn.Close()
	isDirtyInt := 0
	if isDirty {
		isDirtyInt = 1
	}
	now := time.Now().UTC().Format(time.RFC3339)
	query := `UPDATE HeatmapLocalTracking SET SyncStatus = ?, IsDirty = ?, LastSyncAttemptAt = ?, UpdatedAt = ?
		WHERE RepoSlug = ? AND AiName = ?`
	if _, execErr := conn.Exec(query, status, isDirtyInt, now, now, repoSlug, aiName); execErr != nil {
		return apperror.WrapSimple(execErr, "heatmap_sync.setSyncStatus")
	}
	return nil
}

// CommitRcHeatmap executes semantic feature commit and records log entry.
func CommitRcHeatmap(repoCachePath, repoSlug, aiName, taskSummary string) (string, error) {
	if repoCachePath == "" {
		return "", apperror.NewValidation("heatmap_sync", "E1001", "repoCachePath is required")
	}
	commitMsg := fmt.Sprintf("Heatmap - %s - %s", repoSlug, taskSummary)
	if err := executeRcStageAndCommit(repoCachePath, commitMsg); err != nil {
		return "", err
	}
	newSha, err := readGitHeadSha(repoCachePath)
	if err != nil {
		return "", err
	}
	_ = logCommitToRootDb(repoCachePath, aiName, repoSlug, taskSummary, newSha, commitMsg)
	_ = UpdateLocalTrackingCommit(repoSlug, aiName, newSha)
	return newSha, nil
}

func executeRcStageAndCommit(repoDir, commitMsg string) error {
	addCmd := exec.Command("git", "add", ".heatmap/")
	addCmd.Dir = repoDir
	_ = addCmd.Run()
	cpfCmd := exec.Command("gitmap", "cpf", commitMsg)
	cpfCmd.Dir = repoDir
	if err := cpfCmd.Run(); err == nil {
		return nil
	}
	commitCmd := exec.Command("git", "commit", "-m", "Feature: "+commitMsg)
	commitCmd.Dir = repoDir
	if out, err := commitCmd.CombinedOutput(); err != nil && !strings.Contains(string(out), "nothing to commit") {
		return apperror.WrapSimple(err, "heatmap_sync.commit: "+string(out))
	}
	pushCmd := exec.Command("git", "push")
	pushCmd.Dir = repoDir
	_ = pushCmd.Run()
	return nil
}

func readGitHeadSha(repoDir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = repoDir
	out, err := cmd.Output()
	if err != nil {
		return "", apperror.WrapSimple(err, "heatmap_sync.readSha")
	}
	return strings.TrimSpace(string(out)), nil
}

func logCommitToRootDb(repoCachePath, aiName, repoSlug, taskSlug, commitSha, commitMsg string) error {
	aiFolder := ResolveAiFolderName(aiName)
	rootDbPath := filepath.Join(repoCachePath, ".heatmap", aiFolder, "root-heatmap.db")
	_ = os.MkdirAll(filepath.Dir(rootDbPath), 0755)
	openRes := store.OpenSQLiteDB(rootDbPath)
	if openRes.IsFailure() {
		return openRes.AppError()
	}
	conn := openRes.Value
	defer conn.Close()
	if _, err := conn.Exec(sqlCreateHeatmapCommitLog); err != nil {
		return apperror.WrapSimple(err, "heatmap_sync.ensureCommitLogTable")
	}
	query := `INSERT INTO HeatmapCommitLog (TaskSlug, RepoSlug, CommitSha, CommitMessage,
		FilesCommittedCount, TransactionsCommittedCount, CommittedAt, PushedSuccessfully, RemoteOrigin)
		VALUES (?, ?, ?, ?, 1, 1, CURRENT_TIMESTAMP, 1, '')`
	if _, err := conn.Exec(query, taskSlug, repoSlug, commitSha, commitMsg); err != nil {
		return apperror.WrapSimple(err, "heatmap_sync.insertCommitLog")
	}
	return nil
}
