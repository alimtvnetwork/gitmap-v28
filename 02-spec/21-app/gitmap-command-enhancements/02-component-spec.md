# Component Specification: Network Adapter Metrics, Multi-Node Task Queue Delegation & Release Ceremony

**Document ID:** `02-spec/21-app/gitmap-command-enhancements/02-component-spec.md`  
**Classification:** Core Application Component Specification (`02-spec/21-app`)  
**Version:** 1.0.0  
**Updated:** 2026-10-16  
**Milestone:** `gitmap-command-enhancements` (Run 83)  
**Status:** `ratified`  
**AI Confidence:** High  
**Ambiguity:** [04-role-based-task-scheduling.md](../../../.ai-memory/ambiguous-questions/01-new-ambiguity/04-role-based-task-scheduling.md)  
**Companion Documents:**  
- [Architecture Specification: GitMap Command Enhancements](01-architecture-spec.md)  
- [Specification Catalog Readme](readme.md)  
- [Subtask 03: OS Network Adapter Metrics](../../../.ai-memory/plans/subtasks/gitmap-command-enhancements/03-os-network-adapter-metrics.md)  
- [Subtask 04: Multi-Node Task Queue & Delegation](../../../.ai-memory/plans/subtasks/gitmap-command-enhancements/04-nodes-tasks-queue-and-delegation.md)  
- [Subtask 05: Documentation, Help Docs & UI Registry](../../../.ai-memory/plans/subtasks/gitmap-command-enhancements/05-docs-help-and-ui-registry.md)  
- [Subtask 06: Verification & Release Ceremony](../../../.ai-memory/plans/subtasks/gitmap-command-enhancements/06-verification-and-release-ceremony.md)  

---

## Keywords

`network-adapters` · `ip-address` · `mac-address` · `link-speed` · `bytes-sent` · `bytes-received` · `nodes-tasks` · `task-delegation` · `agent-task-db` · `heavy-duty-tasks` · `collect-unified-nodes` · `split-db` · `helpdoc` · `commands-ts` · `semver-minor-bump` · `release-ceremony`

---

## Scoring & Compliance Matrix

| Criterion | Status | Notes |
| :--- | :---: | :--- |
| `readme.md` registered | ✅ | Tracked in `02-spec/21-app/gitmap-command-enhancements/readme.md` |
| AI Confidence assigned | ✅ | High |
| Ambiguity assigned | ✅ | Tracked in `.ai-memory/ambiguous-questions/01-new-ambiguity/04-role-based-task-scheduling.md` |
| Keywords present | ✅ | Indexed above |
| Scoring table present | ✅ | Validated |
| Strict relative paths | ✅ | Zero absolute paths or `file:///` URIs; strictly relative Git paths |
| Positive booleans only | ✅ | `isUp`, `hasCompleted`, `isActive`, `isHeavyDuty`, `isOpenOnly`, `isJSON` |
| Error wrapping protocol | ✅ | Standardized on `*apperror.AppError` and `*appfault.AppError` |
| Zero-build / Zero-test | ✅ | Strictly maintained during specification authoring |

---

## 1. Executive Summary & Component Scope

This component specification details the engineering contracts, data schemas, CLI routing mechanisms, and distributed execution workflows for Wave 2 of the `gitmap-command-enhancements` milestone. It defines the implementation blueprint for:

1. **Network Adapter Interrogation Engine (`gitmap os network` / `gitmap network`):**
   - High-speed, non-intrusive query of physical and virtual network adapters across Windows, Linux, and macOS.
   - Cross-platform extraction of interface names, status (`isUp`), MAC addresses, IPv4/IPv6 addresses, link speeds, and cumulative traffic counters (`bytesSent`, `bytesReceived`).
   - Unified Lipgloss table rendering and machine-readable JSON output (`--json`).

2. **Distributed Multi-Node Task Queue & Delegation Subsystem (`gitmap nodes tasks`):**
   - Fleet-wide task visibility aggregating local and remote AI agent tasks from `agent-task.db`.
   - Identification of active, pending, and heavy-duty tasks across cluster nodes.
   - Surgical resolution of the node enumeration bug in `cli/cmdnodes/nodes_resources.go:67` by switching to `collectUnifiedNodes(ctx, true)`.
   - Recording the future role-based task delegation architecture in the project ambiguity ledger.

