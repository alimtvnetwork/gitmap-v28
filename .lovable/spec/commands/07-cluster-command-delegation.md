# Cluster Command Delegation — `servers-clients` / `clients` / `sc`

**Spec Number:** 07  
**Slug:** cluster-command-delegation  
**Status:** pending  
**Created:** 2026-08-17  
**Author:** User (transcribed from voice spec)

---

## 1. Problem Statement

Once `gitmap serve` and `gitmap join` establish a cluster of machines (the Multi-Machine Join feature from Phase 6), operators need a single, unified CLI surface to:
- Broadcast shell commands (PowerShell, cmd, bash) across every node or a selected subset.
- Install packages on nodes.
- Delegate standard Git operations (pull, push, commit, status) to the cluster.
- Manage the lifecycle of remote machines (restart, shutdown, logoff).
- Run project-level automation (`run.ps1` / `run.sh`) on remote machines.
- Persist a full audit trail of every remote execution in SQLite.

All of this must be expressible from the **operator's** machine alone — one command, many targets.

---

## 2. Cluster Node Roles (recap from Phase 6)

| Role | Command to start | DB identity |
|------|-----------------|-------------|
| **Server / Orchestrator** | `gitmap serve` | `NodeRole = 'server'` |
| **Client / Worker** | `gitmap join <IP>:<PORT> --token <TOKEN>` | `NodeRole = 'client'` |

Each node registers itself in the server's SQLite `ClusterNode` table with:
- `NodeId` (UUID, auto-generated at first join)
- `Alias` (human-readable, defaults to `hostname`)
- `DisplayId` (sequential 1-based integer, assigned by server)
- `IPAddress`
- `NodeRole` (`server` | `client`)
- `OS` (`windows` | `linux` | `darwin`)
- `JoinedAt`
- `LastHeartbeat`
- `Status` (`online` | `offline` | `unreachable`)
- `PasswordHash` (bcrypt, optional, for restart/logoff operations)

---

## 3. Command Grammar

### 3.1 Target Selectors

There are two primary target groups:

| Selector | Long form | Short hand | Targets |
|----------|-----------|-----------|---------|
| `servers-clients` | all nodes including server | `sc` | Server + all clients |
| `clients` | clients only | *(none)* | All clients, excluding server |

#### Filtering — `--except` / `-except`

Any command may append `--except` (or `-except`) to exclude a subset of nodes. Nodes can be excluded by:

| Exclusion type | Example | Resolves to |
|---------------|---------|-------------|
| Display ID (integer) | `--except 2,4` | Nodes with `DisplayId` 2 and 4 |
| IP address | `--except 192.168.1.24,192.168.1.151` | Nodes matching those IPs |
| Partial trailing IP octet | `--except 24,151` | Nodes whose last octet of IP is 24 or 151 |
| Range | `--except 2-5` | Nodes with `DisplayId` 2 through 5 |

Exclusion tokens are comma-separated, whitespace-tolerant. Mixed types (IDs and IPs) in the same `--except` list are resolved independently and merged.

#### Targeting — `--ip` / `--id`

The inverse of exclusion: run on **only** the named subset.

| Flag | Example | Resolves to |
|------|---------|-------------|
| `--ip` | `--ip 192.168.1.10,192.168.1.11` | Those IP addresses only |
| `--id` | `--id 1,3,5` | Nodes with those `DisplayId` values only |

`--except`, `--ip`, and `--id` are mutually exclusive. Providing two or more raises a validation error.

---

### 3.2 Sub-Commands

#### 3.2.1 `ps` — PowerShell Execution

```
gitmap servers-clients ps "<powershell command>" [--except <list>]
gitmap sc              ps "<powershell command>" [--except <list>]
gitmap clients         ps "<powershell command>" [--except <list>] [--ip <list>] [--id <list>]
```

