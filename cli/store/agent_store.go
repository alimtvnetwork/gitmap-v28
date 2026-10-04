package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/types"
	_ "modernc.org/sqlite"
)

// ResolveAgentTempDir resolves the root temporary directory for agent split databases.
func ResolveAgentTempDir(repoRoot string) string {
	envDir := os.Getenv("GITMAP_AGENT_TEMP_DIR")
	hasEnv := len(strings.TrimSpace(envDir)) > 0
	if hasEnv {
		return filepath.ToSlash(envDir)
	}
	root := resolveTargetRoot(repoRoot)
	dir := filepath.Join(root, ".ai-memory", "temp-agents")
	_ = os.MkdirAll(dir, 0755)

	return filepath.ToSlash(dir)
}

// ResolveMasterAgentDbPath resolves the repository root master agent database path (Tier 1).
func ResolveMasterAgentDbPath(repoRoot string) string {
	tempDir := ResolveAgentTempDir(repoRoot)

	return filepath.ToSlash(filepath.Join(tempDir, "ai_agents.db"))
}

// ResolveParentTaskDir resolves the directory for a parent task run (Tier 2).
func ResolveParentTaskDir(repoRoot, taskSlug string) string {
	cleanSlug := SanitizeSlug(taskSlug)
	tempDir := ResolveAgentTempDir(repoRoot)
	dir := filepath.Join(tempDir, cleanSlug)
	_ = os.MkdirAll(dir, 0755)

	return filepath.ToSlash(dir)
}

// ResolveTaskDbPath resolves the parent task database path (Tier 2).
func ResolveTaskDbPath(repoRoot, taskSlug string) string {
	taskDir := ResolveParentTaskDir(repoRoot, taskSlug)

	return filepath.ToSlash(filepath.Join(taskDir, "agent-task.db"))
}

// ResolveAgentDbPath resolves a specific agent split database path (Tier 3).
func ResolveAgentDbPath(repoRoot, taskSlug, agentSlug string) string {
	taskDir := ResolveParentTaskDir(repoRoot, taskSlug)
	agentsDir := filepath.Join(taskDir, "agents")
	_ = os.MkdirAll(agentsDir, 0755)
	cleanAgent := SanitizeSlug(agentSlug)

	return filepath.ToSlash(filepath.Join(agentsDir, cleanAgent+".db"))
}

func ensureDirExists(filePath string) *appfault.AppError {
	dir := filepath.Dir(filePath)
	mkErr := os.MkdirAll(dir, 0755)
	hasErr := mkErr != nil
	if hasErr {
		return appfault.WrapSimple(mkErr, "ensureDirExists")
	}

	return nil
}

func openAgentDB(dbPath string) (*sql.DB, *appfault.AppError) {
	dirErr := ensureDirExists(dbPath)
	hasDirErr := dirErr != nil
	if hasDirErr {
		return nil, dirErr
	}

	return OpenSQLiteDB(dbPath)
}

func executeAgentDDL(conn *sql.DB, queries []string, op string) *appfault.AppError {
	for _, q := range queries {
		res := ExecWrapper(conn, q)
		if res.IsFailure {
			return appfault.WrapSimple(res.Error, op)
		}
	}

	return nil
}

var masterAgentDDL = []string{
	sqlCreateParentTaskRegistry,
	sqlCreateAgentRegistry,
	sqlCreateGlobalLifecycleMetrics,
}

// InitMasterAgentDB initializes the Tier 1 master agent database schema.
func InitMasterAgentDB(dbPath string) (*sql.DB, *appfault.AppError) {
	conn, openErr := openAgentDB(dbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return nil, openErr
	}
	execErr := executeAgentDDL(conn, masterAgentDDL, "InitMasterAgentDB")
	hasExecErr := execErr != nil
	if hasExecErr {
		_ = conn.Close()
		return nil, execErr
	}

	return conn, nil
}

