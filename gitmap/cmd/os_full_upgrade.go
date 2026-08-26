// Package cmd — os_full_upgrade.go handles full OS version upgrades.
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// ExecuteOSFullUpgrade executes full OS distribution version upgrades.
func ExecuteOSFullUpgrade(ctx context.Context) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		fmt.Println("→ Triggering Windows Feature Upgrade...")
		cmd = exec.CommandContext(ctx, "powershell", "-Command", "Start-Process ms-settings:windowsupdate")
	case "darwin":
		fmt.Println("→ Triggering macOS Major Version Upgrade...")
		cmd = exec.CommandContext(ctx, "softwareupdate", "--fetch-full-installer")
	default: // Linux (Ubuntu/Debian)
		fmt.Println("→ Triggering Linux Full Release Upgrade...")
		cmd = exec.CommandContext(ctx, "sh", "-c", "sudo do-release-upgrade || sudo apt-get dist-upgrade -y")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
