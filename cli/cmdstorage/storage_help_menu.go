package cmdstorage

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderStorageHelp displays the styled two-column Storage help menu.
func RenderStorageHelp() {
	termout.RenderMenu(buildStorageHelpMenu())
}

func buildStorageHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "Storage & Database Manager (gitmap storage)",
		UsageLines: []string{
			"gitmap storage [command] [flags]",
			"gitmap disk [command] [flags]",
			"gitmap df [command] [flags]",
		},
		Sections: []termout.HelpSection{
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

func buildStorageFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "-a, --all", Description: "Include hidden, system, and virtual drives"},
		{Command: "-j, --json", Description: "Output storage statistics in structured JSON format"},
		{Command: "-f, --force", Description: "Bypass interactive confirmation for restore/clean"},
		{Command: "-v, --verbose", Description: "Print detailed information during operations"},
		{Command: "--dry-run", Description: "Simulate cleanup without deleting any files"},
		{Command: "-h, --help", Description: "Show this storage help menu"},
	}
}

func buildStorageInspectionSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Drive & Space Inspection",
		Entries: []termout.CommandEntry{
			{Command: "status (st, info)", Description: "Inspect drive capacity, filesystem metrics, and SQLite summary"},
			{Command: "space [ls]", Description: "Inspect repository database sizes and space consumption", HasSubcommands: true},
		},
	}
}

func buildStorageDBSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Database Inventory",
		Entries: []termout.CommandEntry{
			{Command: "ls (list, db)", Description: "List all system, split, and repository SQLite databases"},
			{Command: "restore-db [name]", Description: "Restore or auto-heal database from backup snapshot"},
		},
	}
}

func buildStorageMaintenanceSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Maintenance & Cleanup",
		Entries: []termout.CommandEntry{
			{Command: "clean (clear, prune)", Description: "Clean pipeline logs, temp files, and vacuum databases"},
			{Command: "reset-errors", Description: "Purge pipeline error caches, failure telemetry, and reports"},
		},
	}
}
