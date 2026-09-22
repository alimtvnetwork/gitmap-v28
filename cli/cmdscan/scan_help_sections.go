package cmdscan

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

func buildScanWalkerSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Scanner Walk & Performance",
		Entries: []termhelp.CommandEntry{
			{Command: "--workers <n>", Description: "Parallel directory walker pool size (1-16, auto: NumCPU)"},
			{Command: "--max-depth <n>", Description: "Max folder depth to descend (default: 4, -1: unlimited)"},
			{Command: "--default-branch <b>", Description: "Fallback branch name when HEAD detection returns empty"},
			{Command: "--config <path>", Description: "Path to configuration file (default: ./data/config.json)"},
		},
	}
}

func buildScanIntegrationSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Integrations & Probing",
		Entries: []termhelp.CommandEntry{
			{Command: "--fix", Description: "Run database reconciliation (removes stale rows, adds missing)"},
			{Command: "--no-vscode-sync", Description: "Skip syncing discovered repos into VS Code Project Manager"},
			{Command: "--no-auto-tags", Description: "Skip auto-derived project tags (git, node, go, etc.)"},
			{Command: "--github-desktop", Description: "Auto-register discovered repositories with GitHub Desktop"},
			{Command: "--no-probe", Description: "Skip background remote branch version probe entirely"},
		},
	}
}
