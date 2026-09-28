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
			buildActiveLapAndRerunSection(),
			buildDiagnosticsSection(),
			buildAutomationSection(),
		},
		FooterFlags: []termhelp.CommandEntry{
			{Command: "-h, --help", Description: "Show help for agy"},
			{Command: "-j, --json", Description: "Output structured JSON where supported (rp, lap, telegram, settings)"},
			{Command: "-f, --file <path>", Description: "Export JSON output directly to <path> (rp, lap)"},
		},
		Tips: []string{
			"Run 'gitmap agy <command> --help' for details on any subcommand.",
			"Use 'gitmap agy lap 24 --limit 10 --page 1 --wc 200' to inspect active projects & prompt trees.",
		},
	}
}

func buildProjectMgmtSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Project Management & Onboarding",
		Entries: []termhelp.CommandEntry{
			{Command: "open [path]", Description: "Open Antigravity Desktop IDE on path (default: .)"},
			{Command: "ls", Description: "List projects in status table grouped by root folder"},
			{Command: "add . / add <path>", Description: "Add current repo (.) or <path> to Antigravity workspace registry"},
			{Command: "add-read . / add-read <path> (ar)", Description: "Add current repo (.) or <path> and run Read Memory prompt"},
			{Command: "rm <id>", Description: "Remove project configuration (files on disk preserved)"},
			{Command: "update <id>", Description: "Update project updatedAt timestamp to now"},
			{Command: "scan [path]", Description: "Scan directory and register Antigravity projects"},
			{Command: "reconcile (recon)", Description: "Reconcile missing projects with active paths"},
			{Command: "remove-missing", Description: "Remove stale references to missing directories"},
			{Command: "pin-projects (pins)", Description: "Manage pinned priority projects"},
			{Command: "recreate-project (recreate, rp)", Description: "Purge cache/convs, re-register project, and start read-memory conv"},
			{Command: "group", Description: "Group and categorize projects by tag or folder", HasSubcommands: true},
		},
	}
}

func buildActiveLapAndRerunSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Running Projects, Last-Active Tree (LAP), Rerun-With-ID & Bot Settings",
		Entries: []termhelp.CommandEntry{
			{Command: "rp ls", Description: "List running projects with 24h sequence (#1), ID, alias, path & [convID]"},
			{Command: "rp prompts ls [--wc 200]", Description: "Tree view of running projects and active prompts (200 words default)"},
			{Command: "last-active-projects (lap) [N] [ls]", Description: "Projects active in last N hours (default 24h, --limit 10, --offset, --page, --wc 200)"},
			{Command: "rerun-with-id (rwi) <id|seq> <c> <p>", Description: "Inject prompt/file into project (<id|alias|seq|path>) and conversation (<c>|P1)"},
			{Command: "rerun-with-convid (rwc) <c|P1> <p>", Description: "Resolve project + conversation from <convid|P1> and inject prompt/file"},
			{Command: "account-switch (asw)", Description: "Fast-forward switch account on 98% used / 15% remaining credit threshold"},
			{Command: "telegram [setup|status|send|poll]", Description: "Two-way Telegram bot setup, notifications, and remote chat commands"},
			{Command: "email [setup|status|test]", Description: "Speed SMTP Email notification setup and verification"},
			{Command: "settings [--lap-hours 24]", Description: "Unified speed settings for LAP default hours (24), threshold (15%), Telegram & Email"},
		},
	}
}

func buildDiagnosticsSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Diagnostics & Optimization",
		Entries: []termhelp.CommandEntry{
			{Command: "active (running)", Description: "List conversations with active prompts currently running"},
			{Command: "queues (queue)", Description: "Inspect pending prompt queues across all workspaces"},
			{Command: "ping (check)", Description: "Ping Antigravity IDE and inspect environment health"},
			{Command: "find-duplicates (fdp)", Description: "Detect projects sharing identical filesystem paths"},
			{Command: "optimize-projects", Description: "Deduplicate and keep newest project per path"},
			{Command: "cache-clear (cc)", Description: "Clean runtime cache & prune convs (--keep 10, --pre)"},
			{Command: "ccko", Description: "Clean cache keeping only 1 conversation"},
			{Command: "cckf", Description: "Clean cache keeping top 5 conversations"},
			{Command: "undo", Description: "Revert last cache-clear from temp backup"},
			{Command: "remove-empty-convs", Description: "Purge projects with 0 conversation steps"},
			{Command: "status", Description: "Show Antigravity installation and profile status"},
			{Command: "stats", Description: "Display Antigravity workspace, running, and queue statistics"},
		},
	}
}
