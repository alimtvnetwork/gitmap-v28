package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ThemeModeType specifies desktop visual appearance.
type ThemeModeType string

const (
	ThemeModeDark  ThemeModeType = "dark"
	ThemeModeLight ThemeModeType = "light"
)

// IsValid reports whether the theme mode matches a known variant.
func (t ThemeModeType) IsValid() bool {
	switch t {
	case ThemeModeDark, ThemeModeLight:
		return true
	default:
		return false
	}
}

// String returns the string representation of the theme mode.
func (t ThemeModeType) String() string {
	return string(t)
}

// ParseThemeModeType parses a raw string into a validated ThemeModeType.
func ParseThemeModeType(s string) (ThemeModeType, *apperror.AppError) {
	val := ThemeModeType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseThemeModeType", "E1000", fmt.Sprintf("invalid theme mode: %q", s)).WithContext("value", s)
}

// AllThemeModeTypes returns all valid ThemeModeType variants.
func AllThemeModeTypes() []ThemeModeType {
	return []ThemeModeType{
		ThemeModeDark,
		ThemeModeLight,
	}
}

// IsDark reports whether the theme mode is dark.
func (t ThemeModeType) IsDark() bool {
	return t == ThemeModeDark
}

// IsLight reports whether the theme mode is light.
func (t ThemeModeType) IsLight() bool {
	return t == ThemeModeLight
}

// HasThemeMode reports whether the theme mode is non-empty and valid.
func (t ThemeModeType) HasThemeMode() bool {
	return t.IsValid()
}
