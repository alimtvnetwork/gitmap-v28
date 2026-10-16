# gitmap nodes — Unified Fleet & Infrastructure Nodes

Aggregates and orchestrates fleet operations across SSH Cluster Nodes (`SSHConnection`), Cluster Fleet Database (`ClusterNode`), and Server-Clients Networks.

## Synopsis

```bash
gitmap nodes [flags] [target]
gitmap nodes cp [target|all] "<message>" [flags]
gitmap nodes cm [target|all] "<message>" [flags]
gitmap nodes resources [flags]
gitmap nodes tasks [flags] [target]
gitmap nodes deploy agm-accounts [flags]
gitmap nodes sync-agm-accounts [flags]
gitmap nodes deploy repo <slug|path> [flags]
gitmap nodes deploy repos [targets|all] [flags]
gitmap nodes scan [target] [flags]
gitmap nodes rescan [target] [flags]
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
| `cp [target\|all]` | Broadcast flat commit and push across fleet nodes in parallel |
| `cm [target\|all]` | Broadcast flat commit without push across fleet nodes in parallel |
| `tasks [target]` | Interrogate active and pending AI agent tasks across fleet nodes |
| `resources` | Interrogate CPU, RAM, and disk utilization across unified fleet nodes |
| `deploy agm-accounts` | Broadcast Antigravity Manager accounts & credentials across fleet nodes |
| `sync-agm-accounts` | Shorthand alias for `deploy agm-accounts` |
| `deploy repo <slug>` | Deploy repository, run remote scan, and register across VS Code, Cursor, Antigravity, GitHub Desktop |
| `deploy repos [list]` | Batch deploy multiple repositories across fleet nodes |
| `scan [target]` | Broadcast remote repository discovery scanner across fleet nodes |
| `rescan [target]` | Broadcast remote repository rescan across fleet nodes |
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

## Examples

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

---

## Deploy Repository & Multi-IDE Sync (`deploy repo`, `deploy repos`)

Deploy a repository from the local machine to remote fleet nodes, trigger remote `gitmap rescan`, and automatically register the repository across all installed IDEs on the target nodes:
- **VS Code**: Added to Project Manager `projects.json`.
- **Cursor**: Added to Project Manager `projects.json`.
- **Antigravity IDE**: Added to `~/.gemini/config/projects/<uuid>.json`.
- **GitHub Desktop**: Registered via the `github <path>` CLI shim.

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--target <alias>` | `-t` | `""` | Target a single node alias or IP address |
| `--except <list>` | `-e` | `"main"` | Comma-separated list of nodes to exclude |
| `--exclude <list>` | | | Alias for `--except` |
| `--include <list>` | | `""` | Comma-separated whitelist of nodes to target |
| `--accept <list>` | | `""` | Alias for `--include` |
| `--include-main` | | `false` | Explicitly include `main` node |
| `--open-only` | | `true` | Skip offline nodes immediately via fast preflight probe |
| `--dest <path>` | `-d` | `""` | Destination root directory (`D:\work` on Windows, `~/work` on Linux) |
| `--with-ides` | | `true` | Register across VS Code, Cursor, Antigravity, and GitHub Desktop |
| `--with-pinned` | | `false` | Synchronize pinned status in `pinned_projects.json` and tag as pinned |
| `--with-conversations` | | `false` | Bundle, remap, and transfer associated Antigravity conversations & brain logs |
| `--from-local` | | `false` | Force direct local-to-remote archive transfer over SSH |
| `--clone` | | `false` | Delegate remote Git clone from origin URL |
| `--clean` | | `true` | Exclude `node_modules/`, `target/`, `.venv/`, `dist/` from archive |
| `--dry-run` | `-n` | `false` | Simulate deployment without writing to disk |
| `--json` | `-j` | `false` | Output machine-readable JSON array of `DeployRepoResult` |

### Examples

```bash
# Preview deployment of gitmap to fleet nodes
gitmap nodes deploy repo gitmap --dry-run

# Deploy repository to worker-1 with pinned projects and conversation sync
gitmap nodes deploy repo my-app --target worker-1 --with-pinned --with-conversations

# Batch deploy multiple repositories
gitmap nodes deploy repos api-service,web-frontend --open-only

# Deploy using root-level aliases
gitmap deploy-repo gitmap --dry-run
gitmap nodes-deploy-repos all --dry-run
```

---

## Fleet Scan & Rescan (`scan`, `rescan`)

Broadcast repository discovery across all reachable fleet nodes via SSH delegation, updating remote `gitmap.db` indices.

```bash
# Scan default work directories across all online nodes
gitmap nodes scan --open-only

# Rescan repositories on a specific node
gitmap nodes rescan --target worker-1

# Using root aliases
gitmap fleet-scan
gitmap nodes-rescan
```

---

## Fleet Multi-Repo Commits (`cp`, `cm`, `cpf`, `cpb`, `cpr`, `commit-fix`)

