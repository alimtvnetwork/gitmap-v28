---
name: parallel-pull
description: gitmap pull --parallel <N> runs a worker pool over targets; --only-available pre-filters via FindNext so we only pull repos with new tags. v3.10.0 (Phase 2.5).
type: feature
---
# Parallel pull (Phase 2.5, v3.10.0)

## Overview

Two new flags ship together on `gitmap pull`:

- `--parallel <N>` (default `1`) — runs the pull batch through a worker pool of width N. N=1 preserves exact pre-2.5 serial behavior.
- `--only-available` — intersects the resolved target set with `(*DB).FindNext(0)` (every repo whose latest `VersionProbe` row has `IsAvailable = 1`) before pulling. Turns "pull all 42 repos" into "pull the 6 that actually have new tags".

Combined: `gitmap probe --all && gitmap pull --all --only-available --parallel 4`.

## Worker pool design

`cmd/pullparallel.go` implements a classic buffered-channel worker pool:

1. `jobs` channel (buffered to `len(records)`)
2. N workers spawn via `startPullWorkers`, each draining the channel
3. Dispatcher `dispatchPullJobs` feeds records, respecting `stopOnFail` via a shared `*bool` guarded by `progMu`
4. `wg.Wait()` blocks until every worker exits

`BatchProgress` is **not goroutine-safe** — its counters and failure slice are mutated without locks. Rather than refactor it, every `BeginItem`/`Succeed`/`FailWithError`/`Skip` call inside a worker is wrapped in `progMu.Lock() / Unlock()`. This trivially serializes progress mutations without serializing the actual git work.

## Filter design

`cmd/pullfilter.go::filterByAvailableUpdates`:

1. Calls `(*DB).FindNext(0)` to get every repo with an available update across the whole DB
2. Builds a `map[int64]bool` keyed by `Repo.ID`
3. Intersects the resolved target set against that set

**Fail-open**: when the probe DB cannot be opened or `FindNext` errors, the filter logs a warning to stderr and returns the original record set unchanged. Failing closed (returning empty) would surprise users — they'd type `pull --all --only-available` and see nothing happen with no obvious cause. The warning is `WarnPullFilterFallback` in `constants_pull.go`.

## CLI surface

```
gitmap pull <slug>                                  # serial, single repo (unchanged)
gitmap pull --all                                   # serial, all repos (unchanged)
gitmap pull --all --parallel 4                      # 4-way parallel
gitmap pull --all --only-available                  # pre-filtered, serial
gitmap pull --all --only-available --parallel 8     # full combo
gitmap pull --group backend --parallel 4            # group + parallel
```

## Test coverage

`cmd/flags_test.go` was migrated from the old 5-tuple return of `parsePullFlags` to the new `pullOptions` struct. Two new tests cover the new flags: `TestParsePullFlags_Parallel` and `TestParsePullFlags_OnlyAvailable`.

## Files

- `gitmap/cmd/pullparallel.go` — worker pool: `runPullParallel`, `startPullWorkers`, `dispatchPullJobs`, `pullWorker`, `runOnePullJob`
- `gitmap/cmd/pullfilter.go` — `filterByAvailableUpdates`, `loadAvailableRepoIDs`, `intersectByID`
- `gitmap/cmd/pull.go` — refactored `runPull` + new `pullOptions` struct, `beginPullTask`, `executePull`
- `gitmap/cmd/flags_test.go` — migrated tests + new flag coverage
- `gitmap/constants/constants_pull.go` — `FlagDescPullParallel`, `FlagDescPullOnlyAvailable`, `MsgPullNoAvailable`, `WarnPullFilterFallback`
- `gitmap/helptext/pull.md` — flag table extended, new Example 5 added
