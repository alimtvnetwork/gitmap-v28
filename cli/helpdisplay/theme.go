// Package helpdisplay renders themed CLI help output (spec 243.3): one DRY,
// object-oriented help-rendering system with theme inheritance and variable
// suggestions. Color codes are borrowed from cli/constants, so the ANSI
// rewrite filter installed by cli/termout.Install adapts them automatically.
package helpdisplay

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// Color wraps a raw ANSI color code ("" = unset).
//
// Codes come from cli/constants — e.g. constants.ColorMagenta — the same
// source every existing help builder uses, so the --theme rewrite filter
// (bright / standard / mono) applies here for free. No new color engine.
type Color struct {
	code string
}

// NewColor builds a Color from a cli/constants ANSI code.
func NewColor(code string) Color {
	return Color{code: code}
}

// IsSet reports whether this color was explicitly configured.
func (c Color) IsSet() bool {
	return len(c.code) > 0
}

// Wrap colors text, returning it unchanged when the color is unset so
// unthemed output never leaks raw escape sequences.
func (c Color) Wrap(text string) string {
	if !c.IsSet() {
		return text
	}
	return c.code + text + constants.ColorReset
}

// Theme carries the color roles a HelpDisplay renders with.
//
// Inheritance: an explicitly set field wins; otherwise the parent theme's
// resolved value applies; a nil parent chain ends at the zero Color (unset,
// i.e. plain text). This lets per-help tweaks layer on a shared termout.
type Theme struct {
	headerColor      Color
	commandColor     Color
	descriptionColor Color
	hintColor        Color
	parent           *Theme
}

// NewTheme creates a theme inheriting from parent (nil = global root).
func NewTheme(parent *Theme) *Theme {
	return &Theme{parent: parent}
}

// WithHeaderColor overrides the header role; returns t for chaining.
func (t *Theme) WithHeaderColor(c Color) *Theme {
	t.headerColor = c
	return t
}

// WithCommandColor overrides the command (left side) role.
func (t *Theme) WithCommandColor(c Color) *Theme {
	t.commandColor = c
	return t
}

// WithDescriptionColor overrides the description (right side) role.
func (t *Theme) WithDescriptionColor(c Color) *Theme {
	t.descriptionColor = c
	return t
}

// WithHintColor overrides the hint role.
func (t *Theme) WithHintColor(c Color) *Theme {
	t.hintColor = c
	return t
}

// HeaderColor resolves the header role down the inheritance chain.
func (t *Theme) HeaderColor() Color {
	return resolveThemeColor(t, fieldHeader)
}

// CommandColor resolves the command (left side) role down the chain.
func (t *Theme) CommandColor() Color {
	return resolveThemeColor(t, fieldCommand)
}

// DescriptionColor resolves the description (right side) role down the chain.
func (t *Theme) DescriptionColor() Color {
	return resolveThemeColor(t, fieldDescription)
}

// HintColor resolves the hint role down the chain.
func (t *Theme) HintColor() Color {
	return resolveThemeColor(t, fieldHint)
}

// fieldHeader extracts the headerColor role from a termout.
func fieldHeader(t *Theme) Color {
	return t.headerColor
}

// fieldCommand extracts the commandColor role from a termout.
func fieldCommand(t *Theme) Color {
	return t.commandColor
}

// fieldDescription extracts the descriptionColor role from a termout.
func fieldDescription(t *Theme) Color {
	return t.descriptionColor
}

// fieldHint extracts the hintColor role from a termout.
func fieldHint(t *Theme) Color {
	return t.hintColor
}

// resolveThemeColor walks from the theme up through its parents returning the
// first explicitly set color for the role; a nil chain yields the zero Color.
func resolveThemeColor(t *Theme, field func(*Theme) Color) Color {
	if t == nil {
		return Color{}
	}
	if field(t).IsSet() {
		return field(t)
	}
	return resolveThemeColor(t.parent, field)
}
