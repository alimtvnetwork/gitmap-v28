# Cluster Command Delegation — `servers-clients` / `clients` / `sc`

**Slug:** cluster-command-delegation  
**Steps:** 100  
**Status:** completed  
**Created:** 2026-08-17  
**Spec:** `.lovable/spec/commands/07-cluster-command-delegation.md`

---

## Context

After `gitmap serve` and `gitmap join` establish a cluster of machines (Phase 6), operators need a unified CLI surface to broadcast shell commands, install packages, run Git operations, manage project automation, and control machine lifecycle—all from a single workstation. Every execution is persisted to SQLite for a full audit trail.

**Captured inputs:**
- Spec: `.lovable/spec/commands/07-cluster-command-delegation.md`
- Prerequisite: Phase 6 from `05-gitmap-improvements.md` (cluster join/serve) must be complete before steps 50+.

---

## Steps

### PART A — Foundation & Database Schema (Steps 1–15)

1. Confirm Phase 6 (`gitmap serve` / `gitmap join`) is complete and the `ClusterNode` table exists with fields: `NodeId`, `Alias`, `DisplayId`, `IPAddress`, `NodeRole`, `OS`, `Status`, `PasswordHash`, `PackageManager`.
2. Read all coding guidelines: `spec/05-coding-guidelines/**/*.md` and `spec/12-consolidated-guidelines/`. Apply Boolean (`is`/`has` prefix), int-backed enums, and no-swallow-error rules to all code steps below.
3. Write root cause entry in `.lovable/memory/last-failure.md` for this feature: "Cluster has no command delegation surface; operators must SSH individually into each node."
4. [x] Extend `gitmap/db/enums.go` with `CommandKindEnum` (PsCommand=1, CmdCommand=2, Install=3, GitPull=4, GitPush=5, GitCommit=6, GitStatus=7, ProjRun=8, ProjCreateCICD=9, Restart=10, Shutdown=11, Logoff=12) and extend `ResultStatusEnum` with `Deferred=5`, `RequiresAuth=6`. Add `String()` and `ParseCommandKind` / `ParseResultStatus` helpers.
5. [x] Create SQL migration `gitmap/db/migrations/0NN_cluster_run_exec_result.sql` defining `ClusterRun` (RunRef, CommandKind, RawCommand, TargetSelector, ExceptClause, StartedAt, FinishedAt, TotalNodes, SucceededNodes, FailedNodes, SkippedNodes) and `ClusterExecResult` (ClusterRunId FK, NodeId FK, SubCommand, CommandText, ResultStatus, ExitCode, Stdout truncated 64KB, Stderr truncated 16KB, StartedAt, FinishedAt, DurationMs, ErrorMessage).
6. [x] Register the migration in `gitmap/db/migrations.go` (or the existing migration registry mechanism). Run `go test ./gitmap/db/...` against `:memory:` to confirm schema applies cleanly.
7. [x] Implement `gitmap/db/clusterrun.go`: `InsertClusterRun(db, ClusterRun) (id int64, err error)`, `UpdateClusterRun(db, id, FinishedAt, counts)`, `SelectClusterRun(db, RunRef) (ClusterRun, error)`, `ListClusterRuns(db, limit int) ([]ClusterRun, error)`.
8. [x] Implement `gitmap/db/clusterexecresult.go`: `InsertClusterExecResult(db, ClusterExecResult) (id int64, err error)`, `UpdateClusterExecResult(db, id, ResultStatus, ExitCode, Stdout, Stderr, FinishedAt, DurationMs, ErrorMessage)`, `SelectClusterExecResultsByRunId(db, runId) ([]ClusterExecResult, error)`.
9. [x] Add `gitmap/db/clusterrun_test.go`: in-memory SQLite tests for insert, update, and FK cascade on run delete. All green.
10. [x] Centralize all log and error message strings for this feature in `gitmap/constants/constants_cluster.go`: `MsgClusterPreflight`, `MsgClusterCountdown`, `ErrClusterNoNodes`, `ErrClusterMutuallyExclusiveFlags`, `ErrClusterLifecycleRequiresForce`, `ErrClusterPasswordRequired`, `MsgClusterAuditFooter`, etc. No magic strings anywhere.
11. [x] Create `gitmap/cluster/node_resolver.go`: `ResolveTargetNodes(selector TargetSelector, filter NodeFilter, allNodes []ClusterNode) ([]ClusterNode, error)`. `TargetSelector` is an int-backed enum (ServersClients=1, ClientsOnly=2). `NodeFilter` holds `ExceptIDs`, `ExceptIPs`, `OnlyIDs`, `OnlyIPs` slices. Validate mutual exclusivity of `--except` vs `--ip`/`--id` here.
12. [x] Implement IP/ID exclusion matching in `node_resolver.go`: support integer `DisplayId`, full IP match, and trailing-octet partial match (e.g. `24` matches `192.168.1.24`). Range support for IDs (`2-5` → IDs 2,3,4,5).
13. [x] Add `gitmap/cluster/node_resolver_test.go`: table tests for all exclusion types (integer, IP, partial octet, range, mixed) and mutual-exclusivity validation. All green.
14. [x] Implement `RunRefGenerator`: generates `RUN-<YYYYMMDD>-<NNN>` identifiers using the current date and a daily counter read from the `ClusterRun` table. Zero-pad the sequence counter to 3 digits.
15. [x] Write `gitmap/cluster/preflight.go`: `PrintPreflight(selector, effective []ClusterNode, command string, runRef string, autoConfirm bool)` renders the preflight summary box (as per spec §3.3) using existing terminal formatting helpers. Return `(confirmed bool, err error)`.

