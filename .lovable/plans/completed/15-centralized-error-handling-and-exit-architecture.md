# 15. Centralized Error Handling and Exit Architecture

master-plan: 15-centralized-error-handling-and-exit-architecture
status: in-progress
created: 2026-08-29
author: Antigravity Orchestrator

---

## 1. Deep Analysis of the Visual Evidence (Images 1, 2, and 3)

The user submitted three differential screenshots from GitHub Desktop documenting recent modifications:
- **Image 1** (`.lovable/assets/error-handling/01-releasepull-panic-vs-exit.png`): Shows `gitmap/cmd/releasepull.go` replacing `panic("fatal error")` with `os.Exit(1)`.
- **Image 2** (`.lovable/assets/error-handling/02-reinstall-panic-vs-exit.png`): Shows `gitmap/cmd/reinstall.go` replacing two instances of `panic("fatal error")` with `os.Exit(1)`.
- **Image 3** (`.lovable/assets/error-handling/03-rootadd-panic-vs-exit.png`): Shows `gitmap/cmd/rootadd.go` replacing two instances of `panic("fatal error")` with `os.Exit(1)`.

### Why Both Approaches Are Anti-Patterns

#### A. The Legacy Pattern: `panic("fatal error")`

1. **Uncontrolled Crash**: Panics crash the Go runtime and dump internal thread backtraces, frame pointers, and goroutine states directly into terminal output.
2. **Missing Domain Context**: The string `"fatal error"` conveys zero semantic meaning. It does not identify the failing operation, the entity being processed, the reason for the failure, or what the user should do to resolve it.
3. **No Structured Diagnostic Envelope**: Panics bypass structured logging, telemetry, and error storage tables.

#### B. What AI Modified: Bare `os.Exit(1)`

1. **Silent / Destructive Abort**: Calling `os.Exit(1)` immediately terminates the process. Deferred functions (like closing files, flushing buffers, releasing locks, and writing error audit records) are **never executed**.
2. **Loss of Diagnostic Context**: A bare `os.Exit(1)` produces no structured record of who created the error, what type of error occurred, or what variables led to the state.
3. **Violation of Centralized Error Management**: Error policies (whether to log, exit, report, or notify) are hardcoded at dozens of arbitrary call sites instead of being controlled through a single, configurable error handler.

#### C. The Standard / Correct Architectural Pattern

Centralized error management requires two distinct, coordinated responsibilities:
1. **Structured Domain Error Packaging**: Functions construct or wrap errors using `apperror.AppError`, attaching:
   - `Op`: The operation name (e.g. `release.pull`, `cmd.reinstall`, `add.dispatch`).
   - `Code`: A registered error code (e.g. `E9001`, `E9002`).
   - `Type`: A typed category enum (`ErrorTypeValidation`, `ErrorTypePrecondition`, `ErrorTypeExecution`, etc.).
   - `Severity`: An impact level (`SeverityError`, `SeverityFatal`).
   - `Creator`: The component or package responsible for raising the error.
   - `Ctx`: A key-value map capturing all relevant variables, paths, and flags.
   - `Cause`: The underlying root cause error if wrapping.
2. **Centralized Error Dispatcher & Handler (`HandleError`)**:
   - Instead of calling `panic` or `os.Exit(1)` at arbitrary points, the CLI entry point or runner passes the `AppError` to a centralized handler.
   - The handler formats the message uniformly to `stderr` with actionable user guidance.
   - The handler flushes registered terminal pipes and buffers (`cliexit.Drain`).
   - The handler logs the error to the audit/session store.
   - The handler transitions cleanly to `os.Exit(code)` (or panics if debug assertion mode is active).

---

## 2. Task-Specific Rule Set (Non-Negotiable)

1. **RULE-ERR-01 (No Bare Exits or Panics)**: Never write bare `os.Exit(...)` or `panic(...)` in business logic, command runners, or helper functions. All error paths must route through `apperror` and `cliexit.HandleError`.
2. **RULE-ERR-02 (Always Wrap with Metadata)**: Every error must contain an operation name (`Op`), registered code (`Code`), error category (`Type`), and contextual parameters (`Ctx`). Never produce an error without identifying its creator.
3. **RULE-ERR-03 (Never Be Silent)**: Every error must be reported with meaningful diagnostics. Silent terminations or swallowed errors are strictly prohibited.
4. **RULE-ERR-04 (Clean Resource Teardown)**: Centralized exit transitions must flush all registered pipe drainers and close database handles before terminating.
5. **RULE-ERR-05 (No Automatic Releases)**: All commits must remain standard development commits. Do NOT bump versions, touch changelog, or create GitHub releases.

