// Package cmd — open_opener.go: native OS opener command execution.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func launchNativeOpener(target string) error {
	cmd := buildOpenerCommand(target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open target %q: %w", target, err)
	}

	return cmd.Wait()
}

func buildOpenerCommand(target string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", target)
	}

	if runtime.GOOS == "darwin" {
		return exec.Command("open", target)
	}

	return exec.Command("xdg-open", target)
}
