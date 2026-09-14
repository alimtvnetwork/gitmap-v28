package store

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func setupSSHTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}

	err = RegisterSSHHostMigration(db, 1, false)
	if err != nil {
		t.Fatalf("failed to run migration: %v", err)
	}

	return db
}

func createTestHost(id, alias, ip, user string) SSHHost {
	return SSHHost{
		ID:        id,
		Alias:     alias,
		IP:        ip,
		Username:  user,
		CreatedAt: time.Now().Truncate(time.Second).UTC(),
	}
}

func countHostRows(t *testing.T, db *sql.DB, query string) int {
	var count int
	if err := db.QueryRow(query).Scan(&count); err != nil {
		t.Fatalf("failed to query count: %v", err)
	}

	return count
}

func assertRowCount(t *testing.T, db *sql.DB, query string, expected int) {
	if count := countHostRows(t, db, query); count != expected {
		t.Errorf("expected %d records for %q, got %d", expected, query, count)
	}
}

func assertHostFound(t *testing.T, ctx context.Context, db *sql.DB, alias, expectedIP, expectedUser string) {
	found, err := GetHostByAlias(ctx, alias, db)
	if err != nil {
		t.Fatalf("GetHostByAlias %q failed: %v", alias, err)
	}

	if found.IP != expectedIP || found.Username != expectedUser {
		t.Errorf("host %q: expected ip=%s user=%s, got ip=%s user=%s",
			alias, expectedIP, expectedUser, found.IP, found.Username)
	}
}

func insertHostInTx(t *testing.T, ctx context.Context, db *sql.DB, host SSHHost) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}

	if err = InsertSSHHost(ctx, host, tx); err != nil {
		t.Fatalf("InsertSSHHost failed: %v", err)
	}

	if err = tx.Commit(); err != nil {
		t.Fatalf("failed to commit: %v", err)
	}
}

func TestInsertSSHHost(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()

	ctx := context.Background()
	host := createTestHost("host-1", "my-server", "192.168.1.100", "admin")
	insertHostInTx(t, ctx, db, host)
	assertRowCount(t, db, "SELECT COUNT(*) FROM ssh_hosts WHERE id = 'host-1'", 1)
	testDuplicateInsertError(t, ctx, db, host)
}

func assertInsertInternalError(t *testing.T, err error) {
	appErr, isAppError := err.(*apperror.AppError)
	if !isAppError || appErr.Code != "E_INTERNAL_ERROR" {
		t.Fatalf("expected E_INTERNAL_ERROR AppError, got %v", err)
	}
}

func testDuplicateInsertError(t *testing.T, ctx context.Context, db *sql.DB, host SSHHost) {
	tx, _ := db.BeginTx(ctx, nil)
	defer tx.Rollback()

	err := InsertSSHHost(ctx, host, tx)
	if err == nil {
		t.Fatal("expected error on duplicate insert, got nil")
	}
	assertInsertInternalError(t, err)
}

func TestGetHostByAlias(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()

	ctx := context.Background()
	host := createTestHost("host-2", "web-server", "192.168.1.101", "root")
	insertHostInTx(t, ctx, db, host)
	assertHostFound(t, ctx, db, "web-server", "192.168.1.101", "root")
	assertHostNotFoundError(t, ctx, db, "db-server")
}

func assertHostNotFoundError(t *testing.T, ctx context.Context, db *sql.DB, alias string) {
	_, err := GetHostByAlias(ctx, alias, db)
	if err == nil {
		t.Fatal("expected error for non-existent alias, got nil")
	}

	appErr, isAppError := err.(*apperror.AppError)
	if !isAppError {
		t.Fatalf("expected AppError, got %T", err)
	}

	if !errors.Is(appErr.Cause, apperror.ErrNotFound) {
		t.Errorf("expected cause ErrNotFound, got %v", appErr.Cause)
	}
}

func TestDeleteHostByIP(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()

	ctx := context.Background()
	host := createTestHost("host-3", "delete-me", "10.0.0.99", "ubuntu")
	insertHostInTx(t, ctx, db, host)

	if err := DeleteHostByIP(ctx, "10.0.0.99", db); err != nil {
		t.Fatalf("DeleteHostByIP failed: %v", err)
	}

	assertRowCount(t, db, "SELECT COUNT(*) FROM ssh_hosts WHERE ip = '10.0.0.99'", 0)
}

