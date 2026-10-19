package termgeneric

import (
	"strings"
	"sync"
)

// ANSI escape code constants for terminal text styling.
const (
	CodeReset   = "\033[0m"
	CodeBold    = "\033[1m"
	CodeDim     = "\033[2m"
	CodeRed     = "\033[31m"
	CodeGreen   = "\033[32m"
	CodeYellow  = "\033[33m"
	CodeBlue    = "\033[34m"
	CodeMagenta = "\033[35m"
	CodeCyan    = "\033[36m"
	CodeWhite   = "\033[37m"
)

var (
	colorMutex     sync.RWMutex
	isColorEnabled = true
)

// SetColorEnabled toggles global ANSI color output.
func SetColorEnabled(hasColor bool) {
	colorMutex.Lock()
	defer colorMutex.Unlock()

	isColorEnabled = hasColor
}

// IsColorEnabled checks whether ANSI coloring is currently active.
func IsColorEnabled() bool {
	colorMutex.RLock()
	defer colorMutex.RUnlock()

	return isColorEnabled
}

// Red wraps text in ANSI red escape codes.
func Red(s string) string {
	return wrapColor(s, CodeRed)
}

// Green wraps text in ANSI green escape codes.
func Green(s string) string {
	return wrapColor(s, CodeGreen)
}

// Cyan wraps text in ANSI cyan escape codes.
func Cyan(s string) string {
	return wrapColor(s, CodeCyan)
}

// Yellow wraps text in ANSI yellow escape codes.
func Yellow(s string) string {
	return wrapColor(s, CodeYellow)
}

// Magenta wraps text in ANSI magenta escape codes.
func Magenta(s string) string {
	return wrapColor(s, CodeMagenta)
}

// Blue wraps text in ANSI blue escape codes.
func Blue(s string) string {
	return wrapColor(s, CodeBlue)
}

// Dim wraps text in ANSI dim/faint escape codes.
func Dim(s string) string {
	return wrapColor(s, CodeDim)
}

// Bold wraps text in ANSI bold escape codes.
func Bold(s string) string {
	return wrapColor(s, CodeBold)
}

// Reset returns the ANSI reset escape sequence or wraps text with reset codes.
func Reset(s ...string) string {
	hasColor := IsColorEnabled()
	if !hasColor {
		return strings.Join(s, "")
	}

	if len(s) == 0 {
		return CodeReset
	}

	return CodeReset + strings.Join(s, "") + CodeReset
}

func wrapColor(s string, code string) string {
	if len(s) == 0 {
		return ""
	}

	hasColor := IsColorEnabled()
	if !hasColor {
		return s
	}

	return code + s + CodeReset
}

// StripAnsi removes all ANSI escape sequences (CSI and OSC) from s.
func StripAnsi(s string) string {
	if len(s) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.Grow(len(s))
	idx := 0
	for idx < len(s) {
		if s[idx] == 0x1b {
			idx = skipEscapeSequence(s, idx)
			continue
		}

		sb.WriteByte(s[idx])
		idx++
	}

	return sb.String()
}

func skipEscapeSequence(s string, start int) int {
	if start+1 >= len(s) {
		return start + 1
	}

	nextByte := s[start+1]
	if nextByte == '[' {
		return skipCsiSequence(s, start+2)
	}

	if nextByte == ']' {
		return skipOscSequence(s, start+2)
	}

	return start + 2
}

func skipCsiSequence(s string, start int) int {
	curr := start
	for curr < len(s) {
		b := s[curr]
		if b >= 0x40 && b <= 0x7e {
			return curr + 1
		}
		curr++
	}

	return curr
}

func skipOscSequence(s string, start int) int {
	curr := start
	for curr < len(s) {
		if s[curr] == 0x07 {
			return curr + 1
		}
		if s[curr] == 0x1b && curr+1 < len(s) && s[curr+1] == '\\' {
			return curr + 2
		}
		curr++
	}

	return curr
}
