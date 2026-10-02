# 72 — VMware CLI Command Suite, Local Workstation Integration & AI Prompt Architecture Spec

## Status: Pending
- **Spec ID:** 72-02
- **Document:** Component Specification
- **Scope:** VMware CLI Subsystem, Workstation Engine Integration, E2E Lifecycle Testing & AI Prompt Architecture
- **Related Spec:** `02-spec/21-app/72-chrome-ext-test-and-vmware-spec/01-architecture-spec.md`
- **Execution Plan:** `.ai-memory/plans/pending/72-chrome-ext-test-and-vmware-spec.md`

---

## User Request (Verbatim)

```text
# High Priority Instruction

Okay, I think this is a time that you need to test the Chrome extension. What do I mean by that? We have the Chrome import-export feature. Okay? So what you will do is you will try to export one profile, just take any one sample profile, take it as a test case, and if you need to write something secret, like the email or something secret, you use the repo secrets folder. Okay? Inside this, you create the git map folder with the sequence, and then inside this, you create the task name and with the sequence, by the way, and then you try to work on it. Okay? So you can export there and try to re-import by removing the profile. Try to see if that is successful. Pick one profile and start with that. If that is not successful, then pick another one, export, and every step, you try to have a log so that you understand by looking into the log what went wrong. Okay? So you try to go inside the log, so your logs will be saved inside the repo secrets folder by the task. Okay? And by looking into or inspecting the log or error tracing, you would know this is wrong. And also, I want you to first fix the git map be first. Okay? Try to make sure the build fixes. And then you make a minor bump first. Again, then you start with this testing. And then after the testing is done, you understand everything is working fine, then you make another minor bump and release. Okay, also, at the same time, what I do want from you is to understand, do we have the VMware commands? I think we should have. But I want you to confirm this. Yeah, I think I asked for the VMware commands, like if we can have the VMware. VMware commands will work like the SSH command. So it would have sub commands like add VMs. I should be able to add VM. I should be able to install VMware. Okay? Add VMs to this. Not all the VMs, but selected VMs to my, let's say, default group, and work with those VMs. Like, what is the status? How many snapshots are there? What is their network? Do you want to change their, let's say, MAC addresses? Do you want to expand the hard drive? Do you want to fix the VMX files? Okay? All kinds of things we want to have this. So at the end, I want you to create a detailed, let's say, AI instruction that will have those CLI commands for the VMware. Add, remove, add groups, default groups, LS, help, what each section would do, so that if any AI sees this instruction, they could actually implement all these things. And for the end-to-end testing, I want you to install the VMware inside this VM as well, and you try to create a new VM, add that. So we should have these type of functionalities, actually. Let me check. Yeah, we do have the VMware inside. I'm not sure if the latest one needs to be installed. Let me check. The CTU cycle needs to be changed. Okay. Install the VMware first so that you can do the end-to-end testing. You can do anything inside. No worries. Do not think of anything. Maybe I could revert back because this is a VM, and everything will be good as new. Okay. VMware is installed. There is no issue on this. Okay, so I will install or enable the hypervisor so that you can test this. Okay, and you can create a VM for now, just any VM you can create. Okay. It doesn't have to be the real one. The idea is we should be able to change the MAC. Yeah, we wanted to start the machine as well. We want to see if the machine is running. We also want to shut down the machine. Also, we want to shut down the scheduler with the machine as well. Okay? All kinds of things we want. For this, I want you to write the spec and also the AI instruction so that I share with the AI, they would understand everything. Is it clear? Do you understand? I have installed VM Ware latest everything is working fine.
```

---

## 1. System Architecture & Component Boundaries

