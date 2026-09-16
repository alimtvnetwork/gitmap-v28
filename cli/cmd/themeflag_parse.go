package cmd

// tryConsumeThemeArg inspects argv at index i for theme flag variants.
func tryConsumeThemeArg(args []string, i int, short, long string) (int, bool) {
	a := args[i]
	if val, isFound := parseThemeEqual(a, short, long); isFound {
		applyThemeChoice(val)
		return i, true
	}

	if a == short || a == long {
		return consumeThemeVal(args, i)
	}

	return i, false
}

// consumeThemeVal handles positional theme value following -theme/--theme.
func consumeThemeVal(args []string, i int) (int, bool) {
	if i+1 < len(args) {
		applyThemeChoice(args[i+1])
		return i + 1, true
	}

	applyThemeChoice("")
	return i, true
}

// parseThemeEqual checks for `-theme=val` or `--theme=val`.
func parseThemeEqual(a, short, long string) (string, bool) {
	if val, hasPrefix := stripThemePrefix(a, short+"="); hasPrefix {
		return val, true
	}

	return stripThemePrefix(a, long+"=")
}

// stripThemePrefix helper returns the remainder if a begins with prefix.
func stripThemePrefix(a, prefix string) (string, bool) {
	if len(a) < len(prefix) || a[:len(prefix)] != prefix {
		return "", false
	}

	return a[len(prefix):], true
}