3. **Documentation, Modern Help UI & Web Studio Catalog:**
   - Authoring terminal help topics in `cli/helpdoc/` (`cpu.md`, `ram.md`, `os.md`, `nodes.md`).
   - Registering newly introduced telemetry and fleet commands in `src/data/commands.ts`.

4. **Autonomous Minor Version Bump Release Ceremony (`v6.554.0`):**
   - Standardized release protocol via `03-ai-scripts/37-bump-version.py`.
   - Cross-manifest synchronization across `version.json`, `package.json`, `readme.md`, and `changelog.md`.

---

## 2. Network Adapter Telemetry Component

### 2.1 Component Overview

Developers and AI agents frequently need real-time network interface metrics to verify node reachability, diagnose network congestion, validate IP configuration, and monitor bandwidth consumption without launching heavy GUI tools.

The Network Adapter Telemetry component introduces `gitmap os network` (with top-level alias `gitmap network`).

```bash
# Human-readable Lipgloss terminal table
gitmap os network
gitmap network

# Machine-readable JSON output
gitmap os network --json
gitmap network --json
```

### 2.2 Domain Models & Data Structures

The network data models are located in `cli/cmdos/os_network_types.go`:

```go
package cmdos

// NetworkAdapterMetrics represents normalized interface telemetry across operating systems.
type NetworkAdapterMetrics struct {
	Name                 string   `json:"name"`
	InterfaceAlias       string   `json:"interfaceAlias"`
	Status               string   `json:"status"`
	IsUp                 bool     `json:"isUp"`
	MACAddress           string   `json:"macAddress"`
	IPv4Addresses        []string `json:"ipv4Addresses"`
	IPv6Addresses        []string `json:"ipv6Addresses"`
	LinkSpeedBitsPerSec  uint64   `json:"linkSpeedBitsPerSec"`
	LinkSpeedDisplay     string   `json:"linkSpeedDisplay"`
	BytesSent            uint64   `json:"bytesSent"`
	BytesReceived        uint64   `json:"bytesReceived"`
	PacketsSent          uint64   `json:"packetsSent"`
	PacketsReceived      uint64   `json:"packetsReceived"`
}

// NetworkTelemetryPayload encapsulates the full collection response.
type NetworkTelemetryPayload struct {
	MachineName string                  `json:"machineName"`
	Timestamp   string                  `json:"timestamp"`
	Adapters    []NetworkAdapterMetrics `json:"adapters"`
}
```

### 2.3 Cross-Platform Interrogation Strategy

#### 2.3.1 Windows Platform (`cli/cmdos/os_network_windows.go`)

Windows interrogation leverages PowerShell CIM cmdlets executed with `-NoProfile` and formatted as JSON:

```powershell
powershell -NoProfile -Command "
$adapters = Get-NetAdapter | Select-Object Name, InterfaceDescription, Status, MacAddress, LinkSpeed;
$ips = Get-NetIPAddress | Select-Object InterfaceAlias, IPAddress, AddressFamily;
$stats = Get-NetAdapterStatistics | Select-Object Name, SentBytes, ReceivedBytes, SentPackets, ReceivedPackets;
@{ Adapters = $adapters; IPs = $ips; Stats = $stats } | ConvertTo-Json -Depth 3 -Compress
"
```

The Go collector decodes the compressed JSON stream into Windows-specific intermediate structs, merges IP addresses and byte statistics by interface name/alias, normalizes boolean `isUp` (`Status == "Up"`), and builds `NetworkAdapterMetrics`.

#### 2.3.2 Linux & macOS Platform (`cli/cmdos/os_network_linux.go`)

