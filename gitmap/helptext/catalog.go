// Package helptext — catalog.go defines topic catalog metadata.
package helptext

import (
	"strings"
)

// GetTopicDetailedSummary returns documentation content for a specific command topic.
func GetTopicDetailedSummary(topic string) string {
	switch strings.ToLower(topic) {
	case "commit-in", "commitin", "commit-write":
		return "Comprehensive commit automation engine with JSON author rotation, SEO templates, deduplication heuristics, and AST function intelligence."
	case "os", "os-update", "fix-mirrors":
		return "Cross-platform OS configuration, updates, mirror auto-repair, and symlink diagnostic utilities."
	case "fix-link", "fixlink", "fl":
		return "Inspect, validate, and repair broken symlinks, VMware shared directories, and workstation desktop links."
	case "installer", "in":
		return "Multi-OS installer management, universal Unix execution ordering, Git-direct auto-committing exports, and versioning."
	case "power", "pw", "pwr":
		return "Cross-platform OS screen timeout and sleep management framework with SQLite state tracking and restore profiles."
	default:
		return "Gitmap command line utilities."
	}
}
