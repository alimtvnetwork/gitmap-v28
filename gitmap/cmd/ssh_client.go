package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"golang.org/x/term"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

var sshExecutor = exec.CommandContext

type InteractiveSSHClient struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (c *InteractiveSSHClient) Run(ctx context.Context, target string) error {
	cmd := sshExecutor(ctx, "ssh", target)
	cmd.Stdin = c.Stdin
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr

	if err := cmd.Run(); err != nil {
		return &apperror.AppError{
			Op:    "InteractiveSSHClient.Run",
			Code:  "E_INTERNAL_ERROR",
			Cause: err,
			Ctx:   map[string]any{"target": target},
		}
	}
	return nil
}

func SpawnSSH(ctx context.Context, target SSHTarget, args []string) error {
	cmdArgs := []string{target.String()}
	cmdArgs = append(cmdArgs, args...)

	cmd := sshExecutor(ctx, "ssh", cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return &apperror.AppError{
			Op:    "SpawnSSH",
			Code:  "E_INTERNAL_ERROR",
			Cause: err,
			Ctx:   map[string]any{"target": target.String(), "args": args},
		}
	}

	return nil
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
