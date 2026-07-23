// Package glyphs resolves the active glyph-rendering mode (rich vs safe)
// and installs an os.Stdout / os.Stderr byte-stream filter that rewrites
// emoji into ASCII fallbacks when "safe" is active.
//
// Composition: chains AFTER gitmap/theme.Install — both packages wrap
// the standard handles with a pipe; stacking is fine and order-
// independent since each filter only touches its own byte patterns.
package glyphs

import (
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// Mode is the resolved glyph selection.
type Mode int

const (
	// ModeRich passes UTF-8 emoji through unchanged.
	ModeRich Mode = iota
	// ModeSafe rewrites emoji to ASCII fallbacks.
	ModeSafe
)

// Parse maps a user label to a Mode. Unknown values resolve via auto.
func Parse(label string) Mode {
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
func Resolve() Mode {
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
func autoDetect() Mode {
	if isLegacyWindowsHost() {
		return ModeSafe
	}

	return ModeRich
}

// isLegacyWindowsHost is set by the per-OS file (autodetect_windows.go
// vs autodetect_other.go). Stub here keeps the symbol resolvable for
// callers that import glyphs from non-Windows builds.
var isLegacyWindowsHost = func() bool { return false }
