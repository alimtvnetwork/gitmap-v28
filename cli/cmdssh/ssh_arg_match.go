package cmdssh

import (
	"slices"
)

// argMatches reports whether arg equals any of the given candidates.
// It replaces chains of arg == comparisons with a single slice lookup.
func argMatches(arg string, candidates ...string) bool {
	return slices.Contains(candidates, arg)
}

// argMatchesFlag reports whether arg matches a long flag ("--"+longFlag)
// or its short form ("-"+shortFlag). An empty short flag never matches,
// so a bare "-" cannot match by accident.
func argMatchesFlag(arg, longFlag, shortFlag string) bool {
	if argMatches(arg, "--"+longFlag) {
		return true
	}
	return shortFlag != "" && argMatches(arg, "-"+shortFlag)
}
