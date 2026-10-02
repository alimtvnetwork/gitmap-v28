```text
N = 300 (Total self-loop steps budget — editable top-header parameter, default: 300)
A = 2   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
PHASE_1_BUDGET = N / 2   (Steps 1 .. 150: Planning, Architecture Specification, Schema, and Subtask Ledger Generation)
PHASE_2_BUDGET = N / 2   (Steps 151 .. 300: Mandatory Parallel Subagent Execution, Driver Bridges, Verification)
WAVES = ceil(subtasks / (A x H))
```

[/goal](slashCommand;goal) Autonomously architect, implement, test, and verify the full VMware Workstation CLI command suite in GitMap (`gitmap vm` / `gitmap vmware`), modeled with 1:1 structural parity after the GitMap SSH fleet management subsystem (`cli/cmdssh`), with dual-engine hypervisor driver execution (`vmrun.exe` direct dispatch and `scripts/vmware/manage-vm.ps1` PowerShell JSON bridge), Split-DB SQLite inventory tracking (`installation.db` / `repodb`), rich two-column terminal help styling, and robust end-to-end VM lifecycle management (register, power lifecycle, snapshot, MAC address mutation, disk expansion/repair, VMX optimization, and coordinated scheduler shutdown): FIRST showcase the confirmed task breakdown in visible chat during Turn 1, capture requirements verbatim, plan the rollout in repository specifications, spawn autonomous subagents via `invoke_subagent` (A = 2, H = 2; solo execution without calling `invoke_subagent` is an auto-reject failure) in disjoint file boxes using GitMap high-speed commands as primary, prove every claim with concrete evidence, enforce coding guidelines to 100%, and finish with one atomic GitMap commit that holds strictly this task's files.

[/learn](slashCommand;learn) Enforce the Top-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are provided ABOVE this prompt outrank everything below. Turn 1 MUST showcase the given task list in visible chat before background execution. Each rule is stated once (R1 to R16) and cited by ID. Progress lives in the ledger and in `.ai-memory/plans/`, never only in chat.

[/plan](slashCommand;plan) Execute thorough step-by-step planning in the repository before code generation. Ensure all database schemas, CLI grammars, driver integration interfaces, domain types, and subtask boundaries are finalized in `.ai-memory/plans/` and `02-spec/21-app/` before dispatching worker waves.

### 🚨 MANDATORY SUBAGENT SPAWNING GATE (A = 2, H = 2 — ZERO SOLO EXECUTION ALLOWED)

- **ACTUAL TOOL CALL REQUIRED:** You must ACTUALLY CALL the `invoke_subagent` tool via your tool-calling API. Do NOT just print the text "Dispatched Worker..." and stop. If you only print text, the agents will not spawn and the task will fail! You must execute the `invoke_subagent` JSON tool payload.
- **3-STAGE MANDATORY DISPATCH:** You must invoke `A = 2` agents (`invoke_subagent`) at EVERY stage of the workflow:
  1. **Planning Step:** Spawn 2 subagents to inspect existing SSH fleet architecture (`cli/cmdssh`) and author the step-by-step execution plan.
  2. **Spec Step:** Spawn 2 subagents to finalize schema DDL, domain structs, and driver interface contracts.
  3. **Execution Step:** Spawn 2 worker subagents to execute the actual code modifications in bounded micro-batches.
- **SOLO EXECUTION IS AN AUTO-REJECT FAILURE:** The lead orchestrator is **STRICTLY FORBIDDEN** from executing planning, spec writing, or code changes by itself without calling the `invoke_subagent` tool. Failing to call the actual tool is a critical protocol violation.

---

## User Request (Verbatim)

