package cmd

import (
	"os"
	"runtime"
	"strings"
)

// OSType represents the specific operating system or distribution.
type OSType string

const (
	OSUnknown       OSType = "unknown"
	OSWindows       OSType = "windows"
	OSWindows11     OSType = "windows11"
	OSWindows10     OSType = "windows10"
	OSWindowsServer OSType = "windowsserver"
	OSMacOS         OSType = "macos"
	OSUbuntu        OSType = "ubuntu"
	OSDebian        OSType = "debian"
	OSFedora        OSType = "fedora"
	OSCentOS        OSType = "centos"
)

// CurrentOS detects and returns the specific operating system.
func CurrentOS() OSType {
	switch runtime.GOOS {
	case "windows":
		return detectWindowsVersion()
	case "darwin":
		return OSMacOS
	case "linux":
		return detectLinuxDistro()
	default:
		return OSUnknown
	}
}

// IsUbuntu returns true if the current OS is Ubuntu.
func IsUbuntu() bool { return CurrentOS() == OSUbuntu }

// IsDebian returns true if the current OS is Debian.
func IsDebian() bool { return CurrentOS() == OSDebian }

// IsFedora returns true if the current OS is Fedora.
func IsFedora() bool { return CurrentOS() == OSFedora }

// IsCentOS returns true if the current OS is CentOS.
func IsCentOS() bool { return CurrentOS() == OSCentOS }

// IsWindows returns true if the current OS is any Windows version.
func IsWindows() bool {
	osType := CurrentOS()
	return osType == OSWindows || osType == OSWindows10 || osType == OSWindows11 || osType == OSWindowsServer
}

// IsWindows11 returns true if the current OS is Windows 11.
func IsWindows11() bool { return CurrentOS() == OSWindows11 }

// IsWindowsServer returns true if the current OS is Windows Server.
func IsWindowsServer() bool { return CurrentOS() == OSWindowsServer }

func detectLinuxDistro() OSType {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return OSUnknown
	}
	content := strings.ToLower(string(data))
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "id=") {
			id := strings.Trim(strings.TrimPrefix(line, "id="), "\"'")
			switch id {
			case "ubuntu":
				return OSUbuntu
			case "debian":
				return OSDebian
			case "fedora":
				return OSFedora
			case "centos":
				return OSCentOS
			}
		}
	}
	return OSUnknown
}
