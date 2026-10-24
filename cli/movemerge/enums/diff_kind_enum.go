package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// DiffKindType classifies a path across LEFT and RIGHT.
type DiffKindType int

const (
	// DiffMissingLeft = present on RIGHT only.
	DiffMissingLeft DiffKindType = 0
	// DiffMissingRight = present on LEFT only.
	DiffMissingRight DiffKindType = 1
	// DiffConflict = present on both with different content.
	DiffConflict DiffKindType = 2
	// DiffIdentical = present on both with byte-equal content.
	DiffIdentical DiffKindType = 3
)

// IsValid reports whether the value matches a known variant.
func (d DiffKindType) IsValid() bool {
	switch d {
	case DiffMissingLeft, DiffMissingRight, DiffConflict, DiffIdentical:
		return true
	default:
		return false
	}
}

// String returns the numeric representation of the value.
func (d DiffKindType) String() string {
	return strconv.Itoa(int(d))
}

// ParseDiffKindType parses a raw string into a validated DiffKindType.
func ParseDiffKindType(s string) (DiffKindType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	if n, err := strconv.Atoi(trimmed); err == nil {
		val := DiffKindType(n)
		if val.IsValid() {
			return val, nil
		}
	}
	for _, val := range AllDiffKindTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "diffmissingleft":
		return DiffMissingLeft, nil
	case "missingleft":
		return DiffMissingLeft, nil
	case "diffmissingright":
		return DiffMissingRight, nil
	case "missingright":
		return DiffMissingRight, nil
	case "diffconflict":
		return DiffConflict, nil
	case "conflict":
		return DiffConflict, nil
	case "diffidentical":
		return DiffIdentical, nil
	case "identical":
		return DiffIdentical, nil
	}
	return 0, apperror.NewValidation("ParseDiffKindType", "E1000", fmt.Sprintf("invalid diffkindtype: %q", s)).WithContext("value", s)
}

// AllDiffKindTypes returns all valid variants.
func AllDiffKindTypes() []DiffKindType {
	return []DiffKindType{
		DiffMissingLeft,
		DiffMissingRight,
		DiffConflict,
		DiffIdentical,
	}
}

// IsMissingLeft reports whether the value is DiffMissingLeft.
func (d DiffKindType) IsMissingLeft() bool {
	return d == DiffMissingLeft
}

// IsMissingRight reports whether the value is DiffMissingRight.
func (d DiffKindType) IsMissingRight() bool {
	return d == DiffMissingRight
}

// IsConflict reports whether the value is DiffConflict.
func (d DiffKindType) IsConflict() bool {
	return d == DiffConflict
}

// IsIdentical reports whether the value is DiffIdentical.
func (d DiffKindType) IsIdentical() bool {
	return d == DiffIdentical
}

// HasKind reports whether the value is non-empty and valid.
func (d DiffKindType) HasKind() bool {
	return d.IsValid()
}
