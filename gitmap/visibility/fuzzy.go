// Package visibility — fuzzy.go: forgiving fallbacks when a literal
// (no `*`) pattern matches zero repos. Two strategies:
//
//  1. AutoFixVDigitPatterns — silently rewrites `<base>-<N>` to
//     `<base>-v<N>` (covers the common `macro-ahk-51` typo where the
//     user dropped the `v` separator) and adds the highest-vN sibling
//     of a bare `<base>` token.
//  2. NearMisses — surfaces a small ranked list (≤ topN) of candidate
//     repo names whose Levenshtein distance to the pattern is ≤ 3 so
//     the caller can prompt the user to pick one.
//
// Pure functions, zero I/O.
package visibility

import "strings"

// AutoFixVDigitPatterns inspects each compiled pattern; for every
// literal (no `*`) raw form it derives candidate rewrites that have
// at least one match in `names` but the original did not. Returns the
// additions in input order, deduplicated against the original list.
func AutoFixVDigitPatterns(patterns []Pattern, names []string) []Pattern {
	have := make(map[string]bool, len(patterns))
	for _, p := range patterns {
		have[p.Raw] = true
	}
	nameSet := make(map[string]bool, len(names))
	for _, n := range names {
		nameSet[n] = true
	}

	out := make([]Pattern, 0, len(patterns))
	for _, p := range patterns {
		if strings.Contains(p.Raw, "*") {
			continue
		}
		if nameSet[p.Raw] {
			continue // original already matches; no fix needed
		}
		for _, cand := range vDigitCandidates(p.Raw, nameSet, names) {
			if have[cand] {
				continue
			}
			fixed, err := ParsePattern(cand)
			if err != nil {
				continue
			}
			have[cand] = true
			out = append(out, fixed)
		}
	}

	return out
}

// vDigitCandidates returns rewrites worth trying for a single literal
// raw form. Two cases:
//   - "<base>-<digits>" → "<base>-v<digits>"
//   - "<base>"          → highest "<base>-vN" present in `names`
func vDigitCandidates(raw string, nameSet map[string]bool, names []string) []string {
	out := make([]string, 0, 2)
	if base, digits, ok := SplitTrailingDigits(raw); ok {
		fixed := base + "-v" + digits
		if nameSet[fixed] {
			out = append(out, fixed)
		}
	}
	if best, _, ok := HighestVersionedMatch(names, raw); ok {
		out = append(out, best)
	}

	return out
}

// NearMisses returns up to topN repo names from `names` whose
// Levenshtein distance to any literal pattern is ≤ maxDist, sorted
// by ascending distance. Skips patterns that contain `*`.
func NearMisses(patterns []Pattern, names []string, maxDist, topN int) []string {
	type scored struct {
		name string
		dist int
	}
	scoredAll := make([]scored, 0, topN*2)
	for _, p := range patterns {
		if strings.Contains(p.Raw, "*") {
			continue
		}
		for _, n := range names {
			d := levenshtein(p.Raw, n)
			if d > 0 && d <= maxDist {
				scoredAll = append(scoredAll, scored{name: n, dist: d})
			}
		}
	}

	// Stable insertion sort by dist, dedup names.
	seen := make(map[string]bool, len(scoredAll))
	out := make([]string, 0, topN)
	for len(out) < topN && len(scoredAll) > 0 {
		best := 0
		for i := 1; i < len(scoredAll); i++ {
			if scoredAll[i].dist < scoredAll[best].dist {
				best = i
			}
		}
		pick := scoredAll[best]
		scoredAll = append(scoredAll[:best], scoredAll[best+1:]...)
		if seen[pick.name] {
			continue
		}
		seen[pick.name] = true
		out = append(out, pick.name)
	}

	return out
}

// levenshtein returns the edit distance between a and b. Plain DP,
// O(len(a)*len(b)) time, O(len(b)) space. Used only on short repo
// names so the cost is irrelevant.
func levenshtein(a, b string) int {
	if a == b {
		return 0
	}
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}
	prev := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	curr := make([]int, len(b)+1)
	for i := 1; i <= len(a); i++ {
		curr[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = min3(curr[j-1]+1, prev[j]+1, prev[j-1]+cost)
		}
		prev, curr = curr, prev
	}

	return prev[len(b)]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}

		return c
	}
	if b < c {
		return b
	}

	return c
}
