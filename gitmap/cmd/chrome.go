// Package cmd — chrome.go: umbrella dispatcher for `gitmap chrome` (and alias `gitmap cp`).
package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runChrome dispatches `gitmap chrome <subcommand>`.
func runChrome(args []string) error {
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
	if handleChromeInstallOps(sub, tail) {
		return true
	}
	if handleChromeProfileOps(sub, tail) {
		return true
	}
	if handleChromeBatchOps(sub, tail) {
		return true
	}
	return handleChromeArchiveOps(sub, tail)
}

func handleChromeInstallOps(sub string, tail []string) bool {
	if sub == constants.SubCmdChromeInstall || sub == constants.SubCmdChromeInstallAlias {
		installTool(installOptions{Tool: constants.ToolChrome, Manager: "", DryRun: hasDryRunFlag(tail)})
		return true
	}
	return false
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
		_ = runChromeProfileCopy(tail)
		return true
	case constants.SubCmdChromeExport, constants.SubCmdChromeExportAlias, constants.SubCmdChromeExportAlias2:
		_ = runChromeProfileExport(tail)
		return true
	case constants.SubCmdChromeImport, constants.SubCmdChromeImportAlias, constants.SubCmdChromeImportAlias2:
		_ = runChromeProfileImport(tail)
		return true
	case constants.SubCmdChromeList, constants.SubCmdChromeListAlias, constants.SubCmdChromeListAlias2, constants.SubCmdChromeListAlias3:
		_ = runChromeProfileList(tail)
		return true
	case constants.SubCmdChromeDelete, constants.SubCmdChromeDeleteAlias, constants.SubCmdChromeDeleteAlias2, constants.SubCmdChromeDeleteAlias3:
		_ = runChromeProfileDelete(tail)
		return true
	case constants.SubCmdChromeMerge, constants.SubCmdChromeMergeAlias:
		_ = runChromeProfileMerge(tail)
		return true
	}
	return false
}

func handleChromeBatchOps(sub string, tail []string) bool {
	switch sub {
	case constants.SubCmdChromeCopyAll, constants.SubCmdChromeCopyAllAlias, constants.SubCmdChromeCopyAllAlias2, constants.SubCmdChromeCopyAllAlias3:
		_ = runChromeCopyAll(tail)
		return true
	case constants.SubCmdChromeExportAll, constants.SubCmdChromeExportAllAlias, constants.SubCmdChromeExportAllAlias2, constants.SubCmdChromeExportAllAlias3:
		_ = runChromeExportAll(tail)
		return true
	case constants.SubCmdChromeImportAll, constants.SubCmdChromeImportAllAlias, constants.SubCmdChromeImportAllAlias2, constants.SubCmdChromeImportAllAlias3:
		_ = runChromeImportAll(tail)
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

func printChromeUsage() {
	fmt.Fprintln(os.Stderr, "usage: gitmap chrome <install|copy|copy-all|export|export-all|import|import-all|list|delete|merge|backup|restore|diff|which> [args]")
}

func printChromeRichUsage() {
	fmt.Printf("\n\033[1;96mGitMap Chrome Management\033[0m\n\n")
	fmt.Printf("Usage: gitmap chrome <command> [arguments]\n")
	fmt.Printf("Alias: gitmap cprof <command> [arguments]\n\n")
	printChromeInstallationTable()
	printChromeProfileTable()
	printChromeBatchTable()
	printChromeMaintenanceTable()
	fmt.Println()
}

func printChromeInstallationTable() {
	fmt.Printf("\033[1;94mInstallation & Setup:\033[0m\n")
	fmt.Printf("  install (in)                        Install Chrome via system package manager\n\n")
}

func printChromeProfileTable() {
	fmt.Printf("\033[1;94mSingle Profile Operations:\033[0m\n")
	fmt.Printf("  copy (cpc) <src> <dst>              Copy a Chrome profile into an offline profile\n")
	fmt.Printf("  export (cpe) <name> [out]           Export profile snapshot (json, zip, sqlite)\n")
	fmt.Printf("  import (cpi) <file> [name]          Import profile from snapshot file\n")
	fmt.Printf("  list (ls, cpl)                      List Chrome profiles known to gitmap\n")
	fmt.Printf("  delete (rm, cpd) <name>             Remove a profile & stored artifacts\n")
	fmt.Printf("  merge (cpm) <src> <dst>             Merge selected pieces of one profile into another\n\n")
}

func printChromeBatchTable() {
	fmt.Printf("\033[1;94mBatch Profile Operations:\033[0m\n")
	fmt.Printf("  copy-all (cpc-all) <dst-dir>        Copy ALL discovered Chrome profiles to a directory\n")
	fmt.Printf("  export-all (cpe-all) [out-dir]      Export ALL Chrome profiles in batch\n")
	fmt.Printf("  import-all (cpi-all) <dir>          Import ALL profile snapshots from a directory\n\n")
}

func printChromeMaintenanceTable() {
	fmt.Printf("\033[1;94mBackup, Diff & Inspection:\033[0m\n")
	fmt.Printf("  backup [out]                        Snapshot all Chrome profiles into a tar.gz archive\n")
	fmt.Printf("  restore <tarball>                   Restore Chrome profiles from tar.gz archive\n")
	fmt.Printf("  diff <prof1> <prof2>                Compare two Chrome profiles\n")
	fmt.Printf("  bookmarks (export-bookmarks) <name> Export bookmarks to HTML or JSON\n")
	fmt.Printf("  which                               Print Chrome User Data & binary paths\n")
}
