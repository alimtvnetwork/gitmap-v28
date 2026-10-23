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

// QueryActionLogs executes a read query and returns scanned HeatmapActionLogRecord records.
func (h *HeatmapDB) QueryActionLogs(query string, args ...any) ([]HeatmapActionLogRecord, error) {
	rows, err := h.conn.Query(query, args...)
	if err != nil {
		return nil, apperror.WrapSimple(err, "heatmap.queryActionLogs")
	}

	defer rows.Close()

	return scanActionLogs(rows)
}

// QueryHeatmapActionLogs executes a read query on any sql.DB and returns records.
func QueryHeatmapActionLogs(db *sql.DB, query string, args ...any) ([]HeatmapActionLogRecord, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, apperror.WrapSimple(err, "heatmap.queryActionLogs")
	}

	defer rows.Close()

	return scanActionLogs(rows)
}

type actionLogScanTargets struct {
	record   HeatmapActionLogRecord
	parentID sql.NullInt64
	isMet    int
	code     sql.NullString
	cmd      sql.NullString
	reason   sql.NullString
	exp      sql.NullString
	act      sql.NullString
	model    sql.NullString
	tool     sql.NullString
	notes    sql.NullString
	comments sql.NullString
}

func (t *actionLogScanTargets) bindCoreTarget(col string) any {
	switch strings.ToLower(col) {
	case "heatmapactionlogid": return &t.record.HeatmapActionLogId
	case "parentlogid": return &t.parentID
	case "taskslug": return &t.record.TaskSlug
	case "reponame", "reposlug": return &t.record.RepoName
	case "actiontype": return &t.record.ActionType
	case "filepath": return &t.record.FilePath
	case "startline": return &t.record.StartLine
	case "endline": return &t.record.EndLine
	default: return nil
	}
}

func (t *actionLogScanTargets) bindDetailsTarget(col string) any {
	switch strings.ToLower(col) {
	case "codesnippet": return &t.code
	case "commandtext": return &t.cmd
	case "reasoningstate": return &t.record.ReasoningState
	case "reasoningnotes": return &t.reason
	case "expectedoutcome": return &t.exp
	case "actualoutcome": return &t.act
	case "isexpectationmet": return &t.isMet
	default: return nil
	}
}

func (t *actionLogScanTargets) bindMetaTarget(col string) any {
	switch strings.ToLower(col) {
	case "executiondepth": return &t.record.ExecutionDepth
	case "executionorder": return &t.record.ExecutionOrder
	case "modelname": return &t.model
	case "toolname": return &t.tool
	case "notes": return &t.notes
	case "comments": return &t.comments
	case "createdat": return &t.record.CreatedAt
	default: return nil
	}
}

func (t *actionLogScanTargets) bindDetailsOrMeta(col string) any {
	if p := t.bindDetailsTarget(col); p != nil {
		return p
	}

	return t.bindMetaTarget(col)
}

func (t *actionLogScanTargets) bindTarget(col string) any {
	if p := t.bindCoreTarget(col); p != nil {
		return p
	}

	if p := t.bindDetailsOrMeta(col); p != nil {
		return p
	}

	var discard any

	return &discard
}

func (t *actionLogScanTargets) finalizeStrings() {
	if t.code.Valid { t.record.CodeSnippet = t.code.String }
	if t.cmd.Valid { t.record.CommandText = t.cmd.String }
	if t.reason.Valid { t.record.ReasoningNotes = t.reason.String }
	if t.exp.Valid { t.record.ExpectedOutcome = t.exp.String }
	if t.act.Valid { t.record.ActualOutcome = t.act.String }
}

func (t *actionLogScanTargets) finalizeMeta() {
	if t.model.Valid { t.record.ModelName = t.model.String }
	if t.tool.Valid { t.record.ToolName = t.tool.String }
	if t.notes.Valid { t.record.Notes = t.notes.String }
	if t.comments.Valid { t.record.Comments = t.comments.String }
}

func (t *actionLogScanTargets) finalize() HeatmapActionLogRecord {
	if t.parentID.Valid {
		pid := t.parentID.Int64
		t.record.ParentLogId = &pid
	}

	t.record.IsExpectationMet = t.isMet == 1
	t.finalizeStrings()
	t.finalizeMeta()

	return t.record
}

func scanSingleActionLog(rows *sql.Rows, cols []string) (HeatmapActionLogRecord, error) {
	var targets actionLogScanTargets
	dest := make([]any, len(cols))
	for i, col := range cols {
		dest[i] = targets.bindTarget(col)
	}

	if err := rows.Scan(dest...); err != nil {
		return targets.record, apperror.WrapSimple(err, "heatmap.scanRow")
	}

	return targets.finalize(), nil
}

func scanActionLogsRows(rows *sql.Rows, cols []string) ([]HeatmapActionLogRecord, error) {
	var records []HeatmapActionLogRecord
	for rows.Next() {
		rec, err := scanSingleActionLog(rows, cols)
		if err != nil {
			return nil, err
		}

		records = append(records, rec)
	}

	return records, nil
}

func scanActionLogs(rows *sql.Rows) ([]HeatmapActionLogRecord, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, apperror.WrapSimple(err, "heatmap.scanColumns")
	}

	return scanActionLogsRows(rows, cols)
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

func buildTreeNodeMap(records []HeatmapActionLogRecord) map[int64]*HeatmapExecutionTreeNode {
	nodeMap := make(map[int64]*HeatmapExecutionTreeNode, len(records))
	for i := range records {
		rec := records[i]
		nodeMap[rec.HeatmapActionLogId] = &HeatmapExecutionTreeNode{
			HeatmapActionLogRecord: rec,
			Children:               make([]HeatmapExecutionTreeNode, 0),
		}
	}

	return nodeMap
}

func collectRootTreeNodes(records []HeatmapActionLogRecord, nodeMap map[int64]*HeatmapExecutionTreeNode) []HeatmapExecutionTreeNode {
	var roots []HeatmapExecutionTreeNode
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

func assembleTreeNodes(records []HeatmapActionLogRecord) []HeatmapExecutionTreeNode {
	nodeMap := buildTreeNodeMap(records)

	return collectRootTreeNodes(records, nodeMap)
}
