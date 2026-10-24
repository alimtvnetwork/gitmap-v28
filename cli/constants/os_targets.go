// Package constants — os_targets.go defines supported target operating system identifiers.
package constants

import (
	"runtime"
)

// OSTargetType defines supported installer operating system targets.
type OSTargetType string

// Backwards-compatible type alias.
type OSTarget = OSTargetType

const (
	OSTargetInvalid OSTargetType = ""
	OSTargetWin     OSTargetType = "win"
	OSTargetUbuntu  OSTargetType = "ubuntu"
	OSTargetCentOS  OSTargetType = "centos"
	OSTargetDebian  OSTargetType = "debian"
	OSTargetFedora  OSTargetType = "fedora"
	OSTargetArch    OSTargetType = "arch-linux"
	OSTargetMac     OSTargetType = "macos"
	OSTargetUnix    OSTargetType = "unix"
	OSTargetAll     OSTargetType = "all"
)

// OSGroupType categorizes operating system families.
type OSGroupType string

// Backwards-compatible type alias.
type OSGroup = OSGroupType

const (
	OSGroupInvalid OSGroupType = ""
	OSGroupWindows OSGroupType = "windows"
	OSGroupUnix    OSGroupType = "unix"
	OSGroupMac     OSGroupType = "macos"
	OSGroupPOSIX   OSGroupType = "posix"
)

// OrderModeType specifies installer execution ordering.
type OrderModeType string

// Backwards-compatible type alias.
type OrderMode = OrderModeType

const (
	OrderModeInvalid   OrderModeType = ""
	OrderModeUnixFirst OrderModeType = "unix-first"
	OrderModeOSFirst   OrderModeType = "os-first"
	OrderModeOSOnly    OrderModeType = "os-only"
	OrderModeFallback  OrderModeType = "fallback"
)

// SupportedOSTargets returns the list of all supported installer operating systems.
var SupportedOSTargets = []string{
	string(OSTargetWin),
	string(OSTargetUbuntu),
	string(OSTargetDebian),
	string(OSTargetCentOS),
	string(OSTargetFedora),
	string(OSTargetArch),
	string(OSTargetMac),
	string(OSTargetUnix),
	string(OSTargetAll),
}

// IsValid verifies whether the target is a known supported operating system.
func (t OSTargetType) IsValid() bool {
	hasTarget := t.IsWindows() || t.IsUbuntu() || t.IsDebian() ||
		t.IsCentOS() || t.IsFedora() || t.IsArch() ||
		t.IsMac() || t.IsUnix() || t.IsAll()

	return hasTarget
}

func (t OSTargetType) IsWindows() bool { return t == OSTargetWin }

func (t OSTargetType) IsUbuntu() bool { return t == OSTargetUbuntu }

func (t OSTargetType) IsDebian() bool { return t == OSTargetDebian }

func (t OSTargetType) IsCentOS() bool { return t == OSTargetCentOS }

func (t OSTargetType) IsFedora() bool { return t == OSTargetFedora }

func (t OSTargetType) IsArch() bool { return t == OSTargetArch }

func (t OSTargetType) IsMac() bool { return t == OSTargetMac }

func (t OSTargetType) IsUnix() bool { return t == OSTargetUnix }

func (t OSTargetType) IsAll() bool { return t == OSTargetAll }

// String returns the canonical serialized string representation.
func (t OSTargetType) String() string {
	return string(t)
}

// Label returns the human-readable display label.
func (t OSTargetType) Label() string {
	return t.String()
}

// IsValid verifies whether the OS group is recognized.
func (g OSGroupType) IsValid() bool {
	hasGroup := g.IsWindows() || g.IsUnix() || g.IsMac() || g.IsPOSIX()

	return hasGroup
}

func (g OSGroupType) IsWindows() bool { return g == OSGroupWindows }

func (g OSGroupType) IsUnix() bool { return g == OSGroupUnix }

func (g OSGroupType) IsMac() bool { return g == OSGroupMac }

func (g OSGroupType) IsPOSIX() bool { return g == OSGroupPOSIX }

// String returns the canonical serialized string representation.
func (g OSGroupType) String() string {
	return string(g)
}

// Label returns the human-readable display label.
func (g OSGroupType) Label() string {
	return g.String()
}

// IsValid verifies whether the order mode is valid.
func (m OrderModeType) IsValid() bool {
	hasMode := m == OrderModeUnixFirst || m == OrderModeOSFirst ||
		m == OrderModeOSOnly || m == OrderModeFallback

	return hasMode
}

// String returns the canonical serialized string representation.
func (m OrderModeType) String() string {
	return string(m)
}

// Label returns the human-readable display label.
func (m OrderModeType) Label() string {
	return m.String()
}

// IsWindowsPlatform reports whether the current host OS runtime is Windows.
func IsWindowsPlatform() bool {
	return runtime.GOOS == "windows"
}

// IsDarwinPlatform reports whether the current host OS runtime is macOS (Darwin).
func IsDarwinPlatform() bool {
	return runtime.GOOS == "darwin"
}

// IsLinuxPlatform reports whether the current host OS runtime is Linux.
func IsLinuxPlatform() bool {
	return runtime.GOOS == "linux"
}
