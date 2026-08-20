// Package render — prettypost.go: cosmetic post-processing layer applied
// on top of Render() output when emitting to a terminal.
package render

import (
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// applyANSIPost runs cosmetic passes on the output of RenderANSI.
func applyANSIPost(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = transformLine(line)
	}

	return strings.Join(lines, "\n")
}

// transformLine dispatches per-line transforms. Headings short-circuit;
// everything else runs the inline cosmetic cascade.
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
	line = renderExamplePrompt(line)

	return line
}

var headingRe = regexp.MustCompile(`^(#{1,6})\s+(.+?)\s*$`)

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
	default:
		return constants.ColorWhite + text + constants.ColorReset, true
	}
}

func stripInlineMarkers(s string) string {
	s = strings.ReplaceAll(s, "**", "")
	s = strings.ReplaceAll(s, "`", "")

	return s
}

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

var tableSepRe = regexp.MustCompile(`^\s*\|(\s*:?-{2,}:?\s*\|)+\s*$`)

func colorTableSeparator(line string) string {
	if !tableSepRe.MatchString(line) {
		return line
	}

	return constants.ColorDim + line + constants.ColorReset
}

func colorTablePipes(line string) string {
	t := strings.TrimLeft(line, " ")
	if !strings.HasPrefix(t, "|") || tableSepRe.MatchString(line) {
		return line
	}

	return strings.ReplaceAll(line, "|",
		constants.ColorDim+"│"+constants.ColorReset)
}
