# Subtask 01: `gitmap pending-commits` & `gitmap sends` Local & SSH Suite

> **Parent Plan:** `.ai-memory/plans/232-pending-commits-sends-and-nodes-commit-suite.md`  
> **Spec Reference:** `02-spec/21-app/232-pending-commits-sends-and-nodes-commit-suite/01-architecture-spec.md`  
> **Component Spec Reference:** `02-spec/21-app/232-pending-commits-sends-and-nodes-commit-suite/02-component-spec.md`  
> **Status:** `PENDING`  
> **Worker Assignment:** Worker 01  
> **Target Subsystems:** CLI Command Dispatcher, Status Inspection Engine, Multi-Repository Semantic Commit Dispatcher, Fleet SSH Aggregator  
> **Relative Affected Paths:**  
> - `cli/constants/constants_cli.go`  
> - `cli/cmd/rootcore.go`  
> - `cli/cmd/roottooling.go`  
> - `cli/cmd/pending_commits_cmd.go`  
> - `cli/cmd/pending_commits_types.go`  
> - `cli/cmd/pending_commits_help.go`  
> - `cli/cmd/pending_commits_test.go`  
> - `cli/cmd/sends_cmd.go`  
> - `cli/cmd/sends_types.go`  
> - `cli/cmd/sends_help.go`  
> - `cli/cmd/sends_test.go`  
> **Execution Constraint:** Pure Specification & Subtask Plan Authoring (Strict No-Build / No-Test Mode)

---

## 1. Technical Objective

Implement the complete local inspection and commit dispatch engine comprising:
1. **`gitmap pending-commits` (alias `pc`):**
   - Discovers all repositories across the active workspace and Split-DB registry.
   - Inspects both uncommitted working tree modifications (`git status --porcelain`) and unpushed local commits ahead of upstream (`git rev-list @{u}..HEAD`).
   - Renders a clean summary box first, with sorting by `priority` (default), `name`, or `count`.
   - Supports detailed inspection (`--detail all` or `--detail 1` / 1-by-1 interactive step-through).
   - Supports `-ssh` / `--ssh` flag to query and aggregate pending commits across all registered cluster SSH nodes.
   - Provides full JSON telemetry via `--json` (`-j`).
2. **`gitmap sends <verb> <target> "<message>"`:**
   - Dispatches conventional semantic commits (`cpf` -> `Feature: `, `cpb` -> `Bug: `, `cpr` -> `Release: `, `commit-fix` -> `Fix: `, `cp`/`cm` -> clean/no prefix).
   - Resolves target scope: specific repository name/slug or `all` dirty repositories.
   - Bypasses clean repositories automatically when targeting `all`.
   - Supports `--dry-run` (`-n`), `--no-push`, and `--json`.
3. **Comprehensive Unit & E2E Tests:**
   - Author thorough test coverage in `cli/cmd/pending_commits_test.go` and `cli/cmd/sends_test.go` using mock runners and table-driven scenarios.

---

## 2. Step-by-Step Implementation Plan for Worker 01

### Step 1: Command Constants & Type Definitions
1. In `cli/constants/constants_cli.go`:
   - Define command and alias constants:
     - `CmdPendingCommits = "pending-commits"`
     - `CmdPendingCommitsAlias = "pc"`
     - `CmdSends = "sends"`
2. In `cli/cmd/pending_commits_types.go`:
   - Declare data structures strictly enforcing positive boolean nomenclature:
     - `PendingCommitsPayload` (`TotalReposScanned`, `TotalDirtyRepos`, `TotalUncommittedFiles`, `TotalUnpushedCommits`, `SortMode`, `DetailMode`, `Repositories []RepoPendingCommitRecord`)
     - `RepoPendingCommitRecord` (`RepoName`, `RelativePath`, `CurrentBranch`, `IsDirty`, `IsClean`, `HasUncommitted`, `HasUnpushed`, `HasUpstream`, `UntrackedFilesCount`, `ModifiedFilesCount`, `StagedFilesCount`, `UnpushedCommitsCount`, `PendingFiles []string`, `UnpushedCommitSHAs []string`)
     - `PendingCommitsOptions` (`SortMode`, `DetailMode`, `IsSSH`, `IsJSON`, `IsDirtyOnly`, `IsAll`)