**Behavior:**
- Dispatches the quoted PowerShell command string to all resolved target nodes.
- On Windows nodes: invoked as `powershell.exe -NonInteractive -Command "<cmd>"` (or `pwsh -Command` if pwsh is available, preferred).
- On Linux/macOS nodes: `pwsh -Command "<cmd>"` if PowerShell Core is installed, otherwise the agent logs a warning and marks the node as `Skipped (no pwsh)`.
- stdout and stderr are streamed back to the operator's terminal in real-time, prefixed with `[<DisplayId>/<Alias>]`.
- Exit code from each node is recorded in `ClusterExecResult`.

**Examples:**
```
gitmap servers-clients ps "Get-Service | Where Status -eq Running"
gitmap sc ps "Get-Date" --except 2,4
gitmap clients ps "Get-DiskUsage C:\" --except 24,151
gitmap clients ps "Restart-Service sshd" --id 3,7
```

---

#### 3.2.2 `cmd` — Windows Command Prompt Execution

```
gitmap servers-clients cmd "<cmd command>" [--except <list>]
gitmap clients         cmd "<cmd command>" [--except <list>] [--ip <list>] [--id <list>]
```

**Behavior:**
- Dispatches the quoted command to `cmd.exe /C "<cmd>"` on Windows nodes.
- On Linux/macOS nodes: treated as a shell command, invoked via `/bin/sh -c "<cmd>"`.
- Same streaming, prefix, and audit behavior as `ps`.

**Examples:**
```
gitmap servers-clients cmd "ipconfig /all" --except 24,151
gitmap clients cmd "dir C:\Projects" --id 1,2,3
```

---

#### 3.2.3 Chained Sub-commands on Specific Targets

Multiple sub-commands (`ps`, `cmd`, `install`) can be chained **on the same `gitmap client` line** when targeting a specific subset:

```
gitmap client cmd "<cmd-command>", ps "<ps-command>" --ip <ip1>,<ip2>
gitmap client cmd "<cmd-command>", ps "<ps-command>" --id <id1>,<id2>
gitmap client cmd "<cmd-command>", ps "<ps-command>", install "<pkg1>,<pkg2>" --id <id1>,<id2>
```

**Rules:**
- Sub-commands separated by `,` are executed **sequentially** on each node (not in parallel).
- Each sub-command's result is captured separately in `ClusterExecResult`.
- If any sub-command fails on a node, subsequent sub-commands on that node are **skipped** (not canceled) and recorded with `ResultStatus = Skipped`.

**Examples:**
```
gitmap client cmd "whoami", ps "Get-Date" --ip 192.168.1.10,192.168.1.11
gitmap client ps "Set-ExecutionPolicy Bypass", install "git,nodejs" --id 1,3
```

---

#### 3.2.4 `install` — Package Installation

```
gitmap servers-clients install "<pkg1>,<pkg2>,<pkg3>" [--except <list>]
gitmap clients         install "<pkg1>,<pkg2>" [--except <list>] [--ip <list>] [--id <list>]
```

**Behavior:**
- Parses the comma-separated package list.
- Detects package manager per node:
  - Windows: `winget install <pkg>` (preferred), fallback to `choco install <pkg>`, fallback to `scoop install <pkg>`.
  - Linux: `apt-get install -y <pkg>` (Debian/Ubuntu), `yum install -y <pkg>` (RHEL), `pacman -S --noconfirm <pkg>` (Arch).
  - macOS: `brew install <pkg>`.
- Each package installed individually so partial success is captured.
- Package manager detection is cached per `NodeId` in `ClusterNode.PackageManager`.

**Examples:**
```
gitmap servers-clients install "git,nodejs,dotnet" --except 2
gitmap clients install "python3,pip3" --id 5,6,7
```

---

#### 3.2.5 Git Operations — `pull`, `push`, `commit`, `status`

These are wrappers that run the equivalent `gitmap` operation on every target node.

```
gitmap servers-clients pull   --all [--except <list>]
gitmap servers-clients push   --all [--except <list>]
gitmap servers-clients commit --all [--except <list>]
gitmap servers-clients status --all [--except <list>] [--ip <list>] [--id <list>]
```

