//go:build !windows

package cmd

import "os/exec"

// setHiddenProcessAttr is a no-op on non-Windows platforms.
func setHiddenProcessAttr(_ *exec.Cmd) {}
