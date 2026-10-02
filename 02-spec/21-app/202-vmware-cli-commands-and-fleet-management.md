# 202 — VMware Workstation CLI Command Suite, Fleet Management, Dual-Engine Hypervisor Driver, and Coordinated Scheduler Shutdown

## Status: Active
- **Spec ID:** 202
- **Scope:** Application CLI, VMware Fleet Automation, Dual-Engine Hypervisor Driver, Split-DB Persistence
- **Created At:** 2026-10-02
- **Reference Prompt:** `01-prompts/42-vmware-cli-commands.md`

---

## 1. Executive Summary & Architectural Motivation

Virtual machine fleet orchestration across development and testing workstations requires standard, repeatable, and high-speed CLI interfaces. In modern containerized and virtualized workflows, operators require the ability to:
1. Manage virtual machine inventories with logical grouping and active default contexts, eliminating repetitive target flag typing.
2. Control full VM lifecycles (power on, graceful shutdown, hard termination, suspend, reset) across local Workstation and Player hypervisors.
3. Perform offline virtual hardware modifications, such as dynamic MAC address regeneration, network interface rebinding, virtual disk expansion, and `.vmx` configuration cleanup.
4. Execute snapshot tree operations (create, list, revert, clone, delete) safely with lock detection.
5. Coordinate automated shutdown procedures across all running guest VMs and background scheduler pipelines prior to host maintenance or reboots.

To maintain architectural cohesion across GitMap, the VMware CLI subsystem (`gitmap vm` / `gitmap vmware`) achieves **1:1 structural and behavioral parity** with GitMap's existing SSH fleet management subsystem (`cli/cmdssh`). It adopts the same target selector grammar, Split-DB SQLite inventory persistence (`installation.db` / `repodb`), two-column boxed terminal help rendering, and strict coding guidelines.

### 1.1 Core Design Pillars

```text
┌────────────────────────────────────────────────────────────────────────┐
│                      GitMap CLI Routing Layer                          │
│               [ gitmap vm ]     or     [ gitmap vmware ]               │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │
                                    ▼
┌────────────────────────────────────────────────────────────────────────┐
│                        Target Selector Resolver                        │
│   <alias> | <a1>,<a2> | @<group> | --group <g> | all | [Omitted: Def] │
└───────────────────┬────────────────────────────────┬───────────────────┘
                    │                                │
                    ▼                                ▼
┌──────────────────────────────────────┐ ┌───────────────────────────────┐
│     SQLite Split-DB Storage Layer    │ │   Dual-Engine Driver Layer    │
│  - vm_instances                      │ │  - Native vmrun.exe Driver    │
│  - vm_groups                         │ │  - PowerShell JSON Bridge     │
│  - vm_group_memberships              │ │  - vdiskmanager Utility       │
└──────────────────────────────────────┘ └───────────────────────────────┘
```

1. **Dual-Engine Hypervisor Driver:** Combines high-speed direct binary dispatch (`vmrun.exe`) for power and snapshot operations with a robust PowerShell engine bridge (`scripts/vmware/manage-vm.ps1 -Json`) for offline `.vmx` attribute transformations and virtual disk resizing.
2. **Flexible Target Resolution:** Operators can target single machines, comma-separated lists, logical groups via `@group`, explicit flags, `all`, or rely on an active default group fallback when the target is omitted.
3. **Split-DB SQLite Persistence:** All VM registrations, logical groups, hardware snapshots, and default states are persisted in GitMap's split database structure, ensuring immediate querying without scanning the disk filesystem.
4. **Resilient Coordinated Shutdown:** Introduces `scheduler shutdown` which pauses active job queues, soft-stops running virtual guests with an ACPI timeout window, and cleanly terminates background tasks.

---

## 2. Target Selector Mechanics & Default Group Fallback

Target resolution for all actionable subcommands (`status`, `start`, `stop`, `mac`, `disk`, `snapshot`, `vmx`) is standardized through a shared selector parser (`cli/cmdvmware/selector.go`).

### 2.1 Resolution Hierarchy

The selector resolves targets in the following order of precedence:

