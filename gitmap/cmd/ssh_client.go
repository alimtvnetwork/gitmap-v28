package cmd

import (
	"context"
	"errors"
	"io"
	"os"
	"os/exec"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

type InteractiveSSHClient struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (c *InteractiveSSHClient) Run(ctx context.Context, target string) error {
	cmd := exec.CommandContext(ctx, "ssh", target)
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

	cmd := exec.CommandContext(ctx, "ssh", cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return &apperror.AppError{
				Op:    "SpawnSSH",
				Code:  "E_INTERNAL_ERROR",
				Cause: exitErr,
				Ctx:   map[string]any{"target": target.String(), "args": args},
			}
		}
		return &apperror.AppError{
			Op:    "SpawnSSH",
			Code:  "E_INTERNAL_ERROR",
			Cause: err,
			Ctx:   map[string]any{"target": target.String(), "args": args},
		}
	}

	return nil
}
