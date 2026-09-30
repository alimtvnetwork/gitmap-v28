# Subtask [05]: Interactive VMware Settings Terminal TUI
Traceability ID: Task-56-05
Spec Reference: [02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md](../../../../02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md)
Target Files: [cli/cmdvm/vm_tui.go, cli/cmdvm/vm_tui_model.go, cli/cmdvm/vm_tui_view.go]

Action:
- Author interactive Bubbletea TUI for discovering and reconfiguring VMware virtual machines.
- Include interactive tabs: Hardware, Network, Disks, Snapshots, and Diagnostics.
- Support key bindings for quick editing and clipboard log sharing.

Acceptance Criteria:
- TUI runs seamlessly on both Windows Command Prompt/PowerShell and Linux terminals.
- Exit and redraw transitions restore terminal buffer cleanly.
