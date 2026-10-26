# Subtask 332.1 — Show Log Folder Path Normalization and Single-Path Diff

> **Parent Plan:** `.ai-memory/plans/pending/332-gitmap-ui-and-functionality-improvements.md`  
> **Target Packages:** `cli/cmdlogviewer/`, `cli/cmddiff/`  
> **Status:** Pending  

---

## 1. Objectives
1. In `cli/cmdlogviewer/server.go`:
   - Normalize target path in `ResolveTargetRepoAndPath`: strip trailing slashes, resolve directory vs file correctly, and sanitize relative paths.
   - In `QueryCommits()` (`git_log_parser.go`), ensure directory path filtering uses normalized relative path from repository root.
2. In `cli/cmddiff/dispatchdiff.go` & `cli/cmddiff/diff.go`:
   - Add single-argument diff support: `gitmap diff <path>` or `gitmap diff` (zero args).
   - If single path or zero args provided inside a Git repo, launch working tree diff via `cmdgit.FindCompareTool()` or output formatted working tree diff with complete `apperror` stack traces.
3. In `cli/cmdlogviewer/server_handlers.go`:
   - When `/api/log/diff` fails, return structured `DiffResponse` with `stackTrace` and clear error diagnostics.
4. Add unit tests in `cli/cmddiff/cmddiff_test.go` and `cli/cmdlogviewer/server_handlers_test.go`.
