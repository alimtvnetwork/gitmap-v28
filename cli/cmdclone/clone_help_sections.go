package cmdclone

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

func buildCloneSubcmdSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Clone Subcommands",
		Entries: []termhelp.CommandEntry{
			{Command: "next", Description: "Clone next pending repository in active manifest"},
			{Command: "now <url>", Description: "Immediate priority clone bypassing worker queue"},
			{Command: "only-missing", Description: "Skip already cloned repos, only clone missing"},
			{Command: "sync", Description: "Synchronize remote URLs and tracking branches"},
			{Command: "pick [pattern]", Description: "Filter and interactively select repos to clone"},
			{Command: "fix-repo (fr)", Description: "Rewrite legacy versioned folder tokens after clone"},
		},
	}
}

func buildCloneConcurrencySection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Concurrency & Performance",
		Entries: []termhelp.CommandEntry{
			{Command: "--max-concurrency <n>", Description: "Run up to N parallel clones (default: NumCPU)"},
			{Command: "--target-dir <dir>", Description: "Base directory for cloned repositories"},
			{Command: "--default-branch <b>", Description: "Fallback branch when HEAD detection fails"},
			{Command: "--output <mode>", Description: "Stream output mode (terminal, json, text)"},
		},
	}
}

func buildCloneEnvSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Environment & Integrations",
		Entries: []termhelp.CommandEntry{
			{Command: "--github-desktop", Description: "Auto-register cloned repos with GitHub Desktop"},
			{Command: "--key <name>", Description: "Use specific managed SSH key for Git operations"},
			{Command: "--reclone", Description: "Alias for --force to cleanly recreate targets"},
		},
	}
}
