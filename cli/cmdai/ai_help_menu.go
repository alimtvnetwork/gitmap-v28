package cmdai

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderAiHelp renders the styled two-column help menu for gitmap ai.
func RenderAiHelp() {
	termhelp.RenderMenu(buildAiHelpMenu())
}

func buildAiHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "Native AI Automation & Scripting (gitmap ai)",
		UsageLines: []string{
			"gitmap ai <subcommand|script-token> [args...] [flags]",
			"gitmap ai ls [--copy] [--frequent]",
			"gitmap ai pwsh <powershell-command> [--ai]",
			"gitmap ai run <script-token> [args...]",
			"gitmap ai fix [target]",
		},
		Sections: []termhelp.HelpSection{
			buildAiCommandsSection(),
			buildAiHistorySection(),
			buildAiExecutionSection(),
		},
		FooterFlags: buildAiFooterFlags(),
		Tips: []string{
			"Run 'gitmap ai ls --copy' to view and copy top frequent commands to clipboard.",
			"Add '--ai' to any PowerShell or automation execution to track it in ai split DB.",
			"Use 'gitmap ai run <number|slug>' to execute catalog automation scripts directly.",
		},
	}
}

func buildAiCommandsSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Catalog & Discovery",
		Entries: []termhelp.CommandEntry{
			{Command: "ls, list, catalog", Description: "List registered scripts and top frequent AI commands"},
			{Command: "history, hist, freq", Description: "Inspect top frequently executed AI commands"},
			{Command: "run, exec <token>", Description: "Execute an AI automation script with live streaming"},
			{Command: "create, scaffold", Description: "Scaffold a new standardized repository AI script"},
		},
	}
}

func buildAiHistorySection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Split DB History & PowerShell",
		Entries: []termhelp.CommandEntry{
			{Command: "pwsh, ps <command>", Description: "Execute PowerShell command and record execution history"},
			{Command: "--copy", Description: "Copy frequent AI commands to clipboard for LLM reuse"},
			{Command: "--ai", Description: "Force recording command into ~/.gitmap/ai/instructions.db"},
		},
	}
}

func buildAiExecutionSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Auto-Remediation",
		Entries: []termhelp.CommandEntry{
			{Command: "fix [target]", Description: "Run standard autofix targets (guidelines, newlines, etc.)"},
			{Command: "fix all", Description: "Execute full suite of repository remediation targets"},
		},
	}
}

func buildAiFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "--copy", Description: "Copy frequent AI commands to system clipboard"},
		{Command: "--frequent", Description: "Filter list view to show only frequent AI commands"},
		{Command: "-c, --category <cat>", Description: "Filter script catalog by taxonomy category"},
		{Command: "-s, --search <query>", Description: "Search scripts by number, slug, filename, or alias"},
		{Command: "-h, --help", Description: "Show this AI help menu"},
	}
}
