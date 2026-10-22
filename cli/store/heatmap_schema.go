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
    Notes              TEXT NULL,
    Comments           TEXT NULL,
    CreatedAt          TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ParentLogId) REFERENCES HeatmapActionLog(HeatmapActionLogId) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS IdxHeatmapActionLog_TaskSlug ON HeatmapActionLog(TaskSlug);
CREATE INDEX IF NOT EXISTS IdxHeatmapActionLog_FilePath ON HeatmapActionLog(FilePath);
CREATE INDEX IF NOT EXISTS IdxHeatmapActionLog_ParentId ON HeatmapActionLog(ParentLogId);
CREATE INDEX IF NOT EXISTS IdxHeatmapActionLog_ActionType ON HeatmapActionLog(ActionType);
CREATE INDEX IF NOT EXISTS IdxHeatmapActionLog_CreatedAt ON HeatmapActionLog(CreatedAt);`
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

// EnsureHeatmapTables ensures the heatmap parent task, reasoning state, and action log tables exist.
func EnsureHeatmapTables(conn *sql.DB) error {
	if conn == nil {
		return apperror.NewSimple("database connection is nil", "E_NIL_CONN")
	}

	statements := []string{
		sqlCreateHeatmapParentTask,
		sqlCreateHeatmapReasoningState,
		sqlSeedHeatmapReasoningStates,
		sqlCreateHeatmapActionLog,
	}

	for _, stmt := range statements {
		if strings.TrimSpace(stmt) == "" {
			continue
		}
		if _, err := conn.Exec(stmt); err != nil {
			return apperror.WrapSimple(err, "store.ensureHeatmapTables")
		}
	}

	return nil
}

func scanActionLogs(rows *sql.Rows) ([]HeatmapActionLogRecord, error) {
	var records []HeatmapActionLogRecord
	for rows.Next() {
		var r HeatmapActionLogRecord
		var parentID sql.NullInt64
		var isMet int
		if err := rows.Scan(
			&r.HeatmapActionLogId, &parentID, &r.TaskSlug, &r.RepoName, &r.ActionType,
			&r.FilePath, &r.StartLine, &r.EndLine, &r.CodeSnippet, &r.CommandText,
			&r.ReasoningState, &r.ReasoningNotes, &r.ExpectedOutcome,
			&r.ActualOutcome, &isMet, &r.ExecutionDepth, &r.ExecutionOrder,
			&r.Notes, &r.Comments, &r.CreatedAt,
		); err != nil {
			return nil, apperror.WrapSimple(err, "heatmap.scanRow")
		}
		if parentID.Valid {
			pid := parentID.Int64
			r.ParentLogId = &pid
		}
		r.IsExpectationMet = isMet == 1
		records = append(records, r)
	}
	return records, nil
}

func tryAttachNodeToParent(node *HeatmapExecutionTreeNode, parentID *int64, nodeMap map[int64]*HeatmapExecutionTreeNode) bool {
	hasParent := parentID != nil && *parentID > 0
	if !hasParent {
		return false
	}

	parent, hasNode := nodeMap[*parentID]
	if !hasNode {
		return false
	}

	parent.Children = append(parent.Children, *node)

	return true
}

func assembleTreeNodes(records []HeatmapActionLogRecord) []HeatmapExecutionTreeNode {
	nodeMap := make(map[int64]*HeatmapExecutionTreeNode, len(records))
	var roots []HeatmapExecutionTreeNode

	for i := range records {
		rec := records[i]
		nodeMap[rec.HeatmapActionLogId] = &HeatmapExecutionTreeNode{
			HeatmapActionLogRecord: rec,
			Children:               make([]HeatmapExecutionTreeNode, 0),
		}
	}

	for _, rec := range records {
		node := nodeMap[rec.HeatmapActionLogId]
		hasAttached := tryAttachNodeToParent(node, rec.ParentLogId, nodeMap)
		if hasAttached {
			continue
		}

		roots = append(roots, *node)
	}

	return roots
}
