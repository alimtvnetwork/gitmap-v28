# Component and CLI Specification: Child-Path Inventory, Muse Installer, Smart Repo Create, AGM Secure Deploy & Git Commit Cache

- **Spec ID:** `02-spec/21-app/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/02-component-and-cli-spec.md`
- **Architecture Spec:** [01-architecture-spec.md](01-architecture-spec.md)
- **Status:** `APPROVED`
- **Version:** `1.0.0`
- **Target Subsystems:**
  - `cli/cmdchildpath/` (`gitmap child-path`: High-speed directory inventory, globbing, AST replacement for PowerShell `Get-ChildItem` and `git grep`)
  - `cli/cmdinstall/` (`gitmap install muse` / `gitmap muse install`: Cross-platform Meta Muse companion installer with screenshot verification)
  - `cli/cmd/` (`gitmap repo create` / `gitmap cr`: Safe remote existence pre-flight, interactive clone fallback, headless `--clone-on-exists`)
  - `cli/cmdnodes/` (`gitmap agm deploy`: AES-256-GCM in-transit credential transport, cluster sync, cryptographic zero-fill shredding of local accounts)
  - `cli/cmdgitcache/` (`gitmap log --cached` & `gitmap branch compare` / `gitmap cmp`: SQLite commit graph caching in `.gitmap/data/git-cache/<slug>/sql.db`, 0ms fast path, ahead/behind calculation, terminal rendering)
- **Dependencies:** `cli/store`, `cli/crypto`, `cli/termtable`, `cli/termpad`, `cli/apperror`, `cli/gitutil`
- **Target Version:** `v6.501.0`

---

## 1. Executive Summary & Component Topology

This specification defines the component models, CLI argument contracts, data structures, SQLite Split-DB schemas, and execution workflows for five mission-critical subsystems in GitMap `v6.501.0`. These subsystems extend GitMap's autonomous capabilities across local developer workflows, cloud node orchestration, and AI agent execution:

1. **Child-Path Fast Inventory (`gitmap child-path`):** A blazing-fast, cross-platform directory scanner that eliminates slow, brittle PowerShell `Get-ChildItem` and repo-bound `git grep` invocations for AI coding agents.
2. **Meta Muse AI Companion Installer (`gitmap install muse` / `gitmap muse install`):** Unified cross-platform installation engine supporting Windows, Linux, and macOS with verified integrity checks against `assets/screenshots/240-muse-meta-installer.png`.
3. **Smart Safe Repository Creation (`gitmap repo create` / `gitmap cr`):** Eliminates dangerous force-push overwrites on existing remote repositories by introducing remote pre-flight probes, interactive clone confirmations, headless safety gates (`--clone-on-exists`), and fallback to `cmdclone`.
4. **AGM Cluster Fleet Deploy with Secure Shredding (`gitmap agm deploy`):** Distributes sensitive Antigravity tool credentials to remote cluster nodes via AES-256-GCM encrypted bundles and executes 3-pass zero-fill cryptographic shredding of local disk artifacts upon verified receipt.
5. **Git Commit Graph & Branch Compare Cache (`gitmap log` & `gitmap branch compare` / `gitmap cmp`):** Implements an isolated Split SQLite cache (`.gitmap/data/git-cache/<slug>/sql.db`) providing sub-millisecond (0ms fast path) commit log queries and visual branch ahead/behind diffs via `termtable`/`termpad`.

```
+------------------------------------------------------------------------------------------------------------------+
|                                        GITMAP COMPONENT TOPOLOGY                                                 |
+------------------------------------------------------------------------------------------------------------------+
|                                                                                                                  |
|   +-----------------------+   +-----------------------+   +--------------------+   +-------------------------+   |
|   |   gitmap child-path   |   |  gitmap install muse  |   | gitmap repo create |   |   gitmap agm deploy     |   |
|   |    (cmdchildpath)     |   |     (cmdinstall)      |   |       (cmd)        |   |       (cmdnodes)        |   |
|   +-----------+-----------+   +-----------+-----------+   +---------+----------+   +------------+------------+   |
|               |                           |                         |                           |                |
|               v                           v                         v                           v                |
|       [os.ReadDir AST]             [OS Installers]          [Remote Probes]            [AES-256-GCM Bundle]      |
|       - MaxDepth Filter           - Win: irm | iex          - git ls-remote            - Ephemeral Key           |
|       - Ext & Type Gates          - Nix: curl | bash        - gh repo view             - SSH Fleet Transfer      |
|       - JSON Output Envelope      - Verify Screenshot       - Clone Fallback           - 3-Pass Zero Shred       |
|               |                           |                         |                           |                |
|               +---------------------------+------------+------------+---------------------------+                |
|                                                        |                                                         |
|                                                        v                                                         |
|                               +----------------------------------------------------+                             |
|                               |       gitmap log & gitmap branch compare           |                             |
|                               |                  (cmdgitcache)                     |                             |
|                               +------------------------+---------------------------+                             |
|                                                        |                                                         |
|                                                        v                                                         |
|                                       +----------------------------------+                                       |
|                                       | Split SQLite Commit Cache Engine |                                       |
|                                       | .gitmap/data/git-cache/<slug>/   |                                       |
|                                       |             sql.db               |                                       |
|                                       +----------------------------------+                                       |
+------------------------------------------------------------------------------------------------------------------+
```