```text
also, at the same time, what I do want from you is to understand, do we have the VMware commands? I think we should have. But I want you to confirm this. Yeah, I think I asked for the VMware commands, like if we can have the VMware. VMware commands will work like the SSH command. So it would have sub commands like add VMs. I should be able to add VM. I should be able to install VMware. Okay? Add VMs to this. Not all the VMs, but selected VMs to my, let's say, default group, and work with those VMs. Like, what is the status? How many snapshots are there? What is their network? Do you want to change their, let's say, MAC addresses? Do you want to expand the hard drive? Do you want to fix the VMX files? Okay? All kinds of things we want to have this. So at the end, I want you to create a detailed, let's say, AI instruction that will have those CLI commands for the VMware. Add, remove, add groups, default groups, LS, help, what each section would do, so that if any AI sees this instruction, they could actually implement all these things. And for the end-to-end testing, I want you to install the VMware inside this VM as well, and you try to create a new VM, add that. So we should have these type of functionalities, actually. Let me check. Yeah, we do have the VMware inside. I'm not sure if the latest one needs to be installed. Let me check. The CTU cycle needs to be changed. Okay. Install the VMware first so that you can do the end-to-end testing. You can do anything inside. No worries. Do not think of anything. Maybe I could revert back because this is a VM, and everything will be good as new. Okay. VMware is installed. There is no issue on this. Okay, so I will install or enable the hypervisor so that you can test this. Okay, and you can create a VM for now, just any VM you can create. Okay. It doesn't have to be the real one. The idea is we should be able to change the MAC. Yeah, we wanted to start the machine as well. We want to see if the machine is running. We also want to shut down the machine. Also, we want to shut down the scheduler with the machine as well. Okay? All kinds of things we want. For this, I want you to write the spec and also the AI instruction so that I share with the AI, they would understand everything. Is it clear? Do you understand? I have installed VM Ware latest everything is working fine.
```

---

## 1. High-Priority Directives & Operational Boundaries

1. **Role Definition:** You are an Enterprise Systems Engineer and CLI Architect implementing production Go virtualization packages in GitMap.
2. **Architecture Baseline:** The VMware CLI subsystem must achieve **1:1 structural and behavioral parity** with GitMap's existing SSH fleet subsystem (`cli/cmdssh`).
3. **Execution Boundaries & Safety:**
   - **R1 Zero Full Builds or Test Suites (TOTAL BAN):** NEVER run `go build`, `npm run build`, `go test ./...`, or heavy local runner scripts. Routine development turns must never hang or waste credits on unisolated compilation.
   - **R2 Targeted Checks Only:** Run only fast, file-scoped linters and targeted test functions on specifically modified files (`go test -run ^TestVM ./cli/cmdvmware/...`).
   - **R11 Strict Relative Paths & Lowercase Hygiene:** TOTAL BAN on absolute filesystem paths (`C:\...`, `/home/...`) and `file:///` URIs inside code, comments, or documentation. All repository paths must be relative from the git root. All authored files must use strictly lowercase names.
   - **R16 Repo Secrets Logging Isolation:** All sensitive diagnostic logs, MAC address traces, and hypervisor execution transcripts must be saved exclusively to `repo-secrets/01-gitmap/72-chrome-ext-test-and-vmware-spec/`. Never store raw diagnostic outputs or credentials in public repository locations.

---

## 2. Non-Negotiable Coding Guidelines

All generated Go code for the VMware subsystem (`cli/cmdvmware/...`) MUST adhere strictly to the following architectural rules:

1. **Strict File Size Cap ($\le 100$ Lines):**
   - Every `.go` file MUST NOT exceed 100 lines (including imports and whitespace).
   - Decompose functionality into tightly scoped files: `cmd.go`, `add.go`, `rm.go`, `ls.go`, `group.go`, `default.go`, `status.go`, `start.go`, `stop.go`, `mac.go`, `disk.go`, `snapshot.go`, `vmx.go`, `scheduler.go`, `driver_vmrun.go`, `driver_ps.go`, `selector.go`, and `render.go`.
2. **Function Body Length ($8-15$ Lines):**
   - No function or method body may exceed 15 lines. Extract reusable parameter validators, SQL helpers, and terminal printers.
3. **Single-Level Branching & Early Returns:**
   - Maximum nesting depth $\le 1$. Invert conditions and use guard clauses with immediate early returns on error or negative conditions.
