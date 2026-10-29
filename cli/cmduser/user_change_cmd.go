package cmduser

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func applyUserChangeIdentity(opts UserChangeOptions) error {
	if err := applyUserChangeName(opts); err != nil {
		return err
	}

	return applyUserChangeEmail(opts)
}

func applyUserChangeName(opts UserChangeOptions) error {
	if len(opts.Name) == 0 {
		return nil
	}

	return writeGitConfigKey("user.name", opts.Name, opts.IsGlobal)
}

func applyUserChangeEmail(opts UserChangeOptions) error {
	if opts.HasEmptyEmail || len(opts.Email) == 0 {
		return unsetGitConfigKey("user.email", opts.IsGlobal)
	}

	return writeGitConfigKey("user.email", opts.Email, opts.IsGlobal)
}

func writeGitConfigKey(key, value string, isGlobal bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	args := buildConfigSetArgs(key, value, isGlobal)
	out, err := exec.CommandContext(ctx, "git", args...).CombinedOutput()
	if err != nil {
		return apperror.WrapSimple(err, "git config set failed: "+strings.TrimSpace(string(out)))
	}

	return nil
}

func buildConfigSetArgs(key, value string, isGlobal bool) []string {
	if isGlobal {
		return []string{"config", "--global", key, value}
	}

	return []string{"config", key, value}
}

func unsetGitConfigKey(key string, isGlobal bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	args := buildConfigUnsetArgs(key, isGlobal)
	_ = exec.CommandContext(ctx, "git", args...).Run()

	return nil
}

func buildConfigUnsetArgs(key string, isGlobal bool) []string {
	if isGlobal {
		return []string{"config", "--global", "--unset", key}
	}

	return []string{"config", "--unset", key}
}

func renderChangeSuccess(opts UserChangeOptions) {
	scopeName := resolveScopeDisplayName(opts)
	if opts.HasEmptyEmail || len(opts.Email) == 0 {
		fmt.Printf("✔ Changed %s Git author identity to: %s (email unset)\n", scopeName, opts.Name)
	} else {
		fmt.Printf("✔ Changed %s Git author identity to: %s <%s>\n", scopeName, opts.Name, opts.Email)
	}

	if len(opts.Role) > 0 {
		fmt.Printf("  Role: %s\n", opts.Role)
	}
}

func resolveScopeDisplayName(opts UserChangeOptions) string {
	if opts.IsGlobal {
		return "Global"
	}

	if opts.IsProject {
		return "Project"
	}

	return "Local"
}

// RunChangeUser changes the Git user identity based on arguments.
func RunChangeUser(args []string) error {
	opts, err := ParseChangeUserArgs(args)
	if err != nil {
		return logUserChangeError("Error parsing user change options: %v\n", err)
	}

	if err := applyUserChangeIdentity(opts); err != nil {
		return logUserChangeError("Error updating Git author identity: %v\n", err)
	}

	renderChangeSuccess(opts)

	return nil
}

// RunUserChange is an alias for RunChangeUser.
func RunUserChange(args []string) error {
	return RunChangeUser(args)
}

// RunChangeLogin is an alias for RunChangeUser.
func RunChangeLogin(args []string) error {
	return RunChangeUser(args)
}

// RunUserAutobot switches to the autobot identity.
func RunUserAutobot(args []string) error {
	opts, err := ParseAutobotArgs(args)
	if err != nil {
		return logUserChangeError("Error parsing autobot options: %v\n", err)
	}

	if err := applyUserChangeIdentity(opts); err != nil {
		return logUserChangeError("Error applying autobot identity: %v\n", err)
	}

	renderChangeSuccess(opts)

	return nil
}

// RunAutobot is an alias for RunUserAutobot.
func RunAutobot(args []string) error {
	return RunUserAutobot(args)
}

func logUserChangeError(format string, err error) error {
	fmt.Fprintf(os.Stderr, format, err)

	return err
}
