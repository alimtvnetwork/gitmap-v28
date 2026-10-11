# CI/CD RCA 115: File Size Gate Failure on `cli/cmd/rootutility.go` (300-Line Cap)

## 1. Issue Summary & Failure Telemetry
- **Failed Workflow Runs:**
  - Fast Gate / CI: Workflow job "File Size Gate (300-line cap)" on commit `186d79d` (release `v6.536.0`)
- **Extracted Error Trace:**
  ```text
  ❌ Go file-size check failed: 1 file(s) over 300 lines (1 scanned):
     • cli/cmd/rootutility.go: 559 lines (over by 259)
  Split oversized files by concern (spec 243.4) before committing.
  Error: Process completed with exit code 1.
  ```

---

## 2. Root Cause Analysis (4-Part RCA)

### A. Immediate Cause
The pre-commit / CI gate `python3 03-ai-scripts/52-file-size-check.py --diff HEAD~1 --limit 300` failed because `cli/cmd/rootutility.go` contained 559 lines, exceeding the 300-line cap by 259 lines.

### B. Underlying Mechanism & Propagation
`cli/cmd/rootutility.go` was historically grandfathered as an oversized file (>300 lines). When the prior fix for CI RCA 114 touched `cli/cmd/rootutility.go` to remove an unused import, the diff ratchet triggered: any `.go` file modified in the latest commit must strictly satisfy the 300-line size limit. Because `cli/cmd/rootutility.go` was part of the commit diff, the file size checker evaluated it and blocked CI completion.

### C. Detection Latency & Tooling Defect
The file size check runs against modified files (`--diff HEAD~1`). When removing the import, `52-file-size-check.py --diff HEAD` was not run prior to committing `v6.536.0`.

### D. Prevention & Mitigation
Decompose `cli/cmd/rootutility.go` into focused, cohesive submodules under 300 lines each:
1. `cli/cmd/rootutility.go` (269 lines): utility dispatch tables, categorization, and help routing.
2. `cli/cmd/rootutility_update.go` (138 lines): update command execution and orchestration.
3. `cli/cmd/rootutility_update_args.go` (189 lines): update argument parsing, validation, suggestions, and usage help.
Always run `python 03-ai-scripts/52-file-size-check.py --diff HEAD` locally prior to version bumps and tagging.

---

## 3. Remediation & Code Fix
- Extracted update execution flow into `cli/cmd/rootutility_update.go`.
- Extracted update argument parsing and suggestions into `cli/cmd/rootutility_update_args.go`.
- Removed dead imports (`cmdinstaller`, `store`) from `cli/cmd/rootutility.go`.
- All three files now measure $\le 269$ lines, fully satisfying spec 243.4.

---

## 4. Verification & Prevention
- Verified with `python 03-ai-scripts/52-file-size-check.py --files cli/cmd/rootutility.go cli/cmd/rootutility_update.go cli/cmd/rootutility_update_args.go`: exits 0 with 0 violations.
- Verified with `go build -p 2 ./cmd` and unit tests in `cmdpull` and `cmdpending`: all pass cleanly.
- Bumped version to `v6.537.0` per user instruction:
  "If there is any error, fix that error and release again with the minor bump always."