---

### PART B — Argument Parsing & CLI Wiring (Steps 16–30)

16. [x] Add CLI ID constants to `gitmap/constants/constants_cli.go`: `CmdServersClients = "servers-clients"`, `CmdSC = "sc"`, `CmdClients = "clients"`, `CmdClusterHistory = "cluster history"`, `CmdClusterExport = "cluster export"`, `CmdClusterImport = "cluster import"`, `CmdClusterSetPassword = "cluster set-password"`. Add `// gitmap:cmd top-level` markers for all primary commands.
17. [x] Create `gitmap/cmd/clusterflags.go`: define `ClusterFlags` struct with `Selector TargetSelector`, `ExceptClause string`, `OnlyIPs []string`, `OnlyIDs []int`, `AutoConfirm bool`, `ForceLifecycle bool`, `NoPreflight bool`. Implement `ParseClusterFlags(args []string) (ClusterFlags, []string, error)` — returns flags, remaining positional args, and error.
18. [x] Create `gitmap/cmd/clusterflags_test.go`: table tests covering `--except 2,4`, `--except 192.168.1.24,151`, `--except 2-5`, `--ip ...`, `--id ...`, mutual exclusivity error, and `-Y` auto-confirm. All green.
19. [x] Create `gitmap/cmd/clustersubcmd.go`: define `ClusterSubCommand` struct with `Kind CommandKindEnum`, `RawArg string` (the quoted command or package list). Implement `ParseSubCommands(tokens []string) ([]ClusterSubCommand, error)` — parses comma-separated chained sub-commands from remaining args after flag parsing.
20. [x] Register `servers-clients` and `sc` in the dispatcher: both route to `runClusterCommand(selector=ServersClients, args)`.
21. [x] Register `clients` in the dispatcher: routes to `runClusterCommand(selector=ClientsOnly, args)`.
22. [x] Register `cluster history`, `cluster export`, `cluster import`, `cluster set-password` in the dispatcher.
23. [x] Implement `runClusterCommand(selector, args)` orchestrator: parse flags → resolve nodes → print preflight → confirm → insert `ClusterRun` → dispatch sub-commands → update `ClusterRun` final counts → print audit footer.
24. [x] Implement `gitmap/cluster/dispatcher.go`: `Dispatch(ctx, node ClusterNode, subCmd ClusterSubCommand) ClusterExecResult` — routes to the correct executor (ps, cmd, install, git-op, proj, lifecycle) based on `subCmd.Kind`.
25. [x] Implement the bounded worker pool in `gitmap/cluster/pool.go` (reuse/adapt from Phase 4 clone pool): `RunPool(ctx, nodes, subCmds []ClusterSubCommand, db, runId, maxWorkers int, resultCh chan<- ClusterExecResult)`. Each worker processes one node's sub-command chain sequentially. Pool size default: 10.
26. [x] Wire the parallel-spinner UI (Phase 4 pattern): each node gets a UI row. Row format: `[<DisplayId>/<Alias>] ⠋ Running ps…`. On finish: `[1/dev-01] ✓ ps (0ms, exit 0)`. On fail: `[2/dev-02] ✗ cmd (exit 1)`.
27. [x] Implement streaming output: node stdout lines prefixed with `[<DisplayId>/<Alias>] ` forwarded to the operator's terminal in real-time when `--verbose` is passed. Without `--verbose`, stdout captured and stored in `ClusterExecResult.Stdout`.
28. [x] Implement `Ctrl+C` handler: cancel the context, mark all in-flight nodes as `ResultStatus = Skipped`, call `UpdateClusterRun` with partial counts, print audit footer.
29. [x] Implement final summary box after pool completes: `┌ Cluster Run RUN-YYYYMMDD-NNN ─────────────────┐ │ Nodes: 6  OK: 5  Failed: 1  Skipped: 0       │ └────────────────────────────────────────────────┘`. List error details below.
30. [x] Add `go test ./gitmap/cmd/... -run TestClusterCommand` with stubbed node executors to validate end-to-end orchestration without network calls.

