// Package cmd — prompt_exclusion_filter.go filters out excluded target repositories.
package cmd

import (
	"path/filepath"
	"strings"
)

func FilterPromptExclusions(targets []string, excludeCSV string) []string {
	if excludeCSV == "" {
		return targets
	}

	excludes := strings.Split(excludeCSV, ",")
	excludeMap := make(map[string]bool)
	for _, e := range excludes {
		excludeMap[strings.TrimSpace(e)] = true
	}

	var filtered []string
	for _, t := range targets {
		name := filepath.Base(t)
		if !excludeMap[name] && !excludeMap[t] {
			filtered = append(filtered, t)
		}
	}
	return filtered
}
