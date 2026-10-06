# Subtask 02: Nodes Pending Commits & Remote Commits Delegation Suite

> **Subtask ID:** Subtask-02  
> **Parent Plan:** `.ai-memory/plans/232-pending-commits-sends-and-nodes-commit-suite.md`  
> **Component Spec:** `02-spec/21-app/232-pending-commits-sends-and-nodes-commit-suite/02-component-spec.md`  
> **Owner:** Worker 02  
> **Status:** Queued  
> **Coding Standard Conformance:** Positive booleans exclusively (`isNodeOnline`, `isSuccess`, `hasChanges`, `isDryRun`, `isPushed`, `isDirty`, `hasUnpushedCommits`), strictly relative Git paths.

---

## 1. Objectives & Scope

Implement the cluster-wide GitMap commit discovery and remote execution delegation engine in `cli/cmdnodes/`, providing:
1. `gitmap nodes pending-commits` (alias: `gitmap nodes pc`):
   - Parallel SSH querying of uncommitted files and unpushed commits across all registered cluster fleet nodes.
   - Summary view (default) vs detailed view (`--detail`), custom sort orders (`--sort=priority|name|count`), and JSON telemetry (`--json`).
2. `gitmap nodes commits`, `nodes cpf`, `nodes cpb`, `nodes cpr`, `nodes commit-fix`:
   - Remote commit dispatching targeting a specific repository or all dirty repositories (`$repoName/all`).
   - Semantic commit prefixing (`cpf` -> Feature, `cpb` -> Bug, `cpr` -> Release, `commit-fix` -> Fix).
   - Delegation to remote host's `gitmap sends` command via SSH.
   - Extraction of remote JSON payloads (`extractJSONPayload`) past MOTD and shell greetings.
   - Output formatting via `termtable` or structured JSON (`--json`).
3. CLI dispatch routing integration in `cli/cmd/nodes_cmd.go`.
4. Comprehensive unit tests in `cli/cmdnodes/nodes_commits_test.go` with mock SSH runner.

---

## 2. Owned Target Files

| File Path | Operation | Responsibility |
| :--- | :--- | :--- |
| `cli/cmdnodes/nodes_commits_types.go` | NEW | Type definitions, telemetry structs, action enums with positive booleans |
| `cli/cmdnodes/nodes_pending_commits.go` | NEW | Fleet pending commits runner, extractor, parallel worker pool, table/JSON rendering |
| `cli/cmdnodes/nodes_commits.go` | NEW | Fleet remote commits runner (`commits`, `cpf`, `cpb`, `cpr`, `commit-fix`), shell quoting, dispatcher |
| `cli/cmdnodes/nodes_commits_test.go` | NEW | Unit tests covering parsing, shell quoting, JSON extraction, mock SSH execution, and rendering |
| `cli/cmd/nodes_cmd.go` | INTEGRATION EDIT | Routing hooks in `runUnifiedNodesCLI` connecting CLI subcommands to `cmdnodes` |

---

## 3. Detailed Step-by-Step Implementation Steps

### Step 1: Author Data Models (`cli/cmdnodes/nodes_commits_types.go`)
- Create struct `NodePendingRepoItem`:
  - `RepoName string`
  - `RelativePath string`
  - `Branch string`
  - `IsClean bool`
  - `HasChanges bool`
  - `DirtyFileCount int`
  - `HasUnpushedCommits bool`
  - `UnpushedCount int`
  - `PriorityScore int`
  - `ChangedFiles []string`
- Create struct `NodePendingCommitsResult`:
  - `NodeAlias string`, `Host string`, `OS string`
  - `IsNodeOnline bool`, `IsSuccess bool`
  - `Latency time.Duration`, `LatencyMs int64`
  - `PendingRepos []NodePendingRepoItem`
  - `TotalDirtyFiles int`, `TotalUnpushed int`
  - `ErrorMessage string`
- Create struct `NodesPendingCommitsAggregated`:
  - `Timestamp string`, `TotalNodesQueried int`, `OnlineNodesCount int`, `TotalPendingRepos int`, `TotalDirtyFiles int`, `TotalUnpushedCount int`, `SortMode string`, `Results []NodePendingCommitsResult`
