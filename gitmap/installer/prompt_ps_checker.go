// Package installer — prompt_ps_checker.go verifies PowerShell executable presence.
package installer

import "os/exec"

func HasPowerShell() bool {
	if _, err := exec.LookPath("pwsh"); err == nil {
		return true
	}

	if _, err := exec.LookPath("powershell"); err == nil {
		return true
	}

	return false
}
