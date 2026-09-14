package cmdinstall

import (
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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
	return resolveSpecialToolsMap()[tool]
}

func resolveSpecialToolsMap() map[string]func(installOptions) {
	return map[string]func(installOptions){
		constants.ToolVSCodeCtx:        func(installOptions) { runVSCodeContextMenu() },
		constants.ToolPwshCtx:          func(installOptions) { runPwshContextMenu() },
		constants.ToolCtx:              func(opts installOptions) { runInstallCtx(opts.Explain) },
		constants.ToolAllDevTools:      func(opts installOptions) { runAllDevTools(opts) },
		constants.ToolGitmapOneliner:   func(installOptions) { runInstallGitmapOneliner() },
		constants.ToolScriptsFixer:     func(installOptions) { handleCustomToolInstall("scripts-fixer") },
		constants.ToolCodingGuidelines: func(installOptions) { handleCustomToolInstall("coding-guidelines") },
		constants.ToolMacroAhk:         func(installOptions) { handleCustomToolInstall("macro-ahk") },
		constants.ToolAgManager:        func(opts installOptions) { handleAgManagerInstall(opts) },
		constants.ToolAgCtx:            func(opts installOptions) { runInstallCtx(opts.Explain) },
		constants.ToolAntigravity:      func(opts installOptions) { handleAntigravityInstall(opts) },
		constants.ToolAgy:              func(opts installOptions) { handleAgyInstall(opts) },
		constants.ToolGitCompact:       func(opts installOptions) { handleGitCompactInstall(opts) },
	}
}

func handleAntigravityInstall(opts installOptions) {
	if err := runInstallAntigravityWithOpts(opts); err != nil {
		cliexit.HandleError(apperror.WrapSimple(err, "cmdinstall.installAntigravity"), 1)
	}
}

func handleAgyInstall(opts installOptions) {
	if err := runInstallAgyWithOpts(opts); err != nil {
		cliexit.HandleError(apperror.WrapSimple(err, "cmdinstall.installAgy"), 1)
	}
}

func handleAgManagerInstall(opts installOptions) {
	if err := runInstallAgManagerWithOpts(opts); err != nil {
		cliexit.HandleError(apperror.WrapSimple(err, "cmdinstall.installAgManager"), 1)
	}
}

func handleCustomToolInstall(tool string) {
	if err := runInstallCustomTool(tool); err != nil {
		cliexit.HandleError(apperror.WrapSimple(err, "cmdinstall.installCustomTool"), 1)
	}
}

func specialLinuxHandler(tool string) func(installOptions) {
	if runtime.GOOS != "linux" {
		return nil
	}

	return resolveLinuxToolHandler(tool)
}

func resolveLinuxToolHandler(tool string) func(installOptions) {
	if h := resolveLinuxAppHandler(tool); h != nil {
		return h
	}

	return resolveLinuxJsHandler(tool)
}

func resolveLinuxAppHandler(tool string) func(installOptions) {
	switch tool {
	case constants.ToolGitHubDesktop:

		return func(opts installOptions) { runInstallGitHubDesktopLinux(opts) }
	case constants.ToolVSCode:

		return func(opts installOptions) { runInstallVSCodeLinux(opts) }
	case constants.ToolChrome, constants.ToolGoogleChrome:

		return func(opts installOptions) { _ = runInstallChromeLinux(opts) }
	default:

		return nil
	}
}

func resolveLinuxJsHandler(tool string) func(installOptions) {
	switch tool {
	case constants.ToolPnpm:

		return func(opts installOptions) { _ = runInstallPnpmLinux(opts) }
	case constants.ToolYarn:

		return func(opts installOptions) { _ = runInstallYarnLinux(opts) }
	case constants.ToolBun:

		return func(opts installOptions) { _ = runInstallBunLinux(opts) }
	default:

		return nil
	}
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
