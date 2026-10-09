package cmdschedule

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

func buildScheduleTaskSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Task Scheduling & Creation",
		Entries: []termout.CommandEntry{
			{Command: "add <name> [cmds...]", Description: "Create scheduled task with isolated split DB"},
			{Command: "list (ls)", Description: "List all scheduled tasks, intervals, run counts & status"},
			{Command: "status [name|*]", Description: "View detailed status and metadata of schedule(s)"},
			{Command: "run (exec) <name>", Description: "Trigger on-demand immediate execution of task"},
			{Command: "test <name>", Description: "Test task execution with simulated runs (--times N)"},
		},
	}
}

func buildScheduleLifecycleSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "State & Lifecycle",
		Entries: []termout.CommandEntry{
			{Command: "enable <name>", Description: "Enable a disabled scheduled task"},
			{Command: "disable <name>", Description: "Disable a scheduled task (preserves logs)"},
			{Command: "rm (delete) <name>", Description: "Remove scheduled task and purge split database"},
		},
	}
}

func buildScheduleLogsSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Execution Logs & Reset",
		Entries: []termout.CommandEntry{
			{Command: "logs (history) <name>", Description: "View execution history and logs from split DB"},
			{Command: "reset <name>", Description: "Purge execution logs for a specific schedule"},
			{Command: "reset-all", Description: "Reset execution logs across all schedule split DBs"},
		},
	}
}

func buildScheduleExportImportSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Import & Export",
		Entries: []termout.CommandEntry{
			{Command: "export [name]", Description: "Export schedule(s) to JSON, YAML, SQLite, or ZIP", HasSubcommands: true},
			{Command: "export-all", Description: "Export all schedules with exclusion filter"},
			{Command: "import <file>", Description: "Import schedule(s) from file or archive", HasSubcommands: true},
			{Command: "import-all -f <file>", Description: "Batch import schedules into root & split DBs"},
		},
	}
}

func buildSchedulePowerSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "System Power Operations",
		Entries: []termout.CommandEntry{
			{Command: "shutdown [<dur>]", Description: "Schedule OS shutdown (e.g. 1h30m, status, cancel)"},
			{Command: "restart [<dur>]", Description: "Schedule OS restart (e.g. 1h30m, status, cancel)"},
		},
	}
}
