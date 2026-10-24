package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// AlignmentType defines the horizontal text alignment within the box.
type AlignmentType int

const (
	AlignLeft   AlignmentType = 0
	AlignCenter AlignmentType = 1
	AlignRight  AlignmentType = 2
)

// IsValid reports whether the value matches a known variant.
func (a AlignmentType) IsValid() bool {
	switch a {
	case AlignLeft, AlignCenter, AlignRight:
		return true
	default:
		return false
	}
}

// String returns the numeric representation of the value.
func (a AlignmentType) String() string {
	return strconv.Itoa(int(a))
}

// ParseAlignmentType parses a raw string into a validated AlignmentType.
func ParseAlignmentType(s string) (AlignmentType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	if n, err := strconv.Atoi(trimmed); err == nil {
		val := AlignmentType(n)
		if val.IsValid() {
			return val, nil
		}
	}
	for _, val := range AllAlignmentTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "alignleft":
		return AlignLeft, nil
	case "left":
		return AlignLeft, nil
	case "aligncenter":
		return AlignCenter, nil
	case "center":
		return AlignCenter, nil
	case "alignright":
		return AlignRight, nil
	case "right":
		return AlignRight, nil
	}
	return 0, apperror.NewValidation("ParseAlignmentType", "E1000", fmt.Sprintf("invalid alignmenttype: %q", s)).WithContext("value", s)
}

// AllAlignmentTypes returns all valid variants.
func AllAlignmentTypes() []AlignmentType {
	return []AlignmentType{
		AlignLeft,
		AlignCenter,
		AlignRight,
	}
}

// IsLeft reports whether the value is AlignLeft.
func (a AlignmentType) IsLeft() bool {
	return a == AlignLeft
}

// IsCenter reports whether the value is AlignCenter.
func (a AlignmentType) IsCenter() bool {
	return a == AlignCenter
}

// IsRight reports whether the value is AlignRight.
func (a AlignmentType) IsRight() bool {
	return a == AlignRight
}

// HasAlignment reports whether the value is non-empty and valid.
func (a AlignmentType) HasAlignment() bool {
	return a.IsValid()
}
