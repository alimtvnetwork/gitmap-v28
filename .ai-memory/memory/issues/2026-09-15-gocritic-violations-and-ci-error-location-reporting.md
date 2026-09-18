# 2026-09-15: Gocritic Violations and CI/CD Error Reporting Enhancement

- **Category**: CI/CD / Static Analysis / Error Reporting
- **Status**: Resolved
- **Issue Reference**: `.ai-memory/cicd-issues/45-gocritic-violations-and-ci-error-location-reporting.md`
- **Pipeline Run**: #34937936831

## Context
During CI run #34937936831, the `gocritic` baseline diff gate failed with 4 new violations:
1. `cli/cmd/join.go:34`: `appendAssign: append result not assigned to the same slice`
2. `cli/cmd/mkdir.go:113`: `ifElseChain: rewrite if-else to switch statement`
3. `cli/cmdpurge/purge.go:16`: `ifElseChain: rewrite if-else to switch statement`
4. `cli/cmdschedule/schedule_debug.go:174`: `ifElseChain: rewrite if-else to switch statement`

In addition, the CI logs did not output the file paths and line numbers because GitHub Actions stripped the `::error file=...` annotations, causing severe friction in root-cause diagnosis.

## Remediation
1. **Gocritic Fixes**:
   - Refactored `allPositional := append(positional, ...)` to assign to `positional`.
   - Refactored 3 chained `if-else` blocks in `mkdir.go`, `purge.go`, and `schedule_debug.go` to `switch` statements.
2. **ErrorWrapper Fixes**:
   - Changed all method receivers in `cli/result/error_wrapper.go` from pointer `(ew *ErrorWrapper)` to value `(ew ErrorWrapper)` to eliminate `cannot call pointer method AsError on result.ErrorWrapper`.
3. **CI/CD Error Reporting Enhancement**:
   - `.github/scripts/check-single-linter-diff.py`: Emits structured finding boxes with exact file path (`file:line:col`), location, rule, script name, and local reproduction command, alongside explicit console lines (`❌ file:line:col: [linter] message`).
   - `.github/scripts/lint-diff.py`: Added explicit console prints for new findings.
   - `cli/cmdpipeline/pipeline_error_extract.go`: Added `cleanAnnotationError` to decode workflow annotations into `file:line:col: message`, enhanced `isStrongerSummary` to prioritize file location errors, and updated `formatSectionMetadata` to display `Script:` and `File:`.
