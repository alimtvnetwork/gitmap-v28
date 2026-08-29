// Package cmd — chromeprofile_delete.go: implements `gitmap
// chrome-profile-delete` (alias `cpd`). Removes the SQLite row and
// optionally rm()s the on-disk artifacts that were tracked alongside
// it. Refuses to run without --yes to avoid accidental destruction.
package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
)

// runChromeProfileDelete is the dispatch entry point.
func runChromeProfileDelete(args []string) error {
	checkHelp(constants.CmdChromeProfileDelete, args)
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrChromeProfileUsageDelete)
		cliexit.HandleError(nil, constants.ExitChromeProfileUsage)
	}
	name, confirmed := parseChromeDeleteArgs(args)
	if !confirmed {
		fmt.Fprint(os.Stderr, constants.MsgChromeProfileDeleteAbort)
		cliexit.HandleError(nil, constants.ExitChromeProfileUsage)
	}
	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrChromeProfileDeleteFail, err)
		cliexit.HandleError(nil, constants.ExitChromeProfileCopyFailed)
	}
	defer db.Close()
	paths, err := db.DeleteChromeProfile(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrChromeProfileDeleteFail, err)
		cliexit.HandleError(nil, constants.ExitChromeProfileCopyFailed)
	}
	if len(paths) == 0 && !db.ChromeProfileExists(name) {
		fmt.Fprintf(os.Stderr, constants.ErrChromeProfileNotInDB, name)
		cliexit.HandleError(nil, constants.ExitChromeProfileNotFound)
	}
	removed := removeChromeArtifactFiles(paths)
	fmt.Printf(constants.MsgChromeProfileDeleteOk, name, removed)
	return nil
}

// parseChromeDeleteArgs extracts <name> and the --yes confirmation flag.
func parseChromeDeleteArgs(args []string) (string, bool) {
	name := ""
	confirmed := false
	for _, a := range args {
		if a == "--yes" || a == "-y" {
			confirmed = true
			continue
		}
		if name == "" {
			name = a
		}
	}
	return name, confirmed && name != ""
}

// removeChromeArtifactFiles best-effort rm()s each artifact path. Missing
// files are silently skipped so re-runs stay idempotent.
func removeChromeArtifactFiles(paths []string) int {
	count := 0
	for _, p := range paths {
		if p == "" {
			continue
		}
		fmt.Printf(constants.MsgChromeProfileDeleteRm, p)
		if err := os.Remove(p); err == nil {
			count++
		}
	}
	return count
}
