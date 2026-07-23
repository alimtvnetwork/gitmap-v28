package committransfer

import (
	"fmt"
	"os"
	"sort"
)

// RunBothInterleaved replays LEFT and RIGHT commits in a single
// chronological stream sorted by author date. Each commit gets
// replayed onto the *opposite* side at the moment its timestamp comes
// up in the unified timeline.
//
// Compared with RunBoth (sequential L→R then R→L):
//   - More faithful to "what actually happened first" across both sides.
//   - Matches the original spec §5 wording before sequential simplification.
//   - Slightly higher risk: per-commit failures abort mid-stream, leaving
//     a partial state on whichever side was being written when the error
//     hit. The user's confirmation prompt warns about this.
//
// Activated by `gitmap commit-both --interleave`.
//
// Phase 4 (v3.104.0): added by user request after sequential RunBoth shipped.
func RunBothInterleaved(leftDir, rightDir string, opts Options) error {
	leftToRight, err := BuildPlan(leftDir, rightDir, opts)
	if err != nil {
		return fmt.Errorf("interleave plan L→R: %w", err)
	}
	rightToLeft, err := BuildPlan(rightDir, leftDir, opts)
	if err != nil {
		return fmt.Errorf("interleave plan R→L: %w", err)
	}

	stream := buildInterleavedStream(leftToRight, rightToLeft)
	if len(stream) == 0 {
		fmt.Fprintf(os.Stdout, "%s nothing to replay (interleave).\n", opts.LogPrefix)

		return nil
	}

	return executeInterleaveStream(stream, leftToRight, rightToLeft, opts)
}

// interleaveStep is one ordered entry in the merged timeline.
type interleaveStep struct {
	Commit    SourceCommit
	Direction string // "L→R" or "R→L"
}

// buildInterleavedStream merges both directional plans and sorts by
// author date. Stable sort keeps within-side ordering for ties.
func buildInterleavedStream(leftToRight, rightToLeft ReplayPlan) []interleaveStep {
	stream := make([]interleaveStep, 0, len(leftToRight.Commits)+len(rightToLeft.Commits))
	for _, c := range leftToRight.Commits {
		stream = append(stream, interleaveStep{Commit: c, Direction: "L→R"})
	}
	for _, c := range rightToLeft.Commits {
		stream = append(stream, interleaveStep{Commit: c, Direction: "R→L"})
	}
	sort.SliceStable(stream, func(i, j int) bool {
		return stream[i].Commit.AuthorAt.Before(stream[j].Commit.AuthorAt)
	})

	return stream
}

// executeInterleaveStream prints the unified plan, prompts once, then
// replays each step onto its target side in author-date order.
func executeInterleaveStream(stream []interleaveStep, ltr, rtl ReplayPlan, opts Options) error {
	printInterleavedPlan(stream, opts.LogPrefix)
	if !opts.DryRun && !opts.Yes && !Confirm(opts.LogPrefix) {
		fmt.Fprintf(os.Stderr, "%s aborted by user (interleave).\n", opts.LogPrefix)

		return nil
	}
	if opts.DryRun {
		return nil
	}

	return replayInterleaveSteps(stream, ltr, rtl, opts)
}

// printInterleavedPlan shows the unified author-date order so the user
// can audit before any commit is written.
func printInterleavedPlan(stream []interleaveStep, prefix string) {
	fmt.Fprintf(os.Stdout, "%s interleave plan: %d commits in author-date order\n", prefix, len(stream))
	for i, step := range stream {
		fmt.Fprintf(os.Stdout, "%s [%d/%d] %s  %s  %s\n",
			prefix, i+1, len(stream), step.Direction, step.Commit.ShortSHA, step.Commit.Subject)
	}
}

// replayInterleaveSteps walks the sorted stream and replays each
// commit onto its opposite side. Stops at the first error; reports
// progress per side via PrintSummary.
func replayInterleaveSteps(stream []interleaveStep, ltr, rtl ReplayPlan, opts Options) error {
	results := map[string]*ReplayResult{"L→R": {}, "R→L": {}}
	for i, step := range stream {
		plan := ltr
		if step.Direction == "R→L" {
			plan = rtl
		}
		if step.Commit.SkipCause != "" {
			results[step.Direction].SkippedDrop++

			continue
		}
		newSHA, _, err := replayOne(plan, step.Commit, opts)
		if err != nil {
			return fmt.Errorf("interleave step %d (%s %s): %w",
				i+1, step.Direction, step.Commit.ShortSHA, err)
		}
		if newSHA == "" {
			results[step.Direction].SkippedEmpty++

			continue
		}
		results[step.Direction].NewSHAs = append(results[step.Direction].NewSHAs, newSHA)
		results[step.Direction].Replayed++
	}
	finalizeInterleavePush(ltr, rtl, results, opts)

	return nil
}

// finalizeInterleavePush pushes each side that received commits and
// prints both per-side summaries. Skips push when --no-push or
// --no-commit is set, or when a side received zero commits.
func finalizeInterleavePush(ltr, rtl ReplayPlan, results map[string]*ReplayResult, opts Options) {
	results["L→R"].Pushed = maybePush(ltr.TargetDir, opts, len(results["L→R"].NewSHAs))
	results["R→L"].Pushed = maybePush(rtl.TargetDir, opts, len(results["R→L"].NewSHAs))
	PrintSummary(os.Stdout, fmt.Sprintf("%s (L→R)", opts.LogPrefix), *results["L→R"])
	PrintSummary(os.Stdout, fmt.Sprintf("%s (R→L)", opts.LogPrefix), *results["R→L"])
}
