//go:build !windows
// +build !windows

package cmdchrome

import (
	"os/exec"
)

func configureDetachedProcess(_ *exec.Cmd) {}
