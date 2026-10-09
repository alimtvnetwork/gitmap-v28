package cmdmacro

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderMacroHelp displays the styled two-column Macro help menu.
func RenderMacroHelp() {
	termout.RenderMenu(buildMacroHelpMenu())
}

func buildMacroHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "Macro Automation & Replay Engine (gitmap macro)",
		UsageLines: []string{
			"gitmap macro [command] [args]",
			"gitmap m [command] [args]",
			"gitmap macro add <name> <steps...>",
			"gitmap macro run <name> [flags]",
		},
		Sections: []termout.HelpSection{
			buildMacroCreationSection(),
			buildMacroExecutionSection(),
			buildMacroMgmtSection(),
			buildMacroExportImportSection(),
			buildMacroSystemSection(),
		},
		FooterFlags: buildMacroFooterFlags(),
		Tips: []string{
			"Run 'gitmap macro record <name>' to record a live shell workflow.",
			"Use 'gitmap macro run-until-succeed <cmd> --sleep 5s' to retry until exit 0.",
		},
	}
}

func buildMacroFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "--dry-run", Description: "Preview planned steps without executing"},
		{Command: "--json", Description: "Output macro data in structured JSON format"},
		{Command: "--yaml", Description: "Output macro data in structured YAML format"},
		{Command: "-f, --file <path>", Description: "Specify custom input or output file path"},
		{Command: "-h, --help", Description: "Show this macro help menu"},
	}
}

func buildMacroCreationSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Creation & Editing",
		Entries: []termout.CommandEntry{
			{Command: "add <name> <steps...>", Description: "Create a new macro directly from arguments or interactively"},
			{Command: "edit <name>", Description: "Interactively edit, insert, or reorder macro steps"},
			{Command: "record (rec) <name>", Description: "Record interactive shell commands as a macro"},
		},
	}
}

func buildMacroExecutionSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Replay & Execution",
		Entries: []termout.CommandEntry{
			{Command: "run (exec) <name>", Description: "Replay a recorded macro sequence"},
			{Command: "run-until-succeed <name>", Description: "Execute repeatedly until success (with backoff & AI)"},
		},
	}
}

func buildMacroMgmtSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Inspection & Management",
		Entries: []termout.CommandEntry{
			{Command: "list (ls)", Description: "List all saved macros with step counts and tags"},
			{Command: "show <name>", Description: "Inspect steps and details of a recorded macro"},
			{Command: "rm (delete) <name>", Description: "Delete a saved macro from SQLite storage"},
		},
	}
}

func buildMacroExportImportSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Import & Export",
		Entries: []termout.CommandEntry{
			{Command: "export [name|all]", Description: "Export macro(s) to JSON, YAML, SQLite DB, or ZIP", HasSubcommands: true},
			{Command: "export-all", Description: "Export all stored macros to file or stdout"},
			{Command: "import <file>", Description: "Import macro(s) safely with format auto-detection", HasSubcommands: true},
			{Command: "import-all <file>", Description: "Batch import all macros from archive or database"},
		},
	}
}

func buildMacroSystemSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "System & Fleet Integration",
		Entries: []termout.CommandEntry{
			{Command: "startup <sub> [name]", Description: "Manage macro execution on OS login / reboot", HasSubcommands: true},
			{Command: "schedule <sub> [name]", Description: "Manage recurring scheduled execution of macros", HasSubcommands: true},
			{Command: "sync [target]", Description: "Synchronize local macros across remote SSH fleet nodes"},
		},
	}
}