const (
	sqlCreateParentTaskRegistry = `CREATE TABLE IF NOT EXISTS ParentTaskRegistry (
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
);
CREATE INDEX IF NOT EXISTS IdxParentTaskRegistry_Slug ON ParentTaskRegistry(TaskSlug);
CREATE INDEX IF NOT EXISTS IdxParentTaskRegistry_Status ON ParentTaskRegistry(Status);`

	sqlCreateAgentRegistry = `CREATE TABLE IF NOT EXISTS AgentRegistry (
    AgentId TEXT PRIMARY KEY,
    ParentTaskId TEXT NOT NULL,
    AgentRole TEXT NOT NULL,
    AgentSlug TEXT NOT NULL,
    SplitDbPath TEXT NOT NULL,
    Status TEXT NOT NULL DEFAULT 'IDLE',
    LastHeartbeatAt TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS IdxAgentRegistry_Parent ON AgentRegistry(ParentTaskId);`

	sqlCreateGlobalLifecycleMetrics = `CREATE TABLE IF NOT EXISTS GlobalLifecycleMetrics (
    MetricKey TEXT PRIMARY KEY,
    MetricValue TEXT NOT NULL,
    UpdatedAt TEXT NOT NULL
);`

	sqlUpsertParentTaskRegistry = `INSERT INTO ParentTaskRegistry (
ParentTaskId, TaskSlug, TaskName, RunDirectory, RootDbPath, Status, TotalStepsBudget, CompletedSteps, SpawnedAgentCount, CreatedAt, UpdatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(ParentTaskId) DO UPDATE SET
TaskSlug=excluded.TaskSlug, TaskName=excluded.TaskName, RunDirectory=excluded.RunDirectory,
RootDbPath=excluded.RootDbPath, Status=excluded.Status, TotalStepsBudget=excluded.TotalStepsBudget,
CompletedSteps=excluded.CompletedSteps, SpawnedAgentCount=excluded.SpawnedAgentCount, UpdatedAt=excluded.UpdatedAt`
)

// RegisterParentTask registers or updates a parent task record in the Tier 1 master DB.
func RegisterParentTask(masterDbPath string, task types.ParentTask) *appfault.AppError {
	conn, openErr := InitMasterAgentDB(masterDbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return openErr
	}
	defer conn.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	res := ExecWrapper(conn, sqlUpsertParentTaskRegistry, task.ParentTaskId, task.TaskSlug, task.TaskName, task.RunDirectory, task.RootDbPath, task.Status, task.TotalStepsBudget, task.CompletedSteps, task.SpawnedAgentCount, task.CreatedAt, now)
	if res.IsFailure {
		return appfault.WrapSimple(res.Error, "RegisterParentTask")
	}

	return nil
}

// ListParentTasks retrieves parent tasks from Tier 1 master DB matching status.
func ListParentTasks(masterDbPath string, limit int, status string) ([]types.ParentTask, *appfault.AppError) {
	conn, openErr := InitMasterAgentDB(masterDbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return nil, openErr
	}
	defer conn.Close()

	rows, queryErr := queryParentTasksRows(conn, limit, status)
	hasQueryErr := queryErr != nil
	if hasQueryErr {
		return nil, queryErr
	}
	defer rows.Close()

	return scanParentTasks(rows)
}

func queryParentTasksRows(conn *sql.DB, limit int, status string) (*sql.Rows, *appfault.AppError) {
	boundedLimit := sanitizeTaskLimit(limit)
	query, args := buildListParentTasksQuery(boundedLimit, status)
	res := QueryWrapper(conn, query, args...)
	if res.IsFailure {
		return nil, appfault.WrapSimple(res.Error, "queryParentTasksRows")
	}

	return res.Data, nil
}

func sanitizeTaskLimit(limit int) int {
	if limit <= 0 {
		return 50
	}

	return limit
}

