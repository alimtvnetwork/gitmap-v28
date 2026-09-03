// Package installer — prompt_runner_windows.go executes the remote PowerShell installer on Windows.
package installer

import (
	"context"
	"os"
	"os/exec"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func RunWindowsPromptInstaller(targetDir string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	psExe := "powershell"
	if _, errPwsh := exec.LookPath("pwsh"); errPwsh == nil {
		psExe = "pwsh"
	}

	psScript := "Invoke-Expression \"& { $(Invoke-RestMethod " + constants.PromptArchitectPowerShellURL + ") }\""
	cmd := exec.CommandContext(ctx, psExe, "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", psScript)
	cmd.Dir = targetDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
