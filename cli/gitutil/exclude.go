package gitutil

import (
	"os"
	"path/filepath"
	"strings"
)

// EnsureGitMapExcluded ensures that .gitmap/ is listed in the repository's .git/info/exclude.
func EnsureGitMapExcluded(repoRoot string) {
	excludeFile := filepath.Join(repoRoot, ".git", "info", "exclude")
	data, err := os.ReadFile(excludeFile)
	if err != nil {
		infoDir := filepath.Join(repoRoot, ".git", "info")
		if stat, statErr := os.Stat(filepath.Join(repoRoot, ".git")); statErr == nil && stat.IsDir() {
			_ = os.MkdirAll(infoDir, 0o755)
			_ = os.WriteFile(excludeFile, []byte(".gitmap/\n"), 0o644)
		}

		return
	}

	if strings.Contains(string(data), ".gitmap") {
		return
	}

	_ = os.WriteFile(excludeFile, append(data, []byte("\n.gitmap/\n")...), 0o644)
}
