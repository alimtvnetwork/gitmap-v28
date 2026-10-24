package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// DriftActionType answers: given the *current* visibility we just read
// from the provider, the *expected* visibility we persisted as
// NewVisibility on the original run, and the caller's --force flag,
// should reverseOneRepo proceed, force-override, or skip-with-drift?
//
//	action == "proceed"  → no drift, apply PrevVisibility normally.
//	action == "force"    → --force set, log override line then apply.
//	action == "skip"     → drift detected, leave repo untouched.
type DriftActionType string

const (
	DriftActionTypeProceed DriftActionType = "proceed"
	DriftActionTypeForce   DriftActionType = "force"
	DriftActionTypeSkip    DriftActionType = "skip"
)

// IsValid reports whether the drift action matches a known action.
func (d DriftActionType) IsValid() bool {
	switch d {
	case DriftActionTypeProceed, DriftActionTypeForce, DriftActionTypeSkip:
		return true
	default:
		return false
	}
}

// String returns the string representation of the drift action.
func (d DriftActionType) String() string {
	return string(d)
}

// ParseDriftActionType parses a raw string into a validated DriftActionType.
func ParseDriftActionType(s string) (DriftActionType, *apperror.AppError) {
	val := DriftActionType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseDriftActionType", "E1000", fmt.Sprintf("invalid drift action: %q", s)).WithContext("value", s)
}

// AllDriftActionTypes returns all valid DriftActionType variants.
func AllDriftActionTypes() []DriftActionType {
	return []DriftActionType{DriftActionTypeProceed, DriftActionTypeForce, DriftActionTypeSkip}
}

// IsProceed reports whether the action is proceed.
func (d DriftActionType) IsProceed() bool { return d == DriftActionTypeProceed }

// IsForce reports whether the action is force.
func (d DriftActionType) IsForce() bool { return d == DriftActionTypeForce }

// IsSkip reports whether the action is skip.
func (d DriftActionType) IsSkip() bool { return d == DriftActionTypeSkip }
