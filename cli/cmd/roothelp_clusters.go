// Package cmd provides root CLI command routing, semantic help clustering,
// and automated suggestion engine dispatch interception.
package cmd

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// GetSemanticClusterSections returns the 5 canonical Muse semantic help clusters.
func GetSemanticClusterSections() []termhelp.HelpSection {
	return []termhelp.HelpSection{
		clusterCoreRepoOperations(),
		clusterReleaseCommits(),
		clusterFleetRemoteSSH(),
		clusterAIAutomation(),
		clusterSystemOSTooling(),
	}
}

// BuildSemanticClustersMenu constructs the unified 5-cluster HelpMenu.
func BuildSemanticClustersMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "gitmap — 5 Semantic Clusters",
		UsageLines: []string{
			"gitmap <command> [flags]",
			"gitmap help <topic>",
		},
		Sections:    GetSemanticClusterSections(),
		FooterFlags: buildSemanticClusterFooterFlags(),
		Tips:        buildClusterTips(),
	}
}

// RenderSemanticClustersHelp prints the 5 semantic clusters menu to standard output.
func RenderSemanticClustersHelp() {
	menu := BuildSemanticClustersMenu()
	termhelp.RenderMenu(menu)
}

func buildSemanticClusterFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "-h, --help", Description: "Show detailed help card and usage guide"},
		{Command: "-v, --verbose", Description: "Enable verbose debug logs during execution"},
		{Command: "--json", Description: "Emit results formatted as structured JSON"},
	}
}

func buildClusterTips() []string {
	return []string{
		"Run 'gitmap <command> --help' to display command-specific options.",
		"Use 'gitmap pe -t' to inspect pipeline error logs and dynamic telemetry.",
		"Use 'gitmap s' and 'gitmap c' as fast shorthands for status and clone.",
	}
}

func clusterCoreRepoOperations() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title:   "1. Core & Repo Operations",
		Color:   constants.ColorCyan,
		Entries: resolveClusterEntriesCore(),
	}
}

func resolveClusterEntriesCore() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "scan", Description: "Fast repository scanner with parallel filesystem discovery"},
		{Command: "list, ls", Description: "List recorded repositories with multi-format output"},
		{Command: "clone, c", Description: "Clone repositories from manifest, URL, or selective batch"},
		{Command: "pull, p", Description: "Batch pull remote branches with conflict and error isolation"},
		{Command: "status, s", Description: "Check working tree status across all repositories"},
		{Command: "reconcile", Description: "Reconcile divergent git branches and database reality"},
		{Command: "stash", Description: "Manage stash snapshots across repositories"},
		{Command: "wip", Description: "Save temporary work-in-progress checkpoint"},
		{Command: "group", Description: "Organize repositories into workspaces"},
		{Command: "cd", Description: "Fast interactive directory jumper to repository paths"},
	}
}

func clusterReleaseCommits() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title:   "2. Release & Commits",
		Color:   constants.ColorGreen,
		Entries: resolveClusterEntriesRelease(),
	}
}

func resolveClusterEntriesRelease() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "release, r", Description: "Automated semantic release ceremony, tagging, and publishing"},
		{Command: "changelog", Description: "Generate or view repository changelog from commits"},
		{Command: "list-versions", Description: "Display repository version history and release metadata"},
		{Command: "commit, co", Description: "Structured conventional commit runner with validation"},
		{Command: "cpf", Description: "Commit, push, and create feature branch in one pass"},
		{Command: "cpb", Description: "Commit and push bugfix branch with automated telemetry"},
		{Command: "cpr", Description: "Commit, push, and trigger release workflow orchestration"},
		{Command: "cpar", Description: "Commit and push across all connected repository workspaces"},
	}
}

func clusterFleetRemoteSSH() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title:   "3. Fleet & Remote SSH",
		Color:   constants.ColorYellow,
		Entries: resolveClusterEntriesFleet(),
	}
}

func resolveClusterEntriesFleet() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "ssh", Description: "Manage remote SSH keys, nodes, and terminal sessions"},
		{Command: "nodes", Description: "Fleet node discovery, health monitoring, and latency probes"},
		{Command: "cluster", Description: "Multi-node cluster delegation and remote job execution"},
		{Command: "sc", Description: "Servers and clients connection manager and topology view"},
		{Command: "deploy", Description: "Orchestrate multi-node fleet deployments and asset distribution"},
		{Command: "vhost", Description: "Inspect and manage virtual hosts across infrastructure"},
		{Command: "nginx", Description: "Nginx reverse proxy configurations and automated SSL renew"},
	}
}

func clusterAIAutomation() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title:   "4. AI & Automation",
		Color:   constants.ColorMagenta,
		Entries: resolveClusterEntriesAI(),
	}
}

func resolveClusterEntriesAI() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "agy", Description: "Antigravity developer tools, prompts, and queue manager"},
		{Command: "agm", Description: "Manage Antigravity agent instances and logs"},
		{Command: "aum", Description: "Automation Unit Manager for macros, scripts, and AST transforms"},
		{Command: "ai", Description: "AI analysis, session tracking, and model training"},
		{Command: "macro", Description: "Record and replay CLI macros and workflows"},
		{Command: "pipeline, pl", Description: "CI/CD pipeline diagnostics, error extraction, and telemetry"},
		{Command: "pe", Description: "Pipeline error analyzer with dynamic runner log extraction"},
		{Command: "rerun, rr", Description: "Rerun failed agent steps and pipeline jobs"},
		{Command: "sug", Description: "Shutdown execution until quality gates are green"},
	}
}

func clusterSystemOSTooling() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title:   "5. System, OS & Developer Tooling",
		Color:   constants.ColorBlue,
		Entries: resolveClusterEntriesSystem(),
	}
}

func resolveClusterEntriesSystem() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "os", Description: "Cross-platform OS desktop and taskbar integrations"},
		{Command: "apps", Description: "Discover, inventory, and cleanly uninstall desktop applications"},
		{Command: "install", Description: "Install CLI companion tools and agents"},
		{Command: "uninstall", Description: "Uninstall CLI companion tools and agents"},
		{Command: "storage", Description: "Inspect disk usage, large repo locator, and volume cleanup"},
		{Command: "clean", Description: "Clean working tree and purge build artifacts"},
		{Command: "vscode", Description: "Synchronize VS Code workspaces and settings"},
		{Command: "sync", Description: "Synchronize prompts and repositories across fleet"},
	}
}

// IsSemanticClusterName reports whether the input string matches any cluster identifier.
func IsSemanticClusterName(name string) bool {
	low := strings.ToLower(strings.TrimSpace(name))
	switch low {
	case "core", "release", "fleet", "ai", "system", "clusters", "5-clusters":
		return true
	default:
		return false
	}
}
