package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// PrettyModeType is the tri-state requested by the caller.
type PrettyModeType int

const (
	// PrettyAuto is the default: render when stdout is a real TTY and the
	// shared GITMAP_NO_PRETTY env opt-out is unset.
	PrettyAuto PrettyModeType = 0
	// PrettyOn forces ANSI rendering regardless of TTY / env. Useful when
	// piping into a renderer that understands ANSI (e.g. `less -R`).
	PrettyOn PrettyModeType = 1
	// PrettyOff suppresses ANSI rendering regardless of TTY / env. Use
	// when redirecting to a file or feeding the output into tools that
	// choke on escape codes (diff, sha256sum, grep -F).
	PrettyOff PrettyModeType = 2
)

// IsValid reports whether the value matches a known variant.
func (p PrettyModeType) IsValid() bool {
	switch p {
	case PrettyAuto, PrettyOn, PrettyOff:
		return true
	default:
		return false
	}
}

// String returns the numeric representation of the value.
func (p PrettyModeType) String() string {
	return strconv.Itoa(int(p))
}

// ParsePrettyModeType parses a raw string into a validated PrettyModeType.
func ParsePrettyModeType(s string) (PrettyModeType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	if n, err := strconv.Atoi(trimmed); err == nil {
		val := PrettyModeType(n)
		if val.IsValid() {
			return val, nil
		}
	}
	for _, val := range AllPrettyModeTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "prettyauto":
		return PrettyAuto, nil
	case "auto":
		return PrettyAuto, nil
	case "prettyon":
		return PrettyOn, nil
	case "on":
		return PrettyOn, nil
	case "prettyoff":
		return PrettyOff, nil
	case "off":
		return PrettyOff, nil
	}
	return 0, apperror.NewValidation("ParsePrettyModeType", "E1000", fmt.Sprintf("invalid prettymodetype: %q", s)).WithContext("value", s)
}

// AllPrettyModeTypes returns all valid variants.
func AllPrettyModeTypes() []PrettyModeType {
	return []PrettyModeType{
		PrettyAuto,
		PrettyOn,
		PrettyOff,
	}
}

// IsAuto reports whether the value is PrettyAuto.
func (p PrettyModeType) IsAuto() bool {
	return p == PrettyAuto
}

// IsOn reports whether the value is PrettyOn.
func (p PrettyModeType) IsOn() bool {
	return p == PrettyOn
}

// IsOff reports whether the value is PrettyOff.
func (p PrettyModeType) IsOff() bool {
	return p == PrettyOff
}

// HasMode reports whether the value is non-empty and valid.
func (p PrettyModeType) HasMode() bool {
	return p.IsValid()
}
