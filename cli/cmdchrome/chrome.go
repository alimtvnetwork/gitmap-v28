// Package cmd — chrome.go: umbrella dispatcher for `gitmap chrome` (and alias `gitmap cprof`).
package cmdchrome

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdchromeprofile"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func runChrome(args []string) error {
	if len(args) == 0 || isHelpFlag(args[0]) {
		printChromeRichUsage()

		return nil
	}

	if args[0] == "profile" || args[0] == "profiles" {
		args = args[1:]
	}

	if len(args) == 0 || isHelpFlag(args[0]) {
		printChromeRichUsage()

		return nil
	}

	sub, tail := args[0], args[1:]
	if isHandled := dispatchChromeSubcommand(sub, tail); isHandled {
		return nil
	}

	fmt.Fprintf(os.Stderr, "chrome: ERROR unknown subcommand %q\n", sub)
	printChromeUsage()
	cliexit.HandleError(nil, 2)

	return nil
}

func isHelpFlag(arg string) bool {
	return arg == "help" || arg == "--help" || arg == "-h"
}

func dispatchChromeSubcommand(sub string, tail []string) bool {
	if handleChromeLaunchOps(sub, tail) || handleChromeExtensionOps(sub, tail) {
		return true
	}

	if handleChromeFlagAndResetOps(sub, tail) || handleChromeInstallOps(sub, tail) {
		return true
	}

	if handleChromeProfileOps(sub, tail) || handleChromeBatchOps(sub, tail) {
		return true
	}

	return handleChromeArchiveOps(sub, tail)
}

func handleChromeLaunchOps(sub string, tail []string) bool {
	switch sub {
	case constants.SubCmdChromeOpen, constants.SubCmdChromeOpenAlias, constants.SubCmdChromeOpenAlias2:
		_ = runChromeOpen(tail)

		return true
	case constants.SubCmdChromeObserve, constants.SubCmdChromeObserveAlias, constants.SubCmdChromeObserveAlias2, constants.SubCmdChromeObserveAlias3:
		_ = runChromeObserve(tail)

		return true
	}

	return false
}

func handleChromeExtensionOps(sub string, tail []string) bool {
	switch sub {
	case constants.SubCmdChromeExtensions, constants.SubCmdChromeExtensionsAlias, constants.SubCmdChromeExtensionsAlias2:
		_ = runChromeExtensions(tail)

		return true
	case constants.SubCmdChromeExtInstall, constants.SubCmdChromeExtInstallAlias, constants.SubCmdChromeExtInstallAlias2:
		_ = runChromeExtensionInstall(tail)

		return true
	case constants.SubCmdChromeExtEnable, constants.SubCmdChromeExtEnableAlias, constants.SubCmdChromeExtEnableAlias2:
		_ = runChromeExtensionEnable(tail)

		return true
	case constants.SubCmdChromeExtDisable, constants.SubCmdChromeExtDisableAlias, constants.SubCmdChromeExtDisableAlias2:
		_ = runChromeExtensionDisable(tail)

		return true
	case constants.SubCmdChromeExtDisableAll, constants.SubCmdChromeExtDisableAllAlias:
		_ = runChromeExtensionDisableAll(tail)

		return true
	}

	return false
}

func handleChromeFlagAndResetOps(sub string, tail []string) bool {
	switch sub {
	case constants.SubCmdChromeFlags, constants.SubCmdChromeFlagsAlias:
		_ = runChromeFlags(tail)

		return true
	case constants.SubCmdChromeReset, constants.SubCmdChromeResetAlias:
		_ = runChromeReset(tail)

		return true
	}

	return false
}

func handleChromeInstallOps(sub string, tail []string) bool {
	if sub != constants.SubCmdChromeInstall && sub != constants.SubCmdChromeInstallAlias {
		return false
	}

	isDryRun := hasDryRunFlag(tail)
	if constants.IsLinuxPlatform() {
		runLinuxChromeInstall(isDryRun)

		return true
	}

	if InstallToolFn != nil {
		InstallToolFn(constants.ToolChrome, isDryRun)
	}

	return true
}

func runLinuxChromeInstall(isDryRun bool) {
	if InstallChromeLinuxFn != nil {
		_ = InstallChromeLinuxFn(isDryRun)
	}
}

func hasDryRunFlag(args []string) bool {
	for _, a := range args {
		if a == "--dry-run" {
			return true
		}
	}

	return false
}

