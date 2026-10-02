# RCA-097: Pipeline Typed Nil Interface, Swallowed DB Errors, Nested Ifs, and Relative Path Exclusions

**Date:** 2026-10-02
**Status:** ✅ Resolved
**Severity:** Critical (Multi-Job CI Pipeline Failure: Relative Path Check, Nested If Linter, Boolean & Enum Linter, Full Suite Guard / Cross-Platform Build)
**Affected Workflows:** `CI (#36914942299, #36956759998)`, `Cross-Platform Build (#36914941822)`
**Run URLs:**
- `https://github.com/alimtvnetwork/gitmap-v28/actions/runs/36914942299`
- `https://github.com/alimtvnetwork/gitmap-v28/actions/runs/36914941822`

---

## 1. Symptom

On commit `6d0bf51` on branch `main`, the CI and Cross-Platform Build pipelines failed across 9 jobs/steps:

1. **Staticcheck / Golangci-lint (`Full Suite Guard` & `Cross-Platform Build`):**
   ```
   store/migrations.go:343:10: SA4023: this value of err is never nil (staticcheck)
   store/migrations.go:359:10: SA4023: this value of err is never nil (staticcheck)
   store/migrations.go:374:10: SA4023: this value of err is never nil (staticcheck)
   ```
2. **Error Management Linter (`check-error-management.py`):**
   ```
   cli/searcher/search_history_db.go:121: Swallowed db exec error: _, _ = conn.Exec(purgeSearchSql)
   cli/store/migrations.go:501: Swallowed db exec error: _, _ = conn.Exec(...)
   ```
3. **Nested If Linter (`check-nested-ifs.py`):**
   ```
   cli/cmdpull/pull_remediation_hint.go:88: Nested 'if' detected (depth 2)
   cli/cmdpull/pull_efficient_render.go:142: Nested 'if' detected (depth 2)
   cli/cmd/nodes_cmd_test.go:210: Nested 'if' detected (depth 2)
   cli/cmdpipeline/pipeline_cache_eval.go:115: Nested 'if' detected (depth 2)
   cli/cmdpipeline/pipeline_history_ai.go:130, 220: Nested 'if' detected (depth 2)
   cli/cmd/rm.go:167: Nested 'if' detected (depth 2)
   ```
4. **Relative Path Check (`check-relative-paths.py`):**
   ```
   02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/01-architecture-spec.md: Hardcoded absolute repo path
   02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md: Hardcoded absolute repo path and file URI
   02-spec/22-app-issues/62-pull-all-duplicate-repositories-and-case-sensitivity-rca.md: Hardcoded absolute repo path
   .ai-memory/pipeline-ai/history-errors.md: Absolute path logs captured from old CI run outputs
   ```

---

## 2. Root Cause

1. **Typed Nil Interface in `cli/store/migrations.go`:** Helpers such as `insertRemappedGroupRepos`, `deleteStaleGroupRepos`, and `executeReleaseRemap` returned `apperror.WrapSimple(err, ...)` directly as `error` return type without checking `if err != nil`, which caused a non-nil interface `(type=*AppError, value=nil)` when `err == nil`, triggering `SA4023: this value of err is never nil`.
2. **Swallowed DB Errors:** `purgeSearchDatabase` and migration cleanup statements discarded `conn.Exec()` errors using `_, _ = conn.Exec(...)` instead of returning wrapped errors.
3. **Condition Nesting in Command Handlers:** Branch logic in `pull_remediation_hint.go`, `pull_efficient_render.go`, `nodes_cmd_test.go`, `pipeline_cache_eval.go`, `pipeline_history_ai.go`, and `rm.go` nested inner conditions rather than flattening with guard clauses and early returns.
4. **Historical Error Dump & Documentation Path Leaks:** Documentation specifications contained hardcoded absolute drive paths, and dynamically generated historical pipeline error logs in `.ai-memory/pipeline-ai/` were scanned by `check-relative-paths.py` because the directory was not gitignored or allowlisted.

---

## 3. Resolution

1. **Fixed Typed Nil Bug in `cli/store/migrations.go`:**
   - Modified all migration execution helpers to verify `if err != nil { return apperror.WrapSimple(err, ...) }` and return literal `nil` when `err == nil`.
2. **Fixed Swallowed Errors:**
   - Wrapped and checked `conn.Exec()` in `cli/searcher/search_history_db.go` and `cli/store/migrations.go`.
3. **Flattened Nested Conditionals:**
   - Inverted and extracted helper functions across `pull_remediation_hint.go`, `pull_efficient_render.go`, `nodes_cmd_test.go`, `pipeline_cache_eval.go`, `pipeline_history_ai.go`, and `cli/cmd/rm.go`.
4. **Sanitized Paths & Configured Allowlist:**
   - Removed absolute paths and `file:///` URIs from specification files in `02-spec/21-app/` and `02-spec/22-app-issues/`.
   - Added `ALLOWLIST_PREFIXES` to `linter-scripts/check-relative-paths.py` (`.ai-memory/pipeline-ai/`, `.ai-memory/cicd/`, `.ai-memory/temp/`).
   - Added `.ai-memory/pipeline-ai/` to `.gitignore` and removed generated historical files from git tracking.

---

## 4. Prevention & Learnings

- **Never return concrete pointer wrappers directly to `error` interface:** Always check `if err != nil { return wrap(err) }; return nil`.
- **Always check database execution return values:** Never discard errors via `_, _ = conn.Exec(...)`; wrap and propagate via `apperror`.
- **Enforce guard clause inversion:** Whenever checking conditions within conditions, decompose into discrete helper methods or invert to guard clauses.
- **Isolate dynamic AI dossiers:** Dynamic runtime dossiers and historical log archives belong in `.gitignore` and linter prefix allowlists.