---

## 2. Subsystem 1: Child-Path Fast Inventory & Search (`gitmap child-path`)

### 2.1 Problem Statement & Architectural Justification
In autonomous coding loops (Google Antigravity, Cursor, Cline), agents frequently explore repository structure and locate source files. Historically, workflows relied on two mechanisms that introduce severe instability and latency:
1. **PowerShell `Get-ChildItem` Limitations:**
   - **Performance Degradation:** Allocates high-overhead .NET `FileInfo` objects for every file visited, taking seconds on deep trees (such as mono-repos or nested build caches).
   - **Cross-Platform Inconsistency:** Invocations behave differently between Windows PowerShell 5.1 and PowerShell Core 7 (`pwsh`), breaking script portability.
   - **Encoding & Parsing Hazards:** Piped output suffers from UTF-16/ANSI encoding issues, line truncations, and localized date-time formatting that breaks AI regex scrapers (documented in `assets/screenshots/240-get-childitem-filter.png`).
2. **`git grep` Limitations:**
   - **Git Index Binding:** Fails completely on untracked files, uncommitted directories, scaffolding templates, or repos undergoing initial construction.
   - **Inability to Walk Outside Repos:** Cannot inspect local scratch folders, caches, or `.ai-memory` paths excluded by `.gitignore`.
   - **Lack of Metadata:** Cannot yield file sizes, modification times, or structural directory depth in a structured format (documented in `assets/screenshots/240-git-grep-func-runfix.png`).

`gitmap child-path` replaces both with a dedicated, dependency-free Go scanner utilizing kernel-level directory iterators (`os.ReadDir`), custom glob matching, and unified JSON serialization.

### 2.2 CLI Command Contract & Invocation Syntax
```bash
gitmap child-path [dir] [glob] [flags]
# Shorthand aliases:
gitmap cp [dir] [glob] [flags]
gitmap childpath [dir] [glob] [flags]
```

#### Arguments
- `[dir]`: Target root directory to inspect. Defaults to `.` (current working directory).
- `[glob]`: Optional filename or pattern match (e.g., `*.go`, `*test*`, `*.md`). Defaults to `*`.

#### Flags
| Flag | Shorthand | Type | Default | Description |
| :--- | :--- | :--- | :--- | :--- |
| `--depth <N>` | `-d <N>` | int | `0` (unlimited) | Maximum directory recursion depth (`0` = unlimited, `1` = immediate children only). |
| `--ext <list>` | `-e <list>` | string | `""` | Comma-separated list of allowed file extensions (e.g., `.go,.ts,.json`). Leading dot optional. |
| `--type <t>` | `-t <t>` | string | `"f"` | Entry type filter: `"f"` (files only), `"d"` (directories only), `"all"` (both). |
| `--limit <N>` | `-n <N>`, `-l <N>` | int | `0` (unlimited) | Maximum number of matched entries to emit. |
| `--json` | `-j` | bool | `false` | Emits structured JSON envelope to stdout instead of ANSI table. |
| `--include-hidden` | | bool | `false` | Includes hidden files/directories (starting with `.`). Default excludes `.git`, `.cache`, `node_modules`. |
| `--size` | `-s` | bool | `false` | Displays human-readable file sizes and modification timestamps in table view. |

### 2.3 Data Structures & Go Models (`cli/cmdchildpath/childpath_types.go`)

