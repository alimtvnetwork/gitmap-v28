package cmdfix

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdremediation"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// IsValidGitRepository verifies that the target path contains a verified .git directory or worktree.
func IsValidGitRepository(targetPath string) bool {
	return gitutil.IsValidGitRepository(targetPath)
}

// IsWorkspaceContainerDirectory checks whether a path represents a parent workspace folder rather than a leaf repository.
func IsWorkspaceContainerDirectory(targetPath string) bool {
	return gitutil.IsWorkspaceContainerDirectory(targetPath)
}

// GuardRepoExecution verifies repo validity; returns an AppError if targetPath is invalid or a container directory.
func GuardRepoExecution(targetPath string, opName string) *apperror.AppError {
	if IsWorkspaceContainerDirectory(targetPath) {
		return apperror.NewValidationError("cannot execute " + opName + ": directory is a workspace container, not a Git repository")
	}

	if !isGitDirectory(targetPath) {
		return apperror.NewValidationError("cannot execute " + opName + ": directory is not a Git repository")
	}

	return nil
}

func resolveRepoPathForFix(repoQuery string) string {
	cleanQuery := strings.TrimSpace(repoQuery)
	if cleanQuery == "" || IsWorkspaceContainerDirectory(cleanQuery) {
		return ""
	}

	if info, err := os.Stat(cleanQuery); err == nil && info.IsDir() {
		return cleanQuery
	}

	return queryDBRepoDirFix(cleanQuery)
}

func queryDBRepoDirFix(clean string) string {
	if special := querySpecialRepoDirFix(clean); special != "" {
		return special
	}

	return queryStandardRepoDirFix(clean)
}

func querySpecialRepoDirFix(clean string) string {
	rec := fetchSpecialRepoRecord(clean)
	if rec == nil || rec.LocalPath == "" || IsWorkspaceContainerDirectory(rec.LocalPath) {
		return ""
	}

	info, statErr := os.Stat(rec.LocalPath)
	if statErr == nil && info.IsDir() {
		return rec.LocalPath
	}

	return ""
}

func fetchSpecialRepoRecord(clean string) *store.SpecialRepoRecord {
	db, err := store.OpenSpecialReposSplitDB()
	if err != nil {
		return nil
	}

	defer db.Close()
	rec, _ := db.GetSpecialRepo(clean)

	return rec
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

	return findFirstValidDir(recs)
}

func findFirstValidDir(recs []model.ScanRecord) string {
	for _, r := range recs {
		if r.AbsolutePath == "" || IsWorkspaceContainerDirectory(r.AbsolutePath) {
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
	cloneSteps := resolveCloneOrInitSteps(diag, cleanPath)
	linkSteps := buildLinkRemoteSteps(cleanPath, diag.RemoteURL)

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

func resolveCloneOrInitSteps(diag cloner.NonRepoDiagnosis, cleanPath string) []gitutil.RemediationStep {
	entries, readErr := os.ReadDir(diag.Path)
	if readErr == nil && len(entries) == 0 {
		return []gitutil.RemediationStep{
			{Name: "git", Args: []string{"clone", diag.RemoteURL, cleanPath}},
		}
	}

	return buildLinkRemoteSteps(cleanPath, diag.RemoteURL)
}

func buildLinkRemoteSteps(cleanPath, remoteURL string) []gitutil.RemediationStep {
	return []gitutil.RemediationStep{
		{Name: "git", Args: []string{"-C", cleanPath, "init"}},
		{Name: "git", Args: []string{"-C", cleanPath, "remote", "add", "origin", remoteURL}},
		{Name: "git", Args: []string{"-C", cleanPath, "fetch", "origin"}},
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
	all = appendValidRemediationItems(all, existing, seen)
	liveItems := scanTrackedReposForIssues()
	all = appendValidRemediationItems(all, liveItems, seen)
	all = appendSpecialRepoItems(all, seen)

	return all
}

func appendValidRemediationItems(all, items []cmdremediation.RemediationItem, seen map[string]bool) []cmdremediation.RemediationItem {
	for _, item := range items {
		cleanPath := filepath.Clean(item.RepoPath)
		if !isGitDirectory(cleanPath) || IsWorkspaceContainerDirectory(cleanPath) {
			continue
		}

		if !seen[cleanPath] {
			seen[cleanPath] = true
			all = append(all, item)
		}
	}

	return all
}

func appendSpecialRepoItems(all []cmdremediation.RemediationItem, seen map[string]bool) []cmdremediation.RemediationItem {
	for _, key := range []string{"repo-cache", "repo-secrets"} {
		item, path := collectSpecialRepoItem(key, seen)
		if item != nil && isGitDirectory(path) && !IsWorkspaceContainerDirectory(path) {
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