**Behavior:**
- Each target node executes its local `gitmap pull --all`, `gitmap push --all`, etc.
- Results (per-repo OK/FAIL counts) streamed back and aggregated on the server terminal.
- `status --all` returns per-repo dirty/clean state for all nodes in a unified view.

**Examples:**
```
gitmap servers-clients pull --all
gitmap servers-clients push --all --except 3
gitmap servers-clients status --all --id 1,2,4
```

---

#### 3.2.6 `proj` — Project-Level Automation

```
gitmap servers-clients proj "<project-names-csv>" run     [--except <list>]
gitmap servers-clients proj "<project-names-csv>" create-cicd [--except <list>]
```

**`run` behavior:**
- Resolves the named project(s) on each target node by scanning registered repo paths.
- On Windows: executes `run.ps1` in the project root.
- On Linux/macOS: executes `run.sh` in the project root.
- Captures stdout/stderr; on failure, includes the last 20 lines of output in the ClusterExecResult error field.
- Streams progress back to the operator.

**`create-cicd` behavior:**
- Marked as `?` (placeholder — full specification TBD in a dedicated spec file).
- Placeholder implementation: records a `ClusterExecResult` with `ResultStatus = Deferred`.

**Examples:**
```
gitmap servers-clients proj "api-backend,web-frontend" run
gitmap servers-clients proj "api-backend" run --except 2
```

---

#### 3.2.7 Machine Lifecycle — `restart`, `shutdown`, `logoff`

```
gitmap clients "restart"  [--except <list>] [--ip <list>] [--id <list>]
gitmap clients "shutdown" [--except <list>] [--ip <list>] [--id <list>]
gitmap clients "logoff"   [--except <list>] [--ip <list>] [--id <list>]
```

**Behavior:**
- `restart`: gracefully restarts the OS.
  - Windows: `Restart-Computer -Force`
  - Linux/macOS: `sudo shutdown -r now`
- `shutdown`: initiates a full power-off.
  - Windows: `Stop-Computer -Force`
  - Linux/macOS: `sudo shutdown -h now`
- `logoff`: signs out the current user session.
  - Windows: `logoff`
  - Linux: `pkill -KILL -u $USER`
- If the node has a stored `PasswordHash`, the client agent uses the saved credential to authenticate the privileged command. If no password is stored, logs a warning and marks the node as `RequiresAuth`.
- The server machine (`NodeRole = 'server'`) is **NEVER** included in lifecycle commands regardless of the target selector.
- A mandatory 5-second countdown with cancellation prompt is shown before execution: `⚠ Restarting 3 nodes in 5s… Press Ctrl+C to abort`.

**Examples:**
```
gitmap clients "restart" --except 1
gitmap clients "shutdown" --id 5
```

---

### 3.3 Pre-flight Checks (displayed before any execution)

Before dispatching any command, `gitmap` displays a preflight summary:

```
┌─ Preflight ────────────────────────────────────────────────┐
│ Target   : servers-clients (8 nodes)                       │
│ Excluded : Node 2 (192.168.1.24), Node 4 (192.168.1.151)  │
│ Effective: 6 nodes                                         │
│ Command  : ps "Get-Date"                                   │
│ Audit ID : RUN-20260817-001                                │
└────────────────────────────────────────────────────────────┘
Proceed? [y/N] _
```

- The `--yes` / `-Y` flag bypasses the preflight prompt.
- Preflight is skipped for read-only commands (`status`, `proj run` in dry-run).

---

## 4. SQLite Audit Schema

### 4.1 `ClusterRun` Table

Records each top-level command invocation.

