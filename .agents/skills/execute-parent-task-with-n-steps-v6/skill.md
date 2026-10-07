---
name: execute-parent-task-with-n-steps-v6
description: >-
  Use this skill when the user asks you to execute a parent task with N steps using the V6 prompt (parameter-driven execution, SQLite task manager, GitMap hyphen-separated atomic commits, secrets gate, and ledger resume).
---

# [V6] Parent Task N-Step Continuous Loop & Mandatory Multi-Agent Subagent Orchestration — Workflow (must follow)

```text
N = 300 (Total self-loop steps budget — editable top-header parameter, default: 300)
A = 2   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
PHASE_1_BUDGET = N / 2   (Steps 1 .. 150: Planning, Parallel Discovery Subagents, Detailed Spec, and Lean Subtask Generation)
PHASE_2_BUDGET = N / 2   (Steps 151 .. 300: Mandatory Parallel Subagent Execution, Self-Looping, Targeted Quality Linting)
WAVES = ceil(subtasks / (A x H))
```

> [!IMPORTANT]
> Prompt Version: 6.0.0
> Runtime: Google Antigravity 2.0 (IDE and CLI)
> Invoke: /execute-parent-task-with-n-steps-v6 <task>
>
> **Top-Instruction Priority Mandate (Above Precedence / Preamble Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given ABOVE this prompt (including in the user preamble, header blocks, or incoming user request above) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below. The agent MUST inspect and follow the instructions above with absolute precedence.

[/goal](slashCommand;goal) Autonomously orchestrate and execute the parent task end-to-end: FIRST showcase and list out the given task in visible chat during Turn 1, capture it verbatim, plan it in the repo, spawn autonomous subagents via `invoke_subagent` (A = 2, H = 2; solo execution without calling `invoke_subagent` is an auto-reject failure) in disjoint file boxes using GitMap high-speed commands as primary, prove every single claim with concrete evidence, enforce coding guidelines to 100%, and finish with one atomic GitMap commit (`gitmap cpf "<module> - <summary>"` / `gitmap cpb "<module> - <summary>"`) using hyphen format (no colons in GitMap arguments, as the colon is already provided by GitMap) that holds strictly this task's files.

[/learn](slashCommand;learn) Enforce the Top-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are provided ABOVE this prompt outrank everything below. Turn 1 MUST showcase the given task list in visible chat before any background execution. Each rule is stated once (R1 to R16) and cited by ID. Progress lives in the ledger and in `.ai-memory/plans/`, never only in chat.

[/plan](slashCommand;plan) Execute thorough step-by-step planning in the repository before execution. Ensure all deliverables, architecture boundaries, and requirements are clearly defined in the audit ledger and subtask plans before dispatching worker waves.

### 🚨 MANDATORY SUBAGENT SPAWNING GATE (A = 2, H = 2 — ZERO SOLO EXECUTION ALLOWED)

- **ACTUAL TOOL CALL REQUIRED:** You must ACTUALLY CALL the `invoke_subagent` tool via your tool-calling API. Do NOT just print the text "Dispatched Worker..." and stop. If you only print text, the agents will not spawn and the task will fail! You must execute the `invoke_subagent` JSON tool payload.
- **3-STAGE MANDATORY DISPATCH:** You must invoke `A = 2` agents (`invoke_subagent`) at EVERY stage of the workflow:
  1. **Planning Step (A = 2 `research` subagents):** Spawn 2 read-only discovery subagents (`TypeName: "research"`) to research the codebase and return findings; lead writes the unified plan.
  2. **Spec Step (A = 2 `self` subagents or lead):** Spawn 2 subagents (`TypeName: "self"`) to author modular, disjoint spec files and subtasks (never writing the same file).
  3. **Execution Step (A = 2 `self` worker subagents):** Spawn 2 worker subagents (`TypeName: "self"`) to execute code modifications in strictly disjoint file boxes.
- **SOLO EXECUTION IS AN AUTO-REJECT FAILURE:** The lead orchestrator is **STRICTLY FORBIDDEN** from executing planning, spec writing, or code changes by itself without calling the `invoke_subagent` tool. Failing to call the actual tool is a critical protocol violation.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/<slug>/skill.md` does not exist in the workspace, create it now.
2. Extract the core instructions of this prompt and save it into that `skill.md` using the standard YAML frontmatter (with `name` and `description`).
3. Once installed, rely on progressive disclosure for future runs. Do not keep the entire prompt in active memory if not needed.

---

## The Unified Master Pipeline (Atomic Numbered Steps)

Execute this task via a strict 3-Phase pipeline. Do not skip steps.

### Phase 1A: Verbatim Capture, Task Extraction & Chat Output Gate (Step 0)

Before executing any file searches, scans, spec writing, or code changes, you must execute Phase 1A:

1. **Top-Instruction Priority Verification:** Whatever directives, constraints, checklists, or instructions are given before this section or prompt (user preamble, header constraints, prior instructions) must be verified as highest priority and non-negotiable.
2. **Showcase Given Task First (Turn 1 Action):** In your VERY FIRST response turn upon receiving the prompt, you MUST output the confirmed task breakdown directly in visible chat. Never execute tools silently without displaying the task breakdown to the user first!
3. **Lossless Verbatim Capture:** Store incoming prompt losslessly under `## User Request (Verbatim)` in canonical spec and parent plan.
4. **Screenshots & Media:** Decode base64/screenshots immediately into `assets/screenshots/<slug>-<NN>.png`. Reference via relative markdown links (`![Screenshot](assets/screenshots/<slug>-<NN>.png)`).
5. **Discrete Deliverables Extraction:** Break down whatever user requirements were given into discrete, actionable items with ordered traceable IDs (`Task-01`, `Task-02`, ...).
6. **Mandatory Same-Turn Tool Chaining (TOTAL BAN ON TURNING OFF):** Emit the breakdown in chat with clean vertical formatting, and in the **EXACT SAME TURN**, invoke your first tool call (e.g. `write_to_file` to initialize ledger/spec, or run preflight). NEVER emit text alone (which ends the turn prematurely), and never ask "Should I proceed?".

```markdown
### 📋 Confirmed Task Breakdown & Requirement Ingestion

1. **Task-01: [Descriptive Task Title]**
   - **State:** `[IN PROGRESS — EXECUTING IMMEDIATELY]`
   - **Understood:** `[YES]` — [1-2 concise sentences proving understanding of intent, scope, and verified constraints]
   - **Actionable Scope:** [Precise technical deliverable and implementation scope]
   - **Target Files / Area:** `[relative/path/or/module]`

2. **Task-02: [Descriptive Task Title]**
   - **State:** `[QUEUED — EXECUTING NOW WITHOUT USER PROMPT]`
   - **Understood:** `[YES]` — [1-2 concise sentences proving understanding of intent, scope, and verified constraints]
   - **Actionable Scope:** [Precise technical deliverable and implementation scope]
   - **Target Files / Area:** `[relative/path/or/module]`

Proceeding directly to Preflight & Phase 1B Spec Generation (Active Tool Call Running Below).
```

---

## 1. Precedence Hierarchy & Scope (Highest First)

1. **User Instructions & Preamble:** Directives and parameters ABOVE this prompt outrank everything below.
2. **Platform Limits:** Native tools, Artifact Review Policy, permission prompts, hooks. Never claim to override them.
3. **Repo Rules:** `agents.md`, `.ai-memory/strictly-avoid.md`, and `coding-guidelines.md`.
4. **This Prompt.**

If sources conflict, follow stricter one and record under `Conflicts:` in ledger.

### Scope Control Rules
- **Turn 1 Task Showcase:** In your very first response turn, you MUST showcase and list out the given task in visible chat. Running tools silently without presenting the task breakdown is strictly banned.
- **Read budget:** Read only requested paths, search hits, and Step 0 context.
- **History read-only:** Never edit past events, changelogs, completed plans, release notes, or `06-old-prompts/` and `19-old-execute-prompts/`.
- **Out-of-scope:** Log under `Follow-ups:` in plan; never fix in this run.
- **Minimal diff:** Change only lines required; never reflow unaffected lines.
- **Indexes:** Update `01-prompts/readme.md` and `.ai-memory/prompts.md` only when adding/modifying prompts. Update `.ai-memory/plans/readme.md` and `02-spec/21-app/readme.md` every run. Register recent completed tasks before push.
- **Mandatory Multi-Agent Partitioning:** Even for tasks touching few files, work MUST be partitioned across A workers (e.g. Worker 01 implements changes, Worker 02 implements verification/linters/companion tests). Solo execution is strictly banned.
- **Finish early:** When all Task-IDs are `DONE`, proceed directly to consolidation.
- **Zero releases:** Never bump versions or edit changelogs unless requested (R10).

---

## 2. Core Operational Rules (Cite by ID)

- **R1 Zero Builds or Test Suites (TOTAL BAN).** NEVER run `go build`, `npm run build`, `vite build`, `go test ./...`, `pytest`, `npm test`, or `03-ai-scripts/06-cicd-local-runner.py`. CI verifies builds and suites. Routine turns must never waste time on heavy compilation/tests. Only explicit user command lifts this.
- **R2 Targeted Checks Only.** Run only fast, file-scoped checks on specifically modified files (see Section 10). A check scanning 0 files is a **FAIL**.
- **R3 Evidence or It Did Not Happen.** Every `DONE`, `PASS`, or "verified" claim MUST cite a concrete file path, git diffstat, or command exit code (`exit 0`). Vague assurances are auto-rejected.
- **R4 Never Invent Commands, Flags, or Paths.** Verify commands with a harmless call (`gitmap lf readme.md`), not `--help`. Use documented fallbacks and log in ledger.
- **R5 Mandatory Subagents (`invoke_subagent`).** Spawning subagents via `invoke_subagent` (`A = 2`, `H = 2`) is an **ABSOLUTE MUST** (`research` for discovery in Phase 1, `self` for edits in Phase 2). The lead agent is STRICTLY FORBIDDEN from executing all reads or edits solo. Solo execution without calling `invoke_subagent` is an auto-reject failure on the same tier as Rule 0.
- **R6 One Owner Per File (Disjoint Bounding Boxes).** Within every worker wave, each file has exactly one owner. Shared indexes (`.ai-memory/plans/readme.md`, `.ai-memory/prompts.md`, `.ai-memory/what-to-read.md`, `02-spec/21-app/readme.md`, directory `readme.md`) belong exclusively to lead.
- **R7 Git Safety & Isolation (Worker Git Ban).** Subagents NEVER run git commands (`git add`, `git commit`, `git push`, `git status`, `git diff`, `git checkout`). In shared workspaces (`Workspace: "inherit"`), worker git calls create `.git/index.lock` collisions that immediately crash parallel agents. Nobody runs `git reset --hard`, `git checkout --`, `git clean`, `git stash`, or force pushes.
- **R8/R9 Atomic Commit & Push via GitMap (TOTAL BAN ON RAW GIT COMMITS).** The run ends with ONE GitMap call:
  - For features: `gitmap cpf "<module> - <feature summary>"` (e.g. `gitmap cpf "CBF - implement user profile dashboard"` or `gitmap cpf "aum-agent-db - implement sqlite task tracking engine"`). GitMap automatically prepends `Feature: ` (which already provides the colon), so mention and format as a hyphen `-`; there is NO need to provide a colon in the GitMap `cpf` (or `cpb`/`cpr`) message argument. DO NOT provide a colon; use a hyphen `-` to separate module/scope from summary. NEVER include `feat(...)` or `feature:` in your message.
  - For bug fixes: `gitmap cpb "<module> - <fix summary>"` (e.g. `gitmap cpb "AUM - validate regex without nil fallback"` or `gitmap cpb "aum-validate-regex - prevent nil fallback on malformed patterns"`). GitMap automatically prepends `Bug: ` (which already provides the colon), so mention and format as a hyphen `-`; there is NO need to provide a colon in the GitMap message argument. DO NOT provide a colon; use a hyphen `-` to separate module/scope from summary. NEVER include `fix(...)` or `bug:` in your message.
  GitMap stages, formats, commits, and pushes atomically. TOTAL BAN on raw git commits (`git commit`, `git commit -m "..."`, `git add -A`, raw `git push`), colons inside the GitMap message argument, and conventional prefixes (`docs(...)`, `feat(...)`, `fix(...)`, `chore(...)`). ZERO intermediate commits: never commit during Phase 1 (plans/specs) or Phase 2; all files across the turn MUST be committed together at the final step of Phase 3. Before GitMap, all push gates must pass (targeted checks, secrets gate, and `.gitignore` hygiene; untrack any ignored files: `git rm --cached`). Push rejected: `git pull --rebase`, re-run command. Miss after push: allow one follow-up `gitmap cpb "<module> - <fix summary>"`, logged as `FOLLOW_UP_PUSH: <sha>`. Never amend pushed commits. Workers never run git commands or GitMap commit tools; only lead does.
- **R10 Zero Unauthorized Releases.** Never bump versions, edit `version.json`, update changelogs, or trigger release scripts unless user explicitly requested release.
- **R11 Strict Relative Git Paths & Lowercase Hygiene.** Strict ban on absolute paths (`C:\...`, `/home/...`) and `file:///` URIs. Only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well. Paths relative from git root. All filenames, documentation, and specs strictly lowercase (e.g. `readme.md`, `agents.md`, `skill.md`).
- **R12 No Polling / Immediate Turn Yielding.** When calling `invoke_subagent`, make it the sole tool action at turn end, print progress line (`Dispatched Worker 01 .. Worker <A> (wave k / WAVES); waiting for their results.`) and **STOP CALLING TOOLS**. Never poll in loop. Check `manage_subagents` once if wave runs long.
- **R13 Two-Strike Retry Cap & Anti-Looping.** Tool failing twice: worker replies `STATUS: BLOCKED` with exact error and stops. Lead takes over and logs `LEAD_FALLBACK: <reason>`. Subtask failing two remediation rounds is marked `FAILED` with RCA (Section 12).
- **R14 100% Ambiguity & Decision Boundaries.** Non-blocking: choose conservative option, log in ledger `Assumptions:`, proceed. Blocking: `ask_question` once, log in `.ai-memory/ambiguous-questions/01-new-ambiguity/`, continue unblocked tasks.
- **R15 Zero Generated Artifacts Committed.** Never commit build caches, logs, temp scripts, or newly generated code (Hard Rule 1) unless repository already tracked them.
- **R16 Zero Secrets in Standard Repos.** Never write credentials, tokens, passwords, or `.env` contents into tracked files, commits, ledger, plans, specs, or prompts (`agents.md` section 9).
  - *Secrets Gate (lead, before GitMap call):*
    1. Check changed/new files via `git status --porcelain`.
    2. Run `python linter-scripts/check-forbidden-strings.py`.
    3. Search files via `gitmap aum search -r "(BEGIN [A-Z ]*PRIVATE KEY|AKIA[0-9A-Z]{16}|gh[pousr]_[A-Za-z0-9]{36}|sk-[A-Za-z0-9]{20,}|xox[baprs]-[A-Za-z0-9-]{10,})"` for private keys, AWS, GitHub, OpenAI, Slack tokens, or secret/token assignments (NEVER run `git grep` or `Select-String`).
    4. On hit: if `repo-secrets` exists in default work directory, store via `gitmap rs text "<value>" --slug <slug>` (or `gitmap rs file <path>`) and replace with env var/placeholder; if not, remove value and `ask_question` once. Log `SECRET_OFFLOADED: <file>:<line>` without value.
    5. Never print secrets in chat/logs; refer to file:line only. Workers finding a secret report `BLOCKED: secret at <file>:<line>`. Never put repository URLs or absolute paths into secrets instructions (`agents.md` section 9).

---

## 3. GitMap High-Speed Command Primacy (Run Everything Faster)

GitMap is your **PRIMARY** acceleration engine. NEVER use `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, `Get-ChildItem -Recurse`, or `findstr`. Execute all searches and operations through GitMap:

| Operation | Primary GitMap Command | High-Speed Alias | Purpose & Advantage |
| :--- | :--- | :--- | :--- |
| **Live Streaming Search** | `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r] [-i]` | `gitmap aum grep "<pat>"` | Multi-core streaming live file search (replaces `Select-String`, `git grep`) |
| **Indexed Symbol Search** | `gitmap search "<query>" [--limit <n>]` | `gitmap search` | Instant SQLite cached keyword/symbol search across indexed repos |
| **Wildcard File Search** | `gitmap find "<pattern>" [-ext <ext>]` | `gitmap f "<pat>"` | Index-accelerated multi-core glob filename finder |
| **Directory Inventory** | `gitmap list-files [pattern] [-ext <ext>]` | `gitmap lf [pat]` | Instant indexed repository inventory |
| **Substring File Search** | `gitmap find-files-any "<substring>"` | `gitmap ffa "<str>"` | High-speed partial filename matcher |
| **Stream File** | `gitmap cat <filepath>` | `gitmap cat` | Zero-disk memory streaming to stdout |
| **PowerShell Runner** | `gitmap pwsh "<command>"` | `gitmap ps "<cmd>"` | High-speed PowerShell execution with `-NoProfile` |
| **Bash Runner** | `gitmap bash "<command>"` | `gitmap sh "<cmd>"` | Standard cross-platform Bash command execution |
| **Offload Secrets** | `gitmap rs file <filepath>` / `folder` / `text` | `gitmap rs` | Auto-commits into `repo-secrets` in work directory |
| **Offload Scripts** | `gitmap rc file <file.ps1>` / `text` | `gitmap rc` | Auto-commits reusable scripts into `repo-cache` |
| **Atomic Commits** | `gitmap cpf "<module> - <msg>"` (Feature) / `cpb` (Bug) | `gitmap cpf` | Stages, formats with prefix, and pushes atomically. Mention as a hyphen `-` (no need to provide a colon in GitMap `cpf`/`cpb`/CVF commit arguments because the colon is already automatically provided by GitMap in `Feature: ` or `Bug: `). |
| **Pipeline Waiting** | `gitmap pipeline-ai status --json` | `gitmap pl-ai` | Non-polling dynamic ETA CI/CD monitor |

### 🔍 Code & Symbol Search Protocol (TOTAL BAN ON `rg`, `ripgrep`, `Select-String` & `git grep`)
- **Live Disk Search (Default for discovery, symbol tracking & blast radius):**
  - Search string/symbol: `gitmap aum search "<symbol>" [dir] [-e <.ext>]` (e.g. `gitmap aum search "RunFleetPASCommand" cli -e .go`)
  - Search regex: `gitmap aum search -r "<regex>" [dir] [-e <.ext>]` (e.g. `gitmap aum search -r "(\"pas\"|\"pa\")" cli/cmd -e .go`)
  - Case-insensitive: `gitmap aum search -i "<query>" [dir]`
- **TOTAL BAN:** NEVER run `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, `Get-ChildItem -Recurse`, or `findstr`. Running raw grep, ripgrep, or slow shell search pipelines wastes execution steps, causes Windows process hangs, and violates GitMap primacy. Search exclusively through `gitmap aum search` or `gitmap find`.

