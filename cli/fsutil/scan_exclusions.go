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

// canonicalScanExclusions lists path segments that mark cache, tool, or
// package-manager subtrees which repository scans must never descend into.
var canonicalScanExclusions = []string{
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

// matchesCanonicalScanExclusion reports whether a normalized slash path
// contains one of the canonical exclusion segments.
func matchesCanonicalScanExclusion(norm string) bool {
	for _, seg := range canonicalScanExclusions {
		if strings.Contains(norm, seg) || strings.HasSuffix(norm, strings.TrimSuffix(seg, "/")) {
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

	return matchesCanonicalScanExclusion(norm)
}

// IsExcludedBelowScanRoot reports whether path should be excluded from a scan
// rooted at rootDir. Unlike IsExcludedScanPath, only the portion of path
// below rootDir is judged: the operator explicitly chose rootDir, so an
// excluded segment in rootDir's own ancestors (for example a temp root under
// a ".cache" directory) must not silently disqualify every repository the
// scan was pointed at. Excluded subtrees below the root are still pruned.
func IsExcludedBelowScanRoot(rootDir, path string) bool {
	rel, err := filepath.Rel(rootDir, path)
	if err != nil || rel == "." {
		return IsExcludedScanPath(path)
	}

	if strings.HasPrefix(rel, "..") {
		return IsExcludedScanPath(path)
	}

	norm := "/" + filepath.ToSlash(strings.ToLower(rel)) + "/"

	return matchesCanonicalScanExclusion(norm)
}

// IsValidGitRepo verifies that a given directory contains a structurally sound Git repository.
func IsValidGitRepo(repoPath string) bool {
	gitPath := filepath.Join(repoPath, ".git")
	fi, err := os.Stat(gitPath)
	if err != nil {
		return false
	}

	if fi.IsDir() {
		return true
	}

	data, errRead := os.ReadFile(gitPath)
	if errRead != nil {
		return false
	}

	trimmed := strings.TrimSpace(string(data))

	return strings.HasPrefix(trimmed, "gitdir:")
}
