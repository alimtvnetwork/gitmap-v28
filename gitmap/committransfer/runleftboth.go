package committransfer

import (
	"fmt"
	"os"
)

// RunLeft replays RIGHT → LEFT (writes commits on LEFT). It is the
// directional mirror of RunRight: same plan/replay primitives, with
// source/target swapped. The CLI passes leftDir as the destination
// because `commit-left` writes to LEFT (mirrors `merge-left` semantics
// in spec/01-app/97-move-and-merge.md).
//
// Phase 2 (v3.102.0): wired in by gitmap/cmd/committransfer.go.
func RunLeft(leftDir, rightDir string, opts Options) error {
	// LEFT is the destination, RIGHT is the source.
	return runOneDirection(rightDir, leftDir, opts)
}

// RunBoth replays in both directions: first LEFT → RIGHT, then
// RIGHT → LEFT. The two passes share Options but each one builds its
// own plan and prints its own preview/summary. We deliberately run
// them sequentially (not interleaved by author date) so each side's
// final state is deterministic and the user can audit each direction
// independently in the printed summary.
//
// If the first pass fails the second is skipped — partial commit-both
// is worse than half-done because the second direction's merge-base
// would have shifted. The caller exits with the first error.
//
// Phase 3 (v3.102.0): wired in by gitmap/cmd/committransfer.go.
func RunBoth(leftDir, rightDir string, opts Options) error {
	leftToRightOpts := withDirectionLabel(opts, "(left→right)")
	if err := runOneDirection(leftDir, rightDir, leftToRightOpts); err != nil {
		return fmt.Errorf("left→right pass: %w", err)
	}

	rightToLeftOpts := withDirectionLabel(opts, "(right→left)")
	if err := runOneDirection(rightDir, leftDir, rightToLeftOpts); err != nil {
		return fmt.Errorf("right→left pass: %w", err)
	}

	return nil
}

// runOneDirection is the single-pass engine extracted from RunRight so
// RunLeft / RunBoth can reuse it without duplicating the prompt+replay
// orchestration. The legacy RunRight wraps this for backward compat.
func runOneDirection(sourceDir, targetDir string, opts Options) error {
	plan, err := BuildPlan(sourceDir, targetDir, opts)
	if err != nil {
		return fmt.Errorf("build plan: %w", err)
	}
	willReplay := PrintPlan(os.Stdout, plan, opts.LogPrefix)
	if willReplay == 0 {
		fmt.Fprintf(os.Stdout, "%s nothing to replay.\n", opts.LogPrefix)

		return nil
	}
	if !opts.DryRun && !opts.Yes && !Confirm(opts.LogPrefix) {
		fmt.Fprintf(os.Stderr, "%s aborted by user.\n", opts.LogPrefix)

		return nil
	}
	res, replayErr := Replay(plan, opts)
	if replayErr != nil {
		PrintSummary(os.Stderr, opts.LogPrefix, res)
		PrintReconciliation(os.Stderr, os.Stderr, opts.LogPrefix, plan, res)

		return replayErr
	}
	res.Pushed = maybePush(targetDir, opts, len(res.NewSHAs))
	PrintSummary(os.Stdout, opts.LogPrefix, res)
	PrintReconciliation(os.Stdout, os.Stderr, opts.LogPrefix, plan, res)

	return nil
}

// withDirectionLabel returns a copy of opts with " <suffix>" appended
// to LogPrefix so commit-both's two passes are visually distinguishable
// in the user's terminal output. We do not mutate the caller's struct.
func withDirectionLabel(opts Options, suffix string) Options {
	out := opts
	out.LogPrefix = fmt.Sprintf("%s %s", opts.LogPrefix, suffix)

	return out
}
