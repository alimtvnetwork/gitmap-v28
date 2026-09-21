# CI/CD Issue 72: SSH Build Failures and gitmap pe Stack Trace Leak RCA

## 1. Why It Happened
CI run #35598235164 on commit `4e53b9e` failed across multiple quality gates:
1. **Lint [Go vet]:** Unused `"strings"` imports in `cli/cmdssh/sshjoin_rm_cmd.go` and `cli/cmdssh/sshexisting.go`, and redeclared `var SJRmCmd` in `cli/cmdssh/sshjoin_rm_cmd.go` (previously declared in `cli/cmdssh/sshjoin_cmd.go:173`).
2. **Error Management Policy Check:** Swallowed database execution error `_, _ = conn.Exec(sqlCreateSSHHistory)` in `cli/cmdssh/ssh_history_db.go:33`.
3. **Test Failures / Undefined Symbols:** `cli/cmdssh/sshjoin_rm_cmd_test.go` referenced `resolveRmTarget` and `validateRmTarget`, which were missing after prior file size decomposition.
4. **gitmap pe Stack Trace Leak:** When `gh run view --log-failed` fails (e.g. for cancelled runs or missing log streams), `cli/cmdpipeline/pipeline_query.go` captured `apperror.CaptureStackTrace(...)`, which `pipeline_error_extract.go` matched as an error marker (`"Stack Trace:"`), displaying raw internal Go stack traces in the user terminal instead of a clean, helpful error summary.

---

## 2. How It Happened
- During file size and function size refactoring in `cli/cmdssh/`, helper functions and variables were moved without cleaning up unused package imports or duplicate declarations across split files.
- `resolveRmTarget` and `validateRmTarget` were removed from `sshjoin_rm_cmd.go` to keep the file under 100 lines, leaving unit tests in `sshjoin_rm_cmd_test.go` broken.
- `openSSHHistoryDB` in `ssh_history_db.go` used `_ = store.ConfigureSQLiteConn(conn)` and `_, _ = conn.Exec(...)` to suppress SQLite initialization errors, violating the repository's strict zero-swallowed-errors policy.
- `formatGHFailedError` in `pipeline_query.go` called `apperror.CaptureStackTrace(apperror.DefaultStackTraceSkip)` on `gh` subprocess failures, polluting error output with Go runtime frames.

---

## 3. Root Cause
- `cli/cmdssh/sshjoin_rm_cmd.go:6`: Unused `"strings"` import.
- `cli/cmdssh/sshexisting.go:5`: Unused `"strings"` import.
- `cli/cmdssh/sshjoin_cmd.go:173`: Duplicate declaration of `SJRmCmd`.
- `cli/cmdssh/ssh_rm_target.go`: Missing dedicated module for `resolveRmTarget` and `validateRmTarget`.
- `cli/cmdssh/ssh_history_db.go:33`: Swallowed error on `conn.Exec(sqlCreateSSHHistory)`.
- `cli/cmdpipeline/pipeline_query.go:246`: Unnecessary Go stack trace capture injected into user-facing pipeline error diagnostics.
- `cli/cmdpipeline/pipeline_query.go:13`: Unused `apperror` import after removing stack trace capture.
- `cli/cmdssh/ssh_help_check.go:39`: Unused `checkHelp` function triggering `golangci-lint` `unused` rule.

---

## 4. Code Fix
- Removed unused `"strings"` imports from `cli/cmdssh/sshjoin_rm_cmd.go` and `cli/cmdssh/sshexisting.go`.
- Removed duplicate `SJRmCmd` declaration from `cli/cmdssh/sshjoin_cmd.go`.
- Created `cli/cmdssh/ssh_rm_target.go` containing `resolveRmTarget` and `validateRmTarget`, strictly adhering to affirmative boolean and $\le 15$ line limits.
- Refactored `cli/cmdssh/ssh_history_db.go` to extract `initSSHHistorySchema(conn *sql.DB) error`, properly checking and wrapping all SQLite configuration and execution errors via `apperror.WrapSimple`.
- Removed `stack` capture and formatting from `formatGHFailedError` and `buildGHFailedErrorMessage` in `cli/cmdpipeline/pipeline_query.go`.
- Removed unused `apperror` import from `cli/cmdpipeline/pipeline_query.go`.
- Removed unused `checkHelp` function from `cli/cmdssh/ssh_help_check.go`.
