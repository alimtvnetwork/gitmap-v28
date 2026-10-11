# CI/CD RCA 122: Unformatted Go Code in `cli/cmd/schemaregistry_assert_test.go`

## 1. Issue Summary & Failure Telemetry
- **Failed Workflow Runs:**
  - CI: Run `#38160410685` (`Lint Script Unit Tests` on commit `5437403`)
- **Extracted Error Trace:**
  ```text
  TestGoFormatCheck.test_gofmt_check_clean_repo:
  Dry run detected unformatted .go file(s):
    cmd/schemaregistry_assert_test.go
  FAIL: test_gofmt_check_clean_repo
  ```

---

## 2. Root Cause Analysis (4-Part RCA)

### A. Immediate Cause
The CI step `Run lint-script unit tests` failed during `TestGoFormatCheck.test_gofmt_check_clean_repo` because `cli/cmd/schemaregistry_assert_test.go` had non-canonical indentation/blank lines that failed `gofmt -l`.

### B. Underlying Mechanism & Propagation
When editing `cli/cmd/schemaregistry_assert_test.go` in commit `5437403` to remove dead helpers, manual removals introduced whitespace deviations that were not automatically formatted with `gofmt -w`. The test suite under `.github/scripts/tests/test_ci_scripts.py` asserts full-tree `gofmt` compliance, failing the unit test.

### C. Detection Latency & Tooling Defect
`python .github/scripts/go-format-check.py --check-only` was not run locally on the modified file before tagging `v6.544.0`.

### D. Prevention & Mitigation
Run `gofmt -w` and `python .github/scripts/go-format-check.py --check-only` locally on all modified `.go` files before staging and releasing.

---

## 3. Remediation & Code Fix
- Formatted `cli/cmd/schemaregistry_assert_test.go` with `gofmt -w`.
- Verified with `python .github/scripts/tests/test_ci_scripts.py TestGoFormatCheck`: exits 0 in 0.4s.
- Verified with `python .github/scripts/go-format-check.py --check-only`: 4499 Go files scanned, 0 unformatted.

---

## 4. Verification & Prevention
- Verified with `go-format-check.py`: exits 0.
- Verified file size gate: exits 0.
- Executed minor release bump to `v6.545.0` per continuous loop instruction:
  "If there is any error, fix that error and release again with the minor bump always."
