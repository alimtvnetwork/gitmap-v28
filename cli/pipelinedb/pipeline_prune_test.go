package pipelinedb

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func createTestPruneSplitDb(t *testing.T) *PipelineSplitDb {
	t.Helper()
	db, err := OpenPipelineSplitDB("test-prune-owner/test-prune-repo-" + t.Name())
	if err != nil {
		t.Fatalf("failed to open test pipeline db: %v", err)
	}

	return db
}

func TestPipelinePrune_DiskSizeBytes(t *testing.T) {
	db := createTestPruneSplitDb(t)
	defer db.Close()

	size := db.DiskSizeBytes()
	if size <= 0 {
		t.Errorf("expected positive disk size bytes, got %d", size)
	}
}

func TestPipelinePrune_HasCompletedRunAndSha(t *testing.T) {
	db := createTestPruneSplitDb(t)
	defer db.Close()

	rec := initSamplePruneRun(1001, "abc1234567890")
	if err := db.RecordRun(rec); err != nil {
		t.Fatalf("failed to record run: %v", err)
	}

	assertCompletedRunMatch(t, db, 1001, "abc1234567890")
}

func initSamplePruneRun(id uint64, sha string) PipelineRunRecord {
	return PipelineRunRecord{
		RunId:        id,
		RepoSlug:     "test/prune-repo",
		WorkflowName: "CI",
		Status:       "completed",
		Conclusion:   "success",
		Branch:       "main",
		Sha:          sha,
		CreatedAt:    "2026-09-20T00:00:00Z",
		UpdatedAt:    "2026-09-20T00:01:00Z",
	}
}

func assertCompletedRunMatch(t *testing.T, db *PipelineSplitDb, id uint64, sha string) {
	if !db.HasCompletedRun(id, sha) {
		t.Errorf("expected HasCompletedRun to be true for runId %d and sha %s", id, sha)
	}
	foundRun, err := db.QueryRunBySha(sha[:7])
	if err != nil || foundRun == nil || foundRun.RunId != id {
		t.Errorf("expected QueryRunBySha to find run %d, got: %v, err: %v", id, foundRun, err)
	}
	foundById, err := db.QueryRunByRunId(id)
	if err != nil || foundById == nil || foundById.Sha != sha {
		t.Errorf("expected QueryRunByRunId to find run %d, got: %v, err: %v", id, foundById, err)
	}
}

func TestPipelinePrune_PruneIfExceedsSize(t *testing.T) {
	db := createTestPruneSplitDb(t)
	defer db.Close()

	populateTestPruneRuns(db, 15)

	didPrune, count, err := db.PruneIfExceedsSize(100)
	if err != nil {
		t.Fatalf("prune failed: %v", err)
	}
	if !didPrune || count <= 0 {
		t.Errorf("expected didPrune true and count > 0, got %v, %d", didPrune, count)
	}
}

func populateTestPruneRuns(db *PipelineSplitDb, count uint64) {
	for i := uint64(1); i <= count; i++ {
		rec := PipelineRunRecord{
			RunId: 2000 + i, RepoSlug: "test/prune-repo",
			WorkflowName: "CI", Status: "completed", Conclusion: "failure",
			Branch: "main", Sha: fmt.Sprintf("sha%04d", i),
			CreatedAt: fmt.Sprintf("2026-09-20T00:%02d:00Z", i),
			UpdatedAt: fmt.Sprintf("2026-09-20T00:%02d:30Z", i),
		}
		_ = db.RecordRun(rec)
		detail := PipelineErrorRecord{
			RunId: rec.RunId, RepoSlug: rec.RepoSlug, WorkflowName: rec.WorkflowName,
			StepName: "Test Step", ErrorText: "test failure", RawLogs: "simulated long raw log",
		}
		_ = db.RecordDetailErrorLog(detail)
	}
}

func TestPipelinePrune_RemovesLegacyPipelineDb(t *testing.T) {
	db := createTestPruneSplitDb(t)
	defer db.Close()

	legacyDb := filepath.Join(filepath.Dir(db.Path), "pipeline.db")
	_ = os.WriteFile(legacyDb, []byte("fake-legacy-data"), 0644)

	db.cleanLegacyDbIfPresent()

	if isFileExisting(legacyDb) {
		t.Errorf("expected legacy pipeline.db to be removed")
	}
}

func TestPipelinePrune_TotalPipelineDiskBytes(t *testing.T) {
	db := createTestPruneSplitDb(t)
	defer db.Close()

	totalBytes := db.TotalPipelineDiskBytes()
	if totalBytes <= 0 {
		t.Errorf("expected positive total pipeline disk bytes, got %d", totalBytes)
	}
}
