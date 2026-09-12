package cmd

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

func setupTestSSHDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	initTestSSHTables(t, db)

	return db
}

func initTestSSHTables(t *testing.T, db *sql.DB) {
	createHosts := `CREATE TABLE IF NOT EXISTS ssh_hosts (
		id TEXT PRIMARY KEY, alias TEXT, ip TEXT, username TEXT, created_at DATETIME
	);`
	createHist := `CREATE TABLE IF NOT EXISTS ssh_history (
		id TEXT PRIMARY KEY, host_ip TEXT, joined_at DATETIME, user TEXT
	);`
	if _, err := db.Exec(createHosts); err != nil {
		t.Fatalf("failed to create ssh_hosts: %v", err)
	}

	if _, err := db.Exec(createHist); err != nil {
		t.Fatalf("failed to create ssh_history: %v", err)
	}
}

func TestExecuteSSHJoinAtomicity(t *testing.T) {
	db := setupTestSSHDB(t)
	defer db.Close()

	ctx := context.Background()
	hist := store.SSHHistory{
		ID:       "host-1",
		HostIP:   "192.168.1.100",
		JoinedAt: time.Now(),
		User:     "admin",
	}

	if err := runJoinTransaction(ctx, db, "my-server", hist); err != nil {
		t.Fatalf("runJoinTransaction failed: %v", err)
	}

	assertJoinCounts(t, db, "host-1", 1, 1)
}

func assertJoinCounts(t *testing.T, db *sql.DB, id string, wantHosts, wantHist int) {
	var hostCount, histCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM ssh_hosts WHERE id = ?", id).Scan(&hostCount); err != nil {
		t.Fatalf("query ssh_hosts failed: %v", err)
	}

	if err := db.QueryRow("SELECT COUNT(*) FROM ssh_history WHERE id = ?", id).Scan(&histCount); err != nil {
		t.Fatalf("query ssh_history failed: %v", err)
	}

	if hostCount != wantHosts || histCount != wantHist {
		t.Errorf("counts mismatch: hosts=%d (want %d), hist=%d (want %d)", hostCount, wantHosts, histCount, wantHist)
	}
}

func TestExecuteSSHJoinRollback(t *testing.T) {
	db := setupTestSSHDB(t)
	defer db.Close()

	ctx := context.Background()
	seedCollisionHist(t, db, "collision-id")

	hist := store.SSHHistory{ID: "collision-id", HostIP: "1.1.1.1", JoinedAt: time.Now(), User: "root"}
	if err := runJoinTransaction(ctx, db, "server-fail", hist); err == nil {
		t.Fatalf("expected failure on duplicate history id, got nil")
	}

	assertJoinCounts(t, db, "collision-id", 0, 1)
}

func seedCollisionHist(t *testing.T, db *sql.DB, id string) {
	query := "INSERT INTO ssh_history (id, host_ip, joined_at, user) VALUES (?, ?, ?, ?)"
	if _, err := db.Exec(query, id, "1.1.1.1", time.Now(), "root"); err != nil {
		t.Fatalf("seedCollisionHist failed: %v", err)
	}
}

func TestExecuteSSHJoinSkeleton(t *testing.T) {
	_ = executeSSHJoin
}
