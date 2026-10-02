---
name: vmware-cli
description: Comprehensive reference, command grammar, and execution workflows for VMware Workstation CLI commands in GitMap (`gitmap vm`), including VM fleet management, default groups, lifecycle operations, MAC mutation, snapshot management, disk expansion, and coordinated scheduler shutdown.
---

# VMware Workstation CLI & Fleet Management Skill (`vmware-cli`)

## Mission & Architectural Scope
This skill provides authoritative reference documentation, CLI command grammar, operational workflows, and Go implementation specifications for VMware Workstation fleet management in GitMap (`gitmap vm` and its canonical alias `gitmap vmware`).

Modeled with 1:1 structural parity after the GitMap SSH fleet management subsystem (`cli/cmdssh`), this subsystem enables autonomous agents and operators to register virtual machines, organize them into logical clusters, designate default working groups, execute power lifecycles, safely mutate hardware parameters (such as MAC addresses and disk allocations), manage snapshot trees, optimize `.vmx` configurations, and coordinate graceful scheduler shutdowns.

---

## 1. Architectural Foundation & Component Map

The VMware CLI subsystem is organized into clean, single-responsibility components with strict separation between user-facing CLI routing, target selection, driver abstractions, and Split-DB persistence.

```
┌────────────────────────────────────────────────────────────────────────┐
│                      GitMap CLI Subsystem Entry                        │
│                     (gitmap vm / gitmap vmware)                        │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │
                                    ▼
┌────────────────────────────────────────────────────────────────────────┐
│                   Target Selector & Group Resolver                     │
│      (<alias>, <list>, @<group>, all, or active default group)         │
└───────────────────┬───────────────────────────────┬────────────────────┘
                    │                               │
                    ▼                               ▼
┌──────────────────────────────────────┐  ┌──────────────────────────────┐
│        SQLite Split-DB Storage       │  │   Hypervisor Driver Layer    │
│           (installation.db)          │  │     (HypervisorDriver)       │
├──────────────────────────────────────┤  ├──────────────────────────────┤
│ • vm_instances                       │  │ • Native vmrun.exe Runner    │
│ • vm_groups                          │  │ • PowerShell Automation JSON │
│ • vm_group_memberships               │  │   Bridge (manage-vm.ps1)     │
└──────────────────────────────────────┘  └──────────────┬───────────────┘
                                                         │
                                                         ▼
                                          ┌──────────────────────────────┐
                                          │      VMware Workstation      │
                                          │     Guest VMs & .vmx files   │
                                          └──────────────────────────────┘
```

### Component Code Map

