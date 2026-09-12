// Package helptext — catalog.go defines topic catalog metadata.
package helptext

import (
	"strings"
)

var topicSummaries = map[string]string{
	"commit-in":    "Comprehensive commit automation engine with JSON author rotation, SEO templates, deduplication heuristics, and AST function intelligence.",
	"commitin":     "Comprehensive commit automation engine with JSON author rotation, SEO templates, deduplication heuristics, and AST function intelligence.",
	"commit-write": "Comprehensive commit automation engine with JSON author rotation, SEO templates, deduplication heuristics, and AST function intelligence.",
	"os":           "Cross-platform OS configuration, updates, mirror auto-repair, and symlink diagnostic utilities.",
	"os-update":    "Cross-platform OS configuration, updates, mirror auto-repair, and symlink diagnostic utilities.",
	"fix-mirrors":  "Cross-platform OS configuration, updates, mirror auto-repair, and symlink diagnostic utilities.",
	"fix-link":     "Inspect, validate, and repair broken symlinks, VMware shared directories, and workstation desktop links.",
	"fixlink":      "Inspect, validate, and repair broken symlinks, VMware shared directories, and workstation desktop links.",
	"fl":           "Inspect, validate, and repair broken symlinks, VMware shared directories, and workstation desktop links.",
	"install":      "Cross-platform developer tools, runtimes, AI models, custom scripts, and workstation profiles installation manager.",
	"in":           "Cross-platform developer tools, runtimes, AI models, custom scripts, and workstation profiles installation manager.",
	"installer":    "Multi-OS installer management, universal Unix execution ordering, Git-direct auto-committing exports, and versioning.",
	"power":        "Cross-platform OS screen timeout and sleep management framework with SQLite state tracking and restore profiles.",
	"pw":           "Cross-platform OS screen timeout and sleep management framework with SQLite state tracking and restore profiles.",
	"pwr":          "Cross-platform OS screen timeout and sleep management framework with SQLite state tracking and restore profiles.",
	"vmware":       "VMware guest shared folder mounting, open-vm-tools management, desktop symlinks, and crontab persistence.",
	"vm":           "VMware guest shared folder mounting, open-vm-tools management, desktop symlinks, and crontab persistence.",
	"nginx":        "High-performance HTTP server, reverse proxy, virtual host management, configuration testing, and reload operations.",
	"ngx":          "High-performance HTTP server, reverse proxy, virtual host management, configuration testing, and reload operations.",
}

// GetTopicDetailedSummary returns documentation content for a specific command topic.
func GetTopicDetailedSummary(topic string) string {
	norm := strings.ToLower(topic)
	if summary, ok := topicSummaries[norm]; ok {
		return summary
	}

	return "Gitmap command line utilities."
}