func handleChromeProfileOps(sub string, tail []string) bool {
	switch sub {
	case constants.SubCmdChromeCopy, constants.SubCmdChromeCopyAlias, constants.SubCmdChromeCopyAlias2:
		_ = cmdchromeprofile.RunProfileCopy(tail)

		return true
	case constants.SubCmdChromeExport, constants.SubCmdChromeExportAlias, constants.SubCmdChromeExportAlias2:
		_ = cmdchromeprofile.RunProfileExport(tail)

		return true
	case constants.SubCmdChromeImport, constants.SubCmdChromeImportAlias, constants.SubCmdChromeImportAlias2:
		_ = cmdchromeprofile.RunProfileImport(tail)

		return true
	case constants.SubCmdChromeImportCheck, constants.SubCmdChromeImportCheckAlias, constants.SubCmdChromeImportCheckAlias2:
		_ = cmdchromeprofile.RunProfileImportCheck(tail)

		return true
	case constants.SubCmdChromeList, constants.SubCmdChromeListAlias, constants.SubCmdChromeListAlias2, constants.SubCmdChromeListAlias3:
		_ = cmdchromeprofile.RunProfileList(tail)

		return true
	case constants.SubCmdChromeDelete, constants.SubCmdChromeDeleteAlias, constants.SubCmdChromeDeleteAlias2, constants.SubCmdChromeDeleteAlias3, "clear":
		_ = cmdchromeprofile.RunProfileClear(tail)

		return true
	case "optimize-projects", "optimize", "--repeat-fix", "-r", "dedupe":
		_ = cmdchromeprofile.RunProfileOptimize(tail)

		return true
	case "find-duplicates", "duplicates", "dups", "find-dups":
		if RunFindDuplicatesFn != nil {
			_ = RunFindDuplicatesFn("chrome", tail)
		}

		return true
	case constants.SubCmdChromeMerge, constants.SubCmdChromeMergeAlias:
		_ = cmdchromeprofile.RunProfileMerge(tail)

		return true
	case constants.SubCmdChromeReconcile, constants.SubCmdChromeReconcileAlias, constants.SubCmdChromeReconcileAlias2, constants.SubCmdChromeReconcileAlias3:
		_ = cmdchromeprofile.RunProfileReconcile(tail)

		return true
	case "undo":
		profile := "Default"
		if len(tail) > 0 {
			profile = tail[0]
		}

		_ = cmdchromeprofile.RunProfileUndo(profile)

		return true
	case "redo":
		profile := "Default"
		if len(tail) > 0 {
			profile = tail[0]
		}

		_ = cmdchromeprofile.RunProfileRedo(profile)

		return true
	case "group", "grp", "groups":
		_ = cmdchromeprofile.RunProfileGroupDispatch(tail)

		return true
	}

	return false
}

func handleChromeBatchOps(sub string, tail []string) bool {
	switch sub {
	case constants.SubCmdChromeCopyAll, constants.SubCmdChromeCopyAllAlias, constants.SubCmdChromeCopyAllAlias2, constants.SubCmdChromeCopyAllAlias3:
		_ = cmdchromeprofile.RunCopyAll(tail)

		return true
	case constants.SubCmdChromeExportAll, constants.SubCmdChromeExportAllAlias, constants.SubCmdChromeExportAllAlias2, constants.SubCmdChromeExportAllAlias3:
		_ = cmdchromeprofile.RunExportAll(tail)

		return true
	case constants.SubCmdChromeImportAll, constants.SubCmdChromeImportAllAlias, constants.SubCmdChromeImportAllAlias2, constants.SubCmdChromeImportAllAlias3:
		_ = cmdchromeprofile.RunImportAll(tail)

		return true
	}

	return false
}

func handleChromeArchiveOps(sub string, tail []string) bool {
	switch sub {
	case constants.SubCmdChromeBackup:
		runChromeBackup(tail)

		return true
	case constants.SubCmdChromeRestore:
		runChromeRestore(tail)

		return true
	case constants.SubCmdChromeDiff:
		runChromeDiff(tail)

		return true
	case constants.SubCmdChromeExportBookmrk, constants.SubCmdChromeBookmarks:
		runChromeExportBookmarks(tail)

		return true
	case constants.SubCmdChromeWhich:
		runChromeWhich(tail)

		return true
	}

	return false
}
