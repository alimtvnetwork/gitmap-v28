package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// openTestDB creates a temp DB with migrations and seeded task types.
func openTestDB(t *testing.T) *store.DB {
	t.Helper()
	dir := t.TempDir()
	db, err := store.Open(dir)
	if err != nil {
		t.Fatalf("failed to open test DB: %v", err)
	}
	if err := db.Migrate(); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

// --- buildCommandArgs ---

func TestBuildCommandArgs_JoinsArgs(t *testing.T) {
	result := buildCommandArgs([]string{"scan", "--all", "/repos"})
	if result != "scan --all /repos" {
		t.Errorf("expected 'scan --all /repos', got %q", result)
	}
}

func TestBuildCommandArgs_Empty(t *testing.T) {
	result := buildCommandArgs([]string{})
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestBuildCommandArgs_SingleArg(t *testing.T) {
	result := buildCommandArgs([]string{"pull"})
	if result != "pull" {
		t.Errorf("expected 'pull', got %q", result)
	}
}

// --- findDuplicate ---

func TestFindDuplicate_DeleteMatchesByTypePath(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID, _ := db.GetTaskTypeID(constants.TaskTypeDelete)
	db.InsertPendingTask(typeID, "/repos/old", "", "clone-next", "")

	dup := findDuplicate(db, constants.TaskTypeDelete, typeID, "/repos/old", "")
	if dup == 0 {
		t.Error("expected duplicate found for Delete type")
	}
}

func TestFindDuplicate_DeleteIgnoresCmdArgs(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID, _ := db.GetTaskTypeID(constants.TaskTypeDelete)
	id, _ := db.InsertPendingTask(typeID, "/repos/old", "", "clone-next", "some args")

	dup := findDuplicate(db, constants.TaskTypeDelete, typeID, "/repos/old", "different args")
	if dup != id {
		t.Errorf("expected Delete to match by path only, got dup=%d", dup)
	}
}

func TestFindDuplicate_RemoveMatchesByTypePath(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID, _ := db.GetTaskTypeID(constants.TaskTypeRemove)
	id, _ := db.InsertPendingTask(typeID, "/repos/old", "", "remove", "")

	dup := findDuplicate(db, constants.TaskTypeRemove, typeID, "/repos/old", "")
	if dup != id {
		t.Error("expected duplicate found for Remove type")
	}
}

func TestFindDuplicate_ScanMatchesByTypePathCmd(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID, _ := db.GetTaskTypeID(constants.TaskTypeScan)
	id, _ := db.InsertPendingTask(typeID, "/repos", "/home", "scan", "scan --all")

	dup := findDuplicate(db, constants.TaskTypeScan, typeID, "/repos", "scan --all")
	if dup != id {
		t.Error("expected duplicate found for Scan with same args")
	}
}

func TestFindDuplicate_ScanDifferentArgsNoDuplicate(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID, _ := db.GetTaskTypeID(constants.TaskTypeScan)
	db.InsertPendingTask(typeID, "/repos", "/home", "scan", "scan --mode ssh")

	dup := findDuplicate(db, constants.TaskTypeScan, typeID, "/repos", "scan --mode https")
	if dup != 0 {
		t.Errorf("expected no duplicate for different args, got %d", dup)
	}
}

func TestFindDuplicate_CloneMatchesByTypePathCmd(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID, _ := db.GetTaskTypeID(constants.TaskTypeClone)
	id, _ := db.InsertPendingTask(typeID, "/repos/app", "/work", "clone", "clone src.json")

	dup := findDuplicate(db, constants.TaskTypeClone, typeID, "/repos/app", "clone src.json")
	if dup != id {
		t.Error("expected duplicate found for Clone with same args")
	}
}

func TestFindDuplicate_NoDuplicateOnDifferentPath(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID, _ := db.GetTaskTypeID(constants.TaskTypeScan)
	db.InsertPendingTask(typeID, "/repos/a", "/home", "scan", "scan /repos/a")

	dup := findDuplicate(db, constants.TaskTypeScan, typeID, "/repos/b", "scan /repos/a")
	if dup != 0 {
		t.Errorf("expected no duplicate for different path, got %d", dup)
	}
}

// --- completePendingTask ---

func TestCompletePendingTask_MovesToCompleted(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID, _ := db.GetTaskTypeID(constants.TaskTypeScan)
	taskID, _ := db.InsertPendingTask(typeID, "/repos", "/home", "scan", "scan /repos")

	completePendingTask(db, taskID)

	// Should be gone from pending.
	_, err := db.FindPendingTaskByID(taskID)
	if err == nil {
		t.Error("expected task removed from pending after completion")
	}

	// Should exist in completed.
	completed, _ := db.ListCompletedTasks()
	if len(completed) != 1 {
		t.Fatalf("expected 1 completed, got %d", len(completed))
	}
	if completed[0].OriginalTaskId != taskID {
		t.Errorf("expected original ID %d, got %d", taskID, completed[0].OriginalTaskId)
	}
}

func TestCompletePendingTask_NilDB_NoOp(t *testing.T) {
	// Should not panic.
	completePendingTask(nil, 42)
}

func TestCompletePendingTask_ZeroID_NoOp(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	// Should not panic.
	completePendingTask(db, 0)
}

// --- failPendingTask ---

func TestFailPendingTask_RecordsReason(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID, _ := db.GetTaskTypeID(constants.TaskTypeClone)
	taskID, _ := db.InsertPendingTask(typeID, "/repos/app", "/work", "clone", "clone src.json")

	failPendingTask(db, taskID, "permission denied at /repos/app")

	rec, _ := db.FindPendingTaskByID(taskID)
	if rec.FailureReason != "permission denied at /repos/app" {
		t.Errorf("expected failure reason, got %q", rec.FailureReason)
	}
}

func TestFailPendingTask_TaskRemainsPending(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID, _ := db.GetTaskTypeID(constants.TaskTypeExec)
	taskID, _ := db.InsertPendingTask(typeID, "/work", "/work", "exec", "exec fetch --prune")

	failPendingTask(db, taskID, "exec failed")

	pending, _ := db.ListPendingTasks()
	if len(pending) != 1 {
		t.Fatalf("expected 1 pending, got %d", len(pending))
	}
	if pending[0].ID != taskID {
		t.Errorf("expected ID %d, got %d", taskID, pending[0].ID)
	}
}

func TestFailPendingTask_NilDB_NoOp(t *testing.T) {
	// Should not panic.
	failPendingTask(nil, 42, "some reason")
}

func TestFailPendingTask_ZeroID_NoOp(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	// Should not panic.
	failPendingTask(db, 0, "some reason")
}

// --- End-to-end: create → fail → complete lifecycle ---

func TestPendingTask_FullLifecycle_FailThenComplete(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID, _ := db.GetTaskTypeID(constants.TaskTypePull)
	taskID, _ := db.InsertPendingTask(typeID, "/repos/lib", "/home", "pull", "pull mylib --verbose")

	// First attempt fails.
	failPendingTask(db, taskID, "network timeout")
	rec, _ := db.FindPendingTaskByID(taskID)
	if rec.FailureReason != "network timeout" {
		t.Errorf("expected reason 'network timeout', got %q", rec.FailureReason)
	}

	// Retry succeeds.
	completePendingTask(db, taskID)
	_, err := db.FindPendingTaskByID(taskID)
	if err == nil {
		t.Error("expected task removed from pending after completion")
	}

	completed, _ := db.ListCompletedTasks()
	if len(completed) != 1 {
		t.Fatalf("expected 1 completed, got %d", len(completed))
	}
	if completed[0].CommandArgs != "pull mylib --verbose" {
		t.Errorf("expected cmdargs preserved, got %q", completed[0].CommandArgs)
	}
}
