package cmd

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TestRenderInlineMarkdownHighlightsDoubleQuotes locks in that changelog
// bullet bodies route their text through the shared pretty-markdown
// quote-highlight rule (render.HighlightQuotesANSI) so the formatting
// matches `gitmap help` output. Regression guard: previously bullets only
// rendered **bold** / `code` and ignored "double quotes" entirely.
func TestRenderInlineMarkdownHighlightsDoubleQuotes(t *testing.T) {
	got := renderInlineMarkdown(`Renamed "old" to "new"`, 0, true)

	if !strings.Contains(got, constants.ColorCyan) {
		t.Fatalf("expected cyan ANSI for double-quoted spans, got %q", got)
	}
	if !strings.Contains(got, constants.ColorReset) {
		t.Fatalf("expected ColorReset to close cyan span, got %q", got)
	}
	// Sentinel tokens must not leak into terminal output.
	if strings.Contains(got, "[C]") || strings.Contains(got, "[/C]") {
		t.Fatalf("token sentinels leaked into ANSI output: %q", got)
	}
}

// TestRenderInlineMarkdownLeavesApostrophesAlone protects the single-quote
// passthrough rule — bullets like "user's repo" must not gain stray ANSI.
func TestRenderInlineMarkdownLeavesApostrophesAlone(t *testing.T) {
	got := renderInlineMarkdown(`user's repo`, 0, true)

	if strings.Contains(got, constants.ColorCyan) {
		t.Fatalf("apostrophes must not trigger cyan styling: %q", got)
	}
}

// TestRenderInlineMarkdownPreservesBoldAndCode keeps the existing
// bold/code behavior intact after wiring in the quote-highlight pass.
func TestRenderInlineMarkdownPreservesBoldAndCode(t *testing.T) {
	got := renderInlineMarkdown("Use **bold** and `code` here", 0, true)

	if !strings.Contains(got, constants.ChangelogPrettyBoldOpen) {
		t.Fatalf("bold open marker missing: %q", got)
	}
	if !strings.Contains(got, constants.ChangelogPrettyCodeOpen) {
		t.Fatalf("code open marker missing: %q", got)
	}
}

// TestRenderInlineMarkdownPlainModeStripsAllANSI is the new --no-pretty
// regression guard. Output redirected to a file or piped into `grep -F`
// must contain zero ESC bytes — that's the entire point of the flag.
// Also asserts the markdown punctuation is stripped so plain readers
// don't see literal `**` / backticks bleeding through.
func TestRenderInlineMarkdownPlainModeStripsAllANSI(t *testing.T) {
	got := renderInlineMarkdown(`Use **bold**, `+"`code`"+` and "quotes"`, 0, false)

	if strings.ContainsRune(got, '\x1b') {
		t.Fatalf("plain mode must not emit ESC bytes: %q", got)
	}
	if strings.Contains(got, "**") {
		t.Fatalf("plain mode must strip bold delimiters: %q", got)
	}
	if strings.Contains(got, "`") {
		t.Fatalf("plain mode must strip code delimiters: %q", got)
	}
	// Quotes are part of the prose, not formatting — leave them in.
	if !strings.Contains(got, `"quotes"`) {
		t.Fatalf("plain mode must preserve literal quote characters: %q", got)
	}
}

// TestColorOrEmptyTogglesByPrettyFlag verifies the central ANSI-strip
// helper used by every header / marker site. A single broken case here
// would break headers, bullet markers, and the rule line all at once.
func TestColorOrEmptyTogglesByPrettyFlag(t *testing.T) {
	if got := colorOrEmpty(constants.ColorCyan, true); got != constants.ColorCyan {
		t.Errorf("pretty=true should preserve color, got %q", got)
	}
	if got := colorOrEmpty(constants.ColorCyan, false); got != "" {
		t.Errorf("pretty=false should return empty string, got %q", got)
	}
}
