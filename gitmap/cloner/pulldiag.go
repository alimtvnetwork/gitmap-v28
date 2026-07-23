// Package cloner — pulldiag.go provides pull diagnosis and file-lock remediation.
package cloner

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
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
	hints := make([]string, 0, 3)
	if hasUnlinkFailure(output) {
		hints = append(hints, "file lock/read-only attribute blocked replacing old files")
	}
	if hasPathLengthRisk(repoDir, output) {
		hints = append(hints, "Windows path length risk detected; use a shorter base path like C:\\src")
	}
	if strings.Contains(strings.ToLower(repoDir), "onedrive") {
		hints = append(hints, "repo is under a synced folder (OneDrive), which often locks files")
	}

	return hints
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
		if _, ok := seen[item]; ok {
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
