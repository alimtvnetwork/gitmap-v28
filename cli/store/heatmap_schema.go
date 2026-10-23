package store

import (
	"database/sql"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const (
	sqlCreateHeatmapParentTask = `CREATE TABLE IF NOT EXISTS HeatmapParentTask (
    HeatmapParentTaskId INTEGER PRIMARY KEY AUTOINCREMENT,
    TaskSlug            TEXT NOT NULL UNIQUE,
    TaskName            TEXT NOT NULL,
    RunDirectory        TEXT NOT NULL DEFAULT '',
    RootDbPath          TEXT NOT NULL DEFAULT '',
    Status              TEXT NOT NULL DEFAULT 'IN_PROGRESS',
    Notes               TEXT NULL,
    Comments            TEXT NULL,
    CreatedAt           TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt           TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS IdxHeatmapParentTask_Slug ON HeatmapParentTask(TaskSlug);`

	sqlCreateHeatmapReasoningState = `CREATE TABLE IF NOT EXISTS HeatmapReasoningState (
    HeatmapReasoningStateId INTEGER PRIMARY KEY AUTOINCREMENT,
    StateCode               TEXT NOT NULL UNIQUE,
    Description             TEXT NOT NULL,
    Notes                   TEXT NULL,
    IsActive                INTEGER NOT NULL DEFAULT 1,
    CreatedAt               TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS IdxHeatmapReasoningState_Code ON HeatmapReasoningState(StateCode);`

	sqlSeedHeatmapReasoningStates = `INSERT OR IGNORE INTO HeatmapReasoningState (StateCode, Description, Notes, IsActive) VALUES
('EXPLORATION', 'Initial exploration, file discovery, and codebase inspection', NULL, 1),
('SPEC_REFINEMENT', 'Spec review, requirement alignment, and task decomposition', NULL, 1),
('CODE_EDITING', 'Direct file modification, creation, and refactoring', NULL, 1),
('VERIFICATION', 'Reviewing diffs, linting checks, and quality validation', NULL, 1),
('REMEDIATION', 'Bug fixing, error resolution, and exception handling', NULL, 1),
('COMPLETED', 'Step or subtask completed successfully', NULL, 1);`

	sqlCreateHeatmapActionLog = `CREATE TABLE IF NOT EXISTS HeatmapActionLog (
    HeatmapActionLogId INTEGER PRIMARY KEY AUTOINCREMENT,
    ParentLogId        INTEGER NULL,
    TaskSlug           TEXT NOT NULL,
    RepoName           TEXT NOT NULL DEFAULT '',
    ActionType         TEXT NOT NULL,
    FilePath           TEXT NOT NULL DEFAULT '',
    StartLine          INTEGER NOT NULL DEFAULT 0,
    EndLine            INTEGER NOT NULL DEFAULT 0,
    CodeSnippet        TEXT NULL,
    CommandText        TEXT NULL,
    ReasoningState     TEXT NOT NULL DEFAULT 'EXPLORATION',
    ReasoningNotes     TEXT NULL,
    ExpectedOutcome    TEXT NULL,
    ActualOutcome      TEXT NULL,
    IsExpectationMet   INTEGER NOT NULL DEFAULT 1,
    ExecutionDepth     INTEGER NOT NULL DEFAULT 0,
    ExecutionOrder     INTEGER NOT NULL DEFAULT 1,
    ModelName          TEXT NULL,
    ToolName           TEXT NULL,
    Notes              TEXT NULL,
    Comments           TEXT NULL,
    CreatedAt          TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ParentLogId) REFERENCES HeatmapActionLog(HeatmapActionLogId) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS IdxHeatmapActionLog_TaskSlug ON HeatmapActionLog(TaskSlug);
CREATE INDEX IF NOT EXISTS IdxHeatmapActionLog_FilePath ON HeatmapActionLog(FilePath);
CREATE INDEX IF NOT EXISTS IdxHeatmapActionLog_ParentId ON HeatmapActionLog(ParentLogId);
CREATE INDEX IF NOT EXISTS IdxHeatmapActionLog_ActionType ON HeatmapActionLog(ActionType);
CREATE INDEX IF NOT EXISTS IdxHeatmapActionLog_CreatedAt ON HeatmapActionLog(CreatedAt);
CREATE INDEX IF NOT EXISTS IdxHeatmapActionLog_ModelTool ON HeatmapActionLog(ModelName, ToolName);`

	sqlInsertActionLogWithAttribution = `INSERT INTO HeatmapActionLog (
		ParentLogId, TaskSlug, RepoName, ActionType, FilePath, StartLine, EndLine,
		CodeSnippet, CommandText, ReasoningState, ReasoningNotes, ExpectedOutcome,
		ActualOutcome, IsExpectationMet, ExecutionDepth, ExecutionOrder,
		ModelName, ToolName, Notes, Comments
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
)

// HeatmapParentTaskRecord represents a parent task tracked in the heatmap.
type HeatmapParentTaskRecord struct {
	HeatmapParentTaskId int64  `json:"heatmapParentTaskId"`
	TaskSlug            string `json:"taskSlug"`
	TaskName            string `json:"taskName"`
	RunDirectory        string `json:"runDirectory"`
	RootDbPath          string `json:"rootDbPath"`
	Status              string `json:"status"`
	Notes               string `json:"notes,omitempty"`
	Comments            string `json:"comments,omitempty"`
	CreatedAt           string `json:"createdAt"`
	UpdatedAt           string `json:"updatedAt"`
}

// HeatmapReasoningStateRecord represents a reasoning state entry.
type HeatmapReasoningStateRecord struct {
	HeatmapReasoningStateId int64  `json:"heatmapReasoningStateId"`
	StateCode               string `json:"stateCode"`
	Description             string `json:"description"`
	Notes                   string `json:"notes,omitempty"`
	IsActive                bool   `json:"isActive"`
	CreatedAt               string `json:"createdAt"`
}

// HeatmapActionLogRecord represents an individual action log entry.
type HeatmapActionLogRecord struct {
	HeatmapActionLogId int64  `json:"heatmapActionLogId"`
	ParentLogId        *int64 `json:"parentLogId,omitempty"`
	TaskSlug           string `json:"taskSlug"`
	RepoName           string `json:"repoName"`
	ActionType         string `json:"actionType"`
	FilePath           string `json:"filePath"`
	StartLine          int    `json:"startLine"`
	EndLine            int    `json:"endLine"`
	CodeSnippet        string `json:"codeSnippet,omitempty"`
	CommandText        string `json:"commandText,omitempty"`
	ReasoningState     string `json:"reasoningState"`
	ReasoningNotes     string `json:"reasoningNotes,omitempty"`
	ExpectedOutcome    string `json:"expectedOutcome,omitempty"`
	ActualOutcome      string `json:"actualOutcome,omitempty"`
	IsExpectationMet   bool   `json:"isExpectationMet"`
	ExecutionDepth     int    `json:"executionDepth"`
	ExecutionOrder     int    `json:"executionOrder"`
	ModelName          string `json:"modelName,omitempty"`
	ToolName           string `json:"toolName,omitempty"`
	Notes              string `json:"notes,omitempty"`
	Comments           string `json:"comments,omitempty"`
	CreatedAt          string `json:"createdAt"`
}

// HeatmapExecutionTreeNode represents a hierarchical step node with nested children.
type HeatmapExecutionTreeNode struct {
	HeatmapActionLogRecord
	Children []HeatmapExecutionTreeNode `json:"children,omitempty"`
}

// HeatmapStatusSummary provides an aggregated summary of heatmap actions.
type HeatmapStatusSummary struct {
	TaskSlug             string           `json:"taskSlug"`
	TotalActions         int              `json:"totalActions"`
	ExpectationsMetCount int              `json:"expectationsMetCount"`
	ExpectationsMetPct   float64          `json:"expectationsMetPct"`
	ReasoningStateCounts map[string]int   `json:"reasoningStateCounts"`
	TopTouchedFiles      []FileTouchCount `json:"topTouchedFiles"`
}

// HeatmapTaskSummary aggregates overall action volume and categorizes actions by operation touch type for a task.
type HeatmapTaskSummary struct {
	TaskSlug     string `json:"taskSlug"`
	TotalActions int    `json:"totalActions"`
	ReadCount    int    `json:"readCount"`
	EditCount    int    `json:"editCount"`
	CreateCount  int    `json:"createCount"`
	OtherCount   int    `json:"otherCount"`
}

// FileTouchType represents the categorical operation applied to a touched file.
type FileTouchType string

const (
	FileTouchRead    FileTouchType = "READ"
	FileTouchEdit    FileTouchType = "EDIT"
	FileTouchCreated FileTouchType = "CREATED"
)

// IsValid reports whether the touch type is a recognized enumeration value.
func (t FileTouchType) IsValid() bool {
	switch t {
	case FileTouchRead, FileTouchEdit, FileTouchCreated:
		return true
	default:
		return false
	}
}

// FileTouchCount represents file touch frequency with granular operation metrics.
type FileTouchCount struct {
	FilePath     string `json:"filePath"`
	Count        int    `json:"count"`
	ReadCount    int    `json:"readCount,omitempty"`
	EditCount    int    `json:"editCount,omitempty"`
	CreatedCount int    `json:"createdCount,omitempty"`
}

// InitHeatmapSchema ensures all heatmap tables and migrated columns exist.
func InitHeatmapSchema(conn *sql.DB) error {
	return EnsureHeatmapTables(conn)
}

func execHeatmapAlter(conn *sql.DB, stmt string) error {
	_, err := conn.Exec(stmt)
	if err != nil && !isBenignAlterError(err) {
		return apperror.WrapSimple(err, "heatmap.alterColumn")
	}

	return nil
}

func ensureHeatmapActionLogColumns(conn *sql.DB) error {
	alterStatements := []string{
		"ALTER TABLE HeatmapActionLog ADD COLUMN ModelName TEXT NULL;",
		"ALTER TABLE HeatmapActionLog ADD COLUMN ToolName TEXT NULL;",
	}

	for _, stmt := range alterStatements {
		if err := execHeatmapAlter(conn, stmt); err != nil {
			return err
		}
	}

	return nil
}

func executeHeatmapStatements(conn *sql.DB, stmts []string) error {
	for _, stmt := range stmts {
		if strings.TrimSpace(stmt) == "" {
			continue
		}

		if _, err := conn.Exec(stmt); err != nil {
			return apperror.WrapSimple(err, "store.execHeatmapDDL")
		}
	}

	return nil
}

func buildHeatmapTableStatements() []string {
	return []string{
		sqlCreateHeatmapParentTask, sqlCreateHeatmapReasoningState,
		sqlSeedHeatmapReasoningStates, sqlCreateHeatmapActionLog,
	}
}

// EnsureHeatmapTables ensures the heatmap parent task, reasoning state, and action log tables exist.
func EnsureHeatmapTables(conn *sql.DB) error {
	if conn == nil {
		return apperror.NewSimple("database connection is nil", "E_NIL_CONN")
	}

	if err := executeHeatmapStatements(conn, buildHeatmapTableStatements()); err != nil {
		return err
	}

	return ensureHeatmapActionLogColumns(conn)
}

func resolveMetInt(isMet bool) int {
	if isMet {
		return 1
	}

	return 0
}

func (h *HeatmapDB) execInsertActionLog(log HeatmapActionLogRecord, depth, order int) (int64, error) {
	res, err := h.conn.Exec(sqlInsertActionLogWithAttribution,
		log.ParentLogId, log.TaskSlug, log.RepoName, log.ActionType, log.FilePath,
		log.StartLine, log.EndLine, log.CodeSnippet, log.CommandText, log.ReasoningState,
		log.ReasoningNotes, log.ExpectedOutcome, log.ActualOutcome, resolveMetInt(log.IsExpectationMet),
		depth, order, log.ModelName, log.ToolName, log.Notes, log.Comments,
	)
	if err != nil {
		return 0, apperror.WrapSimple(err, "heatmap.insertActionLog")
	}

	return res.LastInsertId()
}

// InsertActionLog inserts an action record with full provenance attribution.
func (h *HeatmapDB) InsertActionLog(log HeatmapActionLogRecord) (int64, error) {
	depth, order, err := h.calculateDepthAndOrder(log.ParentLogId, log.TaskSlug)
	if err != nil {
		return 0, err
	}

	return h.execInsertActionLog(log, depth, order)
}
