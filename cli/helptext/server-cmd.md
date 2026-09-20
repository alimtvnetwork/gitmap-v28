# gitmap server-cmd

Delegate remote shell commands and on-the-fly bash scripts across cluster nodes (servers, workers, all).

## Usage

```bash
gitmap server-cmd [flags] "<command>"
gitmap server-cmds [flags] "<command>"
gitmap scmd [flags] "<command>"
```

---

## Server Command Connection Architecture & Step-by-Step Flow

When `gitmap server-cmd` (or `gitmap sc`) dispatches commands across remote server nodes, it executes the following six-step connection protocol:

```
[gitmap server-cmd --all "uname -a"]
                 │
                 ▼ Step 1: Node Target Resolution
                 │   • Reads target nodes from SQLite (SSHConnection / ssh_hosts)
                 │   • Filters by role: --control, --workers, or --all
                 ▼ Step 2: Credential Decryption
                 │   • Decrypts stored passwords in-memory using SSH RSA private key
                 │   • Prepares sudo -S elevation buffers if --sudo is specified
                 ▼ Step 3: Transport & Host Key Handshake
                 │   • Connects over SSH port with auto-trusted known_hosts
                 │   • Non-blocking parallel worker pool (default 4 workers)
                 ▼ Step 4: OS-Aware Shell & Command Wrapping
                 │   • Linux/macOS nodes: wrapped in bash -c or sh -c
                 │   • Windows nodes: wrapped in powershell.exe or cmd.exe
                 ▼ Step 5: Real-Time Streaming Execution
                 │   • Captures interleaved stdout/stderr prefixed by node: [node-1] ...
                 ▼ Step 6: Telemetry & Summary Aggregation
                 │   • Aggregates exit codes and displays completion status
                 ▼
[Execution Complete across all Active Nodes]
```

### Detailed Step Breakdown:

1. **Step 1: Node Target Resolution**
   GitMap reads the local database to find all active nodes matching your selector (`--all`, `--control`, `--workers`, `--target <alias|ip>`). Nodes currently offline are flagged and reported cleanly.

2. **Step 2: Credential & Password Decryption**
   If remote nodes require password authentication or sudo escalation (`--sudo`), GitMap decrypts the node's vaulted password in-memory using your local SSH RSA private key. The password is never logged or exposed.

3. **Step 3: Transport & Host Key Handshake**
   Establishes secure SSH client connections using the auto-trusted `known_hosts` ledger. Handshakes occur concurrently using a goroutine worker pool.

4. **Step 4: OS-Aware Shell Wrapping**
   GitMap checks each node's probed OS metadata (`linux`, `windows`, `darwin`) saved during enrollment. Commands on Linux are routed to `/bin/bash` or `/bin/sh`; commands on Windows are routed to `powershell.exe -NoProfile` or `cmd.exe`.

5. **Step 5: Streaming Execution**
   Commands run remotely in real-time. Stdout and stderr lines are prefixed with the node's friendly alias (e.g. `[u2 | 192.168.1.5]`) as they arrive.

6. **Step 6: Completion Telemetry**
   GitMap verifies the process exit code across all nodes and displays a summary table detailing success or failure for each machine.

---

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--all` | false | Target all registered cluster nodes |
| `--control` | false | Target control plane / server nodes only |
| `--workers` | false | Target worker / client nodes only |
| `--target <selector>` | all | Specific target node name, IP, or role |
| `--exclude <list>` | - | Comma-separated list of node IDs or IPs to exclude |
| `--sudo` | false | Run remote command with sudo privilege escalation |
| `--script` | false | Deploy command as an on-the-fly script in /tmp and execute |
| `--dry-run` | false | Preview target nodes and command payload without executing |

---

## Password Review & Management

To review or verify saved passwords for any cluster node:

```bash
gitmap ssh pass show <alias>
gitmap ssh pass ls
```

---

## Examples

### Run command on all nodes

```bash
gitmap server-cmd --all "uname -a"
```

Output:

```text
▶ Delegating command to 3 node(s)...
  ✓ [node-1 (192.168.1.10)] Executed: uname -a
  ✓ [node-2 (192.168.1.11)] Executed: uname -a
  ✓ [node-3 (192.168.1.12)] Executed: uname -a
✓ Cluster execution completed successfully.
```

### Run script on worker nodes with sudo

```bash
gitmap server-cmd --workers --sudo --script "systemctl restart docker"
```

Output:

```text
▶ Delegating command to 2 node(s)...
  ✓ [node-2 (192.168.1.11)] Executed: sudo -S sh -c "mkdir -p /tmp/on-the-fly-cmd..."
  ✓ [node-3 (192.168.1.12)] Executed: sudo -S sh -c "mkdir -p /tmp/on-the-fly-cmd..."
✓ Cluster execution completed successfully.
```

See also: `gitmap cluster join`, `gitmap cluster exec`, `gitmap servers-clients`, `gitmap ssh pass`