| Column | Type | Notes |
|--------|------|-------|
| `ClusterRunId` | INTEGER PK AUTOINCREMENT | |
| `RunRef` | TEXT NOT NULL | Human-readable ID, e.g. `RUN-20260817-001` |
| `CommandKind` | INTEGER NOT NULL | FK → `CommandKindEnum` |
| `RawCommand` | TEXT NOT NULL | The full raw CLI command string |
| `TargetSelector` | TEXT NOT NULL | `servers-clients` / `clients` |
| `ExceptClause` | TEXT | Raw except string |
| `StartedAt` | DATETIME NOT NULL | |
| `FinishedAt` | DATETIME | NULL until complete |
| `TotalNodes` | INTEGER | |
| `SucceededNodes` | INTEGER | |
| `FailedNodes` | INTEGER | |
| `SkippedNodes` | INTEGER | |

### 4.2 `ClusterExecResult` Table

One row per node per sub-command.

| Column | Type | Notes |
|--------|------|-------|
| `ClusterExecResultId` | INTEGER PK AUTOINCREMENT | |
| `ClusterRunId` | INTEGER NOT NULL | FK → `ClusterRun` |
| `NodeId` | TEXT NOT NULL | FK → `ClusterNode` |
| `SubCommand` | TEXT NOT NULL | `ps`, `cmd`, `install`, `pull`, `push`, `commit`, `status`, `restart`, `shutdown`, `logoff` |
| `CommandText` | TEXT | The exact string dispatched |
| `ResultStatus` | INTEGER NOT NULL | FK → `ResultStatusEnum` |
| `ExitCode` | INTEGER | NULL if not applicable |
| `Stdout` | TEXT | Truncated to 64KB |
| `Stderr` | TEXT | Truncated to 16KB |
| `StartedAt` | DATETIME | |
| `FinishedAt` | DATETIME | |
| `DurationMs` | INTEGER | |
| `ErrorMessage` | TEXT | Human-readable failure summary |

### 4.3 Enums

**`CommandKindEnum`:**
```
PsCommand      = 1
CmdCommand     = 2
Install        = 3
GitPull        = 4
GitPush        = 5
GitCommit      = 6
GitStatus      = 7
ProjRun        = 8
ProjCreateCICD = 9
Restart        = 10
Shutdown       = 11
Logoff         = 12
```

**`ResultStatusEnum`** (reuse from bulk-visibility plan):
```
Pending   = 1
Succeeded = 2
Failed    = 3
Skipped   = 4
Deferred  = 5
RequiresAuth = 6
```

---

## 5. Node Registry — Export & Import

```
gitmap cluster export [--format json|csv] [--output <file>]
gitmap cluster import <file>
```

- `export`: dumps the `ClusterNode` table. Passwords are **never** exported (the `PasswordHash` column is omitted).
- `import`: merges node records; existing `NodeId` records are updated, new ones inserted.

---

## 6. Security & Credentials

- All client passwords stored as bcrypt hashes (cost ≥ 12) in the `ClusterNode` table.
- Passwords are set via: `gitmap cluster set-password --id <id>` (prompts securely, never echoed).
- Cluster join tokens (from Phase 6) are one-time-use and stored as SHA-256 hashes.
- All cluster communication over TLS (self-signed cert generated at `gitmap serve` startup, fingerprint printed for client verification).
- On any lifecycle command (`restart`, `shutdown`, `logoff`), the operator must re-confirm even with `-Y` — lifecycle commands require an additional `--force-lifecycle` flag.

---

## 7. Output & UI

- All cluster commands share the same parallel-spinner UI introduced in Phase 4 (bounded worker pool, per-node progress lines, streaming stdout).
- Node lines are prefixed with their display ID and alias: `[1/dev-machine-01]`.
- A final audit summary is printed after every run: `Audit ID: RUN-20260817-001 — view with: gitmap cluster history RUN-20260817-001`.
- `gitmap cluster history` lists past runs; `gitmap cluster history <RunRef>` shows per-node results.

---

## 8. Help Text Structure

