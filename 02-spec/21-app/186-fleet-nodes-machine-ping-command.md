# Specification 186: Fleet Nodes Machine Ping Command (`gitmap nodes ping` & `gitmap ping`)

## 1. Executive Summary & Architectural Mission

Operational fleets manage heterogeneous nodes across local networks, cloud providers, and remote SSH environments. While port-level status checks provide high-level reachability flags, diagnosing raw network connectivity requires invoking the host operating system's native machine `ping` command (ICMP echo transmission), capturing packet loss rates, and measuring round-trip times (RTT).

Furthermore, Windows machines on local subnets frequently enforce Windows Defender Firewall rules that silently drop inbound ICMP echo requests while allowing SSH (TCP port 22) connections. A pure ICMP ping misclassifies such healthy nodes as offline, whereas a pure TCP check fails to exercise the OS ping utility requested by network operators.

This specification introduces:
1. **`gitmap nodes ping`** & **`gitmap ping`**: Top-level and subsystem commands executing concurrent native machine `ping` commands across all registered fleet nodes.
2. **Dual-Stack ICMP & TCP Reachability**: Concurrently executes OS ping and tests port reachability to accurately distinguish `● ONLINE` (ICMP response), `● REACHABLE (TCP)` (ICMP firewall filtered, port open), and `○ OFFLINE`.
3. **Structured Telemetry & Terminal UI**: Outputs ANSI boxed tables with packet counts, loss percentages, RTT metrics, summary statistics, raw command logs (`--raw`), and machine-readable JSON (`--json`).

---

## 2. CLI Command Syntax

### 2.1 Invocation Patterns

- `gitmap ping [target] [flags]`
- `gitmap nodes ping [target] [flags]`
- `gitmap ping help`

### 2.2 Flags & Options

| Flag | Short | Purpose | Default |
| :--- | :--- | :--- | :--- |
| `--count <n>` | `-c <n>`, `-n <n>` | Number of ping packets to transmit per host | `2` |
| `--timeout <dur\|ms>` | `-t <dur>`, `-w <ms>` | Probe timeout per node (e.g., `2s`, `1500ms`, `1500`) | `1500ms` |
| `--raw` | | Displays exact executed machine command lines and stdout/stderr | `false` |
| `--json` | `-j` | Emits structured JSON array of node ping telemetry | `false` |
| `--ssh` | | Limits ping targets to nodes enrolled in SSH connection vault | `false` |
| `--cluster` | `--clst` | Limits ping targets to nodes enrolled in Cluster Fleet DB | `false` |
| `--sc` | `--server-clients` | Limits ping targets to nodes enrolled in Server-Clients network | `false` |
| `[target]` | | Optional target filter by Alias, Host IP, or ad-hoc IP/hostname | All Nodes |

### 2.3 Command Aliases
- `gitmap ping`
- `gitmap nodes ping`
- `gitmap node ping`
- `gitmap fleet-ping`
- `gitmap nodes-ping`

---

## 3. Architecture & Execution Pipeline

```
                     ┌────────────────────────────────┐
                     │ gitmap nodes ping / gitmap ping│
                     └───────────────┬────────────────┘
                                     │
                        [Target specified in args?]
                             ├── No  ──► Collect all registered nodes
                             └── Yes ──► Filter fleet nodes / Ad-hoc node fallback
                                     │
                                     ▼
                     ┌────────────────────────────────┐
                     │ Concurrent Worker Pool (10)    │
                     └───────────────┬────────────────┘
                                     │
           ┌─────────────────────────┴─────────────────────────┐
           ▼                                                   ▼
┌──────────────────────────────┐              ┌──────────────────────────────┐
│ Machine Command Execution    │              │ TCP Port Probe               │
│ Windows: ping -n <c> -w <ms> │              │ net.DialTimeout(tcp, host:22)│
│ Linux:   ping -c <c> -W <s>  │              │                              │
└──────────────┬───────────────┘              └──────────────┬───────────────┘
               │                                             │
               └───────────────────────┬─────────────────────┘
                                       ▼
                     ┌────────────────────────────────┐
                     │ Telemetry Normalization &      │
                     │ Status Resolution:             │
                     │ - recv > 0: ● ONLINE           │
                     │ - recv=0 & tcp: ● REACHABLE    │
                     │ - recv=0 & !tcp: ○ OFFLINE     │
                     └───────────────┬────────────────┘
                                     │
                     ┌───────────────┴────────────────┐
                     ▼                                ▼
       [--json flag active?]             [ANSI Box Table Display]
       Emits JSON Array to stdout        Columns: ALIAS, ROLE, HOST,
                                         ICMP STATS, LOSS, AVG RTT,
                                         TCP (PORT), STATUS
```

