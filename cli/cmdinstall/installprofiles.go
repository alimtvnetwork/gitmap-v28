package cmdinstall

import (
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// InstallProfile defines a bundled suite of developer tools.
type InstallProfile struct {
	Name        string
	Title       string
	Description string
	Tools       []string
	Aliases     []string
}

// AllInstallProfiles returns the registered installation profiles.
func AllInstallProfiles() []InstallProfile {
	profiles := make([]InstallProfile, 0, 18)
	profiles = append(profiles, getCoreWorkstationProfiles()...)

	return append(profiles, getSpecializedWorkstationProfiles()...)
}

func getCoreWorkstationProfiles() []InstallProfile {
	return []InstallProfile{
		buildMinimalProfile(),
		buildBaseProfile(),
		buildGitCompactProfile(),
		buildAdvanceProfile(),
		buildCppDxProfile(),
		buildSmallDevProfile(),
		buildDevProfile(),
		buildDevAdvanceProfile(),
	}
}

func buildMinimalProfile() InstallProfile {
	return InstallProfile{
		Name:        "minimal",
		Title:       "Minimal dev workstation",
		Description: "Editor, Git, Node.js, and Python",
		Tools:       []string{constants.ToolVSCode, constants.ToolGit, constants.ToolNodeJS, constants.ToolPython},
		Aliases:     []string{"min", "basic"},
	}
}

func buildBaseProfile() InstallProfile {
	return InstallProfile{
		Name:        "base",
		Title:       "Base Windows workstation",
		Description: "Daily-driver: Git, media, archivers, fonts, editor, and browser",
		Tools: []string{
			constants.ToolGit, constants.ToolVLC, constants.Tool7Zip,
			constants.ToolWinRAR, constants.ToolUbuntuFont, constants.ToolXMind,
			constants.ToolNpp, constants.ToolChrome, constants.ToolConemu,
		},
		Aliases: []string{"workstation", "daily"},
	}
}

func buildGitCompactProfile() InstallProfile {
	return InstallProfile{
		Name:        "git-compact",
		Title:       "Git compact workstation",
		Description: "Git version control and GitHub Desktop",
		Tools:       []string{constants.ToolGit, constants.ToolGitHubDesktop, constants.ToolGitCompact},
		Aliases:     []string{"gitcompact", "profile-git", "profile-git-compact"},
	}
}

func buildAdvanceProfile() InstallProfile {
	return InstallProfile{
		Name:        "advance",
		Title:       "Advance workstation",
		Description: "Base + git-compact + wordweb, beyondcompare, obs, whatsapp, vscode",
		Tools: []string{
			constants.ToolGit, constants.ToolVLC, constants.Tool7Zip, constants.ToolWinRAR,
			constants.ToolUbuntuFont, constants.ToolXMind, constants.ToolNpp, constants.ToolChrome,
			constants.ToolConemu, constants.ToolGitHubDesktop, constants.ToolGitCompact, constants.ToolWordWeb,
			constants.ToolBeyondCompare, constants.ToolOBS, constants.ToolWhatsApp,
			constants.ToolVSCode, constants.ToolVSCodeSync,
		},
		Aliases: []string{"advanced"},
	}
}

func buildCppDxProfile() InstallProfile {
	return InstallProfile{
		Name:        "cpp-dx",
		Title:       "C++ and DirectX development",
		Description: "VC++ runtimes, DirectX runtime, and DirectX SDK",
		Tools: []string{
			constants.ToolVcRedist, constants.ToolDirectX, constants.ToolDirectXSdk,
		},
		Aliases: []string{"cppdx", "directx"},
	}
}

func buildSmallDevProfile() InstallProfile {
	return InstallProfile{
		Name:        "small-dev",
		Title:       "Small dev workstation",
		Description: "Advance profile + Go programming language",
		Tools:       resolveSmallDevTools(),
		Aliases:     []string{"smalldev", "slim-dev", "simple-dev", "simpledev"},
	}
}

func resolveSmallDevTools() []string {
	if runtime.GOOS == "windows" {
		return resolveWindowsSmallDevTools()
	}

	return resolveLinuxSmallDevTools()
}

func resolveWindowsSmallDevTools() []string {
	return []string{
		constants.ToolGit, constants.ToolVLC, constants.Tool7Zip, constants.ToolWinRAR,
		constants.ToolUbuntuFont, constants.ToolXMind, constants.ToolNpp, constants.ToolChrome,
		constants.ToolConemu, constants.ToolGitHubDesktop, constants.ToolGitCompact, constants.ToolWordWeb,
		constants.ToolBeyondCompare, constants.ToolOBS, constants.ToolWhatsApp,
		constants.ToolVSCode, constants.ToolVSCodeSync, constants.ToolGo,
	}
}

func resolveLinuxSmallDevTools() []string {
	return []string{
		constants.ToolGit, constants.ToolZsh, constants.ToolAria2,
		constants.ToolBuildEssential, constants.ToolVSCode, constants.ToolVSCodeSync,
		constants.ToolGitHubDesktop, constants.ToolGitCompact, constants.ToolGo,
		constants.ToolRust, constants.ToolPHP, constants.ToolPython,
	}
}

func buildDevProfile() InstallProfile {
	return InstallProfile{
		Name:        "dev",
		Title:       "Dev workstation with AI",
		Description: "Standard dev workstation + runtimes + AI suite",
		Tools:       resolveDevTools(),
		Aliases:     []string{"developer", "dev-stack"},
	}
}

func resolveDevTools() []string {
	if runtime.GOOS == "windows" {
		return resolveWindowsDevTools()
	}

	return resolveLinuxDevTools()
}

func resolveWindowsDevTools() []string {
	extra := []string{
		constants.ToolPython, constants.ToolNodeJS, constants.ToolPnpm,
		constants.ToolRust, constants.ToolPHP, constants.ToolAntigravity,
		constants.ToolAgManager,
	}

	return append(resolveSmallDevTools(), extra...)
}

func resolveLinuxDevTools() []string {
	extra := []string{
		constants.ToolNodeJS, constants.ToolPnpm, constants.ToolYarn,
		constants.ToolAntigravity, constants.ToolAgManager,
	}

	return append(resolveSmallDevTools(), extra...)
}

func buildDevAdvanceProfile() InstallProfile {
	return InstallProfile{
		Name:        "dev-advance",
		Title:       "Dev advance polyglot workstation",
		Description: "Dev profile + .NET SDK + C++/DirectX suite",
		Tools: []string{
			constants.ToolVSCode, constants.ToolGit, constants.ToolPython,
			constants.ToolNodeJS, constants.ToolPnpm, constants.ToolGo,
			constants.ToolRust, constants.ToolPHP, constants.ToolAntigravity,
			constants.ToolAgManager, constants.ToolDotnet, constants.ToolVcRedist,
			constants.ToolDirectX, constants.ToolDirectXSdk,
		},
		Aliases: []string{"devadvance", "dev-plus"},
	}
}

// FindInstallProfile resolves a profile by name or alias.
func FindInstallProfile(name string) (InstallProfile, bool) {
	low := strings.ToLower(strings.TrimSpace(name))
	for _, p := range AllInstallProfiles() {
		if matchesProfile(p, low) {
			return p, true
		}
	}

	return InstallProfile{}, false
}

func matchesProfile(p InstallProfile, low string) bool {
	pName := strings.ToLower(p.Name)
	if pName == low || pName == stripProfileAffixes(low) {
		return true
	}

	return matchesProfileAliases(p.Aliases, low)
}

func matchesProfileAliases(aliases []string, low string) bool {
	for _, a := range aliases {
		aLow := strings.ToLower(a)
		if aLow == low || aLow == stripProfileAffixes(low) {
			return true
		}
	}

	return false
}

func stripProfileAffixes(s string) string {
	s = strings.TrimPrefix(s, "profile-")
	s = strings.TrimPrefix(s, "profile:")

	return strings.TrimSuffix(s, "-profile")
}

// IsInstallProfile checks whether the given string names an installation profile.
func IsInstallProfile(name string) bool {
	_, isFound := FindInstallProfile(name)

	return isFound
}
