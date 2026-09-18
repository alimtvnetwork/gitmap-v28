package termhelp

// CalculateMaxCommandWidth calculates the optimal left column width based on entries.
func CalculateMaxCommandWidth(
	sections []HelpSection,
	minWidth int,
	maxWidth int,
) int {
	maxLen := minWidth
	for _, sec := range sections {
		maxLen = evaluateSectionEntries(sec.Entries, maxLen)
	}
	if maxWidth > 0 && maxLen > maxWidth {
		return maxWidth
	}

	return maxLen
}

func evaluateSectionEntries(entries []CommandEntry, currentMax int) int {
	result := currentMax
	for _, e := range entries {
		cmdLen := len(e.Command)
		if e.HasSubcommands {
			cmdLen += 3
		}
		if cmdLen > result {
			result = cmdLen
		}
	}

	return result
}
