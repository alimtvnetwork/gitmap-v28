# Master Plan 54: VMware Automation, Macro Idempotent Removal & Edit UX, Audit Task Logging, and Chained Installer

Spec Reference: [02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md](../../../02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md)
RCA Reference: [02-spec/22-app-issues/53-macro-phantom-steps-and-removal-failure-rca.md](../../../02-spec/22-app-issues/53-macro-phantom-steps-and-removal-failure-rca.md)

## Status: Active
- **Steps Budget:** N = 200, Phase 1 = 100 steps, Phase 2 = 100 steps
- **Concurrency Mode:** A = 2, H = 2 (up to 4 concurrent subtask operations)
- **Constraint Mode:** Strict no-build and no-test execution

---

## 1. Domain Context & Blast Radius Analysis

### 1.1 Context
The user identified two major functional gaps and requested an advanced virtualization automation toolchain:
1. Macro execution failed on Windows PowerShell when running `rm test` because PowerShell's `Remove-Item` fails if the target is absent. In addition, interactive macro edit (`macro edit`) failed to persist user-entered commands like `mkdir -p test` and `ls` because they were intercepted by in-builder helpers and discarded upon exit.
2. In-depth task/audit tracking is needed across macros, SSH commands, and installers, along with a dedicated, polished terminal task history view.
3. Macro export and import must be tested and connected to multi-machine deployment.
4. Installer needs OS targeting (Windows, Unix, Ubuntu, CentOS) and chained sub-installer delegation.
5. A comprehensive single-file PowerShell automation script for VMware Workstation (`manage-vm.ps1`) must be created in the `repo-secrets` folder to control MAC addresses, network adapters, disks, RAM, vCPUs, printers, and snapshots with rich stack traces and clipboard log export.

### 1.2 Blast Radius
- `cli/macro/`: Step execution and shell command transformation.
- `cli/cmdmacro/`: Macro interactive editor, helper dispatch, export/import, and deploy.
- `cli/cmdtask/`: Task history terminal table formatting and inspection.
- `cli/cmdssh/`, `cli/cmdinstall/`: Audit interception hooks for SSH and installer events.
- `secrets/vmware/`: Standalone PowerShell automation engine.

---

## 2. Actionable Deliverable Traceability Matrix

| Deliverable ID | Subtask File | Description | Target Files |
|---|---|---|---|
| **Task-01** | `01-macro-safe-rm-and-editor-fix.md` | Idempotent removal shim (`rm-if-exists`) & interactive editor step appending | `cli/macro/execute.go`, `cli/cmdmacro/macro_edit.go`, `cli/cmdmacro/macro_add_helpers.go` |
| **Task-02** | `02-macro-export-import-and-deploy.md` | Macro export, import, and fleet SSH deployment | `cli/cmdmacro/macro_export.go`, `cli/cmdmacro/macro_import.go`, `cli/cmdmacro/macro_deploy_ssh.go`, `cli/cmdssh/ssh_deploy_router.go` |
| **Task-03** | `03-audit-task-logging-and-tui-viewer.md` | Universal audit interceptor & polished terminal task history view | `cli/cmdtask/task_history_cmd.go`, `cli/cmdmacro/macro_cmd.go`, `cli/cmdssh/ssh_history_tasks.go`, `cli/store/tasks_split_db.go` |
| **Task-04** | `04-modular-installer-and-subinstaller-chain.md` | OS-specific installer & sub-installer command chaining | `cli/cmdinstall/install_archive_deploy.go`, `cli/cmdinstall/installer_types.go` |
| **Task-05** | `05-vmware-powershell-automation-script.md` | Comprehensive single-file VMware automation PowerShell script | `D:/work/repo-secrets/vmware/manage-vm.ps1`, `scripts/vmware/manage-vm.ps1` |
| **Task-06** | `06-future-roadmap-and-specs-consolidation.md` | Architectural specs, RCA reports, and plan consolidation | `02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md`, `02-spec/22-app-issues/53-macro-phantom-steps-and-removal-failure-rca.md` |
