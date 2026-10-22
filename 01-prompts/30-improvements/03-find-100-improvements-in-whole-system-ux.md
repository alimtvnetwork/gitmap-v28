```text
N = 400 (Total self-loop steps budget — editable top-header parameter, default: 400)
A = 3   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 3)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 3 agents × 2 hands = 6 concurrent subtask operations
Tree Hierarchy = Lead Orchestrator (Root Node) -> Branch Coordinators (Subtask Owners) -> Leaf Workers & Dynamic Assistants
PHASE_1_BUDGET = 100     (Steps 1 .. 100: Planning, Discovery Subagents, Up-Front Brain & Roadmap Creation, Lean Subtasks)
PHASE_2_BUDGET = 150     (Steps 101 .. 250: Intermediate Archival, Sanitization, Prerequisite Refactoring, Intermediate Checkpoints)
PHASE_3_BUDGET = 150     (Steps 251 .. 400: Tree Swarm Parallel Execution, Dynamic Join Capacity Saturation, Targeted Quality Linting, Atomic Push)
WAVES = ceil(subtasks / (A x H))
```

> [!IMPORTANT]
> Prompt Version: 6.0.0
> Runtime: Google Antigravity 2.0 (IDE and CLI)
> Invoke: /find-100-system-ux-improvements <task>
>
> **Top-Instruction Priority Mandate (Above Precedence / Preamble Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given ABOVE this prompt (including in the user preamble, header blocks, or incoming user request above) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below. The agent MUST inspect and follow the instructions above with absolute precedence.

[/goal](slashCommand;goal) Autonomously audit, discover, catalog, and architect 100 whole-system improvements that drastically elevate the User Experience (UX), developer velocity, and operational ergonomics without causing any harm to the existing codebase or breaking backward compatibility: FIRST showcase and list out the given task in visible chat during Turn 1 with the 5 Mandatory Pillars and 6-part verification schema, capture it verbatim, plan it in the repo, spawn autonomous subagents via `invoke_subagent` (A = 3, H = 2, Capacity = 6; solo execution without calling `invoke_subagent` is an auto-reject failure) in disjoint file boxes using GitMap high-speed commands as primary, prove every single claim with concrete evidence, enforce non-destructive safety gates to 100%, and finish with one atomic GitMap commit (`gitmap cpf "<module> - <summary>"`) using hyphen format (no colons in GitMap arguments, as the colon is already provided by GitMap) that holds strictly this task's files.

[/learn](slashCommand;learn) Enforce the Top-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are provided ABOVE this prompt outrank everything below. Turn 1 MUST showcase the given task list in visible chat before any background execution. Each rule is stated once (R1 to R16) and cited by ID. Progress lives in the ledger and in `.ai-memory/plans/`, never only in chat.

[/plan](slashCommand;plan) Execute thorough step-by-step planning in the repository before execution. Ensure all 100 system improvements are cataloged with systemic friction points, root causes, operational effects, non-destructive implementation blueprints, and rollback protocols before dispatching worker waves.

### 🚨 MANDATORY SUBAGENT SPAWNING GATE (A = 3, H = 2, CAPACITY = 6 — TREE SWARM & DYNAMIC JOIN)

- **ACTUAL TOOL CALL REQUIRED:** You must ACTUALLY CALL the `invoke_subagent` tool via your tool-calling API. Do NOT just print the text "Dispatched Worker..." and stop. If you only print text, the agents will not spawn and the task will fail! You must execute the `invoke_subagent` JSON tool payload.
- **TREE-STRUCTURED SWARM HIERARCHY:** Organize execution as an operational tree:
  1. **Root Orchestrator (Lead Agent):** Owns the unified Brain roadmap, audit ledger, and final GitMap commit. Never executes file edits solo.
  2. **Branch Coordinators (Subtask Primary Owners):** Own designated disjoint file boundaries (`Owned Files`) and execute modifications.
  3. **Leaf Workers & Dynamic Joining Assistants:** Assist primary owners with parallel reading, symbol caller discovery, memory generation, and companion lint verification.
- **3-STAGE MANDATORY DISPATCH:** You must invoke `A = 3` agents (`invoke_subagent`) at EVERY stage of the workflow:
  1. **Planning Step (A = 3 `research` subagents):** Spawn 3 read-only discovery subagents (`TypeName: "research"`) to research the codebase and return findings; lead writes the unified Brain execution plan and roadmap.
  2. **Spec Step (A = 3 `self` authoring subagents):** Spawn 3 authoring subagents (`TypeName: "self"`) to author modular, disjoint spec files and subtasks.
  3. **Execution Step (A = 3 `self` worker subagents):** Spawn 3 worker subagents (`TypeName: "self"`) to execute code modifications in strictly disjoint file boxes.
- **DYNAMIC AGENT JOIN & IMMEDIATE RESPAWN (NO IDLE WORKERS):** Worker capacity `A` must remain pinned at 3 active agents (System Concurrency Capacity = 6) throughout execution. As soon as ANY single worker completes, the lead immediately spawns a replacement worker for the next pending subtask to keep concurrency at exactly 3 active subagents (100% duty cycle).
- **SOLO EXECUTION IS AN AUTO-REJECT FAILURE:** The lead orchestrator is **STRICTLY FORBIDDEN** from executing planning, spec writing, or code changes by itself without calling the `invoke_subagent` tool.

---

## Non-Destructive System Evolution & Safety Invariants

Every system-wide improvement proposed or implemented under this prompt MUST respect the **Zero-Harm Principle**:

### 1. Data Loss Prevention & Working Tree Protection
- **No Unstashed Discards:** Batch commands (`pull all`, `sync`) must never destroy uncommitted modifications. Always enforce autostash (`--autostash`) and safe merge aborts (`git merge --abort`) upon conflict detection.
- **Non-Destructive Deletions:** File removal operations must move files to system trash or a staging backup folder rather than unrecoverable unlinking.

### 2. Backward Compatibility & Script Stability
- **Additive Evolution:** New features, parameters, or flags must be additive. Existing CLI arguments, JSON telemetry schemas, and exit codes must remain 100% backward compatible for downstream automation scripts.
- **Graceful Fallbacks:** External service failures (network drops, SSH timeouts, missing API keys) must degrade gracefully with human-readable diagnostic hints rather than abrupt panics or crashes.

### 3. Database Longevity & Zero-Lock Concurrency
- **Additive Schema Migrations:** Never alter, rename, or drop existing SQLite tables or columns destructively. Only add new nullable columns or additive tables.
- **WAL Mode & Busy Timeout Preservation:** Preserve `PRAGMA journal_mode=WAL;` and `PRAGMA busy_timeout=5000;` to prevent database locking across concurrent CLI processes.

---

## The 100 Whole-System Improvements: 10 Subsystem Domains

The AI must partition its discovery across exactly 10 subsystem domains, uncovering 10 distinct, grounded improvements per domain (10 × 10 = 100 items):

### Domain 1: Repository Discovery & Filesystem Walker (Items 01–10)
- Lock-free directory dispatch queues eliminating lock contention on NVMe SSDs.
- Native OS fast-path crawling (`FindFirstFileExW` / `FindExInfoBasic` on Windows).
- Incremental filesystem watcher mtime caching skipping un-modified directory trees.
- User-level and workspace-level `.gitmapignore` file support.
- Non-blocking ring-buffer progress channels preventing slow TUI renderers from throttling the scanner.

### Domain 2: Parallel Batch Cloner & Restoration Engine (Items 11–20)
- Dynamic concurrency auto-throttling when remote hosts return HTTP 429 rate limits.
- Clonenow SQLite progress journal enabling resume-after-network-drop (`--resume`).
- Shallow initial clones (`--depth 1`) with asynchronous background unshallow fetches.
- Git Protocol v2 injection (`GIT_PROTOCOL=version=2`) across all spawned subprocesses.
- Pre-flight volume disk space verification preventing out-of-disk clone failures.

### Domain 3: Multi-Node SSH Fleet Orchestration & Host Management (Items 21–30)
- OpenSSH `ControlMaster` socket connection pooling reducing connection latency.
- Non-blocking concurrent fleet heartbeat updating live reachability badges.
- Automated host key fingerprint validation (`SHA256:...`) protecting against MITM.
- One-command public key deployment (`gitmap ssh copy-id <alias>`).
- Soft-draining node states allowing in-flight tasks to finish cleanly before maintenance.

### Domain 4: Cluster Task Delegation, Cross-Node Remote Execution & RPC (Items 31–40)
- Wiring cluster task dispatchers (`ExecInstall`, `ExecGit`, `ExecProj`, `ExecLifecycle`).
- Streaming remote stdout/stderr chunks over RPC rather than buffering whole outputs.
- Distributed lock manager (DLM) coordinating exclusive multi-node sync operations.
- Work-stealing worker pools balancing batch execution across heterogeneous hosts.
- Node placement scoring based on available CPU, RAM, and network latency.

### Domain 5: Three-Tier Split-DB Architecture & Schema Governance (Items 41–50)
- Automated cross-DB attachment (`ATTACH DATABASE`) for unified multi-DB reporting.
- Hard size ceilings (10MB) on telemetry databases (`repodb/pipeline.db`) with FIFO pruning.
- Split-database integrity checker (`PRAGMA integrity_check`) in `gitmap db check`.
- Automated database snapshots prior to running forward schema migrations.
- PascalCase DDL verification ensuring consistent table and column naming.

### Domain 6: Typed SQL Engine, WAL Mode Concurrency & Query Optimization (Items 51–60)
- Prepared statement LRU cache in `dbengine` avoiding SQL re-parsing overhead.
- Monadic `Result[T]` query wrapping standardizing error handling across all queries.
- Slow query profiler logging queries exceeding 50ms with parameter values.
- Code-generated zero-allocation row scanners eliminating reflection overhead.
- Safe transaction rollbacks wrapping all database transactions with `recover()`.

### Domain 7: CI/CD Pipeline Telemetry, Triage & Self-Healing RCA Engine (Items 61–70)
- Persistent commit SHA cache (`PipelineShaCache`) avoiding redundant GitHub API fetches.
- Bounded concurrent worker pools for downloading failure logs in parallel.
- Single-pass zero-allocation ANSI escape stripping on gigabyte log streams.
- Dynamic ETA sleep loops adapting backoff based on rolling job duration history.
- Targeted smart re-run commands (`gitmap pe rerun --failed-only`) triggering failed jobs only.

### Domain 8: Git Credential Management, DPAPI Resilience & Headless SSH (Items 71–80)
- Automatic detection of broken Windows Credential Manager with auto-switch to DPAPI.
- Headless non-interactive Ed25519 SSH keypair generation with 0600 file permissions.
- Repository-scoped SSH command pinning (`core.sshCommand`) avoiding global config edits.
- Automated public key clipboard integration across Windows, macOS, and Linux.
- Stale `.git/index.lock` detection with PID inspection and safe recovery.

### Domain 9: Autonomous AI Agent Task Engine & Heatmap Telemetry (Items 81–90)
- Atomic task claiming with subtask leases in SQLite preventing multi-agent collisions.
- Real-time file touch categorization (read, edit, create) with line-count metrics.
- Heartbeat timeouts and automatic reclamation for abandoned subtasks.
- Worktree collision matrix warning before an agent touches an in-flight file.
- Self-contained interactive HTML playback export (`--html-replay`).

### Domain 10: Interactive Macro Engine, Terminal Styling & Power/Loss Prevention (Items 91–100)
- Dynamic variable interpolation (`${REPO_DIR}`, `${BRANCH}`) in recorded macro steps.
- Reverse action undo journals (`record_undo.go`) enabling safe macro rollbacks.
- OS sleep prevention daemons inhibiting display sleep during multi-hour batch tasks.
- Argon2id salted key derivation for local encrypted secrets vaults.
- Automatic secret masking redacting tokens and passwords from terminal streams.

---

## Mandatory Reporting Schema for Each Systemic Improvement

Every single improvement produced by this prompt MUST adhere strictly to the following 7-field schema:

```markdown
### Item [NNN]: [Systemic Improvement Title]
- **Target Subsystem & Touchpoint:** [Subsystem path, e.g. `cli/cmdpull/` batch pull or `cli/cmdssh/` fleet engine]
- **User Experience Friction Point / Symptom:** [Detailed real-world scenario describing how developer velocity, reliability, or clarity is currently degraded]
- **Root Cause Analysis (Why it happens):** [Underlying technical cause in code, architecture, or environment]
- **Systemic Effect & Impact on User Workflow:** [Exact consequences on execution time, error rates, mental overhead, or automation friction]
- **Non-Destructive Implementation Blueprint:**
  - *Safe Default Strategy:* [How the feature defaults safely, e.g. opt-in flag, non-destructive trash, autostash guard]
  - *Backward Compatibility Guarantee:* [Proof that existing commands, flags, scripts, and database records remain unbroken]
  - *Technical Enhancement Architecture:* [Concrete Go structs, interfaces, algorithms, or CLI flags to be implemented]
- **Safety Verification & Rollback Protocol:** [Exact procedure to verify functionality and roll back cleanly without data loss]
```

---

## Concrete Example Item

```markdown
### Item 022: Clonenow Resume State Journal & Interruption Recovery
- **Target Subsystem & Touchpoint:** `cli/clonenow/clonenow.go` and `cli/cloner/cache.go`
- **User Experience Friction Point / Symptom:** When a developer runs `gitmap clone-now` to restore 80+ repositories from a CSV/JSON manifest and the network disconnects halfway through, restarting the command re-scans and re-attempts every repository from the beginning, wasting gigabytes of bandwidth and 15+ minutes of time.
- **Root Cause Analysis (Why it happens):** `clonenow` currently maintains batch progress only in volatile memory slices. When the CLI process terminates, all completion state is lost, forcing the next run to re-probe all repositories.
- **Systemic Effect & Impact on User Workflow:** Developers facing unstable network links or battery interruptions experience extreme friction, repeated rate limiting from GitHub/GitLab, and inconsistent working tree states.
- **Non-Destructive Implementation Blueprint:**
  - *Safe Default Strategy:* Store progress in an isolated SQLite journal table `ClonenowProgress` in `installation.db`. If the command completes successfully, the journal is automatically cleaned up.
  - *Backward Compatibility Guarantee:* Standard invocation `gitmap clone-now <file>` continues working identically. A new additive `--resume` flag is introduced; without `--resume`, the command runs fresh as before.
  - *Technical Enhancement Architecture:*
    ```go
    type ClonenowJournal struct {
        RepoName    string
        Status      string // "pending", "success", "failed"
        TargetPath  string
        CompletedAt time.Time
    }
    func (j *ClonenowJournal) RecordSuccess(repoName string) error
    func (j *ClonenowJournal) GetPendingRepos(all []model.ScanRecord) []model.ScanRecord
    ```
- **Safety Verification & Rollback Protocol:** Simulate a SIGINT signal during a 10-repo batch clone; re-run with `--resume` and verify that already cloned repositories are skipped instantly, while remaining repositories resume cleanly without touching existing files.
```

---

## Execution Rules & Guidelines

1. **Strict 100-Item Completeness:** The output MUST contain exactly 100 numbered improvements across the 10 domains.
2. **Grounding in Code & Workflows:** Every improvement must target an existing subsystem in the repository with real technical mechanisms.
3. **Strict Non-Destructive Invariants:** Every item must explicitly define its Safe Default Strategy and Backward Compatibility Guarantee. Proposing destructive changes without autostash or safe fallbacks is an immediate failure.
4. **Zero Full Builds or Heavy Tests:** Do not invoke `go build` or `go test ./...` during discovery. Inspect files using GitMap high-speed commands (`gitmap aum search`, `gitmap lf`, `gitmap cat`).
5. **Atomic Commit Format:** On completion, commit atomically using `gitmap cpf "<module> - <summary>"`.
