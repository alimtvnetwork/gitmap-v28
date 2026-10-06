# Component Specification: Nodes Pending Commits & Remote Commits Delegation Suite

> **Document Version:** 1.0.0  
> **Status:** Active  
> **Scope:** Core Go Engine (`cli/cmdnodes/`, `cli/cmd/nodes_cmd.go`, `cli/cmdssh/`)  
> **Traceability:** Task-232  
> **Coding Guideline Conformance:** Positive booleans exclusively (`isNodeOnline`, `isSuccess`, `hasChanges`, `isDryRun`, `isPushed`, `isDirty`, `hasUnpushedCommits`), strict relative paths.

---

## 1. Executive Summary & Component Topology

The Fleet Commits and Delegation Subsystem extends GitMap's cluster orchestration capabilities with two primary command families:
1. **Fleet Pending Commits (`gitmap nodes pc` / `gitmap nodes pending-commits`)**:
   - Queries all registered cluster SSH nodes in parallel.
   - Inspects both uncommitted working tree modifications (`git status --porcelain`) and unpushed commits (`git rev-list @{u}..HEAD`).
   - Supports summary views, detailed views (`--detail`), custom sort orders (`--sort=priority|name|count`), and node filtering.
2. **Fleet Remote Commits Suite (`gitmap nodes commits`, `cpf`, `cpb`, `cpr`, `commit-fix`)**:
   - Dispatches atomic commit and push operations across remote cluster machines for a specified repository or all dirty repositories (`$repoName/all`).
   - Semantic commit prefixes:
     - `cpf`: Feature commits (`Feature: <msg>`)
     - `cpb`: Bug fix commits (`Bug: <msg>`)
     - `cpr`: Release commits (`Release: <msg>`)
     - `commit-fix`: Merge and repair commits (`Fix: <msg>`)
     - `commits`: Standard commits without automatic semantic prefix
   - Delegates execution to the remote machine's local GitMap binary via SSH, captures structured JSON telemetry, and aggregates results.

```
+-----------------------------------------------------------------------------------+
|                           GitMap Cluster CLI Dispatcher                           |
|                             cli/cmd/nodes_cmd.go                                  |
+-----------------------------------------+-----------------------------------------+
                                          |
                      +-------------------+-------------------+
                      |                                       |
                      v                                       v
        +----------------------------+         +----------------------------+
        |   Fleet Pending Commits    |         |   Fleet Remote Commits     |
        |  nodes_pending_commits.go  |         |      nodes_commits.go      |
        +--------------+-------------+         +--------------+-------------+
                       |                                      |
                       +-------------------+------------------+
                                           |
                                           v
                        +-------------------------------------+
                        |     Fleet Node Filter & Planner     |
                        |       cli/cmdnodes/nodes_filter.go   |
                        +------------------+------------------+
                                           |
                                           v
                        +-------------------------------------+
                        |     SSH Execution & Transport       |
                        |      cli/cmdssh/ & cli/crypto/      |
                        +------------------+------------------+
                                           |
                   +-----------------------+-----------------------+
                   | (Parallel SSH)                                | (Parallel SSH)
                   v                                               v
        +----------------------+                       +----------------------+
        | Node 01: Linux Host  |                       | Node 02: Windows Host|
        | Shell: bash / sh     |                       | Shell: powershell.exe|
        | Workspace: ~/work/   |                       | Workspace: work/     |
        | Executes: gitmap ... |                       | Executes: gitmap ... |
        +----------+-----------+                       +----------+-----------+
                   |                                               |
                   | (stdout with potential MOTD / shell noise)    |
                   v                                               v
        +---------------------------------------------------------------------+
        |                     Robust JSON Extraction Engine                   |
        |            extractJSONPayload() & JSON Envelope Decoder             |
        +----------------------------------+----------------------------------+
                                           |
                   +-----------------------+-----------------------+
                   | (Human Output)                                | (Machine Output)
                   v                                               v
        +----------------------+                       +----------------------+
        |  termtable Renderer  |                       | Structured JSON Out  |
        |   Summary & Badges   |                       |  --json / -j Flag    |
        +----------------------+                       +----------------------+
```