| Target Syntax | Grammar Pattern | Resolution Behavior | Concrete CLI Example |
| :--- | :--- | :--- | :--- |
| **Single VM Alias** | `<alias>` | Resolves a single unique VM record by alias. Errors with `E4001` if not found. | `gitmap vm start win11-dev` |
| **Comma-Separated List** | `<alias1>,<alias2>,...` | Resolves multiple distinct VMs. Ignores duplicates and processes in order. | `gitmap vm start win11-dev,ubuntu-ci` |
| **Group Prefix Selector** | `@<group-name>` | Queries all enabled VMs linked to `<group-name>` in `vm_group_memberships`. | `gitmap vm status @staging` |
| **Explicit Group Flag** | `--group <name>` | Filters execution to all enabled VMs within the specified group. | `gitmap vm stop --group runners` |
| **All Inventory Keyword** | `all` or `--all` | Targets every active VM record registered in `vm_instances`. | `gitmap vm status all` |
| **Omitted Target (Fallback)** | *[No target specified]* | Queries SQLite for the active default group (`is_default = 1`) and resolves all member VMs. | `gitmap vm status` |

### 2.2 Default Group Fallback Behavior

When an operator executes a command without specifying a target:
1. The resolver executes `SELECT id, name FROM vm_groups WHERE is_default = 1 LIMIT 1`.
2. If a default group is found:
   - GitMap fetches all enabled VMs associated with the group.
   - If the group has no members, the command outputs a helpful message: `Default group '<group>' has no registered VMs. Use 'gitmap vm group assign <alias> <group>' to add members.`
   - Otherwise, execution proceeds against all member VMs.
3. If no default group is configured:
   - GitMap halts execution with validation code `E4007`.
   - The CLI renders available groups and instructs the operator: `No target specified and no default group configured. Set one using: gitmap vm group set-default <group-name>`.

---

## 3. Exhaustive Subcommand Suite & CLI Grammar

The command entry point is registered under `gitmap vm` and aliased to `gitmap vmware`.

```text
gitmap vm <subcommand> [target] [flags]
```

### 3.1 `gitmap vm add`

Registers an existing virtual machine `.vmx` file into the SQLite inventory.

- **CLI Grammar:** `gitmap vm add <alias> <vmx-path> [flags]`
- **Required Arguments:**
  - `<alias>`: Unique alphanumeric identifier (lowercase, hyphens allowed; e.g. `win11-node1`).
  - `<vmx-path>`: Path to the virtual machine `.vmx` file (relative or environment variable paths supported).
- **Flags:**
  - `--group <name>`: Initial group assignment (defaults to `default`).
  - `--description <text>`: Free-form description of VM role and specifications.
  - `--set-default`: Immediately designate this VM as the default standalone target.
- **Operational Flow:**
  1. Validates that `<alias>` does not already exist in `vm_instances`.
  2. Verifies that the `.vmx` configuration file exists and contains valid VMware markers (`config.version`).
  3. Inspects hardware configuration (RAM, vCPUs, primary MAC address) via configuration file inspection or PowerShell bridge.
  4. Inserts record into `vm_instances` and creates group linkage in `vm_group_memberships`.
- **Examples:**
  ```bash
  gitmap vm add win11-dev vms/win11/win11.vmx --group dev --description "Windows 11 Dev Box"
  gitmap vm add ubuntu-runner vms/ubuntu/ubuntu.vmx --set-default
  ```

### 3.2 `gitmap vm rm`

Unregisters a virtual machine from the inventory.

- **CLI Grammar:** `gitmap vm rm <alias> [flags]`
- **Required Arguments:**
  - `<alias>`: Target virtual machine alias.
- **Flags:**
  - `--force`: Bypass interactive confirmation prompt.
  - `--delete-files`: Permanently purge the `.vmx` directory and virtual disk files from storage.
- **Operational Flow:**
  1. Verifies the VM is currently powered off. If running, rejects removal unless `--force` is set.
  2. Removes group linkages in `vm_group_memberships`.
  3. Removes record from `vm_instances`.
  4. If `--delete-files` is provided, safely deletes VM storage directories.
- **Examples:**
  ```bash
  gitmap vm rm win11-dev
  gitmap vm rm old-node --force --delete-files
  ```

### 3.3 `gitmap vm install`

Verifies or initiates installation of VMware hypervisor tools and utilities.

