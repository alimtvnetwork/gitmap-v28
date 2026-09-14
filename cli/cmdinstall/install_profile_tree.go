package cmdinstall

import (
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstaller"
)

// ToolEntry represents an individual tool included in a profile composition.
type ToolEntry struct {
	Slug        string
	Description string
}

// ProfileComposition represents a profile configuration and its inheritance chain.
type ProfileComposition struct {
	Name        string
	Alias       string
	Description string
	Base        *ProfileComposition
	Tools       []ToolEntry
}

func buildUbuntuBasicProfile() ProfileComposition {
	return ProfileComposition{
		Name:        "ubuntu-basic",
		Alias:       "ub",
		Description: "Basic Ubuntu toolchain and essential utilities",
		Base:        nil,
		Tools: []ToolEntry{
			{Slug: "curl", Description: "Command line tool for transferring data with URLs"},
			{Slug: "git", Description: "Fast, scalable, distributed revision control system"},
			{Slug: "build-essential", Description: "Meta-package for compiling software (gcc, make)"},
			{Slug: "zsh", Description: "Z shell with modern command line features"},
		},
	}
}

func buildUbuntuVscodeProfile() ProfileComposition {
	baseProfile := buildUbuntuBasicProfile()

	return ProfileComposition{
		Name:        "ubuntu+vscode",
		Alias:       "ub+code",
		Description: "Ubuntu basic plus Visual Studio Code editor",
		Base:        &baseProfile,
		Tools: []ToolEntry{
			{Slug: "code", Description: "Visual Studio Code editor binaries and extensions"},
		},
	}
}

func buildUbuntuSmallDevProfile() ProfileComposition {
	baseProfile := buildUbuntuVscodeProfile()

	return ProfileComposition{
		Name:        "ubuntu+small-dev",
		Alias:       "ub+sdev",
		Description: "Ubuntu lightweight development suite with runtime interpreters",
		Base:        &baseProfile,
		Tools:       buildUbuntuSmallDevTools(),
	}
}

func buildUbuntuSmallDevTools() []ToolEntry {
	return []ToolEntry{
		{Slug: "github-desktop", Description: "GitHub Desktop application"},
		{Slug: "git-compact", Description: "Git compact repository tool and cleaner"},
		{Slug: "golang", Description: "Go programming language runtime and tools"},
		{Slug: "rust", Description: "Rust programming language and Cargo toolchain"},
		{Slug: "php", Description: "PHP programming language and interpreter"},
		{Slug: "python3", Description: "Python 3 runtime environment"},
	}
}

func buildGitCompactProfileComposition() ProfileComposition {
	return ProfileComposition{
		Name:        "git-compact",
		Alias:       "gitcompact",
		Description: "Git compact repository tool and version control suite",
		Base:        nil,
		Tools: []ToolEntry{
			{Slug: "git", Description: "Fast, scalable, distributed revision control system"},
			{Slug: "github-desktop", Description: "GitHub Desktop application"},
			{Slug: "git-compact", Description: "Git compact repository tool and cleaner"},
		},
	}
}

func buildUbuntuDevProfile() ProfileComposition {
	baseProfile := buildUbuntuSmallDevProfile()

	return ProfileComposition{
		Name:        "ubuntu+dev",
		Alias:       "ub+dev",
		Description: "Full Ubuntu developer workstation suite",
		Base:        &baseProfile,
		Tools:       buildUbuntuDevTools(),
	}
}

func buildUbuntuDevTools() []ToolEntry {
	return []ToolEntry{
		{Slug: "nodejs", Description: "Node.js JavaScript runtime and npm"},
		{Slug: "pnpm", Description: "Fast, disk space efficient package manager"},
		{Slug: "yarn", Description: "Fast, reliable, and secure dependency management"},
		{Slug: "antigravity", Description: "Google Antigravity Desktop IDE environment"},
		{Slug: "ag-manager", Description: "Antigravity manager GUI and toolchain"},
	}
}

func buildAntigravityProfile() ProfileComposition {
	return ProfileComposition{
		Name:        "antigravity",
		Alias:       "ag",
		Description: "Antigravity IDE and Agent CLI Environment",
		Base:        nil,
		Tools: []ToolEntry{
			{Slug: "antigravity", Description: "Google Antigravity Desktop IDE environment"},
			{Slug: "ag-manager", Description: "Antigravity manager GUI and toolchain"},
		},
	}
}