---

## 5. Step 0: Preflight, Platform Handshake & SQLite Task DB Initialization (Phase 1 Budget)

1. **Platform Handshake:** Confirm tools (`invoke_subagent`, `send_message`, `manage_subagents`, `ask_question`, `write_to_file`, `replace_file_content`, `run_command`). If `task_boundary` exists: set `PLANNING` (Phase 1), `EXECUTION` (Phase 2), `VERIFICATION` (Phase 3).
2. **Commands & Directory:** Confirm `gitmap --version` and `python --version` exit 0. Verify GitMap with harmless call (`gitmap lf readme.md`), not `--help`. `run_command` uses `Cwd` in workspace root, paths relative. Never cd to other drives or tool folders.
3. **Working Tree Cleanliness:** Run `git status --porcelain`. Record modified files in ledger; never touch them. Confirm root `readme.md` is lowercase. Read `.ai-memory/what-to-read.md`, `strictly-avoid.md`, `coding-guidelines.md`.
4. **SQLite Task DB & Deterministic Slug Initialization (Check Before Creating):**
   - Initialize or inspect task state via the Antigravity SQLite task manager:
     `python 03-ai-scripts/46-agent-sqlite-task-manager.py init --name "<task name>" --budget 300`
   - Inspect Task DB Schema & Data Integrity Rules (Mandatory Preflight):
     `python 03-ai-scripts/46-agent-sqlite-task-manager.py schema`
     Verify tables (`ParentTask`, `Subtask`, `AgentActionLog`), positive boolean flags (`IsActive`, `HasCompleted`, `IsBlocked`), and payload constraints before populating subtasks. (Machine-readable schema: `python 03-ai-scripts/46-agent-sqlite-task-manager.py schema --json`; DDL: `python 03-ai-scripts/46-agent-sqlite-task-manager.py schema --ddl`).
   - Output Analysis & Crash Forensics:
     - If `action: "RESUME_FOUND"`: A matching or similar slug exists in `.ai-memory/temp-agents/`! If `diagnostics.hasCrashesDetected: true`, inspect the forensic report (`diagnostics.crashedAgents`) to identify which agent crashed, what file it was touching, and the last logged action. Resume execution from the uncompleted subtask.
     - If `action: "INITIALIZED"`: Created dedicated run directory `.ai-memory/temp-agents/<nn>-<slug>/` and SQLite database `agent-task.db` with WAL mode (`PRAGMA journal_mode=WAL; PRAGMA busy_timeout=5000;`).
