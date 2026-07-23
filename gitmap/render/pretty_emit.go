package render

import (
	"regexp"
	"strings"
)

// emitBlock writes one parsed block in pretty form.
func emitBlock(out *strings.Builder, b block) {
	switch b.kind {
	case bkHeading:
		out.WriteString(b.text)
		out.WriteByte('\n')
	case bkSubtitle:
		out.WriteString(bodyIndent)
		out.WriteString(TokMutedOpen)
		out.WriteString(b.text)
		out.WriteString(TokMutedClose)
		out.WriteByte('\n')
	case bkParagraph:
		out.WriteString(bodyIndent)
		out.WriteString(highlightInline(HighlightQuotes(b.text)))
		out.WriteByte('\n')
	case bkFence:
		for _, l := range b.lines {
			out.WriteString(bodyIndent)
			out.WriteString(highlightFenceLine(l))
			out.WriteByte('\n')
		}
	case bkList:
		for _, l := range b.lines {
			out.WriteString(bodyIndent)
			out.WriteString(highlightInline(HighlightQuotes(l)))
			out.WriteByte('\n')
		}
	case bkBlank:
		out.WriteByte('\n')
	}
}

// HighlightQuotes wraps every "double-quoted span" in cyan tokens
// (TokCyanOpen/TokCyanClose). Single quotes are left untouched (they're
// usually apostrophes). The output uses sentinel tokens — call
// HighlightQuotesANSI to get a string with real ANSI escape codes, or
// run the result through RenderANSI's swap layer.
//
// Exported so other CLI surfaces (e.g. the changelog pretty-printer) can
// share the exact same quote-rendering rule that the help-text renderer
// applies, keeping formatting consistent across commands.
func HighlightQuotes(s string) string {
	var b strings.Builder
	inQuote := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c != '"' {
			b.WriteByte(c)

			continue
		}
		if !inQuote {
			b.WriteString(TokCyanOpen)
			b.WriteByte('"')
			inQuote = true

			continue
		}
		b.WriteByte('"')
		b.WriteString(TokCyanClose)
		inQuote = false
	}
	if inQuote { // unterminated quote: close the token defensively
		b.WriteString(TokCyanClose)
	}

	return b.String()
}

// shellCommentRe matches a line that is entirely a shell-style comment
// (leading `#` after optional indent), e.g. `# sync settings`. The whole
// line renders green so command examples read like annotated transcripts.
var shellCommentRe = regexp.MustCompile(`^(\s*)(#.*)$`)

// keyTokenRe matches uppercase identifiers that read like credentials or
// env-var names (API_KEY, GITMAP_TOKEN, GH_PAT, SOMETHING_SECRET, …) so
// they pop out of dense help text. Three-char minimum keeps it from
// painting random short words.
var keyTokenRe = regexp.MustCompile(`\b[A-Z][A-Z0-9_]{2,}(?:KEY|TOKEN|SECRET|PASSWORD|PAT|API)[A-Z0-9_]*\b|\b(?:API|GITMAP|GH|GITHUB|OPENAI)_[A-Z0-9_]+\b`)

// hdAliasRe matches the standalone `hd` help-doc alias so it pops next
// to its full `help-docs` name. Word-bounded so it doesn't eat HD inside
// other words.
var hdAliasRe = regexp.MustCompile(`\bhd\b`)

// highlightFenceLine paints fenced/code lines:
//   - whole-line `# comments` go green
//   - inline KEY=value identifiers go magenta
//   - the `hd` alias goes magenta
func highlightFenceLine(line string) string {
	if m := shellCommentRe.FindStringSubmatch(line); m != nil {
		return m[1] + TokGreenOpen + m[2] + TokGreenClose
	}

	return highlightInline(line)
}

// highlightInline applies the credential + `hd` highlight rules to any
// non-comment string. Idempotent against already-tokenized output: the
// regexes only match raw ASCII identifiers, not sentinel tokens.
func highlightInline(s string) string {
	s = keyTokenRe.ReplaceAllStringFunc(s, func(m string) string {
		return TokMagentaOpen + m + TokMagentaClose
	})
	s = hdAliasRe.ReplaceAllStringFunc(s, func(m string) string {
		return TokMagentaOpen + m + TokMagentaClose
	})

	return s
}
