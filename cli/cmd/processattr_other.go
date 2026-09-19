//go:build !windows

package cmd

import (
	"os/exec"
	"syscall"
)

// setHiddenProcessAttr is a no-op on non-Windows platforms.
func setHiddenProcessAttr(_ *exec.Cmd) {}

func configureDetachedProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}
