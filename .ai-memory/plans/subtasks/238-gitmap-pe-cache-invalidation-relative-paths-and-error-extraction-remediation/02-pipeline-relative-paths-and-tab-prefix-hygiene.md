# Subtask 02: Pipeline Relative Paths & Tab Prefix Hygiene

> **Subtask ID:** Subtask-02  
> **Parent Plan:** `.ai-memory/plans/238-gitmap-pe-cache-invalidation-relative-paths-and-error-extraction-remediation.md`  
> **Owned Files:**  
> - `cli/cmdpipeline/pipeline_logs.go`  
> - `cli/cmdpipeline/pipeline_details.go`  
> - `cli/cmdpipeline/pipeline_error_extract.go`  

---

## 1. Objectives

1. In `cli/cmdpipeline/pipeline_logs.go`:
   - Use `FormatRelativeDbPath` for `SavedLogFile` in `renderSingleFailureTerminal` and `renderFailedRunCard`.
   - In `queryWorkflowRunsForTarget`, resolve deep negative index offsets via git commit history.
2. In `cli/cmdpipeline/pipeline_details.go`:
   - Use `FormatRelativeDbPath` for `SavedLogFile` in `renderDetailedSectionFailure`.
3. In `cli/cmdpipeline/pipeline_error_extract.go`:
   - In `toRelativeGitPath`, call `FormatRelativeDbPath` when path is outside repository root.
   - Clean redundant `job\tstep\t` prefixes in failure summary formatting.
4. Run linters:
   - `python linter-scripts/check-nested-ifs.py`
   - `python linter-scripts/check-enum-and-boolean.py`
   - `python linter-scripts/check-relative-paths.py`
