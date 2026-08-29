// Package cmd — os_update.go executes system package updates across platforms.
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// ExecuteOSUpdate runs platform-specific system update commands.

func ExecuteOSUpdate(ctx context.Context) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		fmt.Println("→ Triggering Windows Update check...")
		cmd = exec.CommandContext(ctx, "powershell", "-Command", "Get-WindowsUpdate -Install -AcceptAll -IgnoreReboot")
	case "darwin":
		fmt.Println("→ Triggering macOS softwareupdate and brew upgrade...")
		cmd = exec.CommandContext(ctx, "sh", "-c", "softwareupdate -ia --verbose || brew upgrade")
	default: // Linux
		fmt.Println("→ Triggering Linux package manager update...")
		cmd = exec.CommandContext(ctx, "sh", "-c", "sudo apt-get update && sudo apt-get upgrade -y || sudo dnf upgrade -y || sudo pacman -Syu --noconfirm")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
