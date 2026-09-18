//go:build windows

package cmdagy

import (
	"os/exec"
	"syscall"
)

// configureBackgroundProcess detaches the spawned process from the parent console on Windows.
func configureBackgroundProcess(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	// 0x00000008 = DETACHED_PROCESS
	cmd.SysProcAttr.CreationFlags |= 0x00000008
}
