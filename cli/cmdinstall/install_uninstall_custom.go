package cmdinstall

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// IsCustomStandaloneTool reports whether tool is installed outside package managers.
func IsCustomStandaloneTool(tool string) bool {
	canonical := resolveToolAlias(tool)
	switch canonical {
	case constants.ToolAgy, constants.ToolAntigravity,
		constants.ToolAgManager, constants.ToolScriptsFixer,
		constants.ToolCodingGuidelines, constants.ToolMacroAhk,
		constants.ToolCtx, constants.ToolVSCodeCtx,
		constants.ToolPwshCtx, constants.ToolAgCtx,
		constants.ToolScripts:
		return true
	default:
		return false
	}
}

func dispatchCustomRemoval(canonical string, purge bool) {
	switch canonical {
	case constants.ToolAgy:
		uninstallAgyFiles(purge)
	case constants.ToolAntigravity:
		uninstallAntigravityFiles(purge)
	case constants.ToolAgManager:
		uninstallAgManagerFiles()
	case constants.ToolCtx, constants.ToolVSCodeCtx, constants.ToolPwshCtx, constants.ToolAgCtx:
		runUninstallCtx()
	case constants.ToolScripts:
		uninstallScriptsDirectory()
	default:
		uninstallScriptToolFiles(canonical)
	}
}

// RunUninstallCustomTool removes files for custom tools and purges DB entries.
func RunUninstallCustomTool(tool string, purge bool) error {
	canonical := resolveToolAlias(tool)
	fmt.Printf("Removing %s...\n", canonical)

	dispatchCustomRemoval(canonical, purge)

	cleanToolFromDatabases(canonical)
	if tool != canonical {
		cleanToolFromDatabases(tool)
	}

	fmt.Printf(constants.MsgUninstallSuccess, tool)
	return nil
}

func removeFileIfExists(path string) {
	if path == "" {
		return
	}
	_ = os.Remove(path)
}

func removeDirIfExists(path string) {
	if path == "" {
		return
	}
	_ = os.RemoveAll(path)
}

func cleanToolFromDatabases(name string) {
	splitDB, err := store.OpenInstallationSplitDB()
	if err == nil {
		defer splitDB.Close()
		_ = splitDB.RemoveInstalledTool(name)
	}
}
