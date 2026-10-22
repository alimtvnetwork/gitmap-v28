# Research Report 02: 200 Comprehensive Architectural & Functional Improvements for GitMap & 200 Concrete Implementation Steps

**Parent Task:** 287-golang-method-return-and-gitmap-improvement-instructions  
**Subagent Role:** Research 02 (`TypeName: "research"`, Read-Only)  
**Target Recipient:** 144baa42-d371-4ace-b1cf-84a3fb65d85d (Parent Orchestrator)

---

## Executive Summary & Core Subsystem Survey

A comprehensive survey of GitMap's codebase was conducted across its core subsystems:
- **Subsystem A (Discovery, Scanner & Cloner):** `cli/scanpipe/`, `cli/cloner/`, `cli/clonenow/`
- **Subsystem B (Multi-Node Fleet & Cluster):** `cli/cmdssh/`, `cli/cmdnodes/`, `cli/cluster/`
- **Subsystem C (Split-DB Architecture & Engine):** `cli/store/`, `cli/cmddb/`, `cli/dbengine/`
- **Subsystem D (Pipeline Telemetry & Diagnostics):** `cli/cmdpe/`, `cli/cmdpipeline/`, `cli/cmdpullerror/`
- **Subsystem E (Credential Management & Auth):** `cli/cmdfixauth/`, `cli/gitutil/`
- **Subsystem F (AI Agent Engine & Heatmap):** `cli/cmdagent/`, `cli/heatmap/`
- **Subsystem G (Macro Automation & CLI Styling):** `cli/macro/`, `cli/cmdmacro/`, `cli/termout/`
- **Subsystem H (Security, Offloading & ADLP):** `cli/secrets/`, `cli/power/`

Below is the complete catalog of **EXACTLY 200 Architectural & Functional Improvements** organized into 10 thematic categories of 20 improvements each, followed by **EXACTLY 200 Concrete Implementation Steps** defining the exact Go structs, methods, algorithms, and database changes required.

---

# Part 1: Catalog of 200 Architectural & Functional Improvements

### Category 1: High-Speed Repository Discovery & Filesystem Walker (`cli/scanpipe/`)
1. **Lock-Free Atomic Directory Dispatch Queue:** Replace mutex-guarded queue in `scanner.go` with lock-free atomic pointer ring buffer for `dirJob`.
2. **Adaptive Worker Pool Sizing:** Dynamically scale `scanWorkersMax` based on disk I/O latency (NVMe SSD vs HDD vs network share) rather than fixed 16.
3. **Windows Fast I/O Traversal:** Add native Windows `FindFirstFileExW` / `FindNextFileW` with `FindExInfoBasic` to bypass Go runtime `os.ReadDir` allocation overhead.
4. **Selective Symlink Cycle Detection:** Add lightweight inode/device ID tracking to safely follow directory symlinks without infinite recursion.
5. **Worktree Sniffer Zero-Alloc Parser:** Optimize `gitFileSniffBytes` parser in `scanner.go` using `bytes.HasPrefix` and zero-copy string slice for `gitdir:`.
6. **Configurable Ignore Rules (.gitmapignore):** Add repository-root and home-directory `.gitmapignore` file support complementing hardcoded `excludedDirMap`.
7. **Incremental Filesystem Watcher Cache:** Implement persistent mtime/inode cache in `repodb` to skip directory trees whose mtime has not changed since last scan.
8. **Submodule Recursion Budgeting:** Allow controlled recursion into nested submodules up to configurable depth via `--submodule-depth`.
9. **Memory-Mapped Binary Probe in Walker:** Replace 8KB read in `walker.go` (`probeMaxBytes`) with memory-mapped byte buffer for ultra-fast binary detection.
10. **Parallel Sibling Directory Enqueueing:** Enqueue all child directories simultaneously in `sibling_parallel.go` to saturate multi-core traversal queues.
11. **Content-Hash Deduplication:** Compute xxHash64 of discovered files during walking to instantly detect identical cloned repos across disparate folders.
12. **Worktree Detection for Absorbed Submodules:** Detect submodule worktrees pointing to `.git/modules/<name>` and map them back to root project models.
13. **Relative Path Normalization Engine:** Centralize path normalization across OS boundaries using `filepath.ToSlash` with zero-allocation byte buffers.
14. **Depth-Boundary Audit Eventing:** Emit structured telemetry events when a repository is discovered at `DefaultMaxDepth` boundary.
15. **Non-Blocking Live Progress Channel:** Decouple `ScanProgress` callbacks using a ring-buffer channel to prevent slow UI renderers from throttling the crawler.
16. **File Size Threshold Filter:** Expose configurable file size cap via CLI flag `--max-file-size` (defaulting to 300KB) in `walker.go`.
17. **Concurrent Write-Time Batch Queries:** Optimize `scanWriteTimes` in `walker.go` by chunking SQLite SELECT statements into batches of 500 paths.
18. **Ephemeral Scan Results Storage:** Buffer scan results in an in-memory SQLite table before bulk flushing to persistent storage.
19. **Fast Pre-Filter by Directory Prefix:** Add prefix-trie index to `excludedDirMap` to instantly skip common build dirs (`.build`, `target`, `node_modules`).
20. **Graceful Context Cancellation & Drain:** Implement bounded cancellation drain timeout (500ms) in `pool.go` so in-flight workers release open file handles cleanly.

### Category 2: Parallel Batch Cloner, Round-Trip Restoration & Clonenow Engine (`cli/cloner/`, `cli/clonenow/`)
21. **Dynamic Concurrency Throttling on Git Cloner:** Automatically throttle `ConcurrentRunParams.Workers` when remote host rate-limiting (HTTP 429) is detected.
22. **Clonenow Resume State Journal:** Record per-row clone status in a SQLite journal during `gitmap clone-now` to enable `--resume` after network drop.
23. **Shallow-to-Full Clone Progression:** Implement `--depth 1` initial clone with background asynchronous unshallow fetch (`git fetch --unshallow`).
24. **Git Protocol v2 Enforcement:** Inject `GIT_PROTOCOL=version=2` environment variable across all spawned `git clone` and `git fetch` processes.
25. **Automated LFS Smudge Isolation:** Delay `git lfs pull` to post-clone stage with automatic retry and timeout handling in `lfs_retry.go`.
26. **Multi-Transport Coercion Pipeline:** Support dynamic regex-based transport rewriting in `CoerceURL` (e.g. converting SSH to authenticated HTTPS or vice versa).
27. **Fast-Forward vs Rebase Policy Engine:** Implement automated branch divergence resolution (`safe_pull_conflicts.go`) with configurable merge policies.
28. **Clonenow CSV/JSON Schema Validator:** Validate input files with schema checksums and strict headers before dispatching batch clones in `parse_schema.go`.
29. **Zero-Copy ScanRecord Serialization:** Stream scan records directly into clone worker channels without intermediate slice allocations.
30. **Persistent Clone Failure Quarantine:** Move failed clones to a quarantine table in `gitmap.db` with detailed stderr triage for one-click retry.
31. **Smart Bare-Repo Mirroring:** Add `--mirror` and `--bare` support in `cloner.go` for local backup and multi-node fleet caching.
32. **Clone-Now Pre-Flight Disk Space Checker:** Inspect target disk volume space before starting multi-gigabyte batch clones to prevent out-of-disk failures.
33. **Authenticated Git Credential Helper Injection:** Inject ephemeral `git -c credential.helper=...` arguments into cloner subcommands to prevent interactive prompts.
34. **Streamlined Batch Progress Multiplexer:** Group terminal progress outputs into multi-line ANSI widgets showing real-time transfer speeds and bytes transferred.
35. **Automated Worktree Recreation in Clonenow:** Parse worktree linkages from scan exports and automatically execute `git worktree add` during restoration.
36. **Sparse-Checkout Configuration Generator:** Allow batch cloning with `--sparse-paths` using Git 2.25+ sparse-checkout cone mode.
37. **Clone Cache Invalidation on Commit Drift:** Invalidate `CloneCache` entries when remote branch SHA advances ahead of recorded cache SHA.
38. **Parallel Submodule Initialization:** Trigger `git submodule update --init --recursive --jobs=N` concurrently with parent repo cloning.
39. **Network Bandwidth Rate Limiter:** Implement token bucket bandwidth limiter for `git clone` subprocesses via `git config http.lowSpeedLimit`.
40. **Clean Working Tree Safety Guard:** Check for uncommitted changes in existing target folders before executing `OnExists == "update"` in `clonenow`.

