package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// PackageManagerType represents the detected managing package manager.
type PackageManagerType string

const (
	// ManagerApt represents Debian/Ubuntu APT.
	ManagerApt PackageManagerType = "apt"
	// ManagerDpkg represents Debian DPKG.
	ManagerDpkg PackageManagerType = "dpkg"
	// ManagerSnap represents Canonical Snap.
	ManagerSnap PackageManagerType = "snap"
	// ManagerFlatpak represents Flatpak.
	ManagerFlatpak PackageManagerType = "flatpak"
	// ManagerWinget represents Windows Package Manager (winget).
	ManagerWinget PackageManagerType = "winget"
	// ManagerReg represents Windows Registry Uninstall entry.
	ManagerReg PackageManagerType = "registry"
	// ManagerNpm represents Global Node Package Manager.
	ManagerNpm PackageManagerType = "npm"
	// ManagerCustom represents a standalone custom installed binary.
	ManagerCustom PackageManagerType = "custom"
)

// IsValid reports whether the value matches a known variant.
func (p PackageManagerType) IsValid() bool {
	switch p {
	case ManagerApt, ManagerDpkg, ManagerSnap, ManagerFlatpak, ManagerWinget, ManagerReg, ManagerNpm, ManagerCustom:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (p PackageManagerType) String() string {
	return string(p)
}

// ParsePackageManagerType parses a raw string into a validated PackageManagerType.
func ParsePackageManagerType(s string) (PackageManagerType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllPackageManagerTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "managerapt":
		return ManagerApt, nil
	case "apt":
		return ManagerApt, nil
	case "managerdpkg":
		return ManagerDpkg, nil
	case "dpkg":
		return ManagerDpkg, nil
	case "managersnap":
		return ManagerSnap, nil
	case "snap":
		return ManagerSnap, nil
	case "managerflatpak":
		return ManagerFlatpak, nil
	case "flatpak":
		return ManagerFlatpak, nil
	case "managerwinget":
		return ManagerWinget, nil
	case "winget":
		return ManagerWinget, nil
	case "managerreg":
		return ManagerReg, nil
	case "reg":
		return ManagerReg, nil
	case "managernpm":
		return ManagerNpm, nil
	case "npm":
		return ManagerNpm, nil
	case "managercustom":
		return ManagerCustom, nil
	case "custom":
		return ManagerCustom, nil
	}
	return "", apperror.NewValidation("ParsePackageManagerType", "E1000", fmt.Sprintf("invalid packagemanagertype: %q", s)).WithContext("value", s)
}

// AllPackageManagerTypes returns all valid variants.
func AllPackageManagerTypes() []PackageManagerType {
	return []PackageManagerType{
		ManagerApt,
		ManagerDpkg,
		ManagerSnap,
		ManagerFlatpak,
		ManagerWinget,
		ManagerReg,
		ManagerNpm,
		ManagerCustom,
	}
}

// IsApt reports whether the value is ManagerApt.
func (p PackageManagerType) IsApt() bool {
	return p == ManagerApt
}

// IsDpkg reports whether the value is ManagerDpkg.
func (p PackageManagerType) IsDpkg() bool {
	return p == ManagerDpkg
}

// IsSnap reports whether the value is ManagerSnap.
func (p PackageManagerType) IsSnap() bool {
	return p == ManagerSnap
}

// IsFlatpak reports whether the value is ManagerFlatpak.
func (p PackageManagerType) IsFlatpak() bool {
	return p == ManagerFlatpak
}

// IsWinget reports whether the value is ManagerWinget.
func (p PackageManagerType) IsWinget() bool {
	return p == ManagerWinget
}

// IsReg reports whether the value is ManagerReg.
func (p PackageManagerType) IsReg() bool {
	return p == ManagerReg
}

// IsNpm reports whether the value is ManagerNpm.
func (p PackageManagerType) IsNpm() bool {
	return p == ManagerNpm
}

// IsCustom reports whether the value is ManagerCustom.
func (p PackageManagerType) IsCustom() bool {
	return p == ManagerCustom
}

// HasManager reports whether the value is non-empty and valid.
func (p PackageManagerType) HasManager() bool {
	return p.IsValid()
}
