package cmdssh

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

func setupExecClusterTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}
	_ = store.EnsureSSHTables(db)
	return db
}

func setupClusterTestStoreDB(t *testing.T) (string, *store.DB) {
	dbPath := filepath.Join(t.TempDir(), "test_cluster_exec.db")
	dbConn, err := store.OpenAt(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	_ = dbConn.Migrate()
	_ = store.EnsureSSHHostsTable(dbConn.SQL())
	return dbPath, dbConn
}

func withMockClusterStoreDB(t *testing.T, fn func(db *store.DB)) {
	dbPath, testDB := setupClusterTestStoreDB(t)
	defer testDB.Close()
	prevOpener := openClusterDBFunc
	openClusterDBFunc = func() (*store.DB, error) { return store.OpenAt(dbPath) }
	defer func() { openClusterDBFunc = prevOpener }()
	fn(testDB)
}

func seedClusterTestHosts(t *testing.T, ctx context.Context, db *sql.DB) {
	h1 := store.SSHHost{ID: "h-c1", Alias: "master", IP: "192.168.1.10", Username: "root", ClusterRole: "control"}
	h2 := store.SSHHost{ID: "h-w1", Alias: "worker-1", IP: "192.168.1.11", Username: "root", ClusterRole: "worker"}
	h3 := store.SSHHost{ID: "h-w2", Alias: "worker-2", IP: "192.168.1.12", Username: "root", ClusterRole: "worker"}
	_ = store.UpsertSSHHost(ctx, h1, db)
	_ = store.UpsertSSHHost(ctx, h2, db)
	_ = store.UpsertSSHHost(ctx, h3, db)
}

func TestParseClusterExecArgs_Basic(t *testing.T) {
	args := []string{"all", "echo hi"}
	opts, err := parseClusterExecArgs(args)
	if err != nil || opts.target != "all" || opts.command != "echo hi" {
		t.Fatalf("unexpected parse result: %+v err: %v", opts, err)
	}
	if opts.isSudo || opts.parallel != 4 {
		t.Fatalf("expected defaults isSudo=false parallel=4, got %+v", opts)
	}
}

func TestParseClusterExecArgs_SudoShort(t *testing.T) {
	args := []string{"-s", "control", "uptime"}
	opts, err := parseClusterExecArgs(args)
	if err != nil || opts.target != "control" || opts.command != "uptime" {
		t.Fatalf("unexpected parse result: %+v err: %v", opts, err)
	}
	if !opts.isSudo {
		t.Fatalf("expected isSudo=true, got %+v", opts)
	}
}

func TestParseClusterExecArgs_SudoLong(t *testing.T) {
	args := []string{"workers", "df -h", "--sudo"}
	opts, err := parseClusterExecArgs(args)
	if err != nil || opts.target != "workers" || opts.command != "df -h" {
		t.Fatalf("unexpected parse result: %+v err: %v", opts, err)
	}
	if !opts.isSudo {
		t.Fatalf("expected isSudo=true, got %+v", opts)
	}
}

func TestParseClusterExecArgs_Parallel(t *testing.T) {
	args := []string{"-p", "8", "master", "whoami"}
	opts, err := parseClusterExecArgs(args)
	if err != nil || opts.parallel != 8 || opts.target != "master" {
		t.Fatalf("unexpected parse result: %+v err: %v", opts, err)
	}
}

func TestParseClusterExecArgs_ParallelEquals(t *testing.T) {
	args := []string{"--parallel=2", "-s", "all", "date"}
	opts, err := parseClusterExecArgs(args)
	if err != nil || opts.parallel != 2 || !opts.isSudo || opts.target != "all" {
		t.Fatalf("unexpected parse result: %+v err: %v", opts, err)
	}
}

func TestParseClusterExecArgs_MissingTarget(t *testing.T) {
	args := []string{}
	_, err := parseClusterExecArgs(args)
	if err == nil {
		t.Fatal("expected error for empty args, got nil")
	}
}

func TestParseClusterExecArgs_MissingCommand(t *testing.T) {
	args := []string{"all"}
	_, err := parseClusterExecArgs(args)
	if err == nil {
		t.Fatal("expected error for missing command, got nil")
	}
}

func TestParseClusterExecArgs_Help(t *testing.T) {
	args := []string{"--help"}
	opts, err := parseClusterExecArgs(args)
	if err != nil || !opts.isShowHelp {
		t.Fatalf("expected isShowHelp=true, got %+v err: %v", opts, err)
	}
}

func TestParseClusterScriptArgs_Basic(t *testing.T) {
	args := []string{"control", "setup.sh"}
	opts, err := parseClusterScriptArgs(args)
	if err != nil || opts.target != "control" || opts.scriptPath != "setup.sh" {
		t.Fatalf("unexpected parse result: %+v err: %v", opts, err)
	}
	if opts.isSudo || opts.parallel != 4 {
		t.Fatalf("expected defaults isSudo=false parallel=4, got %+v", opts)
	}
}

func TestParseClusterScriptArgs_Flags(t *testing.T) {
	args := []string{"-s", "-p", "10", "workers", "/tmp/deploy.sh"}
	opts, err := parseClusterScriptArgs(args)
	if err != nil || !opts.isSudo || opts.parallel != 10 || opts.scriptPath != "/tmp/deploy.sh" {
		t.Fatalf("unexpected parse result: %+v err: %v", opts, err)
	}
}

func TestParseClusterScriptArgs_MissingArgs(t *testing.T) {
	args := []string{"control"}
	_, err := parseClusterScriptArgs(args)
	if err == nil {
		t.Fatal("expected error for missing script path, got nil")
	}
}

func TestFormatClusterCommand_NonSudo(t *testing.T) {
	cmd := FormatClusterCommand("uptime", "secret", false)
	expected := "bash -c 'uptime'"
	if cmd != expected {
		t.Fatalf("expected %q, got %q", expected, cmd)
	}
}

func TestFormatClusterCommand_Sudo(t *testing.T) {
	cmd := FormatClusterCommand("systemctl restart k8s", "myPass123", true)
	expected := "echo 'myPass123' | sudo -S bash -c 'systemctl restart k8s'"
	if cmd != expected {
		t.Fatalf("expected %q, got %q", expected, cmd)
	}
}

func TestFormatClusterCommand_EscapedQuotes(t *testing.T) {
	cmd := FormatClusterCommand("echo 'hello'", "p'ss", true)
	expected := "echo 'p'\\''ss' | sudo -S bash -c 'echo '\\''hello'\\'''"
	if cmd != expected {
		t.Fatalf("expected %q, got %q", expected, cmd)
	}
}

func TestFormatClusterScriptCommand(t *testing.T) {
	nonSudo := FormatClusterScriptCommand("/tmp/test.sh", "", false)
	if nonSudo != "bash /tmp/test.sh" {
		t.Fatalf("expected bash /tmp/test.sh, got %q", nonSudo)
	}
	sudo := FormatClusterScriptCommand("/tmp/test.sh", "pass123", true)
	if sudo != "echo 'pass123' | sudo -S bash /tmp/test.sh" {
		t.Fatalf("expected sudo script command, got %q", sudo)
	}
}

func TestFormatClusterScriptCleanup(t *testing.T) {
	nonSudo := FormatClusterScriptCleanup("/tmp/test.sh", "", false)
	if nonSudo != "rm -f /tmp/test.sh" {
		t.Fatalf("expected rm -f /tmp/test.sh, got %q", nonSudo)
	}
	sudo := FormatClusterScriptCleanup("/tmp/test.sh", "pass123", true)
	if sudo != "echo 'pass123' | sudo -S rm -f /tmp/test.sh" {
		t.Fatalf("expected sudo cleanup command, got %q", sudo)
	}
}

func TestTargetGroupExpansion_All(t *testing.T) {
	db := setupExecClusterTestDB(t)
	defer db.Close()
	ctx := context.Background()
	seedClusterTestHosts(t, ctx, db)
	hosts, err := store.ListHostsByTarget(ctx, "all", db)
	if err != nil || len(hosts) != 3 {
		t.Fatalf("expected 3 hosts for target 'all', got %d err: %v", len(hosts), err)
	}
}

func TestTargetGroupExpansion_Roles(t *testing.T) {
	db := setupExecClusterTestDB(t)
	defer db.Close()
	ctx := context.Background()
	seedClusterTestHosts(t, ctx, db)
	ctrls, err := store.ListHostsByTarget(ctx, "control", db)
	if err != nil || len(ctrls) != 1 {
		t.Fatalf("expected 1 control host, got %d err: %v", len(ctrls), err)
	}
	workers, err := store.ListHostsByTarget(ctx, "workers", db)
	if err != nil || len(workers) != 2 {
		t.Fatalf("expected 2 worker hosts, got %d err: %v", len(workers), err)
	}
}

func TestTargetGroupExpansion_AliasAndComma(t *testing.T) {
	db := setupExecClusterTestDB(t)
	defer db.Close()
	ctx := context.Background()
	seedClusterTestHosts(t, ctx, db)
	single, err := store.ListHostsByTarget(ctx, "master", db)
	if err != nil || len(single) != 1 {
		t.Fatalf("expected 1 host for 'master', got %d err: %v", len(single), err)
	}
	pair, err := store.ListHostsByTarget(ctx, "master,worker-1", db)
	if err != nil || len(pair) != 2 {
		t.Fatalf("expected 2 hosts for 'master,worker-1', got %d err: %v", len(pair), err)
	}
}

func createSampleSummaryResults() []ClusterRunResult {
	return []ClusterRunResult{
		{Host: store.SSHHost{Alias: "master", IP: "10.0.0.1", ClusterRole: "control"}, ExitCode: 0, Duration: 100 * time.Millisecond},
		{Host: store.SSHHost{Alias: "worker-1", IP: "10.0.0.2", ClusterRole: "worker"}, ExitCode: 1, Duration: 250 * time.Millisecond},
	}
}

func TestRenderClusterSummaryTable(t *testing.T) {
	results := createSampleSummaryResults()
	table := RenderClusterSummaryTable(results)
	hasColumns := strings.Contains(table, "NODE") && strings.Contains(table, "ROLE") && strings.Contains(table, "STATUS")
	if !hasColumns {
		t.Fatalf("expected table header columns in %q", table)
	}
	hasMaster := strings.Contains(table, "master") && strings.Contains(table, "SUCCESS")
	hasWorker := strings.Contains(table, "worker-1") && strings.Contains(table, "FAILED")
	if !hasMaster || !hasWorker {
		t.Fatalf("expected row entries in table, got %q", table)
	}
}

func TestDispatchClusterRun_WorkerPool(t *testing.T) {
	hosts := []store.SSHHost{
		{ID: "1", Alias: "n1", IP: "10.0.0.1"},
		{ID: "2", Alias: "n2", IP: "10.0.0.2"},
	}
	prevFn := executeNodeFn
	executeNodeFn = func(ctx context.Context, h store.SSHHost, cmd string, isSudo bool) ClusterRunResult {
		return ClusterRunResult{Host: h, ExitCode: 0, Stdout: "ok\n", Duration: 10 * time.Millisecond}
	}
	defer func() { executeNodeFn = prevFn }()
	res := DispatchClusterRun(context.Background(), hosts, "echo ok", false, 2)
	if len(res) != 2 || res[0].ExitCode != 0 || res[1].ExitCode != 0 {
		t.Fatalf("expected 2 successful results, got %+v", res)
	}
}

func TestDispatchClusterScript_WorkerPool(t *testing.T) {
	hosts := []store.SSHHost{
		{ID: "1", Alias: "n1", IP: "10.0.0.1"},
	}
	prevFn := executeNodeScriptFn
	executeNodeScriptFn = func(ctx context.Context, h store.SSHHost, content []byte, isSudo bool) ClusterRunResult {
		return ClusterRunResult{Host: h, ExitCode: 0, Stdout: "script-ok\n", Duration: 5 * time.Millisecond}
	}
	defer func() { executeNodeScriptFn = prevFn }()
	res := DispatchClusterScript(context.Background(), hosts, []byte("echo 1"), false, 1)
	if len(res) != 1 || res[0].Stdout != "script-ok\n" {
		t.Fatalf("expected script-ok, got %+v", res)
	}
}

func TestRunClusterExecCLI_Flow(t *testing.T) {
	withMockClusterStoreDB(t, func(db *store.DB) {
		ctx := context.Background()
		seedClusterTestHosts(t, ctx, db.SQL())
		prevDispatch := dispatchClusterRunFn
		dispatchClusterRunFn = func(ctx context.Context, hosts []store.SSHHost, cmd string, isSudo bool, concurrency int) []ClusterRunResult {
			return []ClusterRunResult{{Host: hosts[0], ExitCode: 0}}
		}
		defer func() { dispatchClusterRunFn = prevDispatch }()
		err := RunClusterExecCLI([]string{"master", "hostname"})
		if err != nil {
			t.Fatalf("expected RunClusterExecCLI success, got %v", err)
		}
	})
}

func createTempTestScript(t *testing.T) string {
	tmpFile := filepath.Join(t.TempDir(), "run.sh")
	_ = os.WriteFile(tmpFile, []byte("echo test"), 0700)
	return tmpFile
}

func TestRunClusterScriptCLI_Flow(t *testing.T) {
	withMockClusterStoreDB(t, func(db *store.DB) {
		seedClusterTestHosts(t, context.Background(), db.SQL())
		tmpFile := createTempTestScript(t)
		prevDispatch := dispatchClusterScriptFn
		dispatchClusterScriptFn = func(ctx context.Context, hosts []store.SSHHost, content []byte, isSudo bool, concurrency int) []ClusterRunResult {
			return []ClusterRunResult{{Host: hosts[0], ExitCode: 0}}
		}
		defer func() { dispatchClusterScriptFn = prevDispatch }()
		if err := RunClusterScriptCLI([]string{"master", tmpFile}); err != nil {
			t.Fatalf("expected RunClusterScriptCLI success, got %v", err)
		}
	})
}
