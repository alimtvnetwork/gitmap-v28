package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// BackendType enumerates the Windows-specific add targets. Linux/macOS
// have one canonical backend each (XDG .desktop / LaunchAgents
// .plist) so they ignore this value entirely. The zero value
// (BackendUnspecified) means "let the dispatcher pick the OS default"
// which is `BackendRegistry` on Windows.
type BackendType int

const (
	// BackendUnspecified is the zero value. The dispatcher
	// translates it to the per-OS default rather than failing —
	// keeps the public Add(opts) API ergonomic for non-Windows
	// callers that have no opinion on backend.
	BackendUnspecified BackendType = 0
	// BackendRegistry writes the HKCU Run-key value + tracking
	// subkey (per-user; the long-standing default).
	BackendRegistry BackendType = 1
	// BackendStartupFolder writes the .lnk + tracking subkey.
	BackendStartupFolder BackendType = 2
	// BackendRegistryHKLM writes the HKLM Run-key value +
	// tracking subkey (machine-wide; requires admin). Same on-
	// disk shape as BackendRegistry, just rooted under
	// HKEY_LOCAL_MACHINE so every interactive user on the
	// machine triggers the autostart at login.
	BackendRegistryHKLM BackendType = 3
)

// IsValid reports whether the value matches a known variant.
func (b BackendType) IsValid() bool {
	switch b {
	case BackendUnspecified, BackendRegistry, BackendStartupFolder, BackendRegistryHKLM:
		return true
	default:
		return false
	}
}

func (b BackendType) String() string {
	switch b {
	case BackendRegistry:
		return constants.StartupBackendRegistry
	case BackendStartupFolder:
		return constants.StartupBackendStartupFolder
	case BackendRegistryHKLM:
		return constants.StartupBackendRegistryHKLM
	case BackendUnspecified:
		return ""
	default:
		return ""
	}
}

func ParseBackend(s string) (BackendType, *apperror.AppError) {
	switch s {
	case "":
		return BackendUnspecified, nil
	case constants.StartupBackendRegistry:
		return BackendRegistry, nil
	case constants.StartupBackendStartupFolder:
		return BackendStartupFolder, nil
	case constants.StartupBackendRegistryHKLM:
		return BackendRegistryHKLM, nil
	default:
		return BackendUnspecified, apperror.NewValidationError(fmt.Sprintf(constants.ErrStartupAddBadBackend, s))
	}
}

// AllBackendTypes returns all valid variants.
func AllBackendTypes() []BackendType {
	return []BackendType{
		BackendUnspecified,
		BackendRegistry,
		BackendStartupFolder,
		BackendRegistryHKLM,
	}
}

// IsUnspecified reports whether the value is BackendUnspecified.
func (b BackendType) IsUnspecified() bool {
	return b == BackendUnspecified
}

// IsRegistry reports whether the value is BackendRegistry.
func (b BackendType) IsRegistry() bool {
	return b == BackendRegistry
}

// IsStartupFolder reports whether the value is BackendStartupFolder.
func (b BackendType) IsStartupFolder() bool {
	return b == BackendStartupFolder
}

// IsRegistryHKLM reports whether the value is BackendRegistryHKLM.
func (b BackendType) IsRegistryHKLM() bool {
	return b == BackendRegistryHKLM
}

// HasBackend reports whether the value is non-empty and valid.
func (b BackendType) HasBackend() bool {
	return b.IsValid()
}
