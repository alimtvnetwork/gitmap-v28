package cmdfix

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdremediation"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func resolveRepoPathForFix(repoQuery string) string {
	cleanQuery := strings.TrimSpace(repoQuery)
	if cleanQuery == "" {
		return ""
	}
	if info, err := os.Stat(cleanQuery); err == nil && info.IsDir() {
		return cleanQuery
	}
	if special := querySpecialRepoDirFix(cleanQuery); special != "" {
		return special
	}
	if standard := queryStandardRepoDirFix(cleanQuery); standard != "" {
		return standard
	}
	return cleanQuery
}

func querySpecialRepoDirFix(clean string) string {
	db, err := store.OpenSpecialReposSplitDB()
	if err != nil {
		return ""
	}
	defer db.Close()
	rec, err := db.GetSpecialRepo(clean)
	if err != nil || rec == nil || rec.LocalPath == "" {
		return ""
	}
	info, statErr := os.Stat(rec.LocalPath)
	if statErr == nil && info.IsDir() {
		return rec.LocalPath
	}
	return ""
}

func queryStandardRepoDirFix(clean string) string {
	db, err := store.OpenDefault()
	if err != nil {
		return ""
	}
	defer db.Close()
	recs, err := db.FindBySlug(clean)
	if err != nil || len(recs) == 0 {
		return ""
	}
	for _, r := range recs {
		if r.AbsolutePath == "" {
			continue
		}
		info, statErr := os.Stat(r.AbsolutePath)
		if statErr == nil && info.IsDir() {
			return r.AbsolutePath
		}
	}
	return ""
}

func resolveNonRepoRemediationItem(repoQuery, action string) (*cmdremediation.RemediationItem, error) {
	targetPath := resolveRepoPathForFix(repoQuery)
	if targetPath == "" {
		return nil, nil
	}

	diag := cloner.ClassifyNonRepoFolder(targetPath, repoQuery)
	if !diag.IsNonRepoFolder {
		return nil, nil
	}

	return buildNonRepoRemediationItem(diag, action)
}

func buildNonRepoRemediationItem(diag cloner.NonRepoDiagnosis, action string) (*cmdremediation.RemediationItem, error) {
	cleanPath := filepath.ToSlash(filepath.Clean(diag.Path))
	var recipes []gitutil.RemediationRecipe

	if diag.HasRemote {
		recipes = buildRemoteRecipes(diag, cleanPath)
	} else {
		recipes = buildLocalOnlyRecipes(diag, cleanPath)
	}

	return &cmdremediation.RemediationItem{
		RepoName:      diag.RepoName,
		RepoPath:      diag.Path,
		SummaryReason: diag.Reason,
		Recipes:       recipes,
	}, nil
}

func buildRemoteRecipes(diag cloner.NonRepoDiagnosis, cleanPath string) []gitutil.RemediationRecipe {
	entries, readErr := os.ReadDir(diag.Path)
	isEmpty := readErr == nil && len(entries) == 0

	var cloneSteps []gitutil.RemediationStep
	if isEmpty {
		cloneSteps = []gitutil.RemediationStep{
			{Name: "git", Args: []string{"clone", diag.RemoteURL, cleanPath}},
		}
	} else {
		cloneSteps = []gitutil.RemediationStep{
			{Name: "git", Args: []string{"-C", cleanPath, "init"}},
			{Name: "git", Args: []string{"-C", cleanPath, "remote", "add", "origin", diag.RemoteURL}},
			{Name: "git", Args: []string{"-C", cleanPath, "fetch", "origin"}},
		}
	}

	linkSteps := []gitutil.RemediationStep{
		{Name: "git", Args: []string{"-C", cleanPath, "init"}},
		{Name: "git", Args: []string{"-C", cleanPath, "remote", "add", "origin", diag.RemoteURL}},
		{Name: "git", Args: []string{"-C", cleanPath, "fetch", "origin"}},
	}

	return []gitutil.RemediationRecipe{
		{
			Title:       "Clone from Remote",
			Command:     diag.Option1,
			Description: fmt.Sprintf("Clone %s into %s", diag.RemoteURL, cleanPath),
			Steps:       cloneSteps,
		},
		{
			Title:       "Init & Link Remote",
			Command:     diag.Option2,
			Description: fmt.Sprintf("Initialize local git repo in %s and link origin to %s", cleanPath, diag.RemoteURL),
			Steps:       linkSteps,
		},
	}
}

func buildLocalOnlyRecipes(diag cloner.NonRepoDiagnosis, cleanPath string) []gitutil.RemediationRecipe {
	initSteps := []gitutil.RemediationStep{
		{Name: "git", Args: []string{"-C", cleanPath, "init"}},
	}
	removeSteps := []gitutil.RemediationStep{
		{Name: "gitmap", Args: []string{"rm", diag.RepoName, "--db-only"}},
	}
	return []gitutil.RemediationRecipe{
		{
			Title:       "Initialize Local Repo",
			Command:     diag.Option1,
			Description: fmt.Sprintf("Initialize local git repository in %s", cleanPath),
			Steps:       initSteps,
		},
		{
			Title:       "Remove from Registry",
			Command:     diag.Option2,
			Description: fmt.Sprintf("Remove %s from database registry", diag.RepoName),
			Steps:       removeSteps,
		},
	}
}

func collectAllRemediationItems(existing []cmdremediation.RemediationItem) []cmdremediation.RemediationItem {
	seen := make(map[string]bool)
	var all []cmdremediation.RemediationItem
	for _, item := range existing {
		cleanPath := filepath.Clean(item.RepoPath)
		if !seen[cleanPath] {
			seen[cleanPath] = true
			all = append(all, item)
		}
	}

	liveItems := scanTrackedReposForIssues()
	for _, item := range liveItems {
		cleanPath := filepath.Clean(item.RepoPath)
		if !seen[cleanPath] {
			seen[cleanPath] = true
			all = append(all, item)
		}
	}

	for _, key := range []string{"repo-cache", "repo-secrets"} {
		item, path := collectSpecialRepoItem(key, seen)
		if item != nil {
			seen[path] = true
			all = append(all, *item)
		}
	}

	return all
}

func collectSpecialRepoItem(key string, seen map[string]bool) (*cmdremediation.RemediationItem, string) {
	path := resolveRepoPathForFix(key)
	if path == "" {
		return nil, ""
	}
	cleanPath := filepath.Clean(path)
	if seen[cleanPath] {
		return nil, ""
	}
	diag := cloner.ClassifyNonRepoFolder(path, key)
	if !diag.IsNonRepoFolder {
		return nil, ""
	}
	item, err := buildNonRepoRemediationItem(diag, "")
	if err != nil || item == nil {
		return nil, ""
	}
	return item, cleanPath
}
