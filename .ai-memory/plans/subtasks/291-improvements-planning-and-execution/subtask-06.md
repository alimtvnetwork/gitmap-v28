# Subtask 06: Database Dispatch Streamlining & Vacuum Error Propagation

> **Subtask ID:** Subtask-06  
> **Parent Plan:** `.ai-memory/plans/291-improvements-planning-and-execution.md`  
> **Target Path:** `.ai-memory/plans/subtasks/291-improvements-planning-and-execution/subtask-06.md`  
> **Status:** PENDING  

---

## 1. Objectives & Scope

1. Eliminate repetitive, verbose error-checking boilerplate across the command dispatch switch in `cli/cmddb/cmddb_dispatch.go`.
2. Author `cli/cmddb/db_vacuum.go` to house modular SQLite `VACUUM` and `PRAGMA optimize` routines with discrete error returns.
3. Stop silent swallowing of errors in `tryOptimizeDBFile` / `cmddb_optimize.go` where database locked or corruption errors are discarded as `return 0, 0`.
4. Enforce strict coding guidelines:
   - $\le 15$ lines per function.
   - Zero nested conditionals.
   - Blank line before every `return` and after `}`.
   - Positive boolean naming (`hasError`, `isOptimized`).

---

## 2. Detailed Technical Plan

### 2.1 Streamline Command Dispatcher in `cli/cmddb/cmddb_dispatch.go`
- Locate the subcommand switch cases (`doctor`, `listing`, `read`, `add`, `rm`, `edit`, `rest`):
  - Replace repetitive 5-line patterns:
    ```go
    // Legacy Anti-Pattern:
    appErr := RunDbDoctor(args[1:])
    if appErr != nil {
        return appErr
    }
    return nil

    // Compliant Refactor:
    return RunDbDoctor(args[1:])
    ```
  - Apply this guard simplification across all relevant switch arms in `cmddb_dispatch.go`.

### 2.2 Author Modular Vacuum Engine in `cli/cmddb/db_vacuum.go`
- Create `cli/cmddb/db_vacuum.go` with functions strictly $\le 15$ lines:
  - `executeVacuumOnDB(path string) (int64, error)`: executes `VACUUM` and `PRAGMA optimize` on a single SQLite database and computes reclaimed bytes.
  - `recordVacuumTelemetry(path string, reclaimed int64, label string)`: formats and outputs success telemetry line.
  - `handleVacuumError(path string, label string, err error) error`: outputs non-fatal warnings to `os.Stderr` while returning a structured error if critical.

### 2.3 Refactor Error Propagation in `runDBOptimize`
- In `cli/cmddb/cmddb_optimize.go`:
  - When `optimizeSingleDBFile(path)` fails:
    - Log diagnostic error indicating which database failed to optimize.
    - Return non-zero error count or aggregate error if master DB fails.
  - Ensure `runDBOptimize` reports partial failure summaries if individual databases fail to vacuum.

---

## 3. Verification Criteria

- [ ] Verbose `if appErr != nil { return appErr } return nil` patterns eliminated in `cli/cmddb/cmddb_dispatch.go`.
- [ ] `cli/cmddb/db_vacuum.go` created with functions $\le 15$ lines.
- [ ] Database vacuum errors are logged and bubbled rather than silently swallowed with `return 0, 0`.
- [ ] Blank line precedes every `return` and follows every `}` in `cmddb_dispatch.go` and `db_vacuum.go`.
- [ ] Zero nested `if` statements throughout database maintenance handlers.
