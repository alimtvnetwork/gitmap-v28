# gitmap server-cmd

Delegate remote shell commands and on-the-fly bash scripts across cluster nodes (servers, workers, all).

## Usage

```bash
gitmap server-cmd [flags] "<command>"
gitmap server-cmds [flags] "<command>"
gitmap scmd [flags] "<command>"
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --all | false | Target all registered cluster nodes |
| --control | false | Target control plane / server nodes only |
| --workers | false | Target worker / client nodes only |
| --target <selector> | all | Specific target node name, IP, or role |
| --exclude <list> | - | Comma-separated list of node IDs or IPs to exclude |
| --sudo | false | Run remote command with sudo privilege escalation |
| --script | false | Deploy command as an on-the-fly script in /tmp and execute |
| --dry-run | false | Preview target nodes and command payload without executing |

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
