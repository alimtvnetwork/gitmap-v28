package cmd

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func suggestCDRepos(name string) []string {
	norm := strings.ToLower(strings.TrimSpace(name))
	if norm == "" {
		return nil
	}

	candidates := collectCDRepoCandidates()
	if len(candidates) == 0 {
		return nil
	}

	scores := rankCandidateRepos(norm, candidates)

	return selectBestRepoSuggestions(scores)
}

func collectCDRepoCandidates() []string {
	db, err := store.OpenDefault()
	if err != nil {
		return nil
	}

	defer db.Close()

	candidates := appendRepoCandidates(db)

	return appendWorkDirCandidates(db, candidates)
}

func appendRepoCandidates(db *store.DB) []string {
	var candidates []string
	repos, err := db.ListRepos()
	if err != nil {
		return nil
	}

	for _, r := range repos {
		candidates = appendCandidateFromRecord(candidates, r)
	}

	return candidates
}

func appendCandidateFromRecord(candidates []string, r model.ScanRecord) []string {
	base := filepath.Base(r.AbsolutePath)
	candidates = appendUniqueCandidate(candidates, base)
	if r.Slug == "" {
		return candidates
	}

	candidates = appendUniqueCandidate(candidates, r.Slug)
	parts := strings.Split(r.Slug, "/")
	if len(parts) > 1 {
		candidates = appendUniqueCandidate(candidates, parts[len(parts)-1])
	}

	return candidates
}

func appendWorkDirCandidates(db *store.DB, candidates []string) []string {
	dirs, err := db.ListWorkDirs()
	if err != nil {
		return candidates
	}

	for _, d := range dirs {
		if d.Label != "" {
			candidates = appendUniqueCandidate(candidates, d.Label)
		}
		candidates = appendUniqueCandidate(candidates, filepath.Base(d.AbsolutePath))
	}

	return candidates
}

func appendUniqueCandidate(list []string, item string) []string {
	clean := strings.TrimSpace(item)
	if clean == "" {
		return list
	}

	for _, existing := range list {
		if strings.EqualFold(existing, clean) {
			return list
		}
	}

	return append(list, clean)
}

func rankCandidateRepos(input string, candidates []string) []repoScore {
	var scores []repoScore
	for _, c := range candidates {
		score, hasMatch := evaluateCandidateRepo(input, c)
		if hasMatch {
			scores = append(scores, score)
		}
	}

	sort.SliceStable(scores, func(i, j int) bool {
		return scores[i].dist < scores[j].dist
	})

	return scores
}

func evaluateCandidateRepo(input, candidate string) (repoScore, bool) {
	cLower := strings.ToLower(candidate)
	isSubstring := strings.Contains(cLower, input) || strings.Contains(input, cLower)
	if isSubstring {
		return repoScore{name: candidate, dist: 1}, true
	}

	d := levenshtein(input, cLower)
	if d <= 3 {
		return repoScore{name: candidate, dist: d}, true
	}

	return repoScore{}, false
}

func selectBestRepoSuggestions(scores []repoScore) []string {
	var result []string
	seen := make(map[string]bool)
	for _, sc := range scores {
		lower := strings.ToLower(sc.name)
		if !seen[lower] {
			seen[lower] = true
			result = append(result, sc.name)
		}
		if len(result) >= 3 {
			break
		}
	}

	return result
}

func formatCDNotFoundMessage(name string, suggestions []string) string {
	msg := fmt.Sprintf("no repo found matching '%s'", name)
	if len(suggestions) > 0 {
		msg += fmt.Sprintf("\n  Did you mean: %s?", strings.Join(suggestions, ", "))
	}

	return msg
}
