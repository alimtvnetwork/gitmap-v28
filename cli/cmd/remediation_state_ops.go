package cmd

import (
	"os"
	"path/filepath"
	"strings"
)

func resolveTargetBase(query string) string {
	base := extractQueryBase(query)
	if base != "" {
		return base
	}

	isCurrentDir := query == "." || query == "./" || query == ".\\"
	if !isCurrentDir {
		return ""
	}

	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	return filepath.Base(cwd)
}

func extractQueryBase(query string) string {
	trimmed := strings.TrimRight(query, "/\\")
	if trimmed == "" {
		return ""
	}

	norm := strings.ReplaceAll(trimmed, "\\", "/")
	lastSlash := strings.LastIndex(norm, "/")
	if lastSlash >= 0 {
		return norm[lastSlash+1:]
	}

	base := filepath.Base(filepath.Clean(trimmed))
	isSpecial := base == "." || base == "/" || base == "\\"
	if isSpecial {
		return ""
	}

	return base
}

// RemoveRemediationItem removes a repo from the persisted remediation state.
func RemoveRemediationItem(repoName string) {
	items := LoadRemediationState()
	var remaining []RemediationItem
	for _, item := range items {
		isTarget := strings.EqualFold(item.RepoName, repoName)
		if !isTarget {
			remaining = append(remaining, item)
		}
	}

	hasNone := len(remaining) == 0
	if hasNone {
		removeRemediationStateFile()

		return
	}

	if err := SaveRemediationState(remaining); err != nil {
		return
	}
}

func removeRemediationStateFile() {
	statePath := getRemediationStateFile()
	if err := os.Remove(statePath); err != nil && !os.IsNotExist(err) {
		return
	}
}
