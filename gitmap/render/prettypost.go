// Package render — prettypost.go: cosmetic post-processing layer applied
// on top of Render() output when emitting to a terminal. Lives outside
// Render() so the token-fixture tests (testdata/pretty/*) keep passing —
// they assert on the sentinel-token stream, while this layer only runs
// inside RenderANSI just before the final string returns.
//
// Transforms applied (in order, all idempotent against already-ANSI text):
//  1. Strip backslash escapes used by markdown ('\<', '\>', '\|', '\\').
//  2. Heading lines ('# Title', '## Section', …) → colored, marker-free
//     with a left bar so sections pop on a dense help screen.
//  3. Inline `code` spans → magenta, backticks dropped.
//  4. Inline **bold** spans → bright white, asterisks dropped.
//  5. Markdown links '[text](url)' → cyan text + dim '(url)'.
//  6. Table separator rows ('|---|---|') → dim.
//  7. Table header/body pipe characters → dim, so the columns visually
//     recede behind the actual content.
package render

import (
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// applyANSIPost runs every cosmetic pass on the swap-layer output of
// RenderANSI. Operating on the already-ANSI string keeps the sentinel
// fixture suite unaffected and centralizes terminal-only polish.
func applyANSIPost(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = transformLine(line)
	}

	return strings.Join(lines, "\n")
}

// transformLine dispatches per-line transforms. Headings short-circuit
// (they own the whole line); everything else runs the inline cascade.
func transformLine(line string) string {
	if h, ok := renderHeadingLine(line); ok {
		return h
	}
	line = unescapeMarkdown(line)
	line = colorTableSeparator(line)
	line = colorTablePipes(line)
	line = renderInlineCode(line)
	line = renderInlineBold(line)
	line = renderInlineLinks(line)
	line = renderCommandAliasRow(line)
	line = renderBareFlagTokens(line)
	line = renderAnglePlaceholders(line)
	line = renderDefaultParen(line)

	return line
}

// commandAliasRowRe matches bulleted command-list rows shaped like:
//
//	"  - clone (cl, cln)  Clone repos from gitmap.json"
//
// Captures: leading bullet, command name, alias list, trailing tail.
// Command → green/bold, aliases → yellow/bold. Tail keeps its colour
// (may already carry inline-code magenta from earlier passes).
var commandAliasRowRe = regexp.MustCompile(`^(\s*[-*]\s+)([a-z][a-z0-9-]*)\s+\(([^)]+)\)(\s+.*)?$`)

func renderCommandAliasRow(line string) string {
	m := commandAliasRowRe.FindStringSubmatch(line)
	if m == nil {
		return line
	}
	tail := ""
	if len(m) >= 5 {
		tail = m[4]
	}

	return m[1] +
		constants.ColorGreen + m[2] + constants.ColorReset +
		" (" + constants.ColorYellow + m[3] + constants.ColorReset + ")" +
		tail
}

// bareFlagRe matches `--long-flag` or `-x` tokens NOT already wrapped
// in an ANSI escape (guarded by negative-lookbehind via prefix check).
// Word-boundary on both sides so URLs like `https://x?--a=b` won't
// trigger. We avoid an outright lookbehind (Go RE2 has none) by
// excluding ESC and word chars before the dash via a non-capturing
// prefix that is preserved verbatim.
var bareFlagRe = regexp.MustCompile(`(^|[\s(\[,/])(-{1,2}[A-Za-z][A-Za-z0-9-]*)`)

func renderBareFlagTokens(line string) string {
	return bareFlagRe.ReplaceAllString(line,
		"$1"+constants.ColorCyan+"$2"+constants.ColorReset)
}

// anglePlaceholderRe matches `<value>` placeholders used in usage
// strings (e.g. `--filter <query>`). Green so the eye pairs them with
// the cyan flag immediately to their left.
var anglePlaceholderRe = regexp.MustCompile(`<([a-zA-Z][a-zA-Z0-9_-]*)>`)

func renderAnglePlaceholders(line string) string {
	return anglePlaceholderRe.ReplaceAllString(line,
		constants.ColorGreen+"<$1>"+constants.ColorReset)
}

// defaultParenRe matches "(default: …)" and "(default …)" trailers.
// Dim so they sit visually behind the flag description.
var defaultParenRe = regexp.MustCompile(`\((default:?\s+[^)]+)\)`)

func renderDefaultParen(line string) string {
	return defaultParenRe.ReplaceAllString(line,
		constants.ColorDim+"($1)"+constants.ColorReset)
}

