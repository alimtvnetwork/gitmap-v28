// Package store — pendingtask.go manages PendingTask and CompletedTask CRUD.
package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// InsertPendingTask creates a new pending task and returns its ID.
func (db *DB) InsertPendingTask(
	taskTypeID int64,
	targetPath,
	workDir,
	sourceCmd,
	cmdArgs string,
) (int64, error) {
	result, err := ExecWrapper(db.conn, constants.SQLInsertPendingTask,
		taskTypeID, targetPath, workDir, sourceCmd, cmdArgs).Destruct()
	if err != nil {
		return 0, apperror.WrapSimple(err, "failed to insert pending task (operation: insert)")
	}

	return result.LastInsertId()
}

// FindPendingTaskDuplicate checks if a pending task already exists for the given type and path.
// Returns the existing task ID or 0 if none found.
func (db *DB) FindPendingTaskDuplicate(taskTypeID int64, targetPath string) int64 {
	row := QueryRowWrapper(db.conn, constants.SQLSelectPendingTaskByTypePath,
		taskTypeID, targetPath)

	var id int64

	err := row.Scan(&id)
	if err != nil {
		return 0
	}

	return id
}

// FindPendingTaskDuplicateWithCmd checks if a pending task exists for the given type, path, and command args.
// Returns the existing task ID or 0 if none found.
func (db *DB) FindPendingTaskDuplicateWithCmd(taskTypeID int64, targetPath, cmdArgs string) int64 {
	row := QueryRowWrapper(db.conn, constants.SQLSelectPendingTaskByTypePathCmd,
		taskTypeID, targetPath, cmdArgs)

	var id int64

	err := row.Scan(&id)
	if err != nil {
		return 0
	}

	return id
}

// ListPendingTasks returns all pending tasks ordered by ID.
func (db *DB) ListPendingTasks() ([]model.PendingTaskRecord, error) {
	rows, err := QueryWrapper(db.conn, constants.SQLSelectAllPendingTasks).Destruct()
	if err != nil {
		return nil, apperror.WrapSimple(err, "failed to query pending tasks (operation: query)")
	}

	defer rows.Close()

	return scanPendingTaskRows(rows)
}

// FindPendingTaskByID returns a single pending task by ID.
func (db *DB) FindPendingTaskByID(id int64) (model.PendingTaskRecord, error) {
	row := QueryRowWrapper(db.conn, constants.SQLSelectPendingTaskByID, id)

	var r model.PendingTaskRecord

	err := row.Scan(&r.ID, &r.TaskTypeId, &r.TaskTypeName, &r.TargetPath,
		&r.WorkingDirectory, &r.SourceCommand, &r.CommandArgs,
		&r.FailureReason, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return r, apperror.WrapSimple(err, "failed to query pending tasks (operation: query)")
	}

	return r, nil
}

// CompleteTask moves a pending task to the completed table in a transaction.
func (db *DB) CompleteTask(taskID int64) error {
	wrap, appErr := dbengine.WrapDB(db.conn, dbengine.DBSQLite)
	if appErr != nil {
		return apperror.WrapSimple(appErr, "failed to complete task (operation: complete)")
	}

	ctx := context.Background()
	appErr = wrap.WithTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError {
		return executeCompleteTaskTx(ctx, tx, taskID)
	})
	if appErr != nil {
		return apperror.WrapSimple(appErr, "failed to complete task (operation: complete)")
	}

	return nil
}

func executeCompleteTaskTx(ctx context.Context, tx *dbengine.TxWrapper, taskID int64) *apperror.AppError {
	taskRes := findPendingTaskInTx(ctx, tx, taskID)
	appErr := taskRes.AppError()
	if isPendingTaskAlreadyCompleted(ctx, tx, taskID, appErr) {
		return nil
	}

	if appErr != nil {
		return appErr
	}

	if insertErr := insertCompletedTaskInTx(ctx, tx, taskRes.Value); insertErr != nil {
		return insertErr
	}

	return deletePendingTaskInTx(ctx, tx, taskID)
}