Each new command variant needs a help MD file:
- `gitmap/cmd/help/servers-clients.md` — main help with all sub-commands
- `gitmap/cmd/help/sc.md` — 5-line stub pointing at `servers-clients`
- `gitmap/cmd/help/clients.md` — clients-only variant
- `gitmap/cmd/help/cluster-history.md` — audit history browser
- `gitmap/cmd/help/cluster-export.md`
- `gitmap/cmd/help/cluster-set-password.md`

---

## 10. Node Listing Commands

### 10.1 `servers ls`

```
gitmap servers ls
```

Lists all server nodes (nodes with `NodeRole = 'server'`).

**Output columns:** `ID | IP Address | Machine Name | OS | Status | Last Heartbeat`

**Example output:**
```
  ID  IP Address       Machine Name      OS       Status    Last Heartbeat
  ──  ───────────────  ────────────────  ───────  ────────  ──────────────
   1  192.168.1.10     build-server-01   windows  online    2s ago
```

---

### 10.2 `clients ls`

```
gitmap clients ls [--except <list>]
```

Lists all client nodes (`NodeRole = 'client'`), supporting all `--except` / `--ip` / `--id` filter flags.

**Example:**
```
gitmap clients ls
gitmap clients ls --except 2,4
```

**Output columns:** `ID | IP Address | Machine Name | OS | Status | Last Heartbeat`

---

### 10.3 `servers-clients ls` / `sc ls`

```
gitmap servers-clients ls [--except <list>] [--ip <list>] [--id <list>]
gitmap sc ls
```

Lists ALL nodes (servers and clients combined).

**Example:**
```
gitmap servers-clients ls
gitmap sc ls --except 5
```

**Output columns:** `ID | IP Address | Machine Name | Role | OS | Status | Last Heartbeat`

**`--json` flag**: emit machine-readable JSON array of node objects.

---

## 11. Remote Clone / CFR Commands

### 11.1 Grammar

```
gitmap clients         clone  <url1>[,<url2>,...] [--workdir "<path>"] [--except <list>] [--ip <list>] [--id <list>]
gitmap servers-clients clone  <url1>[,<url2>,...] [--workdir "<path>"] [--except <list>]
gitmap clients         cfrp   <url1>[,<url2>,...] [--workdir "<path>"] [--except <list>]
gitmap servers-clients cfrp   <url1>[,<url2>,...] [--workdir "<path>"] [--except <list>]
gitmap clients         cfr    <url1>[,<url2>,...] [--workdir "<path>"] [--except <list>]
gitmap servers-clients cfr    <url1>[,<url2>,...] [--workdir "<path>"] [--except <list>]
```

### 11.2 `--workdir` Flag

| Value | Behavior |
|-------|-----------|
| Omitted | Uses the `DefaultPath` stored for each node (see §12). If no default is set, uses the node's home directory (`~` on Unix, `%USERPROFILE%` on Windows). |
| Absolute path | Clones into that exact directory on the remote node. |
| Path alias (e.g. `p1`) | Resolves the registered alias (see §12) for each target node and clones there. |

### 11.3 Behavior

- URLs are comma-separated, whitespace-tolerant.
- Each URL is cloned (or `cfrp`/`cfr` processed) sequentially per node, in parallel across nodes.
- The same bounded worker-pool + spinner UI from Phase 4 is used.
- After all clones complete, `gitmap status` is automatically triggered on the `--workdir` path of each node and its output streamed back prefixed with `[<ID>/<Alias>]`.
- Each clone result is persisted to `ClusterExecResult`.

### 11.4 Examples

```
gitmap clients clone https://github.com/org/repo1.git,https://github.com/org/repo2.git
gitmap clients cfrp https://github.com/org/myrepo.git --workdir "C:\Projects"
gitmap servers-clients cfr https://github.com/org/app.git --workdir p1 --except 3
gitmap sc clone https://github.com/org/api.git --id 1,2
```

---

## 12. Path Alias System

Operators frequently target the same directories across many commands. The path alias system avoids typing long paths repeatedly and ensures the right directory is used on each node.

