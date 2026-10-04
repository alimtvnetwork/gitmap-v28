// Package cmd — help.go provides terminal help renderers and documentation
// for compact GitMap workflows including modular multi-tool installation
// (`gitmap install --tools`) and desktop application discovery & uninstallation (`gitmap apps`).
package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// AppsHelpSummary is the concise overview of the apps management subsystem.
const AppsHelpSummary = "Discover, inspect, inventory, and cleanly uninstall desktop and system applications across Linux and Windows."

// InstallToolsHelpSummary is the concise overview of the modular multi-tool installer.
const InstallToolsHelpSummary = "Batch install multiple developer tools sequentially with comma-delimited flags, variadic arguments, and JSON telemetry."

// RenderAppsHelpMenu renders the modern framed-box help menu for the `gitmap apps` subsystem.
func RenderAppsHelpMenu() {
	menu := termhelp.HelpMenu{
		Title: "gitmap apps",
		UsageLines: []string{
			"gitmap apps [command] [flags]",
			"gitmap apps list [--filter <pattern>] [--json] [--system] [--user] [--all]",
			"gitmap apps uninstall <app-id> [--purge] [--force] [--json]",
			"gitmap apps help",
		},
		Sections: []termhelp.HelpSection{
			{
				Title: "Description",
				Entries: []termhelp.CommandEntry{
					{Command: "apps", Description: AppsHelpSummary},
				},
			},
			{
				Title: "Subcommands",
				Entries: []termhelp.CommandEntry{
					{Command: "list, ls", Description: "Inventory desktop apps (.desktop, registry), package managers (apt, snap, winget, choco, npm), and CLI binaries"},
					{Command: "uninstall, rm", Description: "Surgically remove an app, purging desktop launchers, icon assets, symlinks, and cache directories"},
					{Command: "help, -h", Description: "Display this comprehensive apps help documentation"},
				},
			},
			{
				Title: "List Flags",
				Entries: []termhelp.CommandEntry{
					{Command: "-f, --filter <str>", Description: "Fuzzy match by application name, executable path, or package ID"},
					{Command: "--json", Description: "Emit machine-readable JSON array envelope for automated auditing"},
					{Command: "--system", Description: "Filter exclusively system-wide applications (/usr/share/applications, Program Files)"},
					{Command: "--user", Description: "Filter exclusively user-scoped applications (~/.local/share/applications, AppData)"},
					{Command: "-a, --all", Description: "Include system utilities, hidden entries, and NoDisplay launchers"},
				},
			},
			{
				Title: "Uninstall Flags",
				Entries: []termhelp.CommandEntry{
					{Command: "--purge", Description: "Completely purge associated application configuration, caches, and icon files"},
					{Command: "--force", Description: "Bypass confirmation prompts and forcefully execute uninstallation"},
					{Command: "--json", Description: "Emit structured AppUninstallResponse JSON telemetry to stdout"},
				},
			},
			{
				Title: "Examples",
				Entries: []termhelp.CommandEntry{
					{Command: "gitmap apps list", Description: "List all installed desktop applications in interactive table format"},
					{Command: "gitmap apps list --filter antigravity --json", Description: "Discover applications matching 'antigravity' in structured JSON"},
					{Command: "gitmap apps uninstall antigravity-tools --purge --force", Description: "Completely purge an obsolete app and its launcher icon"},
					{Command: "gitmap apps uninstall clot --json", Description: "Remove an npm global CLI tool with JSON telemetry"},
				},
			},
		},
		FooterFlags: []termhelp.CommandEntry{
			{Command: "-h, --help", Description: "Display this formatted help screen"},
			{Command: "--json", Description: "Emit results formatted as JSON"},
		},
		Tips: []string{
			"Use 'gitmap apps list --json' for CI/CD and automated workstation audits.",
			"Run 'gitmap apps uninstall <id> --purge' to clean up orphaned desktop entries and caches.",
		},
	}

	termhelp.RenderMenu(menu)
	printUsageFooterShort()
}

