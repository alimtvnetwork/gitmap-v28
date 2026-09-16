// Package cmdagy — agy_read_memory_partition.go partitions projects into target and excluded lists.
package cmdagy

import (
	"path/filepath"
	"strings"
)

func partitionPromptProjects(projects []AgyProject, tokens []string) ([]AgyProject, []AgyProject) {
	targets := make([]AgyProject, 0)
	excluded := make([]AgyProject, 0)

	for _, p := range projects {
		targets, excluded = sortProjectIntoPromptPartition(p, tokens, targets, excluded)
	}

	return targets, excluded
}

func sortProjectIntoPromptPartition(p AgyProject, tokens []string, targets, excluded []AgyProject) ([]AgyProject, []AgyProject) {
	if isIgnoredAgyProject(p) {
		return targets, excluded
	}

	if isExcludedPromptProject(p, tokens) {
		return targets, append(excluded, p)
	}

	return append(targets, p), excluded
}

func isIgnoredAgyProject(p AgyProject) bool {
	return p.ID == "outside-of-project"
}

func isExcludedPromptProject(p AgyProject, tokens []string) bool {
	path := p.GetPath()
	if path != "" && !checkDirExists(path) {
		return true
	}

	return isMatchPrefixOrSlugExcept(p, tokens)
}

func isMatchPrefixOrSlugExcept(p AgyProject, tokens []string) bool {
	if len(tokens) == 0 {
		return false
	}

	for _, t := range tokens {
		if isProjectTokenMatch(p, t) {
			return true
		}
	}

	return false
}

func isProjectTokenMatch(p AgyProject, token string) bool {
	pID := strings.ToLower(p.ID)
	pName := strings.ToLower(p.Name)
	pSlug := strings.ToLower(filepath.Base(p.GetPath()))
	pPath := strings.ToLower(filepath.Clean(p.GetPath()))

	if isTokenExactMatch(token, pID, pName, pSlug, pPath) {
		return true
	}

	return isTokenPrefixMatch(token, pID, pName, pSlug)
}

func isTokenExactMatch(token, pID, pName, pSlug, pPath string) bool {
	return token == pID || token == pName || token == pSlug || token == pPath
}

func isTokenPrefixMatch(token, pID, pName, pSlug string) bool {
	return strings.HasPrefix(pID, token) || strings.HasPrefix(pName, token) || strings.HasPrefix(pSlug, token)
}
