# CI/CD RCA 121: Unused Schema Registry Assertion Helpers in `cli/cmd/`

## 1. Issue Summary & Failure Telemetry
- **Failed Workflow Runs:**
  - CI: Run `#38159623635` (`Lint Baseline Diff` on commit `f81d531`)
- **Extracted Error Trace:**
  ```text
  [unused] func `assertSchemaKeysArray` is unused (NEW in /tmp/lint-current/report.json)
      ❌ cmd/schemaregistry_assert_test.go:44: [unused] func `assertSchemaKeysArray` is unused
  [unused] func `assertSchemaKeysSlice` is unused (NEW in /tmp/lint-current/report.json)
      ❌ cmd/schemaregistry_assert_test.go:84: [unused] func `assertSchemaKeysSlice` is unused
  [unused] func `readEveryObjectKeys` is unused (NEW in /tmp/lint-current/report.json)
      ❌ cmd/schemaregistry_assert_test.go:205: [unused] func `readEveryObjectKeys` is unused
  ```

---

## 2. Root Cause Analysis (4-Part RCA)

### A. Immediate Cause
The `Lint Baseline Diff` job failed with exit status 1 because `assertSchemaKeysArray`, `assertSchemaKeysSlice`, and `readEveryObjectKeys` in `cli/cmd/schemaregistry_assert_test.go` were uncalled within package `cmd` and flagged as new dead code by the `unused` analyzer.

### B. Underlying Mechanism & Propagation
During the modular split of schema registry test helpers, `assertSchemaKeysArray`, `assertSchemaKeysSlice`, and `readEveryObjectKeys` were created in `cli/cmd/schemaregistry_assert_test.go` anticipating callers in `cmd`. However, existing JSON snapshot tests in other packages (`cmdstartup`, `cmdscan`) already defined their own package-local helpers, leaving these three helpers uncalled in package `cmd`.

### C. Detection Latency & Tooling Defect
Standard `go test ./cmd` compiles and passes without checking for dead test helpers. The baseline diff linter detects unreferenced helpers across `_test.go` files and blocks the build.

### D. Prevention & Mitigation
Remove uncalled test helpers from `_test.go` files in package `cmd`. Verify with `go test` and baseline diff analysis before tagging releases.

---

## 3. Remediation & Code Fix
- Removed unused `assertSchemaKeysArray`, `assertSchemaKeysSlice`, and `readEveryObjectKeys` from `cli/cmd/schemaregistry_assert_test.go`.
- Retained the actively used `assertSchemaKeysFirstObject` (exercised by `cli/cmd/sshlistjson_contract_test.go`).
- Verified with `go test -run TestSchema ./cmd`: exits 0 with all contract tests passing.

---

## 4. Verification & Prevention
- Verified Go test suite in `cmd`: exits 0.
- Verified Go vet: exits 0.
- Verified file size gate: exits 0.
- Executed minor release bump to `v6.544.0` per continuous loop instruction:
  "If there is any error, fix that error and release again with the minor bump always."