```go
package cmdchildpath

import "time"

// EntryTypeFilter defines entry selection (files, directories, or all).
type EntryTypeFilter string

const (
	EntryTypeFilesOnly EntryTypeFilter = "f"
	EntryTypeDirsOnly  EntryTypeFilter = "d"
	EntryTypeAll       EntryTypeFilter = "all"
)

// ChildPathOptions encapsulates parsed command-line flags.
type ChildPathOptions struct {
	TargetDir       string          `json:"targetDir"`
	GlobPattern     string          `json:"globPattern"`
	MaxDepth        int             `json:"maxDepth"`
	Extensions      []string        `json:"extensions"`
	TypeFilter      EntryTypeFilter `json:"typeFilter"`
	Limit           int             `json:"limit"`
	IsJSON          bool            `json:"isJson"`
	IsIncludeHidden bool            `json:"isIncludeHidden"`
	IsShowSize      bool            `json:"isShowSize"`
}

// ChildPathItem represents a single discovered filesystem entry.
type ChildPathItem struct {
	RelPath      string    `json:"relPath"`
	Name         string    `json:"name"`
	IsDir        bool      `json:"isDir"`
	SizeBytes    int64     `json:"sizeBytes"`
	SizeHuman    string    `json:"sizeHuman,omitempty"`
	ModTime      time.Time `json:"modTime"`
	Extension    string    `json:"extension,omitempty"`
	Depth        int       `json:"depth"`
}

// ChildPathResponse defines the top-level JSON telemetry envelope.
type ChildPathResponse struct {
	TargetDir     string          `json:"targetDir"`
	TotalScanned  int             `json:"totalScanned"`
	TotalMatched  int             `json:"totalMatched"`
	DurationMs    int64           `json:"durationMs"`
	Items         []ChildPathItem `json:"items"`
}
```

### 2.4 Error Codes & Exit Behavior
| Error Code | Constant | Trigger Condition | Exit Code |
| :--- | :--- | :--- | :--- |
| `E1081` | `ErrInvalidDirectory` | Target directory does not exist or permission denied | 1 |
| `E1082` | `ErrInvalidTypeFilter` | `--type` value is not `"f"`, `"d"`, or `"all"` | 1 |
| `E1083` | `ErrInvalidGlobPattern`| Glob pattern contains illegal syntax | 1 |
| `E1084` | `ErrDepthOutOfBounds`  | `--depth` is negative | 1 |

---

## 3. Subsystem 2: Meta Muse AI Companion Installer (`gitmap install muse`)

### 3.1 Overview & Visual Verification Protocol
The `cli/cmdinstall` package coordinates dependency provisioning across developer workstations. `gitmap install muse` (and alias `gitmap muse install`) automates the setup of Meta's Muse AI development environment. 

To ensure complete alignment with official releases, the installer verifies steps against reference screenshot `assets/screenshots/240-muse-meta-installer.png`:
- Renders branded installation banners via `termpad`.
- Detects the host operating system (`runtime.GOOS`).
- Dispatches platform-native curl/PowerShell pipelines.
- Validates the resulting binary in `PATH` via `muse --version`.

### 3.2 CLI Command Contract
```bash
gitmap install muse [flags]
gitmap muse install [flags]
gitmap install-muse [flags]
```

#### Flags
| Flag | Shorthand | Type | Default | Description |
| :--- | :--- | :--- | :--- | :--- |
| `--dry-run` | `-n` | bool | `false` | Prints the exact installation script command without executing. |
| `--verify-only`| `-v` | bool | `false` | Checks if Muse is already installed and outputs version and path. |
| `--force` | `-f` | bool | `false` | Overwrites existing Muse installation. |
| `--install-dir <path>` | | string | `""` | Custom binary destination path. |

### 3.3 Cross-Platform Workflows
```
+----------------------------------------------------------------------------------------------------+
|                                    MUSE INSTALLATION DISPATCH                                      |
+----------------------------------------------------------------------------------------------------+
|                                                                                                    |
|                                    [OS Detection: runtime.GOOS]                                    |
|                                                 |                                                  |
|                   +-----------------------------+-----------------------------+                    |
|                   |                             |                             |                    |
|                   v                             v                             v                    |
|             [windows]                        [linux]                       [darwin]                |
|                   |                             |                             |                    |
|                   v                             v                             v                    |
|        powershell -NoProfile         curl -fsSL dev.meta.ai         curl -fsSL dev.meta.ai         |
|         -ExecutionPolicy Bypass            | bash                        | bash                    |
|     -Command "irm dev.meta.ai/install.ps1 | iex"                                                   |
|                   |                             |                             |                    |
|                   +-----------------------------+-----------------------------+                    |
|                                                 |                                                  |
|                                                 v                                                  |
|                                     [Post-Install Verification]                                    |
|                                    - Check Binary in PATH / Root                                   |
|                                    - Run 'muse --version'                                          |
|                                    - Compare with 240-muse-meta-installer.png                     |
+----------------------------------------------------------------------------------------------------+
```

1. **Windows Workflow:**
   ```powershell
   powershell.exe -NoProfile -ExecutionPolicy Bypass -Command "irm https://dev.meta.ai/install.ps1 | iex"
   ```
   - Target binary location: `$env:LOCALAPPDATA\Programs\Muse\muse.exe` or `$env:USERPROFILE\.muse\bin\muse.exe`.
   - Adds directory to Windows User `Path` environment registry key if not already present.
2. **Linux Workflow:**
   ```bash
   bash -c "curl -fsSL https://dev.meta.ai/install.sh | bash"
   ```
   - Fallback: if `curl` is missing, checks for `wget -qO- https://dev.meta.ai/install.sh | bash`.
   - Target binary location: `~/.local/bin/muse` or `/usr/local/bin/muse`.
