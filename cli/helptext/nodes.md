# gitmap nodes — Unified Fleet & Infrastructure Nodes

Aggregates and orchestrates fleet operations across SSH Cluster Nodes (`SSHConnection`), Cluster Fleet Database (`ClusterNode`), and Server-Clients Networks.

## Synopsis

```bash
gitmap nodes [flags] [target]
gitmap nodes deploy agm-accounts [flags]
gitmap nodes sync-agm-accounts [flags]
gitmap nodes ping [target]
gitmap nodes push-settings <node>
gitmap nodes sync-settings
gitmap nodes send-projects [node]
gitmap nodes agy prompt <node> <project> "<prompt>"
gitmap nodes agy query [--ssh]
gitmap nodes agy [ui]
gitmap nodes clone [flags] <targets> [dest]
gitmap nodes cfr [flags] [targets] [dest]
gitmap nodes cfrp [flags] [targets] [dest]
```

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `deploy agm-accounts` | Broadcast Antigravity Manager accounts & credentials across fleet nodes |
| `sync-agm-accounts` | Shorthand alias for `deploy agm-accounts` |
| `push-settings <node>` | Export and push Antigravity settings to target node |
| `sync-settings` | Broadcast Antigravity settings across all fleet nodes |
| `send-projects [node]` | Forward VS Code / Cursor Project Manager workspaces to nodes |
| `ping [target]` | Run machine ping and reachability probe across fleet nodes |
| `clone <targets>` | Fan-out repository cloning asynchronously across local and remote nodes |
| `cfr <targets>` | Clone, fix, and auto-setup repository across fleet nodes |
| `cfrp <targets>` | Clone, fix, and promote public across fleet nodes |
| `agy prompt <node>` | Dispatch prompt directly to target project on node |
| `agy query` | Query active Antigravity instances and prompts |
| `agy ui` | Launch Antigravity web studio dashboard |

---

## Deploy AGM Accounts (`deploy agm-accounts`)

Deploy and synchronize Antigravity Manager (AGM) account credentials, active sessions, and tokens from the local machine across all remote fleet nodes over secure SSH channels.

### How It Works

1. **Local Account Discovery**:
   - Inspects `~/.antigravity_tools/`.
   - Packages `accounts.json`, all account files under `accounts/*.json`, and `user_tokens.db` into an in-memory `.tar.gz` archive.
   - Computes total count of available accounts.
2. **Target Filtering**:
   - Skips local machine (`127.0.0.1`, `localhost`, local interfaces).
   - Excludes the `main` node by default to protect primary control plane, unless `--include-main` is explicitly passed.
   - Applies `--target` alias or `--except` exclusions.
3. **Liveness & Resilience**:
   - Probes node reachability via rapid 1.5s TCP/SSH check.
   - Unreachable nodes are marked as `OFFLINE` without aborting deployment to remaining nodes.
4. **Remote Streaming & Extraction**:
   - Streams the in-memory archive over SSH to staging path:
     - Linux: `/tmp/agm_accounts_sync.tar.gz`
     - Windows: `C:\Windows\Temp\agm_accounts_sync.tar.gz`
   - Unpacks directly into `~/.antigravity_tools`.
   - Hardens filesystem permissions (mode `700` on directories, `600` on secrets).

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--target <alias>` | `-t` | `""` | Target single node alias or IP address |
| `--except <list>` | `-e` | `"main"` | Comma-separated list of node aliases/IPs to exclude |
| `--include-main` | | `false` | Include `main` node in deployment (overrides default exclusion) |
| `--open-only` | | `false` | Only return / display online and reachable nodes in output |
| `--dry-run` | `-n` | `false` | Probe nodes and simulate packaging without remote changes |
| `--json` | `-j` | `false` | Output structured JSON array of `DeployAGMResult` |
| `-h`, `--help` | | `false` | Display command help |

### Table Columns

| Column | Description |
|--------|-------------|
| `NODE ALIAS` | Registered fleet node identifier |
| `HOST` | Target IP address or hostname |
| `OS` | Remote operating system (`linux`, `windows`, `darwin`) |
| `ACCOUNTS` | Number of accounts successfully deployed |
| `STATUS` | Status badge (`● SUCCESS`, `◌ DRY-RUN`, `○ OFFLINE`, `▲ FAILED`) |
| `LATENCY` | Round-trip deployment latency |

### Examples

```bash
# Broadcast AGM accounts to all worker nodes (excluding main and local machine)
gitmap nodes deploy agm-accounts

# Deploy to a specific worker node
gitmap nodes deploy agm-accounts --target worker-1

# Deploy to all nodes including main
gitmap nodes deploy agm-accounts --include-main

# Exclude specific nodes
gitmap nodes deploy agm-accounts --except worker-2,192.168.1.15

# Simulate deployment without modifying remote nodes
gitmap nodes deploy agm-accounts --dry-run

# Output machine-readable JSON
gitmap nodes deploy agm-accounts --json

# Using top-level aliases
gitmap sync-agm-accounts
gitmap deploy-agm-accounts --open-only
gitmap nodes-deploy-agm-accounts --target worker-1
```

## See Also

- `gitmap agm` — Antigravity Manager GUI & tools management
- `gitmap ssh` — SSH fleet keys and machine management
- `gitmap cluster` — Cluster topology and multi-node orchestration
