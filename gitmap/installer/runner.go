// Package installer — runner.go executes shell commands for installer steps.
package installer

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// RunInstallerCommand executes a shell command string with standard outputs attached.
func RunInstallerCommand(command string) error {
	trimmed := strings.TrimSpace(command)
	if trimmed == "" {
		return apperror.New("RunInstallerCommand", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "command cannot be empty"})
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", trimmed)
	} else {
		cmd = exec.Command("sh", "-c", trimmed)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
