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
		{Command: "prompt-project (p)", Description: "Target project by prefix with -name, -txt, --pf/--sf (prefix + 2 newlines)"},
		{Command: "prompt [-n name] [-t txt]", Description: "Dispatch to current repo (auto-adds project, queues read-all first)"},
		{Command: "prompt-with-name (pwn)", Description: "Dispatch named template (same as prompt -name; queues read-all first)"},
		{Command: "prompt-txt (pt)", Description: "Dispatch direct text (defaults to read-all prefix with 2 newlines)"},
		{Command: "prompt ls (list prompts)", Description: "List all available prompt templates in formatted table"},
		{Command: "rerun last [N]", Description: "Replay last N prompts with prefix template (-p is-done)"},
		{Command: "list-prompts [N]", Description: "List prompts or diff in VS Code (--projects <N>)"},
		{Command: "read-all-projects-with-read-prompts (rprp)", Description: "Discover repos, sync, and broadcast Read Memory"},
		{Command: "all-projects-read-memory-prompt", Description: "Broadcast Read Memory prompt to active projects"},
		{Command: "fix-pipeline (fp)", Description: "Diagnose and fix CI/CD pipeline issues", HasSubcommands: true},
		{Command: "sync", Description: "Sync Antigravity settings and project states"},
		{Command: "export / import", Description: "Export or import Antigravity projects JSON", HasSubcommands: true},
		{Command: "plugins", Description: "Inspect, enable, or disable Antigravity plugins", HasSubcommands: true},
	}
}
