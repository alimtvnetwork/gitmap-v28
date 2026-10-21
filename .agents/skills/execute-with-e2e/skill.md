---
name: execute-with-e2e
description: >-
  Autonomously execute end-to-end verification, run all local test suites to 100% green pass, inspect recent git commits and test cases, verify audit logs, and finalize with an automated minor version bump release ceremony.
---

# High Priority Instruction: Execute Parent Task with End-to-End (E2E) Testing & Automated Minor Release

```text
N = 400 (Total self-loop steps budget — editable top-header parameter, default: 400)
A = 3   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 3)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 3 agents × 2 hands = 6 concurrent subtask operations
Tree Hierarchy = Lead Orchestrator (Root Node) -> Branch Coordinators (Subtask Owners) -> Leaf Workers & Dynamic Joining Assistants
PHASE_1_BUDGET = 100     (Steps 1 .. 100: Planning, Discovery Subagents, Up-Front Brain & Roadmap Creation, Lean Subtasks)
PHASE_2_BUDGET = 150     (Steps 101 .. 250: Intermediate Archival, Sanitization, Prerequisite Refactoring, Intermediate Checkpoints)
PHASE_3_BUDGET = 150     (Steps 251 .. 400: Tree Swarm Parallel Execution, Dynamic Join Capacity Saturation, Full Local Test Suites, E2E Verification, Audit Log Review, Minor Bump Release Ceremony, Atomic Push)
WAVES = ceil(subtasks / (A x H))
```

