# RCA: Pipeline Summary Table Release Column Offset and gofmt Formatting Drift

- **Date:** 2026-09-21
- **Affected Run:** [CI #35584286561](https://github.com/alimtvnetwork/gitmap-v28/actions/runs/35584286561) and [Cross-Platform Build #35584286225](https://github.com/alimtvnetwork/gitmap-v28/actions/runs/35584286225)
- **Target Commit:** `3a629d90` (Release `v6.291.0`)
- **Scope:** CI/CD Quality Gates (`go test` unit tests, `test_ci_scripts.py`, and `gofmt`)

---

## 1. Symptom

Pipeline runs [#35584286561](https://github.com/alimtvnetwork/gitmap-v28/actions/runs/35584286561) and [#35584286225](https://github.com/alimtvnetwork/gitmap-v28/actions/runs/35584286225) failed across multiple jobs:
1. **Full Suite Guard** (`go test ./... -count=1`):
   ```text
   --- FAIL: TestRecentCommitsSummaryTableAlignment (0.00s)
       pipeline_history_test.go:313: expected Status column to be PASS, got -
   FAIL
   FAIL github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline 0.557s
   ```
2. **Lint Script Unit Tests** (`python3 .github/scripts/tests/test_ci_scripts.py`):
   `test_gofmt_check_clean_repo` failed because `go-format-check.py --check-only` detected unformatted file:
   ```text
   ::notice::Dry run detected unformatted .go file(s):
     cmdpipeline\pipeline_history_test.go
   Fix locally with: cd gitmap && gofmt -w .
   ```
3. **Cross-Platform Build** (`ubuntu-latest` and `macos-latest`):
   Failed in `go test ./... (no cache)` due to `TestRecentCommitsSummaryTableAlignment`.

---

## 2. Root Cause

1. **Table Column Offset Shift**: In commit `3a629d90`, a new `Release` column was added to the Recent Commits Pipeline Summary table between `Branch` and `Status`:
   `Offset [0]`, `Commit [1]`, `Branch [2]`, `Release [3]`, `Status [4]`, `Workflows [5]`, `Failures [6]`.
   However, `TestRecentCommitsSummaryTableAlignment` in `cli/cmdpipeline/pipeline_history_test.go` was still asserting against the old 6-column layout:
   - It checked `len(parts) < 6` instead of `len(parts) < 7`.
   - It asserted `parts[3] != "PASS"`. Since `parts[3]` is now the `Release` column (value `"-"` for unreleased test commits), the test failed with `expected Status column to be PASS, got -`.
2. **gofmt Formatting Drift**: `cli/cmdpipeline/pipeline_history_test.go` had not been formatted with `gofmt -w` prior to commit `3a629d90`, causing the `go-format-check.py` baseline check to fail in `test_ci_scripts.py`.

---

## 3. Resolution

1. **Update Test Alignment Assertions**:
   - In `cli/cmdpipeline/pipeline_history_test.go`:
     - Updated column count check from `len(parts) < 6` to `len(parts) < 7`.
     - Added validation that `parts[3] == "-"` (asserting correct release column formatting).
     - Updated status badge check from `parts[3]` to `parts[4] != "PASS"`.
2. **Apply gofmt Formatting**:
   - Executed `gofmt -w cli/cmdpipeline/pipeline_history_test.go`.
   - Verified `go-format-check.py --check-only` confirms 3,157 files checked with 0 unformatted files.
3. **Local Test & Linter Verification**:
   - Ran `go test ./cmdpipeline -v -run "Test.*History|TestRecent|TestFormat|TestExtract"` - all passed.
   - Ran `python linter-scripts/check-relative-paths.py` - passed across 7,191 files.
   - Ran `python linter-scripts/check-nested-ifs.py --all` - passed across 3,375 files.
   - Ran `python linter-scripts/check-boolean-guidelines.py` - passed.
   - Ran `python linter-scripts/check-newline-styling.py --all` - passed.

---

## 4. Prevention & Learnings

- **Column Change Audit**: Whenever a table layout or column ordering is modified in CLI renderers, immediately update all associated unit test assertions and column index offsets.
- **Pre-Release gofmt Check**: Always run `python .github/scripts/go-format-check.py --check-only` before tagging and pushing a release to catch any unformatted Go files.
