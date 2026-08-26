package cmd

import (
	"context"
	"os/exec"
	"runtime"
	"testing"

	"github.com/spf13/cobra"
)

func fakeSSHCommand(ctx context.Context, name string, args ...string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.CommandContext(ctx, "cmd", "/c", "exit 0")
	}
	return exec.CommandContext(ctx, "true")
}

func TestRunSSHLogin(t *testing.T) {
	cmd := &cobra.Command{}
	ctx := context.Background()

	oldExecutor := sshExecutor
	sshExecutor = fakeSSHCommand
	defer func() { sshExecutor = oldExecutor }()

	err := runSSHLogin(cmd, []string{"my-target"}, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = runSSHLogin(cmd, []string{}, ctx)
	if err == nil {
		t.Errorf("Expected error for missing arguments")
	}
}
