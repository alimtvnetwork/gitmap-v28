# RCA-110: Unused Symbols in AliasOps, CLI Helpers Argv, LLM Train, and CmdNodes Hook Shim

**Date:** 2026-10-11
**Status:** ✅ Resolved
**Severity:** High (CI Pipeline Failure in Lint Baseline Diff)
**Affected Workflows:** `CI (#38113040806)`
**Run URLs:**
- https://github.com/alimtvnetwork/gitmap-v28/actions/runs/38113040806

---

## 1. Symptom & Error Signature

Following commit `c8c6d90`, GitHub Actions CI workflow `#38113040806` failed on job `Lint Baseline Diff` (Step: `Diff vs baseline (fail only on NEW findings)`):

```text
[unused] func `isLegacyDataError` is unused (NEW in /tmp/lint-current/report.json)
    ❌ cmd/aliasops.go:186: [unused] func `isLegacyDataError` is unused
[unused] func `isTerminalInput` is unused (NEW in /tmp/lint-current/report.json)
    ❌ cmd/clihelpers_argv.go:15: [unused] func `isTerminalInput` is unused
[unused] func `isExistingFile` is unused (NEW in /tmp/lint-current/report.json)
    ❌ cmd/clihelpers_argv.go:24: [unused] func `isExistingFile` is unused
[unused] func `isFileFlagWithArg` is unused (NEW in /tmp/lint-current/report.json)
    ❌ cmd/clihelpers_argv.go:29: [unused] func `isFileFlagWithArg` is unused
[unused] func `printCurriculumSummary` is unused (NEW in /tmp/lint-current/report.json)
    ❌ cmd/llm/llm_train.go:90: [unused] func `printCurriculumSummary` is unused
[unused] func `checkHelp` is unused (NEW in /tmp/lint-current/report.json)
    ❌ cmdnodes/hook_shim.go:6: [unused] func `checkHelp` is unused
Process completed with exit code 1.
```

---

## 2. Root Cause

1. **Dead Legacy Error Checker in `cmd/aliasops.go`:**
   - `isLegacyDataError(err error) bool` was defined at line 186 but never referenced anywhere in `cmd/aliasops.go` or package `cmd`.
   - The unused helper also held an unnecessary import of package `strings`.

2. **Unused Argv Predicates in `cmd/clihelpers_argv.go`:**
   - During the `cli/cmd` modularization wave (moving command handlers into subpackages like `cmdclone`, `cmdcopy`, `cmdinstall`, `cmdmacro`), subpackages implemented their own local helpers.
   - `isTerminalInput`, `isExistingFile`, `isFileFlagWithArg`, `hasArgFlag`, `extractFlagVal`, and `printJSON` remained in `cli/cmd/clihelpers_argv.go` as unexported dead functions never called within package `cmd`.
   - The unused helpers also held unnecessary imports of `encoding/json`, `fmt`, `os`, and `strings`.

3. **Orphaned Wrapper in `cmd/llm/llm_train.go`:**
   - `printCurriculumSummary()` was defined at line 90 as a wrapper around `printChainedSequence()` and `printOperationalGuidelines()`, but calls in `llm_phases.go` invoke those functions directly, leaving `printCurriculumSummary()` orphaned and unreferenced.

4. **Unused DI Shim Function in `cmdnodes/hook_shim.go`:**
   - `checkHelp(command string, args []string)` was defined at line 6, but `cmdnodes` uses direct help rendering and never invoked `checkHelp`. Package `cmdcode` had previously removed this helper, but `cmdnodes` still had it.

---

## 3. Surgical Fix

1. **`cli/cmd/aliasops.go`:**
   - Removed unused `isLegacyDataError` function and pruned unused `"strings"` import.

2. **`cli/cmd/clihelpers_argv.go`:**
   - Removed unused `isTerminalInput`, `isExistingFile`, `isFileFlagWithArg`, `hasArgFlag`, `extractFlagVal`, and `printJSON`.
   - Pruned unused imports (`"encoding/json"`, `"fmt"`, `"os"`, `"strings"`).
   - Retained only the used `isHelpFlag` helper with proper vertical spacing.

3. **`cli/cmd/llm/llm_train.go`:**
   - Removed unused `printCurriculumSummary` function.

4. **`cli/cmdnodes/hook_shim.go`:**
   - Removed unused `checkHelp` function while preserving `CheckHelpFn` variable declaration for DI wiring.

---

## 4. Verification

1. `golangci-lint run --config .golangci.yml --issues-exit-code=0 --out-format=json` in `cli/` verified: 0 new or existing unused findings across all 4 modified files.
2. `python 03-ai-scripts/52-file-size-check.py --files cli/cmd/aliasops.go cli/cmd/clihelpers_argv.go cli/cmd/llm/llm_train.go cli/cmdnodes/hook_shim.go --limit 300` -> all 4 files passed under the 300-line cap.
3. `python 03-ai-scripts/05-guideline-autofixer.py <file> --check-only` -> 0 newline or boolean convention violations.
4. `python linter-scripts/check-nested-ifs.py --changed-only` -> 0 nested if violations.
5. `python linter-scripts/check-relative-paths.py` -> 0 absolute path violations.
