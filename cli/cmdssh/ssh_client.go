package cmdssh

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"

	"golang.org/x/term"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// SSHExecutor is the command factory used for executing SSH commands.
var SSHExecutor = exec.CommandContext

type InteractiveSSHClient struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func executeClientCmd(cmd *exec.Cmd, op string, ctx map[string]any) error {
	if err := cmd.Run(); err != nil {
		return &apperror.AppError{
			Op:    op,
			Code:  "E_INTERNAL_ERROR",
			Cause: err,
			Ctx:   ctx,
		}
	}

	return nil
}

func (c *InteractiveSSHClient) Run(ctx context.Context, target string) error {
	cmd := SSHExecutor(ctx, "ssh", target)
	cmd.Stdin = c.Stdin
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr

	return executeClientCmd(cmd, "InteractiveSSHClient.Run", map[string]any{"target": target})
}

func isCustomSSHPort(port int) bool {
	return port > 0 && port != 22
}

func buildSSHArgs(target SSHTarget, args []string) []string {
	var cmdArgs []string
	if isCustomSSHPort(target.Port) {
		cmdArgs = append(cmdArgs, "-p", strconv.Itoa(target.Port))
	}

	cmdArgs = append(cmdArgs, target.String())
	cmdArgs = append(cmdArgs, args...)

	return cmdArgs
}

func SpawnSSH(ctx context.Context, target SSHTarget, args []string) error {
	cmdArgs := buildSSHArgs(target, args)
	cmd := SSHExecutor(ctx, "ssh", cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return executeClientCmd(cmd, "SpawnSSH", map[string]any{"target": target.String(), "args": args})
}

func PromptSSHPassword(ctx context.Context, prompt string, fd int) (string, error) {
	fmt.Print(prompt)
	password, err := term.ReadPassword(fd)
	fmt.Println()
	if err != nil {
		return "", &apperror.AppError{
			Op:    "PromptSSHPassword",
			Code:  "E_INTERNAL_ERROR",
			Cause: err,
		}
	}

	return string(password), nil
}