> [!IMPORTANT]
> Prompt Version: 6.0.0
> Runtime: Google Antigravity 2.0 (IDE and CLI)
> Invoke: /execute-with-e2e <task>
>
> **Top-Instruction Priority Mandate (Above Precedence / Preamble Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given ABOVE this prompt (including in the user preamble, header blocks, or incoming user request above) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below. The agent MUST inspect and follow the instructions above with absolute precedence.

[/goal](slashCommand;goal) Autonomously orchestrate and execute the parent task end-to-end with comprehensive E2E testing, full local test suite verification, and automated minor release: FIRST showcase and list out the given task in visible chat during Turn 1 with the 5 Mandatory Pillars, the last 5 commits (`git log -n 5 --oneline`), recent test cases/suites inventory, and the 6-part verification schema, capture it verbatim, plan it in the repo, spawn autonomous subagents via `invoke_subagent` (A = 3, H = 2, Capacity = 6; solo execution without calling `invoke_subagent` is an auto-reject failure) in disjoint file boxes using GitMap high-speed commands as primary, run and pass all local test suites (`go test ./...`, `npm test`, `pytest`, etc.) to 100% green, verify cognitive audit logs in `.ai-memory/heatmap.db` and SQLite task engine, execute an automated minor version bump release ceremony (`npm run bump minor` or `python 03-ai-scripts/37-bump-version.py --tier minor`), and finish with an atomic GitMap release commit (`gitmap cpr "<module> - <summary>"` / `gitmap cpf "<module> - <summary>"`) using hyphen format (no colons in GitMap arguments, as the colon is already provided by GitMap) holding strictly this task's files.

[/learn](slashCommand;learn) Enforce the Top-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are provided ABOVE this prompt outrank everything below. Turn 1 MUST showcase the given task list, recent git commits (`git log -n 5 --oneline`), and recent test cases/suites inventory in visible chat before any background execution. Unlike standard V6, this prompt INVERTS Rule R1 into mandatory local test suite execution and INVERTS Rule R10 into an automated minor release finalization once 100% green. Each rule is stated once (R1 to R16) and cited by ID. Progress lives in the ledger and in `.ai-memory/plans/`, never only in chat.

[/plan](slashCommand;plan) Execute thorough step-by-step planning in the repository before execution. Ensure all deliverables, architecture boundaries, test matrices, and requirements are clearly defined in the audit ledger and subtask plans before dispatching worker waves.

### 🚨 MANDATORY SUBAGENT SPAWNING GATE (A = 3, H = 2, CAPACITY = 6 — TREE SWARM & DYNAMIC JOIN)

- **ACTUAL TOOL CALL REQUIRED:** You must ACTUALLY CALL the `invoke_subagent` tool via your tool-calling API. Do NOT just print the text "Dispatched Worker..." and stop. If you only print text, the agents will not spawn and the task will fail! You must execute the `invoke_subagent` JSON tool payload.
- **TREE-STRUCTURED SWARM HIERARCHY:** Organize execution as an operational tree:
  1. **Root Orchestrator (Lead Agent):** Owns the unified Brain roadmap, audit ledger, and final GitMap release commit. Never executes file edits solo.
  2. **Branch Coordinators (Subtask Primary Owners):** Own designated disjoint file boundaries (`Owned Files`) and execute modifications.
  3. **Leaf Workers & Dynamic Joining Assistants:** Assist primary owners with parallel reading, symbol caller discovery, memory generation, test authoring, and companion lint verification.
- **3-STAGE MANDATORY DISPATCH:** You must invoke `A = 3` agents (`invoke_subagent`) at EVERY stage of the workflow:
  1. **Planning Step (A = 3 `research` subagents):** Spawn 3 read-only discovery subagents (`TypeName: "research"`) to research the codebase, callers, and test suites; lead writes the unified Brain execution plan and roadmap.
  2. **Spec Step (A = 3 `self` authoring subagents — ZERO SOLO SPEC AUTHORING):** Spawn 3 authoring subagents (`TypeName: "self"`) to author modular, disjoint spec files, test matrices, and subtasks (never writing the same file). Solo authoring by the lead orchestrator is STRICTLY FORBIDDEN.
  3. **Execution Step (A = 3 `self` worker subagents):** Spawn 3 worker subagents (`TypeName: "self"`) to execute code modifications and test additions in strictly disjoint file boxes.
- **MANDATORY 3 SUBAGENTS IN SINGLE TOOL CALL:** Every call to invoke_subagent MUST contain an array of EXACTLY 3 subagent objects: "Subagents": [agent1, agent2, agent3]. Calling invoke_subagent with fewer than 3 subagent objects, or executing solo without calling invoke_subagent, is an immediate protocol violation and auto-reject failure.
- **DYNAMIC AGENT JOIN & IMMEDIATE RESPAWN (NO IDLE WORKERS):** Worker capacity `A` must remain pinned at 3 active agents (System Concurrency Capacity = 6) throughout execution. The lead orchestrator MUST NOT wait for an entire 3-agent wave to finish. As soon as ANY single worker completes, the lead immediately spawns a replacement worker for the next pending subtask, or commands the worker to dynamically join an active subtask in collaborative assistance mode (parallel reading, companion linters, test runner isolation, memory logging) to keep concurrency at exactly 3 active subagents (100% duty cycle).
- **SOLO EXECUTION IS AN AUTO-REJECT FAILURE:** The lead orchestrator is **STRICTLY FORBIDDEN** from executing planning, spec writing, or code changes by itself without calling the `invoke_subagent` tool. Failing to call the actual tool is a critical protocol violation.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/execute-with-e2e/skill.md` does not exist in the workspace, create it now.
2. Extract the core instructions of this prompt and save it into that `skill.md` using the standard YAML frontmatter (with `name: execute-with-e2e` and descriptive `description`).
3. Once installed, rely on progressive disclosure for future runs. Do not keep the entire prompt in active memory if not needed.

---

## The Unified Master Pipeline (Atomic Numbered Steps)

Execute this task via a strict 3-Phase pipeline. Do not skip steps.

### Phase 1A: Verbatim Capture, Task Extraction & Chat Output Gate (Step 0)

Before executing any file searches, scans, spec writing, or code changes, you must execute Phase 1A:

1. **Top-Instruction Priority Verification:** Whatever directives, constraints, checklists, or instructions are given before this section or prompt (user preamble, header constraints, prior instructions) must be verified as highest priority and non-negotiable.
2. **Extreme "Understood" Turn-1 Showcase (Turn 1 Action):** In your VERY FIRST response turn upon receiving the prompt, you MUST output a comprehensive, in-depth understanding block directly in visible chat before any background tool execution. Any generic boilerplate or superficial 1-sentence summaries (e.g. `Understood: [YES] — I will do this task`) are protocol auto-reject failures. Turn 1 MUST display the top-level mandate containing all 5 Mandatory Pillars, the last 5 commits recitation, the recent test case inventory, and the confirmed task breakdown using the 6-part verification schema (mandating explicit **Understood / AI Comprehension** and **Azure AI / Model Runtime Context** fields for every task):
   * **Pillar 1: Strategic Intent & Systems Reasoning:** Articulate overarching architectural goal, user intent, operational rationale, and systems integration mechanics. Explain *why* this change is necessary, how downstream components interact, and how E2E test verification proves complete system integrity.
   * **Pillar 2: Explicit Scope Boundaries:** Strict two-column partition defining In-Scope vs. strictly Excluded or Deferred items.
   * **Pillar 3: Non-Negotiable Constraints & Negative Guardrails Ingestion:** Explicit recitation of active constraints: raw search ban (`rg`, `ripgrep`, `Select-String`, `git grep`), worker git ban (no git commands in subagents), mandatory full local test suites (`go test ./...`, `npm test`, `pytest`) and E2E verification to 100% green, audit log review (`heatmap.db` and SQLite task engine), automated minor version bump release ceremony, strict relative Git paths, lowercase naming, zero secrets in repo (`repo-secrets` offloading), GitMap hyphen-separated commit format without colons.
   * **Pillar 4: Blast Radius & Dependency Topology:** Map impacted subsystems, files, downstream callers, test suites, and prompt/skill indexes with explicit risk containment strategies.
   * **Pillar 5: Multi-Agent Concurrency Topology:** Detail swarm structure (`A = 3`, `H = 2`, Capacity = 6), disjoint file assignments, and immediate dynamic respawn.
   * **Turn 1 Git Commit Recitation Mandate:** In Turn 1 visible chat, execute and display the output of `git log -n 5 --oneline` (the last 5 git commits) to anchor context in the repository's recent git history.
   * **Turn 1 Test Cases & Suites Inventory Mandate:** In Turn 1 visible chat, scan and display an inventory of all discovered unit tests, integration tests, E2E fixtures, and test suites across the repository that will be executed and verified 100% green.
   * **Task Comprehension & Model Runtime Invariants:** In the task breakdown, every single task MUST explicitly state `- **Understood / AI Comprehension:** [YES] — [1–2 concise sentences proving deep understanding of intent, scope, and verified constraints; answering did the AI understand the requirements or not]` and `- **Azure AI / Model Runtime Context:** [Explicitly document Azure AI, Azure OpenAI, or local model endpoint context if applicable, or state 'N/A — Local CLI/GitMap']`.
3. **Lossless Verbatim Capture:** Store incoming prompt losslessly under `## User Request (Verbatim)` in canonical spec and parent plan.
4. **Screenshots & Media:** Decode base64/screenshots immediately into `assets/screenshots/<slug>-<NN>.png`. Reference via relative markdown links (`![Screenshot](assets/screenshots/<slug>-<NN>.png)`).
5. **Discrete Deliverables Extraction:** Break down whatever user requirements were given into discrete, actionable items with ordered traceable IDs (`Task-01`, `Task-02`, ...).
6. **Mandatory Same-Turn Tool Chaining (TOTAL BAN ON TURNING OFF):** Emit the comprehension mandate and breakdown in chat with clean vertical formatting, and in the **EXACT SAME TURN**, invoke your first tool call (e.g. `write_to_file` to initialize ledger/spec, or run preflight). NEVER emit text alone (which ends the turn prematurely), and never ask "Should I proceed?".

```markdown
### 🧠 Deep Comprehension & Understood Mandate

#### Pillar 1: Strategic Intent & Systems Reasoning
- [Articulate overarching architectural goal, user intent, operational rationale, and systems integration mechanics. Explain why this change is necessary, how downstream components interact, and how E2E testing validates full system behavior.]

#### Pillar 2: Explicit Scope Boundaries (In-Scope vs. Excluded / Deferred)
- **In-Scope:** [Strict enumeration of files, modules, test suites, and changes executed in this turn.]
- **Excluded / Deferred:** [Explicitly deferred items, preserved untouched areas, and follow-up backlog.]

#### Pillar 3: Non-Negotiable Constraints & Negative Guardrails Ingestion
- [Explicit recitation of active constraints: raw search ban (no rg/ripgrep/Select-String), worker git ban (no git commands in subagents), MANDATORY full local test suites and E2E verification to 100% green pass, audit log review, automated minor version bump release ceremony upon 100% green, strict relative Git paths, lowercase naming, zero secrets in repo (repo-secrets offloading), GitMap hyphen-separated commit format without colons.]

#### Pillar 4: Blast Radius & Dependency Topology
- **Impacted Subsystems & Files:** [Files, downstream callers, test fixtures, and prompt/skill indexes affected.]
- **Containment & Risk Mitigation:** [File lock prevention via disjoint ownership, backward compatibility invariants.]

#### Pillar 5: Multi-Agent Concurrency Topology & Wave Orchestration
- **Swarm Structure:** Lead Orchestrator (Root Node) -> Branch Coordinators (Subtask Owners) -> Leaf Workers & Dynamic Joining Assistants
- **Capacity:** A = 3 agents × H = 2 hands = 6 concurrent operations, WAVES = ceil(subtasks / 6)
- **Disjoint Bounding Boxes:** Zero file ownership overlap; dynamic immediate respawn on worker completion (100% duty cycle pinned at 3 active subagents)

#### 📜 Recent Git Commits Recitation (`git log -n 5 --oneline`)
```text
[Output of git log -n 5 --oneline]
```

#### 🧪 Discovered Test Cases & Test Suites Inventory
- **Local Test Suites:** [e.g. go test ./..., npm test, pytest, vitest]
- **E2E Scenarios & Fixtures:** [Discovered test files, Playwright/Cypress/integration specs, and fixtures]
- **Verification Target:** 100% green exit code 0 across all discovered suites.

### 📋 Confirmed Task Breakdown & Requirement Ingestion

1. **Task-01: [Descriptive Task Title]**
   - **State:** `[IN PROGRESS — EXECUTING IMMEDIATELY]`
   - **Understood / AI Comprehension:** `[YES]` — [1–2 concise sentences proving deep understanding of intent, scope, and verified constraints; answering did the AI understand the requirements or not]
   - **Azure AI / Model Runtime Context:** [Explicitly document Azure AI, Azure OpenAI, or local model endpoint context if applicable, or state 'N/A — Local CLI/GitMap']
   - **Intent & Technical Deliverables:**
     - [Clear statement of technical goal, implementation artifacts, and test coverage]
     - [Precise implementation scope: functions, types, schemas, or prompt sections being added/edited]
     - [Quantifiable success metrics: 100% test pass rate, line count, field completeness, schema compliance]
   - **Negative Boundaries & Invariants:**
     - [Explicit list of things this task MUST NOT do or alter]
     - [Untouchable files, reserved packages, or excluded dependencies]
     - [Invariants that must remain preserved after modification]
   - **Disjoint File Boundaries:**
     - `Owned Files (Exclusive):` `[strictly/relative/path/to/file-01.md]`, `[strictly/relative/path/to/file-02.md]`
     - `Read-Only Context Files:` `[strictly/relative/path/to/reference.md]`
   - **Verification Evidence Contract:**
     - [Exact test commands to be executed: e.g. go test ./..., npm test, pytest]
     - [Expected exit code: exit 0, 100% green]
     - [Required diffstat evidence and audit log entry]

2. **Task-02: [Descriptive Task Title]**
   - **State:** `[QUEUED — EXECUTING NOW WITHOUT USER PROMPT]`
   - **Understood / AI Comprehension:** `[YES]` — [1–2 concise sentences proving deep understanding of intent, scope, and verified constraints; answering did the AI understand the requirements or not]
   - **Azure AI / Model Runtime Context:** [Explicitly document Azure AI, Azure OpenAI, or local model endpoint context if applicable, or state 'N/A — Local CLI/GitMap']
   - **Intent & Technical Deliverables:**
     - [Clear statement of technical goal and deliverable artifacts]
     - [Precise implementation scope: functions, types, schemas, or prompt sections being added/edited]
     - [Quantifiable success metrics: line count, field completeness, schema compliance]
   - **Negative Boundaries & Invariants:**
     - [Explicit list of things this task MUST NOT do or alter]
     - [Untouchable files, reserved packages, or excluded dependencies]
     - [Invariants that must remain preserved after modification]
   - **Disjoint File Boundaries:**
     - `Owned Files (Exclusive):` `[strictly/relative/path/to/file-03.md]`
     - `Read-Only Context Files:` `[strictly/relative/path/to/reference.md]`
   - **Verification Evidence Contract:**
     - [Exact test command or check to be executed to prove completion]
     - [Expected exit code: exit 0]
     - [Required diffstat evidence and audit log entry]

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
- **Turn 1 Task Showcase:** In your very first response turn, you MUST showcase and list out the given task in visible chat with commits and test suite inventory. Running tools silently without presenting the task breakdown is strictly banned.
- **Read budget:** Read only requested paths, search hits, and Step 0 context.
- **History read-only:** Never edit past events, changelogs, completed plans, release notes, or `06-old-prompts/` and `19-old-execute-prompts/`.
- **Out-of-scope:** Log under `Follow-ups:` in plan; never fix in this run.
- **Minimal diff:** Change only lines required; never reflow unaffected lines.
- **Indexes:** Update `01-prompts/readme.md` and `.ai-memory/prompts.md` only when adding/modifying prompts. Update `.ai-memory/plans/readme.md` and `02-spec/21-app/readme.md` every run. Register recent completed tasks before release.
- **Mandatory Multi-Agent Partitioning:** Even for tasks touching few files, work MUST be partitioned across A workers (e.g. Worker 01 implements changes, Worker 02 implements verification/linters/memory docs). Solo execution is strictly banned.
- **Automated Minor Release Finalization:** Unlike standard V6, `execute-with-e2e` mandates executing an automated minor version bump release ceremony as soon as all tests pass 100% green and audit logs are verified.

---

## 2. Core Operational Rules (Cite by ID)

- **R1 Mandatory Full Local Test Suites & E2E Verification (INVERTED FROM V6).** Unlike standard lean V6 turns where full test suites are banned, in `08-execute-with-e2e` you MUST run all relevant local test suites (`go test ./...`, `npm test`, `pytest`, vitest, etc.) and all recent or mentioned E2E test cases locally. Every single test MUST pass 100% green with exit code 0 before proceeding to release. In-flight intermediate build/test during code editing by workers can be run or verified locally to guarantee zero regressions.
- **R2 Targeted Checks AND Full Test Suites.** Run both targeted file-scoped linters (Section 10) AND full project test suites. A check scanning 0 files is a **FAIL**.
- **R3 Evidence or It Did Not Happen.** Every `DONE`, `PASS`, or "verified" claim MUST cite concrete test execution logs, exit code (`exit 0`), and diffstat. Vague assurances are auto-rejected.
- **R4 Never Invent Commands, Flags, or Paths.** Verify commands with a harmless call (`gitmap lf readme.md`), not `--help`. Use documented fallbacks and log in ledger.
- **R5 Mandatory Subagents (`invoke_subagent`), Concurrency Scaling & Immediate Dynamic Respawn.** Spawning subagents via `invoke_subagent` (`A = 3`, `H = 2`, System Concurrency Capacity = 6 concurrent operations) is an **ABSOLUTE MUST** (`research` for discovery in Phase 1, `self` for spec authoring and code edits in Phase 2/3). The lead agent is STRICTLY FORBIDDEN from executing all reads or edits solo. Organize execution as a Tree Swarm Hierarchy (Root Orchestrator -> Branch Coordinators -> Leaf Workers & Dynamic Assistants). Under the Immediate Dynamic Respawn & Join Protocol, the lead orchestrator MUST NOT wait for an entire 3-agent wave to finish. As soon as ANY single worker completes, the lead immediately spawns a replacement worker for the next pending subtask, or commands the worker to dynamically join an active subtask in collaborative assistance mode (parallel reading, companion linters, test fixture execution, memory logging) to keep concurrency at exactly 3 active subagents (100% duty cycle) without file write collisions (respecting R6). Solo execution without calling `invoke_subagent` is an auto-reject failure on the same tier as Rule 0.
- **R6 One Owner Per File (Disjoint Bounding Boxes).** Within every worker wave, each file has exactly one owner. Shared indexes (`.ai-memory/plans/readme.md`, `.ai-memory/prompts.md`, `.ai-memory/what-to-read.md`, `02-spec/21-app/readme.md`, directory `readme.md`) belong exclusively to lead.
- **R7 Git Safety & Isolation (Worker Git Ban).** Subagents NEVER run git commands (`git add`, `git commit`, `git push`, `git status`, `git diff`, `git checkout`). In shared workspaces (`Workspace: "inherit"`), worker git calls create `.git/index.lock` collisions that immediately crash parallel agents. Nobody runs `git reset --hard`, `git checkout --`, `git clean`, `git stash`, or force pushes.
- **R8/R9 Atomic Commit & Push via GitMap (TOTAL BAN ON RAW GIT COMMITS).** The run ends with ONE GitMap call:
  - For releases: `gitmap cpr "<module> - <release summary>"` (e.g. `gitmap cpr "E2E - release minor bump with 100 percent green tests"`). GitMap automatically prepends `Release: ` (which already provides the colon), so mention and format as a hyphen `-`; there is NO need to provide a colon in the GitMap `cpr` (or `cpf`/`cpb`) message argument. DO NOT provide a colon; use a hyphen `-` to separate module/scope from summary. NEVER include `release(...)` or `chore:` in your message.
  - For features: `gitmap cpf "<module> - <feature summary>"`.
  - For bug fixes: `gitmap cpb "<module> - <fix summary>"`.
  GitMap stages, formats, commits, and pushes atomically. TOTAL BAN on raw git commits (`git commit`, `git commit -m "..."`, `git add -A`, raw `git push`), colons inside the GitMap message argument, and conventional prefixes (`docs(...)`, `feat(...)`, `fix(...)`, `chore(...)`). ZERO intermediate commits: never commit during Phase 1 (plans/specs) or Phase 2; all files across the turn MUST be committed together at the final step of Phase 3. Before GitMap, all push gates must pass (targeted checks, secrets gate, and `.gitignore` hygiene; untrack any ignored files: `git rm --cached`). Push rejected: `git pull --rebase`, re-run command. Miss after push: allow one follow-up `gitmap cpb "<module> - <fix summary>"`, logged as `FOLLOW_UP_PUSH: <sha>`. Never amend pushed commits. Workers never run git commands or GitMap commit tools; only lead does.
- **R10 Automated Minor Release Finalization (INVERTED FROM V6).** Unlike standard V6 where releases are banned unless requested, `08-execute-with-e2e` mandates an automated minor release ceremony once all test suites and test cases are confirmed 100% green and audit logs are verified. Execute `npm run bump minor` or `python 03-ai-scripts/37-bump-version.py --tier minor`, update changelog and root readme version pin, and push via GitMap release commit (`gitmap cpr "<module> - <summary>"` or `gitmap cpf "<module> - <summary>"`).
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

GitMap is your **PRIMARY** acceleration engine. Banned ALL variants of `Get-ChildItem` (with or without flags/parameters), `dir`, `ls` in shell, `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, and `findstr`. Execute all searches and operations through GitMap:

| Operation | Primary GitMap Command | High-Speed Alias | Purpose & Advantage |
| :--- | :--- | :--- | :--- |
| **Live Streaming Search** | `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r] [-i]` | `gitmap aum grep "<pat>"` | Multi-core streaming live file search (replaces `Select-String`, `git grep`) |
| **Indexed Symbol Search** | `gitmap search "<query>" [--limit <n>]` | `gitmap search` | Instant SQLite cached keyword/symbol search across indexed repos |
| **Wildcard File Search** | `gitmap find "<pattern>" [-ext <ext>]` | `gitmap f "<pat>"` | Index-accelerated multi-core glob filename finder |
| **Directory Inventory** | `gitmap lf <dir>` / `gitmap find "*.*" <dir>` | `gitmap lf [pat]` | Instant indexed repository inventory (replaces `Get-ChildItem`, `dir`, `ls`) |
| **Substring File Search** | `gitmap find-files-any "<substring>"` | `gitmap ffa "<str>"` | High-speed partial filename matcher |
| **Stream File** | `gitmap cat <filepath>` | `gitmap cat` | Direct zero-disk memory streaming to stdout (replaces `cat`, `type`, `Get-Content`) |
| **PowerShell Runner** | `gitmap pwsh "<command>"` | `gitmap ps "<cmd>"` | High-speed PowerShell execution with `-NoProfile` |
| **Bash Runner** | `gitmap bash "<command>"` | `gitmap sh "<cmd>"` | Standard cross-platform Bash command execution |
| **Offload Secrets** | `gitmap rs file <filepath>` / `folder` / `text` | `gitmap rs` | Auto-commits into `repo-secrets` in work directory |
| **Offload Scripts** | `gitmap rc file <file.ps1>` / `text` | `gitmap rc` | Auto-commits reusable scripts into `repo-cache` |
| **Atomic Commits** | `gitmap cpr "<module> - <msg>"` (Release) / `cpf` / `cpb` | `gitmap cpr` | Stages, formats with prefix, and pushes atomically. Mention as a hyphen `-` (no colons in GitMap commit arguments). |
| **Pipeline Waiting** | `gitmap pipeline-ai status --json` | `gitmap pl-ai` | Non-polling dynamic ETA CI/CD monitor (allowed if requested) |
| **Agent Task Engine** | `gitmap task <init|add|claim|complete|fail|status|schema>` | `gitmap task` | Sub-millisecond compiled Go SQLite task manager for multi-agent workflows |
| **Heatmap Action Log** | `gitmap heatmap log --task <task> --action <READ\|EDIT\|CREATE\|DELETE> --file <path> --state <STATE> --expected "<goal>"` | `gitmap heatmap log` | Concurrency-safe cognitive action & reasoning telemetry into `.ai-memory/heatmap.db` |
| **Contention & Leases** | `gitmap agent heatmap [--json]` / `gitmap heatmap status` | `gitmap agent heatmap` | Pre-edit file lease conflict & contention hotspot analysis across active agents |
| **RC Heatmap Suite** | `gitmap rc heatmap <status|tables|query|view|tree|live|sync>` | `gitmap rc heatmap` | Unified RC heatmap commands with automatic vault auto-provisioning & roadmap linkage |
| **Heatmap Diagnostics** | `gitmap heatmap query|tables|view` / `gitmap heatmap rc ...` | `gitmap heatmap` | Non-blocking read-only SQL queries, tree view, and live activity stream |
| **Fleet Synchronization** | `gitmap sync [--workers 8] [--projects <json>]` | `gitmap sync` | Parallel Go sync engine across 43 repositories in <5s with 6-stage ceremony |

### 🔍 Code & Symbol Search Protocol (TOTAL BAN ON `rg`, `ripgrep`, `Select-String`, `git grep`, `Get-ChildItem`, `dir`, `ls`)
- **GitMap High-Speed Commands & Explicit Examples:**
  - **Directory file inventory:** `gitmap lf <dir>` (e.g. `gitmap lf src-tauri/src/modules/repo_db`) or `gitmap find "*.*" <dir>`
  - **Live streaming code/symbol search:** `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r] [-i]` (e.g. `gitmap aum search "RunFleetPASCommand" cli -e .go`)
  - **Finding files:** `gitmap find "<pattern>" [dir]` (e.g. `gitmap find "*.md" 02-spec`)
  - **Stream file:** `gitmap cat <filepath>` (e.g. `gitmap cat 02-spec/21-app/readme.md`)
- **TOTAL BAN (AUTO-REJECT FAILURE):** Banned ALL variants of `Get-ChildItem` (with or without flags/parameters), `dir`, `ls` in shell, `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, and `findstr`. Search exclusively through `gitmap aum search`, `gitmap find`, `gitmap lf`, or `gitmap cat`.

---

## 5. Step 0: Preflight, Platform Handshake & SQLite Task DB Initialization (Phase 1 Budget)

1. **Platform Handshake:** Confirm tools (`invoke_subagent`, `send_message`, `manage_subagents`, `ask_question`, `write_to_file`, `replace_file_content`, `run_command`). If `task_boundary` exists: set `PLANNING` (Phase 1), `EXECUTION` (Phase 2), `VERIFICATION` (Phase 3).
2. **Commands & Directory:** Confirm `gitmap --version` and `python --version` exit 0. Verify GitMap with harmless call (`gitmap lf readme.md`), not `--help`. `run_command` uses `Cwd` in workspace root, paths relative. Never cd to other drives or tool folders.
3. **Turn 1 Recitation Inspections:**
   - Run `git log -n 5 --oneline` and record recent 5 commits.
   - Scan for existing test suites and test fixtures (`gitmap find "*test*"`, `gitmap lf tests`, etc.).
4. **Working Tree Cleanliness:** Run `git status --porcelain`. Record modified files in ledger; never touch them. Confirm root `readme.md` is lowercase. Read `.ai-memory/what-to-read.md`, `strictly-avoid.md`, `coding-guidelines.md`.
5. **Step 0 Heatmap & RC Preflight Checks & Automatic Creation (Mandatory Inspection):**
   - **Automatic RC / Heatmap Creation:** If `repo-cache` (`rc`) or `.ai-memory/heatmap.db` does not exist, GitMap automatically provisions the required split-db structure linked to the active mission and workspace folder.
   - Run preflight status & inspect master catalog tables:
     `gitmap heatmap status` and `gitmap heatmap tables`
     (or `gitmap rc heatmap status` and `gitmap rc heatmap tables`)
   - Check cognitive action metrics: `gitmap heatmap status`
   - Check file contention hotspots: `gitmap agent heatmap`
   - Render hierarchical roadmap & touched file tree: `gitmap heatmap tree` (or `gitmap rc heatmap tree`)
   - **Mandatory Pre-Action Heatmap Logging:** Before any file read or edit:
     `gitmap heatmap log --task <task> --action <READ|EDIT|CREATE|DELETE> --file <path> --state <STATE> --expected "<goal>"`
6. **SQLite Task DB & Deterministic Slug Initialization (Check Before Creating):**
   - Initialize or inspect task state via GitMap SQLite task manager (compiled Go, sub-millisecond):
     `gitmap agent task init --name "<task name>" --budget 400`
   - **Task Already Created Detection & Existence Report:**
     When initializing or requesting a task, if the task was already created:
     * GitMap reports that the task is already created:
       - Displays existing Parent Task ID and created timestamp.
       - Lists sub-agents currently assigned or actively working.
       - Lists sub-tasks in progress, claimed, or completed (`gitmap agent subtask ls`).
       - Lists all files that have already been touched or worked on.
       - Reports the location of the RC database (`repo-cache/.heatmap/<slug>/heatmap.db` or `.ai-memory/heatmap.db`).
   - **Prompt MD5 Hash Verification:**
     * The agent records/verifies the MD5 hash of the incoming prompt with the task database (`--prompt-hash <md5>` or in task metadata).
     * If the prompt hash matches the prior task record: confirm exact task identity and resume seamlessly from the first uncompleted subtask!
     * If the prompt hash differs: detect that the prompt instructions have evolved, update the task roadmap while preserving previously logged file touches and completed subtasks.
   - Inspect Task DB Schema & Data Integrity Rules (Mandatory Preflight):
     `gitmap agent task status`
     Verify tables (`ParentTask`, `Subtask`, `AgentActionLog`), positive boolean flags (`IsActive`, `HasCompleted`, `IsBlocked`), and payload constraints before populating subtasks. (DDL: `gitmap agent task status --ddl`; JSON: `gitmap agent task status --json`).
   - TOTAL BAN ON RAW PYTHON SQLITE: NEVER execute `python -c "import sqlite3..."` or ad-hoc Python DB scripts. Inspect status via `gitmap agent task status` or list subtasks via `gitmap agent subtask ls`.
   - Output Analysis & Crash Forensics:
     - If `action: "RESUME_FOUND"`: A matching or similar slug exists in `.ai-memory/temp-agents/`! If `diagnostics.hasCrashesDetected: true`, inspect the forensic report (`diagnostics.crashedAgents`) to identify which agent crashed, what file it was touching, and the last logged action. Resume execution from the uncompleted subtask.
     - If `action: "INITIALIZED"`: Created dedicated run directory `.ai-memory/temp-agents/<nn>-<slug>/` and SQLite database `agent-task.db` with WAL mode (`PRAGMA journal_mode=WAL; PRAGMA busy_timeout=5000;`).
7. **Ledger Creation:** Also create `.ai-memory/temp-agents/<nn>-<slug>/ledger.md` mirroring the DB initialization for human readability:

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

1. **Planning Step (A = 3 `research` Discovery Subagents):** Lead agent calls `invoke_subagent` to spawn 3 read-only discovery subagents (`TypeName: "research"`, `Role: "Research 01: Architecture & Blast Radius"`, `Role: "Research 02: Specs, Guidelines & Test Case Mapping"`, `Role: "Research 03: Test Suites, Dependencies & Release Targets"`). Their prompt MUST instruct them to research the codebase using GitMap high-speed search (`gitmap aum search "<symbol>" [dir] [-e <.ext>]`, `gitmap find`, `gitmap lf`, `gitmap cat`) — NEVER `rg`, `ripgrep`, PowerShell `Select-String`, or `git grep` — inventory existing test suites and fixtures, define symbol boundaries and caller dependencies, and return their structured findings to the lead in their final message.
   - *Master Plan Generation (Lead Agent):* The lead orchestrator receives all three discovery reports, synthesizes findings, and writes the unified Execution Plan (`.ai-memory/plans/pending/nn-<slug>.md`) and the Root Task JSON Manifest including full test coverage matrix.
   - *Tool Call:* Lead must execute the `invoke_subagent` tool as the final action in the turn, print `Dispatched Planning Agents`, and then STOP CALLING TOOLS to wait for `<SYSTEM_MESSAGE>` reactive wakeup.
2. **Spec Step (A = 3 `self` Authoring Subagents):** Once planning is synthesized, the lead agent calls `invoke_subagent` to spawn 3 authoring subagents (`TypeName: "self"`). Their prompt MUST instruct them to author modular, strictly disjoint spec files and subtask plans:
   - Subagent 1 writes `02-spec/21-app/nn-<slug>/01-architecture-spec.md` and subtasks `.ai-memory/plans/subtasks/nn-<slug>/01-<name>.md`.
   - Subagent 2 writes `02-spec/21-app/nn-<slug>/02-component-and-e2e-spec.md` and subtasks `.ai-memory/plans/subtasks/nn-<slug>/02-<name>.md`.
   - Subagent 3 writes `02-spec/21-app/nn-<slug>/03-test-matrix-and-release-spec.md` and subtasks `.ai-memory/plans/subtasks/nn-<slug>/03-<name>.md`.
   - NEVER have subagents write to the same file path!
   - *Tool Call:* Lead must execute the `invoke_subagent` tool as the final action in the turn, print `Dispatched Spec Agents`, and then STOP CALLING TOOLS to wait for `<SYSTEM_MESSAGE>` reactive wakeup.
3. **Brain & Comprehensive Roadmap Creation (Mandatory Planning Gate):** Before entering Phase 2 execution, the lead orchestrator synthesizes discovery findings into a unified, concrete "Brain" roadmap authored in the specification (`02-spec/21-app/`) and execution plan (`.ai-memory/plans/`):
   - **Target Files Inventory:** Exhaustive list of every file to read, create, or modify.
   - **Modification Sites & Functions:** Exact function names, structs, interfaces, and logic blocks targeted for changes.
   - **E2E Test Suites & Verification Criteria:** Concrete list of local test suites (`go test ./...`, `npm test`, `pytest`) and specific test cases to be executed and passed 100% green.
   - **Data Contracts & Error Wrapping:** Concrete type definitions and structured errors (`*appfault.AppError`).
   - **Caller Blast Radius:** Callers and dependent packages mapped via GitMap search.
   - **Disjoint Subtask Allocation:** Clear file-ownership boundaries assigned to subtask codes (`Task-01`, `Task-02`) so parallel workers never collide.
4. **Populate Subtasks in SQLite Task DB:** Once subtasks are decomposed, inspect schema constraints (`gitmap agent task status`) and populate them into the SQLite database for atomic worker claiming:
   `gitmap agent subtask add --json '[{"code": "Task-01", "title": "<title>", "owned_files": ["<paths>"], "agent_role": "Worker 01"}]'`
5. **Readiness Gate:** Complete Phase 1 planning, Brain roadmap generation, and spec authoring within `PHASE_1_BUDGET` steps, then proceed **UNCONDITIONALLY** into Phase 2. ZERO intermediate git commits during Phase 1!

---

## 7. Phase 2: Execution Step (Worker Waves) (Steps (PHASE_1_BUDGET + 1) .. N)

> [!CRITICAL]
> **MANDATORY `invoke_subagent` DISPATCH (ZERO SOLO EXECUTION):**
> You MUST ACTUALLY CALL the `invoke_subagent` tool to spawn A workers (`A = 3`, `TypeName: "self"`, up to H subtasks per worker) in parallel. Executing all subtasks solo in main agent without the tool call is an immediate auto-reject failure.

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
      "Prompt": "You are Research 01 for task <nn>-<slug>. You have no prior chat context; this brief is your complete specification.\n\n### Core Objective:\nExplore the codebase, map symbol callers, trace dependencies, and discover relevant source files for: <task title and objectives>.\n\n### Tool Capabilities & Strict Read-Only Boundary:\n- You are a READ-ONLY subagent (`TypeName: 'research'`).\n- You have read tools: `run_command`, `view_file`, `search_web`, `read_url_content`.\n- You DO NOT have `write_to_file` or `replace_file_content`. NEVER attempt to create or edit files.\n\n### Mandatory GitMap Search Primacy & Total Ban on Raw Grep/Shell:\n- Execute ALL searches and symbol discoveries exclusively via GitMap:\n  * Directory file inventory: `gitmap lf <dir>` or `gitmap find \"*.*\" <dir>`\n  * Live symbol search: `gitmap aum search \"<symbol>\" [dir] [-e <.ext>] [-r] [-i]`\n  * Finding files: `gitmap find \"<pattern>\" [dir]`\n  * Stream file: `gitmap cat <filepath>` or native `view_file`\n- TOTAL BAN (AUTO-REJECT FAILURE): Banned ALL variants of `Get-ChildItem` (with or without flags/parameters), `dir`, `ls` in shell, `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, and `findstr`.\n- TOTAL BAN ON GIT COMMANDS: NEVER run `git` commands (`git add`, `git commit`, `git status`, `git diff`, etc.).\n\n### Required Findings Report (Send via send_message to Parent):\nWhen your research is complete, send a message to the caller containing:\n1. Key symbol definitions, structures, and entry points.\n2. Call sites and blast radius (files that will be affected by modifications).\n3. Recommended modular file ownership boundaries for worker subtasks.\n4. Verification commands and test suites to validate changes."
    },
    {
      "TypeName": "research",
      "Role": "Research 02: Specs, Guidelines & Test Case Mapping",
      "Model": "inherit",
      "Workspace": "inherit",
      "Prompt": "You are Research 02 for task <nn>-<slug>. You have no prior chat context; this brief is your complete specification.\n\n### Core Objective:\nInspect existing specifications under `02-spec/`, coding guidelines under `02-spec/02-coding-guidelines/`, prior plans, and existing unit/integration/E2E test files for: <task title and objectives>.\n\n### Tool Capabilities & Strict Read-Only Boundary:\n- You are a READ-ONLY subagent (`TypeName: 'research'`).\n- You DO NOT have file authoring tools. NEVER attempt to create or edit files.\n\n### Mandatory GitMap Search Primacy & Total Ban on Raw Grep/Shell:\n- Execute ALL searches exclusively via GitMap:\n  * Directory file inventory: `gitmap lf <dir>` or `gitmap find \"*.*\" <dir>`\n  * Live search in specs & tests: `gitmap aum search \"<term>\" 02-spec`\n  * Finding test files: `gitmap find \"*test*\" [dir]`\n  * Stream file: `gitmap cat <filepath>` or native `view_file`\n- TOTAL BAN: Banned ALL variants of `Get-ChildItem` (with or without flags/parameters), `dir`, `ls` in shell, `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, and `findstr`.\n- TOTAL BAN ON GIT COMMANDS: NEVER run `git` commands.\n\n### Required Findings Report (Send via send_message to Parent):\nSend a message to the caller detailing existing specs, test fixtures, architectural invariants, acceptance criteria, and positive boolean rules."
    },
    {
      "TypeName": "research",
      "Role": "Research 03: Test Suites, Dependencies & Release Targets",
      "Model": "inherit",
      "Workspace": "inherit",
      "Prompt": "You are Research 03 for task <nn>-<slug>. You have no prior chat context; this brief is your complete specification.\n\n### Core Objective:\nExplore test suites (`go test ./...`, `npm test`, `pytest`), bump version scripts (`03-ai-scripts/37-bump-version.py`), package manifests (`package.json`, `version.json`), skills, and release dependencies for: <task title and objectives>.\n\n### Tool Capabilities & Strict Read-Only Boundary:\n- You are a READ-ONLY subagent (`TypeName: 'research'`).\n- You DO NOT have file authoring tools. NEVER attempt to create or edit files.\n\n### Mandatory GitMap Search Primacy & Total Ban on Raw Grep/Shell:\n- Execute ALL searches exclusively via GitMap:\n  * Directory file inventory: `gitmap lf <dir>` or `gitmap find \"*.*\" <dir>`\n  * Live search: `gitmap aum search \"<term>\" [dir] [-e <.ext>]`\n  * Finding files: `gitmap find \"<pattern>\" [dir]`\n  * Stream file: `gitmap cat <filepath>`\n- TOTAL BAN: Banned ALL variants of `Get-ChildItem` (with or without flags/parameters), `dir`, `ls` in shell, `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, and `findstr`.\n- TOTAL BAN ON GIT COMMANDS: NEVER run `git` commands.\n\n### Required Findings Report (Send via send_message to Parent):\nSend a message to the caller detailing test runner configurations, test fixture paths, version bump scripts, and minor release prerequisites."
    }
  ]
}
```

#### Phase 1 Spec Authoring Payload (`TypeName: "self"` — Read/Write Specification Authors):
```json
{
  "Subagents": [
    {
      "TypeName": "self",
      "Role": "Spec Author 01: Architecture & Data Contracts",
      "Model": "inherit",
      "Workspace": "inherit",
      "Prompt": "You are Spec Author 01 for task <nn>-<slug>. You have no prior chat context; this brief is your complete specification.\n\n### Core Objective:\nAuthor the core architectural specification and subtask plans for architecture and preflight.\n\n### Owned Files (Strictly Disjoint — No Other Agent Touches These):\n- `02-spec/21-app/<slug>/01-architecture-spec.md`\n- `.ai-memory/plans/subtasks/<slug>/01-subtask-git-sync-and-preflight.md`\n\n### Invariants & Guardrails:\n- Strict Relative Git Paths ONLY (`02-spec/...`, `.ai-memory/...`). TOTAL BAN on absolute filesystem paths and `file:///` URIs.\n- Strict lowercase filenames across all files.\n- TOTAL BAN ON GIT COMMANDS: Subagents never run git commands (`git add`, `git commit`, `git push`, `git status`, etc.).\n- Code & Symbol Search: Use GitMap exclusively (`gitmap lf`, `gitmap find`, `gitmap aum search`, `gitmap cat`). Banned ALL variants of Get-ChildItem, dir, ls in shell, rg, ripgrep, grep, git grep, Select-String, and findstr.\n- Cognitive Heatmap Telemetry: Before writing files, log cognitive actions via `gitmap heatmap log --task <slug> --action CREATE --file <path> --state SPECIFICATION --expected \"<goal>\"`.\n- Use `write_to_file` to author files.\n\n### Reporting Contract:\nWhen files are written and verified, report completed files and diffstats via `send_message` to parent."
    },
    {
      "TypeName": "self",
      "Role": "Spec Author 02: Component & E2E Testing Specification",
      "Model": "inherit",
      "Workspace": "inherit",
      "Prompt": "You are Spec Author 02 for task <nn>-<slug>. You have no prior chat context; this brief is your complete specification.\n\n### Core Objective:\nAuthor the detailed component and E2E testing specification and subtask plans for implementation and tests.\n\n### Owned Files (Strictly Disjoint — No Other Agent Touches These):\n- `02-spec/21-app/<slug>/02-component-and-e2e-spec.md`\n- `.ai-memory/plans/subtasks/<slug>/02-subtask-core-implementation-and-unit-tests.md`\n\n### Invariants & Guardrails:\n- Strict Relative Git Paths ONLY (`02-spec/...`, `.ai-memory/...`). TOTAL BAN on absolute filesystem paths and `file:///` URIs.\n- Strict lowercase filenames across all files.\n- TOTAL BAN ON GIT COMMANDS: Subagents never run git commands (`git add`, `git commit`, `git push`, `git status`, etc.).\n- Code & Symbol Search: Use GitMap exclusively (`gitmap lf`, `gitmap find`, `gitmap aum search`, `gitmap cat`). Banned ALL variants of Get-ChildItem, dir, ls in shell, rg, ripgrep, grep, git grep, Select-String, and findstr.\n- Cognitive Heatmap Telemetry: Before writing files, log cognitive actions via `gitmap heatmap log --task <slug> --action CREATE --file <path> --state SPECIFICATION --expected \"<goal>\"`.\n- Use `write_to_file` to author files.\n\n### Reporting Contract:\nWhen files are written and verified, report completed files and diffstats via `send_message` to parent."
    },
    {
      "TypeName": "self",
      "Role": "Spec Author 03: Test Matrix, Release & Heatmap Specification",
      "Model": "inherit",
      "Workspace": "inherit",
      "Prompt": "You are Spec Author 03 for task <nn>-<slug>. You have no prior chat context; this brief is your complete specification.\n\n### Core Objective:\nAuthor the test matrix, minor bump release ceremony specification, and subtask plans for validation and release.\n\n### Owned Files (Strictly Disjoint — No Other Agent Touches These):\n- `02-spec/21-app/<slug>/03-test-matrix-and-release-spec.md`\n- `.ai-memory/plans/subtasks/<slug>/03-subtask-e2e-test-execution-and-minor-release.md`\n\n### Invariants & Guardrails:\n- Strict Relative Git Paths ONLY (`02-spec/...`, `.ai-memory/...`). TOTAL BAN on absolute filesystem paths and `file:///` URIs.\n- Strict lowercase filenames across all files.\n- TOTAL BAN ON GIT COMMANDS: Subagents never run git commands (`git add`, `git commit`, `git push`, `git status`, etc.).\n- Code & Symbol Search: Use GitMap exclusively (`gitmap lf`, `gitmap find`, `gitmap aum search`, `gitmap cat`). Banned ALL variants of Get-ChildItem, dir, ls in shell, rg, ripgrep, grep, git grep, Select-String, and findstr.\n- Cognitive Heatmap Telemetry: Before writing files, log cognitive actions via `gitmap heatmap log --task <slug> --action CREATE --file <path> --state SPECIFICATION --expected \"<goal>\"`.\n- Use `write_to_file` to author files.\n\n### Reporting Contract:\nWhen files are written and verified, report completed files and diffstats via `send_message` to parent."
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
      "Role": "Worker 01: Core Architecture & Implementation",
      "Model": "inherit",
      "Workspace": "inherit",
      "Prompt": "You are Worker 01 for task <nn>-<slug>. You have no prior chat context; this brief is your complete specification.\n\n### Owned Files (Exclusive Modification Rights):\n- `<relative/path/to/module1-file1.ext>`\n- `<relative/path/to/module1-file2.ext>`\n\n### Assigned Subtasks:\n- Subtask: `.ai-memory/plans/subtasks/<slug>/01-subtask-*.md`\n\n### Instructions & Invariants:\n- Follow the subtask plan and spec strictly. Implement required changes.\n- Strict Relative Git Paths ONLY. TOTAL BAN on absolute filesystem paths and file:/// URIs.\n- Strict lowercase filenames across all files.\n- TOTAL BAN ON GIT COMMANDS: Subagents never run git commands (git add, commit, push, status, checkout, etc.).\n- Code & Symbol Search: Use GitMap exclusively (`gitmap lf <dir>`, `gitmap find \"*.*\" <dir>`, `gitmap aum search \"<pattern>\" [dir]`, `gitmap cat <filepath>`). Banned ALL variants of Get-ChildItem, dir, ls in shell, rg, ripgrep, grep, git grep, Select-String, and findstr.\n- Cognitive Heatmap Telemetry: Log every read/edit via `gitmap heatmap log --task <slug> --action <READ|EDIT> --file <path> --state EXECUTION --expected \"<goal>\"`.\n- Targeted & Local Verification: Verify modified files with targeted linters and run local tests for your module to ensure zero regressions.\n- Report completion and evidence via `send_message` to parent."
    },
    {
      "TypeName": "self",
      "Role": "Worker 02: Secondary Components & Tests",
      "Model": "inherit",
      "Workspace": "inherit",
      "Prompt": "You are Worker 02 for task <nn>-<slug>. You have no prior chat context; this brief is your complete specification.\n\n### Owned Files (Exclusive Modification Rights):\n- `<relative/path/to/module2-file1.ext>`\n- `<relative/path/to/module2-file2.ext>`\n\n### Assigned Subtasks:\n- Subtask: `.ai-memory/plans/subtasks/<slug>/02-subtask-*.md`\n\n### Instructions & Invariants:\n- Follow the subtask plan and spec strictly. Implement required changes.\n- Strict Relative Git Paths ONLY. TOTAL BAN on absolute filesystem paths and file:/// URIs.\n- Strict lowercase filenames across all files.\n- TOTAL BAN ON GIT COMMANDS: Subagents never run git commands (git add, commit, push, status, checkout, etc.).\n- Code & Symbol Search: Use GitMap exclusively (`gitmap lf <dir>`, `gitmap find \"*.*\" <dir>`, `gitmap aum search \"<pattern>\" [dir]`, `gitmap cat <filepath>`). Banned ALL variants of Get-ChildItem, dir, ls in shell, rg, ripgrep, grep, git grep, Select-String, and findstr.\n- Cognitive Heatmap Telemetry: Log every read/edit via `gitmap heatmap log --task <slug> --action <READ|EDIT> --file <path> --state EXECUTION --expected \"<goal>\"`.\n- Targeted & Local Verification: Verify modified files with targeted linters and run local tests for your module to ensure zero regressions.\n- Report completion and evidence via `send_message` to parent."
    },
    {
      "TypeName": "self",
      "Role": "Worker 03: Validation, E2E Test Suites & Support Verification",
      "Model": "inherit",
      "Workspace": "inherit",
      "Prompt": "You are Worker 03 for task <nn>-<slug>. You have no prior chat context; this brief is your complete specification.\n\n### Owned Files (Exclusive Modification Rights):\n- `<relative/path/to/module3-file1.ext>`\n- `<relative/path/to/module3-file2.ext>`\n\n### Assigned Subtasks:\n- Subtask: `.ai-memory/plans/subtasks/<slug>/03-subtask-*.md`\n\n### Instructions & Invariants:\n- Follow the subtask plan and spec strictly. Implement required changes and test additions.\n- Strict Relative Git Paths ONLY. TOTAL BAN on absolute filesystem paths and file:/// URIs.\n- Strict lowercase filenames across all files.\n- TOTAL BAN ON GIT COMMANDS: Subagents never run git commands (git add, commit, push, status, checkout, etc.).\n- Code & Symbol Search: Use GitMap exclusively (`gitmap lf <dir>`, `gitmap find \"*.*\" <dir>`, `gitmap aum search \"<pattern>\" [dir]`, `gitmap cat <filepath>`). Banned ALL variants of Get-ChildItem, dir, ls in shell, rg, ripgrep, grep, git grep, Select-String, and findstr.\n- Cognitive Heatmap Telemetry: Log every read/edit via `gitmap heatmap log --task <slug> --action <READ|EDIT> --file <path> --state EXECUTION --expected \"<goal>\"`.\n- Targeted Verification: Run targeted linters and companion tests.\n- Report completion and evidence via `send_message` to parent."
    }
  ]
}
```

### 7.2 Self-Contained Worker Brief (Eliminate Context Blindness)

Subagents spawn with clean context. The prompt envelope MUST inject complete instructions:

```text
You are Worker <NN> for task nn-<slug>. You have no prior chat context; this brief is your complete specification.

