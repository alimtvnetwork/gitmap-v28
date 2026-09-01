// Package cmd — schedule_split_test.go: unit tests for per-schedule split databases,
// run history logging, enable/disable states, and reset operations.
package cmd

import (
	"os"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func TestScheduleSplitDBAndLogs(t *testing.T) {
	// 1. Add a scheduled task
	taskName := "split-test-task"
	addArgs := []string{taskName, "echo Testing Split DB Execution", "--every=2h", "--delay=10ms"}
	if err := runScheduleAdd(addArgs); err != nil {
		t.Fatalf("runScheduleAdd failed: %v", err)
	}

	slug := store.ScheduleSlug(taskName)
	splitPath := store.ScheduleDBPath(slug)
	if _, err := os.Stat(splitPath); err != nil {
		t.Fatalf("expected split DB file %s to exist: %v", splitPath, err)
	}

	// 2. Run the task to populate split DB logs
	if err := runScheduleRun([]string{taskName}); err != nil {
		t.Fatalf("runScheduleRun failed: %v", err)
	}

	// 3. Inspect logs directly from split DB
	splitDB, err := store.OpenScheduleSplitDB(slug)
	if err != nil {
		t.Fatalf("open split DB failed: %v", err)
	}
	runs, err := splitDB.GetRuns(10)
	splitDB.Close()
	if err != nil {
		t.Fatalf("get runs failed: %v", err)
	}
	if len(runs) != 1 || !runs[0].IsSuccess {
		t.Fatalf("unexpected run logs: %+v", runs)
	}

	// 4. Test logs CLI (table, JSON, YAML)
	if err := runScheduleLogs([]string{taskName}); err != nil {
		t.Errorf("runScheduleLogs table failed: %v", err)
	}
	if err := runScheduleLogs([]string{taskName, "--json"}); err != nil {
		t.Errorf("runScheduleLogs json failed: %v", err)
	}
	if err := runScheduleLogs([]string{taskName, "--yaml"}); err != nil {
		t.Errorf("runScheduleLogs yaml failed: %v", err)
	}

	// 5. Test disable and enable
	if err := runScheduleSetEnabled([]string{taskName}, false); err != nil {
		t.Fatalf("disable failed: %v", err)
	}
	db, _ := openSchedulerDB()
	task, _ := db.GetSchedule(taskName)
	db.Close()
	if task.IsEnabled {
		t.Errorf("expected task to be disabled")
	}

	if err := runScheduleSetEnabled([]string{taskName}, true); err != nil {
		t.Fatalf("enable failed: %v", err)
	}
	db, _ = openSchedulerDB()
	task, _ = db.GetSchedule(taskName)
	db.Close()
	if !task.IsEnabled {
		t.Errorf("expected task to be enabled")
	}

	// 6. Test reset
	if err := runScheduleReset([]string{taskName}); err != nil {
		t.Fatalf("reset failed: %v", err)
	}
	splitDB, _ = store.OpenScheduleSplitDB(slug)
	runsAfterReset, _ := splitDB.GetRuns(10)
	splitDB.Close()
	if len(runsAfterReset) != 0 {
		t.Errorf("expected 0 runs after reset, got %d", len(runsAfterReset))
	}

	// 7. Cleanup
	_ = runScheduleDelete([]string{taskName})
}
