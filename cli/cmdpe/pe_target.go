package cmdpe

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ResolveWorkdir determines the target working directory from flags or default repo path.
func ResolveWorkdir(flagVal string) string {
	clean := strings.TrimSpace(flagVal)
	if clean != "" {
		return clean
	}

	if envVal := strings.TrimSpace(os.Getenv("GITMAP_WORKDIR")); envVal != "" {
		return envVal
	}

	return fallbackWorkdir()
}

func fallbackWorkdir() string {
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "."
	}

	return filepath.Join(home, "workspace", "repos")
}

// IsGitRepoDir checks whether a directory contains a .git metadata entry.
func IsGitRepoDir(dirPath string) bool {
	stat, err := os.Stat(filepath.Join(dirPath, constants.DotGit))
	if err != nil {
		return false
	}

	return stat.IsDir()
}

// MatchRepoFilter evaluates if a repository name matches a user-specified filter.
func MatchRepoFilter(repoName, filter string) bool {
	cleanFilter := strings.TrimSpace(filter)
	if cleanFilter == "" {
		return true
	}

	return strings.Contains(strings.ToLower(repoName), strings.ToLower(cleanFilter))
}