func TestUpsertSSHHost_Insert(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()

	ctx := context.Background()
	host := createTestHost("h-up-1", "alias-1", "10.0.1.1", "admin")
	if err := UpsertSSHHost(ctx, host, db); err != nil {
		t.Fatalf("UpsertSSHHost insert failed: %v", err)
	}

	assertHostFound(t, ctx, db, "alias-1", "10.0.1.1", "admin")
}

func TestUpsertSSHHost_UpdateByIP(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()

	ctx := context.Background()
	_ = UpsertSSHHost(ctx, createTestHost("h-up-2", "old-alias", "10.0.1.2", "admin"), db)
	host2 := createTestHost("h-up-2-new", "new-alias", "10.0.1.2", "superadmin")
	if err := UpsertSSHHost(ctx, host2, db); err != nil {
		t.Fatalf("UpsertSSHHost update by IP failed: %v", err)
	}

	assertHostFound(t, ctx, db, "new-alias", "10.0.1.2", "superadmin")
	assertRowCount(t, db, "SELECT COUNT(*) FROM ssh_hosts WHERE ip = '10.0.1.2'", 1)
}

func TestUpsertSSHHost_UpdateByAlias(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()

	ctx := context.Background()
	_ = UpsertSSHHost(ctx, createTestHost("h-up-3", "shared-alias", "10.0.1.3", "admin"), db)
	host2 := createTestHost("h-up-3-new", "shared-alias", "10.0.1.99", "newuser")
	if err := UpsertSSHHost(ctx, host2, db); err != nil {
		t.Fatalf("UpsertSSHHost update by alias failed: %v", err)
	}

	assertHostFound(t, ctx, db, "shared-alias", "10.0.1.99", "newuser")
}

func TestDeleteHostByAliasOrIP_Alias(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()

	ctx := context.Background()
	_ = UpsertSSHHost(ctx, createTestHost("h-del-1", "target-alias", "10.0.2.1", "root"), db)
	affected, err := DeleteHostByAliasOrIP(ctx, "target-alias", db)
	if err != nil || affected != 1 {
		t.Fatalf("DeleteHostByAliasOrIP failed: err=%v affected=%d", err, affected)
	}

	assertRowCount(t, db, "SELECT COUNT(*) FROM ssh_hosts WHERE alias = 'target-alias'", 0)
}

func TestDeleteHostByAliasOrIP_IP(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()

	ctx := context.Background()
	_ = UpsertSSHHost(ctx, createTestHost("h-del-2", "ip-alias", "10.0.2.2", "root"), db)
	affected, err := DeleteHostByAliasOrIP(ctx, "10.0.2.2", db)
	if err != nil || affected != 1 {
		t.Fatalf("DeleteHostByAliasOrIP failed: err=%v affected=%d", err, affected)
	}

	assertRowCount(t, db, "SELECT COUNT(*) FROM ssh_hosts WHERE ip = '10.0.2.2'", 0)
}

func TestEnrollSSHHost_Success(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()

	ctx := context.Background()
	host := createTestHost("h-en-1", "enroll-box", "10.0.3.1", "dev")
	hist := SSHHistory{ID: "hist-en-1", HostIP: "10.0.3.1", User: "dev", JoinedAt: time.Now().UTC()}
	if err := EnrollSSHHost(ctx, host, hist, db); err != nil {
		t.Fatalf("EnrollSSHHost failed: %v", err)
	}

	assertRowCount(t, db, "SELECT COUNT(*) FROM ssh_hosts WHERE id = 'h-en-1'", 1)
	assertRowCount(t, db, "SELECT COUNT(*) FROM ssh_history WHERE id = 'hist-en-1'", 1)
}

func seedTestHistory(t *testing.T, ctx context.Context, db *sql.DB, id string) {
	hist := SSHHistory{ID: id, HostIP: "10.0.4.1", User: "dev", JoinedAt: time.Now().UTC()}
	if err := LogSSHHistory(ctx, hist, db); err != nil {
		t.Fatalf("failed to seed test history: %v", err)
	}
}

func verifyEnrollRollback(t *testing.T, ctx context.Context, db *sql.DB) {
	host := createTestHost("h-rb-1", "rollback-box", "10.0.4.1", "dev")
	hist := SSHHistory{ID: "hist-rb-1", HostIP: "10.0.4.1", User: "dev", JoinedAt: time.Now().UTC()}
	if err := EnrollSSHHost(ctx, host, hist, db); err == nil {
		t.Fatal("expected error due to duplicate history id, got nil")
	}
	assertRowCount(t, db, "SELECT COUNT(*) FROM ssh_hosts WHERE id = 'h-rb-1'", 0)
}

