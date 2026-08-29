# 28-legacy-refs-bare-err-and-lint-keyerror

## Error Summary
The CI/CD pipeline encountered three failures across parallel jobs:
1. `Legacy Refs Check`: `check-legacy-refs.sh` identified 4 historical mentions of `gitmap-v6` in `.lovable/cicd-issues/20-cicd-four-headed-hydra.md` and `.lovable/cicd-issues/index.md`. <!-- gitmap-legacy-ref-allow -->
2. `Lint Baseline Guard`: `check-bare-stderr-err.sh` identified 8 bare error prints (`fmt.Fprintln(os.Stderr, err)`) across 6 commands in `gitmap/cmd/`.
3. `Lint Baseline Diff`: `lint-suggest.py` and `lint-diff.py` raised `KeyError: 'is_failure'` because the `query_wrapper` return dictionary uses key `is_fail`.

## Root Cause Analysis
1. `.lovable/cicd-issues/20-cicd-four-headed-hydra.md` documented a past bug fix where `gitmap-v6` references were replaced. Because `check-legacy-refs.sh` runs grep across all files without `gitmap-legacy-ref-allow` markers, mentioning historical version names in CI logs tripped the check.
2. Older CLI commands in `gitmap/cmd` (`dedupe.go`, `orphans.go`, `size.go`, `stale.go`, `visibilityhistory.go`, `visibilityundo.go`) used bare `fmt.Fprintln(os.Stderr, err)` instead of the project's structured `cliexit.Fail` helper.
3. `lint-suggest.py` and `lint-diff.py` accessed `res["is_failure"]` directly. `query_wrapper` in `.github/scripts/query_wrapper.py` returns `is_fail` and `is_success`, not `is_failure`.

## Solution Applied
1. Appended `<!-- gitmap-legacy-ref-allow -->` to historical lines in `.lovable/cicd-issues/20-cicd-four-headed-hydra.md` and `.lovable/cicd-issues/index.md`.
2. Replaced all 8 bare error print occurrences with `cliexit.Fail(command, op, subject, err, exitCode)`.
3. Updated `lint-suggest.py` and `lint-diff.py` to check `res.get("is_fail") or res.get("is_failure")`.
4. Generated `.lovable/ai-fix-scripts/03-cicd-local-runner.py` and verified all checks pass with exit code 0.

## What NOT to Repeat
- NEVER write bare `fmt.Fprintln(os.Stderr, err)` in `gitmap/cmd/`. Always use `cliexit.Fail` or `cliexit.Reportf`.
- NEVER access dictionary result flags with direct index `res["is_failure"]`. Always use `.get()` matching the canonical `is_fail` / `is_success` contract.
- When referencing older sibling version strings in markdown memory or issues, always annotate with `<!-- gitmap-legacy-ref-allow -->`.
