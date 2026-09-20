# RCA 65: CI/CD gofmt Drift, Nested If Linter, and locate_test Directory Probe

## 1. Symptom
GitHub Actions CI run `#35511230421` on commit `a3074bb` failed across 7 pipeline sections:
1. `Lint Script Unit Tests`:
   - `FAIL: test_gofmt_check_clean_repo (__main__.TestGoFormatCheck)` (1 unformatted file in `cli/apperror/apperror_test.go`).
2. `Boolean & Enum Linter` / `Nested If Linter`:
   - `cli/cmdpipeline/pipeline_ai_terminal.go:35: Nested 'if' detected (depth 2): 'if len(fj.FailureSummary) > 0 {'`
   - `cli/cmdpipeline/pipeline_ai_errors.go:68: Nested 'if' detected (depth 2): 'if len(items) > 0 {'`
3. `Full Suite Guard` & `Cross-Platform Build` (Windows, macOS, Ubuntu):
   - `--- FAIL: TestScanRootForFile (0.00s)`
   - `locate_test.go:42: expected sample_tool.exe to be found in scanRootForFile`

## 2. Root Cause
1. **gofmt whitespace drift:** An extra trailing newline was present at the end of `cli/apperror/apperror_test.go`, which violated strict `gofmt` compliance verified by `test_ci_scripts.py`.
2. **Nested if statements:**
   - In `cli/cmdpipeline/pipeline_ai_terminal.go`, `if len(fj.FailureSummary) > 0` was nested inside `if len(p.FailedJobs) > 0`.
   - In `cli/cmdpipeline/pipeline_ai_errors.go`, `if len(items) > 0` was nested inside `if len(rawLogs) > 0`.
3. **Directory vs File existence bug:**
   - In `cli/cmdautomation/locate.go`, `pathExists(p)` was explicitly defined as `err == nil && !info.IsDir()` (i.e. only true for non-directories).
   - In `scanRootForFile(root, filename, maxDepth)`, the function guarded entry with `if !pathExists(root) { return "", false }`. Because `root` is a directory, `pathExists(root)` evaluated to `false`, causing `scanRootForFile` to always abort immediately.

## 3. Resolution
1. **Directory probe:** Added `dirExists(p string) bool` returning `err == nil && info.IsDir()`, and updated `scanRootForFile` to use `!dirExists(root)`.
2. **Nested if flattening:**
   - Extracted `renderSingleFailedJobInfo(fj FailedJobItem)` in `cli/cmdpipeline/pipeline_ai_terminal.go`.
   - Guarded against empty `rawLogs` with an early return in `resolveFailedJobItems` in `cli/cmdpipeline/pipeline_ai_errors.go`.
3. **gofmt formatting:** Formatted `cli/apperror/apperror_test.go` with `gofmt -w`.

## 4. Verification
- `python linter-scripts/check-nested-ifs.py`: Passed cleanly with zero nested ifs across 3,361 files.
- `python linter-scripts/check-enum-and-boolean.py`: Passed cleanly.
- `python linter-scripts/check-boolean-guidelines.py`: Passed cleanly.
- `python .github/scripts/go-format-check.py`: Passed with 0 unformatted files across all 3,143 Go files.
- `python .github/scripts/tests/run_tests.py`: All 18 tests passed in 57s.
- `go test -v ./cmdautomation/...`: All unit tests passed (including `TestScanRootForFile`).
