// Package installer — runner_lang.go executes scripts in specified language runtimes.
package installer

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// RunLanguageScript executes a script payload using the appropriate interpreter engine.

func RunLanguageScript(ctx context.Context, script, lang string) error {
	trimmed := strings.TrimSpace(script)
	if trimmed == "" {
		return apperror.New("RunLanguageScript", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "script cannot be empty"})
	}

	var cmd *exec.Cmd
	switch strings.ToLower(lang) {
	case "python", "py", "python3":
		cmd = exec.CommandContext(ctx, "python", "-c", trimmed)
	case "node", "js", "javascript":
		cmd = exec.CommandContext(ctx, "node", "-e", trimmed)
	case "powershell", "ps1", "pwsh":
		cmd = exec.CommandContext(ctx, "powershell", "-NoProfile", "-Command", trimmed)
	case "bash":
		cmd = exec.CommandContext(ctx, "bash", "-c", trimmed)
	default:
		if runtime.GOOS == "windows" {
			cmd = exec.CommandContext(ctx, "cmd", "/C", trimmed)
		} else {
			cmd = exec.CommandContext(ctx, "sh", "-c", trimmed)
		}
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