3. **macOS (Darwin) Workflow:**
   ```bash
   bash -c "curl -fsSL https://dev.meta.ai/install.sh | bash"
   ```
   - Validates architecture (`arm64` vs `x86_64`) to ensure Rosetta 2 is not needlessly invoked.

### 3.4 Data Structures & Go Models (`cli/cmdinstall/install_muse.go`)

```go
package cmdinstall

import "time"

// MuseInstallOptions holds user options for installing Meta Muse.
type MuseInstallOptions struct {
	IsDryRun     bool   `json:"isDryRun"`
	IsVerifyOnly bool   `json:"isVerifyOnly"`
	IsForce      bool   `json:"isForce"`
	InstallDir   string `json:"installDir,omitempty"`
}

// MuseInstallResult represents the execution outcome.
type MuseInstallResult struct {
	IsSuccess    bool          `json:"isSuccess"`
	Platform     string        `json:"platform"`
	Version      string        `json:"version"`
	BinaryPath   string        `json:"binaryPath"`
	DurationMs   int64         `json:"durationMs"`
	ErrorMessage string        `json:"errorMessage,omitempty"`
}
```

### 3.5 Error Codes
| Error Code | Constant | Trigger Condition | Exit Code |
| :--- | :--- | :--- | :--- |
| `E1085` | `ErrMuseDownloadFailed` | Installation script failed to download (HTTP error/network down) | 1 |
| `E1086` | `ErrMissingDownloadTool`| Neither curl nor powershell/irm found on system | 1 |
| `E1087` | `ErrMuseProbeFailed`   | Post-install verification command `muse --version` failed | 1 |
| `E1088` | `ErrUnsupportedArch`   | Host processor architecture unsupported | 1 |

---

## 4. Subsystem 3: Smart Safe Repository Creation & Remote Clone Fallback (`gitmap repo create`)

### 4.1 Problem Statement & Dangerous Force-Push Elimination
In previous implementations (`cli/cmd/repo_create_remote.go`), when creating a repository with `gitmap repo create` (alias `gitmap cr`), the system invoked `gh repo create`. If GitHub reported that the repository already existed, `handleExistingRemoteRepo` immediately performed:
```go
cmdPush := exec.Command("git", "-C", absDir, "push", "-u", "origin", "main", "--force")
```
This behavior was **catastrophically dangerous**:
- If a developer or AI agent inadvertently typed an existing repository name, their local initialized directory would force-push over remote `main`, destroying production commit history and branches.
- In headless CI/CD and AI agent environments, there was no prompt, leading to silent code erasure.

### 4.2 Safe Architecture & Fallback Protocol
The upgraded `cli/cmd/` architecture introduces a four-step defensive gate:
1. **Pre-Flight Remote Existence Check:** Prior to modifying local git remotes, GitMap checks if the remote exists using `git ls-remote` (or `gh repo view <slug>`).
2. **Interactive Terminal Prompt:**
   - If interactive (`isInteractiveStdin() == true`):
     ```text
     Remote repository 'org/slug' already exists.
     Clone existing repo instead? [y/N]: 
     ```
   - If user answers `y` / `Y`: Cancels local creation and seamlessly hands off to `cmdclone.ExecuteClone(...)`.
   - If user answers `n` / `N`: Aborts safely with exit code 0.
3. **Headless Safety Gate (`--clone-on-exists`):**
   - If headless (`isInteractiveStdin() == false`):
     - If `--clone-on-exists` is provided: Automatically falls back to cloning the remote repo into the destination folder.
     - If `--clone-on-exists` is NOT provided: Hard aborts with structured error `E1089` and refusal message.
4. **Permanent Removal of Force-Push:** The `push --force` instruction in `repo_create_remote.go` is completely removed.

```
+----------------------------------------------------------------------------------------------------+
|                                SMART REPO CREATE DECISION MATRIX                                   |
+----------------------------------------------------------------------------------------------------+
|                                                                                                    |
|                                    gitmap repo create <slug>                                       |
|                                                 |                                                  |
|                                                 v                                                  |
|                                    [Remote Probe: git ls-remote]                                   |
|                                                 |                                                  |
|                        +------------------------+------------------------+                         |
|                        |                                                 |                         |
|                 Does NOT Exist                                        Exists                       |
|                        |                                                 |                         |
|                        v                                                 v                         |
|              [Standard Creation]                               [Interactive Check]                 |
|              - git init                                                  |                         |
|              - gh repo create                           +----------------+----------------+        |
|              - Initial push                             |                                 |        |
|                                                    Interactive                         Headless    |
|                                                         |                                 |        |
|                                                         v                                 v        |
|                                                Prompt: "Clone instead?"         [--clone-on-exists?]
|                                                    /         \                       /         \   |
|                                                  [Y]         [N]                   [Yes]       [No]|
|                                                   |           |                      |          |  |
|                                                   v           v                      v          v  |
|                                              [cmdclone]    [Abort]              [cmdclone]   [E1089]
+----------------------------------------------------------------------------------------------------+
```

