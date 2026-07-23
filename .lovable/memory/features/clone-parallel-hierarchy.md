---
name: clone-parallel-hierarchy
description: gitmap clone now preserves nested folder hierarchy and supports --max-concurrency for bounded parallel execution. Default 1 = sequential. Workers, progress, cache all thread-safe.
type: feature
---

# Clone Parallel + Hierarchy (v3.101.0)

## Hierarchy preservation

Every record clones into `filepath.Join(targetDir, rec.RelativePath)` —
no flattening, no path rewriting. Holds for CSV, JSON, and text inputs
and for both runner modes. Locked in by
`TestCloneAllPreservesNestedHierarchy` in
`gitmap/cloner/concurrent_test.go` (sequential + parallel subtests).

## --max-concurrency flag

Opt-in parallel runner. Default `1` keeps the legacy sequential
ordering for stderr progress lines. When `N > 1`:

- A bounded worker pool of N goroutines drains a buffered job channel
  (`gitmap/cloner/concurrent.go`).
- Cache hits short-circuit before workers receive them (so a fully
  cached run is a no-op regardless of N).
- A single `↪ parallel clone enabled: N workers` header line lands
  before per-repo lines so scripts can detect the mode.
- Progress lines arrive in completion order (not input order).
- On-disk hierarchy is unchanged because each worker still uses
  `rec.RelativePath` verbatim.

Invalid values (≤ 0) exit 1 with `ErrCloneMaxConcurrencyInvalid` —
never silently degrade to a default.

## Thread-safety contract

- `Progress` (`progress.go`) — single mutex guards counters + stderr writes.
- `CloneCache` (`cache.go`) — single mutex guards the entries map.
- `cloneOrPullOne` and `isGitRepo` are pure (exec + stat) → safe.
- `runConcurrent` collector is the only writer of summary, cache.Record,
  and progress.Done/Skip/Fail in the parallel path.

## File layout (post-split)

`gitmap/cloner/` package now splits into:
- `cloner.go` — entry points (`CloneFromFile*`), parsers, `cloneOne`/`runClone`
- `runners.go` — `cloneAll` dispatcher, `runSequential`, `normalizeWorkers`
- `concurrent.go` — `runConcurrent` worker pool + collector
- `summary.go` — `recordTag`, `pickURL`, `updateSummary*`
- `progress.go` — thread-safe `Progress`
- `cache.go` — thread-safe `CloneCache`
- `concurrent_test.go` — hierarchy + worker-clamp tests

All files under the 200-line guideline. Each function under 15 lines.
