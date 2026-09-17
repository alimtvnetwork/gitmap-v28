package cmdssh

import (
	"context"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestBuildGitmapInstallOneLiner_LinuxDefault(t *testing.T) {
	cmd := BuildGitmapInstallOneLiner("linux", "latest")
	expected := "curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.sh | bash"
	if cmd != expected {
		t.Fatalf("expected %q, got %q", expected, cmd)
	}
}

func TestBuildGitmapInstallOneLiner_LinuxEmpty(t *testing.T) {
	cmd := BuildGitmapInstallOneLiner("linux", "")
	expected := "curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.sh | bash"
	if cmd != expected {
		t.Fatalf("expected %q, got %q", expected, cmd)
	}
}

func TestBuildGitmapInstallOneLiner_LinuxPinned(t *testing.T) {
	cmd := BuildGitmapInstallOneLiner("linux", "v2.55.0")
	expected := "curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.sh | bash -s -- --version v2.55.0"
	if cmd != expected {
		t.Fatalf("expected %q, got %q", expected, cmd)
	}
}

func TestBuildGitmapInstallOneLiner_DarwinPinned(t *testing.T) {
	cmd := BuildGitmapInstallOneLiner("darwin", "v3.0.0")
	expected := "curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.sh | bash -s -- --version v3.0.0"
	if cmd != expected {
		t.Fatalf("expected %q, got %q", expected, cmd)
	}
}

func TestBuildGitmapInstallOneLiner_WindowsDefault(t *testing.T) {
	cmd := BuildGitmapInstallOneLiner("windows", "latest")
	expected := `powershell -NoProfile -ExecutionPolicy Bypass -Command "irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.ps1 | iex"`
	if cmd != expected {
		t.Fatalf("expected %q, got %q", expected, cmd)
	}
}

func TestBuildGitmapInstallOneLiner_WindowsPinned(t *testing.T) {
	cmd := BuildGitmapInstallOneLiner("windows", "v2.48.0")
	expected := `powershell -NoProfile -ExecutionPolicy Bypass -Command "irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.ps1 | iex" -Version v2.48.0`
	if cmd != expected {
		t.Fatalf("expected %q, got %q", expected, cmd)
	}
}

func TestParseClusterInstallArgs_Defaults(t *testing.T) {
	opts := parseClusterInstallArgs([]string{"gitmap"})
	if opts.component != "gitmap" || opts.target != "all" {
		t.Fatalf("expected gitmap all, got %s %s", opts.component, opts.target)
	}
	if opts.version != "latest" || !opts.isSudo || opts.parallel != 4 {
		t.Fatalf("expected default flags, got %+v", opts)
	}
}

func TestParseClusterInstallArgs_CustomTargetAndFlags(t *testing.T) {
	args := []string{"gitmap", "workers", "--version", "v1.2.3", "--parallel", "8", "--sudo=false", "--os", "windows"}
	opts := parseClusterInstallArgs(args)
	if opts.component != "gitmap" || opts.target != "workers" {
		t.Fatalf("expected gitmap workers, got %s %s", opts.component, opts.target)
	}
	if opts.version != "v1.2.3" || opts.isSudo || opts.parallel != 8 || opts.osType != "windows" {
		t.Fatalf("unexpected options: %+v", opts)
	}
}

func TestParseClusterInstallArgs_ShortFlags(t *testing.T) {
	args := []string{"control", "-v", "v2.0.0", "-p", "2", "-s"}
	opts := parseClusterInstallArgs(args)
	if opts.target != "control" || opts.version != "v2.0.0" || opts.parallel != 2 || !opts.isSudo {
		t.Fatalf("unexpected options: %+v", opts)
	}
}

func TestRunClusterInstallCLI_Help(t *testing.T) {
	if err := RunClusterInstallCLI([]string{}); err != nil {
		t.Fatalf("expected nil for empty args help, got: %v", err)
	}
	if err := RunClusterInstallCLI([]string{"--help"}); err != nil {
		t.Fatalf("expected nil for --help, got: %v", err)
	}
	if err := RunClusterInstallCLI([]string{"-h"}); err != nil {
		t.Fatalf("expected nil for -h, got: %v", err)
	}
}

func mockSuccessDispatch(t *testing.T) func(context.Context, []store.SSHHost, string, bool, int) []ClusterRunResult {
	return func(ctx context.Context, hosts []store.SSHHost, cmd string, isSudo bool, concurrency int) []ClusterRunResult {
		if !strings.Contains(cmd, "curl") || !isSudo {
			t.Fatalf("unexpected dispatch: cmd=%s isSudo=%v", cmd, isSudo)
		}
		return []ClusterRunResult{{Host: hosts[0], ExitCode: 0}}
	}
}

func TestRunClusterInstallCLI_Success(t *testing.T) {
	withMockClusterStoreDB(t, func(db *store.DB) {
		seedClusterTestHosts(t, context.Background(), db.SQL())
		prevDispatch := dispatchClusterRunFn
		dispatchClusterRunFn = mockSuccessDispatch(t)
		defer func() { dispatchClusterRunFn = prevDispatch }()
		if err := RunClusterInstallCLI([]string{"gitmap", "master"}); err != nil {
			t.Fatalf("expected RunClusterInstallCLI success, got: %v", err)
		}
	})
}

func TestRunClusterInstallCLI_Failure(t *testing.T) {
	withMockClusterStoreDB(t, func(db *store.DB) {
		seedClusterTestHosts(t, context.Background(), db.SQL())
		prevDispatch := dispatchClusterRunFn
		dispatchClusterRunFn = func(ctx context.Context, hosts []store.SSHHost, cmd string, isSudo bool, concurrency int) []ClusterRunResult {
			return []ClusterRunResult{{Host: hosts[0], ExitCode: 1}}
		}
		defer func() { dispatchClusterRunFn = prevDispatch }()
		err := RunClusterInstallCLI([]string{"gitmap", "master"})
		if err == nil {
			t.Fatal("expected execution error on failure, got nil")
		}
	})
}
