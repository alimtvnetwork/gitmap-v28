// Package installer — prompt_runner_dispatch.go dispatches installer execution based on runtime.GOOS.
package installer

import (
	"runtime"
	"time"
)

func RunPromptInstallerForHost(targetDir string, timeout time.Duration) error {
	if timeout <= 0 {
		timeout = 2 * time.Minute
	}

	if runtime.GOOS == "windows" {
		return RunWindowsPromptInstaller(targetDir, timeout)
	}
	return RunUnixPromptInstaller(targetDir, timeout)
}
