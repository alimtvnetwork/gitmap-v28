package store

import (
	"path/filepath"
	"testing"
	"time"
)

func TestPurgeHistory_MigrationFromLegacySchema(t *testing.T) {
	tempDir := t.TempDir()
	db, err := OpenAt(filepath.Join(tempDir, "legacy_purge.db"))
	if err != nil {
		t.Fatalf("OpenAt failed: %v", err)
	}
	defer db.Close()

	legacySchema := `CREATE TABLE PurgeHistoryLog (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		RepoPath TEXT NOT NULL,
		Pattern TEXT NOT NULL,
		BackupBranch TEXT NOT NULL,
		TempDir TEXT NOT NULL,
		Files TEXT NOT NULL,
		Timestamp INTEGER NOT NULL,
		Restored BOOLEAN NOT NULL DEFAULT 0
	)`

	if _, err := db.conn.Exec(legacySchema); err != nil {
		t.Fatalf("failed to create legacy PurgeHistoryLog table: %v", err)
	}

	assertLegacyColumns(t, db)

	insertSQL := `INSERT INTO PurgeHistoryLog (RepoPath, Pattern, BackupBranch, TempDir, Files, Timestamp, Restored)
		VALUES ('/path/to/legacy/repo', '*.orig', 'backup-purge-100', '/tmp/p100', '["file.orig"]', 1700000000, 0)`
	if _, err := db.conn.Exec(insertSQL); err != nil {
		t.Fatalf("failed to insert legacy row: %v", err)
	}

	if err := db.EnsurePurgeHistoryTable(); err != nil {
		t.Fatalf("EnsurePurgeHistoryTable failed on legacy DB: %v", err)
	}

	assertMigratedColumns(t, db)

	log, err := db.GetLastPurgeHistoryLog("/path/to/legacy/repo")
	if err != nil || log == nil {
		t.Fatalf("GetLastPurgeHistoryLog failed or returned nil: %v", err)
	}
	if log.PurgeHistoryLogId != 1 {
		t.Errorf("expected PurgeHistoryLogId 1, got %d", log.PurgeHistoryLogId)
	}
	if log.IsRestored {
		t.Errorf("expected unrestored state")
	}

	if err := db.EnsurePurgeHistoryTable(); err != nil {
		t.Fatalf("second EnsurePurgeHistoryTable call failed: %v", err)
	}
}

func assertLegacyColumns(t *testing.T, db *DB) {
	if !db.columnExists("PurgeHistoryLog", "ID") || !db.columnExists("PurgeHistoryLog", "Restored") {
		t.Fatalf("expected legacy columns ID and Restored to exist")
	}
}

func assertMigratedColumns(t *testing.T, db *DB) {
	cols := []string{"PurgeHistoryLogId", "IsRestored", "Notes", "Comments"}
	for _, c := range cols {
		if !db.columnExists("PurgeHistoryLog", c) {
			t.Errorf("expected column %s after migration", c)
		}
	}
}

func TestPurgeHistory_RestoredState(t *testing.T) {
	tempDir := t.TempDir()
	db, err := OpenAt(filepath.Join(tempDir, "purge_compat.db"))
	if err != nil {
		t.Fatalf("OpenAt failed: %v", err)
	}
	defer db.Close()

	entry := &PurgeHistoryLog{
		RepoPath:     "/workspace/alias-repo",
		Pattern:      "*.bak",
		BackupBranch: "backup-bak",
		TempDir:      "/tmp/purge-bak",
		Files:        `[]`,
		Timestamp:    time.Now().Unix(),
		IsRestored:   true,
	}

	if err := db.InsertPurgeHistoryLog(entry); err != nil {
		t.Fatalf("InsertPurgeHistoryLog failed: %v", err)
	}

	byId, err := db.GetPurgeHistoryLogById(entry.PurgeHistoryLogId)
	if err != nil {
		t.Fatalf("GetPurgeHistoryLogById failed: %v", err)
	}
	if !byId.IsRestored {
		t.Errorf("expected IsRestored=true")
	}
}
