# Plan 78: DetectedProject Foreign Key RCA, Comprehensive DB Reset, SSH Lifecycle & Cross-OS Firewall Hardening

## Status: Completed
- **Plan ID:** 78
- **Spec Reference:** [02-spec/21-app/208-detectedproject-foreign-key-and-ssh-lifecycle-hardening/01-architecture-spec.md](../../02-spec/21-app/208-detectedproject-foreign-key-and-ssh-lifecycle-hardening/01-architecture-spec.md)
- **Scope:** Scanner Subsystem, SQLite Split-DB, OpenSSH Lifecycle, Cross-OS Firewall, Web UI
- **Completed At:** 2026-10-04
- **Parent Goal:** Fix all causes of SQLite Error 787 during project detection, verify comprehensive DB reset & reseed, reorder SSH output layout, add dedicated `gitmap ssh view`, safeguard `gitmap ssh create` with overwrite prompt and undo/redo, align scripts-fixer, implement cross-OS SSH port and firewall automation, and provide SSH UI & 7-step diagnostics.

---

## 1. Subtasks Overview

1. **Subtask 78.1: DetectedProject FOREIGN KEY Constraint (787) In-Depth RCA & Fix**
   - File: `.ai-memory/plans/subtasks/208-detectedproject-foreign-key-and-ssh-lifecycle-hardening/01-detectedproject-foreign-key-rca-and-fix.md`
   - Scope: `cli/cmdscan/scanprojects.go`, `cli/store/project.go`, `cli/store/repo.go`, `cli/constants/constants_project_sql.go`

2. **Subtask 78.2: Comprehensive Database Reset & Reseed Engine**
   - File: `.ai-memory/plans/subtasks/208-detectedproject-foreign-key-and-ssh-lifecycle-hardening/02-database-reset-and-reseed-engine.md`
   - Scope: `cli/cmddb/cmddb_reset.go`, `cli/cmd/reset.go`, `cli/store/storage_inventory.go`

3. **Subtask 78.3: SSH Output Layout & Dedicated Key Viewer (`gitmap ssh view`)**
   - File: `.ai-memory/plans/subtasks/208-detectedproject-foreign-key-and-ssh-lifecycle-hardening/03-ssh-output-layout-and-ssh-view.md`
   - Scope: `cli/cmdssh/ssh.go`, `cli/cmdssh/sshcat.go`, `cli/cmdssh/ssh_help_menu.go`

4. **Subtask 78.4: `gitmap ssh create` Overwrite Guard, Backup & Undo/Redo Enqueue**
   - File: `.ai-memory/plans/subtasks/208-detectedproject-foreign-key-and-ssh-lifecycle-hardening/04-ssh-create-overwrite-guard-and-backup-undo-redo.md`
   - Scope: `cli/cmdssh/sshgen.go`, `cli/cmdssh/sshexisting.go`, `cli/cmdssh/ssh_history_tasks.go`

5. **Subtask 78.5: Scripts-Fixer Project SSH Key Regeneration Alignment**
   - File: `.ai-memory/plans/subtasks/208-detectedproject-foreign-key-and-ssh-lifecycle-hardening/05-scripts-fixer-ssh-regeneration-alignment.md`
   - Scope: `scripts-fixer/`, `03-ai-scripts/`, `scripts/`

6. **Subtask 78.6: Cross-OS SSH Port Management & Firewall Automation**
   - File: `.ai-memory/plans/subtasks/208-detectedproject-foreign-key-and-ssh-lifecycle-hardening/06-cross-os-ssh-port-and-firewall-automation.md`
   - Scope: `cli/firewall/`, `cli/cmdssh/ssh_port_config.go`

7. **Subtask 78.7: Interactive SSH Web UI & 7-Step Troubleshooting Engine**
   - File: `.ai-memory/plans/subtasks/208-detectedproject-foreign-key-and-ssh-lifecycle-hardening/07-guided-ssh-troubleshooting-and-diagnostics.md`
   - Scope: `cli/cmdssh/ssh_ui.go`, `cli/cmdssh/ssh_troubleshoot.go`, `cli/cmdssh/ssh_target_exec.go`
