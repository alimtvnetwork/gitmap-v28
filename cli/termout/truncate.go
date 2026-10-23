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

// TruncateMiddleWithEllipsis truncates text exceeding maxWidth by placing a "..." ellipsis
// in the middle, calculated via rune visual display width to avoid character tearing.
func TruncateMiddleWithEllipsis(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	if VisualStringWidth(text) <= maxWidth {
		return text
	}
	if maxWidth <= 3 {
		return "..."[:maxWidth]
	}

	return assembleMiddleEllipsis(text, maxWidth)
}

func assembleMiddleEllipsis(text string, maxWidth int) string {
	avail := maxWidth - 3
	prefixBudget := (avail + 1) / 2
	suffixBudget := avail - prefixBudget
	runes := []rune(text)
	prefix := slicePrefixRunes(runes, prefixBudget)
	suffix := sliceSuffixRunes(runes, suffixBudget)

	return prefix + "..." + suffix
}

func slicePrefixRunes(runes []rune, budget int) string {
	var sb strings.Builder
	curWidth := 0
	for _, r := range runes {
		rw := resolveRuneVisualWidth(r)
		if curWidth+rw > budget {
			break
		}
		sb.WriteRune(r)
		curWidth += rw
	}

	return sb.String()
}

func sliceSuffixRunes(runes []rune, budget int) string {
	var collected []rune
	curWidth := 0
	for i := len(runes) - 1; i >= 0; i-- {
		r := runes[i]
		rw := resolveRuneVisualWidth(r)
		if curWidth+rw > budget {
			break
		}
		collected = append(collected, r)
		curWidth += rw
	}

	return reverseRunesToString(collected)
}

func reverseRunesToString(runes []rune) string {
	var sb strings.Builder
	for i := len(runes) - 1; i >= 0; i-- {
		sb.WriteRune(runes[i])
	}

	return sb.String()
}

// TruncateLeftWithEllipsis shortens text exceeding maxWidth by placing "..." at the left,
// preserving the rightmost characters based on rune visual display width.
func TruncateLeftWithEllipsis(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	if VisualStringWidth(text) <= maxWidth {
		return text
	}
	if maxWidth <= 3 {
		return "..."[:maxWidth]
	}

	return assembleLeftEllipsis(text, maxWidth)
}

func assembleLeftEllipsis(text string, maxWidth int) string {
	budget := maxWidth - 3
	runes := []rune(text)
	suffix := sliceSuffixRunes(runes, budget)

	return "..." + suffix
}
