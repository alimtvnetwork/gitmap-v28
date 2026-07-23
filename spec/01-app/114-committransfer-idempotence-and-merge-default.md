# Spec 114 — Commit-Transfer Idempotence + Merge-Default Hardening

Status: PUBLISHED (Gap A resolved v5.78.0/v5.83.0; Gap B flag v5.62.0, default flip v6.0.0)
Version target: v5.62.0 – v6.0.0
Owners: committransfer package

## Problem

Two latent bugs in `gitmap/committransfer` were identified during the
v5.52.0 commit-transfer post-mortem but never tracked in a spec.

### Gap A — `recentLogSubjectsAndBodies` 200-commit cap

`gitmap/committransfer/plan.go:27` calls:

    recentLogSubjectsAndBodies(targetDir, 200)

The cap was chosen for short-lived targets, but on long-history
destinations the idempotence check only inspects the last 200 commits.
Any source commit that already landed earlier (e.g. cherry-picked weeks
ago into a busy `main`) is mis-classified as "fresh" and re-applied,
producing duplicate commits with new SHAs.

**Repro envelope:** target with >200 commits since the cherry-pick of
any source subject; source range includes that same subject.

### Gap B — Planner default `IncludeMerges=false`

`gitmap/committransfer/git.go:50-54` appends `--no-merges` when
`opts.IncludeMerges == false`, which is the zero-value default in
`types.go:81`. The CLI surface therefore silently drops merge commits
during transfer. This is the safe v1 default but is a footgun for users
transferring release branches whose merge-commits encode actual history.

Flipping the default is a **semver behaviour change** and requires a
deprecation window.

## Resolution

### Gap A — Unbounded log + early-exit set lookup

1. Replace the constant `200` with a sentinel meaning "all reachable":
   pass `n <= 0` to `recentLogSubjectsAndBodies` and skip the
   `--max-count` flag in that branch.
2. Convert the substring scan in the caller to a hash-set keyed on the
   full normalised subject+body block. Memory cost: O(N) strings; for a
   100k-commit target that is ~10 MB — acceptable for a one-shot CLI.
3. Add a regression test `TestPlanIdempotenceBeyond200Commits` that
   builds a target with 250+ commits, cherry-picks subject at index 5,
   and asserts the planner reports `AlreadyApplied`.
4. v5.83.0 follow-up: add `--max-history-scan N` (Go field
   `Options.MaxHistoryScan int`, default `0` = unbounded) as an escape
   hatch for operators running against pathologically large targets
   (e.g. mirrored monorepos in the tens of millions of commits) where
   the unbounded `git log` of the v5.78.0 fix is prohibitive. Default
   behaviour is unchanged; the knob is opt-in. Pinned by
   `gitmap/committransfer/maxhistoryscan_test.go`.

### Gap B — Default flip with deprecation

1. v5.62.0: keep `IncludeMerges=false` default; emit a one-line
   `os.Stderr` notice when the source range contains ≥1 merge that is
   being stripped, including the new flag name to opt-in.
2. v5.62.0: add `--include-merges` / `--no-include-merges` CLI flags;
   document them in `helptext/commit-in.md` and `helptext/commit-out.md`.
3. v6.0.0 (next major): flip default to `IncludeMerges=true`. Stderr
   notice inverts: warn when `--no-include-merges` strips commits.

## Out of scope

- Re-architecting plan.go around a streaming git-log reader. Current
  in-memory model is acceptable up to ~1M target commits.
- Persisting the dedup hash-set to SQLite. Not needed for a single CLI
  invocation.

## Files touched

- `gitmap/committransfer/plan.go` — call-site change
- `gitmap/committransfer/git.go` — `--max-count` branch
- `gitmap/committransfer/types.go` — doc updates only in v5.62.0
- `gitmap/cmd/committransfer.go` — `--include-merges` flag wiring
- `gitmap/helptext/commit-in.md`, `commit-out.md` — doc the flag
- `gitmap/committransfer/plan_idempotence_test.go` — NEW regression test

## Acceptance

- `go test ./gitmap/committransfer/...` green including the new
  beyond-200 regression test.
- Manual smoke: cherry-pick a commit 500 entries deep into target,
  re-run `gitmap commit-in`, observe `AlreadyApplied` (not duplicated).
- v5.62.0 changelog entry under "Fixes" (Gap A) and "Deprecations"
  (Gap B opt-in flag + planned v6 default flip).