### Boundaries & Crash Prevention:
- Read any file in the workspace; edit ONLY your Owned Files: <relative paths>.
- TOTAL BAN ON GIT COMMANDS (LOCK COLLISION PREVENTION): NEVER run ANY git commands (`git add`, `git commit`, `git push`, `git status`, `git diff`, `git checkout`). In shared workspaces, worker git calls create `.git/index.lock` collisions that immediately crash parallel agents. Only the lead orchestrator runs git commands after workers complete.
- TOTAL BAN ON COMMITS: Workers NEVER commit, stage, or push. Committing is exclusively reserved for the Lead Agent at Phase 3 via GitMap (`gitmap cpr` or `gitmap cpf` using hyphen `-`; no colon needed in GitMap commit arguments).
- Code & Symbol Search: Use GitMap exclusively:
  * Directory file inventory: `gitmap lf <dir>` (e.g. `gitmap lf src-tauri/src/modules/repo_db`) or `gitmap find "*.*" <dir>`
  * Live streaming code/symbol search: `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r] [-i]`
  * Finding files: `gitmap find "<pattern>" [dir]`
  * Stream file: `gitmap cat <filepath>`
  TOTAL BAN: Banned ALL variants of `Get-ChildItem` (with or without flags/parameters), `dir`, `ls` in shell, `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, and `findstr`.
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
6. Local Verification: Run module-level tests and targeted linters for modified code to confirm zero regressions.
7. Targeted Verification: Run only fast file-scoped linters (e.g. `python 03-ai-scripts/05-guideline-autofixer.py <folder> --check-only`). A check scanning 0 files is a FAIL.
8. GitMap Search Primacy: Banned ALL variants of `Get-ChildItem` (with or without flags/parameters), `dir`, `ls` in shell, `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, and `findstr`. Always use `gitmap lf <dir>`, `gitmap find "*.*" <dir>`, `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r]`, or `gitmap cat <filepath>`.
9. Mandatory Heatmap Telemetry: In multi-agent autonomous sessions, agents MUST record all file read/edit decisions with reasoning via gitmap heatmap log or gitmap agent log, and inspect active file leases and contention hotspots via gitmap agent heatmap (or gitmap heatmap status) before editing shared files.

