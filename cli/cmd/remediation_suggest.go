package cmd

import (
	"sort"
	"strings"
)

type repoScore struct {
	name string
	dist int
}

// FindRemediationSuggestions calculates close repo matches for a mistyped query.
func FindRemediationSuggestions(items []RemediationItem, query string) []string {
	cleanQuery := strings.ToLower(resolveTargetBase(query))
	if cleanQuery == "" {
		cleanQuery = strings.ToLower(strings.TrimSpace(query))
	}

	hasItems := cleanQuery != "" && len(items) > 0
	if !hasItems {
		return nil
	}

	scores := scoreRemediationItems(items, cleanQuery)

	return pickTopSuggestions(scores)
}

func scoreRemediationItems(items []RemediationItem, query string) []repoScore {
	var scores []repoScore
	for _, it := range items {
		d := levenshtein(query, strings.ToLower(it.RepoName))
		isSub := strings.Contains(strings.ToLower(it.RepoName), query)
		if isSub {
			d = 1
		}

		if d <= 3 {
			scores = append(scores, repoScore{name: it.RepoName, dist: d})
		}
	}

	sort.SliceStable(scores, func(i, j int) bool {
		return scores[i].dist < scores[j].dist
	})

	return scores
}

func pickTopSuggestions(scores []repoScore) []string {
	var result []string
	seen := make(map[string]bool)
	for _, sc := range scores {
		if !seen[sc.name] {
			seen[sc.name] = true
			result = append(result, sc.name)
		}

		if len(result) >= 3 {
			break
		}
	}

	return result
}