---

## 3. Master Architecture & Components

```
┌────────────────────────────────────────────────────────┐
│                   gitmap/cmd (CLI Callers)             │
│                                                        │
│   err := apperror.New("reinstall.confirm", "E9001",    │
│            map[string]any{"mode": mode},               │
│            apperror.ErrorTypeAbort,                    │
│            apperror.SeverityWarn)                      │
│                                                        │
│   cliexit.HandleError(err) ────────────────────────┐   │
└────────────────────────────────────────────────────┼───┘
                                                     │
                                                     ▼
┌────────────────────────────────────────────────────────┐
│               gitmap/cliexit.HandleError()             │
│                                                        │
│  1. Unwraps AppError or creates default envelope.      │
│  2. Formats uniform output:                            │
│     "gitmap <op> [<code/type>]: <msg> (creator: ...)"  │
│  3. Drains pipe buffers (theme, glyphs).               │
│  4. Records error in audit/log store.                  │
│  5. Dispatches exit according to configured strategy:  │
│     - CLI Strategy: os.Exit(code)                      │
│     - Debug Strategy: panic with diagnostic report     │
└────────────────────────────────────────────────────────┘
```

---

## 4. Subtasks Breakdown

- [ ] **Subtask 1 (`01-apperror-enhancement.md`)**:
  - Enhance `gitmap/apperror/apperror.go`:
    - Add `ErrorType` enum (`ErrorTypeValidation`, `ErrorTypeNotFound`, `ErrorTypeExecution`, `ErrorTypeAbort`, `ErrorTypeInternal`).
    - Add `SeverityType` enum (`SeverityInfo`, `SeverityWarn`, `SeverityError`, `SeverityFatal`).
    - Add `Creator` field to `AppError`.
    - Provide constructors: `NewWithDetails`, `WrapWithDetails`, `NewSimple`, `WrapSimple`.
- [ ] **Subtask 2 (`02-cliexit-central-handler.md`)**:
  - Enhance `gitmap/cliexit/cliexit.go`:
    - Implement `HandleError(err error, options ...HandleOption)`.
    - Format errors with full diagnostic context (code, op, creator, context variables).
    - Provide pluggable strategies: `ExitStrategyProcess` (default) vs. `ExitStrategyPanic` (debug mode).
    - Ensure all flushers are drained before exit.
- [ ] **Subtask 3 (`03-cmd-refactor.md`)**:
  - Refactor command entry points across `gitmap/cmd/*.go` (`root.go`, `releasepull.go`, `reinstall.go`, `rootadd.go`, `releaserebase.go`, `reset.go`, `revert.go`, `revertscript.go`, `reverttxn.go`, `reverttxn_lastn.go`, `scanresolve.go`, `selfuninstallhandoff.go`, `seowrite.go`, `seowritecsv.go`, `seowriteloop.go`, `seowritetemplate.go`, `sshcat.go`) to use `apperror` and `cliexit.HandleError`.
- [ ] **Subtask 4 (`04-generic-guidelines-and-prompt.md`)**:
  - Create File 1: Generic Error Handling Guideline Document at `.lovable/coding-guidelines/centralized-error-handling-architecture.md` and link in `spec/03-error-manage/`.
  - Create File 2: Reusable System Prompt at `.lovable/prompts/05-coding-guidelines/02-centralized-error-handling-checklist.md`.
- [ ] **Subtask 5 (`05-verification-and-ci.md`)**:
  - Write comprehensive unit tests for `apperror.AppError` and `cliexit.HandleError`.
  - Run `.lovable/ai-fix-scripts/03-cicd-local-runner.py` to ensure all 11 checks pass with exit code 0.

---

## 5. Acceptance Criteria

- [ ] Zero instances of `panic("fatal error")` or untyped `panic("...")` in production code.
- [ ] Zero instances of bare `os.Exit(...)` bypassing `cliexit.HandleError` or `cliexit.Exit`.
- [ ] `apperror.AppError` captures `Op`, `Code`, `Type`, `Severity`, `Creator`, and `Ctx`.
- [ ] Centralized error handler flushes buffers and renders uniform, user-actionable error messages.
- [ ] Generic error guideline file created and shareable with any project.
- [ ] Reusable AI prompt file created in `.lovable/prompts/` with full error checklist.
- [ ] Unit tests verify both exit and diagnostic formatting.
- [ ] All 11 local runner checks pass 100% green.
- [ ] No version bump or release is performed.
