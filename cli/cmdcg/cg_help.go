package cmdcg

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderCGHelp displays the styled two-column Coding Guidelines help menu.
func RenderCGHelp() {
	termout.RenderMenu(buildCGHelpMenu())
}

func buildCGHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title:       "Coding Guidelines, Prompts & Quality Standards (gitmap cg)",
		BannerColor: termout.RoleSGRBold(termout.RolePrimary),
		UsageLines:  buildCGUsageLines(),
		Sections: []termout.HelpSection{
			buildCGOpsSection(),
			buildCGPromptSection(),
			buildCGVersioningSection(),
		},
		FooterFlags: buildCGFooterFlags(),
		Tips:        buildCGTips(),
	}
}

func buildCGUsageLines() []string {
	return []string{
		"gitmap cg <command> [targets...] [flags]",
		"gitmap cg sync [targets...] [--all] [--dry-run]",
		"gitmap cg prompt sync [targets...] [--dry-run]",
		"gitmap cg skill sync [targets...] [--dry-run]",
		"gitmap cg execute <v6|v7> [--raw]",
		"gitmap cg status [--all]",
		"gitmap cg audit [targets...]",
		"gitmap cg doctor [targets...]",
		"gitmap cg install-version-json --version 1.0.0",
	}
}

func buildCGOpsSection() termout.HelpSection {
	return termout.HelpSection{
		Title:      "Coding Guidelines Operations",
		Color:      termout.RoleSGRBold(termout.RoleOK),
		BadgeText:  "OPS",
		BadgeColor: termout.RoleSGR(termout.RoleFGSubtle),
		Entries: []termout.CommandEntry{
			{Command: "status, stat, ls", Description: "View tabulated Coding Guidelines install status & version across repos"},
			{Command: "audit, check", Description: "Deep audit of coding standards compliance (specs, rules, memory)"},
			{Command: "doctor, doc", Description: "Health check and configuration repair diagnostics for coding guidelines"},
			{Command: "install", Description: "Install official Coding Guidelines standards into target repositories"},
			{Command: "update", Description: "In-place update of Coding Guidelines for installed repositories"},
			{Command: "diff, drift", Description: "Compare local guidelines version with latest upstream canonical standards"},
		},
	}
}

func buildCGPromptSection() termout.HelpSection {
	return termout.HelpSection{
		Title:      "Prompt Architect & AI Standards",
		Color:      termout.RoleSGRBold(termout.RoleAccent),
		BadgeText:  "AI",
		BadgeColor: termout.RoleSGR(termout.RoleAccent),
		Entries: []termout.CommandEntry{
			{Command: "sync, sy", Description: "Unified 3-way synchronization across specs, prompts, skills, and scripts with human conflict protection"},
			{Command: "prompt sync", Description: "Synchronize canonical 01-prompts/ templates using 3-way hash verification"},
			{Command: "skill sync", Description: "Synchronize .agents/skills/ and .cursor/skills/ with human edit preservation"},
			{Command: "execute <v6|v7>", Description: "Stream master V6 or V7 autonomous AI execution directives to terminal for model learning"},
			{Command: "install-prompts", Description: "Install Prompt Architect v2 suite and canonical prompt library"},
			{Command: "update-prompts", Description: "In-place update of Prompt Architect templates across target repos"},
			{Command: "prompts-status", Description: "Check installed status and version of Prompt Architect"},
			{Command: "prompts-version", Description: "Show Prompt Architect version installed in repositories"},
		},
	}
}

func buildCGVersioningSection() termout.HelpSection {
	return termout.HelpSection{
		Title:      "Canonical Versioning & SSOT",
		Color:      termout.RoleSGRBold(termout.RoleWarn),
		BadgeText:  "SSOT",
		BadgeColor: termout.RoleSGR(termout.RoleWarn),
		Entries: []termout.CommandEntry{
			{Command: "version, ver, -v", Description: "Display installed Coding Guidelines version"},
			{Command: "install-version-json, init-version", Description: "Initialize canonical version.json SSOT with version tag"},
			{Command: "general, general-install", Description: "Multi-repo installation suite across discovered workspaces"},
		},
	}
}

func buildCGFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "--all", Description: "Execute across all discovered repositories in database"},
		{Command: "--dry-run", Description: "Simulate operations without modifying files"},
		{Command: "--exclude, --except <csv>", Description: "Comma-separated list of repository names to skip"},
		{Command: "--version <ver>", Description: "Starting version value for version.json (default: 1.0.0)"},
		{Command: "-j, --json", Description: "Emit results formatted as JSON telemetry"},
		{Command: "-h, --help", Description: "Display this formatted help screen"},
	}
}

func buildCGTips() []string {
	return []string{
		"Run 'gitmap cg audit' anytime to inspect coding guideline compliance.",
		"Use 'gitmap cg sync --all' to keep your entire local fleet synchronized with upstream.",
	}
}