On Linux environments:
1. Interface enumeration uses Go's native `net.Interfaces()`, extracting hardware MAC addresses and IP addresses (`net.Interface.Addrs()`).
2. Status flags evaluate `intf.Flags&net.FlagUp != 0` to set `IsUp: true`.
3. Link speed is queried directly from `/sys/class/net/<name>/speed` (expressed in Mbps).
4. Cumulative bytes and packet statistics are parsed efficiently from `/proc/net/dev`:
   ```text
   Inter-|   Receive                                                |  Transmit
    face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
     eth0: 10485760   12450    0    0    0     0          0         0  5242880    8900    0    0    0     0       0          0
   ```

On macOS / Darwin:
1. Go's native `net.Interfaces()` provides base enumeration.
2. Interface throughput counters fall back to parsing `netstat -ibn` or default zero counters gracefully without erroring.

### 2.4 Terminal Rendering & Formatting

The Lipgloss renderer formats network telemetry into a structured terminal table:

```text
┌────────────────────────────────────────────────────────────────────────────────────────┐
│ NETWORK ADAPTERS                                                      Host: dev-win-01 │
├─────────────────┬──────────┬───────────────────┬─────────────────────┬─────────────────┤
│ Adapter         │ Status   │ MAC Address       │ IPv4 / IPv6         │ Speed & Traffic │
├─────────────────┼──────────┼───────────────────┼─────────────────────┼─────────────────┤
│ Ethernet 1      │ UP       │ 00:15:5D:8A:2B:01 │ 192.168.1.150       │ 1.0 Gbps        │
│                 │          │                   │ fe80::492a:...      │ ▲ 42.1 MB       │
│                 │          │                   │                     │ ▼ 184.8 MB      │
├─────────────────┼──────────┼───────────────────┼─────────────────────┼─────────────────┤
│ vEthernet (WSL) │ UP       │ 00:15:5D:F3:11:00 │ 172.28.16.1         │ 10.0 Gbps       │
│                 │          │                   │                     │ ▲ 1.2 MB        │
│                 │          │                   │                     │ ▼ 4.5 MB        │
└─────────────────┴──────────┴───────────────────┴─────────────────────┴─────────────────┘
```

- Status `UP` rendered in bold green (`#50fa7b`).
- Status `DOWN` rendered in dim gray (`#6272a4`).
- Traffic formatted with human-readable binary prefixes (`KB`, `MB`, `GB`).

### 2.5 CLI Wiring

1. **`cli/cmdos/os.go`:** Added subcommand dispatch case:
   ```go
   case "network", "net", "adapter", "adapters", "interface", "interfaces":
       return RunNetworkCLI(subArgs)
   ```
2. **`cli/cmd/roottooling.go`:** Added root-level command mapping:
   ```go
   {[]string{"network", "net", "adapters", "net-adapters", "netadapters"}, func() error { return cmdos.RunNetworkCLI(argsTail()) }},
   ```

---

## 3. Multi-Node Task Queue & Fleet Delegation Component

### 3.1 Component Overview

In distributed multi-agent clusters, developers and orchestrator AI agents require immediate visibility into which tasks are active, enqueued, or pending across all fleet machines.

The fleet command `gitmap nodes tasks` (with aliases `gitmap nodes queue`, `gitmap nodes jobs`, `gitmap nodes task`) interrogates all cluster nodes in parallel, reads task state from the run-scoped Tier 2 SQLite task databases (`agent-task.db`), and aggregates workload distribution.

```bash
# Interrogate task queues across all fleet nodes
gitmap nodes tasks

# Interrogate specific node task queue
gitmap nodes tasks vm-02

# Structured JSON output
gitmap nodes tasks --json
```

### 3.2 Bug Resolution in `nodes_resources.go:67`

#### 3.2.1 Root Cause Analysis

In `cli/cmdnodes/nodes_resources.go:67`, the existing implementation called:
```go
conns, _ := cmdssh.FetchAllSSHConnections()
records := executeFleetResourcesProbe(conns, opts)
```

**Architectural Defect:**
1. `cmdssh.FetchAllSSHConnections()` only scans remote SSH connections from legacy configuration files.
2. It completely omits the **local machine** node from the resource audit.
3. It bypasses the unified cluster store (`store.OpenDefault()`), ignoring nodes discovered via `ListClusterNodes()` and `ListHosts()`.
4. The error return was discarded (`_`), silencing database connection failures.

