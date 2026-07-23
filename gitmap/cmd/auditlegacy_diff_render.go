// Package cmd: unified-diff rendering for audit-legacy --diffs.
//
// Pure-Go: reads the source file once, rewrites every regex match to
// DefaultAuditLegacyReplace, and emits one single-line hunk per
// changed line in standard `diff --unified=0` format.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// auditDiffLine is one line that changed (1-based line number).
type auditDiffLine struct {
	Line int
	Old  string
	New  string
}

// renderAuditDiff loads file, computes per-line substitutions, and
// returns a unified-diff body. Empty string means "no actual change"
// (e.g. all matches already equal the replacement).
func renderAuditDiff(file string, pats []*regexp.Regexp) (string, error) {
	changes, err := computeAuditDiffLines(file, pats)
	if err != nil || len(changes) == 0 {
		return "", err
	}

	return formatAuditUnifiedDiff(file, changes), nil
}

// computeAuditDiffLines reads file and returns every line whose
// substitution actually differs from the original.
func computeAuditDiffLines(file string, pats []*regexp.Regexp) ([]auditDiffLine, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	return collectAuditDiffLines(scanner, pats), nil
}

// collectAuditDiffLines is the scanner loop split out for the 15-line cap.
func collectAuditDiffLines(scanner *bufio.Scanner, pats []*regexp.Regexp) []auditDiffLine {
	var out []auditDiffLine
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		original := scanner.Text()
		rewritten := applyAuditReplacements(original, pats)
		if rewritten != original {
			out = append(out, auditDiffLine{Line: lineNo, Old: original, New: rewritten})
		}
	}

	return out
}

// applyAuditReplacements substitutes every pattern with the legacy→v8
// replacement string.
func applyAuditReplacements(line string, pats []*regexp.Regexp) string {
	for _, re := range pats {
		line = re.ReplaceAllString(line, constants.DefaultAuditLegacyReplace)
	}

	return line
}

// formatAuditUnifiedDiff emits standard unified-diff text with one
// hunk per changed line (context = 0). Compatible with `patch -p0`.
func formatAuditUnifiedDiff(file string, changes []auditDiffLine) string {
	var b strings.Builder
	fmt.Fprintf(&b, "--- %s\n", file)
	fmt.Fprintf(&b, "+++ %s\n", file)
	for _, c := range changes {
		writeAuditDiffHunk(&b, c)
	}

	return b.String()
}

// writeAuditDiffHunk writes one @@ hunk for a single-line change.
func writeAuditDiffHunk(b *strings.Builder, c auditDiffLine) {
	fmt.Fprintf(b, "@@ -%d,1 +%d,1 @@\n", c.Line, c.Line)
	fmt.Fprintf(b, "-%s\n", c.Old)
	fmt.Fprintf(b, "+%s\n", c.New)
}
