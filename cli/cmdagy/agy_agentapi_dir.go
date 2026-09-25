// Package cmdagy — agy_agentapi_dir.go resolves repository directories and defaults for agentapi operations.
package cmdagy

import (
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

// resolveAgentAPIWorkingDir determines the working directory for executing agentapi commands.
func resolveAgentAPIWorkingDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		return resolveDefaultGitmapRepoRoot()
	}
	root, rErr := gitutil.RepoRoot(cwd)
	if rErr == nil && len(root) > 0 {
		return root
	}

	return resolveDefaultGitmapRepoRoot()
}

// resolveDefaultGitmapRepoRoot finds the gitmap repository root or fallback directory.
func resolveDefaultGitmapRepoRoot() string {
	projects, err := loadActiveSortedProjects()
	if err == nil {
		return scanProjectsForGitmap(projects)
	}

	return scanFallbackPathsForGitmap()
}

func scanProjectsForGitmap(projects []AgyProject) string {
	for _, p := range projects {
		if isGitmapProjectCandidate(p) {
			return p.GetPath()
		}
	}

	return scanFallbackPathsForGitmap()
}

func scanFallbackPathsForGitmap() string {
	candidates := []string{"d:\\work\\gitmap", "c:\\work\\gitmap", "e:\\work\\gitmap"}
	for _, c := range candidates {
		if checkDirExists(c) {
			return c
		}
	}

	return resolveProjectRootDir()
}

func isGitmapProjectCandidate(p AgyProject) bool {
	nameLow := strings.ToLower(p.Name)
	pathLow := strings.ToLower(p.GetPath())
	hasGitmap := strings.Contains(nameLow, "gitmap") || strings.Contains(pathLow, "gitmap")

	return hasGitmap && checkDirExists(p.GetPath())
}
