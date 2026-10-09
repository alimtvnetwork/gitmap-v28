package cmdautofix

import (
	"bytes"
	"path/filepath"
	"regexp"
	"strings"
)

// codeExts mirrors 02-shared-engine.py DEFAULT_CODE_EXTENSIONS.
var codeExts = []string{".ts", ".tsx", ".js", ".jsx", ".go", ".py", ".php", ".cs"}

// isCodeFile reports whether relPath is in the naming category's scope.
func isCodeFile(relPath string) bool {
	return extInList(strings.ToLower(filepath.Ext(relPath)), codeExts)
}

var (
	// Port of 02-shared-engine.py explicit-true patterns.
	explicitDoubleTrueRe = regexp.MustCompile(`(?i)==\s*true\b`)
	explicitTripleTrueRe = regexp.MustCompile(`(?i)===\s*true\b`)
	explicitPythonTrueRe = regexp.MustCompile(`==\s*True\b`)
	// Port of COMMENT_PREFIX: ^\s*(//|#|\*|/\*).
	commentPrefixRe = regexp.MustCompile(`^\s*(?://|#|\*|/\*)`)
)

// namingCheck flags explicit boolean-true comparisons (`== true`,
// `=== true`, `== True`), skipping comment lines via the script's
// line-based heuristic — no parser is built.
// Fix is intentionally nil: this category is report-only by design (D11),
// even under --apply. A style violation is a judgment call for the author,
// not something a machine should rewrite.
func namingCheck(relPath string, src []byte, opts Options) []Violation {
	content := string(bytesNormalizeLF(src))
	// Fast substring pre-filter mirrors the script (5x-10x speedup).
	if !strings.Contains(strings.ToLower(content), "true") {
		return nil
	}
	violations := []Violation{}
	for idx, line := range strings.Split(content, "\n") {
		if commentPrefixRe.MatchString(line) {
			continue
		}
		if !hasExplicitTrue(line) {
			continue
		}
		violations = append(violations, Violation{
			Path:     relPath,
			Category: "naming",
			Line:     idx + 1,
			Detail:   "explicit true comparison: " + strings.TrimSpace(line),
		})
	}
	return violations
}

func hasExplicitTrue(line string) bool {
	if explicitDoubleTrueRe.MatchString(line) {
		return true
	}
	if explicitTripleTrueRe.MatchString(line) {
		return true
	}
	return explicitPythonTrueRe.MatchString(line)
}

// bytesNormalizeLF returns a copy with CRLF folded to LF for line scanning.
func bytesNormalizeLF(src []byte) []byte {
	return bytes.ReplaceAll(src, crlfByte, lfByte)
}
