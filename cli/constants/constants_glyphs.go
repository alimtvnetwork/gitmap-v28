// Package constants — glyphs.go: global `--glyphs` switch that controls
// whether emoji-rich (📦 🎉 ✓ →) or ASCII-safe ([pkg] [done] v ->) output
// is rendered. Mirrors the `--theme` pattern so subcommands stay decoupled
// from the rendering decision.
//
// Rationale: legacy Windows PowerShell 5.1 hosts (and any terminal whose
// font lacks the required glyphs) render multi-byte UTF-8 sequences as
// mojibake (Γ£ô, ≡ƒôª). The safe set replaces every non-BMP / problem
// glyph with an ASCII fallback at the byte-stream level, so existing
// Print sites need no changes.
package constants

// Flag + env names.
const (
	FlagGlyphs = "glyphs"
	EnvGlyphs  = "GITMAP_GLYPHS"
)

// Mode labels.
const (
	GlyphsAuto = "auto"
	GlyphsRich = "rich"
	GlyphsSafe = "safe"

	GlyphsDefault = GlyphsAuto
)