func TestEnrollSSHHost_Rollback(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()

	ctx := context.Background()
	seedTestHistory(t, ctx, db, "hist-rb-1")
	verifyEnrollRollback(t, ctx, db)
}

func TestSSHRepoIntegration(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()

	ctx := context.Background()
	host := createTestHost("host-int-1", "int-server", "10.0.0.50", "int-user")
	insertHostInTx(t, ctx, db, host)

	hist := SSHHistory{ID: "hist-1", HostIP: "10.0.0.50", User: "int-user", JoinedAt: host.CreatedAt}
	_ = LogSSHJoin(ctx, hist, db)

	assertRowCount(t, db, "SELECT COUNT(*) FROM ssh_hosts WHERE id = 'host-int-1'", 1)
	assertRowCount(t, db, "SELECT COUNT(*) FROM ssh_history WHERE id = 'hist-1'", 1)
}

func assertTableExists(t *testing.T, db *sql.DB, tableName string) {
	query := `SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = ?`
	var count int
	if err := db.QueryRow(query, tableName).Scan(&count); err != nil {
		t.Fatalf("failed to query sqlite_master for table %s: %v", tableName, err)
	}
	if count != 1 {
		t.Errorf("expected table %s to exist in fresh db, got count %d", tableName, count)
	}
}

func setupFreshMemoryDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}
	return db
}

func TestListHosts_FreshDB(t *testing.T) {
	db := setupFreshMemoryDB(t)
	defer db.Close()

	ctx := context.Background()
	hosts, err := ListHosts(ctx, db)
	if err != nil || len(hosts) != 0 {
		t.Fatalf("expected 0 hosts on fresh db, got %d (err: %v)", len(hosts), err)
	}
	assertTableExists(t, db, "ssh_hosts")
	assertTableExists(t, db, "ssh_history")
}

func TestGetHostByIP(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()

	ctx := context.Background()
	insertHostInTx(t, ctx, db, createTestHost("host-ip-1", "ip-server", "192.168.1.150", "admin"))
	found, err := GetHostByIP(ctx, "192.168.1.150", db)
	if err != nil || found.ID != "host-ip-1" {
		t.Fatalf("GetHostByIP failed: id=%s err=%v", found.ID, err)
	}
}

func TestGetHostByID(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()

	ctx := context.Background()
	insertHostInTx(t, ctx, db, createTestHost("host-id-1", "id-server", "192.168.1.151", "root"))
	found, err := GetHostByID(ctx, "host-id-1", db)
	if err != nil || found.Alias != "id-server" {
		t.Fatalf("GetHostByID failed: alias=%s err=%v", found.Alias, err)
	}
}

func TestLogSSHHistory(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()

	ctx := context.Background()
	hist := SSHHistory{ID: "hist-def-1", HostIP: "10.0.9.1", User: "tester", JoinedAt: time.Now().UTC()}
	if err := LogSSHHistory(ctx, hist, db); err != nil {
		t.Fatalf("LogSSHHistory failed: %v", err)
	}

	assertRowCount(t, db, "SELECT COUNT(*) FROM ssh_history WHERE id = 'hist-def-1'", 1)
}

func createEncryptedHost() SSHHost {
	return SSHHost{
		ID:                "host-pass-1",
		Alias:             "pass-server",
		IP:                "192.168.1.199",
		Username:          "deploy",
		Port:              2222,
		EncryptedPassword: "rsa:encrypted-secret-token",
		CreatedAt:         time.Now().UTC(),
	}
}

func TestSSHHost_EncryptedPasswordAndPort(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()

	ctx := context.Background()
	insertHostInTx(t, ctx, db, createEncryptedHost())
	found, err := GetHostByAlias(ctx, "pass-server", db)
	if err != nil || found.Port != 2222 || found.EncryptedPassword != "rsa:encrypted-secret-token" {
		t.Fatalf("expected port 2222 and encrypted token, got %+v (err: %v)", found, err)
	}
}

func seedClusterHosts(t *testing.T, ctx context.Context, db *sql.DB) {
	h1 := SSHHost{ID: "h-c1", Alias: "ctrl-1", IP: "10.1.0.1", Username: "root", ClusterRole: "control"}
	h2 := SSHHost{ID: "h-w1", Alias: "work-1", IP: "10.1.0.2", Username: "root", ClusterRole: "worker"}
	h3 := SSHHost{ID: "h-w2", Alias: "work-2", IP: "10.1.0.3", Username: "root", ClusterRole: "worker"}
	insertHostInTx(t, ctx, db, h1)
	insertHostInTx(t, ctx, db, h2)
	insertHostInTx(t, ctx, db, h3)
}

