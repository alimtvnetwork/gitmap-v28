package store

import (
	"context"
	"database/sql"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}

	err = RegisterSSHHistoryMigration(db, 1, false)
	if err != nil {
		t.Fatalf("failed to run migration: %v", err)
	}
	return db
}

func TestLogSSHJoin(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()
	now := time.Now().Truncate(time.Second).UTC() // Use UTC to avoid timezone issues with SQLite

	h := SSHHistory{
		ID:       "log-test-1",
		HostIP:   "10.0.0.1",
		JoinedAt: now,
		User:     "testuser",
	}

	err := LogSSHJoin(ctx, h, db)
	if err != nil {
		t.Fatalf("LogSSHJoin failed: %v", err)
	}

	// Verify insertion
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM ssh_history WHERE id = 'log-test-1'").Scan(&count)
	if err != nil {
		t.Fatalf("failed to query count: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 record, got %d", count)
	}
}

func TestListSSHHistory(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()
	now := time.Now().Truncate(time.Second).UTC()

	// Insert some records
	records := []SSHHistory{
		{ID: "list-test-1", HostIP: "10.0.0.1", JoinedAt: now.Add(-1 * time.Hour), User: "user1"},
		{ID: "list-test-2", HostIP: "10.0.0.2", JoinedAt: now.Add(-2 * time.Hour), User: "user2"},
		{ID: "list-test-3", HostIP: "10.0.0.3", JoinedAt: now.Add(-3 * time.Hour), User: "user3"},
	}

	for _, r := range records {
		err := LogSSHJoin(ctx, r, db)
		if err != nil {
			t.Fatalf("LogSSHJoin failed: %v", err)
		}
	}

	// Test List
	results, err := ListSSHHistory(ctx, 2, 0, db)
	if err != nil {
		t.Fatalf("ListSSHHistory failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 records, got %d", len(results))
	}

	// Should be ordered by joined_at DESC
	if results[0].ID != "list-test-1" {
		t.Errorf("expected list-test-1 to be first, got %s", results[0].ID)
	}
	if results[1].ID != "list-test-2" {
		t.Errorf("expected list-test-2 to be second, got %s", results[1].ID)
	}
	
	// Test Empty
	emptyResults, err := ListSSHHistory(ctx, 10, 10, db)
	if err != nil {
		t.Fatalf("ListSSHHistory failed on empty offset: %v", err)
	}
	if emptyResults == nil {
		t.Error("expected non-nil empty slice, got nil")
	}
	if len(emptyResults) != 0 {
		t.Errorf("expected empty results, got %d", len(emptyResults))
	}
}
