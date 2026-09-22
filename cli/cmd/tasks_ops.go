// Package cmd — tasks_ops.go implements undo and redo operations for tasks.
package cmd

import (
	"fmt"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func runTasksUndo(args []string) error {
	db, err := openTasksDB()
	if err != nil {
		return apperror.WrapSimple(err, constants.WarnPendingDBOpen)
	}

	defer db.Close()

	task, findErr := findTaskToRevert(db, args)
	if findErr != nil {
		return findErr
	}

	return executeTaskUndo(db, task)
}

func findTaskToRevert(db *store.DB, args []string) (*model.CompletedTaskRecord, error) {
	completed, err := db.ListCompletedTasks()
	if err != nil {
		return nil, apperror.WrapSimple(err, "list completed tasks failed")
	}
	if len(completed) == 0 {
		return nil, apperror.NewSimple("no completed tasks available to undo", "E1096")
	}

	if len(args) == 0 {
		return &completed[0], nil
	}

	matched := findCompletedTaskByID(completed, args[0])
	if matched != nil {
		return matched, nil
	}

	return &completed[0], nil
}

func findCompletedTaskByID(completed []model.CompletedTaskRecord, idStr string) *model.CompletedTaskRecord {
	id, parseErr := strconv.ParseInt(idStr, 10, 64)
	if parseErr != nil {
		return nil
	}

	for _, t := range completed {
		if t.ID == id {
			return &t
		}
	}

	return nil
}

func executeTaskUndo(db *store.DB, t *model.CompletedTaskRecord) error {
	_, insErr := db.InsertPendingTask(t.TaskTypeId, t.TargetPath, t.WorkingDirectory, t.SourceCommand, t.CommandArgs)
	if insErr != nil {
		return apperror.WrapSimple(insErr, "re-insert pending task")
	}

	_, _ = store.ExecWrapper(db.Conn(), "DELETE FROM CompletedTask WHERE CompletedTaskId = ?", t.ID).Destruct()
	fmt.Printf("  %s✓ Task %d (%s: %s) reverted to pending tasks queue.%s\n\n",
		constants.ColorGreen, t.ID, t.TaskTypeName, t.SourceCommand, constants.ColorReset)

	return nil
}

func runTasksRedo(args []string) error {
	db, err := openTasksDB()
	if err != nil {
		return apperror.WrapSimple(err, constants.WarnPendingDBOpen)
	}

	defer db.Close()

	task, findErr := findTaskToRevert(db, args)
	if findErr != nil {
		return findErr
	}

	return executeTaskRedo(db, task)
}

func executeTaskRedo(db *store.DB, t *model.CompletedTaskRecord) error {
	taskID, insErr := db.InsertPendingTask(t.TaskTypeId, t.TargetPath, t.WorkingDirectory, t.SourceCommand, t.CommandArgs)
	if insErr != nil {
		return apperror.WrapSimple(insErr, "queue task replay")
	}

	fmt.Printf("  %s✓ Task %d queued for re-execution (New Task ID: %d): %s %s%s\n\n",
		constants.ColorGreen, t.ID, taskID, t.SourceCommand, t.CommandArgs, constants.ColorReset)

	return nil
}
