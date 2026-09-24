//go:build windows

package cmdos

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"golang.org/x/sys/windows/registry"
)

func probeLocalPlatformOS() LocalOSProbe {
	k, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`,
		registry.QUERY_VALUE,
	)
	if err != nil {
		return defaultWindowsProbe()
	}
	defer k.Close()

	productName, _, _ := k.GetStringValue("ProductName")
	displayVersion, _, _ := k.GetStringValue("DisplayVersion")
	currentBuild, _, _ := k.GetStringValue("CurrentBuild")
	fullVer := assembleWindowsVersion(productName, displayVersion, currentBuild)

	return LocalOSProbe{
		OSType:       constants.OSTargetWin,
		OSGroup:      constants.OSGroupWindows,
		OSVersion:    fullVer,
		BuildVersion: currentBuild,
		Kernel:       "",
	}
}

func defaultWindowsProbe() LocalOSProbe {
	return LocalOSProbe{
		OSType:       constants.OSTargetWin,
		OSGroup:      constants.OSGroupWindows,
		OSVersion:    "Windows",
		BuildVersion: "",
		Kernel:       "",
	}
}

func assembleWindowsVersion(name, dispVer, build string) string {
	var parts []string
	if name != "" {
		parts = append(parts, name)
	}
	if dispVer != "" {
		parts = append(parts, dispVer)
	}
	if build != "" {
		parts = append(parts, fmt.Sprintf("(Build %s)", build))
	}
	full := strings.Join(parts, " ")
	if full == "" {
		return "Windows"
	}
	return full
}
