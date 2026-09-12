package constants

// SQL: pending task operations (v15: PendingTaskId / CompletedTaskId / TaskTypeId PKs).
const (
	SQLInsertPendingTask = `INSERT INTO PendingTask
		(TaskTypeId, TargetPath, WorkingDirectory, SourceCommand, CommandArgs)
		VALUES (?, ?, ?, ?, ?)`

	SQLSelectAllPendingTasks = `SELECT p.PendingTaskId, p.TaskTypeId, t.Name, p.TargetPath,
		p.WorkingDirectory, p.SourceCommand, p.CommandArgs,
		p.FailureReason, p.CreatedAt, p.UpdatedAt
		FROM PendingTask p JOIN TaskType t ON p.TaskTypeId = t.TaskTypeId
		ORDER BY p.PendingTaskId`

	SQLSelectPendingTaskByID = `SELECT p.PendingTaskId, p.TaskTypeId, t.Name, p.TargetPath,
		p.WorkingDirectory, p.SourceCommand, p.CommandArgs,
		p.FailureReason, p.CreatedAt, p.UpdatedAt
		FROM PendingTask p JOIN TaskType t ON p.TaskTypeId = t.TaskTypeId
		WHERE p.PendingTaskId = ?`

	SQLSelectPendingTaskByTypePath = `SELECT p.PendingTaskId FROM PendingTask p
		WHERE p.TaskTypeId = ? AND p.TargetPath = ?`

	SQLSelectPendingTaskByTypePathCmd = `SELECT p.PendingTaskId FROM PendingTask p
		WHERE p.TaskTypeId = ? AND p.TargetPath = ? AND p.CommandArgs = ?`

	SQLUpdatePendingTaskFailure = `UPDATE PendingTask
		SET FailureReason = ?, UpdatedAt = CURRENT_TIMESTAMP
		WHERE PendingTaskId = ?`

	SQLDeletePendingTask = `DELETE FROM PendingTask WHERE PendingTaskId = ?`
)

// SQL: completed task operations (v15: CompletedTaskId PK).
const (
	SQLInsertCompletedTask = `INSERT INTO CompletedTask
		(OriginalTaskId, TaskTypeId, TargetPath, WorkingDirectory,
		 SourceCommand, CommandArgs, CreatedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	SQLSelectAllCompletedTasks = `SELECT c.CompletedTaskId, c.OriginalTaskId, c.TaskTypeId, t.Name,
		c.TargetPath, c.WorkingDirectory, c.SourceCommand, c.CommandArgs,
		c.CompletedAt, c.CreatedAt
		FROM CompletedTask c JOIN TaskType t ON c.TaskTypeId = t.TaskTypeId
		ORDER BY c.CompletedAt DESC`
)

// SQL: task type lookup (v15: TaskTypeId PK).
const SQLSelectTaskTypeByName = `SELECT TaskTypeId FROM TaskType WHERE Name = ?`
