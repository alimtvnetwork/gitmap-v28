package cmd

import (
	"os/exec"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// resolvePackageManager detects or uses the specified package manager.
func resolvePackageManager(override string) string {
	if override != "" {
		return override
	}

	return detectPackageManager()
}

// detectPackageManager finds the available package manager.
func detectPackageManager() string {
	if runtime.GOOS == "windows" {
		return detectWindowsManager()
	}
	if runtime.GOOS == "darwin" {
		return detectDarwinManager()
	}

	return detectLinuxManager()
}

// detectWindowsManager checks for Chocolatey then Winget.
func detectWindowsManager() string {
	if isCommandAvailable(constants.PkgMgrChocolatey) {
		return constants.PkgMgrChocolatey
	}
	if isCommandAvailable(constants.PkgMgrWinget) {
		return constants.PkgMgrWinget
	}

	return constants.PkgMgrChocolatey
}

// detectDarwinManager checks for Homebrew on macOS.
func detectDarwinManager() string {
	if isCommandAvailable(constants.PkgMgrBrew) {
		return constants.PkgMgrBrew
	}

	return constants.PkgMgrBrew
}

// detectLinuxManager checks for apt, dnf, then pacman.
func detectLinuxManager() string {
	if isCommandAvailable(constants.PkgMgrApt) {
		return constants.PkgMgrApt
	}
	if isCommandAvailable(constants.PkgMgrDnf) {
		return constants.PkgMgrDnf
	}
	if isCommandAvailable(constants.PkgMgrPacman) {
		return constants.PkgMgrPacman
	}

	return constants.PkgMgrApt
}

// isCommandAvailable checks if a command exists in PATH.
func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)

	return err == nil
}
