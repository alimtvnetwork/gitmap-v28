package orchestrator

import (
	"errors"
	"fmt"
	"io"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/finalize"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/replay"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/walk"
)

// errConflictAborted is the sentinel recordFail receives when the
// user aborted under Prompt mode. Kept package-private so callers
// outside the orchestrator cannot fabricate the abort signal.
var errConflictAborted = errors.New("conflict aborted by user")

// conflictCheck inspects the planned replay for HEAD-vs-source blob
// clobbers and consults the resolved ConflictMode. Returns:
//   - shouldAbortRun=true  → caller must propagate ConflictAborted
//     exit code AND stop the entire run (no further commits).
//   - shouldSkipCommit=true → caller should record this commit as a
//     Failed/skipped row and continue with the next one.
//   - both false           → no clobber (or ForceMerge) → proceed.
//
// Errors from the clobber probe are non-fatal: we log + treat as
// "no clobber detected" so a flaky `git rev-parse` cannot wedge the
// pipeline (zero-swallow: error is still printed via stdout).
func conflictCheck(ctx *runContext, plan replay.Plan, c walk.SourceCommit, stdout io.Writer) (shouldAbortRun, shouldSkipCommit bool) {
	clobbers, err := replay.DetectClobbers(plan)
	if err != nil {
		fmt.Fprintf(stdout, "commit-in: conflict probe %s: %v\n", c.Sha, err)
		return false, false
	}
	if len(clobbers) == 0 {
		return false, false
	}
	decision := finalize.Resolve(ctx.Resolved.ConflictMode, c.Sha, stdout)
	if decision == finalize.ConflictDecisionAbort {
		ctx.aborted = true
		return true, true
	}
	// ForceMerge: log clobber list at info-level so audits can see it.
	fmt.Fprintf(stdout, "commit-in: %s force-merge clobbering %d file(s)\n", c.Sha, len(clobbers))
	return false, false
}