- **CLI Grammar:** `gitmap vm install [flags]`
- **Flags:**
  - `--silent`: Execute non-interactive unattended installation.
  - `--verify-only`: Inspect system without altering configuration or downloading packages.
- **Operational Flow:**
  1. On Windows: Probes for `vmrun.exe` in the discovery hierarchy. If absent and not in verify mode, invokes `winget install VMware.WorkstationPro` or `VMware.WorkstationPlayer`.
  2. On Linux: Inspects `open-vm-tools` or `vmware` workstation installation via package manager (`apt-get install -y open-vm-tools`).
  3. Renders diagnostic report verifying hypervisor utility paths and kernel module availability.
- **Examples:**
  ```bash
  gitmap vm install --verify-only
  gitmap vm install --silent
  ```

### 3.4 `gitmap vm group`

Manages logical clusters of virtual machines for batch targeting.

- **CLI Grammar:** `gitmap vm group <subaction> [args] [flags]`
- **Subactions:**
  - `add <name>`: Creates a new logical group. Optional `--description <desc>`.
  - `rm <name>`: Deletes a logical group (does not delete registered VMs).
  - `ls`: Lists all registered groups, member counts, and indicates the default group.
  - `assign <alias> <group>`: Adds a VM to the specified group.
  - `remove <alias> <group>`: Removes a VM from the group.
  - `set-default <name>`: Marks `<name>` as the default fallback group.
- **Flags:**
  - `--description <desc>`: Group description text.
- **Examples:**
  ```bash
  gitmap vm group add staging --description "Staging test nodes"
  gitmap vm group assign win11-dev staging
  gitmap vm group set-default staging
  gitmap vm group ls
  ```

### 3.5 `gitmap vm default`

Inspects or mutates active default target context.

- **CLI Grammar:** `gitmap vm default [target]`
- **Operational Flow:**
  - When called with no arguments: Displays the current default group and any default VM.
  - When called with a group name or VM alias: Sets the default group or default VM context in SQLite.
- **Examples:**
  ```bash
  gitmap vm default
  gitmap vm default staging
  ```

### 3.6 `gitmap vm ls`

Displays registered virtual machines in a formatted terminal table.

- **CLI Grammar:** `gitmap vm ls [flags]`
- **Flags:**
  - `--group <name>`: Restrict listing to VMs belonging to `<name>`.
  - `--all`: Show all VMs across all groups.
  - `--json`: Output machine-readable JSON array.
- **Columns Rendered:**
  - `ALIAS`: Unique machine identifier.
  - `STATUS`: Live power status (`RUNNING`, `STOPPED`, `SUSPENDED`).
  - `GROUPS`: Comma-separated list of group memberships.
  - `MAC`: Primary network interface MAC address.
  - `RAM / VCPUS`: Configured memory in MB and vCPU count.
  - `DEFAULT`: `[DEFAULT]` indicator if marked.
- **Examples:**
  ```bash
  gitmap vm ls
  gitmap vm ls --group staging --json
  ```

### 3.7 `gitmap vm status`

Queries real-time hypervisor execution status and guest networking for target VMs.

- **CLI Grammar:** `gitmap vm status [target] [flags]`
- **Flags:**
  - `--detailed`: Include disk controllers, network type, and snapshot count.
  - `--json`: Output structured status JSON.
- **Operational Flow:**
  1. Resolves target list via selector mechanics.
  2. Invocations query `vmrun.exe -T ws list` to detect running `.vmx` files.
  3. For running machines, interrogates guest IP via VMware Tools or ARP table cache.
  4. For stopped machines, inspects `.lck` directories to ensure clean shutdown.
- **Examples:**
  ```bash
  gitmap vm status
  gitmap vm status win11-dev --detailed
  gitmap vm status @dev --json
  ```

### 3.8 `gitmap vm start`

Powers on target virtual machines.

- **CLI Grammar:** `gitmap vm start [target] [flags]`
- **Flags:**
  - `--nogui`: Start headless in the background without UI window (default: `true`).
  - `--gui`: Launch with the VMware graphical console window.
  - `--timeout <secs>`: Hypervisor launch confirmation timeout in seconds (default: `60`).
- **Operational Flow:**
  1. Resolves targets.
  2. For each VM, checks if already running via `vmrun.exe -T ws list`.
  3. Executes `vmrun.exe -T ws start "<vmx-path>" nogui` (or `gui`).
  4. Polls running list until VM path appears or timeout expires.
