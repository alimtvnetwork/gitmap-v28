package gitutil

import (
	"context"
	"errors"
	"os/exec"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// DefaultGitTimeout is the standard fail-safe timeout for git metadata commands.
const DefaultGitTimeout = 15 * time.Second

// ExecGitWithTimeout executes a git command in dir with a bounded context timeout.
func ExecGitWithTimeout(timeout time.Duration, dir string, args ...string) ([]byte, error) {
	if timeout <= 0 {
		timeout = DefaultGitTimeout
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return ExecGitWithTimeoutContext(ctx, dir, args...)
}

// ExecGitWithTimeoutContext executes a git command respecting a caller-supplied context.
func ExecGitWithTimeoutContext(ctx context.Context, dir string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, constants.GitBin, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Env = BuildSafeGitEnv(constants.EnvGitSSHCommandBatchYes)

	out, err := cmd.CombinedOutput()
	if err == nil {
		return out, nil
	}

	return handleGitTimeoutError(ctx, err)
}

func handleGitTimeoutError(ctx context.Context, err error) ([]byte, error) {
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return nil, apperror.NewSimple("git command timed out after deadline", "E_GIT_TIMEOUT")
	}
	if ctx.Err() != nil {
		return nil, apperror.WrapSimple(ctx.Err(), "git command cancelled")
	}

	return nil, err
}

// ExecGitCheck executes a git command with timeout and returns true if it exits with code 0.
func ExecGitCheck(timeout time.Duration, dir string, args ...string) bool {
	_, err := ExecGitWithTimeout(timeout, dir, args...)
	return err == nil
}
