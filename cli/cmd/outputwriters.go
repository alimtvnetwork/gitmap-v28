package cmd

import (
	"io"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/glyphs"
	"github.com/alimtvnetwork/gitmap-v28/cli/output"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// filteredWriters builds the dispatch's synchronous UI writers for
// the resolved theme/glyph modes. This is the SINGLE canonical
// composition order: the theme SGR rewrite runs first, the glyph
// safe-mode rewrite second — exactly the order the old
// termout-then-glyphs pipe stack applied.
//
// Bright theme + rich glyphs is a zero-cost passthrough: both
// filters are identity there, so the raw handles are returned
// unwrapped.
func filteredWriters(themeMode termout.ModeType, glyphMode glyphs.ModeType) (stdout, stderr io.Writer) {
	if themeMode == termout.ModeBright && glyphMode == glyphs.ModeRich {
		return os.Stdout, os.Stderr
	}

	filters := []func([]byte) []byte{
		func(p []byte) []byte { return termout.Filter(p, themeMode) },
		func(p []byte) []byte { return glyphs.Filter(p, glyphMode) },
	}

	return output.NewFilterWriter(os.Stdout, filters...),
		output.NewFilterWriter(os.Stderr, filters...)
}
