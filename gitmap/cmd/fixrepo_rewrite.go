package cmd

// Token-rewrite engine for `gitmap fix-repo`. Mirrors
// scripts/fix-repo/Rewrite-Engine.ps1: replace literal `{base}-v{N}`
// with `{base}-v{current}` for every N in targets, guarded by a
// negative-lookahead so `-v1` does not match inside `-v10`.
//
// Go's RE2 has no native `(?!...)` so the guard is implemented as a
// hand-rolled scan that walks each occurrence of the literal token
// and inspects the next byte before substituting.

import (
	"os"
	"strconv"
	"strings"
)

// rewriteFixRepoFile reads fullPath, applies every target rewrite,
// and (unless dryRun) writes the result back. Returns the total
// replacement count across all targets, or an error on read/write
// failure.
func rewriteFixRepoFile(fullPath, base string, current int, targets []int, dryRun bool) (int, error) {
	return rewriteFixRepoFileR(fullPath, base, current, targets, dryRun, false)
}

// rewriteFixRepoFileR is the restrict-aware variant (v5.39.0+).
func rewriteFixRepoFileR(fullPath, base string, current int, targets []int, dryRun, restrictNoVersion bool) (int, error) {
	original, err := os.ReadFile(fullPath)
	if err != nil {
		return 0, err
	}
	updated, count := applyAllTargetsR(string(original), base, current, targets, restrictNoVersion)
	if count == 0 {
		return 0, nil
	}
	if dryRun {
		return count, nil
	}
	if err := os.WriteFile(fullPath, []byte(updated), 0o644); err != nil {
		return 0, err
	}

	return count, nil
}

// applyAllTargets is the unrestricted entry preserved for tests and
// existing callers. Delegates to applyAllTargetsR with restrict=false.
func applyAllTargets(text, base string, current int, targets []int) (string, int) {
	return applyAllTargetsR(text, base, current, targets, false)
}

// applyAllTargetsR folds every target rewrite over text. When
// restrictNoVersion is true, the v1→v2 bare-base sweep is skipped so
// ONLY `{base}-vN` tokens are rewritten — see spec
// 27-fix-repo-command.md §"Restrict modes".
func applyAllTargetsR(text, base string, current int, targets []int, restrictNoVersion bool) (string, int) {
	total := 0
	for _, n := range targets {
		updated, added := applyOneTarget(text, base, n, current)
		text = updated
		total += added
		// Bare-base sweep is gated on the v1→v2 transition (see the
		// v5.38.0 scope rule). v5.39.0+ adds a second gate:
		// `--restrict no-version` suppresses the sweep entirely so the
		// caller can guarantee bare `{base}` tokens are never touched
		// even during a v1→v2 bump.
		if n == 1 && current == 2 && !restrictNoVersion {
			updated, added = applyBareBase(text, base, current)
			text = updated
			total += added
		}
	}

	return text, total
}

// applyBareBase rewrites every bare `{base}` occurrence to
// `{base}-v{current}` when it appears between non-word boundaries.
// "Word char" includes ASCII alnum, `_`, `-`, `.` so versioned forms
// like `{base}-v2` (next byte `-`) and dotted forms like `{base}.js`
// are left alone — the bare-base rewrite targets ONLY the standalone
// repo-name token left over from the pre-versioned (v1 == bare) era.
func applyBareBase(text, base string, current int) (string, int) {
	replacement := base + "-v" + strconv.Itoa(current)
	if base == "" || !strings.Contains(text, base) {
		return text, 0
	}

	return scanBareBase(text, base, replacement)
}

// scanBareBase is the inner loop — split out so applyBareBase stays
// under the 15-line cap.
func scanBareBase(text, base, replacement string) (string, int) {
	var b strings.Builder
	count := 0
	tlen := len(base)
	pos := 0
	for pos <= len(text) {
		rel := strings.Index(text[pos:], base)
		if rel < 0 {
			b.WriteString(text[pos:])
			break
		}
		idx := pos + rel
		b.WriteString(text[pos:idx])
		end := idx + tlen
		count += writeBareBaseHit(&b, text, idx, end, base, replacement)
		pos = end
	}

	return b.String(), count
}

// writeBareBaseHit emits either replacement or the literal base.
func writeBareBaseHit(b *strings.Builder, text string, start, end int,
	base, replacement string,
) int {
	if isBareBaseBoundary(text, start, end) {
		b.WriteString(replacement)

		return 1
	}
	b.WriteString(base)

	return 0
}

// isBareBaseBoundary returns true iff the bytes immediately before
// `start` and at `end` are not word chars (alnum / `_` / `-` / `.`).
// Start/end of string count as non-word boundaries.
func isBareBaseBoundary(text string, start, end int) bool {
	if start > 0 && isBareBaseWordByte(text[start-1]) {
		return false
	}
	if end < len(text) && isBareBaseWordByte(text[end]) {
		return false
	}

	return true
}

// isBareBaseWordByte reports whether c continues an identifier-like
// token (so the bare-base scan must NOT treat the adjacent position
// as a boundary).
func isBareBaseWordByte(c byte) bool {
	switch {
	case c >= 'a' && c <= 'z':
		return true
	case c >= 'A' && c <= 'Z':
		return true
	case c >= '0' && c <= '9':
		return true
	case c == '_' || c == '-' || c == '.':
		return true
	}

	return false
}

// applyOneTarget walks every literal `{base}-vN` occurrence and
// substitutes it with `{base}-v{current}` when the next byte is not
// an ASCII digit (so `-v1` does not match inside `-v10`).
func applyOneTarget(text, base string, n, current int) (string, int) {
	token := base + "-v" + strconv.Itoa(n)
	replacement := base + "-v" + strconv.Itoa(current)
	if !strings.Contains(text, token) {
		return text, 0
	}

	return rewriteToken(text, token, replacement)
}

// rewriteToken is the inner scan loop. Extracted so applyOneTarget
// stays well under the 15-line cap and the loop can be unit-tested
// directly without going through the file-IO layer.
func rewriteToken(text, token, replacement string) (string, int) {
	var b strings.Builder
	count := 0
	tlen := len(token)
	for {
		idx := strings.Index(text, token)
		if idx < 0 {
			b.WriteString(text)
			break
		}
		b.WriteString(text[:idx])
		count += writeOneTokenHit(&b, text, idx, tlen, token, replacement)
		text = text[idx+tlen:]
	}

	return b.String(), count
}

// writeOneTokenHit emits either the replacement (when the byte after
// the token is not an ASCII digit) or the literal token unchanged
// (when it IS a digit, i.e. we matched a prefix of -v10/-v123/etc).
// Returns 1 on substitution, 0 on guarded skip.
func writeOneTokenHit(b *strings.Builder, text string, idx, tlen int,
	token, replacement string,
) int {
	nextOff := idx + tlen
	if nextOff < len(text) && isASCIIDigit(text[nextOff]) {
		b.WriteString(token)

		return 0
	}
	b.WriteString(replacement)

	return 1
}

// isASCIIDigit returns true when c is in '0'..'9'. Inlined helper
// keeps the hot-path readable without pulling in unicode tables.
func isASCIIDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
