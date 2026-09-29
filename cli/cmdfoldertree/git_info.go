package cmdfoldertree

import (
	"os"
	"path/filepath"
	"strings"
)

// DetectGitInfo checks if dirPath is a Git repo and returns (isGit, repoName, branch).
func DetectGitInfo(dirPath string) (bool, string, string) {
	gitDir := filepath.Join(dirPath, ".git")
	info, err := os.Stat(gitDir)
	if err != nil {
		return false, "", ""
	}
	repoName := filepath.Base(dirPath)
	branch := resolveGitBranch(gitDir, info.IsDir())
	remoteRepo := resolveGitRemoteRepo(gitDir, info.IsDir())
	if remoteRepo != "" {
		repoName = remoteRepo
	}
	return true, repoName, branch
}

func resolveGitBranch(gitDir string, isDir bool) string {
	headPath := resolveHeadPath(gitDir, isDir)
	data, err := os.ReadFile(headPath)
	if err != nil {
		return ""
	}
	trimmed := strings.TrimSpace(string(data))
	if strings.HasPrefix(trimmed, "ref: refs/heads/") {
		return strings.TrimPrefix(trimmed, "ref: refs/heads/")
	}
	if len(trimmed) >= 7 {
		return trimmed[:7]
	}
	return trimmed
}

func resolveHeadPath(gitDir string, isDir bool) string {
	if isDir {
		return filepath.Join(gitDir, "HEAD")
	}
	content, err := os.ReadFile(gitDir)
	if err != nil {
		return ""
	}
	line := strings.TrimSpace(string(content))
	realGitDir := strings.TrimPrefix(line, "gitdir: ")
	return filepath.Join(realGitDir, "HEAD")
}

func resolveGitRemoteRepo(gitDir string, isDir bool) string {
	configPath := resolveConfigPath(gitDir, isDir)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return ""
	}
	return parseRemoteFromConfig(string(data))
}

func resolveConfigPath(gitDir string, isDir bool) string {
	if isDir {
		return filepath.Join(gitDir, "config")
	}
	content, err := os.ReadFile(gitDir)
	if err != nil {
		return ""
	}
	line := strings.TrimSpace(string(content))
	realGitDir := strings.TrimPrefix(line, "gitdir: ")
	return filepath.Join(realGitDir, "config")
}

func parseRemoteFromConfig(cfg string) string {
	lines := strings.Split(cfg, "\n")
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if !strings.HasPrefix(trimmed, "url = ") {
			continue
		}
		rawURL := strings.TrimPrefix(trimmed, "url = ")
		return extractRepoNameFromURL(rawURL)
	}
	return ""
}

func extractRepoNameFromURL(rawURL string) string {
	clean := strings.TrimSuffix(strings.TrimSpace(rawURL), ".git")
	idx := strings.LastIndexAny(clean, "/:")
	if idx >= 0 && idx+1 < len(clean) {
		return clean[idx+1:]
	}
	return clean
}
