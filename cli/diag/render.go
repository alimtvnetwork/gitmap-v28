package diag

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strings"
)

const (
	boxWidth          = 72
	catppuccinBorder  = "\033[38;2;137;220;235m" // Pastel Cyan
	catppuccinHeader  = "\033[38;2;203;166;247m" // Pastel Magenta
	catppuccinCommand = "\033[38;2;166;227;161m" // Pastel Green
	catppuccinDim     = "\033[2;37m"             // Muted Dim
	catppuccinPrompt  = "\033[38;2;249;226;175m" // Pastel Yellow
	catppuccinReset   = "\033[0m"
)

// Render dispatches rendering according to requested RenderFormat.
func Render(w io.Writer, group SuggestionGroup, format RenderFormat) error {
	switch format {
	case RenderFormatCompact:
		return RenderCompact(w, group)
	case RenderFormatJSON:
		return RenderJSON(w, group)
	default:
		return RenderBox(w, group)
	}
}

// RenderBox renders an elegant Catppuccin terminal card for suggestions.
func RenderBox(w io.Writer, group SuggestionGroup) error {
	if !group.HasSuggestions() {
		return nil
	}
	title := resolveBoxTitle(group.Title)
	writeBoxTop(w, title, boxWidth)
	writeReasonIfPresent(w, group.Reason, boxWidth)
	writeSuggestionItems(w, group.Suggestions, boxWidth)
	writeActionHintIfPresent(w, group.Suggestions, boxWidth)
	writeBoxBottom(w, boxWidth)
	return nil
}

func resolveBoxTitle(title string) string {
	if len(title) > 0 {
		return title
	}
	return "Did you mean?"
}

func writeReasonIfPresent(w io.Writer, reason string, width int) {
	if len(reason) == 0 {
		return
	}
	writeBoxLine(w, reason, width)
	writeBoxLine(w, "", width)
}

func writeActionHintIfPresent(w io.Writer, items []Suggestion, width int) {
	if len(items) == 0 {
		return
	}
	first := items[0]
	hint := fmt.Sprintf("Run: %s%s --help%s for options", catppuccinPrompt, first.Command, catppuccinReset)
	writeBoxLine(w, hint, width)
}

func writeBoxTop(w io.Writer, title string, width int) {
	prefix := fmt.Sprintf("╭─ %s%s%s ", catppuccinHeader, title, catppuccinBorder)
	plainPrefix := fmt.Sprintf("╭─ %s ", title)
	rem := width - len(plainPrefix) - 1
	if rem < 2 {
		rem = 2
	}
	fmt.Fprintf(w, "%s%s%s%s%s\n", catppuccinBorder, prefix, strings.Repeat("─", rem), "╮", catppuccinReset)
}

func writeBoxBottom(w io.Writer, width int) {
	fmt.Fprintf(w, "%s╰%s╯%s\n", catppuccinBorder, strings.Repeat("─", width-2), catppuccinReset)
}

func writeBoxLine(w io.Writer, content string, width int) {
	cleanLen := len(stripAnsiRunes(content))
	pad := width - 4 - cleanLen
	if pad < 0 {
		pad = 0
	}
	fmt.Fprintf(w, "%s│%s %s%s %s│%s\n", catppuccinBorder, catppuccinReset, content, strings.Repeat(" ", pad), catppuccinBorder, catppuccinReset)
}

func writeSuggestionItems(w io.Writer, items []Suggestion, width int) {
	for _, s := range items {
		pct := int(math.Round(s.Confidence * 100))
		header := fmt.Sprintf(" %s•%s %s%s%s  %s(Confidence: %d%%)%s", catppuccinHeader, catppuccinReset, catppuccinCommand, s.Command, catppuccinReset, catppuccinDim, pct, catppuccinReset)
		writeBoxLine(w, header, width)
		writeSuggestionDescription(w, s.Description, width)
		writeBoxLine(w, "", width)
	}
}

func writeSuggestionDescription(w io.Writer, desc string, width int) {
	if len(desc) == 0 {
		return
	}
	line := fmt.Sprintf("   %s%s%s", catppuccinDim, desc, catppuccinReset)
	writeBoxLine(w, line, width)
}

func stripAnsiRunes(s string) string {
	var b strings.Builder
	inSeq := false
	for _, r := range s {
		if r == '\033' {
			inSeq = true
			continue
		}
		if inSeq && r == 'm' {
			inSeq = false
			continue
		}
		if !inSeq {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// RenderCompact writes a minimal inline suggestion line.
func RenderCompact(w io.Writer, group SuggestionGroup) error {
	if !group.HasSuggestions() {
		return nil
	}
	parts := buildCompactParts(group.Suggestions)
	line := fmt.Sprintf("Did you mean: %s\n", strings.Join(parts, " | "))
	_, err := fmt.Fprint(w, line)
	return err
}

func buildCompactParts(items []Suggestion) []string {
	var out []string
	for _, s := range items {
		pct := int(math.Round(s.Confidence * 100))
		out = append(out, fmt.Sprintf("%s (%d%%)", s.Command, pct))
	}
	return out
}

// RenderJSON writes structured JSON formatted suggestions.
func RenderJSON(w io.Writer, group SuggestionGroup) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(group)
}
