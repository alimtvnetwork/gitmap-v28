package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// OutputMode selects the rendering for Report / FailWith. Zero value
// (OutputHuman) keeps the existing terminal-friendly shape so callers
// that don't care about JSON pay no migration cost.
type OutputModeType int

const (
	// OutputHuman renders the canonical text line + indented
	// context fields. This is the default.
	OutputHuman OutputModeType = 0
	// OutputJSON renders one single-line JSON object per failure.
	OutputJSON OutputModeType = 1
)

// IsValid reports whether the value matches a known variant.
func (o OutputModeType) IsValid() bool {
	switch o {
	case OutputHuman, OutputJSON:
		return true
	default:
		return false
	}
}

// String returns the numeric representation of the value.
func (o OutputModeType) String() string {
	return strconv.Itoa(int(o))
}

// ParseOutputModeType parses a raw string into a validated OutputModeType.
func ParseOutputModeType(s string) (OutputModeType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	if n, err := strconv.Atoi(trimmed); err == nil {
		val := OutputModeType(n)
		if val.IsValid() {
			return val, nil
		}
	}
	for _, val := range AllOutputModeTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "outputhuman":
		return OutputHuman, nil
	case "human":
		return OutputHuman, nil
	case "outputjson":
		return OutputJSON, nil
	case "json":
		return OutputJSON, nil
	}
	return 0, apperror.NewValidation("ParseOutputModeType", "E1000", fmt.Sprintf("invalid outputmodetype: %q", s)).WithContext("value", s)
}

// AllOutputModeTypes returns all valid variants.
func AllOutputModeTypes() []OutputModeType {
	return []OutputModeType{
		OutputHuman,
		OutputJSON,
	}
}

// IsHuman reports whether the value is OutputHuman.
func (o OutputModeType) IsHuman() bool {
	return o == OutputHuman
}

// IsJSON reports whether the value is OutputJSON.
func (o OutputModeType) IsJSON() bool {
	return o == OutputJSON
}

// HasMode reports whether the value is non-empty and valid.
func (o OutputModeType) HasMode() bool {
	return o.IsValid()
}