5. **Ledger Creation:** Also create `.ai-memory/temp-agents/<nn>-<slug>/ledger.md` mirroring the DB initialization for human readability:

```markdown
# Ledger: nn-<slug>
Request slug: <slug>
Request first line: <verbatim first line>
Status: ACTIVE
Phase: 1    Wave: 0 / WAVES    Step: 1 / N
Last completed action: Phase 1A Capture & Task Breakdown
Next action: Phase 1B Spec & Plan
Workers in flight: none
Commits: none    Pushed: no
Branch: <branch> | Tree at start: clean (or dirty with <paths>)
Tools: invoke_subagent=yes send_message=yes ask_question=yes gitmap=yes sqlite_db=<databasePath>
| Task-ID | Subtask | Owner | Owned files | Status | Evidence |
|---|---|---|---|---|---|
| Task-01 | 01-<name> | Worker 01 | <paths> | PENDING | - |
Assumptions: <list or none>
Conflicts: <list or none>
Stage list: <every path this run creates or modifies>
```

---

## 6. Phase 1: Planning Step & Spec Step (Steps 1 .. PHASE_1_BUDGET)

You must use `invoke_subagent` to delegate both planning discovery and spec authoring.

> [!CAUTION]
> **SUBAGENT CRASH PREVENTION (MANDATORY ARCHITECTURAL RULES):**
> 1. **Tool Capabilities (`research` vs `self`):** Subagents with `TypeName: "research"` have READ-ONLY tools (`view_file`, `run_command`, web search). They DO NOT have `write_to_file` or `replace_file_content`. Commanding a `research` subagent to write files will cause an immediate tool failure and crash. Subagents that modify or author files MUST use `TypeName: "self"`.
> 2. **File Write Collisions & Windows Locks:** Two subagents MUST NEVER be assigned to create or write the same file path simultaneously in `Workspace: "inherit"`. Windows file locks and race conditions will crash both agents. In Phase 1 Planning, discovery subagents return findings via structured messages to the lead, and the lead alone writes the unified plan. In Phase 1 Spec, subagents author strictly disjoint modular files.
> 3. **Valid TypeName Strings:** `TypeName` MUST strictly be `"research"` (read-only discovery) or `"self"` (read-write execution). Passing custom strings (e.g. `"TypeName": "Research 01"` or `"TypeName": "Worker 01"`) is an invalid schema and immediately crashes `invoke_subagent`. Put descriptive labels in `Role` only.
> 4. **Worker Git Ban (Index Lock Prevention):** Subagents NEVER run git commands (`git add`, `git commit`, `git push`, `git status`, `git diff`, `git checkout`). In shared workspaces, worker git calls create `.git/index.lock` collisions that immediately crash parallel agents with exit code 128.
> 5. **Clean Turn-Yielding:** When calling `invoke_subagent`, make it the final tool action of the turn. Emit the status line and IMMEDIATELY STOP CALLING TOOLS to allow platform reactive wakeup. Chaining additional tools or tight polling crashes the message queue.

