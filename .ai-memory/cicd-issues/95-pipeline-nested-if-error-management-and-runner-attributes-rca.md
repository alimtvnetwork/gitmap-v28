# RCA-095: CI Pipeline Failures - Nested Ifs, Swallowed DB Errors, and CI Runner Attribute Regression

**Date:** 2026-10-02
**Status:** Resolved
**Severity:** Critical (CI Pipeline Failure across Error Management Linter, Nested If Linter, Boolean & Enum Linter, and Lint Script Unit Tests)
**Affected Workflows:** `CI` (`Error Management Linter`, `Nested If Linter`, `Boolean & Enum Linter`, `Lint Script Unit Tests`)
**Run ID:** `36887486133`

---

## 1. Root Cause Analysis
Four distinct failure categories caused the CI pipeline on commit `6e9ed360` to fail:

1. **Lint Script Unit Tests (`test_ci_scripts.py`):**
   - `03-ai-scripts/06-cicd-local-runner.py` was missing `CICD_DIR`, `normalize_repo_rel`, `format_banner_metadata`, and `format_log_locations_section`.
   - `JobResult` was defined as a rigid 5-field dataclass without the 6-positional argument constructor (`JobResult(name, cmd, return_code, out, err, duration)`), causing `TypeError` during test initialization.

2. **Error Management Linter (`check-error-management.py`):**
   - `cli/store/split_db_cache.go:334-335`: `scanRepoSummaryDetails` ignored scan errors using `_ = db.QueryRow(...).Scan(...)`. `linter-scripts/check-error-management.py` flagged this as swallowed database errors.

3. **Nested If Linter (`check-nested-ifs.py`):**
   - `cli/cmd/dopendingretry.go:132`: Nested `if isIgnorablePendingTaskError(err)` inside `if err != nil`.
   - `cli/cmd/pendingtaskhelper.go:97`: Nested `if isIgnorablePendingTaskError(err)` inside `if err != nil`.
   - `cli/cmdignore/ignore_cmd.go:175`: Nested `if days == 1` inside `if days >= 1 && ...`.
   - `cli/store/pendingtask.go:115`: Nested `if isPendingTaskAlreadyCompleted(...)` inside `if appErr != nil`.

4. **Boolean & Enum Linter (`check-enum-and-boolean.py`):**
   - Flagged the identical 4 nested-if violations under its condition depth check rules.

---

## 2. Why Not Caught Earlier?
The changes to `cli/store/split_db_cache.go`, `cli/store/pendingtask.go`, `cli/cmdignore/ignore_cmd.go`, and `cli/cmd/pendingtaskhelper.go` were introduced during recent fast refactoring commits without running local linter scripts prior to push. In addition, an external sync reverted `06-cicd-local-runner.py` path attributes.

---

## 3. Remediation
1. **Restored `03-ai-scripts/06-cicd-local-runner.py` Attributes & Constructor:**
   - Re-added `REPO_ROOT`, `CICD_DIR`, `normalize_repo_rel`, `format_banner_metadata`, `format_log_locations_section`.
   - Restored the backward-compatible 6-argument `JobResult.__init__`.
2. **Proper Error Handling in `cli/store/split_db_cache.go`:**
   - Captured and checked errors from `RepoMetadata` scans, ignoring expected `sql.ErrNoRows` cleanly.
3. **Guard Clauses / Inverted Returns in Nested Ifs:**
   - `cli/cmd/dopendingretry.go`: Inverted to `if isIgnorablePendingTaskError(err) { return }`.
   - `cli/cmd/pendingtaskhelper.go`: Inverted to `if isIgnorablePendingTaskError(err) { return }`.
   - `cli/cmdignore/ignore_cmd.go`: Decomposed with `hasExactDays` guard clause.
   - `cli/store/pendingtask.go`: Reordered checks so `isPendingTaskAlreadyCompleted` guards early.

---

## 4. Verification
- `python -m unittest .github.scripts.tests.test_ci_scripts` (TestCicdLocalRunnerPaths): **PASS** (4/4 passed).
- `python linter-scripts/check-error-management.py`: **PASS** (zero swallowed errors across 4,163 files).
- `python linter-scripts/check-nested-ifs.py`: **PASS** (zero violations across 42 files).
- `python linter-scripts/check-enum-and-boolean.py`: **PASS** (zero violations across 3,128 files).
- `python .github/scripts/go-format-check.py --check-only`: **PASS** (3,874 files clean).
- `go test ./cmdignore/... ./cmdpull/... ./store/...` (in `cli/`): **PASS**.
- `go build ./...` (in `cli/`): **PASS**.

---

## 5. Prevention
- Run `check-nested-ifs.py` and `check-error-management.py` locally before committing Go store/cmd updates.
- Keep runner path helper definitions in `06-cicd-local-runner.py` locked against blanket sync overwrites.
