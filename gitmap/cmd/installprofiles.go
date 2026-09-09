package cmd

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
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

	return []InstallProfile{
		buildMinimalProfile(),
		buildDevProfile(),
		buildUbuntuProfile(),
		buildUbuntuDevAIProfile(),
		buildAIProfile(),
		buildBackendProfile(),
		buildFullstackProfile(),
	}
}

func buildMinimalProfile() InstallProfile {

	return InstallProfile{
		Name:        "minimal",
		Title:       "Minimal dev workstation",
		Description: "Essential developer workstation (editor + git + node + python)",
		Tools:       []string{constants.ToolVSCode, constants.ToolGit, constants.ToolNodeJS, constants.ToolPython},
		Aliases:     []string{"min", "basic"},
	}
}

func buildDevProfile() InstallProfile {

	return InstallProfile{
		Name:        "dev",
		Title:       "Dev workstation with AI",
		Description: "Standard dev workstation + runtimes + Antigravity AI suite",
		Tools: []string{
			constants.ToolVSCode, constants.ToolGit, constants.ToolPython,
			constants.ToolNodeJS, constants.ToolPnpm, constants.ToolGo,
			constants.ToolRust, constants.ToolPHP, constants.ToolAntigravity,
			constants.ToolAgManager,
		},
		Aliases: []string{"developer", "dev-stack"},
	}
}

func buildUbuntuProfile() InstallProfile {

	return InstallProfile{
		Name:        "ubuntu",
		Title:       "Ubuntu developer workstation",
		Description: "Compiler toolchain, shell, browsers, dev runtimes & Antigravity",
		Tools: []string{
			constants.ToolBuildEssential, constants.ToolGit, constants.ToolZsh,
			constants.ToolVSCode, constants.ToolChrome, constants.ToolNodeJS,
			constants.ToolPython, constants.ToolGo, constants.ToolAntigravity,
			constants.ToolAgManager,
		},
		Aliases: []string{"ubuntu-dev", "linux-dev"},
	}
}

func buildUbuntuDevAIProfile() InstallProfile {

	return InstallProfile{
		Name:        "ubuntu-dev-ai",
		Title:       "Ubuntu AI / ML developer workstation",
		Description: "Full Ubuntu dev workstation + Ollama LLM + Antigravity AI suite",
		Tools: []string{
			constants.ToolBuildEssential, constants.ToolGit, constants.ToolZsh,
			constants.ToolVSCode, constants.ToolChrome, constants.ToolNodeJS,
			constants.ToolPython, constants.ToolGo, constants.ToolOllama,
			constants.ToolLlamaCpp, constants.ToolPythonLibs, constants.ToolAntigravity,
			constants.ToolAgManager,
		},
		Aliases: []string{"ubuntu-ai", "linux-ai"},
	}
}

func buildAIProfile() InstallProfile {

	return InstallProfile{
		Name:        "ai",
		Title:       "AI / ML workstation",
		Description: "Local LLM runners, Python ML libs, Antigravity & AG-Manager",
		Tools: []string{
			constants.ToolPython, constants.ToolOllama, constants.ToolLlamaCpp,
			constants.ToolPythonLibs, constants.ToolAntigravity, constants.ToolAgManager,
		},
		Aliases: []string{"ai-dev", "ml", "llm"},
	}
}

func buildBackendProfile() InstallProfile {

	return InstallProfile{
		Name:        "backend",
		Title:       "Backend developer workstation",
		Description: "Minimal stack + databases + docker + backend languages",
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
		Description: "Backend + pnpm, php, composer, mongodb & CI/CD",
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
	if strings.ToLower(p.Name) == low {

		return true
	}
	for _, a := range p.Aliases {
		if strings.ToLower(a) == low {

			return true
		}
	}

	return false
}

// IsInstallProfile checks whether the given string names an installation profile.
func IsInstallProfile(name string) bool {
	_, found := FindInstallProfile(name)

	return found
}
