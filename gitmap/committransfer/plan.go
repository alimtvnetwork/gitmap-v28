package committransfer

import (
	"fmt"
)

// BuildPlan computes the replay set for one direction. It does NOT
// mutate either repo — it only reads source history and the target's
// recent log (for the idempotence check).
func BuildPlan(sourceDir, targetDir string, opts Options) (ReplayPlan, error) {
	sourceHead, err := currentRefName(sourceDir)
	if err != nil {
		return ReplayPlan{}, fmt.Errorf("read source HEAD ref: %w", err)
	}
	base, err := resolveBase(sourceDir, targetDir, opts.Since)
	if err != nil {
		return ReplayPlan{}, err
	}
	shas, err := revListReverse(sourceDir, base, "HEAD", opts.IncludeMerges)
	if err != nil {
		return ReplayPlan{}, fmt.Errorf("rev-list source: %w", err)
	}
	mergeExcluded := countMergeExcluded(sourceDir, base, opts.IncludeMerges, len(shas))
	if opts.Limit > 0 && len(shas) > opts.Limit {
		shas = shas[:opts.Limit]
	}
	// Unbounded by default (spec 114 Gap A — prevents false-fresh
	// classification on targets with >200 commits since the
	// already-applied source commit). opts.MaxHistoryScan > 0 lets
	// callers cap the scan on very large targets.
	recentTargetLog, _ := recentLogSubjectsAndBodies(targetDir, opts.MaxHistoryScan)
	replayedSet := BuildReplayedSet(recentTargetLog)

	plan, err := assemblePlan(sourceDir, targetDir, sourceHead, base, shas, replayedSet, opts)
	if err == nil {
		plan.MergeExcluded = mergeExcluded
		plan.IncludeMerges = opts.IncludeMerges
	}

	return plan, err
}

// countMergeExcluded reports how many merge commits the planner
// stripped from the source range. When IncludeMerges is true, merges
// are kept and the count is always 0. When false, we re-run rev-list
// without --no-merges and subtract.
func countMergeExcluded(sourceDir, base string, includeMerges bool, mainlineCount int) int {
	if includeMerges {
		return 0
	}
	withMerges, err := revListReverse(sourceDir, base, "HEAD", true)
	if err != nil {
		return 0
	}
	delta := len(withMerges) - mainlineCount
	if delta < 0 {
		return 0
	}

	return delta
}

// resolveBase honors --since when set; otherwise asks git for the
// merge-base. An unrelated history yields "" (use full source history).
func resolveBase(sourceDir, targetDir, since string) (string, error) {
	if since != "" {
		return since, nil
	}
	targetHead, err := gitOut(targetDir, "rev-parse", "HEAD")
	if err != nil {
		// Empty target repo — no base, replay full source history.
		return "", nil
	}

	return mergeBase(sourceDir, "HEAD", targetHead)
}

// assemblePlan turns raw SHAs into hydrated SourceCommit entries with
// the message pipeline + idempotence check applied.
func assemblePlan(sourceDir, targetDir, sourceHead, base string,
	shas []string, replayedSet map[string]struct{}, opts Options,
) (ReplayPlan, error) {
	plan := ReplayPlan{
		SourceDir: sourceDir, TargetDir: targetDir,
		SourceHEAD: sourceHead, BaseSHA: base,
	}
	for _, sha := range shas {
		entry, err := hydrateCommit(sourceDir, sha, replayedSet, opts)
		if err != nil {
			return plan, err
		}
		if entry.SkipCause == "drop-pattern" || isDropSkip(entry.SkipCause) {
			plan.SkippedDrop++
		}
		plan.Commits = append(plan.Commits, entry)
	}

	return plan, nil
}

// hydrateCommit reads one source commit, runs the message pipeline, and
// flags it as skipped when the pipeline says so or when the target
// already carries its provenance footer.
func hydrateCommit(sourceDir, sha string, replayedSet map[string]struct{}, opts Options) (SourceCommit, error) {
	subject, body, author, shortSHA, when, err := readCommit(sourceDir, sha)
	if err != nil {
		return SourceCommit{}, fmt.Errorf("read commit %s: %w", sha, err)
	}
	entry := SourceCommit{
		SHA: sha, ShortSHA: shortSHA, Subject: subject, Body: body,
		Author: author, AuthorAt: when,
	}
	if !opts.ForceReplay && opts.Message.Provenance &&
		SetHasReplayed(replayedSet, opts.Message.SourceDisplayName, shortSHA) {
		entry.SkipCause = "already-replayed"

		return entry, nil
	}
	cleaned := CleanMessage(subject, body, opts.Message, shortSHA, when)
	if cleaned.Skipped != "" {
		entry.SkipCause = cleaned.Skipped

		return entry, nil
	}
	entry.Cleaned = cleaned.Final

	return entry, nil
}

// isDropSkip reports whether a SkipCause originated from the drop filter.
func isDropSkip(cause string) bool {
	return len(cause) >= len("drop-pattern") && cause[:len("drop-pattern")] == "drop-pattern"
}
