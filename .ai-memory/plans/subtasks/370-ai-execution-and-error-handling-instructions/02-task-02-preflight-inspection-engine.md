# Subtask: TASK-02 - PowerShell Argument Pre-Flight Inspection & AI Reflection Engine

- **Parent Task:** `370-ai-execution-and-error-handling-instructions`
- **Subtask ID:** `sub-121928-task-02`
- **Owned Files:** `cli/cmd/cmd_powershell.go`, `cli/cmd/powershell_advice.go`

## Objectives
1. Implement `isPowerShellPathTarget(arg string)` in Go to detect directory paths and non-executable document files (`.md`, `.json`, `.txt`, `.csv`, `.yml`, `.yaml`, etc.).
2. Implement `renderPowerShellPathDiagnostic(target string, isDir bool)` to render Tokyo Night styled diagnostic card with diagnostic questions and replacements.
3. Hook into `RunPowerShell` before process invocation to intercept invalid targets immediately.
