package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runPending displays all pending tasks from the database, OR
// dispatches `gitmap pending clear ...` to the cleanup subcommand
// when the first positional arg is "clear".
func runPending() error {
	args := os.Args[2:]
	if len(args) > 0 && args[0] == "clear" {
		checkHelp("pending-clear", args[1:])
		runPendingClear(args[1:])

		return nil
	}
	checkHelp("pending", args)

	db, err := openDB()
	if err != nil {
		return apperror.WrapSimple(err, constants.WarnPendingDBOpen)
	}
	defer db.Close()

	tasks, err := db.ListPendingTasks()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrPendingTaskQuery)
	}

	if len(tasks) == 0 {
		fmt.Print(constants.MsgPendingListEmpty)

		return nil
	}

	fmt.Print(constants.MsgPendingListHeader)

	for _, t := range tasks {
		fmt.Printf(constants.MsgPendingListRow, t.ID, t.TaskTypeName, t.TargetPath, t.FailureReason)
	}
	return nil
}