- **Examples:**
  ```bash
  gitmap vm start win11-dev
  gitmap vm start @dev --gui
  gitmap vm start --timeout 90
  ```

### 3.9 `gitmap vm stop`

Initiates virtual machine shutdown, power-off, or suspension.

- **CLI Grammar:** `gitmap vm stop [target] [flags]`
- **Flags:**
  - `--soft`: Graceful OS shutdown via VMware Tools ACPI signal (default: `true`).
  - `--hard` / `--force`: Immediate hard power cut (unclean power off).
  - `--suspend`: Suspend execution state to disk (`.vmss`).
  - `--timeout <secs>`: Graceful shutdown waiting window in seconds (default: `90`).
- **Operational Flow:**
  1. Resolves targets.
  2. If `--soft` (default): Executes `vmrun.exe -T ws stop "<vmx-path>" soft`.
  3. Polls process list up to timeout. If VM does not stop and `--force` is specified, escalates to `vmrun.exe -T ws stop "<vmx-path>" hard`.
  4. If `--suspend`: Executes `vmrun.exe -T ws suspend "<vmx-path>"`.
- **Examples:**
  ```bash
  gitmap vm stop win11-dev
  gitmap vm stop @staging --soft --timeout 120
  gitmap vm stop win11-dev --hard
  ```

### 3.10 `gitmap vm mac`

Safely updates or regenerates virtual network adapter MAC addresses.

- **CLI Grammar:** `gitmap vm mac <alias> <mac|generate> [flags]`
- **Required Arguments:**
  - `<alias>`: Target virtual machine alias.
  - `<mac|generate>`: Explicit MAC address (format: `00:50:56:XX:YY:ZZ`) or literal keyword `generate`.
- **Flags:**
  - `--adapter <name>`: Network adapter key (default: `ethernet0`).
  - `--type <static|generated>`: MAC allocation type.
  - `--force`: Bypass lock checks.
- **Operational Flow:**
  1. Verifies target VM is completely powered off (offline mutation requirement).
  2. Creates automatic backup of configuration file (`<vmx-path>.bak`).
  3. Dispatches mutation to `scripts/vmware/manage-vm.ps1` with `-Action Set-MAC` and `-Json`.
  4. Verifies updated `.vmx` contents on disk.
  5. Updates `primary_mac` in `vm_instances` SQLite record.
- **Examples:**
  ```bash
  gitmap vm mac win11-dev generate
  gitmap vm mac win11-dev 00:50:56:38:57:7B --adapter ethernet0
  ```

### 3.11 `gitmap vm disk`

Manipulates virtual machine storage volumes (`.vmdk`).

- **CLI Grammar:** `gitmap vm disk <alias> <subaction> [flags]`
- **Subactions:**
  - `list`: Enumerates virtual disks, adapter channels, sparse/preallocated mode, and sizes.
  - `expand`: Expands virtual disk capacity. Requires `--size <size>`.
  - `repair`: Checks and repairs corrupted virtual disk structures via `vmware-vdiskmanager.exe -R`.
- **Flags:**
  - `--disk <path|channel>`: Specific virtual disk path or SCSI/NVMe channel.
  - `--size <size>`: Target expanded disk capacity (e.g. `120GB`, `250GB`).
- **Operational Flow:**
  1. Verifies VM is stopped.
  2. For `expand`: Invokes `vmware-vdiskmanager.exe -x <size> "<vmdk-path>"`.
  3. For `repair`: Invokes `vmware-vdiskmanager.exe -R "<vmdk-path>"`.
  4. Parses results and updates database records.
- **Examples:**
  ```bash
  gitmap vm disk win11-dev list
  gitmap vm disk win11-dev expand --size 120GB
  gitmap vm disk win11-dev repair --disk win11.vmdk
  ```

### 3.12 `gitmap vm snapshot`

Governs virtual machine snapshot trees.

- **CLI Grammar:** `gitmap vm snapshot <alias> <subaction> [flags]`
- **Subactions:**
  - `list`: Displays snapshot tree hierarchy with timestamps.
  - `create <name>`: Captures live or offline snapshot state.
  - `revert <name>`: Restores machine state to named snapshot.
  - `delete <name>`: Deletes and merges snapshot delta disks.
  - `clone <dest>`: Clones VM or snapshot to a new target directory.
