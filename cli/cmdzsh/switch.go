// Package cmdzsh provides shell binary resolution and default shell switching logic.
package cmdzsh

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// ResolveZshBinary searches for the local zsh executable path.
func ResolveZshBinary() ZshStringResult {
	path, err := exec.LookPath("zsh")
	if err == nil && path != "" {
		return result.Ok(path)
	}

	candidates := []string{"/bin/zsh", "/usr/bin/zsh", "/usr/local/bin/zsh", "/opt/homebrew/bin/zsh"}
	found := searchCandidateShells(candidates)
	if found != "" {
		return result.Ok(found)
	}

	return result.Fail[string](apperror.NewExecutionError("zsh binary not found on system"))
}

func searchCandidateShells(candidates []string) string {
	for _, c := range candidates {
		if HasFile(c) {
			return c
		}
	}

	return ""
}

// ChangeDefaultShell updates the user default login shell using chsh -s.
func ChangeDefaultShell(opts ZshSwitchOptions) *apperror.AppError {
	if runtime.GOOS == "windows" {
		return apperror.NewExecutionError("changing login shell is only supported on UNIX platforms")
	}

	return applyChsh(opts)
}

func applyChsh(opts ZshSwitchOptions) *apperror.AppError {
	zshRes := ResolveZshBinary()
	if zshRes.IsFailure() {
		return zshRes.AppError()
	}

	args := buildChshArgs(zshRes.Value, opts.TargetUser)
	if opts.IsDryRun {
		fmt.Printf("[dry-run] Would execute: %v%s", args, constants.NewLineUnix)

		return nil
	}

	return executeChshCommand(args)
}

func buildChshArgs(zshPath, targetUser string) []string {
	if targetUser != "" {
		return []string{"chsh", "-s", zshPath, targetUser}
	}

	return []string{"chsh", "-s", zshPath}
}

func executeChshCommand(args []string) *apperror.AppError {
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return apperror.Wrap(err, "cmdzsh.executeChshCommand", map[string]any{"output": string(out)})
	}

	return nil
}
