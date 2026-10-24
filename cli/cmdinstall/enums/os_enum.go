package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
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

// IsValid reports whether the OS matches a known operating system.
func (o OSType) IsValid() bool {
	switch o {
	case OSUnknown, OSWindows, OSWindows11, OSWindows10, OSWindowsServer,
		OSMacOS, OSUbuntu, OSDebian, OSFedora, OSCentOS:
		return true
	default:
		return false
	}
}

// String returns the string representation of the operating system.
func (o OSType) String() string {
	return string(o)
}

// ParseOSType parses a raw string into a validated OSType.
func ParseOSType(s string) (OSType, *apperror.AppError) {
	val := OSType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseOSType", "E1000", fmt.Sprintf("invalid OS: %q", s)).WithContext("value", s)
}

// AllOSTypes returns all valid OSType variants.
func AllOSTypes() []OSType {
	return []OSType{
		OSUnknown,
		OSWindows,
		OSWindows11,
		OSWindows10,
		OSWindowsServer,
		OSMacOS,
		OSUbuntu,
		OSDebian,
		OSFedora,
		OSCentOS,
	}
}

// IsWindows reports whether the OS is any Windows version.
func (o OSType) IsWindows() bool {
	return o == OSWindows || o == OSWindows10 || o == OSWindows11 || o == OSWindowsServer
}

// IsWindows11 reports whether the OS is Windows 11.
func (o OSType) IsWindows11() bool {
	return o == OSWindows11
}

// IsWindowsServer reports whether the OS is Windows Server.
func (o OSType) IsWindowsServer() bool {
	return o == OSWindowsServer
}

// IsMacOS reports whether the OS is macOS.
func (o OSType) IsMacOS() bool {
	return o == OSMacOS
}

// IsUbuntu reports whether the OS is Ubuntu.
func (o OSType) IsUbuntu() bool {
	return o == OSUbuntu
}

// IsDebian reports whether the OS is Debian.
func (o OSType) IsDebian() bool {
	return o == OSDebian
}

// IsFedora reports whether the OS is Fedora.
func (o OSType) IsFedora() bool {
	return o == OSFedora
}

// IsCentOS reports whether the OS is CentOS.
func (o OSType) IsCentOS() bool {
	return o == OSCentOS
}