### Category 3: Multi-Node SSH Fleet Orchestration & Host Management (`cli/cmdssh/`)
41. **SSH Connection Pooling & Multiplexing:** Implement OpenSSH `ControlMaster` / `ControlPath` socket reuse in `cmdssh` to reduce connection latency.
42. **Resilient Fallback Authentication Chain:** Sequence Ed25519 -> RSA -> Password -> Pageant/SSH-Agent automatically in `sshjoin_dial.go`.
43. **Concurrent Fleet Liveness Heartbeat:** Background daemon pinging all nodes in `probeHostsLiveStatus` with exponential backoff on unreachable nodes.
44. **Host Identity Auto-Fingerprinting:** Compute and verify host key SHA256 fingerprints in `ssh_known_hosts_repo.go` to protect against MITM.
45. **Automated Public Key Provisioning:** Add `gitmap ssh copy-id` that pushes local SSH public keys to remote `~/.ssh/authorized_keys` with correct permissions (0600).
46. **Dynamic Port Forwarding & Tunneling:** Add `gitmap ssh tunnel` command to bridge remote split databases and telemetry ports to localhost.
47. **Multi-Hop Bastion Host Routing:** Support `ProxyJump` and jump host directives in `sshconfig.go` for air-gapped server fleets.
48. **Node Resource Telemetry Probe:** Gather CPU, RAM, disk space, and load averages during `gitmap ssh probe` and display in `termtable`.
49. **Zero-Downtime Node Deprecation:** Implement soft-drain status in `ssh_repo.go` allowing in-flight tasks to complete before node removal.
50. **Cross-Platform Path Aliasing Engine:** Map Windows UNC / POSIX paths transparently between fleet nodes via `pathalias.go`.
51. **SSH Config Bidirectional Sync:** Keep `~/.ssh/config` and GitMap's SQLite `SSHHost` table synchronized with conflict resolution.
52. **Encrypted Vault for Remote Sudo Passwords:** Store machine credentials encrypted with AES-256-GCM using host-bound hardware keys in `ssh_vmpass.go`.
53. **Remote Process Tree Termination:** Implement `kill -9 -- -$PGID` signal broadcasting to kill rogue spawned processes on fleet nodes.
54. **Fleet Grouping & Tag Filtering:** Allow targeting subsets of nodes with `--tags linux,gpu,prod` across all `gitmap nodes` commands.
55. **Interactive Remote Shell with PTY:** Add `gitmap ssh attach <alias>` that allocates pseudo-terminal (`ssh -t`) with terminal window resize sync (`SIGWINCH`).
56. **SSH Key Expiry & Rotation Alerts:** Track creation and expiration dates for enrolled SSH keys and notify via `gitmap ssh status`.
57. **Remote OS & Shell Detection Matrix:** Auto-detect Windows OpenSSH, POSIX bash, zsh, PowerShell 7, and ash in `sshjoin_probe.go`.
58. **Audit Logging for Fleet Executions:** Record every remote command execution, invoker, timestamp, and exit code in `installation_split_log.go`.
59. **Automated Agent Bootstrap via SSH:** Push standalone `gitmap` binary to remote nodes via SCP/SFTP and execute silent initialization.
60. **Network Jitter & MTU Diagnostic:** Run ping/traceroute/path-MTU discovery during node enrollment to detect flaky network links.

### Category 4: Cluster Task Delegation, Cross-Node Remote Execution & RPC (`cli/cluster/`, `cli/cmdnodes/`)
61. **Dispatch Wiring for ExecInstall:** Wire `ExecInstall` implementation into `cluster/dispatcher.go` to replace the existing `ExecInstall stubbed` placeholder.
62. **Dispatch Wiring for ExecGit:** Wire `ExecGit` implementation into `cluster/dispatcher.go` replacing `ExecGit stubbed`.
63. **Dispatch Wiring for ExecProj:** Wire `ExecProj` implementation into `cluster/dispatcher.go` replacing `ExecProj stubbed`.
64. **Dispatch Wiring for ExecLifecycle:** Wire `ExecLifecycle` implementation into `cluster/dispatcher.go` replacing `ExecLifecycle stubbed`.
65. **mTLS Authentication for Cluster RPC:** Secure cluster agent communications (`cluster/tls.go`) with mutual TLS certificate validation.
66. **Streamed Subcommand Output Protocol:** Stream remote stdout/stderr in real-time over WebSocket/gRPC rather than buffering whole outputs in `AgentExecReply`.
67. **Distributed Lock Manager (DLM):** Coordinate exclusive multi-node tasks (e.g. repo sync) using SQLite-backed distributed locks.
68. **Node Affinity & Task Scheduling:** Route tasks to nodes based on platform tags (`windows`, `linux`, `darwin`) and current CPU load.
69. **Automatic Re-routing on Node Failure:** Retry failed tasks on secondary healthy nodes when primary node drops connection during execution.
70. **Centralized Job Queue & Work Stealing:** Implement worker pool with work stealing across cluster agents for high-throughput batch operations.
71. **Remote File Synchronization Engine:** Add delta-rsync mechanism in `cluster/distribution.go` to push modified repositories to cluster nodes.
72. **Fleet-Wide Workspace Synchronization:** Sync VS Code workspace definitions across nodes using `nodes_send_projects.go`.
73. **Cluster Health Status Dashboard:** Terminal UI showing node availability, active jobs, completed jobs, and error counts in real time.
74. **Remote Git Pull Errors Aggregation:** Aggregate remote pull errors across all nodes into central root DB with `nodes_pe_all.go`.
75. **Sandboxed Remote Process Execution:** Run remote commands within restricted cgroups (Linux) or Windows Job Objects to prevent CPU starvation.
76. **Remote Git Config Synchronization:** Push global git configuration and identity (`user.name`, `user.email`) to all fleet nodes in a single command.
77. **Bounded Concurrent Node Dispatcher:** Throttle parallel broadcasts to large fleets using a semaphore pool in `cluster/broadcast.go`.
78. **Remote Macro Delegation:** Allow replaying local macros on remote cluster nodes with environment variable overrides.
79. **Node Resource Quota Enforcement:** Reject job assignment if remote node disk free space falls below a safe threshold (e.g. 5GB).
80. **Cluster Agent Auto-Update Protocol:** Distribute new GitMap binary releases to all cluster nodes with automatic rollback on checksum mismatch.

### Category 5: Three-Tier Split-DB Architecture & Schema Migrations (`cli/store/`, `cli/cmddb/`)
81. **Automated Cross-DB Attachment Engine:** Implement dynamic `ATTACH DATABASE` helper in `store.go` to query master and repodb in a single join.
82. **Multi-DB WAL Mode Verification:** Ensure `PRAGMA journal_mode=WAL` and `PRAGMA synchronous=NORMAL` are verified on every split-db connection open.
83. **Split Database Integrity Verifier:** Add `gitmap db check` that runs `PRAGMA integrity_check` on `gitmap.db`, `installation.db`, and all `repodb/*.db`.
84. **Dynamic Schema Migration Versioning:** Track applied migrations in `SchemaMigration` table with forward/rollback migration scripts.
85. **Automated Pipeline DB Size Governance:** Enforce hard 10MB ceiling on `repodb/pipeline.db` with automated rolling FIFO pruning of oldest logs.
86. **Centralized Split Registry Auto-Discovery:** Discover and auto-register detached split databases in `split_database_registry.go` upon startup.
87. **Split DB Transaction Coordinator:** Implement two-phase commit coordinator for transactions spanning `gitmap.db` and split databases.
88. **Safe Database Compaction & Vacuum:** Add `gitmap db vacuum` that performs incremental vacuum (`PRAGMA incremental_vacuum`) during idle periods.
89. **In-Memory SQLite Mode for Unit Tests:** Support `:memory:` connection strings with schema bootstrapping in `dbengine_test_helpers_test.go`.
90. **PascalCase DDL Enforcer in Generator:** Lint generated DDLs to ensure all table names are singular PascalCase and columns adhere to camelCase/PascalCase.
91. **Non-Null Description/Notes Standardization:** Automatically inject nullable `Description`, `Notes`, and `Comments` columns into all new entity tables.
92. **Schema Drift Detector:** Compare live SQLite table schemas against Go struct tags and alert developers during development.
93. **Split DB Encryption at Rest:** Support SQLCipher or AES-256 page-level encryption for sensitive split databases.
94. **Automated Backup & Snapshotting:** Create atomic snapshot copies of `gitmap.db` before running destructive schema migrations.
95. **Orphaned Record Garbage Collector:** Clean up orphaned child records when parent projects or repositories are removed.
96. **Split DB Connection Pool Capping:** Limit open database handles per split DB to prevent `too many open files` errors during mass operations.
97. **REST Visualizer Schema Endpoint:** Expose full split DB schema metadata and table counts in the Swagger/REST visualizer (`rest_visualizer_js1.go`).
98. **Cross-Platform Database Path Normalizer:** Guarantee case-insensitive path comparisons for SQLite file paths on Windows while preserving POSIX cases.
99. **Split Database Export to JSON/Parquet:** Allow exporting split DB tables to compressed JSONL or Parquet for external analytics.
100. **Hot Schema Reloading:** Support adding new nullable columns without restarting background daemons or closing active connections.