---

### PART C — Sub-Command Executors (Steps 31–55)

31. [x] Implement `gitmap/cluster/exec_ps.go`: `ExecPS(ctx, node ClusterNode, command string) (stdout, stderr string, exitCode int, err error)`. Windows: `pwsh -NonInteractive -Command <cmd>`, fallback `powershell`. Non-Windows: `pwsh -Command <cmd>`, skip with warning if not found.
32. [x] Implement `gitmap/cluster/exec_cmd.go`: `ExecCmd(ctx, node ClusterNode, command string) (stdout, stderr string, exitCode int, err error)`. Windows: `cmd.exe /C <cmd>`. Non-Windows: `/bin/sh -c <cmd>`.
33. [x] Implement `gitmap/cluster/exec_install.go`: `ExecInstall(ctx, node ClusterNode, packages []string) (results []PackageResult, err error)`. Detect package manager from `node.PackageManager` or auto-detect via node OS. Install each package individually. `PackageResult` holds `PackageName`, `Succeeded bool`, `Stderr string`.
34. [x] Implement package manager detection: probe in priority order per OS. Cache detected manager back to DB via `UpdateClusterNodePackageManager`.
35. [x] Implement `gitmap/cluster/exec_git.go`: `ExecGitPull(ctx, node ClusterNode) (stdout, stderr string, exitCode int, err error)` — runs `gitmap pull --all` on the remote node (dispatched as a `ps`/`sh` command calling the remote gitmap binary). Same pattern for `ExecGitPush`, `ExecGitCommit`, `ExecGitStatus`.
36. [x] Implement `gitmap/cluster/exec_proj.go`: `ExecProjRun(ctx, node ClusterNode, projectNames []string) (results []ProjRunResult, err error)`. Scans registered repo paths on the node, finds `run.ps1` or `run.sh`, executes, captures last 20 lines on failure.
37. [x] Implement `ExecProjCreateCICD` as a stub returning `ResultStatus = Deferred` with message `"create-cicd: reserved for future spec"`.
38. Implement `gitmap/cluster/exec_lifecycle.go`: `ExecRestart`, `ExecShutdown`, `ExecLogoff`. Each checks `node.NodeRole != 'server'` (safety guard — server node is never a target of lifecycle commands). Uses stored password hash to authenticate if required.
39. Add `--force-lifecycle` guard: lifecycle commands fail with a clear error if `--force-lifecycle` is not also present, preventing accidental shutdown of the cluster.
40. Implement 5-second countdown with abort: `printCountdown(nodes []string, action string, seconds int)` — prints `⚠ <action> <N> nodes in <S>s… Press Ctrl+C to abort` for each second. Uses `time.Tick` and `select` on context.
41. Add unit tests for `exec_ps.go` with injectable `execRunner` seam (no real subprocess). Validate PowerShell vs pwsh fallback logic. All green.
42. Add unit tests for `exec_install.go`: validate package manager detection and per-package result capture. All green.
43. Add unit tests for `exec_lifecycle.go`: validate server-node guard and `--force-lifecycle` requirement. All green.
44. Add unit tests for `exec_git.go`: validate the `gitmap pull --all` invocation string on Windows vs Unix. All green.
45. Wire `ExecPS` and `ExecCmd` into the remote execution transport (Phase 6 gRPC/HTTP channel) so commands are actually dispatched over the network to the client agent.
46. Implement the client-side agent handler: receives a `ClusterSubCommand` proto/JSON payload, executes locally, streams stdout/stderr back, and returns exit code.
47. Add TLS to the cluster transport channel (reuse the self-signed cert from Phase 6 `gitmap serve`). Verify client fingerprint matches on join.
48. Implement `--verbose` flag: when set, forward all node stdout to the operator terminal in real-time rather than buffering. Prefix every line.
49. Implement `--dry-run` flag: print the preflight and sub-command list but do NOT execute. Record `ResultStatus = Skipped` for all nodes. Useful for validation.
50. Add integration test: spin up two in-process fake cluster nodes, dispatch a `ps "echo hello"` command, assert both `ClusterExecResult` rows are `Succeeded`.

