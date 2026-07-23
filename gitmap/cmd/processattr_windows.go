//go:build windows

package cmd

import (
	"os/exec"
	"syscall"
)

// setHiddenProcessAttr hides transient helper processes on Windows.
func setHiddenProcessAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
