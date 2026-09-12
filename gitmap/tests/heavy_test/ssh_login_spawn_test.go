package heavy_test

import (
	"context"
	"os/exec"
	"runtime"
	"testing"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd"
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

	oldExecutor := cmd.SSHExecutor
	cmd.SSHExecutor = fakeSSHCommand
	defer func() { cmd.SSHExecutor = oldExecutor }()

	err := cmd.RunSSHLogin(c, []string{"my-target"}, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = cmd.RunSSHLogin(c, []string{}, ctx)
	if err == nil {
		t.Errorf("Expected error for missing arguments")
	}
}
