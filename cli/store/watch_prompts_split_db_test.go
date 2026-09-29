package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWatchPromptsSplitDB_Lifecycle(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gitmap-wpr-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	origOverride := binaryDataDirOverride
	SetBinaryDataDirForTesting(filepath.Join(tempDir, "data"))
	defer SetBinaryDataDirForTesting(origOverride)

	slug := "test-repo"
	db, err := OpenWatchPromptsSplitDB(slug)
	if err != nil {
		t.Fatalf("OpenWatchPromptsSplitDB failed: %v", err)
	}
	defer db.Close()

	rec1 := WatchPromptRecord{
		RecordId:     "rec-1",
		RepoSlug:     slug,
		ProjectName:  "Test Repo",
		ProjectPath:  tempDir,
		ProjectId:    "proj-1",
		PromptText:   "Initial prompt",
		PromptStatus: "queued",
		WordCount:    2,
		MediaPaths:   []string{"assets/test.png"},
		IsActive:     true,
		UpdatedAt:    "2026-09-29T10:00:00Z",
		CreatedAt:    "2026-09-29T10:00:00Z",
	}

	if saveErr := db.SavePromptRecord(rec1); saveErr != nil {
		t.Fatalf("SavePromptRecord failed: %v", saveErr)
	}

	records, listErr := db.ListPromptRecords(slug)
	if listErr != nil || len(records) != 1 {
		t.Fatalf("expected 1 record, got %d (err: %v)", len(records), listErr)
	}
	if len(records[0].MediaPaths) != 1 || records[0].MediaPaths[0] != "assets/test.png" {
		t.Errorf("expected media path assets/test.png, got %+v", records[0].MediaPaths)
	}

	// Test prune previous records so only recent prompt is preserved
	if pruneErr := db.PrunePreviousPromptRecords(slug); pruneErr != nil {
		t.Fatalf("PrunePreviousPromptRecords failed: %v", pruneErr)
	}

	rec2 := WatchPromptRecord{
		RecordId:     "rec-2",
		RepoSlug:     slug,
		ProjectName:  "Test Repo",
		ProjectPath:  tempDir,
		ProjectId:    "proj-1",
		PromptText:   "Recent prompt",
		PromptStatus: "running",
		WordCount:    2,
		IsActive:     true,
		UpdatedAt:    "2026-09-29T10:05:00Z",
		CreatedAt:    "2026-09-29T10:05:00Z",
	}
	_ = db.SavePromptRecord(rec2)

	afterPrune, _ := db.ListPromptRecords(slug)
	if len(afterPrune) != 1 || afterPrune[0].RecordId != "rec-2" {
		t.Fatalf("expected 1 recent record 'rec-2', got %d", len(afterPrune))
	}

	// Test WatchLog
	if logErr := db.InsertWatchLog(slug, "heartbeat", "Watchdog active", "healthy"); logErr != nil {
		t.Fatalf("InsertWatchLog failed: %v", logErr)
	}
	logs, logsErr := db.ListWatchLogs(slug, 10)
	if logsErr != nil || len(logs) == 0 {
		t.Fatalf("expected logs, got error: %v", logsErr)
	}

	// Test WatchConfig
	if cfgErr := db.SaveWatchConfig("interval", "2m"); cfgErr != nil {
		t.Fatalf("SaveWatchConfig failed: %v", cfgErr)
	}
	val, getErr := db.GetWatchConfig("interval")
	if getErr != nil || val != "2m" {
		t.Errorf("expected config '2m', got %q (err: %v)", val, getErr)
	}
}
