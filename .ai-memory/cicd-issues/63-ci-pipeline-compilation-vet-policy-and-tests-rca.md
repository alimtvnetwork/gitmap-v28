# RCA 63: CI Pipeline Compilation, Vet, Policy Linters, and Unit Tests

## 1. Symptom
GitHub Actions CI run `#35483242224` on commit `fe357c2` failed across 24 pipeline sections:
1. `Go vet`:
   - `cli/cmdssh/ssh_pass_cmd.go:138:20: not enough arguments in call to db.GetSSHConnectionByAlias`
   - `cli/cmdssh/ssh_pass_cmd.go:142:17: not enough arguments in call to db.GetSSHConnectionByIP`
   - `cli/cmdagy/agy_ping.go: undefined: result`
   - `cli/cmdagy/agy_ping.go: res.Error undefined (type Result has no field or method Error, but does have Err)`
2. `Go Compile Gate`:
   - `cli/store/ssh_repo.go: undefined: GetSSHHostByAlias, GetSSHHostByIP, ListSSHHosts`
3. `Nested If Linter`:
   - `cli/cmdssh/ssh_pass_cmd.go:175: nested if statement at depth 2 exceeds maximum allowed depth of 1`
   - `cli/cmdagy/agy_fix_pipeline_inject.go:227: nested if statement at depth 2`
   - `cli/cmdagy/agy_conv_scanner.go:215, 283: nested if statement at depth 2`
4. `Boolean Guidelines Linter`:
   - Inverted boolean expressions and negative identifiers (`hasNoPass`, `hasNoTranscript`, `hasNoTag`) in `cli/cmdssh/ssh_pass_cmd.go`, `cli/cmdssh/sshjoin_enroll.go`, `cli/cmdagy/agy_fix_pipeline_inject.go`, and `cli/cmdagy/agy_transcript_parse.go`.
5. `Relative Path Check`:
   - Hardcoded local Windows filesystem paths found in `cli/helptext/agy.md`.
6. `Unused Code Guard & golangci-lint`:
   - Unused helper functions `calcWideAvailableWidth` and `calcCompactAvailableWidth` in `cli/cmdpull/pull_table_layout.go`.
7. `Unit Tests & Smoke Failures`:
   - `cli/cmdpull/pull_table_test.go`: `TestFormatRepoNameLeadTruncatePresentation` expected suffix mismatch.
   - `cli/cmdmacro/macro_add_interactive.go`: Typed-nil `*apperror.AppError` returned as Go `error` interface evaluated as non-nil.
   - `cli/cmdssh/ssh_exec_install_test.go`: Shell expectation mismatch between Windows and Unix.
   - `cli/cmdssh/ssh_login_cmd_test.go`: SQLite database connection closed by deferred cleanup across consecutive test cases.
   - `cli/db/sshconnection.go`: SQLite `COALESCE` query returning text string failing type scan into `time.Time`.

## 2. Root Cause
- Missing imports and incorrect field accesses (`result` import omitted; `res.Error` accessed instead of `res.Err`).
- Missing helper forwarding functions in `cli/store/ssh_repo.go`.
- Conditionals exceeding maximum allowed nesting depth (depth > 1) violating project coding guidelines.
- Negative boolean names (`hasNoXxx`) violating boolean naming standards.
- Absolute Windows path examples copied directly into markdown documentation.
- Go typed-nil pitfall: concrete pointer `*apperror.AppError(nil)` returned to interface `error` evaluates to `err != nil`.
- SQLite `COALESCE` expressions producing text strings in `modernc.org/sqlite`, causing scan failure when scanning into `*time.Time`.
- Reused singleton `*store.DB` instance closed prematurely by `defer db.Close()` inside individual command runs.

## 3. Resolution
1. **Compilation & Vet:**
   - Added missing `"github.com/alimtvnetwork/gitmap-v28/cli/result"` import and fixed `res.Err` in `cli/cmdagy/agy_ping.go`.
   - Added `GetSSHHostByAlias`, `GetSSHHostByIP`, and `ListSSHHosts` alias functions to `cli/store/ssh_repo.go`.
   - Updated `db.GetSSHConnectionByAlias/IP` unpack signatures in `cli/cmdssh/ssh_pass_cmd.go`.
2. **Policy & Guidelines:**
   - Flattened nested conditionals in `cli/cmdssh/ssh_pass_cmd.go`, `cli/cmdagy/agy_fix_pipeline_inject.go`, and `cli/cmdagy/agy_conv_scanner.go` using extracted single-responsibility functions.
   - Replaced negative booleans (`hasNoPass`, `hasNoTranscript`, `hasNoTag`) with positive identifiers (`hasPass`, `hasTranscript`, `hasTag`).
   - Sanitized hardcoded Windows absolute paths in `cli/helptext/agy.md` with portable relative placeholders.
   - Removed dead/unused layout calculation functions in `cli/cmdpull/pull_table_layout.go`.
3. **Unit Tests & Bug Fixes:**
   - Changed error return types in `cli/cmdmacro/macro_add_interactive.go` from `*apperror.AppError` to `error` interface to eliminate typed-nil bugs.
   - Updated `cli/db/sshconnection.go` queries and row scanning to query `FirstRunAt` directly and scan into `sql.NullTime`.
   - Updated `cli/cmdpull/pull_table_test.go` assertion for 24-char lead truncation.
   - Updated `cli/cmdssh/ssh_exec_install_test.go` and `cli/cmdssh/ssh_login_cmd_test.go` test fixtures and mock database lifecycle.
   - Formatted all modified Go files with `gofmt -w`.

## 4. Verification
- `python 03-ai-scripts/06-cicd-local-runner.py`: All 46 quality gates passed (exit code 0).
- All linters (`check-error-management.py`, `check-enum-and-boolean.py`, `check-boolean-guidelines.py`, `check-nested-ifs.py`, `check-relative-paths.py`, `go-format-check.py`) passed cleanly with 0 violations.
- `go vet -C cli ./...` and `golangci-lint` passed with 0 issues.
- `go test ./cmdagy/... ./cmdssh/... ./cmdmacro/... ./cmdpull/...` passed with 100% success.
