package cluster

import (
	"context"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func TestExecGit_Commands(t *testing.T) {
	origRunCmd := runCmdFunc
	defer func() { runCmdFunc = origRunCmd }()

	var lastCmd string
	runCmdFunc = func(cmd *exec.Cmd) error {
		lastCmd = strings.Join(cmd.Args, " ")
		return nil
	}

	node := ClusterNode{}
	ctx := context.Background()

	tests := []struct {
		name     string
		fn       func(context.Context, ClusterNode) (string, string, int, error)
		expected string
	}{
		{"Pull", ExecGitPull, "gitmap pull --all"},
		{"Push", ExecGitPush, "gitmap push --all"},
		{"Commit", ExecGitCommit, "gitmap commit --all"},
		{"Status", ExecGitStatus, "gitmap status --all"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, _, _, err := tc.fn(ctx, node)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(lastCmd, tc.expected) {
				t.Errorf("expected command %q to contain %q", lastCmd, tc.expected)
			}

			// Validate Windows vs Unix shell wrapping
			if runtime.GOOS == "windows" && (!strings.Contains(strings.ToLower(lastCmd), "cmd.exe") || !strings.Contains(strings.ToLower(lastCmd), "/c")) {
				t.Errorf("expected windows shell cmd.exe /c, got %q", lastCmd)
			}
			if runtime.GOOS != "windows" && (!strings.Contains(lastCmd, "sh") || !strings.Contains(lastCmd, "-c")) {
				t.Errorf("expected unix shell sh -c, got %q", lastCmd)
			}
		})
	}
}