### Concurrency-Safe SQLite Action Logging & Heatmap Telemetry (CRASH FORENSICS MANDATE):
- Worker subtasks are tracked in the run database: `<databasePath>`.
- Inspect Database Schema & Constraints:
  `gitmap agent task status`
  Inspect tables, positive boolean fields, and JSON payload specifications before claiming or logging actions.
- Claim assigned subtask atomically:
  `gitmap agent subtask claim --agent "Worker <NN>"`

#### 1. Heatmap File Contention & Pre-Edit Blast Radius Check:
- BEFORE modifying any shared files, you MUST inspect active file leases and collision contention hotspots:
  `gitmap agent heatmap [--json]` (or `gitmap heatmap status [--json]`)
  If high write contention is reported on a file, coordinate disjoint boundaries or sequence edits to prevent file lock crashes.

#### 2. In-Flight Cognitive Decision Logging (Heatmap & Action Log):
- BEFORE touching, inspecting, or modifying any file, you MUST log your in-flight action and reasoning for GitMap heat map tracking and crash forensics:
  * Heatmap File Touch & Exploration (reasoning):
    `gitmap heatmap log --task <task-slug> --action READ --file "<path>" --state EXPLORATION --expected "<what information is sought>"`
    (or via agent engine: `gitmap agent log --subtask <id> --agent "Worker <NN>" --action "read_file" --file "<path>" --details "Reasoning: <why this file is inspected>"`)
  * Heatmap File Modification (editing decision & rationale):
    `gitmap heatmap log --task <task-slug> --action EDIT --file "<path>" --state EXECUTION --expected "<planned refactor/fix>"`
    (or via agent engine: `gitmap agent log --subtask <id> --agent "Worker <NN>" --action "write_to_file" --file "<path>" --details "Decision: <why this file is being modified>"`)
  * Heatmap Outcome Update (after subagent completion report):
    `gitmap heatmap update --id <id> --outcome "<outcome>" --met true`
  *(Note: This guarantees that GitMap indexes full cognitive heat map telemetry into `.ai-memory/heatmap.db` and agent task DB, and if a tool execution crashes or the session restarts, the database permanently records the exact file and reasoning trace).*

