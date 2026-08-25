package cmd

import (
	"context"
	"io"
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
