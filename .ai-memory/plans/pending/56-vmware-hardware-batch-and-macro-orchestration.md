# Plan 56: VMware Hardware Customization, Multi-VM Batch Operations, and Macro Orchestration

## Status: Pending
- **Plan ID:** 56
- **Spec Reference:** [02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md](../../../02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md)
- **Scope:** CLI Subsystems, Automation Scripting, Macro Integration, Terminal UI
- **Created At:** 2026-09-30
- **Parent Goal:** Expand VMware Workstation automation from single-file PowerShell scripting (`manage-vm.ps1`) into native GitMap CLI subcommands, multi-VM batch hardware mutations, and macro-executable automation steps with high-aesthetic Bubbletea TUI settings.

---

## 1. Architectural Context & Problem Statement

In Spec 190, the initial single-file PowerShell automation engine (`scripts/vmware/manage-vm.ps1` and `repo-secrets/vmware/manage-vm.ps1`) provided core functionality:
- VM discovery (`Scan`), hardware inspection (`Inspect`), static MAC configuration (`Set-MAC`), network adapter configuration (`Set-Network`), basic RAM/vCPU adjustments (`Set-Hardware`), virtual disk repair (`Repair-Disks`), printer toggle (`Set-Printer`), snapshot management (`List-Snapshots`, `Create-Snapshot`, `Clone-Snapshot`), and diagnostic log clipboard sharing (`-CopyLog`).

However, the user explicitly mandated future-phase capabilities:
1. **Dynamic Hardware Expansion:** Support adding new virtual disks (SCSI, NVMe, SATA), expanding existing virtual disk capacities without data loss, and adjusting controller topologies.
2. **Multi-VM Batch Operations:** Apply hardware specifications (e.g. allocating 16 GB RAM and 8 vCPUs, or updating network adapters) across multiple virtual machines in a folder or cluster fleet simultaneously.
3. **Macro Engine Orchestration:** Expose VMware management commands as first-class, idempotent GitMap CLI commands (`gitmap vm ...`) so that macros can sequence VM lifecycle events (e.g., configure VM -> start VM -> execute SSH deploy -> take snapshot).
4. **Interactive TUI Hardware Customization:** Provide a Bubbletea-based terminal interface (`gitmap vm ui` or `gitmap os vm`) to interactively inspect and reconfigure VM hardware settings.

---

## 2. Subtasks Breakdown

### Subtask 56.1: Virtual Disk Lifecycle & Capacity Expansion
- **Spec Reference:** Spec 190 §4.5
- **Files:** `scripts/vmware/manage-vm.ps1`, `cli/cmdvm/vm_disk.go`
- **Actions:**
  - Implement `-Action Add-Disk`: Attach new `.vmdk` disk images to `.vmx` files with user-defined bus type (NVMe, SCSI, SATA) and capacity.
  - Implement `-Action Expand-Disk`: Automate `vmware-vdiskmanager.exe -x <size>GB <path>` with pre-expansion backup safety (`.vmdk.bak`).
  - Add native Go wrapper `gitmap vm disk add <vmx> <sizeGB>` and `gitmap vm disk expand <vmx> <sizeGB>`.

### Subtask 56.2: Dynamic RAM, CPU & Hotplug Resource Management
- **Spec Reference:** Spec 190 §4.5
- **Files:** `scripts/vmware/manage-vm.ps1`, `cli/cmdvm/vm_hardware.go`
- **Actions:**
  - Extend `-Action Set-Hardware` to support core topology tuning (`cpuid.coresPerSocket`, `numvcpus`, `memsize`).
  - Add memory and CPU hot-add flags (`mem.hotadd = "TRUE"`, `vcpu.hotadd = "TRUE"`).
  - Provide validation against host physical boundaries to prevent over-allocation crashes.

### Subtask 56.3: Multi-VM Discovery and Batch Operations
- **Spec Reference:** Spec 190 §4.5
- **Files:** `cli/cmdvm/vm_batch.go`, `cli/cmdvm/vm_scan.go`
- **Actions:**
  - Implement batch targeting via wildcard globbing or folder scanning: `gitmap vm batch --dir <path> --set-ram 8192 --set-cpu 8`.
  - Concurrently process multiple `.vmx` configurations while acquiring file locks to prevent configuration corruption.
  - Render progress bars and execution summaries using standard boxed table rendering.

### Subtask 56.4: GitMap CLI Subsystem & Macro Orchestration
- **Spec Reference:** Spec 190 §4.1, §4.2
- **Files:** `cli/cmd/root.go`, `cli/cmdvm/vm_root.go`, `cli/macro/execute.go`
- **Actions:**
  - Register `gitmap vm` command root in CLI dispatcher with subcommands (`scan`, `inspect`, `mac`, `net`, `hardware`, `disk`, `snap`, `batch`).
  - Seamlessly enable macros to include `gitmap vm ...` steps with cross-platform argument quoting and error propagation.
  - Log all VM mutation operations to `TaskHistory` in SQLite split-DB (`tasks_root`).

### Subtask 56.5: Interactive VMware Settings Terminal TUI
- **Spec Reference:** Spec 190 §4.3
- **Files:** `cli/cmdvm/vm_tui.go`, `cli/cmdvm/vm_tui_model.go`, `cli/cmdvm/vm_tui_view.go`
- **Actions:**
  - Build interactive Bubbletea UI to navigate discovered VMs, inspect hardware settings, and tweak MAC/RAM/CPU/Network/Disks visually.
  - Provide instant keyboard shortcuts (`m` for MAC, `r` for RAM, `c` for CPU, `s` for Snapshot, `b` for Backup).
  - Include single-key clipboard copy of diagnostic logs (`c`).

### Subtask 56.6: End-to-End Verification & Acceptance Testing
- **Spec Reference:** Spec 190 §5
- **Files:** `cli/cmdvm/vm_test.go`, `scripts/vmware/manage-vm.tests.ps1`
- **Actions:**
  - Create mock `.vmx` test fixtures under `cli/cmdvm/testdata/`.
  - Verify disk addition, expansion, RAM modification, CPU modification, and snapshot operations preserve existing `.vmx` keys without data loss.
  - Validate macro execution containing `gitmap vm` steps executes successfully and audits to `TaskHistory`.

---

## 3. Acceptance Criteria

- [ ] All virtual disk additions and capacity expansions safely modify `.vmx` and `.vmdk` files with automated `.bak` backups.
- [ ] Multi-VM batch operations execute reliably across multiple targets with unified progress reporting.
- [ ] `gitmap vm` subcommands can be sequenced within macros and run via `gitmap <macro_name>` with zero regressions.
- [ ] Mutation operations log structured records into `TaskHistory` viewable via `gitmap task history --section vm`.
- [ ] Full unit test suite passes with 100% green coverage on mock VM configurations.
