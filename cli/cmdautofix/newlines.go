package cmdautofix

import (
	"bytes"
	"strings"
	"unicode"
)

// fixNewlines mirrors 04-newline-fixer.py clean_file_content: CRLF→LF, trim
// trailing whitespace per line, exactly one final newline. An empty file
// stays empty — no newline is invented.
func fixNewlines(src []byte) []byte {
	if len(src) == 0 {
		return src
	}
	normalized := bytes.ReplaceAll(src, crlfByte, lfByte)
	parts := bytes.Split(normalized, lfByte)
	lines := make([][]byte, 0, len(parts))
	for _, line := range parts {
		lines = append(lines, trimTrailingSpace(line))
	}
	for len(lines) > 0 && len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	out := bytes.Join(lines, lfByte)
	return append(out, '\n')
}

// trimTrailingSpace strips trailing Unicode whitespace (Python rstrip() semantics).
func trimTrailingSpace(line []byte) []byte {
	s := string(line)
	trimmed := strings.TrimRightFunc(s, unicode.IsSpace)
	if len(trimmed) == len(s) {
		return line
	}
	return []byte(trimmed)
}

// newlinesCheck flags CRLF, trailing-whitespace lines, and files not ending
// in exactly one newline (script 04 port).
func newlinesCheck(relPath string, src []byte, opts Options) []Violation {
	violations := []Violation{}
	if bytes.Contains(src, crlfByte) {
		violations = append(violations, Violation{Path: relPath, Category: "newlines", Detail: "CRLF line endings present"})
	}
	for idx, line := range bytes.Split(bytes.ReplaceAll(src, crlfByte, lfByte), lfByte) {
		if !bytes.Equal(line, trimTrailingSpace(line)) {
			violations = append(violations, Violation{Path: relPath, Category: "newlines", Line: idx + 1, Detail: "trailing whitespace"})
		}
	}
	if len(src) > 0 && !endsWithSingleNewline(src) {
		violations = append(violations, Violation{Path: relPath, Category: "newlines", Detail: "file does not end with exactly one newline"})
	}
	return violations
}

func endsWithSingleNewline(src []byte) bool {
	if len(src) == 0 || src[len(src)-1] != '\n' {
		return false
	}
	return len(src) < 2 || src[len(src)-2] != '\n'
}

// newlinesFix rewrites only when the bytes actually change. It reports no
// violations: the scan-time Check already recorded every finding, and Fix
// never encounters a new unfixable condition.
func newlinesFix(relPath string, src []byte, opts Options) ([]byte, []Violation) {
	fixed := fixNewlines(src)
	if bytes.Equal(fixed, src) {
		return src, nil
	}
	return fixed, nil
}
