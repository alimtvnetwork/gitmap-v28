package startup

// Tiny shim so remove.go can call parseDesktopFields without
// importing bufio directly — keeps that file's import list short
// and lets desktop.go own the scanner construction details.

import (
	"bufio"
	"io"
)

// newScanner wraps an io.Reader in a bufio.Scanner with default
// (line-based) tokenization. .desktop files are tiny (kilobytes) so
// the default 64KiB buffer is always sufficient.
func newScanner(r io.Reader) *bufio.Scanner {
	return bufio.NewScanner(r)
}
