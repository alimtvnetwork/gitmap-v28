package cmdssh

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderDeployRichHelp displays the styled two-column Deploy help menu.
func RenderDeployRichHelp() {
	termout.RenderMenu(buildDeployHelpMenu())
}

func buildDeployHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "Smart File, Configuration & Fleet Deployment (gitmap deploy)",
		UsageLines: []string{
			"gitmap deploy config ssh [all|<target>] [flags]",
			"gitmap deploy config ssh --file <json-file> [all|<target>] [flags]",
			"gitmap deploy <target> <source> <dest> [flags]",
			"gitmap deploy-right <target> <source> <dest> [flags]",
			"gitmap deploy-left <target> <source> <dest> [flags]",
			"gitmap sc deploy <target> <source> <dest> [flags]",
		},
		Sections: []termout.HelpSection{
			buildDeployConfigSection(),
			buildDeploySyncModesSection(),
			buildDeployTargetsSection(),
		},
		FooterFlags: buildDeployFooterFlags(),
		Tips: []string{
			"Use 'gitmap deploy keys all' or 'gitmap deploy-keys-all' to configure passwordless SSH between all nodes.",
			"Use 'gitmap ssh export-json' to export cluster topology to JSON for migration to another machine.",
			"Use 'gitmap ssh import-json <file>' to import cluster topology on a new machine.",
			"Use 'gitmap deploy config ssh' to synchronize fleet topology and credentials across nodes.",
			"Use 'gitmap deploy-right' to push local files only when newer than remote.",
			"Use 'gitmap deploy-left' to pull remote files only when newer than local.",
		},
	}
}

func buildDeployConfigSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Configuration & Fleet Deployment",
		Entries: []termout.CommandEntry{
			{Command: "deploy keys [all]", Description: "Deploy SSH public keys across fleet (alias: deploy-keys-all)"},
			{Command: "deploy config ssh [target]", Description: "Deploy local SSH topology & credentials across fleet nodes"},
			{Command: "deploy config ssh --file <f>", Description: "Deploy SSH nodes and passwords from a JSON file to fleet"},
			{Command: "deploy import [file]", Description: "Import SSH topology / cluster configuration from JSON file or stdin"},
			{Command: "deploy node-config [all]", Description: "Alias for deploy config ssh (deploy cluster topology)"},
			{Command: "deploy bin [target]", Description: "Deploy latest GitMap executable binary to remote fleet nodes"},
		},
	}
}

func buildDeploySyncModesSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Sync Modes & Transfer Direction",
		Entries: []termout.CommandEntry{
			{Command: "deploy-right", Description: "Unidirectional push: transfer when local is newer"},
			{Command: "deploy-left", Description: "Unidirectional pull: transfer when remote is newer"},
			{Command: "--sync", Description: "Bidirectional sync: transfer newer file by mtime"},
			{Command: "--overwrite (-o)", Description: "Unconditionally overwrite existing remote files"},
			{Command: "--skip (-s)", Description: "Skip existing files without prompting"},
		},
	}
}

func buildDeployTargetsSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Target Resolution",
		Entries: []termout.CommandEntry{
			{Command: "<alias>", Description: "Target node by alias (e.g. w1, w2, alpha-win)"},
			{Command: "<ip>", Description: "Target node by IP address (e.g. 192.168.1.3)"},
			{Command: "<seq>", Description: "Target node by sequence index from 'ssh ls'"},
			{Command: "<id>", Description: "Target node by database host ID (e.g. host-1)"},
		},
	}
}

func buildDeployFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "-p, --parallel <N>", Description: "Parallel workers for directory transfer (1-16, default 4)"},
		{Command: "-n, --dry-run", Description: "Preview transfer actions without modifying files"},
		{Command: "-j, --json", Description: "Output machine-readable JSON telemetry"},
		{Command: "-h, --help", Description: "Show this deployment help menu"},
	}
}
