package store

import (
	"database/sql"
	"path/filepath"
	"strings"
	"time"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/types"
	_ "modernc.org/sqlite"
)

func queryParentTask(conn *sql.DB, parentId string) types.ParentTask {
	var p types.ParentTask
	query := "SELECT ParentTaskId, TaskSlug, TaskName, RunDirectory, RootDbPath, Status, TotalStepsBudget, CompletedSteps, SpawnedAgentCount, CreatedAt, UpdatedAt FROM ParentTask LIMIT 1"
	hasParent := len(parentId) > 0
	if hasParent {
		query = "SELECT ParentTaskId, TaskSlug, TaskName, RunDirectory, RootDbPath, Status, TotalStepsBudget, CompletedSteps, SpawnedAgentCount, CreatedAt, UpdatedAt FROM ParentTask WHERE ParentTaskId = ? LIMIT 1"
		row := QueryRowWrapper(conn, query, parentId)
		_ = row.Scan(&p.ParentTaskId, &p.TaskSlug, &p.TaskName, &p.RunDirectory, &p.RootDbPath, &p.Status, &p.TotalStepsBudget, &p.CompletedSteps, &p.SpawnedAgentCount, &p.CreatedAt, &p.UpdatedAt) // lint-allow: ignore-db-error
		return p
	}
	row := QueryRowWrapper(conn, query)
	_ = row.Scan(&p.ParentTaskId, &p.TaskSlug, &p.TaskName, &p.RunDirectory, &p.RootDbPath, &p.Status, &p.TotalStepsBudget, &p.CompletedSteps, &p.SpawnedAgentCount, &p.CreatedAt, &p.UpdatedAt) // lint-allow: ignore-db-error

	return p
}

func buildTaskStatusSummary(parent types.ParentTask, subtasks []types.Subtask) *types.TaskStatusSummary {
	s := &types.TaskStatusSummary{
		ParentTask: parent,
		Subtasks:   subtasks,
	}
	for _, sub := range subtasks {
		incrementStatusCount(s, sub.Status)
	}
	isAllDone := len(subtasks) > 0 && s.CompletedCount == len(subtasks)
	s.IsCompleted = isAllDone

	return s
}

func incrementStatusCount(s *types.TaskStatusSummary, status string) {
	switch status {
	case "PENDING":
		s.PendingCount++
	case "IN_PROGRESS", "WORKING":
		s.InProgressCount++
	case "DONE", "COMPLETED":
		s.CompletedCount++
	case "FAILED":
		s.FailedCount++
	}
}

// InitAgentSplitDB initializes the Tier 3 agent split database schema.
func InitAgentSplitDB(dbPath string) (*sql.DB, *appfault.AppError) {
	conn, openErr := openAgentDB(dbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return nil, openErr
	}
	ddl := []string{sqlCreateAgentActionLog}
	execErr := executeAgentDDL(conn, ddl, "InitAgentSplitDB")
	hasExecErr := execErr != nil
	if hasExecErr {
		_ = conn.Close()
		return nil, execErr
	}

	return conn, nil
}

