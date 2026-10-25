package store

import (
	"database/sql"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// SQLCreateAiConfusionTable defines the SQLite DDL for the AiConfusion entity.
const SQLCreateAiConfusionTable = `
CREATE TABLE IF NOT EXISTS AiConfusion (
    AiConfusionId     INTEGER PRIMARY KEY AUTOINCREMENT,
    ConfusionUuid     TEXT NOT NULL UNIQUE,
    AgentSlug         TEXT NOT NULL,
    FileNames         TEXT NOT NULL DEFAULT '',
    MdFiles           TEXT NOT NULL DEFAULT '',
    Spec              TEXT NOT NULL DEFAULT '',
    Task              TEXT NOT NULL DEFAULT '',
    Subtask           TEXT NOT NULL DEFAULT '',
    HeatmapInfo       TEXT NOT NULL DEFAULT '',
    LineNumbers       TEXT NOT NULL DEFAULT '',
    CodeSnippet       TEXT NOT NULL DEFAULT '',
    Reason            TEXT NOT NULL DEFAULT '',
    InitialIntention  TEXT NOT NULL DEFAULT '',
    ReasoningThinking TEXT NOT NULL DEFAULT '',
    ConfusionQuestion TEXT NOT NULL DEFAULT '',
    IsCompleted       INTEGER NOT NULL DEFAULT 0,
    CreatedAt         TEXT NOT NULL,
    UpdatedAt         TEXT NOT NULL
);`

// SQLCreateIdxAiConfusionUuid creates the unique index on ConfusionUuid.
const SQLCreateIdxAiConfusionUuid = `
CREATE UNIQUE INDEX IF NOT EXISTS IdxAiConfusion_Uuid ON AiConfusion(ConfusionUuid);`

// SQLCreateIdxAiConfusionAgentSlug creates the index on AgentSlug.
const SQLCreateIdxAiConfusionAgentSlug = `
CREATE INDEX IF NOT EXISTS IdxAiConfusion_AgentSlug ON AiConfusion(AgentSlug);`

// SQLCreateIdxAiConfusionTask creates the index on Task.
const SQLCreateIdxAiConfusionTask = `
CREATE INDEX IF NOT EXISTS IdxAiConfusion_Task ON AiConfusion(Task);`

// SQLCreateIdxAiConfusionIsCompleted creates the index on IsCompleted.
const SQLCreateIdxAiConfusionIsCompleted = `
CREATE INDEX IF NOT EXISTS IdxAiConfusion_IsCompleted ON AiConfusion(IsCompleted);`

var aiConfusionSchemaStatements = []string{
	SQLCreateAiConfusionTable,
	SQLCreateIdxAiConfusionUuid,
	SQLCreateIdxAiConfusionAgentSlug,
	SQLCreateIdxAiConfusionTask,
	SQLCreateIdxAiConfusionIsCompleted,
}

// InitAiConfusionSchema initializes the AiConfusion table and all performance indexes.
func InitAiConfusionSchema(db *sql.DB) *apperror.AppError {
	for _, stmt := range aiConfusionSchemaStatements {
		res := ExecWrapper(db, stmt)
		if res.IsFailure {
			return apperror.WrapSimple(res.Error, "init ai confusion schema")
		}
	}

	return nil
}

// ResolveAiConfusionDbPath returns the canonical path to the ai-confusion.db file.
func ResolveAiConfusionDbPath(repoCacheDir string) string {
	baseDir := repoCacheDir
	hasEmptyBase := len(strings.TrimSpace(baseDir)) == 0
	if hasEmptyBase {
		baseDir = ResolveRepoCacheRoot()
	}

	isDbFile := strings.HasSuffix(strings.ToLower(baseDir), ".db")
	if isDbFile {
		return filepath.ToSlash(baseDir)
	}

	return filepath.ToSlash(filepath.Join(baseDir, "ai-confusion.db"))
}