func buildListParentTasksQuery(limit int, status string) (string, []any) {
	hasStatus := len(status) > 0 && status != "all"
	if hasStatus {
		q := `SELECT ParentTaskId, TaskSlug, TaskName, RunDirectory, RootDbPath, Status, TotalStepsBudget, CompletedSteps, SpawnedAgentCount, CreatedAt, UpdatedAt FROM ParentTaskRegistry WHERE Status = ? ORDER BY CreatedAt DESC LIMIT ?`
		return q, []any{status, limit}
	}
	q := `SELECT ParentTaskId, TaskSlug, TaskName, RunDirectory, RootDbPath, Status, TotalStepsBudget, CompletedSteps, SpawnedAgentCount, CreatedAt, UpdatedAt FROM ParentTaskRegistry ORDER BY CreatedAt DESC LIMIT ?`

	return q, []any{limit}
}

func scanParentTasks(rows *sql.Rows) ([]types.ParentTask, *appfault.AppError) {
	var list []types.ParentTask
	for rows.Next() {
		var item types.ParentTask
		scanErr := rows.Scan(
			&item.ParentTaskId, &item.TaskSlug, &item.TaskName, &item.RunDirectory,
			&item.RootDbPath, &item.Status, &item.TotalStepsBudget, &item.CompletedSteps,
			&item.SpawnedAgentCount, &item.CreatedAt, &item.UpdatedAt,
		)
		if scanErr == nil {
			list = append(list, item)
		}
	}

	return list, nil
}

// GetGlobalMetrics computes aggregate operational metrics from the Tier 1 master DB.
func GetGlobalMetrics(masterDbPath string) (map[string]string, *appfault.AppError) {
	conn, openErr := InitMasterAgentDB(masterDbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return nil, openErr
	}
	defer conn.Close()

	metrics := make(map[string]string)
	populateGlobalMetrics(conn, metrics)

	return metrics, nil
}

func populateGlobalMetrics(conn *sql.DB, metrics map[string]string) {
	populateCountMetric(conn, metrics, "total_tasks", "SELECT COUNT(*) FROM ParentTaskRegistry")
	populateCountMetric(conn, metrics, "active_tasks", "SELECT COUNT(*) FROM ParentTaskRegistry WHERE Status = 'ACTIVE'")
	populateCountMetric(conn, metrics, "completed_tasks", "SELECT COUNT(*) FROM ParentTaskRegistry WHERE Status IN ('COMPLETED', 'DONE')")
	populateCountMetric(conn, metrics, "failed_tasks", "SELECT COUNT(*) FROM ParentTaskRegistry WHERE Status = 'FAILED'")
	populateCountMetric(conn, metrics, "total_agents", "SELECT COUNT(*) FROM AgentRegistry")
}

func populateCountMetric(conn *sql.DB, target map[string]string, key, query string) {
	row := QueryRowWrapper(conn, query)
	var count int
	if err := row.Scan(&count); err == nil {
		target[key] = fmt.Sprintf("%d", count)
	}
}

var taskDDL = []string{
	sqlCreateParentTask,
	sqlCreateSubtask,
	sqlCreateSubtaskAuditRollup,
}

