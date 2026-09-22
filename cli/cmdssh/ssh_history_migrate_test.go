package cmdssh

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

const sqlInitTaskHistoryTable = `CREATE TABLE IF NOT EXISTS TaskHistory (
    TaskHistoryId  INTEGER PRIMARY KEY AUTOINCREMENT,
    TaskId         TEXT NOT NULL UNIQUE,
    Section        TEXT NOT NULL DEFAULT '',
    Action         TEXT NOT NULL,
    Target         TEXT NOT NULL,
    ForwardPayload TEXT NOT NULL DEFAULT '',
    InversePayload TEXT NOT NULL DEFAULT '',
    Status         TEXT NOT NULL DEFAULT 'completed',
    RestoredAt     TEXT NOT NULL DEFAULT '',
    ExecutedAt     TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CreatedAt      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);`

func TestMigrateLegacySSHHistory_WithoutForwardPayload(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec(sqlInitTaskHistoryTable); err != nil {
		t.Fatalf("failed to create TaskHistory: %v", err)
	}

	// Legacy table with payload_json only (no forward_payload)
	sqlLegacy := `CREATE TABLE ssh_task_history (
		task_id TEXT PRIMARY KEY,
		action TEXT NOT NULL,
		target TEXT NOT NULL,
		payload_json TEXT NOT NULL,
		created_at TEXT NOT NULL,
		restored_at TEXT DEFAULT ''
	);`
	if _, err := db.Exec(sqlLegacy); err != nil {
		t.Fatalf("failed to create legacy table: %v", err)
	}

	insertLegacy := `INSERT INTO ssh_task_history VALUES ('t1', 'rm', 'w2', '[{"alias":"w2"}]', '2026-09-21 17:00:00', '');`
	if _, err := db.Exec(insertLegacy); err != nil {
		t.Fatalf("failed to insert legacy row: %v", err)
	}

	migrateLegacySSHHistoryTables(db)

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM TaskHistory WHERE TaskId = 't1'").Scan(&count); err != nil {
		t.Fatalf("query migrated row failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 migrated row, got %d", count)
	}

	if checkTableExists(db, "ssh_task_history") {
		t.Errorf("expected ssh_task_history table to be dropped post-migration")
	}

	// Idempotent re-run
	migrateLegacySSHHistoryTables(db)
}

func TestMigrateLegacySSHHistory_WithForwardPayload(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec(sqlInitTaskHistoryTable); err != nil {
		t.Fatalf("failed to create TaskHistory: %v", err)
	}

	sqlModern := `CREATE TABLE ssh_task_history (
		task_id TEXT PRIMARY KEY,
		action TEXT NOT NULL,
		target TEXT NOT NULL,
		forward_payload TEXT NOT NULL,
		inverse_payload TEXT NOT NULL,
		created_at TEXT NOT NULL,
		restored_at TEXT DEFAULT ''
	);`
	if _, err := db.Exec(sqlModern); err != nil {
		t.Fatalf("failed to create modern table: %v", err)
	}

	insertModern := `INSERT INTO ssh_task_history VALUES ('t2', 'rm', 'w3', 'w3', '[{"alias":"w3"}]', '2026-09-22 10:00:00', '');`
	if _, err := db.Exec(insertModern); err != nil {
		t.Fatalf("failed to insert modern row: %v", err)
	}

	migrateLegacySSHHistoryTables(db)

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM TaskHistory WHERE TaskId = 't2'").Scan(&count); err != nil {
		t.Fatalf("query migrated row failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 migrated row, got %d", count)
	}
}
