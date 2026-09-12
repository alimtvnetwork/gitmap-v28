package store

import (
	"path/filepath"
	"testing"
)

func TestPipelineStoreLifecycle(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := OpenAt(dbPath)
	if err != nil {
		t.Fatalf("OpenAt failed: %v", err)
	}

	defer db.Close()

	err = db.InitPipelineTable()
	if err != nil {
		t.Fatalf("InitPipelineTable failed: %v", err)
	}

	run := PipelineRun{
		RunID:        12345,
		Repo:         "alimtvnetwork/gitmap-v28",
		WorkflowName: "CI/CD Pipeline",
		Status:       "in_progress",
		Conclusion:   "",
		Branch:       "main",
		Sha:          "abc1234",
		EtaSeconds:   45,
		ErrorLog:     "",
		URL:          "https://github.com/alimtvnetwork/gitmap-v28/actions/runs/12345",
	}

	err = db.InsertOrUpdatePipelineRun(run)
	if err != nil {
		t.Fatalf("InsertOrUpdatePipelineRun failed: %v", err)
	}

	latest, err := db.GetLatestPipelineRun("alimtvnetwork/gitmap-v28")
	if err != nil {
		t.Fatalf("GetLatestPipelineRun failed: %v", err)
	}

	if latest.RunID != 12345 || latest.Status != "in_progress" {
		t.Fatalf("unexpected latest run: %+v", latest)
	}

	// Update to failure with error log
	run.Status = "completed"
	run.Conclusion = "failure"
	run.ErrorLog = "Compile error at main.go:42"
	run.EtaSeconds = 0

	err = db.InsertOrUpdatePipelineRun(run)
	if err != nil {
		t.Fatalf("InsertOrUpdatePipelineRun update failed: %v", err)
	}

	errRun, err := db.GetLatestPipelineError("alimtvnetwork/gitmap-v28")
	if err != nil {
		t.Fatalf("GetLatestPipelineError failed: %v", err)
	}

	if errRun.ErrorLog != "Compile error at main.go:42" {
		t.Fatalf("unexpected error log: %s", errRun.ErrorLog)
	}

	runs, err := db.ListRecentPipelineRuns("alimtvnetwork/gitmap-v28", 5)
	if err != nil {
		t.Fatalf("ListRecentPipelineRuns failed: %v", err)
	}

	if len(runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(runs))
	}
}