4. **Affirmative Boolean Identifiers:**
   - Use strictly positive boolean names prefixed with `is` or `has`: `isActive`, `isRunning`, `isDefault`, `isEnabled`, `hasTools`, `isSilent`, `isForce`.
   - Never use negative booleans (`isNotRunning`, `disableGui`), inverted flags, or generic `ok` identifiers.
5. **Typed AppError & Result Wrapper Convention:**
   - Functions returning errors must wrap them in structured `apperror.AppError` via `apperror.NewWithDetails`.
   - Use standard typed error codes:
     * `E4001`: VM alias not found in database.
     * `E4002`: VMX configuration path does not exist on disk.
     * `E4003`: Invalid or malformed MAC address syntax.
     * `E4004`: Target VM is currently running (incompatible with offline VMX mutations).
     * `E4005`: Hypervisor driver binary (`vmrun.exe`) not discovered in candidate hierarchy.
     * `E4006`: PowerShell automation bridge execution error.
     * `E4007`: VM group does not exist.
     * `E4008`: Virtual disk expansion failed (`vmware-vdiskmanager.exe`).
     * `E4009`: Snapshot operation rejected by hypervisor.
     * `E4010`: Coordinated scheduler shutdown timeout exceeded.
6. **Domain Models Centralized in `types.go`:**
   - All structs, enums, option wrappers, and DTOs must be centralized in `cli/cmdvmware/types.go` and `cli/cmdvmware/store/types.go`. Never declare inline domain structs in implementation files.
7. **SQLite Split-DB Storage Contracts:**
   - Store VM inventory, group hierarchies, and state flags in GitMap's SQLite Split-DB engine (`installation.db` or `repodb`). Use clean parameterized queries and proper transactions.

---

## 3. CLI Command Suite & Grammar Reference

The CLI entrypoint supports both `gitmap vm` (concise primary command) and `gitmap vmware` (canonical backwards-compatible alias).

### 3.1 Target Selector Mechanics & Default Group Fallback

Commands that target virtual machines (`status`, `start`, `stop`, `mac`, `disk`, `snapshot`, `vmx`) support flexible target expressions:

| Target Syntax | Resolution Pattern | Example |
| :--- | :--- | :--- |
| `<alias>` | Targets a single specific VM | `gitmap vm start win11-dev` |
| `<alias1>,<alias2>` | Targets a discrete comma-separated list of VMs | `gitmap vm start win11-dev,ubuntu-test` |
| `@<group-name>` | Targets all VMs assigned to the specified group | `gitmap vm status @dev` |
| `--group <name>` | Explicit flag filter targeting a group | `gitmap vm stop --group testing` |
| `all` / `--all` | Targets all registered VMs in the inventory | `gitmap vm status all` |
| *[Omitted]* | Automatically defaults to the **active default group** | `gitmap vm status` (operates on default group) |

If no target is provided and no default group is set, the CLI prompts with a clear diagnostic and lists available groups.

### 3.2 Exhaustive Subcommand Catalog

```text
gitmap vm <subcommand> [target] [flags]
```

