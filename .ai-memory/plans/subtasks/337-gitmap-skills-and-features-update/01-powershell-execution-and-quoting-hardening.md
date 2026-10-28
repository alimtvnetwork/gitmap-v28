# Subtask 01: PowerShell Argument Execution Hardening

**Master Plan:** `.ai-memory/plans/pending/337-gitmap-skills-and-features-update.md`  
**Target:** `cli/cmd/cmd_powershell.go`, `cli/cmd/cmd_powershell_test.go`  

## Objectives
- Coalesce `-Command` and `-c` arguments into a safe script invocation on Windows.
- Ensure `-NoProfile` is always applied when running commands.
- Prevent Windows argument parsing from breaking variable assignments (`$a = 123`).
- Add comprehensive unit tests in `cli/cmd/cmd_powershell_test.go`.
