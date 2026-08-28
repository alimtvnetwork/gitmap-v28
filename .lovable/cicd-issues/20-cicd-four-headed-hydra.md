# 20-cicd-four-headed-hydra

## Error Summary
The CI/CD pipeline failed with four separate errors:
1. `cmd/envplatform_unix.go:142:10`: Too many return values (`have *apperror.AppError`, `want ()`).
2. `gofmt -d`: Several `.go` files were not formatted cleanly.
3. `lint-issue-summary.py`: `KeyError: 'is_failure'` because the JSON response field was typed as `is_fail` by the query wrapper.
4. `policy-check`: Found legacy references to `gitmap-v6` in `spec/` markdown files.

## Root Cause Analysis
1. A function (`writeProfileContent`) in `cmd/envplatform_unix.go` was attempting to return an error (`*apperror.AppError`), but its signature declared no return type. The strict error management rule ("never swallow errors, always wrap with context") was applied inconsistently, causing a compilation error.
2. An earlier AI step ran `sed` or Python regex replaces on `.go` files, bypassing `gofmt -w .`, which left dangling trailing spaces and missing blank lines.
3. The query wrapper standard returns `is_fail` (boolean rule constraint), but the downstream lint summary Python script still expected `is_failure`.
4. Legacy markdown files hardcoded previous sibling repository names (e.g. `gitmap-v6`), violating the policy scanner for `gitmap-v[567]\b`.

## Solution Applied
1. Cascaded the `error` return type from `writeProfileContent` all the way up through `appendToProfile`, `setEnvPersistent`, and into `applyEnvSet` in both Windows and Unix implementations.
2. Ran `gofmt -w .` across the repository to perfectly standardize all `.go` files.
3. Updated `.github/scripts/lint-issue-summary.py` to check for `is_fail` instead of `is_failure`.
4. Replaced all occurrences of `gitmap-v6` in `spec/**/*.md` with `gitmap-v28`.

## What NOT to Repeat
- NEVER run text-replacements on `.go` files without subsequently running `gofmt -w .`.
- NEVER change a return type inside a deeply nested utility without cascading it all the way to the top-level invoker.
- NEVER rename API JSON boolean keys without verifying downstream Python/Bash scripts that consume them.
