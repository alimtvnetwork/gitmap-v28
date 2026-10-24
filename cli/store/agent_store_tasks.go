package store

import (
	"database/sql"
	"time"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/types"
	_ "modernc.org/sqlite"
)

// InitTaskDB initializes the Tier 2 parent task database schema.
func InitTaskDB(dbPath string) (*sql.DB, *appfault.AppError) {
	conn, openErr := openAgentDB(dbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return nil, openErr
	}

	return prepareTaskDBConn(conn)
}

// prepareTaskDBConn migrates legacy columns before applying the task
// DDL, closing the connection whenever a step fails.
func prepareTaskDBConn(conn *sql.DB) (*sql.DB, *appfault.AppError) {
	// Migrate legacy DBs before DDL: IdxSubtask_Slug in taskDDL requires the
	// TaskSlug column, so the column must exist before the DDL executes.
	migErr := ensureSubtaskSlugColumnConn(conn)
	hasMigErr := migErr != nil
	if hasMigErr {
		_ = conn.Close()

		return nil, migErr
	}

	return applyTaskDDL(conn)
}

// applyTaskDDL executes the Tier 2 task schema DDL against conn,
// closing the connection on failure.
func applyTaskDDL(conn *sql.DB) (*sql.DB, *appfault.AppError) {
	execErr := executeAgentDDL(conn, taskDDL, "InitTaskDB")
	hasExecErr := execErr != nil
	if hasExecErr {
		_ = conn.Close()

		return nil, execErr
	}

	return conn, nil
}

const (
	sqlCreateParentTask = `CREATE TABLE IF NOT EXISTS ParentTask (
    ParentTaskId TEXT PRIMARY KEY,
    TaskSlug TEXT NOT NULL,
    TaskName TEXT NOT NULL,
    RunDirectory TEXT NOT NULL,
    RootDbPath TEXT NOT NULL,
    Status TEXT NOT NULL DEFAULT 'ACTIVE',
    TotalStepsBudget INTEGER NOT NULL DEFAULT 0,
    CompletedSteps INTEGER NOT NULL DEFAULT 0,
    SpawnedAgentCount INTEGER NOT NULL DEFAULT 0,
    CreatedAt TEXT NOT NULL,
    UpdatedAt TEXT NOT NULL
);`

	sqlCreateSubtask = `CREATE TABLE IF NOT EXISTS Subtask (
    SubtaskId TEXT PRIMARY KEY,
    ParentTaskId TEXT NOT NULL,
    TaskCode TEXT NOT NULL,
    TaskSlug TEXT NOT NULL DEFAULT '',
    Title TEXT NOT NULL,
    AssignedAgentRole TEXT NOT NULL DEFAULT '',
    OwnedFilesJson TEXT NOT NULL DEFAULT '[]',
    Status TEXT NOT NULL DEFAULT 'PENDING',
    Evidence TEXT NOT NULL DEFAULT '',
    CreatedAt TEXT NOT NULL,
    UpdatedAt TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS IdxSubtask_Parent ON Subtask(ParentTaskId);
CREATE INDEX IF NOT EXISTS IdxSubtask_Status ON Subtask(Status);
CREATE INDEX IF NOT EXISTS IdxSubtask_Slug ON Subtask(TaskSlug);`

	sqlCreateSubtaskAuditRollup = `CREATE TABLE IF NOT EXISTS SubtaskAuditRollup (
    RollupId INTEGER PRIMARY KEY AUTOINCREMENT,
    SubtaskId TEXT NOT NULL,
    TotalActions INTEGER NOT NULL DEFAULT 0,
    TotalDurationMs INTEGER NOT NULL DEFAULT 0,
    LastAction TEXT NOT NULL DEFAULT '',
    UpdatedAt TEXT NOT NULL
);`

	sqlUpsertParentTask = `INSERT INTO ParentTask (
ParentTaskId, TaskSlug, TaskName, RunDirectory, RootDbPath, Status, TotalStepsBudget, CompletedSteps, SpawnedAgentCount, CreatedAt, UpdatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(ParentTaskId) DO UPDATE SET
TaskSlug=excluded.TaskSlug, TaskName=excluded.TaskName, RunDirectory=excluded.RunDirectory,
RootDbPath=excluded.RootDbPath, Status=excluded.Status, TotalStepsBudget=excluded.TotalStepsBudget,
CompletedSteps=excluded.CompletedSteps, SpawnedAgentCount=excluded.SpawnedAgentCount, UpdatedAt=excluded.UpdatedAt`

	sqlUpsertSubtask = `INSERT INTO Subtask (
SubtaskId, ParentTaskId, TaskCode, TaskSlug, Title, AssignedAgentRole, OwnedFilesJson, Status, Evidence, CreatedAt, UpdatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(SubtaskId) DO UPDATE SET
TaskSlug=excluded.TaskSlug, Title=excluded.Title, AssignedAgentRole=excluded.AssignedAgentRole, OwnedFilesJson=excluded.OwnedFilesJson,
Status=excluded.Status, Evidence=excluded.Evidence, UpdatedAt=excluded.UpdatedAt`

	sqlCompleteSubtask = `UPDATE Subtask SET Status = 'DONE', AssignedAgentRole = ?, Evidence = ?, UpdatedAt = ? WHERE SubtaskId = ?`

	sqlFailSubtask = `UPDATE Subtask SET Status = 'FAILED', AssignedAgentRole = ?, Evidence = ?, UpdatedAt = ? WHERE SubtaskId = ?`
)

