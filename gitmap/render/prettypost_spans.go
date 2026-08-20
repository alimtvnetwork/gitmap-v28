package render

import (
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

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

var bareFlagRe = regexp.MustCompile(`(^|[\s(\[,/])(-{1,2}[A-Za-z][A-Za-z0-9-]*)`)

func renderBareFlagTokens(line string) string {
	return bareFlagRe.ReplaceAllString(line,
		"$1"+constants.ColorCyan+"$2"+constants.ColorReset)
}

var anglePlaceholderRe = regexp.MustCompile(`<([a-zA-Z][a-zA-Z0-9_-]*)>`)

func renderAnglePlaceholders(line string) string {
	return anglePlaceholderRe.ReplaceAllString(line,
		constants.ColorGreen+"<$1>"+constants.ColorReset)
}

var defaultParenRe = regexp.MustCompile(`\((default:?\s+[^)]+)\)`)

func renderDefaultParen(line string) string {
	return defaultParenRe.ReplaceAllString(line,
		constants.ColorDim+"($1)"+constants.ColorReset)
}

var inlineCodeRe = regexp.MustCompile("`([^`]+?)`")

func renderInlineCode(s string) string {
	return inlineCodeRe.ReplaceAllString(s,
		constants.ColorMagenta+"$1"+constants.ColorReset)
}

var inlineBoldRe = regexp.MustCompile(`\*\*([^*]+?)\*\*`)

func renderInlineBold(s string) string {
	return inlineBoldRe.ReplaceAllString(s,
		constants.ColorWhite+"$1"+constants.ColorReset)
}

var inlineLinkRe = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

func renderInlineLinks(s string) string {
	return inlineLinkRe.ReplaceAllString(s,
		constants.ColorCyan+"$1"+constants.ColorReset+
			constants.ColorDim+" ($2)"+constants.ColorReset)
}

var examplePromptRe = regexp.MustCompile(`^\s*(\$\s+)(gitmap\s+.*)$`)

func renderExamplePrompt(line string) string {
	m := examplePromptRe.FindStringSubmatch(line)
	if m == nil {
		return line
	}
	return "  " + constants.ColorDim + "$ " + constants.ColorReset +
		constants.ColorGreen + "gitmap " + constants.ColorReset +
		constants.ColorWhite + strings.TrimPrefix(m[2], "gitmap ") + constants.ColorReset
}
