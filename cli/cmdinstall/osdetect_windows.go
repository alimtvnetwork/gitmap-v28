//go:build windows

package cmdinstall

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall/enums"
	"strconv"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func detectWindowsVersion() enums.OSType {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return enums.OSWindows
	}

	defer k.Close()

	productName, _, err := k.GetStringValue("ProductName")
	if err == nil && strings.Contains(productName, "Server") {
		return enums.OSWindowsServer
	}

	buildStr, _, err := k.GetStringValue("CurrentBuild")
	if err != nil {
		return enums.OSWindows
	}

	build, err := strconv.Atoi(buildStr)
	if err != nil {
		return enums.OSWindows
	}

	if build >= 22000 {
		return enums.OSWindows11
	}

	return enums.OSWindows10
}