### Category 6: Typed SQL Engine, WAL Mode Concurrency & Query Optimization (`cli/dbengine/`, `cli/store/`)
101. **Type-Safe Query Builder Join Compilation:** Add compile-time join validation to ensure joined tables and foreign keys exist in schema constants.
102. **Prepared Statement Cache with LRU Eviction:** Cache compiled SQLite prepared statements in `query_cache.go` to avoid query re-parsing overhead.
103. **Bulk Insert & Upsert Generator:** Add `BatchInsert` / `BatchUpsert` helper to `Repository[T, F]` that chunks records into multi-value INSERT statements.
104. **Monadic Result[T] Query Wrapping:** Standardize all query terminal methods (`FindOne`, `FindAll`) to return `result.Result[T]` Monads.
105. **Query Profiling & Slow Query Logging:** Measure execution times of queries and log queries exceeding 50ms with parameter values.
106. **Zero-Allocation Row Scanner:** Use code-generated column scanning functions instead of reflection-based scanners in `scan_helpers.go`.
107. **Composite Index Optimization:** Add index coverage recommendations based on `WHERE` and `ORDER BY` clauses in `dbengine`.
108. **Safe Transaction Rollback on Panic:** Wrap database transactions with `defer` functions that catch panics and roll back uncommitted changes.
109. **Parameterized Subquery Builder:** Support type-safe nested subqueries (`WHERE col IN (SELECT ...)`) in `query_builder.go`.
110. **Window Function & CTE Support:** Add fluent API for SQLite Common Table Expressions (`WITH ...`) and window functions (`ROW_NUMBER() OVER (...)`).
111. **Soft-Delete Filter Interceptor:** Automatically append `WHERE IsDeleted = 0` to queries for entities that support soft deletion.
112. **SQLite Busy Timeout Configurator:** Ensure all SQLite connections set `PRAGMA busy_timeout = 5000;` to prevent transient locked errors.
113. **Typed Aggregation Helpers:** Add type-safe `Count()`, `Sum()`, `Avg()`, `Min()`, `Max()` methods to `QueryBuilder`.
114. **Read-Only Transaction Replica:** Open separate read-only connection handles (`?mode=ro`) for concurrent read queries to prevent write locks.
115. **Query Explain Plan Analyzer:** Add `builder.Explain()` method that runs `EXPLAIN QUERY PLAN` and returns human-readable execution indexes.
116. **Memory-Efficient Cursor Iteration:** Implement iterator/cursor pattern (`Rows.Next()`) for large table scans to avoid loading thousands of structs into memory.
117. **Automatic Column Renaming Mapping:** Support mapping between database snake_case/PascalCase columns and Go struct fields via struct tags.
118. **Optimistic Locking Support:** Add automatic version column increment (`Version = Version + 1 WHERE Version = ?`) for concurrent updates.
119. **Raw SQL Sanitization & Injection Gate:** Reject raw SQL strings containing unbalanced quotes or suspicious tokens in `compiler_sql.go`.
120. **Query Execution Event Hooks:** Add `BeforeQuery` and `AfterQuery` lifecycle hooks for auditing and distributed tracing.

### Category 7: CI/CD Pipeline Telemetry, Triage & Self-Healing RCA Engine (`cli/cmdpe/`, `cli/cmdpipeline/`, `cli/cmdpullerror/`)
121. **Persistent Commit SHA Cache (`PipelineShaCache`):** Persist parsed pipeline run summaries keyed by Git commit SHA in `repodb/pipeline.db` for instant lookups.
122. **Parallel Pipeline Log Downloading:** Concurrently fetch logs for failing GitHub Actions jobs using bounded worker pool in `pipeline_query_logs.go`.
123. **Zero-Allocation ANSI Log Strip:** Replace heavy regex stripping with single-pass byte scanner in `pipeline_extract_clean.go`.
124. **Automated 4-Part RCA Formatter:** Format pipeline errors into Root Cause, Impact, Proposed Remediation, and Prevention in `pipeline_extract_summarize.go`.
125. **Dynamic ETA Wait Sleep Loop:** Implement progressive backoff and intelligent sleep intervals based on average job duration in `pipeline_ai.go`.
126. **GitHub Token Auto-Rotation & Fallback:** Switch between multiple personal access tokens in `pe_gh.go` when API rate limit (5000 req/hr) is approached.
127. **Targeted Smart Re-run Command:** Add `gitmap pe rerun --failed-only` to trigger re-runs only for failed jobs rather than entire workflows.
128. **Pipeline Log Truncation & Failure Tree:** Render hierarchical tree of failed steps and compact stack traces in `pipeline_failure_tree.go`.
129. **Stale TTL Cache Invalidation on Commit Push:** Invalidate in-memory and SQLite cache immediately when local HEAD moves or new commit is pushed.
130. **Flaky Test Detection Matrix:** Track historical test failure frequencies across commits to flag intermittent and non-deterministic tests.
131. **Pull Request Pipeline Status Correlator:** Correlate open pull request branch heads with workflow run IDs and surface in terminal cards.
132. **Local Runner vs CI Parity Checker:** Compare local test execution flags with `.github/workflows` YAML steps to detect environment discrepancies.
133. **Compact Error Digest for AI Contexts:** Generate ultra-lean markdown error summaries (<100 lines) optimized for AI model context windows.
134. **Real-Time Pipeline Progress Indicator:** Render ANSI spinner with dynamic countdown timer during `gitmap pipeline-ai status -ud`.
135. **Automated Bug Ticket Creation:** Format pipeline failures into structured GitHub issue templates with stack traces and reproduction steps.
136. **Step-Level Duration Profiling:** Track execution duration trends for individual CI/CD steps to detect performance regressions.
137. **Artifact Download & Inspection:** Add `gitmap pe artifacts` to fetch and unpack build artifacts or test reports from GitHub Actions.
138. **Multi-Provider CI/CD Support:** Extend telemetry extractors to support GitLab CI, GitHub Actions, and Woodpecker CI pipelines.
139. **Git Pull Error Remediation Advisor:** Suggest concrete git commands (`git stash`, `git merge --abort`, `git rebase --continue`) on pull conflicts.
140. **Continuous Pipeline Health Score:** Compute repository pipeline health score (pass rate, MTTR, flakiness index) displayed in `gitmap status`.

