package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

type packageResolver func(string) string

var managerResolvers = map[string]packageResolver{
	constants.PkgMgrWinget: resolveWingetPackage,
	constants.PkgMgrApt:    resolveAptPackage,
	constants.PkgMgrBrew:   resolveBrewPackage,
	constants.PkgMgrSnap:   resolveSnapPackage,
}

// resolvePackageName maps tool name to package ID for a manager.
func resolvePackageName(manager, tool string) string {
	if resolver, isFound := managerResolvers[manager]; isFound {

		return resolver(tool)
	}

	return resolveChocoPackage(tool)
}

var toolAliasMap = map[string]string{
	"k8s":               constants.ToolKubernetes,
	"kubectl":           constants.ToolKubernetes,
	"dotnet-sdk":        constants.ToolDotnet,
	"jdk":               constants.ToolJava,
	"openjdk":           constants.ToolJava,
	"rustup":            constants.ToolRust,
	"cargo":             constants.ToolRust,
	"llamacpp":          constants.ToolLlamaCpp,
	"ngx":               constants.ToolNginx,
	"engine-x":          constants.ToolNginx,
	"wp":                constants.ToolWordPress,
	"wp-cli":            constants.ToolWordPress,
	"wpcli":             constants.ToolWordPress,
	"artisan":           constants.ToolLaravel,
	"laravel-installer": constants.ToolLaravel,
	"open-vm-tools":       constants.ToolVMware,
	"vmtools":             constants.ToolVMware,
	"vmware-tools":        constants.ToolVMware,
	"vm":                  constants.ToolVMware,
	"agy":                 constants.ToolAntigravity,
	"antigravity-cli":     constants.ToolAntigravity,
	"ag-m":                constants.ToolAgManager,
	"ag-manager":          constants.ToolAgManager,
	"antigravity-manager": constants.ToolAgManager,
	"manager":             constants.ToolAgManager,
	"agy-manager":         constants.ToolAgManager,
	"qtorrent":            constants.ToolQBittorrent,
	"qbittorrent":         constants.ToolQBittorrent,
	"qbit":                constants.ToolQBittorrent,
	"utorrent":            constants.ToolUTorrent,
	"u-torrent":           constants.ToolUTorrent,
	"uttorrent":           constants.ToolUTorrent,
}

// resolveToolAlias normalizes known tool aliases to their canonical tool name.
func resolveToolAlias(tool string) string {
	if canonical, isFound := toolAliasMap[tool]; isFound {

		return canonical
	}

	return tool
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
	constants.ToolNginx:         constants.ChocoPkgNginx,
	constants.ToolWordPress:     constants.ChocoPkgWordPress,
	constants.ToolLaravel:       constants.ChocoPkgLaravel,
	constants.ToolVMware:        constants.ChocoPkgVMware,
	constants.ToolQBittorrent:   constants.ChocoPkgQBittorrent,
	constants.ToolUTorrent:      constants.ChocoPkgUTorrent,
}

// resolveChocoPackage maps tool names to Chocolatey package IDs.
func resolveChocoPackage(tool string) string {
	if pkg, isFound := chocoPackageMap[tool]; isFound {

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
	constants.ToolNginx:         constants.WingetPkgNginx,
	constants.ToolWordPress:     constants.WingetPkgWordPress,
	constants.ToolLaravel:       constants.WingetPkgLaravel,
	constants.ToolVMware:        constants.WingetPkgVMware,
	constants.ToolQBittorrent:   constants.WingetPkgQBittorrent,
	constants.ToolUTorrent:      constants.WingetPkgUTorrent,
}

// resolveWingetPackage maps tool names to Winget package IDs.
func resolveWingetPackage(tool string) string {
	if pkg, isFound := wingetPackageMap[tool]; isFound {

		return pkg
	}

	return tool
}
