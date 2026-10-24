package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// PreferPolicyType mirrors movemerge.PreferPolicyType for file-level conflicts
// during the snapshot-copy step. We re-export the values rather than
// importing the type to keep the package decoupled in tests.
type PreferPolicyType int

const (
	// PreferNone means use the existing file (no overwrite on conflict).
	PreferNone PreferPolicyType = 0
	// PreferSource overwrites target with source on every conflict.
	PreferSource PreferPolicyType = 1
	// PreferTarget keeps target on every conflict (effectively skip).
	PreferTarget PreferPolicyType = 2
)

// IsValid reports whether the value matches a known variant.
func (p PreferPolicyType) IsValid() bool {
	switch p {
	case PreferNone, PreferSource, PreferTarget:
		return true
	default:
		return false
	}
}

// String returns the numeric representation of the value.
func (p PreferPolicyType) String() string {
	return strconv.Itoa(int(p))
}

// ParsePreferPolicyType parses a raw string into a validated PreferPolicyType.
func ParsePreferPolicyType(s string) (PreferPolicyType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	if n, err := strconv.Atoi(trimmed); err == nil {
		val := PreferPolicyType(n)
		if val.IsValid() {
			return val, nil
		}
	}
	for _, val := range AllPreferPolicyTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "prefernone":
		return PreferNone, nil
	case "none":
		return PreferNone, nil
	case "prefersource":
		return PreferSource, nil
	case "source":
		return PreferSource, nil
	case "prefertarget":
		return PreferTarget, nil
	case "target":
		return PreferTarget, nil
	}
	return 0, apperror.NewValidation("ParsePreferPolicyType", "E1000", fmt.Sprintf("invalid preferpolicytype: %q", s)).WithContext("value", s)
}

// AllPreferPolicyTypes returns all valid variants.
func AllPreferPolicyTypes() []PreferPolicyType {
	return []PreferPolicyType{
		PreferNone,
		PreferSource,
		PreferTarget,
	}
}

// IsNone reports whether the value is PreferNone.
func (p PreferPolicyType) IsNone() bool {
	return p == PreferNone
}

// IsSource reports whether the value is PreferSource.
func (p PreferPolicyType) IsSource() bool {
	return p == PreferSource
}

// IsTarget reports whether the value is PreferTarget.
func (p PreferPolicyType) IsTarget() bool {
	return p == PreferTarget
}

// HasPolicy reports whether the value is non-empty and valid.
func (p PreferPolicyType) HasPolicy() bool {
	return p.IsValid()
}
