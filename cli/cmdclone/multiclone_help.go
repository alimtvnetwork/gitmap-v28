package cmdclone

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderMultiCloneHelp renders the rich two-column styled help menu for multiclone.
func RenderMultiCloneHelp() {
	termhelp.RenderMenu(buildMultiCloneHelpMenu())
}

func buildMultiCloneHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "Batch Multi-Repository Cloner (gitmap multiclone / mc)",
		UsageLines: []string{
			"gitmap mc <paste|file|urls> [target-dir] [flags]",
			"gitmap multiclone \"```\\nowner/repo\\nhttps://github.com/...\\n```\" [flags]",
			"gitmap mc repos.txt -d ./workspace",
			"gitmap mutliclone <paste-with-descriptions> [flags]",
		},
		Sections: []termhelp.HelpSection{
			buildMultiCloneInputSection(),
			buildMultiCloneFeaturesSection(),
			buildMultiCloneTargetSection(),
		},
		FooterFlags: buildMultiCloneFooterFlags(),
		Tips: []string{
			"Paste multiline blocks directly enclosed in ``` codeblocks or raw text.",
			"Deduplication is automatic: identical URLs and owner/repo entries are cloned once.",
			"Run with '--dry-run' to preview repositories and destinations without touching disk.",
		},
	}
}

func buildMultiCloneInputSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Supported Input Modes",
		Entries: []termhelp.CommandEntry{
			{Command: "``` <markdown> ```", Description: "Paste Markdown codeblock with URLs or repos"},
			{Command: "<file.txt>", Description: "Read repository list from local text file"},
			{Command: "-f, --file <path>", Description: "Explicit path to repository list file"},
			{Command: "cat list.txt | gitmap mc", Description: "Read newline-separated URLs from standard input"},
			{Command: "owner/repo", Description: "GitHub shorthand auto-expanded to full Git URL"},
		},
	}
}

func buildMultiCloneFeaturesSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Parsing & Processing",
		Entries: []termhelp.CommandEntry{
			{Command: "Auto Deduplication", Description: "Case-insensitive dedup; normalizes .git suffixes"},
			{Command: "Inline Descriptions", Description: "Strips ': description' after repository slugs"},
			{Command: "Universal URLs", Description: "Accepts HTTPS, HTTP, SSH, and git@host:owner/repo"},
			{Command: "VS Code PM Sync", Description: "Automatically registers cloned repos to Project Manager"},
		},
	}
}

func buildMultiCloneTargetSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Destination & Concurrency",
		Entries: []termhelp.CommandEntry{
			{Command: "-d, --dir <dir>", Description: "Clone all repositories into target directory"},
			{Command: "<target-dir>", Description: "Positional destination directory argument"},
			{Command: "-j, --concurrency <N>", Description: "Number of parallel clone workers (default: 1)"},
		},
	}
}

func buildMultiCloneFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "--dry-run", Description: "Preview parsed repositories and destinations without cloning"},
		{Command: "--no-replace", Description: "Fail rather than replace existing destination directories"},
		{Command: "--no-vscode-sync", Description: "Skip registering repos in VS Code Project Manager"},
		{Command: "--desktop", Description: "Register newly cloned repositories with GitHub Desktop"},
		{Command: "-h, --help", Description: "Show this multiclone help menu"},
	}
}