1. **Planning Step (A = 2 `research` Discovery Subagents):** Lead agent calls `invoke_subagent` to spawn 2 read-only discovery subagents (`TypeName: "research"`, `Role: "Research 01: Architecture & Blast Radius"`, `Role: "Research 02: Specs & Dependency Mapping"`). Their prompt MUST instruct them to research the codebase using GitMap high-speed search (`gitmap aum search "<symbol>" [dir] [-e <.ext>]`, `gitmap find`, `gitmap lf`, `gitmap cat`) — NEVER `rg`, `ripgrep`, PowerShell `Select-String`, or `git grep` — define symbol boundaries and caller dependencies, and return their structured findings to the lead in their final message.
   - *Master Plan Generation (Lead Agent):* The lead orchestrator receives both discovery reports, synthesizes findings, and writes the unified Execution Plan (`.ai-memory/plans/pending/nn-<slug>.md`) and the Root Task JSON Manifest.
   - *Tool Call:* Lead must execute the `invoke_subagent` tool as the final action in the turn, print `Dispatched Planning Agents`, and then STOP CALLING TOOLS to wait for `<SYSTEM_MESSAGE>` reactive wakeup.
2. **Spec Step (A = 2 `self` Authoring Subagents):** Once planning is synthesized, the lead agent calls `invoke_subagent` to spawn 2 authoring subagents (`TypeName: "self"`). Their prompt MUST instruct them to author modular, strictly disjoint spec files and subtask plans:
   - Subagent 1 writes `02-spec/21-app/nn-<slug>/01-architecture-spec.md` and subtasks `.ai-memory/plans/subtasks/nn-<slug>/01-<name>.md`.
   - Subagent 2 writes `02-spec/21-app/nn-<slug>/02-component-spec.md` and subtasks `.ai-memory/plans/subtasks/nn-<slug>/02-<name>.md`.
   - NEVER have both subagents write to the same file path!
   - *Tool Call:* Lead must execute the `invoke_subagent` tool as the final action in the turn, print `Dispatched Spec Agents`, and then STOP CALLING TOOLS to wait for `<SYSTEM_MESSAGE>` reactive wakeup.
3. **Populate Subtasks in SQLite Task DB:** Once subtasks are decomposed, inspect schema constraints (`python 03-ai-scripts/46-agent-sqlite-task-manager.py schema`) and populate them into the SQLite database for atomic worker claiming:
   `python 03-ai-scripts/46-agent-sqlite-task-manager.py add-subtasks --db <databasePath> --tasks-json '[{"code": "Task-01", "title": "<title>", "owned_files": ["<paths>"], "agent_role": "Worker 01"}]'`
4. **Readiness Gate:** Complete Phase 1 planning and spec authoring within `PHASE_1_BUDGET` steps, then proceed **UNCONDITIONALLY** into Phase 2. ZERO intermediate git commits during Phase 1!

---

## 7. Phase 2: Execution Step (Worker Waves) (Steps (PHASE_1_BUDGET + 1) .. N)

> [!CRITICAL]
> **MANDATORY `invoke_subagent` DISPATCH (ZERO SOLO EXECUTION):**
> You MUST ACTUALLY CALL the `invoke_subagent` tool to spawn A workers (`TypeName: "self"`, up to H subtasks per worker) in parallel. Executing all subtasks solo in main agent without the tool call is an immediate auto-reject failure.

### 7.1 Dispatch Payload Examples (`invoke_subagent`)

