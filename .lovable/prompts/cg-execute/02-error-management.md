# Instruction (must follow): Coding Guideline Execution — Error Management & Architecture

Trigger Keywords & Aliases: `cg-error`, `cg-execute error`, `fix error management`, `execute error architecture`

```text
N = 100
```

N = total self-loop steps budget. The user may override this number when triggering the prompt.

- [ ] /goal First `N/2` steps (Phase 1) are dedicated to auditing all error handling sites, finding bare panics, bare exits, swallowed errors, and un-wrapped errors, writing the master execution plan to `.lovable/plans/pending/XX-error-management-audit.md`, and decomposing it into subtasks in `.lovable/plans/subtasks/XX-error-management/`.
- [ ] /goal Second `N/2` steps (Phase 2) are dedicated to executing those subtasks in an autonomous self-loop until all error paths route through structured `AppError` wrappers and centralized dispatchers.
- [ ] /goal **Linter Mandate**: If an automated linter for bare panics, bare exits, and swallowed errors does not exist in CI/CD, you MUST create an advanced scanner script (e.g. in `scripts/` or `linters/`) and wire it directly into the CI/CD local runner and workflows.
- [ ] /learn Ingest `.lovable/coding-guidelines/centralized-error-handling-architecture.md`, `spec/03-error-manage/`, `.lovable/coding-guidelines/coding-guidelines.md`, and `.lovable/strictly-avoid.md`.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after the user sets them. Never change them mid-execution.

---

## 1. Error Management Non-Negotiable Checklist

You MUST audit and strictly enforce every rule below across the entire codebase:

### A. The Double Anti-Pattern Ban
1. **No Bare Panics**: Never call `panic("...")`, `panic(err)`, or `throw "string"` in business logic, command runners, or helpers.
2. **No Bare Exits**: Never call `os.Exit(...)` or `sys.exit(...)` directly from command handlers. All exits must route through the central error dispatcher.

### B. Structured Domain Errors (`AppError`)
1. **Always Wrap with Full Context**: Construct or wrap all errors using `apperror.AppError` (or language equivalent), ensuring the following fields are populated:
   - `Op`: Operation name (e.g. `cmd.reinstall`, `db.query`, `auth.verify`).
   - `Code`: Stable registered error code (e.g. `E1001`, `E2002`).
   - `Type`: Classified error category enum (`VALIDATION`, `NOT_FOUND`, `PRECONDITION`, `EXECUTION`, `ABORT`, `INTERNAL`).
   - `Severity`: Level enum (`INFO`, `WARN`, `ERROR`, `FATAL`).
   - `Creator`: Attributed component/package that created the error.
   - `Message`: Human-readable description with actionable remediation.
   - `Ctx`: Key-value map of relevant runtime parameters, paths, and flags.
   - `Cause`: The underlying wrapped error.
2. **Preserve Root Cause**: Always preserve the underlying error stack using `Unwrap() error`.

### C. Centralized Dispatch & Never-Be-Silent Rule
1. **Central Dispatcher (`HandleError`)**: Route error handling to the centralized exit/dispatch function (e.g. `cliexit.HandleError(err, code)`).
2. **Buffer & Pipe Flushing**: The central handler must flush all registered pipe drainers and buffer sinks before process exit.
3. **Pluggable Strategies**: Support clean CLI exit, debug mode panic (`APP_ERROR_PANIC=1`), and API envelope serialization.
4. **Never Swallow**: Every catch/error block must log full context and handle or rethrow. Silent error swallowing is a build failure.

---

## 2. Phase 1: Planning, Audit & Subtask Decomposition (Steps 1 .. N/2)

1. **Global Codebase Audit**: Search for bare `os.Exit`, untyped `panic(`, swallowed errors, and generic `errors.New` without domain metadata.
2. **Master Plan**: Write a detailed execution plan to `.lovable/plans/pending/XX-error-management-audit.md`.
3. **Subtask Files**: Decompose into subtask files in `.lovable/plans/subtasks/XX-error-management/` (e.g. `01-task.md`, `02-task.md`, ...).
4. **Linter Connection**: Implement an automated error linter (e.g. `python scripts/lint-errors.py`) that checks for unapproved `os.Exit` and `panic` calls, and hook it into `.lovable/ai-fix-scripts/03-cicd-local-runner.py`.

---

## 3. Phase 2: Autonomous Execution Loop (Steps N/2+1 .. N)

1. **Execute Subtasks**: Refactor error handling across all identified files to use `AppError` and the centralized dispatcher.
2. **Unit Testing**: Add tests for error constructors, error formatting, unwrapping, and centralized exit handlers.
3. **Verify CI**: Run the local CI/CD runner and verify that all error linters and tests pass cleanly with exit code 0.
4. **Update Status**: Mark completed tasks as `DONE`, move completed plans to `.lovable/plans/completed/`, and update `.lovable/plans/index.md`.

---

## 4. Pre-Commit Verification Checklist

- [ ] Zero bare `panic(...)` calls exist in production code.
- [ ] Zero bare `os.Exit(...)` calls exist outside the central error dispatcher.
- [ ] Every error path constructs/wraps an `AppError` with `Op`, `Code`, `Type`, `Severity`, `Creator`, and `Ctx`.
- [ ] Central error dispatcher flushes buffers before exiting.
- [ ] Error linters are integrated into CI/CD and pass with exit 0.
