// Package cmdagy — agy_target_match.go matches targets by ID or slug.
package cmdagy

import "strings"

func matchByID(token string, projects []AgyProject) []AgyProject {
	var matches []AgyProject
	lower := strings.ToLower(token)
	for _, p := range projects {
		if strings.ToLower(p.ID) == lower || strings.HasPrefix(strings.ToLower(p.ID), lower) {
			matches = append(matches, p)
		}
	}

	return matches
}

func matchBySlug(token string, projects []AgyProject) []AgyProject {
	var matches []AgyProject
	lower := strings.ToLower(token)
	for _, p := range projects {
		nameLower := strings.ToLower(p.Name)
		if strings.HasPrefix(nameLower, lower) || strings.Contains(nameLower, lower) {
			matches = append(matches, p)
		}
	}

	return matches
}
