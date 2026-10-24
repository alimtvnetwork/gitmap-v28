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
	"github.com/alimtvnetwork/gitmap-v28/cli/strutil"
)
import "github.com/alimtvnetwork/gitmap-v28/cli/glyphs/enums"

// ModeType is the resolved glyph selection.
// ModeType aliases the strongly-typed enum from the enums package.
type ModeType = enums.ModeType

const (
	ModeRich = enums.ModeRich
	ModeSafe = enums.ModeSafe
)

var glyphModeByLabel = map[string]ModeType{
	constants.GlyphsRich: ModeRich,
	constants.GlyphsSafe: ModeSafe,
}

// Parse maps a user label to a ModeType. Unknown values resolve via auto.
func Parse(label string) ModeType {
	if mode, ok := glyphModeByLabel[strings.ToLower(strings.TrimSpace(label))]; ok {
		return mode
	}

	return autoDetect()
}

// Resolve picks the mode from the GITMAP_GLYPHS env var (populated by
// the flag stripper or the user's shell).
func Resolve() ModeType {
	return Parse(os.Getenv(constants.EnvGlyphs))
}

// IsValidLabel reports whether label is a recognized glyph mode.
func IsValidLabel(label string) bool {
	return strutil.EqualFoldAnyTrim(label, constants.GlyphsAuto, constants.GlyphsRich, constants.GlyphsSafe)
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