func isPendingTaskAlreadyCompleted(ctx context.Context, tx *dbengine.TxWrapper, taskID int64, appErr *apperror.AppError) bool {
	if appErr == nil {
		return false
	}
	isNoRows := errors.Is(appErr.Unwrap(), sql.ErrNoRows) || strings.Contains(appErr.Error(), "no rows in result set")
	if !isNoRows {
		return false
	}
	row, qErr := tx.QueryRow(ctx, "SELECT 1 FROM CompletedTask WHERE OriginalTaskId = ? LIMIT 1", taskID)
	if qErr != nil {
		return false
	}
	var dummy int
	return row.Scan(&dummy) == nil
}

func insertCompletedTaskInTx(ctx context.Context, tx *dbengine.TxWrapper, task model.PendingTaskRecord) *apperror.AppError {
	_, appErr := tx.Exec(ctx, constants.SQLInsertCompletedTask,
		task.ID, task.TaskTypeId, task.TargetPath, task.WorkingDirectory,
		task.SourceCommand, task.CommandArgs, task.CreatedAt)
	if appErr != nil {
		return appErr
	}

	return nil
}

func deletePendingTaskInTx(ctx context.Context, tx *dbengine.TxWrapper, taskID int64) *apperror.AppError {
	_, appErr := tx.Exec(ctx, constants.SQLDeletePendingTask, taskID)
	if appErr != nil {
		return appErr
	}

	return nil
}

// FailTask updates the failure reason for a pending task.
func (db *DB) FailTask(taskID int64, reason string) error {
	result, err := ExecWrapper(db.conn, constants.SQLUpdatePendingTaskFailure, reason, taskID).Destruct()
	if err != nil {
		return apperror.WrapSimple(err, "failed to update task failure (operation: update)")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apperror.NewNotFoundError(fmt.Sprintf("pending task not found: %d", taskID))
	}

	return nil
}

// ListCompletedTasks returns all completed tasks ordered by completion time.
func (db *DB) ListCompletedTasks() ([]model.CompletedTaskRecord, error) {
	rows, err := QueryWrapper(db.conn, constants.SQLSelectAllCompletedTasks).Destruct()
	if err != nil {
		return nil, apperror.WrapSimple(err, "failed to query pending tasks (operation: query)")
	}

	defer rows.Close()

	return scanCompletedTaskRows(rows)
}

// DeletePendingTask removes a single pending task by ID without
// recording it in CompletedTask. Used by `gitmap pending clear` to
// drop orphaned/illegal entries that should never have been queued.
// Returns ErrPendingTaskNotFound when the ID does not exist so the
// caller can surface a precise message instead of a silent no-op.
func (db *DB) DeletePendingTask(id int64) error {
	result, err := ExecWrapper(db.conn, constants.SQLDeletePendingTask, id).Destruct()
	if err != nil {
		return apperror.WrapSimple(err, "failed to complete task (operation: complete)")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apperror.NewNotFoundError(fmt.Sprintf("pending task not found: %d", id))
	}

	return nil
}

// findPendingTaskInTx reads a pending task within an existing transaction.
func findPendingTaskInTx(ctx context.Context, tx *dbengine.TxWrapper, id int64) PendingTaskRecordWrap {
	row, appErr := tx.QueryRow(ctx, constants.SQLSelectPendingTaskByID, id)
	if appErr != nil {
		return result.Fail[model.PendingTaskRecord](appErr)
	}

	var r model.PendingTaskRecord
	err := row.Scan(&r.ID, &r.TaskTypeId, &r.TaskTypeName, &r.TargetPath,
		&r.WorkingDirectory, &r.SourceCommand, &r.CommandArgs,
		&r.FailureReason, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return result.Fail[model.PendingTaskRecord](apperror.WrapSimple(err, "scan pending task in tx"))
	}

	return result.Ok(r)
}
