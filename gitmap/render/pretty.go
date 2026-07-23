// Package render contains terminal-output helpers shared across CLI
// commands. The pretty markdown renderer (pretty.go) is the primary
// caller-facing API.
package render

import (
	"bufio"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// Token sentinels used by Render() so unit tests can assert on a stable,
// ANSI-free string. RenderANSI swaps these for real escape codes.
const (
	TokYellowOpen  = "[Y]"
	TokYellowClose = "[/Y]"
	TokCyanOpen    = "[C]"
	TokCyanClose   = "[/C]"
	TokMutedOpen   = "[M]"
	TokMutedClose  = "[/M]"
	TokGreenOpen   = "[G]"
	TokGreenClose  = "[/G]"
	TokMagentaOpen = "[P]"
	TokMagentaClose = "[/P]"

	collapseArrow = "→ "
	bodyIndent    = "  "
)

// Render applies the four pretty-print rules to markdown input and returns
// a token-annotated string (no ANSI escape codes). Use RenderANSI when
// emitting to a terminal.
//
// Rules:
//  1. A fenced code block whose body matches the immediately-preceding
//     paragraph's text collapses to a single yellow "→ <content>" line and
//     the fence is dropped.
//  2. Double-quoted strings ("like this") become cyan; single quotes are
//     left alone (apostrophes).
//  3. An italic line directly under a heading renders as a muted subtitle.
//  4. Body content under a heading is indented by two spaces.
func Render(md string) string {
	lines := splitLines(md)
	doc := parse(lines)

	var out strings.Builder
	prevBlank := false
	for _, b := range doc {
		if b.kind == bkBlank && prevBlank {
			continue
		}
		emitBlock(&out, b)
		prevBlank = b.kind == bkBlank
	}

	return strings.TrimRight(out.String(), "\n") + "\n"
}

// RenderANSI is Render with ANSI codes substituted for the token sentinels
// and the cosmetic post-processor (headings/bold/backticks/links/tables)
// applied so terminal output looks like a polished help screen, not raw
// markdown. The post pass runs only here — Render() stays token-pure so
// the fixture suite is unaffected.
func RenderANSI(md string) string {
	return applyANSIPost(tokenToANSI().Replace(Render(md)))
}

// HighlightQuotesANSI applies the cyan double-quote rule (rule 2 of the
// pretty renderer) to a single string and returns it with real ANSI escape
// codes already substituted. Useful for callers that have their own
// block-level layout (e.g. the changelog bullet renderer) but still want
// quote highlighting consistent with the help-text pretty renderer.
func HighlightQuotesANSI(s string) string {
	return tokenToANSI().Replace(HighlightQuotes(s))
}

// tokenToANSI builds the shared sentinel→ANSI replacer used by
// RenderANSI and HighlightQuotesANSI. Centralized so the mapping stays
// in lockstep across pretty-render entry points.
func tokenToANSI() *strings.Replacer {
	return strings.NewReplacer(
		TokYellowOpen, constants.ColorYellow,
		TokYellowClose, constants.ColorReset,
		TokCyanOpen, constants.ColorCyan,
		TokCyanClose, constants.ColorReset,
		TokMutedOpen, constants.ColorDim,
		TokMutedClose, constants.ColorReset,
		TokGreenOpen, constants.ColorGreen,
		TokGreenClose, constants.ColorReset,
		TokMagentaOpen, constants.ColorMagenta,
		TokMagentaClose, constants.ColorReset,
	)
}

// splitLines splits on \n and drops the trailing empty token from a final
// newline. Preserves blank lines in between.
func splitLines(s string) []string {
	var out []string
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		out = append(out, scanner.Text())
	}

	return out
}
