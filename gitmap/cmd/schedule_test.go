// Package cmd — schedule_test.go: unit tests for scheduler CLI and REST API.
package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScheduleLifecycle(t *testing.T) {
	// 1. Test adding scheduled task with command and interval
	addArgs := []string{"daily-build", "echo Build Finished", "--every=1d", "--delay=10ms"}
	if err := runScheduleAdd(addArgs); err != nil {
		t.Fatalf("runScheduleAdd failed: %v", err)
	}

	// 2. Test retrieving from store
	db, err := openSchedulerDB()
	if err != nil {
		t.Fatalf("open db failed: %v", err)
	}
	defer db.Close()

	task, err := db.GetSchedule("daily-build")
	if err != nil {
		t.Fatalf("get schedule failed: %v", err)
	}
	if task.Name != "daily-build" || task.IntervalVal != "1d" || task.CommandLine != "echo Build Finished" {
		t.Errorf("unexpected task data: %+v", task)
	}

	// 3. Test running the task
	if err := runScheduleRun([]string{"daily-build"}); err != nil {
		t.Fatalf("runScheduleRun failed: %v", err)
	}

	// 4. Test testing the task with 2 iterations
	if err := runScheduleTest([]string{"daily-build", "--times=2"}); err != nil {
		t.Fatalf("runScheduleTest failed: %v", err)
	}

	// Verify run count was updated
	taskAfterRuns, _ := db.GetSchedule("daily-build")
	if taskAfterRuns.RunCount < 3 {
		t.Errorf("expected at least 3 runs, got %d", taskAfterRuns.RunCount)
	}

	// 5. Test deleting the task
	if err := runScheduleDelete([]string{"daily-build"}); err != nil {
		t.Fatalf("runScheduleDelete failed: %v", err)
	}
	if _, err := db.GetSchedule("daily-build"); err == nil {
		t.Errorf("expected schedule to be deleted")
	}
}

func TestScheduleTimeUnitsParsing(t *testing.T) {
	opts := parseScheduleAddOpts([]string{"test-hours", "--hour=2", "--delay=5s"})
	if opts.Interval != "2h" {
		t.Errorf("expected 2h, got %s", opts.Interval)
	}

	optsMin := parseScheduleAddOpts([]string{"test-min", "--minute=30"})
	if optsMin.Interval != "30m" {
		t.Errorf("expected 30m, got %s", optsMin.Interval)
	}

	optsSec := parseScheduleAddOpts([]string{"test-sec", "--second=15"})
	if optsSec.Interval != "15s" {
		t.Errorf("expected 15s, got %s", optsSec.Interval)
	}
}

func TestTerminalCommandExecAPI(t *testing.T) {
	mux := http.NewServeMux()
	mountTerminalHandlers(mux)

	reqBody := []byte(`{"command": "echo API_TEST_SUCCESS"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/command/exec", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("API_TEST_SUCCESS")) {
		t.Errorf("expected response to contain output, got %s", rec.Body.String())
	}
}
