package store

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestSQLCreateSSHHosts(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}
	defer db.Close()

	err = RegisterSSHHostMigration(db, 1, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify table exists
	row := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='ssh_hosts'")
	var name string
	err = row.Scan(&name)
	if err != nil {
		t.Fatalf("expected table ssh_hosts to exist: %v", err)
	}
}
