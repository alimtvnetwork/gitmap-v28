# Subtask 01: SSH Multi-Command & GitMap Remote Execution Verification

## Scope
- Inspect `cli/cmdssh/` implementation for running multiple chained commands, subcommands, and remote GitMap calls.
- Verify `gitmap ssh exec <host> "cmd1 && cmd2"` parsing, shell wrapping, execution streaming, and error handling.
- Verify remote GitMap CLI detection, command invocation, and output forwarding.
- Confirm detailed help text in `cli/cmdssh/` commands and `cli/helptext/ssh.md`.

## Acceptance Criteria
- [x] SSH execution command supports single and multiple commands with quotes and operators (`&&`, `;`, `|`).
- [x] SSH execution handles remote `gitmap` command execution cleanly without double `gitmap gitmap` prefix bug.
- [x] Expanded `isGitmapCommand` to recognize all top-level GitMap subcommands.
- [x] Timeout and dial guards prevent hung connections via `CheckConnLiveness`.
- [x] Terminal help text accurately reflects usage with rich examples for single, multiple, and remote GitMap commands.

## Completed Changes
- Fixed `runSSHWorker` in `cli/cmdssh/sshexec.go`: resolved `commandStr` via `resolveGitmapCommandString(args)` instead of blind prepending.
- Expanded `isGitmapCommand` in `cli/cmdssh/ssh_exec_command.go` to cover all GitMap subcommands with clean <= 15 line switch.
- Added `printSSHExecHelp` and `printSSHExecExamples` with `-h/--help` flag interception and validation in `cli/cmdssh/sshexec.go`.
- Added documentation and examples in `cli/helptext/ssh.md` and `src/data/commands.ts`.