#### 3.2.2 Ratified Architectural Fix

Both `RunNodesResources` and `RunNodesTasks` must utilize `collectUnifiedNodes(ctx, true)` defined in `cli/cmdnodes/nodes_collect.go`:

```go
// Correct invocation in cli/cmdnodes/nodes_resources.go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

unifiedNodes, err := collectUnifiedNodes(ctx, true)
if err != nil {
    return apperror.WrapSimple(err, "RunNodesResources.collectUnifiedNodes")
}
```

`collectUnifiedNodes(ctx, true)` correctly:
- Opens `store.OpenDefault()`.
- Unifies SSH hosts, configured SSH connections, and cluster nodes.
- Preserves local node representation.
- Provides standard `UnifiedFleetNode` structs with IP, port, alias, and credentials.

### 3.3 Task Database Interrogation (`agent-task.db`)

Each AI agent execution runs within an isolated workspace run directory:
`.ai-memory/temp-agents/<slug>/agent-task.db`

The task database schema includes:
```sql
CREATE TABLE IF NOT EXISTS ParentTask (
    ParentTaskId INTEGER PRIMARY KEY AUTOINCREMENT,
    TaskName TEXT NOT NULL,
    TaskSlug TEXT NOT NULL,
    RunDirectory TEXT NOT NULL,
    Status TEXT NOT NULL DEFAULT 'ACTIVE',
    IsActive INTEGER NOT NULL DEFAULT 1,
    HasCompleted INTEGER NOT NULL DEFAULT 0,
    TotalStepsBudget INTEGER NOT NULL DEFAULT 300,
    CurrentStep INTEGER NOT NULL DEFAULT 1,
    Notes TEXT NULL,
    CreatedAt TEXT NOT NULL,
    UpdatedAt TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS Subtask (
    SubtaskId INTEGER PRIMARY KEY AUTOINCREMENT,
    ParentTaskId INTEGER NOT NULL,
    TaskCode TEXT NOT NULL,
    TaskSlug TEXT NOT NULL DEFAULT '',
    Title TEXT NOT NULL,
    AssignedAgentRole TEXT NULL,
    OwnedFilesJson TEXT NOT NULL DEFAULT '[]',
    Status TEXT NOT NULL DEFAULT 'PENDING',
    IsBlocked INTEGER NOT NULL DEFAULT 0,
    HasCompleted INTEGER NOT NULL DEFAULT 0,
    Evidence TEXT NULL,
    CreatedAt TEXT NOT NULL,
    UpdatedAt TEXT NOT NULL
);
```

#### 3.3.1 Heavy-Duty Task Classification

A task is identified as **Heavy-Duty** (`isHeavyDuty: true`) if:
1. `ParentTask.TotalStepsBudget >= 100`, OR
2. Task slug or notes contain heavy-load keywords (`heavy`, `full-audit`, `deep-scan`, `migration`, `build-all`), OR
3. Number of owned files in `OwnedFilesJson` exceeds 8.

Heavy-duty tasks are highlighted in the UI and marked for potential node-to-node work stealing and load rebalancing.

### 3.4 Data Models for Multi-Node Tasks

Located in `cli/cmdnodes/nodes_tasks.go`:

