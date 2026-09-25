package cmdagy

import (
	"os"
	"path/filepath"
	"strings"
)

func collectCandidateWorkspaces() map[string]string {
	wsMap := make(map[string]string)
	seen := make(map[string]bool)
	projects, _ := loadActiveSortedProjects()
	for _, p := range projects {
		addCandidateWorkspace(wsMap, seen, p.GetPath(), p.Name)
	}
	cwd, _ := os.Getwd()
	if cwd != "" {
		addCandidateWorkspace(wsMap, seen, cwd, filepath.Base(cwd))
	}
	return wsMap
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
