package cmdssh

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

func setupBootstrapTestStore(t *testing.T) *store.DB {
	dbPath := filepath.Join(t.TempDir(), "test_cluster_bs.db")
	dbConn, err := store.OpenAt(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	_ = dbConn.Migrate()
	_ = store.EnsureSSHHostsTable(dbConn.SQL())
	return dbConn
}

func withMockBootstrapStore(t *testing.T, fn func(db *store.DB)) {
	testDB := setupBootstrapTestStore(t)
	defer testDB.Close()
	prevOpener := openClusterDBFunc
	openClusterDBFunc = func() (*store.DB, error) { return testDB, nil }
	defer func() { openClusterDBFunc = prevOpener }()
	fn(testDB)
}

func withMockKeyDir(t *testing.T, fn func(dir string)) {
	dir := t.TempDir()
	prevDir := clusterKeypairDir
	clusterKeypairDir = dir
	defer func() { clusterKeypairDir = prevDir }()
	fn(dir)
}

func resetExecutionHooks(pInj func(context.Context, store.SSHHost, string, string) error, pSud func(context.Context, store.SSHHost, string) error, pVer func(context.Context, string, store.SSHHost) error) {
	bootstrapInjectKeyFn = pInj
	bootstrapInjectSudoersFn = pSud
	bootstrapVerifyKeyAuthFn = pVer
}

func withMockExecutionHooks(fn func()) {
	pInj, pSud, pVer := bootstrapInjectKeyFn, bootstrapInjectSudoersFn, bootstrapVerifyKeyAuthFn
	bootstrapInjectKeyFn = func(ctx context.Context, h store.SSHHost, k, p string) error { return nil }
	bootstrapInjectSudoersFn = func(ctx context.Context, h store.SSHHost, p string) error { return nil }
	bootstrapVerifyKeyAuthFn = func(ctx context.Context, k string, h store.SSHHost) error { return nil }
	defer resetExecutionHooks(pInj, pSud, pVer)
	fn()
}

func TestParseClusterBootstrapArgs_BasicPositional(t *testing.T) {
	args := []string{"192.168.1.100", "mysecret"}
	opts, err := parseClusterBootstrapArgs(args)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if opts.target != "192.168.1.100" || opts.password != "mysecret" || !opts.isSudoEnabled {
		t.Fatalf("unexpected opts: %+v", opts)
	}
}

func TestParseClusterBootstrapArgs_Flags(t *testing.T) {
	args := []string{"workers", "--password", "p1", "--user", "admin", "--port", "2222", "--sudo=false"}
	opts, err := parseClusterBootstrapArgs(args)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if opts.target != "workers" || opts.password != "p1" || opts.user != "admin" || opts.port != 2222 || opts.isSudoEnabled {
		t.Fatalf("unexpected opts: %+v", opts)
	}
}

func TestParseClusterBootstrapArgs_NoSudoAndHelp(t *testing.T) {
	opts1, err1 := parseClusterBootstrapArgs([]string{"all", "--no-sudo"})
	if err1 != nil || opts1.isSudoEnabled {
		t.Fatalf("unexpected opts1: %+v err: %v", opts1, err1)
	}
	opts2, err2 := parseClusterBootstrapArgs([]string{"--help"})
	if err2 != nil || !opts2.isShowHelp {
		t.Fatalf("unexpected opts2: %+v err: %v", opts2, err2)
	}
}

func TestParseClusterBootstrapArgs_Empty(t *testing.T) {
	_, err := parseClusterBootstrapArgs([]string{})
	if err == nil {
		t.Fatal("expected error on empty args")
	}
}

func TestEnsureClusterRSAKeypair_GenerationAndDiscovery(t *testing.T) {
	withMockKeyDir(t, func(dir string) {
		privPath, pubKey, err := EnsureClusterRSAKeypair()
		if err != nil || privPath == "" || !strings.HasPrefix(pubKey, "ssh-rsa") {
			t.Fatalf("generation failed: %v, %s, %s", err, privPath, pubKey)
		}
		priv2, pub2, err2 := EnsureClusterRSAKeypair()
		if err2 != nil || priv2 != privPath || pub2 != pubKey {
			t.Fatalf("discovery mismatch: %v, %s, %s", err2, priv2, pub2)
		}
	})
}

func TestBuildScripts(t *testing.T) {
	keyScript := buildInjectPublicKeyScript("ssh-rsa KEY123")
	if !strings.Contains(keyScript, "authorized_keys") || !strings.Contains(keyScript, "ssh-rsa KEY123") {
		t.Fatalf("unexpected key script: %s", keyScript)
	}
	sudoScript := buildSudoersScript("deployer")
	if !strings.Contains(sudoScript, "deployer ALL=(ALL) NOPASSWD:ALL") || !strings.Contains(sudoScript, "0440") {
		t.Fatalf("unexpected sudoers script: %s", sudoScript)
	}
}

func TestBuildBatchVerifyArgs(t *testing.T) {
	h := store.SSHHost{IP: "10.0.0.1", Username: "root", Port: 2200}
	args := buildBatchVerifyArgs("/path/to/key", h)
	argsStr := strings.Join(args, " ")
	if !strings.Contains(argsStr, "BatchMode=yes") || !strings.Contains(argsStr, "-p 2200") || !strings.Contains(argsStr, "echo ssh_ok") {
		t.Fatalf("unexpected verify args: %s", argsStr)
	}
}

func TestBootstrapSingleTarget_Success(t *testing.T) {
	withMockExecutionHooks(func() {
		withMockBootstrapStore(t, func(db *store.DB) {
			ctx := context.Background()
			h := store.SSHHost{ID: "h-1", IP: "192.168.1.100", Username: "root"}
			opts := &clusterBootstrapOptions{isSudoEnabled: true, password: "pass"}
			res := bootstrapSingleTarget(ctx, h, "/key", "ssh-rsa PUB", opts, db.Conn())
			if res.Status != "SUCCESS" || !res.HasSudo || !res.HasKey {
				t.Fatalf("unexpected res: %+v", res)
			}
		})
	})
}

func setMockInjectFailure() {
	bootstrapInjectKeyFn = func(ctx context.Context, h store.SSHHost, k, p string) error {
		return apperror.NewExecutionError("mock inject fail")
	}
}

func TestBootstrapSingleTarget_InjectFailure(t *testing.T) {
	withMockExecutionHooks(func() {
		setMockInjectFailure()
		withMockBootstrapStore(t, func(db *store.DB) {
			h := store.SSHHost{ID: "h-1", IP: "192.168.1.100", Username: "root"}
			opts := &clusterBootstrapOptions{isSudoEnabled: true, password: "pass"}
			res := bootstrapSingleTarget(context.Background(), h, "/key", "ssh-rsa PUB", opts, db.Conn())
			if res.Status != "FAILED" || res.Err == nil {
				t.Fatalf("expected failure, got: %+v", res)
			}
		})
	})
}

func setMockVerifyFailure() {
	bootstrapVerifyKeyAuthFn = func(ctx context.Context, k string, h store.SSHHost) error {
		return apperror.NewExecutionError("mock verify fail")
	}
}

func TestBootstrapSingleTarget_VerifyFailure(t *testing.T) {
	withMockExecutionHooks(func() {
		setMockVerifyFailure()
		withMockBootstrapStore(t, func(db *store.DB) {
			h := store.SSHHost{ID: "h-1", IP: "192.168.1.100", Username: "root"}
			opts := &clusterBootstrapOptions{isSudoEnabled: true, password: "pass"}
			res := bootstrapSingleTarget(context.Background(), h, "/key", "ssh-rsa PUB", opts, db.Conn())
			if res.Status != "FAILED" || res.Err == nil {
				t.Fatalf("expected failure, got: %+v", res)
			}
		})
	})
}

func assertBootstrapHostEnrolled(t *testing.T, db *store.DB, ip string) {
	h, err := store.GetHostByIP(context.Background(), ip, db.Conn())
	if err != nil || h.IP != ip {
		t.Fatalf("host enrollment check failed: %v", err)
	}
}

func TestRunClusterBootstrapCLI_Success(t *testing.T) {
	withMockKeyDir(t, func(dir string) {
		withMockExecutionHooks(func() {
			withMockBootstrapStore(t, func(db *store.DB) {
				err := RunClusterBootstrapCLI([]string{"192.168.1.200", "pass123"})
				if err != nil {
					t.Fatalf("RunClusterBootstrapCLI failed: %v", err)
				}
				assertBootstrapHostEnrolled(t, db, "192.168.1.200")
			})
		})
	})
}

func TestRunClusterBootstrapCLI_Help(t *testing.T) {
	err := RunClusterBootstrapCLI([]string{"--help"})
	if err != nil {
		t.Fatalf("unexpected help error: %v", err)
	}
}

func TestRenderBootstrapSummaryTable(t *testing.T) {
	res := []BootstrapResult{
		{Node: "node-1", IP: "192.168.1.1", HasSudo: true, HasKey: true, Status: "SUCCESS", Duration: 50 * time.Millisecond},
		{Node: "node-2", IP: "192.168.1.2", HasSudo: false, HasKey: false, Status: "FAILED", Duration: 10 * time.Millisecond},
	}
	table := RenderBootstrapSummaryTable(res)
	if !strings.Contains(table, "NODE") || !strings.Contains(table, "KEY_AUTH") || !strings.Contains(table, "node-1") {
		t.Fatalf("unexpected table: %s", table)
	}
}

func TestResolveStoredPassword(t *testing.T) {
	hNoPass := store.SSHHost{IP: "1.1.1.1"}
	pass, err := resolveStoredPassword(hNoPass)
	if err != nil || pass != "" {
		t.Fatalf("unexpected pass: %s, %v", pass, err)
	}
}

func TestWipePasswordString(t *testing.T) {
	s := "secretPassword"
	wipePasswordString(&s)
	if s != "" {
		t.Fatalf("expected wiped string, got: %s", s)
	}
}
