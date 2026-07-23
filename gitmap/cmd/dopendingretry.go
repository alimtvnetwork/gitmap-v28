package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// retryPendingTask executes a single pending task based on its type.
func retryPendingTask(db *store.DB, taskID int64, typeName, targetPath, workDir, cmdArgs string) {
	if isDeleteTaskType(typeName) {
		retryDeleteTask(db, taskID, targetPath)

		return
	}

	if isReplayableTaskType(typeName) {
		retryReplayTask(db, taskID, workDir, cmdArgs)

		return
	}

	fmt.Fprintf(os.Stderr, constants.ErrTaskTypeNotFound, typeName)
}

// isDeleteTaskType returns true for Delete or Remove task types.
func isDeleteTaskType(typeName string) bool {
	return typeName == constants.TaskTypeDelete || typeName == constants.TaskTypeRemove
}

// isReplayableTaskType returns true for task types that can be replayed via CLI.
func isReplayableTaskType(typeName string) bool {
	return typeName == constants.TaskTypeScan ||
		typeName == constants.TaskTypeClone ||
		typeName == constants.TaskTypePull ||
		typeName == constants.TaskTypeExec
}

// retryDeleteTask attempts to delete the target path for a pending task.
func retryDeleteTask(db *store.DB, taskID int64, targetPath string) {
	if !pathExists(targetPath) {
		fmt.Printf(constants.MsgPendingSkipNotExist, taskID)
		completeTaskWithLog(db, taskID)

		return
	}

	err := os.RemoveAll(targetPath)
	if err != nil {
		reason := formatPathError(targetPath, "delete", err)
		fmt.Printf(constants.MsgPendingTaskFailed, taskID, reason)
		if failErr := db.FailTask(taskID, reason); failErr != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Could not mark task %d as failed: %v\n", taskID, failErr)
		}

		return
	}

	fmt.Printf(constants.MsgPendingTaskCompleted, taskID, targetPath)
	completeTaskWithLog(db, taskID)
}

// retryReplayTask re-executes a stored CLI command.
func retryReplayTask(db *store.DB, taskID int64, workDir, cmdArgs string) {
	args := strings.Fields(cmdArgs)
	if len(args) == 0 {
		reason := fmt.Sprintf(constants.ReasonReplayFailed, "empty command args")
		fmt.Printf(constants.MsgPendingTaskFailed, taskID, reason)
		if failErr := db.FailTask(taskID, reason); failErr != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Could not mark task %d as failed: %v\n", taskID, failErr)
		}

		return
	}

	if workDir != "" && !pathExists(workDir) {
		reason := fmt.Sprintf(constants.ReasonWorkDirNotFound, workDir)
		fmt.Printf(constants.MsgPendingTaskFailed, taskID, reason)
		if failErr := db.FailTask(taskID, reason); failErr != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Could not mark task %d as failed: %v\n", taskID, failErr)
		}

		return
	}

	fmt.Printf(constants.MsgPendingReplaying, cmdArgs)

	binaryPath, err := os.Executable()
	if err != nil {
		binaryPath = "gitmap"
	}

	cmd := exec.Command(binaryPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if workDir != "" {
		cmd.Dir = workDir
	}

	err = cmd.Run()
	if err != nil {
		reason := fmt.Sprintf(constants.ReasonReplayFailed, err)
		fmt.Printf(constants.MsgPendingTaskFailed, taskID, reason)
		if failErr := db.FailTask(taskID, reason); failErr != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Could not mark task %d as failed: %v\n", taskID, failErr)
		}

		return
	}

	fmt.Printf(constants.MsgPendingTaskCompleted, taskID, cmdArgs)
	completeTaskWithLog(db, taskID)
}

// pathExists returns true if the given path exists on the file system.
func pathExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

// formatPathError returns a Code Red formatted error with path context.
func formatPathError(path, operation string, err error) string {
	if errors.Is(err, os.ErrPermission) {
		return fmt.Sprintf(constants.ReasonPermissionDenied, path, operation, err)
	}

	return fmt.Sprintf(constants.ReasonRetryFailed, err)
}

// completeTaskWithLog marks a task complete and logs any transactional failure.
func completeTaskWithLog(db *store.DB, taskID int64) {
	err := db.CompleteTask(taskID)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnPendingCompleteFail, taskID, err)
	}
}
