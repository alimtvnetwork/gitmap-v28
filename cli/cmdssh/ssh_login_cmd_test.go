package cmdssh

import (
	"bufio"
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"golang.org/x/crypto/ssh"
)

var testLoginDBPath string

func setupTestDB(t *testing.T) *store.DB {
	testLoginDBPath = filepath.Join(t.TempDir(), "test.db")
	testDB, err := store.OpenAt(testLoginDBPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	_ = store.EnsureSSHHostsTable(testDB.SQL())
	_ = store.EnsureSSHHistoryTable(testDB.SQL())
	t.Cleanup(func() { testDB.Close() })
	return testDB
}

func hookTestDB(t *testing.T, db *store.DB) {
	orig := openSSHDB
	openSSHDB = func() (*store.DB, error) {
		if testLoginDBPath != "" {
			return store.OpenAt(testLoginDBPath)
		}
		return db, nil
	}
	t.Cleanup(func() { openSSHDB = orig })
	t.Cleanup(SetConnectProbeClientForTesting(func(*SSHTarget, string) *ssh.Client { return nil }))
	t.Cleanup(SetAutoTrustTargetHostForTesting(func(context.Context, *SSHTarget) {}))
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

func checkConsentAffirmativeCases(t *testing.T, cases []struct {
	input string
	want  bool
}) {
	for _, tc := range cases {
		if got := isConsentAffirmative(tc.input); got != tc.want {
			t.Errorf("isConsentAffirmative(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestIsConsentAffirmative(t *testing.T) {
	checkConsentAffirmativeCases(t, []struct {
		input string
		want  bool
	}{
		{"y", true}, {"Y", true}, {"yes", true}, {"YES", true}, {"  y  ", true},
		{"n", false}, {"no", false}, {"", false}, {"invalid", false},
	})
}

func TestFormatPrompts(t *testing.T) {
	passPrompt := formatPasswordPrompt("admin", "192.168.1.10")
	if passPrompt != "Enter password for admin@192.168.1.10: " {
		t.Errorf("unexpected password prompt: %q", passPrompt)
	}

	consentPrompt := formatConsentPrompt()
	if !strings.Contains(consentPrompt, "RSA algorithm") {
		t.Errorf("expected consent prompt to mention RSA, got: %q", consentPrompt)
	}
}

func hookKeyAuthSuccess(t *testing.T) *bool {
	isPrompted := false
	origKey, origPass := tryConnectKeyHook, readPasswordHook
	tryConnectKeyHook = func(*SSHTarget) *ssh.Client { return &ssh.Client{} }
	readPasswordHook = func(int) ([]byte, error) {
		isPrompted = true
		return []byte("pass"), nil
	}
	t.Cleanup(func() {
		tryConnectKeyHook, readPasswordHook = origKey, origPass
	})

	return &isPrompted
}

func TestInterceptSSHPassword_KeyAuthBypass(t *testing.T) {
	isPrompted := hookKeyAuthSuccess(t)
	target := &SSHTarget{Username: "root", IP: "10.0.0.1", Port: 22}

	pass, err := interceptSSHPasswordIfNeeded(context.Background(), "t1", target, "")
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if pass != "" || *isPrompted {
		t.Errorf("expected empty password and no prompt, got pass=%q prompted=%v", pass, *isPrompted)
	}
}

func setupInteractiveMockHooks(t *testing.T, consent string, pass string) {
	origTerm, origKey := isTerminalHook, tryConnectKeyHook
	origPass, origDial, origConsent := readPasswordHook, dialTargetPassHook, readConsentHook
	isTerminalHook = func() bool { return true }
	tryConnectKeyHook = func(*SSHTarget) *ssh.Client { return nil }
	dialTargetPassHook = func(*SSHTarget, string) (*ssh.Client, error) { return nil, nil }
	readPasswordHook = func(int) ([]byte, error) { return []byte(pass), nil }
	readConsentHook = func(byte) (string, error) { return consent, nil }
	t.Cleanup(func() {
		isTerminalHook, tryConnectKeyHook = origTerm, origKey
		readPasswordHook, dialTargetPassHook, readConsentHook = origPass, origDial, origConsent
	})
}

func TestInterceptSSHPassword_ConsentAffirmative(t *testing.T) {
	testDB := setupTestDB(t)
	hookTestDB(t, testDB)
	setupInteractiveMockHooks(t, "y\n", "secret123")
	target := &SSHTarget{Username: "admin", IP: "10.0.0.88", Port: 22}

	pass, err := interceptSSHPasswordIfNeeded(context.Background(), "vault-box", target, "")
	if err != nil || pass != "secret123" {
		t.Fatalf("expected pass secret123, got: %s (err: %v)", pass, err)
	}

	verifyEncryptedHostPassword(t, testDB, "vault-box")
}

func assertHostPasswordNotSaved(t *testing.T, db *store.DB, alias string) {
	host, err := store.GetHostByAlias(context.Background(), alias, db.Conn())
	if err == nil && host.EncryptedPassword != "" {
		t.Errorf("expected no saved password for host %s, got: %s", alias, host.EncryptedPassword)
	}
}

func TestInterceptSSHPassword_ConsentNegative(t *testing.T) {
	testDB := setupTestDB(t)
	hookTestDB(t, testDB)
	setupInteractiveMockHooks(t, "n\n", "sessionOnly456")
	target := &SSHTarget{Username: "admin", IP: "10.0.0.89", Port: 22}

	pass, err := interceptSSHPasswordIfNeeded(context.Background(), "ephemeral-box", target, "")
	if err != nil || pass != "sessionOnly456" {
		t.Fatalf("expected pass sessionOnly456, got: %s (err: %v)", pass, err)
	}

	assertHostPasswordNotSaved(t, testDB, "ephemeral-box")
}

func hookNonInteractiveTerminal(t *testing.T) *bool {
	origTerm := isTerminalHook
	isTerminalHook = func() bool { return false }
	t.Cleanup(func() { isTerminalHook = origTerm })

	isPrompted := hookKeyAuthSuccess(t)
	*isPrompted = false

	return isPrompted
}

func TestInterceptSSHPassword_NonInteractive(t *testing.T) {
	isPrompted := hookNonInteractiveTerminal(t)
	target := &SSHTarget{Username: "ci-user", IP: "10.0.0.90", Port: 22}

	pass, err := interceptSSHPasswordIfNeeded(context.Background(), "ci-box", target, "")
	if err != nil {
		t.Fatalf("expected nil error for non-interactive pass-through, got: %v", err)
	}

	if pass != "" || *isPrompted {
		t.Errorf("expected empty password and no prompt, got pass=%q prompted=%v", pass, *isPrompted)
	}
}

func TestReadInteractiveConsent(t *testing.T) {
	rYes := bufio.NewReader(strings.NewReader("yes\n"))
	if isOk, err := readInteractiveConsent(rYes); err != nil || !isOk {
		t.Errorf("expected yes to be true, got %v (err: %v)", isOk, err)
	}

	rNo := bufio.NewReader(strings.NewReader("no\n"))
	if isOk, err := readInteractiveConsent(rNo); err != nil || isOk {
		t.Errorf("expected no to be false, got %v (err: %v)", isOk, err)
	}
}

func setupVerifyTargetTestHook(t *testing.T) {
	origDial := dialTargetPassHook
	dialTargetPassHook = func(target *SSHTarget, pass string) (*ssh.Client, error) {
		if pass == "valid" {
			return nil, nil
		}

		return nil, fmt.Errorf("auth fail")
	}
	t.Cleanup(func() { dialTargetPassHook = origDial })
}

func TestVerifyTargetPassword(t *testing.T) {
	setupVerifyTargetTestHook(t)
	target := &SSHTarget{Username: "admin", IP: "10.0.0.1"}

	if !verifyTargetPassword(target, "valid") {
		t.Errorf("expected true for valid password")
	}

	if verifyTargetPassword(target, "wrong") {
		t.Errorf("expected false for wrong password")
	}
}
