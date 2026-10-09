// Package output owns the CLI's stdout/stderr writer plumbing.
//
// Design (programs 250/265/267): filtering — the theme SGR rewrite,
// then the glyph safe-mode rewrite — happens synchronously inside
// FilterWriter.Write. No pipes, no forwarder goroutines, no Drain.
// cmd.Run builds the writers once after --theme/--glyphs stripping
// (see cli/cmd/outputwriters.go for the single canonical theme-then-
// glyphs composition order) and registers them as the dispatch
// writers (see dispatch.go); UI code takes the explicit writer
// instead of relying on a globally swapped os.Stdout.
//
// Byte-faithful commands (cat/view/type) write file payloads through
// Raw()/RawErr() — the true *os.File handles captured at package
// init — so emoji/ANSI bytes in file content are never rewritten.
//
// os.Stdout and os.Stderr are NEVER reassigned by this package. The
// package is deliberately dependency-free so cli/termout can write
// through UI()/UIErr() without an import cycle.
package output

import (
	"io"
)

// FilterWriter is an io.Writer that applies a chain of byte filters
// synchronously on every Write, in order, then forwards the result
// to dest.
type FilterWriter struct {
	dest    io.Writer
	filters []func([]byte) []byte
}

// NewFilterWriter wraps dest, applying each filter in order on every
// Write. A nil/empty chain is a plain passthrough.
func NewFilterWriter(dest io.Writer, filters ...func([]byte) []byte) *FilterWriter {
	return &FilterWriter{dest: dest, filters: filters}
}

// Write runs p through the whole filter chain, then writes the
// result to dest. It reports len(p) on success so callers observe a
// complete write; filter functions must not retain the returned
// slice beyond the call.
func (w *FilterWriter) Write(p []byte) (int, error) {
	out := p
	for _, f := range w.filters {
		out = f(out)
	}

	if _, err := w.dest.Write(out); err != nil {
		return 0, err
	}

	return len(p), nil
}