### 12.1 `set-default-path`

```
gitmap servers-clients set-default-path "<path> [as <alias>]" [--except <list>]
gitmap clients         set-default-path "<path> [as <alias>]" [--id <list>]
```

Sets a **default working directory** for the targeted nodes. The optional `as <alias>` registers a reusable short name simultaneously.

**Storage:** `ClusterNode.DefaultPath` column + `NodePathAlias` table.

**Examples:**
```
gitmap servers-clients set-default-path "C:\Projects as p1"
gitmap clients set-default-path "/home/dev/repos as work" --id 2,3
```

After running this, every subsequent command that omits `--workdir` will automatically use the registered default on those nodes.

**Tip printed after set:**
```
✓ DefaultPath set on 4 nodes → "C:\Projects" (alias: p1)
  Use with: --workdir p1
  Examples:
    gitmap sc clone <url> --workdir p1
    gitmap clients cat p1/config.yaml
    gitmap clients ps "ls p1"
```

---

### 12.2 `set-path-alias`

```
gitmap servers-clients set-path-alias "<path1> as <alias1>[, <path2> as <alias2>, ...]" [--except <list>]
gitmap clients         set-path-alias "<path1> as <alias1>[, <path2> as <alias2>, ...]" [--id <list>]
```

Registers one or more **named path aliases** that can be used in place of full paths in any subsequent cluster command. Each alias is stored per-node in the `NodePathAlias` table.

**SQLite table `NodePathAlias`:**

| Column | Type | Notes |
|--------|------|-------|
| `NodePathAliasId` | INTEGER PK | |
| `NodeId` | TEXT NOT NULL | FK → `ClusterNode` |
| `Alias` | TEXT NOT NULL | Short name, e.g. `p1` |
| `AbsolutePath` | TEXT NOT NULL | The full path on that node |
| `CreatedAt` | DATETIME | |

**Examples:**
```
gitmap servers-clients set-path-alias "C:\Projects as p1, D:\Backup as p2"
gitmap clients set-path-alias "/opt/apps as apps, /home/dev as home" --id 3,5
```

**Tips printed after registration:**
```
✓ Registered 2 aliases on 6 nodes:
  p1 → C:\Projects
  p2 → D:\Backup

  Use in any cluster command with --workdir or inline:
    gitmap sc clone <url> --workdir p1
    gitmap sc cat p1\config.json
    gitmap sc write p2\notes.txt "content here"
    gitmap clients ps "dir p1" --id 1,3
```

---

### 12.3 Using Path Aliases in Commands

Any command that accepts a path argument or `--workdir` may use a registered alias instead of a full path:

```
gitmap sc clone  https://github.com/org/api.git --workdir p1
gitmap sc clone  https://github.com/org/api.git,https://github.com/org/web.git --path "p1, p2"  # clone first to p1, second to p2
gitmap sc cat    p1/config.yaml
gitmap sc write  p1/notes.txt "deployment notes here"
gitmap clients   ps "Get-ChildItem p1" --id 2
```

**`--path` flag** (multi-alias clone distribution):
- Accepts a comma-separated list of aliases (or full paths).
- Maps the Nth URL to the Nth path alias.
- If the URL list is longer than the path list, remaining URLs use the last alias in the list.
- Post-clone: `gitmap status` is triggered on each resolved path and streamed back to the operator.

**Example with status output:**
```
gitmap sc clone url1,url2 --path "p1, p2"

  Cloning url1 → p1 on 4 nodes…
  Cloning url2 → p2 on 4 nodes…

  [1/build-01] ✓ url1 → C:\Projects (2.3s)   ✓ url2 → D:\Backup (1.8s)
  [2/dev-02]   ✓ url1 → C:\Projects (2.7s)   ✗ url2 → D:\Backup (git auth failed)

  ── Status: [1/build-01] ────────────────────────────────────
  REPO         STATUS   SYNC    BRANCH
  api          clean    ✓       main
  web          clean    ✓       main
  ────────────────────────────────────────────────────────────
```

