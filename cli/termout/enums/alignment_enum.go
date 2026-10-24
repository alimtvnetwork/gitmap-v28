package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// AlignmentType defines column text alignment for terminal tables.
type AlignmentType string

const (
	// AlignLeft positions text flush with the left margin.
	AlignLeft AlignmentType = "left"
	// AlignRight positions text flush with the right margin (ideal for numeric metrics).
	AlignRight AlignmentType = "right"
	// AlignCenter centers text symmetrically within the column width.
	AlignCenter AlignmentType = "center"
)

// IsValid reports whether the value matches a known variant.
func (a AlignmentType) IsValid() bool {
	switch a {
	case AlignLeft, AlignRight, AlignCenter:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (a AlignmentType) String() string {
	return string(a)
}

// ParseAlignmentType parses a raw string into a validated AlignmentType.
func ParseAlignmentType(s string) (AlignmentType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
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
	case "alignright":
		return AlignRight, nil
	case "right":
		return AlignRight, nil
	case "aligncenter":
		return AlignCenter, nil
	case "center":
		return AlignCenter, nil
	}
	return "", apperror.NewValidation("ParseAlignmentType", "E1000", fmt.Sprintf("invalid alignmenttype: %q", s)).WithContext("value", s)
}

// AllAlignmentTypes returns all valid variants.
func AllAlignmentTypes() []AlignmentType {
	return []AlignmentType{
		AlignLeft,
		AlignRight,
		AlignCenter,
	}
}

// IsLeft reports whether the value is AlignLeft.
func (a AlignmentType) IsLeft() bool {
	return a == AlignLeft
}

// IsRight reports whether the value is AlignRight.
func (a AlignmentType) IsRight() bool {
	return a == AlignRight
}

// IsCenter reports whether the value is AlignCenter.
func (a AlignmentType) IsCenter() bool {
	return a == AlignCenter
}

// HasAlignment reports whether the value is non-empty and valid.
func (a AlignmentType) HasAlignment() bool {
	return a.IsValid()
}
