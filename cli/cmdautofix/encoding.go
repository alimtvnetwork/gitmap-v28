package cmdautofix

import (
	"bytes"
	"unicode/utf8"
)

var (
	utf8BOM  = []byte{0xEF, 0xBB, 0xBF}
	utf16LEB = []byte{0xFF, 0xFE}
	utf16BEB = []byte{0xFE, 0xFF}
	crlfByte = []byte{0x0D, 0x0A}
	lfByte   = []byte{0x0A}
)

// encodingCheck flags a BOM prefix or invalid UTF-8 (script 10 port).
func encodingCheck(relPath string, src []byte, opts Options) []Violation {
	violations := []Violation{}
	if hasBOM(src) {
		violations = append(violations, Violation{Path: relPath, Category: "encoding", Detail: "UTF-16/UTF-8 BOM present"})
	}
	if !utf8.Valid(stripBOM(src)) {
		violations = append(violations, Violation{Path: relPath, Category: "encoding", Detail: "bytes are not valid UTF-8"})
	}
	return violations
}

// encodingFix strips the BOM and normalizes CRLF→LF. When the result is still
// not valid UTF-8 it returns the ORIGINAL bytes plus a violation: undecodable
// input is never lossy-written.
func encodingFix(relPath string, src []byte, opts Options) ([]byte, []Violation) {
	fixed := stripBOM(src)
	fixed = bytes.ReplaceAll(fixed, crlfByte, lfByte)
	if bytes.Equal(fixed, src) {
		return src, nil
	}
	if !utf8.Valid(fixed) {
		return src, []Violation{{Path: relPath, Category: "encoding", Detail: "not valid UTF-8 after BOM/CRLF handling; skipped, not lossy-written"}}
	}
	return fixed, nil
}

func hasBOM(src []byte) bool {
	return bytes.HasPrefix(src, utf8BOM) || bytes.HasPrefix(src, utf16LEB) || bytes.HasPrefix(src, utf16BEB)
}

// stripBOM removes a single leading BOM (mirrors utf-8-sig decode semantics).
func stripBOM(src []byte) []byte {
	if bytes.HasPrefix(src, utf8BOM) {
		return src[len(utf8BOM):]
	}
	if bytes.HasPrefix(src, utf16LEB) {
		return src[len(utf16LEB):]
	}
	if bytes.HasPrefix(src, utf16BEB) {
		return src[len(utf16BEB):]
	}
	return src
}
