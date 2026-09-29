package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderLlmHelp displays the styled two-column LLM / Train help menu.
func RenderLlmHelp() {
	termhelp.RenderMenu(buildLlmHelpMenu())
}

func buildLlmHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "Autonomous LLM Training & AI Pairing Suite (gitmap llm)",
		UsageLines: []string{
			"gitmap llm [flags]",
			"gitmap llm train [--loop] [--self-loop N] [--url] [flags]",
			"gitmap train [--loop] [--self-loop N]",
			"gitmap llm-train [--loop]",
		},
		Sections: []termhelp.HelpSection{
			buildLlmCommandsSection(),
			buildLlmLifecycleSection(),
			buildLlmLoopSection(),
		},
		FooterFlags: buildLlmFooterFlags(),
		Tips: []string{
			"Run 'gitmap train --loop' for continuous autonomous 5-phase execution.",
			"Use 'gitmap llm --url' to download raw public instruction spec for AI models.",
		},
	}
}

func buildLlmFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "--loop", Description: "Execute autonomous 5-phase AI self-looping execution cycle"},
		{Command: "--self-loop <N>", Description: "Run N consecutive iterations of the AI self-loop"},
		{Command: "--url", Description: "Output raw public URL to llm.md instruction specification"},
		{Command: "-j, --json", Description: "Output structured machine-readable command specifications"},
		{Command: "--skill-path <p>", Description: "Target path for generated Antigravity skill (SKILL.md)"},
		{Command: "--text-only", Description: "Output curriculum markdown without writing to filesystem"},
		{Command: "-h, --help", Description: "Show this LLM help menu"},
	}
}

func buildLlmCommandsSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Curriculum & Directives",
		Entries: []termhelp.CommandEntry{
			{Command: "train / chain", Description: "Execute chained curriculum and generate Antigravity skill"},
			{Command: "llm-train", Description: "Direct alias for 'gitmap llm train'"},
			{Command: "llm-docs (ld)", Description: "Generate complete LLM Markdown reference document"},
		},
	}
}

func buildLlmLifecycleSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "5-Phase AI Agent Lifecycle",
		Entries: []termhelp.CommandEntry{
			{Command: "1. Discovery", Description: "find-files, find-files-any, search, list-files"},
			{Command: "2. Refactoring", Description: "replace, replace-regex, targeted zero-nesting edits"},
			{Command: "3. Verification", Description: "python linters, go test, smart test runner"},
			{Command: "4. Semantic Commit", Description: "commit-push-feature (cpf), commit-push-bug (cpb)"},
			{Command: "5. Telemetry & Heal", Description: "pipeline-ai status --json, dynamic ETA sleep, RCA logs"},
		},
	}
}

func buildLlmLoopSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Autonomous Loop & Self-Training",
		Entries: []termhelp.CommandEntry{
			{Command: "--loop", Description: "Runs iterative self-training simulation through all 5 phases"},
			{Command: "--self-loop N", Description: "Specifies loop repetition count (1 to N iterations)"},
		},
	}
}
