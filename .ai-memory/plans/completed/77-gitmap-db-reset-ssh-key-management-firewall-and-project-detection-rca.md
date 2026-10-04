# Plan 77: GitMap Database Reset, SSH Key Lifecycle, Cross-OS Firewall, and DetectedProject Foreign Key RCA

## Status: Completed
- **Plan ID:** 77
- **Spec Reference:** [02-spec/21-app/207-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca/01-architecture-spec.md](../../02-spec/21-app/207-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca/01-architecture-spec.md)
- **Scope:** Database Engine, SSH Key Lifecycle, Cross-OS Firewall, Web UI, Scanner Project Detection
- **Completed At:** 2026-10-04
- **Parent Goal:** Implement full database purge and reseed (`gitmap reset`), reorder SSH terminal output, add `gitmap ssh view`, safeguard `gitmap ssh create` with overwrite prompt and undo/redo, fix DetectedProject foreign key bug, implement cross-OS SSH port and firewall automation, introduce `gitmap ssh ui`, and build intelligent SSH connection troubleshooting.

---

## 1. Executive Summary

This plan coordinates the end-to-end execution of Spec 207 across seven key areas:
1. Root cause analysis and remediation of foreign key constraint errors in project detection.
2. Complete SQLite database reset and clean reseed pipeline.
3. SSH terminal output reordering and dedicated key viewer command (`gitmap ssh view`).
4. SSH key creation overwrite guard with backup and undo/redo task registration.
5. Cross-platform SSH port and firewall automation (Windows, Ubuntu, CentOS, macOS).
6. Local browser-based SSH management UI (`gitmap ssh ui`).
7. Guided 7-step SSH connection diagnostics and troubleshooting engine.

---

## 2. Subtasks Breakdown

### Subtask 77.1: DetectedProject Foreign Key Constraint RCA & Fix
- **Spec Reference:** Spec 207 §2
- **Scope:** `cli/cmdscanner/`, `cli/store/`, `cli/models/`
- **Actions:**
  - Audit `DetectedProject` table definition and foreign keys (`RepoId` and `ProjectTypeId`).
  - Pre-seed all required `ProjectType` records before scanning.
  - Verify that `RepoId` is committed and valid prior to inserting child `DetectedProject` records.
  - Add transactional safety and defensive checks.

### Subtask 77.2: Comprehensive Database Reset & Reseed Engine
- **Spec Reference:** Spec 207 §3
- **Scope:** `cli/cmd/`, `cli/store/`
- **Actions:**
  - Implement `gitmap reset` / `gitmap db reset` command.
  - Close active connections and purge all database files in `.gitmap/`.
  - Report all removed databases with size metrics.
  - Trigger automatic reseeding of clean schemas and default data.

### Subtask 77.3: SSH Output Layout Reorder & `gitmap ssh view` Command
- **Spec Reference:** Spec 207 §4.1, §4.2
- **Scope:** `cli/cmdssh/`
- **Actions:**
  - Reorder `gitmap ssh` terminal output so SSH key info appears at the very bottom.
  - Implement `gitmap ssh view` (`v`, `show`, `key`) showing only SSH key info.
  - Implement automated clipboard copy of the public key.
  - Add guidance tips for recreating/regenerating the key.

### Subtask 77.4: `gitmap ssh create` Overwrite Guard, Backup Undo/Redo & Scripts Fixer Sync
- **Spec Reference:** Spec 207 §4.3
- **Scope:** `cli/cmdssh/`, `cli/system/`, `scripts-fixer`
- **Actions:**
  - Check for existing key and prompt user `[y/N]` unless `-y`/`--confirm` passed.
  - Create timestamped backup file (`id_rsa.bak.<timestamp>`).
  - Register undo/redo task in GitMap journal.
  - Update `scripts-fixer` SSH regeneration logic.

### Subtask 77.5: Cross-OS SSH Port Management & Firewall Automation
- **Spec Reference:** Spec 207 §5
- **Scope:** `cli/cmdssh/`, `cli/firewall/`
- **Actions:**
  - Implement `gitmap ssh port add <port>`, `gitmap ssh port rm <port>`, `gitmap ssh port ls`.
  - Implement `gitmap ssh enable [--port <port>]` and `gitmap ssh disable`.
  - Implement `gitmap ssh enable-public <port>` handling firewall and binding.
  - Build cross-OS adapters: Windows (Netsh/PowerShell), Ubuntu (ufw), CentOS (firewalld), macOS (pfctl).

### Subtask 77.6: Interactive SSH Web UI Command (`gitmap ssh ui`)
- **Spec Reference:** Spec 207 §6
- **Scope:** `cli/cmdssh/`, `cli/tui/`, web assets
- **Actions:**
  - Implement `gitmap ssh ui` HTTP server and browser auto-launcher.
  - Build interactive web dashboard displaying keys, ports, firewall rules, and fleet nodes.

### Subtask 77.7: Guided SSH Connection Failure Diagnostics & 7-Step Troubleshooting Engine
- **Spec Reference:** Spec 207 §7
- **Scope:** `cli/cmdssh/`, `cli/apperror/`
- **Actions:**
  - Intercept SSH dial/connect errors.
  - Render formatted 7-step actionable troubleshooting guide.