- **Flags:**
  - `--name <name>`: Snapshot label identifier.
  - `--description <desc>`: Informational note for snapshot.
- **Operational Flow:**
  - Invocations map directly to `vmrun.exe -T ws snapshot`, `revertToSnapshot`, `deleteSnapshot`, and `clone`.
- **Examples:**
  ```bash
  gitmap vm snapshot win11-dev list
  gitmap vm snapshot win11-dev create clean-base --description "Initial clean OS state"
  gitmap vm snapshot win11-dev revert clean-base
  ```

### 3.13 `gitmap vm vmx`

Performs low-level configuration maintenance, sanitization, and optimization on `.vmx` files.

- **CLI Grammar:** `gitmap vm vmx <alias> <subaction> [flags]`
- **Subactions:**
  - `inspect`: Pretty-prints parsed `.vmx` key-value pairs categorized by subsystem (CPU, RAM, Network, Display, Disks).
  - `repair`: Removes duplicate keys, repairs syntax corruption, and normalizes line endings.
  - `optimize`: Injects enterprise performance settings (disables host memory trimming `MemTrimRate = "0"`, disables background paging `mainMem.useNamedFile = "FALSE"`).
- **Examples:**
  ```bash
  gitmap vm vmx win11-dev inspect
  gitmap vm vmx win11-dev repair
  gitmap vm vmx win11-dev optimize
  ```

### 3.14 `gitmap vm scheduler shutdown`

Coordinates host maintenance by cleanly draining scheduler queues and gracefully stopping all running virtual machines.

- **CLI Grammar:** `gitmap vm scheduler shutdown [flags]`
- **Flags:**
  - `--all-vms`: Ensure all running VMs are stopped before scheduler termination (default: `true`).
  - `--force`: Force kill running VMs if graceful shutdown times out.
  - `--timeout <secs>`: Coordinated drain and shutdown timeout in seconds (default: `120`).
- **Operational Flow:**
  1. Pauses GitMap background task scheduler and rejects incoming job dispatches.
  2. Queries all active VMs currently in `RUNNING` state via `vmrun.exe -T ws list`.
  3. Dispatches parallel graceful soft-shutdown signals (`vmrun stop <path> soft`) to all running machines.
  4. Continuously polls VM process states concurrently with a progress spinner.
  5. If timeout expires and `--force` is active, terminates remaining machines with hard stop.
  6. Flushes scheduler queue state and logs coordinated shutdown telemetry to `repo-secrets`.
- **Examples:**
  ```bash
  gitmap vm scheduler shutdown
  gitmap vm scheduler shutdown --timeout 180 --force
  ```

---

## 4. Dual-Engine Driver Architecture & Execution Bridges

The subsystem relies on an abstracted driver interface implemented across two underlying execution engines.

```text
                     ┌───────────────────────────┐
                     │     HypervisorDriver      │
                     │        (Interface)        │
                     └─────────────┬─────────────┘
                                   │
                 ┌─────────────────┴─────────────────┐
                 ▼                                   ▼
   ┌───────────────────────────┐       ┌───────────────────────────┐
   │    Native vmrun Driver    │       │  PowerShell Script Bridge │
   │ (cli/cmdvmware/driver_vm) │       │ (cli/cmdvmware/driver_ps) │
   └─────────────┬─────────────┘       └─────────────┬─────────────┘
                 │                                   │
                 ▼                                   ▼
          vmrun.exe -T ws              scripts/vmware/manage-vm.ps1
```

### 4.1 Go HypervisorDriver Interface Contract

```go
type HypervisorDriver interface {
    IsInstalled() bool
    ListRunningVMs() ([]string, *apperror.AppError)
    StartVM(vmxPath string, isHeadless bool, timeoutSecs int) *apperror.AppError
    StopVM(vmxPath string, mode StopMode, timeoutSecs int) *apperror.AppError
    GetVMStatus(vmxPath string) (VMStatusResult, *apperror.AppError)
    ModifyMAC(vmxPath string, adapter string, mac string, isGenerated bool) *apperror.AppError
    ManageDisk(vmxPath string, action DiskAction, size string) *apperror.AppError
    ManageSnapshot(vmxPath string, action SnapshotAction, name string) *apperror.AppError
    RepairVMX(vmxPath string) *apperror.AppError
}
```