#### Phase 1 Planning Payload (`TypeName: "research"` — Read-Only Discovery):
```json
{
  "Subagents": [
    {
      "TypeName": "research",
      "Role": "Research 01: Architecture & Blast Radius Discovery",
      "Model": "inherit",
      "Workspace": "inherit",
      "Prompt": "You are Research 01 for task <nn>-<slug>. You have no prior chat context; this brief is your complete specification.\n\n### Core Objective:\nExplore the codebase, map symbol callers, trace dependencies, and discover relevant source files for: <task title and objectives>.\n\n### Tool Capabilities & Strict Read-Only Boundary:\n- You are a READ-ONLY subagent (`TypeName: 'research'`).\n- You have read tools: `run_command`, `view_file`, `search_web`, `read_url_content`.\n- You DO NOT have `write_to_file` or `replace_file_content`. NEVER attempt to create or edit files.\n\n### Mandatory Search Primacy & Total Ban on Raw Grep:\n- Execute ALL searches and symbol discoveries exclusively via GitMap:\n  * Live symbol search: `gitmap aum search \"<symbol>\" [dir] [-e <.ext>]`\n  * Path-scoped search: `gitmap aum search \"<symbol>\" --path <relative-dir>`\n  * Indexed keyword search: `gitmap search \"<query>\"`\n  * File discovery: `gitmap find \"<pattern>\"`\n  * Directory inventory: `gitmap lf [dir]`\n  * View files: `gitmap cat <path>` or native `view_file`\n- TOTAL BAN (AUTO-REJECT FAILURE): NEVER run `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, `Get-ChildItem -Recurse`, or `findstr`.\n- TOTAL BAN ON GIT COMMANDS: NEVER run `git` commands (`git add`, `git commit`, `git status`, `git diff`, etc.).\n\n### Required Findings Report (Send via send_message to Parent):\nWhen your research is complete, send a message to the caller containing:\n1. Key symbol definitions, structures, and entry points.\n2. Call sites and blast radius (files that will be affected by modifications).\n3. Recommended modular file ownership boundaries for worker subtasks.\n4. Verification commands to validate changes."
    },
    {
      "TypeName": "research",
      "Role": "Research 02: Specs & Dependency Mapping",
      "Model": "inherit",
      "Workspace": "inherit",
      "Prompt": "You are Research 02 for task <nn>-<slug>. You have no prior chat context; this brief is your complete specification.\n\n### Core Objective:\nInspect existing specifications under `02-spec/`, coding guidelines under `02-spec/02-coding-guidelines/`, and prior plan tasks for: <task title and objectives>.\n\n### Tool Capabilities & Strict Read-Only Boundary:\n- You are a READ-ONLY subagent (`TypeName: 'research'`).\n- You DO NOT have file authoring tools. NEVER attempt to create or edit files.\n\n### Mandatory Search Primacy & Total Ban on Raw Grep:\n- Execute ALL searches exclusively via GitMap:\n  * Live search in specs: `gitmap aum search \"<term>\" 02-spec`\n  * Indexed search: `gitmap search \"<term>\"`\n  * File discovery: `gitmap find \"*.md\"`\n  * View files: `gitmap cat <file>` or native `view_file`\n- TOTAL BAN: NEVER run `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, `Get-ChildItem`, or `findstr`.\n- TOTAL BAN ON GIT COMMANDS: NEVER run `git` commands.\n\n### Required Findings Report (Send via send_message to Parent):\nSend a message to the caller detailing existing specs, architectural invariants, acceptance criteria, and positive boolean rules."
    }
  ]
}
```

#### Phase 2 Execution Payload (`TypeName: "self"` — Read/Write Workers):
```json
{
  "Subagents": [
    {
      "TypeName": "self",
      "Role": "Worker 01: [Assigned Subtask]",
      "Model": "inherit",
      "Workspace": "inherit",
      "Prompt": "<Worker Brief Below>"
    },
    {
      "TypeName": "self",
      "Role": "Worker 02: [Assigned Subtask]",
      "Model": "inherit",
      "Workspace": "inherit",
      "Prompt": "<Worker Brief Below>"
    }
  ]
}
```

### 7.1.B Self-Contained Research Discovery Brief (TypeName: 'research')

Read-only discovery subagents spawn with zero prior chat context. The prompt envelope MUST inject complete instructions, strictly enforce GitMap search primacy, and impose an absolute ban on raw search tools:

```text
You are Research <NN> for task nn-<slug>. You have no prior chat context; this brief is your complete specification.

### Boundaries & Crash Prevention:
- Read-Only Mode: You possess read-only tools (`view_file`, `run_command`, web search). You DO NOT have `write_to_file` or `replace_file_content`. NEVER attempt to create or modify files. Commanding a research agent to write files crashes the subagent.
- TOTAL BAN ON RAW SEARCH TOOLS: NEVER run `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, `Get-ChildItem -Recurse`, or `findstr`. Spawning external search binaries hangs execution, causes OS lock collisions, and wastes steps.
- GitMap Search Primacy: Execute all code, symbol, and pattern searches exclusively via GitMap commands:
  - Live streaming code/regex search: `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r] [-i]`
  - Fast file finder: `gitmap find "<pattern>" [-ext <ext>]`
  - High-speed directory inventory: `gitmap lf [pat]`
  - Stream file content: `gitmap cat <filepath>`
- TOTAL BAN ON GIT COMMANDS: NEVER run `git add`, `git commit`, `git push`, `git status`, `git diff`, `git checkout`.
- Strict Relative Git Paths: All paths cited in your report must be relative to the repository root; only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well (TOTAL BAN on absolute filesystem paths and `file:///` URIs).

### Assigned Research Scope:
- Research Target: <symbol / architectural area / module>
- Primary Objective: Map callers, define type contracts, identify blast radius, and discover existing specifications.

### Output Contract:
Return your structured findings via send_message to the Lead Orchestrator with this JSON envelope, then stop:
{
  "agent": "Research <NN>",
  "target": "<symbol or module>",
  "filesExamined": ["<path1>", "<path2>"],
  "callers": ["<caller1>", "<caller2>"],
  "architecturalBoundaries": "<findings>",
  "existingSpecs": ["<path>"],
  "risksAndBlastRadius": "<risks>"
}
```

### 7.2 Self-Contained Worker Brief (Eliminate Context Blindness)

Subagents spawn with clean context. The prompt envelope MUST inject complete instructions:

```text
You are Worker <NN> for task nn-<slug>. You have no prior chat context; this brief is your complete specification.

