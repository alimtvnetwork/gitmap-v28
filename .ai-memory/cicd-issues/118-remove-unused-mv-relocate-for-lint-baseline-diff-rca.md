# CI/CD RCA 118: Remove Unused `cli/cmd/mv_relocate.go` for `Lint Baseline Diff`

## 1. Issue Summary & Failure Telemetry
- **Failed Workflow Runs:**
  - CI: Run `#38155311078` (`Lint Baseline Diff`)
- **Extracted Error Trace:**
  ```text
  GOLANGCI-LINT DIFF vs LAST SUCCESSFUL MAIN
  current  : /tmp/lint-current/report.json
  baseline : /tmp/lint-baseline/report.json

  + NEW       :    3
  - FIXED     :    3
  = UNCHANGED :   47

  NEW findings:
    cmd/mv_relocate.go
      + L14 [unused] func `executeMove` is unused
      + L30 [unused] func `syncExternalMove` is unused
      + L40 [unused] func `printMoveSuccess` is unused

  FIXED findings:
    cmdsummary/summary_pe_all.go
      - L48 [unused] func `runPipelineErrorsAllLegacy` is unused
      - L126 [unused] func `inspectPipelineForItem` is unused
      - L197 [unused] func `renderPipelineAllTerminal` is unused

  ##[error][unused] func `executeMove` is unused (NEW in /tmp/lint-current/report.json)
  ##[error][unused] func `syncExternalMove` is unused (NEW in /tmp/lint-current/report.json)
  ##[error][unused] func `printMoveSuccess` is unused (NEW in /tmp/lint-current/report.json)
  Process completed with exit code 1.
  ```

---

## 2. Root Cause Analysis (4-Part RCA)

### A. Immediate Cause
The CI workflow job `Lint Baseline Diff` runs `golangci-lint` and diffs findings against `/tmp/lint-baseline/report.json` via `.github/scripts/lint-diff.py`.
Three functions in `cli/cmd/mv_relocate.go` (`executeMove`, `syncExternalMove`, `printMoveSuccess`) were reported as NEW unused functions in `current` because they were not present in the baseline report.

### B. Underlying Mechanism & Propagation
`cli/cmd/mv_relocate.go` contained legacy duplicate move relocation logic in `package cmd`. The actual `gitmap mv` command dispatches directly to `cmdmerge.RunMove(argsTail())`, which uses the implementations in `package cmdmerge` (`cmdmerge/move.go` and `cmdmerge/helpers_shim.go`). The copies in `cli/cmd/mv_relocate.go` were uncalled dead code. When `cmdsummary/summary_pe_all.go` was removed, the diff against baseline caused `cmd/mv_relocate.go` to be flagged as new findings.

### C. Detection Latency & Tooling Defect
`mv_relocate.go` in `cmd/` had duplicated functions already present in `cmdmerge/`.

### D. Prevention & Mitigation
Remove dead and duplicate implementations when features are modularized into dedicated subpackages like `cmdmerge`.

---

## 3. Remediation & Code Fix
Deleted obsolete `cli/cmd/mv_relocate.go`. The production `gitmap mv` workflow in `cmdmerge` remains fully intact and operational.

---

## 4. Verification & Prevention
- Verified with `golangci-lint run --out-format=json --issues-exit-code=0`: 0 findings for `mv_relocate.go`.
- Verified with `go vet ./...`: exits 0 with zero warnings.
- Verified relative paths, newline styling, boolean guidelines, and enum linters: all pass.
- Executed minor release bump per continuous self-healing instructions.