func buildAiToolsProfileTree() ProfileComposition {
	return ProfileComposition{
		Name:        "ai-tools",
		Alias:       "all-ai",
		Description: "Comprehensive AI development suite and local LLMs",
		Base:        nil,
		Tools:       buildAiToolsEntries(),
	}
}

func buildAiToolsEntries() []ToolEntry {
	return []ToolEntry{
		{Slug: "python3", Description: "Python 3 runtime environment"},
		{Slug: "ollama", Description: "Ollama local large language model runner"},
		{Slug: "llama-cpp", Description: "llama.cpp LLM inference engine in C/C++"},
		{Slug: "python-libs", Description: "Python AI, ML, and data science libraries"},
		{Slug: "antigravity", Description: "Google Antigravity Desktop IDE environment"},
		{Slug: "ag-manager", Description: "Antigravity manager GUI and toolchain"},
	}
}

func buildAntigravitySuiteProfileTree() ProfileComposition {
	return ProfileComposition{
		Name:        "antigravity-suite",
		Alias:       "ag-suite",
		Description: "Antigravity Desktop IDE, Manager, and CLI suite",
		Base:        nil,
		Tools: []ToolEntry{
			{Slug: "antigravity", Description: "Google Antigravity Desktop IDE environment"},
			{Slug: "ag-manager", Description: "Antigravity manager GUI and toolchain"},
			{Slug: "agy", Description: "Antigravity CLI autonomous coding assistant"},
		},
	}
}

func buildVscodeSettingsProfile() ProfileComposition {
	return ProfileComposition{
		Name:        "vscode-settings",
		Alias:       "vscode+settings",
		Description: "Visual Studio Code with synchronized settings",
		Base:        nil,
		Tools: []ToolEntry{
			{Slug: "code", Description: "Visual Studio Code editor installation"},
			{Slug: "settings.json", Description: "VS Code user settings configuration"},
			{Slug: "keybindings.json", Description: "VS Code keyboard shortcuts mapping"},
			{Slug: "extensions", Description: "Recommended VS Code workspace extensions"},
		},
	}
}

func buildUbuntuBuildEssentialProfile() ProfileComposition {
	return ProfileComposition{
		Name:        "build-essential",
		Alias:       "ubuntu-common",
		Description: "Ubuntu build-essential compiler toolchain & common utilities",
		Base:        nil,
		Tools: []ToolEntry{
			{Slug: "build-essential", Description: "Meta-package for compiling software (gcc, g++, make)"},
			{Slug: "gcc", Description: "GNU C compiler"},
			{Slug: "g++", Description: "GNU C++ compiler"},
			{Slug: "make", Description: "GNU make utility"},
			{Slug: "git", Description: "Git distributed version control"},
			{Slug: "git-lfs", Description: "Git Large File Storage"},
			{Slug: "curl", Description: "Command line data transfer utility"},
			{Slug: "wget", Description: "Non-interactive network downloader"},
			{Slug: "vim", Description: "Vi IMproved text editor"},
			{Slug: "nano", Description: "Terminal text editor"},
			{Slug: "zsh", Description: "Z shell command environment"},
			{Slug: "file", Description: "File type identification utility"},
			{Slug: "sshpass", Description: "Non-interactive ssh password provider"},
			{Slug: "snapd", Description: "Snap package management daemon"},
		},
	}
}

func resolveProfileTree(slug string) (ProfileComposition, bool) {
	norm := strings.ToLower(strings.TrimSpace(slug))
	if prof, ok := resolveUbuntuProfileTree(norm); ok {
		return prof, true
	}

	return resolveSpecialProfileTree(norm)
}

func resolveUbuntuProfileTree(norm string) (ProfileComposition, bool) {
	switch norm {
	case "ubuntu-basic", "ubuntu+basic", "ub":
		return buildUbuntuBasicProfile(), true
	case "ubuntu+vscode", "ubuntu-vscode", "ub+code":
		return buildUbuntuVscodeProfile(), true
	case "ubuntu+small-dev", "ubuntu-small-dev", "ub+sdev", "small-dev", "simple-dev", "simpledev":
		return buildUbuntuSmallDevProfile(), true
	case "ubuntu+dev", "ubuntu-dev", "ub+dev", "dev":
		return buildUbuntuDevProfile(), true
	default:
		return ProfileComposition{}, false
	}
}