---

## 4. Output Formats

### 4.1 ANSI Box Table Output

```text
╔══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗
║ GITMAP FLEET NODES PING (MACHINE ICMP & TCP PROBE)                                                               ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝
  Probing 6 node(s) across fleet with machine command 'ping' (packets: 2, timeout: 1500ms)...

  ALIAS          ROLE         HOST (IP:PORT)         ICMP STATS      LOSS     AVG RTT    TCP (PORT)       STATUS
  ────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────
  w1             worker       node-w1:22         2/2 rcvd        0%       <1ms       ● 22 (<1ms)      ● ONLINE
  main           worker       node-main:22        0/2 rcvd        100%     -          ● 22 (2ms)       ● REACHABLE (TCP)
  w3             worker       node-w3:22        0/2 rcvd        100%     -          ● 22 (2ms)       ● REACHABLE (TCP)
  u1             worker       node-u1:22        0/2 rcvd        100%     -          ○ timeout        ○ OFFLINE
  w2             worker       node-w2:22         0/2 rcvd        100%     -          ○ timeout        ○ OFFLINE
  w4             worker       node-w4:22        0/2 rcvd        100%     -          ○ timeout        ○ OFFLINE

  ────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────
  Fleet Ping: 3/6 reachable (1 online ICMP, 2 reachable via TCP) | 3 offline | Max Elapsed: 2.01s

==============================================================================================================================
 TIP: Run 'gitmap ssh <alias>' for shell access, or 'gitmap ping <alias>' to ping an individual node.
==============================================================================================================================
```

### 4.2 Raw Command Execution Output (`--raw`)

```text
--- [w1 | node-w1] command: C:\Windows\system32\ping.exe -n 2 -w 1500 node-w1 ---
Pinging node-w1 with 32 bytes of data:
Reply from node-w1: bytes=32 time<1ms TTL=128
Reply from node-w1: bytes=32 time<1ms TTL=128

Ping statistics for node-w1:
    Packets: Sent = 2, Received = 2, Lost = 0 (0% loss),
Approximate round trip times in milli-seconds:
    Minimum = 0ms, Maximum = 0ms, Average = 0ms
```

### 4.3 JSON Telemetry (`--json`)

```json
[
  {
    "alias": "w1",
    "role": "worker",
    "host": "node-w1",
    "port": 22,
    "user": "Administrator",
    "subsystems": ["SSH", "Cluster", "SC"],
    "command": "C:\\Windows\\system32\\ping.exe -n 2 -w 1500 node-w1",
    "packets_sent": 2,
    "packets_received": 2,
    "packets_lost": 0,
    "packet_loss_percent": 0,
    "min_rtt": "0ms",
    "avg_rtt": "<1ms",
    "max_rtt": "0ms",
    "icmp_ok": true,
    "tcp_ok": true,
    "tcp_rtt": "<1ms",
    "status": "● ONLINE",
    "duration_ms": 1022061000
  }
]
```

---

## 5. Verification & Testing

- Unit tests in `cli/cmd/nodes_ping_cmd_test.go` validate:
  - Option parsing and flag overrides (`--count`, `--timeout`, `--raw`, `--json`, `--ssh`, target).
  - Windows ping output regex parser (success, timeout, packet loss, RTT).
  - Linux ping output regex parser (transmitted, received, loss, RTT).
  - Reachability status resolution.
  - ANSI table rendering.
  - JSON serialization.
  - Help text flags.
- Live verification executed across test fleet and loopback target.
