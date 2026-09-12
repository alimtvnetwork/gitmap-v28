package cmdinstall

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var aptPackageMap = map[string]string{
	constants.ToolNodeJS:          constants.AptPkgNodeJS,
	constants.ToolPython:          constants.AptPkgPython,
	constants.ToolGo:              constants.AptPkgGo,
	constants.ToolGit:             constants.AptPkgGit,
	constants.ToolGitLFS:          constants.AptPkgGitLFS,
	constants.ToolCPP:             constants.AptPkgCPP,
	constants.ToolPHP:             constants.AptPkgPHP,
	constants.ToolMySQL:           constants.AptPkgMySQL,
	constants.ToolMariaDB:         constants.AptPkgMariaDB,
	constants.ToolPostgreSQL:      constants.AptPkgPostgreSQL,
	constants.ToolSQLite:          constants.AptPkgSQLite,
	constants.ToolMongoDB:         constants.AptPkgMongoDB,
	constants.ToolCouchDB:         constants.AptPkgCouchDB,
	constants.ToolRedis:           constants.AptPkgRedis,
	constants.ToolCassandra:       constants.AptPkgCassandra,
	constants.ToolElasticsearch:   constants.AptPkgElasticsearch,
	constants.ToolChrome:          constants.AptPkgChrome,
	constants.ToolGoogleChrome:    constants.AptPkgChrome,
	constants.ToolRust:            constants.AptPkgRust,
	constants.ToolDotnet:          constants.AptPkgDotnet,
	constants.ToolJava:            constants.AptPkgJava,
	constants.ToolFlutter:         constants.AptPkgFlutter,
	constants.ToolOllama:          constants.AptPkgOllama,
	constants.ToolLlamaCpp:        constants.AptPkgLlamaCpp,
	constants.ToolPythonLibs:      constants.AptPkgPythonLibs,
	constants.ToolDocker:          constants.AptPkgDocker,
	constants.ToolKubernetes:      constants.AptPkgKubernetes,
	constants.ToolJenkins:         constants.AptPkgJenkins,
	constants.ToolZsh:             constants.AptPkgZsh,
	constants.ToolFlameshot:       constants.AptPkgFlameshot,
	constants.ToolConemu:          constants.AptPkgConemu,
	constants.ToolVLC:             constants.AptPkgVLC,
	constants.ToolNginx:           constants.AptPkgNginx,
	constants.ToolWordPress:       constants.AptPkgWordPress,
	constants.ToolLaravel:         constants.AptPkgLaravel,
	constants.ToolVMware:          constants.AptPkgVMware,
	constants.ToolQBittorrent:     constants.AptPkgQBittorrent,
	constants.ToolUTorrent:        constants.AptPkgUTorrent,
	constants.ToolOBS:             constants.AptPkgOBS,
	constants.ToolDbeaver:         constants.AptPkgDbeaver,
	constants.ToolUbuntuFont:      constants.AptPkgUbuntuFont,
	constants.ToolWhatsApp:        constants.AptPkgWhatsApp,
	constants.ToolOneNote:         constants.AptPkgOneNote,
	constants.ToolLightshot:       constants.AptPkgLightshot,
	constants.ToolWindowsTerminal: constants.AptPkgWindowsTerminal,
	constants.ToolAria2:           constants.AptPkgAria2,
	constants.Tool7Zip:            constants.AptPkg7Zip,
	constants.ToolWinRAR:          constants.AptPkgWinRAR,
	constants.ToolXMind:           constants.AptPkgXMind,
	constants.ToolWordWeb:         constants.AptPkgWordWeb,
	constants.ToolBeyondCompare:   constants.AptPkgBeyondCompare,
	constants.ToolVcRedist:        constants.AptPkgVcRedist,
	constants.ToolDirectX:         constants.AptPkgDirectX,
	constants.ToolDirectXSdk:      constants.AptPkgDirectXSdk,
	constants.ToolStarship:        constants.AptPkgStarship,
	constants.ToolOhMyPosh:        constants.AptPkgOhMyPosh,
	constants.ToolScoop:           constants.AptPkgScoop,
}

// resolveAptPackage maps tool names to apt package IDs.
func resolveAptPackage(tool string) string {
	if pkg, isFound := aptPackageMap[tool]; isFound {
		return pkg
	}

	return tool
}