// UpsertParentTask inserts or updates the parent task record in Tier 2 DB.
func UpsertParentTask(taskDbPath string, task types.ParentTask) *appfault.AppError {
	conn, openErr := InitTaskDB(taskDbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return openErr
	}
	defer conn.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	res := ExecWrapper(conn, sqlUpsertParentTask, task.ParentTaskId, task.TaskSlug, task.TaskName, task.RunDirectory, task.RootDbPath, task.Status, task.TotalStepsBudget, task.CompletedSteps, task.SpawnedAgentCount, task.CreatedAt, now)
	if res.IsFailure {
		return appfault.WrapSimple(res.Error, "UpsertParentTask")
	}

	return nil
}

// AddSubtasks enqueues multiple subtasks into the Tier 2 parent task database.
func AddSubtasks(taskDbPath string, parentId string, subtasks []types.Subtask) *appfault.AppError {
	conn, openErr := InitTaskDB(taskDbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return openErr
	}
	defer conn.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	for _, sub := range subtasks {
		insErr := insertSubtaskRow(conn, sub, parentId, now)
		hasInsErr := insErr != nil
		if hasInsErr {
			return insErr
		}
	}

	return nil
}

func insertSubtaskRow(conn *sql.DB, sub types.Subtask, parentId, now string) *appfault.AppError {
	pId := resolveSubtaskParentId(sub.ParentTaskId, parentId)
	cAt := resolveSubtaskCreatedAt(sub.CreatedAt, now)
	res := ExecWrapper(conn, sqlUpsertSubtask, sub.SubtaskId, pId, sub.TaskCode, sub.TaskSlug, sub.Title, sub.AssignedAgentRole, sub.OwnedFilesJson, sub.Status, sub.Evidence, cAt, now)
	if res.IsFailure {
		return appfault.WrapSimple(res.Error, "insertSubtaskRow")
	}

	return nil
}

func resolveSubtaskParentId(current, fallback string) string {
	hasCurrent := len(current) > 0
	if hasCurrent {
		return current
	}

	return fallback
}

func resolveSubtaskCreatedAt(current, fallback string) string {
	hasCurrent := len(current) > 0
	if hasCurrent {
		return current
	}

	return fallback
}

// ClaimSubtask claims the next available pending subtask atomically.
func ClaimSubtask(
	taskDbPath string,
	parentId string,
	agentRole string,
) (*types.Subtask, *appfault.AppError) {
	conn, openErr := InitTaskDB(taskDbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return nil, openErr
	}
	defer conn.Close()

	return executeTxClaim(conn, parentId, agentRole)
}

func executeTxClaim(conn *sql.DB, parentId, agentRole string) (*types.Subtask, *appfault.AppError) {
	tx, txErr := conn.Begin()
	hasTxErr := txErr != nil
	if hasTxErr {
		return nil, appfault.WrapSimple(txErr, "executeTxClaim.begin")
	}
	sub, claimErr := executeClaimInTx(tx, parentId, agentRole)
	hasClaimErr := claimErr != nil
	if hasClaimErr {
		_ = tx.Rollback()
		return nil, claimErr
	}
	_ = tx.Commit()

	return sub, nil
}

func executeClaimInTx(tx *sql.Tx, parentId, agentRole string) (*types.Subtask, *appfault.AppError) {
	sub, findErr := selectPendingSubtask(tx, parentId)
	hasFindErr := findErr != nil
	if hasFindErr {
		return nil, findErr
	}
	hasSub := sub != nil
	if !hasSub {
		return nil, nil
	}
	now := time.Now().UTC().Format(time.RFC3339)
	res := ExecWrapper(tx, "UPDATE Subtask SET Status = 'IN_PROGRESS', AssignedAgentRole = ?, UpdatedAt = ? WHERE SubtaskId = ?", agentRole, now, sub.SubtaskId)
	if res.IsFailure {
		return nil, appfault.WrapSimple(res.Error, "executeClaimInTx.update")
	}
	sub.Status = "IN_PROGRESS"
	sub.AssignedAgentRole = agentRole
	sub.UpdatedAt = now

	return sub, nil
}

func selectPendingSubtask(tx *sql.Tx, parentId string) (*types.Subtask, *appfault.AppError) {
	query, args := buildPendingSubtaskQuery(parentId)
	res := QueryWrapper(tx, query, args...)
	if res.IsFailure {
		return nil, appfault.WrapSimple(res.Error, "selectPendingSubtask")
	}
	defer res.Data.Close()

	hasRow := res.Data.Next()
	if !hasRow {
		return nil, nil
	}
	var sub types.Subtask
	scanErr := res.Data.Scan(
		&sub.SubtaskId, &sub.ParentTaskId, &sub.TaskCode, &sub.Title,
		&sub.AssignedAgentRole, &sub.OwnedFilesJson, &sub.Status,
		&sub.Evidence, &sub.CreatedAt, &sub.UpdatedAt,
	)
	if scanErr != nil {
		return nil, appfault.WrapSimple(scanErr, "selectPendingSubtask.scan")
	}

	return &sub, nil
}

func buildPendingSubtaskQuery(parentId string) (string, []any) {
	hasParent := len(parentId) > 0
	if hasParent {
		q := `SELECT SubtaskId, ParentTaskId, TaskCode, Title, AssignedAgentRole, OwnedFilesJson, Status, Evidence, CreatedAt, UpdatedAt FROM Subtask WHERE Status = 'PENDING' AND ParentTaskId = ? ORDER BY rowid ASC LIMIT 1`
		return q, []any{parentId}
	}
	q := `SELECT SubtaskId, ParentTaskId, TaskCode, Title, AssignedAgentRole, OwnedFilesJson, Status, Evidence, CreatedAt, UpdatedAt FROM Subtask WHERE Status = 'PENDING' ORDER BY rowid ASC LIMIT 1`

	return q, []any{}
}