| Component / Responsibility | Target File Path | Description |
|---|---|---|
| **CLI Entrypoint & Routing** | `cli/cmdvmware/cmd.go` | Subcommand registration, alias mapping (`vm` / `vmware`), root flag parsing. |
| **Domain Models & Enums** | `cli/cmdvmware/types.go` | Centralized Go structs (`VmInstance`, `VmGroup`, etc.), enums, option wrappers. |
| **Database Schema & Migrations** | `cli/cmdvmware/store/schema.go` | DDL creation, migration logic, and indices in SQLite Split-DB. |
| **Instance Storage Store** | `cli/cmdvmware/store/instance_store.go` | CRUD queries and state updates for VM instances. |
| **Group Storage Store** | `cli/cmdvmware/store/group_store.go` | Group creation, deletion, listing, and default group toggle. |
| **Membership Storage Store** | `cli/cmdvmware/store/membership_store.go` | Association mapping between instances and groups. |
| **Target Selector Resolver** | `cli/cmdvmware/selector.go` | Target grammar resolution (`alias`, `@group`, `all`, default group fallback). |
| **Driver Discovery & Runner** | `cli/cmdvmware/discovery.go` | Candidate executable lookup for `vmrun.exe` and `vmware-vdiskmanager.exe`. |
| **Native `vmrun.exe` Driver** | `cli/cmdvmware/driver_vmrun.go` | Direct process execution with `-T ws` for power, status, and snapshot trees. |
| **PowerShell Driver Bridge** | `cli/cmdvmware/driver_ps.go` | JSON bridge executing `scripts/vmware/manage-vm.ps1` for MAC and disk ops. |
| **Inventory Subcommands** | `cli/cmdvmware/add.go`, `rm.go`, `ls.go` | VM registration, unregistration, and table listing. |
| **Group & Default Commands** | `cli/cmdvmware/group.go`, `default.go` | Group clustering and default target inspection/assignment. |
| **Power Lifecycle Commands** | `cli/cmdvmware/start.go`, `stop.go`, `status.go` | Headless/GUI start, graceful soft shutdown, hard kill, suspension, status. |
| **Hardware & VMX Maintenance** | `cli/cmdvmware/mac.go`, `disk.go`, `vmx.go` | MAC mutation, disk expansion/repair, VMX validation and line normalization. |
| **Snapshot Subcommands** | `cli/cmdvmware/snapshot.go` | Snapshot tree viewing, creation, rollback, deletion, and cloning. |
| **Scheduler Shutdown** | `cli/cmdvmware/scheduler.go` | Coordinated fleet drain and background scheduler daemon termination. |
| **Terminal UI & Help Styling** | `cli/cmdvmware/render.go`, `help.go` | Two-column help formatting, color status badges, and table formatting. |

---

## 2. Hypervisor Detection & Candidate Discovery

GitMap interfaces with VMware Workstation through native executables. Because installations vary across drive letters, editions, and architecture tiers, GitMap utilizes an ordered candidate discovery hierarchy.

### Candidate Search Hierarchy

1. **System `PATH` Lookup:** Probes `exec.LookPath("vmrun.exe")` and `exec.LookPath("vmware-vdiskmanager.exe")`.
2. **64-bit Workstation Directory:** `%ProgramFiles%/VMware/VMware Workstation/vmrun.exe`.
3. **32-bit Compatibility Directory:** `%ProgramFiles(x86)%/VMware/VMware Workstation/vmrun.exe`.
4. **VMware Player Fallback:** `%ProgramFiles(x86)%/VMware/VMware Player/vmrun.exe`.
5. **VMware VIX SDK Fallback:** `%ProgramFiles(x86)%/VMware/VMware VIX/vmrun.exe`.

### Dual-Driver Integration Strategy

1. **Native Executable Dispatch (`vmrun.exe`):**
   - High-performance, direct execution with `-T ws` (Workstation hypervisor target).
   - Used for lifecycle events: `start`, `stop`, `reset`, `suspend`, `snapshot`, `revertToSnapshot`, `deleteSnapshot`, and `list`.
2. **PowerShell Automation Engine (`scripts/vmware/manage-vm.ps1`):**
   - Dedicated bridge executed with `-NoProfile -ExecutionPolicy Bypass -File scripts/vmware/manage-vm.ps1 -Json`.
   - Used for complex `.vmx` regex replacements, automated static/generated MAC mutation with `.vmx.bak` safety backups, and virtual disk verification.
   - Emits structured JSON on stdout, allowing direct Go unmarshaling into domain structures.

---

## 3. Target Selector Mechanics & Default Group Fallback

All lifecycle, hardware, and diagnostic commands accept flexible target expressions to avoid repetitive flag passing:

| Expression Syntax | Target Resolution | Concrete Example |
|---|---|---|
| `<alias>` | Exactly one registered VM instance | `gitmap vm start win11-dev` |
| `<alias1>,<alias2>` | Discrete comma-separated list of VMs | `gitmap vm start win11-dev,ubuntu-ci` |
| `@<group-name>` | All VMs belonging to the specified group | `gitmap vm status @dev` |
| `--group <name>` | Flag filter targeting a specific group | `gitmap vm stop --group testing` |
| `all` or `--all` | Every registered VM in the inventory | `gitmap vm status all` |
| *[Omitted]* | Resolves automatically to the **active default group** | `gitmap vm status` |