### 4.3 CLI Command Contract
```bash
gitmap repo create <name> [folder] [slug] [flags]
# Shorthand aliases:
gitmap cr <name> [folder] [slug] [flags]
gitmap repoc <name> [folder] [slug] [flags]
```

#### Flags
| Flag | Shorthand | Type | Default | Description |
| :--- | :--- | :--- | :--- | :--- |
| `--clone-on-exists`| | bool | `false` | Headless safety gate: clones existing remote repo instead of failing if it already exists. |
| `--private` | `-p` | bool | `true` | Creates private repository on GitHub (default). |
| `--public` | | bool | `false` | Creates public repository on GitHub. |
| `--description <d>`| `-m <d>` | string | `""` | Sets remote repository description. |
| `--no-sync` | | bool | `false` | Skips automatic workspace synchronization. |

### 4.4 Data Structures & Types (`cli/cmd/repo_create_remote.go`)

```go
package cmd

// RemoteProbeStatus defines remote existence check states.
type RemoteProbeStatus string

const (
	RemoteStatusNotFound RemoteProbeStatus = "not_found"
	RemoteStatusExists   RemoteProbeStatus = "exists"
	RemoteStatusAuthFail RemoteProbeStatus = "auth_failed"
)

// RemoteProbeResult contains remote repository discovery details.
type RemoteProbeResult struct {
	Status        RemoteProbeStatus `json:"status"`
	RemoteURL     string            `json:"remoteUrl"`
	DefaultBranch string            `json:"defaultBranch,omitempty"`
	ProbeDuration int64             `json:"probeDurationMs"`
}

// RepoCreateOptions encapsulates creation flags and parameters.
type RepoCreateOptions struct {
	Name            string `json:"name"`
	LocalDir        string `json:"localDir"`
	Slug            string `json:"slug"`
	IsPrivate       bool   `json:"isPrivate"`
	IsCloneOnExists bool   `json:"isCloneOnExists"`
	Description     string `json:"description"`
	IsNoSync        bool   `json:"isNoSync"`
}
```

### 4.5 Error Codes
| Error Code | Constant | Trigger Condition | Exit Code |
| :--- | :--- | :--- | :--- |
| `E1089` | `ErrRemoteRepoExists` | Remote repository already exists and `--clone-on-exists` not set in headless environment | 1 |
| `E1090` | `ErrRemoteProbeFailed` | Network or remote host error while probing repository | 1 |
| `E1091` | `ErrGitHubAuthMissing`| GitHub CLI (`gh`) or Git credentials missing | 1 |

---

## 5. Subsystem 4: AGM Cluster Fleet Deploy with Secure Zero-Fill Purge (`gitmap agm deploy`)

### 5.1 Overview & In-Transit Cryptographic Transport
GitMap cluster management coordinates worker nodes across local and remote networks. Deploying Antigravity Manager credentials (`~/.antigravity_tools/accounts/*.json` and `accounts.json`) to remote machines requires strict defense-in-depth:
1. **AES-256-GCM In-Transit Encryption:**
   - Account artifacts are packaged into an in-memory `tar.gz` archive.
   - An ephemeral 32-byte cryptographically secure session key is generated (`crypto/rand`).
   - The archive is encrypted using AES-256-GCM via `cli/crypto.Encrypt`.
   - The encrypted payload is transmitted over authenticated SSH connections to target cluster nodes (`cli/cmdnodes/nodes_deploy_agm.go`).
2. **Remote Deployment & Cryptographic Receipt:**
   - Target nodes receive the payload, decrypt it into memory, and unpack it to destination paths.
   - Remote node computes a SHA-256 checksum of the deployed files and returns a verified acknowledgment to the coordinator.

### 5.2 Local Zero-Fill Purge & Shredding Algorithm
When deploying credentials from an ephemeral staging workstation or public laptop, operators require automatic, verified deletion of local credentials once fleet nodes are live.

#### The 3-Pass Cryptographic Shred Protocol
If `--purge-accounts` (or `--shred`) is supplied:
1. **Pre-Purge Verification Invariant:** The purge procedure **ONLY** begins if 100% of target cluster nodes report `status = "success"` and verified SHA-256 receipts. If even a single node fails, local account files are preserved and `E1092` is thrown.
2. **File Discovery:** Locates:
   - `~/.antigravity_tools/accounts/*.json`
   - `~/.antigravity_tools/accounts.json`
