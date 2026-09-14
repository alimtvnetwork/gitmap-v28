package cmdinstall

import (
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func getSpecializedWorkstationProfiles() []InstallProfile {
	return []InstallProfile{
		buildTerminalProfile(),
		buildWebDevProfile(),
		buildDevopsProfile(),
		buildUbuntuProfile(),
		buildAIProfile(),
		buildAiToolsProfile(),
		buildAntigravitySuiteProfile(),
		buildBackendProfile(),
		buildFullstackProfile(),
	}
}

func buildTerminalProfile() InstallProfile {
	return InstallProfile{
		Name:        "terminal",
		Title:       "Terminal essentials workstation",
		Description: "ConEmu, Notepad++, PowerShell, Chrome, Ubuntu font, 7-Zip, WinRAR",
		Tools:       resolveTerminalProfileTools(),
		Aliases: []string{
			"term", "cli", "terminal-profile", "profile-terminal",
			"terminalprofile", "terminal-essentials",
		},
	}
}

func resolveTerminalProfileTools() []string {
	if runtime.GOOS == "windows" {
		return resolveWindowsTerminalTools()
	}

	return resolveUnixTerminalTools()
}

func resolveWindowsTerminalTools() []string {
	return []string{
		constants.ToolConemu, constants.ToolNpp, constants.ToolPowerShell,
		constants.ToolChrome, constants.ToolUbuntuFont, constants.Tool7Zip,
		constants.ToolWinRAR,
	}
}

func resolveUnixTerminalTools() []string {
	return []string{
		constants.ToolZsh, constants.ToolPowerShell, constants.ToolChrome,
		constants.ToolUbuntuFont, constants.ToolJq, constants.ToolYq,
		constants.ToolZellij,
	}
}

func buildWebDevProfile() InstallProfile {
	return InstallProfile{
		Name:        "web-dev",
		Title:       "Web developer workstation",
		Description: "VS Code, Node.js, pnpm, Git, and VS Code settings",
		Tools: []string{
			constants.ToolVSCode, constants.ToolNodeJS, constants.ToolPnpm,
			constants.ToolGit, constants.ToolVSCodeSync,
		},
		Aliases: []string{"webdev", "frontend"},
	}
}

func buildDevopsProfile() InstallProfile {
	return InstallProfile{
		Name:        "devops",
		Title:       "DevOps and infrastructure workstation",
		Description: "Git, Docker container platform, and Kubernetes CLI",
		Tools:       []string{constants.ToolGit, constants.ToolDocker, constants.ToolKubernetes},
		Aliases:     []string{"infra", "cloud"},
	}
}

func buildUbuntuProfile() InstallProfile {
	return InstallProfile{
		Name:        "ubuntu",
		Title:       "Ubuntu developer workstation",
		Description: "Compiler toolchain, shell, browsers, runtimes",
		Tools: []string{
			constants.ToolBuildEssential, constants.ToolGit, constants.ToolZsh,
			constants.ToolVSCode, constants.ToolChrome, constants.ToolNodeJS,
			constants.ToolPython, constants.ToolGo, constants.ToolAntigravity,
			constants.ToolAgManager,
		},
		Aliases: []string{"ubuntu-dev", "linux-dev"},
	}
}

func buildAIProfile() InstallProfile {
	return InstallProfile{
		Name:        "ai",
		Title:       "AI / ML workstation",
		Description: "Local LLM runners, Python ML libs & Antigravity",
		Tools: []string{
			constants.ToolPython, constants.ToolOllama, constants.ToolLlamaCpp,
			constants.ToolPythonLibs, constants.ToolAntigravity, constants.ToolAgManager,
		},
		Aliases: []string{"ai-dev", "ml", "llm"},
	}
}

func buildAiToolsProfile() InstallProfile {
	return InstallProfile{
		Name:        "ai-tools",
		Title:       "AI tools suite",
		Description: "Local LLM runners, Python ML libs & Antigravity tools",
		Tools: []string{
			constants.ToolPython, constants.ToolOllama, constants.ToolLlamaCpp,
			constants.ToolPythonLibs, constants.ToolAntigravity, constants.ToolAgManager,
		},
		Aliases: []string{"all-ai", "aitools"},
	}
}

func buildAntigravitySuiteProfile() InstallProfile {
	return InstallProfile{
		Name:        "antigravity-suite",
		Title:       "Antigravity workstation suite",
		Description: "Antigravity Desktop IDE, Manager, and Agent CLI suite",
		Tools: []string{
			constants.ToolAntigravity, constants.ToolAgManager, constants.ToolAgy,
		},
		Aliases: []string{"ag-suite", "antigravitysuite"},
	}
}

func buildBackendProfile() InstallProfile {
	return InstallProfile{
		Name:        "backend",
		Title:       "Backend developer workstation",
		Description: "Minimal stack + databases, Docker & languages",
		Tools: []string{
			constants.ToolVSCode, constants.ToolGit, constants.ToolNodeJS,
			constants.ToolPython, constants.ToolDocker, constants.ToolMySQL,
			constants.ToolPostgreSQL, constants.ToolRedis, constants.ToolGo,
			constants.ToolDotnet, constants.ToolJava,
		},
		Aliases: []string{"back", "server"},
	}
}

func buildFullstackProfile() InstallProfile {
	return InstallProfile{
		Name:        "fullstack",
		Title:       "Full-stack web workstation",
		Description: "Backend + pnpm, PHP, Composer, MongoDB & CI/CD",
		Tools: []string{
			constants.ToolVSCode, constants.ToolGit, constants.ToolNodeJS,
			constants.ToolPython, constants.ToolPnpm, constants.ToolPHP,
			constants.ToolComposer, constants.ToolDocker, constants.ToolMySQL,
			constants.ToolPostgreSQL, constants.ToolMongoDB, constants.ToolRedis,
			constants.ToolGo, constants.ToolJenkins,
		},
		Aliases: []string{"full", "web"},
	}
}