### Boundaries & Crash Prevention:
- Read any file in the workspace; edit ONLY your Owned Files: <relative paths>.
- TOTAL BAN ON GIT COMMANDS (LOCK COLLISION PREVENTION): NEVER run ANY git commands (`git add`, `git commit`, `git push`, `git status`, `git diff`, `git checkout`). In shared workspaces, worker git calls create `.git/index.lock` collisions that immediately crash parallel agents. Only the lead orchestrator runs git commands after workers complete.
- TOTAL BAN ON COMMITS: Workers NEVER commit, stage, or push. Committing is exclusively reserved for the Lead Agent at Phase 3 via GitMap (`gitmap cpf "<module> - <summary>"` using hyphen `-`; no colon needed in GitMap cpf as colon is already provided).
- Code & Symbol Search: Use GitMap exclusively: `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r] [-i]` or `gitmap find`. TOTAL BAN on `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, `Get-ChildItem -Recurse`, or `findstr`.
- After C tool calls, stop and report what you have.
- A tool failing twice: reply "STATUS: BLOCKED" with exact error and stop. Never guess paths and never troubleshoot machine.
- Workers that find a secret stop and report "BLOCKED: secret at <file>:<line>". They do not handle it themselves.
- Adhere to R1, R2, and R11 by ID.

### Assigned Subtasks (up to H subtasks):
- Subtask 1: .ai-memory/plans/subtasks/nn-<slug>/01-<name>.md
- Subtask 2: .ai-memory/plans/subtasks/nn-<slug>/02-<name>.md (if assigned)

### 100% Non-Negotiable Coding Guidelines (AUTO-REJECT ON VIOLATION):
1. Positive booleans ONLY: use `is` and `has` prefixes exclusively. NEVER evaluate explicit `== true`. NEVER combine positive and negative checks in the same condition (`if isA && !isB` is BANNED).
2. Go Structured Errors: return `*appfault.AppError`, never bare `error`.
3. Function Sizing: <= 8 lines preferred, hard cap 15 lines. Extract domain structs and raw generics to `types.go`.
4. Strict Relative Git Paths & Lowercase: zero absolute filesystem paths and zero `file:///` URIs. Only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well. All new files strictly lowercase.
5. Repo Secrets: if any credentials or private tokens are needed, store them in the `repo-secrets` folder in the default work directory (via `gitmap rs`). Never commit secrets.
6. Zero Builds or Tests: NEVER run `go build`, `npm run build`, `go test`, or `pytest`.
7. Targeted Verification: Run only fast file-scoped linters (e.g. `python 03-ai-scripts/05-guideline-autofixer.py <folder> --check-only`). A check scanning 0 files is a FAIL.
8. GitMap Search Primacy (TOTAL BAN on rg / ripgrep / Select-String / git grep): NEVER execute `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, `Get-ChildItem -Recurse`, or `findstr`. Always use `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r]` for live symbol/regex discovery.

### Concurrency-Safe SQLite Action Logging (CRASH FORENSICS MANDATE):
- Worker subtasks are tracked in the run database: `<databasePath>`.
- Inspect Database Schema & Constraints:
  `python 03-ai-scripts/46-agent-sqlite-task-manager.py schema`
  Inspect tables, positive boolean fields, and JSON payload specifications before claiming or logging actions.
- Claim assigned subtask atomically:
  `python 03-ai-scripts/46-agent-sqlite-task-manager.py claim --db <databasePath> --agent "Worker <NN>"`
- BEFORE touching or modifying any owned file, you MUST log your in-flight action:
  `python 03-ai-scripts/46-agent-sqlite-task-manager.py log-action --db <databasePath> --subtask-id <id> --agent "Worker <NN>" --action "write_to_file" --file "<path>" --details "<action description>"`
  *(Note: This guarantees that if a tool execution crashes or the session is interrupted, the database permanently records the exact file you were touching and what caused the crash!)*
- When your subtask passes targeted checks, mark completion in the database:
  `python 03-ai-scripts/46-agent-sqlite-task-manager.py complete --db <databasePath> --subtask-id <id> --agent "Worker <NN>" --evidence "PASS exit 0, <files>"`
- If blocked or failing, record the failure:
  `python 03-ai-scripts/46-agent-sqlite-task-manager.py fail --db <databasePath> --subtask-id <id> --agent "Worker <NN>" --reason "<reason>"`

### Output Contract:
Write your subtask output to .ai-memory/plans/subtasks/nn-<slug>/01-<name>.json and reply with this JSON block, once per subtask, then stop:
{
  "task": "Task-01",
  "status": "DONE",
  "filesChanged": ["<path1>", "<path2>"],
  "checks": "<command> -> exit <code>, <files scanned>",
  "acceptance": { "ac1": "PASS <evidence>" },
  "assumptions": [],
  "blockers": []
}
```

### 7.3 Turn-Yielding, Verification & Crash Forensics Protocol

1. **Invoke & Yield:** You must ACTUALLY CALL the `invoke_subagent` tool as the final action in your turn. Print the progress line (`Dispatched Worker 01 .. Worker <A> (wave k / WAVES); waiting for their results.`) and **STOP CALLING TOOLS** to end your turn.
2. **Automated Crash Forensics & Status Inspection:**
   - If any worker fails to report, crashes, or times out, lead immediately runs:
     `python 03-ai-scripts/46-agent-sqlite-task-manager.py diagnose --db <databasePath>`
     This pinpoints the autopsy: which agent crashed, on which subtask, targeting which file, and the exact action that was executing when it failed.
   - Lead inspects overall completion status at any time:
     `python 03-ai-scripts/46-agent-sqlite-task-manager.py status --db <databasePath>`
3. **Verify Worker Reports Independently:** Confirm `git diff --stat -- <owned files>` matches `filesChanged`, no files outside owned files modified, re-run targeted checks for `exit 0` on non-zero files.
4. **Reject Violations:** Send failures via `send_message`. On `BLOCKED`, lead does work and logs `LEAD_FALLBACK: <reason>`. After two failed rounds, mark `FAILED`, write RCA, continue (R13).
5. **Update Ledger:** Record status, evidence, changed paths in `ledger.md` via `replace_file_content`.
6. **Loop:** Dispatch subsequent waves via `invoke_subagent` until all subtasks are `DONE` or `FAILED`.

---

## 8. Phase 3: Consolidation, Evidence Verification & Atomic GitMap Push

1. **Consolidate Subtasks:** Merge completed subtasks into `.ai-memory/plans/completed/nn-<slug>.md`, logging real steps from ledger; link to canonical spec. Delete `.ai-memory/plans/subtasks/nn-<slug>/` and pending plan. Canonical spec in `02-spec/21-app/` stays permanently.
2. **Update Registers:** Update `.ai-memory/plans/readme.md` and `02-spec/21-app/readme.md`. Update `01-prompts/readme.md` and `.ai-memory/prompts.md` only when adding/modifying prompts. Register recent completed tasks before push gate.
3. **Push Gate Verification:** Before GitMap call, verify: (a) targeted checks exit 0 (>0 files), (b) secrets gate clean, (c) `.gitignore` covers caches, build outputs, logs, reports, `.env*` (untrack any tracked ignored files via `git rm --cached <file>` or `git rm -r --cached <dir>`). On failure: abort GitMap call; mark task `FAILED` with RCA.
4. **Atomic Commit & Push:** Call `gitmap cpf "<module> - <feature summary>"` (features, e.g. `gitmap cpf "CBF - implement user profile dashboard"`) or `gitmap cpb "<module> - <fix summary>"` (fixes, e.g. `gitmap cpb "AUM - validate regex without nil fallback"`). Mention as a hyphen `-` because the colon is already going to be provided by GitMap (which automatically prepends `Feature: ` or `Bug: `). So, there is no need to provide the colon in GitMap `cpf`/`cpb`/CVF message arguments (TOTAL BAN on colons `:` inside the message argument); always use a hyphen `-` to separate module from summary (e.g. `<module> - <summary>`). NEVER include redundant conventional prefixes (`feat(...)`, `fix(...)`, `docs(...)`) in the commit string. Push rejected: `git pull --rebase` and re-run. Miss after push: allow one follow-up `gitmap cpb "<module> - <fix summary>"`, logged as `FOLLOW_UP_PUSH: <sha>`. Never amend pushed commits.

