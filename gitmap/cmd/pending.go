package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runPending displays all pending tasks from the database, OR
// dispatches `gitmap pending clear ...` to the cleanup subcommand
// when the first positional arg is "clear".
func runPending() {
	args := os.Args[2:]
	if len(args) > 0 && args[0] == "clear" {
		checkHelp("pending-clear", args[1:])
		runPendingClear(args[1:])

		return
	}
	checkHelp("pending", args)

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

	fmt.Print(constants.MsgPendingListHeader)

	for _, t := range tasks {
		fmt.Printf(constants.MsgPendingListRow, t.ID, t.TaskTypeName, t.TargetPath, t.FailureReason)
	}
}
