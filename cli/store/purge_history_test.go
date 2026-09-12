package store

import (
	"path/filepath"
	"testing"
	"time"
)

func TestPurgeHistory_TableCreationAndMigrate(t *testing.T) {
	tempDir := t.TempDir()
	db, err := OpenAt(filepath.Join(tempDir, "test_purge.db"))
	if err != nil {
		t.Fatalf("OpenAt failed: %v", err)
	}

	defer db.Close()

	if err := db.Migrate(); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	if !db.tableExists("PurgeHistoryLog") {
		t.Fatalf("expected table PurgeHistoryLog to exist after Migrate")
	}

	expectedColumns := []string{
		"PurgeHistoryLogId", "RepoPath", "Pattern", "BackupBranch",
		"TempDir", "Files", "Timestamp", "IsRestored", "Notes", "Comments",
	}

	for _, col := range expectedColumns {
		if !db.columnExists("PurgeHistoryLog", col) {
			t.Errorf("expected column %s to exist in PurgeHistoryLog", col)
		}
	}
}

func TestPurgeHistory_InsertRetrievalAndRestoration(t *testing.T) {
	tempDir := t.TempDir()
	db, err := OpenAt(filepath.Join(tempDir, "purge_crud.db"))
	if err != nil {
		t.Fatalf("OpenAt failed: %v", err)
	}

	defer db.Close()

	now := time.Now().Unix()
	logEntry := &PurgeHistoryLog{
		RepoPath:     "/workspace/sample-repo",
		Pattern:      "*.secrets",
		BackupBranch: "backup-purge-secrets",
		TempDir:      "/tmp/purge-temp",
		Files:        `["keys.secrets", "cert.secrets"]`,
		Timestamp:    now,
		IsRestored:   false,
		Notes:        "Automated purge operation",
		Comments:     "Tested in unit test",
	}

	if err := db.InsertPurgeHistoryLog(logEntry); err != nil {
		t.Fatalf("InsertPurgeHistoryLog failed: %v", err)
	}

	if logEntry.PurgeHistoryLogId <= 0 {
		t.Fatalf("expected positive PurgeHistoryLogId, got %d", logEntry.PurgeHistoryLogId)
	}

	retrieved, err := db.GetLastPurgeHistoryLog("/workspace/sample-repo")
	if err != nil || retrieved == nil {
		t.Fatalf("GetLastPurgeHistoryLog failed or returned nil: %v", err)
	}

	if retrieved.PurgeHistoryLogId != logEntry.PurgeHistoryLogId {
		t.Errorf("expected PurgeHistoryLogId %d, got %d", logEntry.PurgeHistoryLogId, retrieved.PurgeHistoryLogId)
	}

	if retrieved.IsRestored {
		t.Errorf("expected unrestored state")
	}

	if err := db.MarkPurgeHistoryRestored(retrieved.PurgeHistoryLogId); err != nil {
		t.Fatalf("MarkPurgeHistoryRestored failed: %v", err)
	}

	unrestored, err := db.GetLastPurgeHistoryLog("/workspace/sample-repo")
	if err != nil || unrestored != nil {
		t.Fatalf("expected nil for unrestored search after restore, got %+v (err: %v)", unrestored, err)
	}

	restoredEntry, err := db.GetPurgeHistoryLogById(retrieved.PurgeHistoryLogId)
	if err != nil || restoredEntry == nil {
		t.Fatalf("GetPurgeHistoryLogById failed or returned nil: %v", err)
	}

	if !restoredEntry.IsRestored {
		t.Errorf("expected IsRestored=true")
	}
}