```go
package cmdnodes

// TaskItemSummary represents a single active or pending subtask.
type TaskItemSummary struct {
	TaskID          int64  `json:"taskId"`
	TaskSlug        string `json:"taskSlug"`
	TaskCode        string `json:"taskCode"`
	Title           string `json:"title"`
	AssignedRole    string `json:"assignedRole"`
	Status          string `json:"status"`
	IsHeavyDuty     bool   `json:"isHeavyDuty"`
	CurrentStep     int    `json:"currentStep"`
	TotalSteps      int    `json:"totalSteps"`
	CreatedAt       string `json:"createdAt"`
}

// NodeTaskReport represents aggregated task queue state for an individual fleet node.
type NodeTaskReport struct {
	NodeAlias       string            `json:"nodeAlias"`
	NodeHost        string            `json:"nodeHost"`
	NodePort        int               `json:"nodePort"`
	IsLocal         bool              `json:"isLocal"`
	Status          string            `json:"status"`
	ActiveTaskCount int               `json:"activeTaskCount"`
	PendingTaskCount int              `json:"pendingTaskCount"`
	HeavyTaskCount  int               `json:"heavyTaskCount"`
	ActiveTasks     []TaskItemSummary `json:"activeTasks"`
	PendingTasks    []TaskItemSummary `json:"pendingTasks"`
	ErrorMessage    string            `json:"errorMessage,omitempty"`
}

// FleetTasksPayload represents the fleet-wide task queue response.
type FleetTasksPayload struct {
	Timestamp      string           `json:"timestamp"`
	TotalNodes     int              `json:"totalNodes"`
	TotalActive    int              `json:"totalActive"`
	TotalPending   int              `json:"totalPending"`
	TotalHeavyDuty int              `json:"totalHeavyDuty"`
	NodeReports    []NodeTaskReport `json:"nodeReports"`
}
```

### 3.5 Remote SSH Execution Protocol

For remote nodes:
1. The collector invokes the remote GitMap binary via SSH:
   `gitmap tasks --json --limit 10`
2. If `gitmap tasks` is unavailable on the remote host, it executes a fallback SQLite query against `.ai-memory/temp-agents/*/agent-task.db`:
   ```bash
   sqlite3 .ai-memory/temp-agents/*/agent-task.db "SELECT SubtaskId, TaskCode, Title, Status, AssignedAgentRole FROM Subtask WHERE HasCompleted = 0 LIMIT 10;"
   ```
3. Queries execute concurrently using a bounded worker pool (`maxConcurrency = 8`) with a strict per-node timeout of 5 seconds to prevent hanging SSH channels.

### 3.6 Terminal UI Table Rendering

```text
┌─────────────────────────────────────────────────────────────────────────────────────────────────┐
│ FLEET TASK QUEUES & DELEGATION                                                Nodes Active: 2/3 │
├────────────┬─────────┬────────┬─────────┬───────┬───────────────────────────────────────────────┤
│ Node Alias │ Host    │ Active │ Pending │ Heavy │ Running Task / Subtask                        │
├────────────┼─────────┼────────┼─────────┼───────┼───────────────────────────────────────────────┤
│ local (me) │ 127.0.0.1│ 1      │ 3       │ 1     │ Task-03-Net (OS Network Adapter Telemetry)    │
│ vm-02      │ 192.168.1│ 1      │ 1       │ 0     │ Task-02-CIM (Fast Hardware Telemetry)         │
│ cluster-03 │ 192.168.1│ 0      │ 0       │ 0     │ [IDLE] No tasks enqueued                      │
└────────────┴─────────┴────────┴─────────┴───────┴───────────────────────────────────────────────┘
```

---

## 4. Documentation & Web UI Registration Component

### 4.1 Terminal Markdown Help Topics (`cli/helpdoc/`)

Help topics in `cli/helpdoc/` are rendered dynamically in the terminal and embedded into GitMap's help system.

1. **`cli/helpdoc/os.md`:**
   - Update subcommand table to include `network` / `net`.
   - Add usage examples for `gitmap os network` and `gitmap os network --json`.
2. **`cli/helpdoc/cpu.md`:**
   - Author standalone help topic for `gitmap cpu` / `gitmap processor`.
   - Document 5-second cross-process telemetry caching, CIM LoadPercentage on Windows, and `/proc/stat` on Linux.
3. **`cli/helpdoc/ram.md`:**
   - Author standalone help topic for `gitmap ram`.
   - Document visible physical memory, used memory, free memory, and percentage utilization.
4. **`cli/helpdoc/nodes.md`:**
   - Document `gitmap nodes tasks` and `gitmap nodes resources`.
   - Detail task queue aggregation, heavy-duty task flags, and remote execution options.

### 4.2 Web UI Command Registry (`src/data/commands.ts`)

