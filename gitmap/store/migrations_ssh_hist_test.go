package store

import (
	"database/sql"
	"testing"
)

func TestSQLCreateSSHHistory(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}
	defer db.Close()

	err = RegisterSSHHistoryMigration(db, 1, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	rows, err := db.Query("SELECT id, host_ip, joined_at, user FROM ssh_history LIMIT 1")
	if err != nil {
		t.Errorf("expected table to exist, got query error: %v", err)
	} else {
		rows.Close()
	}

	// Idempotent test
	err = RegisterSSHHistoryMigration(db, 1, false)
	if err != nil {
		t.Fatalf("expected no error on second run, got %v", err)
	}
}