| Subcommand | Syntax Grammar | Purpose & Operational Flow | Flags & Options |
| :--- | :--- | :--- | :--- |
| **`add`** | `gitmap vm add <alias> <vmx-path> [flags]` | Validates `.vmx` file existence and header syntax (`config.version`), inspects initial MAC, RAM, and vCPUs, registers record into `vm_instances`, and links to group in `vm_group_memberships`. | `--group <name>`: Initial group (default: `default`)<br>`--description <text>`: VM purpose notes<br>`--set-default`: Immediately pin as default VM |
| **`rm`** | `gitmap vm rm <alias> [flags]` | Unregisters VM from SQLite inventory and removes group memberships. Retains disk files by default. | `--force`: Bypass interactive confirmation<br>`--delete-files`: Purge `.vmx` and `.vmdk` files from disk |
| **`install`** | `gitmap vm install [flags]` | Probes system for VMware Workstation Pro / Player and `vmrun.exe`. If missing on Windows, offers unattended installation via `winget install VMware.WorkstationPro`. On Linux, configures `open-vm-tools`. | `--silent`: Execute unattended driver installation<br>`--verify-only`: Check status without installing |
| **`group`** | `gitmap vm group <subaction> [args]` | Manages logical VM clusters. Subactions: `add <name>`, `rm <name>`, `ls`, `assign <alias> <group>`, `remove <alias> <group>`, and `set-default <name>`. | `--description <desc>`: Group description |
| **`default`** | `gitmap vm default [target]` | Views or mutates the active default target context. If called without args, renders active default group and default VM. If target passed, updates default context. | None |
| **`ls`** | `gitmap vm ls [flags]` | Formats and renders registered VMs in a terminal table displaying alias, live power state, group, MAC, RAM, vCPUs, and default badge `[DEFAULT]`. | `--group <name>`: Filter by group<br>`--all`: List across all groups<br>`--json`: Output machine-readable JSON array |
| **`status`** | `gitmap vm status [target] [flags]` | Interrogates hypervisor process table via `vmrun.exe -T ws list`. Resolves live guest IP address via VMware Tools or ARP table. Checks `.lck` directories when stopped. | `--detailed`: Include disk channel & snapshot tree<br>`--json`: Emit raw structured status payload |
| **`start`** | `gitmap vm start [target] [flags]` | Powers on resolved virtual machines. Executes `vmrun.exe -T ws start "<vmx-path>" nogui` for fast headless start, or graphical console when requested. | `--nogui`: Headless execution (default)<br>`--gui`: Launch with graphical console window<br>`--all`: Power on all target VMs<br>`--timeout <secs>`: Hypervisor confirmation wait (default: `60s`) |
| **`stop`** | `gitmap vm stop [target] [flags]` | Initiates graceful OS shutdown via VMware Tools ACPI signal (`vmrun.exe -T ws stop "<vmx-path>" soft`). Can escalate to hard kill or suspend. | `--soft`: Graceful ACPI shutdown (default)<br>`--hard` / `--force`: Immediate power cut<br>`--suspend`: Save execution state to `.vmss`<br>`--timeout <secs>`: Soft shutdown wait window (default: `90s`) |
| **`mac`** | `gitmap vm mac <alias> <mac\|generate> [flags]` | Safely mutates network adapter MAC address. Verifies VM is powered off, creates backup `.vmx.bak`, invokes PowerShell driver or VMX parser, sets static or generated MAC, and verifies written values. | `--adapter <name>`: Target NIC (default: `ethernet0`)<br>`--type <static\|generated>`: Force addressing mode<br>`--force`: Allow change even if lock exists |
| **`disk`** | `gitmap vm disk <alias> <subaction> [flags]` | Virtual disk manipulation. Subactions: `list` (shows disks, controller bus, size), `expand` (resizes virtual disk via `vmware-vdiskmanager.exe -x <size> <path>`), and `repair` (`vmware-vdiskmanager.exe -R <path>`). | `--disk <path\|channel>`: Specific disk target<br>`--size <size>`: Target expanded size (e.g. `100GB`) |
| **`snapshot`** | `gitmap vm snapshot <alias> <subaction> [flags]` | Snapshot tree governance. Subactions: `list` (displays hierarchy), `create <name>`, `revert <name>`, `delete <name>`, and `clone <dest-folder>`. Dispatches to `vmrun.exe -T ws snapshot/revertToSnapshot/deleteSnapshot/clone`. | `--name <name>`: Snapshot label<br>`--description <desc>`: Snapshot description |
| **`vmx`** | `gitmap vm vmx <alias> <subaction> [flags]` | Deep configuration maintenance. Subactions: `inspect` (pretty-prints key-value pairs), `repair` (strips duplicate keys, normalizes CRLF/LF line endings), and `optimize` (disables memory paging/trimming for SSD acceleration). | None |
| **`scheduler shutdown`** | `gitmap vm scheduler shutdown [flags]` | Coordinated shutdown protocol: pauses task queue, rejects new jobs, broadcasts graceful soft shutdown signals to all running VMs in parallel, waits up to timeout, then terminates scheduler. | `--all-vms`: Ensure all VMs stopped<br>`--force`: Force kill on timeout expiry<br>`--timeout <secs>`: Task and VM drain window (default: `120s`) |