3. **Execution Passes:**
   - **Pass 1 (Zero-Fill):** Reads file length $S$. Overwrites file content with $S$ bytes of `0x00`. Flushes to storage via `file.Sync()`.
   - **Pass 2 (CSPRNG Shredding - when `--shred`):** Generates $S$ bytes of cryptographic pseudo-random noise (`crypto/rand.Read`) and writes across the file. Flushes via `file.Sync()`.
   - **Pass 3 (Final Zero-Fill):** Overwrites $S$ bytes with `0x00`. Flushes via `file.Sync()`.
   - **Truncate & Remove:** Truncates file length to 0 bytes via `file.Truncate(0)`, closes file descriptor, and unlinks path via `os.Remove()`.

```
+----------------------------------------------------------------------------------------------------+
|                                  CRYPTOGRAPHIC SHRED LIFECYCLE                                     |
+----------------------------------------------------------------------------------------------------+
|                                                                                                    |
|    Local File: accounts.json [Size: S bytes]                                                       |
|                                                                                                    |
|    Pass 1: Zero-Fill (0x00 * S)                  file.Sync() [OS Disk Flush]                       |
|    --------------------------------------------> [Disk Blocks Overwritten with 0x00]               |
|                                                                                                    |
|    Pass 2: CSPRNG Shred (crypto/rand * S)        file.Sync() [OS Disk Flush]                       |
|    --------------------------------------------> [Disk Blocks Overwritten with Random Data]        |
|                                                                                                    |
|    Pass 3: Final Zero-Fill (0x00 * S)            file.Sync() [OS Disk Flush]                       |
|    --------------------------------------------> [Disk Blocks Overwritten with 0x00]               |
|                                                                                                    |
|    Truncation & Removal                          os.Remove(filePath)                               |
|    --------------------------------------------> [File Truncated to 0 & Removed from Filesystem]   |
+----------------------------------------------------------------------------------------------------+
```

### 5.3 CLI Command Contract
```bash
gitmap agm deploy [target-node] [flags]
# Shorthand alias:
gitmap nodes deploy-agm [target-node] [flags]
```

#### Flags
| Flag | Shorthand | Type | Default | Description |
| :--- | :--- | :--- | :--- | :--- |
| `--purge-accounts`| | bool | `false` | Wipes local account files after 100% verified deployment across cluster nodes. |
| `--shred` | | bool | `false` | Enforces 3-pass DoD/CSPRNG zero-fill shredding prior to unlinking local account files. |
| `--except <list>` | | string | `""` | Comma-separated list of node aliases to exclude from deployment. |
| `--open-only` | | bool | `false` | Deploys only to reachable/responsive nodes. |
| `--include-main` | | bool | `false` | Deploys credentials to main cluster controller as well. |
| `--dry-run` | `-n` | bool | `false` | Simulates deployment and shredding without transferring or deleting files. |
| `--json` | `-j` | bool | `false` | Returns structured JSON deployment and purge telemetry. |

### 5.4 Data Structures & Types (`cli/cmdnodes/nodes_shred.go`)

```go
package cmdnodes

import "time"

// ShredAuditEntry records zero-fill shred details for a specific file.
type ShredAuditEntry struct {
	FilePath        string    `json:"filePath"`
	OriginalBytes   int64     `json:"originalBytes"`
	PassCount       int       `json:"passCount"`
	IsVerifiedZero  bool      `json:"isVerifiedZero"`
	CompletedAt     time.Time `json:"completedAt"`
}

// AGMDeployReport encapsulates overall deploy and purge execution results.
type AGMDeployReport struct {
	SessionId      string            `json:"sessionId"`
	TotalNodes     int               `json:"totalNodes"`
	SucceededNodes int               `json:"succeededNodes"`
	FailedNodes    int               `json:"failedNodes"`
	IsPurgeRun     bool              `json:"isPurgeRun"`
	ShreddedFiles  []ShredAuditEntry `json:"shreddedFiles,omitempty"`
	DurationMs     int64             `json:"durationMs"`
}
```

### 5.5 Error Codes
| Error Code | Constant | Trigger Condition | Exit Code |
| :--- | :--- | :--- | :--- |
| `E1092` | `ErrFleetDeployIncomplete`| One or more nodes failed deployment; local account purge aborted | 1 |
| `E1093` | `ErrShredIOFailure`       | File permission or disk I/O failure during shredding | 1 |
| `E1094` | `ErrEncryptionFailed`     | Session key generation or AES-256-GCM cipher initialization failed | 1 |

---

## 6. Subsystem 5: Git Commit Graph & Branch Compare Cache (`gitmap log` & `gitmap branch compare`)

### 6.1 Problem Statement & 0ms Fast Path Strategy
Repeatedly running `git log` or calculating ahead/behind branch divergences in multi-agent environments creates severe bottlenecks:
- Subshell spawning (`exec.Command("git", ...)`) incurs 50–150ms of process start latency on Windows.
- Running branch graph diffs (`git rev-list --left-right branchA...branchB`) requires walking thousands of packfile objects repeatedly.

