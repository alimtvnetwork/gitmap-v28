// Package glyphs resolves the active glyph-rendering mode (rich vs safe)
// and provides the byte-level emoji → ASCII filter applied when "safe"
// is active.
//
// Filtering runs synchronously inside the cli/output FilterWriter that
// cmd.Run builds (theme first, glyphs second) — os.Stdout/os.Stderr are
// never reassigned. Each filter only touches its own byte patterns, so
// the composition order is interchangeable in effect.
package glyphs

import (
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ModeType is the resolved glyph selection.
type ModeType int

const (
	// ModeRich passes UTF-8 emoji through unchanged.
	ModeRich ModeType = iota
	// ModeSafe rewrites emoji to ASCII fallbacks.
	ModeSafe
)

// Parse maps a user label to a ModeType. Unknown values resolve via auto.
func Parse(label string) ModeType {
	switch strings.ToLower(strings.TrimSpace(label)) {
	case constants.GlyphsRich:
		return ModeRich
	case constants.GlyphsSafe:
		return ModeSafe
	default:
		return autoDetect()
	}
}

// Resolve picks the mode from the GITMAP_GLYPHS env var (populated by
// the flag stripper or the user's shell).
func Resolve() ModeType {
	return Parse(os.Getenv(constants.EnvGlyphs))
}

// IsValidLabel reports whether label is a recognized glyph mode.
func IsValidLabel(label string) bool {
	switch strings.ToLower(strings.TrimSpace(label)) {
	case constants.GlyphsAuto, constants.GlyphsRich, constants.GlyphsSafe:
		return true
	default:
		return false
	}
}

// autoDetect returns ModeRich on terminals known to render emoji well
// (Windows Terminal, VS Code, ConEmu, iTerm2, every *nix TTY) and
// ModeSafe on legacy Windows ConsoleHost (powershell.exe 5.1, cmd.exe)
// where the host font typically lacks the required glyphs.
func autoDetect() ModeType {
	if isLegacyWindowsHost() || isLegacyUnixHost() {
		return ModeSafe
	}

	return ModeRich
}

// isLegacyWindowsHost is set by autodetect_windows.go on Windows.
var isLegacyWindowsHost = func() bool { return false }

// isLegacyUnixHost is set by autodetect_other.go on non-Windows platforms.
var isLegacyUnixHost = func() bool { return false }
