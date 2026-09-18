package termtable

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
