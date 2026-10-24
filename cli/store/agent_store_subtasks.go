package store

import (
	"database/sql"
	"strings"
	"time"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/types"
	_ "modernc.org/sqlite"
)

// StartSubtask marks a specific subtask as IN_PROGRESS by an agent.
func StartSubtask(taskDbPath string, subtaskId string, agentRole string) *appfault.AppError {
	conn, openErr := InitTaskDB(taskDbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return openErr
	}
	defer conn.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	res := ExecWrapper(conn, "UPDATE Subtask SET Status = 'IN_PROGRESS', AssignedAgentRole = ?, UpdatedAt = ? WHERE SubtaskId = ?", agentRole, now, subtaskId)
	if res.IsFailure {
		return appfault.WrapSimple(res.Error, "StartSubtask")
	}

	return nil
}

// CompleteSubtask marks a subtask as DONE with verification evidence and increments completed steps.
func CompleteSubtask(
	taskDbPath string,
	subtaskId string,
	agentRole string,
	evidence string,
) *appfault.AppError {
	conn, openErr := InitTaskDB(taskDbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return openErr
	}
	defer conn.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	res := ExecWrapper(conn, sqlCompleteSubtask, agentRole, evidence, now, subtaskId)
	if res.IsFailure {
		return appfault.WrapSimple(res.Error, "CompleteSubtask")
	}
	_ = ExecWrapper(conn, "UPDATE ParentTask SET CompletedSteps = CompletedSteps + 1, UpdatedAt = ?", now)

	return nil
}

// FailSubtask marks a subtask as FAILED with failure reasoning.
func FailSubtask(
	taskDbPath string,
	subtaskId string,
	agentRole string,
	reason string,
) *appfault.AppError {
	conn, openErr := InitTaskDB(taskDbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return openErr
	}
	defer conn.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	res := ExecWrapper(conn, sqlFailSubtask, agentRole, reason, now, subtaskId)
	if res.IsFailure {
		return appfault.WrapSimple(res.Error, "FailSubtask")
	}

	return nil
}

// ListSubtasks lists all subtasks belonging to a parent task.
func ListSubtasks(taskDbPath string, parentId string) ([]types.Subtask, *appfault.AppError) {
	conn, openErr := InitTaskDB(taskDbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return nil, openErr
	}
	defer conn.Close()

	rows, queryErr := querySubtaskRows(conn, parentId)
	hasQueryErr := queryErr != nil
	if hasQueryErr {
		return nil, queryErr
	}
	defer rows.Close()

	return scanSubtasks(rows)
}

func querySubtaskRows(conn *sql.DB, parentId string) (*sql.Rows, *appfault.AppError) {
	query, args := buildListSubtasksQuery(parentId)
	res := QueryWrapper(conn, query, args...)
	if res.IsFailure {
		return nil, appfault.WrapSimple(res.Error, "querySubtaskRows")
	}

	return res.Data, nil
}

func buildListSubtasksQuery(parentId string) (string, []any) {
	hasParent := len(parentId) > 0
	if hasParent {
		q := `SELECT SubtaskId, ParentTaskId, TaskCode, TaskSlug, Title, AssignedAgentRole, OwnedFilesJson, Status, Evidence, CreatedAt, UpdatedAt FROM Subtask WHERE ParentTaskId = ? ORDER BY rowid ASC`
		return q, []any{parentId}
	}
	q := `SELECT SubtaskId, ParentTaskId, TaskCode, TaskSlug, Title, AssignedAgentRole, OwnedFilesJson, Status, Evidence, CreatedAt, UpdatedAt FROM Subtask ORDER BY rowid ASC`

	return q, []any{}
}

func scanSubtasks(rows *sql.Rows) ([]types.Subtask, *appfault.AppError) {
	var list []types.Subtask
	for rows.Next() {
		var s types.Subtask
		scanErr := rows.Scan(
			&s.SubtaskId, &s.ParentTaskId, &s.TaskCode, &s.TaskSlug, &s.Title,
			&s.AssignedAgentRole, &s.OwnedFilesJson, &s.Status,
			&s.Evidence, &s.CreatedAt, &s.UpdatedAt,
		)
		if scanErr == nil {
			list = append(list, s)
		}
	}

	return list, nil
}

