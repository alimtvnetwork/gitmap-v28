// Package render — prettypost.go: cosmetic post-processing layer applied
// on top of Render() output when emitting to a terminal.
package render

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// applyANSIPost runs cosmetic passes on the output of RenderANSI.
func applyANSIPost(s string) string {
	lines := strings.Split(s, "\n")
	bannerRendered := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !bannerRendered && trimmed != "" {
			if strings.HasPrefix(trimmed, "# ") {
				lines[i] = formatBoxBanner(strings.TrimPrefix(trimmed, "# "))
				bannerRendered = true
				continue
			} else if !strings.HasPrefix(trimmed, "#") {
				lines[i] = formatBoxBanner(trimmed)
				bannerRendered = true
				continue
			}
		}
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
	if t, ok := transformTableRow(line); ok {
		return t
	}

	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "```") {
		return ""
	}
	if strings.HasPrefix(trimmed, "#") && !headingRe.MatchString(line) {
		return "    " + constants.ColorDim + trimmed + constants.ColorReset
	}
	if strings.HasPrefix(trimmed, "gitmap ") || strings.HasPrefix(trimmed, "$ gitmap ") {
		cmdStr := strings.TrimPrefix(trimmed, "$ ")
		return "    " + constants.ColorCyan + cmdStr + constants.ColorReset
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
		return formatBoxBanner(text), true
	case 2:
		title := strings.TrimSuffix(text, ":") + ":"
		return "\n  " + constants.ColorYellow + title + constants.ColorReset, true
	case 3:
		return "\n  " + constants.ColorMagenta + "› " + text + constants.ColorReset, true
	default:
		return "  " + constants.ColorWhite + text + constants.ColorReset, true
	}
}

func formatBoxBanner(text string) string {
	text = strings.Trim(text, "`'\" ")
	width := len(text) + 6
	if width < 50 {
		width = 50
	}
	if width > 80 {
		width = 80
	}
	border := strings.Repeat("═", width)
	pad := (width - len(text)) / 2
	if pad < 1 {
		pad = 1
	}
	leftSpace := strings.Repeat(" ", pad)
	rightSpace := strings.Repeat(" ", width-len(text)-pad)
	if len(rightSpace) < 0 {
		rightSpace = ""
	}

	return fmt.Sprintf("\n  %s╔%s╗%s\n  %s║%s%s%s║%s\n  %s╚%s╝%s\n",
		constants.ColorCyan, border, constants.ColorReset,
		constants.ColorCyan, leftSpace, text, rightSpace, constants.ColorReset,
		constants.ColorCyan, border, constants.ColorReset,
	)
}

func transformTableRow(line string) (string, bool) {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "|") || !strings.HasSuffix(trimmed, "|") {
		return "", false
	}
	if tableSepRe.MatchString(trimmed) {
		return "", true
	}
	parts := strings.Split(trimmed, "|")
	if len(parts) < 3 {
		return "", false
	}
	col1 := stripInlineMarkers(strings.TrimSpace(parts[1]))
	col2 := stripInlineMarkers(strings.TrimSpace(parts[2]))
	if len(parts) >= 4 && strings.TrimSpace(parts[3]) != "" {
		col3 := stripInlineMarkers(strings.TrimSpace(parts[3]))
		if strings.EqualFold(col1, "Flag") || strings.EqualFold(col1, "Command") {
			return fmt.Sprintf("    %s%-24s %-12s %s%s", constants.ColorDim, strings.ToUpper(col1), strings.ToUpper(col2), strings.ToUpper(col3), constants.ColorReset), true
		}
		return fmt.Sprintf("    %s%-24s%s %s%-12s%s %s", constants.ColorCyan, col1, constants.ColorReset, constants.ColorYellow, col2, constants.ColorReset, col3), true
	}
	if strings.EqualFold(col1, "Flag") || strings.EqualFold(col1, "Command") || strings.EqualFold(col1, "Action") {
		return fmt.Sprintf("    %s%-28s%s %s%s%s", constants.ColorDim, strings.ToUpper(col1), constants.ColorReset, constants.ColorDim, strings.ToUpper(col2), constants.ColorReset), true
	}
	return fmt.Sprintf("    %s%-28s%s %s", constants.ColorCyan, col1, constants.ColorReset, col2), true
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
