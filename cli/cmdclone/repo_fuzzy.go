package cmdclone

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func findClosestRepoSuggestions(db *store.DB, target string) []string {
	repos, err := db.ListRepos()
	if err != nil {
		return nil
	}
	clean := strings.ToLower(strings.TrimSpace(target))
	return collectClosestMatches(repos, clean)
}

func collectClosestMatches(repos []model.ScanRecord, clean string) []string {
	var matches []string
	seen := make(map[string]bool)
	for _, r := range repos {
		name := strings.ToLower(r.RepoName)
		slug := strings.ToLower(r.Slug)
		base := strings.Split(name, "-v")[0]
		if !isSuggestionMatch(clean, name, slug, base) || seen[r.Slug] {
			continue
		}
		seen[r.Slug] = true
		matches = append(matches, r.Slug)
		if len(matches) >= 5 {
			break
		}
	}
	return matches
}

func isSuggestionMatch(clean, name, slug, base string) bool {
	if strings.Contains(name, clean) || (len(base) > 2 && strings.Contains(clean, base)) {
		return true
	}
	if levenshtein(clean, name) <= 4 || levenshtein(clean, slug) <= 4 || levenshtein(clean, base) <= 3 {
		return true
	}
	return false
}

func levenshtein(a, b string) int {
	ar, br := []rune(a), []rune(b)
	la, lb := len(ar), len(br)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}
	dp := initLevenshteinDP(lb)
	fillLevenshteinDP(dp, ar, br, la, lb)
	return dp[lb]
}

func initLevenshteinDP(lb int) []int {
	dp := make([]int, lb+1)
	for j := 0; j <= lb; j++ {
		dp[j] = j
	}
	return dp
}

func fillLevenshteinDP(dp []int, ar, br []rune, la, lb int) {
	for i := 1; i <= la; i++ {
		prev := dp[0]
		dp[0] = i
		updateLevenshteinRow(dp, ar[i-1], br, lb, &prev)
	}
}

func updateLevenshteinRow(dp []int, charA rune, br []rune, lb int, prev *int) {
	for j := 1; j <= lb; j++ {
		temp := dp[j]
		cost := 0
		if charA != br[j-1] {
			cost = 1
		}
		dp[j] = min3(dp[j]+1, dp[j-1]+1, *prev+cost)
		*prev = temp
	}
}

func min3(a, b, c int) int {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	return m
}
