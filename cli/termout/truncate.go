package termout

import (
	"strings"
)

// TruncateMiddle shortens text exceeding maxWidth by placing ellipsis in the middle.
func TruncateMiddle(text string, maxWidth int, ellipsis string) string {
	if maxWidth <= 0 {
		return ""
	}
	if len(text) <= maxWidth {
		return text
	}
	if len(ellipsis) == 0 {
		ellipsis = "..."
	}
	if maxWidth <= len(ellipsis) {
		return ellipsis[:maxWidth]
	}

	return assembleMiddleTruncation(text, maxWidth, ellipsis)
}

func assembleMiddleTruncation(text string, maxWidth int, ellipsis string) string {
	avail := maxWidth - len(ellipsis)
	prefixLen := (avail + 1) / 2
	suffixLen := avail - prefixLen
	prefix := text[:prefixLen]
	suffix := text[len(text)-suffixLen:]

	return prefix + ellipsis + suffix
}

func resolveRuneVisualWidth(r rune) int {
	if isDoubleWidthRune(r) {
		return 2
	}

	return 1
}

func assembleRightTruncation(text string, avail int) string {
	var sb strings.Builder
	curWidth := 0
	for _, r := range text {
		rw := resolveRuneVisualWidth(r)
		if curWidth+rw > avail {
			break
		}
		sb.WriteRune(r)
		curWidth += rw
	}
	sb.WriteString("…")

	return sb.String()
}

// TruncateRightWithEllipsis safely truncates text to maxWidth columns appending an ellipsis.
func TruncateRightWithEllipsis(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	if VisualStringWidth(text) <= maxWidth {
		return text
	}
	if maxWidth <= 1 {
		return "…"
	}

	return assembleRightTruncation(text, maxWidth-1)
}

