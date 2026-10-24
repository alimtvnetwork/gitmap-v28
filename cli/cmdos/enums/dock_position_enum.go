package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// DockPositionType represents cardinal screen edge orientation for desktop dock.
type DockPositionType string

// DockPosition is an alias for DockPositionType.
type DockPosition = DockPositionType

const (
	DockPositionBottom DockPositionType = "bottom"
	DockPositionLeft   DockPositionType = "left"
	DockPositionRight  DockPositionType = "right"
	DockPositionTop    DockPositionType = "top"
	DockPositionCenter DockPositionType = "center"
)

// dockPositionAliases maps shorthand and directional spellings to canonical
// dock positions, preserving the aliases previously handled by the cmdos
// NormalizeDockPosition switch.
var dockPositionAliases = map[string]DockPositionType{
	"bottom": DockPositionBottom,
	"b":      DockPositionBottom,
	"down":   DockPositionBottom,
	"left":   DockPositionLeft,
	"l":      DockPositionLeft,
	"right":  DockPositionRight,
	"r":      DockPositionRight,
	"top":    DockPositionTop,
	"t":      DockPositionTop,
	"up":     DockPositionTop,
	"center": DockPositionCenter,
	"c":      DockPositionCenter,
}

// IsValid reports whether the dock position matches a known variant.
func (d DockPositionType) IsValid() bool {
	switch d {
	case DockPositionBottom, DockPositionLeft, DockPositionRight,
		DockPositionTop, DockPositionCenter:
		return true
	default:
		return false
	}
}

// String returns the string representation of the dock position.
func (d DockPositionType) String() string {
	return string(d)
}

// ParseDockPositionType parses a raw string into a validated DockPositionType,
// accepting canonical positions and their shorthand aliases.
func ParseDockPositionType(s string) (DockPositionType, *apperror.AppError) {
	if val, ok := dockPositionAliases[strings.ToLower(strings.TrimSpace(s))]; ok {
		return val, nil
	}

	return "", apperror.NewValidation("ParseDockPositionType", "E1000", fmt.Sprintf("invalid dock position: %q", s)).WithContext("value", s)
}

// AllDockPositionTypes returns all valid DockPositionType variants.
func AllDockPositionTypes() []DockPositionType {
	return []DockPositionType{
		DockPositionBottom,
		DockPositionLeft,
		DockPositionRight,
		DockPositionTop,
		DockPositionCenter,
	}
}

// IsBottom reports whether the dock position is bottom.
func (d DockPositionType) IsBottom() bool {
	return d == DockPositionBottom
}

// IsLeft reports whether the dock position is left.
func (d DockPositionType) IsLeft() bool {
	return d == DockPositionLeft
}

// IsRight reports whether the dock position is right.
func (d DockPositionType) IsRight() bool {
	return d == DockPositionRight
}

// IsTop reports whether the dock position is top.
func (d DockPositionType) IsTop() bool {
	return d == DockPositionTop
}

// IsCenter reports whether the dock position is center.
func (d DockPositionType) IsCenter() bool {
	return d == DockPositionCenter
}

// HasDockPosition reports whether the dock position is non-empty and valid.
func (d DockPositionType) HasDockPosition() bool {
	return d.IsValid()
}
