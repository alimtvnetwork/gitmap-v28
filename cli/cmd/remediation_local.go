package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

// InspectLocalRepoItem inspects a local filesystem path for git status.
func InspectLocalRepoItem(query string) (*RemediationItem, string) {
	cleanQuery := strings.TrimSpace(query)
	if cleanQuery == "" {
		return nil, ""
	}

	absPath, err := filepath.Abs(cleanQuery)
	if err != nil {
		return nil, ""
	}

	gitDir := filepath.Join(absPath, ".git")
	if _, statErr := os.Stat(gitDir); statErr != nil {
		return nil, ""
	}

	return buildLocalRemediationItem(absPath)
}

func buildLocalRemediationItem(absPath string) (*RemediationItem, string) {
	diag := gitutil.InspectDirtyState(absPath)
	repoName := filepath.Base(absPath)
	if diag.IsDirty == false {
		msg := fmt.Sprintf("%s Repository %q is clean. No remediation needed.", constants.ColorGreen+"✓"+constants.ColorReset, repoName)

		return nil, msg
	}

	recipes := gitutil.GenerateRemediationRecipes(absPath, diag)
	if len(recipes) == 0 {
		return nil, ""
	}

	item := RemediationItem{
		RepoPath:      absPath,
		RepoName:      repoName,
		SummaryReason: diag.SummaryReason,
		Recipes:       recipes,
	}

	return &item, ""
}

// FindRemediationSuggestions calculates close repo matches for a mistyped query.
func FindRemediationSuggestions(items []RemediationItem, query string) []string {
	cleanQuery := strings.ToLower(resolveTargetBase(query))
	if cleanQuery == "" {
		cleanQuery = strings.ToLower(strings.TrimSpace(query))
	}

	hasNoItems := cleanQuery == "" || len(items) == 0
	if hasNoItems {
		return nil
	}

	scores := scoreRemediationItems(items, cleanQuery)

	return pickTopSuggestions(scores)
}

type repoScore struct {
	name string
	dist int
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
		if seen[sc.name] == false {
			seen[sc.name] = true
			result = append(result, sc.name)
		}
		if len(result) >= 3 {
			break
		}
	}

	return result
}