---

### PART D — Audit, History & Export (Steps 51–65)

51. Implement `gitmap cluster history` command: queries `ClusterRun` table ordered by `StartedAt` DESC, prints a table with columns `RunRef | CommandKind | TargetSelector | Nodes | OK | FAIL | StartedAt`.
52. Implement `gitmap cluster history <RunRef>`: expands a single run to show per-node `ClusterExecResult` rows with `[DisplayId/Alias] SubCommand ResultStatus ExitCode DurationMs`.
53. Implement `gitmap cluster export [--format json|csv] [--output <file>]`: dumps `ClusterNode` table. Omit `PasswordHash`. Default format: JSON.
54. Implement `gitmap cluster import <file>`: merges nodes by `NodeId` (upsert). Reports `Inserted: N, Updated: M, Skipped: K`.
55. Implement `gitmap cluster set-password --id <id>`: prompts for password twice (hidden echo), bcrypt-hashes (cost=12), updates `ClusterNode.PasswordHash`. Prints `✓ Password updated for node <Alias>`.
56. Add `gitmap cluster reset-password --id <id>`: clears the stored hash. Requires `--confirm` flag.
57. Implement `gitmap cluster nodes` (alias: `gitmap cluster ls`): lists all registered nodes with `DisplayId | Alias | IP | OS | Role | Status | LastHeartbeat`.
58. Add `--json` flag to `gitmap cluster nodes` for machine-readable output.
59. Implement `gitmap cluster remove --id <id>`: removes a node from the registry. Requires `--confirm`.
60. Add `gitmap cluster audit-clean --before <date>`: purges `ClusterRun` and `ClusterExecResult` rows older than the given date. Shows row counts before purging. Requires `--confirm`.
61. Add indexes: `ClusterRun(StartedAt)`, `ClusterExecResult(ClusterRunId)`, `ClusterExecResult(NodeId)`. Add to migration SQL.
62. Ensure Stdout/Stderr storage is capped at 64KB / 16KB in the insert helpers. Truncate with a `[truncated…]` suffix.
63. Add `ClusterRun.RawCommand` storage to capture the exact CLI string entered by the operator.
64. Write `gitmap/db/clusterrun_test.go` full suite: insert run, insert 5 exec results, query by RunRef, verify counts, verify FK cascade. All green.
65. Add `gitmap cluster stats`: prints aggregate statistics — total runs, total commands dispatched, success rate, most-targeted node, most-used sub-command.

---

### PART E — Help Text & Documentation (Steps 66–80)