3. In `cli/cmd/sends_types.go`:
   - Declare data structures:
     - `SendsExecutionPayload` (`Verb`, `PrefixApplied`, `RawMessage`, `FinalCommitMessage`, `TargetScope`, `IsDryRun`, `IsPushed`, `TotalProcessed`, `TotalCommitted`, `TotalSkipped`, `Results []RepoSendResultRecord`)
     - `RepoSendResultRecord` (`RepoName`, `RelativePath`, `Status`, `Branch`, `HeadSHA`, `FilesStaged`, `IsSuccess`, `ErrorMessage`)
     - `SendsOptions` (`Verb`, `Target`, `RawMessage`, `IsDryRun`, `IsPushed`, `IsJSON`, `IsVerbose`)

### Step 2: Local Pending Commits Engine Implementation
1. In `cli/cmd/pending_commits_cmd.go`:
   - Implement `RunPendingCommits(args []string) error`:
     - Parse flags: `--sort`, `-s`, `--detail`, `-d`, `--ssh`, `--json`, `-j`, `--dirty-only`, `--all`.
     - Route `--help` / `-h` to help renderer.
     - If `--ssh` is present, invoke fleet aggregator `aggregateSSHPendingCommits(localRecords, opts)`.
   - Implement repository resolution:
     - Discover top-level git repos using `cmdpull.ResolvePullDirectoryTargets(cwd)`.
     - Cross-reference with `store.DB.ListRepos()` for enriched repository metadata.
   - Implement repository status inspection:
     - Execute `git status --porcelain` to capture untracked, modified, and staged file counts.
     - Execute `git rev-list --count @{u}..HEAD` to capture unpushed commit counts.
     - Populate `RepoPendingCommitRecord` ensuring `IsClean = !IsDirty && !HasUnpushed`.

### Step 3: Sorting & View Modes
1. Implement sorting algorithms in `sortPendingCommits(records []RepoPendingCommitRecord, sortMode string)`:
   - `priority`: `(IsDirty && HasUnpushed)` > `IsDirty` > `HasUnpushed` > `IsClean`.
   - `count`: Descending total changes `(untracked + modified + staged + unpushed)`.
   - `name`: Ascending alphabetical slug comparison.
2. Implement rendering:
   - If `opts.IsJSON`: Marshal `PendingCommitsPayload` to formatted JSON on stdout.
   - If terminal mode:
     - Summary First: Print summary box card with totals and tabular overview.
     - Detail mode (`--detail all`): Print file changes and unpushed commit SHAs.
     - Detail 1-by-1 (`--detail 1`): Interactive pagination stepping through dirty repositories.

### Step 4: Fleet SSH Aggregation (`-ssh` / `--ssh`)
1. Implement `aggregateSSHPendingCommits(localPayload PendingCommitsPayload, opts PendingCommitsOptions) error`:
   - Query `cmdssh.FetchAllSSHConnections()` for registered cluster nodes.
   - Execute parallel goroutines dialing each node via `crypto.ConnectWithFallback` / `crypto.RunCommand`.
   - Remote command string: `gitmap pending-commits --json`.
   - Parse remote stdout JSON into `PendingCommitsPayload`.
   - If node is unreachable or encounters timeout, record error gracefully without aborting other nodes.
   - Render multi-node aggregated summary or unified JSON.

### Step 5: `gitmap sends` Implementation
1. In `cli/cmd/sends_cmd.go`:
   - Implement `RunSends(args []string) error`:
     - Syntax: `gitmap sends <verb> <target> "<message>" [flags]`.
     - Validate verb (`cpf`, `cpb`, `cpr`, `commit-fix`, `cp`, `cm`).
     - Resolve semantic prefix:
       - `cpf`: `Feature: `
       - `cpb`: `Bug: `
       - `cpr`: `Release: `
       - `commit-fix`: `Fix: `
       - `cp` / `cm`: `""`
     - Parse execution flags: `--dry-run` (`-n`), `--no-push`, `--json` (`-j`), `--verbose` (`-v`).
