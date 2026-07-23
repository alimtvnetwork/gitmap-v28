package clonepick

// execute.go: the side-effecting branch of clone-pick. Runs the
// sparse-checkout pipeline, persists the selection, and returns a
// Result the cmd layer can translate into an exit code.

import (
	"fmt"
	"io"
	"os"
)

// Execute runs the sparse-checkout pipeline for plan and persists
// the selection (when not dry-run) via p. Returns a Result whose
// Status drives the cmd-layer exit code.
//
// progress is where per-step git output is streamed (typically
// os.Stderr; io.Discard when --quiet was passed).
func Execute(plan Plan, p Persister, progress io.Writer) Result {
	dest, err := runSparseCheckout(plan, progress)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		// Promote either consumed PreClonedSrc (rename) or left
		// it intact (failed before promote). Best-effort cleanup
		// covers both -- a missing dir is not an error.
		if len(plan.PreClonedSrc) > 0 {
			os.RemoveAll(plan.PreClonedSrc)
		}

		return Result{Status: StatusFailed, Detail: err.Error()}
	}

	id, saveErr := SaveSelection(p, plan)
	if saveErr != nil {
		// Clone succeeded -- treat persistence failure as a soft
		// warning, not a clone failure. The user has the files.
		fmt.Fprintln(os.Stderr, saveErr)
	}

	return Result{
		Status:      StatusOK,
		SelectionId: id,
		Detail:      dest,
	}
}
