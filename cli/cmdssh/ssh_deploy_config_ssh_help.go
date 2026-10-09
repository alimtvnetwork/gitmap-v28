// Package cmdssh — ssh_deploy_config_ssh_help.go renders help menus for deploy config ssh.
package cmdssh

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderDeployConfigSSHHelp displays the styled help menu for deploy config ssh.
func RenderDeployConfigSSHHelp() {
	termout.RenderMenu(buildDeployConfigSSHHelpMenu())
}

func buildDeployConfigSSHHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "Deploy SSH Fleet Configuration (gitmap deploy config ssh)",
		UsageLines: []string{
			"gitmap deploy config ssh [all|<target>] [flags]",
			"gitmap deploy-config ssh [all|<target>] [flags]",
			"gitmap ssh deploy config [all|<target>] [flags]",
			"gitmap deploy node-config [all|<target>] [flags]",
		},
		Sections: []termout.HelpSection{
			buildDeployConfigSSHSourceSection(),
			buildDeployConfigSSHSyncSection(),
			buildDeployConfigSSHTargetsSection(),
		},
		FooterFlags: buildDeployConfigSSHFooterFlags(),
		Tips: []string{
			"Run 'gitmap deploy config ssh' to deploy local CL config to all remote machines automatically.",
			"Run 'gitmap deploy config ssh w1' to deploy specifically to node w1.",
			"Remote machines match existing nodes by Alias or IP and update them without wiping passwords.",
			"Use '--file <path>' to deploy from an exported JSON file instead of local CL config.",
		},
	}
}

func buildDeployConfigSSHSourceSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Configuration Sources",
		Entries: []termout.CommandEntry{
			{Command: "(default)", Description: "Deploy from local GitMap CL configuration / database (SSHConnection & ssh_hosts)"},
			{Command: "-f, --file <path>", Description: "Deploy from specified JSON file (auto-probes repo-secrets if relative/short)"},
		},
	}
}

func buildDeployConfigSSHSyncSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Smart Match & Update Logic",
		Entries: []termout.CommandEntry{
			{Command: "Match by Alias/IP", Description: "Checks existing nodes on remote machines before writing"},
			{Command: "Update on Match", Description: "Updates matching nodes without creating duplicates"},
			{Command: "Preserve Passwords", Description: "Retains existing encrypted passwords if incoming is blank"},
			{Command: "Idempotent Insert", Description: "Inserts only new nodes discovered in the configuration"},
		},
	}
}

func buildDeployConfigSSHTargetsSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Target Resolution",
		Entries: []termout.CommandEntry{
			{Command: "all", Description: "Deploy across all enrolled fleet nodes (default)"},
			{Command: "<alias>", Description: "Target node by alias (e.g. w1, w2, w3)"},
			{Command: "<ip>", Description: "Target node by IP address (e.g. 192.168.1.3)"},
			{Command: "<worker-id>", Description: "Target node by worker ID (e.g. worker-1)"},
			{Command: "<seq>", Description: "Target node by sequence index from 'gitmap ssh ls'"},
		},
	}
}

func buildDeployConfigSSHFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "-f, --file <path>", Description: "Deploy from JSON file instead of local CL config"},
		{Command: "-e, --except <tokens>", Description: "Exclude nodes by alias, IP, or worker ID"},
		{Command: "-n, --dry-run", Description: "Preview deployment actions without executing"},
		{Command: "-j, --json", Description: "Output machine-readable JSON metrics"},
		{Command: "-h, --help", Description: "Display this configuration deployment help menu"},
	}
}
