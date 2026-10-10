// Package cloner — pulldiag.go provides pull diagnosis and file-lock remediation.
package cloner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func clearReadOnlyAttrs(repoDir, output string) bool {
	if runtime.GOOS != constants.OSWindows {
		return false
	}

	paths := extractUnlinkPaths(output)
	if len(paths) == 0 {
		return false
	}

	cleared := false
	for _, relativePath := range paths {
		fullPath := filepath.Join(repoDir, filepath.FromSlash(relativePath))
		if clearReadOnly(fullPath) {
			cleared = true
		}
	}

	return cleared
}

func clearReadOnly(path string) bool {
	cmd := exec.Command("attrib", "-R", path)
	if err := cmd.Run(); err == nil {
		return true
	}

	return os.Chmod(path, 0o666) == nil
}

func buildPullDiagnosis(repoDir, output string) string {
	hints := collectDiagnosisHints(repoDir, output)
	if len(hints) == 0 {
		hints = append(hints, "non-unlink git pull failure (check auth/merge or run pull manually for full output)")
	}

	return strings.Join(hints, "; ")
}

func collectDiagnosisHints(repoDir, output string) []string {
	hints := collectFailureHints(output)
	hints = appendPathHints(hints, repoDir, output)
	hints = appendNonGitRepoHints(hints, repoDir, output)

	return hints
}

func collectFailureHints(output string) []string {
	hints := make([]string, 0, 4)
	hints = appendLockFailureHints(hints, output)
	hints = appendConflictFailureHints(hints, output)

	return hints
}

func appendLockFailureHints(hints []string, output string) []string {
	if hasUnlinkFailure(output) {
		hints = append(hints, "file lock/read-only attribute blocked replacing old files")
	}

	return hints
}

func appendConflictFailureHints(hints []string, output string) []string {
	if hasUnmergedFailure(output) {
		hints = append(hints, "unresolved merge conflict detected; run 'gitmap fix-git' or 'git merge --abort'")
	}
	if hasUntrackedOverwriteFailure(output) {
		hints = append(hints, "untracked files conflict with incoming commits; run 'gitmap fix-git' to backup & pull")
	}
	if hasMultipleBranchesFailure(output) {
		hints = append(hints, "multiple upstream branches or concurrent fetch conflict in FETCH_HEAD; run pull targeting specific branch")
	}

	return hints
}

func appendPathHints(hints []string, repoDir, output string) []string {
	if hasPathLengthRisk(repoDir, output) {
		hints = append(hints, "Windows path length risk detected; use a shorter base path like C:\\src")
	}
	if strings.Contains(strings.ToLower(repoDir), "onedrive") {
		hints = append(hints, "repo is under a synced folder (OneDrive), which often locks files")
	}

	return hints
}

func hasMultipleBranchesFailure(output string) bool {
	lower := strings.ToLower(output)

	return strings.Contains(lower, "cannot fast-forward to multiple branches")
}

func hasUnmergedFailure(output string) bool {
	lower := strings.ToLower(output)

	return strings.Contains(lower, "you have unmerged files") || strings.Contains(lower, "unresolved conflict")
}

func hasUntrackedOverwriteFailure(output string) bool {
	lower := strings.ToLower(output)

	return strings.Contains(lower, "untracked working tree files would be overwritten")
}

func hasUnlinkFailure(output string) bool {
	lower := strings.ToLower(output)

	return strings.Contains(lower, "unable to unlink old") || strings.Contains(lower, "unlink of file")
}

func hasPathLengthRisk(repoDir, output string) bool {
	if runtime.GOOS != constants.OSWindows {
		return false
	}

	for _, relativePath := range extractUnlinkPaths(output) {
		fullPath := filepath.Join(repoDir, filepath.FromSlash(relativePath))
		if len(fullPath) >= constants.WindowsPathWarnThreshold {
			return true
		}
	}

	return false
}

func extractUnlinkPaths(output string) []string {
	matches := collectRegexMatches(output)

	return deduplicateStrings(matches)
}

