package cmdremediation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
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
	if _, statErr := os.Stat(gitDir); statErr == nil {
		return buildLocalRemediationItem(absPath)
	}

	resolved := resolveTargetGitRepoPath(cleanQuery)
	if resolved != "" {
		return buildLocalRemediationItem(resolved)
	}

	return nil, ""
}

func resolveTargetGitRepoPath(query string) string {
	cwd, err := os.Getwd()
	if err != nil {
		return resolveRepoFromStoreDB(query)
	}

	if isCwdRepoMatch(cwd, query) {
		return cwd
	}

	if childPath := resolveChildRepoPath(cwd, query); childPath != "" {
		return childPath
	}

	return resolveRepoFromStoreDB(query)
}

func isCwdRepoMatch(cwd, query string) bool {
	if !strings.EqualFold(filepath.Base(cwd), query) {
		return false
	}
	_, statErr := os.Stat(filepath.Join(cwd, ".git"))
	return statErr == nil
}

func resolveChildRepoPath(cwd, query string) string {
	childPath := filepath.Join(cwd, query)
	if _, statErr := os.Stat(filepath.Join(childPath, ".git")); statErr == nil {
		return childPath
	}
	return ""
}

func resolveRepoFromStoreDB(query string) string {
	db, errDB := store.OpenDefault()
	if errDB != nil {
		return ""
	}
	defer db.Close()

	repos, errList := db.ListRepos()
	if errList != nil {
		return ""
	}

	for _, r := range repos {
		if matchesRepoQuery(r, query) {
			return r.AbsolutePath
		}
	}
	return ""
}

func matchesRepoQuery(r model.ScanRecord, query string) bool {
	if !strings.EqualFold(r.RepoName, query) && !strings.EqualFold(r.Slug, query) {
		return false
	}
	_, statErr := os.Stat(filepath.Join(r.AbsolutePath, ".git"))
	return statErr == nil
}

func buildLocalRemediationItem(absPath string) (*RemediationItem, string) {
	diag := gitutil.InspectDirtyState(absPath)
	repoName := filepath.Base(absPath)
	if !diag.IsDirty {
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
		Files:         diag.AllFiles,
	}

	return &item, ""
}