var brewPackageMap = map[string]string{
	constants.ToolNodeJS:          constants.BrewPkgNodeJS,
	constants.ToolPython:          constants.BrewPkgPython,
	constants.ToolGo:              constants.BrewPkgGo,
	constants.ToolGit:             constants.BrewPkgGit,
	constants.ToolGitLFS:          constants.BrewPkgGitLFS,
	constants.ToolGHCLI:           constants.BrewPkgGHCLI,
	constants.ToolCPP:             constants.BrewPkgCPP,
	constants.ToolPHP:             constants.BrewPkgPHP,
	constants.ToolMySQL:           constants.BrewPkgMySQL,
	constants.ToolMariaDB:         constants.BrewPkgMariaDB,
	constants.ToolPostgreSQL:      constants.BrewPkgPostgreSQL,
	constants.ToolSQLite:          constants.BrewPkgSQLite,
	constants.ToolMongoDB:         constants.BrewPkgMongoDB,
	constants.ToolCouchDB:         constants.BrewPkgCouchDB,
	constants.ToolRedis:           constants.BrewPkgRedis,
	constants.ToolNeo4j:           constants.BrewPkgNeo4j,
	constants.ToolElasticsearch:   constants.BrewPkgElasticsearch,
	constants.ToolDuckDB:          constants.BrewPkgDuckDB,
	constants.ToolDbeaver:         constants.BrewPkgDbeaver,
	constants.ToolOBS:             constants.BrewPkgOBS,
	constants.ToolChrome:          constants.BrewPkgChrome,
	constants.ToolGoogleChrome:    constants.BrewPkgChrome,
	constants.ToolRust:            constants.BrewPkgRust,
	constants.ToolDotnet:          constants.BrewPkgDotnet,
	constants.ToolJava:            constants.BrewPkgJava,
	constants.ToolFlutter:         constants.BrewPkgFlutter,
	constants.ToolOllama:          constants.BrewPkgOllama,
	constants.ToolLlamaCpp:        constants.BrewPkgLlamaCpp,
	constants.ToolPythonLibs:      constants.BrewPkgPythonLibs,
	constants.ToolDocker:          constants.BrewPkgDocker,
	constants.ToolKubernetes:      constants.BrewPkgKubernetes,
	constants.ToolJenkins:         constants.BrewPkgJenkins,
	constants.ToolZsh:             constants.BrewPkgZsh,
	constants.ToolFlameshot:       constants.BrewPkgFlameshot,
	constants.ToolConemu:          constants.BrewPkgConemu,
	constants.ToolVLC:             constants.BrewPkgVLC,
	constants.ToolNginx:           constants.BrewPkgNginx,
	constants.ToolWordPress:       constants.BrewPkgWordPress,
	constants.ToolLaravel:         constants.BrewPkgLaravel,
	constants.ToolVMware:          constants.BrewPkgVMware,
	constants.ToolQBittorrent:     constants.BrewPkgQBittorrent,
	constants.ToolUTorrent:        constants.BrewPkgUTorrent,
	constants.ToolUbuntuFont:      constants.BrewPkgUbuntuFont,
	constants.ToolWhatsApp:        constants.BrewPkgWhatsApp,
	constants.ToolOneNote:         constants.BrewPkgOneNote,
	constants.ToolLightshot:       constants.BrewPkgLightshot,
	constants.ToolWindowsTerminal: constants.BrewPkgWindowsTerminal,
	constants.ToolAria2:           constants.BrewPkgAria2,
	constants.Tool7Zip:            constants.BrewPkg7Zip,
	constants.ToolWinRAR:          constants.BrewPkgWinRAR,
	constants.ToolXMind:           constants.BrewPkgXMind,
	constants.ToolWordWeb:         constants.BrewPkgWordWeb,
	constants.ToolBeyondCompare:   constants.BrewPkgBeyondCompare,
	constants.ToolVcRedist:        constants.BrewPkgVcRedist,
	constants.ToolDirectX:         constants.BrewPkgDirectX,
	constants.ToolDirectXSdk:      constants.BrewPkgDirectXSdk,
	constants.ToolStarship:        constants.BrewPkgStarship,
	constants.ToolOhMyPosh:        constants.BrewPkgOhMyPosh,
	constants.ToolScoop:           constants.BrewPkgScoop,
}

// resolveBrewPackage maps tool names to Homebrew package IDs.
func resolveBrewPackage(tool string) string {
	if pkg, isFound := brewPackageMap[tool]; isFound {
		return pkg
	}

	return tool
}

var snapPackageMap = map[string]string{
	constants.ToolCouchDB:         constants.SnapPkgCouchDB,
	constants.ToolRedis:           constants.SnapPkgRedis,
	constants.ToolVSCode:          "code",
	constants.ToolFlutter:         constants.SnapPkgFlutter,
	constants.ToolKubernetes:      constants.SnapPkgKubernetes,
	constants.ToolVLC:             constants.SnapPkgVLC,
	constants.ToolDotnet:          constants.SnapPkgDotnet,
	constants.ToolRust:            constants.SnapPkgRust,
	constants.ToolOBS:             constants.SnapPkgOBS,
	constants.ToolDbeaver:         constants.SnapPkgDbeaver,
	constants.ToolUbuntuFont:      constants.SnapPkgUbuntuFont,
	constants.ToolWhatsApp:        constants.SnapPkgWhatsApp,
	constants.ToolOneNote:         constants.SnapPkgOneNote,
	constants.ToolLightshot:       constants.SnapPkgLightshot,
	constants.ToolWindowsTerminal: constants.SnapPkgWindowsTerminal,
	constants.ToolAria2:           constants.SnapPkgAria2,
	constants.Tool7Zip:            constants.SnapPkg7Zip,
	constants.ToolWinRAR:          constants.SnapPkgWinRAR,
	constants.ToolXMind:           constants.SnapPkgXMind,
	constants.ToolWordWeb:         constants.SnapPkgWordWeb,
	constants.ToolBeyondCompare:   constants.SnapPkgBeyondCompare,
	constants.ToolVcRedist:        constants.SnapPkgVcRedist,
	constants.ToolDirectX:         constants.SnapPkgDirectX,
	constants.ToolDirectXSdk:      constants.SnapPkgDirectXSdk,
	constants.ToolStarship:        constants.SnapPkgStarship,
	constants.ToolOhMyPosh:        constants.SnapPkgOhMyPosh,
	constants.ToolScoop:           constants.SnapPkgScoop,
}

// resolveSnapPackage maps tool names to Snap package IDs.
func resolveSnapPackage(tool string) string {
	if pkg, isFound := snapPackageMap[tool]; isFound {
		return pkg
	}

	return tool
}
