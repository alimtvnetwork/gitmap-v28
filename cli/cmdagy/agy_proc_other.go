//go:build !windows

package cmdagy

import "os/exec"

// configureBackgroundProcess is a no-op on non-Windows platforms.
func configureBackgroundProcess(cmd *exec.Cmd) {
}