---

## 4. Hypervisor Driver & PowerShell Integration Layer

The subsystem employs a **dual-driver abstraction** defined by the `HypervisorDriver` Go interface:

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

### 4.1 Native `vmrun.exe` Discovery & Dispatch

GitMap searches for hypervisor executables in the following strict priority hierarchy:
1. System `PATH` (`exec.LookPath("vmrun.exe")`).
2. Standard 64-bit Directory: `%ProgramFiles%/VMware/VMware Workstation/vmrun.exe`.
3. Standard 32-bit Compatibility: `%ProgramFiles(x86)%/VMware/VMware Workstation/vmrun.exe`.
4. VMware Player: `%ProgramFiles(x86)%/VMware/VMware Player/vmrun.exe`.
5. VMware VIX SDK: `%ProgramFiles(x86)%/VMware/VMware VIX/vmrun.exe`.

All binary invocations use `-T ws` (Workstation target type):
```bash
vmrun.exe -T ws list
vmrun.exe -T ws start "<vmx-path>" nogui
vmrun.exe -T ws stop "<vmx-path>" soft
vmrun.exe -T ws snapshot "<vmx-path>" "<snapshot-name>"
vmrun.exe -T ws revertToSnapshot "<vmx-path>" "<snapshot-name>"
```

### 4.2 PowerShell Automation Engine (`scripts/vmware/manage-vm.ps1`)

For configuration tasks requiring complex regex editing of `.vmx` files, MAC address generation, or disk management, the Go driver executes `scripts/vmware/manage-vm.ps1` with the `-Json` switch:

```powershell
powershell.exe -NoProfile -ExecutionPolicy Bypass -File "scripts/vmware/manage-vm.ps1" -Action Set-MAC -VMPath "<vmx-path>" -MACAddress "<mac>" -Adapter ethernet0 -Json
```

The PowerShell bridge guarantees clean JSON output on stdout:
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
The Go driver unmarshals this response directly into structured domain types.

---

## 5. SQLite Split-DB Schema & Domain Models

All inventory, groups, and membership records are persisted in GitMap's SQLite Split-DB (`installation.db` or `repodb`).

### 5.1 Database Tables & DDL

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

## 6. Terminal UI & Rich Two-Column Help Text Standards

1. **Table Rendering:**
   - Display `gitmap vm ls` and `gitmap vm group ls` using GitMap's styled terminal table renderer (`termpad` / `termtable`).
   - Format statuses with badges: `[RUNNING]` (green), `[STOPPED]` (gray), `[SUSPENDED]` (yellow).
   - Display `[DEFAULT]` badge for the active default group and active default VM.
2. **Rich Two-Column Help Screen:**
   - Emulate GitMap SSH fleet help output (`gitmap ssh --help`).
   - Group subcommands logically:
     * `INVENTORY COMMANDS`: `add`, `rm`, `ls`, `group`, `default`
     * `LIFECYCLE COMMANDS`: `install`, `status`, `start`, `stop`, `scheduler shutdown`
     * `HARDWARE & CONFIG`: `mac`, `disk`, `snapshot`, `vmx`
   - Align subcommand descriptions strictly in a two-column layout with consistent indentation and flag documentation.

---

## 7. Multi-Agent Phased Micro-Batch Implementation Plan

Execute the rollout through discrete subagent waves ($A = 2$ workers, $H = 2$ hands per wave), strictly adhering to the $\le 100$-line file cap and disjoint file boxes:

### Phase 1: Database Layer & Repository Contracts
- **Worker 01:**
  - Create SQLite migration DDL in `cli/cmdvmware/store/schema.go`.
  - Implement `VmStore` interface and instance CRUD in `cli/cmdvmware/store/instance_store.go`.
