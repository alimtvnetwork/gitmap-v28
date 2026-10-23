package cmdcommit

import (
	"database/sql"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func updateHeatmapTelemetrySafely(rec store.CommitHistoryRecord) {
	recordHeatmapActionSafely(rec)
	recordHeatmapCommitLogSafely(rec)
}

func recordHeatmapActionSafely(rec store.CommitHistoryRecord) {
	hdb, err := store.OpenHeatmapRootDB("")
	if err != nil || hdb == nil {
		return
	}
	defer hdb.Close()

	action := buildHeatmapActionRecord(rec)
	_, _ = hdb.InsertActionLog(action)
}

func buildHeatmapActionRecord(rec store.CommitHistoryRecord) store.HeatmapActionLogRecord {
	return store.HeatmapActionLogRecord{
		TaskSlug:         rec.TaskSlug,
		RepoName:         rec.RepoName,
		ActionType:       "COMMIT",
		CommandText:      "gitmap commit " + rec.ShortSha,
		ReasoningState:   "EXECUTION",
		ReasoningNotes:   rec.CommitMessage,
		ExpectedOutcome:  "Commit recorded in git history",
		ActualOutcome:    rec.CommitSha,
		IsExpectationMet: true,
		ModelName:        rec.AiModelName,
		ToolName:         rec.AiToolName,
	}
}

func recordHeatmapCommitLogSafely(rec store.CommitHistoryRecord) {
	dataDir := store.EffectiveDataDir()
	dbPath := filepath.Join(dataDir, "heatmap.db")
	conn, err := store.OpenSQLiteDB(dbPath)
	if err != nil || conn == nil {
		return
	}
	defer conn.Close()

	ensureHeatmapCommitLogTable(conn)
	insertHeatmapCommitLogRow(conn, rec)
}

func ensureHeatmapCommitLogTable(conn *sql.DB) {
	const sqlCreate = `CREATE TABLE IF NOT EXISTS HeatmapCommitLog (
    HeatmapCommitLogId INTEGER PRIMARY KEY AUTOINCREMENT, TaskSlug TEXT NOT NULL, RepoSlug TEXT NOT NULL,
    CommitSha TEXT NOT NULL, CommitMessage TEXT NOT NULL, FilesCommittedCount INTEGER NOT NULL DEFAULT 0,
    TransactionsCommittedCount INTEGER NOT NULL DEFAULT 0, CommittedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PushedSuccessfully INTEGER NOT NULL DEFAULT 1, RemoteOrigin TEXT NOT NULL DEFAULT '',
    Notes TEXT NULL, Comments TEXT NULL, CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP);`
	_, _ = conn.Exec(sqlCreate)
}

func insertHeatmapCommitLogRow(conn *sql.DB, rec store.CommitHistoryRecord) {
	query := `INSERT INTO HeatmapCommitLog (TaskSlug, RepoSlug, CommitSha, CommitMessage,
		FilesCommittedCount, TransactionsCommittedCount, CommittedAt, PushedSuccessfully, RemoteOrigin)
		VALUES (?, ?, ?, ?, ?, 1, CURRENT_TIMESTAMP, ?, '')`
	hasPushedInt := 0
	if rec.HasPushed {
		hasPushedInt = 1
	}
	_, _ = conn.Exec(query, rec.TaskSlug, rec.RepoSlug, rec.CommitSha, rec.CommitMessage, rec.FilesChanged, hasPushedInt)
}
