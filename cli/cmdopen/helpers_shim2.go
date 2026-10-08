package cmdopen

import (
	"os/exec"
)

// ConfigureDetachedProcessFn is wired by cmd/di_hooks.go.
var ConfigureDetachedProcessFn func(cmd *exec.Cmd)

func configureDetachedProcess(cmd *exec.Cmd) {
	if ConfigureDetachedProcessFn != nil {
		ConfigureDetachedProcessFn(cmd)
	}
}