### Category 8: Git Credential Management, DPAPI Resilience & Headless SSH (`cli/cmdfixauth/`, `cli/gitutil/`)
141. **Windows DPAPI vs WinCredMan Resilience:** Detect broken Windows Credential Manager and auto-switch to `credential.credentialStore dpapi`.
142. **Headless SSH Key Generation Engine:** Generate Ed25519 keypairs non-interactively with empty passphrase and correct file permissions (0600/0700).
143. **Repository-Scoped SSH Command Pinning:** Pin repository to specific SSH key via `git config core.sshCommand` without altering global SSH config.
144. **Automated GitHub SSH Key Upload via API:** Upload generated public key directly to GitHub account using user OAuth token in `fixauth.go`.
145. **OS Clipboard Integration for Public Keys:** Copy generated public keys to Windows, macOS, or Linux X11/Wayland clipboard automatically.
146. **SSH Config Sanitizer & Deduplicator:** Clean corrupted `~/.ssh/config` files, removing duplicate blocks and fixing invalid whitespace.
147. **Uncommitted Changes Classifier:** Classify uncommitted working tree changes into tracked modifications, untracked files, and staged changes in `uncommitted_classifier.go`.
148. **Deterministic Repo ID Generator:** Generate stable, collision-free repository identifiers from remote origin URLs in `repoid.go`.
149. **Automated Stash & Discard Remediators:** Provide safe, dry-run-capable remediation helpers (`remediation_stash.go`, `remediation_discard.go`).
150. **PR Target Branch Autodetector:** Detect target pull request branch (main, master, develop) from remote tracking branches in `pr_detector.go`.
151. **Git Execution Helper Centralization:** Centralize all subprocess `exec.Command("git", ...)` invocations into unified, auditable `gitutil.Exec` wrapper.
152. **SSH Host Key Verification Bypass for CI:** Provide secure flag (`--insecure-accept-host-key`) for non-interactive CI environments.
153. **Multi-Identity Git Profile Switcher:** Add `gitmap profile switch <profile>` that swaps `user.name`, `user.email`, and signing keys across repos.
154. **GPG / SSH Commit Signing Validator:** Verify commit signatures and report unverified or expired cryptographic keys.
155. **Divergent Upstream Delta Calculator:** Compute exact commit counts ahead/behind remote tracking branch in `upstream_delta.go`.
156. **Safe Remote URL Rewriter:** Convert HTTPS URLs with embedded plaintext passwords to standard SSH URLs or credential-helper-backed HTTPS.
157. **Headless Git Auth Health Check:** Add `gitmap check-auth` that tests read/write push permissions against remote origins without mutating code.
158. **Git Credential Helper Unset Safety Guard:** Safeguard against accidental global unsetting of credential helpers when fixing individual repositories.
159. **Windows SSH Agent Service Manager:** Check if `ssh-agent` Windows service is running and offer to start and set startup type to Automatic.
160. **Git Lock File (`index.lock`) Auto-Cleaner:** Detect stale `.git/index.lock` files from killed processes and offer safe removal with PID verification.

### Category 9: Autonomous AI Agent Task Engine & Heatmap Telemetry (`cli/cmdagent/`, `cli/heatmap/`)
161. **Native SQLite Subtask Queue Management:** Enforce atomic task claiming (`gitmap task claim`) with subtask leases to prevent worker collisions.
162. **Heatmap File Telemetry Categorization:** Track modified, created, and deleted files with line counts in `heatmap.db` per agent turn.
163. **AI Agent Action Journaling:** Log every agent action (`read`, `edit`, `create`, `command`) into SQLite task ledger with structured evidence.
164. **Relative Path Sanitation for Heatmap:** Eliminate hardcoded drive letters (`D:\work`) in `heatmap_conflict.go` using workspace-relative resolution.
165. **Agent Collision Detection Matrix:** Check if multiple agents are modifying overlapping files in `agent_claim_collision.go` and warn before claim.
166. **Structured Task Initiation with Slug Helpers:** Generate slugified task IDs and initialize parent task records in `agent_task_sqlite.go`.
167. **Heatmap Tree & Status Terminal Visualizer:** Render terminal UI showing task tree, agent assignments, and progress percentages in `heatmap_tree.go`.
168. **Automated Task Evidence Verification:** Validate that claimed subtasks include verifiable evidence (commit SHA or test pass) before completion.
169. **Agent Heartbeat & Timeout Reclamation:** Automatically reclaim claimed subtasks if agent heartbeat fails for more than 15 minutes.
170. **Heatmap Conflict Resolution Engine:** Reconcile parallel changes from different agents using three-way merge in `heatmap_conflict.go`.
171. **Subtask Export to Markdown Checklist:** Export SQLite subtask state to markdown roadmap documents in `.ai-memory/plans/`.
172. **Agent Execution Budget & Token Tracker:** Track turn counts and model budget consumption per task in `agent_stats.go`.
173. **Web UI Dashboard for Agent Fleets:** Serve interactive HTML/JS dashboard (`agent_ui_server.go`) displaying agent swarm status in real-time.
174. **Structured Error Enveloping in Agent CLI:** Replace raw `fmt.Errorf` in `agent_task_sqlite.go` with strongly typed `appfault.AppError`.
175. **Task Rollup Progress Calculator:** Compute parent task completion percentage dynamically from child subtask states in `agent_task_status.go`.
176. **Heatmap SQL Query CLI Interface:** Expose safe read-only SQL querying on heatmap databases via `gitmap heatmap query "<sql>"`.
177. **Agent Role & Skill Matching:** Route pending subtasks to agents based on matching capability tags (`research`, `coder`, `qa`).
178. **Subtask Failure Reason Classification:** Categorize subtask failures into compile error, test failure, linter violation, or timeout.
179. **Ephemeral Snapshot DB Synchronization:** Export pending uncommitted deltas into temporary SQLite DB for cross-agent synchronization.
180. **Agent Swarm Dynamic Join & Leave:** Allow subagents to register and deregister dynamically from task pool with state persistence.

### Category 10: Interactive Macro Engine, Terminal UI Styling & Power/Loss Prevention (`cli/macro/`, `cli/termout/`, `cli/secrets/`, `cli/power/`)
181. **Macro Variable Interpolation Engine:** Support dynamic parameter expansion (`${REPO_DIR}`, `${BRANCH}`) within recorded macro steps in `expand.go`.
182. **Atomic Macro Storage & Sudo Ownership:** Preserve file ownership when writing macros atomically with `restoreSudoOwnership` in `storage.go`.
183. **Macro Undo & Rollback Journal:** Record reverse actions (e.g. `mkdir` -> `rmdir`, `git checkout` -> previous commit) in `record_undo.go`.
184. **Non-Destructive Safe-RM Command:** Move deleted files to system trash or backup directory instead of permanent unlinking in `safe_rm.go`.
185. **Macro Interactive Recording Shell:** Record shell commands into executable macro with step-by-step confirmation in `macro_record.go`.
186. **Cross-Platform Directory Tracking in Macros:** Maintain correct working directory state across macro steps using `DirTracker` in `execute.go`.
187. **Macro Export/Import via SQLite & YAML:** Allow sharing macros across developers in structured YAML or SQLite formats in `export_sqlite.go`.
188. **Macro Step Timeout & Error Policy:** Configure per-step timeouts and failure behaviors (`continue-on-error`, `abort`, `retry`) in `types.go`.
189. **Dynamic CLI Help Menu Theming:** Render responsive, ANSI-colored help cards with category grouping in `termout/markdown_help.go`.
190. **High-Performance Terminal Table Formatter:** Format data tables with dynamic column width calculation and text truncation in `termout/table.go`.
191. **Salted Cipher Key Derivation (PBKDF2/Argon2):** Upgrade key derivation in `salted_cipher.go` to Argon2id with random salt for secrets storage.
192. **Secret Masking in CLI Terminal Output:** Redact passwords, API keys, and GitHub tokens automatically from terminal output and log files.
193. **Accidental Data Loss Prevention (ADLP) Pre-Push Hook:** Scan commits for hardcoded secrets, AWS keys, and private keys before allowing push.
194. **OS Sleep Prevention Daemon:** Inhibit OS display and system sleep during long-running batch clones and builds via `power/manager.go`.
195. **Windows PowerCfg Timeout Parser:** Parse and restore system sleep timeouts accurately using `powercfg` in `parser_powercfg.go`.
196. **macOS Pmset and Linux Gnome Power Drivers:** Implement native power drivers for macOS (`pmset`) and Linux (`systemd-inhibit` / `gsettings`).
197. **Secure Secrets Resolver with DPAPI:** Store and retrieve GitHub tokens securely on Windows using DPAPI in `secrets/ghtoken_windows.go`.
198. **Terminal Padded Column Layout Engine:** Calculate terminal column widths with ANSI escape code awareness in `termpad.go`.
199. **Interactive Macro TUI Selector:** Provide curses-like terminal menu for selecting, inspecting, and running saved macros in `macro_tui.go`.
200. **Clean Shutdown & Cleanup on SIGINT:** Catch termination signals (`SIGINT`, `SIGTERM`) across all CLI commands to restore power state and flush DBs.

---

# Part 2: 200 Concrete, Actionable Implementation Steps

