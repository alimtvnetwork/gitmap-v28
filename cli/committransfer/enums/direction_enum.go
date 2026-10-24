package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// DirectionType names which side receives the new commits.
type DirectionType int

const (
	// DirRight replays LEFT → RIGHT (writes commits on RIGHT).
	DirRight DirectionType = 0
	// DirLeft replays RIGHT → LEFT (writes commits on LEFT).
	DirLeft DirectionType = 1
	// DirBoth replays in both directions sequentially: LEFT → RIGHT
	// first, then RIGHT → LEFT. Each pass builds its own plan and
	// shares the same Options (with a directional log-prefix suffix).
	DirBoth DirectionType = 2
)

// IsValid reports whether the value matches a known variant.
func (d DirectionType) IsValid() bool {
	switch d {
	case DirRight, DirLeft, DirBoth:
		return true
	default:
		return false
	}
}

// String returns the numeric representation of the value.
func (d DirectionType) String() string {
	return strconv.Itoa(int(d))
}

// ParseDirectionType parses a raw string into a validated DirectionType.
func ParseDirectionType(s string) (DirectionType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	if n, err := strconv.Atoi(trimmed); err == nil {
		val := DirectionType(n)
		if val.IsValid() {
			return val, nil
		}
	}
	for _, val := range AllDirectionTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "dirright":
		return DirRight, nil
	case "right":
		return DirRight, nil
	case "dirleft":
		return DirLeft, nil
	case "left":
		return DirLeft, nil
	case "dirboth":
		return DirBoth, nil
	case "both":
		return DirBoth, nil
	}
	return 0, apperror.NewValidation("ParseDirectionType", "E1000", fmt.Sprintf("invalid directiontype: %q", s)).WithContext("value", s)
}

// AllDirectionTypes returns all valid variants.
func AllDirectionTypes() []DirectionType {
	return []DirectionType{
		DirRight,
		DirLeft,
		DirBoth,
	}
}

// IsRight reports whether the value is DirRight.
func (d DirectionType) IsRight() bool {
	return d == DirRight
}

// IsLeft reports whether the value is DirLeft.
func (d DirectionType) IsLeft() bool {
	return d == DirLeft
}

// IsBoth reports whether the value is DirBoth.
func (d DirectionType) IsBoth() bool {
	return d == DirBoth
}

// HasDirection reports whether the value is non-empty and valid.
func (d DirectionType) HasDirection() bool {
	return d.IsValid()
}
