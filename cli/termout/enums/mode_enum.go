package enums

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ModeType is the resolved palette selection.
type ModeType int

const (
	// ModeBright is the loud bright+bold palette (default).
	ModeBright ModeType = 0
	// ModeStandard is the muted pre-v5.13 plain palette.
	ModeStandard ModeType = 1
	// ModeMono strips all ANSI SGR codes.
	ModeMono ModeType = 2
)

// IsValid reports whether the value matches a known variant.
func (m ModeType) IsValid() bool {
	switch m {
	case ModeBright, ModeStandard, ModeMono:
		return true
	default:
		return false
	}
}

func (m ModeType) String() string {
	switch m {
	case ModeStandard:
		return constants.ThemeStandard
	case ModeMono:
		return constants.ThemeMonochrome
	case ModeBright:
		return constants.ThemeBright
	default:
		return constants.ThemeBright
	}
}

// AllModeTypes returns all valid variants.
func AllModeTypes() []ModeType {
	return []ModeType{
		ModeBright,
		ModeStandard,
		ModeMono,
	}
}

// IsBright reports whether the value is ModeBright.
func (m ModeType) IsBright() bool {
	return m == ModeBright
}

// IsStandard reports whether the value is ModeStandard.
func (m ModeType) IsStandard() bool {
	return m == ModeStandard
}

// IsMono reports whether the value is ModeMono.
func (m ModeType) IsMono() bool {
	return m == ModeMono
}

// HasMode reports whether the value is non-empty and valid.
func (m ModeType) HasMode() bool {
	return m.IsValid()
}