66. Create `gitmap/cmd/help/servers-clients.md`: full help with all sub-commands, flag reference table, and 8–12 realistic simulation examples matching spec §3.2.
67. Create `gitmap/cmd/help/sc.md`: 5-line stub with description and `See: gitmap servers-clients --help`.
68. Create `gitmap/cmd/help/clients.md`: clients-only variant help with examples.
69. Create `gitmap/cmd/help/cluster-history.md`: explains RunRef format, table columns, per-node expansion.
70. Create `gitmap/cmd/help/cluster-export.md`: explains JSON/CSV format, omitted password field.
71. Create `gitmap/cmd/help/cluster-import.md`: explains merge semantics.
72. Create `gitmap/cmd/help/cluster-set-password.md`: security notes, bcrypt details.
73. Create `gitmap/cmd/help/cluster-nodes.md`: column definitions, `--json` flag.
74. Regenerate LLM.md golden fixture via `gitmap regoldens` (or the project's golden fixture tool). Confirm no diff regressions.
75. Add `servers-clients`, `sc`, `clients`, `cluster history`, `cluster nodes`, `cluster export`, `cluster import`, `cluster set-password` to `src/data/commands.ts` with example invocations matching the help MDs.
76. Add a **Cluster Command Delegation** section to the docs-site `commands.md` or equivalent page. Describe target selectors, chaining, `--except` / `--ip` / `--id` semantics, and the audit trail.
77. Add a **Security** note to the docs: credentials stored as bcrypt, passwords never exported, TLS required.
78. Update the main `README.md` with a `## Cluster` section showing the 3 most common command patterns as a quick-start.
79. Update `changelog.md` with a new `## Added` block under the next version describing all cluster delegation commands.
80. Bump version constant across unified sites (Go constants, `version.json`, README pins, `src/constants/index.ts`).

---

### PART F — Testing, Linting & Verification (Steps 81–95)

81. Run `go vet ./...` — zero findings required.
82. Run `golangci-lint run ./...` (v1.64.8 per Core memory) — zero findings required.
83. Run `go test ./gitmap/cluster/... -count=1 -race` — all green.
84. Run `go test ./gitmap/db/... -count=1 -race` — all green (including new migration).
85. Run `go test ./gitmap/cmd/... -count=1 -race` — all green (including new flag parser tests).
86. Run `go test ./... -count=1` full suite — note and document any pre-existing failures (Python env) vs. regressions from this change.
87. Run `gofmt -l .` — no files with formatting differences.
88. Add `TestClusterSubCommandParser` covering all chaining forms: `cmd "...", ps "..."`, `ps "...", install "pkg1,pkg2"`, error on unknown sub-command token.
89. Add `TestNodeResolverExclusion` covering: all nodes included (no filter), exclude by ID, exclude by IP, exclude by partial octet, exclude by range, mutual exclusivity error.
90. Add `TestPreflightBoxRender`: assert the preflight box string contains `Target`, `Excluded`, `Effective`, `Command`, `Audit ID` for a known input. Snapshot test.
91. Add `TestRunRefGenerator`: assert format `RUN-YYYYMMDD-NNN`, assert daily sequence increments correctly in a `:memory:` DB.
92. Add `TestClusterPoolCancellation`: start a pool of 5 nodes, cancel context after 2 complete, assert remaining are `Skipped`.
93. Add `TestLifecycleServerGuard`: verify that when `NodeRole = 'server'` is in the effective targets, lifecycle commands return `ErrClusterServerProtected`.
94. Add `TestForceLifecycleGuard`: verify that without `--force-lifecycle`, restart/shutdown/logoff return `ErrClusterLifecycleRequiresForce`.
95. Manual smoke test (document in verification notes): start `gitmap serve`, join 2 client nodes (can be localhost on different ports), run `gitmap sc ps "hostname"`, confirm 3 rows in `ClusterExecResult` (server + 2 clients), run `gitmap cluster history` to see the run.

---

### PART G — Release & Plan Closeout (Steps 96–100)

96. Ensure `go build ./...` succeeds with no warnings across all packages.
97. Run `gitmap fix-repo --strict` over touched packages to auto-rewrite `{base}-vN` references if applicable (per project Fix-Repo Strict Mode memory).
98. Commit all changes: `feat: cluster command delegation (servers-clients, clients, sc, audit trail)`.
99. Push to `origin main`. Verify CI passes (excluding known pre-existing Python env failure).
100. Move this file to `.lovable/plans/completed/06-cluster-command-delegation.md`, flip `Status: completed`, and update `.lovable/plans/index.md`.

---

## Verification

- `go test ./gitmap/cluster/... ./gitmap/db/... ./gitmap/cmd/... -count=1 -race` — all green.
- `gitmap --version` prints the new bumped version.
- Preflight box renders correctly in terminal.
- Manual smoke (step 95) shows 3 `ClusterExecResult` rows.
- `gitmap cluster history` lists the run.
- `gitmap cluster export` produces valid JSON without passwords.
- Help golden fixtures match (no diff after `gitmap regoldens`).
- `gitmap sc ls` lists all nodes in aligned table format.
- `gitmap sc cat p1/config.yaml` prints file content prefixed per node.
- `gitmap sc write p1/note.txt "hello"` writes atomically and reports bytes written per node.
- `gitmap sc set-path-alias "C:\Projects as p1"` prints the tips block.
- `gitmap clients clone url1,url2 --workdir p1` clones into the registered alias path and auto-runs status.
- `gitmap sc update` updates gitmap.
- `gitmap sc update antigravity vscode` updates gitmap and specified packages.

---

### PART H — Node Listing, Remote Clone, Path Aliases, Cat/Write (Steps 101–120)

101. Add CLI constants to `constants_cli.go`: `CmdServersLS = "servers ls"`, `CmdClientsLS = "clients ls"`, `CmdSCLS = "servers-clients ls"`, `CmdSCCat = "servers-clients cat"`, `CmdSCWrite = "servers-clients write"`, `CmdSCSetDefaultPath = "servers-clients set-default-path"`, `CmdSCSetPathAlias = "servers-clients set-path-alias"`. Add `// gitmap:cmd top-level` markers.
102. Register `servers ls`, `clients ls`, `servers-clients ls` (and `sc ls`) in the dispatcher: route to `runClusterLS(selector, args)`. This command reads the local `ClusterNode` SQLite table — no Phase 6 network call needed.
103. Implement `runClusterLS(selector, args)`: resolve nodes by selector + filters → render aligned table: `ID | IP Address | Machine Name | Role | OS | Status | Last Heartbeat`. `servers ls` hides the Role column. Add `--json` flag. No preflight (read-only).
104. Implement `gitmap/cluster/lsrender.go`: `RenderNodeTable(nodes []ClusterNode, showRole bool) string` — dynamic column widths, minimum 2 spaces padding between columns, reusing the `statusTableContext` pattern from Phase 2.
105. Add `gitmap/cluster/lsrender_test.go`: snapshot tests for 0 nodes, 1 node, 3 nodes with mixed roles and OSes.
106. Create `gitmap/db/migrations/0NN_node_path_alias.sql`: adds `DefaultPath TEXT NOT NULL DEFAULT ''` to `ClusterNode`; creates `NodePathAlias(NodePathAliasId INTEGER PK, NodeId TEXT NOT NULL FK, Alias TEXT NOT NULL, AbsolutePath TEXT NOT NULL, CreatedAt DATETIME)` with unique index on `(NodeId, Alias)`.
107. Implement `gitmap/db/nodepath.go`: `UpsertPathAlias(db, nodeId, alias, absPath string) error`, `GetPathAlias(db, nodeId, alias string) (string, error)`, `SetDefaultPath(db, nodeId, path string) error`, `GetDefaultPath(db, nodeId string) (string, error)`, `ListPathAliases(db, nodeId string) ([]NodePathAlias, error)`.
108. Implement `gitmap/cluster/pathalias.go`: `ParseSetPathAliasArg(raw string) ([]AliasEntry, error)` — parses `"C:\Projects as p1, D:\Backup as p2"` into `[]AliasEntry`. Validate: alias must be alphanumeric + hyphens only; path must be non-empty.
109. Implement `runClusterSetDefaultPath(selector, args)`: parse optional `as <alias>` from arg, resolve nodes, call `SetDefaultPath` (and optionally `UpsertPathAlias`) per node in DB, print the tips block (per spec §12.1). Commit change locally — no network dispatch needed at this stage.
110. Implement `runClusterSetPathAlias(selector, args)`: parse multi-alias string via `ParseSetPathAliasArg`, resolve nodes, call `UpsertPathAlias` for each alias × each node, print the tips block (per spec §12.2).
111. Implement `gitmap/cluster/pathresolver.go`: `ResolvePath(db, nodeId, pathOrAlias string) (string, error)` — try `GetPathAlias` first; if alias not found and the string contains no path separator, return `ErrAliasNotFound` with a clear message; otherwise treat as literal path.
112. Register `clients clone`, `servers-clients clone`, `clients cfrp`, `servers-clients cfrp`, `clients cfr`, `servers-clients cfr` in the dispatcher → `runClusterClone(selector, subCmd, args)`.
113. Implement `runClusterClone(selector, subCmd, args)`: parse comma-separated URLs, resolve `--workdir` or `--path` aliases per node via `ResolvePath`, preflight, insert `ClusterRun`, dispatch each URL to each node using the Phase 4 bounded worker pool + spinner UI. After all clones complete, trigger `gitmap status` on the resolved workdir on each node and stream results prefixed with `[<ID>/<Alias>]`.
114. Implement `gitmap/cluster/exec_clone.go`: `ExecCloneURL(ctx, node ClusterNode, url, workdir, subCmd string) (stdout, stderr string, exitCode int, err error)` — builds the correct `gitmap clone`, `gitmap cfrp`, or `gitmap cfr` command string and dispatches it over the cluster transport (SSH/gRPC).
115. Handle `--path "p1, p2"` flag in `runClusterClone`: split by comma, resolve each alias per node via `ResolvePath`, map Nth URL → Nth path; extra URLs reuse the last resolved path.
116. Register `servers-clients cat` (and `sc cat`, `clients cat`) in the dispatcher → `runClusterCat(selector, args)`.
117. Implement `runClusterCat(selector, args)`: resolve path/alias per node via `ResolvePath`; dispatch `Get-Content` (Windows) or `cat` (Unix) over cluster transport; prefix output `[<ID>/<Alias>] ── <resolved-path> ──`; truncate at 64 KB with `[truncated…]`; reject binary (non-UTF-8 bytes); store result in `ClusterExecResult` (SubCommand = `cat`).
118. Register `servers-clients write` (and `sc write`, `clients write`) in the dispatcher → `runClusterWrite(selector, args)`.
119. Implement `runClusterWrite(selector, args)`: validate content length ≤ 1 MB (reject before any network call); require `--force-write` if targeting `--all` or more than 5 effective nodes; resolve path/alias; dispatch atomic write (write to temp file then rename) over cluster transport; report `[<ID>/<Alias>] ✓ wrote <N> bytes to <resolved-path>` per node; store in `ClusterExecResult` (SubCommand = `write`).
120. Add help MD files: `servers-ls.md`, `clients-ls.md`, `servers-clients-ls.md`, `cluster-set-default-path.md`, `cluster-set-path-alias.md`, `cluster-cat.md`, `cluster-write.md`. Run `gitmap regoldens`. Update `src/data/commands.ts`. Commit: `feat: cluster ls, remote clone/cfr, path alias system, cat/write (steps 101-120)`. Push to `origin main`. Move plan to `completed/`.

---

### PART I — Remote Update Commands (Steps 121–125)

121. Add CLI constants to `constants_cli.go`: `CmdSCUpdate = "servers-clients update"`, `CmdSCUpdateAll = "servers-clients update-all"`. Register `servers update`, `clients update`, `servers-clients update` (and `sc update`) → `runClusterUpdate(selector, false, args)`. Register `update-all` variants → `runClusterUpdate(selector, true, args)`.
122. Implement `runClusterUpdate(selector, isAll, args)`: if `isAll`, show preflight warning requiring 'y' to proceed. Insert `ClusterRun` and use bounded worker pool spinner UI to dispatch. 
123. Implement `gitmap/cluster/exec_update.go`: `ExecUpdate(ctx, node, isAll, packages...)`. If `isAll`, build OS-specific command (`winget upgrade --all` or `choco upgrade all -y` for Windows; `apt-get update && apt-get upgrade -y` or `brew upgrade` for Unix).
124. Handle specific package updates in `ExecUpdate`: always prepend `gitmap` to the package list if not present. Build string: `winget upgrade <packages>` (Windows) or `apt-get upgrade <packages>` (Unix). Dispatch over transport and persist to `ClusterExecResult`.
125. Add help MD files: `cluster-update.md`, `cluster-update-all.md`. Run `gitmap regoldens`. Update `src/data/commands.ts`. Commit: `feat: cluster remote update commands (steps 121-125)`. Push to `origin main`.

---

## Appended from prior pending tasks

- `01-bulk-visibility-mapub-mapri.md` — still pending; unblocked by this plan.
- `05-gitmap-improvements.md` (Phases 5 & 6) — Phase 6 (`serve`/`join`) is a prerequisite for Part G, Part H (network commands), and Part I. Steps 1-80, 101-111 can be implemented prior to Phase 6.
