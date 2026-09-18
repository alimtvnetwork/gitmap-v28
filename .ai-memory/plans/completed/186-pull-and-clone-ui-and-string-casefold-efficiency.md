# Plan 186: Pull & Clone UI Line Wrap Fix and String Case-Fold Efficiency

> **Task Origin & Objective**:
> - Fix terminal output line wrapping and scrollback pollution during `gitmap pull` and parallel multi-repo operations (`media_1789541608964.png`).
> - Fix clone spinner line clearing on frame refresh in `cli/cmdclone/clonespinner.go`.
> - Implement zero-allocation case-folding (`strings.EqualFold`) and lazy short-circuiting in `cli/cmdagy/agy_ls_match.go` (`media_1789542237413.png` and `media_1789543475985.png`) and `cli/cmdinstall/install_archive_detect.go`.
> - Decompose oversized `pull_progress_bar.go` (516 lines) into 9 modular, single-responsibility files under the 100-line cap (all <= 91 lines).
> - Defer build checks to final step, 100% go vet pass, updated test inventory.

---

## 1. Problem Analysis & Root Cause

| Issue Area | Symptom | Root Cause | Remediation Plan |
|:---:|---|---|---|
| **Pull Progress UI** | Repeated broken progress bars filling terminal (`media_1789541608964.png`) | `formatActiveWorkersList()` joined all active parallel workers (up to 20+), creating 500-800 character lines that wrapped across 4-6 terminal rows. `\r\033[K` only returned to column 0 of the last wrapped line, leaving previous rows in scrollback | Bounded active worker display to at most 2 items with `(+N more)` indicator. Added `clampLineToTermWidth()` utilizing `detectTerminalWidth()` to clamp all TTY lines strictly to `termWidth - 1` with ellipsis. |
| **Clone Spinner UI** | Ghost trailing characters on spinner updates in `clonespinner.go` | Ticker frame render lacked `\033[K` clear-to-end-of-line escape code, only having it on exit stop channel | Added `\033[K` directly after `\r` on each ticker frame. |
| **String Efficiency** | `strings.ToLower(name)` & `strings.ToLower(path)` allocated lowercased strings repeatedly | In `matchesFilter`, both strings were eagerly converted to lowercase and re-checked every iteration without short-circuiting | Added zero-allocation `strings.EqualFold` fast path for exact matches. Implemented lazy evaluation: if `name` matches `lowerFilter`, `path` is never allocated or searched. Replaced `strings.ToLower == strings.ToLower` with `strings.EqualFold` in `install_archive_detect.go`. |
| **File Sizing Cap** | `pull_progress_bar.go` was 516 lines | Monolithic file contained types, ticker, renderers, worker registry, visual formatting, and state metrics | Decomposed into 9 files under 100 lines each (all <= 91 lines). |

---

## 2. File Decomposition Breakdown

- `cli/cmdpull/pull_progress_bar_types.go`: 68 lines (WorkerSlotState, PullProgressBar structs, SetOutput, SetTTY)
- `cli/cmdpull/pull_progress_bar.go`: 49 lines (NewPullProgressBar, Start, Stop)
- `cli/cmdpull/pull_progress_bar_ticker.go`: 59 lines (initTickerLifecycle, runTickerLoop, tick, shutdownTicker, flushFinalOutput)
- `cli/cmdpull/pull_progress_bar_render.go`: 91 lines (Render, renderLocked, currentSpinner, renderTTY, renderSingleRepoTTY, renderMultiRepoTTY, clampLineToTermWidth, formatActiveDescription)
- `cli/cmdpull/pull_progress_bar_workers.go`: 53 lines (RegisterWorker, UpdateWorkerProgress, UnregisterWorker, ActiveWorkers)
- `cli/cmdpull/pull_progress_bar_workers_summary.go`: 62 lines (formatWorkersSummary, formatFallbackSummary, formatActiveWorkersList, formatWorkerSlotDesc, sortedWorkerIDs)
- `cli/cmdpull/pull_progress_bar_metrics.go`: 87 lines (CompleteRepo, applyStateMetrics, incrementStepCounters, incrementSuccessCounters, incrementFailureOrSkipCounters, IsStopped, States, Failed, Succeeded)
- `cli/cmdpull/pull_progress_bar_milestones.go`: 38 lines (UpdateRepoStep, SetActivePercent, SetSubStep, SubStepMilestone)
- `cli/cmdpull/pull_progress_bar_visual.go`: 71 lines (renderNonTTY, resolveSingleRepoPercent, calcProgressPercent, FormatVisualBar, calcFilledBarLength, resolveBarChars)

---

## 3. Verification & Quality Gates

- `go vet ./cmdpull/... ./cmdagy/... ./cmdclone/... ./cmdinstall/...`: Clean exit code 0.
- `03-ai-scripts/26-go-code-formatter.py`: All files formatted with gofmt.
- `03-ai-scripts/33-test-inventory-generator.py`: 3,536 tests indexed.
- Added unit tests in `cli/cmdpull/pull_progress_bar_test.go`:
  - `TestClampLineToTermWidth`: verifies short lines pass through and long lines (> termWidth) clamp cleanly with `...`.
  - `TestFormatActiveWorkersListBounded`: verifies 3+ active workers display at most 2 with `(+N more)` suffix.
