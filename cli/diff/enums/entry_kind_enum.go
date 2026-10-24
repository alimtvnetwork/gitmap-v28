package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// EntryKindType classifies how a single relative path differs.
type EntryKindType int

const (
	// MissingRight = present on LEFT, absent on RIGHT.
	MissingRight EntryKindType = 0
	// MissingLeft = present on RIGHT, absent on LEFT.
	MissingLeft EntryKindType = 1
	// Identical = present on both, byte-equal.
	Identical EntryKindType = 2
	// Conflict = present on both, byte-different.
	Conflict EntryKindType = 3
)

// IsValid reports whether the value matches a known variant.
func (e EntryKindType) IsValid() bool {
	switch e {
	case MissingRight, MissingLeft, Identical, Conflict:
		return true
	default:
		return false
	}
}

// String returns the numeric representation of the value.
func (e EntryKindType) String() string {
	return strconv.Itoa(int(e))
}

// ParseEntryKindType parses a raw string into a validated EntryKindType.
func ParseEntryKindType(s string) (EntryKindType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	if n, err := strconv.Atoi(trimmed); err == nil {
		val := EntryKindType(n)
		if val.IsValid() {
			return val, nil
		}
	}
	for _, val := range AllEntryKindTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "missingright":
		return MissingRight, nil
	case "missingleft":
		return MissingLeft, nil
	case "identical":
		return Identical, nil
	case "conflict":
		return Conflict, nil
	}
	return 0, apperror.NewValidation("ParseEntryKindType", "E1000", fmt.Sprintf("invalid entrykindtype: %q", s)).WithContext("value", s)
}

// AllEntryKindTypes returns all valid variants.
func AllEntryKindTypes() []EntryKindType {
	return []EntryKindType{
		MissingRight,
		MissingLeft,
		Identical,
		Conflict,
	}
}

// IsMissingRight reports whether the value is MissingRight.
func (e EntryKindType) IsMissingRight() bool {
	return e == MissingRight
}

// IsMissingLeft reports whether the value is MissingLeft.
func (e EntryKindType) IsMissingLeft() bool {
	return e == MissingLeft
}

// IsIdentical reports whether the value is Identical.
func (e EntryKindType) IsIdentical() bool {
	return e == Identical
}

// IsConflict reports whether the value is Conflict.
func (e EntryKindType) IsConflict() bool {
	return e == Conflict
}

// HasKind reports whether the value is non-empty and valid.
func (e EntryKindType) HasKind() bool {
	return e.IsValid()
}
