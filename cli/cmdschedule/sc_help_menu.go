package cmdschedule

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderSCHelp displays the styled two-column SC help menu.
func RenderSCHelp() {
	termout.RenderMenu(buildSCHelpMenu())
}

func buildSCHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "Distributed Fleet Broadcast Engine (gitmap sc)",
		UsageLines: []string{
			"gitmap sc <subcommand> [args] [flags]",
			"gitmap servers-clients <subcommand> [args] [flags]",
			"gitmap sc bash \"<command>\" [--except <list>]",
		},
		Sections: []termout.HelpSection{
			buildSCBroadcastSection(),
			buildSCInventorySection(),
			buildSCGitSection(),
			buildSCPowerSection(),
		},
		FooterFlags: buildSCFooterFlags(),
		Tips: []string{
			"Use 'gitmap sc bash \"uname -a\" --except 2' to execute across fleet.",
			"Run 'gitmap sc install gitmap' to deploy GitMap remotely across nodes.",
		},
	}
}

func buildSCFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "--except <list>", Description: "Exclude nodes by ID, IP, or trailing octet"},
		{Command: "--ip <list>", Description: "Target specific node IP addresses"},
		{Command: "--id <list>", Description: "Target specific node Display IDs"},
		{Command: "--json", Description: "Output machine list or results in JSON format"},
		{Command: "-y, --yes", Description: "Bypass preflight confirmation prompt"},
		{Command: "--dry-run", Description: "Preview execution plan without running commands"},
		{Command: "-h, --help", Description: "Show this servers-clients help menu"},
	}
}

func buildSCBroadcastSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Shell Command Broadcast",
		Entries: []termout.CommandEntry{
			{Command: "bash <cmd>", Description: "Execute a Bash command across Linux/macOS nodes"},
			{Command: "sh (shell) <cmd>", Description: "Execute a POSIX shell command across all nodes"},
			{Command: "ps <cmd>", Description: "Execute a PowerShell command on Windows cluster nodes"},
			{Command: "cmd <cmd>", Description: "Execute a Windows Command Prompt command on nodes"},
		},
	}
}

func buildSCInventorySection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Machine Enrollment & Inventory",
		Entries: []termout.CommandEntry{
			{Command: "nodes (ls)", Description: "List all registered machines joined to the cluster"},
			{Command: "join <t> [alias]", Description: "Join and enroll a machine into the cluster registry"},
			{Command: "remove (rm) <t>", Description: "Remove a machine from the cluster registry"},
			{Command: "ping (status)", Description: "Check reachability and round-trip ping latency"},
		},
	}
}

func buildSCGitSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Fleet Synchronization & Git",
		Entries: []termout.CommandEntry{
			{Command: "pull --all", Description: "Run git pull --all across all repositories on all nodes"},
			{Command: "push --all", Description: "Run git push --all across all repositories on all nodes"},
			{Command: "status --all", Description: "Show combined dirty/clean status across all nodes"},
			{Command: "install <pkgs>", Description: "Install packages or GitMap binary across all nodes"},
		},
	}
}

func buildSCPowerSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Power & Architecture",
		Entries: []termout.CommandEntry{
			{Command: "restart [dur]", Description: "Trigger or schedule system restart across nodes"},
			{Command: "shutdown [dur]", Description: "Trigger or schedule system shutdown across nodes"},
			{Command: "compare (matrix)", Description: "Display architecture matrix: SSH vs Cluster vs SC"},
		},
	}
}