1. In `cli/scanpipe/scanner.go`, implement `type lockFreeDirQueue struct` using `atomic.Pointer[node]` with head/tail CAS operations to eliminate mutex contention during high-speed directory dispatch.
2. In `cli/scanpipe/pool.go`, add `probeDiskLatency(rootPath string) time.Duration` and dynamically set `workers = clamp(int(targetThroughput / latency), 4, 32)` based on benchmarked seek time.
3. In `cli/scanpipe/scanner_windows.go` (build tag `windows`), invoke `windows.FindFirstFileEx` passing `FindExInfoBasic` and `FIND_FIRST_EX_LARGE_FETCH` to crawl folders at kernel speed without stat calls.
4. In `cli/scanpipe/scanner.go`, introduce `type inodeKey struct { dev, ino uint64 }` and maintain a thread-safe `map[inodeKey]bool` to detect and skip cyclic symlinks.
5. In `cli/scanpipe/scanner.go`, replace `gitdirPrefix` string slice allocations with `bytes.Equal(buf[:7], []byte("gitdir:"))` for zero-allocation git worktree detection.
6. In `cli/scanpipe/scanner.go`, implement `loadGitmapIgnore(root string) []gitignore.Pattern` reading `.gitmapignore` and evaluating paths with standard glob rules.
7. In `cli/scanpipe/walker.go`, store `(path, mtime, inode)` in SQLite `FileMtimeIndex`; query mtime on traversal and skip re-scanning directories when directory mtime is unchanged.
8. In `cli/scanpipe/mapper_options.go`, add `SubmoduleDepth int` to `ScanOptions`; when discovering `.git` files, descend into submodule paths if current depth < `SubmoduleDepth`.
9. In `cli/scanpipe/walker.go`, replace `os.ReadFile` in `isBinaryFile` with `golang.org/x/exp/mmap` to inspect the first 512 bytes without memory allocation.
10. In `cli/scanpipe/sibling_parallel.go`, batch child directories into `[]dirJob` slices and push them to worker channels in a single lock acquisition.
11. In `cli/scanpipe/walker.go`, compute `cespare/xxhash/v2` checksum for scanned source files and populate `FileHash` in `FileInfo` struct.
12. In `cli/scanpipe/scanner.go`, add regex `gitdir: (.*)/\.git/modules/(.*)` to detect absorbed submodules and resolve the root parent workspace path.
13. In `cli/scanpipe/paths.go`, replace string allocations in `filepath.ToSlash` with in-place ASCII replacement `byteReplace(p, '\\', '/')`.
14. In `cli/scanpipe/scanner.go`, add `onDepthBoundary(depth int, path string)` callback triggering `AuditLogger.Warn("Repository at depth boundary", path)`.
15. In `cli/scanpipe/progress.go`, buffer `ScanProgress` snapshots into a 64-capacity channel with drop-oldest semantic to guarantee workers never block on terminal rendering.
16. In `cli/scanpipe/walker.go`, add `--max-file-size` CLI flag binding to `Walker.MaxFileSize` and guard file reading with `if info.Size() > w.MaxFileSize { return nil }`.
17. In `cli/scanpipe/walker.go`, rewrite `scanWriteTimes` to construct parameterized query `SELECT RelativePath, WriteTime FROM RepoFiles WHERE RelativePath IN (?, ?, ...)` chunked by 500 items.
18. In `cli/scanpipe/pool.go`, initialize an in-memory SQLite connection (`file::memory:?cache=shared`) to collect discovered repos before dumping to target format.
19. In `cli/scanpipe/walker.go`, replace `excludedDirMap` with a trie-based lookup struct `dirTrie` for O(k) prefix matching against excluded directory names.
20. In `cli/scanpipe/pool.go`, wrap worker loops with `select { case <-ctx.Done(): return }` and enforce a 500ms `sync.WaitGroup` wait with timeout.
21. In `cli/cloner/concurrent.go`, detect `HTTP 429` or `fatal: early EOF` in `cloneOutcome.result.Stderr` and dynamically decrement `p.Workers` down to 1 with exponential backoff.
22. In `cli/clonenow/clonenow.go`, add `InitClonenowJournal(dbPath string)` creating `ClonenowProgress` table recording `(RepoName, Status, CompletedAt)` to resume interrupted runs.
23. In `cli/cloner/strategy.go`, execute `git clone --depth 1` and enqueue a background goroutine executing `git fetch --unshallow` upon successful checkout.
24. In `cli/cloner/runners.go`, inject `cmd.Env = append(os.Environ(), "GIT_PROTOCOL=version=2")` into all Git exec invocations.
25. In `cli/cloner/lfs_retry.go`, set `GIT_LFS_SKIP_SMUDGE=1` on initial clone and run `git lfs pull` as an isolated retryable post-step.
26. In `cli/clonenow/clonenow.go`, define `type URLRewriter func(raw string) string` and parse regex rewrite flags (e.g. `s#git@github.com:#https://#`).
27. In `cli/cloner/safe_pull_conflicts.go`, implement `AutoResolveConflict(repo string, policy MergePolicyType) error` handling `theirs`, `ours`, or `abort`.
28. In `cli/clonenow/parse_schema.go`, calculate SHA256 checksum of input CSV/JSON file and verify required columns (`RepoName`, `HTTPSUrl`, `SSHUrl`, `RelativePath`).
29. In `cli/cloner/concurrent.go`, stream `model.ScanRecord` from an `io.Reader` directly into `cloneJob` channel avoiding `[]model.ScanRecord` heap allocation.
30. In `cli/cloner/cache.go`, create SQLite table `CloneQuarantine` and insert failed records with stderr output and error classification enum.
31. In `cli/cloner/cloner.go`, add `IsMirror bool` and `IsBare bool` flags to `CloneOptions` translating to `--mirror` and `--bare` git clone arguments.
32. In `cli/clonenow/execute.go`, inspect target volume available bytes using `windows.GetDiskFreeSpaceEx` or `unix.Statfs` and verify `freeBytes > totalEstimatedRepoBytes`.
33. In `cli/cloner/runners.go`, inject `-c credential.helper="" -c core.askPass=""` when running in headless mode to immediately fail rather than hang on credentials.
34. In `cli/cloner/progress_stream.go`, implement ANSI multi-progress bar widget using `termout.TerminalWidth()` to show concurrent worker throughput.
35. In `cli/clonenow/execute_hooks.go`, inspect parsed rows for `.git` pointers and run `git worktree add <path> <branch>` automatically.
36. In `cli/cloner/strategy.go`, support `--sparse` flag invoking `git clone --filter=blob:none --sparse <url>` and applying cone patterns.
37. In `cli/cloner/cache.go`, check remote HEAD SHA via `git ls-remote <url> HEAD` and invalidate `CloneCache` when remote commit differs from cached commit.
38. In `cli/cloner/strategy.go`, add `CloneOptions.SubmoduleJobs int` and pass `--recurse-submodules --jobs=<N>` to `git clone`.
39. In `cli/cloner/runners.go`, inject `-c http.lowSpeedLimit=1000 -c http.lowSpeedTime=30` into `git clone` to abort hung transfers automatically.
40. In `cli/clonenow/execute_idempotent.go`, check `git status --porcelain` on existing repositories and return `ErrDirtyWorktree` before running pull.
41. In `cli/cmdssh/ssh_stream.go`, configure OpenSSH client with `-o ControlMaster=auto -o ControlPath=~/.ssh/gitmap-%r@%h:%p -o ControlPersist=10m` to multiplex connections.
42. In `cli/cmdssh/sshjoin_dial.go`, implement `dialNodeWithFallbackChain` iterating through `[Ed25519, RSA, Pageant, Agent, Password]` until successful handshake.
43. In `cli/cmdssh/sshjoin_probe.go`, spawn a persistent ticker goroutine (`15s` period) updating `store.SSHHost.Status` in background.
44. In `cli/store/ssh_known_hosts_repo.go`, store host public key fingerprints as `SHA256:<hash>` and verify matching host key before opening SSH sessions.
45. In `cli/cmdssh/sshcopy.go`, implement `gitmap ssh copy-id <alias>` reading public key and appending to remote `~/.ssh/authorized_keys` with `chmod 600`.
46. In `cli/cmdssh/sshbind.go`, add `Tunnel(localPort, remotePort int, host string)` establishing bidirectional port forwarding via SSH client connection.
47. In `cli/cmdssh/sshconfig.go`, parse and honor `ProxyJump` directives by establishing an SSH client connection to the bastion host before dialing target node.
48. In `cli/cmdssh/sshjoin_probe.go`, execute `cat /proc/loadavg; free -m; df -h /` on remote nodes and parse metrics into `NodeSystemStats` struct.
49. In `cli/store/ssh_repo.go`, add `StatusDraining` enum and reject new command dispatching to draining nodes while allowing active processes to exit.
50. In `cli/cluster/pathalias.go`, implement `TranslateRemotePath(localPath string, targetOS string) string` converting drive letters to mount points.
51. In `cli/cmdssh/ssh_sync_nodes.go`, implement two-way reconciliation between `~/.ssh/config` file blocks and `gitmap.db` `SSHHost` rows.
52. In `cli/cmdssh/ssh_vmpass.go`, encrypt remote credentials using `crypto/aes` GCM mode with machine-guid derived key before storing in `SSHHost.EncryptedPassword`.
53. In `cli/cmdssh/sshexec.go`, wrap remote execution commands in process groups (`setsid` or `JobObject`) to allow clean termination of children via process group kill.
54. In `cli/cmdnodes/nodes_filter.go`, implement `FilterNodesByTags(tags []string)` matching against comma-separated tags in `SSHHost.Tags`.
55. In `cli/cmdssh/ssh.go`, implement `RunInteractivePTY` requesting PTY (`xterm-256color`) and binding `os.Stdin` / `os.Stdout` with raw terminal mode.
56. In `cli/cmdssh/sshstatus.go`, inspect SSH key metadata and emit warnings if key expires within 30 days or uses deprecated 1024-bit RSA.
57. In `cli/cmdssh/sshjoin_probe.go`, run `uname -s; $PSVersionTable.PSVersion` and record `HostOSType` and `DefaultShellType` in node record.
58. In `cli/store/installation_split_log.go`, write audit record `(NodeAlias, CommandText, ExitCode, DurationMs, InvokerUser, Timestamp)` for every fleet execution.
59. In `cli/cmdnodes/nodes_deploy_agm.go`, copy local `gitmap` binary to remote `/usr/local/bin/` or `C:\ProgramData\gitmap\` using SFTP client.
60. In `cli/cmdssh/sshjoin_dial.go`, run ICMP ping and TCP SYN probe measuring round-trip time and packet loss, persisting metrics in `SSHHost.LatencyMs`.
61. In `cli/cluster/dispatcher.go`, wire `case db.CommandKindInstall:` to call `ExecInstall(ctx, node, splitPackages(subCmd.RawArg))` and store package install results.
62. In `cli/cluster/dispatcher.go`, wire `case db.CommandKindGitPull, db.CommandKindGitPush:` to invoke `ExecGit(ctx, node, subCmd.Kind, subCmd.RawArg)`.
63. In `cli/cluster/dispatcher.go`, wire `case db.CommandKindProjRun, db.CommandKindProjCreateCICD:` to call `ExecProj(ctx, node, subCmd.Kind, subCmd.RawArg)`.
64. In `cli/cluster/dispatcher.go`, wire `case db.CommandKindRestart, db.CommandKindShutdown:` to invoke `ExecLifecycle(ctx, node, subCmd.Kind)`.
65. In `cli/cluster/tls.go`, configure `tls.Config` with `ClientAuth: tls.RequireAndVerifyClientCert` and embed cluster root CA certificate.
66. In `cli/cluster/stream.go`, implement `StreamOutput(ctx context.Context, subCmd ClusterSubCommand) (io.Reader, error)` using chunked RPC response streaming.
67. In `cli/cluster/distribution.go`, create SQLite table `ClusterLock` with `(LockName, OwnerNodeId, AcquiredAt, ExpiresAt)` and atomic acquisition queries.
68. In `cli/cluster/node_resolver.go`, calculate placement score `score = (100 - CPUUsage) * 0.5 + (FreeRAMMB / 1024) * 0.5` and select best candidate node.
69. In `cli/cluster/dispatcher.go`, catch node connection timeouts and automatically re-queue `ClusterSubCommand` to the next available healthy node.
70. In `cli/cluster/pool.go`, implement work-stealing deque across cluster worker channels to balance batch execution across heterogeneous hosts.
71. In `cli/cluster/distribution.go`, compute block-level hashes of modified files and transmit only differing byte ranges to remote node targets.
72. In `cli/cmdnodes/nodes_send_projects.go`, serialize local `vscode_project` SQLite table and push workspace settings to remote nodes via RPC.
73. In `cli/cluster/ui.go`, render dynamic terminal dashboard showing active nodes, thread utilization, memory, and completed tasks using `termout`.
74. In `cli/cmdnodes/nodes_pe_all.go`, query `repodb/pipeline.db` on all reachable nodes via RPC and aggregate failed pull reports into local database.
75. In `cli/cluster/exec_bash.go`, prepend `systemd-run --scope -p CPUQuota=80% -p MemoryMax=4G` when spawning cluster jobs on Linux nodes.
76. In `cli/cluster/exec_git.go`, distribute user Git identity (`git config --global user.name` and `user.email`) to remote node prior to pull/push.
77. In `cli/cluster/broadcast.go`, wrap remote dispatch loop in a buffered semaphore channel `chan struct{}` capped at `maxBroadcastConcurrency`.
78. In `cli/cluster/exec_cmd.go`, serialize local macro definitions to JSON and invoke `gitmap macro exec --json-payload` on remote cluster node.
79. In `cli/cluster/preflight.go`, check remote node available disk space before delegating build jobs; fail early if free space < 5GB.
80. In `cli/cluster/exec_update.go`, compare local and remote GitMap version strings; trigger self-update over RPC if remote version is older.
81. In `cli/store/store.go`, add `AttachSplitDatabase(db *sql.DB, alias, path string) error` executing `ATTACH DATABASE ? AS ?` with parameter sanitation.
82. In `cli/store/sqlite_config.go`, execute `PRAGMA journal_mode=WAL; PRAGMA synchronous=NORMAL;` inside `OpenSQLiteDB` helper.
83. In `cli/store/store.go`, add `CheckIntegrity() map[string]string` running `PRAGMA integrity_check` on all registered split databases.
84. In `cli/store/migrations.go`, create `SchemaMigration` table recording `(MigrationId, Name, AppliedAt, BatchId)` with transactional execution.
85. In `cli/pipelinedb/pipeline_db_purge.go`, check database file size on startup; if `size > 10MB`, execute rolling prune query `DELETE FROM PipelineRuns WHERE Timestamp < ...`.
86. In `cli/store/split_database_registry_sync.go`, scan known directories for `*.db` files and automatically register discovered databases into `SplitDatabaseRegistry`.
87. In `cli/store/transaction.go`, implement `TwoPhaseCommit(masterTx, splitTx *sql.Tx) error` ensuring both transactions commit or both roll back.
88. In `cli/store/store.go`, add `VacuumIncremental(db *sql.DB, pages int) error` running `PRAGMA incremental_vacuum(?)` during low-load intervals.
89. In `cli/dbengine/dbengine_test_helpers_test.go`, initialize in-memory SQLite connection `file::memory:?cache=shared&mode=memory` with full schema migrations.
90. In `03-ai-scripts/35-db-struct-enum-generator.py`, add linter rule verifying that all SQLite `CREATE TABLE` statements use singular PascalCase names.
91. In `cli/store/migrations.go`, add validation rule checking that all new table DDLs include `Description TEXT NULL`, `Notes TEXT NULL`, and `Comments TEXT NULL`.
92. In `cli/store/store.go`, implement `VerifySchemaDrift(entity any, table string) []string` comparing struct fields against `PRAGMA table_info`.
93. In `cli/store/store.go`, support SQLCipher connection URI `file:path.db?_crypto_key=hexkey` when building with cipher tag.
94. In `cli/store/store.go`, implement `SnapshotDatabase(srcPath, backupPath string) error` using SQLite online backup API (`sqlite3_backup`).
95. In `cli/store/repo_sanitize.go`, execute cascade cleanup query removing records in `ProjectLocation` where `ProjectId NOT IN (SELECT Id FROM Project)`.
96. In `cli/store/store.go`, configure `sql.DB.SetMaxOpenConns(5)` and `SetMaxIdleConns(2)` per split-database connection pool.
97. In `cli/cmddb/rest_server.go`, add HTTP endpoint `GET /api/v1/schemas` returning registered split databases, table names, and column counts in JSON.
98. In `cli/store/split_db_path.go`, normalize all stored database paths using `filepath.Clean(filepath.ToSlash(strings.ToLower(path)))` on Windows.
99. In `cli/store/store.go`, add `ExportTableToJSONL(tableName string, out io.Writer) error` streaming table rows as JSON lines.
100. In `cli/store/store.go`, implement `AddColumnIfNotExists(db *sql.DB, table, colName, colType string) error` inspecting table columns prior to `ALTER TABLE`.
101. In `cli/dbengine/query_build_clauses.go`, validate join target strings against registered repository table constants; return `appfault.NewValidation` on unknown joins.
102. In `cli/dbengine/query_cache.go`, implement `type StatementCache struct` backed by a fixed-size LRU cache storing `*sql.Stmt` instances.
103. In `cli/dbengine/repository.go`, add `BatchInsert(ctx context.Context, entities []T) result.Result[int]` generating multi-row insert statements.
104. In `cli/dbengine/query_builder.go`, refactor `FindOne` and `FindAll` to return `result.Result[T]` and `result.ResultList[T]` Monads.
105. In `cli/dbengine/executor.go`, measure query duration with `time.Since(start)`; if duration > 50ms, emit slow query log event with SQL and args.
106. In `cli/dbengine/scan_helpers.go`, implement code-generated `ScanRowDirect(rows *sql.Rows, dest *T) error` eliminating `reflect.ValueOf` overhead.
107. In `cli/dbengine/query_builder_compile.go`, analyze `Where` clause columns and suggest `CREATE INDEX idx_table_col ON table(col)` in debug logs.
108. In `cli/dbengine/executor.go`, implement `WithTx(ctx context.Context, fn func(tx *sql.Tx) error) (err error)` wrapping call with `recover()` that executes `tx.Rollback()`.
109. In `cli/dbengine/query_builder_filters.go`, add `WhereInSubquery(field F, sub *QueryBuilder[any, string]) *QueryBuilder[T, F]`.
110. In `cli/dbengine/query_builder.go`, add `WithCTE(name string, sub *QueryBuilder[any, string]) *QueryBuilder[T, F]` compiling to `WITH <name> AS (...)`.
111. In `cli/dbengine/query_builder_compile.go`, check if `T` implements `SoftDeletable`; if true, automatically inject `AND IsDeleted = 0` into where clauses.
112. In `cli/dbengine/executor.go`, execute `PRAGMA busy_timeout = 5000;` immediately upon opening connection to prevent immediate `sqlite3_busy` failures.
113. In `cli/dbengine/query_builder.go`, add `Count(ctx context.Context) (int64, error)` compiling to `SELECT COUNT(*) FROM ...` with active filters.
114. In `cli/dbengine/repository.go`, support opening separate `readConn *sql.DB` with `?mode=ro&_journal_mode=WAL` for concurrent non-blocking reads.
115. In `cli/dbengine/query_builder.go`, add `Explain(ctx context.Context) (string, error)` prepending `EXPLAIN QUERY PLAN` to generated SQL.
116. In `cli/dbengine/repository.go`, implement `ForEach(ctx context.Context, fn func(item T) bool) error` streaming results without loading all records into memory.
117. In `cli/dbengine/scan_helpers.go`, cache struct field names and db column tags in a static `sync.Map` during package initialization.
118. In `cli/dbengine/repository.go`, implement `UpdateOptimistic(ctx context.Context, entity T) error` appending `WHERE Version = ?` and checking `RowsAffected() == 1`.
119. In `cli/dbengine/compiler_sql.go`, add AST token scanner verifying that raw SQL inputs contain no unescaped semicolons or unbalanced quotes.
120. In `cli/dbengine/repository.go`, define `type QueryHook func(ctx context.Context, sql string, args []any)` invoked before and after every execution.
121. In `cli/cmdpipeline/pipeline_sha_cache.go`, store parsed `PipelineErrorLogsPayload` in SQLite table `PipelineShaCache` keyed by `(RepoSlug, CommitSha)`.
122. In `cli/cmdpipeline/pipeline_query_logs.go`, spawn bounded goroutine pool (`maxWorkers = 4`) fetching failed job logs concurrently via GitHub API.
123. In `cli/cmdpipeline/pipeline_extract_clean.go`, replace regex `\x1b\[[0-9;]*[a-zA-Z]` with single-pass byte scanner skipping ANSI escape sequences.
124. In `cli/cmdpipeline/pipeline_extract_summarize.go`, implement `Build4PartRCA(errs []PipelineError)` formatting Root Cause, Impact, Remediation, and Prevention.
125. In `cli/cmdpipeline/pipeline_ai.go`, calculate dynamic wait time `min(max(etaSeconds * 0.8, 15), 180)` and sleep before querying workflow status.
126. In `cli/cmdpe/pe_gh.go`, implement token bucket pool switching to secondary `GITHUB_TOKEN` when `X-RateLimit-Remaining < 100`.
127. In `cli/cmdpe/pe_rerun.go`, invoke `POST /repos/{owner}/{repo}/actions/runs/{run_id}/rerun-failed-jobs` rather than re-running entire workflow.
128. In `cli/cmdpipeline/pipeline_failure_tree.go`, render ASCII failure hierarchy connecting workflow -> failing job -> step -> error summary line.
129. In `cli/cmdpipeline/pipeline_cache_eval.go`, check current local commit SHA; if `localSha != cachedSha`, invalidate cache immediately.
130. In `cli/cmdpipeline/pipeline_history.go`, compute test failure variance `failureRate = failedRuns / totalRuns` and tag tests with `flaky` badge if `0.05 < failureRate < 0.95`.
131. In `cli/cmdpipeline/pipeline_logs_repo.go`, query open PRs via GitHub API and map PR commit SHAs to active GitHub Actions workflow runs.
132. In `cli/cmdpipeline/pipeline_cicd_checker.go`, parse `.github/workflows/*.yml` test commands and verify that local test scripts execute identical flags.
133. In `cli/cmdpipeline/pipeline_format_cmd.go`, generate compact markdown digest (<100 lines) formatting failing test names and stack traces.
134. In `cli/cmdpipeline/pipeline_dynamic_timeline.go`, implement ANSI countdown timer rendering `[ETA: %02dm %02ds]` updated every 1s.
135. In `cli/cmdpe/pe_triage.go`, add `--create-issue` flag formatting RCA findings into GitHub issue payload and posting via `POST /repos/.../issues`.
136. In `cli/cmdpipeline/pipeline_timing_test.go`, record duration per workflow step across runs and alert if any step duration increases by > 25%.
137. In `cli/cmdpe/pe_dispatch.go`, add `pe artifacts` subcommand invoking `GET /repos/{owner}/{repo}/actions/artifacts` and extracting downloaded zip files.
138. In `cli/cmdpipeline/pipeline.go`, define `type PipelineProvider interface` and implement adapters for GitLab CI (`/api/v4/projects/.../pipelines`) and Woodpecker.
139. In `cli/cmdpullerror/pull_error_cmd.go`, parse git pull conflict markers and output step-by-step resolution script (`git stash; git pull; git stash pop`).
140. In `cli/cmdpipeline/pipeline_status.go`, compute rolling 14-day pipeline health score `(passedRuns / totalRuns) * 100` and display in `gitmap status`.
141. In `cli/cmdfixauth/fixcredential.go`, test Windows Credential Manager functionality; on failure, execute `git config --global credential.credentialStore dpapi`.
142. In `cli/cmdfixauth/fixauth.go`, invoke `ssh-keygen -t ed25519 -N "" -f <path> -C <email>` non-interactively with `0600` file permissions.
143. In `cli/cmdfixauth/fixauth.go`, execute `git config core.sshCommand "ssh -i <path> -F /dev/null -o IdentitiesOnly=yes"` in current working repository.
144. In `cli/cmdfixauth/fixauth.go`, add `--upload-to-github` flag invoking `POST /user/keys` with generated public key content.
145. In `cli/cmdfixauth/fixauth.go`, invoke OS clipboard utilities (`clip.exe` on Windows, `pbcopy` on macOS, `xclip` on Linux) to copy public key text.
146. In `cli/cmdssh/sshconfig.go`, implement `SanitizeSSHConfigFile(path string)` removing duplicated `Host` blocks and standardizing indentation.
147. In `cli/gitutil/uncommitted_classifier.go`, parse `git status --porcelain=v2` into structured `UncommittedChanges{Tracked, Untracked, Staged}` slices.
148. In `cli/gitutil/repoid.go`, generate deterministic SHA256 repo ID from normalized remote URL (`git remote get-url origin`).
149. In `cli/gitutil/remediation_stash.go`, implement `SafeStashChanges(repoPath string) (stashRef string, err error)` verifying clean state before returning.
150. In `cli/gitutil/pr_detector.go`, inspect `refs/remotes/origin/HEAD` to determine default branch (`main` or `master`) without network calls.
151. In `cli/gitutil/gitutil.go`, replace fragmented `exec.Command("git", ...)` invocations with `gitutil.Exec(ctx, args...)` enforcing structured error wrapping.
152. In `cli/cmdssh/sshjoin_dial.go`, support `StrictHostKeyChecking=no` flag when environment variable `GITMAP_SSH_INSECURE_HOST_KEY=1` is set.
153. In `cli/store/profile.go`, implement `SwitchGitProfile(name string)` updating local repository `user.name`, `user.email`, and `user.signingkey`.
154. In `cli/gitutil/gitutil.go`, execute `git verify-commit HEAD` and return boolean `isValidSignature` and signer fingerprint.
155. In `cli/gitutil/upstream_delta.go`, parse `git rev-list --left-right --count HEAD...@{u}` returning `(ahead int, behind int)`.
156. In `cli/store/repo_sanitize.go`, strip embedded basic auth credentials (`http://user:pass@host`) from remote URLs and save sanitized URL.
157. In `cli/cmdfixauth/fixauth.go`, add `gitmap check-auth` executing `git ls-remote --exit-code <url>` to verify read access without altering local files.
158. In `cli/cmdfixauth/fixcredential.go`, add guard checking if `--global` was explicitly requested before running `git config --unset`.
159. In `cli/cmdfixauth/fixcredential.go`, query Windows service `ssh-agent`; if stopped, execute `Set-Service ssh-agent -StartupType Automatic; Start-Service ssh-agent`.
160. In `cli/gitutil/gitutil.go`, inspect `.git/index.lock`; if process holding lock is no longer alive, prompt or remove stale lock file.
161. In `cli/cmdagent/agent_task_sqlite.go`, implement `gitmap task claim` using atomic transaction `UPDATE Task SET AssignedAgent = ?, Status = 'claimed' WHERE Id = ? AND Status = 'pending'`.
162. In `cli/heatmap/heatmap_rc_update.go`, parse `git diff --stat` into categorized records `(FilePath, Insertions, Deletions, OperationType)` stored in `heatmap.db`.
163. In `cli/cmdagent/agent_log.go`, record each agent turn into `AgentActionLog` table storing `(TaskId, SubtaskId, ActionType, FileTarget, EvidenceText, Timestamp)`.
164. In `cli/heatmap/heatmap_conflict.go`, replace hardcoded `D:\work\repo-cache` strings with `filepath.Join(os.Getenv("GITMAP_WORKSPACE"), "repo-cache")`.
165. In `cli/cmdagent/agent_claim_collision.go`, query active subtasks for overlapping `FileTarget` paths and return warning before claim.
166. In `cli/cmdagent/agent_task_slug.go`, implement `GenerateTaskSlug(title string) string` converting titles to lowercase hyphenated slugs.
167. In `cli/heatmap/heatmap_tree.go`, render ASCII hierarchical tree showing parent task -> subtasks -> assigned agents and status colors.
168. In `cli/cmdagent/agent_subtask_resolve.go`, require non-empty `Evidence` string matching commit SHA format or test execution log before marking complete.
169. In `cli/cmdagent/agent_task_status.go`, spawn monitor checking for claimed subtasks with `last_heartbeat < NOW() - 15 minutes` and revert to `pending`.
170. In `cli/heatmap/heatmap_conflict.go`, implement automated git three-way merge resolution for non-overlapping hunk modifications across agents.
171. In `cli/cmdagent/agent_task_export.go`, serialize SQLite subtask hierarchy to markdown checklist format compatible with `.ai-memory/plans/*.md`.
172. In `cli/cmdagent/agent_stats.go`, calculate total execution time, turn count, and file mutation count per agent task.
173. In `cli/cmdagent/agent_ui_server.go`, serve embedded static assets (`agent_ui_assets.go`) over HTTP on `localhost:3842` for visual swarm telemetry.
174. In `cli/cmdagent/agent_task_sqlite.go`, replace `fmt.Errorf("unknown task subcommand: %s", sub)` with `appfault.NewValidation("agent_task", "E_UNKNOWN_SUBCOMMAND", msg)`.
175. In `cli/cmdagent/agent_task_status.go`, update parent task status to `completed` automatically when all child subtasks reach `completed`.
176. In `cli/heatmap/heatmap_cmd.go`, implement `ExecuteHeatmapQuery` validating that query string starts with `SELECT` before executing against `heatmap.db`.
177. In `cli/cmdagent/agent_config.go`, define `type AgentRoleConfig struct` specifying capabilities (`research`, `coder`, `qa`) and match against subtask tags.
178. In `cli/cmdagent/agent_task_sqlite_subtask.go`, add `FailureReasonType` enum column to `Subtask` table (`CompileError`, `TestFailure`, `Timeout`).
179. In `cli/heatmap/heatmap_conflict.go`, implement `ExportPendingChangesToTempDB` writing delta tables into temporary SQLite database for peer sync.
180. In `cli/cmdagent/agent_ps.go`, implement dynamic agent registration heartbeat storing process ID, hostname, and active claimed subtask.
181. In `cli/macro/expand.go`, implement variable expander replacing `${VAR_NAME}` with environment variables or macro execution arguments.
182. In `cli/macro/storage.go`, invoke `restoreSudoOwnership` on written macro files when executed under `SUDO_UID` to prevent root-owned user files.
183. In `cli/macro/record_undo.go`, implement reverse action generator mapping `mkdir` -> `rmdir`, `touch` -> `rm`, and `git checkout` -> `git checkout -`.
184. In `cli/macro/safe_rm.go`, move deleted files into `.gitmap/trash/<timestamp>` with restore metadata instead of direct unlinking.
185. In `cli/cmdmacro/macro_record.go`, launch subshell recording entered commands and prompt `Save to macro? [y/n]` upon subshell exit.
186. In `cli/macro/execute.go`, ensure `DirTracker.Chdir` tracks directory stack and restores initial working directory upon macro completion or abort.
187. In `cli/macro/export_sqlite.go`, implement `ExportMacroToSQLite(m *Macro, db *sql.DB)` and `ImportMacroFromSQLite(name string, db *sql.DB)`.
188. In `cli/macro/types.go`, add `StepTimeout time.Duration` and `OnErrorPolicyType` (`Abort`, `Continue`, `Retry`) to `MacroStep` struct.
189. In `cli/termout/markdown_help.go`, parse markdown headings into ANSI styled card boxes with colored borders and categorized command listings.
190. In `cli/termout/table.go`, calculate max column widths dynamically based on `termout.TerminalWidth()` and truncate middle of long paths with `...`.
191. In `cli/secrets/salted_cipher.go`, replace static password hashing with `golang.org/x/crypto/argon2` generating 32-byte AES key from password and random salt.
192. In `cli/termout/writer.go`, wrap terminal output writer with regex filter redacting patterns matching `ghp_[A-Za-z0-9]{36}` or `AKIA[0-9A-Z]{16}`.
193. In `cli/gitutil/gitutil.go`, implement pre-push scanner inspecting commit diffs for high-entropy strings and known API key signatures before pushing.
194. In `cli/power/manager.go`, implement `power.InhibitSleep(reason string) (restore func(), err error)` holding OS power assertion during long batch operations.
195. In `cli/power/parser_powercfg.go`, parse `powercfg /q` output extracting current AC/DC display and system sleep timeouts, restoring on exit.
196. In `cli/power/driver_darwin.go` and `driver_linux.go`, implement macOS `caffeinate` process holder and Linux `systemd-inhibit` D-Bus calls.
197. In `cli/secrets/ghtoken_windows.go`, use Windows DPAPI `CryptProtectData` and `CryptUnprotectData` to store and retrieve GitHub PAT tokens securely.
198. In `cli/termout/termpad.go`, strip ANSI color escape sequences before calculating string display width to ensure perfect table padding.
199. In `cli/cmdmacro/macro_tui.go`, implement interactive terminal cursor menu allowing users to browse, view steps, and press Enter to execute macros.
200. In `cli/root.go`, register `signal.Notify(c, os.Interrupt, syscall.SIGTERM)` calling cleanup registry that flushes open SQLite DBs and restores OS power settings.

---

## Verification & Read-Only Invariant Compliance

- **Read-Only Verification:** All exploration was performed via read-only tools (`view_file`, `pwsh Get-ChildItem`, `git log`, `gitmap lf`). Zero files were created, edited, or deleted in the target workspace.
- **Rule Verification:** All formulate improvements adhere strictly to GitMap guidelines:
  1. Sovereign error handling via `appfault.AppError` and Monadic `Result[T]` wrappers.
  2. Affirmative boolean prefixes (`is`, `has`) and zero negative booleans.
  3. 15-line function decomposition constraints.
  4. Multi-tier SQLite Split-DB architecture (`gitmap.db`, `installation.db`, `repodb/pipeline.db`) with WAL mode and singular PascalCase tables.

*Report complete and ready for Parent Task 287 execution plan synthesis.*