---

## 2. Command Interface & Routing Integration

### 2.1 Fleet Pending Commits Routing
- **Command**: `gitmap nodes pending-commits`
- **Aliases**: `gitmap nodes pc`, `gitmap nodes pendingcommits`
- **Dispatcher Location**: `cli/cmd/nodes_cmd.go` inside `runUnifiedNodesCLI(args []string)`
- **Function Signature**: `cmdnodes.RunNodesPendingCommits(args []string) error`

#### Supported Flags:
| Flag | Short | Default | Description |
| :--- | :--- | :--- | :--- |
| `--target <node>` | `-t` | `""` | Target single node alias or IP address |
| `--except <nodes>`| `-e` | `""` | Comma-separated list of nodes to exclude |
| `--include <nodes>`| | `""` | Comma-separated list of nodes to whitelist |
| `--include-main` | | `false` | Include main control node |
| `--open-only` | | `false` | Probe SSH liveness before running queries |
| `--sort <mode>` | `-s` | `"priority"`| Sort repositories by `priority`, `name`, or `count` |
| `--detail` | `-d` | `false` | Display granular per-file change breakdown |
| `--json` | `-j` | `false` | Emit raw aggregated JSON telemetry to stdout |
| `--help` | `-h` | `false` | Display rich terminal UI help box |

### 2.2 Fleet Remote Commits Routing
- **Commands**:
  - `gitmap nodes commits <repoName|all> "<msg>" [flags]`
  - `gitmap nodes cpf <repoName|all> "<msg>" [flags]`
  - `gitmap nodes cpb <repoName|all> "<msg>" [flags]`
  - `gitmap nodes cpr <repoName|all> "<msg>" [flags]`
  - `gitmap nodes commit-fix <repoName|all> "<msg>" [flags]`
- **Dispatcher Location**: `cli/cmd/nodes_cmd.go` inside `runUnifiedNodesCLI(args []string)`
- **Function Signatures**:
  - `cmdnodes.RunNodesCommits(args []string) error`
  - `cmdnodes.RunNodesCommitPushFeature(args []string) error`
  - `cmdnodes.RunNodesCommitPushBug(args []string) error`
  - `cmdnodes.RunNodesCommitPushRelease(args []string) error`
  - `cmdnodes.RunNodesCommitFix(args []string) error`
  - `cmdnodes.RunNodesCommitsWithAction(action string, args []string) error`

#### Positional Argument Validation:
1. `args[0]`: `<repoName|all>` - Target repository name, directory slug, or keyword `all` for all dirty repositories on the node.
2. `args[1]`: `"<msg>"` - Commit message string. Required unless `--dry-run` is active or help is requested.

#### Supported Flags:
| Flag | Short | Default | Description |
| :--- | :--- | :--- | :--- |
| `--target <node>` | `-t` | `""` | Filter to single node alias or IP address |
| `--except <nodes>`| `-e` | `""` | Comma-separated node aliases to exclude |
| `--include <nodes>`| | `""` | Comma-separated node aliases to whitelist |
| `--include-main` | | `false` | Include main control node |
| `--open-only` | | `false` | Probe SSH connection before dispatching |
| `--dry-run` | `-n` | `false` | Simulate staging and committing without making modifications |
| `--no-push` | | `false` | Stage and commit locally on remote node without running `git push` |
| `--json` | `-j` | `false` | Emit aggregated JSON telemetry |
| `--help` | `-h` | `false` | Display rich terminal UI help box |

---

## 3. Data Models & Type Definitions

All structs reside in `cli/cmdnodes/nodes_commits_types.go` or `cli/cmdnodes/nodes_pending_commits.go` adhering strictly to positive boolean nomenclature.

