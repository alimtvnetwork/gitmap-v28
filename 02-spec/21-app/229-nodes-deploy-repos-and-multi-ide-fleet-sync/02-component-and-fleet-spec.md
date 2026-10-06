# Component & Fleet Specification: Nodes Deploy Repos & Multi-IDE Sync Engine

> **Document Version:** 1.0.0  
> **Status:** Active  
> **Scope:** Core Go Engine (`cli/cmdnodes`, `cli/cmdssh`, `cli/workspacesync`, `cli/cmd`)  
> **Traceability:** Task-229  

---

## 1. Component Boundaries & Responsibilities

The fleet repository deployment subsystem consists of four primary components:

### Component 1: Fleet Node Filter Engine (`cli/cmdnodes/nodes_filter.go`)
- **Responsibilities**:
  - Unifies node filtering across all fleet commands (`nodes clone`, `nodes deploy agm-accounts`, `nodes deploy repo`, `nodes scan`).
  - Implements `NodeFilterOptions`:
    - `Target string`: Single node alias or IP.
    - `Except []string`: Blacklisted node aliases (case-insensitive). Defaults to `["main"]`.
    - `Include []string`: Whitelisted node aliases (case-insensitive).
    - `IncludeMain bool`: Flag to override default `main` exclusion.
    - `OpenOnly bool`: Probe liveness before returning list of nodes.
  - Excludes local machine connection (`isLocalMachineConnection(c)`).

### Component 2: Fleet Repo Deploy Engine (`cli/cmdnodes/nodes_deploy_repo.go`)
- **Responsibilities**:
  - Resolves target repository from local database (`store.OpenDefault()`) or filesystem path.
  - Evaluates deployment strategy:
    - **Local Archive Streaming (`--from-local`)**: Bundles repository into in-memory `.tar.gz` stream (excluding build artifacts and ephemeral folders), streams over SSH via `cmdssh.StreamFileToRemote`, and extracts on the remote host.
    - **Remote Clone Delegation (`--clone`)**: Delegates `gitmap clone <url>` or `gitmap cfr <url>` directly on the remote machine via SSH.
  - Executes remote post-deployment hooks:
    - Invokes remote `gitmap rescan` to refresh remote `gitmap.db`.
    - Invokes remote multi-IDE registration across VS Code, Cursor, Antigravity, and GitHub Desktop.
    - If `--with-pinned` is set, updates remote `~/.gemini/config/pinned_projects.json`.
    - If `--with-conversations` is set, streams and unpacks matching Antigravity conversations and brain logs.
  - Implements `RunNodesDeployRepo(args []string) error` and `RunNodesDeployRepos(args []string) error`.

### Component 3: Multi-IDE Remote Registration Hook (`cli/cmdnodes/nodes_deploy_ide.go`)
- **Responsibilities**:
  - Translates local filesystem paths to target node operating system paths:
    - Windows: `D:\work\<slug>` (or configured `--dest` path).
    - Linux: `~/work/<slug>` (or `/home/<user>/work/<slug>`).
  - Constructs remote execution script for the target operating system:
    - **Windows**: PowerShell commands to register into VS Code `projects.json`, Cursor `projects.json`, Antigravity `projects/<uuid>.json`, and probe `github.bat`.
    - **Linux**: POSIX shell commands to create config folders, write JSON entries, and update IDE project registries.

### Component 4: Remote Scanning Delegation & SSH Core Routing (`cli/cmdnodes/nodes_scan.go`, `cli/cmdssh/ssh_exec_command.go`)
- **Responsibilities**:
  - Adds `"scan"` and `"rescan"` to `isGitmapCoreCommand` in `cli/cmdssh/ssh_exec_command.go`, enabling `gitmap ssh exec <node> scan` without requiring explicit `"gitmap "` prefix.
  - Implements `RunNodesScan(args []string) error` and `RunNodesRescan(args []string) error` to broadcast repository scanning across all reachable fleet nodes and print structured status tables.

---

## 2. Data Structures & Contract Interfaces

