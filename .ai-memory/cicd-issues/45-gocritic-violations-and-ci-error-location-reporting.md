# CI/CD Issue 45: Gocritic Violations and CI/CD Error Reporting Enhancement (Missing File Paths & Script Context)

- **Job**: Gocritic diff (baseline-diff, full-path only)
- **Type**: FAIL
- **Detected**: 2026-09-15
- **Status**: resolved
- **Pipeline Run**: #34937936831

## 1. Why It Happened (High-Level Architectural Reason)
1. In CI run #34937936831, step `Gocritic diff (baseline-diff, full-path only)` failed due to 4 new `gocritic` static analyzer findings introduced during cluster, join, and schedule enhancements:
   - `appendAssign` in `cli/cmd/join.go:34`
   - `ifElseChain` in `cli/cmd/mkdir.go:113`
   - `ifElseChain` in `cli/cmdpurge/purge.go:16`
   - `ifElseChain` in `cli/cmdschedule/schedule_debug.go:174`
2. **Defective CI/CD Error Reporting**: When GitHub Actions runner executes `.github/scripts/check-single-linter-diff.py`, new violations were emitted *only* via workflow annotations: `::error file={filename},line={line},col={col}::[{label}] {text}`. GitHub Actions runner intercepts this annotation format for the PR UI and strips the file metadata (`file=...`) from plain console stdout logs. As a result, the console log only displayed:
   `[gocritic] ifElseChain: rewrite if-else to switch statement (NEW vs baseline)`
   This made it impossible for developers and AI assistants to locate the offending file without downloading raw JSON artifacts.
3. **Pipeline Error Log Parser Gaps**: `cli/cmdpipeline/pipeline_error_extract.go` stripped `##[error]` prefixes but did not parse legacy raw annotations `::error file=...,line=...::` into compiler-standard `filename:line:col: message`. Furthermore, `isStrongerSummary` did not prioritize file location lines over generic exit status strings.

## 2. How It Happened (Exact Execution Flow)
1. **Gocritic Failure**:
   - `cli/cmd/join.go`: `allPositional := append(positional, fs.Args()...)` triggered `appendAssign: append result not assigned to the same slice`.
   - `cli/cmd/mkdir.go:113`, `cli/cmdpurge/purge.go:16`, `cli/cmdschedule/schedule_debug.go:174`: Multiple sequential `if ... else if ... else if` branches checking argument flags triggered `ifElseChain: rewrite if-else to switch statement`.
2. **Path-Less CI Logs**:
   - `check-single-linter-diff.py` computed `new_keys = sorted(current_keys - baseline_keys)`.
   - In the failure branch, it executed only `print(f"::error file={filename},line={line},col={col}::[{label}] {text} (NEW vs baseline)")`.
   - The GitHub Actions runner stripped `::error file=...::` from the terminal stream, leaving only the rule name and message.
3. **Compiler Error on Result.ErrorWrapper**:
   - Methods on `result.ErrorWrapper` had pointer receivers `(ew *ErrorWrapper)`. When functions returned `result.ErrorWrapper` by value as unaddressable temporary expressions (e.g. `routeCluster(args).AsError()`), the Go compiler failed with `cannot call pointer method AsError on result.ErrorWrapper`.

## 3. Root Cause
- **Gocritic Code Violations**:
  - `cli/cmd/join.go`: Assigning `append` result to a new slice identifier instead of the original slice.
  - `cli/cmd/mkdir.go`, `cli/cmdpurge/purge.go`, `cli/cmdschedule/schedule_debug.go`: Using 3+ chained `if-else` blocks for string comparison instead of `switch`.
- **Linter Diff Script Reporting Gap**:
  - `.github/scripts/check-single-linter-diff.py` and `.github/scripts/lint-diff.py` did not print explicit stdout lines containing `{filename}:{line}:{col}`.
- **Value vs. Pointer Receiver Mismatch**:
  - `cli/result/error_wrapper.go`: Methods were defined on pointer receiver `*ErrorWrapper` instead of value receiver `ErrorWrapper`.
- **Pipeline Parser Missing Annotation Support**:
  - `cli/cmdpipeline/pipeline_error_extract.go` lacked an annotation parser to extract `file`, `line`, and `col` attributes into readable location strings.

## 4. Code Fix
1. **Fixed All 4 Gocritic Violations**:
   - `cli/cmd/join.go`: Reassigned to `positional = append(positional, fs.Args()...)`.
   - `cli/cmd/mkdir.go`: Refactored to `switch arg`.
   - `cli/cmdpurge/purge.go`: Refactored to `switch`.
   - `cli/cmdschedule/schedule_debug.go`: Refactored to `switch`.
2. **Fixed ErrorWrapper Receivers**:
   - `cli/result/error_wrapper.go`: Changed all method receivers from `(ew *ErrorWrapper)` to `(ew ErrorWrapper)`.
   - `cli/cmdssh/ssh_e2e_test.go`: Modernized assertions to check `res.IsSuccess()`.
   - `cli/cmd/profiles_cmd.go`: Added missing `cli/result` package import.
3. **Enhanced CI/CD Error Reporting in Scripts**:
   - In `.github/scripts/check-single-linter-diff.py`:
     - Emits formatted finding boxes with `File: {filename}:{line}:{col}`, `Location: ...`, `Linter: {label}`, `Message: {text}`, `Script: ...`, and copy-pasteable local reproduction command (`golangci-lint run --no-config --disable-all --enable={linter} ./{lint_dir}/...`).
     - Emits explicit terminal line `❌ {filename}:{line}:{col}: [{label}] {text}`.
   - In `.github/scripts/lint-diff.py`:
     - Added `print(f"  ❌ {file}:{line}: [{linter}] {message}")`.
4. **Enhanced Pipeline Error Extractor**:
   - `cli/cmdpipeline/pipeline_error_extract.go`:
     - Implemented `cleanAnnotationError` to decode `::error file={file},line={line},col={col}::{msg}` into standard compiler format `{file}:{line}:{col}: {msg}`.
     - Implemented `isLocationSummary` to prioritize file location lines over generic `exit status` summaries.
     - Enhanced `formatSectionMetadata` to extract and display `Script:` and `File:` directly in section failure banners.
