package usercontext

import (
	"context"
	"errors"
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

	return newGitIdentity(name, email), nil
}

func newGitIdentity(name, email string) GitIdentity {
	hasName := len(name) > 0
	hasEmail := len(email) > 0

	return GitIdentity{
		Name:         name,
		Email:        email,
		IsConfigured: hasName && hasEmail,
		HasName:      hasName,
		HasEmail:     hasEmail,
	}
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

// UnsetGitConfigKey removes a Git config key, ignoring code 5 if missing.
func UnsetGitConfigKey(ctx context.Context, key string, isGlobal bool) error {
	args := buildConfigUnsetArgs(key, isGlobal)
	out, err := execGitCommand(ctx, args...)
	if err == nil || isExitStatusFive(err) {
		return nil
	}

	return apperror.WrapSimple(err, "git config unset failed: "+strings.TrimSpace(string(out)))
}

func isExitStatusFive(err error) bool {
	if err == nil {
		return false
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 5 {
		return true
	}

	return strings.Contains(err.Error(), "exit status 5")
}

func buildConfigUnsetArgs(key string, isGlobal bool) []string {
	if isGlobal {
		return []string{"config", "--global", "--unset", key}
	}

	return []string{"config", "--unset", key}
}

// SetGitIdentity writes user.name and user.email to Git configuration.
func SetGitIdentity(name, email string, isGlobal bool) error {
	hasEmpty := len(email) == 0
	opts := GitIdentityOptions{
		Name:          name,
		Email:         email,
		IsGlobal:      isGlobal,
		IsUnsetEmail:  hasEmpty,
		HasEmptyEmail: hasEmpty,
	}

	return SetGitIdentityWithOptions(opts)
}

// SetGitIdentityWithOptions writes or unsets Git identity based on options.
func SetGitIdentityWithOptions(opts GitIdentityOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), gitCommandTimeout)
	defer cancel()

	if err := applyIdentityName(ctx, opts); err != nil {
		return err
	}

	return applyIdentityEmail(ctx, opts)
}

func applyIdentityName(ctx context.Context, opts GitIdentityOptions) error {
	if len(opts.Name) == 0 {
		return nil
	}

	return writeConfigKey(ctx, "user.name", opts.Name, opts.IsGlobal)
}

func applyIdentityEmail(ctx context.Context, opts GitIdentityOptions) error {
	if opts.IsUnsetEmail {
		return UnsetGitConfigKey(ctx, "user.email", opts.IsGlobal)
	}

	if opts.HasEmptyEmail || len(opts.Email) > 0 {
		return writeConfigKey(ctx, "user.email", opts.Email, opts.IsGlobal)
	}

	return nil
}

func writeConfigKey(ctx context.Context, key, value string, isGlobal bool) error {
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