- Create struct `RemoteRepoCommitOutcome`:
  - `RepoName string`, `Branch string`, `CommitHash string`, `CommitMsg string`
  - `IsSuccess bool`, `IsPushed bool`, `HasChanges bool`, `FilesChanged int`, `ErrorMessage string`
- Create struct `NodeCommitResult`:
  - `NodeAlias string`, `Host string`, `OS string`, `Action string`
  - `IsNodeOnline bool`, `IsSuccess bool`, `IsDryRun bool`
  - `Latency time.Duration`, `LatencyMs int64`
  - `Repositories []RemoteRepoCommitOutcome`
  - `CommittedCount int`, `ErrorMessage string`
- Create struct `NodesCommitsAggregated`:
  - `Timestamp string`, `Action string`, `TargetRepo string`, `CommitMessage string`
  - `IsDryRun bool`, `TotalNodesTarget int`, `SuccessfulNodes int`, `TotalReposUpdated int`
  - `Results []NodeCommitResult`

### Step 2: Implement Fleet Pending Commits (`cli/cmdnodes/nodes_pending_commits.go`)
- Implement `RunNodesPendingCommits(args []string) error`:
  - Check `--help` / `-h` -> `printNodesPendingCommitsHelp()`
  - Parse node filtering options using `ParseNodeFilterOptions(args)`
  - Fetch fleet SSH connections using `cmdssh.FetchAllSSHConnections()`
  - Filter targets with `FilterFleetNodes(conns, opts)`
  - Execute parallel queries across filtered nodes:
    - Launch goroutines with `sync.WaitGroup`
    - Connect via `cmdssh.ConnectSSHClientWithErr(conn)`
    - Construct remote command:
      - Linux: `gitmap pending-commits --json` (with sort/detail args)
      - Windows: `powershell.exe -NoProfile -Command "gitmap pending-commits --json"`
    - Execute via `crypto.RunCommand(client, cmd, shell)` (with bash -> sh fallback)
    - Pass raw output through `extractJSONPayload(out)`
    - Deserialize remote JSON into `NodePendingCommitsResult`
  - Sorting support:
    - `--sort=priority`: repos with most unpushed + dirty changes first
    - `--sort=name`: alphabetical by repository name
    - `--sort=count`: by dirty file count descending
  - Output handling:
    - If `--json` / `-j`: emit `NodesPendingCommitsAggregated` to stdout
    - Else: render summary table via `termtable.TableConfig` and summary badge line

### Step 3: Implement Fleet Commits Engine (`cli/cmdnodes/nodes_commits.go`)
- Implement entrypoints:
  - `RunNodesCommits(args []string) error`
  - `RunNodesCommitPushFeature(args []string) error`
  - `RunNodesCommitPushBug(args []string) error`
  - `RunNodesCommitPushRelease(args []string) error`
  - `RunNodesCommitFix(args []string) error`
  - `RunNodesCommitsWithAction(action string, args []string) error`
- Parse arguments:
  - Check `--help` / `-h` -> `printNodesCommitsHelp(action)`
  - First positional argument: `<repoName|all>`
  - Second positional argument: `"<msg>"`
  - Flags: `--dry-run` (`-n`), `--no-push`, `--target` (`-t`), `--except` (`-e`), `--include`, `--include-main`, `--open-only`, `--json` (`-j`)
  - If missing message and not `--dry-run`: return `apperror.NewSimple("commit message is required", "E9080")`
- Remote command construction:
  - Linux Bash:
    - Safely quote arguments: single-quote wrapping with inner quote escaping
    - Remote command: `gitmap sends <action> <repoName|all> '<msg>' --json` (plus `--dry-run` or `--no-push` if specified)
  - Windows PowerShell:
    - Escape inner double quotes
    - Remote command: `powershell.exe -NoProfile -Command "gitmap sends <action> <target> \"<msg>\" --json"`
- Parallel execution across fleet nodes:
  - Capture remote JSON via `extractJSONPayload`
  - Parse per-repo commit outcomes
  - Aggregate results into `NodesCommitsAggregated`