---

## 9. Final Report Format (Strict Vertical Lines)

```markdown
### Task Completion Summary

- ✅ **Task-01: [Descriptive Task Title]** — `[Completed]` — [diff/check evidence]
- ❌ **Task-02: [Descriptive Task Title]** — `[Failed]` — [RCA link]

### Modified Files Summary

- [relative/path/to/modified/file1.ext]
- [relative/path/to/modified/file2.ext]

### Steps Used

- Step x / N (Phase 1: y / PHASE_1_BUDGET, Phase 2: z / PHASE_2_BUDGET), Wave k / WAVES

### Implementation Confidence Score

- Confidence: [passed checks / total checks]
- Rationale: [Verified evidence across all criteria, passing targeted linters, zero regressions]

Independent check: run /verify-parent-task-run <slug>

### 🤖 Independent AI Verification & Audit Prompt

(Emit self-contained audit prompt linking spec, plan, and modified files)
```

---

## 10. Targeted Verification Checks (R2)

Confirm scripts exist via harmless workspace call before invoking (R4). Run on changed files/folders only:

- **Coding Guidelines & Boolean Linter:** `python 03-ai-scripts/05-guideline-autofixer.py <folder> --check-only --ext <.ext>`
- **Relative Path Linter:** `python linter-scripts/check-relative-paths.py`
- **Prompts & Spec Index Linter:** `python linter-scripts/check-prompts-loaded.py`
- **Markdown Link & Doc Path Linter:** `python 03-ai-scripts/22-doc-path-linter.py <folder>`
- **Sequence Integrity Linter:** `python linter-scripts/check-sequence-integrity.py`
- **Forbidden Strings Check:** `python linter-scripts/check-forbidden-strings.py`

---

## 11. AI Fix Scripts Memory (Reusable Tooling)

- [ ] [/goal](slashCommand;goal) Reuse First: Scanned and learned `03-ai-scripts/readme.md` before writing temporary code.
- [ ] Strict In-Repository Execution: All Python scripts executed strictly within the codebase repository root.
- [ ] Strict .ai-memory/ Folder Storage: All helper scripts, local runners, and linters stored in `03-ai-scripts/`.
- [ ] Native File Manipulator: Use `python 03-ai-scripts/03-file-manipulator.py <command>` for mass file operations.
- [ ] Go Generate Sync: If Go constants or enums are modified, run `go generate ./...` in the relevant package and commit generated files.

---

## 12. Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] TOP-INSTRUCTION PRIORITY MANDATE: Whatever directives, constraints, checklists, or instructions are given before this section or prompt (user preamble, header constraints, prior instructions) are verified as highest priority and non-negotiable, overriding all lower-level guidelines below.
- [ ] MANDATORY SUBAGENT SPAWNING GATE (ZERO SOLO EXECUTION): Verified that `invoke_subagent` (`A = 2, H = 2`) was explicitly called in Phase 1 (parallel discovery/spec modules) and Phase 2 (`TypeName: "self"` parallel subtask execution). Executing all reads or code modifications solo without calling `invoke_subagent` is an immediate auto-reject failure.
- [ ] NO SUBAGENT FILE WRITE COLLISIONS (TOTAL BAN): Never assign two subagents to create or modify the same file. Shared files belong exclusively to the lead orchestrator. Subagents must have strictly disjoint, non-overlapping file boundaries to prevent file lock crashes.
- [ ] NO READ-ONLY SUBAGENT FILE WRITES (TOTAL BAN): Subagents of `TypeName: "research"` have read-only tools and cannot write files. Never command a research subagent to write files; research subagents discover and report back, and the lead writes the unified plan.
- [ ] NO INVALID SUBAGENT TYPENAMES (TOTAL BAN): Never pass custom strings like "Research 01" or "Worker 01" in the `TypeName` field of `invoke_subagent`. `TypeName` MUST be strictly `"research"` (read-only discovery) or `"self"` (read/write execution). Put descriptive names in `Role` only.
- [ ] NO SUBAGENT GIT EXECUTION (TOTAL BAN): Subagents must NEVER run any git commands (`git add`, `git commit`, `git push`, `git status`, `git diff`, `git checkout`) in shared workspaces. Worker git execution creates `.git/index.lock` collisions that immediately crash parallel workers with exit code 128.
- [ ] ISSUE & RCA DESTINATION ROUTING: Whenever resolving an issue or performing a fix with RCA, verified that CI/CD failures are documented in `.ai-memory/cicd-issues/nn-<slug>.md` (indexed in `.ai-memory/cicd-index.md`), while non-CI/CD issues (application bugs, logic/runtime defects) are documented in `02-spec/22-app-issues/nn-<slug>.md` (indexed in `02-spec/22-app-issues/readme.md`).
- [ ] NO TEST RUNNING (TOTAL BAN): Never run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] NO BUILD CHECKING (TOTAL BAN): Never run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- [ ] NO RUNNER SCRIPTS (TOTAL BAN): Never launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] NO AUTOMATIC RELEASES (TOTAL BAN): Never bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] NO RAW GIT COMMITS OR CONVENTIONAL COMMIT PREFIXES (TOTAL BAN): Never execute `git commit`, `git commit -m "..."`, `git add -A`, `git add .`, or raw `git push`. Never use conventional commit prefixes (`docs(...):`, `feat(...):`, `fix(...):`, `chore(...):`) in raw git commands. All staging, committing, and pushing MUST be executed exclusively by GitMap: `gitmap cpf "<module> - <summary>"` (features) or `gitmap cpb "<module> - <summary>"` (fixes). GitMap automatically formats, stages, commits, and pushes atomically. Always mention and format as a hyphen (`-`); no need to provide a colon in GitMap commit commands (`gitmap cpf` / `cpb` / CVF) because the colon is already automatically provided by GitMap (`Feature: ` or `Bug: `).
- [ ] NO INTERMEDIATE COMMITS (TOTAL BAN): Never commit after Phase 1 (e.g. committing plans or specs) or mid-Phase 2 (committing individual files or tests). Committing early pollutes git history, creates race conditions, and breaks atomicity. All changes (specs, plans, code modifications, index updates) MUST be committed together in ONE single atomic GitMap commit at the final step of Phase 3.
- [ ] NO PER-FILE COMMITTING (TOTAL BAN): Never commit each file individually as you work (e.g. running `git commit` or `gitmap cpf "<module> - <summary>"` after editing File 1, then another commit after File 2). Committing file-by-file pollutes git history, creates subagent lock collisions, and breaks atomic changes. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.
- [ ] NO RAPID CI/CD POLLING (TOTAL BAN): Never query or loop rapidly (`gh run view` in tight loops) when inspecting remote CI/CD pipelines. Agents must query pipeline state using GitMap Pipeline-AI (`gitmap pipeline-ai status --json` or `gitmap pl-ai status -t <sec>`) and strictly wait or sleep based on `etaSeconds` to eliminate credit waste.
- [ ] NO PREMATURE TURN CLOSING BEFORE EXECUTION (TOTAL BAN): Never halt execution, conclude the turn, or ask the user for permission after generating specs or subtasks. Planning constitutes only 50% of the task budget; you must proceed unconditionally to Phase 2 code execution. (Note: When dispatching asynchronous background subagents via `invoke_subagent`, yielding control to allow platform reactive wakeup is mandatory and is exempt from this ban).
- [ ] NO HORIZONTAL TASK CONCATENATION (TOTAL BAN): Never concatenate tasks horizontally in the Task Completion Summary (e.g. NEVER `✅ #1... ✅ #2...` run-on). Every completed task MUST be rendered on its OWN SEPARATE LINE starting with an individual markdown list bullet (`- ✅`).
- [ ] INDEPENDENT AI VERIFICATION PROMPT MANDATE: Emitted the self-contained independent AI verification and audit prompt linking to the canonical spec, consolidated plan, and modified files with verbatim score audit criteria.
- [ ] GITMAP HEAVY USAGE & ROUTINE PULL BAN: Heavily leveraged GitMap commands (`cpf`, `cpb`, `cpr`, `search`, `find`, `pwsh`) for discovery, execution, and commits. Never ran `pull-all` (`gitmap pa` or `gitmap pae`) unconditionally during routine turns; only ran `gitmap pae --json` when explicitly commanded by the user.
- [ ] NO RG, RIPGREP, OR RAW SHELL SEARCHES (TOTAL BAN): Never run `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, `Get-ChildItem -Recurse`, `findstr`, or slow shell search pipelines to search code. All code searching, symbol discovery, and regex scans MUST use GitMap high-speed search tools: `gitmap aum search "<query>" [dir] [-e <.ext>] [-r]` (streaming multi-core text/regex search) or `gitmap search "<query>"` (indexed symbol search). Running `rg`, `ripgrep`, `Select-String`, or `git grep` is an immediate auto-reject failure.

---

## 13. Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

[/goal](slashCommand;goal) You must verify every item on this checklist before committing any code. If a subagent violated one of these rules, you must reject their work.

- [ ] Master Guidelines: Fully enforced every file in `02-spec/02-coding-guidelines/` and `.ai-memory/coding-guidelines.md`.
- [ ] Concrete Types Centralization (`types.go`): Extracted domain structs, raw generic instantiations, and Result wrappers into dedicated `types.go` files as single reusable named types with follow-through comments (never leak raw generics like `result.Result[*Config]`).
- [ ] Error Management: Enforced `02-spec/03-error-manage/` using domain-specific `*appfault.AppError`, never generic error.
- [ ] Boolean Conventions: All booleans begin with is or has only (all other prefixes like can, should, was, will, did, must are banned). No negatives (`!isSuccess` is banned; use `isFail`).
- [ ] Semantic Naming: Zero generic garbage names (`temp`, `data`, `obj`). Behavior-driven unit test names.
- [ ] Multi-Line Arguments (Rule 9a/9b): Signatures and call sites with >2 arguments formatted one argument per line with trailing commas.
- [ ] Line Endings & Encoding: Strictly Unix LF (`\n`) and UTF-8 without BOM.
- [ ] Function Sizing: Functions <= 8 lines preferred (hard cap 15 lines).
- [ ] Strict Relative Git Paths & Lowercase: Zero absolute paths (`/absolute/path/to/...`) or `file:///` URIs anywhere in the diff, ledger, plans, release notes, or release page; strictly relative paths and lowercase filenames.

