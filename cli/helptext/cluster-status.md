# gitmap cluster status

Inspect cluster connectivity, node health, heartbeat records, and ping latency.

## Usage

```bash
gitmap cluster status [target] [flags]
gitmap cluster ping [target] [flags]
gitmap cluster health [target] [flags]
```

## Description

`gitmap cluster status` displays current heartbeat states, registry health, and connection status for nodes across the cluster.

`gitmap cluster ping` (or `health`) performs active network connectivity probes against registered cluster machines, measuring round-trip latency (RTT) and identifying unreachable nodes.

## Target Selector

- Omit target: Probes all registered cluster nodes.
- `<target>`: Probes a specific node by alias (`k8s-w1`), IP address (`192.168.0.101`), or group (`control`, `workers`).

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--port` | `-p` | 22 | Target SSH port to probe |
| `--timeout` | `-t` | `1.5s` | Probe timeout duration |
| `--json` | | false | Output status results in JSON format |
| `--help` | `-h` | | Show command help |

## Examples

```bash
# 1. Display cluster status and heartbeat overview
gitmap cluster status

# 2. Ping all registered cluster nodes
gitmap cluster ping

# 3. Check health of a specific node by alias
gitmap cluster status k8s-w1
gitmap cluster ping k8s-w1

# 4. Probe specific node by IP address
gitmap cluster ping 192.168.0.101

# 5. Probe with non-standard port and custom timeout
gitmap cluster ping 10.0.0.12 --port 2222 --timeout 2s

# 6. Check health of control plane nodes
gitmap cluster health control

# 7. Output health telemetry in JSON format
gitmap cluster ping --json
```

See also: `gitmap cluster nodes`, `gitmap cluster exec`, `gitmap cluster`