func collectRegexMatches(output string) []string {
	matches := make([]string, 0, 2)
	for _, m := range unlinkOldRegex.FindAllStringSubmatch(output, -1) {
		if len(m) > 1 {
			matches = append(matches, m[1])
		}
	}

	for _, m := range unlinkPromptRegex.FindAllStringSubmatch(output, -1) {
		if len(m) > 1 {
			matches = append(matches, m[1])
		}
	}

	return matches
}

func deduplicateStrings(items []string) []string {
	seen := map[string]struct{}{}
	unique := make([]string, 0, len(items))
	for _, item := range items {
		_, isSeen := seen[item]
		if isSeen {
			continue
		}

		seen[item] = struct{}{}
		unique = append(unique, item)
	}

	return unique
}

func trimOutput(output string) string {
	trimmed := strings.TrimSpace(output)
	if len(trimmed) <= 1200 {
		return trimmed
	}

	return trimmed[:1200] + "..."
}

// NonRepoDiagnosis represents the diagnosis and remediation options for a non-git directory.
type NonRepoDiagnosis struct {
	IsNonRepoFolder bool   `json:"is_non_repo_folder"`
	IsSpecialRepo   bool   `json:"is_special_repo"`
	RepoName        string `json:"repo_name"`
	Path            string `json:"path"`
	RemoteURL       string `json:"remote_url,omitempty"`
	HasRemote       bool   `json:"has_remote"`
	Reason          string `json:"reason"`
	Option1         string `json:"option_1"`
	Option1Title    string `json:"option_1_title,omitempty"`
	Option1Cmd      string `json:"option_1_cmd,omitempty"`
	Option2         string `json:"option_2"`
	Option2Title    string `json:"option_2_title,omitempty"`
	Option2Cmd      string `json:"option_2_cmd,omitempty"`
}

// ClassifyNonRepoFolder inspects a path and classifies non-repo folder diagnosis and remediation.
func ClassifyNonRepoFolder(path string, repoName string) NonRepoDiagnosis {
	if path == "" {
		return NonRepoDiagnosis{}
	}
	stat, err := os.Stat(path)
	if err != nil || !stat.IsDir() {
		return NonRepoDiagnosis{}
	}
	if IsGitRepo(path) {
		return NonRepoDiagnosis{}
	}

	cleanPath := filepath.ToSlash(filepath.Clean(path))
	name := repoName
	if name == "" {
		name = filepath.Base(cleanPath)
	}

	if isSpecialInfrastructureRepo(name) {
		canonicalName, _ := canonicalizeSpecialRepoKey(name)
		remoteURL := resolveSpecialRepoRemoteURL(canonicalName)
		hasRemote := remoteURL != ""

		reason := fmt.Sprintf("directory exists but is not a Git repository (known infrastructure repository: %s)", name)
		opt1Cmd := fmt.Sprintf("gitmap clone %s", canonicalName)
		opt2Cmd := fmt.Sprintf("git -C \"%s\" init", cleanPath)

		return NonRepoDiagnosis{
			IsNonRepoFolder: true,
			IsSpecialRepo:   true,
			RepoName:        name,
			Path:            cleanPath,
			RemoteURL:       remoteURL,
			HasRemote:       hasRemote,
			Reason:          reason,
			Option1:         opt1Cmd,
			Option1Title:    "Clone from Remote",
			Option1Cmd:      opt1Cmd,
			Option2:         opt2Cmd,
			Option2Title:    "Initialize Local Repo",
			Option2Cmd:      opt2Cmd,
		}
	}

	// Standard repository
	remoteURL, hasRemote := resolveStandardRepoRemoteURL(name, cleanPath)
	reason := "directory exists but is not a Git repository (missing .git)"
	if hasRemote {
		opt1Cmd := fmt.Sprintf("gitmap clone %s", name)
		opt2Cmd := fmt.Sprintf("cd \"%s\" && git init && git remote add origin %s && git fetch", cleanPath, remoteURL)

		return NonRepoDiagnosis{
			IsNonRepoFolder: true,
			IsSpecialRepo:   false,
			RepoName:        name,
			Path:            cleanPath,
			RemoteURL:       remoteURL,
			HasRemote:       true,
			Reason:          reason,
			Option1:         opt1Cmd,
			Option1Title:    "Clone from Remote",
			Option1Cmd:      opt1Cmd,
			Option2:         opt2Cmd,
			Option2Title:    "Init & Link Remote",
			Option2Cmd:      opt2Cmd,
		}
	}

	opt1Cmd := fmt.Sprintf("git -C \"%s\" init", cleanPath)
	opt2Cmd := fmt.Sprintf("gitmap rm \"%s\" --db-only", name)

	return NonRepoDiagnosis{
		IsNonRepoFolder: true,
		IsSpecialRepo:   false,
		RepoName:        name,
		Path:            cleanPath,
		RemoteURL:       "",
		HasRemote:       false,
		Reason:          reason,
		Option1:         opt1Cmd,
		Option1Title:    "Initialize Local Repo",
		Option1Cmd:      opt1Cmd,
		Option2:         opt2Cmd,
		Option2Title:    "Remove from Registry",
		Option2Cmd:      opt2Cmd,
	}
}

