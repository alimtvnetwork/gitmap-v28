package cmdllm

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderLlmHelp displays the styled two-column LLM / Train help menu.
func RenderLlmHelp() {
	termout.RenderMenu(buildLlmHelpMenu())
}

func buildLlmHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "Autonomous LLM Training & AI Pairing Suite (gitmap llm)",
		UsageLines: []string{
			"gitmap llm [flags]",
			"gitmap llm train [--loop] [--self-loop N] [--url] [flags]",
			"gitmap train [--loop] [--self-loop N]",
			"gitmap llm-train [--loop]",
		},
		Sections: []termout.HelpSection{
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

func buildLlmFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "--loop", Description: "Execute autonomous 5-phase AI self-looping execution cycle"},
		{Command: "--self-loop <N>", Description: "Run N consecutive iterations of the AI self-loop"},
		{Command: "--url", Description: "Output raw public URL to llm.md instruction specification"},
		{Command: "-j, --json", Description: "Output structured machine-readable command specifications"},
		{Command: "--skill-path <p>", Description: "Target path for generated Antigravity skill (SKILL.md)"},
		{Command: "--text-only", Description: "Output curriculum markdown without writing to filesystem"},
		{Command: "-h, --help", Description: "Show this LLM help menu"},
	}
}

func buildLlmCommandsSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Curriculum & Directives",
		Entries: []termout.CommandEntry{
			{Command: "train / chain", Description: "Execute chained curriculum and generate Antigravity skill"},
			{Command: "llm-train", Description: "Direct alias for 'gitmap llm train'"},
			{Command: "llm-docs (ld)", Description: "Generate complete LLM Markdown reference document"},
			{Command: "pe history-ai [N]", Description: "Dumps historical failure dossiers to prevent AI repeating mistakes"},
		},
	}
}

func buildLlmLifecycleSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "5-Phase AI Agent Lifecycle",
		Entries: []termout.CommandEntry{
			{Command: "1. Discovery", Description: "find-files, find-files-any, search, list-files"},
			{Command: "2. Refactoring", Description: "replace, replace-regex, targeted zero-nesting edits"},
			{Command: "3. Verification", Description: "python linters, go test, smart test runner"},
			{Command: "4. Semantic Commit", Description: "commit-push-feature (cpf), commit-push-bug (cpb)"},
			{Command: "5. Telemetry & Heal", Description: "pipeline-ai status --json, dynamic ETA sleep, pe history-ai"},
		},
	}
}

func buildLlmLoopSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Autonomous Loop & Self-Training",
		Entries: []termout.CommandEntry{
			{Command: "--loop", Description: "Runs iterative self-training simulation through all 5 phases"},
			{Command: "--self-loop N", Description: "Specifies loop repetition count (1 to N iterations)"},
		},
	}
}
