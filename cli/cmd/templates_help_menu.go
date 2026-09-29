package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderTemplatesHelp displays the styled two-column Templates help menu.
func RenderTemplatesHelp() {
	termhelp.RenderMenu(buildTemplatesHelpMenu())
}

func buildTemplatesHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "Template Engine & Variable Scaffolding (gitmap templates)",
		UsageLines: []string{
			"gitmap templates <subcommand> [flags]",
			"gitmap tpl <subcommand> [flags]",
			"gitmap templates init go node --lfs",
		},
		Sections: []termhelp.HelpSection{
			buildTemplatesScaffoldSection(),
			buildTemplatesStateDBSection(),
			buildTemplatesVariablesSection(),
		},
		FooterFlags: buildTemplatesFooterFlags(),
		Tips: []string{
			"Run 'gitmap templates init go' to scaffold .gitignore and .gitattributes.",
			"Run 'gitmap templates ui' to manage template variables in a local browser UI.",
		},
	}
}

func buildTemplatesFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "--kind <ignore|attributes|lfs>", Description: "Filter templates to specific artifact kind"},
		{Command: "--lang <name>", Description: "Filter templates to target language (go, node, python)"},
		{Command: "--lfs", Description: "Include Git LFS attributes during init scaffolding"},
		{Command: "--dry-run", Description: "Preview generated templates without touching disk"},
		{Command: "--force", Description: "Overwrite existing target files outright"},
		{Command: "--port <n>", Description: "HTTP port for Templates Web UI (default: 8787)"},
		{Command: "-h, --help", Description: "Show this templates help menu"},
	}
}

func buildTemplatesScaffoldSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Built-In Scaffolding",
		Entries: []termhelp.CommandEntry{
			{Command: "list (tl)", Description: "List all available embedded templates with language filter"},
			{Command: "show (ts) <kind> <lang>", Description: "Print formatted template content directly to stdout"},
			{Command: "init (ti) <lang...>", Description: "Scaffold .gitignore / .gitattributes into current repo"},
			{Command: "diff (td) --lang <l>", Description: "Preview what template initialization would change"},
		},
	}
}

func buildTemplatesStateDBSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Database State Templates",
		Entries: []termhelp.CommandEntry{
			{Command: "ls [--category <c>]", Description: "List state templates stored in gitmap-templates.db"},
			{Command: "add [flags]", Description: "Add or upsert custom template with category and slug"},
			{Command: "edit <id|slug>", Description: "Update title, body, or metadata of existing template"},
			{Command: "remove <id|slug>", Description: "Delete template record from templates database"},
			{Command: "import / export", Description: "Transfer templates and variables via JSON manifests"},
		},
	}
}

func buildTemplatesVariablesSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Variables & Interactive UI",
		Entries: []termhelp.CommandEntry{
			{Command: "var <set|ls|rm>", Description: "Manage template substitution variables ($VAR syntax)"},
			{Command: "ui [--port 8787]", Description: "Launch the local reactive Templates & Variables Web UI"},
		},
	}
}
