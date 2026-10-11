# CI/CD RCA 120: Unused Legacy Status Formatting in `cli/cmd/`

## 1. Issue Summary & Failure Telemetry
- **Failed Workflow Runs:**
  - CI: Run `#38158570691` (`Lint Baseline Diff` on commit `615574b`)
- **Extracted Error Trace:**
  ```text
  [unused] func `formatStateIcon` is unused (NEW in /tmp/lint-current/report.json)
      ❌ cmd/statusformat.go:44: [unused] func `formatStateIcon` is unused
  [unused] func `formatSyncText` is unused (NEW in /tmp/lint-current/report.json)
      ❌ cmd/statusformat.go:57: [unused] func `formatSyncText` is unused
  [unused] func `formatStashText` is unused (NEW in /tmp/lint-current/report.json)
      ❌ cmd/statusformat.go:86: [unused] func `formatStashText` is unused
  [unused] func `formatFileCounts` is unused (NEW in /tmp/lint-current/report.json)
      ❌ cmd/statusformat.go:97: [unused] func `formatFileCounts` is unused
  ```

---

## 2. Root Cause Analysis (4-Part RCA)

### A. Immediate Cause
The `Lint Baseline Diff` job failed with exit status 1 because `formatStateIcon`, `formatSyncText`, `formatStashText`, and `formatFileCounts` in `cli/cmd/statusformat.go` were flagged as new unused functions by the `unused` analyzer.

### B. Underlying Mechanism & Propagation
`cli/cmd/statusformat.go`, `cli/cmd/statusmissing.go`, and `cli/cmd/statusprint.go` were grandfathered legacy status printing files left behind when the status commands were migrated to `cli/cmdstatus/status.go`. Because the top-level status printers were never called by the active command dispatcher, the entire cluster of unexported helpers was detected as dead code by the baseline diff linter.

### C. Detection Latency & Tooling Defect
Local builds succeed because dead code is syntactically valid Go. The baseline diff linter runs against the cached main branch and flags any unreferenced identifiers.

### D. Prevention & Mitigation
Purge all unreferenced legacy status files from `cli/cmd/`:
- `cli/cmd/statusformat.go`
- `cli/cmd/statusmissing.go`
- `cli/cmd/statusprint.go`
- Remove unused `statusSummary` declaration from `cli/cmd/dbhelpers.go`.
Run `golangci-lint run` / baseline diff checks locally before tagging.

---

## 3. Remediation & Code Fix
- Deleted dead legacy files `cli/cmd/statusformat.go`, `cli/cmd/statusmissing.go`, and `cli/cmd/statusprint.go`.
- Cleaned up unreferenced `statusSummary` struct from `cli/cmd/dbhelpers.go`.
- Verified clean compilation with `go build -p 2 ./cmd` and `go vet ./...`.

---

## 4. Verification & Prevention
- Verified Go compilation: exits 0 with no errors.
- Verified Go vet: exits 0 with zero warnings.
- Verified file size gate: exits 0.
- Executed minor release bump to `v6.543.0` per continuous loop instruction:
  "If there is any error, fix that error and release again with the minor bump always."
