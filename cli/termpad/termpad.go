package termpad

import (
	"fmt"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	padMu          sync.Mutex
	paddingApplied bool
	lastWasRule    bool
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

// ResetPaddingState resets the global padding and boundary rule tracking.
func ResetPaddingState() {
	padMu.Lock()
	defer padMu.Unlock()
	paddingApplied = false
	lastWasRule = false
}

// EnsureBottomPadding prints a closing newline conditionally if not already padded.
func EnsureBottomPadding(lastOutput string) {
	padMu.Lock()
	defer padMu.Unlock()
	if paddingApplied {
		return
	}
	paddingApplied = true

	if strings.HasSuffix(lastOutput, "\n\n") {
		return
	}
	if strings.HasSuffix(lastOutput, "\n") {
		fmt.Println()
		return
	}
	fmt.Print("\n\n")
}

// PrintBlueRule prints a blue separator rule unless one was just printed.
func PrintBlueRule(width int) {
	padMu.Lock()
	defer padMu.Unlock()
	if lastWasRule {
		return
	}
	if width <= 0 {
		width = 80
	}
	fmt.Printf("%s%s%s\n", constants.ColorBlue, strings.Repeat("─", width), constants.ColorReset)
	lastWasRule = true
}

// PrintSeparator prints a separator line unless a boundary rule was already printed.
func PrintSeparator(line string) {
	padMu.Lock()
	defer padMu.Unlock()
	if lastWasRule && IsBoundaryRule(line) {
		return
	}
	fmt.Println(line)
	lastWasRule = IsBoundaryRule(line)
}

// FormatPadded formats raw text with a 2-space left indentation for all lines,
// avoiding duplicate boundary rules and honoring existing left margins.
func FormatPadded(text string) string {
	lines := strings.Split(text, "\n")
	var sb strings.Builder
	var prevWasRule bool
	for i, line := range lines {
		appendPaddedLine(&sb, line, i, len(lines), &prevWasRule)
	}

	return sb.String()
}

func appendPaddedLine(sb *strings.Builder, line string, idx, total int, prevWasRule *bool) {
	isRule := IsBoundaryRule(line)
	if isRule && *prevWasRule {
		return
	}
	if !hasExistingMargin(line) {
		sb.WriteString("  ")
	}
	sb.WriteString(line)
	if idx < total-1 {
		sb.WriteString("\n")
	}
	*prevWasRule = isRule
}

// PrintPadded outputs left-padded text directly to standard output.
func PrintPadded(text string) {
	fmt.Print(FormatPadded(text))
}

// IsBoundaryRule reports whether line consists entirely of repeating separator characters.
func IsBoundaryRule(line string) bool {
	plain := strings.TrimSpace(StripAnsi(line))
	if len(plain) < 3 {
		return false
	}

	return isAllRuleRunes(plain)
}

func isAllRuleRunes(plain string) bool {
	for _, r := range plain {
		if !isRuleRune(r) {
			return false
		}
	}

	return true
}

func isRuleRune(r rune) bool {
	return r == '─' || r == '-' || r == '=' || r == '_' || r == '~' || r == '━' || r == '═' || r == '—'
}

// HasExistingMargin reports whether line has leading whitespace or indentation.
func HasExistingMargin(line string) bool {
	return hasExistingMargin(line)
}

func hasExistingMargin(line string) bool {
	if len(line) == 0 {
		return true
	}
	plain := StripAnsi(line)
	if len(plain) == 0 {
		return true
	}

	return strings.HasPrefix(plain, "  ") || strings.HasPrefix(plain, "\t")
}

// StripAnsi removes ANSI escape sequences from s.
func StripAnsi(s string) string {
	var sb strings.Builder
	sb.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if isAnsiEscapeStart(s, i) {
			i = skipAnsiSequence(s, i)
			continue
		}
		sb.WriteByte(s[i])
	}

	return sb.String()
}

func isAnsiEscapeStart(s string, i int) bool {
	return s[i] == 0x1b && i+1 < len(s) && s[i+1] == '['
}

func skipAnsiSequence(s string, start int) int {
	j := start + 2
	for j < len(s) && s[j] != 'm' {
		j++
	}

	return j
}