### Default Group Fallback Logic
1. When no target is specified on the command line:
   - Query `vm_groups` for `is_default = 1`.
   - If a default group exists: query `vm_group_memberships` for active members and execute across that set.
   - If no default group is configured: abort execution with error code `E4007`, outputting available groups and suggesting `gitmap vm group set-default <name>`.

---

## 4. Subcommand Cheat Sheet & Operational Syntax

```text
gitmap vm <subcommand> [target] [flags]
```

### Complete Subcommand Matrix

| Command | Full Syntax | Primary Purpose & Behavior | Key Flags |
|---|---|---|---|
| **`add`** | `gitmap vm add <alias> <vmx-path>` | Inspects `.vmx`, validates syntax, extracts RAM/vCPU/MAC, persists record into SQLite, and maps to initial group. | `--group <name>`: Target group (default: `default`)<br>`--description <text>`: Notes<br>`--set-default`: Make active default VM |
| **`rm`** | `gitmap vm rm <alias>` | Removes VM record from database and group mappings. Retains disk files by default. | `--force`: Bypass interactive confirmation<br>`--delete-files`: Purge `.vmx` and `.vmdk` from disk |
| **`install`** | `gitmap vm install` | Detects hypervisor binaries. Offers unattended installation via `winget` (Windows) or `open-vm-tools` (Linux). | `--silent`: Headless execution<br>`--verify-only`: Diagnostic check without install |
| **`group`** | `gitmap vm group <action> [args]` | Logical clustering management (`add`, `rm`, `ls`, `assign`, `remove`, `set-default`). | `--description <desc>`: Group purpose |
| **`default`** | `gitmap vm default [target]` | Without args: renders active default group and default VM. With target: sets new default context. | None |
| **`ls`** | `gitmap vm ls` | Renders registered inventory table with power state, group, MAC, RAM, and `[DEFAULT]` badges. | `--group <name>`: Filter by group<br>`--all`: Full inventory<br>`--json`: Output JSON array |
| **`status`** | `gitmap vm status [target]` | Queries live hypervisor state via `vmrun list`, guest IP from VMware Tools or ARP, lock files. | `--detailed`: Include disk controllers and snapshot tree<br>`--json`: Raw JSON |
| **`start`** | `gitmap vm start [target]` | Powers on VM. Executes `vmrun -T ws start "<path>" nogui` (or GUI console if requested). | `--nogui`: Headless start (default)<br>`--gui`: Launch with GUI window<br>`--timeout <secs>`: Wait timeout (default: `60s`) |
| **`stop`** | `gitmap vm stop [target]` | Initiates graceful ACPI shutdown via VMware Tools (`vmrun stop "<path>" soft`). Can escalate to hard kill. | `--soft`: Graceful ACPI shutdown (default)<br>`--hard` / `--force`: Immediate power cut<br>`--suspend`: Save state to `.vmss`<br>`--timeout <secs>`: Timeout (default: `90s`) |
| **`mac`** | `gitmap vm mac <alias> <mac\|generate>` | Mutates network adapter MAC. Verifies VM is powered off, backs up `.vmx.bak`, writes new MAC, verifies. | `--adapter <name>`: Target NIC (default: `ethernet0`)<br>`--type <static\|generated>`: MAC address mode<br>`--force`: Override lock checks |
| **`disk`** | `gitmap vm disk <alias> <action>` | Virtual disk maintenance: `list` (show channels/sizes), `expand` (resizes disk via `vmware-vdiskmanager`), `repair` (repairs `.vmdk`). | `--disk <path\|channel>`: Specific disk target<br>`--size <size>`: Target size (e.g. `120GB`) |
| **`snapshot`** | `gitmap vm snapshot <alias> <action>` | Snapshot management: `list` (tree view), `create <name>`, `revert <name>`, `delete <name>`, `clone <dest>`. | `--name <name>`: Snapshot label<br>`--description <desc>`: Snapshot notes |
| **`vmx`** | `gitmap vm vmx <alias> <action>` | Configuration maintenance: `inspect` (key-value dump), `repair` (deduplicate keys, CRLF normalization), `optimize` (SSD memory paging flags). | None |
| **`scheduler shutdown`** | `gitmap vm scheduler shutdown` | Pauses job queue, broadcasts soft ACPI shutdown to all running VMs in parallel, waits, exits scheduler daemon. | `--all-vms`: Ensure all VMs stopped<br>`--force`: Force kill on timeout<br>`--timeout <secs>`: Drain wait (default: `120s`) |

