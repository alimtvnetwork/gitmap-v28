//go:build !windows

package cmdos

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func probeLocalPlatformOS() LocalOSProbe {
	if runtime.GOOS == "darwin" {
		return probeDarwinOS()
	}
	return probeLinuxOS()
}

func probeDarwinOS() LocalOSProbe {
	out, _ := exec.Command("sw_vers", "-productVersion").Output()
	ver := strings.TrimSpace(string(out))
	return LocalOSProbe{
		OSType:       constants.OSTargetMac,
		OSGroup:      constants.OSGroupMac,
		OSVersion:    "macOS " + ver,
		BuildVersion: ver,
		Kernel:       "",
	}
}

func probeLinuxOS() LocalOSProbe {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return defaultUnixProbe()
	}
	osType, pretty, buildVer := parseOSReleaseContent(string(data))
	unameOut, _ := exec.Command("uname", "-r").Output()
	return LocalOSProbe{
		OSType:       osType,
		OSGroup:      constants.OSGroupUnix,
		OSVersion:    pretty,
		BuildVersion: buildVer,
		Kernel:       strings.TrimSpace(string(unameOut)),
	}
}

func defaultUnixProbe() LocalOSProbe {
	return LocalOSProbe{
		OSType:       constants.OSTargetUnix,
		OSGroup:      constants.OSGroupPOSIX,
		OSVersion:    runtime.GOOS,
		BuildVersion: "",
		Kernel:       "",
	}
}

func parseOSReleaseContent(data string) (string, string, string) {
	lines := strings.Split(data, "\n")
	osType := constants.OSTargetUnix
	pretty := "Linux"
	buildVer := ""

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "ID=") {
			osType = mapDistroToOSTarget(trimmed)
		}
		if strings.HasPrefix(trimmed, "PRETTY_NAME=") {
			pretty = strings.Trim(strings.TrimPrefix(trimmed, "PRETTY_NAME="), "\"")
		}
		if strings.HasPrefix(trimmed, "VERSION_ID=") {
			buildVer = strings.Trim(strings.TrimPrefix(trimmed, "VERSION_ID="), "\"")
		}
	}
	return osType, pretty, buildVer
}

func mapDistroToOSTarget(line string) string {
	raw := strings.ToLower(strings.Trim(strings.TrimPrefix(line, "ID="), "\""))
	switch raw {
	case "ubuntu":
		return constants.OSTargetUbuntu
	case "debian":
		return constants.OSTargetDebian
	case "centos":
		return constants.OSTargetCentOS
	case "fedora":
		return constants.OSTargetFedora
	case "arch":
		return constants.OSTargetArch
	default:
		return raw
	}
}
