package message

import "strings"

func applyTitleAffix(msg, prefix, suffix string) string {
	if prefix == "" && suffix == "" {
		return msg
	}

	idx := strings.IndexByte(msg, '\n')
	if idx < 0 {
		return prefix + msg + suffix
	}

	return prefix + msg[:idx] + suffix + msg[idx:]
}

func applyBodyAffix(msg string, prefixPool, suffixPool []string, pick func(int) int) string {
	chosenPrefix, chosenSuffix := pickAffixPair(prefixPool, suffixPool, pick)
	return applyChosenBodyAffix(msg, chosenPrefix, chosenSuffix)
}

func pickAffixPair(prefixPool, suffixPool []string, pick func(int) int) (string, string) {
	var chosenPrefix, chosenSuffix string
	if len(prefixPool) > 0 {
		chosenPrefix = pickOne(prefixPool, pick)
	}
	if len(suffixPool) > 0 {
		chosenSuffix = pickOne(suffixPool, pick)
	}

	return chosenPrefix, chosenSuffix
}

func applyChosenBodyAffix(msg, chosenPrefix, chosenSuffix string) string {
	if chosenPrefix != "" {
		msg = prependBodyBlock(msg, chosenPrefix)
	}
	if chosenSuffix != "" {
		msg = appendBodyBlock(msg, chosenSuffix)
	}

	return msg
}

func prependBodyBlock(msg, prefix string) string {
	isMarkdownBlock := strings.HasPrefix(strings.TrimSpace(prefix), "#") || strings.Contains(prefix, "\n")
	if !isMarkdownBlock {
		return prefix + "\n" + msg
	}
	idx := strings.IndexByte(msg, '\n')
	if idx < 0 {
		return strings.TrimRight(msg, " \t") + "\n\n" + strings.TrimSpace(prefix)
	}
	title := strings.TrimRight(msg[:idx], " \t")
	body := strings.TrimSpace(msg[idx+1:])
	if body == "" {
		return title + "\n\n" + strings.TrimSpace(prefix)
	}

	return title + "\n\n" + strings.TrimSpace(prefix) + "\n\n" + body
}

func appendBodyBlock(msg, suffix string) string {
	trimmedMsg := strings.TrimRight(msg, " \t\n")
	trimmedSuffix := strings.TrimSpace(suffix)
	if trimmedMsg == "" {
		return trimmedSuffix
	}

	return trimmedMsg + "\n\n" + trimmedSuffix
}

func pickOne(pool []string, pick func(int) int) string {
	if pick == nil {
		return pool[0]
	}

	i := pick(len(pool))
	if i < 0 || i >= len(pool) {
		i = 0
	}

	return pool[i]
}
