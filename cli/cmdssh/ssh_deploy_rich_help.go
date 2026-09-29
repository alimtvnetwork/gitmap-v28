package cmdssh

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderDeployRichHelp displays the styled two-column Deploy help menu.
func RenderDeployRichHelp() {
	termhelp.RenderMenu(buildDeployHelpMenu())
}

func buildDeployHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "Smart File, Configuration & Fleet Deployment (gitmap deploy)",
		UsageLines: []string{
			"gitmap deploy config ssh [all|<target>] [flags]",
			"gitmap deploy config ssh --file <json-file> [all|<target>] [flags]",
			"gitmap deploy <target> <source> <dest> [flags]",
			"gitmap deploy-right <target> <source> <dest> [flags]",
			"gitmap deploy-left <target> <source> <dest> [flags]",
			"gitmap sc deploy <target> <source> <dest> [flags]",
		},
		Sections: []termhelp.HelpSection{
			buildDeployConfigSection(),
			buildDeploySyncModesSection(),
			buildDeployTargetsSection(),
		},
		FooterFlags: buildDeployFooterFlags(),
		Tips: []string{
			"Use 'gitmap deploy config ssh' to synchronize fleet topology and credentials across nodes.",
			"Use 'gitmap deploy config ssh --file <f>' to deploy from an exported JSON configuration.",
			"Use 'gitmap deploy-right' to push local files only when newer than remote.",
			"Use 'gitmap deploy-left' to pull remote files only when newer than local.",
			"Add '--json' for clean programmatic pipeline telemetry.",
		},
	}
}

func buildDeployConfigSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Configuration & Fleet Deployment",
		Entries: []termhelp.CommandEntry{
			{Command: "deploy config ssh [target]", Description: "Deploy local SSH topology & credentials across fleet nodes"},
			{Command: "deploy config ssh --file <f>", Description: "Deploy SSH nodes and passwords from a JSON file to fleet"},
			{Command: "deploy import [file]", Description: "Import SSH topology / cluster configuration from JSON file or stdin"},
			{Command: "deploy node-config [all]", Description: "Alias for deploy config ssh (deploy cluster topology)"},
			{Command: "deploy bin [target]", Description: "Deploy latest GitMap executable binary to remote fleet nodes"},
			{Command: "deploy keys [all]", Description: "Gather, deduplicate, and deploy SSH public keys across fleet"},
		},
	}
}

func buildDeploySyncModesSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Sync Modes & Transfer Direction",
		Entries: []termhelp.CommandEntry{
			{Command: "deploy-right", Description: "Unidirectional push: transfer when local is newer"},
			{Command: "deploy-left", Description: "Unidirectional pull: transfer when remote is newer"},
			{Command: "--sync", Description: "Bidirectional sync: transfer newer file by mtime"},
			{Command: "--overwrite (-o)", Description: "Unconditionally overwrite existing remote files"},
			{Command: "--skip (-s)", Description: "Skip existing files without prompting"},
		},
	}
}

func buildDeployTargetsSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Target Resolution",
		Entries: []termhelp.CommandEntry{
			{Command: "<alias>", Description: "Target node by alias (e.g. w1, w2, alpha-win)"},
			{Command: "<ip>", Description: "Target node by IP address (e.g. 192.168.1.3)"},
			{Command: "<seq>", Description: "Target node by sequence index from 'ssh ls'"},
			{Command: "<id>", Description: "Target node by database host ID (e.g. host-1)"},
		},
	}
}

func buildDeployFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "-p, --parallel <N>", Description: "Parallel workers for directory transfer (1-16, default 4)"},
		{Command: "-n, --dry-run", Description: "Preview transfer actions without modifying files"},
		{Command: "-j, --json", Description: "Output machine-readable JSON telemetry"},
		{Command: "-h, --help", Description: "Show this deployment help menu"},
	}
}
