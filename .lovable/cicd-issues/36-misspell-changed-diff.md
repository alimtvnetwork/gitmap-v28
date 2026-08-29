# CI/CD Issue 36: Misspell Changed Files Gate Failure & US English Standardization

- **Job**: Spell Check (`misspell-changed.sh` / `misspell-changed.py`)
- **Type**: FAIL
- **Detected**: 2026-08-29
- **Status**: resolved

## Error
```text
.github/scripts/lint-issue-summary.py:14:0: "Behaviour" is a misspelling of "Behavior"
.github/scripts/lint-suggest.py:148:23: "colour" is a misspelling of "color"
.github/scripts/pr-summary.sh:10:53: "cancelled" is a misspelling of "canceled"
gitmap/clonepick/picker.go:29:56: "cancelled" is a misspelling of "canceled"
...
FAIL: misspellings found in repository files.
```

## Root Cause
Several documentation, script, and codebase files contained British English spellings (e.g. `behaviour`, `cancelled`, `colour`, `flavoured`, `optimise`, `organise`, `catalogue`, `labelled`), which violated the project's strict US-English spelling requirement enforced by `misspell -locale US`. In addition, `pickerModel` in `clonepick` had inconsistent field naming (`isCanceled`).

## Fix Applied
1. Bounded the spelling fixer strictly to the repository root using `git ls-files` without any external directory traversal.
2. Ran `misspell -w -locale US` across all tracked files in the repository.
3. Standardized `pickerModel.isCanceled` in `gitmap/clonepick/picker.go`, `picker_run.go`, and `picker_test.go`.
4. Removed redundant `//nolint:misspell` directive in `picker.go`.
5. Wired `Spell Check (misspell)` directly into `.lovable/ai-fix-scripts/03-cicd-local-runner.py`.
6. Verified with `python .lovable/ai-fix-scripts/03-cicd-local-runner.py` (all 20 checks passed with exit 0).