const (
	sqlCreateAgentActionLog = `CREATE TABLE IF NOT EXISTS AgentActionLog (
    ActionLogId INTEGER PRIMARY KEY AUTOINCREMENT,
    SubtaskId TEXT NOT NULL,
    AgentRole TEXT NOT NULL,
    ActionType TEXT NOT NULL,
    TargetFile TEXT NOT NULL DEFAULT '',
    StartLine INTEGER NOT NULL DEFAULT 0,
    EndLine INTEGER NOT NULL DEFAULT 0,
    QueryOrCommand TEXT NOT NULL DEFAULT '',
    ActionDetails TEXT NOT NULL DEFAULT '',
    DurationMs INTEGER NOT NULL DEFAULT 0,
    Status TEXT NOT NULL DEFAULT 'SUCCESS',
    ErrorMessage TEXT NOT NULL DEFAULT '',
    CreatedAt TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS IdxActionLog_Subtask ON AgentActionLog(SubtaskId);
CREATE INDEX IF NOT EXISTS IdxActionLog_Type ON AgentActionLog(ActionType);
CREATE INDEX IF NOT EXISTS IdxActionLog_Agent ON AgentActionLog(AgentRole);`

	sqlInsertAgentActionLog = `INSERT INTO AgentActionLog (
SubtaskId, AgentRole, ActionType, TargetFile, StartLine, EndLine, QueryOrCommand, ActionDetails, DurationMs, Status, ErrorMessage, CreatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
)

// LogAgentAction records high-frequency telemetry in Tier 3 agent split DB.
func LogAgentAction(agentDbPath string, log types.AgentActionLog) *appfault.AppError {
	conn, openErr := InitAgentSplitDB(agentDbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return openErr
	}
	defer conn.Close()

	now := resolveActionCreatedAt(log.CreatedAt)
	res := ExecWrapper(conn, sqlInsertAgentActionLog, log.SubtaskId, log.AgentRole, string(log.ActionType), log.TargetFile, log.StartLine, log.EndLine, log.QueryOrCommand, log.ActionDetails, log.DurationMs, log.Status, log.ErrorMessage, now)
	if res.IsFailure {
		return appfault.WrapSimple(res.Error, "LogAgentAction")
	}

	return nil
}

func resolveActionCreatedAt(raw string) string {
	hasRaw := len(raw) > 0
	if hasRaw {
		return raw
	}

	return time.Now().UTC().Format(time.RFC3339)
}

// ListAgentActions retrieves telemetry records from the Tier 3 agent split DB.
func ListAgentActions(
	agentDbPath string,
	subtaskId string,
	limit int,
) ([]types.AgentActionLog, *appfault.AppError) {
	conn, openErr := InitAgentSplitDB(agentDbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return nil, openErr
	}
	defer conn.Close()

	rows, queryErr := queryAgentActionRows(conn, subtaskId, limit)
	hasQueryErr := queryErr != nil
	if hasQueryErr {
		return nil, queryErr
	}
	defer rows.Close()

	return scanAgentActions(rows)
}

func queryAgentActionRows(
	conn *sql.DB,
	subtaskId string,
	limit int,
) (*sql.Rows, *appfault.AppError) {
	boundedLimit := sanitizeTaskLimit(limit)
	query, args := buildListActionsQuery(subtaskId, boundedLimit)
	res := QueryWrapper(conn, query, args...)
	if res.IsFailure {
		return nil, appfault.WrapSimple(res.Error, "queryAgentActionRows")
	}

	return res.Data, nil
}

func buildListActionsQuery(subtaskId string, limit int) (string, []any) {
	hasSubtask := len(subtaskId) > 0
	if hasSubtask {
		q := `SELECT ActionLogId, SubtaskId, AgentRole, ActionType, TargetFile, StartLine, EndLine, QueryOrCommand, ActionDetails, DurationMs, Status, ErrorMessage, CreatedAt FROM AgentActionLog WHERE SubtaskId = ? ORDER BY ActionLogId DESC LIMIT ?`
		return q, []any{subtaskId, limit}
	}
	q := `SELECT ActionLogId, SubtaskId, AgentRole, ActionType, TargetFile, StartLine, EndLine, QueryOrCommand, ActionDetails, DurationMs, Status, ErrorMessage, CreatedAt FROM AgentActionLog ORDER BY ActionLogId DESC LIMIT ?`

	return q, []any{limit}
}

func scanAgentActions(rows *sql.Rows) ([]types.AgentActionLog, *appfault.AppError) {
	var list []types.AgentActionLog
	for rows.Next() {
		var log types.AgentActionLog
		var actionType string
		scanErr := rows.Scan(
			&log.ActionLogId, &log.SubtaskId, &log.AgentRole, &actionType,
			&log.TargetFile, &log.StartLine, &log.EndLine, &log.QueryOrCommand,
			&log.ActionDetails, &log.DurationMs, &log.Status,
			&log.ErrorMessage, &log.CreatedAt,
		)
		if scanErr == nil {
			log.ActionType = types.AgentActionType(actionType)
			list = append(list, log)
		}
	}

	return list, nil
}

// FindTaskDbPath locates the SQLite database path for a given task ID or slug.
func FindTaskDbPath(repoRoot, taskIdOrSlug string) (string, *appfault.AppError) {
	hasFile := strings.HasSuffix(taskIdOrSlug, ".db") && isFileExisting(taskIdOrSlug)
	if hasFile {
		return filepath.ToSlash(taskIdOrSlug), nil
	}
	masterDb := ResolveMasterAgentDbPath(repoRoot)
	resolved, findErr := lookupRootDbInMaster(masterDb, taskIdOrSlug)
	hasResolved := findErr == nil && len(resolved) > 0
	if hasResolved {
		return resolved, nil
	}

	return ResolveTaskDbPath(repoRoot, taskIdOrSlug), nil
}

func lookupRootDbInMaster(masterDbPath, taskIdOrSlug string) (string, *appfault.AppError) {
	conn, openErr := InitMasterAgentDB(masterDbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return "", openErr
	}
	defer conn.Close()

	hasTask := len(taskIdOrSlug) > 0
	if hasTask {
		return scanRootDbByTask(conn, taskIdOrSlug)
	}

	return scanLatestRootDb(conn)
}

func scanRootDbByTask(conn *sql.DB, taskIdOrSlug string) (string, *appfault.AppError) {
	var rootDbPath string
	query := "SELECT RootDbPath FROM ParentTaskRegistry WHERE ParentTaskId = ? OR TaskSlug = ? LIMIT 1"
	row := QueryRowWrapper(conn, query, taskIdOrSlug, taskIdOrSlug)
	if err := row.Scan(&rootDbPath); err == nil {
		return rootDbPath, nil
	}

	return scanLatestRootDb(conn)
}

func scanLatestRootDb(conn *sql.DB) (string, *appfault.AppError) {
	var rootDbPath string
	query := "SELECT RootDbPath FROM ParentTaskRegistry ORDER BY CreatedAt DESC LIMIT 1"
	row := QueryRowWrapper(conn, query)
	_ = row.Scan(&rootDbPath) // lint-allow: ignore-db-error

	return rootDbPath, nil
}
