package cmdinstall

import (
	"os"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall/enums"
)

// OSType is an alias for the canonical enum in the package's enums
// subpackage, kept so existing callers keep compiling.
type OSType = enums.OSType

// OS constants, aliased from the enums subpackage.
const (
	OSUnknown       = enums.OSUnknown
	OSWindows       = enums.OSWindows
	OSWindows11     = enums.OSWindows11
	OSWindows10     = enums.OSWindows10
	OSWindowsServer = enums.OSWindowsServer
	OSMacOS         = enums.OSMacOS
	OSUbuntu        = enums.OSUbuntu
	OSDebian        = enums.OSDebian
	OSFedora        = enums.OSFedora
	OSCentOS        = enums.OSCentOS
)

// CurrentOS detects and returns the specific operating system.
func CurrentOS() enums.OSType {
	switch runtime.GOOS {
	case "windows":
		return detectWindowsVersion()
	case "darwin":
		return enums.OSMacOS
	case "linux":
		return detectLinuxDistro()
	default:
		return enums.OSUnknown
	}
}

// IsUbuntu returns true if the current OS is Ubuntu.
func IsUbuntu() bool { return CurrentOS() == enums.OSUbuntu }

// IsDebian returns true if the current OS is Debian.
func IsDebian() bool { return CurrentOS() == enums.OSDebian }

// IsFedora returns true if the current OS is Fedora.
func IsFedora() bool { return CurrentOS() == enums.OSFedora }

// IsCentOS returns true if the current OS is CentOS.
func IsCentOS() bool { return CurrentOS() == enums.OSCentOS }

// IsWindows returns true if the current OS is any Windows version.
func IsWindows() bool {
	osType := CurrentOS()

	return osType == enums.OSWindows || osType == enums.OSWindows10 || osType == enums.OSWindows11 || osType == enums.OSWindowsServer
}

// IsWindows11 returns true if the current OS is Windows 11.
func IsWindows11() bool { return CurrentOS() == enums.OSWindows11 }

// IsWindowsServer returns true if the current OS is Windows Server.
func IsWindowsServer() bool { return CurrentOS() == enums.OSWindowsServer }

func detectLinuxDistro() enums.OSType {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return enums.OSUnknown
	}

	content := strings.ToLower(string(data))
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "id=") {
			id := strings.Trim(strings.TrimPrefix(line, "id="), "\"'")
			switch id {
			case "ubuntu":
				return enums.OSUbuntu
			case "debian":
				return enums.OSDebian
			case "fedora":
				return enums.OSFedora
			case "centos":
				return enums.OSCentOS
			}
		}
	}

	return enums.OSUnknown
}
