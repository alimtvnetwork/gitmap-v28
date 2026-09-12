//go:build windows
// +build windows

package cmdchrome

import (
	"os/exec"
	"syscall"
)

func configureDetachedProcess(cmd *exec.Cmd) {
	if cmd == nil {
		return
	}

	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}

	cmd.SysProcAttr.CreationFlags |= 0x00000008 | 0x00000200
}
