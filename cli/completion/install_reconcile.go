package completion

import (
	"strings"
)

func reconcileUnixProfile(normExisting, normSnippet, rawExisting, eol string) (string, bool) {
	if strings.Contains(normExisting, compUnixMarkerStart) && strings.Contains(normExisting, compUnixMarkerEnd) {
		replaced := replaceExistingBlock(normExisting, normSnippet, compUnixMarkerStart, compUnixMarkerEnd)
		final := toNativeEOL(replaced, eol)

		return final, final != rawExisting
	}

	appended := appendSnippetToProfile(normExisting, normSnippet)
	final := toNativeEOL(appended, eol)

	return final, final != rawExisting
}

func replaceExistingBlock(content, snippet, startMarker, endMarker string) string {
	startIndex := strings.Index(content, startMarker)
	if startIndex < 0 {
		return content
	}

	endIndex := strings.Index(content, endMarker)
	if endIndex < 0 {
		return content
	}

	endIndex += len(endMarker)

	return content[:startIndex] + snippet + content[endIndex:]
}

func appendSnippetToProfile(base, snippet string) string {
	trimmed := strings.TrimRight(base, " \t\r\n")
	if len(trimmed) == 0 {
		return snippet + "\n"
	}

	return trimmed + "\n\n" + snippet + "\n"
}

func stripLegacyPowerShellBlocks(content string) string {
	lines := strings.Split(content, "\n")
	cleanLines := make([]string, 0, len(lines))
	isInsideLegacy := false
	braceDepth := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == compLegacyMarker {
			isInsideLegacy = true
			braceDepth = 0
			continue
		}
		if !isInsideLegacy {
			cleanLines = append(cleanLines, line)
			continue
		}
		braceDepth += strings.Count(trimmed, "{") - strings.Count(trimmed, "}")
		if isLegacyLineToDrop(trimmed, braceDepth) {
			continue
		}
		isInsideLegacy = false
		cleanLines = append(cleanLines, line)
	}

	return strings.Join(cleanLines, "\n")
}

func isLegacyLineToDrop(trimmed string, braceDepth int) bool {
	if strings.HasPrefix(trimmed, ".") && strings.Contains(trimmed, "gitmap") {
		return true
	}

	if strings.HasPrefix(trimmed, "if ((Get-Module") {
		return true
	}

	if braceDepth > 0 {
		return true
	}

	if strings.HasPrefix(trimmed, "}") {
		return true
	}

	return false
}

func collapseBlankLines(text string) string {
	for strings.Contains(text, "\n\n\n") {
		text = strings.ReplaceAll(text, "\n\n\n", "\n\n")
	}

	return text
}

func toNativeEOL(content, eol string) string {
	normalized := strings.ReplaceAll(content, "\r\n", "\n")
	if eol == "\r\n" {
		return strings.ReplaceAll(normalized, "\n", "\r\n")
	}

	return normalized
}