#### 3. Real-Time Telemetry, Roadmap Linkage & Heatmap Stream:
- **Roadmap-Connected Heatmap for Dynamic Join & Continuation:**
  * When subagents spawn or dynamically join (`gitmap agent subtask join`), the agent inspects `gitmap heatmap status` (or `gitmap rc heatmap status`) and `gitmap agent heatmap` to immediately discover:
    1. Where previous agents have read files and why (`READ` events + rationale).
    2. What files were modified and what decisions were made (`EDIT` events + rationale).
    3. Which files currently hold active leases (preventing write collisions).
    4. What subtasks remain in the roadmap (`.ai-memory/plans/<slug>.md`).
  * This eliminates duplicate reading steps, allowing newly joined subagents to immediately continue from where previous agents left off and accelerate execution.
- Inspect live heat map, tree, and action status:
  * Hierarchical tree view linked to roadmap: `gitmap heatmap tree [-t <task-slug>]` / `gitmap rc heatmap tree`
  * Action summary and metrics: `gitmap heatmap status [-t <task-slug>]` / `gitmap rc heatmap status`
  * Real-time continuous stream: `gitmap heatmap live` / `gitmap rc heatmap live`
  * Action search: `gitmap heatmap search [query] [-f <file>] [-t <task>]`
  * 48-hour active repository activity: `gitmap full-summary` (alias `gitmap fs`)

