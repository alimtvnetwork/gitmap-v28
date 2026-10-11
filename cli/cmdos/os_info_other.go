//go:build !windows

package cmdos

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func probeDetailedOS() OSVersionInfo {
	if runtime.GOOS == "darwin" {
		return probeDetailedDarwinOS()
	}
	return probeDetailedLinuxOS()
}

func probeDetailedDarwinOS() OSVersionInfo {
	verOut, _ := exec.Command("sw_vers", "-productVersion").Output()
	unameOut, _ := exec.Command("uname", "-r").Output()

	productVer := strings.TrimSpace(string(verOut))
	kernelVer := strings.TrimSpace(string(unameOut))

	major, minor, build := parseNumericVersion(productVer)
	releaseName := mapDarwinMajorToRelease(major)
	edition := "macOS " + productVer

	return OSVersionInfo{
		Type:            "darwin",
		IsWindowsServer: false,
		IsWindowsPro:    false,
		IsWindowsHome:   false,
		Version:         productVer,
		Major:           major,
		Minor:           minor,
		Build:           build,
		Release:         releaseName,
		Edition:         edition,
		KernelVersion:   kernelVer,
	}
}

func mapDarwinMajorToRelease(major int) string {
	switch major {
	case 15:
		return "Sequoia"
	case 14:
		return "Sonoma"
	case 13:
		return "Ventura"
	case 12:
		return "Monterey"
	case 11:
		return "Big Sur"
	default:
		return "macOS"
	}
}

func probeDetailedLinuxOS() OSVersionInfo {
	data, err := readOSReleaseData()
	if err != nil {
		return defaultLinuxOSVersion()
	}

	pretty, verID, codename := parseOSReleaseMap(data)
	unameOut, _ := exec.Command("uname", "-r").Output()
	kernelVer := strings.TrimSpace(string(unameOut))
	major, minor, build := parseNumericVersion(verID)

	release := codename
	if release == "" {
		release = verID
	}

	return OSVersionInfo{
		Type:            "linux",
		IsWindowsServer: false,
		IsWindowsPro:    false,
		IsWindowsHome:   false,
		Version:         verID,
		Major:           major,
		Minor:           minor,
		Build:           build,
		Release:         release,
		Edition:         pretty,
		KernelVersion:   kernelVer,
	}
}

func readOSReleaseData() (string, error) {
	data, err := os.ReadFile("/etc/os-release")
	if err == nil {
		return string(data), nil
	}
	data2, err2 := os.ReadFile("/usr/lib/os-release")
	if err2 == nil {
		return string(data2), nil
	}
	return "", err2
}

func parseOSReleaseMap(data string) (string, string, string) {
	pretty := "Linux"
	verID := ""
	codename := ""

	lines := strings.Split(data, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "PRETTY_NAME=") {
			pretty = strings.Trim(strings.TrimPrefix(trimmed, "PRETTY_NAME="), "\"")
		}
		if strings.HasPrefix(trimmed, "VERSION_ID=") {
			verID = strings.Trim(strings.TrimPrefix(trimmed, "VERSION_ID="), "\"")
		}
		if strings.HasPrefix(trimmed, "VERSION_CODENAME=") {
			codename = strings.Trim(strings.TrimPrefix(trimmed, "VERSION_CODENAME="), "\"")
		}
	}
	return pretty, verID, codename
}

func parseNumericVersion(ver string) (int, int, int) {
	parts := strings.Split(ver, ".")
	major := 0
	minor := 0
	build := 0

	if len(parts) > 0 {
		major, _ = strconv.Atoi(parts[0])
	}
	if len(parts) > 1 {
		minor, _ = strconv.Atoi(parts[1])
	}
	if len(parts) > 2 {
		build, _ = strconv.Atoi(parts[2])
	}
	return major, minor, build
}

func defaultLinuxOSVersion() OSVersionInfo {
	return OSVersionInfo{
		Type:            "linux",
		IsWindowsServer: false,
		IsWindowsPro:    false,
		IsWindowsHome:   false,
		Version:         "Linux",
		Major:           0,
		Minor:           0,
		Build:           0,
		Release:         "",
		Edition:         "Linux",
		KernelVersion:   "",
	}
}

func probeLocalPlatformOS() LocalOSProbe {
	det := probeDetailedOS()
	if det.Type == "darwin" {
		return LocalOSProbe{
			OSType:           constants.OSTargetMac,
			OSGroup:          constants.OSGroupMac,
			OSVersion:        det.Edition,
			BuildVersion:     det.Version,
			Kernel:           det.KernelVersion,
			OSVersionDetails: det,
		}
	}

	return LocalOSProbe{
		OSType:           constants.OSTargetUnix,
		OSGroup:          constants.OSGroupUnix,
		OSVersion:        det.Edition,
		BuildVersion:     det.Version,
		Kernel:           det.KernelVersion,
		OSVersionDetails: det,
	}
}
