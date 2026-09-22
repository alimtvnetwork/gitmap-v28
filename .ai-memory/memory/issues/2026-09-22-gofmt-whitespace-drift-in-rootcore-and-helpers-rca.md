# RCA: Gofmt Whitespace Drift in rootcore.go and helpers.go

## 1. Why It Happened
CI run #35679200221 (`v6.299.0`) failed on the `Lint Script Unit Tests` job (`.github/scripts/tests/test_ci_scripts.py:171` `test_gofmt_check_clean_repo`):
```text
FAIL: test_gofmt_check_clean_repo (__main__.TestGoFormatCheck.test_gofmt_check_clean_repo)
AssertionError: 1 != 0
```
Running `.github/scripts/go-format-check.py --check-only` identified unformatted files:
- `cli/cmd/rootcore.go`
- `cli/cmdschedule/helpers.go`

## 2. How It Happened
When removing unused symbols `resolveSCTopic` and `checkHelp` in commit `7a66ede4`, consecutive blank lines remained in the edited files. `gofmt` collapses multiple blank lines into single blank lines.

## 3. Root Cause
`cli/cmd/rootcore.go` and `cli/cmdschedule/helpers.go` had extra blank lines violating `gofmt` canonical format.

## 4. Code Fix
- Formatted both files using `gofmt -w cli/cmd/rootcore.go cli/cmdschedule/helpers.go`.
- Verified `python .github/scripts/go-format-check.py --check-only` passes with 0 unformatted files.
- Verified `python .github/scripts/tests/test_ci_scripts.py` passes 18/18 tests.
- Verified all 4 policy linters pass.
- Verified cross-platform `go vet` passes on Linux, Darwin, and Windows.
