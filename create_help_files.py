import os

help_dir = r"d:\work\gitmap\gitmap\cmd\help"

files = {
    "servers-clients.md": """# gitmap servers-clients

Broadcast shell commands, git operations, or installations across all server and client nodes in the cluster.

## Usage

```
gitmap servers-clients <subcommand> [args] [flags]
```

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `ps` | Execute a PowerShell command |
| `cmd` | Execute a Windows Command Prompt command |
| `install` | Install packages (comma-separated list) |
| `pull` | Run git pull --all on all nodes |
| `push` | Run git push --all on all nodes |
| `commit` | Run git commit --all on all nodes |
| `status` | Show combined dirty/clean status across nodes |
| `proj` | Run project-level automation (run, create-cicd) |

## Flags

| Flag | Description |
|------|-------------|
| `--except <list>` | Exclude nodes by ID, IP, or trailing IP octet |
| `--yes`, `-Y` | Bypass preflight confirmation prompt |

## Examples

### PowerShell & Command Prompt
Execute PowerShell and cmd commands across the cluster:
```
gitmap servers-clients ps "Get-Service | Where Status -eq Running"
gitmap servers-clients cmd "ipconfig /all" --except 24,151
```

### Installation
Install packages simultaneously:
```
gitmap servers-clients install "git,nodejs,dotnet" --except 2
```

### Git Operations
Delegate git operations to all joined machines:
```
gitmap servers-clients pull --all
gitmap servers-clients push --all --except 3
gitmap servers-clients status --all
```

### Project Automation
Run the local `run.ps1` or `run.sh` for specific projects:
```
gitmap servers-clients proj "api-backend,web-frontend" run
gitmap servers-clients proj "api-backend" run --except 2
```

See also: `gitmap sc`, `gitmap clients`
""",

    "sc.md": """# gitmap sc

Alias for `gitmap servers-clients`. Broadcasts shell commands and git operations across the entire cluster.

See: `gitmap servers-clients --help`
""",

    "clients.md": """# gitmap clients

Broadcast shell commands, git operations, or lifecycle actions across all **client** nodes in the cluster (excludes the server/orchestrator).

## Usage

```
gitmap clients <subcommand> [args] [flags]
```

## Subcommands

Supports all standard execution commands (`ps`, `cmd`, `install`, `pull`, `push`, `status`, `proj`) plus lifecycle commands:
| Subcommand | Description |
|------------|-------------|
| `restart` | Gracefully restart the OS |
| `shutdown` | Initiate a full power-off |
| `logoff` | Sign out the current user session |

## Target Filtering Flags

| Flag | Description |
|------|-------------|
| `--except <list>` | Exclude nodes by ID, IP, or trailing IP octet |
| `--ip <list>` | Run *only* on the named IPs |
| `--id <list>` | Run *only* on the named Display IDs |

*Note: `--except`, `--ip`, and `--id` are mutually exclusive.*

## Examples

### Targeting Specific Clients
```
gitmap clients ps "Get-DiskUsage C:\\" --except 24,151
gitmap clients cmd "dir C:\\Projects" --id 1,2,3
```

### Chaining Commands
Execute multiple sub-commands sequentially on a target:
```
gitmap clients cmd "whoami", ps "Get-Date" --ip 192.168.1.10,192.168.1.11
gitmap clients ps "Set-ExecutionPolicy Bypass", install "git,nodejs" --id 1,3
```

### Machine Lifecycle
Restart or shutdown machines (requires `--force-lifecycle` and stored password):
```
gitmap clients restart --except 1
gitmap clients shutdown --id 5
```
""",

    "cluster-history.md": """# gitmap cluster history

View the audit trail of past cluster executions.

## Usage

```
gitmap cluster history [RunRef]
```

## Description

Every cluster execution (e.g., `gitmap sc ps "..."`) is assigned a unique `RunRef` format ID (e.g., `RUN-20260817-001`). These runs and their results are stored locally in the SQLite `ClusterRun` and `ClusterExecResult` tables.

- Omitting `[RunRef]` lists all past cluster runs in a table.
- Providing `[RunRef]` expands the results to show execution outcome, exit code, and stdout/stderr per-node.

## Table Columns (Listing Mode)

| Column | Description |
|--------|-------------|
| RUN ID | The unique RunRef string |
| COMMAND| The raw command string executed |
| TARGET | The target selector used (e.g., servers-clients) |
| SUCCESS| Number of nodes that successfully executed |
| FAILED | Number of nodes that failed |
| TIME   | Started at / Duration |

## Examples

```
gitmap cluster history
gitmap cluster history RUN-20260817-001
```
""",

    "cluster-export.md": """# gitmap cluster export

Export the cluster node registry to a file.

## Usage

```
gitmap cluster export [--format json|csv] [--output <file>]
```

## Description

Dumps the `ClusterNode` SQLite table. 

**Security Note:** Passwords are never exported. The `PasswordHash` field is strictly omitted from all export formats.

## Formats

- `json` (default): Exports an array of JSON objects.
- `csv`: Exports a CSV table.

## Examples

```
gitmap cluster export --format json --output nodes.json
gitmap cluster export --format csv > nodes.csv
```
""",

    "cluster-import.md": """# gitmap cluster import

Import a cluster node registry file (JSON or CSV).

## Usage

```
gitmap cluster import <file>
```

## Description

Merges node records from the file into the local SQLite `ClusterNode` table.

**Merge Semantics:**
- Records with an existing `NodeId` are updated.
- Records with a new `NodeId` are inserted as new nodes.
- Unchanged records are safely skipped.

## Examples

```
gitmap cluster import nodes.json
gitmap cluster import backup.csv
```
""",

    "cluster-set-password.md": """# gitmap cluster set-password

Store a client node password for privileged lifecycle commands (restart, shutdown, logoff).

## Usage

```
gitmap cluster set-password --id <id>
```

## Description

Prompts for a password securely (without echoing to the terminal). The password is encrypted as a bcrypt hash (cost ≥ 12) and stored in the local SQLite `ClusterNode` table. It is used transparently when dispatching lifecycle commands to that node.

**Security Notes:**
- Passwords are only stored as bcrypt hashes.
- The raw text is never written to disk or exported.
- Lifecycle commands using this password require TLS communication.

## Examples

```
gitmap cluster set-password --id 5
```
""",

    "cluster-nodes.md": """# gitmap cluster nodes

List registered nodes in the cluster.

## Usage

```
gitmap cluster nodes [--json]
```

## Description

Displays the list of nodes currently registered in the server's database.

## Table Columns

| Column | Description |
|--------|-------------|
| ID | Sequential 1-based integer Display ID |
| ALIAS | Human-readable alias (usually hostname) |
| IP | IP address of the node |
| ROLE | `server` or `client` |
| OS | `windows`, `linux`, or `darwin` |
| STATUS | `online`, `offline`, or `unreachable` |

## Flags

- `--json`: Emit the listing as a machine-readable JSON array of node objects.

## Examples

```
gitmap cluster nodes
gitmap cluster nodes --json
```
"""
}

for filename, content in files.items():
    filepath = os.path.join(help_dir, filename)
    with open(filepath, "w", encoding="utf-8") as f:
        f.write(content)

print("Created help files.")
