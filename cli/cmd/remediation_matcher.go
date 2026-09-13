package cmd

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// FindRemediationItem locates a matching item in items by exact name,
// normalized path/base, absolute path, substring, or 1-based index.
func FindRemediationItem(items []RemediationItem, query string) *RemediationItem {
	cleanQuery := strings.TrimSpace(query)
	if cleanQuery == "" {
		return nil
	}

	if item := matchByNameOrPath(items, cleanQuery); item != nil {
		return item
	}

	if item := matchBySubstring(items, cleanQuery); item != nil {
		return item
	}

	return matchByIndex(items, cleanQuery)
}

func matchByNameOrPath(items []RemediationItem, query string) *RemediationItem {
	base := resolveTargetBase(query)
	absQuery, absErr := filepath.Abs(query)
	hasAbs := absErr == nil

	for i := range items {
		if isItemMatchingNameOrBase(&items[i], query, base) {
			return &items[i]
		}

		if isItemMatchingPath(&items[i], query, absQuery, hasAbs) {
			return &items[i]
		}
	}

	return nil
}

func isItemMatchingNameOrBase(item *RemediationItem, query, base string) bool {
	if strings.EqualFold(item.RepoName, query) || strings.EqualFold(filepath.Base(item.RepoPath), query) {
		return true
	}

	hasBase := base != ""
	if hasBase && (strings.EqualFold(item.RepoName, base) || strings.EqualFold(filepath.Base(item.RepoPath), base)) {
		return true
	}

	return false
}

func isItemMatchingPath(item *RemediationItem, query, absQuery string, hasAbs bool) bool {
	cleanItemPath := filepath.Clean(item.RepoPath)
	if strings.EqualFold(cleanItemPath, filepath.Clean(query)) {
		return true
	}

	if hasAbs && strings.EqualFold(cleanItemPath, absQuery) {
		return true
	}

	return false
}

func resolveTargetBase(query string) string {
	base := extractQueryBase(query)
	if base != "" {
		return base
	}

	isCurrentDir := query == "." || query == "./" || query == ".\\"
	if isCurrentDir == false {
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

func matchBySubstring(items []RemediationItem, query string) *RemediationItem {
	qLower := strings.ToLower(query)
	base := strings.ToLower(resolveTargetBase(query))

	for i := range items {
		nameLower := strings.ToLower(items[i].RepoName)
		if strings.Contains(nameLower, qLower) {
			return &items[i]
		}

		hasBase := base != ""
		if hasBase && strings.Contains(nameLower, base) {
			return &items[i]
		}
	}

	return nil
}

func matchByIndex(items []RemediationItem, query string) *RemediationItem {
	num, err := strconv.Atoi(query)
	isValidIndex := err == nil && num > 0 && num <= len(items)
	if isValidIndex {
		return &items[num-1]
	}

	return nil
}

// RemoveRemediationItem removes a repo from the persisted remediation state.
func RemoveRemediationItem(repoName string) {
	items := LoadRemediationState()
	var remaining []RemediationItem
	for _, item := range items {
		isTarget := strings.EqualFold(item.RepoName, repoName)
		if isTarget == false {
			remaining = append(remaining, item)
		}
	}

	hasNone := len(remaining) == 0
	if hasNone {
		statePath := getRemediationStateFile()
		if err := os.Remove(statePath); err != nil && !os.IsNotExist(err) {
			return
		}

		return
	}

	if err := SaveRemediationState(remaining); err != nil {
		return
	}
}
