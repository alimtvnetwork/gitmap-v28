package cmdinstall

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var desktopToolAliases = map[string]string{
	"github-desktop": constants.ToolGitHubDesktop,
	"gh-desktop":     constants.ToolGitHubDesktop,
	"gd":             constants.ToolGitHubDesktop,
}

// GenerateAutoAliasFn generates an auto alias from a repository name or slug.
var GenerateAutoAliasFn func(name string) string

// PopulateRepoAliasesFn checks and populates aliases for untracked/unaliased repositories.
var PopulateRepoAliasesFn func(db *store.DB) (int, error)

// EnsureTrackedRepoAliases ensures all tracked repositories have aliases.
func EnsureTrackedRepoAliases(db *store.DB) (int, error) {
	if PopulateRepoAliasesFn != nil {
		return PopulateRepoAliasesFn(db)
	}

	return populateFallbackAliases(db)
}

func populateFallbackAliases(db *store.DB) (int, error) {
	repos, err := db.ListUnaliasedRepos()
	if err != nil {
		return 0, err
	}

	return countCreatedFallbackAliases(db, repos), nil
}

func countCreatedFallbackAliases(db *store.DB, repos []store.UnaliasedRepo) int {
	created := 0
	for _, r := range repos {
		if createSingleFallbackAlias(db, r) {
			created++
		}
	}

	return created
}

func createSingleFallbackAlias(db *store.DB, r store.UnaliasedRepo) bool {
	alias := r.RepoName
	if GenerateAutoAliasFn != nil {
		if gen := GenerateAutoAliasFn(r.RepoName); gen != "" {
			alias = gen
		}
	}

	unique := resolveUniqueCandidate(alias, db.AliasExists)
	if unique == "" {
		return false
	}

	_, err := db.CreateAliasWithDetails(unique, r.ID, true, "auto")

	return err == nil
}

func resolveUniqueCandidate(base string, isTaken func(string) bool) string {
	if base == "" {
		return ""
	}
	if !isTaken(base) {
		return base
	}
	for i := 2; i <= 99; i++ {
		candidate := fmt.Sprintf("%s%d", base, i)
		if !isTaken(candidate) {
			return candidate
		}
	}

	return ""
}
