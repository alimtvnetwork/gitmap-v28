// Package store_test — unit tests for PendingTask and CompletedTask CRUD.
package store_test

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// --- Insert ---

// TestInsertPendingTask_ReturnsID verifies a new pending task gets a valid ID.
func TestInsertPendingTask_ReturnsID(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	id, err := db.InsertPendingTask(scanTypeID(t, db), "/repos/myapp", "/home", "scan", "scan --all")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id == 0 {
		t.Error("expected non-zero task ID")
	}
}

// TestInsertPendingTask_MultipleGetUniqueIDs verifies sequential inserts get unique IDs.
func TestInsertPendingTask_MultipleGetUniqueIDs(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID := scanTypeID(t, db)
	id1, _ := db.InsertPendingTask(typeID, "/repos/a", "/home", "scan", "scan /repos/a")
	id2, _ := db.InsertPendingTask(typeID, "/repos/b", "/home", "scan", "scan /repos/b")

	if id1 == id2 {
		t.Errorf("expected unique IDs, got %d and %d", id1, id2)
	}
}

// --- FindByID ---

// TestFindPendingTaskByID_ReturnsCorrectFields verifies all fields are stored.
func TestFindPendingTaskByID_ReturnsCorrectFields(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID := cloneTypeID(t, db)
	id, _ := db.InsertPendingTask(typeID, "/repos/app", "/work", "clone", "clone source.json /repos/app")

	rec, err := db.FindPendingTaskByID(id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rec.TargetPath != "/repos/app" {
		t.Errorf("expected target /repos/app, got %s", rec.TargetPath)
	}
	if rec.WorkingDirectory != "/work" {
		t.Errorf("expected workdir /work, got %s", rec.WorkingDirectory)
	}
	if rec.CommandArgs != "clone source.json /repos/app" {
		t.Errorf("expected cmdargs, got %s", rec.CommandArgs)
	}
	if rec.TaskTypeName != constants.TaskTypeClone {
		t.Errorf("expected type Clone, got %s", rec.TaskTypeName)
	}
	if rec.CreatedAt == "" {
		t.Error("expected CreatedAt to be populated")
	}
}

// TestFindPendingTaskByID_NotFound returns error for missing ID.
func TestFindPendingTaskByID_NotFound(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	_, err := db.FindPendingTaskByID(9999)
	if err == nil {
		t.Error("expected error for non-existent ID")
	}
}

// --- Duplicate Detection ---

// TestFindPendingTaskDuplicate_ExactMatch detects duplicate by type+path.
func TestFindPendingTaskDuplicate_ExactMatch(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID := deleteTypeID(t, db)
	id, _ := db.InsertPendingTask(typeID, "/repos/old", "", "clone-next", "")

	dup := db.FindPendingTaskDuplicate(typeID, "/repos/old")
	if dup != id {
		t.Errorf("expected duplicate ID %d, got %d", id, dup)
	}
}

// TestFindPendingTaskDuplicate_NoMatch returns 0 for different path.
func TestFindPendingTaskDuplicate_NoMatch(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID := deleteTypeID(t, db)
	db.InsertPendingTask(typeID, "/repos/old", "", "clone-next", "")

	dup := db.FindPendingTaskDuplicate(typeID, "/repos/different")
	if dup != 0 {
		t.Errorf("expected 0, got %d", dup)
	}
}

// TestFindPendingTaskDuplicateWithCmd_MatchesAll detects duplicate by type+path+cmd.
func TestFindPendingTaskDuplicateWithCmd_MatchesAll(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID := scanTypeID(t, db)
	id, _ := db.InsertPendingTask(typeID, "/repos", "/home", "scan", "scan --all")

	dup := db.FindPendingTaskDuplicateWithCmd(typeID, "/repos", "scan --all")
	if dup != id {
		t.Errorf("expected duplicate ID %d, got %d", id, dup)
	}
}

// TestFindPendingTaskDuplicateWithCmd_DifferentArgs allows same path with different args.
func TestFindPendingTaskDuplicateWithCmd_DifferentArgs(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID := scanTypeID(t, db)
	db.InsertPendingTask(typeID, "/repos", "/home", "scan", "scan --mode ssh")

	dup := db.FindPendingTaskDuplicateWithCmd(typeID, "/repos", "scan --mode https")
	if dup != 0 {
		t.Errorf("expected 0 for different args, got %d", dup)
	}
}

// --- CompleteTask ---

// TestCompleteTask_MovesToCompleted verifies task moves from pending to completed.
func TestCompleteTask_MovesToCompleted(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID := scanTypeID(t, db)
	id, _ := db.InsertPendingTask(typeID, "/repos/app", "/home", "scan", "scan /repos/app")

	err := db.CompleteTask(id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Pending should be gone.
	_, err = db.FindPendingTaskByID(id)
	if err == nil {
		t.Error("expected pending task to be removed after completion")
	}

	// Completed should exist.
	completed, err := db.ListCompletedTasks()
	if err != nil {
		t.Fatalf("failed to list completed: %v", err)
	}
	if len(completed) != 1 {
		t.Fatalf("expected 1 completed task, got %d", len(completed))
	}
	if completed[0].OriginalTaskId != id {
		t.Errorf("expected original ID %d, got %d", id, completed[0].OriginalTaskId)
	}
	if completed[0].TargetPath != "/repos/app" {
		t.Errorf("expected target /repos/app, got %s", completed[0].TargetPath)
	}
	if completed[0].CommandArgs != "scan /repos/app" {
		t.Errorf("expected cmdargs preserved, got %s", completed[0].CommandArgs)
	}
}

// TestCompleteTask_NonExistentID returns error.
func TestCompleteTask_NonExistentID(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	err := db.CompleteTask(9999)
	if err == nil {
		t.Error("expected error for non-existent task ID")
	}
}

// --- FailTask ---

// TestFailTask_UpdatesReason verifies failure reason is stored.
func TestFailTask_UpdatesReason(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID := cloneTypeID(t, db)
	id, _ := db.InsertPendingTask(typeID, "/repos/app", "/work", "clone", "clone src.json")

	err := db.FailTask(id, "permission denied at /repos/app")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	rec, _ := db.FindPendingTaskByID(id)
	if rec.FailureReason != "permission denied at /repos/app" {
		t.Errorf("expected failure reason, got %q", rec.FailureReason)
	}
}

// TestFailTask_NonExistentID returns error.
func TestFailTask_NonExistentID(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	err := db.FailTask(9999, "some reason")
	if err == nil {
		t.Error("expected error for non-existent task ID")
	}
}

// TestFailTask_TaskRemainsInPending verifies failed task stays pending.
func TestFailTask_TaskRemainsInPending(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID := scanTypeID(t, db)
	id, _ := db.InsertPendingTask(typeID, "/repos", "/home", "scan", "scan /repos")
	db.FailTask(id, "scan error")

	pending, err := db.ListPendingTasks()
	if err != nil {
		t.Fatalf("failed to list: %v", err)
	}
	if len(pending) != 1 {
		t.Fatalf("expected 1 pending, got %d", len(pending))
	}
	if pending[0].ID != id {
		t.Errorf("expected ID %d, got %d", id, pending[0].ID)
	}
}

// --- ListPendingTasks ---

// TestListPendingTasks_Empty returns empty slice on fresh DB.
func TestListPendingTasks_Empty(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	tasks, err := db.ListPendingTasks()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}
}

// TestListPendingTasks_ReturnsAll lists all pending tasks.
func TestListPendingTasks_ReturnsAll(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	scanID := scanTypeID(t, db)
	cloneID := cloneTypeID(t, db)
	db.InsertPendingTask(scanID, "/repos/a", "/home", "scan", "scan /repos/a")
	db.InsertPendingTask(cloneID, "/repos/b", "/work", "clone", "clone src.json")

	tasks, err := db.ListPendingTasks()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

// --- ListCompletedTasks ---

// TestListCompletedTasks_AfterComplete returns completed tasks.
func TestListCompletedTasks_AfterComplete(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID := scanTypeID(t, db)
	id, _ := db.InsertPendingTask(typeID, "/repos", "/home", "scan", "scan /repos")
	db.CompleteTask(id)

	completed, err := db.ListCompletedTasks()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(completed) != 1 {
		t.Errorf("expected 1 completed, got %d", len(completed))
	}
	if completed[0].CompletedAt == "" {
		t.Error("expected CompletedAt to be set")
	}
}

// --- Pull type ---

// TestInsertPendingTask_PullType verifies Pull task type works correctly.
func TestInsertPendingTask_PullType(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID := pullTypeID(t, db)
	id, err := db.InsertPendingTask(typeID, "/repos/mylib", "/home", "pull", "pull mylib --verbose")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	rec, _ := db.FindPendingTaskByID(id)
	if rec.TaskTypeName != constants.TaskTypePull {
		t.Errorf("expected type Pull, got %s", rec.TaskTypeName)
	}
	if rec.CommandArgs != "pull mylib --verbose" {
		t.Errorf("expected cmdargs, got %s", rec.CommandArgs)
	}
}

// --- Helpers ---

func scanTypeID(t *testing.T, db interface{ GetTaskTypeID(string) (int64, error) }) int64 {
	t.Helper()
	id, err := db.GetTaskTypeID(constants.TaskTypeScan)
	if err != nil {
		t.Fatalf("failed to get Scan type ID: %v", err)
	}
	return id
}

func cloneTypeID(t *testing.T, db interface{ GetTaskTypeID(string) (int64, error) }) int64 {
	t.Helper()
	id, err := db.GetTaskTypeID(constants.TaskTypeClone)
	if err != nil {
		t.Fatalf("failed to get Clone type ID: %v", err)
	}
	return id
}

func deleteTypeID(t *testing.T, db interface{ GetTaskTypeID(string) (int64, error) }) int64 {
	t.Helper()
	id, err := db.GetTaskTypeID(constants.TaskTypeDelete)
	if err != nil {
		t.Fatalf("failed to get Delete type ID: %v", err)
	}
	return id
}

func pullTypeID(t *testing.T, db interface{ GetTaskTypeID(string) (int64, error) }) int64 {
	t.Helper()
	id, err := db.GetTaskTypeID(constants.TaskTypePull)
	if err != nil {
		t.Fatalf("failed to get Pull type ID: %v", err)
	}
	return id
}