// RenderInstallToolsHelpMenu renders the modern framed-box help menu for `gitmap install --tools`.
func RenderInstallToolsHelpMenu() {
	menu := termhelp.HelpMenu{
		Title: "gitmap install --tools",
		UsageLines: []string{
			"gitmap install --tools <tool1,tool2,...> [flags]",
			"gitmap in <tool1> <tool2> <tool3>... [flags]",
			"gitmap install --tools <list> --json [--ignore-errors]",
		},
		Sections: []termhelp.HelpSection{
			{
				Title: "Description",
				Entries: []termhelp.CommandEntry{
					{Command: "install --tools", Description: InstallToolsHelpSummary},
				},
			},
			{
				Title: "Execution Syntax & Flags",
				Entries: []termhelp.CommandEntry{
					{Command: "--tools <list>", Description: "Comma-separated list of developer tools to install sequentially (e.g. antigravity,chrome,vscode)"},
					{Command: "<tool1> <tool2>...", Description: "Variadic multiple arguments for fast multi-tool setup (e.g. gitmap in chrome vscode)"},
					{Command: "--json", Description: "Output BatchInstallResponse JSON envelope with per-tool status, version, and durations"},
					{Command: "--ignore-errors", Description: "Continue batch installation even if an individual tool installation encounters an error"},
					{Command: "-y, --yes", Description: "Automatic yes to installation confirmation prompts"},
					{Command: "--manager <mgr>", Description: "Specify explicit package manager override (apt, snap, winget, choco, direct)"},
				},
			},
			{
				Title: "Examples",
				Entries: []termhelp.CommandEntry{
					{Command: "gitmap install --tools antigravity,chrome,vscode,flameshot", Description: "Install 4 essential developer tools sequentially in one command"},
					{Command: "gitmap in antigravity chrome vscode flameshot", Description: "Variadic multi-tool installation syntax without commas"},
					{Command: "gitmap install --tools chrome,vscode --manager apt", Description: "Install tools sequentially enforcing apt package manager"},
					{Command: "gitmap install --tools chrome,vscode --json", Description: "Batch install emitting machine-readable JSON envelope to stdout"},
					{Command: "gitmap install --tools node,python,go --ignore-errors", Description: "Batch install continuing past any tool failures"},
				},
			},
		},
		FooterFlags: []termhelp.CommandEntry{
			{Command: "-h, --help", Description: "Display this formatted help screen"},
			{Command: "--json", Description: "Emit results formatted as JSON"},
			{Command: "--ignore-errors", Description: "Continue past tool failures"},
			{Command: "-y, --yes", Description: "Automatic yes to prompts"},
		},
		Tips: []string{
			"Combine with 'gitmap ssh <host>' to migrate entire toolchains to remote nodes in 1 command:",
			"  $ gitmap ssh u1 \"gitmap install --tools antigravity,chrome,vscode,flameshot --json\"",
			"Use '--json' for automated pipeline verification and telemetry tracking.",
		},
	}

	termhelp.RenderMenu(menu)
	printUsageFooterShort()
}

// PrintAppsTerminalHelp writes the plain-text manual help page for gitmap apps to os.Stdout.
func PrintAppsTerminalHelp() {
	helpText := fmt.Sprintf(`gitmap apps - %s

Usage:
  gitmap apps [command] [flags]
  gitmap apps list [--filter <pattern>] [--json] [--system] [--user] [--all]
  gitmap apps uninstall <app-id> [--purge] [--force] [--json]

Available Commands:
  list, ls         List all installed applications, desktop launchers, and CLI tools
  uninstall, rm    Surgically remove or purge an application and its desktop entries
  help             Show this help message

Flags:
  --filter, -f     Fuzzy match name, executable, or package identifier
  --json           Output structured JSON array
  --system         Filter only system-wide applications
  --user           Filter only user-scoped applications
  --all, -a        Include utilities and NoDisplay desktop entries
  --purge          Delete associated application configuration and icon files
  --force          Bypass confirmation prompts
  --help, -h       Show help for apps command

Examples:
  # List all desktop applications in interactive table format:
  $ gitmap apps list

  # Discover applications matching 'antigravity' in structured JSON:
  $ gitmap apps list --filter antigravity --json

  # Completely purge an obsolete app and its launcher icon:
  $ gitmap apps uninstall antigravity-tools --purge --force

  # Remove an npm global CLI tool on Windows:
  $ gitmap apps uninstall clot --json
`, AppsHelpSummary)
	fmt.Fprint(os.Stdout, helpText)
}

// PrintInstallToolsTerminalHelp writes the plain-text manual help page for gitmap install --tools to os.Stdout.
func PrintInstallToolsTerminalHelp() {
	helpText := fmt.Sprintf(`gitmap install --tools - %s

Usage:
  gitmap install --tools <tool1,tool2,...> [flags]
  gitmap in <tool1> <tool2> <tool3>... [flags]
  gitmap install --tools <list> --json [--ignore-errors]

Flags:
  --tools          Comma-separated list of developer tools to install
  --json           Output installation progress and summary in structured JSON
  --ignore-errors  Continue batch install even if an individual tool fails
  --manager        Specify package manager override (apt, snap, winget, choco)
  --yes, -y        Automatic yes to prompts
  --help, -h       Show this install tools help

Examples:
  $ gitmap install --tools antigravity,chrome,vscode,flameshot
  $ gitmap in antigravity chrome vscode flameshot
  $ gitmap install --tools chrome,vscode --json
  $ gitmap install --tools node,python,go --ignore-errors
`, InstallToolsHelpSummary)
	fmt.Fprint(os.Stdout, helpText)
}

// AgyDeployHelpSummary is the concise overview of the Antigravity fleet deployment subsystem.
const AgyDeployHelpSummary = "Deploy Antigravity IDE themes, presets, official plugins (4), and skills (43) to remote SSH fleet nodes."

