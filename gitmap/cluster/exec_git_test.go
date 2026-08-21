package cluster

import (
	"context"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

type gitExecTestCase struct {
	name     string
	fn       func(context.Context, ClusterNode) (string, string, int, error)
	expected string
}

func gitExecTestCases() []gitExecTestCase {
	return []gitExecTestCase{
		{"Pull", ExecGitPull, constants.ClusterGitPullCmd},
		{"Push", ExecGitPush, constants.ClusterGitPushCmd},
		{"Commit", ExecGitCommit, constants.ClusterGitCommitCmd},
		{"Status", ExecGitStatus, constants.ClusterGitStatusCmd},
	}
}

func TestExecGit_Commands(t *testing.T) {
	origRunCmd := runCmdFunc
	defer func() { runCmdFunc = origRunCmd }()

	var lastCmd string
	runCmdFunc = func(cmd *exec.Cmd) error {
		lastCmd = strings.Join(cmd.Args, " ")
		return nil
	}

	runAllGitExecTests(t, &lastCmd)
}

func runAllGitExecTests(t *testing.T, lastCmd *string) {
	node := ClusterNode{}
	ctx := context.Background()
	for _, tc := range gitExecTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			runSingleGitExecTest(t, tc, lastCmd, node, ctx)
		})
	}
}

func runSingleGitExecTest(t *testing.T, tc gitExecTestCase, lastCmd *string, node ClusterNode, ctx context.Context) {
	if _, _, _, err := tc.fn(ctx, node); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(*lastCmd, tc.expected) {
		t.Errorf("expected command %q to contain %q", *lastCmd, tc.expected)
	}
	assertPlatformShellWrapping(t, *lastCmd)
}

func assertPlatformShellWrapping(t *testing.T, cmdStr string) {
	if runtime.GOOS == constants.PlatformWindows {
		assertWindowsShellWrapping(t, cmdStr)
		return
	}
	assertUnixShellWrapping(t, cmdStr)
}

func assertWindowsShellWrapping(t *testing.T, cmdStr string) {
	lower := strings.ToLower(cmdStr)
	if !strings.Contains(lower, constants.WindowsShell) || !strings.Contains(lower, strings.ToLower(constants.WindowsShellArg)) {
		t.Errorf("expected windows shell %s %s, got %q", constants.WindowsShell, constants.WindowsShellArg, cmdStr)
	}
}

func assertUnixShellWrapping(t *testing.T, cmdStr string) {
	if !strings.Contains(cmdStr, "sh") || !strings.Contains(cmdStr, constants.UnixShellArg) {
		t.Errorf("expected unix shell sh %s, got %q", constants.UnixShellArg, cmdStr)
	}
}
