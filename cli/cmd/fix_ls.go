// Package cmd — fix_ls.go discovers and lists all git repositories with issues.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func isFixLsRequest(args []string) bool {
	if len(args) == 0 {
		return false
	}
	low := strings.ToLower(args[0])

	return low == "ls" || low == "list" || low == "-l" || low == "--list"
}

func runFixLs(args []string) error {
	checkHelp("fix", args)
	records := discoverReposForFixLs()
	items := collectReposWithIssues(records)
	if len(items) == 0 {
		printCleanFixLsOutput(len(records))

		return nil
	}

	_ = SaveRemediationState(items)
	renderFixLsTable(items)

	return nil
}

func discoverReposForFixLs() []model.ScanRecord {
	if HasAlias() {
		return []model.ScanRecord{{
			RepoName:     GetAliasSlug(),
			Slug:         GetAliasSlug(),
			AbsolutePath: GetAliasPath(),
		}}
	}
	records := loadAllRecordsDB()
	if len(records) > 0 {
		return mergeCachedRemediationItems(records)
	}

	return discoverFallbackRecords()
}

func discoverFallbackRecords() []model.ScanRecord {
	records := loadRecordsJSONFallback()
	if len(records) > 0 {
		return mergeCachedRemediationItems(records)
	}
	cwdRecords := discoverCwdRecord()

	return mergeCachedRemediationItems(cwdRecords)
}

func discoverCwdRecord() []model.ScanRecord {
	root, err := gitutil.RepoRoot(".")
	if err != nil || root == "" {
		return nil
	}

	return []model.ScanRecord{{
		RepoName:     filepath.Base(root),
		Slug:         filepath.Base(root),
		AbsolutePath: root,
	}}
}

func mergeCachedRemediationItems(records []model.ScanRecord) []model.ScanRecord {
	cached := LoadRemediationState()
	if len(cached) == 0 {
		return records
	}
	known := make(map[string]bool)
	for _, r := range records {
		known[filepath.Clean(r.AbsolutePath)] = true
	}

	return appendUnknownCached(records, cached, known)
}

func appendUnknownCached(
	records []model.ScanRecord,
	cached []RemediationItem,
	known map[string]bool,
) []model.ScanRecord {
	for _, c := range cached {
		cleanPath := filepath.Clean(c.RepoPath)
		if !known[cleanPath] {
			records = append(records, model.ScanRecord{
				RepoName:     c.RepoName,
				Slug:         c.RepoName,
				AbsolutePath: c.RepoPath,
			})
			known[cleanPath] = true
		}
	}

	return records
}

func collectReposWithIssues(records []model.ScanRecord) []RemediationItem {
	var items []RemediationItem
	for _, rec := range records {
		item := inspectRepoIssue(rec)
		if item != nil {
			items = append(items, *item)
		}
	}

	return items
}

func inspectRepoIssue(rec model.ScanRecord) *RemediationItem {
	if !isGitDirectory(rec.AbsolutePath) {
		return nil
	}
	diag := gitutil.InspectDirtyState(rec.AbsolutePath)
	rs := gitutil.Status(rec.AbsolutePath)
	isConflict := hasMergeConflict(rec.AbsolutePath)
	isLocked := hasGitLockFile(rec.AbsolutePath)
	if !isRepoIssue(diag, rs, isConflict, isLocked) {
		return nil
	}

	return buildRemediationIssueItem(rec, diag, rs, isConflict, isLocked)
}

func isGitDirectory(path string) bool {
	if len(path) == 0 {
		return false
	}
	gitDir := filepath.Join(path, ".git")
	info, err := os.Stat(gitDir)

	return err == nil && info.IsDir()
}

func hasMergeConflict(repoPath string) bool {
	gitDir := filepath.Join(repoPath, ".git")
	markers := []string{"MERGE_HEAD", "rebase-apply", "rebase-merge", "CHERRY_PICK_HEAD"}
	for _, m := range markers {
		if _, err := os.Stat(filepath.Join(gitDir, m)); err == nil {
			return true
		}
	}

	return false
}

func hasGitLockFile(repoPath string) bool {
	lockFile := filepath.Join(repoPath, ".git", "index.lock")
	_, err := os.Stat(lockFile)

	return err == nil
}

func isRepoIssue(diag gitutil.DirtyDiagnosis, rs gitutil.RepoStatus, isConflict, isLocked bool) bool {
	if diag.IsDirty || isConflict || isLocked {
		return true
	}
	if rs.Behind > 0 || rs.Ahead > 0 || rs.StashCount > 0 {
		return true
	}

	return false
}

func buildRemediationIssueItem(
	rec model.ScanRecord,
	diag gitutil.DirtyDiagnosis,
	rs gitutil.RepoStatus,
	isConflict bool,
	isLocked bool,
) *RemediationItem {
	return &RemediationItem{
		RepoPath:      rec.AbsolutePath,
		RepoName:      resolveRepoDisplayName(rec),
		SummaryReason: buildIssueSummary(diag, rs, isConflict, isLocked),
		Recipes:       resolveIssueRecipes(rec.AbsolutePath, diag),
		Files:         diag.AllFiles,
	}
}

func resolveRepoDisplayName(rec model.ScanRecord) string {
	if rec.RepoName != "" {
		return rec.RepoName
	}
	if rec.Slug != "" {
		return rec.Slug
	}

	return filepath.Base(rec.AbsolutePath)
}

func resolveIssueRecipes(repoPath string, diag gitutil.DirtyDiagnosis) []gitutil.RemediationRecipe {
	if diag.IsDirty {
		return gitutil.GenerateRemediationRecipes(repoPath, diag)
	}

	return []gitutil.RemediationRecipe{
		gitutil.GenerateStashRecipe(repoPath),
		gitutil.GenerateCommitRecipe(repoPath),
		gitutil.GenerateDiscardRecipe(repoPath),
	}
}

func buildIssueSummary(diag gitutil.DirtyDiagnosis, rs gitutil.RepoStatus, isConflict, isLocked bool) string {
	var parts []string
	if isConflict {
		parts = append(parts, "merge conflict")
	}
	if isLocked {
		parts = append(parts, "lock (index.lock)")
	}
	parts = appendDirtySummary(parts, diag)
	parts = appendSyncSummary(parts, rs)
	if len(parts) == 0 {
		return "issue detected"
	}

	return strings.Join(parts, ", ")
}

func appendDirtySummary(parts []string, diag gitutil.DirtyDiagnosis) []string {
	if !diag.IsDirty {
		return parts
	}
	if diag.SummaryReason != "" {
		return append(parts, "dirty: "+diag.SummaryReason)
	}

	return append(parts, "uncommitted changes")
}

func appendSyncSummary(parts []string, rs gitutil.RepoStatus) []string {
	switch {
	case rs.Behind > 0 && rs.Ahead > 0:
		parts = append(parts, fmt.Sprintf("diverged (+%d/-%d)", rs.Ahead, rs.Behind))
	case rs.Behind > 0:
		parts = append(parts, fmt.Sprintf("behind (%d)", rs.Behind))
	case rs.Ahead > 0:
		parts = append(parts, fmt.Sprintf("ahead (%d)", rs.Ahead))
	}
	if rs.StashCount > 0 {
		parts = append(parts, fmt.Sprintf("%d stash(es)", rs.StashCount))
	}

	return parts
}

func scanTrackedReposForIssues() []RemediationItem {
	records := discoverReposForFixLs()

	return collectReposWithIssues(records)
}