func TestClusterRole_Persistence(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()
	ctx := context.Background()

	host := SSHHost{ID: "h-role-1", Alias: "role-ctrl", IP: "10.2.0.1", Username: "root", ClusterRole: "control"}
	insertHostInTx(t, ctx, db, host)
	found, err := GetHostByAlias(ctx, "role-ctrl", db)
	if err != nil || found.ClusterRole != "control" {
		t.Fatalf("expected control role, got %s (err: %v)", found.ClusterRole, err)
	}
}

func TestClusterRole_DefaultWorker(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()
	ctx := context.Background()

	host := createTestHost("h-role-def", "role-def", "10.2.0.2", "root")
	insertHostInTx(t, ctx, db, host)
	found, err := GetHostByAlias(ctx, "role-def", db)
	if err != nil || found.ClusterRole != "worker" {
		t.Fatalf("expected default worker role, got %s (err: %v)", found.ClusterRole, err)
	}
}

func TestClusterRole_Update(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()
	ctx := context.Background()

	host := SSHHost{ID: "h-role-up", Alias: "role-up", IP: "10.2.0.3", Username: "root", ClusterRole: "control"}
	_ = UpsertSSHHost(ctx, host, db)
	host.ClusterRole = "worker"
	_ = UpsertSSHHost(ctx, host, db)
	found, err := GetHostByAlias(ctx, "role-up", db)
	if err != nil || found.ClusterRole != "worker" {
		t.Fatalf("expected updated worker role, got %s (err: %v)", found.ClusterRole, err)
	}
}

func TestListHostsByRole(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()
	ctx := context.Background()
	seedClusterHosts(t, ctx, db)

	ctrls, err := ListHostsByRole(ctx, "control", db)
	if err != nil || len(ctrls) != 1 {
		t.Fatalf("expected 1 control host, got %d (err: %v)", len(ctrls), err)
	}
	workers, err := ListHostsByRole(ctx, "worker", db)
	if err != nil || len(workers) != 2 {
		t.Fatalf("expected 2 worker hosts, got %d (err: %v)", len(workers), err)
	}
}

func assertTargetCount(t *testing.T, ctx context.Context, db *sql.DB, target string, expected int) {
	hosts, err := ListHostsByTarget(ctx, target, db)
	if err != nil {
		t.Fatalf("ListHostsByTarget %q failed: %v", target, err)
	}
	if len(hosts) != expected {
		t.Errorf("target %q: expected %d hosts, got %d", target, expected, len(hosts))
	}
}

func TestListHostsByTarget_Groups(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()
	ctx := context.Background()
	seedClusterHosts(t, ctx, db)

	assertTargetCount(t, ctx, db, "all", 3)
	assertTargetCount(t, ctx, db, "", 3)
	assertTargetCount(t, ctx, db, "control", 1)
	assertTargetCount(t, ctx, db, "master", 1)
	assertTargetCount(t, ctx, db, "workers", 2)
	assertTargetCount(t, ctx, db, "worker", 2)
	assertTargetCount(t, ctx, db, "nodes", 2)
}

func TestListHostsByTarget_AliasAndIP(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()
	ctx := context.Background()
	seedClusterHosts(t, ctx, db)

	assertTargetCount(t, ctx, db, "ctrl-1", 1)
	assertTargetCount(t, ctx, db, "10.1.0.2", 1)
}

func TestListHostsByTarget_CommaAndDeduplication(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()
	ctx := context.Background()
	seedClusterHosts(t, ctx, db)

	assertTargetCount(t, ctx, db, "ctrl-1,work-1", 2)
	assertTargetCount(t, ctx, db, "control,ctrl-1", 1)
	assertTargetCount(t, ctx, db, "work-1, work-2", 2)
}

func TestListHostsByTarget_NotFound(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()
	ctx := context.Background()

	_, err := ListHostsByTarget(ctx, "nonexistent-target", db)
	if err == nil {
		t.Fatal("expected error for nonexistent target, got nil")
	}
	appErr, isAppError := err.(*apperror.AppError)
	if !isAppError || !errors.Is(appErr.Cause, apperror.ErrNotFound) {
		t.Fatalf("expected ErrNotFound AppError, got %v", err)
	}
}