Delegate flat WIP or semantic conventional commits and pushes across remote fleet nodes in parallel over SSH.

```bash
gitmap nodes cp all "<message>" [flags]
gitmap nodes cm all "<message>" [flags]
gitmap nodes cpf <repoName|all> "<message>" [flags]
```

### How It Works

1. **Fleet Discovery & Filtering**:
   - Queries all registered SSH connections from `gitmap.db`.
   - Filters candidate nodes via `--target <alias>`, `--except <nodes>`, `--include <nodes>`, `--open-only`.
   - Probes SSH reachability; offline nodes are marked `○ OFFLINE` without failing the batch execution.
2. **Parallel SSH Dispatch**:
   - Dispatches `gitmap sends cp all "<message>" --json` (or `gitmap sends cm all` for local commit) concurrently across all targeted fleet nodes.
   - Transparently forwards `--dry-run` and `--no-push` flags to remote nodes.
3. **Outcome Parsing & Aggregation**:
   - Parses structured remote JSON outcome envelopes.
   - Renders a unified status table detailing node alias, repository name, commit hash, duration, and status badge (`● PUSHED`, `● COMMITTED`, `○ NO CHANGES`, `○ OFFLINE`, `▲ FAILED`).

### Examples

```bash
# Broadcast WIP commit and push across all dirty repositories on all fleet nodes
gitmap nodes cp all "wip: save changes"

# Commit locally on remote nodes without pushing
gitmap nodes cm all "wip: checkpoint"

# Simulate fleet WIP commits without changing remote Git refs
gitmap nodes cp all "wip: preview" --dry-run

# Target a specific node for conventional feature commit
gitmap nodes cpf my-repo "Feature: add cluster health monitor" -t worker-1
```

---

## Fleet Node Resources Telemetry (`resources`)

Probes and visualizes CPU and RAM utilization across cluster fleet nodes and local machine to determine idle capacity.

```bash
gitmap nodes resources [flags]
```

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--target <alias>` | `-t` | `""` | Target single node alias or IP address |
| `--threshold <N>` | | `30` | Resource benchmark utilization percentage limit |
| `--open-only` | | `false` | Only query online reachable nodes |
| `--json` | `-j` | `false` | Output machine-readable JSON telemetry |

### Benchmark Utilization Limit

Fleet nodes with CPU and RAM below the configured benchmark threshold (default: 30%) are considered idle. When local machine resources are idle and no other heavy node processes are active, GitMap's adaptive concurrency booster dynamically scales parallel worker routines up to maximize pull and clone throughput.

---

## Fleet AI Agent Task Queues (`tasks`)

Interrogates active, pending, and heavy-duty AI agent tasks enqueued in `agent-task.db` across all cluster nodes and local workspaces.

```bash
gitmap nodes tasks [target] [flags]
```

### How It Works

1. **Unified Node Discovery**:
   - Queries both local and remote nodes using unified fleet enumeration (`collectUnifiedNodes`).
   - Discovers nodes across registered SSH connections, SSH hosts, and cluster database entries.
2. **Task Database Interrogation**:
   - Interrogates `.ai-memory/temp-agents/*/agent-task.db` for task and subtask status.
   - Summarizes active subtasks, pending queues, and step progression.
3. **Heavy-Duty Task Detection**:
   - Automatically tags tasks as heavy-duty (`isHeavyDuty: true`) if step budget >= 100, task slug contains heavy keywords (`heavy`, `full-audit`, `deep-scan`, `migration`, `build-all`), or owned files exceed 8.
   - Highlights candidates for node-to-node work stealing and cluster load rebalancing.
4. **Parallel SSH Dispatch**:
   - Queries remote nodes concurrently via SSH using a bounded worker pool (`maxConcurrency = 8`) and a 5-second per-node timeout to prevent hanging connections.
   - Executes remote `gitmap tasks --json --limit 10` with fallback SQLite queries.

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--target <alias>` | `-t` | `""` | Target single node alias or IP address |
| `--heavy-only` | | `false` | Filter only heavy-duty tasks exceeding 100 steps budget |
| `--open-only` | | `false` | Only query online and reachable nodes |
| `--json` | `-j` | `false` | Output machine-readable JSON array of `FleetTasksPayload` |
| `-h`, `--help` | | `false` | Display command help |

### Examples

```bash
# Display aggregated task queues across all fleet nodes
gitmap nodes tasks

# Interrogate task queue on a specific remote node
gitmap nodes tasks vm-02

# Filter only heavy-duty tasks across the fleet
gitmap nodes tasks --heavy-only

# Output machine-readable JSON
gitmap nodes tasks --json
```

## See Also

- `gitmap agm` — Antigravity Manager GUI & tools management
- `gitmap ssh` — SSH fleet keys and machine management
- `gitmap cluster` — Cluster topology and multi-node orchestration
