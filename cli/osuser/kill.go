package osuser

import (
	"errors"
	"fmt"
	"os/exec"
	"os/user"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// KillUserProcesses terminates all active processes belonging to the target user.
func KillUserProcesses(opts UserKillOptions) *apperror.AppError {
	appErr := validateKillTarget(opts)
	if appErr != nil {
		return appErr
	}
	if opts.IsDryRun {
		return nil
	}
	return dispatchPlatformKill(opts)
}

func validateKillTarget(opts UserKillOptions) *apperror.AppError {
	hasUser := len(strings.TrimSpace(opts.Username)) > 0
	if !hasUser {
		return apperror.NewValidationError("target username is required for kill")
	}
	return checkProtectionGuards(opts)
}

func checkProtectionGuards(opts UserKillOptions) *apperror.AppError {
	isProtected := isProtectedUser(opts.Username)
	isSelf := isCurrentProcessUser(opts.Username)
	isGuarded := (isProtected || isSelf) && !opts.IsForce
	if isGuarded {
		return apperror.NewValidationError(fmt.Sprintf("cannot terminate protected or current user %q without explicit force", opts.Username))
	}
	return nil
}

func isProtectedUser(username string) bool {
	lower := strings.ToLower(strings.TrimSpace(username))
	isRoot := lower == "root" || lower == "0"
	isAdmin := lower == "administrator" || lower == "system"
	return isRoot || isAdmin
}

func isCurrentProcessUser(username string) bool {
	curr, err := user.Current()
	if err != nil {
		return false
	}
	target := strings.ToLower(strings.TrimSpace(username))
	isUsernameMatch := strings.ToLower(curr.Username) == target
	isUidMatch := curr.Uid == target
	return isUsernameMatch || isUidMatch
}

func dispatchPlatformKill(opts UserKillOptions) *apperror.AppError {
	switch runtime.GOOS {
	case "windows":
		return killWindowsUserProcesses(opts.Username)
	case "linux":
		return killLinuxUserProcesses(opts.Username, opts.IsForce)
	default:
		return apperror.NewExecutionError(fmt.Sprintf("unsupported os for process termination: %s", runtime.GOOS))
	}
}

func killWindowsUserProcesses(username string) *apperror.AppError {
	filter := fmt.Sprintf("USERNAME eq %s", username)
	cmd := exec.Command("taskkill", "/F", "/FI", filter)
	out, err := cmd.CombinedOutput()
	if err != nil && !isWindowsTaskkillHarmless(string(out)) {
		return apperror.Wrap(err, "osuser.killWindowsUserProcesses", map[string]any{"output": string(out)})
	}
	return nil
}

func isWindowsTaskkillHarmless(output string) bool {
	lower := strings.ToLower(output)
	hasIdleNotice := strings.Contains(lower, "no tasks running")
	hasMissingNotice := strings.Contains(lower, "not found")
	return hasIdleNotice || hasMissingNotice
}

func killLinuxUserProcesses(username string, isForce bool) *apperror.AppError {
	args := buildPkillArgs(username, isForce)
	cmd := exec.Command("pkill", args...)
	out, err := cmd.CombinedOutput()
	if err != nil && !isLinuxPkillHarmless(err) {
		return apperror.Wrap(err, "osuser.killLinuxUserProcesses", map[string]any{"output": string(out)})
	}
	return nil
}

func buildPkillArgs(username string, isForce bool) []string {
	if isForce {
		return []string{"-9", "-u", username}
	}
	return []string{"-u", username}
}

func isLinuxPkillHarmless(err error) bool {
	var exitErr *exec.ExitError
	isExit := errors.As(err, &exitErr)
	if isExit && exitErr.ExitCode() == 1 {
		return true
	}
	return false
}
