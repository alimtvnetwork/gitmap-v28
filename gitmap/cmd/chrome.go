// Package cmd — chrome.go: umbrella dispatcher for the `gitmap chrome`
// command group (backup, restore, diff, export-bookmarks, which).
// Added in v6.69.0. See helptext/chrome.md.
package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

func runChrome(args []string) {
	if len(args) == 0 {
		printChromeUsage()
		os.Exit(2)
	}
	sub, tail := args[0], args[1:]
	switch sub {
	case constants.SubCmdChromeBackup:
		runChromeBackup(tail)
	case constants.SubCmdChromeRestore:
		runChromeRestore(tail)
	case constants.SubCmdChromeDiff:
		runChromeDiff(tail)
	case constants.SubCmdChromeExportBookmrk:
		runChromeExportBookmarks(tail)
	case constants.SubCmdChromeWhich:
		runChromeWhich(tail)
	default:
		fmt.Fprintf(os.Stderr, "chrome: ERROR unknown subcommand %q\n", sub)
		printChromeUsage()
		os.Exit(2)
	}
}

func printChromeUsage() {
	fmt.Fprintln(os.Stderr, "usage: gitmap chrome <backup|restore|diff|export-bookmarks|which> [args]")
}