```go
package cmdnodes

import "time"

// DeployRepoOptions defines the runtime configuration for repository fleet deployment.
type DeployRepoOptions struct {
    RepoSlug          string
    TargetNode        string
    ExceptNodes       []string
    IncludeNodes      []string
    IncludeMain       bool
    OpenOnly          bool
    DestPath          string
    WithIDEs          bool
    WithPinned        bool
    WithConversations bool
    FromLocal         bool
    FromClone         bool
    Clean             bool
    DryRun            bool
    JSONOutput        bool
}

// DeployRepoResult captures per-node deployment outcomes.
type DeployRepoResult struct {
    NodeAlias      string        `json:"node_alias"`
    Host           string        `json:"host"`
    OS             string        `json:"os"`
    RepoPath       string        `json:"repo_path"`
    RepoDeployed   bool          `json:"repo_deployed"`
    IDEsRegistered []string      `json:"ides_registered"`
    IsPinned       bool          `json:"is_pinned"`
    ConvsDeployed  int           `json:"convs_deployed"`
    Status         string        `json:"status"` // SUCCESS, DRY-RUN, OFFLINE, FAILED
    Latency        time.Duration `json:"latency"`
    LatencyMs      int64         `json:"latency_ms"`
    ErrorMessage   string        `json:"error_message,omitempty"`
}
```

---

## 3. Remote IDE Script Templates

### 3.1 Windows Target Execution Script (PowerShell)
```powershell
# Setup paths
$dest = "D:\work\<repoName>"
$appdata = [Environment]::GetFolderPath('ApplicationData')
$userprofile = $env:USERPROFILE

# 1. VS Code & Cursor Project Manager
$dirs = @(
    (Join-Path $appdata 'Code\User\globalStorage\alefragnani.project-manager'),
    (Join-Path $appdata 'Cursor\User\globalStorage\alefragnani.project-manager')
)
foreach ($d in $dirs) {
    if (-not (Test-Path $d)) { New-Item -ItemType Directory -Force -Path $d | Out-Null }
    $pJson = Join-Path $d 'projects.json'
    # Upsert entry
}

# 2. Antigravity Project Config
$agyDir = Join-Path $userprofile '.gemini\config\projects'
if (-not (Test-Path $agyDir)) { New-Item -ItemType Directory -Force -Path $agyDir | Out-Null }
# Write <uuid>.json if not existing

# 3. GitHub Desktop
$ghBat = Join-Path $env:LOCALAPPDATA 'GitHubDesktop\bin\github.bat'
if (Test-Path $ghBat) { & $ghBat $dest }

# 4. Remote GitMap Rescan
gitmap rescan
```

### 3.2 Linux / POSIX Target Execution Script (POSIX Shell)
```sh
dest="$HOME/work/<repoName>"

# 1. VS Code & Cursor Project Manager
for d in "$HOME/.config/Code/User/globalStorage/alefragnani.project-manager" "$HOME/.config/Cursor/User/globalStorage/alefragnani.project-manager"; do
    mkdir -p "$d"
    # Upsert entry into projects.json
done

# 2. Antigravity Project Config
mkdir -p "$HOME/.gemini/config/projects"
# Write <uuid>.json

# 3. GitHub Desktop (if on PATH)
if command -v github >/dev/null 2>&1; then
    github "$dest" >/dev/null 2>&1 || true
fi

# 4. Remote GitMap Rescan
gitmap rescan >/dev/null 2>&1 || true
```

---

## 4. Verification & Testing Strategy

1. **Unit Tests**:
   - `cli/cmdnodes/nodes_filter_test.go`: Test node exclusion logic, comma-separated tokens, and default exclusion of `main`.
   - `cli/cmdnodes/nodes_deploy_repo_test.go`: Test flag parsing, archive packaging, and dry-run output formatting.
2. **Integration & E2E Validation**:
   - Run `gitmap nodes deploy repo gitmap --dry-run` to verify node discovery and simulated actions.
   - Run `gitmap nodes scan --open-only` to verify remote scanning delegation across live nodes.
3. **Coding Standards & Linters**:
   - Enforce zero nested `if` statements via `check-nested-ifs.py`.
   - Enforce enum and boolean standards via `check-enum-and-boolean.py`.
   - Enforce relative path rules via `check-relative-paths.py`.
