package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// PreferPolicyType is how -y / --prefer-* resolve conflicts non-interactively.
type PreferPolicyType int

const (
	// PreferNone means use the interactive prompt.
	PreferNone PreferPolicyType = 0
	// PreferLeft makes LEFT always win.
	PreferLeft PreferPolicyType = 1
	// PreferRight makes RIGHT always win.
	PreferRight PreferPolicyType = 2
	// PreferNewer compares mtime; newer side wins.
	PreferNewer PreferPolicyType = 3
	// PreferSkip skips every conflict (only missing files copied).
	PreferSkip PreferPolicyType = 4
)

// IsValid reports whether the value matches a known variant.
func (p PreferPolicyType) IsValid() bool {
	switch p {
	case PreferNone, PreferLeft, PreferRight, PreferNewer, PreferSkip:
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
	case "preferleft":
		return PreferLeft, nil
	case "left":
		return PreferLeft, nil
	case "preferright":
		return PreferRight, nil
	case "right":
		return PreferRight, nil
	case "prefernewer":
		return PreferNewer, nil
	case "newer":
		return PreferNewer, nil
	case "preferskip":
		return PreferSkip, nil
	case "skip":
		return PreferSkip, nil
	}
	return 0, apperror.NewValidation("ParsePreferPolicyType", "E1000", fmt.Sprintf("invalid preferpolicytype: %q", s)).WithContext("value", s)
}

// AllPreferPolicyTypes returns all valid variants.
func AllPreferPolicyTypes() []PreferPolicyType {
	return []PreferPolicyType{
		PreferNone,
		PreferLeft,
		PreferRight,
		PreferNewer,
		PreferSkip,
	}
}

// IsNone reports whether the value is PreferNone.
func (p PreferPolicyType) IsNone() bool {
	return p == PreferNone
}

// IsLeft reports whether the value is PreferLeft.
func (p PreferPolicyType) IsLeft() bool {
	return p == PreferLeft
}

// IsRight reports whether the value is PreferRight.
func (p PreferPolicyType) IsRight() bool {
	return p == PreferRight
}

// IsNewer reports whether the value is PreferNewer.
func (p PreferPolicyType) IsNewer() bool {
	return p == PreferNewer
}

// IsSkip reports whether the value is PreferSkip.
func (p PreferPolicyType) IsSkip() bool {
	return p == PreferSkip
}

// HasPolicy reports whether the value is non-empty and valid.
func (p PreferPolicyType) HasPolicy() bool {
	return p.IsValid()
}