To deliver **sub-millisecond (0ms rounded) fast path performance**, GitMap introduces an isolated Split SQLite commit graph cache located in:
```text
.gitmap/data/git-cache/<slug>/sql.db
```

#### Fast-Path Cache Verification
1. Reads current commit hashes for target branches directly from `.git/refs/heads/<branch>` (or `.git/packed-refs`) in < 1ms.
2. Forms a unique cache composite key: `<branchA>...<branchB>:<hashA>:<hashB>`.
3. If entry exists in `GitBranchCompareCache` table: Returns parsed ahead/behind counts and commit lists in **< 1ms**.
4. If cache miss occurs: Walks git rev-list once, populates `GitCommit` and `GitBranchCompareCache` tables in SQLite, and serves response.

### 6.2 SQLite Split-DB Schemas (`.gitmap/data/git-cache/<slug>/sql.db`)

```sql
-- Git Commit Cache Table
CREATE TABLE IF NOT EXISTS GitCommit (
    CommitHash TEXT PRIMARY KEY,
    AuthorName TEXT NOT NULL,
    AuthorEmail TEXT NOT NULL,
    AuthorDate DATETIME NOT NULL,
    CommitterName TEXT NOT NULL,
    CommitterDate DATETIME NOT NULL,
    CommitMessage TEXT NOT NULL,
    ParentHashes TEXT, -- Comma-separated parent commit hashes
    TreeHash TEXT NOT NULL,
    IsMergeCommit INTEGER NOT NULL DEFAULT 0,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_git_commit_author_date ON GitCommit(AuthorDate DESC);

-- Git Branch Reference Tracking
CREATE TABLE IF NOT EXISTS GitRef (
    RefName TEXT PRIMARY KEY,
    CommitHash TEXT NOT NULL,
    UpdatedAt DATETIME NOT NULL
);

-- Branch Compare Ahead/Behind Result Cache
CREATE TABLE IF NOT EXISTS GitBranchCompareCache (
    CompareKey TEXT PRIMARY KEY, -- "<branchA>...<branchB>:<hashA>:<hashB>"
    BaseBranch TEXT NOT NULL,
    TargetBranch TEXT NOT NULL,
    AheadCount INTEGER NOT NULL,
    BehindCount INTEGER NOT NULL,
    AheadHashes TEXT NOT NULL,   -- JSON array of commit hashes
    BehindHashes TEXT NOT NULL,  -- JSON array of commit hashes
    CommonMergeBase TEXT NOT NULL,
    CachedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_git_compare_branches ON GitBranchCompareCache(BaseBranch, TargetBranch);
```

### 6.3 CLI Command Contracts

#### 1. Cached Commit Log (`gitmap log`)
```bash
gitmap log [ref] [flags]
```
| Flag | Shorthand | Type | Default | Description |
| :--- | :--- | :--- | :--- | :--- |
| `--cached` | `-c` | bool | `true` | Queries the local SQLite commit cache (falls back to git on miss). |
| `--max <N>`| `-n <N>` | int | `20` | Limits number of displayed commits. |
| `--author <a>` | | string | `""` | Filters commits by author name or email. |
| `--json` | `-j` | bool | `false` | Emits structured JSON commit records. |

#### 2. Branch Ahead/Behind Comparison (`gitmap branch compare` / `gitmap cmp`)
```bash
gitmap branch compare <branchA> <branchB> [flags]
# Shorthand aliases:
gitmap branch cmp <branchA> <branchB> [flags]
gitmap cmp <branchA> <branchB> [flags]
```
| Flag | Shorthand | Type | Default | Description |
| :--- | :--- | :--- | :--- | :--- |
| `--refresh` | `-r` | bool | `false` | Forces re-calculation and updates SQLite cache. |
| `--json` | `-j` | bool | `false` | Emits structured JSON ahead/behind metrics. |

### 6.4 Terminal Output Rendering (`termtable` / `termpad`)
When rendered to an interactive terminal, `gitmap cmp` displays an aligned comparison box:

