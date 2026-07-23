---
Slug: exclusion-parser
Status: pending
Created: 2026-06-06
Parent: 01-bulk-visibility-mapub-mapri
---

# Subtask 03 — Interactive prompt + exclusion parser

## Flow
1. Render matched repos as a numbered list (1-based).
2. Prompt 1: `Make N repos PUBLIC|PRIVATE? [y]es / [n]o / [e]xclude :`
3. On `e`: prompt 2: `Numbers to exclude (e.g. 1,3-5,9 or "none"):`
4. Filter excluded indexes; re-display remaining; prompt 1 again.
5. On `y`: apply. On `n`: abort, persist run with all rows `Skipped`.

## Parser (`gitmap/cmd/visibilitybulkprompt.go`)
- `parseExclusionInput(raw string, max int) (map[int]bool, error)`
- Accept comma-separated tokens; each token is either `N` or `N-M` (inclusive).
- Reject out-of-range; reject `M < N`; reject non-numeric.
- `"none"` / empty → empty set.
- `"all"` → set of `1..max` (rare, but supported for symmetry).

## -Y bypass
- `runMakeAllVisibility(..., opts)` checks `opts.autoConfirm`; if true, skip
  both prompts entirely AND skip the `IsExcluded` write loop (all rows
  `IsExcluded=0`).

## Tests (`gitmap/cmd/visibilitybulkprompt_test.go`)
- `"1,3-5"` → {1,3,4,5}
- `"1,2,2,3"` → {1,2,3} (dedup)
- `""` / `"none"` → {}
- `"all"` with max=4 → {1,2,3,4}
- `"5"` with max=4 → error
- `"3-1"` → error (descending)
- `"abc"` → error
