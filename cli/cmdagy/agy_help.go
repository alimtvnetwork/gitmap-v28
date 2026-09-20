// Package cmdagy — agy_help.go renders colorful, aligned help for the agy command suite.
package cmdagy

import (
	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/helptext"
	"github.com/alimtvnetwork/gitmap-v28/cli/render"
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

func renderAgyHelp(cmd *cobra.Command, args []string) {
	if handleAgySubcommandHelp(cmd) {
		return
	}
	termhelp.RenderMenu(buildAgyHelpMenu())
}

func handleAgySubcommandHelp(cmd *cobra.Command) bool {
	if cmd == nil || cmd == AgyCmd {
		return false
	}
	if renderSubcommandHelp(cmd) {
		return true
	}
	_ = cmd.Usage()

	return true
}

func renderSubcommandHelp(cmd *cobra.Command) bool {
	helpName := cmd.Name()
	if _, err := helptext.ReadRaw(helpName); err == nil {
		helptext.PrintWithMode(helpName, render.PrettyAuto)

		return true
	}

	return false
}

func buildAgyHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "Antigravity CLI Management (gitmap agy)",
		UsageLines: []string{
			"gitmap agy [command] [flags]",
			"agy [command] [flags]",
		},
		Sections: []termhelp.HelpSection{
			buildProjectMgmtSection(),
			buildDiagnosticsSection(),
			buildAutomationSection(),
		},
		FooterFlags: []termhelp.CommandEntry{
			{Command: "-h, --help", Description: "Show help for agy"},
		},
		Tips: []string{
			"Run 'gitmap agy <command> --help' for details on any subcommand.",
		},
	}
}

func buildProjectMgmtSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Project Management",
		Entries: []termhelp.CommandEntry{
			{Command: "open [path]", Description: "Open Antigravity Desktop IDE on path (default: .)"},
			{Command: "ls", Description: "List projects in status table grouped by root folder"},
			{Command: "add <id> <name>", Description: "Add an Antigravity project configuration"},
			{Command: "rm <id>", Description: "Remove project configuration (files on disk preserved)"},
			{Command: "update <id>", Description: "Update project updatedAt timestamp to now"},
			{Command: "scan [path]", Description: "Scan directory and register Antigravity projects"},
			{Command: "reconcile (recon)", Description: "Reconcile missing projects with active paths"},
			{Command: "remove-missing", Description: "Remove stale references to missing directories"},
			{Command: "pin-projects (pins)", Description: "Manage pinned priority projects"},
			{Command: "group", Description: "Group and categorize projects by tag or folder", HasSubcommands: true},
		},
	}
}

func buildDiagnosticsSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Diagnostics & Optimization",
		Entries: []termhelp.CommandEntry{
			{Command: "ping (check)", Description: "Ping Antigravity IDE and inspect environment health"},
			{Command: "find-duplicates (fdp)", Description: "Detect projects sharing identical filesystem paths"},
			{Command: "optimize-projects", Description: "Deduplicate and keep newest project per path"},
			{Command: "clean-cache (cc)", Description: "Clean runtime cache and orphan project artifacts"},
			{Command: "remove-empty-convs", Description: "Purge projects with 0 conversation steps"},
			{Command: "status", Description: "Show Antigravity installation and profile status"},
			{Command: "stats", Description: "Display Antigravity workspace usage statistics"},
		},
	}
}