```text
┌────────────────────────────────────────────────────────────────────────────────────────┐
│ GitMap Branch Comparison: main ... feature/mcp-server-ai-analysis                     │
├────────────────────────────────────────────────────────────────────────────────────────┤
│ Base Branch:   main                      (Commit: c3f91a2)                             │
│ Target Branch: feature/mcp-server-ai-... (Commit: 7e82b41)                             │
│ Merge Base:    c3f91a2                                                                 │
│ Divergence:    Ahead: 4 commits │ Behind: 0 commits (Fast-Forward Capable)             │
│ Cache Status:  HIT (0ms)                                                               │
├────────────────────────────────────────────────────────────────────────────────────────┤
│ Ahead Commits:                                                                         │
│   • 7e82b41 feat(mcp): add json-rpc tool router & stdio handlers (2 hours ago)         │
│   • 4b1a092 feat(ai-db): add split sqlite ai-analysis.db storage engine (3 hours ago)  │
│   • 91d0e81 feat(vault): zero-loss file removal vault & revert engine (4 hours ago)    │
│   • 1928fa0 docs(spec): ratify mcp server and ai analysis specification (5 hours ago)  │
└────────────────────────────────────────────────────────────────────────────────────────┘
```

### 6.5 Data Structures & Go Models (`cli/cmdgitcache/gitcache_types.go`)

```go
package cmdgitcache

import "time"

// GitCommitRecord represents a cached git commit in SQLite.
type GitCommitRecord struct {
	CommitHash    string    `json:"commitHash"`
	AuthorName    string    `json:"authorName"`
	AuthorEmail   string    `json:"authorEmail"`
	AuthorDate    time.Time `json:"authorDate"`
	CommitterName string    `json:"committerName"`
	CommitterDate time.Time `json:"committerDate"`
	CommitMessage string    `json:"commitMessage"`
	ParentHashes  []string  `json:"parentHashes"`
	TreeHash      string    `json:"treeHash"`
	IsMergeCommit bool      `json:"isMergeCommit"`
}

// BranchCompareResult encapsulates the ahead/behind diff between two branches.
type BranchCompareResult struct {
	BaseBranch      string            `json:"baseBranch"`
	TargetBranch    string            `json:"targetBranch"`
	BaseHash        string            `json:"baseHash"`
	TargetHash      string            `json:"targetHash"`
	CommonMergeBase string            `json:"commonMergeBase"`
	AheadCount      int               `json:"aheadCount"`
	BehindCount     int               `json:"behindCount"`
	AheadCommits    []GitCommitRecord `json:"aheadCommits"`
	BehindCommits   []GitCommitRecord `json:"behindCommits"`
	IsCacheHit      bool              `json:"isCacheHit"`
	CalculationMs   int64             `json:"calculationMs"`
}
```

### 6.6 Error Codes
| Error Code | Constant | Trigger Condition | Exit Code |
| :--- | :--- | :--- | :--- |
| `E1095` | `ErrBranchNotFound`   | Specified branch name does not exist in repository | 1 |
| `E1096` | `ErrGitCacheDBOpen`   | Failed to initialize or open `.gitmap/data/git-cache/<slug>/sql.db` | 1 |
| `E1097` | `ErrNoMergeBase`      | Branches do not share a common git history | 1 |

---

## 7. Unified CLI Command Registry & Dispatch Mapping

To maintain clean architecture under 100-line file sizing and strict cyclomatic complexity rules:
1. `cli/cmd/rootdispatch.go` routes top-level commands to modular subcommand runners.
2. All commands adhere to standardized parameter parsing and monadic or structured `apperror.AppError` return envelopes.

| Command | Canonical Handler | Package | Registered Aliases |
| :--- | :--- | :--- | :--- |
| `gitmap child-path` | `cmdchildpath.RunChildPathCommand` | `cli/cmdchildpath` | `cp`, `childpath`, `cp-find` |
| `gitmap install muse` | `cmdinstall.RunInstallMuse` | `cli/cmdinstall` | `muse install`, `install-muse` |
| `gitmap repo create` | `cmd.RunRepoCreateCommand` | `cli/cmd` | `cr`, `repoc`, `repo-create` |
| `gitmap agm deploy` | `cmdnodes.RunDeployAGM` | `cli/cmdnodes` | `nodes deploy-agm` |
| `gitmap log` | `cmdgitcache.RunCachedLogCommand` | `cli/cmdgitcache` | `log` |
| `gitmap branch compare`| `cmdgitcache.RunBranchCompareCommand`| `cli/cmdgitcache` | `branch cmp`, `cmp` |

---

## 8. Verification Strategy & Quality Gate

Every subsystem must pass automated verification suites:
1. **Unit Testing:** Individual unit tests for each package (`childpath_test.go`, `install_muse_test.go`, `repo_create_remote_test.go`, `nodes_shred_test.go`, `gitcache_test.go`).
2. **Linting & Guideline Compliance:** Zero violations of coding standards checked via `python 03-ai-scripts/05-guideline-autofixer.py --check-only`.
3. **Relative Path Hygiene:** Zero absolute filesystem paths detected via `python linter-scripts/check-relative-paths.py`.
4. **Visual Parity:** Confirmed visual consistency with `assets/screenshots/240-muse-meta-installer.png`, `assets/screenshots/240-get-childitem-filter.png`, and `assets/screenshots/240-git-grep-func-runfix.png`.
