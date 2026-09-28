//go:build !windows

package cmdinstall

import "os/exec"

func setHiddenProcessAttr(cmd *exec.Cmd) {
	// No-op on non-Windows platforms.
}