### 4.2 Native `vmrun.exe` Discovery Hierarchy

GitMap searches for hypervisor binaries across Windows and Linux platforms in strict candidate order:
1. System `PATH` via `exec.LookPath("vmrun.exe")` / `exec.LookPath("vmrun")`.
2. Standard 64-bit Directory: `%ProgramFiles%/VMware/VMware Workstation/vmrun.exe`.
3. Standard 32-bit Compatibility: `%ProgramFiles(x86)%/VMware/VMware Workstation/vmrun.exe`.
4. VMware Player: `%ProgramFiles(x86)%/VMware/VMware Player/vmrun.exe`.
5. VMware VIX SDK: `%ProgramFiles(x86)%/VMware/VMware VIX/vmrun.exe`.

If no binary is discovered, all driver calls return structured error `E4005` detailing search candidate paths.

### 4.3 PowerShell Automation Bridge (`scripts/vmware/manage-vm.ps1`)

For operations requiring atomic `.vmx` property updates, regex replacement, or disk analysis, the driver delegates to `scripts/vmware/manage-vm.ps1` using the `-Json` switch.

- **Command Pattern:**
  ```powershell
  powershell.exe -NoProfile -ExecutionPolicy Bypass -File scripts/vmware/manage-vm.ps1 -Action Set-MAC -VMPath "<vmx-path>" -MACAddress "<mac>" -Adapter ethernet0 -Json
  ```
- **Standardized JSON Envelope Output:**
  ```json
  {
    "status": "SUCCESS",
    "action": "Set-MAC",
    "vmPath": "vms/win11/win11.vmx",
    "adapter": "ethernet0",
    "macAddress": "00:50:56:38:57:7B",
    "isGenerated": false,
    "backupPath": "vms/win11/win11.vmx.bak"
  }
  ```
- **Safety Invariant:** The PowerShell script generates an automatic `.bak` backup copy prior to modifying any file. In case of serialization errors or mid-process crashes, the backup file is restored automatically.

---

## 5. SQLite Split-DB Inventory Schema & Data Models

VM inventory and logical group associations are persisted inside GitMap's Split-DB engine (`installation.db` or `repodb`).

### 5.1 Relational DDL Specification

```sql
CREATE TABLE IF NOT EXISTS vm_instances (
    id TEXT PRIMARY KEY,
    alias TEXT NOT NULL UNIQUE,
    vmx_path TEXT NOT NULL,
    description TEXT DEFAULT '',
    primary_mac TEXT DEFAULT '',
    memory_mb INTEGER DEFAULT 2048,
    vcpu_count INTEGER DEFAULT 2,
    is_active INTEGER NOT NULL DEFAULT 1,
    is_running INTEGER NOT NULL DEFAULT 0,
    has_tools INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS vm_groups (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT DEFAULT '',
    is_default INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS vm_group_memberships (
    id TEXT PRIMARY KEY,
    group_id TEXT NOT NULL,
    vm_id TEXT NOT NULL,
    is_enabled INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL,
    FOREIGN KEY(group_id) REFERENCES vm_groups(id) ON DELETE CASCADE,
    FOREIGN KEY(vm_id) REFERENCES vm_instances(id) ON DELETE CASCADE,
    UNIQUE(group_id, vm_id)
);

CREATE INDEX IF NOT EXISTS idx_vm_instances_alias ON vm_instances(alias);
CREATE INDEX IF NOT EXISTS idx_vm_groups_name ON vm_groups(name);
CREATE INDEX IF NOT EXISTS idx_vm_group_memberships_lookup ON vm_group_memberships(group_id, vm_id);
```

### 5.2 Go Domain Structs (`cli/cmdvmware/types.go`)

