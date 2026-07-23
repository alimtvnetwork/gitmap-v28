package cmd

// Path-scope helpers for the regoldens diff summary. Extracted so
// regoldens_diff.go stays under the 200-line file cap.

import "strings"

// goldenDiffPathFragment scopes git output to fixture files. Anything
// outside `testdata/` is filtered because regenerate passes should
// only touch those paths.
const goldenDiffPathFragment = "testdata/"

// goldenDiffBasenameFragment further narrows the report to files
// whose basename contains "golden" (case-insensitive). Excludes
// unrelated `testdata/` fixtures (corpora, sample inputs, schemas)
// that may legitimately be touched by side-effects but are not the
// regeneration target.
const goldenDiffBasenameFragment = "golden"

// isGoldenFixturePath reports whether p is a regeneration-relevant
// golden file: it must live under a `testdata/` directory AND have
// "golden" in its basename (case-insensitive). Used to gate every
// entry shown by the diff summary.
func isGoldenFixturePath(p string) bool {
	if !strings.Contains(p, goldenDiffPathFragment) {
		return false
	}
	base := p
	if i := strings.LastIndex(p, "/"); i >= 0 {
		base = p[i+1:]
	}
	return strings.Contains(strings.ToLower(base), goldenDiffBasenameFragment)
}
