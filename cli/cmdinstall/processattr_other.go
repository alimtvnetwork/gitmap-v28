//go:build !windows

package cmdinstall

import "os/exec"

func SetHiddenProcessAttr(cmd *exec.Cmd) {
	// No-op on non-Windows platforms.
}
