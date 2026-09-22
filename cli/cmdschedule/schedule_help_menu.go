package cmdschedule

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderScheduleHelp displays the styled two-column Schedule help menu.
func RenderScheduleHelp() {
	termhelp.RenderMenu(buildScheduleHelpMenu())
}

func buildScheduleHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "Task Scheduler & Power Manager (gitmap schedule)",
		UsageLines: []string{
			"gitmap schedule [command] [args]",
			"gitmap sched [command] [args]",
			"gitmap schedule add <name> [cmds...] [--every <interval>]",
		},
		Sections: []termhelp.HelpSection{
			buildScheduleTaskSection(),
			buildScheduleLifecycleSection(),
			buildScheduleLogsSection(),
			buildScheduleExportImportSection(),
			buildSchedulePowerSection(),
		},
		FooterFlags: buildScheduleFooterFlags(),
		Tips: []string{
			"Run 'gitmap schedule add backup --every 1d' to schedule a daily task.",
			"Use 'gitmap schedule logs <name> --limit 10' to inspect execution records.",
			"Run 'gitmap schedule shutdown 1h30m' to schedule an OS shutdown.",
		},
	}
}

func buildScheduleFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "--every, -i <interval>", Description: "Execution interval (e.g. 1d, 2h, 30m, 15s)"},
		{Command: "--macro, -m <name>", Description: "Link scheduled execution to a saved macro"},
		{Command: "--startup", Description: "Automatically register task in OS startup autorun"},
		{Command: "--delay, -d <dur>", Description: "Initial sleep/delay before execution starts"},
		{Command: "-f, --file <path>", Description: "File path for export/import or saved reports"},
		{Command: "--json", Description: "Output schedule list, status, or logs as JSON"},
		{Command: "--yaml", Description: "Output schedule list, status, or logs as YAML"},
		{Command: "-h, --help", Description: "Show this scheduler help menu"},
	}
}