func buildTerminalProfileComposition() ProfileComposition {
	return ProfileComposition{
		Name:        "terminal",
		Alias:       "term",
		Description: "Terminal essentials: console emulators, editors, and modern shell",
		Base:        nil,
		Tools:       buildTerminalToolEntries(),
	}
}

func buildTerminalToolEntries() []ToolEntry {
	if runtime.GOOS == "windows" {
		return buildWindowsTerminalToolEntries()
	}

	return buildUnixTerminalToolEntries()
}

func buildWindowsTerminalToolEntries() []ToolEntry {
	return []ToolEntry{
		{Slug: "conemu", Description: "ConEmu Windows console emulator with tabs and splits"},
		{Slug: "notepad++", Description: "Notepad++ source code and text editor"},
		{Slug: "powershell", Description: "PowerShell 7 cross-platform automation shell"},
		{Slug: "chrome", Description: "Google Chrome web browser"},
		{Slug: "ubuntu-font", Description: "Ubuntu typeface font family for programming"},
		{Slug: "7zip", Description: "7-Zip high compression ratio file archiver"},
		{Slug: "winrar", Description: "WinRAR archiver and RAR compression utility"},
	}
}

func buildUnixTerminalToolEntries() []ToolEntry {
	return []ToolEntry{
		{Slug: "zsh", Description: "Z shell command environment"},
		{Slug: "powershell", Description: "PowerShell 7 cross-platform automation shell"},
		{Slug: "chrome", Description: "Google Chrome web browser"},
		{Slug: "ubuntu-font", Description: "Ubuntu typeface font family for programming"},
		{Slug: "jq", Description: "Command-line JSON processor"},
		{Slug: "yq", Description: "Command-line YAML/JSON/XML processor"},
		{Slug: "zellij", Description: "Terminal workspace and multiplexer"},
	}
}

func resolveSpecialProfileTree(norm string) (ProfileComposition, bool) {
	if prof, isFound := resolveWorkstationProfileTree(norm); isFound {
		return prof, true
	}

	return resolveAiSuiteProfileTree(norm)
}

func resolveWorkstationProfileTree(norm string) (ProfileComposition, bool) {
	switch norm {
	case "terminal", "term", "cli", "terminal-profile", "profile-terminal", "terminalprofile", "terminal-essentials":
		return buildTerminalProfileComposition(), true
	case "git-compact", "gitcompact", "profile-git", "profile-git-compact":
		return buildGitCompactProfileComposition(), true
	case "vscode+settings", "vscode-settings", "vscode":
		return buildVscodeSettingsProfile(), true
	case "build-essential", "buildessential", "be", "ubuntu-common", "ub-common":
		return buildUbuntuBuildEssentialProfile(), true
	default:
		return ProfileComposition{}, false
	}
}

func resolveAiSuiteProfileTree(norm string) (ProfileComposition, bool) {
	switch norm {
	case "ag", "antigravity":
		return buildAntigravityProfile(), true
	case "ai-tools", "all-ai", "aitools":
		return buildAiToolsProfileTree(), true
	case "antigravity-suite", "ag-suite", "antigravitysuite":
		return buildAntigravitySuiteProfileTree(), true
	default:
		return ProfileComposition{}, false
	}
}

func profileToTreeNode(profile ProfileComposition) cmdinstaller.InstallerTreeNode {
	rootNode := cmdinstaller.InstallerTreeNode{
		Title:       profile.Name,
		Description: profile.Description,
	}

	if profile.Base != nil {
		rootNode.Children = append(rootNode.Children, profileToTreeNode(*profile.Base))
	}

	return appendToolNodes(rootNode, profile.Tools)
}

func appendToolNodes(parent cmdinstaller.InstallerTreeNode, tools []ToolEntry) cmdinstaller.InstallerTreeNode {
	for _, toolEntry := range tools {
		parent.Children = append(parent.Children, cmdinstaller.InstallerTreeNode{
			Title:       toolEntry.Slug,
			Description: toolEntry.Description,
		})
	}

	return parent
}

func printProfileTree(profile ProfileComposition) {
	rootNode := profileToTreeNode(profile)
	cmdinstaller.PrintInstallerTree(rootNode, "", true)
}

func printProfileInstallSummary(slug string) {
	profile, hasProfile := resolveProfileTree(slug)
	if !hasProfile {
		cmdinstaller.PrintInstallSummaryHeader(slug)

		return
	}

	cmdinstaller.PrintInstallSummaryHeader(profile.Name)
	printProfileTree(profile)
}