#### 4. Multi-AI Companion & repo-cache Synchronization (RC Heatmap Suite):
- **Automatic Vault Creation:** If `repo-cache` (`rc`) or `.ai-memory/heatmap.db` does not exist, GitMap auto-provisions the required split-db structure linked to the active mission and folder.
- Lists all supported AI namespaces: `gitmap heatmap rc list-ai`
- Logs cognitive actions into isolated split-databases: `gitmap heatmap rc log [flags]`
- Checkpoint commit in `repo-cache`: `gitmap heatmap rc commit [-m "<msg>"]`
- Bidirectional reconciliation: `gitmap heatmap rc sync`
- Tree visualization in RC vault: `gitmap heatmap rc tree` / `gitmap rc heatmap tree`
- Live stream in RC vault: `gitmap heatmap rc live` / `gitmap rc heatmap live`

#### 5. Non-Blocking Read-Only Heatmap Queries & Diagnostics:
- Execute read-only SQL queries with automatic `mode=ro` protection:
  `gitmap heatmap query "<sql>" [--format json|yaml|table]`
  `gitmap heatmap rc query "<sql>" [--ai <name>] [--repo <slug>] [--tier <repo|root>]`
- Inspect SQLite master catalog, tables, and row counts:
  `gitmap heatmap tables [--db <root|repo|local>]` / `gitmap heatmap rc tables`
- Query records from specific heatmap table:
  `gitmap heatmap view <table> [--limit <n>]` / `gitmap heatmap rc view <table>`

- When your subtask passes targeted checks and tests, mark completion in the database:
  `gitmap agent subtask complete <id> --agent "Worker <NN>" --evidence "PASS exit 0, <files>"`
- If blocked or failing, record the failure:
  `gitmap agent subtask fail <id> --agent "Worker <NN>" --reason "<reason>"`

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
   - Lead inspects overall completion status at any time:
     `gitmap task status --db <databasePath>`
     *(Legacy fallback: `gitmap agent task status`)*
   - If any worker fails to report, crashes, or times out, inspect status or diagnose:
     `gitmap agent task status` or `gitmap agent diagnose`
     This pinpoints the autopsy: which agent crashed, on which subtask, targeting which file, and the exact action that was executing when it failed.
3. **Verify Worker Reports Independently:** Confirm `git diff --stat -- <owned files>` matches `filesChanged`, no files outside owned files modified, re-run targeted checks for `exit 0` on non-zero files.
4. **Reject Violations:** Send failures via `send_message`. On `BLOCKED`, lead does work and logs `LEAD_FALLBACK: <reason>`. After two failed rounds, mark `FAILED`, write RCA, continue (R13).
5. **Update Heatmap on Subagent Report:** After each subagent reports completion or findings via `send_message`, update the corresponding heatmap entry:
   `gitmap heatmap update --id <id> --outcome "<outcome>" --met true`
6. **Update Ledger:** Record status, evidence, changed paths in `ledger.md` via `replace_file_content`.
7. **Loop:** Dispatch subsequent waves via `invoke_subagent` until all subtasks are `DONE` or `FAILED`.

### 7.4 Immediate Dynamic Respawn & Join Protocol (100% Duty Cycle Concurrency)

