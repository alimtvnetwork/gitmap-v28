# Issue 106: History Purge Shadowing, Recycle Bin Export, Misspell Canceled, and Relative Paths RCA

> **Issue ID:** `CICD-106`
> **Repository:** `gitmap-v28`
> **Status:** Resolved
> **Date:** 2026-10-07

---

## 1. Symptom

Failures detected in GitHub Actions runs (37505869271 and 37505868690) via `gitmap pe`:

1. **Relative Paths Linter Failure:**
   - Absolute `file:///d:/work/gitmap/...` and hardcoded `D:\work\gitmap` URIs detected in `.ai-memory/cicd-issues/105-pipeline-pe-1-helptext-constants-nested-ifs-unused-and-bools-rca.md` and `repo-secrets/01-gitmap/01-chrome-import-export-test/readme.md`, failing `linter-scripts/check-relative-paths.py`.

2. **Golangci-Lint Static Checks (`Full Suite Guard`):**
   - `cmdpurge/purge_engine.go:58:2`: field `fileRecords` is unused (`unused`).
   - `cmdpurge/purge_recycle_other.go:9:6`: func `sendToRecycleBin` is unused (`unused`).
   - `cmd/historyrewrite.go:31:6`: func `runHistoryPurge` is unused (`unused`).
   - `cmdpurge/purge_engine.go:125:26`: `cancelled` is a misspelling of `canceled` (`misspell`).

3. **History Rewrite Smoke Test Failure (`Ubuntu`):**
   - `Scenario 1 — history-purge removes file across all commits` failed with:
     `FAIL: could not parse sandbox path from output`
     Stack trace indicated `RunHistoryPurgeCLI` was invoked with `secret.env --no-push --keep-sandbox`, which failed with `ERR_PURGE_NO_TARGET` because positional paths were expected by the filter-repo mirror-clone runner.

---

## 2. Root Cause

1. **Dispatch Shadowing:** In `cli/cmd/roottooling.go`, the `constants.CmdHistoryPurge` registration was switched from `runHistoryPurge` to `RunHistoryPurgeCLI`. This shadowed the filter-repo mirror-clone rewrite command (`history-purge <path>`), caused `runHistoryPurge` to become dead code (flagged by `golangci-lint`), and broke the smoke test scenario expecting sandbox output.
2. **Unused Struct Field & Unexported Cross-Platform Helper:** `rewriteContext` had a leftover `fileRecords` slice from early development that was never populated or read. In `purge_recycle_other.go` and `purge_recycle_windows.go`, `sendToRecycleBin` was unexported and unreferenced inside package `cmdpurge`.
3. **Misspelling in Terminal Output:** `confirmPurgeInteractively` printed `Purge cancelled by user` using the non-US spelling `cancelled` instead of `canceled`.
4. **Hardcoded Drive Letter / URI Syntax:** Recent documentation files contained copy-pasted `file:///` links and Windows drive letter paths (`d:\<repo-root>`).

---

## 3. Resolution

1. **Restored Dispatch & Added Smart Delegation:**
   - Mapped `constants.CmdHistoryPurge` and `constants.CmdHistoryPurgeAlias` back to `runHistoryPurge` in `cli/cmd/roottooling.go`.
   - Updated `runHistoryPurge` in `cli/cmd/historyrewrite.go` to inspect incoming arguments for `--folder`, `--file`, or `--commit` flags and delegate cleanly to `RunHistoryPurgeCLI`, allowing both generic-cli positional paths and app-level splitdb flag modes.
   - Cleaned duplicate alias in `toolingSystemEntries()`.
2. **Removed Dead Field & Exported Cross-Platform Helper:**
   - Removed unused `fileRecords` field from `rewriteContext` in `cli/cmdpurge/purge_engine.go`.
   - Exported `SendToRecycleBin` in `cli/cmdpurge/purge_recycle_other.go` and `cli/cmdpurge/purge_recycle_windows.go`.
3. **Fixed Spelling:**
   - Changed `cancelled` to `canceled` in `cli/cmdpurge/purge_engine.go`.
4. **Sanitized Paths:**
   - Converted all absolute markdown links and drive letters in `.ai-memory/cicd-issues/105-*` and `repo-secrets/01-gitmap/01-chrome-import-export-test/readme.md` to relative repository paths.
   - Verified that `python linter-scripts/check-relative-paths.py` reports clean across 7,775 repository files.

---

## 4. Prevention & Learnings

- When adding new CLI subcommands that share name roots with existing tooling (e.g. `history purge` vs `history-purge`), never blindly overwrite existing dispatch entries without auditing test suites (`history-rewrite-smoke.yml`).
- Always run `python linter-scripts/check-relative-paths.py` before finalizing any documentation edits to catch OS-specific drive letters or `file:///` scheme references.