// RenderAgyDeployHelpMenu renders the modern framed-box help menu for `gitmap agy deploy` and `gitmap deploy ide`.
func RenderAgyDeployHelpMenu() {
	menu := termhelp.HelpMenu{
		Title: "gitmap agy deploy / gitmap deploy ide",
		UsageLines: []string{
			"gitmap agy deploy <target-node> [flags]",
			"gitmap deploy ide <target-node> [flags]",
			"gitmap agy deploy u1 --all [--json] [--dry-run]",
		},
		Sections: []termhelp.HelpSection{
			{
				Title: "Description",
				Entries: []termhelp.CommandEntry{
					{Command: "agy deploy", Description: AgyDeployHelpSummary},
					{Command: "deploy ide", Description: "Convenience alias for deploying complete Antigravity IDE bundle to fleet nodes"},
				},
			},
			{
				Title: "Deployment Flags",
				Entries: []termhelp.CommandEntry{
					{Command: "--all", Description: "Deploy complete bundle (preset + theme + plugins + skills + sanitization)"},
					{Command: "--preset", Description: "Deploy execution and permission presets (eager, turbo, default)"},
					{Command: "--theme", Description: "Deploy UI theme seeds (#BD93F9, #19191C, #F8F8F2) and settings.json"},
					{Command: "--plugins", Description: "Deploy 4 core official plugins (chrome-devtools, data-agent-kit, sdk, web-guidance)"},
					{Command: "--skills", Description: "Deploy all 43 plugin skills into ~/.gemini/config/plugins/"},
					{Command: "--binaries", Description: "Verify and configure language_server binary permissions"},
					{Command: "--projects", Description: "Normalize per-project descriptors in ~/.gemini/config/projects/"},
					{Command: "--dry-run", Description: "Simulate operations and display planned changes without modifying remote state"},
					{Command: "--json", Description: "Output structured AgyDeployResultJSON telemetry to stdout"},
					{Command: "--restart", Description: "Restart remote Antigravity background services after deployment"},
					{Command: "--force", Description: "Overwrite existing remote configurations without prompting"},
				},
			},
			{
				Title: "Examples",
				Entries: []termhelp.CommandEntry{
					{Command: "gitmap agy deploy u1 --all", Description: "Deploy complete Antigravity configuration to Ubuntu node u1"},
					{Command: "gitmap deploy ide u1", Description: "Deploy full IDE configuration to u1 via the fleet deploy router"},
					{Command: "gitmap agy deploy u1 --all --dry-run --json", Description: "Simulate deployment with machine-readable JSON telemetry"},
					{Command: "gitmap agy deploy u1 --plugins --skills", Description: "Deploy only plugins and agent skills without modifying themes"},
				},
			},
		},
		FooterFlags: []termhelp.CommandEntry{
			{Command: "-h, --help", Description: "Display this formatted help screen"},
			{Command: "--json", Description: "Emit results formatted as JSON"},
			{Command: "--dry-run", Description: "Simulate deployment mutations"},
		},
		Tips: []string{
			"Use 'gitmap deploy ide <target>' for fast turnkey IDE replication to new compute nodes.",
			"Run 'gitmap agy deploy <target> --dry-run --json' to verify configuration diffs in CI/CD pipelines.",
		},
	}

	termhelp.RenderMenu(menu)
	printUsageFooterShort()
}

// PrintAgyDeployTerminalHelp writes the plain-text manual help page for gitmap agy deploy to os.Stdout.
func PrintAgyDeployTerminalHelp() {
	helpText := fmt.Sprintf(`gitmap agy deploy / gitmap deploy ide - %s

Usage:
  gitmap agy deploy <target-node> [flags]
  gitmap deploy ide <target-node> [flags]

Flags:
  --all            Deploy complete bundle (preset + theme + plugins + skills + sanitization)
  --preset         Deploy execution and permission presets (eager, turbo)
  --theme          Deploy UI theme seeds (#BD93F9, #19191C) and settings.json
  --plugins        Deploy the 4 core official plugins
  --skills         Deploy all 43 plugin skills
  --binaries       Verify and set language_server binary permissions
  --projects       Normalize per-project descriptors in ~/.gemini/config/projects/
  --dry-run        Simulate deployment and display planned mutations without changes
  --json           Output structured AgyDeployResultJSON telemetry to stdout
  --restart        Restart remote Antigravity services after deployment
  --force          Overwrite existing configurations without prompting
  --target <node>  Explicit target node selector
  --help, -h       Show help for deployment command

Examples:
  # Deploy complete Antigravity configuration to Ubuntu fleet node u1:
  $ gitmap agy deploy u1 --all

  # Turnkey IDE deployment via fleet router:
  $ gitmap deploy ide u1

  # Simulate deployment with machine-readable JSON output:
  $ gitmap agy deploy u1 --all --dry-run --json

  # Deploy only plugins and skills:
  $ gitmap agy deploy u1 --plugins --skills
`, AgyDeployHelpSummary)
	fmt.Fprint(os.Stdout, helpText)
}

