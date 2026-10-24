package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// InstallerOSType defines supported operating systems for modular recipes.
type InstallerOSType string

const (
	InstallerOSWindows InstallerOSType = "windows"
	InstallerOSUnix    InstallerOSType = "unix"
	InstallerOSUbuntu  InstallerOSType = "ubuntu"
	InstallerOSCentOS  InstallerOSType = "centos"
	InstallerOSAll     InstallerOSType = "all"
)

// IsValid reports whether the OS matches a known installer OS.
func (o InstallerOSType) IsValid() bool {
	switch o {
	case InstallerOSWindows, InstallerOSUnix, InstallerOSUbuntu,
		InstallerOSCentOS, InstallerOSAll:
		return true
	default:
		return false
	}
}

// String returns the string representation of the installer OS.
func (o InstallerOSType) String() string {
	return string(o)
}

// ParseInstallerOSType parses a raw string into a validated InstallerOSType.
func ParseInstallerOSType(s string) (InstallerOSType, *apperror.AppError) {
	val := InstallerOSType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseInstallerOSType", "E1000", fmt.Sprintf("invalid installer OS: %q", s)).WithContext("value", s)
}

// AllInstallerOSTypes returns all valid InstallerOSType variants.
func AllInstallerOSTypes() []InstallerOSType {
	return []InstallerOSType{
		InstallerOSWindows,
		InstallerOSUnix,
		InstallerOSUbuntu,
		InstallerOSCentOS,
		InstallerOSAll,
	}
}

// IsWindows reports whether the installer OS is Windows.
func (o InstallerOSType) IsWindows() bool {
	return o == InstallerOSWindows
}

// IsUnix reports whether the installer OS is Unix.
func (o InstallerOSType) IsUnix() bool {
	return o == InstallerOSUnix
}

// IsUbuntu reports whether the installer OS is Ubuntu.
func (o InstallerOSType) IsUbuntu() bool {
	return o == InstallerOSUbuntu
}

// IsCentOS reports whether the installer OS is CentOS.
func (o InstallerOSType) IsCentOS() bool {
	return o == InstallerOSCentOS
}

// IsAll reports whether the installer OS matches all systems.
func (o InstallerOSType) IsAll() bool {
	return o == InstallerOSAll
}
