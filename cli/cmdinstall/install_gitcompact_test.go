package cmdinstall

import (
	"errors"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestGitCompactScripts(t *testing.T) {
	if !strings.Contains(gitCompactWinScript, "git-compact/main/install.ps1") {
		t.Errorf("expected win script to contain install.ps1, got %s", gitCompactWinScript)
	}

	if !strings.Contains(gitCompactUnixScript, "git-compact/main/install.sh") {
		t.Errorf("expected unix script to contain install.sh, got %s", gitCompactUnixScript)
	}
}

func TestBuildGitCompactCommand(t *testing.T) {
	cmd := buildGitCompactCommand()
	if cmd == nil {
		t.Fatal("expected non-nil exec.Cmd")
	}

	if len(cmd.Args) == 0 {
		t.Fatal("expected non-empty args in exec.Cmd")
	}
}

func TestRunInstallGitCompactDryRun(t *testing.T) {
	opts := installOptions{DryRun: true}
	if err := runInstallGitCompact(opts); err != nil {
		t.Errorf("expected dry run to succeed without error, got %v", err)
	}
}

func TestGitCompactInToolProbeMap(t *testing.T) {
	cfg, isFound := toolProbeMap[constants.ToolGitCompact]
	if !isFound {
		t.Fatal("expected ToolGitCompact to exist in toolProbeMap")
	}

	if len(cfg.bins) == 0 || cfg.bins[0] != "git-compact" {
		t.Errorf("expected first candidate binary to be 'git-compact', got %v", cfg.bins)
	}
}

func TestGitCompactAliasResolution(t *testing.T) {
	if res := resolveToolAlias("git-compact"); res != constants.ToolGitCompact {
		t.Errorf("expected git-compact to resolve to canonical name, got %s", res)
	}

	if res := resolveToolAlias("gitcompact"); res != constants.ToolGitCompact {
		t.Errorf("expected gitcompact to resolve to canonical name, got %s", res)
	}
}

func TestResolveGitCompactErrorDetails(t *testing.T) {
	code0, msg0 := resolveGitCompactErrorDetails(nil)
	if code0 != 0 || msg0 != "" {
		t.Errorf("expected (0, '') for nil error, got (%d, %q)", code0, msg0)
	}

	testErr := errors.New("simulated error")
	code1, msg1 := resolveGitCompactErrorDetails(testErr)
	if code1 != 1 || msg1 != "simulated error" {
		t.Errorf("expected (1, 'simulated error'), got (%d, %q)", code1, msg1)
	}
}
