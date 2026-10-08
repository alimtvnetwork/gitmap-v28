package heavy_test

import (
	"context"
	"os/exec"
	"runtime"
	"testing"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
)

func fakeSSHCommand(ctx context.Context, name string, args ...string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.CommandContext(ctx, "cmd", "/c", "exit 0")
	}

	return exec.CommandContext(ctx, "true")
}

func TestRunSSHLogin(t *testing.T) {
	c := &cobra.Command{}
	ctx := context.Background()

	oldExecutor := cmdssh.SSHExecutor
	cmdssh.SSHExecutor = fakeSSHCommand
	defer func() { cmdssh.SSHExecutor = oldExecutor }()

	err := cmdssh.RunSSHLogin(c, []string{"127.0.0.1"}, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = cmdssh.RunSSHLogin(c, []string{}, ctx)
	if err == nil {
		t.Errorf("Expected error for missing arguments")
	}
}
