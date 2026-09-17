// Package cmdos — os_full_upgrade.go handles full OS version upgrades.
package cmdos

import (
	"context"
	"os"
	"os/exec"
	"runtime"
)

// ExecuteOSFullUpgrade executes full OS distribution version upgrades.
func ExecuteOSFullUpgrade(ctx context.Context) error {
	cmd := buildOSFullUpgradeCmd(ctx)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return defaultOSCommandRunner(cmd)
}

func buildOSFullUpgradeCmd(ctx context.Context) *exec.Cmd {
	switch runtime.GOOS {
	case "windows":
		return exec.CommandContext(ctx, "powershell", "-Command", "Start-Process ms-settings:windowsupdate")
	case "darwin":
		return exec.CommandContext(ctx, "softwareupdate", "--fetch-full-installer")
	default:
		return exec.CommandContext(ctx, "sh", "-c", "sudo do-release-upgrade || sudo apt-get dist-upgrade -y")
	}
}