// FindParentTaskBySlug retrieves a parent task from the Tier 1 master DB by
// sanitized slug. Returns (nil, nil) when no task matches. The input is
// sanitized before comparison so Title Case slugs from prompts match storage.
func FindParentTaskBySlug(
	masterDbPath string,
	slug string,
) (*types.ParentTask, *appfault.AppError) {
	conn, openErr := InitMasterAgentDB(masterDbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return nil, openErr
	}
	defer conn.Close()

	clean := SanitizeSlug(slug)
	query := `SELECT ParentTaskId, TaskSlug, TaskName, RunDirectory, RootDbPath, Status, TotalStepsBudget, CompletedSteps, SpawnedAgentCount, CreatedAt, UpdatedAt FROM ParentTaskRegistry WHERE TaskSlug = ? LIMIT 1`
	res := QueryWrapper(conn, query, clean)
	if res.IsFailure {
		return nil, appfault.WrapSimple(res.Error, "FindParentTaskBySlug")
	}
	defer res.Data.Close()

	tasks, scanErr := scanParentTasks(res.Data)
	hasScanErr := scanErr != nil
	if hasScanErr {
		return nil, scanErr
	}
	hasTask := len(tasks) > 0
	if !hasTask {
		return nil, nil
	}

	return &tasks[0], nil
}

// ensureSubtaskSlugColumnConn adds TaskSlug to Subtask on databases created
// before the column existed. Idempotent: no-op when the table or the column
// is already present.
func ensureSubtaskSlugColumnConn(conn *sql.DB) *appfault.AppError {
	hasTable, tblErr := agentTableExists(conn, "Subtask")
	hasTblErr := tblErr != nil
	if hasTblErr {
		return tblErr
	}
	if !hasTable {
		return nil
	}
	res := QueryWrapper(conn, "PRAGMA table_info(Subtask)")
	if res.IsFailure {
		return appfault.WrapSimple(res.Error, "ensureSubtaskSlugColumn")
	}
	defer res.Data.Close()

	hasColumn := false
	for res.Data.Next() {
		var cid int
		var name, colType string
		var notNull int
		var dflt sql.NullString
		var pk int
		if scanErr := res.Data.Scan(&cid, &name, &colType, &notNull, &dflt, &pk); scanErr != nil {
			continue
		}
		if name == "TaskSlug" {
			hasColumn = true
			break
		}
	}
	if hasColumn {
		return nil
	}
	alter := ExecWrapper(conn, "ALTER TABLE Subtask ADD COLUMN TaskSlug TEXT NOT NULL DEFAULT ''")
	if alter.IsFailure {
		return appfault.WrapSimple(alter.Error, "ensureSubtaskSlugColumn")
	}

	return nil
}

func agentTableExists(conn *sql.DB, table string) (bool, *appfault.AppError) {
	res := QueryWrapper(conn, "SELECT 1 FROM sqlite_master WHERE type='table' AND name=? LIMIT 1", table)
	if res.IsFailure {
		return false, appfault.WrapSimple(res.Error, "agentTableExists")
	}
	defer res.Data.Close()

	return res.Data.Next(), nil
}

// ListPendingSubtasks retrieves PENDING subtasks from a Tier 2 task DB,
// optionally scoped to parent IDs (covers both stored ID and slug forms).
func ListPendingSubtasks(
	taskDbPath string,
	parentIds ...string,
) ([]types.Subtask, *appfault.AppError) {
	conn, openErr := InitTaskDB(taskDbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return nil, openErr
	}
	defer conn.Close()

	query, args := buildPendingSubtasksQuery(parentIds)
	res := QueryWrapper(conn, query, args...)
	if res.IsFailure {
		return nil, appfault.WrapSimple(res.Error, "ListPendingSubtasks")
	}
	defer res.Data.Close()

	return scanSubtasks(res.Data)
}

func buildPendingSubtasksQuery(parentIds []string) (string, []any) {
	base := `SELECT SubtaskId, ParentTaskId, TaskCode, TaskSlug, Title, AssignedAgentRole, OwnedFilesJson, Status, Evidence, CreatedAt, UpdatedAt FROM Subtask WHERE Status = 'PENDING'`
	hasParents := len(parentIds) > 0
	if !hasParents {
		return base + ` ORDER BY rowid ASC`, []any{}
	}
	placeholders := make([]string, len(parentIds))
	args := make([]any, len(parentIds))
	for i, id := range parentIds {
		placeholders[i] = "?"
		args[i] = id
	}

	return base + ` AND ParentTaskId IN (` + strings.Join(placeholders, ",") + `) ORDER BY rowid ASC`, args
}

// GetTaskStatus retrieves the overall status summary for a parent task.
func GetTaskStatus(
	taskDbPath string,
	parentId string,
) (*types.TaskStatusSummary, *appfault.AppError) {
	subtasks, listErr := ListSubtasks(taskDbPath, parentId)
	hasListErr := listErr != nil
	if hasListErr {
		return nil, listErr
	}
	conn, openErr := InitTaskDB(taskDbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return nil, openErr
	}
	defer conn.Close()

	parent := queryParentTask(conn, parentId)

	return buildTaskStatusSummary(parent, subtasks), nil
}
