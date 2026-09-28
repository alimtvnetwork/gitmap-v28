package cmdagy

import (
	"path/filepath"
	"strings"
)

func collectCandidateWorkspaces() map[string]string {
	wsMap := make(map[string]string)
	seen := make(map[string]bool)
	projects, _ := loadActiveSortedProjects()
	for _, p := range projects {
		if isValidCandidateProject(p) {
			addCandidateWorkspace(wsMap, seen, p.GetPath(), p.Name)
		}
	}
	return wsMap
}

func isValidCandidateProject(p AgyProject) bool {
	path := p.GetPath()
	if path == "" || IsRestrictedSystemOrHomeDir(path) {
		return false
	}
	return isGitRepo(path)
}

func addCandidateWorkspace(wsMap map[string]string, seen map[string]bool, path, name string) {
	if path == "" {
		return
	}
	clean := filepath.Clean(path)
	low := strings.ToLower(clean)
	if !seen[low] {
		seen[low] = true
		wsMap[clean] = name
	}
}
