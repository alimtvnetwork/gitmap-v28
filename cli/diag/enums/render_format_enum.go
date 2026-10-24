package enums

import (
	"fmt"
	"strings"
)

// RenderFormatType defines the rendering style for suggestions.
type RenderFormatType string

const (
	// RenderFormatBox renders an elegant Catppuccin terminal box.
	RenderFormatBox RenderFormatType = "box"
	// RenderFormatCompact renders a minimal single-line inline hint.
	RenderFormatCompact RenderFormatType = "compact"
	// RenderFormatJSON renders machine-readable JSON for agents and tooling.
	RenderFormatJSON RenderFormatType = "json"
)

// IsValid reports whether the value matches a known variant.
func (r RenderFormatType) IsValid() bool {
	switch r {
	case RenderFormatBox, RenderFormatCompact, RenderFormatJSON:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (r RenderFormatType) String() string {
	return string(r)
}

// ParseRenderFormatType parses a raw string into a validated RenderFormatType.
// ParseRenderFormatType returns a plain error (not *apperror.AppError)
// by recorded exemption: apperror imports cli/diag, so an AppError
// Parse here would close an import cycle (diag -> diag/enums -> apperror -> diag).
func ParseRenderFormatType(s string) (RenderFormatType, error) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllRenderFormatTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "renderformatbox":
		return RenderFormatBox, nil
	case "box":
		return RenderFormatBox, nil
	case "renderformatcompact":
		return RenderFormatCompact, nil
	case "compact":
		return RenderFormatCompact, nil
	case "renderformatjson":
		return RenderFormatJSON, nil
	case "json":
		return RenderFormatJSON, nil
	}
	return "", fmt.Errorf("invalid renderformattype: %q", s)
}

// AllRenderFormatTypes returns all valid variants.
func AllRenderFormatTypes() []RenderFormatType {
	return []RenderFormatType{
		RenderFormatBox,
		RenderFormatCompact,
		RenderFormatJSON,
	}
}

// IsBox reports whether the value is RenderFormatBox.
func (r RenderFormatType) IsBox() bool {
	return r == RenderFormatBox
}

// IsCompact reports whether the value is RenderFormatCompact.
func (r RenderFormatType) IsCompact() bool {
	return r == RenderFormatCompact
}

// IsJSON reports whether the value is RenderFormatJSON.
func (r RenderFormatType) IsJSON() bool {
	return r == RenderFormatJSON
}

// HasFormat reports whether the value is non-empty and valid.
func (r RenderFormatType) HasFormat() bool {
	return r.IsValid()
}
