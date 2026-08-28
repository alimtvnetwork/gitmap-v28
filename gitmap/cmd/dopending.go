package cmd

import (
	"fmt"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runDoPending retries all pending tasks or a single task by ID.
func runDoPending(args []string) error {
	checkHelp("do-pending", args)

	if len(args) > 0 {
		runDoPendingSingle(args[0])

		return nil
	}

	runDoPendingAll()
	return nil
}

// runDoPendingAll retries all pending tasks.
func runDoPendingAll() error {
	db, err := openDB()
	if err != nil {
		return apperror.Wrap(err, constants.WarnPendingDBOpen, nil)
	}
	defer db.Close()

	tasks, err := db.ListPendingTasks()
	if err != nil {
		return apperror.Wrap(err, constants.ErrPendingTaskQuery, nil)
	}

	if len(tasks) == 0 {
		fmt.Print(constants.MsgPendingListEmpty)

		return nil
	}

	fmt.Printf(constants.MsgPendingRetryAll, len(tasks))

	for _, t := range tasks {
		retryPendingTask(db, t.ID, t.TaskTypeName, t.TargetPath, t.WorkingDirectory, t.CommandArgs)
	}
	return nil
}

// runDoPendingSingle retries a single pending task by its ID string.
func runDoPendingSingle(idStr string) error {
	taskID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return apperror.New(constants.ErrPendingTaskNotFound, "E9000", nil)
	}

	db, err := openDB()
	if err != nil {
		return apperror.Wrap(err, constants.WarnPendingDBOpen, nil)
	}
	defer db.Close()

	task, err := db.FindPendingTaskByID(taskID)
	if err != nil {
		return apperror.New(constants.ErrPendingTaskNotFound, "E9000", nil)
	}

	fmt.Printf(constants.MsgPendingRetryOne, taskID)
	retryPendingTask(db, task.ID, task.TaskTypeName, task.TargetPath, task.WorkingDirectory, task.CommandArgs)
	return nil
}
