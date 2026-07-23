// Package cmd — visibilityundoflags.go: flag parsing + dry-run
// rendering for `gitmap visibility-undo` / `visibility-redo`.
// Extracted from visibilityundo.go to honor the 200-line per-file cap.
//
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §undo-redo.
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// parseVisUndoArgs accepts --verbose, --dry-run, --force, --json, and --run <id>.
// Unknown tokens are tolerated (mirrors parseBulkArgs).
func parseVisUndoArgs(args []string) undoFlags {
	flags := undoFlags{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--verbose":
			flags.Verbose = true
		case "--dry-run":
			flags.DryRun = true
		case "--force":
			flags.Force = true
		case "--json":
			flags.JSON = true
		case "--run":
			flags.RunID = mustParseRunID(args, i)
			i++
		}
	}

	return flags
}

// mustParseRunID validates the `--run <id>` pairing and exits on bad
// input (zero-swallow — a typo here would silently undo the wrong run).
func mustParseRunID(args []string, i int) int64 {
	if i+1 >= len(args) {
		fmt.Fprintf(os.Stderr, constants.ErrUndoBadRunFlagFmt, "", fmt.Errorf("missing value"), "no value after --run")
		os.Exit(constants.ExitVisBadFlag)
	}
	raw := args[i+1]
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		fmt.Fprintf(os.Stderr, constants.ErrUndoBadRunFlagFmt, raw, err, "must be positive integer")
		os.Exit(constants.ExitVisBadFlag)
	}

	return id
}

// printVisDryRun lists the planned per-repo reversals without mutating.
func printVisDryRun(cmdName string, run model.MakeAllVisibilityRunRecord, rs []model.MakeAllVisibilityResultRecord) {
	fmt.Fprintf(os.Stdout, constants.MsgDryRunHeaderFmt, cmdName, run.ID, run.Provider, run.Owner, len(rs))
	total := len(rs)
	for i, r := range rs {
		fmt.Fprintf(os.Stdout, constants.MsgDryRunRowFmt, i+1, total, r.RepoName, r.PrevVisibility)
	}
	fmt.Fprintf(os.Stdout, constants.MsgDryRunFooterFmt, cmdName)
}
