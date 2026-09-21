package cmdssh

import (
	"context"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func restoreNodesToDB(ctx context.Context, nodes []store.SSHHost) error {
	dbConn, err := openSSHDBFunc()
	hasErr := err != nil
	if hasErr {
		return apperror.WrapSimple(err, "restoreNodesToDB.openDB")
	}
	defer dbConn.Close()

	for _, node := range nodes {
		errUpsert := store.UpsertSSHHost(ctx, node, dbConn.SQL())
		hasUpsertErr := errUpsert != nil
		if hasUpsertErr {
			return apperror.WrapSimple(errUpsert, "restoreNodesToDB.Upsert")
		}
	}

	return nil
}

func printRestoreSuccess(task *SSHHistoryTask) {
	fmt.Printf("✓ Restored %d SSH node(s) from task '%s' (%s: %s).\n\n",
		len(task.Nodes), task.TaskID, task.Action, task.Target)
	_ = RenderSSHHostsTable(os.Stdout, task.Nodes)
}

func executeTaskRestore(ctx context.Context, task *SSHHistoryTask) error {
	errRestore := restoreNodesToDB(ctx, task.Nodes)
	hasRestoreErr := errRestore != nil
	if hasRestoreErr {
		return errRestore
	}

	_ = MarkSSHHistoryTaskRestored(ctx, task.TaskID)
	printRestoreSuccess(task)

	return nil
}

// RunSSHUndoCLI restores nodes from the most recent destructive operation.
func RunSSHUndoCLI(args []string) error {
	ctx := context.Background()
	task, err := GetLastSSHHistoryTask(ctx)
	hasErr := err != nil
	if hasErr {
		fmt.Println("No reversible SSH node operations found.")

		return nil
	}

	return executeTaskRestore(ctx, task)
}

// RunSSHRestoreCLI restores nodes from a specific history task ID.
func RunSSHRestoreCLI(args []string) error {
	hasArgs := len(args) > 0
	if !hasArgs {
		return apperror.NewValidationError("missing task ID to restore\n\nUsage: gitmap ssh restore <task-id>")
	}

	ctx := context.Background()
	task, err := GetSSHHistoryTaskByID(ctx, args[0])
	hasErr := err != nil
	if hasErr {
		return apperror.NewValidationError(fmt.Sprintf("task '%s' not found in SSH history", args[0]))
	}

	return executeTaskRestore(ctx, task)
}
