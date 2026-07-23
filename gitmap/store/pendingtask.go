// Package store — pendingtask.go manages PendingTask and CompletedTask CRUD.
package store

import (
	"database/sql"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// InsertPendingTask creates a new pending task and returns its ID.
func (db *DB) InsertPendingTask(taskTypeID int64, targetPath, workDir, sourceCmd, cmdArgs string) (int64, error) {
	result, err := db.conn.Exec(constants.SQLInsertPendingTask,
		taskTypeID, targetPath, workDir, sourceCmd, cmdArgs)
	if err != nil {
		return 0, fmt.Errorf(constants.ErrPendingTaskInsert, err)
	}

	return result.LastInsertId()
}

// FindPendingTaskDuplicate checks if a pending task already exists for the given type and path.
// Returns the existing task ID or 0 if none found.
func (db *DB) FindPendingTaskDuplicate(taskTypeID int64, targetPath string) int64 {
	row := db.conn.QueryRow(constants.SQLSelectPendingTaskByTypePath,
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
	row := db.conn.QueryRow(constants.SQLSelectPendingTaskByTypePathCmd,
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
	rows, err := db.conn.Query(constants.SQLSelectAllPendingTasks)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrPendingTaskQuery, err)
	}
	defer rows.Close()

	return scanPendingTaskRows(rows)
}

// FindPendingTaskByID returns a single pending task by ID.
func (db *DB) FindPendingTaskByID(id int64) (model.PendingTaskRecord, error) {
	row := db.conn.QueryRow(constants.SQLSelectPendingTaskByID, id)

	var r model.PendingTaskRecord

	err := row.Scan(&r.ID, &r.TaskTypeId, &r.TaskTypeName, &r.TargetPath,
		&r.WorkingDirectory, &r.SourceCommand, &r.CommandArgs,
		&r.FailureReason, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return r, fmt.Errorf(constants.ErrPendingTaskQuery, err)
	}

	return r, nil
}

// CompleteTask moves a pending task to the completed table in a transaction.
func (db *DB) CompleteTask(taskID int64) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf(constants.ErrPendingTaskComplete, err)
	}

	task, err := findPendingTaskInTx(tx, taskID)
	if err != nil {
		_ = tx.Rollback()

		return fmt.Errorf(constants.ErrPendingTaskComplete, err)
	}

	_, err = tx.Exec(constants.SQLInsertCompletedTask,
		task.ID, task.TaskTypeId, task.TargetPath, task.WorkingDirectory,
		task.SourceCommand, task.CommandArgs, task.CreatedAt)
	if err != nil {
		_ = tx.Rollback()

		return fmt.Errorf(constants.ErrPendingTaskComplete, err)
	}

	_, err = tx.Exec(constants.SQLDeletePendingTask, taskID)
	if err != nil {
		_ = tx.Rollback()

		return fmt.Errorf(constants.ErrPendingTaskComplete, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf(constants.ErrPendingTaskComplete, err)
	}

	return nil
}

// FailTask updates the failure reason for a pending task.
func (db *DB) FailTask(taskID int64, reason string) error {
	result, err := db.conn.Exec(constants.SQLUpdatePendingTaskFailure, reason, taskID)
	if err != nil {
		return fmt.Errorf(constants.ErrPendingTaskFail, err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf(constants.ErrPendingTaskNotFound, taskID)
	}

	return nil
}

// ListCompletedTasks returns all completed tasks ordered by completion time.
func (db *DB) ListCompletedTasks() ([]model.CompletedTaskRecord, error) {
	rows, err := db.conn.Query(constants.SQLSelectAllCompletedTasks)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrPendingTaskQuery, err)
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
	result, err := db.conn.Exec(constants.SQLDeletePendingTask, id)
	if err != nil {
		return fmt.Errorf(constants.ErrPendingTaskComplete, err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf(constants.ErrPendingTaskNotFound, id)
	}

	return nil
}

// findPendingTaskInTx reads a pending task within an existing transaction.
func findPendingTaskInTx(tx *sql.Tx, id int64) (model.PendingTaskRecord, error) {
	row := tx.QueryRow(constants.SQLSelectPendingTaskByID, id)

	var r model.PendingTaskRecord

	err := row.Scan(&r.ID, &r.TaskTypeId, &r.TaskTypeName, &r.TargetPath,
		&r.WorkingDirectory, &r.SourceCommand, &r.CommandArgs,
		&r.FailureReason, &r.CreatedAt, &r.UpdatedAt)

	return r, err
}