To maximize execution velocity and prevent idle worker slots during Phase 2 and Phase 3:
1. **Zero-Idle Immediate Respawn Mandate:** The lead orchestrator MUST NOT wait for an entire 3-agent wave to finish. As soon as ANY single worker completes (`send_message` or `gitmap agent subtask complete`), the lead immediately checks the queue for the next pending subtask and spawns a replacement worker (`invoke_subagent`), keeping concurrency pinned at exactly `A = 3` active subagents (System Concurrency Capacity = 6).
2. **Dynamic Collaborative Assistant Mode:** If all pending subtasks are already claimed or in-flight, any newly freed agent MUST NOT terminate, halt, or idle. The lead immediately commands or spawns the worker to dynamically join an active subtask in collaborative assistance mode:
   - **Mode A: Parallel Reading & Blast Radius Discovery:** The joining agent uses GitMap (`gitmap aum search`, `gitmap cat`, `gitmap find`) to read caller files, inspect imports, and verify type signatures, feeding findings to the primary worker via structured message.
   - **Mode B: Cognitive Memory & Spec Documentation:** The joining agent records architectural discoveries, updates memory registers, and generates companion audit specs.
   - **Mode C: Companion Lint & Test Suite Verification:** The joining agent executes targeted linters (`check-relative-paths.py`, `check-forbidden-strings.py`, guideline autofixer check) and runs isolated test packages across related workspaces while the primary worker completes editing.
   - **Mode D: Disjoint File Partitioning:** If the ongoing subtask touches multiple files, the subtask is subdivided dynamically so the joining agent takes exclusive ownership of File 2 while the primary worker focuses on File 1.
3. **Disjoint Ownership Invariant (Rule R6):** Dynamic assistants NEVER write to the primary worker's owned files. File ownership remains strictly disjoint to eliminate file write collisions.
4. **GitMap Telemetry Tracking:**
   - Log join actions: `gitmap agent subtask join --subtask <id> --agent "Worker <NN>" --role "Assistant: <focus>"`
   - Inspect tree hierarchy: `gitmap agent task tree` (renders live Root -> Branch -> Leaf/Assistant hierarchy).

---

## 8. Phase 3: Consolidation, E2E Testing, Audit Log Verification, Minor Release & Atomic Push

1. **Consolidate Subtasks:** Merge completed subtasks into `.ai-memory/plans/completed/nn-<slug>.md`, logging real steps from ledger; link to canonical spec. Delete `.ai-memory/plans/subtasks/nn-<slug>/` and pending plan. Canonical spec in `02-spec/21-app/` stays permanently.
2. **Comprehensive Local Test Suite & E2E Execution (Mandatory 100% Green Pass):**
   - Run ALL project test suites locally to ensure zero regressions across the codebase:
     * For Go: `go test ./...` (must exit 0 with all packages passing).
     * For Node/TypeScript: `npm test` or `pnpm test` or `npm run test:e2e` (must exit 0 with 100% tests green).
     * For Python: `pytest` (must exit 0 with all test cases passing).
   - Run all specific E2E scenarios and test fixtures inventoried during Turn 1.
   - If any test fails: stop immediately, perform 4-part RCA, apply surgical fix via subagents, and re-run test suites until 100% green exit code 0 is achieved.
3. **Audit Log Inspection & Heatmap Verification:**
   - Inspect `.ai-memory/heatmap.db` and SQLite task engine:
     * Run `gitmap agent log` or `gitmap task status` to review in-flight agent action logs.
     * Run `gitmap heatmap tree --task <task>` and `gitmap heatmap status --task <task>` to verify 100% action closure and ensure that all planned operations met their expected outcomes.
     * Confirm zero unaccounted touches or unclosed file modifications across the session.
4. **Automated Minor Version Bump Release Ceremony:**
   - As soon as all local test suites pass 100% green and audit logs are confirmed, execute the minor version bump release ceremony:
     * Run minor bump command: `python 03-ai-scripts/37-bump-version.py --tier minor` (or `npm run bump minor` where package.json scripts define it).
     * Verify updated `version.json` (or `package.json`).
     * Generate / update changelog entry with release notes under `02-spec/24-releases/` or root `changelog.md`.
     * Update root `readme.md` version pin to match the newly bumped minor SemVer.
5. **Update Registers:**
   - Update `.ai-memory/plans/readme.md` and `02-spec/21-app/readme.md`.
   - Update `01-prompts/readme.md` and `.ai-memory/prompts.md` only when adding/modifying prompts.
   - Register recent completed tasks before push gate.
6. **Push Gate Verification:**
   - Before GitMap call, verify:
     (a) All local test suites (`go test ./...`, `npm test`, `pytest`) exit 0 with 100% green pass,
     (b) Targeted linters exit 0 (>0 files),
     (c) Secrets gate clean,
     (d) `.gitignore` covers caches, build outputs, logs, reports, `.env*` (untrack any tracked ignored files via `git rm --cached <file>` or `git rm -r --cached <dir>`),
     (e) Phase 3 cognitive heatmap verification: run `gitmap heatmap tree --task <task>` and `gitmap heatmap status --task <task>` to verify all actions and expected outcomes are 100% met before committing.
   - On failure: abort GitMap call; mark task `FAILED` with RCA.
7. **Atomic Release Commit & Push:**
   - Call `gitmap cpr "<module> - <release summary>"` (e.g. `gitmap cpr "E2E - minor release bump with 100 percent green test suite"`) or `gitmap cpf "<module> - <feature summary>"`.
   - Mention as a hyphen `-` because the colon is already going to be provided by GitMap (which automatically prepends `Release: ` or `Feature: `). So, there is no need to provide the colon in GitMap `cpr`/`cpf`/CVF message arguments (TOTAL BAN on colons `:` inside the message argument); always use a hyphen `-` to separate module from summary. NEVER include redundant conventional prefixes (`feat(...)`, `fix(...)`, `docs(...)`, `chore(...)`) in the commit string.
   - Push rejected: `git pull --rebase` and re-run.
   - Miss after push: allow one follow-up `gitmap cpb "<module> - <fix summary>"`, logged as `FOLLOW_UP_PUSH: <sha>`. Never amend pushed commits.

---

## 9. Final Report Format (Strict Vertical Lines)

```markdown
### Task Completion Summary

- ✅ **Task-01: [Descriptive Task Title]** — `[Completed]` — [diff/check evidence]
- ✅ **Task-02: [Descriptive Task Title]** — `[Completed]` — [diff/check evidence]

### Test Suites & E2E Verification Summary

- ✅ **Local Test Suite (go test / npm test / pytest):** 100% Green (`exit 0`, [N] tests passed, 0 failed)
- ✅ **E2E Scenarios & Fixtures:** 100% Green ([N] test cases passed)
- ✅ **Audit Log Review:** Confirmed via `gitmap agent log` and `gitmap heatmap tree` (100% action closure)

### Automated Release Summary

- **Release Tier:** Minor Version Bump
- **Bumped Version:** `vX.Y.0`
- **Changelog:** Updated and recorded
- **Root README Pin:** Synchronized to `vX.Y.0`
- **Release Commit:** `gitmap cpr "<module> - <summary>"` pushed upstream

### Modified Files Summary

- [relative/path/to/modified/file1.ext]
- [relative/path/to/modified/file2.ext]

### Steps Used

- Step x / N (Phase 1: y / PHASE_1_BUDGET, Phase 2: z / PHASE_2_BUDGET), Wave k / WAVES

### Implementation Confidence Score

- Confidence: 1.0 (100% green test pass rate, verified audit logs, zero regressions)
- Rationale: [Verified evidence across all criteria, passing local test suites, passing targeted linters, minor release pushed]

Independent check: run /verify-parent-task-run <slug>

### 🤖 Independent AI Verification & Audit Prompt

(Emit self-contained audit prompt linking spec, plan, test results, and modified files)
```

---

## 10. Targeted Verification Checks (R2) & Full Test Suites

Confirm scripts exist via harmless workspace call before invoking (R4). Run on changed files/folders AND run full project suites:

- **Local Test Suites (Mandatory 100% Green):**
  * `go test ./...`
  * `npm test` / `pnpm test` / `npm run test:e2e`
  * `pytest`
- **Coding Guidelines & Boolean Linter:** `python 03-ai-scripts/05-guideline-autofixer.py <folder> --check-only --ext <.ext>`
- **Relative Path Linter:** `python linter-scripts/check-relative-paths.py`
- **Prompts & Spec Index Linter:** `python linter-scripts/check-prompts-loaded.py`
- **Markdown Link & Doc Path Linter:** `python 03-ai-scripts/22-doc-path-linter.py <folder>`
- **Sequence Integrity Linter:** `python linter-scripts/check-sequence-integrity.py`
- **Forbidden Strings Check:** `python linter-scripts/check-forbidden-strings.py`
- **Minor Version Bump Script:** `python 03-ai-scripts/37-bump-version.py --tier minor`

---

## 11. AI Fix Scripts Memory (Reusable Tooling)

- [ ] [/goal](slashCommand;goal) Reuse First: Scanned and learned `03-ai-scripts/readme.md` before writing temporary code.
- [ ] Strict In-Repository Execution: All Python scripts executed strictly within the codebase repository root.
- [ ] Strict .ai-memory/ Folder Storage: All helper scripts, local runners, and linters stored in `03-ai-scripts/`.
- [ ] Native File Manipulator: Use `python 03-ai-scripts/03-file-manipulator.py <command>` for mass file operations.
- [ ] Minor Version Bump Script: Use `python 03-ai-scripts/37-bump-version.py --tier minor` or `npm run bump minor` for release finalization.
- [ ] Go Generate Sync: If Go constants or enums are modified, run `go generate ./...` in the relevant package and commit generated files.

---