```go
package cmdvmware

import "time"

type VmInstance struct {
	ID          string    `json:"id"`
	Alias       string    `json:"alias"`
	VmxPath     string    `json:"vmxPath"`
	Description string    `json:"description"`
	PrimaryMac  string    `json:"primaryMac"`
	MemoryMb    int       `json:"memoryMb"`
	VcpuCount   int       `json:"vcpuCount"`
	IsActive    bool      `json:"isActive"`
	IsRunning   bool      `json:"isRunning"`
	HasTools    bool      `json:"hasTools"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type VmGroup struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsDefault   bool      `json:"isDefault"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type VmGroupMembership struct {
	ID        string    `json:"id"`
	GroupID   string    `json:"groupId"`
	VmID      string    `json:"vmId"`
	IsEnabled bool      `json:"isEnabled"`
	CreatedAt time.Time `json:"createdAt"`
}
```

---

## 6. Strict Coding Guidelines & Implementation Rules

Every file authored in `cli/cmdvmware/...` must adhere strictly to repository coding standards:

1. **File Size Cap ($\le 100$ Lines):**
   - No source file may exceed 100 lines (including imports and whitespace).
   - Code is decomposed into discrete single-responsibility files: `cmd.go`, `add.go`, `rm.go`, `ls.go`, `group.go`, `default.go`, `status.go`, `start.go`, `stop.go`, `mac.go`, `disk.go`, `snapshot.go`, `vmx.go`, `scheduler.go`, `driver_vmrun.go`, `driver_ps.go`, `selector.go`, and `render.go`.
2. **Function Body Length ($8-15$ Lines):**
   - Functions and methods must stay within 15 lines. Extract helper functions for flag parsing, table row rendering, and SQL statement execution.
3. **Single-Level Branching & Guard Clauses:**
   - Maximum cyclomatic indentation $\le 1$. Invert conditional checks and return early on errors.
4. **Positive Boolean Naming:**
   - Booleans must use affirmative prefixes (`is` or `has`): `isActive`, `isRunning`, `isDefault`, `isEnabled`, `hasTools`, `isSilent`, `isForce`. Negative booleans like `isNotRunning` or `disabled` are strictly forbidden.
5. **Structured AppError Catalog:**
   - Errors must be constructed via `apperror.NewWithDetails` using standard typed error codes:
     - `E4001`: VM alias not found in database.
     - `E4002`: VMX configuration path does not exist on disk.
     - `E4003`: Invalid or malformed MAC address syntax.
     - `E4004`: Target VM is running (offline mutation rejected).
     - `E4005`: Hypervisor driver binary (`vmrun.exe`) not found.
     - `E4006`: PowerShell automation bridge execution failure.
     - `E4007`: VM group does not exist or has no members.
     - `E4008`: Virtual disk expansion or repair failure.
     - `E4009`: Snapshot operation rejected by hypervisor.
     - `E4010`: Coordinated scheduler shutdown timeout exceeded.
6. **Centralized Domain Types:**
   - All structs, options, and DTOs live in `cli/cmdvmware/types.go` or `cli/cmdvmware/store/types.go`. No inline struct declarations.
7. **Zero Absolute Paths (R11):**
   - Total ban on absolute paths (`C:\...`, `/home/...`, `file:///...`) in source code and documentation. All paths must be relative to repository root.

---

## 7. End-to-End Testing & Verification Protocol

Verification of the VMware fleet management subsystem involves testing command dispatch, selector resolution, and driver execution against fixture VMs.

### 7.1 Fixture VM Generation

For test isolation, a synthetic lightweight `.vmx` fixture is created inside a temporary test directory:
- Includes dummy virtual hardware entries (`config.version = "8"`, `virtualHW.version = "19"`, `memsize = "2048"`, `numvcpus = "2"`).
- Network interface `ethernet0` configured with `addressType = "generated"` and valid MAC seed.

### 7.2 Verification Checklist

- [x] Target selector parses single alias, comma-separated lists, `@group`, and defaults to active default group when omitted.
- [x] Default group mutation via `gitmap vm group set-default <name>` persists in SQLite and governs fallback dispatch.
- [x] Dual-engine driver locates `vmrun.exe` through the 5-step search hierarchy and cleanly falls back to PowerShell bridge.
- [x] `gitmap vm mac` backs up `.vmx` to `.vmx.bak`, regenerates MAC, and updates SQLite inventory.
- [x] `gitmap vm scheduler shutdown` broadcasts soft shutdown signals, awaits termination, and safely exits within timeout.
- [x] All authored Go source files adhere to $\le 100$ lines and functions stay within $8-15$ lines.
