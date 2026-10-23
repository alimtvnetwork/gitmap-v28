# RCA 126: CI Pipeline Failures - Legacy Refs Check, Gofmt Cleanliness & Nested Ifs

## Part 1: Symptom
In GitHub Actions CI workflow (Run #39978916615 on commit `0fe54aa` / `94a294c`), two jobs failed:
1. `Legacy Refs Check`: Step `Run ./.github/actions/policy-check` failed because `.github/scripts/check-legacy-refs.py .` found 3 legacy reference violations in `cli/cmdupdate/update_binary_swap_test.go`.
2. `Lint Script Unit Tests`: Step `Run lint-script unit tests` (`test_ci_scripts.py`) failed in `test_legacy_refs_runs_clean_on_repo` (exit 1) and `test_gofmt_check_clean_repo` (exit 1).

## Part 2: Root Cause
1. **Legacy Refs Pattern Collision:** `check-legacy-refs.py` uses pattern `gitmap-v[567]\b` to detect old repository slugs. In `cli/cmdupdate/update_binary_swap_test.go`, mock release assets used the name `gitmap-v6.593.1-windows-amd64.zip`. Because `.` is a non-word character, `\b` matched `gitmap-v6`, causing the linter to flag false positive legacy references.
2. **Gofmt Trailing Newline:** `cli/cmdpipeline/pipeline_flags.go` had a trailing unformatted newline at the end of the file, causing `go-format-check.py --check-only` to fail.
3. **Nested If Guideline Drift:** `cli/cmdpipeline/pipeline_logs_sections.go:53` and `cli/cmdpipeline/pipeline_query_logs.go:138` contained nested `if` statements, violating the Zero Nested Ifs rule in `check-enum-and-boolean.py`.

## Part 3: Surgical Fix
1. In `cli/cmdupdate/update_binary_swap_test.go`: Marked lines with `// gitmap-legacy-ref-allow` to exempt test release archive fixtures from legacy slug checking.
2. In `cli/cmdpipeline/pipeline_flags.go`: Executed `gofmt -w` to remove trailing blank newline.
3. In `cli/cmdpipeline/pipeline_logs_sections.go`: Extracted `resolveSectionUrl(url string, runId int64) string` helper function to eliminate nested conditional URL assembly.
4. In `cli/cmdpipeline/pipeline_query_logs.go`: Flattened `resolveJobLogConcurrency(total int) int` using early returns and guard clauses (`total < 1`, `total > 8`).

## Part 4: Verification
All targeted linters and test suites pass locally with exit code 0:
- `python .github/scripts/go-format-check.py --check-only`: 100% PASS across 5,046 Go files.
- `python .github/scripts/check-legacy-refs.py .`: 100% PASS across repository.
- `python linter-scripts/check-relative-paths.py`: 100% PASS across 9,959 repository files.
- `python linter-scripts/check-enum-and-boolean.py --changed-only`: 100% PASS across 71 changed files.
- `python linter-scripts/check-forbidden-strings.py`: 100% PASS across all rules.
- `python .github/scripts/tests/test_ci_scripts.py`: 100% PASS across all 18 unit tests.
