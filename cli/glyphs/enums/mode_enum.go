package enums

import (
	"strconv"
)

// ModeType is the resolved glyph selection.
type ModeType int

const (
	// ModeRich passes UTF-8 emoji through unchanged.
	ModeRich ModeType = 0
	// ModeSafe rewrites emoji to ASCII fallbacks.
	ModeSafe ModeType = 1
)

// IsValid reports whether the value matches a known variant.
func (m ModeType) IsValid() bool {
	switch m {
	case ModeRich, ModeSafe:
		return true
	default:
		return false
	}
}

// String returns the numeric representation of the value.
func (m ModeType) String() string {
	return strconv.Itoa(int(m))
}

// AllModeTypes returns all valid variants.
func AllModeTypes() []ModeType {
	return []ModeType{
		ModeRich,
		ModeSafe,
	}
}

// IsRich reports whether the value is ModeRich.
func (m ModeType) IsRich() bool {
	return m == ModeRich
}

// IsSafe reports whether the value is ModeSafe.
func (m ModeType) IsSafe() bool {
	return m == ModeSafe
}

// HasMode reports whether the value is non-empty and valid.
func (m ModeType) HasMode() bool {
	return m.IsValid()
}
