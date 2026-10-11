package fsutil

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var customExclusionCheckers []func(string) bool

// RegisterCustomExclusionChecker registers a dynamic exclusion predicate.
func RegisterCustomExclusionChecker(fn func(string) bool) {
	if fn != nil {
		customExclusionCheckers = append(customExclusionCheckers, fn)
	}
}

// IsDefaultScanExcluded checks whether a directory name matches DefaultScanExcludeDirs.
func IsDefaultScanExcluded(name string) bool {
	lower := strings.ToLower(strings.TrimSpace(name))
	for _, excl := range constants.DefaultScanExcludeDirs {
		if strings.EqualFold(lower, excl) {
			return true
		}
	}
	return false
}

// IsExcludedScanPath checks whether an absolute or relative directory path matches exclusions.
func IsExcludedScanPath(path string) bool {
	norm := filepath.ToSlash(strings.ToLower(filepath.Clean(path)))

	// Check dynamic / custom registered exclusion checkers first
	for _, checker := range customExclusionCheckers {
		if checker(norm) {
			return true
		}
	}

	// Canonical cache and tool path substrings
	canonicalExclusions := []string{
		"/.cursor/",
		"/.vscode/",
		"/.vscode-server/",
		"/.vscode-insiders/",
		"/.antigravity/",
		"/.antigravity_tools/",
		"/.gemini/",
		"/plugins/cache/",
		"/uv/cache/",
		"/.cache/",
		"/.pnpm-store/",
		"/pnpm-store/",
		"/node_modules/",
	}

	for _, seg := range canonicalExclusions {
		if strings.Contains(norm, seg) || strings.HasSuffix(norm, strings.TrimSuffix(seg, "/")) {
			return true
		}
	}

	return false
}

// IsValidGitRepo verifies that a given directory contains a structurally sound Git repository.
func IsValidGitRepo(repoPath string) bool {
	gitPath := filepath.Join(repoPath, ".git")
	fi, err := os.Stat(gitPath)
	if err != nil {
		return false
	}

	// If .git is a file (worktree or submodule pointer), it MUST start with "gitdir:"
	if !fi.IsDir() {
		data, errRead := os.ReadFile(gitPath)
		if errRead != nil {
			return false
		}
		trimmed := strings.TrimSpace(string(data))
		return strings.HasPrefix(trimmed, "gitdir:")
	}

	return true
}
