package cmdautofix

import (
	"bytes"
	"regexp"
)

// gapRunRe mirrors 31-md-gap-fixer.py fix_gaps: 3+ newlines → exactly 2.
var gapRunRe = regexp.MustCompile(`\n{3,}`)

// collapseGaps mirrors fix_gaps: normalize CRLF→LF, then collapse runs of
// 3+ newlines into exactly 2.
func collapseGaps(src []byte) []byte {
	normalized := bytes.ReplaceAll(src, crlfByte, lfByte)
	return gapRunRe.ReplaceAll(normalized, lfByte[:2])
}

// markdownCheck flags any run of 3+ consecutive newlines in .md files.
func markdownCheck(relPath string, src []byte, opts Options) []Violation {
	normalized := bytes.ReplaceAll(src, crlfByte, lfByte)
	violations := []Violation{}
	for _, loc := range gapRunRe.FindAllIndex(normalized, -1) {
		violations = append(violations, Violation{
			Path:     relPath,
			Category: "markdown",
			Line:     lineAt(string(normalized), loc[0]),
			Detail:   "3+ consecutive blank lines (collapse to one)",
		})
	}
	return violations
}

// markdownFix collapses 3+ newlines → 2; no-op rewrites are never written.
func markdownFix(relPath string, src []byte, opts Options) ([]byte, []Violation) {
	fixed := collapseGaps(src)
	if bytes.Equal(fixed, src) {
		return src, nil
	}
	// Check already recorded the gap hits; Fix reports only new findings.
	return fixed, nil
}
