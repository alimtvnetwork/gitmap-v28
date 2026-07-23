package cmd

// Dry-run preview helpers for `gitmap fix-repo` (v5.40.0+). Mirrors
// the rewrite engine but never writes to disk and returns a
// per-target breakdown of would-be substitutions so callers can
// surface "what would change" at the granularity of each
// `{base}-vN` rule (plus the v1→v2 bare-base sweep).

import "os"

// fixRepoTargetHit summarizes one rule's effect on a file. n == -1
// is the sentinel for the bare-base sweep so dry-run output can
// label it distinctly from the numbered targets.
type fixRepoTargetHit struct {
	n     int
	count int
}

// fixRepoBareBaseSentinel marks a hit as coming from the v1→v2
// bare-base sweep rather than a numbered `{base}-vN` rule.
const fixRepoBareBaseSentinel = -1

// previewFixRepoFile reads fullPath and computes (total, hits)
// without mutating disk. Returns an error only on read failure;
// scannable-but-empty files yield (0, nil, nil).
func previewFixRepoFile(fullPath, base string, current int,
	targets []int, restrictNoVersion bool,
) (int, []fixRepoTargetHit, error) {
	original, err := os.ReadFile(fullPath)
	if err != nil {
		return 0, nil, err
	}
	total, hits := previewAllTargets(string(original), base, current, targets, restrictNoVersion)

	return total, hits, nil
}

// previewAllTargets folds every target rule + the conditional
// bare-base sweep, accumulating per-rule hit counts. Pure function:
// safe to call from tests without touching the filesystem.
func previewAllTargets(text, base string, current int, targets []int,
	restrictNoVersion bool,
) (int, []fixRepoTargetHit) {
	hits := make([]fixRepoTargetHit, 0, len(targets))
	total := 0
	for _, n := range targets {
		text, hits, total = previewOneTarget(text, base, n, current, restrictNoVersion, hits, total)
	}

	return total, hits
}

// previewOneTarget applies one numbered rule plus (conditionally)
// the v1→v2 bare-base sweep, appending non-zero hits to the slice.
func previewOneTarget(text, base string, n, current int, restrictNoVersion bool,
	hits []fixRepoTargetHit, total int,
) (string, []fixRepoTargetHit, int) {
	updated, added := applyOneTarget(text, base, n, current)
	if added > 0 {
		hits = append(hits, fixRepoTargetHit{n: n, count: added})
	}
	total += added
	if n == 1 && current == 2 && !restrictNoVersion {
		u2, a2 := applyBareBase(updated, base, current)
		if a2 > 0 {
			hits = append(hits, fixRepoTargetHit{n: fixRepoBareBaseSentinel, count: a2})
		}

		return u2, hits, total + a2
	}

	return updated, hits, total
}
