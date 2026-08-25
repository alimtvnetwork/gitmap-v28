//go:build windows

package cmd

import (
	"strconv"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func detectWindowsVersion() OSType {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return OSWindows
	}
	defer k.Close()

	productName, _, err := k.GetStringValue("ProductName")
	if err == nil && strings.Contains(productName, "Server") {
		return OSWindowsServer
	}

	buildStr, _, err := k.GetStringValue("CurrentBuild")
	if err == nil {
		if build, err := strconv.Atoi(buildStr); err == nil {
			if build >= 22000 {
				return OSWindows11
			}
			return OSWindows10
		}
	}

	return OSWindows
}
