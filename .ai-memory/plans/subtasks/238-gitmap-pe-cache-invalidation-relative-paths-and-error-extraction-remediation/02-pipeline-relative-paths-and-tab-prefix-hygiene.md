# Subtask 02: Pipeline Relative Paths & Tab Prefix Hygiene

> **Subtask ID:** Subtask-02  
> **Parent Plan:** `.ai-memory/plans/238-gitmap-pe-cache-invalidation-relative-paths-and-error-extraction-remediation.md`  
> **Owned Files:**  
> - `cli/cmdpipeline/pipeline_logs.go`  
> - `cli/cmdpipeline/pipeline_details.go`  
> - `cli/cmdpipeline/pipeline_error_extract.go`  

---

## 1. Objectives

- [x] 1. In `cli/cmdpipeline/pipeline_logs.go`:
   - Use `FormatRelativeDbPath` for `SavedLogFile` in `renderSingleSectionFailureRow` and `renderRunCardBranchAndLog`.
   - In `queryWorkflowRunsForTarget`, resolve deep negative index offsets via git commit history.
   - Use `cleanDisplayErrorText` in `renderSingleSectionFailureRow` and `renderFailedJobSection`.
- [x] 2. In `cli/cmdpipeline/pipeline_details.go`:
   - Use `FormatRelativeDbPath` for `SavedLogFile` in `renderSingleDetailsFailure`.
   - Use `cleanDisplayErrorText` in `renderSingleDetailsFailure`.
- [x] 3. In `cli/cmdpipeline/pipeline_error_extract.go`:
   - In `toRelativeGitPath`, call `FormatRelativeDbPath` when path is outside repository root.
- [x] 4. Run linters:
   - `python linter-scripts/check-nested-ifs.py`: PASS (0 violations).
   - `python linter-scripts/check-enum-and-boolean.py`: PASS (0 violations).
   - `python linter-scripts/check-relative-paths.py`: PASS (0 violations across 7,828 files).

---

## 2. Verification Evidence

- `check-nested-ifs.py`: PASS (0 violations across candidate files).
- `check-enum-and-boolean.py`: PASS (0 violations across 3,228 files).
- `check-relative-paths.py`: PASS (0 violations across 7,828 files).
- `gitmap pe` output verified in terminal: All log paths are rendered as `data/pipeline/...` without any `C:/Users/...` leakage.
- Status: **DONE**
