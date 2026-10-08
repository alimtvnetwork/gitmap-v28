package usercontext

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var (
	execGitCommand = defaultGitCommand
)

const gitCommandTimeout = 2 * time.Second

func defaultGitCommand(ctx context.Context, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "git", args...)

	return cmd.CombinedOutput()
}

// GetGitIdentity queries Git configuration for user.name and user.email.
func GetGitIdentity(isGlobal bool) (GitIdentity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gitCommandTimeout)
	defer cancel()

	name := readConfigKey(ctx, "user.name", isGlobal)
	email := readConfigKey(ctx, "user.email", isGlobal)
	isConfigured := len(name) > 0 || len(email) > 0

	return GitIdentity{
		Name:         name,
		Email:        email,
		IsConfigured: isConfigured,
	}, nil
}

func readConfigKey(ctx context.Context, key string, isGlobal bool) string {
	args := buildConfigReadArgs(key, isGlobal)
	out, err := execGitCommand(ctx, args...)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func buildConfigReadArgs(key string, isGlobal bool) []string {
	if isGlobal {
		return []string{"config", "--global", key}
	}

	return []string{"config", key}
}

// SetGitIdentity writes user.name and user.email to Git configuration.
func SetGitIdentity(name, email string, isGlobal bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), gitCommandTimeout)
	defer cancel()

	if err := writeConfigKey(ctx, "user.name", name, isGlobal); err != nil {
		return err
	}

	return writeConfigKey(ctx, "user.email", email, isGlobal)
}

func writeConfigKey(ctx context.Context, key, value string, isGlobal bool) error {
	if len(value) == 0 {
		return nil
	}

	args := buildConfigWriteArgs(key, value, isGlobal)
	out, err := execGitCommand(ctx, args...)
	if err != nil {
		return apperror.WrapSimple(err, "git config set failed: "+strings.TrimSpace(string(out)))
	}

	return nil
}

func buildConfigWriteArgs(key, value string, isGlobal bool) []string {
	if isGlobal {
		return []string{"config", "--global", key, value}
	}

	return []string{"config", key, value}
}