// InitTaskDB initializes the Tier 2 parent task database schema.
func InitTaskDB(dbPath string) (*sql.DB, *appfault.AppError) {
	conn, openErr := openAgentDB(dbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return nil, openErr
	}
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
    Title TEXT NOT NULL,
    AssignedAgentRole TEXT NOT NULL DEFAULT '',
    OwnedFilesJson TEXT NOT NULL DEFAULT '[]',
    Status TEXT NOT NULL DEFAULT 'PENDING',
    Evidence TEXT NOT NULL DEFAULT '',
    CreatedAt TEXT NOT NULL,
    UpdatedAt TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS IdxSubtask_Parent ON Subtask(ParentTaskId);
CREATE INDEX IF NOT EXISTS IdxSubtask_Status ON Subtask(Status);`

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
SubtaskId, ParentTaskId, TaskCode, Title, AssignedAgentRole, OwnedFilesJson, Status, Evidence, CreatedAt, UpdatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(SubtaskId) DO UPDATE SET
Title=excluded.Title, AssignedAgentRole=excluded.AssignedAgentRole, OwnedFilesJson=excluded.OwnedFilesJson,
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
	res := ExecWrapper(conn, sqlUpsertSubtask, sub.SubtaskId, pId, sub.TaskCode, sub.Title, sub.AssignedAgentRole, sub.OwnedFilesJson, sub.Status, sub.Evidence, cAt, now)
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
func ClaimSubtask(taskDbPath string, parentId string, agentRole string) (*types.Subtask, *appfault.AppError) {
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
func CompleteSubtask(taskDbPath string, subtaskId string, agentRole string, evidence string) *appfault.AppError {
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
func FailSubtask(taskDbPath string, subtaskId string, agentRole string, reason string) *appfault.AppError {
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
		q := `SELECT SubtaskId, ParentTaskId, TaskCode, Title, AssignedAgentRole, OwnedFilesJson, Status, Evidence, CreatedAt, UpdatedAt FROM Subtask WHERE ParentTaskId = ? ORDER BY rowid ASC`
		return q, []any{parentId}
	}
	q := `SELECT SubtaskId, ParentTaskId, TaskCode, Title, AssignedAgentRole, OwnedFilesJson, Status, Evidence, CreatedAt, UpdatedAt FROM Subtask ORDER BY rowid ASC`

	return q, []any{}
}

func scanSubtasks(rows *sql.Rows) ([]types.Subtask, *appfault.AppError) {
	var list []types.Subtask
	for rows.Next() {
		var s types.Subtask
		scanErr := rows.Scan(
			&s.SubtaskId, &s.ParentTaskId, &s.TaskCode, &s.Title,
			&s.AssignedAgentRole, &s.OwnedFilesJson, &s.Status,
			&s.Evidence, &s.CreatedAt, &s.UpdatedAt,
		)
		if scanErr == nil {
			list = append(list, s)
		}
	}

	return list, nil
}

// GetTaskStatus retrieves the overall status summary for a parent task.
func GetTaskStatus(taskDbPath string, parentId string) (*types.TaskStatusSummary, *appfault.AppError) {
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

func queryParentTask(conn *sql.DB, parentId string) types.ParentTask {
	var p types.ParentTask
	query := "SELECT ParentTaskId, TaskSlug, TaskName, RunDirectory, RootDbPath, Status, TotalStepsBudget, CompletedSteps, SpawnedAgentCount, CreatedAt, UpdatedAt FROM ParentTask LIMIT 1"
	hasParent := len(parentId) > 0
	if hasParent {
		query = "SELECT ParentTaskId, TaskSlug, TaskName, RunDirectory, RootDbPath, Status, TotalStepsBudget, CompletedSteps, SpawnedAgentCount, CreatedAt, UpdatedAt FROM ParentTask WHERE ParentTaskId = ? LIMIT 1"
		row := QueryRowWrapper(conn, query, parentId)
		_ = row.Scan(&p.ParentTaskId, &p.TaskSlug, &p.TaskName, &p.RunDirectory, &p.RootDbPath, &p.Status, &p.TotalStepsBudget, &p.CompletedSteps, &p.SpawnedAgentCount, &p.CreatedAt, &p.UpdatedAt)
		return p
	}
	row := QueryRowWrapper(conn, query)
	_ = row.Scan(&p.ParentTaskId, &p.TaskSlug, &p.TaskName, &p.RunDirectory, &p.RootDbPath, &p.Status, &p.TotalStepsBudget, &p.CompletedSteps, &p.SpawnedAgentCount, &p.CreatedAt, &p.UpdatedAt)

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
func ListAgentActions(agentDbPath string, subtaskId string, limit int) ([]types.AgentActionLog, *appfault.AppError) {
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

func queryAgentActionRows(conn *sql.DB, subtaskId string, limit int) (*sql.Rows, *appfault.AppError) {
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
	_ = row.Scan(&rootDbPath)

	return rootDbPath, nil
}
