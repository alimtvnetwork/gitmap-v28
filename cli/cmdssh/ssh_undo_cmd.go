package cmdssh

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func restoreNodesToDB(ctx context.Context, nodes []store.SSHHost) error {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return apperror.WrapSimple(err, "restoreNodesToDB.openDB")
	}
	defer dbConn.Close()

	for _, node := range nodes {
		errUpsert := store.UpsertSSHHost(ctx, node, dbConn.SQL())
		if errUpsert != nil {
			return apperror.WrapSimple(errUpsert, "restoreNodesToDB.Upsert")
		}
	}
	return nil
}

func restoreKeyToDB(inversePayload string) error {
	var key model.SSHKey
	if err := json.Unmarshal([]byte(inversePayload), &key); err != nil {
		return apperror.WrapSimple(err, "restoreKeyToDB.Unmarshal")
	}

	dbConn, err := openSSHDBFunc()
	if err != nil {
		return apperror.WrapSimple(err, "restoreKeyToDB.openDB")
	}
	defer dbConn.Close()

	_, errInsert := dbConn.InsertSSHKey(key.Name, key.PrivatePath, key.PublicKey, key.Fingerprint, key.Email)
	if errInsert != nil {
		return apperror.WrapSimple(errInsert, "restoreKeyToDB.InsertSSHKey")
	}
	return nil
}

func removeAddedNodeFromDB(ctx context.Context, target string) error {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return apperror.WrapSimple(err, "removeAddedNodeFromDB.openDB")
	}
	defer dbConn.Close()

	_, _ = store.DeleteHostByAliasOrIP(ctx, target, dbConn.SQL())
	_ = db.DeleteSSHConnectionByTarget(ctx, dbConn.SQL(), target)
	return nil
}

func dispatchActionRestore(ctx context.Context, task *SSHHistoryTask) error {
	if task.Action == ActionAddNode {
		return removeAddedNodeFromDB(ctx, task.Target)
	}
	if task.Action == ActionRmKey {
		return restoreKeyToDB(task.InversePayload)
	}
	return restoreNodesToDB(ctx, task.Nodes)
}

func printRestoreSuccess(task *SSHHistoryTask) {
	fmt.Printf("✓ Restored %d SSH node(s) from task '%s' (%s: %s).\n\n",
		len(task.Nodes), task.TaskID, task.Action, task.Target)
	_ = RenderSSHHostsTable(os.Stdout, task.Nodes)
}

func executeTaskRestore(ctx context.Context, task *SSHHistoryTask) error {
	errRestore := dispatchActionRestore(ctx, task)
	if errRestore != nil {
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
	if err != nil {
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
	if err != nil {
		return apperror.NewValidationError(fmt.Sprintf("task '%s' not found in SSH history", args[0]))
	}

	return executeTaskRestore(ctx, task)
}

func reAddNodeToDB(ctx context.Context, forwardPayload string) error {
	var host store.SSHHost
	if err := json.Unmarshal([]byte(forwardPayload), &host); err != nil {
		return apperror.WrapSimple(err, "reAddNodeToDB.Unmarshal")
	}

	dbConn, err := openSSHDBFunc()
	if err != nil {
		return apperror.WrapSimple(err, "reAddNodeToDB.openDB")
	}
	defer dbConn.Close()

	return store.UpsertSSHHost(ctx, host, dbConn.SQL())
}

func reDeleteKeyFromDB(keyName string) error {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return apperror.WrapSimple(err, "reDeleteKeyFromDB.openDB")
	}
	defer dbConn.Close()

	return dbConn.DeleteSSHKey(keyName)
}

func reDeleteNodesFromDB(ctx context.Context, nodes []store.SSHHost) error {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return apperror.WrapSimple(err, "reDeleteNodesFromDB.openDB")
	}
	defer dbConn.Close()

	return purgeCandidateNodes(ctx, dbConn.SQL(), nodes)
}

func dispatchActionRedo(ctx context.Context, task *SSHHistoryTask) error {
	if task.Action == ActionAddNode {
		return reAddNodeToDB(ctx, task.ForwardPayload)
	}
	if task.Action == ActionRmKey {
		return reDeleteKeyFromDB(task.Target)
	}
	return reDeleteNodesFromDB(ctx, task.Nodes)
}

func printRedoSuccess(task *SSHHistoryTask) {
	fmt.Printf("✓ Redone operation '%s' for task '%s' (%s).\n",
		task.Action, task.TaskID, task.Target)
}

func executeTaskRedo(ctx context.Context, task *SSHHistoryTask) error {
	errRedo := dispatchActionRedo(ctx, task)
	if errRedo != nil {
		return errRedo
	}

	_ = MarkSSHHistoryTaskActive(ctx, task.TaskID)
	printRedoSuccess(task)

	return nil
}

// RunSSHRedoCLI re-executes the most recently undone operation.
func RunSSHRedoCLI(args []string) error {
	ctx := context.Background()
	task, err := GetLastRestoredSSHHistoryTask(ctx)
	if err != nil {
		fmt.Println("No reversible SSH operations found to redo.")
		return nil
	}

	return executeTaskRedo(ctx, task)
}
