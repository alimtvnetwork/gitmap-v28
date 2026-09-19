# Completed Plan 40: Macro Fleet Export/Import & Multi-Node Interactive Resilience

> **Started by:** User prompt demonstrating failure of `gitmap ssh exec gitmap macro add alim1` across multiple SSH nodes and requesting macro import/export sync commands across SSH nodes and interactive shell resilience.
> **Total Steps / Loops:** 12 steps across Phase 1 planning, 2 parallel execution subagents, and Phase 3 consolidation.
> **Outcome:** 100% verified and completed.

---

## 1. Problem & Root Cause Summary
- **Symptom:** Running `gitmap ssh exec gitmap macro add alim1` across nodes failed on remote nodes with `validation error: macro name and at least one command required`.
- **Root Cause:**
  1. Interactive macro creation requires a controlling TTY and interactive standard input. Over non-interactive SSH execution (`session.CombinedOutput`), `os.Stdin` is disconnected/empty.
  2. Running interactive commands concurrently across multiple nodes in parallel goroutines causes stdin race conditions and cannot work interactively from a single terminal.
  3. Lack of fleet-wide macro synchronization commands (`gitmap macro sync --all` / `gitmap ssh macro sync`) to author macros locally in development and push them to all nodes.

---

## 2. Deliverables & Solutions Implemented

### 2.1 Dual-Layer Interactive Detection & Actionable Advice
- **Remote side (`cli/cmdmacro/macro_add_interactive.go`)**: When `!isTerminalInput()` and `os.Getenv("SSH_CLIENT") != ""` or `len(steps) == 0`, displays clear guidance instead of generic usage errors:
  ```text
  Interactive macro creation cannot run over non-interactive SSH exec. Create locally and sync:
    1. gitmap macro add <name> <cmd1> [cmd2...] (non-interactive)
    2. gitmap macro sync --all (sync to all SSH nodes)
  ```
- **Client side (`cli/cmdssh/sshexec.go`)**: Pre-flight inspection in `isInteractiveMacroAdd(args)` intercepts zero-command `macro add <name>` before network dispatch, preventing unnecessary remote failures.

### 2.2 SSH Macro Fleet Command Suite (`cli/cmdssh/`)
- `gitmap ssh macro sync [flags]`: Pushes local macros to remote SSH node(s) via atomic base64 write to `~/.gitmap/macros/<name>.json`.
- `gitmap ssh macro export [target] [flags]`: Exports macros from a remote node or local machine as JSON.
- `gitmap ssh macro import <target> [flags]`: Pulls macros from a remote node and imports them locally.
- Top-level aliases: `gitmap macro sync --all`, `gitmap macro export --ssh [target]`, `gitmap macro import --from <node>`.

### 2.3 Documentation & Helptext Parity
- Updated `cli/helptext/macro.md` and `cli/helptext/ssh.md` with full usage examples.

---

## 3. Verification
- Targeted unit tests passed:
  - `cli/cmdmacro`: `TestIsRemoteSSHSession`, `TestHandleZeroPipedSteps_RemoteSSH`, `TestHandleZeroInteractiveSteps_RemoteSSH`
  - `cli/cmdssh`: `TestIsInteractiveMacroAdd_TrueBasic`, `TestIsInteractiveMacroAdd_TrueFlags`, `TestIsInteractiveMacroAdd_FalseCommands`, `TestIsInteractiveMacroAdd_FalseOther`, `TestRunSSHExec_InterceptsInteractiveMacroAdd`
- All coding guidelines enforced: functions $\le 15$ lines, affirmative booleans only, universal `*apperror.AppError`, Unix LF line endings.
