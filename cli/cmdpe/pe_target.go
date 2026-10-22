package cmdpe

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/envutil"
)

// ResolveWorkdir determines the target working directory from flags or default repo path.
func ResolveWorkdir(flagVal string) string {
	clean := strings.TrimSpace(flagVal)
	if clean != "" {
		return clean
	}

	envVal := envutil.GetTrimmed(constants.EnvWorkdir)
	if envVal != "" {
		return envVal
	}

	return fallbackWorkdir()
}

// fallbackWorkdir resolves the default workspace repository folder (<UserHomeDir>/workspace/repos, e.g. ~/workspace/repos).
// If the user's home directory resolution fails, it falls back to constants.DefaultDir (".").
func fallbackWorkdir() string {
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return constants.DefaultDir
	}

	return filepath.Join(home, constants.DefaultWorkdirWorkspace, constants.DefaultWorkdirRepos)
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
