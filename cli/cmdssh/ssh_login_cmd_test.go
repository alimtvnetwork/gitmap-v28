package cmdssh

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func setupTestDB(t *testing.T) *store.DB {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	testDB, err := store.OpenAt(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	t.Cleanup(func() { testDB.Close() })
	return testDB
}

func hookTestDB(t *testing.T, db *store.DB) {
	orig := openSSHDB
	openSSHDB = func() (*store.DB, error) { return db, nil }
	t.Cleanup(func() { openSSHDB = orig })
}

func hookTestSpawner(t *testing.T) *int {
	count := 0
	orig := spawnSSHFn
	spawnSSHFn = func(ctx context.Context, target SSHTarget, args []string, password string) error {
		count++
		return nil
	}
	t.Cleanup(func() { spawnSSHFn = orig })
	return &count
}

func hookTargetSpawner(t *testing.T, out *SSHTarget) {
	orig := spawnSSHFn
	spawnSSHFn = func(ctx context.Context, target SSHTarget, args []string, password string) error {
		*out = target
		return nil
	}
	t.Cleanup(func() { spawnSSHFn = orig })
}

func hookTestJoin(t *testing.T) (*bool, *[]string) {
	isCalled := false
	var capturedArgs []string
	orig := runSSHJoinFn
	runSSHJoinFn = func(args []string) error {
		isCalled = true
		capturedArgs = args
		return nil
	}
	t.Cleanup(func() { runSSHJoinFn = orig })
	return &isCalled, &capturedArgs
}

func seedTestHost(t *testing.T, db *store.DB, id, alias, ip, user string) {
	host := store.SSHHost{ID: id, Alias: alias, IP: ip, Username: user}
	if err := store.UpsertSSHHost(context.Background(), host, db.Conn()); err != nil {
		t.Fatalf("failed to seed host: %v", err)
	}
}

func assertErrorContains(t *testing.T, msg, substr, label string) {
	if hasSubstr := strings.Contains(msg, substr); !hasSubstr {
		t.Errorf("expected message to contain %s %q, got: %s", label, substr, msg)
	}
}

func assertUnknownAliasError(t *testing.T, err error, target, expectedHost string) {
	if err == nil {
		t.Fatalf("expected error for unknown alias, got nil")
	}
	msg := err.Error()
	pat := fmt.Sprintf("SSH host alias '%s' not found in registry", target)
	assertErrorContains(t, msg, pat, "alias pattern")
	assertErrorContains(t, msg, expectedHost, "host")
	assertErrorContains(t, msg, "gitmap ssh join", "join example")
}

func assertJoinIntercept(t *testing.T, isCalled *bool, args *[]string, spawns *int, want string) {
	if !*isCalled {
		t.Errorf("expected RunSSHJoinCLI to be called")
	}
	if len(*args) != 1 || (*args)[0] != want {
		t.Errorf("unexpected args passed: %v", *args)
	}
	if *spawns != 0 {
		t.Errorf("expected 0 spawns, got %d", *spawns)
	}
}

func TestRunSSHLogin_InterceptsJoin(t *testing.T) {
	spawns := hookTestSpawner(t)
	isCalled, passedArgs := hookTestJoin(t)

	err := runSSHLogin(nil, []string{"join", "192.168.1.50"}, context.Background())
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	assertJoinIntercept(t, isCalled, passedArgs, spawns, "192.168.1.50")
}

func TestRunSSHLogin_InterceptsSJ(t *testing.T) {
	spawns := hookTestSpawner(t)
	isCalled, passedArgs := hookTestJoin(t)

	err := runSSHLogin(nil, []string{"sj", "10.0.0.99"}, context.Background())
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	assertJoinIntercept(t, isCalled, passedArgs, spawns, "10.0.0.99")
}

func TestExecuteSSHLogin_UnknownAliasRejectionWithSuggestions(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestHost(t, testDB, "h1", "prod-web", "10.0.0.1", "admin")
	hookTestDB(t, testDB)
	spawnCount := hookTestSpawner(t)

	err := executeSSHLogin(context.Background(), "unknown-server", false)
	assertUnknownAliasError(t, err, "unknown-server", "prod-web")
	if *spawnCount != 0 {
		t.Errorf("expected 0 spawns, got %d", *spawnCount)
	}
}

func assertEmptyRegistryError(t *testing.T, err error, spawnCount *int) {
	if err == nil {
		t.Fatalf("expected error for empty registry lookup, got nil")
	}
	if *spawnCount != 0 {
		t.Errorf("expected 0 spawns, got %d", *spawnCount)
	}
	assertErrorContains(t, err.Error(), "(none)", "none table indicator")
}

func TestExecuteSSHLogin_UnknownAliasEmptyRegistry(t *testing.T) {
	testDB := setupTestDB(t)
	hookTestDB(t, testDB)
	spawnCount := hookTestSpawner(t)

	err := executeSSHLogin(context.Background(), "ghost", false)
	assertEmptyRegistryError(t, err, spawnCount)
}

func TestExecuteSSHLogin_KnownAliasSpawns(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestHost(t, testDB, "h2", "db-box", "10.0.0.2", "postgres")
	hookTestDB(t, testDB)
	var spawnedTarget SSHTarget
	hookTargetSpawner(t, &spawnedTarget)

	err := executeSSHLogin(context.Background(), "db-box", false)
	if err != nil {
		t.Fatalf("expected nil error for known alias, got: %v", err)
	}
	if spawnedTarget.IP != "10.0.0.2" || spawnedTarget.Username != "postgres" {
		t.Errorf("unexpected target spawned: %+v", spawnedTarget)
	}
}

func TestExecuteSSHLogin_DirectIPSpawnsDirectly(t *testing.T) {
	var spawnedTarget SSHTarget
	hookTargetSpawner(t, &spawnedTarget)

	err := executeSSHLogin(context.Background(), "192.168.1.1", false)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if spawnedTarget.IP != "192.168.1.1" {
		t.Errorf("expected IP 192.168.1.1, got: %s", spawnedTarget.IP)
	}
}

func TestExecuteSSHLogin_UserAtHostSpawnsDirectly(t *testing.T) {
	var spawnedTarget SSHTarget
	hookTargetSpawner(t, &spawnedTarget)

	err := executeSSHLogin(context.Background(), "deploy@remote.lan", false)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if spawnedTarget.IP != "remote.lan" || spawnedTarget.Username != "deploy" {
		t.Errorf("unexpected target: %+v", spawnedTarget)
	}
}
