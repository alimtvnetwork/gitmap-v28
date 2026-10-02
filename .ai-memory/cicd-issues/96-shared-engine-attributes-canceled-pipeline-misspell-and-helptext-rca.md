# RCA-096: CI Pipeline Failures - Shared Engine Attributes, Canceled Pipeline Workflows, US Misspell, and Missing Helptext

**Date:** 2026-10-02
**Status:** Resolved
**Severity:** Critical (CI Pipeline Failure across Spell Check, Lint Script Unit Tests, Full Suite Guard, Relative Path Check, Nested If Linter, Boolean & Enum Linter, and Cross-Platform Build)
**Affected Workflows:** `CI` (`Spell Check`, `Lint Script Unit Tests`, `Full Suite Guard`, `Relative Path Check`, `Nested If Linter`, `Boolean & Enum Linter`), `Cross-Platform Build` (`macos-latest`, `ubuntu-latest`, `windows-latest`)
**Run ID:** `36914942299` / `36914941822`

---

## 1. Symptoms & Error Reports
From GitHub Actions run `36914942299` and `36914941822` on commit `6d0bf51`:

1. **Spell Check (misspell, US locale):**
   ```text
   changelog.md:21:11: "cancelled" is a misspelling of "canceled"
   Process completed with exit code 2.
   ```

2. **Lint Script Unit Tests (`test_ci_scripts.py`):**
   ```text
   ERROR: test_cicd_dir_is_inside_repo (__main__.TestCicdLocalRunnerPaths)
   AttributeError: module '06-cicd-local-runner' has no attribute 'CICD_DIR'
   ERROR: test_format_banner_metadata_contains_no_temp_or_absolute_drive
   TypeError: JobResult.__init__() takes 6 positional arguments but 7 were given
   ERROR: test_format_log_locations_section_all_inside_repo
   AttributeError: module '06-cicd-local-runner' has no attribute 'CICD_DIR'
   ERROR: test_normalize_repo_rel_converts_absolute_paths
   AttributeError: module '06-cicd-local-runner' has no attribute 'normalize_repo_rel'
   ```

3. **Policy Linters (`check-relative-paths.py`, `check-nested-ifs.py`, `check-enum-and-boolean.py`):**
   ```text
   AttributeError: module '02-shared-engine' has no attribute 'chunk_items'
   AttributeError: module '02-shared-engine' has no attribute 'WorkerHeartbeatMonitor'
   ```

4. **Go Test Suite in `cli/cmdpipeline/`:**
   ```text
   --- FAIL: TestGroupRunsByCommit_CanceledWorkflows (0.00s)
       pipeline_commit_groups_test.go:76: unexpected first group conclusion: sha=12e40b1, conclusion=success
   --- FAIL: TestFormatStatusBadge_CanceledAndFailing (0.00s)
       pipeline_history_test.go:436: expected badge "FAIL" for conclusion="canceled" status="completed", got "canceled"
   --- FAIL: TestFormatGroupWorkflowsSummary_PrioritizesFailures (0.00s)
       pipeline_history_test.go:462: expected failing workflow to be prioritized at front, got: "Release [PASS] (+2)"
   --- FAIL: TestIsCommitGroupFailure_Canceled (0.00s)
       pipeline_history_test.go:490: expected isCommitGroupFailure to be true for canceled conclusion
   ```

---

## 2. Root Cause Analysis
Four distinct causes triggered this multi-suite cascade:

1. **External Sync Overwrite:**
   Commit `d6260d40` synced external files and inadvertently reverted `03-ai-scripts/02-shared-engine.py` and `03-ai-scripts/06-cicd-local-runner.py`. This stripped `chunk_items`, `WorkerHeartbeatMonitor`, `CICD_DIR`, `normalize_repo_rel`, `format_banner_metadata`, and `format_log_locations_section`, crashing all dependent linters and unit tests.

2. **Accidental Removal of Cancelled Runs from `isFailingConclusion`:**
   Commit `5438cf18` attempted to prevent canceled runs from showing false-positive failure alerts during pull, but stripped `strings.HasPrefix(lower, "cancel")` from `isFailingConclusion`. In commit groups and pipeline history, canceled runs represent incomplete/failed workflows; removing it caused canceled runs to be treated as `success`, breaking 4 unit tests in `cli/cmdpipeline`.

3. **US Locale Spelling Violation:**
   `changelog.md` line 21 used British English `cancelled` instead of US English `canceled`.

4. **Missing Helptext & Topic Metadata:**
   The user requested complete helptext for newly added commands (`gitmap pull-error`, `gitmap devtools`, `gitmap ignore`). `ignore.md` was an empty 12-line stub, and `pull-error.md` / `devtools.md` were unregistered in `topicSummaries`.

---

## 3. Resolution Applied
1. **Restored `03-ai-scripts/02-shared-engine.py`:**
   Added `chunk_items(items, chunk_size)` and `WorkerHeartbeatMonitor` class with thread-safe snapshot emission.

2. **Restored `03-ai-scripts/06-cicd-local-runner.py`:**
   Restored `REPO_ROOT`, `CICD_DIR`, `normalize_repo_rel`, `format_banner_metadata`, `format_log_locations_section`, and flexible `JobResult` constructor.

3. **Restored Canceled Handling in `cli/cmdpipeline/pipeline_logs.go`:**
   Re-added `if strings.HasPrefix(lower, "cancel") { return true }` to `isFailingConclusion`.

4. **Fixed Spelling in `changelog.md`:**
   Replaced `cancelled` with `canceled`.

5. **Authored Comprehensive Help Text:**
   - Authored `cli/helptext/ignore.md` covering all subcommands (`config`, `cache`), TTL configuration, Split-DB caching, and flags.
   - Authored `cli/helptext/pull-error.md` and `cli/helptext/pulle.md`.
   - Authored `cli/helptext/devtools.md` and `cli/helptext/devtool.md`.
   - Registered all new topics in `topicSummaries` in `cli/helptext/catalog.go`.

---

## 4. Prevention & Learnings
- **Pre-Sync Preservation:** Before running mass external synchronization scripts, critical shared engine utilities (`chunk_items`, `WorkerHeartbeatMonitor`) must be protected or re-verified.
- **US Spelling Linter Pre-Commit Gate:** Always verify `misspell` on changed markdown files before committing release notes.
- **Contract Tests for Pipeline Enums:** `isFailingConclusion` must respect contract tests in `pipeline_history_test.go` and `pipeline_commit_groups_test.go`.
