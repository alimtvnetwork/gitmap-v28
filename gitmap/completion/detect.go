package completion

import (
	"os"
	"runtime"
	"strings"
)

// DetectShell returns the current shell type based on environment.
func DetectShell() string {
	if runtime.GOOS == "windows" {
		return detectWindowsShell()
	}

	return detectUnixShell()
}

// detectWindowsShell checks for PowerShell on Windows.
func detectWindowsShell() string {
	psModule := os.Getenv("PSModulePath")
	if len(psModule) > 0 {
		return "powershell"
	}

	return "powershell"
}

// detectUnixShell inspects $SHELL for bash or zsh.
func detectUnixShell() string {
	shell := os.Getenv("SHELL")
	if strings.Contains(shell, "zsh") {
		return "zsh"
	}
	if strings.Contains(shell, "bash") {
		return "bash"
	}

	psModule := os.Getenv("PSModulePath")
	if len(psModule) > 0 {
		return "powershell"
	}

	return "bash"
}
