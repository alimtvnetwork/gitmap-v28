# Subtask 04 — Testing, Verification & CI Quality Gates

**Parent Plan:** [.lovable/plans/pending/64-remediation-fix-and-chrome-token-export.md](../../pending/64-remediation-fix-and-chrome-token-export.md)  
**Status:** COMPLETED  
**Files:**
- `gitmap/gitutil/remediation_steps_test.go` [NEW]
- `gitmap/cmd/chromeprofile_tokens_test.go`
- `linter-scripts/check-nested-ifs.py`
- `linter-scripts/check-error-management.py`

---

## Objective
Author automated test coverage for remediation step generation and execution, token cipher reversible variations, roundtrip export/import, and pass all repository linter quality gates.

## Requirements
1. Unit tests for `RemediationStep` execution and error diagnostics.
2. Unit tests for double Base64, Caesar cipher text, and Caesar cipher byte shift reversibility.
3. Unit tests for profile export containing `TokenVault` and import restoring `token_service`.
4. Run `go test ./...` across modified packages.
5. Run nested if and error management linters with 0 violations.