---

## 14. Anti-Hallucination & Blast Radius Checklist

- [ ] Echo Back the Spec: Verified Acceptance Criteria from the Spec file verbatim.
- [ ] Pre-Commit Diff Proof: Verified `git status` shows actual modified files before committing.
- [ ] No Placeholder Search: Confirmed zero `TODO` or `\[.*\]` placeholders remain in modified files.
- [ ] Index Sync Deadman Switch: Every new file is explicitly linked in `readme.md` and enqueued in `.ai-memory/what-to-read.md`.
- [ ] Blast Radius Acknowledgment: Global search across codebase performed via `gitmap aum search "<symbol>" [dir]` to update all callers of modified symbols (never `Select-String` or `git grep`).
- [ ] Continuous Loop Maintained: Continuous self-loop executed until 100% complete without running banned test/build commands.
- [ ] Final Step Commit & Push Verified: Staged and committed all changes atomically via GitMap semantic commit commands using `<module> - <summary>` format: `gitmap cpf "<module> - <summary>"` (features, e.g. `gitmap cpf "CBF - implement user profile dashboard"`) or `gitmap cpb "<module> - <summary>"` (bug fixes, e.g. `gitmap cpb "AUM - validate regex without nil fallback"`). Mentioned and formatted with a hyphen `-` to separate module from summary; no need to provide a colon in GitMap `cpf`/`cpb`/CVF commit arguments because the colon is already automatically provided by GitMap (`Bug: ` or `Feature: `). TOTAL BAN on colons `:` inside the GitMap argument and conventional prefixes (`fix(...)`, `feat(...)`, `docs(...)`). TOTAL BAN on raw `git add -A` and `git commit`. Pushed to remote via GitMap in a single final command.

---

## 15. Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] MANDATORY FINAL COMMIT & PUSH VIA GITMAP (ANYHOW): At the final step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, use GitMap semantic commit commands exclusively: `gitmap cpf "<module> - <summary>"` (features, e.g. `gitmap cpf "CBF - implement user profile dashboard"`) or `gitmap cpb "<module> - <summary>"` (bug fixes, e.g. `gitmap cpb "AUM - validate regex without nil fallback"`), which automatically stage, format commit messages, and push directly to remote. Mention and format as a hyphen `-` to separate module and summary; no need to provide a colon `:` in the GitMap `cpf`/`cpb`/CVF message argument as the colon is already provided automatically by GitMap (`Feature: ` or `Bug: `). TOTAL BAN on raw `git commit`, `git commit -m`, `git add -A`, colons `:` in the message argument, or conventional prefixes inside GitMap arguments (`docs(...)`, `feat(...)`, `fix(...)`). ZERO intermediate commits during Phase 1 or Phase 2; all files across the run are committed together at the final step. Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY): You must not create separate git commits for each individual file as you edit them (e.g. running `git commit` or `gitmap cpf` after editing File 1, then committing again after File 2 is strictly forbidden). Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback. All modified files, test change caches, and plan records across the turn must be accumulated in the working tree and committed together in a single grouped atomic commit at the final step before pushing.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase readme.md files left uncorrected, `.ai-memory/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `02-spec/` and `.ai-memory/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.
