# gitmap cluster add

Enroll a new node into the cluster topology and SSH host registry.

## Usage

```bash
gitmap cluster add <user@ip|ip> [alias] [flags]
gitmap cluster node add <user@ip|ip> [alias] [flags]
```

## Description

`gitmap cluster add` registers a new target node into the local cluster database (`ssh_hosts` and `ssh_history`). Once enrolled, the node can be targeted across all cluster operations (`gitmap cluster exec`, `gitmap cluster run-script`, `gitmap cluster node`, `gitmap cluster k8s`) by its memorable alias or IP address.

Target format accepts:
- User and IP: `ubuntu@192.168.0.101` or `root@10.0.0.12` (stores user, default alias `host-<ip>`)
- Plain IP: `192.168.0.101` (uses current user, default alias `host-192.168.0.101`)
- IP with Port: `192.168.0.101:2222`
- User, IP, and Port: `admin@10.0.0.15:2222`

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--user` | `-u` | current user | Remote SSH username override |
| `--name`, `--alias` | `-n` | `host-<ip>` | Memorable alias name for host recall |
| `--port` | `-p` | 22 | Target SSH port |
| `--auth` | | false | Push local public key to remote `~/.ssh/authorized_keys` |
| `--force` | `-f` | false | Overwrite existing alias or host mapping |
| `--json` | | false | Output enrollment result in JSON format |
| `--help` | `-h` | | Show command help |

## Examples

```bash
# 1. Enroll node by user and IP with custom alias
gitmap cluster add ubuntu@192.168.0.101 k8s-w1

# 2. Enroll plain IP with automatic alias generation
gitmap cluster add 192.168.0.102

# 3. Enroll node and authorize local SSH public key in one step
gitmap cluster add kube@192.168.0.103 k8s-w3 --auth

# 4. Enroll node with non-standard SSH port
gitmap cluster add dev@192.168.0.104:2222 worker-4

# 5. Overwrite existing alias mapping
gitmap cluster add ubuntu@192.168.0.105 k8s-w1 --force

# 6. Output structured JSON for automation scripts
gitmap cluster add root@192.168.0.106 master-2 --json

# 7. Use cluster node add syntax
gitmap cluster node add 192.168.0.107 k8s-w7
```

See also: `gitmap cluster join`, `gitmap cluster nodes`, `gitmap cluster node`, `gitmap cluster`
