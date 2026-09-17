package osuser

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// RemoveEnhancedUser removes an OS user, terminates processes, and cleans sudoers entries.
func RemoveEnhancedUser(opts UserRemoveOptions) *apperror.AppError {
	appErr := validateRemoveOptions(opts)
	if appErr != nil {
		return appErr
	}
	return executeUserRemoval(opts)
}

func validateRemoveOptions(opts UserRemoveOptions) *apperror.AppError {
	hasUser := len(strings.TrimSpace(opts.Username)) > 0
	if !hasUser {
		return apperror.NewValidationError("username is required for removal")
	}
	return nil
}

func executeUserRemoval(opts UserRemoveOptions) *apperror.AppError {
	appErr := terminateUserProcessesIfRequested(opts)
	if appErr != nil {
		return appErr
	}
	return dispatchPlatformRemove(opts)
}

func terminateUserProcessesIfRequested(opts UserRemoveOptions) *apperror.AppError {
	if !opts.IsForceKill {
		return nil
	}
	killOpts := UserKillOptions{
		Username:  opts.Username,
		IsForce:   true,
		IsDryRun:  opts.IsDryRun,
		IsVerbose: opts.IsVerbose,
	}
	return KillUserProcesses(killOpts)
}

func dispatchPlatformRemove(opts UserRemoveOptions) *apperror.AppError {
	switch runtime.GOOS {
	case "windows":
		return removeWindowsUser(opts)
	case "linux":
		return removeLinuxUser(opts)
	default:
		return apperror.NewExecutionError(fmt.Sprintf("unsupported os for removal: %s", runtime.GOOS))
	}
}

func removeWindowsUser(opts UserRemoveOptions) *apperror.AppError {
	if opts.IsDryRun {
		return nil
	}
	cmd := exec.Command("net", "user", opts.Username, "/DELETE")
	out, err := defaultOSCommandRunner(cmd)
	if err != nil {
		return apperror.Wrap(err, "osuser.removeWindowsUser", map[string]any{"output": string(out)})
	}
	return nil
}

func removeLinuxUser(opts UserRemoveOptions) *apperror.AppError {
	appErr := cleanLinuxSudoersIfRequested(opts)
	if appErr != nil {
		return appErr
	}
	return deleteLinuxUserAccount(opts)
}

func cleanLinuxSudoersIfRequested(opts UserRemoveOptions) *apperror.AppError {
	if !opts.IsCleanSudoers {
		return nil
	}
	appErr := CleanSudoers(opts.Username, opts.IsDryRun)
	if appErr != nil {
		return appErr
	}
	return purgeUserSudoGroup(opts.Username, opts.IsDryRun)
}

func purgeUserSudoGroup(username string, isDryRun bool) *apperror.AppError {
	if isDryRun {
		return nil
	}
	cmd := exec.Command("gpasswd", "-d", username, "sudo")
	_, _ = defaultOSCommandRunner(cmd)
	return nil
}

func deleteLinuxUserAccount(opts UserRemoveOptions) *apperror.AppError {
	if opts.IsDryRun {
		return nil
	}
	deluserBin, err := exec.LookPath("deluser")
	hasDeluser := err == nil
	if hasDeluser {
		return runDeluserCommand(deluserBin, opts.Username, opts.IsRemoveHome)
	}
	return runUserdelCommand(opts.Username, opts.IsRemoveHome)
}

func runDeluserCommand(binary, username string, isRemoveHome bool) *apperror.AppError {
	args := buildDeluserArgs(username, isRemoveHome)
	cmd := exec.Command(binary, args...)
	out, err := defaultOSCommandRunner(cmd)
	if err != nil {
		return apperror.Wrap(err, "osuser.runDeluserCommand", map[string]any{"output": string(out)})
	}
	return nil
}

func buildDeluserArgs(username string, isRemoveHome bool) []string {
	if isRemoveHome {
		return []string{"--remove-home", username}
	}
	return []string{username}
}

func runUserdelCommand(username string, isRemoveHome bool) *apperror.AppError {
	args := buildUserdelArgs(username, isRemoveHome)
	cmd := exec.Command("userdel", args...)
	out, err := defaultOSCommandRunner(cmd)
	if err != nil {
		return apperror.Wrap(err, "osuser.runUserdelCommand", map[string]any{"output": string(out)})
	}
	return nil
}

func buildUserdelArgs(username string, isRemoveHome bool) []string {
	if isRemoveHome {
		return []string{"-r", username}
	}
	return []string{username}
}
