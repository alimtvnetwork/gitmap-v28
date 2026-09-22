# Plan 81 (Completed): SSH Join Authentication Error Logging, Execution Tracing & Diagnostic Command

## 1. Goal Description
The user reported that `gitmap ssh join a@192.168.1.3 w1` failed authentication with a generic message:
`⚠ Failed to authenticate with remote machine: ssh: authentication failed: invalid password or remote server rejected credentials`
without displaying what actions were run, what auth methods were attempted, or what the internal error was, leaving the user unable to diagnose or share the failure details.

The goal was to:
1. Retain the underlying raw SSH internal error without discarding or replacing it in `normalizeSSHAuthFailure`.
2. Fix `dialNodeWithPassword` to attempt PAM keyboard-interactive fallback when password authentication fails or is rejected.
3. Build `SSHTrace` logger in `cli/cmdssh/ssh_trace.go` and `cli/cmdssh/ssh_trace_report.go` to track every execution milestone:
   - Target details (`user@host:port`)
   - TCP port 22 reachability probe
   - Host key auto-trust in `known_hosts`
   - Public key authentication attempts across default user keys
   - Password and keyboard-interactive authentication attempts
   - Remote bootstrap commands and results
   - Raw internal error and actionable hints
4. Automatically persist diagnostic reports to `data/logs/ssh_error.log` (and mirror to `.ai-memory/temp/ssh_last_error.log` if workspace temp folder exists).
5. Implement `gitmap ssh error-logs` (aliases: `err`, `errors`, `logs`) to inspect, export (`--json`, `--file`, `--tempfile`), or clear (`--clear`) the last SSH error logs directly from the terminal.
6. Enforce strict error management guidelines (wrap in `*apperror.AppError` with `.WithContext()`, zero swallowed errors, zero bare returns).

---

## 2. 4-Part Root Cause Analysis (RCA)

### Part 1: Symptoms
- User entered SSH password for `a@192.168.1.3`.
- Terminal showed:
  `⚠ Failed to authenticate with remote machine: ssh: authentication failed: invalid password or remote server rejected credentials`
- The stacktrace pointed to `checkEnrollAuth (cmdssh/sshjoin_enroll.go:323)`.
- No details were given about what steps were run, which auth methods were attempted, or what the underlying raw SSH error was.
- No log file was written, preventing the user from inspecting or sharing the internal failure.

### Part 2: Underlying Mechanism
1. In `cli/cmdssh/sshjoin_dial.go:normalizeSSHAuthFailure`:
   The function checked `if strings.Contains(msg, "unable to authenticate") || ...` and unconditionally returned a static string `fmt.Errorf("ssh: authentication failed: invalid password or remote server rejected credentials")`, discarding the original `err` completely.
2. In `cli/cmdssh/sshjoin_dial.go:dialNodeWithPassword`:
   The keyboard-interactive fallback was guarded by `if !isSSHPasswordUnsupported(errPass)`. But `golang.org/x/crypto/ssh` formats auth failure as `ssh: unable to authenticate, attempted methods [none password]...`. Because `errPass.Error()` always contained the substring `"password"`, `isSSHPasswordUnsupported` returned `false`, and `!isSSHPasswordUnsupported` returned `true`, causing keyboard-interactive fallback to never be executed.
3. In `cli/cmdssh/sshjoin_enroll.go:checkEnrollAuth`:
   When authentication failed, it printed only `session.err` and called `apperror.WrapSimple(session.err, "ExecuteSSHJoinEnrollment.auth")`. No execution trace was maintained, no diagnostic file was written to disk, and no context map was attached to the error.
4. Missing CLI Diagnostics Command:
   Unlike pipelines which have `gitmap pipeline error-logs` (`gitmap pe`), SSH operations had no equivalent error log reader command.

### Part 3: Solutions Implemented
1. `cli/cmdssh/ssh_trace.go` & `cli/cmdssh/ssh_trace_report.go`:
   - Defined `SSHExecutionTrace` and `SSHTraceStep`.
   - Records each action (TCP probe, host key trust, public key check, password attempt, keyboard-interactive fallback, remote command bootstrap).
   - Generates formatted terminal reports and automatically writes structured JSON logs to `data/logs/ssh_error.log` and `.ai-memory/temp/ssh_last_error.log`.
2. `cli/cmdssh/sshjoin_dial.go`:
   - Updated `normalizeSSHAuthFailure` to wrap the raw error with `: %w` so internal error details are preserved.
   - Updated `dialNodeWithPassword` to attempt keyboard-interactive PAM fallback when password auth is rejected, recording both attempts in `SSHTrace`.
   - Updated `tryConnectDefaultKey` to record public key checks in `SSHTrace`.
3. `cli/cmdssh/sshjoin_enroll.go`:
   - Added `BeginSSHTrace` at start of `ExecuteSSHJoinEnrollment`.
   - In `resolveTargetClient`: recorded TCP reachability and host key trust milestones; set non-nil error if TCP port 22 is unreachable.
   - In `performConnectedBootstrap`: recorded remote command execution steps without discarding errors with `_ =`.
   - In `checkEnrollAuth`: printed structured execution report, raw internal error, and diagnostic log path. Wrapped in `*apperror.AppError` with full context (`target`, `raw_error`, `log_path`, `user`, `host`, `port`).
4. `cli/cmdssh/ssh_error_logs_cmd.go` & `cli/cmdssh/ssh.go`:
   - Implemented `gitmap ssh error-logs` (aliases: `err`, `errors`, `logs`) supporting `--json`, `--file <path>`, `--tempfile <filename>`, and `--clear`.
   - Exported in `cli/cmdssh/exports.go` and documented in `cli/helptext/ssh.md` and help menus.

### Part 4: Verification
- `check-error-management.py` AST scanner: scanned 3643 files -> PASS (zero violations).
- `05-guideline-autofixer.py`, `08-naming-autofixer.py`, `04-newline-fixer.py`: verified all files conform to coding guidelines.
- `go vet ./cmdssh/...`: zero warnings, exit 0.
- Unit test coverage in `cli/cmdssh/ssh_trace_test.go` added.
- All modified files recorded in test inventory via `33-test-inventory-generator.py --record`.

---

## 3. Files Modified
- `cli/cmdssh/ssh_trace.go` (new): trace data model, step tracker, global active/last trace.
- `cli/cmdssh/ssh_trace_report.go` (new): log persistence to `data/logs/ssh_error.log` and terminal formatting.
- `cli/cmdssh/ssh_error_logs_cmd.go` (new): `gitmap ssh error-logs` CLI handler.
- `cli/cmdssh/ssh_trace_test.go` (new): unit tests for tracer and error-logs CLI.
- `cli/cmdssh/sshjoin_dial.go`: raw error preservation with `%w`, keyboard-interactive fallback, tracer hooks.
- `cli/cmdssh/sshjoin_enroll.go`: trace initialization, TCP probe check, structured failure reporting, context-rich `apperror.Wrap`.
- `cli/cmdssh/ssh.go`: dispatch routing for `error-logs`, `err`, `errors`, `logs`.
- `cli/cmdssh/exports.go`: export `RunSSHErrorLogsCLI`.
- `cli/cmdssh/ssh_help_sections.go`: added `error-logs` to help menu.
- `cli/helptext/ssh.md`: documented `error-logs` in helptext.
