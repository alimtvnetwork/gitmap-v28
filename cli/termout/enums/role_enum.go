package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// RoleType is a semantic terminal color role. Help and UI code references
// these — never raw hex or literal SGR codes.
type RoleType string

const (
	RoleFG       RoleType = "fg"
	RoleFGMuted  RoleType = "fg-muted"
	RoleFGSubtle RoleType = "fg-subtle"
	RolePrimary  RoleType = "primary"
	RoleAccent   RoleType = "accent"
	RoleWarn     RoleType = "warn"
	RoleError    RoleType = "error"
	RoleOK       RoleType = "ok"
)

// IsValid reports whether the value matches a known variant.
func (r RoleType) IsValid() bool {
	switch r {
	case RoleFG, RoleFGMuted, RoleFGSubtle, RolePrimary, RoleAccent, RoleWarn, RoleError, RoleOK:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (r RoleType) String() string {
	return string(r)
}

// ParseRoleType parses a raw string into a validated RoleType.
func ParseRoleType(s string) (RoleType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllRoleTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "rolefg":
		return RoleFG, nil
	case "fg":
		return RoleFG, nil
	case "rolefgmuted":
		return RoleFGMuted, nil
	case "fgmuted":
		return RoleFGMuted, nil
	case "rolefgsubtle":
		return RoleFGSubtle, nil
	case "fgsubtle":
		return RoleFGSubtle, nil
	case "roleprimary":
		return RolePrimary, nil
	case "primary":
		return RolePrimary, nil
	case "roleaccent":
		return RoleAccent, nil
	case "accent":
		return RoleAccent, nil
	case "rolewarn":
		return RoleWarn, nil
	case "warn":
		return RoleWarn, nil
	case "roleerror":
		return RoleError, nil
	case "error":
		return RoleError, nil
	case "roleok":
		return RoleOK, nil
	case "ok":
		return RoleOK, nil
	}
	return "", apperror.NewValidation("ParseRoleType", "E1000", fmt.Sprintf("invalid roletype: %q", s)).WithContext("value", s)
}

// AllRoleTypes returns all valid variants.
func AllRoleTypes() []RoleType {
	return []RoleType{
		RoleFG,
		RoleFGMuted,
		RoleFGSubtle,
		RolePrimary,
		RoleAccent,
		RoleWarn,
		RoleError,
		RoleOK,
	}
}

// IsFG reports whether the value is RoleFG.
func (r RoleType) IsFG() bool {
	return r == RoleFG
}

// IsFGMuted reports whether the value is RoleFGMuted.
func (r RoleType) IsFGMuted() bool {
	return r == RoleFGMuted
}

// IsFGSubtle reports whether the value is RoleFGSubtle.
func (r RoleType) IsFGSubtle() bool {
	return r == RoleFGSubtle
}

// IsPrimary reports whether the value is RolePrimary.
func (r RoleType) IsPrimary() bool {
	return r == RolePrimary
}

// IsAccent reports whether the value is RoleAccent.
func (r RoleType) IsAccent() bool {
	return r == RoleAccent
}

// IsWarn reports whether the value is RoleWarn.
func (r RoleType) IsWarn() bool {
	return r == RoleWarn
}

// IsError reports whether the value is RoleError.
func (r RoleType) IsError() bool {
	return r == RoleError
}

// IsOK reports whether the value is RoleOK.
func (r RoleType) IsOK() bool {
	return r == RoleOK
}

// HasRole reports whether the value is non-empty and valid.
func (r RoleType) HasRole() bool {
	return r.IsValid()
}
