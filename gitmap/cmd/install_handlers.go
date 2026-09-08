package cmd

import (
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func specialSyncHandler(tool string) func(installOptions) {
	return map[string]func(installOptions){
		constants.ToolScripts:     func(installOptions) { runInstallScripts() },
		constants.ToolNppSettings: func(installOptions) { runNppSettingsOnly() },
		constants.ToolVSCodeSync:  func(installOptions) { runVSCodeSettingsOnly() },
		constants.ToolOBSSync:     func(installOptions) { runOBSSettingsOnly() },
		constants.ToolWTSync:      func(installOptions) { runWTSettingsOnly() },
	}[tool]
}

func specialToolHandler(tool string) func(installOptions) {
	return map[string]func(installOptions){
		constants.ToolVSCodeCtx:        func(installOptions) { runVSCodeContextMenu() },
		constants.ToolPwshCtx:          func(installOptions) { runPwshContextMenu() },
		constants.ToolCtx:              func(opts installOptions) { runInstallCtx(opts.Explain) },
		constants.ToolAllDevTools:      func(opts installOptions) { runAllDevTools(opts) },
		constants.ToolGitmapOneliner:   func(installOptions) { runInstallGitmapOneliner() },
		constants.ToolScriptsFixer:     func(installOptions) { runInstallCustomTool("scripts-fixer") },
		constants.ToolCodingGuidelines: func(installOptions) { runInstallCustomTool("coding-guidelines") },
		constants.ToolMacroAhk:         func(installOptions) { runInstallCustomTool("macro-ahk") },
		constants.ToolAgManager:        func(installOptions) { runInstallAgManager() },
		constants.ToolAgCtx:            func(opts installOptions) { runInstallCtx(opts.Explain) },
		constants.ToolAntigravity:      func(opts installOptions) { runInstallAntigravity() },
	}[tool]
}

func specialLinuxHandler(tool string) func(installOptions) {
	if tool == constants.ToolGitHubDesktop && runtime.GOOS == "linux" {

		return func(opts installOptions) { runInstallGitHubDesktopLinux(opts) }
	}
	if tool == constants.ToolVSCode && runtime.GOOS == "linux" {

		return func(opts installOptions) { runInstallVSCodeLinux(opts) }
	}
	if (tool == constants.ToolChrome || tool == constants.ToolGoogleChrome) && runtime.GOOS == "linux" {

		return func(opts installOptions) { _ = runInstallChromeLinux(opts) }
	}

	return nil
}

func specialAliasHandler(tool string) func(installOptions) {
	if isCleanCodeAlias(tool) {

		return func(installOptions) { runInstallCleanCode() }
	}
	if isBuildEssentialAlias(tool) {

		return func(opts installOptions) { _ = runInstallBuildEssential(opts) }
	}

	return nil
}

func specialInstallHandler(tool string) func(installOptions) {
	if h := specialAliasHandler(tool); h != nil {

		return h
	}
	if h := specialLinuxHandler(tool); h != nil {

		return h
	}
	if h := specialSyncHandler(tool); h != nil {

		return h
	}

	return specialToolHandler(tool)
}
