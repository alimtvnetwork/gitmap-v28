package pipelinedb

import (
	"testing"
)

func TestPipelineSplitDBLifecycle(t *testing.T) {
	db, err := OpenPipelineSplitDB("test-owner/test-repo")
	if err != nil {
		t.Fatalf("failed to open pipeline split db: %v", err)
	}
	defer db.Close()

	// 1. Record Run
	run := PipelineRunRecord{
		RunId:        12345,
		RepoSlug:     "test-owner/test-repo",
		WorkflowName: "CI",
		Status:       "completed",
		Conclusion:   "failure",
		Branch:       "main",
		Sha:          "abc1234",
		EtaSeconds:   60,
		RunUrl:       "https://github.com/test-owner/test-repo/actions/runs/12345",
		CreatedAt:    "2026-09-03T10:00:00Z",
		UpdatedAt:    "2026-09-03T10:02:00Z",
	}
	if err := db.RecordRun(run); err != nil {
		t.Fatalf("failed to record run: %v", err)
	}

	// 2. Record Error Log
	errLog := PipelineErrorRecord{
		RunId:        12345,
		RepoSlug:     "test-owner/test-repo",
		WorkflowName: "CI",
		StepName:     "Test Step",
		ErrorText:    "##[error]Process exited with code 1",
	}
	if err := db.RecordErrorLog(errLog); err != nil {
		t.Fatalf("failed to record error log: %v", err)
	}

	// 3. Query Stats
	stats, err := db.GetStats()
	if err != nil {
		t.Fatalf("failed to get stats: %v", err)
	}
	if stats.TotalRuns != 1 {
		t.Errorf("expected 1 run, got %d", stats.TotalRuns)
	}
	if stats.ErrorLogCount != 1 {
		t.Errorf("expected 1 error log, got %d", stats.ErrorLogCount)
	}

	// 4. Query Error Logs
	logs, err := db.QueryRecentErrorLogs(5)
	if err != nil || len(logs) != 1 {
		t.Fatalf("expected 1 error log, got %d (err: %v)", len(logs), err)
	}
	if logs[0].StepName != "Test Step" {
		t.Errorf("expected step name 'Test Step', got %s", logs[0].StepName)
	}

	// 5. Optimize
	if _, err := db.Optimize(); err != nil {
		t.Errorf("failed to optimize: %v", err)
	}

	// 6. Clear
	if err := db.Clear(); err != nil {
		t.Fatalf("failed to clear: %v", err)
	}
	statsAfterClear, _ := db.GetStats()
	if statsAfterClear.TotalRuns != 0 || statsAfterClear.ErrorLogCount != 0 {
		t.Errorf("expected 0 runs and 0 error logs after clear, got %d, %d",
			statsAfterClear.TotalRuns, statsAfterClear.ErrorLogCount)
	}

	// 7. Reset
	if err := db.Reset(); err != nil {
		t.Fatalf("failed to reset: %v", err)
	}
}
