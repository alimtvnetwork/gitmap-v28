package termpad

import (
	"fmt"
	"strings"
	"sync"
)

var (
	padMu          sync.Mutex
	paddingApplied bool
)

// SetPaddingApplied updates the global padding state.
func SetPaddingApplied(applied bool) {
	padMu.Lock()
	defer padMu.Unlock()
	paddingApplied = applied
}

// IsPaddingApplied reports whether terminal padding is active.
func IsPaddingApplied() bool {
	padMu.Lock()
	defer padMu.Unlock()

	return paddingApplied
}

// EnsureBottomPadding prints a closing newline conditionally if not already padded.
func EnsureBottomPadding(lastOutput string) {
	if IsPaddingApplied() {
		return
	}
	if strings.HasSuffix(lastOutput, "\n\n") {
		return
	}
	if strings.HasSuffix(lastOutput, "\n") {
		fmt.Println()
		return
	}
	fmt.Print("\n\n")
}

// FormatPadded formats raw text with a 2-space left indentation for all lines.
func FormatPadded(text string) string {
	lines := strings.Split(text, "\n")
	var sb strings.Builder
	for i, line := range lines {
		appendPaddedLine(&sb, line, i, len(lines))
	}

	return sb.String()
}

func appendPaddedLine(sb *strings.Builder, line string, idx, total int) {
	if len(line) > 0 && !strings.HasPrefix(line, "  ") {
		sb.WriteString("  ")
	}
	sb.WriteString(line)
	if idx < total-1 {
		sb.WriteString("\n")
	}
}

// PrintPadded outputs left-padded text directly to standard output.
func PrintPadded(text string) {
	fmt.Print(FormatPadded(text))
}
