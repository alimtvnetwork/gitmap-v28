package orchestrator

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/cmd/commitin/checkpoint"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/cmd/commitin/walk"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/cmd/commitin/workspace"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)


// executePipeline performs the per-input walk + replay loop. Returns
// the exit code; the summary is printed by the caller.
func executePipeline(ctx *runContext, stdout io.Writer) int {
	inputs, code := expandAndStage(ctx, stdout)
	if code != constants.CommitInExitOk {
		return code
	}
	for _, staged := range inputs {
		if code := processOneInput(ctx, staged, stdout); code != constants.CommitInExitOk {
			return code
		}
		if ctx.aborted {
			return constants.CommitInExitConflictAborted
		}
	}
	return constants.CommitInExitOk
}

func expandAndStage(ctx *runContext, stdout io.Writer) ([]workspace.StagedInput, int) {
	resolved, err := workspace.ExpandInputs(ctx.Source.Path, ctx.Raw.Inputs, ctx.Raw.Keyword, ctx.Raw.KeywordTail)
	if err != nil {
		fmt.Fprint(stdout, err.Error())
		return nil, constants.CommitInExitInputUnusable
	}
	fmt.Fprintf(stdout, constants.CommitInMsgPhaseStageInputs, len(resolved), ctx.TempDir)
	staged, err := workspace.CloneInputs(ctx.Paths, ctx.RunID, resolved)
	if err != nil {
		fmt.Fprint(stdout, err.Error())
		return nil, constants.CommitInExitInputUnusable
	}
	return staged, constants.CommitInExitOk
}

func processOneInput(ctx *runContext, staged workspace.StagedInput, stdout io.Writer) int {
	commits, err := walk.WalkFirstParent(staged.WorkPath)
	if err != nil {
		fmt.Fprintf(stdout, constants.CommitInErrInputOpen, staged.Input.Original, err)
		ctx.Counters.Failed++
		return constants.CommitInExitOk
	}
	fmt.Fprintf(stdout, constants.CommitInMsgPhaseWalk, len(commits))
	cp := openCheckpoint(ctx, staged)
	picker := newPicker()
	for _, c := range commits {
		if cp != nil && cp.IsDone(c.Sha) {
			ctx.Counters.Skipped++
			continue
		}
		processOneCommit(ctx, staged, c, picker, stdout)
		if cp != nil {
			_ = cp.MarkDone(c.Sha)
		}
		if ctx.aborted {
			return constants.CommitInExitOk // outer loop sees ctx.aborted and exits
		}
	}
	return constants.CommitInExitOk
}

// openCheckpoint returns the per-input checkpoint handle, or nil when
// the state dir cannot be created (resume becomes a best-effort no-op
// rather than failing the run).
func openCheckpoint(ctx *runContext, staged workspace.StagedInput) *checkpoint.File {
	stateDir := filepath.Join(ctx.Paths.CommitInRoot, checkpoint.DirName)
	fp := inputFingerprint(staged)
	cp, err := checkpoint.Open(stateDir, fp, ctx.RunID)
	if err != nil {
		return nil
	}
	return cp
}

// inputFingerprint stabilizes the state.json filename across runs of
// the same input. Hash of the work path keeps it filesystem-safe.
func inputFingerprint(staged workspace.StagedInput) string {
	sum := sha1.Sum([]byte(staged.Input.Original))
	return hex.EncodeToString(sum[:8])
}


// newPicker returns a deterministic-seeded RNG bound to the wall clock
// (per-run seed satisfies spec §3.4 "deterministic within a run").
func newPicker() func(n int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return func(n int) int {
		if n <= 0 {
			return 0
		}
		return r.Intn(n)
	}
}
