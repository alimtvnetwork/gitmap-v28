# Fleet Nodes Running Only (`gitmap nodes-running-only` / `gitmap nro`)

> **Spec Directory:** `02-spec/21-app/gitmap-nodes-running-only/`  
> **Status:** Production-Ready Specification  
> **Master Plan:** `.ai-memory/plans/gitmap-nodes-running-only.md`  
> **Created:** 2026-10-16  

---

## 1. Overview & Problem Statement

Running `gitmap nodes --json` queries and outputs all registered fleet machines across the cluster. In heterogeneous distributed environments, many registered worker nodes may be powered off, sleeping, or disconnected from the local network. This produces two significant friction points for users and automation scripts:

1. **Output Clutter**: The resulting JSON includes offline nodes with status `"○ offline (timeout)"` and `"is_online": false`, requiring extra filtering passes in user scripts (`jq 'map(select(.is_online == true))'`).
2. **Latency Penalty on Fleet Delegation**: Running batch commands (e.g. `gitmap nodes tasks`, `gitmap nodes resources`, `gitmap ssh exec`) attempts network connections to offline machines, repeatedly suffering 1200ms TCP probe timeouts and cluttering terminal output with warning notices.

This specification package establishes `gitmap nodes-running-only` and its short alias `gitmap nro`, providing:
- A dedicated command returning strictly active, online fleet nodes (`is_online: true`).
- A 10-second cross-process node liveness cache to eliminate repeated TCP connection delays.
- A suite of running-only delegation commands (`nro exec`, `nro tasks`, `nro resources`, `nro ping`, `nro pull`, `nro cp all`, `nro scan`) executing exclusively across reachable machines.

---

## 2. Specification Documents Inventory

| Document | Path | Description |
|---|---|---|
| **Architecture Specification** | `02-spec/21-app/gitmap-nodes-running-only/01-architecture-spec.md` | System topology, architectural goals, command syntax, flag matrix, and high-level routing. |
| **Component Specification** | `02-spec/21-app/gitmap-nodes-running-only/02-component-spec.md` | Component models, `NodesLivenessCache` TTL engine, `RunningOnlySSHExecutor`, data contracts, and error boundaries. |
| **Directory Guide** | `02-spec/21-app/gitmap-nodes-running-only/readme.md` | Package overview, file catalog, command reference, and workflows. |

---

## 3. Subtask Decomposition

The implementation is partitioned across 6 modular subtasks documented in `.ai-memory/plans/subtasks/gitmap-nodes-running-only/`:

| Subtask ID | Document Path | Primary Scope & Owned Files |
|---|---|---|
| **Subtask 01** | `.ai-memory/plans/subtasks/gitmap-nodes-running-only/01-nodes-liveness-cache-and-probe.md` | Liveness caching engine and probe optimization (`cli/cmdnodes/nodes_liveness_cache.go`, `cli/cmdnodes/nodes_probe.go`). |
| **Subtask 02** | `.ai-memory/plans/subtasks/gitmap-nodes-running-only/02-nodes-running-only-dispatcher-and-cli.md` | Top-level dispatcher and alias routing (`cli/cmdnodes/nodes_running_only.go`, `cli/cmd/rootcore.go`). |
| **Subtask 03** | `.ai-memory/plans/subtasks/gitmap-nodes-running-only/03-running-only-delegation-and-sshexec.md` | Streamlined SSH execution without offline delays (`cli/cmdssh/sshexec.go`, `cli/cmdnodes/nodes_filter.go`). |
| **Subtask 04** | `.ai-memory/plans/subtasks/gitmap-nodes-running-only/04-nodes-tasks-and-resources-open-only.md` | Filtering offline nodes in tasks and resources (`cli/cmdnodes/nodes_resources.go`, `cli/cmdnodes/nodes_tasks.go`, `cli/cmdnodes/nodes_tasks_types.go`). |
| **Subtask 05** | `.ai-memory/plans/subtasks/gitmap-nodes-running-only/05-docs-help-completion-and-ui-registry.md` | Markdown help docs, tab completion, constants, and React web catalog (`cli/helpdoc/nodes-running-only.md`, `cli/completion/completion.go`, `src/data/commands.ts`). |
| **Subtask 06** | `.ai-memory/plans/subtasks/gitmap-nodes-running-only/06-verification-release-bump-and-cicd.md` | Linter gates, minor version bump to `v6.555.0`, atomic commit, and CI/CD verification (`version.json`, `changelog.md`). |

---

## 4. Command Reference & Syntax

### 4.1 Primary Commands & Aliases
```bash
# Full command name
gitmap nodes-running-only [flags]

# Canonical short alias
gitmap nro [flags]

# Hyphenated aliases
gitmap nodesrunningonly [flags]
gitmap nodes-running [flags]
gitmap nodes-online [flags]
gitmap nodes-open [flags]
gitmap nodes-active [flags]
```

### 4.2 Supported Flags
- `--json`, `-j`: Emit clean JSON array containing only nodes where `is_online == true`.
- `--refresh`, `-r`: Invalidate liveness cache and trigger immediate probe.
- `--target <alias|ip>`, `-t`: Filter execution to specific node if online.
- `--timeout <duration>`: Custom timeout for probe or command execution.
- `-h`, `--help`: Display formatted terminal help documentation.

### 4.3 Running-Only Delegation Commands
```bash
# Execute remote command across running nodes only (no offline probe delay)
gitmap nro exec "git status"
gitmap nro run "uptime"

# Inspect active task queues on running nodes only
gitmap nro tasks [--json]

# Probe CPU and RAM telemetry on running nodes only
gitmap nro resources [--json]

# Ping active nodes with zero timeout lag
gitmap nro ping

# Run batch pull across online nodes
gitmap nro pull

# Commit and push across dirty repos on online nodes
gitmap nro cp all "chore: sync cluster nodes"

# Trigger repository discovery scan across online nodes
gitmap nro scan
```

---

## 5. User Workflows

### Workflow 1: Machine-Readable Scripting with `--json`
```bash
# Query active worker IPs in bash / pwsh scripts
$onlineNodes = gitmap nro --json | ConvertFrom-Json
$onlineNodes | ForEach-Object { Write-Host "Active: $($_.alias) at $($_.host)" }
```

### Workflow 2: Rapid Ad-Hoc Cluster Command Execution
```bash
# Run uptime across active nodes without waiting for powered-off laptops
gitmap nro exec "uptime"
```

### Workflow 3: Non-Blocking Telemetry Inspection
```bash
# View resources only on responsive machines
gitmap nro resources
```

---

## 6. Architecture & Coding Invariants

1. **Strictly Relative Paths**: All path references must use forward slashes and begin from the repository root (e.g. `cli/cmdnodes/...`). Absolute filesystem paths or `file:///` URIs are forbidden.
2. **Positive Booleans**: Use positive prefix identifiers (`isOnline`, `hasValidCache`, `canExecute`). Prohibit negative boolean comparisons (`!isOffline`).
3. **Structured Errors**: Go error boundaries wrap errors using `*appfault.AppError` or `apperror.WrapSimple`.
4. **Zero-Build & Zero-Test During Spec Authoring**: Document specifications and subtasks thoroughly without intermediate local compilation or test runs.
