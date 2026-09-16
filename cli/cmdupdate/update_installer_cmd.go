package cmdupdate

import (
	"os/exec"
	"runtime"
)

func buildRemoteInstallerCmd(scriptPath, installDir string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return buildRemoteWindowsInstallerCmd(scriptPath, installDir)
	}

	return buildUnixInstallerCmd(scriptPath, installDir)
}

func buildRemoteWindowsInstallerCmd(scriptPath, installDir string) *exec.Cmd {
	args := []string{
		"-ExecutionPolicy", "Bypass",
		"-NoProfile", "-NoLogo",
		"-File", scriptPath,
	}

	if len(installDir) > 0 {
		args = append(args, "-InstallDir", installDir)
	}

	return exec.Command("powershell", args...)
}

func buildUnixInstallerCmd(scriptPath, installDir string) *exec.Cmd {
	args := []string{scriptPath}
	if len(installDir) > 0 {
		args = append(args, "--dir", installDir)
	}

	return exec.Command(getUnixShell(), args...)
}

func getUnixShell() string {
	shell := "bash"
	_, errLook := exec.LookPath(shell)
	if errLook != nil {
		shell = "sh"
	}

	return shell
}
