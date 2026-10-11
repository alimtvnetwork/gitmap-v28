# RCA-109: File Size Gate 300-Line Cap and Unused Symbols in Pending Commits, Code, Create, and Ignore

**Date:** 2026-10-11
**Status:** ✅ Resolved (Commit `eebb9f6`, Release `v6.532.0`)
**Severity:** High (CI Pipeline Failure across File Size Gate, Lint Baseline Diff, and Lint Baseline Guard)
**Affected Workflows:** `CI (#38066317509)`
**Run URLs:**
- https://github.com/alimtvnetwork/gitmap-v28/actions/runs/38066317509

---

## 1. Symptom & Error Signature

Following commit `ab9ac4a`, GitHub Actions CI workflow `#38066317509` failed on three separate jobs:

1. **Job `File Size Gate (300-line cap)` (Step: `Check changed .go files against 300-line cap`):**
   ```text
   ❌ Go file-size check failed: 2 file(s) over 300 lines (8 scanned):
      • cli/cmdpending/pending_commits_cmd.go: 1182 lines (over by 882)
      • cli/cmdpending/pending_commits_test.go: 523 lines (over by 223)
   Split oversized files by concern (spec 243.4) before committing.
   Process completed with exit code 1.
   ```

2. **Job `Lint Baseline Diff` (Step: `Diff vs baseline (fail only on NEW findings)`):**
   ```text
   [unused] func `checkHelp` is unused (NEW in /tmp/lint-current/report.json)
       ❌ cmdcode/hook_shim.go:6: [unused] func `checkHelp` is unused
   [unused] func `applyWorkspaceSync` is unused (NEW in /tmp/lint-current/report.json)
       ❌ cmdcreate/create_ops.go:93: [unused] func `applyWorkspaceSync` is unused
   [unused] func `runIgnore` is unused (NEW in /tmp/lint-current/report.json)
       ❌ cmdignore/ignore_entry.go:9: [unused] func `runIgnore` is unused
   [unused] func `resolveAndInspectLocalRepos` is unused (NEW in /tmp/lint-current/report.json)
       ❌ cmdpending/pending_commits_cmd.go:399: [unused] func `resolveAndInspectLocalRepos` is unused
   [unused] func `inspectRepositoriesPendingCommits` is unused (NEW in /tmp/lint-current/report.json)
       ❌ cmdpending/pending_commits_cmd.go:508: [unused] func `inspectRepositoriesPendingCommits` is unused
   [unused] func `closeTaskDB` is unused (NEW in /tmp/lint-current/report.json)
       ❌ cmdpending/pendingtaskhelper.go:146: [unused] func `closeTaskDB` is unused
   Process completed with exit code 1.
   ```

3. **Job `Lint Baseline Guard (unused, gosec G115, misspell, gocritic, exhaustive)` (Step: `Unused diff (baseline-diff, full-path only)`):**
   ```text
   FAIL: 2 new unused finding(s) detected!
   ❌ cmdpending/pending_commits_cmd.go:399:6: [unused] func `resolveAndInspectLocalRepos` is unused
   ❌ cmdpending/pending_commits_cmd.go:508:6: [unused] func `inspectRepositoriesPendingCommits` is unused
   Process completed with exit code 1.
   ```

---

## 2. Root Cause

1. **Oversized Changed Go Files Exceeding 300-Line Cap:**
   - Commit `ab9ac4a` modified `cli/cmdpending/pending_commits_cmd.go` and `cli/cmdpending/pending_commits_test.go` to add GitHub owner/repo shorthand resolution and tabular tree rendering.
   - `cli/cmdpending/pending_commits_cmd.go` remained at 1048 lines (over by 748 lines) and `cli/cmdpending/pending_commits_test.go` remained at 467 lines (over by 167 lines).
   - The CI File Size Gate enforces that any `.go` file modified in the diff must not exceed 300 lines (`03-ai-scripts/52-file-size-check.py`).

2. **Unused Dead Functions in `cmdpending`:**
   - In `cli/cmdpending/pending_commits_cmd.go`, `resolveAndInspectLocalRepos` and `inspectRepositoriesPendingCommits` were superseded by `resolveAndInspectLocalReposWithDB` and `inspectRepositoriesWithDB`, but were left behind in the file as unexported dead functions.
   - In `cli/cmdpending/pendingtaskhelper.go`, `closeTaskDB` was replaced by the cleanup closure in `ensureDB`, leaving `closeTaskDB` orphaned and unused.

3. **Unused Functions Across Disjoint Shim and Entry Files:**
   - In `cli/cmdcode/hook_shim.go`, unexported `checkHelp` was declared but never referenced by `cmdcode` since `RunCode` directly invokes `cmdvscode.RenderVSCodeHelp()`.
   - In `cli/cmdcreate/create_ops.go`, `applyWorkspaceSync` was authored with conditional flags (`params.IsNoSync`, `params.IsNoDesktop`) but `tryPushRemoteProvisioned` caller invoked `workspacesync.SyncAll` directly.
   - In `cli/cmdignore/ignore_entry.go`, `runIgnore` was declared unexported and unused since `RunIgnoreCLI` in `ignore_cli.go` routes all `gitmap ignore` subcommands.

---

## 3. Surgical Fix

1. **Decompose `cli/cmdpending/pending_commits_cmd.go` under 300 Lines:**
   - `cli/cmdpending/pending_commits_cmd.go`: CLI options, flag parser, and entry dispatch (~240 lines).
   - `cli/cmdpending/pending_commits_inspect.go`: Repository inspection, git status porcelain parsing, unpushed commits, and branch version resolution (~260 lines).
   - `cli/cmdpending/pending_commits_filter.go`: Repository sorting, filtering, payload building, and ANSI/interactive breakdown rendering (~240 lines).
   - `cli/cmdpending/pending_commits_fleet.go`: SSH cluster node execution and fleet table rendering (~170 lines).
   - Remove dead functions `resolveAndInspectLocalRepos` and `inspectRepositoriesPendingCommits`.

2. **Decompose `cli/cmdpending/pending_commits_test.go` under 300 Lines:**
   - `cli/cmdpending/pending_commits_test.go`: Option parsing, sorting, and JSON schema tests (~240 lines).
   - `cli/cmdpending/pending_commits_inspect_test.go`: Status parsing, mock execution, version formatting, and remediation tests (~230 lines).

3. **Remediate Unused Symbols in Disjoint Packages:**
   - In `cli/cmdcode/hook_shim.go`: Remove unused `checkHelp`.
   - In `cli/cmdcreate/create_ops.go`: Wire `applyWorkspaceSync(params, absTarget)`.
   - In `cli/cmdignore/ignore_entry.go`: Remove unused `runIgnore`.
   - In `cli/cmdpending/pendingtaskhelper.go`: Remove unused `closeTaskDB`.

---

## 4. Verification

1. `python 03-ai-scripts/52-file-size-check.py --diff HEAD --limit 300` -> all modified files <= 300 lines (exit code 0).
2. `go test ./cli/cmdpending/... ./cli/cmdcode/... ./cli/cmdcreate/... ./cli/cmdignore/...` -> all tests pass.
3. `python 03-ai-scripts/05-guideline-autofixer.py cli/cmdpending/ --check-only` -> clean.
4. Continuous loop via `gitmap pe -t` until all CI checks are green.
