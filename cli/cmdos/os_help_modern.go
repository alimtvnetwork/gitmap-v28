// Package cmdos — os_help_modern.go renders the modernized boxed enterprise help menu for gitmap os.
package cmdos

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderModernOSHelp prints the boxed enterprise help menu for gitmap os.
func RenderModernOSHelp() {
	termhelp.RenderMenu(buildModernOSHelpMenu())
	printCrossPlatformOSExamples()
}

func buildModernOSHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "Operating System, Machine Identity & Alias Management (gitmap os)",
		UsageLines: []string{
			"gitmap os <subcommand> [flags]",
			"gitmap machine [ls|set|change|revert|help] [--ssh] [-y]",
			"gitmap alias [ls|set|change|revert|help] [--ssh] [-y]",
		},
		Sections: []termhelp.HelpSection{
			buildOSIdentityAndAliasSection(),
			buildOSCoreUpdateAndInfoSection(),
			buildOSSystemAndNetworkSection(),
		},
		FooterFlags: buildOSHelpFooterFlags(),
		Tips: []string{
			"Both 'gitmap machine' and 'gitmap alias' work at top-level or inside 'gitmap os'.",
			"Pass '--ssh' to inspect or update machine aliases across all joined SSH fleet nodes.",
		},
	}
}

func buildOSIdentityAndAliasSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Machine Identity & Alias Management (Local & --ssh Fleet)",
		Entries: []termhelp.CommandEntry{
			{Command: "os machine ls [--ssh]", Description: "Show Machine IP, Alias (auto-IP fallback), Hostname & OS (alias: gitmap machine ls)"},
			{Command: "os machine set <name> [-y] [--ssh]", Description: "Set OS Hostname & GitMap Machine Name (sample: dev-win-01, ubuntu-node-02)"},
			{Command: "os machine change <name> [-y]", Description: "Change machine name and snapshot previous value for instant revert"},
			{Command: "os machine revert [--ssh]", Description: "Restore previous machine name and OS hostname from snapshot"},
			{Command: "os alias ls [--ssh]", Description: "Display Machine Alias and IP across local or joined SSH fleet (alias: gitmap alias ls)"},
			{Command: "os alias set <alias> [-y] [--ssh]", Description: "Set network alias so peer machines recognize this node"},
			{Command: "os alias change <alias> [-y]", Description: "Change machine network alias and store rollback snapshot"},
			{Command: "os alias revert [--ssh]", Description: "Revert machine alias back to previous saved value"},
		},
	}
}

func buildOSCoreUpdateAndInfoSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "OS Inspection, Updates & Package Maintenance",
		Entries: []termhelp.CommandEntry{
			{Command: "os tui (menu)", Description: "Interactive Bubbletea dashboard for tweaks, DNS, auto-login, and clean"},
			{Command: "os info (gitmap os-info)", Description: "Inspect OS distribution, kernel, architecture, CPU, RAM, and hostname"},
			{Command: "os update (gitmap os-update)", Description: "Update OS package repositories and security patches (Windows/macOS/Ubuntu)"},
			{Command: "os upgrade / full-upgrade", Description: "Execute full OS package and distribution upgrade"},
			{Command: "os fix-mirrors", Description: "Switch regional APT/package mirrors to canonical high-speed mirrors"},
			{Command: "os status (st)", Description: "Display OS environment health and symlink diagnostics"},
		},
	}
}

func buildOSSystemAndNetworkSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Network, Storage, Desktop & Cache Hygiene",
		Entries: []termhelp.CommandEntry{
			{Command: "os ip [show|set|change|revert]", Description: "Inspect or configure static/DHCP IPv4 with ping rollback protection"},
			{Command: "os dns [show|set|bench|revert]", Description: "Inspect, benchmark, or configure network DNS nameservers"},
			{Command: "os ai-clean / dev-clean / clean", Description: "Purge AI brain caches, compiler caches, and ephemeral temp directories"},
			{Command: "os storage (disk)", Description: "Inspect disk drive capacities, partitions, and storage utilization"},
			{Command: "os display / theme / tweak", Description: "Configure display resolution, dark/light theme, and power schemes"},
			{Command: "os dock [bottom|left|right|top] [--node]", Description: "Inspect or configure desktop dock / taskbar position (aliases: panel, start-menu)"},
			{Command: "os user / group / autologin", Description: "Manage OS users, groups, and automatic login credentials"},
			{Command: "os change-password [user] [pass]", Description: "Change OS user password cross-platform (Windows / Linux / macOS)"},
		},
	}
}

func buildOSHelpFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "-s, --ssh", Description: "Execute machine/alias operations across all joined SSH fleet nodes"},
		{Command: "-y, --yes", Description: "Apply changes immediately without interactive confirmation"},
		{Command: "-j, --json", Description: "Output structured JSON for automation and scripting"},
		{Command: "-h, --help", Description: "Display this modernized OS help menu"},
	}
}

func printCrossPlatformOSExamples() {
	fmt.Printf("  %sReal-World Examples (Windows, macOS, Ubuntu/Linux):%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    %sgitmap os info%s                              # Inspect OS platform, IP, and hardware\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    %sgitmap os update%s                            # Update OS packages (winget / brew / apt)\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    %sgitmap os change-password Administrator MyPass%s # Set user password directly on Windows/Linux/macOS\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    %sgitmap os dock bottom%s                       # Set dock/panel to bottom (Linux / Win / macOS)\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    %sgitmap os dock left --node u1%s               # Delegate dock position over SSH fleet\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    %sgitmap machine set dev-win-01 -y%s            # Set machine name on Windows/macOS/Ubuntu\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    %sgitmap alias set build-ubuntu-02 -y%s         # Set machine network alias for fleet discovery\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    %sgitmap machine ls --ssh%s                     # List all SSH fleet machines with IPs & aliases\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    %sgitmap machine revert%s                       # Revert machine name/alias to previous value\n\n", constants.ColorCyan, constants.ColorReset)
}