2. Target resolution:
   - If `target == "all"`:
     - Inspect all workspace repositories.
     - Filter only repositories where `isDirty == true`.
     - Clean repositories marked as `status: "clean-skipped"`.
   - If `target == "$repoName"`:
     - Find matching repository by slug, name, or directory.
     - Return error if target repository does not exist.
3. Execution loop:
   - If `opts.IsDryRun`: Simulate stage, commit message, and push; record `status: "dry-run-simulated"`.
   - If real execution:
     - `git add -A`
     - `git commit -m "<finalCommitMessage>"`
     - If `opts.IsPushed`: `git push`
     - Record `status: "pushed"` or `status: "committed"`.
4. Telemetry & Output:
   - If `opts.IsJSON`: Print `SendsExecutionPayload` JSON.
   - If terminal mode: Render execution summary card.

### Step 6: CLI Routing & Root Registration
1. In `cli/cmd/rootcore.go`:
   - Register entries in `coreClusterEntries()` / `coreVisibilityActionEntries()` or appropriate dispatch group:
     - `{[]string{constants.CmdPendingCommits, constants.CmdPendingCommitsAlias}, func() error { return runPendingCommits(argsTail()) }}`
2. In `cli/cmd/roottooling.go`:
   - Register entries for `sends`:
     - `{[]string{constants.CmdSends}, func() error { return runSends(argsTail()) }}`

### Step 7: UI Help Menus & Dashboards
1. In `cli/cmd/pending_commits_help.go`:
   - Implement `PrintPendingCommitsHelp()` rendering ANSI-colored box help, supported flags, and usage examples.
2. In `cli/cmd/sends_help.go`:
   - Implement `PrintSendsHelp()` detailing semantic verbs (`cpf`, `cpb`, `cpr`, `commit-fix`), target resolution (`$repoName`, `all`), and dry-run safety.

### Step 8: Unit & E2E Testing
1. In `cli/cmd/pending_commits_test.go`:
   - Test flag parsing (`--sort`, `--detail`, `--ssh`, `--json`).
   - Test sorting logic (`priority`, `name`, `count`).
   - Test status classification (clean vs dirty vs unpushed).
   - Test JSON output serialization and schema validation.
2. In `cli/cmd/sends_test.go`:
   - Test semantic prefix mapping across all verbs (`cpf`, `cpb`, `cpr`, `commit-fix`, `cp`, `cm`).
   - Test target resolution (`all` vs single repo).
   - Test clean-repository skipping behavior.
   - Test dry-run simulation mode without modifying git state.
   - Test JSON serialization.

---

## 3. Worker 01 Acceptance Verification Matrix

| Gate ID | Requirement | Verification Step | Pass Criteria |
| :--- | :--- | :--- | :--- |
| **VG-01** | Command Routing | Test `gitmap pending-commits`, `gitmap pc`, `gitmap sends` | No panic, correct subcommand execution |
| **VG-02** | Status Classification | Inspect dirty working tree vs unpushed commits | Distinct counts for modified vs ahead |
| **VG-03** | Target Resolution | Test `sends cpf all "msg"` and `sends cpf gitmap "msg"` | `all` skips clean repos; single matches slug |
| **VG-04** | Semantic Prefix | Test verb prefix injection | Correct prefixes (`Feature: `, `Bug: `, etc.) |
| **VG-05** | Fleet SSH Delegation | Test `-ssh` flag handling | Remote JSON unmarshaled and aggregated |
| **VG-06** | UI Help Text | Test `--help` flag on both commands | Box help rendered with examples |
| **VG-07** | Coding Guidelines | Audit paths and boolean nomenclature | 100% relative paths, positive booleans |
