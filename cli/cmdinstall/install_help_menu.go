package cmdinstall

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderInstallHelp displays the styled two-column Install help menu.
func RenderInstallHelp() {
	termhelp.RenderMenu(buildInstallHelpMenu())
}

func buildInstallHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "Tool & Profile Installer (gitmap install)",
		UsageLines: []string{
			"gitmap install <tool|profile> [flags]",
			"gitmap in <tool|profile> [flags]",
			"gitmap in profile [name] [--tree]",
		},
		Sections: []termhelp.HelpSection{
			buildInstallCoreSection(),
			buildInstallDatabaseSection(),
			buildInstallAISection(),
			buildInstallProfileSection(),
			buildInstallCustomSection(),
			buildInstallMgmtSection(),
		},
		FooterFlags: buildInstallFooterFlags(),
		Tips: []string{
			"Run 'gitmap in profile dev --tree' to preview the complete tool tree.",
			"Use 'gitmap in <tool> --check' to verify if a tool is installed.",
			"Run 'gitmap in ls' to list all supported tools with installed status.",
		},
	}
}

func buildInstallFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "-y, --yes", Description: "Auto-confirm installation prompts without asking"},
		{Command: "--dry-run", Description: "Preview resolved install commands without executing"},
		{Command: "--check", Description: "Only check if tool is installed on system"},
		{Command: "-t, --tree", Description: "Preview full tool hierarchy of a profile"},
		{Command: "--manager <name>", Description: "Force package manager (choco, winget, apt, brew)"},
		{Command: "--version <ver>", Description: "Install a specific tool version"},
		{Command: "--explain", Description: "Print resolved command before executing"},
		{Command: "-h, --help", Description: "Show this installer help menu"},
	}
}

func buildInstallCoreSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Core Developer Runtimes",
		Entries: []termhelp.CommandEntry{
			{Command: "vscode (code)", Description: "Visual Studio Code editor"},
			{Command: "git", Description: "Git version control system"},
			{Command: "node", Description: "Node.js JavaScript runtime (npm/npx)"},
			{Command: "python (py)", Description: "Python 3 programming language"},
			{Command: "go (golang)", Description: "Go programming language and toolchain"},
			{Command: "rust (cargo)", Description: "Rust language & Cargo package manager"},
			{Command: "dotnet (sdk)", Description: ".NET SDK and developer platform"},
			{Command: "docker", Description: "Docker container engine & CLI"},
		},
	}
}
