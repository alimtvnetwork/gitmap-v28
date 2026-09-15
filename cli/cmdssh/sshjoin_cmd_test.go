package cmdssh

import (
	"context"
	"database/sql"
	"path/filepath"
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
	if err := store.EnsureSSHTables(db); err != nil {
		t.Fatalf("failed to ensure ssh tables: %v", err)
	}
}

func createTestHistory(id, ip, user string) store.SSHHistory {
	return store.SSHHistory{
		ID:       id,
		HostIP:   ip,
		JoinedAt: time.Now(),
		User:     user,
	}
}

func TestExecuteSSHJoinAtomicity(t *testing.T) {
	db := setupTestSSHDB(t)
	defer db.Close()

	ctx := context.Background()
	hist := createTestHistory("host-1", "192.168.1.100", "admin")

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

func setupTestStoreDB(t *testing.T) (string, *store.DB) {
	dbPath := filepath.Join(t.TempDir(), "test_ssh.db")
	dbConn, err := store.OpenAt(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	if err := dbConn.Migrate(); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	_ = store.EnsureSSHHostsTable(dbConn.SQL())
	_ = store.EnsureSSHHistoryTable(dbConn.SQL())
	return dbPath, dbConn
}

func withMockSSHDB(t *testing.T, fn func(db *store.DB)) {
	dbPath, testDB := setupTestStoreDB(t)
	defer testDB.Close()

	prevOpener := openSSHDBFunc
	openSSHDBFunc = func() (*store.DB, error) { return store.OpenAt(dbPath) }
	defer func() { openSSHDBFunc = prevOpener }()

	fn(testDB)
}

func TestRunSSHJoinCLI_Validation(t *testing.T) {
	if err := runSSHJoinCLI([]string{}); err == nil {
		t.Fatal("expected error for empty args")
	}

	if err := runSSHJoinCLI([]string{"rm"}); err == nil {
		t.Fatal("expected error for rm without target")
	}

	if err := runSSHJoinCLI([]string{"add-auth"}); err == nil {
		t.Fatal("expected error for add-auth without target")
	}
}

func assertEnrolledHost(t *testing.T, db *sql.DB, alias, ip string) {
	ctx := context.Background()
	host, err := store.GetHostByAlias(ctx, alias, db)
	if err != nil {
		t.Fatalf("expected host %q to exist: %v", alias, err)
	}

	if host.IP != ip {
		t.Fatalf("expected IP %s, got %s", ip, host.IP)
	}
}

func TestRunSSHJoinCLI_EnrollmentAndRecall(t *testing.T) {
	withMockSSHDB(t, func(db *store.DB) {
		args := []string{"10.0.10.1", "node-a"}
		if err := runSSHJoinCLI(args); err != nil {
			t.Fatalf("runSSHJoinCLI failed: %v", err)
		}
		assertEnrolledHost(t, db.SQL(), "node-a", "10.0.10.1")
	})
}

func assertHostDeleted(t *testing.T, db *sql.DB, alias string) {
	ctx := context.Background()
	if _, err := store.GetHostByAlias(ctx, alias, db); err == nil {
		t.Fatalf("expected host %q to be deleted", alias)
	}
}

func TestRunSSHJoinCLI_Removal(t *testing.T) {
	withMockSSHDB(t, func(db *store.DB) {
		_ = runSSHJoinCLI([]string{"10.0.10.2", "rm-box"})
		if err := runSSHJoinCLI([]string{"rm", "rm-box"}); err != nil {
			t.Fatalf("rm failed: %v", err)
		}
		assertHostDeleted(t, db.SQL(), "rm-box")
	})
}

func TestRunSSHJoinCLI_Subcommands(t *testing.T) {
	withMockSSHDB(t, func(db *store.DB) {
		if err := runSSHJoinCLI([]string{"ls"}); err != nil {
			t.Fatalf("ls failed: %v", err)
		}
		if err := runSSHJoinCLI([]string{"history"}); err != nil {
			t.Fatalf("history failed: %v", err)
		}
	})
}

func TestRunSSHJoinCLI_AddSubcommand(t *testing.T) {
	withMockSSHDB(t, func(db *store.DB) {
		args := []string{"add", "192.168.1.14", "box"}
		if err := runSSHJoinCLI(args); err != nil {
			t.Fatalf("runSSHJoinCLI 'add' failed: %v", err)
		}
		assertEnrolledHost(t, db.SQL(), "box", "192.168.1.14")
	})
}

func TestRunSSHJoinCLI_DirectPositional(t *testing.T) {
	withMockSSHDB(t, func(db *store.DB) {
		args := []string{"192.168.1.14", "box"}
		if err := runSSHJoinCLI(args); err != nil {
			t.Fatalf("runSSHJoinCLI direct positional failed: %v", err)
		}
		assertEnrolledHost(t, db.SQL(), "box", "192.168.1.14")
	})
}

func setupFreshStoreDB(t *testing.T) (string, *store.DB) {
	dbPath := filepath.Join(t.TempDir(), "fresh_ssh.db")
	dbConn, err := store.OpenAt(dbPath)
	if err != nil {
		t.Fatalf("failed to open fresh db: %v", err)
	}

	return dbPath, dbConn
}

func withFreshSSHDB(t *testing.T, fn func(db *store.DB)) {
	dbPath, testDB := setupFreshStoreDB(t)
	defer testDB.Close()

	prevOpener := openSSHDBFunc
	openSSHDBFunc = func() (*store.DB, error) { return store.OpenAt(dbPath) }
	defer func() { openSSHDBFunc = prevOpener }()

	fn(testDB)
}

func TestRunSSHJoinCLI_LsFreshDB(t *testing.T) {
	withFreshSSHDB(t, func(db *store.DB) {
		if err := runSSHJoinCLI([]string{"ls"}); err != nil {
			t.Fatalf("runSSHJoinCLI 'ls' failed on fresh db: %v", err)
		}
	})
}

func TestSJAddCmd_Execute(t *testing.T) {
	withMockSSHDB(t, func(db *store.DB) {
		cmd := SJAddCmd
		args := []string{"192.168.1.20", "add-box"}
		if err := cmd.RunE(cmd, args); err != nil {
			t.Fatalf("SJAddCmd RunE failed: %v", err)
		}
		assertEnrolledHost(t, db.SQL(), "add-box", "192.168.1.20")
	})
}

func TestSSHJoinCmd_RoutingAddAndPositional(t *testing.T) {
	withMockSSHDB(t, func(db *store.DB) {
		cmd := SSHJoinCmd
		if err := cmd.RunE(cmd, []string{"add", "192.168.1.21", "box-add"}); err != nil {
			t.Fatalf("SSHJoinCmd add routing failed: %v", err)
		}
		assertEnrolledHost(t, db.SQL(), "box-add", "192.168.1.21")

		if err := cmd.RunE(cmd, []string{"192.168.1.22", "box-pos"}); err != nil {
			t.Fatalf("SSHJoinCmd positional routing failed: %v", err)
		}
		assertEnrolledHost(t, db.SQL(), "box-pos", "192.168.1.22")
	})
}

func assertEnrolledHostWithUser(t *testing.T, db *sql.DB, alias, ip, user string) {
	ctx := context.Background()
	host, err := store.GetHostByAlias(ctx, alias, db)
	if err != nil {
		t.Fatalf("expected host %q to exist: %v", alias, err)
	}
	if host.IP != ip || host.Username != user {
		t.Fatalf("mismatch: got %s@%s, want %s@%s", host.Username, host.IP, user, ip)
	}
}

func TestRunSSHJoinCLI_UserAtIP(t *testing.T) {
	withMockSSHDB(t, func(db *store.DB) {
		args := []string{"alim@192.168.1.14"}
		if err := runSSHJoinCLI(args); err != nil {
			t.Fatalf("runSSHJoinCLI failed: %v", err)
		}
		assertEnrolledHostWithUser(t, db.SQL(), "host-192.168.1.14", "192.168.1.14", "alim")
	})
}

func TestRunSSHJoinCLI_UserAtIP_WithAlias(t *testing.T) {
	withMockSSHDB(t, func(db *store.DB) {
		args := []string{"alim@192.168.1.14", "mybox"}
		if err := runSSHJoinCLI(args); err != nil {
			t.Fatalf("runSSHJoinCLI failed: %v", err)
		}
		assertEnrolledHostWithUser(t, db.SQL(), "mybox", "192.168.1.14", "alim")
	})
}

func TestSJAddCmd_UserAtIP(t *testing.T) {
	withMockSSHDB(t, func(db *store.DB) {
		cmd := SJAddCmd
		args := []string{"root@192.168.1.30", "prod-server"}
		if err := cmd.RunE(cmd, args); err != nil {
			t.Fatalf("SJAddCmd RunE failed: %v", err)
		}
		assertEnrolledHostWithUser(t, db.SQL(), "prod-server", "192.168.1.30", "root")
	})
}

func TestSSHJoinCmd_PositionalUserAtIP(t *testing.T) {
	withMockSSHDB(t, func(db *store.DB) {
		cmd := SSHJoinCmd
		args := []string{"ubuntu@192.168.1.40"}
		if err := cmd.RunE(cmd, args); err != nil {
			t.Fatalf("SSHJoinCmd RunE failed: %v", err)
		}
		assertEnrolledHostWithUser(t, db.SQL(), "host-192.168.1.40", "192.168.1.40", "ubuntu")
	})
}
