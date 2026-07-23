// Package orchestrator wires the commit-in sub-packages together.
// Each public function maps 1:1 to a spec §3.1 stage so the data-flow
// is auditable from the spec.
package orchestrator

import (
	"fmt"
	"io"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/cmd/commitin"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/cmd/commitin/finalize"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/cmd/commitin/profile"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/cmd/commitin/runlog"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/cmd/commitin/workspace"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// Run executes the entire commit-in pipeline for one parsed argv.
// Returns the exit code the caller should propagate to os.Exit. Never
// panics; every stage funnels its error through a single mapping seam.
func Run(raw *commitin.RawArgs, stdout, stderr io.Writer) int {
	ctx, code := setUp(raw, stderr)
	if code != constants.CommitInExitOk {
		return code
	}
	defer ctx.Cleanup()
	if code := maybeSaveProfile(ctx, stderr); code != constants.CommitInExitOk {
		_ = runlog.FinishRun(ctx.DB.Conn(), ctx.RunID, constants.CommitInRunStatusFailed, time.Now())
		return code
	}
	if code := executePipeline(ctx, stdout); code != constants.CommitInExitOk {
		_ = runlog.FinishRun(ctx.DB.Conn(), ctx.RunID, constants.CommitInRunStatusFailed, time.Now())
		return code
	}
	_ = runlog.FinishRun(ctx.DB.Conn(), ctx.RunID, finalRunStatus(ctx.Counters), time.Now())
	finalize.PrintSummary(stderr, ctx.Counters)
	if ctx.Raw.IsDryRun {
		finalize.PrintDryRunBanner(stderr)
	}
	return finalize.Outcome(ctx.Counters)
}

// finalRunStatus picks the spec §4 RunStatus enum for the FinishRun
// row based on the per-commit counters.
func finalRunStatus(c finalize.Counters) string {
	if c.Failed == 0 {
		return constants.CommitInRunStatusCompleted
	}
	if c.Created == 0 {
		return constants.CommitInRunStatusFailed
	}
	return constants.CommitInRunStatusPartiallyFailed
}

// setUp performs the no-mutation prerequisites: source resolution,
// workspace creation, lock acquisition, DB open + migrate, run row.
// On any failure it returns a non-zero exit code and a nil ctx; the
// caller MUST NOT defer ctx.Cleanup in that case.
func setUp(raw *commitin.RawArgs, stderr io.Writer) (*runContext, int) {
	src, code := resolveSource(raw, stderr)
	if code != constants.CommitInExitOk {
		return nil, code
	}
	paths, code := ensureWorkspace(src.Path, stderr)
	if code != constants.CommitInExitOk {
		return nil, code
	}
	lock, code := acquireLock(paths, stderr)
	if code != constants.CommitInExitOk {
		return nil, code
	}
	return finishSetUp(raw, src, paths, lock, stderr)
}

// finishSetUp opens the DB, starts the run row, and resolves settings.
// Split out so setUp stays under the 15-line cap.
func finishSetUp(raw *commitin.RawArgs, src *workspace.SourceHandle, paths *workspace.Paths, lock *workspace.LockHandle, stderr io.Writer) (*runContext, int) {
	db, code := openAndMigrate(paths, stderr)
	if code != constants.CommitInExitOk {
		lock.Release()
		return nil, code
	}
	resolved, prof, code := loadProfile(raw, paths, db, stderr)
	if code != constants.CommitInExitOk {
		_ = db.Close()
		lock.Release()
		return nil, code
	}
	runID, code := startRun(db, src, prof, stderr)
	if code != constants.CommitInExitOk {
		_ = db.Close()
		lock.Release()
		return nil, code
	}
	return newContext(raw, src, paths, lock, db, resolved, runID), constants.CommitInExitOk
}

// startRun inserts the CommitInRun row and returns its primary key.
func startRun(db dbCloser, src *workspace.SourceHandle, prof *profile.Profile, stderr io.Writer) (int64, int) {
	var url *string
	if src.Kind == workspace.SourceKindCloned {
		s := src.Path
		url = &s
	}
	var profID *int64
	_ = prof // profile->ID not yet persisted; placeholder for future ProfileId FK
	id, err := runlog.StartRun(db.Conn(), src.Path, url, src.IsFreshlyInit, profID, time.Now())
	if err != nil {
		fmt.Fprintf(stderr, constants.CommitInErrDbWrite, err)
		return 0, constants.CommitInExitDbFailed
	}
	return id, constants.CommitInExitOk
}
