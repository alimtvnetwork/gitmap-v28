package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// DirectionType selects which side(s) the operation writes to.
type DirectionType int

const (
	// DirBoth writes into both sides (merge-both).
	DirBoth DirectionType = 0
	// DirLeftOnly writes only into LEFT (merge-left).
	DirLeftOnly DirectionType = 1
	// DirRightOnly writes only into RIGHT (merge-right).
	DirRightOnly DirectionType = 2
)

// IsValid reports whether the value matches a known variant.
func (d DirectionType) IsValid() bool {
	switch d {
	case DirBoth, DirLeftOnly, DirRightOnly:
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
	case "dirboth":
		return DirBoth, nil
	case "both":
		return DirBoth, nil
	case "dirleftonly":
		return DirLeftOnly, nil
	case "leftonly":
		return DirLeftOnly, nil
	case "dirrightonly":
		return DirRightOnly, nil
	case "rightonly":
		return DirRightOnly, nil
	}
	return 0, apperror.NewValidation("ParseDirectionType", "E1000", fmt.Sprintf("invalid directiontype: %q", s)).WithContext("value", s)
}

// AllDirectionTypes returns all valid variants.
func AllDirectionTypes() []DirectionType {
	return []DirectionType{
		DirBoth,
		DirLeftOnly,
		DirRightOnly,
	}
}

// IsBoth reports whether the value is DirBoth.
func (d DirectionType) IsBoth() bool {
	return d == DirBoth
}

// IsLeftOnly reports whether the value is DirLeftOnly.
func (d DirectionType) IsLeftOnly() bool {
	return d == DirLeftOnly
}

// IsRightOnly reports whether the value is DirRightOnly.
func (d DirectionType) IsRightOnly() bool {
	return d == DirRightOnly
}

// HasDirection reports whether the value is non-empty and valid.
func (d DirectionType) HasDirection() bool {
	return d.IsValid()
}
