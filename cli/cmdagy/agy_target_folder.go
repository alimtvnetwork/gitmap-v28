// Package cmdagy — agy_target_folder.go handles folder root matching and shell completion for targets.
package cmdagy

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ResolveAgyFolderTargets finds all projects whose workspace directory is inside folderRoot.
func ResolveAgyFolderTargets(folderRoot string, projects []AgyProject) ([]AgyProject, error) {
	cleanRoot, err := filepath.Abs(folderRoot)
	if err != nil {
		return nil, fmt.Errorf("resolve folder path %q: %w", folderRoot, err)
	}
	matches := collectProjectsUnderFolder(filepath.Clean(cleanRoot), projects)
	if len(matches) == 0 {
		return nil, fmt.Errorf("no registered Antigravity projects found under folder %q", folderRoot)
	}
	return matches, nil
}

func collectProjectsUnderFolder(cleanRoot string, projects []AgyProject) []AgyProject {
	var matches []AgyProject
	for _, p := range projects {
		if isProjectUnderFolder(p.GetPath(), cleanRoot) {
			matches = append(matches, p)
		}
	}
	return matches
}

func isProjectUnderFolder(projectPath, folderRoot string) bool {
	if projectPath == "" {
		return false
	}
	cleanPath, err := filepath.Abs(projectPath)
	if err != nil {
		cleanPath = filepath.Clean(projectPath)
	}
	pNorm := strings.ToLower(filepath.Clean(cleanPath))
	fNorm := strings.ToLower(filepath.Clean(folderRoot))

	return pNorm == fNorm || strings.HasPrefix(pNorm, fNorm+string(filepath.Separator))
}

// CompleteAgyProjectSuggestions builds completion suggestions (sequences, IDs, and slugs).
func CompleteAgyProjectSuggestions(toComplete string, projects []AgyProject) []string {
	var suggestions []string
	lower := strings.ToLower(toComplete)

	for i, p := range projects {
		seq := fmt.Sprintf("%03d", i+1)
		if strings.HasPrefix(seq, lower) {
			suggestions = append(suggestions, fmt.Sprintf("%s\t%s", seq, p.Name))
		}
		if strings.HasPrefix(strings.ToLower(p.Name), lower) {
			suggestions = append(suggestions, fmt.Sprintf("%s\tseq:%s", p.Name, seq))
		}
		if strings.HasPrefix(strings.ToLower(p.ID), lower) {
			suggestions = append(suggestions, fmt.Sprintf("%s\t%s", shortProjectId(p.ID), p.Name))
		}
	}

	return suggestions
}
