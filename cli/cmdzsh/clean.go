// Package cmdzsh provides cleanup and uninstallation logic for Oh-My-Zsh and .zshrc.
package cmdzsh

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// BackupZshrc copies target ~/.zshrc to a timestamped backup file.
func BackupZshrc(targetHome string) ZshStringResult {
	home := ResolveTargetHome(targetHome)
	zshrcPath := filepath.Join(home, ".zshrc")
	if !HasFile(zshrcPath) {
		return result.Ok("")
	}

	timestamp := time.Now().Format("20060102150405")
	backupPath := filepath.Join(home, fmt.Sprintf(".zshrc.bak.%s", timestamp))
	appErr := CopyFile(zshrcPath, backupPath)
	if appErr != nil {
		return result.Fail[string](appErr)
	}

	return result.Ok(backupPath)
}

// RemoveOhMyZshFiles deletes ~/.oh-my-zsh and ~/.zshrc.
func RemoveOhMyZshFiles(targetHome string) *apperror.AppError {
	home := ResolveTargetHome(targetHome)
	ohMyZshDir := filepath.Join(home, ".oh-my-zsh")
	zshrcPath := filepath.Join(home, ".zshrc")

	_ = os.RemoveAll(ohMyZshDir)
	_ = os.Remove(zshrcPath)

	return nil
}

// CleanOhMyZsh orchestrates backup, directory removal, and optional clean reinstallation.
func CleanOhMyZsh(opts ZshCleanOptions) *apperror.AppError {
	if opts.IsDryRun {
		fmt.Printf("[dry-run] Would clean Oh-My-Zsh and remove .zshrc%s", constants.NewLineUnix)

		return nil
	}

	return executeCleanSequence(opts)
}

func executeCleanSequence(opts ZshCleanOptions) *apperror.AppError {
	appErr := performBackupIfRequested(opts)
	if appErr != nil {
		return appErr
	}

	return purgeAndReinstall(opts)
}

func purgeAndReinstall(opts ZshCleanOptions) *apperror.AppError {
	appErr := RemoveOhMyZshFiles(opts.TargetHome)
	if appErr != nil {
		return appErr
	}

	if opts.IsReinstall {
		return executeCleanReinstall(opts)
	}

	return nil
}

func performBackupIfRequested(opts ZshCleanOptions) *apperror.AppError {
	if !opts.IsBackup {
		return nil
	}

	res := BackupZshrc(opts.TargetHome)
	if res.IsFailure() {
		return res.AppError()
	}

	return nil
}

func executeCleanReinstall(opts ZshCleanOptions) *apperror.AppError {
	suiteOpts := ZshOptions{
		Theme:      opts.Theme,
		TargetHome: opts.TargetHome,
		IsDryRun:   opts.IsDryRun,
	}

	return InstallZshSuite(suiteOpts)
}
