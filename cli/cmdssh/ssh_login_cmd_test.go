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

func TestRunSSHLogin_InterceptsNodes(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestHost(t, testDB, "h3", "prod-node", "10.0.0.3", "admin")
	hookTestDB(t, testDB)

	err := runSSHLogin(nil, []string{"nodes"}, context.Background())
	if err != nil {
		t.Fatalf("expected nil error for nodes intercept, got: %v", err)
	}
}

func TestRunSSHLogin_InterceptsLs(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestHost(t, testDB, "h4", "worker-node", "10.0.0.4", "admin")
	hookTestDB(t, testDB)

	err := runSSHLogin(nil, []string{"ls"}, context.Background())
	if err != nil {
		t.Fatalf("expected nil error for ls intercept, got: %v", err)
	}
}

func hookPasswordSpawner(t *testing.T, outPass *string) {
	orig := spawnSSHFn
	spawnSSHFn = func(ctx context.Context, target SSHTarget, args []string, password string) error {
		*outPass = password
		return nil
	}
	t.Cleanup(func() { spawnSSHFn = orig })
}

func TestRunSSHLogin_WithPasswordSavesEncrypted(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestHost(t, testDB, "h5", "w2", "10.0.0.5", "admin")
	hookTestDB(t, testDB)
	var capturedPass string
	hookPasswordSpawner(t, &capturedPass)

	err := runSSHLogin(nil, []string{"w2", "secretPass123"}, context.Background())
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if capturedPass != "secretPass123" {
		t.Errorf("expected secretPass123, got: %s", capturedPass)
	}
	verifyEncryptedHostPassword(t, testDB, "w2")
}

func verifyEncryptedHostPassword(t *testing.T, db *store.DB, alias string) {
	host, err := store.GetHostByAlias(context.Background(), alias, db.Conn())
	if err != nil {
		t.Fatalf("failed to retrieve host: %v", err)
	}
	isRSA := strings.HasPrefix(host.EncryptedPassword, "rsa:")
	isAES := strings.HasPrefix(host.EncryptedPassword, "aes:")
	if isRSA || isAES {
		return
	}
	t.Errorf("expected rsa: or aes: prefix in encrypted password, got: %s", host.EncryptedPassword)
}

func TestRunSSHLogin_SubsequentLoginWithoutPassword(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestHost(t, testDB, "h6", "w3", "10.0.0.6", "admin")
	hookTestDB(t, testDB)
	var capturedPass string
	hookPasswordSpawner(t, &capturedPass)

	_ = runSSHLogin(nil, []string{"w3", "pass789"}, context.Background())
	capturedPass = ""

	err := runSSHLogin(nil, []string{"w3"}, context.Background())
	if err != nil {
		t.Fatalf("expected nil error on subsequent login, got: %v", err)
	}
	if capturedPass != "pass789" {
		t.Errorf("expected pass789 decrypted automatically, got: %s", capturedPass)
	}
}

func TestRunSSHLogin_DirectIPWithPassword(t *testing.T) {
	testDB := setupTestDB(t)
	hookTestDB(t, testDB)
	var capturedPass string
	hookPasswordSpawner(t, &capturedPass)

	err := runSSHLogin(nil, []string{"192.168.1.88", "ipSecret999"}, context.Background())
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if capturedPass != "ipSecret999" {
		t.Errorf("expected ipSecret999, got: %s", capturedPass)
	}
	testDirectIPRecall(t, &capturedPass)
}

func testDirectIPRecall(t *testing.T, capturedPass *string) {
	*capturedPass = ""
	err := runSSHLogin(nil, []string{"192.168.1.88"}, context.Background())
	if err != nil {
		t.Fatalf("expected nil error on recall, got: %v", err)
	}
	if *capturedPass != "ipSecret999" {
		t.Errorf("expected ipSecret999 recalled, got: %s", *capturedPass)
	}
}
