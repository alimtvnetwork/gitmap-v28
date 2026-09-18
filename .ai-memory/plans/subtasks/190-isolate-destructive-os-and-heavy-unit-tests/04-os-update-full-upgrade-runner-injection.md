# Subtask 04: OS Update & Full Upgrade Runner Injection

## Parent Plan
- Parent Plan: `.ai-memory/plans/pending/190-isolate-destructive-os-and-heavy-unit-tests.md`

## Objectives
1. Refactor `cli/cmdos/os_update.go`:
   - Introduce `type OSCommandRunner func(cmd *exec.Cmd) error` with `var defaultOSCommandRunner OSCommandRunner = func(cmd *exec.Cmd) error { return cmd.Run() }`.
   - Update `ExecuteOSUpdate` to invoke `defaultOSCommandRunner(cmd)`.
2. Refactor `cli/cmdos/os_full_upgrade.go`:
   - Update `ExecuteOSFullUpgrade` to invoke `defaultOSCommandRunner(cmd)`.
3. Create `cli/cmdos/os_update_test.go`:
   - Test `ExecuteOSUpdate` and `ExecuteOSFullUpgrade` with mock `defaultOSCommandRunner`.
   - Verify that the correct commands are generated for each operating system without executing host updates.

## Coding Standards
- Functions <= 15 lines (target <= 8 lines).
- Affirmative booleans only (`is*`, `has*`).
- Strict Unix LF line endings.
