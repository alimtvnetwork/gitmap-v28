package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// resolvePackageName maps tool name to package ID for a manager.
func resolvePackageName(manager, tool string) string {
	switch manager {
	case constants.PkgMgrWinget:

		return resolveWingetPackage(tool)
	case constants.PkgMgrApt:

		return resolveAptPackage(tool)
	case constants.PkgMgrBrew:

		return resolveBrewPackage(tool)
	case constants.PkgMgrSnap:

		return resolveSnapPackage(tool)
	default:

		return resolveChocoPackage(tool)
	}
}

// resolveToolAlias normalizes known tool aliases to their canonical tool name.
func resolveToolAlias(tool string) string {
	switch tool {
	case "k8s", "kubectl":

		return constants.ToolKubernetes
	case "dotnet-sdk":

		return constants.ToolDotnet
	case "jdk", "openjdk":

		return constants.ToolJava
	case "rustup", "cargo":

		return constants.ToolRust
	case "llamacpp":

		return constants.ToolLlamaCpp
	default:

		return tool
	}
}

var chocoPackageMap = map[string]string{
	constants.ToolVSCode:        constants.ChocoPkgVSCode,
	constants.ToolNodeJS:        constants.ChocoPkgNodeJS,
	constants.ToolYarn:          constants.ChocoPkgYarn,
	constants.ToolBun:           constants.ChocoPkgBun,
	constants.ToolPnpm:          constants.ChocoPkgPnpm,
	constants.ToolPython:        constants.ChocoPkgPython,
	constants.ToolGo:            constants.ChocoPkgGo,
	constants.ToolGit:           constants.ChocoPkgGit,
	constants.ToolGitLFS:        constants.ChocoPkgGitLFS,
	constants.ToolGHCLI:         constants.ChocoPkgGHCLI,
	constants.ToolGitHubDesktop: constants.ChocoPkgGitHubDesktop,
	constants.ToolCPP:           constants.ChocoPkgCPP,
	constants.ToolPHP:           constants.ChocoPkgPHP,
	constants.ToolPowerShell:    constants.ChocoPkgPowerShell,
	constants.ToolMySQL:         constants.ChocoPkgMySQL,
	constants.ToolMariaDB:       constants.ChocoPkgMariaDB,
	constants.ToolPostgreSQL:    constants.ChocoPkgPostgreSQL,
	constants.ToolSQLite:        constants.ChocoPkgSQLite,
	constants.ToolMongoDB:       constants.ChocoPkgMongoDB,
	constants.ToolCouchDB:       constants.ChocoPkgCouchDB,
	constants.ToolRedis:         constants.ChocoPkgRedis,
	constants.ToolNeo4j:         constants.ChocoPkgNeo4j,
	constants.ToolElasticsearch: constants.ChocoPkgElasticsearch,
	constants.ToolDuckDB:        constants.ChocoPkgDuckDB,
	constants.ToolNpp:           constants.ChocoPkgNpp,
	constants.ToolNppInstall:    constants.ChocoPkgNpp,
	constants.ToolDbeaver:       constants.ChocoPkgDbeaver,
	constants.ToolOBS:           constants.ChocoPkgOBS,
	constants.ToolChrome:        constants.ChocoPkgChrome,
	constants.ToolGoogleChrome:  constants.ChocoPkgChrome,
	constants.ToolRust:          constants.ChocoPkgRust,
	constants.ToolDotnet:        constants.ChocoPkgDotnet,
	constants.ToolJava:          constants.ChocoPkgJava,
	constants.ToolFlutter:       constants.ChocoPkgFlutter,
	constants.ToolOllama:        constants.ChocoPkgOllama,
	constants.ToolLlamaCpp:      constants.ChocoPkgLlamaCpp,
	constants.ToolPythonLibs:    constants.ChocoPkgPythonLibs,
	constants.ToolDocker:        constants.ChocoPkgDocker,
	constants.ToolKubernetes:    constants.ChocoPkgKubernetes,
	constants.ToolJenkins:       constants.ChocoPkgJenkins,
	constants.ToolZsh:           constants.ChocoPkgZsh,
	constants.ToolFlameshot:     constants.ChocoPkgFlameshot,
	constants.ToolConemu:        constants.ChocoPkgConemu,
	constants.ToolVLC:           constants.ChocoPkgVLC,
}

// resolveChocoPackage maps tool names to Chocolatey package IDs.
func resolveChocoPackage(tool string) string {
	if pkg, exists := chocoPackageMap[tool]; exists {

		return pkg
	}

	return tool
}

var wingetPackageMap = map[string]string{
	constants.ToolVSCode:        constants.WingetPkgVSCode,
	constants.ToolPowerShell:    constants.WingetPkgPowerShell,
	constants.ToolDbeaver:       constants.WingetPkgDbeaver,
	constants.ToolOBS:           constants.WingetPkgOBS,
	constants.ToolStickyNotes:   constants.WingetPkgStickyNotes,
	constants.ToolGitHubDesktop: constants.WingetPkgGitHubDesktop,
	constants.ToolChrome:        constants.WingetPkgChrome,
	constants.ToolGoogleChrome:  constants.WingetPkgChrome,
	constants.ToolRust:          constants.WingetPkgRust,
	constants.ToolDotnet:        constants.WingetPkgDotnet,
	constants.ToolJava:          constants.WingetPkgJava,
	constants.ToolFlutter:       constants.WingetPkgFlutter,
	constants.ToolOllama:        constants.WingetPkgOllama,
	constants.ToolLlamaCpp:      constants.WingetPkgLlamaCpp,
	constants.ToolPythonLibs:    constants.WingetPkgPythonLibs,
	constants.ToolDocker:        constants.WingetPkgDocker,
	constants.ToolKubernetes:    constants.WingetPkgKubernetes,
	constants.ToolJenkins:       constants.WingetPkgJenkins,
	constants.ToolZsh:           constants.WingetPkgZsh,
	constants.ToolFlameshot:     constants.WingetPkgFlameshot,
	constants.ToolConemu:        constants.WingetPkgConemu,
	constants.ToolVLC:           constants.WingetPkgVLC,
}

// resolveWingetPackage maps tool names to Winget package IDs.
func resolveWingetPackage(tool string) string {
	if pkg, exists := wingetPackageMap[tool]; exists {

		return pkg
	}

	return tool
}
