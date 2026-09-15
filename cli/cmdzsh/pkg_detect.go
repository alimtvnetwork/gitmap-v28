// Package cmdzsh provides package manager detection and command dispatching.
package cmdzsh

import (
	"os/exec"
)

// IsCommandAvailable checks if an executable exists on the PATH.
func IsCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)

	return err == nil
}

// DetectPackageManager detects the active system package manager.
func DetectPackageManager() PackageManagerType {
	candidates := []PackageManagerType{PkgMgrApt, PkgMgrDnf, PkgMgrBrew, PkgMgrPacman}
	for _, mgr := range candidates {
		if isPkgMgrInstalled(mgr) {
			return mgr
		}
	}

	return PkgMgrUnknown
}

func isPkgMgrInstalled(mgr PackageManagerType) bool {
	binaryMap := map[PackageManagerType]string{
		PkgMgrApt:    "apt-get",
		PkgMgrDnf:    "dnf",
		PkgMgrBrew:   "brew",
		PkgMgrPacman: "pacman",
	}

	bin, isFound := binaryMap[mgr]

	return isFound && IsCommandAvailable(bin)
}

// BuildInstallCommand resolves the package manager install command and arguments.
func BuildInstallCommand(mgr PackageManagerType, packages ...string) []string {
	prefix := resolveSudoPrefix()

	return assembleInstallArgs(prefix, mgr, packages...)
}

func resolveSudoPrefix() []string {
	if IsCommandAvailable("sudo") {
		return []string{"sudo"}
	}

	return []string{}
}

func assembleInstallArgs(prefix []string, mgr PackageManagerType, pkgs ...string) []string {
	switch mgr {
	case PkgMgrApt:
		return append(append(prefix, "apt-get", "install", "-y"), pkgs...)
	case PkgMgrDnf:
		return append(append(prefix, "dnf", "install", "-y"), pkgs...)
	case PkgMgrBrew:
		return append([]string{"brew", "install"}, pkgs...)
	case PkgMgrPacman:
		return append(append(prefix, "pacman", "-S", "--noconfirm"), pkgs...)
	case PkgMgrUnknown:
		return append(append(prefix, "apt-get", "install", "-y"), pkgs...)
	default:
		return append(append(prefix, "apt-get", "install", "-y"), pkgs...)
	}
}
