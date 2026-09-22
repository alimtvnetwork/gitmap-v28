package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderStorageHelp displays the styled two-column Storage help menu.
func RenderStorageHelp() {
	termhelp.RenderMenu(buildStorageHelpMenu())
}

func buildStorageHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "Storage & Database Manager (gitmap storage)",
		UsageLines: []string{
			"gitmap storage [command] [flags]",
			"gitmap disk [command] [flags]",
			"gitmap df [command] [flags]",
		},
		Sections: []termhelp.HelpSection{
			buildStorageInspectionSection(),
			buildStorageDBSection(),
			buildStorageMaintenanceSection(),
		},
		FooterFlags: buildStorageFooterFlags(),
		Tips: []string{
			"Run 'gitmap storage ls' to inspect all split databases in .gitmap/data/.",
			"Use 'gitmap storage clean --dry-run' to simulate cleanup before deleting.",
		},
	}
}

func buildStorageFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "-a, --all", Description: "Include hidden, system, and virtual drives"},
		{Command: "-j, --json", Description: "Output storage statistics in structured JSON format"},
		{Command: "-f, --force", Description: "Bypass interactive confirmation for restore/clean"},
		{Command: "-v, --verbose", Description: "Print detailed information during operations"},
		{Command: "--dry-run", Description: "Simulate cleanup without deleting any files"},
		{Command: "-h, --help", Description: "Show this storage help menu"},
	}
}

func buildStorageInspectionSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Drive & Space Inspection",
		Entries: []termhelp.CommandEntry{
			{Command: "status (st, info)", Description: "Inspect drive capacity, filesystem metrics, and SQLite summary"},
			{Command: "space [ls]", Description: "Inspect repository database sizes and space consumption", HasSubcommands: true},
		},
	}
}

func buildStorageDBSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Database Inventory",
		Entries: []termhelp.CommandEntry{
			{Command: "ls (list, db)", Description: "List all system, split, and repository SQLite databases"},
			{Command: "restore-db [name]", Description: "Restore or auto-heal database from backup snapshot"},
		},
	}
}

func buildStorageMaintenanceSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Maintenance & Cleanup",
		Entries: []termhelp.CommandEntry{
			{Command: "clean (clear, prune)", Description: "Clean pipeline logs, temp files, and vacuum databases"},
			{Command: "reset-errors", Description: "Purge pipeline error caches, failure telemetry, and reports"},
		},
	}
}
