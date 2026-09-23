//go:build windows

package cmdos

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func probeLocalPlatformOS() (string, string, string) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return "windows", "Windows", ""
	}
	defer k.Close()

	productName, _, _ := k.GetStringValue("ProductName")
	displayVersion, _, _ := k.GetStringValue("DisplayVersion")
	currentBuild, _, _ := k.GetStringValue("CurrentBuild")

	var parts []string
	if productName != "" {
		parts = append(parts, productName)
	}
	if displayVersion != "" {
		parts = append(parts, displayVersion)
	}
	if currentBuild != "" {
		parts = append(parts, fmt.Sprintf("(Build %s)", currentBuild))
	}

	fullVer := strings.Join(parts, " ")
	if fullVer == "" {
		fullVer = "Windows"
	}

	return "windows", fullVer, currentBuild
}
