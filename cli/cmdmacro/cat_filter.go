// Package cmdmacro - cat_filter.go: line-window, query, and token filtering for gitmap cat.
// Split from cat_reader.go on 31-Oct-26 so cat_reader.go stays under the 300-line file-size gate.
package cmdmacro

import (
	"strings"
)

func resolveCatMatches(lines []string, opts CatOptions) []LineMatch {
	// WHY combined guards: the inner Query check inside the
	// HasLineRange branch tripped the nested-if linter.
	if opts.HasLineRange && opts.Query != "" {
		return filterLineMatchesByQuery(extractLinesByWindow(lines, opts), opts.Query, opts.Limit)
	}

	if opts.HasLineRange {
		return extractLinesByWindow(lines, opts)
	}

	if opts.Query != "" {
		return filterLinesWithWindow(lines, opts)
	}

	return extractLinesByWindow(lines, opts)
}

func filterLineMatchesByQuery(window []LineMatch, query string, limit int) []LineMatch {
	var matches []LineMatch
	tokens := extractQueryTokens(query)

	for _, m := range window {
		// WHY guard-continue: the nested limit check inside the
		// token-match check tripped the nested-if linter.
		if !lineMatchesAnyToken(m.Text, tokens) {
			continue
		}
		matches = append(matches, m)
		if limit > 0 && len(matches) >= limit {
			break
		}
	}

	return matches
}

func filterLinesWithWindow(lines []string, opts CatOptions) []LineMatch {
	matches := filterLinesByQuery(lines, opts.Query, 0)
	if opts.TailCount > 0 && len(matches) > opts.TailCount {
		return matches[len(matches)-opts.TailCount:]
	}

	if opts.HeadCount > 0 && len(matches) > opts.HeadCount {
		return matches[:opts.HeadCount]
	}

	if opts.Limit > 0 && len(matches) > opts.Limit {
		return matches[:opts.Limit]
	}

	return matches
}

func filterLinesByQuery(lines []string, query string, limit int) []LineMatch {
	var matches []LineMatch
	tokens := extractQueryTokens(query)

	for idx, line := range lines {
		// WHY guard-continue: same nested-if flattening as above.
		if !lineMatchesAnyToken(line, tokens) {
			continue
		}
		matches = append(matches, LineMatch{Number: idx + 1, Text: line})
		if limit > 0 && len(matches) >= limit {
			break
		}
	}

	return matches
}

// extractPipeTokens splits a pipe-separated query into lowercase tokens,
// skipping blanks. WHY separate helper: the guard inside the loop would
// otherwise sit inside the Contains branch and trip the nested-if linter.
func extractPipeTokens(parts []string) []string {
	var tokens []string
	for _, p := range parts {
		trimmed := strings.ToLower(strings.TrimSpace(p))
		if trimmed == "" {
			continue
		}
		tokens = append(tokens, trimmed)
	}

	return tokens
}

func extractQueryTokens(query string) []string {
	if strings.Contains(query, "|") {
		return extractPipeTokens(strings.Split(query, "|"))
	}

	return []string{strings.ToLower(query)}
}

func lineMatchesAnyToken(line string, tokens []string) bool {
	lLower := strings.ToLower(line)
	for _, token := range tokens {
		if strings.Contains(lLower, token) {
			return true
		}
	}

	return false
}

func extractLinesByWindow(lines []string, opts CatOptions) []LineMatch {
	total := len(lines)
	start, end := resolveWindowBounds(opts, total)
	var matches []LineMatch

	for idx := start; idx <= end && idx <= total; idx++ {
		matches = append(matches, LineMatch{Number: idx, Text: lines[idx-1]})
	}

	return matches
}

func resolveWindowBounds(opts CatOptions, total int) (int, int) {
	if opts.TailCount > 0 {
		return resolveTailBounds(opts.TailCount, total)
	}

	start := opts.StartLine
	if start <= 0 {
		start = 1
	}

	end := total
	if opts.EndLine > 0 && opts.EndLine < total {
		end = opts.EndLine
	}

	if opts.HeadCount > 0 && opts.HeadCount < (end-start+1) {
		end = start + opts.HeadCount - 1
	}

	return start, end
}

func resolveTailBounds(tailCount, total int) (int, int) {
	start := total - tailCount + 1
	if start < 1 {
		start = 1
	}

	return start, total
}
