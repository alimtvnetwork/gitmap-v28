# CI/CD RCA 119: File Size Gate Failure on `cli/cmdpipeline/pipeline_db_ops.go` (300-Line Cap)

## 1. Issue Summary & Failure Telemetry
- **Failed Workflow Runs:**
  - CI: Workflow job "File Size Gate (300-line cap)" on commit `13fef34` (CI Run `#38157934570`)
- **Extracted Error Trace:**
  ```text
  ❌ Go file-size check failed: 1 file(s) over 300 lines (5 scanned):
     • cli/cmdpipeline/pipeline_db_ops.go: 490 lines (over by 190)
  Split oversized files by concern (spec 243.4) before committing.
  Error: Process completed with exit code 1.
  ```

---

## 2. Root Cause Analysis (4-Part RCA)

### A. Immediate Cause
The pre-commit / CI gate `python3 03-ai-scripts/52-file-size-check.py --diff HEAD~1 --limit 300` failed because `cli/cmdpipeline/pipeline_db_ops.go` contained 490 lines, exceeding the 300-line cap by 190 lines.

### B. Underlying Mechanism & Propagation
`cli/cmdpipeline/pipeline_db_ops.go` was historically grandfathered as an oversized file. In commit `13fef34` (`fix(pipeline): prevent false failures on in-progress runs, invalidate hollow pe-all cache, and support clear all`), additions for `hasPositionalAll` and pe-all cache clearing were made to `pipeline_db_ops.go`. Because the file was modified in the commit diff, the file size gate evaluated it and failed.

### C. Detection Latency & Tooling Defect
The file size check evaluates changed files in the diff. `52-file-size-check.py` was not executed locally against the staged changes prior to pushing commit `13fef34`.

### D. Prevention & Mitigation
Decompose `cli/cmdpipeline/pipeline_db_ops.go` into two cohesive submodules strictly under 300 lines each:
1. `cli/cmdpipeline/pipeline_db_ops.go` (236 lines): CLI handlers, user confirmations, summary renderers, database optimize and error logs query/formatting.
2. `cli/cmdpipeline/pipeline_db_purge.go` (264 lines): Directory candidates discovery, orphan purge, recursive deletion, and pipeline file eligibility.
Always run `python 03-ai-scripts/52-file-size-check.py --files ...` locally prior to committing.

---

## 3. Remediation & Code Fix
- Extracted filesystem purging, orphan scanning, and path discovery into `cli/cmdpipeline/pipeline_db_purge.go`.
- Retained core DB operations, clear orchestration, reset, optimize, and error log viewing in `cli/cmdpipeline/pipeline_db_ops.go`.
- Both files now measure $\le 264$ lines, fully satisfying spec 243.4.

---

## 4. Verification & Prevention
- Verified with `python 03-ai-scripts/52-file-size-check.py --files cli/cmdpipeline/pipeline_db_ops.go cli/cmdpipeline/pipeline_db_purge.go`: exits 0 with 0 violations.
- Verified with `go vet ./cmdpipeline/...` and `go vet ./...`: exits 0 with zero warnings.
- Verified with `python linter-scripts/check-relative-paths.py`: exits 0 with no absolute paths.
- Verified unit tests in `cmdpipeline`: all pass cleanly.
