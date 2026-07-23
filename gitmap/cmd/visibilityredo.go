// Package cmd — visibilityredo.go: `gitmap visibility-redo` (`vr`)
// reverses the most recent `VisibilityUndo` run, restoring the
// visibility state that the undo reverted. Pure reuse of the
// shared reverseRunAndExit helper from visibilityundo.go.
//
// Accepts `--run <id>` and `--dry-run`.
//
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §undo-redo.
package cmd

import (
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runVisibilityRedo is the dispatcher entry point.
func runVisibilityRedo(args []string) {
	flags := parseVisUndoArgs(args)
	run, results := loadReversible(flags.RunID, constants.CommandKindVisibilityUndo, constants.ErrRedoNoRunFound)
	if flags.DryRun {
		printVisDryRun(constants.CmdVisibilityRedo, run, results)
		os.Exit(constants.ExitVisOK)
	}
	reverseRunAndExit(run, results, flags, constants.CmdVisibilityRedo)
}
