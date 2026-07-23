// Package theme — filter.go: byte-level ANSI SGR rewriter.
//
// Kept dependency-free (no regexp) so the per-write overhead stays
// tiny — the filter sits on the hot path of every Printf in the
// codebase.
package theme

import "bytes"

const (
	esc    = 0x1b
	sgrEnd = 'm'
)

// Filter applies the active mode's transformation to p. Bright mode is
// a passthrough; other modes scan for SGR escape sequences (ESC `[`
// ... `m`) and rewrite or drop them.
func Filter(p []byte, mode Mode) []byte {
	if mode == ModeBright {
		return p
	}

	var out bytes.Buffer
	out.Grow(len(p))

	i := 0
	for i < len(p) {
		start, end, ok := nextSGR(p, i)
		if !ok {
			out.Write(p[i:])
			break
		}
		out.Write(p[i:start])
		if mode == ModeStandard {
			out.WriteString(downgrade(string(p[start : end+1])))
		}
		i = end + 1
	}

	return out.Bytes()
}

// nextSGR locates the next SGR escape sequence (ESC `[` ... `m`)
// starting at or after off. Returns its [start, end] inclusive byte
// offsets, or ok=false when none is found.
func nextSGR(p []byte, off int) (start, end int, ok bool) {
	for i := off; i < len(p)-1; i++ {
		if p[i] != esc || p[i+1] != '[' {
			continue
		}
		for j := i + 2; j < len(p); j++ {
			if p[j] == sgrEnd {
				return i, j, true
			}
			if !isSGRParam(p[j]) {
				break
			}
		}
	}

	return 0, 0, false
}

// isSGRParam reports whether b is a legal SGR parameter byte (digits
// and `;`). Any other byte terminates the scan — we deliberately do
// not try to rewrite non-SGR CSI sequences (cursor moves, clears,
// etc.), they pass through untouched.
func isSGRParam(b byte) bool {
	return (b >= '0' && b <= '9') || b == ';'
}

// downgrade maps a bright-palette SGR escape to its plain-3X
// equivalent. Unknown escapes pass through unchanged so plain SGRs
// emitted by other surfaces (e.g. the changelog pretty-renderer's
// `\033[36m` / `\033[39m` / `\033[1m`) keep working in standard
// mode — only the explicitly-bright variants are rewritten.
func downgrade(seq string) string {
	if alt, ok := standardDowngrades[seq]; ok {
		return alt
	}

	return seq
}

// standardDowngrades enumerates the bright-palette escapes emitted by
// constants.Color* so the rewriter stays deterministic. Anything not
// listed here passes through unchanged (see downgrade). Adding a new
// bright accent in constants_terminal.go SHOULD add a row here so the
// "standard" theme actually tones it down.
var standardDowngrades = map[string]string{
	"\033[1;92m":     "\033[32m",
	"\033[1;91m":     "\033[31m",
	"\033[1;93m":     "\033[33m",
	"\033[1;96m":     "\033[36m",
	"\033[1;97m":     "\033[97m",
	"\033[2;37m":     "\033[90m",
	"\033[1;95m":     "\033[35m",
	"\033[1;94m":     "\033[34m",
	"\033[38;5;208m": "\033[33m",
}

