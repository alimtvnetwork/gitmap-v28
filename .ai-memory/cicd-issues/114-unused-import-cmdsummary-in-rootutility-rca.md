# CI/CD RCA 114: Unused Import `cmdsummary` in `cli/cmd/rootutility.go`

## 1. Issue Summary & Failure Telemetry
- **Failed Workflow Runs:**
  - Release run: `#38152077804` (Build release binaries)
  - Race Detector: `#38152062126` (`go test -race`)
  - Cross-Platform Build: `#38152062058` (`go build ./...`)
  - History Rewrite Smoke: `#38152061939` (Build gitmap binary)
- **Extracted Error Trace:**
  ```text
  cmd/rootutility.go:35:2: "github.com/alimtvnetwork/gitmap-v28/cli/cmdsummary" imported and not used
  Process completed with exit code 1.
  ```

---

## 2. Root Cause Analysis (4-Part RCA)

### A. Immediate Cause
The compiler halted compilation across all target platforms with exit status 1 due to an unused import statement on line 35 of `cli/cmd/rootutility.go`:
```go
"github.com/alimtvnetwork/gitmap-v28/cli/cmdsummary"
```
In Go, unused imports are compile-time errors.

### B. Underlying Mechanism & Propagation
During earlier refactoring of command dispatch utilities in commit `985f6e4`, the `cmdsummary` package reference was added or left behind without any call sites within `cli/cmd/rootutility.go`. When subsequent builds and CI actions invoked `go build ./...` and `go test -race`, the compiler rejected the package.

### C. Detection Latency & Tooling Defect
Local file-scoped checks on other packages did not evaluate the `cmd` package build graph until full module compilation `go build ./...` was executed across all targets.

### D. Prevention & Mitigation
Run full module compilation `go build ./...` across the entire `cli/` directory prior to release tagging to ensure zero compilation or unused symbol errors.

---

## 3. Remediation & Code Fix
Removed the unused `"github.com/alimtvnetwork/gitmap-v28/cli/cmdsummary"` import from `cli/cmd/rootutility.go`.

---

## 4. Verification & Prevention
- Verified with `go build ./...` across the entire `cli/` module: exits with code 0.
- Executed minor release bump to `v6.536.0` as per user instruction:
  "If there is any error, fix that error and release again with the minor bump always."