The VMware subsystem in GitMap provides complete virtualization management across local workstations and remote hypervisors. It is architected around strict structural parity with the existing GitMap SSH fleet management subsystem (`cli/cmdssh`), offering intuitive commands, target grouping, default context selection, rich terminal visualization, and SQLite-backed inventory persistence.

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                                   GitMap CLI Dispatcher                                 │
│                                  (gitmap vm / gitmap vmware)                            │
└───────────────────────────────────────────┬─────────────────────────────────────────────┘
                                            │
               ┌────────────────────────────┼───────────────────────────┐
               ▼                            ▼                           ▼
    ┌──────────────────────┐    ┌──────────────────────┐    ┌───────────────────────┐
    │  Inventory & Groups  │    │ Lifecycle Execution  │    │ Hardware & VMX Engine │
    │  (cmdvmware/store)   │    │ (cmdvmware/lifecycle)│    │  (cmdvmware/config)   │
    ├──────────────────────┤    ├──────────────────────┤    ├───────────────────────┤
    │ - vm add / vm rm     │    │ - vm start / stop    │    │ - vm mac (set/gen)    │
    │ - vm group / default │    │ - vm status          │    │ - vm disk (expand)    │
    │ - vm ls              │    │ - scheduler shutdown │    │ - vm vmx / snapshot   │
    └──────────┬───────────┘    └──────────┬───────────┘    └───────────┬───────────┘
               │                           │                            │
               └───────────────────────────┼────────────────────────────┘
                                           ▼
                       ┌────────────────────────────────────────┐
                       │     Platform Driver / Adapter Layer    │
                       ├────────────────────┬───────────────────┤
                       │  vmrun.exe Driver  │  PowerShell Engine│
                       │ (Live WS Ops/VIX)  │ (manage-vm.ps1)   │
                       └─────────┬──────────┴─────────┬─────────┘
                                 │                    │
                                 ▼                    ▼
                       ┌────────────────────┬───────────────────┐
                       │ VMware Workstation │ Local .vmx / .vmdk│
                       │ Hypervisor Engine  │ Virtual Filesystem│
                       └────────────────────┴───────────────────┘
