// Package glyphs — filter.go: byte-level emoji → ASCII rewriter.
//
// Implemented as a precomputed bytes.Replacer-equivalent: a flat list
// of (needle, replacement) pairs scanned in a single left-to-right
// pass. Keeps allocations to one buffer per Write call.
package glyphs

import "bytes"

// pair is one emoji → fallback substitution.
type pair struct {
	from []byte
	to   []byte
}

// table is the ordered substitution list. Longer needles first when
// they overlap shorter ones (e.g. ✅ before bare ✓).
var table = buildTable()

// buildTable returns the canonical glyph → ASCII mapping. The fallback
// strings stay short so terminal layouts stay aligned.
func buildTable() []pair {
	raw := [...][2]string{
		// Variation selector — strip silently so "⚠️" → "!" not "!?"
		{"\uFE0F", ""},
		// Status / arrows.
		{"✅", "[OK]"}, {"❌", "[X]"}, {"✔", "v"},
		{"✓", "v"}, {"✗", "x"}, {"→", "->"},
		{"⚠", "!"}, {"ℹ", "i"}, {"⚡", "[!]"},
		// Object emoji.
		{"📦", "[pkg]"}, {"📄", "[doc]"}, {"📁", "[dir]"},
		{"📂", "[dir]"}, {"📝", "[edit]"}, {"🔍", "[find]"},
		{"🎉", "[done]"}, {"📊", "[chart]"}, {"💾", "[save]"},
		{"🏷", "[tag]"}, {"🚀", "[go]"}, {"🔑", "[key]"},
		{"🗺", "[map]"}, {"🧬", "[dna]"}, {"🌳", "[tree]"},
		{"🪄", "[wand]"}, {"🔐", "[lock]"}, {"🖥", "[host]"},
		{"🗄", "[db]"}, {"🧭", "[nav]"}, {"📤", "[out]"},
		{"📰", "[news]"}, {"📖", "[book]"}, {"🪟", "[win]"},
		{"🐧", "[nix]"}, {"📋", "[copy]"}, {"📡", "[net]"},
		{"🔁", "[loop]"}, {"💡", "[tip]"}, {"📎", "[ref]"},
	}

	out := make([]pair, 0, len(raw))
	for _, r := range raw {
		out = append(out, pair{from: []byte(r[0]), to: []byte(r[1])})
	}

	return out
}

// Filter returns p with every glyph rewritten per table when mode is
// ModeSafe. ModeRich returns p unchanged (zero-cost passthrough).
func Filter(p []byte, mode Mode) []byte {
	if mode == ModeRich {
		return p
	}

	out := p
	for _, pr := range table {
		if bytes.Contains(out, pr.from) {
			out = bytes.ReplaceAll(out, pr.from, pr.to)
		}
	}

	return out
}
