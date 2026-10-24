package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// MusePlatformType represents the target operating system platform.
type MusePlatformType string

const (
	MusePlatformWindows MusePlatformType = "windows"
	MusePlatformLinux   MusePlatformType = "linux"
	MusePlatformDarwin  MusePlatformType = "darwin"
)

// IsValid reports whether the platform matches a known Muse platform.
func (p MusePlatformType) IsValid() bool {
	switch p {
	case MusePlatformWindows, MusePlatformLinux, MusePlatformDarwin:
		return true
	default:
		return false
	}
}

// String returns the string representation of the platform.
func (p MusePlatformType) String() string {
	return string(p)
}

// ParseMusePlatformType parses a raw string into a validated MusePlatformType.
// Alias spellings (win, mac, ubuntu) are not accepted here; the alias-aware
// resolver in the parent package handles those.
func ParseMusePlatformType(s string) (MusePlatformType, *apperror.AppError) {
	val := MusePlatformType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseMusePlatformType", "E1000", fmt.Sprintf("invalid muse platform: %q", s)).WithContext("value", s)
}

// AllMusePlatformTypes returns all valid MusePlatformType variants.
func AllMusePlatformTypes() []MusePlatformType {
	return []MusePlatformType{
		MusePlatformWindows,
		MusePlatformLinux,
		MusePlatformDarwin,
	}
}

// IsWindows reports whether the platform is Windows.
func (p MusePlatformType) IsWindows() bool {
	return p == MusePlatformWindows
}

// IsLinux reports whether the platform is Linux.
func (p MusePlatformType) IsLinux() bool {
	return p == MusePlatformLinux
}

// IsDarwin reports whether the platform is macOS.
func (p MusePlatformType) IsDarwin() bool {
	return p == MusePlatformDarwin
}
