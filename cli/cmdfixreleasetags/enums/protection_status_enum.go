package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// ProtectionStatusType indicates whether a tag is immune to deletion.
type ProtectionStatusType string

const (
	StatusEligible        ProtectionStatusType = "ELIGIBLE"
	StatusProtectedActive ProtectionStatusType = "PROTECTED_ACTIVE_VERSION"
	StatusProtectedLatest ProtectionStatusType = "PROTECTED_LATEST_HEALTHY"
	StatusProtectedGrace  ProtectionStatusType = "PROTECTED_GRACE_WINDOW"
)

// IsValid reports whether the value matches a known variant.
func (p ProtectionStatusType) IsValid() bool {
	switch p {
	case StatusEligible, StatusProtectedActive, StatusProtectedLatest, StatusProtectedGrace:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (p ProtectionStatusType) String() string {
	return string(p)
}

// ParseProtectionStatusType parses a raw string into a validated ProtectionStatusType.
func ParseProtectionStatusType(s string) (ProtectionStatusType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllProtectionStatusTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "statuseligible":
		return StatusEligible, nil
	case "eligible":
		return StatusEligible, nil
	case "statusprotectedactive":
		return StatusProtectedActive, nil
	case "protectedactive":
		return StatusProtectedActive, nil
	case "statusprotectedlatest":
		return StatusProtectedLatest, nil
	case "protectedlatest":
		return StatusProtectedLatest, nil
	case "statusprotectedgrace":
		return StatusProtectedGrace, nil
	case "protectedgrace":
		return StatusProtectedGrace, nil
	}
	return "", apperror.NewValidation("ParseProtectionStatusType", "E1000", fmt.Sprintf("invalid protectionstatustype: %q", s)).WithContext("value", s)
}

// AllProtectionStatusTypes returns all valid variants.
func AllProtectionStatusTypes() []ProtectionStatusType {
	return []ProtectionStatusType{
		StatusEligible,
		StatusProtectedActive,
		StatusProtectedLatest,
		StatusProtectedGrace,
	}
}

// IsEligible reports whether the value is StatusEligible.
func (p ProtectionStatusType) IsEligible() bool {
	return p == StatusEligible
}

// IsProtectedActive reports whether the value is StatusProtectedActive.
func (p ProtectionStatusType) IsProtectedActive() bool {
	return p == StatusProtectedActive
}

// IsProtectedLatest reports whether the value is StatusProtectedLatest.
func (p ProtectionStatusType) IsProtectedLatest() bool {
	return p == StatusProtectedLatest
}

// IsProtectedGrace reports whether the value is StatusProtectedGrace.
func (p ProtectionStatusType) IsProtectedGrace() bool {
	return p == StatusProtectedGrace
}

// HasStatus reports whether the value is non-empty and valid.
func (p ProtectionStatusType) HasStatus() bool {
	return p.IsValid()
}
