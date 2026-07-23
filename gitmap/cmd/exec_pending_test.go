package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// --- Exec enqueue ---

func TestExecPending_EnqueueCreatesTask(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID, _ := db.GetTaskTypeID(constants.TaskTypeExec)
	taskID, _ := db.InsertPendingTask(typeID, "/work", "/work", "exec", "exec fetch --prune")

	if taskID == 0 {
		t.Fatal("expected non-zero task ID for Exec enqueue")
	}

	rec, err := db.FindPendingTaskByID(taskID)
	if err != nil {
		t.Fatalf("expected task in pending table: %v", err)
	}
	if rec.TaskTypeName != constants.TaskTypeExec {
		t.Errorf("expected type %q, got %q", constants.TaskTypeExec, rec.TaskTypeName)
	}
	if rec.TargetPath != "/work" {
		t.Errorf("expected target /work, got %q", rec.TargetPath)
	}
	if rec.CommandArgs != "exec fetch --prune" {
		t.Errorf("expected cmdArgs 'exec fetch --prune', got %q", rec.CommandArgs)
	}
}

func TestExecPending_EnqueuePreservesWorkDir(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID, _ := db.GetTaskTypeID(constants.TaskTypeExec)
	taskID, _ := db.InsertPendingTask(typeID, "/projects", "/projects", "exec", "exec status")

	rec, _ := db.FindPendingTaskByID(taskID)
	if rec.WorkingDirectory != "/projects" {
		t.Errorf("expected workDir /projects, got %q", rec.WorkingDirectory)
	}
}

// --- Exec duplicate detection ---

func TestExecPending_DuplicateWithSameArgs(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID, _ := db.GetTaskTypeID(constants.TaskTypeExec)
	id, _ := db.InsertPendingTask(typeID, "/work", "/work", "exec", "exec fetch --all")

	dup := findDuplicate(db, constants.TaskTypeExec, typeID, "/work", "exec fetch --all")
	if dup != id {
		t.Errorf("expected duplicate %d, got %d", id, dup)
	}
}

func TestExecPending_NoDuplicateWithDifferentArgs(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID, _ := db.GetTaskTypeID(constants.TaskTypeExec)
	db.InsertPendingTask(typeID, "/work", "/work", "exec", "exec fetch --all")

	dup := findDuplicate(db, constants.TaskTypeExec, typeID, "/work", "exec status")
	if dup != 0 {
		t.Errorf("expected no duplicate for different args, got %d", dup)
	}
}

// --- Exec completion ---

func TestExecPending_CompleteMovesToCompleted(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID, _ := db.GetTaskTypeID(constants.TaskTypeExec)
	taskID, _ := db.InsertPendingTask(typeID, "/work", "/work", "exec", "exec fetch --prune")

	completePendingTask(db, taskID)

	_, err := db.FindPendingTaskByID(taskID)
	if err == nil {
		t.Error("expected task removed from pending after completion")
	}

	completed, _ := db.ListCompletedTasks()
	if len(completed) != 1 {
		t.Fatalf("expected 1 completed, got %d", len(completed))
	}
	if completed[0].OriginalTaskId != taskID {
		t.Errorf("expected original ID %d, got %d", taskID, completed[0].OriginalTaskId)
	}
	if completed[0].CommandArgs != "exec fetch --prune" {
		t.Errorf("expected cmdArgs preserved, got %q", completed[0].CommandArgs)
	}
}

// --- Exec failure ---

func TestExecPending_FailRecordsReason(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID, _ := db.GetTaskTypeID(constants.TaskTypeExec)
	taskID, _ := db.InsertPendingTask(typeID, "/work", "/work", "exec", "exec fetch --prune")

	failPendingTask(db, taskID, "exec batch failed with exit code 1")

	rec, _ := db.FindPendingTaskByID(taskID)
	if rec.FailureReason != "exec batch failed with exit code 1" {
		t.Errorf("expected failure reason, got %q", rec.FailureReason)
	}
}

func TestExecPending_FailKeepsTaskPending(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID, _ := db.GetTaskTypeID(constants.TaskTypeExec)
	taskID, _ := db.InsertPendingTask(typeID, "/work", "/work", "exec", "exec pull --rebase")

	failPendingTask(db, taskID, "batch failed")

	pending, _ := db.ListPendingTasks()
	if len(pending) != 1 {
		t.Fatalf("expected 1 pending, got %d", len(pending))
	}
	if pending[0].ID != taskID {
		t.Errorf("expected ID %d, got %d", taskID, pending[0].ID)
	}
}

// --- Exec full lifecycle ---

func TestExecPending_FullLifecycle_FailThenComplete(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	typeID, _ := db.GetTaskTypeID(constants.TaskTypeExec)
	taskID, _ := db.InsertPendingTask(typeID, "/work", "/work", "exec", "exec fetch --all --prune")

	// First attempt fails.
	failPendingTask(db, taskID, "network timeout")
	rec, _ := db.FindPendingTaskByID(taskID)
	if rec.FailureReason != "network timeout" {
		t.Errorf("expected 'network timeout', got %q", rec.FailureReason)
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
	if completed[0].CommandArgs != "exec fetch --all --prune" {
		t.Errorf("expected cmdArgs preserved, got %q", completed[0].CommandArgs)
	}
}

// --- buildCommandArgs for exec ---

func TestBuildCommandArgs_ExecWithFlags(t *testing.T) {
	result := buildCommandArgs([]string{"exec", "--group", "backend", "fetch", "--prune"})
	if result != "exec --group backend fetch --prune" {
		t.Errorf("expected 'exec --group backend fetch --prune', got %q", result)
	}
}

func TestBuildCommandArgs_ExecStopOnFail(t *testing.T) {
	result := buildCommandArgs([]string{"exec", "--stop-on-fail", "status"})
	if result != "exec --stop-on-fail status" {
		t.Errorf("expected 'exec --stop-on-fail status', got %q", result)
	}
}
