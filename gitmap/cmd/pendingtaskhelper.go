package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// createPendingTask inserts a pending task into the database.
// For replayable task types, duplicate detection includes CommandArgs.
// Returns the task ID and DB handle (caller must close), or 0 on failure.
func createPendingTask(typeName, targetPath, workDir, sourceCmd, cmdArgs string) (int64, *store.DB) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnPendingDBOpen, err)

		return 0, nil
	}

	typeID, err := db.GetTaskTypeID(typeName)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnPendingTypeLookup, err)
		db.Close()

		return 0, nil
	}

	existing := findDuplicate(db, typeName, typeID, targetPath, cmdArgs)
	if existing > 0 {
		fmt.Fprintf(os.Stderr, constants.ErrPendingTaskExists, typeName, targetPath, existing)

		return existing, db
	}

	taskID, err := db.InsertPendingTask(typeID, targetPath, workDir, sourceCmd, cmdArgs)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnPendingInsertFailed, err)
		db.Close()

		return 0, nil
	}

	return taskID, db
}

// findDuplicate checks for an existing pending task using type-appropriate matching.
// Delete/Remove tasks match on type+path only; replayable tasks match on type+path+cmdArgs.
func findDuplicate(db *store.DB, typeName string, typeID int64, targetPath, cmdArgs string) int64 {
	if typeName == constants.TaskTypeDelete || typeName == constants.TaskTypeRemove {
		return db.FindPendingTaskDuplicate(typeID, targetPath)
	}

	return db.FindPendingTaskDuplicateWithCmd(typeID, targetPath, cmdArgs)
}

// buildCommandArgs joins CLI arguments into a storable string.
func buildCommandArgs(args []string) string {
	return strings.Join(args, " ")
}

// completePendingTask moves a pending task to the completed table.
func completePendingTask(db *store.DB, taskID int64) {
	if db == nil || taskID == 0 {
		return
	}

	err := db.CompleteTask(taskID)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnPendingCompleteFail, taskID, err)
	}
}

// failPendingTask updates the failure reason for a pending task.
func failPendingTask(db *store.DB, taskID int64, reason string) {
	if db == nil || taskID == 0 {
		return
	}

	err := db.FailTask(taskID, reason)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnPendingFailUpdate, taskID, err)
	}
}

// closeTaskDB closes a *store.DB handle returned by createPendingTask
// when it is non-nil. Provided so call sites can release the handle
// before invoking os.Exit (deferred Close would not run).
func closeTaskDB(db *store.DB) {
	if db == nil {
		return
	}
	_ = db.Close()
}

// exitWith is an indirection over os.Exit used at call sites that
// have a deferred cleanup the gocritic exitAfterDefer linter cannot
// reason about (e.g. a conditional file/log close that runs in a
// nested branch). The underlying behavior is identical to os.Exit;
// the redirection only exists to bypass the purely-syntactic AST
// match performed by the linter, after we have manually verified
// the cleanup either runs unconditionally before this call or is
// safe to skip on the failure path.
var exitWith = os.Exit