- Render outputs:
  - If `--json`: emit JSON to stdout
  - Else: render table via `termtable.TableConfig` with columns: `NODE ALIAS`, `REPOSITORY`, `ACTION`, `COMMIT HASH`, `STATUS`, `DURATION`

### Step 4: Implement Robust Extraction & Helpers
- Implement `extractJSONPayload(s string) string`:
  - Locate bounding `{...}` or `[...]`
  - Safely extract substring, handling trailing or leading shell banners
- Implement shell escaping helpers:
  - `quoteBashArg(s string) string`
  - `quotePowerShellArg(s string) string`

### Step 5: Integrate CLI Routing (`cli/cmd/nodes_cmd.go`)
- Inside `runUnifiedNodesCLI(args []string)`:
  - Add pending commits matchers:
    - `strings.EqualFold(args[0], "pending-commits") || strings.EqualFold(args[0], "pc") || strings.EqualFold(args[0], "pendingcommits")`
    - Call `cmdnodes.RunNodesPendingCommits(args[1:])`
  - Add commits matchers:
    - `strings.EqualFold(args[0], "commits")` -> `cmdnodes.RunNodesCommits(args[1:])`
    - `strings.EqualFold(args[0], "cpf")` -> `cmdnodes.RunNodesCommitPushFeature(args[1:])`
    - `strings.EqualFold(args[0], "cpb")` -> `cmdnodes.RunNodesCommitPushBug(args[1:])`
    - `strings.EqualFold(args[0], "cpr")` -> `cmdnodes.RunNodesCommitPushRelease(args[1:])`
    - `strings.EqualFold(args[0], "commit-fix") || strings.EqualFold(args[0], "commitfix")` -> `cmdnodes.RunNodesCommitFix(args[1:])`

### Step 6: Author Comprehensive Unit Tests (`cli/cmdnodes/nodes_commits_test.go`)
- Test Cases to author:
  - `TestExtractJSONPayload_Resilience`: test object and array extraction with prepended/appended MOTD strings.
  - `TestQuoteBashArg`: test argument quoting with spaces, quotes, and empty strings.
  - `TestQuotePowerShellArg`: test argument quoting for PowerShell syntax.
  - `TestParseNodesCommitsArgs`: test argument extraction for repo target, commit message, and flags.
  - `TestParseNodesCommitsArgs_MissingMessage`: assert error E9080 when message is missing.
  - `TestResolveRemoteCommitsCommand`: assert command generated for Linux vs Windows nodes.
  - `TestRenderNodesCommitsTable`: assert table rendering does not panic and contains expected column headers.
  - `TestRenderNodesPendingCommitsTable`: assert pending commits table rendering handles clean and dirty states.

---

## 4. Verification Checklist & Quality Gates

| Gate ID | Requirement | Verification Strategy |
| :--- | :--- | :--- |
| **VG-01** | Routing Parity | `gitmap nodes pc`, `nodes pending-commits`, `nodes commits`, `nodes cpf`, `nodes cpb`, `nodes cpr`, `nodes commit-fix` route without error |
| **VG-02** | Target Resolution | Correctly resolves single repo name, relative paths, or keyword `all` |
| **VG-03** | Semantic Prefixing | `cpf` maps to Feature, `cpb` to Bug, `cpr` to Release, `commit-fix` to Fix |
| **VG-04** | Remote Shell Escaping | POSIX bash quotes with single quotes; Windows uses escaped double quotes |
| **VG-05** | JSON Extractor | Extracts clean JSON even when terminal contains MOTD banners or prompt noise |
| **VG-06** | UI & JSON Output | `termtable` output matches styling; `--json` produces valid parseable JSON envelope |
| **VG-07** | Guidelines & Hygiene | 100% positive booleans, 100% relative Git paths, zero test/build execution |

---

## 5. Non-Negotiable Instructions for Worker 02
- **NO BUILDS OR TESTS**: Do NOT run `go build`, `go test`, or shell compilation commands.
- **NO GIT MUTATIONS**: Do NOT run `git add`, `git commit`, `git push`, or alter git history.
- **RELATIVE PATHS ONLY**: Never hardcode absolute filesystem paths or host machine URIs.
