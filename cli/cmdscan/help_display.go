package cmdscan

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdisplay"
)

// ScanHelpDisplay builds the scan help display: the DRY source of truth for
// the scan help menu (spec 243.3 wave C), bound to helpdisplay.Displayer.
// Content is migrated faithfully from the termhelp builders in
// scan_help_menu.go and scan_help_sections.go; helptext/scan.md is
// generated from this builder.
func ScanHelpDisplay() helpdisplay.Displayer {
	return helpdisplay.NewHelpDisplay(scanHeader(), scanHelpGroups(), scanHelpSuggestions(), scanHelpTheme())
}

// scanHeader composes the display title plus the usage block.
func scanHeader() string {
	lines := []string{
		"Repository Discovery Scanner (gitmap scan)",
		"",
		"  Usage:",
		"    gitmap scan [dir] [flags]",
		"    gitmap s [dir] [flags]",
		"    gitmap scan --fix",
		"    gitmap scan ~ --force-include .oh-my-zsh",
		"    gitmap scan ~ -fi .oh-my-zsh,node_modules",
		"    gitmap scan /home/user --force-include all",
	}
	return strings.Join(lines, "\n")
}

// scanHelpGroups assembles the scan help sections; footer flags and tips
// close the display in the same order the old menu rendered them.
func scanHelpGroups() []helpdisplay.CommandHelpGroup {
	return []helpdisplay.CommandHelpGroup{
		helpdisplay.NewCommandHelpGroup("Output & Formatting:", scanOutputCommands(), nil),
		helpdisplay.NewCommandHelpGroup("Scanner Walk & Performance:", scanWalkerCommands(), nil),
		helpdisplay.NewCommandHelpGroup("Integrations & Probing:", scanIntegrationCommands(), nil),
		helpdisplay.NewCommandHelpGroup("Flags:", scanFooterCommands(), scanHelpTips()),
	}
}

// scanOutputCommands lists the output & formatting flags.
func scanOutputCommands() []helpdisplay.CommandHelper {
	return []helpdisplay.CommandHelper{
		helpdisplay.NewCommandHelper("--output <mode>", "Output format: terminal (default), csv, or json", "", ""),
		helpdisplay.NewCommandHelper("--output-path <dir>", "Directory to write scan artifacts (.gitmap/output)", "", ""),
		helpdisplay.NewCommandHelper("--manifest <dir>", "Alias for --output-path (aligned with reclone)", "", ""),
		helpdisplay.NewCommandHelper("--relative-root <dir>", "Pin base directory for byte-stable relative paths", "", ""),
	}
}

// scanWalkerCommands lists the scanner walk & performance flags.
func scanWalkerCommands() []helpdisplay.CommandHelper {
	return []helpdisplay.CommandHelper{
		helpdisplay.NewCommandHelper("--workers <n>", "Parallel directory walker pool size (1-16, auto: NumCPU)", "", ""),
		helpdisplay.NewCommandHelper("--max-depth <n>", "Max folder depth to descend (default: 4, -1: unlimited)", "", ""),
		helpdisplay.NewCommandHelper("--default-branch <b>", "Fallback branch name when HEAD detection returns empty", "", ""),
		helpdisplay.NewCommandHelper("--force-include <dir>, -fi", "Force scan of default excluded directories (.oh-my-zsh, or 'all')", "", ""),
		helpdisplay.NewCommandHelper("--config <path>", "Path to configuration file (default: ./data/config.json)", "", ""),
	}
}

// scanIntegrationCommands lists the integration & probing flags.
func scanIntegrationCommands() []helpdisplay.CommandHelper {
	return []helpdisplay.CommandHelper{
		helpdisplay.NewCommandHelper("--fix", "Run database reconciliation (removes stale rows, adds missing)", "", ""),
		helpdisplay.NewCommandHelper("--no-vscode-sync", "Skip syncing discovered repos into VS Code Project Manager", "", ""),
		helpdisplay.NewCommandHelper("--no-auto-tags", "Skip auto-derived project tags (git, node, go, etc.)", "", ""),
		helpdisplay.NewCommandHelper("--github-desktop", "Auto-register discovered repositories with GitHub Desktop", "", ""),
		helpdisplay.NewCommandHelper("--no-probe", "Skip background remote branch version probe entirely", "", ""),
	}
}

// scanFooterCommands lists the trailing standard flags.
func scanFooterCommands() []helpdisplay.CommandHelper {
	return []helpdisplay.CommandHelper{
		helpdisplay.NewCommandHelper("--quiet", "Suppress interactive walker spinner and clone hints", "", ""),
		helpdisplay.NewCommandHelper("--open", "Automatically open output directory after scan completes", "", ""),
		helpdisplay.NewCommandHelper("-v, --verbose", "Emit verbose directory traversal and probe logs", "", ""),
		helpdisplay.NewCommandHelper("-h, --help", "Show this scanner help menu", "", ""),
	}
}

// scanHelpTips returns the scan tips rendered after the flag list.
func scanHelpTips() []string {
	return []string{
		"Tip: Run 'gitmap scan . --output json' to generate a machine-readable repo catalog.",
		"Tip: Use 'gitmap scan --fix' to reconcile local disk with GitMap database.",
		"Tip: Pass '--max-depth 2' to quickly scan shallow project directories.",
		"Tip: Pass '--force-include .oh-my-zsh' (alias: -fi) to force scanning excluded directories or 'all' for everything.",
	}
}

// scanHelpSuggestions returns contextual suggestions for scan (none today).
func scanHelpSuggestions() []helpdisplay.Suggestion {
	return nil
}

// scanHelpTheme wires the displayer to the existing help color roles.
func scanHelpTheme() *helpdisplay.Theme {
	return helpdisplay.NewTheme(nil).
		WithHeaderColor(helpdisplay.NewColor(constants.ColorYellow)).
		WithCommandColor(helpdisplay.NewColor(constants.ColorGreen)).
		WithDescriptionColor(helpdisplay.NewColor(constants.ColorWhite)).
		WithHintColor(helpdisplay.NewColor(constants.ColorCyan))
}
