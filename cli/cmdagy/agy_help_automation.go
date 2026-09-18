package cmdagy

import "github.com/alimtvnetwork/gitmap-v28/cli/termhelp"

func buildAutomationSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title:   "Protocols & Automation",
		Entries: buildAutomationEntries(),
	}
}

func buildAutomationEntries() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "rerun last [N]", Description: "Replay last N prompts with prefix template (-p is-done)"},
		{Command: "list-prompts [N]", Description: "List prompts or diff in VS Code (--projects <N>)"},
		{Command: "read-all-projects-with-read-prompts (rprp)", Description: "Discover repos, sync, and broadcast Read Memory"},
		{Command: "all-projects-read-memory-prompt", Description: "Broadcast Read Memory prompt to active projects"},
		{Command: "fix-pipeline (fp)", Description: "Diagnose and fix CI/CD pipeline issues", HasSubcommands: true},
		{Command: "prompt <text>", Description: "Send prompt to active project session"},
		{Command: "sync", Description: "Sync Antigravity settings and project states"},
		{Command: "export / import", Description: "Export or import Antigravity projects JSON", HasSubcommands: true},
		{Command: "plugins", Description: "Inspect, enable, or disable Antigravity plugins", HasSubcommands: true},
	}
}