## 12. Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] TOP-INSTRUCTION PRIORITY MANDATE: Whatever directives, constraints, checklists, or instructions are given before this section or prompt (user preamble, header constraints, prior instructions) are verified as highest priority and non-negotiable, overriding all lower-level guidelines below.
- [ ] TURN 1 RECITATION MANDATE: Verified that Turn 1 chat output displayed the 5 Mandatory Pillars, the output of `git log -n 5 --oneline` (recent 5 commits), and the complete inventory of discovered test cases and test suites before background tool calls.
- [ ] MANDATORY SUBAGENT SPAWNING GATE (ZERO SOLO EXECUTION): Verified that `invoke_subagent` (`A = 3, H = 2, Capacity = 6`) was explicitly called in Phase 1 (parallel discovery/spec modules) and Phase 2 (`TypeName: "self"` parallel subtask execution). Executing all reads or code modifications solo without calling `invoke_subagent` is an immediate auto-reject failure.
- [ ] TREE SWARM & DYNAMIC JOIN MANDATE (ZERO IDLE CAPACITY): Verified that execution followed the Tree Swarm Hierarchy (Root Orchestrator -> Branch Coordinators -> Leaf Workers & Dynamic Assistants), that Phase 1 authored the comprehensive Brain roadmap up front before worker dispatch, and that newly idle workers immediately respawned for pending tasks or dynamically joined in-flight tasks for parallel reading, caller analysis, test authoring, or companion verification to keep concurrency capacity pinned at 3 active subagents (Capacity = 6) until completion.
- [ ] NO SUBAGENT FILE WRITE COLLISIONS (TOTAL BAN): Never assign two subagents to create or modify the same file. Shared files belong exclusively to the lead orchestrator. Subagents must have strictly disjoint, non-overlapping file boundaries to prevent file lock crashes.
- [ ] NO READ-ONLY SUBAGENT FILE WRITES (TOTAL BAN): Subagents of `TypeName: "research"` have read-only tools and cannot write files. Never command a research subagent to write files; research subagents discover and report back, and the lead writes the unified plan.
- [ ] NO INVALID SUBAGENT TYPENAMES (TOTAL BAN): Never pass custom strings like "Research 01" or "Worker 01" in the `TypeName` field of `invoke_subagent`. `TypeName` MUST be strictly `"research"` (read-only discovery) or `"self"` (read/write execution). Put descriptive names in `Role` only.
- [ ] NO SUBAGENT GIT EXECUTION (TOTAL BAN): Subagents must NEVER run any git commands (`git add`, `git commit`, `git push`, `git status`, `git diff`, `git checkout`) in shared workspaces. Worker git execution creates `.git/index.lock` collisions that immediately crash parallel workers with exit code 128.
- [ ] MANDATORY LOCAL TEST SUITES EXECUTION (100% GREEN REQUIREMENT): Unlike standard V6, executing without running all project local test suites (`go test ./...`, `npm test`, `pytest`) and passing 100% green exit code 0 is STRICTLY FORBIDDEN. Skipping tests or proceeding with failing tests is an immediate auto-reject failure.
- [ ] AUDIT LOG REVIEW MANDATE: Verified `.ai-memory/heatmap.db` and SQLite task engine logs (`gitmap agent log`, `gitmap heatmap tree --task <task>`) confirming 100% action closure and zero unaccounted touches.
- [ ] AUTOMATED MINOR RELEASE CEREMONY MANDATE: Verified minor version bump ceremony (`npm run bump minor` or `python 03-ai-scripts/37-bump-version.py --tier minor`), changelog update, and root readme version pin executed as soon as all tests passed 100% green.
- [ ] NO RAW GIT COMMITS OR CONVENTIONAL COMMIT PREFIXES (TOTAL BAN): Never execute `git commit`, `git commit -m "..."`, `git add -A`, `git add .`, or raw `git push`. Never use conventional commit prefixes (`docs(...):`, `feat(...):`, `fix(...):`, `chore(...):`) in raw git commands. All staging, committing, and pushing MUST be executed exclusively by GitMap: `gitmap cpr "<module> - <summary>"` (release), `gitmap cpf "<module> - <summary>"` (features), or `gitmap cpb "<module> - <summary>"` (fixes). GitMap automatically formats, stages, commits, and pushes atomically. Always mention and format as a hyphen (`-`); no need to provide a colon in GitMap commit commands (`gitmap cpr` / `cpf` / `cpb` / CVF) because the colon is already automatically provided by GitMap.
- [ ] NO INTERMEDIATE COMMITS (TOTAL BAN): Never commit after Phase 1 (e.g. committing plans or specs) or mid-Phase 2 (committing individual files or tests). Committing early pollutes git history, creates race conditions, and breaks atomicity. All changes (specs, plans, code modifications, tests, index updates, release bump) MUST be committed together in ONE single atomic GitMap commit at the final step of Phase 3.
- [ ] NO PER-FILE COMMITTING (TOTAL BAN): Never commit each file individually as you work. All modified files across the turn must be accumulated and committed together in a single atomic release commit at the final step.
- [ ] NO UNCOMMITTED WORK OR UNPUSHED CODE (TOTAL BAN — NO PUSH = NOT DONE): Never conclude a turn, claim success, or mark any task as DONE while leaving changes uncommitted or unpushed. If the code is not committed to Git and pushed upstream to GitHub (main/master/tracking branch), the task is strictly considered INCOMPLETE and NOT DONE. Leaving uncommitted dirty changes, untracked files, or unpushed local commits means the execution is unfinished and has failed. Concluding without a confirmed successful GitMap push to remote is an immediate auto-reject failure.
- [ ] NO HORIZONTAL TASK CONCATENATION (TOTAL BAN): Never concatenate tasks horizontally in the Task Completion Summary. Every completed task MUST be rendered on its OWN SEPARATE LINE starting with an individual markdown list bullet (`- ✅`).
- [ ] INDEPENDENT AI VERIFICATION PROMPT MANDATE: Emitted the self-contained independent AI verification and audit prompt linking to the canonical spec, consolidated plan, and modified files with verbatim score audit criteria.
- [ ] GITMAP HEAVY USAGE & ROUTINE PULL BAN: Heavily leveraged GitMap commands (`cpr`, `cpf`, `cpb`, `search`, `find`, `pwsh`) for discovery, execution, and commits. Never ran `pull-all` (`gitmap pa` or `gitmap pae`) unconditionally during routine turns; only ran `gitmap pae --json` when explicitly commanded by the user.
- [ ] NO RG, RIPGREP, OR RAW SHELL SEARCHES (TOTAL BAN): Banned ALL variants of `Get-ChildItem` (with or without flags/parameters), `dir`, `ls` in shell, `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, and `findstr`. All directory inventories, code searching, symbol discovery, and regex scans MUST use GitMap high-speed tools: `gitmap lf <dir>` / `gitmap find "*.*" <dir>` (directory file inventory), `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r] [-i]` (live streaming code/symbol search), `gitmap find "<pattern>" [dir]` (finding files), and `gitmap cat <filepath>` (stream file). Running `rg`, `ripgrep`, `Select-String`, `git grep`, `Get-ChildItem`, `dir`, or `ls` is an immediate auto-reject failure.
- [ ] HEATMAP TELEMETRY & CONTENTION INSPECTION: Verified that Step 0 preflight checks (`gitmap heatmap status` and `gitmap heatmap tables`) were executed, that `gitmap heatmap log --task <task> --action <READ|EDIT|CREATE|DELETE> --file <path> --state <STATE> --expected "<goal>"` was executed before file reads/edits, that `gitmap heatmap update --id <id> --outcome "<outcome>" --met true` was logged after subagent reports, and that Phase 3 verification executed `gitmap heatmap tree --task <task>` and `gitmap heatmap status --task <task>`.

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
- [ ] Local Test Suites 100% Green: Verified that `go test ./...`, `npm test`, or `pytest` passed completely with exit code 0 before initiating release ceremony.
- [ ] Final Step Release Commit & Push Verified: Staged and committed all changes atomically via GitMap semantic release commit commands using `<module> - <summary>` format: `gitmap cpr "<module> - <summary>"` (release) or `gitmap cpf "<module> - <summary>"` (features). Mentioned and formatted with a hyphen `-` to separate module from summary; no need to provide a colon in GitMap `cpr`/`cpf`/CVF commit arguments because the colon is already automatically provided by GitMap. TOTAL BAN on colons `:` inside the GitMap argument and conventional prefixes (`fix(...)`, `feat(...)`, `docs(...)`). TOTAL BAN on raw `git add -A` and `git commit`. Pushed to remote via GitMap in a single final command.
- [ ] Strict Completion Invariant ('No Push = Not Done'): Confirmed that all modifications are staged, committed atomically via GitMap, and pushed upstream to GitHub. If the code is not committed to Git and pushed upstream to GitHub (main/master/tracking branch), the task is strictly considered INCOMPLETE and NOT DONE. Verified that `git status --porcelain` is completely clean (zero uncommitted or untracked files) and that the remote tracking branch reflects the pushed commit before declaring task completion.

---

## 15. Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] MANDATORY FINAL COMMIT & PUSH VIA GITMAP (ANYHOW): At the final step of the turn, after all targeted files have been refactored, verified with targeted linters, verified 100% green across local test suites, minor release bumped, and plans/subtasks consolidated, use GitMap semantic commit commands exclusively: `gitmap cpr "<module> - <summary>"` (release) or `gitmap cpf "<module> - <summary>"` (features), which automatically stage, format commit messages, and push directly to remote. Mention and format as a hyphen `-` to separate module and summary; no need to provide a colon `:` in the GitMap `cpr`/`cpf`/CVF message argument as the colon is already provided automatically by GitMap (`Release: ` or `Feature: `). TOTAL BAN on raw `git commit`, `git commit -m`, `git add -A`, colons `:` in the message argument, or conventional prefixes inside GitMap arguments (`docs(...)`, `feat(...)`, `fix(...)`). ZERO intermediate commits during Phase 1 or Phase 2; all files across the run are committed together at the final step. Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY): You must not create separate git commits for each individual file as you edit them. Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback. All modified files, test change caches, and plan records across the turn must be accumulated in the working tree and committed together in a single grouped atomic release commit at the final step before pushing.
- [ ] MANDATORY COMPLETION INVARIANT (NO PUSH = NOT DONE): If the code is not committed to Git and pushed upstream to GitHub (main/master/tracking branch), the task is strictly considered INCOMPLETE and NOT DONE. Leaving uncommitted dirty changes or unpushed commits means the execution is unfinished and has failed. A task cannot be marked completed or successful until all changes are committed and confirmed pushed to GitHub. The orchestrator must verify that `git status --porcelain` returns completely empty and that the remote tracking branch is up to date before rendering the final completion summary.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase readme.md files left uncorrected, `.ai-memory/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `02-spec/` and `.ai-memory/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run all local test suites to 100% green pass, inspect audit logs, bump minor version, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.
