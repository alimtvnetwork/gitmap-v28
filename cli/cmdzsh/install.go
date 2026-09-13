// Package cmdzsh provides installation logic for ZSH, Oh-My-Zsh, and plugins.
package cmdzsh

import (
	"fmt"
	"os/exec"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// InstallZshPackage installs zsh and essential tools via the detected package manager.
func InstallZshPackage(isDryRun bool) *apperror.AppError {
	mgr := DetectPackageManager()
	if mgr == PkgMgrUnknown {
		return apperror.NewExecutionError("unsupported package manager for zsh install")
	}

	cmdArgs := BuildInstallCommand(mgr, "zsh", "git", "curl", "wget")
	if isDryRun {
		fmt.Printf("[dry-run] Would execute: %v%s", cmdArgs, constants.NewLineUnix)

		return nil
	}

	return executeInstallCommand(cmdArgs)
}

func executeInstallCommand(cmdArgs []string) *apperror.AppError {
	if len(cmdArgs) == 0 {
		return apperror.NewExecutionError("empty install command arguments")
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return apperror.Wrap(err, "cmdzsh.executeInstallCommand", map[string]any{"output": string(out)})
	}

	return nil
}

// InstallZshSuite orchestrates the complete installation of zsh, oh-my-zsh, and plugins.
func InstallZshSuite(opts ZshOptions) *apperror.AppError {
	appErr := installPackageIfRequested(opts)
	if appErr != nil {
		return appErr
	}

	appErr = InstallOhMyZsh(opts.TargetHome, opts.IsDryRun)
	if appErr != nil {
		return appErr
	}

	return completeInstallSuite(opts)
}

func installPackageIfRequested(opts ZshOptions) *apperror.AppError {
	if !opts.IsInstallZsh {
		return nil
	}

	return InstallZshPackage(opts.IsDryRun)
}

func completeInstallSuite(opts ZshOptions) *apperror.AppError {
	appErr := InstallAutosuggestions(opts.TargetHome, opts.IsDryRun)
	if appErr != nil {
		return appErr
	}

	theme := opts.Theme
	if theme == "" {
		theme = string(ZshThemeRobbyrussell)
	}

	return ApplyTheme(opts.TargetHome, theme, opts.IsDryRun)
}
