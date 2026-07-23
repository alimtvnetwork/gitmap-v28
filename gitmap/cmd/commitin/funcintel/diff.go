package funcintel

import (
	"sort"
	"strings"
)

func addedNames(prevSet, newSet map[string]struct{}) []string {
	out := make([]string, 0, len(newSet))
	for name := range newSet {
		if _, had := prevSet[name]; !had {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}

func extractByRegex(src string, lineMatch func(line string) (string, bool)) map[string]struct{} {
	set := make(map[string]struct{})
	for _, line := range strings.Split(src, "\n") {
		if name, ok := lineMatch(line); ok {
			set[name] = struct{}{}
		}
	}
	return set
}