// IsSpecialInfrastructureRepo reports whether name is a companion infrastructure repository.
func IsSpecialInfrastructureRepo(name string) bool {
	return isSpecialInfrastructureRepo(name)
}

func isSpecialInfrastructureRepo(name string) bool {
	cleaned := strings.ToLower(strings.TrimSpace(name))
	switch cleaned {
	case "repo-cache", "repo-secrets", "rc", "rs", "cache", "storage",
		"repo-storage", "secrets", "vault":
		return true
	default:
		return false
	}
}

// CanonicalizeSpecialRepoKey returns primary key and short key for special infrastructure repos.
func CanonicalizeSpecialRepoKey(name string) (string, string) {
	return canonicalizeSpecialRepoKey(name)
}

func canonicalizeSpecialRepoKey(name string) (string, string) {
	cleaned := strings.ToLower(strings.TrimSpace(name))
	switch cleaned {
	case "rc", "repo-cache", "cache", "storage", "repo-storage":
		return "repo-cache", "rc"
	default:
		return "repo-secrets", "rs"
	}
}

func probeRemoteSpecialRepoURL(repoName string) string {
	cmd := exec.Command("gh", "repo", "view", repoName, "--json", "url", "-q", ".url")
	out, err := cmd.Output()
	if err == nil {
		u := strings.TrimSpace(string(out))
		if len(u) > 0 {
			return u
		}
	}
	return ""
}

func resolveSpecialRepoRemoteURL(repoName string) string {
	if db, err := store.OpenSpecialReposSplitDB(); err == nil {
		defer db.Close()
		if rec, err := db.GetSpecialRepo(repoName); err == nil && len(rec.RemoteURL) > 0 {
			return rec.RemoteURL
		}
	}
	return probeRemoteSpecialRepoURL(repoName)
}

func resolveStandardRepoRemoteURL(name, path string) (string, bool) {
	if db, err := store.OpenDefault(); err == nil {
		defer db.Close()
		if name != "" {
			if recs, err := db.FindBySlug(name); err == nil && len(recs) > 0 {
				for _, r := range recs {
					if r.HTTPSUrl != "" {
						return r.HTTPSUrl, true
					}
					if r.SSHUrl != "" {
						return r.SSHUrl, true
					}
				}
			}
		}
		if path != "" {
			if abs, err := filepath.Abs(path); err == nil {
				if recs, err := db.FindByPath(abs); err == nil && len(recs) > 0 {
					for _, r := range recs {
						if r.HTTPSUrl != "" {
							return r.HTTPSUrl, true
						}
						if r.SSHUrl != "" {
							return r.SSHUrl, true
						}
					}
				}
			}
		}
	}

	if name != "" {
		if u := probeRemoteSpecialRepoURL(name); u != "" {
			return u, true
		}
	}

	return "", false
}

func appendNonGitRepoHints(hints []string, repoDir, output string) []string {
	if hasNotGitRepoFailure(output) {
		if isDirPresent(repoDir) {
			hints = append(hints, "directory is not a git repository; run 'gitmap clone' or 'gitmap fix' to initialize")
		} else {
			hints = append(hints, "repository directory missing on disk; run 'gitmap clone'")
		}
	}
	return hints
}

func hasNotGitRepoFailure(output string) bool {
	lower := strings.ToLower(output)
	return strings.Contains(lower, "not a git repository") ||
		strings.Contains(lower, "fatal: not a git repo")
}

func isDirPresent(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

