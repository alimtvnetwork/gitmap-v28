# Master Plan 54: VMware Automation, Macro Idempotent Removal & Edit UX, Audit Task Logging, and Chained Installer

Spec Reference: [02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md](../../../02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md)
RCA Reference: [02-spec/22-app-issues/53-macro-phantom-steps-and-removal-failure-rca.md](../../../02-spec/22-app-issues/53-macro-phantom-steps-and-removal-failure-rca.md)

## Status: Completed
- **Steps Budget:** N = 200, Phase 1 = 100 steps, Phase 2 = 100 steps
- **Concurrency Mode:** A = 2, H = 2 (up to 4 concurrent subtask operations)
- **Constraint Mode:** Full build and test verification passed (zero runtime failures, zero guideline violations)

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
- `cli/macro/`: Step execution and shell command transformation (`safe_rm.go`, `execute.go`).
- `cli/cmdmacro/`: Macro interactive editor, helper dispatch, export/import, and deploy (`macro_edit.go`, `macro_add.go`, `macro_cmd.go`).
- `cli/cmdtask/`: Task history terminal table formatting, audit logging, and inspection (`task_history_cmd.go`, `task_audit.go`).
- `cli/cmdssh/`, `cli/cmdinstall/`: Audit interception hooks for SSH and installer events (`cluster_exec_runner.go`, `install.go`, `install_add.go`).
- `scripts/vmware/`, `D:/work/repo-secrets/vmware/`: Standalone PowerShell automation engine (`manage-vm.ps1`).

---

## 2. Actionable Deliverable Traceability Matrix

| Deliverable ID | Subtask File | Description | Target Files | Status |
|---|---|---|---|---|
| **Task-01** | `01-macro-safe-rm-and-editor-fix.md` | Idempotent removal shim (`rm-if-exists`) & interactive editor step appending | `cli/macro/safe_rm.go`, `cli/macro/execute.go`, `cli/cmdmacro/macro_edit.go` | Completed |
| **Task-02** | `02-macro-export-import-and-deploy.md` | Macro export, import, and fleet SSH deployment | `cli/cmdmacro/macro_export.go`, `cli/cmdmacro/macro_import.go`, `cli/cmdmacro/macro_deploy_ssh.go` | Completed |
| **Task-03** | `03-audit-task-logging-and-tui-viewer.md` | Universal audit interceptor & polished terminal task history view | `cli/cmdtask/task_history_cmd.go`, `cli/cmdtask/task_audit.go`, `cli/store/tasks_split_db.go`, `cli/model/pendingtask.go` | Completed |
| **Task-04** | `04-modular-installer-and-subinstaller-chain.md` | OS-specific installer & sub-installer command chaining | `cli/cmdinstall/install.go`, `cli/cmdinstall/install_add.go` | Completed |
| **Task-05** | `05-vmware-powershell-automation-script.md` | Comprehensive single-file VMware automation PowerShell script | `D:/work/repo-secrets/vmware/manage-vm.ps1`, `scripts/vmware/manage-vm.ps1` | Completed |
| **Task-06** | `06-future-roadmap-and-specs-consolidation.md` | Architectural specs, RCA reports, and plan consolidation | `02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md`, `02-spec/22-app-issues/53-macro-phantom-steps-and-removal-failure-rca.md` | Completed |

---

## 3. Detailed Implementation Outcomes

### 3.1 Task-01: Macro Safe Removal Shim and Interactive Edit Step Recording Fix
- **Idempotent Removal Shim (`cli/macro/safe_rm.go`):** Implemented `AdaptCommandForPlatform` and `transformToWindowsSafeRemoval`. Detects `rm`, `rmdir`, `del`, `rd`, and `Remove-Item` on Windows, tokenizing quoted paths and wrapping them in:
  `foreach ($__target in @(...)) { if (Test-Path -LiteralPath $__target) { Remove-Item -Recurse -Force -LiteralPath $__target } }`
- **Interactive Editor Step Recording (`cli/cmdmacro/macro_edit.go`):**
  - Resolved issue where user-entered commands (e.g., `mkdir -p test`) ran live but were omitted from the macro step list. Now, shell commands are both executed live and automatically appended to `m.Steps`.
  - Refined `isDeleteStepCmd` to require numeric argument tokens so commands starting with `del` or `rm` (e.g. `rm test`) fall through as shell commands rather than step-deletion directives.

### 3.2 Task-02: Macro Export, Import, and Multi-Node Fleet Deployment
- Verified `gitmap macro export <name>` and `gitmap macro import <file>` workflows for JSON and YAML payloads.
- Verified fleet distribution via `gitmap macro deploy <name>` connecting to SSH cluster nodes.

### 3.3 Task-03: Universal Audit Task Logging and Polished Terminal History View
- **Universal Audit Interceptor (`cli/cmdtask/task_audit.go`):** Added `RecordTaskAudit(section, action, target, forwardPayload, status)` interfacing with `TasksSplitDB.InsertTaskHistory`.
- **Wired Subsystems:**
  - Macro: `macro add`, `macro edit`, `macro rm`, `macro run`
  - SSH: `cluster exec`
  - Installer: `install <tool>`, `install add`, `install update`
- **TUI History Table (`cli/cmdtask/task_history_cmd.go`):** Rendered table with cyan headers, rounded box boundaries, colored section badges (`⚡ MACRO`, `🔑 SSH`, `📦 INSTALL`), status indicators (`[✔ COMPLETED]`, `[✖ FAILED]`), humanized timestamps, and `--section` filtering with limit/offset pagination.

### 3.4 Task-04: OS-Specific Modular Installer and Chained Sub-Installer Engine
- Supported multi-OS installer targets (`win`, `unix`, `ubuntu`, `all`) in `cli/cmdinstall/install_add.go`.
- Added audit instrumentation in `executeInstall` and `persistInstallerScript` in `cli/cmdinstall/`.

### 3.5 Task-05: Single-File Enterprise VMware Automation PowerShell Script
- Authored 980+ line production PowerShell engine at `scripts/vmware/manage-vm.ps1` and mirrored to `D:/work/repo-secrets/vmware/manage-vm.ps1`.
- Fully supports: `Scan`, `Inspect`, `Set-MAC` (e.g. `00:50:56:38:57:7B`), `Set-Network` (Bridged, NAT, HostOnly, Custom), `Set-Hardware` (RAM, vCPUs), `Manage-Disk` (List, Repair, Expand), `Manage-Printer` (remove/disable), `Manage-Snapshot` (list, create, clone, revert), and `Optimize-VM`.
- Includes automated clipboard export (`-CopyLog`) and detailed error diagnostics with script stack traces.
- Mirrored and committed cleanly to `repo-secrets` git repository.

### 3.6 Task-06: Future Architecture Specifications and Plan Consolidation
- Authored canonical specification: `02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md` and indexed in `02-spec/21-app/readme.md`.
- Authored RCA: `02-spec/22-app-issues/53-macro-phantom-steps-and-removal-failure-rca.md` and indexed in `02-spec/22-app-issues/01-index.md`.
- Updated master plan registry in `.ai-memory/plans/readme.md`.
