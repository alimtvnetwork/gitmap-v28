package cmdssh

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestRequiresGitmap_Positive(t *testing.T) {
	cases := []string{
		"gitmap", "gitmap status", "gitmap.exe pull",
		"/usr/local/bin/gitmap", "/bin/gitmap status", "sudo gitmap clone",
	}
	for _, cmd := range cases {
		if !RequiresGitmap(cmd) {
			t.Errorf("expected true for %q", cmd)
		}
	}
}

func TestRequiresGitmap_Negative(t *testing.T) {
	cases := []string{"git status", "uptime", "echo gitmap-tool"}
	for _, cmd := range cases {
		if RequiresGitmap(cmd) {
			t.Errorf("expected false for %q", cmd)
		}
	}
}

func TestResolveDetectedOSType(t *testing.T) {
	if got := resolveDetectedOSType("Linux node1"); got != "linux" {
		t.Fatalf("expected linux, got %s", got)
	}
	if got := resolveDetectedOSType("Darwin macbook"); got != "darwin" {
		t.Fatalf("expected darwin, got %s", got)
	}
	if got := resolveDetectedOSType("Windows_NT"); got != "windows" {
		t.Fatalf("expected windows, got %s", got)
	}
	if got := resolveDetectedOSType("unknown"); got != "linux" {
		t.Fatalf("expected default linux, got %s", got)
	}
}

func mockRunNodeSSHProcess(ctx context.Context, host store.SSHHost, alias, cmd, pwd string, isSudo bool, start time.Time) ClusterRunResult {
	return ClusterRunResult{Host: host, ExitCode: 0, Stdout: "executed: " + cmd}
}

func TestExecuteNodeCommand_ProbeSkippedWhenNotGitmap(t *testing.T) {
	probeInvoked := false
	prevProbe := runNodeProbeFn
	runNodeProbeFn = func(ctx context.Context, host store.SSHHost, probeCmd, pwd string) int {
		probeInvoked = true
		return 0
	}
	defer func() { runNodeProbeFn = prevProbe }()

	host := store.SSHHost{ID: "h1", Alias: "node1", IP: "10.0.0.1"}
	_ = runNodeSSHProcess(context.Background(), host, "node1", "uptime", "", false, time.Now())
	if probeInvoked {
		t.Fatal("expected probe to not be invoked for non-gitmap command")
	}
}

func mockInstalledProbe(called *bool) func() {
	prevProbe := runNodeProbeFn
	prevBoot := executeRemoteBootstrapFn
	runNodeProbeFn = func(ctx context.Context, h store.SSHHost, pCmd, pwd string) int { return 0 }
	executeRemoteBootstrapFn = func(ctx context.Context, h store.SSHHost, osType, pwd string, isSudo bool) int {
		*called = true
		return 0
	}
	return func() {
		runNodeProbeFn = prevProbe
		executeRemoteBootstrapFn = prevBoot
	}
}

func TestEnsureNodeGitmap_AlreadyInstalled(t *testing.T) {
	var bootstrapCalled bool
	reset := mockInstalledProbe(&bootstrapCalled)
	defer reset()

	host := store.SSHHost{ID: "h1", Alias: "node1", IP: "10.0.0.1"}
	ensureNodeGitmap(context.Background(), host, "node1", "", false)
	if bootstrapCalled {
		t.Fatal("expected bootstrap to not be called when gitmap is already installed")
	}
}

func mockMissingBootstrap(targetOS string, called *bool) func() {
	prevProbe := runNodeProbeFn
	prevOS := runNodeOSProbeFn
	prevBoot := executeRemoteBootstrapFn
	runNodeProbeFn = func(ctx context.Context, h store.SSHHost, pCmd, pwd string) int { return 1 }
	runNodeOSProbeFn = func(ctx context.Context, h store.SSHHost, pwd string) string { return targetOS }
	executeRemoteBootstrapFn = func(ctx context.Context, h store.SSHHost, osType, pwd string, isSudo bool) int {
		*called = (osType == targetOS)
		return 0
	}
	return func() {
		runNodeProbeFn = prevProbe
		runNodeOSProbeFn = prevOS
		executeRemoteBootstrapFn = prevBoot
	}
}

func TestEnsureNodeGitmap_AutoBootstrapMissing(t *testing.T) {
	var bootstrapCalled bool
	reset := mockMissingBootstrap("linux", &bootstrapCalled)
	defer reset()

	host := store.SSHHost{ID: "h1", Alias: "node1", IP: "10.0.0.1"}
	ensureNodeGitmap(context.Background(), host, "node1", "", true)
	if !bootstrapCalled {
		t.Fatal("expected bootstrap to be called when gitmap is missing")
	}
}

func TestEnsureNodeGitmap_WindowsBootstrap(t *testing.T) {
	var bootstrapCalled bool
	reset := mockMissingBootstrap("windows", &bootstrapCalled)
	defer reset()

	host := store.SSHHost{ID: "h1", Alias: "win-node", IP: "10.0.0.2"}
	ensureNodeGitmap(context.Background(), host, "win-node", "", false)
	if !bootstrapCalled {
		t.Fatal("expected windows bootstrap to be called")
	}
}

func TestBuildNodeSSHArgs(t *testing.T) {
	host := store.SSHHost{IP: "192.168.1.50", Username: "dev", Port: 2222}
	args := BuildNodeSSHArgs(host, "whoami")
	argsStr := strings.Join(args, " ")
	if !strings.Contains(argsStr, "-p 2222") || !strings.Contains(argsStr, "dev@192.168.1.50") {
		t.Fatalf("unexpected args: %s", argsStr)
	}
}
