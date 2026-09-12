//go:build !windows

package cmdupdate

import "os/exec"

func setHiddenProcessAttr(_ *exec.Cmd) {}
