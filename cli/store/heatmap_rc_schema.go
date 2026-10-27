package store

import (
	"database/sql"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const (
	sqlCreateHeatmapAiSession = `CREATE TABLE IF NOT EXISTS HeatmapAiSession (
    HeatmapAiSessionId INTEGER PRIMARY KEY AUTOINCREMENT, SessionUuid TEXT NOT NULL UNIQUE, AiName TEXT NOT NULL, SequencePrefix TEXT NOT NULL,
    MachineHostname TEXT NOT NULL DEFAULT '', OsPlatform TEXT NOT NULL DEFAULT '', AgentRole TEXT NOT NULL DEFAULT 'ORCHESTRATOR',
    StartedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, CompletedAt TEXT NULL, IsActive INTEGER NOT NULL DEFAULT 1, HasCompleted INTEGER NOT NULL DEFAULT 0,
    TotalTasksExecuted INTEGER NOT NULL DEFAULT 0, TotalActionsLogged INTEGER NOT NULL DEFAULT 0, Notes TEXT NULL, Comments TEXT NULL,
    CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, UpdatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
); CREATE UNIQUE INDEX IF NOT EXISTS IdxHeatmapAiSession_Uuid ON HeatmapAiSession(SessionUuid); CREATE INDEX IF NOT EXISTS IdxHeatmapAiSession_AiName ON HeatmapAiSession(AiName); CREATE INDEX IF NOT EXISTS IdxHeatmapAiSession_IsActive ON HeatmapAiSession(IsActive);`

	sqlCreateHeatmapRepoRegistry = `CREATE TABLE IF NOT EXISTS HeatmapRepoRegistry (
    HeatmapRepoRegistryId INTEGER PRIMARY KEY AUTOINCREMENT, RepoSlug TEXT NOT NULL UNIQUE, RepoUrl TEXT NOT NULL DEFAULT '', DefaultBranch TEXT NOT NULL DEFAULT 'main',
    FirstSeenAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, LastActiveAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, IsActive INTEGER NOT NULL DEFAULT 1, TotalActionsRecorded INTEGER NOT NULL DEFAULT 0,
    Notes TEXT NULL, Comments TEXT NULL, CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, UpdatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
); CREATE UNIQUE INDEX IF NOT EXISTS IdxHeatmapRepoRegistry_Slug ON HeatmapRepoRegistry(RepoSlug); CREATE INDEX IF NOT EXISTS IdxHeatmapRepoRegistry_IsActive ON HeatmapRepoRegistry(IsActive);`

	sqlCreateRcHeatmapParentTask = `CREATE TABLE IF NOT EXISTS HeatmapParentTask (
    HeatmapParentTaskId INTEGER PRIMARY KEY AUTOINCREMENT, TaskSlug TEXT NOT NULL UNIQUE, TaskName TEXT NOT NULL, RepoSlug TEXT NOT NULL, Status TEXT NOT NULL DEFAULT 'IN_PROGRESS',
    RootSpecFilePath TEXT NOT NULL DEFAULT '', TaskFolder TEXT NOT NULL DEFAULT '',
    StepBudget INTEGER NOT NULL DEFAULT 300, CompletedSteps INTEGER NOT NULL DEFAULT 0, HasCompleted INTEGER NOT NULL DEFAULT 0, IsActive INTEGER NOT NULL DEFAULT 1,
    StartedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, CompletedAt TEXT NULL, Notes TEXT NULL, Comments TEXT NULL, CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, UpdatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
); CREATE UNIQUE INDEX IF NOT EXISTS IdxHeatmapParentTask_Slug ON HeatmapParentTask(TaskSlug); CREATE INDEX IF NOT EXISTS IdxHeatmapParentTask_RepoSlug ON HeatmapParentTask(RepoSlug); CREATE INDEX IF NOT EXISTS IdxHeatmapParentTask_Status ON HeatmapParentTask(Status); CREATE INDEX IF NOT EXISTS IdxHeatmapParentTask_IsActive ON HeatmapParentTask(IsActive);`

	sqlCreateHeatmapCommitLog = `CREATE TABLE IF NOT EXISTS HeatmapCommitLog (
    HeatmapCommitLogId INTEGER PRIMARY KEY AUTOINCREMENT, TaskSlug TEXT NOT NULL, RepoSlug TEXT NOT NULL, CommitSha TEXT NOT NULL, CommitMessage TEXT NOT NULL,
    FilesCommittedCount INTEGER NOT NULL DEFAULT 0, TransactionsCommittedCount INTEGER NOT NULL DEFAULT 0, CommittedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, PushedSuccessfully INTEGER NOT NULL DEFAULT 1,
    RemoteOrigin TEXT NOT NULL DEFAULT '', Notes TEXT NULL, Comments TEXT NULL, CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
); CREATE INDEX IF NOT EXISTS IdxHeatmapCommitLog_TaskSlug ON HeatmapCommitLog(TaskSlug); CREATE INDEX IF NOT EXISTS IdxHeatmapCommitLog_RepoSlug ON HeatmapCommitLog(RepoSlug); CREATE INDEX IF NOT EXISTS IdxHeatmapCommitLog_CommitSha ON HeatmapCommitLog(CommitSha); CREATE INDEX IF NOT EXISTS IdxHeatmapCommitLog_CommittedAt ON HeatmapCommitLog(CommittedAt);`

	sqlCreateRcHeatmapActionLog = `CREATE TABLE IF NOT EXISTS HeatmapActionLog (
    HeatmapActionLogId INTEGER PRIMARY KEY AUTOINCREMENT, ParentLogId INTEGER NULL, TaskSlug TEXT NOT NULL, RepoSlug TEXT NOT NULL DEFAULT '', ActionType TEXT NOT NULL,
    FilePath TEXT NOT NULL DEFAULT '', StartLine INTEGER NOT NULL DEFAULT 0, EndLine INTEGER NOT NULL DEFAULT 0, CodeSnippet TEXT NULL, CommandText TEXT NULL,
    ReasoningState TEXT NOT NULL DEFAULT 'EXPLORATION', ReasoningNotes TEXT NULL, ExpectedOutcome TEXT NULL, ActualOutcome TEXT NULL, IsExpectationMet INTEGER NOT NULL DEFAULT 1,
    ExecutionDepth INTEGER NOT NULL DEFAULT 0, ExecutionOrder INTEGER NOT NULL DEFAULT 1, DurationMs INTEGER NOT NULL DEFAULT 0,
    ModelName TEXT NULL, ToolName TEXT NULL, Notes TEXT NULL, Comments TEXT NULL, CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
); CREATE INDEX IF NOT EXISTS IdxHeatmapActionLog_TaskSlug ON HeatmapActionLog(TaskSlug); CREATE INDEX IF NOT EXISTS IdxHeatmapActionLog_FilePath ON HeatmapActionLog(FilePath); CREATE INDEX IF NOT EXISTS IdxHeatmapActionLog_ParentId ON HeatmapActionLog(ParentLogId); CREATE INDEX IF NOT EXISTS IdxHeatmapActionLog_ActionType ON HeatmapActionLog(ActionType); CREATE INDEX IF NOT EXISTS IdxHeatmapActionLog_CreatedAt ON HeatmapActionLog(CreatedAt);`

	sqlCreateHeatmapExecutionTree = `CREATE TABLE IF NOT EXISTS HeatmapExecutionTree (
    HeatmapExecutionTreeId INTEGER PRIMARY KEY AUTOINCREMENT, RootLogId INTEGER NOT NULL, NodeLogId INTEGER NOT NULL UNIQUE, ParentLogId INTEGER NULL,
    NodeDepth INTEGER NOT NULL DEFAULT 0, SiblingOrder INTEGER NOT NULL DEFAULT 1, PathVector TEXT NOT NULL DEFAULT '/', IsLeaf INTEGER NOT NULL DEFAULT 1, HasChildren INTEGER NOT NULL DEFAULT 0,
    CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, UpdatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
); CREATE UNIQUE INDEX IF NOT EXISTS IdxHeatmapExecutionTree_Node ON HeatmapExecutionTree(NodeLogId); CREATE INDEX IF NOT EXISTS IdxHeatmapExecutionTree_Root ON HeatmapExecutionTree(RootLogId); CREATE INDEX IF NOT EXISTS IdxHeatmapExecutionTree_Parent ON HeatmapExecutionTree(ParentLogId); CREATE INDEX IF NOT EXISTS IdxHeatmapExecutionTree_Path ON HeatmapExecutionTree(PathVector);`

	sqlCreateHeatmapLocalTracking = `CREATE TABLE IF NOT EXISTS HeatmapLocalTracking (
    HeatmapLocalTrackingId INTEGER PRIMARY KEY AUTOINCREMENT, RepoSlug TEXT NOT NULL, AiName TEXT NOT NULL, SequencePrefix TEXT NOT NULL DEFAULT '01',
    RepoCachePath TEXT NOT NULL DEFAULT '', LastCommittedSha TEXT NOT NULL DEFAULT '', LastCommitTimestamp TEXT NOT NULL DEFAULT '1970-01-01T00:00:00Z',
    UncommittedChangeCount INTEGER NOT NULL DEFAULT 0, SyncStatus TEXT NOT NULL DEFAULT 'IDLE', IsDirty INTEGER NOT NULL DEFAULT 0,
    LastSyncAttemptAt TEXT NULL, LastSuccessfulSyncAt TEXT NULL, ErrorMessage TEXT NULL, CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, UpdatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT UqHeatmapLocalTracking_Repo_AI UNIQUE (RepoSlug, AiName)
); CREATE UNIQUE INDEX IF NOT EXISTS IdxHeatmapLocalTracking_Repo_AI ON HeatmapLocalTracking(RepoSlug, AiName); CREATE INDEX IF NOT EXISTS IdxHeatmapLocalTracking_SyncStatus ON HeatmapLocalTracking(SyncStatus); CREATE INDEX IF NOT EXISTS IdxHeatmapLocalTracking_IsDirty ON HeatmapLocalTracking(IsDirty);`
)

// HeatmapAiSessionRecord maps to table HeatmapAiSession.
type HeatmapAiSessionRecord struct {
	HeatmapAiSessionId int64   `json:"heatmapAiSessionId"`
	SessionUuid        string  `json:"sessionUuid"`
	AiName             string  `json:"aiName"`
	SequencePrefix     string  `json:"sequencePrefix"`
	MachineHostname    string  `json:"machineHostname"`
	OsPlatform         string  `json:"osPlatform"`
	AgentRole          string  `json:"agentRole"`
	StartedAt          string  `json:"startedAt"`
	CompletedAt        *string `json:"completedAt,omitempty"`
	IsActive           bool    `json:"isActive"`
	HasCompleted       bool    `json:"hasCompleted"`
	TotalTasksExecuted int     `json:"totalTasksExecuted"`
	TotalActionsLogged int     `json:"totalActionsLogged"`
	Notes              string  `json:"notes,omitempty"`
	Comments           string  `json:"comments,omitempty"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt"`
}

// HeatmapRepoRegistryRecord maps to table HeatmapRepoRegistry.
type HeatmapRepoRegistryRecord struct {
	HeatmapRepoRegistryId int64  `json:"heatmapRepoRegistryId"`
	RepoSlug              string `json:"repoSlug"`
	RepoUrl               string `json:"repoUrl"`
	DefaultBranch         string `json:"defaultBranch"`
	FirstSeenAt           string `json:"firstSeenAt"`
	LastActiveAt          string `json:"lastActiveAt"`
	IsActive              bool   `json:"isActive"`
	TotalActionsRecorded  int    `json:"totalActionsRecorded"`
	Notes                 string `json:"notes,omitempty"`
	Comments              string `json:"comments,omitempty"`
	CreatedAt             string `json:"createdAt"`
	UpdatedAt             string `json:"updatedAt"`
}

// RcHeatmapParentTaskRecord maps to table HeatmapParentTask in root-heatmap.db.
type RcHeatmapParentTaskRecord struct {
	HeatmapParentTaskId int64   `json:"heatmapParentTaskId"`
	TaskSlug            string  `json:"taskSlug"`
	TaskName            string  `json:"taskName"`
	RepoSlug            string  `json:"repoSlug"`
	RootSpecFilePath    string  `json:"rootSpecFilePath,omitempty"`
	TaskFolder          string  `json:"taskFolder,omitempty"`
	Status              string  `json:"status"`
	StepBudget          int     `json:"stepBudget"`
	CompletedSteps      int     `json:"completedSteps"`
	HasCompleted        bool    `json:"hasCompleted"`
	IsActive            bool    `json:"isActive"`
	StartedAt           string  `json:"startedAt"`
	CompletedAt         *string `json:"completedAt,omitempty"`
	Notes               string  `json:"notes,omitempty"`
	Comments            string  `json:"comments,omitempty"`
	CreatedAt           string  `json:"createdAt"`
	UpdatedAt           string  `json:"updatedAt"`
}

// HeatmapCommitLogRecord maps to table HeatmapCommitLog.
type HeatmapCommitLogRecord struct {
	HeatmapCommitLogId         int64  `json:"heatmapCommitLogId"`
	TaskSlug                   string `json:"taskSlug"`
	RepoSlug                   string `json:"repoSlug"`
	CommitSha                  string `json:"commitSha"`
	CommitMessage              string `json:"commitMessage"`
	FilesCommittedCount        int    `json:"filesCommittedCount"`
	TransactionsCommittedCount int    `json:"transactionsCommittedCount"`
	CommittedAt                string `json:"committedAt"`
	PushedSuccessfully         bool   `json:"pushedSuccessfully"`
	RemoteOrigin               string `json:"remoteOrigin"`
	Notes                      string `json:"notes,omitempty"`
	Comments                   string `json:"comments,omitempty"`
	CreatedAt                  string `json:"createdAt"`
}

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

// RcHeatmapExecutionTreeNode represents a hierarchical step node with nested children.
type RcHeatmapExecutionTreeNode struct {
	HeatmapActionLogRecord
	Children []RcHeatmapExecutionTreeNode `json:"children,omitempty"`
}

func tryAttachRcNodeToParent(node *RcHeatmapExecutionTreeNode, parentID *int64, nodeMap map[int64]*RcHeatmapExecutionTreeNode) bool {
	if parentID == nil || *parentID <= 0 {
		return false
	}

	parent, hasParent := nodeMap[*parentID]
	if !hasParent {
		return false
	}

	parent.Children = append(parent.Children, *node)

	return true
}

// AssembleRcTreeNodes builds tree nodes from flat action records.
func AssembleRcTreeNodes(records []HeatmapActionLogRecord) []RcHeatmapExecutionTreeNode {
	nodeMap := make(map[int64]*RcHeatmapExecutionTreeNode, len(records))
	var roots []RcHeatmapExecutionTreeNode
	for i := range records {
		rec := records[i]
		nodeMap[rec.HeatmapActionLogId] = &RcHeatmapExecutionTreeNode{HeatmapActionLogRecord: rec, Children: make([]RcHeatmapExecutionTreeNode, 0)}
	}
	for _, rec := range records {
		node := nodeMap[rec.HeatmapActionLogId]
		if tryAttachRcNodeToParent(node, rec.ParentLogId, nodeMap) {
			continue
		}

		roots = append(roots, *node)
	}
	return roots
}

// ConfigureRcHeatmapPragmas applies WAL mode, busy timeout, and performance pragmas.
func ConfigureRcHeatmapPragmas(conn *sql.DB) error {
	if conn == nil {
		return apperror.NewSimple("connection is nil", "E_NIL_CONN")
	}
	pragmas := []string{
		"PRAGMA journal_mode = WAL;", "PRAGMA synchronous = NORMAL;", "PRAGMA busy_timeout = 5000;",
		"PRAGMA foreign_keys = ON;", "PRAGMA temp_store = MEMORY;", "PRAGMA cache_size = -64000;",
	}
	for _, p := range pragmas {
		if _, err := conn.Exec(p); err != nil {
			return apperror.WrapSimple(err, "heatmap_rc.pragma")
		}
	}
	return nil
}

// EnsureRcHeatmapRootTables creates all root-heatmap.db tables.
func EnsureRcHeatmapRootTables(conn *sql.DB) error {
	if err := executeSchemaBatch(conn, []string{
		sqlCreateHeatmapAiSession, sqlCreateHeatmapRepoRegistry,
		sqlCreateRcHeatmapParentTask, sqlCreateHeatmapCommitLog,
	}, "heatmap_rc.ensureRoot"); err != nil {
		return err
	}

	alterStatements := []string{
		"ALTER TABLE HeatmapParentTask ADD COLUMN RootSpecFilePath TEXT NOT NULL DEFAULT '';",
		"ALTER TABLE HeatmapParentTask ADD COLUMN TaskFolder TEXT NOT NULL DEFAULT '';",
		"CREATE INDEX IF NOT EXISTS IdxHeatmapParentTask_Spec ON HeatmapParentTask(RootSpecFilePath);",
	}
	for _, stmt := range alterStatements {
		_ = execHeatmapAlter(conn, stmt)
	}

	return nil
}

// EnsureRcHeatmapRepoTables creates all per-repo heatmap.db tables.
func EnsureRcHeatmapRepoTables(conn *sql.DB) error {
	if err := executeSchemaBatch(conn, []string{
		sqlCreateHeatmapReasoningState, sqlSeedHeatmapReasoningStates,
		sqlCreateRcHeatmapActionLog, sqlCreateHeatmapExecutionTree,
	}, "heatmap_rc.ensureRepo"); err != nil {
		return err
	}

	return ensureHeatmapActionLogColumns(conn)
}

// EnsureRcHeatmapLocalTrackingTables creates the HeatmapLocalTracking table in local tracking DB.
func EnsureRcHeatmapLocalTrackingTables(conn *sql.DB) error {
	return executeSchemaBatch(conn, []string{sqlCreateHeatmapLocalTracking}, "heatmap_rc.ensureLocal")
}

func executeSchemaBatch(conn *sql.DB, stmts []string, scope string) error {
	if conn == nil {
		return apperror.NewSimple("connection is nil", "E_NIL_CONN")
	}
	for _, stmt := range stmts {
		clean := strings.TrimSpace(stmt)
		if clean == "" {
			continue
		}
		if _, err := conn.Exec(clean); err != nil {
			return apperror.WrapSimple(err, scope)
		}
	}
	return nil
}

// ScanSessionRows reads session records from sql.Rows.
func ScanSessionRows(rows *sql.Rows) ([]HeatmapAiSessionRecord, error) {
	var sessions []HeatmapAiSessionRecord
	for rows.Next() {
		var r HeatmapAiSessionRecord
		var isActive, hasComp int
		if err := rows.Scan(&r.HeatmapAiSessionId, &r.SessionUuid, &r.AiName, &r.SequencePrefix, &r.MachineHostname, &r.OsPlatform, &r.AgentRole,
			&r.StartedAt, &r.CompletedAt, &isActive, &hasComp, &r.TotalTasksExecuted, &r.TotalActionsLogged, &r.Notes, &r.Comments, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, apperror.WrapSimple(err, "heatmap_rc.scanSession")
		}
		r.IsActive, r.HasCompleted = isActive == 1, hasComp == 1
		sessions = append(sessions, r)
	}
	return sessions, nil
}

// ScanRepoRegistryRows reads repo registry records from sql.Rows.
func ScanRepoRegistryRows(rows *sql.Rows) ([]HeatmapRepoRegistryRecord, error) {
	var repos []HeatmapRepoRegistryRecord
	for rows.Next() {
		var r HeatmapRepoRegistryRecord
		var isActive int
		if err := rows.Scan(&r.HeatmapRepoRegistryId, &r.RepoSlug, &r.RepoUrl, &r.DefaultBranch, &r.FirstSeenAt, &r.LastActiveAt,
			&isActive, &r.TotalActionsRecorded, &r.Notes, &r.Comments, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, apperror.WrapSimple(err, "heatmap_rc.scanRepo")
		}
		r.IsActive = isActive == 1
		repos = append(repos, r)
	}
	return repos, nil
}
