# Subtask [01]: Macro Safe Removal Shim and Interactive Edit Step Recording Fix
Traceability ID: Task-01
Spec Reference: [02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md](../../../02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md)
Target Files: [cli/macro/execute.go, cli/cmdmacro/macro_edit.go, cli/cmdmacro/macro_add_helpers.go]
Action:
- In `cli/macro/execute.go`, add cross-platform idempotent file/folder removal wrapping: if running under Windows PowerShell and the command starts with `rm `, `rmdir `, or `Remove-Item `, safely check path existence before deleting to eliminate non-zero exit codes when deleting non-existent files.
- In `cli/cmdmacro/macro_edit.go`, modify `handleEditSpecialAction` and `processInBuilderCommand` so that standard modification commands (`mkdir`, `touch`, `mkfile`) entered during `macro edit` are recorded directly into the macro step list rather than requiring a secondary `+add` command.
- Update terminal feedback in `macro_edit.go` to clearly confirm appended steps.
Acceptance Criteria:
- Running `rm non_existent_file` via macro step does not abort with exit code 1 on Windows PowerShell.
- Commands typed in `gitmap macro edit <name>` (such as `mkdir -p test`) are appended and persisted upon typing `done` or `exit`.
- Unit and sanity flows pass without compilation errors.
Targeted Verification: [python 03-ai-scripts/05-guideline-autofixer.py cli/macro/execute.go]
