// Package helptext — catalog.go defines topic catalog metadata.
package helptext

import (
	"strings"
)

var topicSummaries = map[string]string{
	"commit-in":              "Comprehensive commit automation engine with JSON author rotation, SEO templates, deduplication heuristics, and AST function intelligence.",
	"commitin":               "Comprehensive commit automation engine with JSON author rotation, SEO templates, deduplication heuristics, and AST function intelligence.",
	"commit-write":           "Comprehensive commit automation engine with JSON author rotation, SEO templates, deduplication heuristics, and AST function intelligence.",
	"os":                     "Cross-platform OS configuration, updates, mirror auto-repair, and symlink diagnostic utilities.",
	"os-update":              "Cross-platform OS configuration, updates, mirror auto-repair, and symlink diagnostic utilities.",
	"fix-mirrors":            "Cross-platform OS configuration, updates, mirror auto-repair, and symlink diagnostic utilities.",
	"fix-link":               "Inspect, validate, and repair broken symlinks, VMware shared directories, and workstation desktop links.",
	"fixlink":                "Inspect, validate, and repair broken symlinks, VMware shared directories, and workstation desktop links.",
	"fl":                     "Inspect, validate, and repair broken symlinks, VMware shared directories, and workstation desktop links.",
	"install":                "Cross-platform developer tools, runtimes, AI models, custom scripts, and workstation profiles installation manager.",
	"in":                     "Cross-platform developer tools, runtimes, AI models, custom scripts, and workstation profiles installation manager.",
	"installer":              "Multi-OS installer management, universal Unix execution ordering, Git-direct auto-committing exports, and versioning.",
	"power":                  "Cross-platform OS screen timeout and sleep management framework with SQLite state tracking and restore profiles.",
	"pw":                     "Cross-platform OS screen timeout and sleep management framework with SQLite state tracking and restore profiles.",
	"pwr":                    "Cross-platform OS screen timeout and sleep management framework with SQLite state tracking and restore profiles.",
	"vmware":                 "VMware guest shared folder mounting, open-vm-tools management, desktop symlinks, and crontab persistence.",
	"vm":                     "VMware guest shared folder mounting, open-vm-tools management, desktop symlinks, and crontab persistence.",
	"nginx":                  "High-performance HTTP server, reverse proxy, virtual host management, configuration testing, and reload operations.",
	"ngx":                    "High-performance HTTP server, reverse proxy, virtual host management, configuration testing, and reload operations.",
	"ssh":                    "SSH key pair generation, credential management, host machine enrollment, remote execution, and cluster connectivity.",
	"ssh-join":               "First-class SSH machine enrollment by user@ip or IP, host alias registration, public key authorization, and instant host recall.",
	"sj":                     "First-class SSH machine enrollment by user@ip or IP, host alias registration, public key authorization, and instant host recall.",
	"ssh-join-add":           "Enroll an SSH target machine by user@ip or IP with an alias and optional public key authorization.",
	"ssh-join-add-with-pass": "Enroll an SSH machine with RSA-encrypted password storage for seamless auto-login.",
	"ssh-join-add-pass":      "Enroll an SSH machine with RSA-encrypted password storage for seamless auto-login.",
	"ssh-join-scan":          "Scan local subnet or CIDR network for machines with open SSH port 22.",
	"ssh-join-status":        "Inspect connectivity, latency, and health of registered SSH machines.",
	"sj-add":                 "Enroll an SSH target machine by user@ip or IP with an alias and optional public key authorization.",
	"sj-add-with-pass":       "Enroll an SSH machine with RSA-encrypted password storage for seamless auto-login.",
	"sj-add-pass":            "Enroll an SSH machine with RSA-encrypted password storage for seamless auto-login.",
	"sj-scan":                "Scan local subnet or CIDR network for machines with open SSH port 22.",
	"sj-status":              "Inspect connectivity, latency, and health of registered SSH machines.",
}

// GetTopicDetailedSummary returns documentation content for a specific command topic.
func GetTopicDetailedSummary(topic string) string {
	norm := strings.ToLower(topic)
	if summary, hasSummary := topicSummaries[norm]; hasSummary {
		return summary
	}

	return "Gitmap command line utilities."
}