```

### 1.1 Structural Parity with SSH Fleet Architecture
In GitMap SSH fleet management:
- Nodes are registered with aliases (`gitmap ssh join add ...`).
- Nodes belong to logical groups and a default context can be selected (`gitmap ssh ...`).
- Inventory lists display liveness, hostnames, and operational states (`gitmap ssh ls`).
- Bulk actions can target all nodes, specific groups, or single aliases.

The VMware CLI (`gitmap vm` and `gitmap vmware`) mirrors this exact pattern:
- Virtual machines are registered by alias with paths to `.vmx` files (`gitmap vm add <alias> <path>`).
- Virtual machines are partitioned into groups (e.g. `default`, `dev`, `staging`, `testing`).
- A default group or default active VM can be pinned (`gitmap vm default <target>`).
- Operations without explicit targets operate automatically on the active default context.

---

## 2. CLI Command Specification (SSH Parity Architecture)

The CLI entrypoint supports both `gitmap vm` (concise primary command) and `gitmap vmware` (canonical backwards-compatible alias).

### 2.1 Inventory Management Commands

#### 2.1.1 `gitmap vm add`
Registers a virtual machine into GitMap's SQLite inventory store.
```bash
gitmap vm add <alias> <vmx-path> [flags]
```
- **Arguments:**
  - `alias`: Unique alphanumeric identifier (e.g. `win11-dev`, `ubuntu-test`).
  - `vmx-path`: Relative or absolute path to the `.vmx` configuration file.
- **Flags:**
  - `--group <name>`: Assign VM to an initial group (defaults to `default`).
  - `--description <text>`: Brief purpose of the VM.
  - `--set-default`: Set this VM as the active single default VM immediately.
- **Behavior:**
  - Validates that `vmx-path` exists and contains a valid `.vmx` header (`config.version`).
  - Inspects initial MAC address, RAM, vCPUs, and disk references.
  - Inserts record into `vm_instances` and adds relationship into `vm_group_memberships`.
  - Emits styled terminal success box with parsed properties.

#### 2.1.2 `gitmap vm rm`
Removes a virtual machine from GitMap inventory.
```bash
gitmap vm rm <alias> [flags]
```
- **Flags:**
  - `--force`: Remove without interactive confirmation.
  - `--delete-files`: Also remove underlying disk/vmx files from filesystem (requires explicit secondary prompt).
- **Behavior:**
  - Removes record from `vm_instances` and associated `vm_group_memberships`.
  - Preserves local `.vmx` files on disk unless `--delete-files` is explicitly asserted.

#### 2.1.3 `gitmap vm ls`
Lists registered virtual machines in a formatted terminal table.
```bash
gitmap vm ls [flags]
```
- **Flags:**
  - `--group <name>`: Filter by specific group (e.g. `--group dev`).
  - `--all`: Show all VMs regardless of active group filter.
  - `--json`: Output full raw inventory in machine-readable JSON format.
- **Columns:**
  - `ALIAS`: Unique name.
  - `STATUS`: Live power state (`running`, `stopped`, `suspended`, `unknown`).
  - `GROUP`: Primary group membership.
  - `MAC`: Primary network adapter MAC address (`ethernet0`).
  - `MEMORY`: Configured RAM (e.g. `8192 MB`).
  - `VCPUS`: Configured CPU cores.
  - `VMX PATH`: Truncated or relative configuration path.
  - `DEFAULT`: Indicator badge `[DEFAULT]` if marked active.

#### 2.1.4 `gitmap vm group`
Manages logical groupings of virtual machines.
```bash
gitmap vm group <subcommand> [args]
```
- **Subcommands:**
  - `gitmap vm group add <group-name>`: Creates a new VM grouping.
  - `gitmap vm group rm <group-name>`: Deletes a VM group (does not delete member VMs).
  - `gitmap vm group ls`: Lists all groups and member count.
  - `gitmap vm group assign <alias> <group-name>`: Adds a VM to a group.
  - `gitmap vm group remove <alias> <group-name>`: Removes a VM from a group.
  - `gitmap vm group set-default <group-name>`: Pins the default group for unattended operations.

#### 2.1.5 `gitmap vm default`
Sets or views the current default target context.
```bash
gitmap vm default [target]
```
- If invoked without arguments: Displays currently active default group and default VM alias.
- If invoked with an argument: Sets either the default VM alias or default group name.

---

### 2.2 Lifecycle & Hypervisor Execution Commands

#### 2.2.1 `gitmap vm install`
Installs or validates the VMware virtualization stack and supporting utilities.
```bash
gitmap vm install [flags]
```
- **Windows Actions:**
  - Probes for VMware Workstation Pro / Player installation directory.
  - Locates `vmrun.exe` in `PATH` or standard installation candidates (`C:\Program Files (x86)\VMware\VMware Workstation\vmrun.exe`).
  - Validates `VMTools` Windows service status.
  - If missing and running with administrator elevation: offers automated silent installer invocation via winget (`winget install VMware.WorkstationPro`) or direct MSI recipe.
- **Linux Actions:**
  - Installs `open-vm-tools` and `open-vm-tools-desktop` via `apt-get` / `dnf`.
  - Verifies `/mnt/hgfs` kernel module loading.

#### 2.2.2 `gitmap vm status`
Inspects real-time operational status of one or all VMs.
```bash
gitmap vm status [alias] [flags]
```
- **Flags:**
  - `--json`: Output status payload as JSON.
  - `--detailed`: Include full disk geometry, snapshot tree, and network topology.
- **Live Discovery:**
  - Queries `vmrun.exe -T ws list` to detect if the target `.vmx` is currently running in the hypervisor process table.
  - If running: probes guest IP address via VMware tools query or ARP table lookup.
  - If stopped: checks file lock state (`.lck` directories) to verify clean dismount.

#### 2.2.3 `gitmap vm start`
Powers on a virtual machine.
```bash
gitmap vm start [alias] [flags]
```
- **Flags:**
  - `--nogui` / `--headless`: Start without VMware Workstation GUI window (via `vmrun.exe start <path> nogui`).
  - `--gui`: Launch with full graphical workstation console.
  - `--all`: Start all VMs in the active default group sequentially or concurrently.
  - `--timeout <seconds>`: Maximum duration to wait for hypervisor confirmation (default: `60s`).

#### 2.2.4 `gitmap vm stop`
Shuts down or powers off a virtual machine.
```bash
gitmap vm stop [alias] [flags]
```
- **Flags:**
  - `--soft`: Initiate graceful OS shutdown via VMware Tools ACPI guest signal (`vmrun.exe stop <path> soft`).
  - `--hard` / `--force`: Immediate hard power-off (`vmrun.exe stop <path> hard`).
  - `--suspend`: Suspend execution and save state to disk (`.vmss`).
  - `--all`: Stop all VMs in the active group.
  - `--timeout <seconds>`: Timeout for soft shutdown before prompting or enforcing hard stop (default: `90s`).

#### 2.2.5 `gitmap vm scheduler shutdown`
Coordinates automated hypervisor power-down with background task scheduler termination.
```bash
gitmap vm scheduler shutdown [flags]
```
- **Flags:**
  - `--all-vms`: Gracefully stop all currently running VMs before scheduler shutdown.
  - `--force`: Force terminate running tasks and hypervisor processes.
  - `--timeout <seconds>`: Maximum wait window for inflight tasks to drain (default: `120s`).
- **Execution Flow:**
  1. Signals the GitMap task queue worker pool to reject incoming jobs and initiate graceful drain.
  2. Queries all running VMs registered in GitMap inventory.
  3. Sends graceful soft shutdown signals to all running VMs in parallel.
  4. Waits up to `--timeout` for all VM processes to terminate cleanly.
  5. Terminates the background scheduler process and logs audit record to SQLite database.

---

### 2.3 Hardware & Configuration Manipulation Commands

#### 2.3.1 `gitmap vm mac`
Queries or updates network adapter MAC addresses.
```bash
gitmap vm mac <alias> [mac-address|generate] [flags]
```
- **Arguments:**
  - `mac-address`: Formatted MAC string (e.g. `00:50:56:38:57:7B`) or keyword `generate`.
- **Flags:**
  - `--adapter <name>`: Target adapter identifier (defaults to `ethernet0`).
  - `--type <generated|static>`: Force MAC addressing type (`generated` vs `static`).
- **Behavior:**
  - If VM is running: Refuses modification unless `--force` is passed (modifying running `.vmx` causes corruption).
  - Modifies `.vmx` properties:
    - For static: `ethernet0.addressType = "static"`, `ethernet0.address = "<mac>"`.
    - For generated: `ethernet0.addressType = "generated"`, `ethernet0.generatedAddress = "<mac>"`.
  - Creates timestamped backup `.vmx.bak` prior to modification.
  - Re-reads and verifies written values.

#### 2.3.2 `gitmap vm disk`
Inspects, expands, or repairs virtual machine disks (`.vmdk`).
```bash
gitmap vm disk <alias> <subcommand> [flags]
```
- **Subcommands:**
  - `list`: Displays all attached virtual disks, controller channels (SCSI/NVMe/SATA), capacity, and sparse allocation.
  - `expand`: Expands virtual disk size (e.g. `--disk scsi0:0 --size 120GB`). Calls `vmware-vdiskmanager.exe -x <size> <path>`.
  - `repair`: Checks and repairs damaged `.vmdk` files. Calls `vmware-vdiskmanager.exe -R <path>`.

#### 2.3.3 `gitmap vm snapshot`
Manages virtual machine snapshot tree and state persistence.
```bash
gitmap vm snapshot <alias> <subcommand> [flags]
```
- **Subcommands:**
  - `list`: Displays snapshot tree hierarchy with timestamps and descriptions.
  - `create <name>`: Creates a snapshot (`vmrun.exe -T ws snapshot <path> <name>`).
  - `revert <name>`: Restores VM state to a specific named snapshot (`vmrun.exe -T ws revertToSnapshot <path> <name>`).
  - `delete <name>`: Merges and removes a snapshot (`vmrun.exe -T ws deleteSnapshot <path> <name>`).
  - `clone <dest-folder>`: Clones snapshot or VM to a new destination folder (`vmrun.exe -T ws clone <path> <dest> full`).

#### 2.3.4 `gitmap vm vmx`
Inspects and repairs `.vmx` configuration file integrity.
```bash
gitmap vm vmx <alias> [inspect|repair|optimize] [flags]
```
- **Subcommands:**
  - `inspect`: Displays key-value entries with syntax highlighting and validation warnings.
  - `repair`: Rebuilds broken line endings, removes invalid duplicate keys, fixes relative paths to `.vmdk` / `.nvram`.
  - `optimize`: Disables aggressive memory trimming and page sharing for faster guest SSD execution.

---

## 3. Local Workstation Integration Engine

### 3.1 Driver Resolution Strategy
The GitMap Go core (`cli/cmdvmware`) interfaces with VMware Workstation through a dual-driver model:
1. **Direct Native `vmrun.exe` Execution:**
   - Used for high-speed, binary lifecycle commands (`start`, `stop`, `reset`, `suspend`, `snapshot`, `list`).
   - Command dispatch formats:
     ```text
     vmrun.exe -T ws start "<vmx-path>" nogui
     vmrun.exe -T ws stop "<vmx-path>" soft
     vmrun.exe -T ws list
     vmrun.exe -T ws snapshot "<vmx-path>" "<snap-name>"
     ```
2. **PowerShell Automation Engine (`scripts/vmware/manage-vm.ps1`):**
   - Used for complex offline `.vmx` parsing, regex configuration manipulation, multi-adapter MAC generation, and hardware resizing.
   - Invocation wrapper:
     ```powershell
     powershell.exe -NoProfile -ExecutionPolicy Bypass -File "scripts/vmware/manage-vm.ps1" -Action Set-MAC -VMPath "<vmx-path>" -MACAddress "<mac>" -Adapter ethernet0 -Json
     ```
   - JSON output mode enables structured Go unmarshaling directly into domain models.

### 3.2 Candidate Path Discovery Hierarchy
When invoking `vmrun.exe` or `vmware-vdiskmanager.exe`, GitMap searches candidate paths in priority order:
1. System `PATH` (`exec.LookPath("vmrun.exe")`).
2. Standard 64-bit Directory: `C:\Program Files\VMware\VMware Workstation\vmrun.exe`.
3. Standard 32-bit Compatibility: `C:\Program Files (x86)\VMware\VMware Workstation\vmrun.exe`.
4. VMware Player: `C:\Program Files (x86)\VMware\VMware Player\vmrun.exe`.
5. VMware VIX SDK: `C:\Program Files (x86)\VMware\VMware VIX\vmrun.exe`.

---

## 4. End-to-End VM Lifecycle Testing Specification

### 4.1 End-to-End Test VM Provisioning
For safe, deterministic end-to-end testing without external ISO downloads or heavy disk overhead:
- A lightweight testing VM template is initialized under a temporary test directory:
  - Configuration file: `test-e2e.vmx`
  - Minimal virtual disk stub: `test-e2e.vmdk`
  - Minimum virtual hardware: 1 vCPU, 512 MB RAM, 1 network adapter (`ethernet0`).
- The test harness supports testing against live existing VMs if available, or isolated dummy `.vmx` definitions for configuration mutators.

### 4.2 Test Execution Protocol Matrix

| Phase | Operation | Command / Invocation | Expected Outcome | Verification Metric |
|---|---|---|---|---|
| **Phase 1** | Register VM | `gitmap vm add test-vm ./test-e2e.vmx --group testing` | Success exit code 0 | Present in `vm_instances` SQLite table |
| **Phase 2** | Set Default Group | `gitmap vm default testing` | Context pinned | Default group set to `testing` |
| **Phase 3** | Set Static MAC | `gitmap vm mac test-vm 00:50:56:11:22:33` | `.vmx` updated with static MAC | Backup `.bak` created; file contains MAC |
| **Phase 4** | Generate Dynamic MAC | `gitmap vm mac test-vm generate` | `.vmx` updated with generated MAC | Valid OUI prefix (`00:0c:29` or `00:50:56`) |
| **Phase 5** | Power On VM | `gitmap vm start test-vm --nogui` | Hypervisor process spawned | Listed in `vmrun.exe list` |
| **Phase 6** | Status Inspection | `gitmap vm status test-vm --json` | JSON output with `isRunning: true` | Valid status payload parsed |
| **Phase 7** | Soft Graceful Stop | `gitmap vm stop test-vm --soft` | Guest shutdown signal dispatched | Removed from `vmrun.exe list` |
| **Phase 8** | Scheduler Shutdown | `gitmap vm scheduler shutdown --timeout 30` | Scheduler process and all VMs stopped | Audit entry logged; zero dangling locks |

### 4.3 Repo Secrets Logging Isolation
All sensitive diagnostic output, VM configuration dumps, MAC addresses, and execution transcripts MUST be saved to the sequenced repository secrets folder:
`repo-secrets/01-gitmap/72-chrome-ext-test-and-vmware-spec/`
- `01-chrome-export-test.log`
- `02-chrome-import-test.log`
- `03-vmware-e2e-lifecycle.log`
- `04-scheduler-shutdown.log`

---

## 5. AI Prompt Instruction Specification Structure

The companion AI instruction prompt file `01-prompts/42-vmware-cli-commands.md` must be constructed as a self-contained, high-authority engineering guide that any AI agent can ingest to implement the complete VMware CLI command suite.

### 5.1 Required Prompt Sections & Blueprints

1. **Title & Frontmatter:**
   - Canonical metadata, prompt ID (`42-vmware-cli-commands`), version, target CLI binary (`gitmap`).
2. **Role & Operational Objective:**
   - Defines the agent as an Enterprise Systems Engineer implementing high-performance CLI commands in Go.
3. **Non-Negotiable Architectural Rules:**
   - File size limits ($\le 100$ lines per file).
   - Function length limits ($8-15$ lines per function).
   - Nesting limits ($\le 1$ level of indentation).
   - Positive boolean naming (`isActive`, `isDefault`, `isEnabled`, `isRunning`, `hasTools`).
   - AppError / Result monadic error envelope usage (`apperror.NewWithDetails`).
4. **Complete CLI Grammar & Subcommand Reference:**
   - Exhaustive markdown table detailing every command, short alias, arguments, flags, and descriptions.
5. **PowerShell Engine & `vmrun.exe` Integration Specs:**
   - Precise invocation signatures and JSON bridge parsing routines.
6. **SQLite Split-DB Schema & Model Contracts:**
   - Data structures for `vm_instances`, `vm_groups`, `vm_group_memberships`.
7. **Subagent Orchestration & Micro-Batch Execution Plan:**
   - Step-by-step phased rollout with bounded 5-8 file micro-batches.
8. **Verification Gates & Quality Checks:**
   - Build, lint, and mock unit tests ensuring zero regressions.

---

## 6. Data Contracts & Database Schema

The VMware subsystem stores VM inventory and group relationships in GitMap's SQLite Split-DB (`gitmap.db` or `installation.db`).

### 6.1 Database Tables & DDL

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

### 6.2 Go Domain Structs

```go
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
```

---

## 7. Quality Gates & Acceptance Criteria

- **QG-1:** Specification contains full lossless verbatim user prompt.
- **QG-2:** All boolean fields and models follow positive naming conventions (`isActive`, `isDefault`, `isEnabled`, `isRunning`, `hasTools`).
- **QG-3:** Complete CLI parity with SSH commands demonstrated across commands, subcommands, flags, and terminal outputs.
- **QG-4:** Integration with `scripts/vmware/manage-vm.ps1` and `vmrun.exe` clearly specified with candidate path discovery hierarchy.
- **QG-5:** Zero absolute filesystem paths inside specification content; all repository paths strictly lowercase and relative.
- **QG-6:** End-to-end VM lifecycle and scheduler shutdown procedures documented with concrete verification metrics.
