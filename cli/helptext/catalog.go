// Package helptext — catalog.go defines topic catalog metadata.
package helptext

import (
	"strings"
)

var topicSummaries = map[string]string{
	"commit-in":              "Comprehensive commit automation engine with JSON author rotation, SEO templates, deduplication heuristics, and AST function intelligence.",
	"clone-only-missing":     "Clone only missing repositories from JSON manifests or URLs, skipping existing directories on disk without pulling.",
	"com":                    "Clone only missing repositories from JSON manifests or URLs, skipping existing directories on disk without pulling.",
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
	"cluster":                "Comprehensive cluster management for multi-node topology, remote execution, script deployment, and Ubuntu provisioning.",
	"cluster-add":            "Enroll new nodes into cluster topology and SSH host registry with alias registration and key authorization.",
	"cluster-join":           "Join nodes to cluster using SSH enrollment and automated host registration.",
	"cluster-node":           "Provision and configure Ubuntu cluster nodes with static Netplan IP, base tools, oh-my-zsh users, and cleanup recipes.",
	"cluster-nodes":          "List registered nodes in the cluster database with display ID, role, OS, and reachability status.",
	"cluster-status":         "Inspect cluster connectivity, node health, heartbeat records, and ping round-trip latency.",
	"cluster-remove":         "Remove and unregister nodes from the cluster topology and SSH host registry.",
	"cluster-bootstrap":      "Bootstrap cluster nodes with RSA keys, passwordless sudo, and automated enrollment.",
	"cluster-exec":           "Execute remote command across targeted cluster nodes",
	"cluster-run-script":     "Deploy and execute local scripts remotely across cluster nodes with base64 staging and automatic cleanup.",
	"cluster-import":         "Import cluster topology from JSON configuration and enroll control plane and worker nodes.",
	"cluster-k8s":            "Manage Kubernetes cluster lifecycle, runtime, and components",
	"cluster-k8s-init":       "Initialize Kubernetes control plane",
	"cluster-k8s-join":       "Join worker nodes to Kubernetes control plane",
	"cluster-k8s-helm":       "Deploy Helm and storage class provisioners",
	"servers-clients":        "Distributed fan-out execution across cluster nodes",
	"sc":                     "Short alias for servers-clients",
	"clients":                "Target clients in servers-clients topology",
	"agy-fix-pipeline":       "Extract failing pipeline error logs and combine with CI/CD fix prompt for Antigravity IDE into clipboard and active temp file.",
	"fix-pipeline":           "Extract failing pipeline error logs and combine with CI/CD fix prompt for Antigravity IDE into clipboard and active temp file.",
	"author":                 "Display credentials, background, and innovations of GitMap's creator & system architect MD ALIM UL KARIM.",
	"sponsor":                "Display information about GitMap's official sponsor and engineering partner RISE UP ASIA LLC.",
	"credits":                "Display creator, sponsor, and engineering partner acknowledgments.",
	"task":                   "Manage and inspect pending and completed task execution queues, history, and undo/redo operations.",
	"tasks":                  "Manage and inspect pending and completed task execution queues, history, and undo/redo operations.",
}

// GetTopicDetailedSummary returns documentation content for a specific command topic.
func GetTopicDetailedSummary(topic string) string {
	norm := strings.ToLower(topic)
	if summary, hasSummary := topicSummaries[norm]; hasSummary {
		return summary
	}

	return "Gitmap command line utilities."
}