```go
package cmdnodes

import "time"

// NodePendingRepoItem represents a single repository's pending state on a remote node.
type NodePendingRepoItem struct {
	RepoName           string   `json:"repo_name"`
	RelativePath       string   `json:"relative_path"`
	Branch             string   `json:"branch"`
	IsClean            bool     `json:"is_clean"`
	HasChanges         bool     `json:"has_changes"`
	DirtyFileCount     int      `json:"dirty_file_count"`
	HasUnpushedCommits bool     `json:"has_unpushed_commits"`
	UnpushedCount      int      `json:"unpushed_count"`
	PriorityScore      int      `json:"priority_score"`
	ChangedFiles       []string `json:"changed_files,omitempty"`
}

// NodePendingCommitsResult represents the pending commits discovery outcome for a single node.
type NodePendingCommitsResult struct {
	NodeAlias       string                `json:"node_alias"`
	Host            string                `json:"host"`
	OS              string                `json:"os"`
	IsNodeOnline    bool                  `json:"is_node_online"`
	IsSuccess       bool                  `json:"is_success"`
	Latency         time.Duration         `json:"latency"`
	LatencyMs       int64                 `json:"latency_ms"`
	PendingRepos    []NodePendingRepoItem `json:"pending_repos,omitempty"`
	TotalDirtyFiles int                   `json:"total_dirty_files"`
	TotalUnpushed   int                   `json:"total_unpushed"`
	ErrorMessage    string                `json:"error_message,omitempty"`
}

// NodesPendingCommitsAggregated is the top-level envelope for gitmap nodes pc --json.
type NodesPendingCommitsAggregated struct {
	Timestamp          string                     `json:"timestamp"`
	TotalNodesQueried  int                        `json:"total_nodes_queried"`
	OnlineNodesCount   int                        `json:"online_nodes_count"`
	TotalPendingRepos  int                        `json:"total_pending_repos"`
	TotalDirtyFiles    int                        `json:"total_dirty_files"`
	TotalUnpushedCount int                        `json:"total_unpushed_count"`
	SortMode           string                     `json:"sort_mode"`
	Results            []NodePendingCommitsResult `json:"results"`
}

// CommitActionKind defines semantic action types.
type CommitActionKind string

const (
	CommitActionStandard CommitActionKind = "commits"
	CommitActionFeature  CommitActionKind = "cpf"
	CommitActionBug      CommitActionKind = "cpb"
	CommitActionRelease  CommitActionKind = "cpr"
	CommitActionFix      CommitActionKind = "commit-fix"
)

// RemoteRepoCommitOutcome captures the commit result for an individual repository on a node.
type RemoteRepoCommitOutcome struct {
	RepoName     string `json:"repo_name"`
	Branch       string `json:"branch"`
	CommitHash   string `json:"commit_hash,omitempty"`
	CommitMsg    string `json:"commit_msg"`
	IsSuccess    bool   `json:"is_success"`
	IsPushed     bool   `json:"is_pushed"`
	HasChanges   bool   `json:"has_changes"`
	FilesChanged int    `json:"files_changed"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// NodeCommitResult captures the fleet commit dispatch outcome for a single node.
type NodeCommitResult struct {
	NodeAlias       string                    `json:"node_alias"`
	Host            string                    `json:"host"`
	OS              string                    `json:"os"`
	Action          string                    `json:"action"`
	IsNodeOnline    bool                      `json:"is_node_online"`
	IsSuccess       bool                      `json:"is_success"`
	IsDryRun        bool                      `json:"is_dry_run"`
	Latency         time.Duration             `json:"latency"`
	LatencyMs       int64                     `json:"latency_ms"`
	Repositories    []RemoteRepoCommitOutcome `json:"repositories,omitempty"`
	CommittedCount  int                       `json:"committed_count"`
	ErrorMessage    string                    `json:"error_message,omitempty"`
}

