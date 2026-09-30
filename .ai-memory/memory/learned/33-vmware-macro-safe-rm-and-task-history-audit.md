# Learned Memory 33: VMware Automation Engine, Macro Idempotent Removal, Interactive Edit Recording, and Task Audit History

Date: 2026-09-30
Related Specs:
- [02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md](../../../02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md)
- [02-spec/22-app-issues/53-macro-phantom-steps-and-removal-failure-rca.md](../../../02-spec/22-app-issues/53-macro-phantom-steps-and-removal-failure-rca.md)
- [.ai-memory/plans/completed/54-vmware-macro-audit-task-and-installer-chain.md](../../plans/completed/54-vmware-macro-audit-task-and-installer-chain.md)

## 1. Architectural Lessons & Principles

### 1.1 Cross-Platform Idempotent File Deletion
On Windows, `rm` maps to PowerShell's `Remove-Item` cmdlet, which throws terminating `ItemNotFoundException` if the path is missing.
- In macro execution (`cli/macro/safe_rm.go`) and live interactive runner (`resolveSingleLiveCmd`), wrap removal commands in:
  ```powershell
  foreach ($__target in @(<paths>)) { if (Test-Path -LiteralPath $__target) { Remove-Item -Recurse -Force -LiteralPath $__target } }
  ```
- Parse arguments with quote-aware tokenization (`tokenizeCommandArgs`) to properly handle quoted paths with spaces without regex corruption or space-split damage.
- Provide native `gitmap safe-rm <path...>` / `gitmap rm-safe` CLI command with exit code 0 on missing targets.

### 1.2 Interactive Macro Editor UX Contract
- Never silently drop shell commands typed into an interactive editor.
- Ephemeral in-builder inspection tools must be explicitly distinguished (e.g. `:ls`, `:cat`, `pwd`, `cd`, `help`).
- All other commands typed by the user must be appended directly to the macro step array.

### 1.3 Subsystem Task Audit Logging
- Unify task mutation tracking across `macro`, `ssh`, and `installer` subsystems into `TaskHistory` in SQLite split-DB (`tasks_root`).
- Use standardized fields: `TaskId`, `Section`, `Action`, `Target`, `ForwardPayload`, `InversePayload`, `Status`, `ExecutedAt`.
- Present records in terminal via boxed Unicode tables with color badges and section filtering (`--section`).

### 1.4 Single-File PowerShell VMware Automation
- Centralize VMware Workstation management in `repo-secrets/vmware/manage-vm.ps1`.
- Directly manipulate `.vmx` configuration keys (MAC addresses, network connection types, RAM, vCPUs, disks, printers).
- Support automatic backup `.vmx.bak` before in-place file modifications.
- Implement `-CopyLog` switch using Windows `Set-Clipboard` for instantaneous diagnostic log sharing.
