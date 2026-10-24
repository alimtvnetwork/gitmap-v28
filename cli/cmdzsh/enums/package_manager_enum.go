package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// PackageManagerType identifies supported system package managers.
type PackageManagerType string

const (
	PkgMgrApt     PackageManagerType = "apt"
	PkgMgrDnf     PackageManagerType = "dnf"
	PkgMgrBrew    PackageManagerType = "brew"
	PkgMgrPacman  PackageManagerType = "pacman"
	PkgMgrUnknown PackageManagerType = "unknown"
)

// IsValid reports whether the value matches a known variant.
func (p PackageManagerType) IsValid() bool {
	switch p {
	case PkgMgrApt, PkgMgrDnf, PkgMgrBrew, PkgMgrPacman, PkgMgrUnknown:
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
	case "pkgmgrapt":
		return PkgMgrApt, nil
	case "apt":
		return PkgMgrApt, nil
	case "pkgmgrdnf":
		return PkgMgrDnf, nil
	case "dnf":
		return PkgMgrDnf, nil
	case "pkgmgrbrew":
		return PkgMgrBrew, nil
	case "brew":
		return PkgMgrBrew, nil
	case "pkgmgrpacman":
		return PkgMgrPacman, nil
	case "pacman":
		return PkgMgrPacman, nil
	case "pkgmgrunknown":
		return PkgMgrUnknown, nil
	case "unknown":
		return PkgMgrUnknown, nil
	}
	return "", apperror.NewValidation("ParsePackageManagerType", "E1000", fmt.Sprintf("invalid packagemanagertype: %q", s)).WithContext("value", s)
}

// AllPackageManagerTypes returns all valid variants.
func AllPackageManagerTypes() []PackageManagerType {
	return []PackageManagerType{
		PkgMgrApt,
		PkgMgrDnf,
		PkgMgrBrew,
		PkgMgrPacman,
		PkgMgrUnknown,
	}
}

// IsApt reports whether the value is PkgMgrApt.
func (p PackageManagerType) IsApt() bool {
	return p == PkgMgrApt
}

// IsDnf reports whether the value is PkgMgrDnf.
func (p PackageManagerType) IsDnf() bool {
	return p == PkgMgrDnf
}

// IsBrew reports whether the value is PkgMgrBrew.
func (p PackageManagerType) IsBrew() bool {
	return p == PkgMgrBrew
}

// IsPacman reports whether the value is PkgMgrPacman.
func (p PackageManagerType) IsPacman() bool {
	return p == PkgMgrPacman
}

// IsUnknown reports whether the value is PkgMgrUnknown.
func (p PackageManagerType) IsUnknown() bool {
	return p == PkgMgrUnknown
}

// HasManager reports whether the value is non-empty and valid.
func (p PackageManagerType) HasManager() bool {
	return p.IsValid()
}