---

## 5. Core Execution Workflows

### 5.1 VM Fleet Registration & Group Management
```bash
# Register a Windows 11 VM and assign to the 'dev' group
gitmap vm add win11-dev "vms/win11/win11.vmx" --group dev --description "Primary Windows Dev Box"

# Register an Ubuntu VM and immediately designate as active default
gitmap vm add ubuntu-runner "vms/ubuntu/ubuntu.vmx" --group ci-runners --set-default

# Create a new testing group and assign multiple VMs
gitmap vm group add testing --description "Automated test fleet"
gitmap vm group assign win11-dev testing
gitmap vm group assign ubuntu-runner testing

# Set the active default group
gitmap vm group set-default testing

# List inventory and groups
gitmap vm ls
gitmap vm group ls
```

### 5.2 Power Lifecycle Workflows
```bash
# Start all VMs in the active default group headlessly
gitmap vm start

# Start a specific VM with graphical console window
gitmap vm start win11-dev --gui

# Query live status with detailed disk and snapshot hierarchy
gitmap vm status win11-dev --detailed

# Soft graceful ACPI shutdown with 60s timeout
gitmap vm stop win11-dev --soft --timeout 60

# Suspend a running VM to disk
gitmap vm stop win11-dev --suspend

# Emergency hard stop across an entire group
gitmap vm stop @testing --hard --force
```

### 5.3 MAC Address Mutation Workflow
```bash
# 1. Verify VM is powered off
gitmap vm status win11-dev

# 2. Mutate to a specific static MAC address (creates automatic .vmx.bak)
gitmap vm mac win11-dev 00:50:56:38:57:7B --adapter ethernet0

# 3. Generate a dynamic randomized VMware OUI MAC address
gitmap vm mac win11-dev generate --adapter ethernet0

# 4. Verify updated MAC in inventory
gitmap vm ls
```

### 5.4 Disk Expansion & Repair Workflow
```bash
# 1. List attached virtual disks and bus channels
gitmap vm disk win11-dev list

# 2. Expand SCSI disk to 120GB (dispatches vmware-vdiskmanager -x 120GB)
gitmap vm disk win11-dev expand --disk scsi0:0 --size 120GB

# 3. Repair fragmented or inconsistent .vmdk file
gitmap vm disk win11-dev repair --disk scsi0:0
```

### 5.5 Snapshot Lifecycle & Tree Management
```bash
# 1. Take a clean pre-test snapshot
gitmap vm snapshot win11-dev create "pre-chrome-extension-test" --description "Clean state before test bundle"

# 2. List the snapshot hierarchy
gitmap vm snapshot win11-dev list

# 3. Revert to earlier snapshot after test run
gitmap vm snapshot win11-dev revert "pre-chrome-extension-test"

# 4. Delete temporary test snapshot to consolidate disks
gitmap vm snapshot win11-dev delete "pre-chrome-extension-test"
```

