# gitmap nodes-running-only — Active Cluster Node Discovery & Execution

Discover and execute operations strictly on online and reachable cluster fleet nodes, skipping offline machines with 10s liveness caching and zero timeout latency.

## Synopsis

```bash
gitmap nodes-running-only [flags]
gitmap nodes-running-only [flags] <subcommand> [args...]
gitmap nro [flags]
gitmap nro [flags] <subcommand> [args...]
```

**Aliases:** `nro`, `nodes-running`, `nodesrunningonly`, `nodes-online`, `nodes-open`, `nodes-active`

## Overview

Standard fleet operations query all registered nodes across SSH, Cluster DB, and Server-Client networks, including offline or sleeping machines. When running automated pipelines or interactive commands, waiting for offline connection timeouts introduces latency and noisy warning banners.

`gitmap nodes-running-only` solves this by:
1. **Strict Online Filtering**: Restricts output and operations exclusively to verified reachable nodes (`is_online: true`).
2. **Dual-Tier Liveness Caching**: Caches probe results in-memory and in the OS temporary directory with a 10-second TTL, enabling instantaneous response times across rapid sequential CLI invocations.
3. **Transparent Subcommand Delegation**: Seamlessly forwards execution verbs (`exec`, `tasks`, `resources`, `ping`, `pull`, `cp`, `scan`) to active targets without attempting connections to offline machines.

## Delegation Subcommands

When a subcommand verb is provided, `gitmap nro` filters the fleet to online nodes first, then delegates execution:

| Subcommand | Aliases | Description |
|---|---|---|
| `exec <cmd>` | `run` | Execute remote shell command exclusively on online cluster nodes |
| `tasks [target]` | `queue` | Query active and pending AI agent task queues on online nodes |
| `resources` | `res`, `top` | Interrogate CPU, RAM, and disk utilization on reachable nodes |
| `ping [target]` | `probe` | Run fast reachability and latency probes against online nodes |
| `pull [target]` | `pa` | Fan-out repository pull operations across active fleet nodes |
| `cp "<msg>"` | `cm` | Broadcast flat commit and push across reachable nodes |
| `scan [target]` | `rescan` | Broadcast remote repository discovery scans across online nodes |

If no subcommand is passed, `gitmap nro` lists all currently active nodes in a formatted ANSI table or JSON array.

## Flags & Options

| Flag | Short | Default | Description |
|---|---|---|---|
| `--json` | `-j` | `false` | Output active nodes as a clean JSON array (`is_online: true` only) |
| `--refresh` | `-r` | `false` | Bypass 10s liveness cache and trigger immediate network probe |
| `--target <alias>` | `-t` | `""` | Filter or target a specific node alias or IP address |
| `--timeout <dur>` | | `1500ms` | Network reachability probe and connection timeout limit |
| `--verbose` | `-v` | `false` | Display verbose probe diagnostics and cache lookup details |
| `--help` | `-h` | `false` | Display command documentation |

## Output Formats

### ANSI Table (Default)

```text
+---------------+-----------------+---------+--------------+---------+---------+
| NODE ALIAS    | HOST            | OS      | ROLE         | STATUS  | LATENCY |
+---------------+-----------------+---------+--------------+---------+---------+
| ai-main       | 127.0.0.1       | windows | controlplane | ONLINE  | 0ms     |
| worker-gpu-01 | 192.168.1.101   | linux   | worker       | ONLINE  | 14ms    |
| worker-cpu-02 | 192.168.1.102   | linux   | worker       | ONLINE  | 18ms    |
+---------------+-----------------+---------+--------------+---------+---------+
3 active nodes online (2 offline skipped) [cached: 8s remaining]
```

### JSON (`--json`, `-j`)

Guaranteed to contain only records where `is_online` is `true`:

```json
[
  {
    "alias": "ai-main",
    "host": "127.0.0.1",
    "os": "windows",
    "role": "controlplane",
    "is_online": true,
    "latency_ms": 0
  },
  {
    "alias": "worker-gpu-01",
    "host": "192.168.1.101",
    "os": "linux",
    "role": "worker",
    "is_online": true,
    "latency_ms": 14
  }
]
```

## Examples

### Discover Active Nodes

```bash
# List active nodes using ANSI table
gitmap nro

# Force cache bypass and probe all nodes immediately
gitmap nro --refresh

# Output online nodes in JSON for scripting
gitmap nro --json
```

### Remote Command Execution

```bash
# Run uptime command across all online cluster nodes
gitmap nro exec "uptime"

# Check available disk space on online nodes
gitmap nro exec "df -h /"

# Target a specific node if online
gitmap nro -t worker-gpu-01 exec "nvidia-smi"
```

### Telemetry & Agent Queues

```bash
# Check CPU and RAM utilization on active nodes
gitmap nro resources

# Interrogate active AI agent task queues on online nodes
gitmap nro tasks

# Probe network latency across active nodes
gitmap nro ping
```

### Scripting & Automation

```bash
# Iterate through active node aliases with jq
for node in $(gitmap nro --json | jq -r '.[].alias'); do
  echo "Deploying update to $node..."
done
```

## See Also

- `gitmap nodes`: Unified fleet node discovery and management
- `gitmap nodes tasks`: Inspect active and pending AI agent tasks across fleet
- `gitmap nodes resources`: Probe CPU and RAM telemetry across fleet
- `gitmap cluster`: Multi-node cluster topology and SSH key enrollment
