# gitmap cluster join

Enroll a remote machine into the GitMap cluster topology with automated SSH admission, known_hosts auto-acceptance, encrypted password vaulting, and OS metadata discovery.

## Usage

```bash
gitmap cluster join <user@ip|ip> [alias] [flags]
gitmap cluster node join <user@ip|ip> [alias] [flags]
```

Aliases: `gitmap cluster add`

---

## Why GitMap Cluster Join?

`gitmap cluster join` is the official admission gateway that transforms isolated remote machines into managed members of a synchronized, role-based computing cluster.

### Comparison: `ssh join` vs `cluster join`

| Dimension | `gitmap ssh join` (`sj`) | `gitmap cluster join` |
|-----------|--------------------------|-----------------------|
| **Primary Scope** | Individual host connectivity & SSH alias catalog | Cohesive cluster topology & multi-node orchestration |
| **Membership** | Ad-hoc machines for point-to-point SSH logins | Managed fleet with role tags (`control` vs `worker`) |
| **Execution Context** | Single-node commands (`gitmap ssh exec <target>`) | Parallel fanout (`gitmap cluster exec`, `gitmap sc`) |
| **Orchestration** | Basic remote command execution | Kubernetes lifecycle (`cluster k8s`), recipes (`node`) |
| **Persistence** | SQLite `ssh_hosts` & `SSHConnection` records | Role-aware cluster registry with full OS telemetry |

---

## Step-by-Step Connection & Admission Flow

When `gitmap cluster join` runs, it executes a seven-step connection and enrollment protocol:

```
[Target: user@ip:port]
        │
        ▼ Step 1: TCP Handshake & Port Probing (800ms non-blocking check)
        │
        ▼ Step 2: Automatic Host Key Trust (known_hosts auto-acceptance, no prompt hang)
        │
        ▼ Step 3: Cryptographic Authentication Discovery (id_rsa key check)
        │         ├── Key accepted ───────────┐
        │         └── Key rejected / Missing  │
        ▼                                     │
  Step 4: Interactive Hidden Password Input   │
          • Saved in RSA-OAEP / AES vault     │
          • Review anytime: gitmap ssh pass   │
        │                                     │
        ▼◄────────────────────────────────────┘
  Step 5: Node OS & Environment Probing
          • OS type (Linux / Windows / macOS)
          • OS version & distribution (Ubuntu, Debian, Win10)
          • First-time run timestamp logged
        │
        ▼ Step 6: GitMap Agent Presence Check & Auto-Bootstrap
        │
        ▼ Step 7: Bidirectional Confirmation & Topology Registration
        │
[Node Enrolled & Available for Cluster Commands]
```

### Detailed Step Breakdown:

1. **Step 1: Network Reachability & TCP Handshake**
   GitMap performs a fast, non-blocking TCP socket check against the target IP and SSH port (`800ms` probe). If unreachable, it warns immediately without blocking the terminal.

2. **Step 2: Automatic Host Key Trust (`known_hosts`)**
   Automatically scans the remote SSH host key and appends it directly to `~/.ssh/known_hosts` and the local SQLite security ledger. You are never stuck on an interactive `Are you sure you want to continue connecting (yes/no)?` prompt.

3. **Step 3: Cryptographic Authentication Discovery**
   Attempts connection using default local SSH keys (`~/.ssh/id_rsa`). If public key authentication succeeds, the node connects immediately.

4. **Step 4: Interactive Password Prompt & Encrypted Vault Storage**
   If key authentication requires a password, GitMap securely prompts for the password in hidden CLI format (`terminal.ReadPassword`).
   - The password is encrypted at rest using RSA-OAEP (or AES-GCM fallback).
   - The CLI informs you:
     `ℹ Saving password as encrypted representation (RSA-OAEP/AES) in local vault.`
     `ℹ Review saved password anytime with: gitmap ssh pass show <alias>`
   - The password is saved and reused automatically across all future `ssh exec`, `cluster exec`, and `servers-clients` commands.

5. **Step 5: Node OS & Environment Probing**
   GitMap probes the remote machine for its operating system (`linux`, `windows`, `darwin`) and specific distribution/version (`Ubuntu 22.04 LTS`, `Windows 10 Pro`, etc.). This metadata is saved to SQLite so future commands automatically optimize shell wrappers and binary paths.

6. **Step 6: GitMap Agent Presence & Auto-Bootstrap**
   GitMap checks whether the `gitmap` binary is present on the remote machine (checking PATH, `$HOME/.local/bin`, `/usr/local/bin`, etc.). If missing, it auto-bootstraps GitMap automatically.

7. **Step 7: Bidirectional Confirmation & Topology Registration**
   Performs a bidirectional loop check, commits the node into the local cluster database, and outputs confirmation.

---

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--user` | `-u` | current user | Remote SSH username override |
| `--name`, `--alias` | `-n` | `host-<ip>` | Memorable alias name for cluster recall |
| `--port` | `-p` | 22 | Target SSH port |
| `--auth` | | false | Push local public key to remote `~/.ssh/authorized_keys` |
| `--force` | `-f` | false | Overwrite existing alias or host mapping |
| `--json` | | false | Output result in JSON format |
| `--help` | `-h` | | Show command help |

---

## Password Review & Management

Once a node is joined, view or review saved passwords anytime:

```bash
# Review saved password for a node
gitmap ssh pass show <alias>

# List all nodes with password vault status
gitmap ssh pass ls
```

---

## Examples

```bash
# 1. Join node by user and IP (auto-prompts password if needed)
gitmap cluster join ubuntu@192.168.1.5 u2

# 2. Join node with SSH public key deployment
gitmap cluster join kube@192.168.1.16 main --auth

# 3. Join Windows machine with memorable alias
gitmap cluster join Administrator@192.168.1.8 w2

# 4. Join node with custom SSH port
gitmap cluster join root@192.168.1.14:2222 u1

# 5. Review saved encrypted password for node
gitmap ssh pass show u2
```

See also: `gitmap cluster nodes`, `gitmap cluster exec`, `gitmap ssh pass`, `gitmap servers-clients`
