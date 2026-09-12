package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// createPendingTask inserts a pending task into the database.
// For replayable task types, duplicate detection includes CommandArgs.
// Returns the task ID and DB handle (caller must close), or 0 on failure.
func createPendingTask(
	typeName,
	targetPath,
	workDir,
	sourceCmd,
	cmdArgs string,
) (int64, *store.DB) {
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
		db.Close()

		return existing, nil
	}

	taskID, err := db.InsertPendingTask(typeID, targetPath, workDir, sourceCmd, cmdArgs)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnPendingInsertFailed, err)
		db.Close()

		return 0, nil
	}

	db.Close()

	return taskID, nil
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

func completePendingTask(db *store.DB, taskID int64) {
	if taskID == 0 {
		return
	}

	activeDB, cleanup := ensureDB(db, "complete", taskID)
	if activeDB == nil {
		return
	}

	defer cleanup()

	err := activeDB.CompleteTask(taskID)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnPendingCompleteFail, taskID, err)
	}
}

func failPendingTask(db *store.DB, taskID int64, reason string) {
	if taskID == 0 {
		return
	}

	activeDB, cleanup := ensureDB(db, "fail", taskID)
	if activeDB == nil {
		return
	}

	defer cleanup()

	err := activeDB.FailTask(taskID, reason)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnPendingFailUpdate, taskID, err)
	}
}

func ensureDB(db *store.DB, action string, taskID int64) (*store.DB, func()) {
	if db != nil {
		return db, func() {}
	}

	opened, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not open database to %s pending task %d: %v\n", action, taskID, err)

		return nil, nil
	}

	return opened, func() { opened.Close() }
}

// closeTaskDB closes a *store.DB handle returned by createPendingTask
// when it is non-nil.
func closeTaskDB(db *store.DB) {
	if db == nil {
		return
	}

	_ = db.Close()
}