// headingRe matches a leading '#' run at column 0. Help bodies are
// indented two spaces by emitBlock, so a column-0 '#' is always a heading.
var headingRe = regexp.MustCompile(`^(#{1,6})\s+(.+?)\s*$`)

// renderHeadingLine swaps a raw '## Heading' line for a colored,
// marker-free presentation. H1=cyan bar, H2=yellow bar, H3=magenta arrow,
// H4+=plain bold. Returns (rendered, true) on hit so the caller skips
// further inline transforms (a heading text shouldn't be re-styled).
func renderHeadingLine(line string) (string, bool) {
	m := headingRe.FindStringSubmatch(line)
	if m == nil {
		return "", false
	}
	text := unescapeMarkdown(m[2])
	text = stripInlineMarkers(text)
	switch len(m[1]) {
	case 1:
		return constants.ColorCyan + "▌ " + text + constants.ColorReset, true
	case 2:
		return constants.ColorYellow + "▌ " + text + constants.ColorReset, true
	case 3:
		return constants.ColorMagenta + "› " + text + constants.ColorReset, true
	}

	return constants.ColorWhite + text + constants.ColorReset, true
}

// stripInlineMarkers removes the markdown asterisks/backticks from heading
// text so '## Parallel execution (`--max-concurrency`)' renders as
// 'Parallel execution (--max-concurrency)' instead of leaking syntax.
func stripInlineMarkers(s string) string {
	s = strings.ReplaceAll(s, "**", "")
	s = strings.ReplaceAll(s, "`", "")

	return s
}

// markdownEscapes pairs each backslash escape used by help authors with
// its rendered literal. Keeping these in one table makes it trivial to
// extend (e.g. \[ \] for square brackets) without touching the dispatcher.
var markdownEscapes = [][2]string{
	{`\<`, `<`}, {`\>`, `>`}, {`\|`, `|`}, {`\\`, `\`},
	{`\*`, `*`}, {`\_`, `_`}, {`\` + "`", "`"},
}

func unescapeMarkdown(s string) string {
	for _, p := range markdownEscapes {
		s = strings.ReplaceAll(s, p[0], p[1])
	}

	return s
}

// inlineCodeRe matches a `…` span without consuming surrounding text.
// Non-greedy body keeps adjacent spans (`a` and `b`) from merging.
var inlineCodeRe = regexp.MustCompile("`([^`]+?)`")

func renderInlineCode(s string) string {
	return inlineCodeRe.ReplaceAllString(s,
		constants.ColorMagenta+"$1"+constants.ColorReset)
}

// inlineBoldRe matches **bold** spans. Non-greedy and bounded to keep
// stray asterisks (e.g. globs) from being interpreted as bold markers.
var inlineBoldRe = regexp.MustCompile(`\*\*([^*]+?)\*\*`)

func renderInlineBold(s string) string {
	return inlineBoldRe.ReplaceAllString(s,
		constants.ColorWhite+"$1"+constants.ColorReset)
}

// inlineLinkRe matches '[text](url)'. The text becomes cyan, the URL
// fades into dim parens so the eye finds the label first.
var inlineLinkRe = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

func renderInlineLinks(s string) string {
	return inlineLinkRe.ReplaceAllString(s,
		constants.ColorCyan+"$1"+constants.ColorReset+
			constants.ColorDim+" ($2)"+constants.ColorReset)
}

// tableSepRe matches '|---|---|' separator rows (any column count, with
// optional alignment colons). They carry no content, so we dim the
// entire line to push it visually behind the data rows.
var tableSepRe = regexp.MustCompile(`^\s*\|(\s*:?-{2,}:?\s*\|)+\s*$`)

func colorTableSeparator(line string) string {
	if !tableSepRe.MatchString(line) {
		return line
	}

	return constants.ColorDim + line + constants.ColorReset
}

// colorTablePipes dims every '|' character on a table row so the columns
// recede and the cell contents pop. Only runs on lines that look like a
// table row (leading whitespace + '|') and skips separators (already
// fully-dimmed) to avoid nested ANSI nonsense.
func colorTablePipes(line string) string {
	t := strings.TrimLeft(line, " ")
	if !strings.HasPrefix(t, "|") || tableSepRe.MatchString(line) {
		return line
	}

	return strings.ReplaceAll(line, "|",
		constants.ColorDim+"│"+constants.ColorReset)
}
