# Subtask 02: Cluster Lifecycle Shutdown & Logoff Mock Tests

## Parent Plan
- Parent Plan: `.ai-memory/plans/pending/190-isolate-destructive-os-and-heavy-unit-tests.md`

## Objectives
1. Update `cli/cluster/exec_lifecycle_test.go`:
   - Add `TestExecShutdown_Mocked` with `runCmdFunc` mock and `defer` restoration.
   - Add `TestExecLogoff_Mocked` with `runCmdFunc` mock and `defer` restoration.
   - Verify that `ExecShutdown` constructs the correct OS-specific command:
     - Windows: `shutdown /s /t 0`
     - Linux/macOS: `shutdown -h now`
   - Verify that `ExecLogoff` constructs the correct OS-specific command:
     - Windows: `logoff`
     - Linux: shell logout
   - Verify that `runCmdFunc` intercepts command execution without invoking system processes.
2. In `cli/cluster/exec_lifecycle.go`:
   - Ensure all functions adhere to <= 15 lines (target <= 8 lines).
   - Ensure affirmative booleans only.

## Coding Standards
- Functions <= 15 lines (target <= 8 lines).
- Affirmative booleans only (`is*`, `has*`).
- Strict Unix LF line endings.