---

## 13. Remote File Read/Write — `cat` / `write`

### 13.1 `cat`

```
gitmap servers-clients cat <path-or-alias>[/<subpath>] [--except <list>] [--id <list>] [--ip <list>]
gitmap sc cat <path-or-alias>[/<subpath>]
```

Reads and prints the content of a file from each target node. Output is prefixed with `[<ID>/<Alias>]`.

**Behavior:**
- `<path-or-alias>` resolves registered aliases first; falls back to treating it as a literal path.
- If the file is not found on a node, prints `[<ID>/<Alias>] ✗ file not found: <resolved-path>` and continues.
- Binary files are rejected: if the file contains non-UTF-8 bytes, prints `[<ID>/<Alias>] ✗ binary file, cannot display`.
- File content is truncated at **64 KB** with a `[truncated…]` suffix.
- Result stored in `ClusterExecResult` (SubCommand = `cat`, Stdout = file content).

**Examples:**
```
gitmap sc cat p1/config.yaml
gitmap sc cat p1/deploy.log --except 3
gitmap clients cat /etc/hosts --id 1,2,4
gitmap sc cat "C:\Projects\api\appsettings.json"
```

**Example output:**
```
[1/build-01] ── p1/config.yaml ──────────────────────────────
version: "3.9"
services:
  api:
    image: my-api:latest
[1/build-01] ─────────────────────────────────────────────────

[2/dev-02] ── p1/config.yaml ────────────────────────────────
version: "3.9"
services:
  api:
    image: my-api:staging
[2/dev-02] ─────────────────────────────────────────────────
```

---

### 13.2 `write`

```
gitmap servers-clients write <path-or-alias>[/<subpath>] "<content>" [--except <list>] [--id <list>] [--ip <list>]
gitmap sc write <path-or-alias>[/<subpath>] "<content>"
```

Writes the quoted `<content>` string to the specified file on each target node. The write is **atomic** — content is written to a temp file then renamed to prevent partial writes.

**Behavior:**
- Parent directories are created automatically (`mkdir -p` / `New-Item -Force`).
- Content length is capped at **1 MB** — larger content is rejected with a clear error before any network dispatch.
- On success, `[<ID>/<Alias>] ✓ wrote <N> bytes to <resolved-path>` is printed per node.
- On failure, `[<ID>/<Alias>] ✗ write failed: <reason>` is printed; other nodes continue.
- Result stored in `ClusterExecResult` (SubCommand = `write`).

**Security:** `write` requires the `--force-write` flag when targeting `--all` nodes or more than 5 nodes simultaneously (prevents accidental mass overwrites).

**Examples:**
```
gitmap sc write p1/deploy-note.txt "Deployed v6.26.0 on 2026-08-17"
gitmap clients write p1/.env "DB_HOST=10.0.0.5\nDB_PORT=5432" --id 1,2
gitmap sc write "C:\Projects\api\version.txt" "6.26.0" --force-write
```

---

## 14. Updated Help Text Manifest

Add help MD files for the new commands:

- `gitmap/cmd/help/servers-ls.md`
- `gitmap/cmd/help/clients-ls.md`
- `gitmap/cmd/help/servers-clients-ls.md`
- `gitmap/cmd/help/cluster-set-default-path.md` — includes the alias tips block
- `gitmap/cmd/help/cluster-set-path-alias.md` — includes the alias reference examples block
- `gitmap/cmd/help/cluster-cat.md`
- `gitmap/cmd/help/cluster-write.md`

---

## 15. Out of Scope (for this spec)

- `proj create-cicd` — placeholder only, full spec TBD.
- Cross-cloud or internet-routed cluster connectivity (LAN-only in this version).
- GUI/web dashboard for cluster management.
- Binary file transfer (only UTF-8 text via `cat`/`write`; large binary transfer is out of scope).

