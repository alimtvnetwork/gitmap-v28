# Subtask 02 - Leaf Error Returns Refactoring (Zero Dual-Handling)

## Parent Specification
[28-error-management-audit.md](.lovable/plans/pending/28-error-management-audit.md)

## Acceptance Criteria & Requirements
- Refactor target command leaf functions to return `error` or `*apperror.AppError` directly rather than calling `cliexit.HandleError` internally and returning `nil`.
- Files to refactor:
  - `gitmap/cmd/vscode_cmd.go`: `runVSCode` and `dispatchVSCodeAction` returns error to caller.
  - `gitmap/cmd/reinstall.go`: On aborted reinstall, return error directly without proceeding.
  - `gitmap/cmd/macro_cmd.go`: Return validation/delete errors directly.
  - `gitmap/cmd/workflow_open_pr.go`: Return validation errors directly.
  - `gitmap/cmd/sshcat.go`: Return errors directly.
  - `gitmap/cmd/sshgen.go`: Return errors directly.
- Maintain function lengths <= 15 lines, zero nested ifs, and proper vertical spacing.
