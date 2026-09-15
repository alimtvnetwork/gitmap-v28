# gitmap cluster join

Join a node to the cluster using SSH enrollment and automated host registration.

## Usage

```bash
gitmap cluster join <user@ip|ip> [alias] [flags]
gitmap cluster node join <user@ip|ip> [alias] [flags]
```

## Description

`gitmap cluster join` performs SSH enrollment to join a remote machine into the cluster topology. It verifies SSH connectivity, records node metadata in the local SQLite database, and assigns a cluster alias for multi-node orchestration.

Joined nodes become immediately available for:
- Role targeting (`all`, `control`, `workers`)
- Remote command execution (`gitmap cluster exec`)
- Provisioning recipes (`gitmap cluster node`)
- Kubernetes cluster orchestration (`gitmap cluster k8s join`)

Target format accepts:
- User and IP: `ubuntu@192.168.0.101`
- Plain IP: `192.168.0.101`
- IP with Port: `192.168.0.101:2222`
- User, IP, and Port: `kube@192.168.0.101:2222`

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--user` | `-u` | current user | Remote SSH username override |
| `--name`, `--alias` | `-n` | `host-<ip>` | Memorable alias name for host recall |
| `--port` | `-p` | 22 | Target SSH port |
| `--auth` | | false | Push local public key to remote `~/.ssh/authorized_keys` |
| `--force` | `-f` | false | Overwrite existing alias or host mapping |
| `--json` | | false | Output result in JSON format |
| `--help` | `-h` | | Show command help |

## Examples

```bash
# 1. Join node by user and IP
gitmap cluster join ubuntu@192.168.0.101 k8s-w1

# 2. Join node with SSH public key authorization
gitmap cluster join kube@192.168.0.102 k8s-w2 --auth

# 3. Join plain IP address with generated alias
gitmap cluster join 192.168.0.103

# 4. Join node with custom SSH port
gitmap cluster join root@192.168.0.104:2222 storage-node

# 5. Join using cluster node join alias
gitmap cluster node join deploy@192.168.0.105 node-5

# 6. Join node and emit JSON output
gitmap cluster join 192.168.0.106 worker-6 --json
```

See also: `gitmap cluster add`, `gitmap cluster nodes`, `gitmap cluster bootstrap`, `gitmap cluster`
