// Package cmd — tasks_list.go renders terminal tables for pending and completed tasks.
package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdtask"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func runTasksList() error {
	db, err := openTasksDB()
	if err != nil {
		return apperror.WrapSimple(err, constants.WarnPendingDBOpen)
	}

	defer db.Close()

	pending, _ := db.ListPendingTasks()
	completed, _ := db.ListCompletedTasks()

	printPendingTasksTable(pending)
	printRecentCompletedTasksTable(completed)

	return nil
}

func printPendingTasksTable(tasks []model.PendingTaskRecord) {
	fmt.Println()
	fmt.Printf("  %s⏳ Pending Tasks Queue (%d)%s\n", constants.ColorYellow, len(tasks), constants.ColorReset)
	if len(tasks) == 0 {
		fmt.Printf("    %s(no pending tasks in queue)%s\n\n", constants.ColorDim, constants.ColorReset)

		return
	}

	fmt.Printf("    %-6s %-12s %-30s %s\n", "ID", "TYPE", "TARGET", "COMMAND")
	fmt.Printf("    %s\n", constants.TermTableRule)
	for _, t := range tasks {
		cmdStr := t.SourceCommand + " " + t.CommandArgs
		fmt.Printf("    %-6d %-12s %-30s %s\n", t.ID, t.TaskTypeName, t.TargetPath, cmdStr)
	}
	fmt.Println()
}

func printRecentCompletedTasksTable(tasks []model.CompletedTaskRecord) {
	limit := len(tasks)
	if limit > 10 {
		limit = 10
	}

	fmt.Printf("  %s✔ Recently Completed Tasks (showing %d of %d)%s\n",
		constants.ColorGreen, limit, len(tasks), constants.ColorReset)
	if limit == 0 {
		fmt.Printf("    %s(no completed tasks recorded yet)%s\n\n", constants.ColorDim, constants.ColorReset)

		return
	}

	fmt.Printf("    %-6s %-12s %-26s %-20s %s\n", "TASK", "TYPE", "COMPLETED_AT", "SOURCE_CMD", "ARGS")
	fmt.Printf("    %s\n", constants.TermTableRule)
	for i := 0; i < limit; i++ {
		t := tasks[i]
		fmt.Printf("    %-6d %-12s %-26s %-20s %s\n",
			t.ID, t.TaskTypeName, t.CompletedAt, t.SourceCommand, t.CommandArgs)
	}
	fmt.Println()
}

func runTasksHistory(args []string) error {
	return cmdtask.RunTaskHistory(args)
}
