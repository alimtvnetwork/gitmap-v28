# [V6] Execute Pending Tasks N-Step Continuous Loop & Mandatory Multi-Agent Subagent Orchestration — Workflow (must follow)

```text
N = 300 (Total self-loop steps budget — editable top-header parameter, default: 300)
A = 2   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
PHASE_1_BUDGET = N / 2   (Steps 1 .. 150: Pending Tasks Inventory, Parallel Discovery Subagents, Plan Alignment)
PHASE_2_BUDGET = N / 2   (Steps 151 .. 300: Mandatory Parallel Subagent Execution, Continuous Self-Looping, Targeted Quality Linting)
WAVES = ceil(pending_subtasks / (A x H))
```

> [!IMPORTANT]
> Prompt Version: 6.0.0
> Runtime: Google Antigravity 2.0 (IDE and CLI)
> Location: `01-prompts/14-execute/01-execute-pending-tasks.md`
> Skill: `pending-tasks` (`.agents/skills/pending-tasks/skill.md`)
> Invoke: /execute-pending-tasks
>
> **Top-Instruction Priority Mandate (Above Precedence / Preamble Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given ABOVE this prompt (including in the user preamble, header blocks, or incoming user request above) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below. The agent MUST inspect and follow the instructions above with absolute precedence.

[/goal](slashCommand;goal) Autonomously orchestrate and execute ALL pending tasks in `.ai-memory/plans/pending/` end-to-end: FIRST showcase and list out every pending task in visible chat during Turn 1, capture requirements verbatim, spawn autonomous subagents via `invoke_subagent` (A = 2, H = 2; solo execution without calling `invoke_subagent` is an auto-reject failure) in disjoint file boxes using GitMap high-speed commands as primary, prove every single claim with concrete evidence, enforce coding guidelines to 100%, and finish with one atomic GitMap commit (`gitmap cpf "<module> - <summary>"` / `gitmap cpb "<module> - <summary>"`) using hyphen format (no colons in GitMap arguments, as the colon is already provided by GitMap) holding strictly this task's files.

[/learn](slashCommand;learn) Enforce the Top-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are provided ABOVE this prompt outrank everything below. Turn 1 MUST showcase the pending task list in visible chat before any background execution. Each rule is stated once (R1 to R16) and cited by ID. Progress lives in the SQLite task manager and ledger, never only in chat.

[/plan](slashCommand;plan) Execute thorough step-by-step planning and pending task inventory in the repository before execution. Ensure all deliverables, architecture boundaries, and requirements are clearly defined in the audit ledger and subtask plans before dispatching worker waves.

### 🚨 MANDATORY SUBAGENT SPAWNING GATE (A = 2, H = 2 — ZERO SOLO EXECUTION ALLOWED)

- **ACTUAL TOOL CALL REQUIRED:** You must ACTUALLY CALL the `invoke_subagent` tool via your tool-calling API. Do NOT just print the text "Dispatched Worker..." and stop. If you only print text, the agents will not spawn and the task will fail! You must execute the `invoke_subagent` JSON tool payload.
- **3-STAGE MANDATORY DISPATCH:** You must invoke `A = 2` agents (`invoke_subagent`) at EVERY stage of the workflow:
  1. **Discovery & Inventory Step (A = 2 `research` subagents):** Spawn 2 read-only discovery subagents (`TypeName: "research"`) to inventory pending plans/subtasks, map affected source files, and return findings; lead writes the unified execution wave plan.
  2. **Spec & Wave Alignment Step (A = 2 `self` subagents or lead):** Spawn 2 subagents (`TypeName: "self"`) to author or verify modular, disjoint spec files and subtask plans.
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

### Phase 1A: Verbatim Capture, Pending Task Inventory & Chat Output Gate (Step 0)

Before executing any file modifications, scans, or code changes, you must execute Phase 1A:

1. **Top-Instruction Priority Verification:** Whatever directives, constraints, checklists, or instructions are given before this section or prompt (user preamble, header constraints, prior instructions) must be verified as highest priority and non-negotiable.
2. **Showcase Given Task First (Turn 1 Action):** In your VERY FIRST response turn upon receiving the prompt, you MUST inspect `.ai-memory/plans/pending/` and output the confirmed pending task breakdown directly in visible chat. Never execute tools silently without displaying the task breakdown to the user first!
3. **Lossless Verbatim Capture:** Store incoming prompt losslessly under `## User Request (Verbatim)` in canonical spec and parent plan.
4. **Screenshots & Media:** Decode base64/screenshots immediately into `assets/screenshots/<slug>-<NN>.png`. Reference via relative markdown links (`![Screenshot](assets/screenshots/<slug>-<NN>.png)`).
5. **Discrete Deliverables Extraction:** Break down all pending items into discrete, actionable items with ordered traceable IDs (`Task-01`, `Task-02`, ...).
6. **Mandatory Same-Turn Tool Chaining (TOTAL BAN ON TURNING OFF):** Emit the breakdown in chat with clean vertical formatting, and in the **EXACT SAME TURN**, invoke your first tool call (e.g. SQLite task manager init or preflight discovery). NEVER emit text alone (which ends the turn prematurely), and never ask "Should I proceed?".

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

Proceeding directly to Preflight & Phase 1B Inventory (Active Tool Call Running Below).
```

---

## 1. Precedence Hierarchy & Scope (Highest First)

1. **User Instructions & Preamble:** Directives and parameters ABOVE this prompt outrank everything below.
2. **Platform Limits:** Native tools, Artifact Review Policy, permission prompts, hooks. Never claim to override them.
3. **Repo Rules:** `agents.md`, `.ai-memory/strictly-avoid.md`, and `coding-guidelines.md`.
4. **This Prompt.**

If sources conflict, follow stricter one and record under `Conflicts:` in ledger.

### Scope Control Rules
- **Turn 1 Task Showcase:** In your very first response turn, you MUST showcase and list out the pending tasks in visible chat. Running tools silently without presenting the task breakdown is strictly banned.
- **Read budget:** Read only requested paths, search hits, and Step 0 context.
- **History read-only:** Never edit past events, changelogs, completed plans, release notes, or `06-old-prompts/` and `19-old-execute-prompts/`.
- **Out-of-scope:** Log under `Follow-ups:` in plan; never fix in this run.
- **Minimal diff:** Change only lines required; never reflow unaffected lines.
- **Indexes:** Update `01-prompts/readme.md` and `.ai-memory/prompts.md` only when adding/modifying prompts. Update `.ai-memory/plans/readme.md` every run. Register recent completed tasks before push.
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
  - For features: `gitmap cpf "<module> - <feature summary>"` (e.g. `gitmap cpf "CBF - implement user profile dashboard"`). GitMap automatically prepends `Feature: ` (which already provides the colon), so mention and format as a hyphen `-`; there is NO need to provide a colon in the GitMap `cpf` (or `cpb`/`cpr`) message argument. DO NOT provide a colon; use a hyphen `-` to separate module/scope from summary. NEVER include `feat(...)` or `feature:` in your message.
  - For bug fixes: `gitmap cpb "<module> - <fix summary>"` (e.g. `gitmap cpb "AUM - validate regex without nil fallback"`). GitMap automatically prepends `Bug: ` (which already provides the colon), so mention and format as a hyphen `-`; there is NO need to provide a colon in the GitMap message argument. DO NOT provide a colon; use a hyphen `-` to separate module/scope from summary. NEVER include `fix(...)` or `bug:` in your message.
  GitMap stages, formats, commits, and pushes atomically. TOTAL BAN on raw git commits (`git commit`, `git commit -m "..."`, `git add -A`, raw `git push`), colons inside the GitMap message argument, and conventional prefixes (`docs(...)`, `feat(...)`, `fix(...)`, `chore(...)`). ZERO intermediate commits: never commit during Phase 1 (plans/specs) or Phase 2; all files across the turn MUST be committed together at the final step of Phase 3. Before GitMap, all push gates must pass (targeted checks, secrets gate, and `.gitignore` hygiene; untrack any ignored files: `git rm --cached`). Push rejected: `git pull --rebase`, re-run command. Miss after push: allow one follow-up `gitmap cpb "<module> - <fix summary>"`, logged as `FOLLOW_UP_PUSH: <sha>`. Never amend pushed commits. Workers never run git commands or GitMap commit tools; only lead does.
- **R10 Zero Unauthorized Releases.** Never bump versions, edit `version.json`, update changelogs, or trigger release scripts unless user explicitly requested release.
- **R11 Strict Relative Git Paths & Lowercase Hygiene.** Strict ban on absolute paths (`C:\...`, `/home/...`) and `file:///` URIs. Only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well. Paths relative from git root. All filenames, documentation, and specs strictly lowercase (e.g. `readme.md`, `agents.md`, `skill.md`).
- **R12 No Polling / Immediate Turn Yielding.** When calling `invoke_subagent`, make it the sole tool action at turn end, print progress line (`Dispatched Worker 01 .. Worker <A> (wave k / WAVES); waiting for their results.`) and **STOP CALLING TOOLS**. Never poll in loop. Check `manage_subagents` once if wave runs long.
- **R13 Two-Strike Retry Cap & Anti-Looping.** Tool failing twice: worker replies `STATUS: BLOCKED` with exact error and stops. Lead takes over and logs `LEAD_FALLBACK: <reason>`. Subtask failing two remediation rounds is marked `FAILED` with RCA.
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
| **Atomic Commits** | `gitmap cpf "<module> - <msg>"` (Feature) / `cpb` (Bug) | `gitmap cpf` | Stages, formats with prefix, and pushes atomically. Mention as a hyphen `-` (no need to provide a colon in GitMap `cpf`/`cpb` commit arguments because the colon is already automatically provided by GitMap in `Feature: ` or `Bug: `). |
| **Pipeline Waiting** | `gitmap pipeline-ai status --json` | `gitmap pl-ai` | Non-polling dynamic ETA CI/CD monitor |

---

## 4. Step 0: Preflight, Platform Handshake & SQLite Task DB Initialization (Phase 1 Budget)

1. **Platform Handshake:** Confirm tools (`invoke_subagent`, `send_message`, `manage_subagents`, `ask_question`, `write_to_file`, `replace_file_content`, `run_command`).
2. **Commands & Directory:** Confirm `gitmap --version` and `python --version` exit 0. Verify GitMap with harmless call (`gitmap lf readme.md`), not `--help`. `run_command` uses `Cwd` in workspace root, paths relative. Never cd to other drives or tool folders.
3. **Working Tree Cleanliness:** Run `git status --porcelain`. Record modified files in ledger; never touch them. Confirm root `readme.md` is lowercase. Read `.ai-memory/what-to-read.md`, `strictly-avoid.md`, `coding-guidelines.md`.
4. **SQLite Task DB & Deterministic Slug Initialization:**
   - Initialize or inspect task state via the Antigravity SQLite task manager:
     `python 03-ai-scripts/46-agent-sqlite-task-manager.py init --name "execute-pending-tasks" --budget 300`
   - Inspect Task DB Schema & Data Integrity Rules:
     `python 03-ai-scripts/46-agent-sqlite-task-manager.py schema`
   - Output Analysis & Crash Forensics:
     - If `action: "RESUME_FOUND"`: Forensic report identifies crashed agent state. Resume execution from the uncompleted subtask.
     - If `action: "INITIALIZED"`: Created dedicated run directory `.ai-memory/temp-agents/<nn>-execute-pending-tasks/` and SQLite database `agent-task.db` with WAL mode.
5. **Ledger Creation:** Also create `.ai-memory/temp-agents/<nn>-execute-pending-tasks/ledger.md` mirroring the DB initialization for human readability:

```markdown
# Ledger: nn-execute-pending-tasks
Request slug: execute-pending-tasks
Request first line: <verbatim first line>
Status: ACTIVE
Phase: 1    Wave: 0 / WAVES    Step: 1 / N
Last completed action: Phase 1A Capture & Pending Task Breakdown
Next action: Phase 1B Inventory & Wave Allocation
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

## 5. Phase 1: Pending Tasks Inventory & Wave Grouping (Steps 1 .. PHASE_1_BUDGET)

You must use `invoke_subagent` to delegate inventory discovery and wave partitioning.

1. **Inventory Step (A = 2 `research` Discovery Subagents):** Lead agent calls `invoke_subagent` to spawn 2 read-only discovery subagents:
   - `Research 01: Pending Plans & Subtasks Inventory` (scans `.ai-memory/plans/pending/` and `.ai-memory/plans/subtasks/` using `gitmap find "*.md"`)
   - `Research 02: Codebase Symbols & Blast Radius` (maps callers, type contracts, and files affected by the pending tasks using `gitmap aum search`)
   - Findings report sent to lead via `send_message`.
2. **Wave Grouping (Lead Orchestrator):** Sequence pending tasks into logical waves:
   - Wave 1: Schemas, DB, and query wrappers
   - Wave 2: Business logic and services
   - Wave 3: UI and documentation
3. **Populate Subtasks in SQLite Task DB:**
   `python 03-ai-scripts/46-agent-sqlite-task-manager.py add-subtasks --db <databasePath> --tasks-json '[{"code": "Task-01", "title": "<title>", "owned_files": ["<paths>"], "agent_role": "Worker 01"}]'`
4. **Readiness Gate:** Complete Phase 1 inventory and wave grouping within `PHASE_1_BUDGET` steps, then proceed **UNCONDITIONALLY** into Phase 2 execution. ZERO intermediate git commits during Phase 1!

---

## 6. Phase 2: Execution Waves (Steps (PHASE_1_BUDGET + 1) .. N)

Execute pending tasks wave by wave using parallel subagents (`invoke_subagent`).

```json
{
  "Subagents": [
    {
      "TypeName": "self",
      "Role": "Worker 01: [Assigned Wave Subtask]",
      "Model": "inherit",
      "Workspace": "inherit",
      "Prompt": "You are Worker 01 for pending tasks wave <K>.\n\n### Core Objective:\nExecute the assigned pending subtask: <title>.\n\n### Strict Boundaries:\n- You possess read-write tools (`write_to_file`, `replace_file_content`, `run_command`).\n- Disjoint file box: ONLY modify assigned files: <owned_files>.\n- WORKER GIT BAN: NEVER run git commands (`git add`, `git commit`, `git status`). Git is exclusively managed by the lead agent.\n- Zero builds and test suites (R1): NEVER run `go build`, `npm run build`, or `go test`.\n- Run only targeted file linters.\n- Enforce coding guidelines: boolean prefixes (`is*`/`has*`, no `== true`), guard clauses, `*appfault.AppError`, lowercase filenames, strict relative git paths.\n\nSend final report to lead citing concrete paths and exit codes."
    },
    {
      "TypeName": "self",
      "Role": "Worker 02: [Assigned Wave Subtask]",
      "Model": "inherit",
      "Workspace": "inherit",
      "Prompt": "You are Worker 02 for pending tasks wave <K>.\n\n### Core Objective:\nExecute the assigned pending subtask: <title>.\n\n### Strict Boundaries:\n- You possess read-write tools (`write_to_file`, `replace_file_content`, `run_command`).\n- Disjoint file box: ONLY modify assigned files: <owned_files>.\n- WORKER GIT BAN: NEVER run git commands.\n- Zero builds and test suites (R1).\n- Run only targeted file linters.\n- Enforce coding guidelines to 100%.\n\nSend final report to lead citing concrete paths and exit codes."
    }
  ]
}
```

- **Two-Strike Retry Cap (R13):** If a worker tool fails twice, reply `STATUS: BLOCKED` with exact error and stop. Lead takes over.
- **Wave Transition:** Move completed subtasks to `.ai-memory/plans/completed/`, update `.ai-memory/plans/readme.md`, and proceed to the next wave.

---

## 7. Phase 3: Task Consolidation, Secrets Gate & Atomic Commit (End of Turn)

When all pending tasks are completed:

1. **Subtask Consolidation:** Combine completed granular subtasks from `.ai-memory/plans/subtasks/xx-<slug>/*.md` into `.ai-memory/plans/completed/xx-<slug>.md`, remove individual subtask files, and update `.ai-memory/plans/readme.md`.
2. **Pre-Commit Secrets Gate (Mandatory):**
   - Check status via `git status --porcelain`.
   - Run `python linter-scripts/check-forbidden-strings.py`.
   - Run `gitmap aum search -r "(BEGIN [A-Z ]*PRIVATE KEY|AKIA[0-9A-Z]{16}|gh[pousr]_[A-Za-z0-9]{36}|sk-[A-Za-z0-9]{20,}|xox[baprs]-[A-Za-z0-9-]{10,})"`.
   - On secret detection: offload via `gitmap rs text "<value>" --slug <slug>` into `repo-secrets`.
3. **Atomic Commit & Push via GitMap (Hyphen Format Mandate):**
   - Feature: `gitmap cpf "<module> - <summary>"`
   - Bug Fix: `gitmap cpb "<module> - <summary>"`
   - Total ban on colons inside the message argument (GitMap already provides `Feature: ` / `Bug: `).
   - TOTAL BAN on raw git commits (`git commit -m`).
4. **Task Completion Summary Output:**
   Emit a clean, line-by-line completion summary:

```markdown
### Task Completion Summary

- ✅ **Task-01: [Descriptive Task Title]** — `[Completed]`
- ✅ **Task-02: [Descriptive Task Title]** — `[Completed]`

### Modified Files Summary

- [relative/path/to/modified/file1.ext]
- [relative/path/to/modified/file2.ext]

### Implementation Confidence Score

- Confidence: 100%
- Rationale: Verified targeted checks passed, zero builds/tests run (R1), secrets gate passed, atomic push completed.
```

---

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.ai-memory/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `02-spec/` and `.ai-memory/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, do NOT run builds or tests during routine turns (build and test verification deferred to CI/CD), group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

---

## Metadata

- slug: execute-pending-tasks
- status: active
- version: 6.0.0
