package cmdpull

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/config"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitignoreagm"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// IgnoreRepoIssue holds ignore configuration and index issues for a repository.
type IgnoreRepoIssue struct {
	RepoName          string
	RepoPath          string
	DuplicatePatterns []string
	TrackedPaths      []string
}

// HasIssues reports whether any ignore or tracking issues were detected.
func (i IgnoreRepoIssue) HasIssues() bool {
	return len(i.DuplicatePatterns) > 0 || len(i.TrackedPaths) > 0
}

// IgnoreScanHandle manages asynchronous ignore scanning across repositories.
type IgnoreScanHandle struct {
	done chan []IgnoreRepoIssue
}

func startAsyncIgnoreScan(records []model.ScanRecord) *IgnoreScanHandle {
	return StartThrottledAsyncIgnoreScan(records, resolvePullIgnoreTTL())
}

func resolvePullIgnoreTTL() time.Duration {
	s, err := store.OpenDefault()
	if err != nil {
		return 24 * time.Hour
	}
	defer s.Close()
	return config.GetGitIgnoreTTL(s)
}

func (h *IgnoreScanHandle) Collect() []IgnoreRepoIssue {
	if h == nil {
		return nil
	}
	return <-h.done
}

// DeduplicateIgnoreIssues eliminates duplicate issues by canonical repo path.
func DeduplicateIgnoreIssues(issues []IgnoreRepoIssue) []IgnoreRepoIssue {
	if len(issues) <= 1 {
		return issues
	}
	seen := make(map[string]bool, len(issues))
	unique := make([]IgnoreRepoIssue, 0, len(issues))
	for _, issue := range issues {
		unique = appendUniqueIgnoreIssue(unique, issue, seen)
	}

	return unique
}

func appendUniqueIgnoreIssue(unique []IgnoreRepoIssue, issue IgnoreRepoIssue, seen map[string]bool) []IgnoreRepoIssue {
	key := resolveIgnoreIssueKey(issue)
	if seen[key] {
		return unique
	}
	seen[key] = true

	return append(unique, issue)
}

func resolveIgnoreIssueKey(issue IgnoreRepoIssue) string {
	key := CanonicalRepoPathKey(issue.RepoPath)
	if key != "" {
		return key
	}

	return strings.ToLower(strings.TrimSpace(issue.RepoName))
}

func sortIgnoreIssues(issues []IgnoreRepoIssue) []IgnoreRepoIssue {
	deduped := DeduplicateIgnoreIssues(issues)
	sorted := make([]IgnoreRepoIssue, len(deduped))
	copy(sorted, deduped)
	sort.Slice(sorted, func(i, j int) bool {
		return strings.ToLower(sorted[i].RepoName) < strings.ToLower(sorted[j].RepoName)
	})

	return sorted
}

func inspectRepoForIgnoreIssues(repoDir, repoName string) IgnoreRepoIssue {
	issue := IgnoreRepoIssue{RepoPath: repoDir, RepoName: repoName}
	if !gitignoreagm.IsGitRepository(repoDir) {
		return issue
	}
	issue.DuplicatePatterns = findGitignoreDuplicatePatterns(repoDir)
	issue.TrackedPaths = findTrackedDefaultIgnorePaths(repoDir)
	return issue
}

func findGitignoreDuplicatePatterns(repoDir string) []string {
	data, err := os.ReadFile(filepath.Join(repoDir, ".gitignore"))
	if err != nil {
		return nil
	}
	lines := strings.Split(string(data), "\n")
	seen := make(map[string]bool, len(lines))
	dupSeen := make(map[string]bool)
	var duplicates []string
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if isIgnorablePatternLine(trimmed) {
			continue
		}
		norm := strings.TrimPrefix(trimmed, "/")
		if !seen[norm] {
			seen[norm] = true
			continue
		}
		if !dupSeen[norm] {
			dupSeen[norm] = true
			duplicates = append(duplicates, trimmed)
		}
	}
	return duplicates
}

func isIgnorablePatternLine(trimmed string) bool {
	return trimmed == "" || strings.HasPrefix(trimmed, "#")
}

func findTrackedDefaultIgnorePaths(repoDir string) []string {
	paths := buildDefaultTrackedPathsToCheck(repoDir)
	var tracked []string
	for _, p := range paths {
		if isPathTrackedInGitIndex(repoDir, p) {
			tracked = append(tracked, p)
		}
	}
	return tracked
}

func buildDefaultTrackedPathsToCheck(repoDir string) []string {
	return []string{
		".gitmap/backup/",
		"antigravity-resume_task.json",
		".antigravity_resume_task.json",
		"antigravity_resume_task.json",
		".antigravity-resume_task.json",
	}
}

func isPathTrackedInGitIndex(repoDir, pathspec string) bool {
	return gitutil.ExecGitCheck(5*time.Second, repoDir, "-C", repoDir, "ls-files", "--error-unmatch", pathspec)
}
