package store

import (
	"path/filepath"
	"testing"
	"time"
)

func TestPullSplitDB_LifecycleAndInactivity(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test-gitmap-pull.db")

	db, err := OpenPullSplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("OpenPullSplitDBAt failed: %v", err)
	}
	defer db.Close()

	runRec := &PullRunRecord{
		CommandType:   "pull-all",
		WorkingDir:    tempDir,
		TotalRepos:    2,
		PulledRepos:   2,
		SkippedRepos:  0,
		SuccessCount:  2,
		FailedCount:   0,
		IsEfficient:   false,
		DurationMs:    150,
		GitMapVersion: "v6.307.0",
		Notes:         "test run",
	}

	runID, err := db.InsertPullRun(runRec)
	if err != nil {
		t.Fatalf("InsertPullRun failed: %v", err)
	}
	if runID <= 0 {
		t.Fatalf("expected positive runID, got %d", runID)
	}

	repoA := filepath.Join(tempDir, "repo-alpha")
	repoB := filepath.Join(tempDir, "repo-beta")

	records := []PullRepoRunRecord{
		{
			RepoPath:      repoA,
			RepoName:      "repo-alpha",
			PullStatus:    "up-to-date",
			FilesChanged:  0,
			LastCommitSha: "abc1234",
			IsActive:      true,
			HasChanges:    false,
			DurationMs:    45,
		},
		{
			RepoPath:      repoB,
			RepoName:      "repo-beta",
			PullStatus:    "success",
			FilesChanged:  3,
			LastCommitSha: "def5678",
			IsActive:      true,
			HasChanges:    true,
			DurationMs:    95,
		},
	}

	if err := db.InsertPullRepoRuns(runID, records); err != nil {
		t.Fatalf("InsertPullRepoRuns failed: %v", err)
	}

	historyA, err := db.GetRecentRepoPullHistory(repoA, 10)
	if err != nil {
		t.Fatalf("GetRecentRepoPullHistory failed: %v", err)
	}
	if len(historyA) != 1 {
		t.Fatalf("expected 1 history record for repoA, got %d", len(historyA))
	}

	// 1 run is insufficient for minRuns=20 -> should be active
	statusA, err := db.EvaluateRepoInactivity(repoA, 20, 24)
	if err != nil {
		t.Fatalf("EvaluateRepoInactivity failed: %v", err)
	}
	if statusA.IsInactive {
		t.Fatalf("expected repoA to be active due to insufficient runs, got inactive")
	}

	// Insert 20 zero-change records for repoA
	var batch []PullRepoRunRecord
	for i := 0; i < 20; i++ {
		batch = append(batch, PullRepoRunRecord{
			RepoPath:      repoA,
			RepoName:      "repo-alpha",
			PullStatus:    "up-to-date",
			FilesChanged:  0,
			LastCommitSha: "abc1234",
			IsActive:      true,
			HasChanges:    false,
			DurationMs:    10,
		})
	}
	if err := db.InsertPullRepoRuns(runID, batch); err != nil {
		t.Fatalf("batch insert failed: %v", err)
	}

	// Now repoA has 21 records with 0 changes -> should be inactive
	statusA2, err := db.EvaluateRepoInactivity(repoA, 20, 24)
	if err != nil {
		t.Fatalf("EvaluateRepoInactivity failed: %v", err)
	}
	if !statusA2.IsInactive {
		t.Fatalf("expected repoA to be inactive after 20 zero-change runs, got active (reason: %s)", statusA2.Reason)
	}

	// repoB had changes -> should be active
	statusB, err := db.EvaluateRepoInactivity(repoB, 1, 24)
	if err != nil {
		t.Fatalf("EvaluateRepoInactivity repoB failed: %v", err)
	}
	if statusB.IsInactive {
		t.Fatalf("expected repoB to be active because it has changes")
	}
}

func TestPullSplitDB_ParseCreatedAtTime(t *testing.T) {
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	parsed, isParsed := parseCreatedAtTime(now)
	if !isParsed {
		t.Fatalf("failed to parse standard format time: %s", now)
	}
	if parsed.IsZero() {
		t.Fatalf("expected non-zero time")
	}

	rfcNow := time.Now().UTC().Format(time.RFC3339)
	parsedRFC, isParsedRFC := parseCreatedAtTime(rfcNow)
	if !isParsedRFC {
		t.Fatalf("failed to parse RFC3339 format time: %s", rfcNow)
	}
	if parsedRFC.IsZero() {
		t.Fatalf("expected non-zero time")
	}

	_, isInvalid := parseCreatedAtTime("invalid-timestamp")
	if isInvalid {
		t.Fatalf("expected invalid timestamp to fail parsing")
	}
}
