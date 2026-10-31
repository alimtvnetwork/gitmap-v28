# Subtask: TASK-03 - PowerShell Error Post-Processing & Fallback Interceptor

- **Parent Task:** `370-ai-execution-and-error-handling-instructions`
- **Subtask ID:** `sub-121928-task-03`
- **Owned Files:** `cli/cmd/powershell_interceptor.go`, `cli/cmd/cmd_powershell_test.go`

## Objectives
1. Implement post-execution error interceptor for PowerShell exit status 1.
2. If PowerShell outputs "is not recognized as a name of a cmdlet, function, script file, or executable program", surface educational diagnostic without raw Go stack trace.
3. Add comprehensive unit tests in `cli/cmd/cmd_powershell_test.go` verifying interception behavior.
