package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runDoPending retries all pending tasks or a single task by ID.
func runDoPending(args []string) {
	checkHelp("do-pending", args)

	if len(args) > 0 {
		runDoPendingSingle(args[0])

		return
	}

	runDoPendingAll()
}

// runDoPendingAll retries all pending tasks.
func runDoPendingAll() {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnPendingDBOpen, err)
		os.Exit(1)
	}
	defer db.Close()

	tasks, err := db.ListPendingTasks()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrPendingTaskQuery, err)
		os.Exit(1)
	}

	if len(tasks) == 0 {
		fmt.Print(constants.MsgPendingListEmpty)

		return
	}

	fmt.Printf(constants.MsgPendingRetryAll, len(tasks))

	for _, t := range tasks {
		retryPendingTask(db, t.ID, t.TaskTypeName, t.TargetPath, t.WorkingDirectory, t.CommandArgs)
	}
}

// runDoPendingSingle retries a single pending task by its ID string.
func runDoPendingSingle(idStr string) {
	taskID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrPendingTaskNotFound, 0)
		os.Exit(1)
	}

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnPendingDBOpen, err)
		os.Exit(1)
	}
	defer db.Close()

	task, err := db.FindPendingTaskByID(taskID)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrPendingTaskNotFound, taskID)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgPendingRetryOne, taskID)
	retryPendingTask(db, task.ID, task.TaskTypeName, task.TargetPath, task.WorkingDirectory, task.CommandArgs)
}