// NodesCommitsAggregated is the top-level envelope for gitmap nodes commits/cpf/... --json.
type NodesCommitsAggregated struct {
	Timestamp         string             `json:"timestamp"`
	Action            string             `json:"action"`
	TargetRepo        string             `json:"target_repo"`
	CommitMessage     string             `json:"commit_message"`
	IsDryRun          bool               `json:"is_dry_run"`
	TotalNodesTarget  int                `json:"total_nodes_target"`
	SuccessfulNodes   int                `json:"successful_nodes"`
	TotalReposUpdated int                `json:"total_repos_updated"`
	Results           []NodeCommitResult `json:"results"`
}
```

---

## 4. Remote SSH Command Construction & Shell Architecture

### 4.1 Cross-Platform Command Delegation
GitMap executes remote fleet operations via SSH by connecting to the target machine and delegating to the target machine's installed `gitmap` binary.

#### Linux / POSIX Shell Architecture:
- Default remote shell: `bash` with automatic fallback to `sh`.
- Default work workspace: `~/work` or active GitMap workspace directory.
- Positional argument quoting: All messages, targets, and repository names must be safely single-quoted or escaped to prevent shell injection and whitespace splitting.

```go
func buildRemoteBashCommand(subCmd string, passArgs []string) (string, string) {
	var parts []string
	parts = append(parts, "gitmap", subCmd)
	for _, arg := range passArgs {
		parts = append(parts, quoteBashArg(arg))
	}
	return strings.Join(parts, " "), "bash"
}

