# Specification 185: Unified Fleet Nodes Command (`gitmap nodes`)

## 1. Executive Summary & Architectural Mission

GitMap environments operate across three distinct infrastructure management tiers:
1. **SSH Connection Vault & Remote Host Inventory**: Managed via `gitmap ssh nodes`, `SSHConnection` records, and `ssh_hosts` tables in `installation.db`.
2. **Distributed Cluster Fleet Database**: Managed via `gitmap cluster nodes`, tracking control plane and worker topology in `ClusterNode` SQLite tables.
3. **Server-Clients Orchestration (`sc`)**: Autonomous peer-to-peer broadcast and RPC network topology managing live node command fan-out.

Previously, operators had to run separate commands (`gitmap ssh nodes`, `gitmap cluster nodes`, `gitmap sc status`) to ascertain fleet membership and status.

This specification introduces **`gitmap nodes`** (aliases: `gitmap node`, `allnodes`, `all-nodes`, `fleet-nodes`), providing a unified, deduplicated aggregation plane across all three subsystems.

---

## 2. CLI Command: `gitmap nodes`

### 2.1 Invocation Syntax

- `gitmap nodes [flags] [target]`
- `gitmap nodes help`

### 2.2 Flags & Options

| Flag | Short | Purpose |
| :--- | :--- | :--- |
| `--fast` | `-f`, `--no-probe` | Skips live network TCP/SSH liveness probing for instantaneous cached output. |
| `--json` | `-j` | Emits structured JSON array of nodes for automation and scripting pipelines. |
| `--ssh` | | Filters results to only include nodes registered in the SSH connection vault. |
| `--cluster` | `--clst` | Filters results to only include nodes participating in the Cluster Fleet DB. |
| `--sc` | `--server-clients` | Filters results to only include nodes participating in the Server-Clients network. |
| `[target]` | | Positional substring filter matching against Node Alias or Host IP. |

### 2.3 Aliases
- `node`
- `allnodes`
- `all-nodes`
- `fleet-nodes`

---

## 3. Architecture & Data Flow

```
                      ┌────────────────────────────┐
                      │    gitmap nodes CLI        │
                      └─────────────┬──────────────┘
                                    │
           ┌────────────────────────┼────────────────────────┐
           ▼                        ▼                        ▼
  ┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐
  │ SSHConnection / │      │  ClusterNode    │      │ Server-Clients  │
  │   ssh_hosts     │      │     Table       │      │   Topology      │
  │ (SSH Subsystem) │      │ (Cluster Subsys)│      │  (SC Subsys)    │
  └────────┬────────┘      └────────┬────────┘      └────────┬────────┘
           │                        │                        │
           └────────────────────────┼────────────────────────┘
                                    ▼
                      ┌────────────────────────────┐
                      │ Unified Deduplication &    │
                      │ Subsystem Normalization    │
                      └─────────────┬──────────────┘
                                    │
                        [--fast flag passed?]
                                ├── Yes ──────────────┐
                                └── No ──┐            │
                                         ▼            │
                      ┌────────────────────────────┐  │
                      │ Concurrent Liveness Probe  │  │
                      │ (ProbeHostsLiveStatus)     │  │
                      └─────────────┬──────────────┘  │
                                    │                 │
                                    ▼                 ▼
                      ┌───────────────────────────────────────┐
                      │ Tabular ANSI Box or JSON Telemetry    │
                      └───────────────────────────────────────┘
```

### 3.1 Deduplication & Merging Key
Nodes discovered across sources are indexed by lowercase Host IP (or alias if IP is unavailable):
- `key = strings.ToLower(ip)`
- Multiple registrations across subsystems merge into a single `UnifiedFleetNode` record with combined `Subsystems: ["SSH", "Cluster", "SC"]`.
- Role precedence favors cluster-specific designations (`control`, `master`, `worker`).

### 3.2 Liveness Verification
Unless `--fast` / `--no-probe` is specified, `gitmap nodes` runs non-blocking concurrent TCP/SSH connection checks via `cmdssh.ProbeHostsLiveStatus(ctx, hosts)`.
- Ready status displays as: `● ready` (Green)
- Unreachable/offline status displays as: `○ offline (timeout)` (Dim)
- Auth failure displays as: `▲ auth failed` (Yellow)

---

## 4. Output Schemas

### 4.1 ANSI Box Table Output

```
╔══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗
║ GITMAP UNIFIED FLEET NODES (SSH, CLUSTER & SERVER-CLIENTS)                                                       ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝
  Discovered: 6 registered node(s) across SSH, Cluster DB & Server-Client (SC) networks

  ALIAS            ROLE           HOST (IP:PORT)         USER           SUBSYSTEMS         STATUS                ENROLLED
  ──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────
  main             worker         192.168.1.20:22        administrator  SSH, Cluster, SC   ● ready               2026-09-29 12:27:18
  w3               worker         192.168.1.12:22        Administrator  SSH, Cluster, SC   ● ready               2026-09-29 12:33:35
  w1               worker         192.168.1.3:22         Administrator  SSH, Cluster, SC   ● ready               2026-09-29 12:33:35
  u1               worker         192.168.1.22:22        a              SSH, Cluster, SC   ○ offline (timeout)   2026-09-29 12:33:35
  w2               worker         192.168.1.7:22         Administrator  SSH, Cluster, SC   ○ offline (timeout)   2026-09-29 12:33:35
  w4               worker         192.168.1.13:22        Administrator  SSH, Cluster, SC   ○ offline (timeout)   2026-09-29 12:33:35

==============================================================================================================================
 TIP: Run 'gitmap ssh <alias>' for shell access, or 'gitmap sc exec <cmd>' to broadcast across all nodes.
==============================================================================================================================
```

### 4.2 JSON Output Format (`--json`)

```json
[
  {
    "alias": "main",
    "role": "worker",
    "host": "192.168.1.20",
    "port": 22,
    "user": "administrator",
    "status": "● ready",
    "subsystems": ["SSH", "Cluster", "SC"],
    "enrolled_at": "2026-09-29 12:27:18",
    "is_online": true
  }
]
```

---

## 5. Testing & Verification

1. **Unit Tests**:
   - `TestUnifiedNodes_Help`: Validates help banner formatting and flag documentation.
   - `TestUnifiedNodes_RenderTable`: Validates column alignment, ANSI escaping, and empty state rendering.
   - `TestUnifiedNodes_FilterSubsystems`: Validates filtering by SSH, Cluster, and SC flags.
   - `TestUnifiedNodes_JSONFormat`: Validates JSON serialization and structure validity.
2. **Live E2E Verification**:
   - Executed against local test fleet (`main`, `w1`, `w2`, `w3`, `w4`, `u1`).
   - Verified fast execution with `--fast` and targeted filtering with `gitmap nodes main`.