### 5.6 VMX Optimization & Repair
```bash
# Inspect raw key-value configuration
gitmap vm vmx win11-dev inspect

# Normalize line endings and strip duplicate entries
gitmap vm vmx win11-dev repair

# Optimize performance for SSD host (disables memory trim and pagefiles)
gitmap vm vmx win11-dev optimize
```

### 5.7 Coordinated Scheduler Shutdown
```bash
# Gracefully stop the task scheduler, drain queues, and shut down all VMs
gitmap vm scheduler shutdown --timeout 120

# Emergency shutdown with forced VM termination
gitmap vm scheduler shutdown --all-vms --force --timeout 30
```

---

## 6. Go Implementation Architecture & Split-DB Contracts

### 6.1 Database Schema (`cli/cmdvmware/store/schema.go`)
All VM metadata resides in GitMap's SQLite Split-DB tier (`installation.db` or `repodb`).

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
```

### 6.2 Domain Models (`cli/cmdvmware/types.go`)
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

### 6.3 Standard AppError Error Codes
All errors are wrapped using `apperror.NewWithDetails`:

| Code | Constant Identifier | Failure Condition |
|---|---|---|
| `E4001` | `ErrVmNotFound` | Specified VM alias does not exist in SQLite database. |
| `E4002` | `ErrVmxNotFound` | Target `.vmx` path does not exist on disk. |
| `E4003` | `ErrInvalidMac` | MAC address syntax fails validation regex (`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`). |
| `E4004` | `ErrVmRunning` | Mutation attempted while VM is currently powered on. |
| `E4005` | `ErrDriverNotFound` | `vmrun.exe` or `vmware-vdiskmanager.exe` missing from candidate paths. |
| `E4006` | `ErrPowerShellBridge` | PowerShell automation script (`scripts/vmware/manage-vm.ps1`) error. |
| `E4007` | `ErrGroupNotFound` | Group alias missing or no default group configured. |
| `E4008` | `ErrDiskOperation` | Virtual disk resizing or repair returned non-zero exit code. |
| `E4009` | `ErrSnapshotOperation` | Snapshot creation, revert, or deletion rejected by hypervisor. |
| `E4010` | `ErrSchedulerTimeout` | Coordinated scheduler shutdown exceeded allotted timeout window. |

---

## 7. Operational Invariants & Engineering Guardrails

1. **R1 Zero Full Builds or Test Suites (TOTAL BAN):** NEVER run `go build`, `npm run build`, or `go test ./...`. Run only targeted, file-scoped checks (`go test -run ^TestVM ./cli/cmdvmware/...`).
2. **R2 Micro-Batch File Cap ($\le 100$ lines):** Every Go source file must strictly stay within 100 lines. Function bodies must remain between 8 and 15 lines.
3. **R11 Strict Relative Paths & Lowercase Hygiene:** NEVER use absolute paths (drive letters, user home directories) or file scheme URIs in source code, comments, or documentation. Always reference files relative to the git repository root. All file and folder names must be strictly lowercase.
4. **R16 Diagnostic Log Isolation:** All hypervisor test execution transcripts, MAC mutation logs, and diagnostic traces must be stored exclusively in `repo-secrets/01-gitmap/72-chrome-ext-test-and-vmware-spec/` (or current task sequence directory under `repo-secrets/`). Never commit live machine credentials or network dumps to public repository paths.
5. **Affirmative Boolean Conventions:** Always name boolean fields and variables with affirmative prefixes (`isActive`, `isRunning`, `isDefault`, `hasTools`, `isForce`). Inverted booleans (`isNotRunning`, `disableGui`) are strictly prohibited.
6. **Offline State Guard for VMX Alterations:** Direct modification of `.vmx` files (MAC address changes, disk channel reconfiguration, memory tuning) must verify that the target VM is completely stopped (`is_running = 0` and absence of `.lck` directories) before proceeding.
7. **Automated Safety Backups:** Any operation that rewrites a `.vmx` file must create an atomic `.vmx.bak` backup file immediately prior to modification.
