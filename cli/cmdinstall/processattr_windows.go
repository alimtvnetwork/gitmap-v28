//go:build windows

package cmdinstall

import (
	"os/exec"
	"syscall"
)

func SetHiddenProcessAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
