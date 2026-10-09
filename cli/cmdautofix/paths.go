package cmdautofix

import (
	"bytes"
	"regexp"
	"strings"
)

// defaultRepoURI is the --uri-pattern default (D10).
const defaultRepoURI = "coding-guidelines"

// devicePathPrefix mirrors 02-shared-engine.py DEVICE_PATH_PREFIX ("\\?\").
const devicePathPrefix = `\\?\`

var (
	// Port of FILE_URI_WIN: file:///[A-Za-z]:/[^\s\)\]"'>]+
	fileURIWinRe = regexp.MustCompile(`file:///[A-Za-z]:/[^\s\)\]"'>]+`)
	// Port of DRIVE_ABS_WIN without the lookbehind (RE2 has none):
	// (?<![A-Za-z0-9_])[A-Za-z]:\\[A-Za-z0-9_\\.-]+
	// The lookbehind is enforced manually in isDriveAbsMatch.
	driveAbsWinRe = regexp.MustCompile(`[A-Za-z]:\\[A-Za-z0-9_\\.-]+`)
)

// repoFileURIRe builds the REPO_FILE_URI pattern for a repo-URI pattern P:
// file:///[A-Za-z]:/[^/]+/P/([^\s\)\]"'>]+) — group 1 is the relative target.
func repoFileURIRe(pattern string) *regexp.Regexp {
	return regexp.MustCompile(`file:///[A-Za-z]:/[^/]+/` + regexp.QuoteMeta(pattern) + `/([^\s\)\]"'>]+)`)
}

// effectiveURIPattern returns Options.URIPattern or the default (D10).
func effectiveURIPattern(opts Options) string {
	if opts.URIPattern != "" {
		return opts.URIPattern
	}
	return defaultRepoURI
}

// isDriveAbsMatch applies the Python lookbehind (?<![A-Za-z0-9_]).
func isDriveAbsMatch(content string, loc []int) bool {
	if loc[0] == 0 {
		return true
	}
	prev := content[loc[0]-1]
	isWordChar := prev == '_' ||
		(prev >= 'A' && prev <= 'Z') ||
		(prev >= 'a' && prev <= 'z') ||
		(prev >= '0' && prev <= '9')
	return !isWordChar
}

// lineAt returns the 1-based line number of a byte offset.
func lineAt(content string, offset int) int {
	line := 1
	for i := 0; i < offset && i < len(content); i++ {
		if content[i] == '\n' {
			line++
		}
	}
	return line
}

// pathsCheck audits for forbidden absolute paths (script 07 port): file:///
// URIs and Windows drive-absolute paths, ignoring device paths and trailing
// dots — mirroring the script's violation loop exactly.
func pathsCheck(relPath string, src []byte, opts Options) []Violation {
	content := string(bytesNormalizeLF(src))
	// Fast substring pre-filter mirrors the script.
	if !strings.Contains(content, "file:") && !strings.Contains(content, `:\`) {
		return nil
	}
	violations := []Violation{}
	for _, loc := range fileURIWinRe.FindAllStringIndex(content, -1) {
		val := content[loc[0]:loc[1]]
		if isIgnoredPathVal(val) {
			continue
		}
		violations = append(violations, Violation{Path: relPath, Category: "paths", Line: lineAt(content, loc[0]), Detail: "absolute path found: " + val})
	}
	for _, loc := range driveAbsWinRe.FindAllStringIndex(content, -1) {
		if !isDriveAbsMatch(content, loc) {
			continue
		}
		val := content[loc[0]:loc[1]]
		if isIgnoredPathVal(val) {
			continue
		}
		violations = append(violations, Violation{Path: relPath, Category: "paths", Line: lineAt(content, loc[0]), Detail: "absolute path found: " + val})
	}
	return violations
}

func isIgnoredPathVal(val string) bool {
	if strings.Contains(val, devicePathPrefix) {
		return true
	}
	return strings.HasSuffix(val, ".")
}

// pathsFix mirrors script 07's sanitize_content_paths exactly: every
// file:///<drive>:/<anything>/<pattern>/<rel> match is replaced by the
// relative target (group 1). Byte-identical behavior to the script by
// default; --uri-pattern swaps the pattern.
func pathsFix(relPath string, src []byte, opts Options) ([]byte, []Violation) {
	re := repoFileURIRe(effectiveURIPattern(opts))
	modified := src
	count := 0
	for _, m := range re.FindAllSubmatch(src, -1) {
		rel := m[1]
		modified = bytes.ReplaceAll(modified, m[0], rel)
		count++
	}
	if count == 0 || bytes.Equal(modified, src) {
		return src, nil
	}
	// No violations reported: the scan-time Check already recorded every
	// finding; Fix never encounters a new unfixable condition.
	return modified, nil
}
