package cmdfoldertree

import (
	"os"
	"path/filepath"
	"strings"
)

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
