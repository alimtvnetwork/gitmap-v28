//go:build windows

package cmdinstall

import (
	"os/exec"
	"syscall"
)

func setHiddenProcessAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
