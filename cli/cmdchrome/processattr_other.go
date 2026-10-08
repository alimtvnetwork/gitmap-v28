//go:build !windows
// +build !windows

package cmdchrome

import (
	"os/exec"
)

func ConfigureDetachedProcess(_ *exec.Cmd) {}
