# 2026-05-09 — commit-transfer count mismatch (212 → 150)

## Symptom

User ran `gitmap commit-right LEFT RIGHT` (or any of `commit-left` /
`commit-both`). LEFT had **212** commits according to `git log`. RIGHT
ended up with **150** commits. No clear breakdown told the user where
the 62 missing commits went.

The commit-in family (`gitmap commit-in`) shares the same family of
silent-skip behaviours and is included below.

## Root causes (in order of impact)

### 1. `--no-merges` is the planner default

`gitmap/committransfer/git.go::revListReverse` adds `--no-merges`
unless `IncludeMerges` is set. The CLI default (`Options.IncludeMerges
= false`) silently strips every merge commit from the source list
before the plan is even printed. The user's `git log` count is the
**total** including merges, so:

    user-counted = mainline + merges
    plan-counted = mainline only

`PrintPlan` then announces `replaying N commits` where N is the
post-strip count, with no mention of how many were excluded.

### 2. Silent `SkippedEmpty` after snapshot

`replay.go::replayOne` checks `hasStagedChanges` after copying the
source-snapshot tree into the target. When the source commit's tree
is byte-identical to the target's current tree (e.g. consecutive
no-op commits, revert/redo cycles, or commits that touched only
ignored paths like `node_modules/`, `.git/`), `git diff --cached
--quiet` exits zero and the commit is dropped with `SkippedEmpty++`.
**No per-commit log line is printed**, so the user only learns the
aggregate count from `PrintSummary`'s `empty=N` field — which is
easy to miss in a long output, and gives no SHAs to investigate.

### 3. `recentLogSubjectsAndBodies(targetDir, 200)` cap

The idempotence check (`AlreadyReplayed`) only reads the last **200**
target commits. A re-run on a target that already carries >200
replayed commits can mis-classify commits as "fresh" and replay them
again — or, with provenance footer matching, can also miss the
already-replayed signal for older commits. This is a secondary
contributor; not the cause of the 212→150 case but worth fixing.

### 4. commit-in: stage-level skips are logged but not totaled

`gitmap/cmd/commitin/orchestrator/conflict.go` and the walk stages
log per-commit skip rows ("DuplicateSourceSha", "EmptyDiff", etc.),
and the final summary prints `skipped=N`. The bug class is the same
shape but lower severity: the user already sees the row.

## Fix shipped (this commit)

1. **Reconciliation in `BuildPlan`** — also runs `rev-list` *with*
   merges to compute `MergeExcluded`. Stored on `ReplayPlan` so
   downstream UI can surface it.
2. **`PrintPlan` headline reformatted** — now prints
   `replaying R commits onto target (source-considered=C,
   merge-excluded=M; pass --include-merges to replay merges)`.
3. **`replayOne` per-commit empty log** — when the snapshot produces
   no staged diff, the loop now prints
   `[i/n] <shortSHA> → empty (snapshot tree unchanged on target)`
   so the user can see *which* SHAs were dropped.
4. **`PrintSummary` reconciliation line** — added a second line:
   `source considered=C, replayed=R, skipped=S; accounted=R+S=C ✓
   (or ✗ — discrepancy)`. The `✗` path also writes to os.Stderr
   so CI scripts can detect drift.
5. **E2E test** —
   `gitmap/committransfer/count_parity_e2e_test.go` builds two real
   git repos via plumbing, runs `RunRight`, and asserts:
   - Pure-mainline 5-commit source → 5 target commits.
   - Source with 1 merge commit → MergeExcluded=1 reported, 0 merge
     commits in target.
   - `Replayed + SkippedDrop + SkippedReplayed + SkippedEmpty ==
     len(plan.Commits)` — the replay-side accounting invariant.

## What was *not* fixed at the time (subsequently resolved)

- ✅ `recentLogSubjectsAndBodies` 200-commit cap — **resolved v5.78.0** (unbounded scan, pass `n <= 0`); **v5.83.0** added `--max-history-scan` escape hatch for pathologically large targets. See `spec/01-app/114-committransfer-idempotence-and-merge-default.md`.
- ✅ Defaulting `IncludeMerges=true` — **resolved v6.0.0** (breaking change). Merge commits are now preserved by default; `--no-include-merges` opts back into legacy strip behaviour. See `spec/01-app/115-v6-migration.md`.

## Files touched

- gitmap/committransfer/types.go — `ReplayPlan.MergeExcluded`
- gitmap/committransfer/plan.go — count merges-excluded
- gitmap/committransfer/replay.go — per-commit empty log + signature
- gitmap/committransfer/log.go — reconciliation summary
- gitmap/committransfer/count_parity_e2e_test.go — new E2E
- .lovable/memory/issues/2026-05-09-commit-transfer-count-mismatch.md

**Status:** All root causes addressed. Issue closed as of v6.0.0.
