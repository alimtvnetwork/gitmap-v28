# Issue 107: Pipeline Telemetry (pe -1) Nested Ifs, Boolean Guidelines, Relative Paths, and Misspell RCA

> **Issue ID:** `CICD-107`
> **Repository:** `gitmap-v28`
> **Status:** Resolved
> **Date:** 2026-10-08

---

## 1. Symptom

Failures detected in GitHub Actions runs for commit `d958c99` (workflows #37754871375 and #37754870863) via `gitmap pe -1`:

1. **Relative Path Check:**
   - 6 absolute repo path / URI violations detected in `02-spec/21-app/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/01-architecture-spec.md` (lines 94, 442, 444) and `cli/helptext/root.md` (lines 37, 38, 39).
2. **Boolean Guidelines Linter:**
   - 15 boolean guideline violations across 6 files:
     - Banned function prefix returning bool (`shouldSkipChildPathDir` in `cli/cmd/child_path.go:383`, `shouldCloneFromAnswer` in `cli/cmd/repo_create_smart.go:168`, `canFallbackToSafeRm` in `cli/cmd/rm.go:119`).
     - Negative boolean variable names (`hasNoArgs`, `hasNoTasks`, `hasNoFiles` in `cli/cmdai/ai_analysis_cmd.go:455, 466, 481, 540`, `hasNoFilter` in `cli/cmdai/ai_remove_revert.go:236`, `hasNoTasks` in `cli/cmdai/ai_train.go:399`).
3. **Nested If Linter & Boolean/Enum Linter:**
   - 40 nested-if / anti-compression violations (depth >= 2) across 10 files:
     - `cli/cmd/child_path.go` (7 violations)
     - `cli/cmd/repo_create_smart.go` (1 violation)
     - `cli/cmdagy/agm_crypto.go` (1 violation)
     - `cli/cmdagy/agm_deploy.go` (5 violations)
     - `cli/cmdai/ai_remove.go` (2 violations)
     - `cli/cmdai/ai_remove_revert.go` (4 violations)
     - `cli/cmdai/ai_train.go` (2 violations)
     - `cli/cmdai/ai_transfer.go` (2 violations)
     - `cli/cmdautomation/child_path.go` (7 violations)
     - `cli/cmdinstall/install_muse.go` (9 violations)
4. **Golangci-Lint Static Checks (`Full Suite Guard` & `Lint`):**
   - `cli/cmdai/ai_analysis_cmd.go:29:2`: var `analysisCmd` is unused (`unused`).
   - `cli/cmdai/ai_analysis_types.go:19:45`: `cancelled` is a misspelling of `canceled` (`misspell`).
5. **Lint Script Unit Tests:**
   - `FAIL: test_gofmt_check_clean_repo (__main__.TestGoFormatCheck.test_gofmt_check_clean_repo)` caused by unformatted files (`cli/cmdagy/agm_deploy.go` and `cli/cmdai/ai_remove.go`).

---

## 2. Root Cause

1. **Absolute Repo Path Hardcoding:** Sample snippets in `01-architecture-spec.md` and `root.md` used local development environment paths rather than generic/relative platform paths.
2. **Boolean Naming and Prefix Violations:** Functions returning booleans were given conversational prefixes (`should*`, `can*`) instead of the canonical `is*`/`has*` predicates mandated by coding guidelines. Variables holding emptiness checks were given negative names (`hasNo*`) violating positive boolean framing.
3. **Deep Conditional Branching:** Error checks and loop conditions were nested inside conditional blocks instead of being flattened via inverted guard clauses, early returns, or extracted single-responsibility helper functions.
4. **Dead Code & British Spelling:** Unused assignment `analysisCmd = AiAnalysisCmd` triggered `unused` lint; `AiTaskStatusCancelled` constant used British spelling `cancelled` triggering `misspell` in US English mode.
5. **Gofmt Whitespace Drift:** Automated or manual edits in `agm_deploy.go` and `ai_remove.go` lacked `gofmt -w` formatting.

---

## 3. Resolution

1. **Relative Paths Sanitization:** Replaced absolute repo paths in `01-architecture-spec.md` and `root.md` with portable relative/generic path examples. Verified `python linter-scripts/check-relative-paths.py` reports 100% clean across 7,970 files.
2. **Boolean Predicates & Positive Framing:**
   - Renamed `shouldSkipChildPathDir` -> `isChildPathDirSkipped`.
   - Renamed `shouldCloneFromAnswer` -> `isCloneConfirmedFromAnswer`.
   - Renamed `canFallbackToSafeRm` -> `isSafeRmFallbackAllowed`.
   - Replaced negative boolean variables (`hasNoArgs`, `hasNoTasks`, `hasNoFiles`, `hasNoFilter`) with `isEmpty`.
3. **Control Flow Flattening:**
   - Extracted flag parsing helpers in `cli/cmd/child_path.go` and `cli/cmdautomation/child_path.go` to eliminate nested error handling.
   - Refactored `agm_deploy.go`, `agm_crypto.go`, `ai_remove.go`, `ai_remove_revert.go`, `ai_train.go`, `ai_transfer.go`, and `install_muse.go` using early return guards and single-responsibility helper functions to bring nested if count to 0.
4. **Linter Remediation:**
   - Removed unused `analysisCmd = AiAnalysisCmd` variable.
   - Renamed `AiTaskStatusCancelled` to `AiTaskStatusCanceled` with value `"canceled"`.
5. **Code Formatting & Prompts Index Sync:**
   - Ran `gofmt -w` across all 14 modified Go source files (`gofmt -l .` output is clean).
   - Synchronized `.ai-memory/prompts.md` via `python linter-scripts/check-prompts-loaded.py --fix`.

---

## 4. Prevention & Learnings

- Always run `python linter-scripts/check-nested-ifs.py`, `python linter-scripts/check-boolean-guidelines.py`, and `gofmt -l .` before committing code changes.
- Ensure all new boolean predicates strictly use `is*` or `has*` prefixes, and variable names avoid `hasNo*` and `isNot*` prefixes.
- Keep documentation examples free of OS-specific absolute paths to ensure cross-platform compliance and prevent CI linter regression.