`src/data/commands.ts` drives the GitMap documentation web portal, command palette, and interactive explorer. Four new command definitions are registered:

1. **`gitmap os network` / `gitmap network`:**
   - Category: `remote` ("Machine & Network Identity")
   - Flags: `--json` (`-j`), `--all` (`-a`)
   - Examples: `gitmap network`, `gitmap os network --json`
2. **`gitmap cpu`:**
   - Category: `remote`
   - Alias: `gitmap processor`
   - Flags: `--json`, `--no-cache`
   - Examples: `gitmap cpu`, `gitmap processor --json`
3. **`gitmap ram`:**
   - Category: `remote`
   - Flags: `--json`, `--no-cache`
   - Examples: `gitmap ram`, `gitmap ram --json`
4. **`gitmap nodes tasks`:**
   - Category: `cluster` ("Cluster & Delegation")
   - Alias: `gitmap nodes queue`
   - Flags: `--json`, `--open-only`, `--heavy-only`
   - Examples: `gitmap nodes tasks`, `gitmap nodes tasks --heavy-only`

---

## 5. Minor Version Bump Release Ceremony (`v6.554.0`)

### 5.1 Release Management Protocol

Following completion of all implementation subtasks (Subtasks 01-05), Subtask 06 executes the canonical minor version bump release ceremony:

1. **Pre-flight Cleanliness Verification:**
   - Confirm all generated files, documentation, and source changes adhere to zero-absolute-path rules.
   - Confirm positive boolean field naming across all Go structs and TypeScript interfaces.
2. **Execute Autonomous Version Bumper:**
   ```bash
   python 03-ai-scripts/37-bump-version.py --tier minor --scope "Hardware telemetry, network adapter metrics, and fleet task queue delegation"
   ```
3. **Automated Manifest Synchronization:**
   - `version.json`: Bump version string from `6.553.0` to `6.554.0`.
   - `package.json`: Update `"version": "6.554.0"`.
   - `readme.md`: Update version badges and header metadata.
   - `changelog.md`: Append `## [6.554.0] - 2026-10-16` release block cataloging hardware metrics, network adapter telemetry, and fleet task queue commands.
4. **Git Finalization:**
   - Stage all modified and added files using relative paths.
   - Formulate atomic commit message conforming to GitMap release conventions.
   - Verify zero uncommitted artifacts remain.

---

## 6. Architecture & Implementation Invariants

### 6.1 Positive Boolean Convention

In strict adherence to repository coding guidelines (`02-spec/02-coding-guidelines/`), all boolean fields and variables MUST use positive prefixes and semantics:

| Forbidden Pattern | Required Replacement | Context |
| :--- | :--- | :--- |
| `isDown`, `disabled`, `isNotUp` | `isUp` | `NetworkAdapterMetrics` |
| `isNotBlocked` | `isBlocked` (or `canProceed`) | `Subtask` |
| `uncompleted`, `isPending` | `hasCompleted` | Task lifecycle |
| `if isSuccess == true` | `if isSuccess` | Conditionals |
| `if !isSuccess` | `if result.IsFail()` | Result wrappers |
| `noCache` | `isBypassCache` / `isLive` | Telemetry flags |

### 6.2 Error Wrapping Protocol

All returned errors across `cmdos` and `cmdnodes` MUST be wrapped using `*apperror.AppError` or `*appfault.AppError`:

```go
if err != nil {
    return apperror.WrapSimple(err, "dispatchNetworkQuery")
}
```

Raw Go `error` returns from top-level CLI entrypoints must adapt `*apperror.AppError` cleanly:

```go
func RunNetworkCLI(args []string) error {
    if err := RunNetwork(args); err != nil {
        return err
    }
    return nil
}
```

### 6.3 Path Hygiene Invariant

Every file path generated in code, tests, documentation, or release notes MUST be strictly relative to the repository root (e.g., `cli/cmdos/os_network_types.go`, `02-spec/21-app/gitmap-command-enhancements/readme.md`). Absolute filesystem paths (e.g., `D:\work\...`, `C:\Users\...`) and `file:///` URIs are strictly prohibited.
