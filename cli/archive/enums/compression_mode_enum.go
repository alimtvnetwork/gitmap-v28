package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"strings"
)

// CompressionModeType is the user-facing knob persisted in ArchiveHistory.CompressionMode.
type CompressionModeType string

const (
	ModeStandard CompressionModeType = constants.CompressionStandard
	ModeBest     CompressionModeType = constants.CompressionBest
	ModeFast     CompressionModeType = constants.CompressionFast
)

// IsValid reports whether the value matches a known variant.
func (c CompressionModeType) IsValid() bool {
	switch c {
	case ModeStandard, ModeBest, ModeFast:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (c CompressionModeType) String() string {
	return string(c)
}

// ParseCompressionModeType parses a raw string into a validated CompressionModeType.
func ParseCompressionModeType(s string) (CompressionModeType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllCompressionModeTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "modestandard":
		return ModeStandard, nil
	case "standard":
		return ModeStandard, nil
	case "modebest":
		return ModeBest, nil
	case "best":
		return ModeBest, nil
	case "modefast":
		return ModeFast, nil
	case "fast":
		return ModeFast, nil
	}
	return "", apperror.NewValidation("ParseCompressionModeType", "E1000", fmt.Sprintf("invalid compressionmodetype: %q", s)).WithContext("value", s)
}

// AllCompressionModeTypes returns all valid variants.
func AllCompressionModeTypes() []CompressionModeType {
	return []CompressionModeType{
		ModeStandard,
		ModeBest,
		ModeFast,
	}
}

// IsStandard reports whether the value is ModeStandard.
func (c CompressionModeType) IsStandard() bool {
	return c == ModeStandard
}

// IsBest reports whether the value is ModeBest.
func (c CompressionModeType) IsBest() bool {
	return c == ModeBest
}

// IsFast reports whether the value is ModeFast.
func (c CompressionModeType) IsFast() bool {
	return c == ModeFast
}

// HasMode reports whether the value is non-empty and valid.
func (c CompressionModeType) HasMode() bool {
	return c.IsValid()
}
