---
name: Commit Transfer Complete
description: commit-left / commit-right / commit-both all live, with --interleave variant for commit-both (v3.104.0). All three reuse committransfer.runOneDirection.
type: feature
---

# Commit Transfer (Phases 1–4 complete)

## Status by phase

| Phase | Command | Shipped | Engine |
|-------|---------|---------|--------|
| 1 | `commit-right` (cmr) | v3.76.0 | `RunRight` (now thin wrapper over `runOneDirection`) |
| 2 | `commit-left` (cml) | v3.102.0 | `RunLeft` (swap source/target, then `runOneDirection`) |
| 3 | `commit-both` (cmb) sequential | v3.102.0 | `RunBoth` = two sequential passes (L→R, then R→L) |
| 4 | `commit-both --interleave` | v3.104.0 | `RunBothInterleaved` = author-date merged stream |

## commit-both algorithm

### Default — sequential (v3.102.0)

1. Pass 1: LEFT → RIGHT (build plan, prompt, replay, push)
2. Pass 2: RIGHT → LEFT (rebuild plan, prompt, replay, push)
3. Pass 1 failure aborts before Pass 2 (avoids merge-base drift)

Each pass appends a directional suffix to LogPrefix (`(left→right)` /
`(right→left)`) so output is visually attributable. Options struct is
copied per pass (`withDirectionLabel`) so the original LogPrefix is
never mutated — locked in by `TestRunBothImmutableOptions`.

### `--interleave` — author-date (v3.104.0)

1. Build BOTH directional plans up front.
2. Merge commit lists into a single stream sorted by `AuthorAt`
   (stable sort; LEFT-side wins exact ties).
3. Print unified plan, prompt once, then walk the stream replaying
   each commit onto its opposite side.
4. After the stream, push each side that received commits and print
   per-side summaries.

Tradeoffs:

- Faithful to "what actually happened first" across both sides.
- One prompt instead of two.
- First per-commit failure aborts mid-stream → partial state on
  whichever side was being written. Use `--dry-run` first.
- No per-direction merge-base recompute between commits → later
  interleaved commits may re-touch files that an opposite-direction
  commit just modified.

CLI guard: `--interleave` is rejected (exit 2) for `commit-left` and
`commit-right`. Only `commit-both` accepts it.

## Files

- `gitmap/committransfer/runleftboth.go` — RunLeft, RunBoth, runOneDirection, withDirectionLabel
- `gitmap/committransfer/interleave.go` — RunBothInterleaved + helpers (v3.104.0)
- `gitmap/committransfer/runleftboth_test.go` — direction-label + immutability tests
- `gitmap/committransfer/interleave_test.go` — sort invariant + tie-breaking + empty cases
- `gitmap/committransfer/run.go` — RunRight (now a 1-line wrapper)
- `gitmap/cmd/committransfer.go` — `dispatchDirection` routes spec.Name → RunX, validates --interleave
- `spec/01-app/106-commit-left-right-both.md` — §5 split into 5.1 sequential + 5.2 --interleave
- `gitmap/helptext/commit-{left,both}.md` — status updated, examples added