func quoteBashArg(s string) string {
	if s == "" {
		return "''"
	}
	// Wrap in single quotes, escaping internal single quotes: ' -> '\''
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
```

#### Windows PowerShell Architecture:
- Remote shell: `powershell.exe -NoProfile -Command "..."` using shell tag `"ps"`.
- Workspace: Configured Windows workspace (e.g., `work\`).
- Argument escaping: Wrap arguments and escape inner double quotes with backticks or double quotation marks.

```go
func buildRemotePowerShellCommand(subCmd string, passArgs []string) (string, string) {
	var parts []string
	parts = append(parts, "gitmap", subCmd)
	for _, arg := range passArgs {
		parts = append(parts, quotePowerShellArg(arg))
	}
	innerCmd := strings.Join(parts, " ")
	return fmt.Sprintf(`powershell.exe -NoProfile -Command "%s"`, innerCmd), "ps"
}

func quotePowerShellArg(s string) string {
	if s == "" {
		return `\"\"`
	}
	escaped := strings.ReplaceAll(s, `"`, `\"`)
	return `\"` + escaped + `\"`
}
```

### 4.2 Mapping to Remote Commands
1. **For `gitmap nodes pc`**:
   - The coordinator executes on remote host: `gitmap pending-commits --json` (with passed sort and filter flags).
   - This reuses Worker 01's local `gitmap pending-commits` implementation on each remote host!
2. **For `gitmap nodes <action> <repoName|all> "<msg>"`**:
   - The coordinator executes on remote host: `gitmap sends <action> <repoName|all> "<msg>" --json` (with `--dry-run` or `--no-push` if passed).
   - This directly delegates to Worker 01's local `gitmap sends` dispatcher on the remote node.

---

## 5. Remote JSON Extraction & Fleet Aggregation Engine

### 5.1 The `extractJSONPayload` Algorithm
When running commands over SSH, remote outputs frequently contain extraneous text:
- SSH login banners (`Welcome to Ubuntu...`)
- Message of the Day (MOTD) notices
- Shell configuration warnings or zsh plugin messages
- PowerShell profile loading notices

To guarantee clean JSON deserialization without syntax errors, the extractor scans for valid JSON object `{...}` or array `[...]` boundaries:

```go
// extractJSONPayload locates and extracts the outermost valid JSON object or array from raw output.
func extractJSONPayload(s string) string {
	trimmed := strings.TrimSpace(s)
	if len(trimmed) == 0 {
		return ""
	}

	objStart := strings.Index(trimmed, "{")
	objEnd := strings.LastIndex(trimmed, "}")

	arrStart := strings.Index(trimmed, "[")
	arrEnd := strings.LastIndex(trimmed, "]")

	// Determine if an object or array appears first and is bounded
	hasObj := objStart >= 0 && objEnd > objStart
	hasArr := arrStart >= 0 && arrEnd > arrStart

	if hasObj && hasArr {
		if objStart < arrStart && objEnd > arrEnd {
			return trimmed[objStart : objEnd+1]
		}
		if arrStart < objStart && arrEnd > objEnd {
			return trimmed[arrStart : arrEnd+1]
		}
		if objStart < arrStart {
			return trimmed[objStart : objEnd+1]
		}
		return trimmed[arrStart : arrEnd+1]
	}

	if hasObj {
		return trimmed[objStart : objEnd+1]
	}

	if hasArr {
		return trimmed[arrStart : arrEnd+1]
	}

	return trimmed
}
```

### 5.2 Error Handling & Resilience
- **Offline Nodes**: If SSH dial fails with connection refused, timeout, or no route to host:
  - `IsNodeOnline = false`
  - `IsSuccess = false`
  - `ErrorMessage = "host unreachable or SSH port closed"`
- **Execution Errors**: If remote GitMap exits non-zero:
  - If output contains extractable JSON error envelope, parse `ErrorMessage` from JSON.
  - If output is raw text, truncate to 200 characters and report as execution failure.
- **Zero Panic Guarantee**: Concurrency is managed via `sync.WaitGroup` with defensive recovery wrappers per goroutine.

---

## 6. Terminal Table Formatting via `termtable` vs Raw JSON Output

### 6.1 `termtable.TableConfig` Layouts

#### Table Layout 1: Fleet Pending Commits (`gitmap nodes pc`)
```go
func renderNodesPendingCommitsTable(results []NodePendingCommitsResult) {
	columns := []termtable.Column{
		{Title: "NODE ALIAS", Align: termtable.AlignLeft, MinWidth: 14},
		{Title: "REPOSITORY", Align: termtable.AlignLeft, MinWidth: 22},
		{Title: "BRANCH", Align: termtable.AlignLeft, MinWidth: 12},
		{Title: "DIRTY FILES", Align: termtable.AlignRight, MinWidth: 12},
		{Title: "UNPUSHED", Align: termtable.AlignRight, MinWidth: 10},
		{Title: "STATUS", Align: termtable.AlignLeft, MinWidth: 14},
	}

	var rows []termtable.Row
	for _, res := range results {
		if !res.IsNodeOnline {
			rows = append(rows, termtable.Row{
				res.NodeAlias, "-", "-", "-", "-", constants.ColorDim + "○ OFFLINE" + constants.ColorReset,
			})
			continue
		}
		if len(res.PendingRepos) == 0 {
			rows = append(rows, termtable.Row{
				res.NodeAlias, "(all clean)", "-", "0", "0", constants.ColorGreen + "● SYNCED" + constants.ColorReset,
			})
			continue
		}
		for _, repo := range res.PendingRepos {
			statusBadge := formatPendingBadge(repo.DirtyFileCount, repo.UnpushedCount)
			rows = append(rows, termtable.Row{
				res.NodeAlias,
				repo.RepoName,
				repo.Branch,
				fmt.Sprintf("%d", repo.DirtyFileCount),
				fmt.Sprintf("%d", repo.UnpushedCount),
				statusBadge,
			})
		}
	}

	tableCfg := termtable.TableConfig{
		Columns: columns,
		Rows:    rows,
	}
	termtable.PrintTable(tableCfg)
}
```

#### Table Layout 2: Fleet Commits Execution (`gitmap nodes cpf / cpb / ...`)
```go
func renderNodesCommitsTable(results []NodeCommitResult, action string) {
	columns := []termtable.Column{
		{Title: "NODE ALIAS", Align: termtable.AlignLeft, MinWidth: 14},
		{Title: "REPOSITORY", Align: termtable.AlignLeft, MinWidth: 20},
		{Title: "ACTION", Align: termtable.AlignLeft, MinWidth: 10},
		{Title: "COMMIT HASH", Align: termtable.AlignLeft, MinWidth: 12},
		{Title: "STATUS", Align: termtable.AlignLeft, MinWidth: 16},
		{Title: "DURATION", Align: termtable.AlignRight, MinWidth: 10},
	}

	var rows []termtable.Row
	for _, res := range results {
		if !res.IsNodeOnline {
			rows = append(rows, termtable.Row{
				res.NodeAlias, "-", strings.ToUpper(action), "-", constants.ColorDim + "○ OFFLINE" + constants.ColorReset, "-",
			})
			continue
		}
		for _, repo := range res.Repositories {
			statusStr := formatCommitStatusBadge(repo.IsSuccess, repo.HasChanges, repo.IsPushed, res.IsDryRun)
			rows = append(rows, termtable.Row{
				res.NodeAlias,
				repo.RepoName,
				strings.ToUpper(action),
				shortHash(repo.CommitHash),
				statusStr,
				res.Latency.Round(time.Millisecond).String(),
			})
		}
	}

	tableCfg := termtable.TableConfig{
		Columns: columns,
		Rows:    rows,
	}
	termtable.PrintTable(tableCfg)
}
```

### 6.2 Raw JSON Output Schema (`--json`)
When `--json` or `-j` is passed, the engine suppresses all human-readable headers, banners, and table output, writing indented JSON (2 spaces) directly to `os.Stdout`:

```json
{
  "timestamp": "2026-10-06T18:30:00Z",
  "action": "cpf",
  "target_repo": "gitmap",
  "commit_message": "Feature: add cluster node commit delegation suite",
  "is_dry_run": false,
  "total_nodes_target": 2,
  "successful_nodes": 2,
  "total_repos_updated": 2,
  "results": [
    {
      "node_alias": "u1",
      "host": "192.168.1.101",
      "os": "linux",
      "action": "cpf",
      "is_node_online": true,
      "is_success": true,
      "is_dry_run": false,
      "latency_ms": 342,
      "repositories": [
        {
          "repo_name": "gitmap",
          "branch": "main",
          "commit_hash": "a1b2c3d4",
          "commit_msg": "Feature: add cluster node commit delegation suite",
          "is_success": true,
          "is_pushed": true,
          "has_changes": true,
          "files_changed": 4
        }
      ],
      "committed_count": 1
    }
  ]
}
```

---

## 7. Rich Terminal UI Box Help Text

Help menus must use GitMap's terminal styling standards (`constants.ColorCyan`, `constants.ColorGreen`, `constants.ColorReset`, and Unicode box drawing characters).

```
+-----------------------------------------------------------------------------------+
|  gitmap nodes pending-commits (alias: pc)                                         |
|  Inspect uncommitted and unpushed changes across all cluster fleet nodes via SSH  |
+-----------------------------------------------------------------------------------+

USAGE:
  gitmap nodes pending-commits [flags]
  gitmap nodes pc [flags]

FLAGS:
  -t, --target <alias|ip>    Target specific node alias or IP address
  -e, --except <nodes>       Exclude comma-separated node aliases
      --include <nodes>      Whitelist comma-separated node aliases
      --include-main         Include main controller node in query
      --open-only            Only query nodes currently reachable
  -s, --sort <mode>          Sort repositories: priority (default), name, count
  -d, --detail               Show granular 1-by-1 file modifications
  -j, --json                 Output raw aggregated JSON telemetry
  -h, --help                 Display this help menu

EXAMPLES:
  gitmap nodes pc
  gitmap nodes pc --sort=count --detail
  gitmap nodes pc -t u1 --json
```

```
+-----------------------------------------------------------------------------------+
|  gitmap nodes <action> <repoName|all> "<msg>" [flags]                             |
|  Delegate semantic commits and pushes across remote cluster nodes via SSH         |
+-----------------------------------------------------------------------------------+

ACTIONS:
  commits       Standard commit (no semantic prefix)
  cpf           Feature commit (prepends 'Feature: ')
  cpb           Bug fix commit (prepends 'Bug: ')
  cpr           Release commit (prepends 'Release: ')
  commit-fix    Merge and conflict fix commit (prepends 'Fix: ')

TARGETS:
  <repoName>    Specific repository directory or slug on target node(s)
  all           All dirty repositories detected on target node(s)

FLAGS:
  -t, --target <alias|ip>    Target specific node alias or IP address
  -e, --except <nodes>       Exclude comma-separated node aliases
      --include <nodes>      Whitelist comma-separated node aliases
      --include-main         Include main controller node
      --open-only            Only dispatch to reachable nodes
  -n, --dry-run              Simulate commit operations without staging/pushing
      --no-push              Stage and commit on remote node without git push
  -j, --json                 Output raw aggregated JSON telemetry
  -h, --help                 Display this help menu

EXAMPLES:
  gitmap nodes cpf gitmap "add nodes commit suite"
  gitmap nodes cpb all "fix nil pointer exception in split-db indexer"
  gitmap nodes commit-fix my-service "resolve merge conflict in main branch"
  gitmap nodes commits all "sync latest submodules" --dry-run -t u1
```

---

## 8. E2E Testing Architecture & Mock SSH Runner

### 8.1 Test Isolation & Mocking Strategy
To adhere to repository testing guidelines (strictly zero real network / OS mutations during testing):
1. **Executor Interface**: Define an internal execution interface `remoteCommandExecutor`:
```go
type remoteCommandExecutor interface {
	Execute(conn db.SSHConnection, command string, shell string) (string, error)
}
```
2. **Default vs Mock Runner**:
   - Production runner delegates to `crypto.RunCommand(client, cmd, shell)`.
   - Test runner injects mock JSON payloads and simulated connection errors.

### 8.2 Comprehensive Test Scenarios (`cli/cmdnodes/nodes_commits_test.go`)
1. **Argument Parsing & Action Mapping**:
   - Verify `cpf` correctly sets `action = "cpf"`.
   - Verify `cpb`, `cpr`, `commit-fix`, and `commits` map accurately.
   - Verify missing commit message returns descriptive `apperror.AppError` with error code `E9080`.
2. **POSIX vs Windows Command Formatting**:
   - Verify Linux nodes produce `bash` command with single-quoted arguments.
   - Verify Windows nodes produce `powershell.exe -NoProfile -Command "..."` with escaped double quotes.
3. **`extractJSONPayload` Resilience**:
   - Input: MOTD banner followed by `{ "status": "ok" }`.
   - Input: Array `[{"repo": "a"}]` surrounded by shell greeting and logout warning.
   - Input: Clean JSON string without extra characters.
   - Input: Malformed string without braces (returns error without panic).
4. **Cross-Node Parallel Aggregation**:
   - Simulate Node 1 returning 2 dirty repos, Node 2 returning 0 dirty repos (clean), Node 3 timing out.
   - Assert aggregated counts: `TotalPendingRepos == 2`, `OnlineNodesCount == 2`.
5. **Dry-Run Flag Propagation**:
   - Verify that passing `-n` or `--dry-run` appends `--dry-run` to the remote command.
6. **No-Push Flag Propagation**:
   - Verify that passing `--no-push` appends `--no-push` to the remote command.
7. **JSON Output Verification**:
   - Capture `os.Stdout` in test and assert valid JSON unmarshalling into `NodesCommitsAggregated` and `NodesPendingCommitsAggregated`.

---

## 9. Non-Negotiable Coding Standards Conformance

1. **Path Representation**:
   - Every file path in documentation and code must be strictly relative (`cli/...`, `02-spec/...`, `.ai-memory/...`).
   - Remote target locations are defined as standard work workspace references (`~/work` or `work/`).
2. **Boolean Conventions**:
   - Positive booleans exclusively: `isNodeOnline`, `isSuccess`, `hasChanges`, `isDryRun`, `isPushed`, `isDirty`, `hasUnpushedCommits`.
   - Absolute prohibition against negative booleans (`isNotOnline`, `noPush` as a struct boolean field, `hasNoChanges`).
3. **Zero Test/Build Executions**:
   - All spec and plan artifacts are authored with zero intermediate `go build` or `go test` invocations.
