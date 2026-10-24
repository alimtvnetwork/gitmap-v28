package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// DMType enumerates supported display managers.
type DMType string

const (
	DMTypeGDM3    DMType = "gdm3"
	DMTypeLightDM DMType = "lightdm"
	DMTypeSDDM    DMType = "sddm"
	DMTypeUnknown DMType = "unknown"
)

// IsValid reports whether the display manager matches a known variant.
func (d DMType) IsValid() bool {
	switch d {
	case DMTypeGDM3, DMTypeLightDM, DMTypeSDDM, DMTypeUnknown:
		return true
	default:
		return false
	}
}

// String returns the string representation of the display manager.
func (d DMType) String() string {
	return string(d)
}

// ParseDMType parses a raw string into a validated DMType.
func ParseDMType(s string) (DMType, *apperror.AppError) {
	val := DMType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseDMType", "E1000", fmt.Sprintf("invalid display manager: %q", s)).WithContext("value", s)
}

// AllDMTypes returns all valid DMType variants.
func AllDMTypes() []DMType {
	return []DMType{
		DMTypeGDM3,
		DMTypeLightDM,
		DMTypeSDDM,
		DMTypeUnknown,
	}
}

// IsGDM3 reports whether the display manager is gdm3.
func (d DMType) IsGDM3() bool {
	return d == DMTypeGDM3
}

// IsLightDM reports whether the display manager is lightdm.
func (d DMType) IsLightDM() bool {
	return d == DMTypeLightDM
}

// IsSDDM reports whether the display manager is sddm.
func (d DMType) IsSDDM() bool {
	return d == DMTypeSDDM
}

// IsUnknown reports whether the display manager is unknown.
func (d DMType) IsUnknown() bool {
	return d == DMTypeUnknown
}

// HasDMType reports whether the display manager is non-empty and valid.
func (d DMType) HasDMType() bool {
	return d.IsValid()
}
