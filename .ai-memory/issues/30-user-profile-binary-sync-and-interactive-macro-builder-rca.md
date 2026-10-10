# RCA-30: Stale User-Profile Binary Drift and Interactive Macro Builder Step Interception

## 1. Executive Summary

- **One-Sentence Root Cause:** Macro execution terminated with exit status 1 (`ItemNotFoundException`) on Windows PowerShell because user Alim's shell invoked a stale runtime binary (v6.447.0) that predated platform-adaptive safe removal shims, and interactive macro creation previously intercepted workspace modification commands as non-recorded helpers.
- **Status:** Resolved
- **Affected Surface:** Multi-user Windows CLI paths (`C:\Users\Alim\AppData\Local\gitmap\`, `AppData\Local\gitmap-cli\`), `cli/cmdmacro/macro_add_interactive.go`, `cli/cmdmacro/macro_edit.go`, and inverted boolean assertion patterns in tests.

---

## 2. Symptom & Evidence

Executing user macro `gitmap alim1` in user `Alim`'s shell:
```text
PS C:\Users\Alim> gitmap alim1
  alim1
  └── rm test

  ▶ Executing Macro: "alim1" (1 steps)

  [ 1/1] ➜ rm test

  rm : Cannot find path 'C:\Users\Alim\test' because it does not exist.
  At line:1 char:1
  + rm test
  + ~~~~~~~
      + CategoryInfo          : ObjectNotFound: (C:\Users\Alim\test:String) [Remove-Item], ItemNotFoundException
      + FullyQualifiedErrorId : PathNotFound,Microsoft.PowerShell.Commands.RemoveItemCommand

  ✖ failed (0.2s)
  ✖ Step 1 failed: exit status 1
```

Followed by interactive macro editing:
```text
PS C:\Users\Alim> gitmap macro edit alim1
  ● Interactive Macro Editor: "alim1" (1 steps)
  Current steps:
    1. rm test
```
Where workspace modification steps like `mkdir -p test` were treated as non-recorded helper commands rather than macro steps unless prefixed with `+add`.

---

## 3. Detailed Root Cause Analysis

1. **Multi-User Windows Binary Drift:**
   On Windows developer environments with multiple user accounts (`Administrator` and `Alim`), local build outputs (`<repo-root>/bin/gitmap.exe`) or single-user deployments (`%LOCALAPPDATA%\gitmap-cli\gitmap.exe`) do not update user `Alim`'s active execution paths.
   User `Alim`'s shell environment executes binaries from:
   - `C:\Users\Alim\AppData\Local\gitmap\gitmap.exe`
   - `C:\Users\Alim\AppData\Local\gitmap-cli\gitmap.exe`
   These binaries remained on `v6.447.0` (compiled October 1, 2026), which predated the introduction of `AdaptCommandForPlatform` in `cli/macro/safe_rm.go` and `cli/macro/execute.go`. As a result, the runtime binary lacked the PowerShell idempotent removal loop guards (`Test-Path`), causing `rm test` to abort with code 1.

2. **In-Builder Helper Interception:**
   In `cli/cmdmacro/macro_add_interactive.go` and `cli/cmdmacro/macro_edit.go`, helper commands were matched too broadly: commands starting with `mkdir`, `cat`, etc., were intercepted by `processInBuilderCommand` without recording them into the macro definition, leaving macros with only destructive removal steps rather than setup + teardown workflows.

3. **Inverted Boolean Checks (`!isSuccess`):**
   Assertions across `cli/cmdos`, `cli/cmdschedule`, and `cli/store` relied on negative boolean logic (`!r.IsSuccess`), violating the repository coding guidelines which mandate affirmative failure checks via explicit `.IsFail()` / `.IsFailed()` methods.

---

## 4. Remediation Implemented

1. **Four-Target Binary Synchronization:**
   Recompiled `gitmap.exe` at v6.527.0 and deployed synchronously across all active Windows targets:
   - `<repo-root>/gitmap.exe` & `<repo-root>/bin/gitmap.exe`
   - `C:\Users\Administrator\AppData\Local\gitmap-cli\gitmap.exe`
   - `C:\Users\Alim\AppData\Local\gitmap-cli\gitmap.exe`
   - `C:\Users\Alim\AppData\Local\gitmap\gitmap.exe`
   Verified that `& "C:\Users\Alim\AppData\Local\gitmap\gitmap.exe" --version` reports `gitmap v6.527.0`.

2. **Interactive Macro Builder Polish:**
   - In `cli/cmdmacro/macro_edit.go`, refined `isExplicitHelperCmd` to properly discern helper directives (`add `, `del `, `replace `, `insert `, `list `, `test `, `help`, `save`, `exit`, `:`-prefixed).
   - In `cli/cmdmacro/macro_add_interactive.go`, guarded `processInBuilderCommand` in `processInteractiveStepLine` so user shell commands execute live AND record as macro steps.

3. **Boolean Modernization & Guard Inversion Elimination:**
   - Added `IsFail()` and `IsFailed()` to `OSTUIActionResult` (`cli/cmdos/os_tui_types.go`).
   - Added `IsFail()` to `MatchResult` and `MatchAllResult` (`cli/lazyregex/match_result.go`, `cli/lazyregex/match_all_result.go`).
   - Converted `!runs[0].IsSuccess` to `runs[0].IsFail()` in `cli/cmdschedule/schedule_split_test.go`.
   - Converted `!logs[0].IsSuccess` to `logs[0].IsFail()` in `cli/store/installation_split_db_test.go`.
   - Converted `!m.Results[0].IsSuccess` to `m.Results[0].IsFail()` in `cli/cmdos/os_tui_test.go`.

4. **Test Suite Verification:**
   - `go test ./cmdmacro ./cmdos ./cmdschedule ./store ./cmdinstall ./lazyregex` all passed cleanly with 0 failures.

---

## 5. Permanent Avoidances & Invariants

- **NEVER** deploy a binary fix exclusively to `Administrator` without updating all user profiles (`C:\Users\Alim\AppData\Local\gitmap\gitmap.exe` and `C:\Users\Alim\AppData\Local\gitmap-cli\gitmap.exe`).
- **NEVER** write inverted boolean assertions (`!isSuccess`, `!r.IsSuccess`). Provide and invoke explicit `.IsFail()` / `.IsFailed()` methods.
- **NEVER** allow interactive macro editors to silently discard user-entered shell commands.
