package cmdpipeline

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderPipelineHelp displays the styled two-column Pipeline help menu.
func RenderPipelineHelp() {
	termhelp.RenderMenu(buildPipelineHelpMenu())
}

func buildPipelineHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "CI/CD Pipeline Diagnostics & Healing (gitmap pipeline)",
		UsageLines: []string{
			"gitmap pipeline [command] [flags]",
			"gitmap pd [commit|-N] [flags]  (runner details table shortcut)",
			"gitmap pe [commit|-N] [flags]  (error logs & clean status shortcut)",
			"gitmap pipeline-ai [cmd]       (AI telemetry & timeout waiting)",
		},
		Sections: []termhelp.HelpSection{
			buildPipelineTelemetrySection(),
			buildPipelineHealingSection(),
			buildPipelineDBSection(),
		},
		FooterFlags: buildPipelineFooterFlags(),
		Tips: []string{
			"Use shortcut 'gitmap pd' for runner target breakdown and step timings.",
			"Use shortcut 'gitmap pe' to inspect failing jobs and actionable snippets.",
			"Pass commit SHA or offset ('gitmap pe 7b1a2c', 'gitmap pe -1', 'gitmap pe HEAD~1') to target past runs.",
			"Run 'gitmap pe clear -y' to reset pipeline failure history for repo.",
		},
	}
}

func buildPipelineFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "-1, -2, -3, HEAD~N", Description: "Inspect pipeline logs for past commit by negative offset or SHA"},
		{Command: "-t, --timeline", Description: "Watch pipeline run until completion with dynamic ETA timeline"},
		{Command: "-f, --fix", Description: "Run internal CI/CD issue diagnostic and auto-repair scripts"},
		{Command: "-c, --check", Description: "Run internal diagnostic checks without modifying files"},
		{Command: "-y, --yes", Description: "Auto-confirm prompts non-interactively"},
		{Command: "--no-release", Description: "Prepare CI/CD fix prompt without triggering release"},
		{Command: "--json", Description: "Output pipeline telemetry in structured JSON format"},
		{Command: "--force, --no-cache", Description: "Bypass local SQLite DB cache and pull fresh from GitHub"},
		{Command: "-h, --help", Description: "Show this pipeline help menu"},
	}
}

func buildPipelineTelemetrySection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Live Telemetry & Diagnostics",
		Entries: []termhelp.CommandEntry{
			{Command: "status (st)", Description: "Check live CI/CD pipeline status, ETA, and pending PRs"},
			{Command: "details (pd)", Description: "Display runner target table, step elapsed timings & diagnostics"},
			{Command: "errors (pe)", Description: "Display failure logs, rerun ETA, and actionable snippets"},
			{Command: "waittime (eta)", Description: "Output remaining ETA seconds for active workflow run"},
			{Command: "history (hist)", Description: "Display recent commits pipeline execution tree"},
			{Command: "logs", Description: "Display consolidated workflow execution logs"},
		},
	}
}

func buildPipelineHealingSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "AI & Auto-Healing",
		Entries: []termhelp.CommandEntry{
			{Command: "fix", Description: "Feed pipeline errors to Antigravity IDE (alias: fix errors agy)"},
			{Command: "pipeline-ai status", Description: "Auto-delay (-t <sec>), stream live errors, and switch to fix"},
			{Command: "pipeline-ai errors", Description: "Extract failing workflow errors with AI remediation guidance"},
		},
	}
}

func buildPipelineDBSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Database & Reset",
		Entries: []termhelp.CommandEntry{
			{Command: "db", Description: "Inspect or manage isolated pipeline split SQLite database"},
			{Command: "clear (clean)", Description: "Clear logs, reports, and database for repository"},
			{Command: "clear-db", Description: "Clear recorded pipeline runs and error logs"},
		},
	}
}
