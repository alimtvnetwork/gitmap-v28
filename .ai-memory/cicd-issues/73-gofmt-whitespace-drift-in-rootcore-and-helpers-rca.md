# CI/CD Issue 73: Gofmt Whitespace Drift in rootcore.go and helpers.go RCA

## 1. Why It Happened
CI run #35679200221 on commit `fa515c2` failed on the `Lint Script Unit Tests` job (`.github/scripts/tests/test_ci_scripts.py:171` `test_gofmt_check_clean_repo`):
```text
FAIL: test_gofmt_check_clean_repo (__main__.TestGoFormatCheck.test_gofmt_check_clean_repo)
AssertionError: 1 != 0
```
Running `.github/scripts/go-format-check.py --check-only` locally pinpointed:
```text
Dry run detected unformatted .go file(s):
  cmd\rootcore.go
  cmdschedule\helpers.go
Fix locally with: cd gitmap && gofmt -w .
```

---

## 2. How It Happened
- In commit `7a66ede4`, unused functions `checkHelp` in `cli/cmdschedule/helpers.go` and `resolveSCTopic` in `cli/cmd/rootcore.go` were deleted to resolve `golangci-lint` `unused` warnings.
- During deletion, extra consecutive blank lines remained in the files.
- While `golangci-lint` and `go vet` passed, `.github/scripts/tests/test_ci_scripts.py` executes `go-format-check.py --check-only` against all 3,280 Go files in `cli/`, which strictly requires `gofmt -w` compliance.

---

## 3. Root Cause
- `cli/cmd/rootcore.go`: Multiple consecutive blank lines between `coreClusterEntries` and `dispatchServersClients`.
- `cli/cmdschedule/helpers.go`: Multiple consecutive blank lines following the package import block.

---

## 4. Code Fix
- Formatted `cli/cmd/rootcore.go` and `cli/cmdschedule/helpers.go` using `gofmt -w`.
- Re-ran `python .github/scripts/go-format-check.py --check-only`, confirming 3,280 files are gofmt-clean (0 unformatted).
- Re-ran `python .github/scripts/tests/test_ci_scripts.py`, confirming 18/18 tests pass in 68s (`OK`).
- Ran all 4 policy linters (`check-nested-ifs.py`, `check-boolean-guidelines.py`, `check-relative-paths.py`, `check-enum-and-boolean.py`), all passing 100% green.
- Cross-platform `go vet` for Linux, macOS, and Windows passed with exit code 0.
- `go test -v ./cmdos/...` passed 100% green.
