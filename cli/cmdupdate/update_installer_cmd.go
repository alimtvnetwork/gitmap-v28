package cmdupdate

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var targetPinnedVersion string

// SetTargetVersion sets an explicit target version to update to.
func SetTargetVersion(ver string) {
	targetPinnedVersion = ver
}

// GetTargetVersion returns the explicit target version if set.
func GetTargetVersion() string {
	return targetPinnedVersion
}

func buildRemoteInstallerCmd(scriptPath, installDir string) *exec.Cmd {
	return buildRemoteInstallerCmdWithVersion(scriptPath, installDir, targetPinnedVersion)
}

func buildRemoteInstallerCmdWithVersion(scriptPath, installDir, version string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return buildRemoteWindowsInstallerCmd(scriptPath, installDir, version)
	}

	return buildUnixInstallerCmd(scriptPath, installDir, version)
}

func sanitizeUnixInstallDir(dir string) string {
	base := filepath.Base(dir)
	if base == "gitmap-cli" || base == "gitmap" {
		return filepath.Dir(dir)
	}
	return dir
}

func buildRemoteWindowsInstallerCmd(scriptPath, installDir, version string) *exec.Cmd {
	args := []string{
		"-ExecutionPolicy", "Bypass",
		"-NoProfile", "-NoLogo",
		"-File", scriptPath,
	}

	if len(installDir) > 0 {
		args = append(args, "-InstallDir", installDir)
	}

	if len(version) > 0 {
		args = append(args, "-Version", version)
	}

	return exec.Command("powershell", args...)
}

func buildUnixInstallerCmd(scriptPath, installDir, version string) *exec.Cmd {
	args := []string{scriptPath}
	cleanDir := sanitizeUnixInstallDir(installDir)
	if len(cleanDir) > 0 {
		args = append(args, "--dir", cleanDir)
	}

	if len(version) > 0 {
		args = append(args, "--version", formatInstallerVersionFlag(version))
	}

	return exec.Command(getUnixShell(), args...)
}

func formatInstallerVersionFlag(version string) string {
	if strings.HasPrefix(version, "v") {
		return version
	}
	return "v" + version
}

func getUnixShell() string {
	shell := "bash"
	_, errLook := exec.LookPath(shell)
	if errLook != nil {
		shell = "sh"
	}

	return shell
}

// BuildRemoteWindowsInstallerCmdForTest exports buildRemoteWindowsInstallerCmd for testing.
func BuildRemoteWindowsInstallerCmdForTest(scriptPath, installDir, version string) *exec.Cmd {
	return buildRemoteWindowsInstallerCmd(scriptPath, installDir, version)
}

// BuildUnixInstallerCmdForTest exports buildUnixInstallerCmd for testing.
func BuildUnixInstallerCmdForTest(scriptPath, installDir, version string) *exec.Cmd {
	return buildUnixInstallerCmd(scriptPath, installDir, version)
}