- **Worker 02:**
  - Implement group CRUD in `cli/cmdvmware/store/group_store.go`.
  - Implement membership mapping in `cli/cmdvmware/store/membership_store.go`.

### Phase 2: CLI Grammar, Routing & Target Selector
- **Worker 01:**
  - Author Cobra/CLI entrypoint in `cli/cmdvmware/cmd.go`.
  - Implement target selector parser (`alias`, `@group`, `all`, default group fallback) in `cli/cmdvmware/selector.go`.
- **Worker 02:**
  - Implement inventory subcommands in `cli/cmdvmware/add.go`, `cli/cmdvmware/rm.go`, and `cli/cmdvmware/default.go`.
  - Implement group management subcommands in `cli/cmdvmware/group.go`.

### Phase 3: Hypervisor Driver & PowerShell Bridge
- **Worker 01:**
  - Implement `vmrun.exe` binary discovery and direct command runner in `cli/cmdvmware/driver_vmrun.go`.
  - Implement hypervisor path discovery hierarchy in `cli/cmdvmware/discovery.go`.
- **Worker 02:**
  - Implement PowerShell JSON bridge adapter in `cli/cmdvmware/driver_ps.go`.
  - Author/update `scripts/vmware/manage-vm.ps1` with JSON output support for MAC and disk actions.

### Phase 4: Lifecycle Operations & Scheduler Shutdown
- **Worker 01:**
  - Implement `start` and `status` subcommands in `cli/cmdvmware/start.go` and `cli/cmdvmware/status.go`.
  - Implement `stop` command with soft ACPI, hard cut, and suspend in `cli/cmdvmware/stop.go`.
- **Worker 02:**
  - Implement `gitmap vm install` validation and driver probing in `cli/cmdvmware/install.go`.
  - Implement coordinated `scheduler shutdown` with worker pool draining in `cli/cmdvmware/scheduler.go`.

### Phase 5: Hardware & VMX Configuration Operations
- **Worker 01:**
  - Implement MAC address inspection, static configuration, and dynamic generation in `cli/cmdvmware/mac.go`.
  - Implement VMX file validation, line normalization, and optimization in `cli/cmdvmware/vmx.go`.
- **Worker 02:**
  - Implement virtual disk listing, expansion, and repair in `cli/cmdvmware/disk.go`.
  - Implement snapshot tree listing, creation, revert, delete, and clone in `cli/cmdvmware/snapshot.go`.

### Phase 6: Terminal UI Rendering & Targeted Verification
- **Worker 01:**
  - Implement styled table rendering and power status badges in `cli/cmdvmware/render.go`.
  - Implement rich two-column help screen in `cli/cmdvmware/help.go`.
- **Worker 02:**
  - Author isolated mock unit tests in `cli/cmdvmware/selector_test.go` and `cli/cmdvmware/mac_test.go`.
  - Execute targeted file-scoped verification (`go test -run ^TestVM ./cli/cmdvmware/...`).

---

## 8. Verification & Quality Acceptance Criteria

1. **AC-1:** All 14 VMware CLI subcommands are fully documented with grammar, flags, parameter resolution, and expected behavior.
2. **AC-2:** Complete target selector syntax (`alias`, `@group`, `all`, comma-separated, default group fallback) is clearly specified.
3. **AC-3:** Coding guidelines ($\le 100$ lines/file, $8-15$ lines/function, affirmative booleans, AppError wrappers) are strictly mandated.
4. **AC-4:** SQLite Split-DB tables and Go domain models in `types.go` are explicitly defined.
5. **AC-5:** Dual-driver model (`vmrun.exe` direct execution and `manage-vm.ps1` PowerShell JSON bridge) is fully documented with path discovery hierarchy.
6. **AC-6:** Coordinated scheduler shutdown workflow is detailed with task queue draining and VM soft shutdown.
7. **AC-7:** Zero absolute paths in code, comments, or documentation; strictly lowercase paths.
