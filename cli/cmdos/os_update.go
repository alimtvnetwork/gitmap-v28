// Package cmdos — os_update.go executes system package updates across platforms.
package cmdos

import (
	"context"
	"os"
	"os/exec"
	"runtime"
)

// OSCommandRunner executes a configured exec.Cmd.
type OSCommandRunner func(cmd *exec.Cmd) error

var defaultOSCommandRunner OSCommandRunner = func(cmd *exec.Cmd) error {
	return cmd.Run()
}

// ExecuteOSUpdate runs platform-specific system update commands.
func ExecuteOSUpdate(ctx context.Context) error {
	cmd := buildOSUpdateCmd(ctx)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return defaultOSCommandRunner(cmd)
}

func buildOSUpdateCmd(ctx context.Context) *exec.Cmd {
	switch runtime.GOOS {
	case "windows":
		return exec.CommandContext(ctx, "powershell", "-Command", "Get-WindowsUpdate -Install -AcceptAll -IgnoreReboot")
	case "darwin":
		return exec.CommandContext(ctx, "sh", "-c", "softwareupdate -ia --verbose || brew upgrade")
	default:
		return exec.CommandContext(ctx, "sh", "-c", "sudo apt-get update && sudo apt-get upgrade -y || sudo dnf upgrade -y || sudo pacman -Syu --noconfirm")
	}
}
