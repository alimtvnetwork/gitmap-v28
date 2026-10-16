//go:build windows

package cmdos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"golang.org/x/sys/windows/registry"
)

func probeDetailedOS() OSVersionInfo {
	k, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`,
		registry.QUERY_VALUE,
	)
	if err != nil {
		return defaultWindowsOSVersion()
	}
	defer k.Close()

	installType, _, _ := k.GetStringValue("InstallationType")
	editionID, _, _ := k.GetStringValue("EditionID")
	productName, _, _ := k.GetStringValue("ProductName")
	displayVersion, _, _ := k.GetStringValue("DisplayVersion")
	releaseID, _, _ := k.GetStringValue("ReleaseId")
	majorNT, _, errMaj := k.GetIntegerValue("CurrentMajorVersionNumber")
	minorNT, _, _ := k.GetIntegerValue("CurrentMinorVersionNumber")
	buildStr, _, _ := k.GetStringValue("CurrentBuildNumber")
	if buildStr == "" {
		buildStr, _, _ = k.GetStringValue("CurrentBuild")
	}
	ubr, _, _ := k.GetIntegerValue("UBR")

	buildNum, _ := strconv.Atoi(buildStr)
	major := int(majorNT)
	if errMaj != nil || major == 0 {
		major = 10
	}

	lowInstall := strings.ToLower(installType)
	lowEdition := strings.ToLower(editionID)
	lowProduct := strings.ToLower(productName)

	isServer := lowInstall == "server" || strings.Contains(lowEdition, "server") || strings.Contains(lowProduct, "server")
	isPro := false
	isHome := false

	if !isServer {
		isPro = strings.Contains(lowEdition, "professional") || strings.Contains(lowEdition, "pro") || strings.Contains(lowProduct, "pro")
		isHome = strings.EqualFold(lowEdition, "core") || strings.Contains(lowEdition, "home") || strings.Contains(lowProduct, "home")
	}

	if buildNum >= 22000 && !isServer {
		major = 11
	}

	edition := resolveWindowsEdition(productName, editionID, isServer, major, buildNum)
	release := displayVersion
	if release == "" {
		release = releaseID
	}

	kernelVer := fmt.Sprintf("%d.%d.%d.%d", majorNT, minorNT, buildNum, ubr)
	version := kernelVer

	return OSVersionInfo{
		Type:            "windows",
		IsWindowsServer: isServer,
		IsWindowsPro:    isPro,
		IsWindowsHome:   isHome,
		Version:         version,
		Major:           major,
		Minor:           int(minorNT),
		Build:           buildNum,
		Release:         release,
		Edition:         edition,
		KernelVersion:   kernelVer,
	}
}

func resolveWindowsEdition(product, editionID string, isServer bool, major, buildNum int) string {
	res := product
	if buildNum >= 22000 && !isServer && strings.Contains(res, "Windows 10") {
		res = strings.Replace(res, "Windows 10", "Windows 11", 1)
	}

	if isServer && strings.HasPrefix(res, "Windows ") {
		res = strings.TrimPrefix(res, "Windows ")
	}

	if res != "" {
		return res
	}

	if isServer {
		return "Server " + editionID
	}

	if major == 11 {
		return "Windows 11 " + editionID
	}

	return "Windows 10 " + editionID
}

func defaultWindowsOSVersion() OSVersionInfo {
	return OSVersionInfo{
		Type:            "windows",
		IsWindowsServer: false,
		IsWindowsPro:    false,
		IsWindowsHome:   false,
		Version:         "10.0.0",
		Major:           10,
		Minor:           0,
		Build:           0,
		Release:         "",
		Edition:         "Windows",
		KernelVersion:   "10.0.0",
	}
}

func probeLocalPlatformOS() LocalOSProbe {
	det := probeDetailedOS()
	buildStr := strconv.Itoa(det.Build)
	fullVer := assembleWindowsVersion(det.Edition, det.Release, buildStr)

	return LocalOSProbe{
		OSType:           constants.OSTargetWin,
		OSGroup:          constants.OSGroupWindows,
		OSVersion:        fullVer,
		BuildVersion:     buildStr,
		Kernel:           det.KernelVersion,
		OSVersionDetails: det,
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
	if build != "" && build != "0" {
		parts = append(parts, fmt.Sprintf("(Build %s)", build))
	}
	full := strings.Join(parts, " ")
	if full == "" {
		return "Windows"
	}
	return full
}
