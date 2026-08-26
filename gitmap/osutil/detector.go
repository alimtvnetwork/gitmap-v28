// Package osutil — detector.go provides cross-platform operating system detection.
package osutil

import (
	"os"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// DetailedOSInfo encapsulates detected operating system and package manager details.
type DetailedOSInfo struct {
	Family         string // "windows", "darwin", "linux"
	Distro         string // "ubuntu", "debian", "arch", "centos", "fedora", "unknown"
	PackageManager string // "apt", "pacman", "dnf", "yum", "brew", "winget"
}

// DetectHostOS determines the exact host OS family and Linux distribution.
func DetectHostOS() DetailedOSInfo {
	switch runtime.GOOS {
	case "windows":
		return DetailedOSInfo{
			Family:         constants.OSTargetWin,
			Distro:         constants.OSTargetWin,
			PackageManager: "winget",
		}
	case "darwin":
		return DetailedOSInfo{
			Family:         constants.OSTargetMac,
			Distro:         constants.OSTargetMac,
			PackageManager: "brew",
		}
	default: // linux / unix
		return detectLinuxDistro()
	}
}

func detectLinuxDistro() DetailedOSInfo {
	info := DetailedOSInfo{
		Family:         constants.OSTargetUnix,
		Distro:         constants.OSTargetUbuntu,
		PackageManager: "apt",
	}

	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return info
	}

	content := strings.ToLower(string(data))
	if strings.Contains(content, "arch") || strings.Contains(content, "manjaro") {
		info.Distro = constants.OSTargetArch
		info.PackageManager = "pacman"
	} else if strings.Contains(content, "fedora") {
		info.Distro = constants.OSTargetFedora
		info.PackageManager = "dnf"
	} else if strings.Contains(content, "centos") || strings.Contains(content, "rhel") {
		info.Distro = constants.OSTargetCentOS
		info.PackageManager = "yum"
	} else if strings.Contains(content, "debian") {
		info.Distro = constants.OSTargetDebian
		info.PackageManager = "apt"
	}

	return info
}
