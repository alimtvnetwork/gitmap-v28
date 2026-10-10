# RCA-108: Unused Symbols in Release Notes and Missing ScanPipe Package in CI Unit Tests

**Date:** 2026-10-10
**Status:** ✅ Resolved
**Severity:** High (CI Pipeline Failure across Lint Baseline Diff and Unit Test steps)
**Affected Workflows:** `CI (#38048066838)`
**Run URLs:**
- `https://github.com/alimtvnetwork/gitmap-v28/actions/runs/38048066838`

---

## 1. Symptom & Error Signature

Following commit `162b7d4`, GitHub Actions CI workflow `#38048066838` failed on two separate jobs:

1. **Job `Lint Baseline Diff` (Step: `Diff vs baseline (fail only on NEW findings)`):**
   ```text
   [unused] const `releaseNotesFormatJSON` is unused (NEW in /tmp/lint-current/report.json)
       ❌ cmdrelease/release_notes_opts.go:37: [unused] const `releaseNotesFormatJSON` is unused
   [unused] func `applyReleaseNotesFlag` is unused (NEW in /tmp/lint-current/report.json)
       ❌ cmdrelease/release_notes_opts.go:40: [unused] func `applyReleaseNotesFlag` is unused
   [unused] func `applyReleaseNotesArg` is unused (NEW in /tmp/lint-current/report.json)
       ❌ cmdrelease/release_notes_opts.go:63: [unused] func `applyReleaseNotesArg` is unused
   Process completed with exit code 1.
   ```

2. **Job `Test (unit)` (Step: `Run tests`):**
   ```text
   # ./mapper/...
   pattern ./mapper/...: lstat ./mapper/: no such file or directory
   FAIL	./mapper/... [setup failed]
   Process completed with exit code 1.
   ```

---

## 2. Root Cause

1. **Unexported & Disconnected Release Notes Implementation in `cmdrelease`:**
   - In `cli/cmdrelease/release_tools.go`, `runReleaseNotes` was declared unexported (`func runReleaseNotes(args []string) error`) and was not invoked anywhere within `cmdrelease`.
   - Concurrently, in `cli/cmd/roottooling.go`, `constants.CmdReleaseNotes` was wired to a stub `cmdchangelog.RunReleaseNotes(argsTail())` which returned `"release-notes not implemented"`.
   - Because `runReleaseNotes` was never called, its entire underlying call graph in `release_notes_opts.go` (`runReleaseNotesV2`, `parseReleaseNotesArgs`, `applyReleaseNotesFlag`, `applyReleaseNotesArg`, `renderReleaseNotes`, `renderJSON`, and `releaseNotesFormatJSON`) became dead/unreachable code.
   - When functions were decomposed into granular helpers, the newly introduced identifiers were flagged by `golangci-lint unused` as NEW violations absent from the baseline.

2. **Package Consolidation Left Outdated Path in CI Matrix:**
   - In commit `171b244c0` (Wave 3 package consolidation), `cli/mapper/` was merged into `cli/scanpipe/`.
   - However, `.github/workflows/ci.yml` retained `./mapper/...` in the `unit` test shard package list (`./clonenext/... ./cloner/... ./cmd/... ./config/... ./formatter/... ./mapper/... ./release/...`).
   - When `go test` executed on the ephemeral CI runner, `lstat ./mapper/` failed because the folder no longer existed.

3. **File Size Cap Hygiene on Refactored Files:**
   - `cli/cmdrelease/release_notes_opts.go` was 314 lines, and `cli/cmd/roottooling.go` was 494 lines.
   - Touching these files required decomposing them under the mandatory 300-line cap (`52-file-size-check.py`).

---

## 3. Applied Solution & Code Changes

1. **Exported & Wired Release Notes Command:**
   - Exported `RunReleaseNotes` in `cli/cmdrelease/release_tools.go` delegating to `runReleaseNotesV2(args)`.
   - Updated `cli/cmd/roottooling.go` to dispatch `constants.CmdReleaseNotes` directly to `cmdrelease.RunReleaseNotes(argsTail())`.
   - Updated `cli/cmdchangelog/release_notes_stub.go` to delegate to `cmdrelease.RunReleaseNotes(args)` as a compatibility shim.
   - Purged dead unexported duplicate flags in `cli/cmdrelease/releaseargs.go` and unused helper `resolveCurrentRepoID` in `cli/cmdrelease/releasepersist.go`.

2. **Split Oversized Files Under 300-Line Limit:**
   - Extracted formatting and JSON serialization from `cli/cmdrelease/release_notes_opts.go` into `cli/cmdrelease/release_notes_render.go` (116 lines), reducing `release_notes_opts.go` to 184 lines.
   - Extracted utility and chrome tool dispatch tables from `cli/cmd/roottooling.go` into `cli/cmd/roottooling_util.go` (203 lines), reducing `roottooling.go` to 283 lines.

3. **Updated CI Workflow Package Matrix:**
   - In `.github/workflows/ci.yml`, updated the `unit` test shard from `./mapper/...` to `./scanpipe/...`.

---

## 4. Prevention Invariant

1. **Dead Code Elimination Prior to Merge:**
   - Command dispatchers in `cli/cmd/` must always point to real implementations rather than stubs when underlying packages provide the concrete engine.
   - Whenever an unexported command runner is created in a subpackage, it must be exported and tested or wired to prevent dead-code decay.
2. **Matrix Package Synchronization with Refactored Packages:**
   - Any package move or rename (such as `mapper` to `scanpipe`) must audit `.github/workflows/ci.yml` and local CI runner scripts (`06-cicd-local-runner.py`) for package paths.
3. **Continuous 300-Line Cap Enforcement:**
   - All modified files must be verified with `python 03-ai-scripts/52-file-size-check.py --diff HEAD --limit 300` before staging.
