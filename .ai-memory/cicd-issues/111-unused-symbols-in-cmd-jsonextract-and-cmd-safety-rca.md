# RCA-111: Unused Symbols in cmd_jsonextract and cmd_safety

**Date:** 2026-10-11
**Status:** ✅ Resolved
**Severity:** High (CI Pipeline Failure in Lint Baseline Diff)
**Affected Workflows:** `CI (#38119396009)`
**Run URLs:**
- https://github.com/alimtvnetwork/gitmap-v28/actions/runs/38119396009

---

## 1. Symptom & Error Signature

Following commit `da35126`, GitHub Actions CI workflow `#38119396009` failed on job `Lint Baseline Diff` (Step: `Diff vs baseline (fail only on NEW findings)`):

```text
[unused] func `extractJSONString` is unused (NEW in /tmp/lint-current/report.json)
    ❌ cmd/cmd_jsonextract.go:4: [unused] func `extractJSONString` is unused
[unused] func `findKeyValue` is unused (NEW in /tmp/lint-current/report.json)
    ❌ cmd/cmd_jsonextract.go:16: [unused] func `findKeyValue` is unused
[unused] func `extractQuotedValue` is unused (NEW in /tmp/lint-current/report.json)
    ❌ cmd/cmd_jsonextract.go:31: [unused] func `extractQuotedValue` is unused
[unused] func `indexOf` is unused (NEW in /tmp/lint-current/report.json)
    ❌ cmd/cmd_jsonextract.go:45: [unused] func `indexOf` is unused
[unused] func `runGuard` is unused (NEW in /tmp/lint-current/report.json)
    ❌ cmd/cmd_safety.go:140: [unused] func `runGuard` is unused
[unused] const `guardHookBody` is unused (NEW in /tmp/lint-current/report.json)
    ❌ cmd/cmd_safety.go:163: [unused] const `guardHookBody` is unused
Process completed with exit code 1.
```

---

## 2. Root Cause

1. **Dead Helpers in `cmd/cmd_jsonextract.go`:**
   - Functions `extractJSONString`, `findKeyValue`, `extractQuotedValue`, and `indexOf` were declared in `cli/cmd/cmd_jsonextract.go` for package `cmd`.
   - Subpackages (`cmddoctor`, `cmdinstall`, `cmdupdate`) have their own local implementations of `extractJSONString`, leaving the copy in `cli/cmd` completely uncalled and unexported.
   - The entire file `cli/cmd/cmd_jsonextract.go` became dead code.

2. **Unused Guard Hook Installer in `cmd/cmd_safety.go`:**
   - In `cli/cmd/cmd_safety.go`, `runGuard` and `guardHookBody` were defined at lines 140–180.
   - The CLI guard command was migrated to `cmdautomation/guard_cmd.go` (`runGuardCmd`), leaving `runGuard` and `guardHookBody` in `cmd/cmd_safety.go` orphaned and never invoked.

---

## 3. Surgical Fix

1. **`cli/cmd/cmd_jsonextract.go`:**
   - Deleted `cli/cmd/cmd_jsonextract.go` completely as all four functions (`extractJSONString`, `findKeyValue`, `extractQuotedValue`, `indexOf`) were unused in package `cmd`.

2. **`cli/cmd/cmd_safety.go`:**
   - Removed unused `runGuard` function and `guardHookBody` constant.
   - Preserved active working-tree snapshot and rollback handlers (`RunSnapshot`, `writeSnapshot`, `RunRollback`, `latestSnapshot`).

---

## 4. Verification

1. `golangci-lint run --timeout=5m --issues-exit-code=0 --out-format=json` in `cli/` confirmed 0 new unused findings in `cli/cmd/`.
2. `python 03-ai-scripts/52-file-size-check.py --files cli/cmd/cmd_safety.go --limit 300` -> passed (138 lines <= 300 cap).
3. `python 03-ai-scripts/05-guideline-autofixer.py cli/cmd/cmd_safety.go --check-only` -> 0 newline or implicit boolean violations.
4. `python linter-scripts/check-nested-ifs.py --changed-only` -> 0 nested if violations.
5. `python linter-scripts/check-relative-paths.py` -> 0 absolute path violations across 8612 files.
6. `go -C cli test ./cmd/... -run TestCmd -v` -> all package tests pass.
