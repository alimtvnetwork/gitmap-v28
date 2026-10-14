// Package cmdpending provides command implementations for pending commits and discovery.
package cmdpending

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/helpdoc"
)

func discoveredCategory(dc discoveredCommand) string {
	for _, n := range dc.names {
		if strings.Contains(strings.ToLower(n), "commit") {
			return "commits"
		}
	}
	return dc.category
}

// resolveDescription picks the best available help text: cobra Short, then the
// helpdoc topic catalog (primary name, then aliases), then a pointer to the
// command's own help.
func resolveDescription(dc discoveredCommand, name string) string {
	if strings.TrimSpace(dc.description) != "" {
		return strings.TrimSpace(dc.description)
	}
	for _, n := range dc.names {
		if summary := helpdoc.GetTopicDetailedSummary(n); isRealSummary(summary) {
			return summary
		}
	}
	return fmt.Sprintf("Run 'gitmap %s --help' for usage details.", name)
}

func isRealSummary(summary string) bool {
	trimmed := strings.TrimSpace(summary)
	return trimmed != "" && trimmed != "Gitmap command line utilities."
}

// dispatchGroupCategory maps a dispatch function (or parent command) name to
// the nc category vocabulary, derived from the router's own grouping.
func dispatchGroupCategory(funcName string) string {
	lower := strings.ToLower(funcName)
	for _, kw := range dispatchCategoryKeywords() {
		if strings.Contains(lower, kw.keyword) {
			return kw.category
		}
	}
	return "general"
}

type categoryKeyword struct {
	keyword  string
	category string
}

func dispatchCategoryKeywords() []categoryKeyword {
	return []categoryKeyword{
		{"tooling", "tooling"},
		{"commit", "commits"},
		{"pending", "commits"},
		{"backup", "commits"},
		{"core", "core"},
		{"data", "data"},
		{"release", "release"},
		{"pipeline", "diagnostics"},
		{"utility", "utility"},
		{"project", "projects"},
		{"diff", "diff"},
		{"compare", "diff"},
		{"merge", "merge"},
		{"sync", "sync"},
		{"common", "sync"},
		{"space", "space"},
		{"spec", "spec"},
		{"template", "templates"},
		{"audit", "diagnostics"},
		{"doctor", "diagnostics"},
		{"ssh", "fleet"},
		{"cluster", "fleet"},
		{"node", "fleet"},
		{"scan", "scanner"},
		{"clone", "scanner"},
		{"os", "os"},
		{"storage", "storage"},
		{"cache", "storage"},
		{"secret", "storage"},
		{"llm", "ai"},
		{"macro", "ai"},
		{"automation", "ai"},
		{"aum", "ai"},
		{"ai", "ai"},
	}
}
