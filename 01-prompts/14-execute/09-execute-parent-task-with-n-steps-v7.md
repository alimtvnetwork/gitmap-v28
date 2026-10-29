# [V7] Parent Task N-Step Continuous Loop & Autonomous Multi-Agent Swarm Orchestration — Workflow (must follow)

```text
N = 400 (Total self-loop steps budget — editable top-header parameter, default: 400)
A = 3   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 3)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 3 agents × 2 hands = 6 concurrent subtask operations
Tree Hierarchy = Lead Orchestrator (Root Node) -> Branch Coordinators (Subtask Owners) -> Leaf Workers & Dynamic Joining Assistants
PHASE_1_BUDGET = 100     (Steps 1 .. 100: Planning, Discovery Subagents, Up-Front Brain & Roadmap Creation, Lean Subtasks, On-Call Skill Check)
PHASE_2_BUDGET = 150     (Steps 101 .. 250: Intermediate Archival, Sanitization, Prerequisite Refactoring, Intermediate Checkpoints)
PHASE_3_BUDGET = 150     (Steps 251 .. 400: Tree Swarm Parallel Execution, Dynamic Join Capacity Saturation, Recommit Gate, Time-Restriction Validation, Atomic Push)
WAVES = ceil(subtasks / (A x H))
```

> [!IMPORTANT]
> Prompt Version: 7.0.0
> Runtime: Google Antigravity 2.0 (IDE and CLI)
> Invoke: /execute-parent-task-with-n-steps-v7 <task>
>
> **Top-Instruction Priority Mandate (Above Precedence / Preamble Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given ABOVE this prompt (including in the user preamble, header blocks, or incoming user request above) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below. The agent MUST inspect and follow the instructions above with absolute precedence.

[/goal](slashCommand;goal) Autonomously orchestrate and execute the parent task end-to-end: FIRST showcase and list out the given task in visible chat during Turn 1 with the 5 Mandatory Pillars and 6-part verification schema, capture it verbatim, plan it in the repo, evaluate On-Call Skill Gates (`gitmap-time-restriction`, `gitmap-automation-bot`), spawn autonomous subagents via `invoke_subagent` (A = 3, H = 2, Capacity = 6; solo execution without calling `invoke_subagent` is an auto-reject failure) in disjoint file boxes using GitMap high-speed commands as primary, prove every single claim with concrete evidence, enforce coding guidelines to 100%, and finish with one atomic GitMap commit (`gitmap cpf "<module> - <summary>"` / `gitmap cpb "<module> - <summary>"` / `gitmap cpc "<module> - <summary>"` / `gitmap cpr "<module> - <summary>"`) using hyphen format with Title Case or lowercase subsystem module and NO third brackets (square brackets) or tautological 'commit' module that holds strictly this task's files.

[/learn](slashCommand;learn) Enforce the Top-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are provided ABOVE this prompt outrank everything below. Turn 1 MUST showcase the given task list in visible chat before any background execution. Each rule is stated once (R1 to R16) and cited by ID. Progress lives in the ledger and in `.ai-memory/plans/`, never only in chat.

[/plan](slashCommand;plan) Execute thorough step-by-step planning in the repository before execution for improved work or less friction, faster work and more pure and effective work. Ensure all deliverables, architecture boundaries, and requirements are clearly defined in the audit ledger and subtask plans before dispatching worker waves.

### 🚨 MANDATORY SUBAGENT SPAWNING GATE (A = 3, H = 2, CAPACITY = 6 — TREE SWARM & DYNAMIC JOIN)

- **ACTUAL TOOL CALL REQUIRED:** You must ACTUALLY CALL the `invoke_subagent` tool via your tool-calling API. Do NOT just print the text "Dispatched Worker..." and stop. If you only print text, the agents will not spawn and the task will fail! You must execute the `invoke_subagent` JSON tool payload.
- **TREE-STRUCTURED SWARM HIERARCHY:** Organize execution as an operational tree:
  1. **Root Orchestrator (Lead Agent):** Owns the unified Brain roadmap, audit ledger, and final GitMap commit. Never executes file edits solo.
  2. **Branch Coordinators (Subtask Primary Owners):** Own designated disjoint file boundaries (`Owned Files`) and execute modifications.
  3. **Leaf Workers & Dynamic Joining Assistants:** Assist primary owners with parallel reading, symbol caller discovery, memory generation, and companion lint verification.
- **3-STAGE MANDATORY DISPATCH:** You must invoke `A = 3` agents (`invoke_subagent`) at EVERY stage of the workflow:
  1. **Planning Step (A = 3 `research` subagents):** Spawn 3 read-only discovery subagents (`TypeName: "research"`) to research the codebase and return findings; lead writes the unified Brain execution plan and roadmap.
  2. **Spec Step (A = 3 self authoring subagents — ZERO SOLO SPEC AUTHORING):** Spawn 3 authoring subagents (TypeName: "self") to author modular, disjoint spec files and subtasks (never writing the same file). Solo authoring by the lead orchestrator is STRICTLY FORBIDDEN.
  3. **Execution Step (A = 3 `self` worker subagents):** Spawn 3 worker subagents (`TypeName: "self"`) to execute code modifications in strictly disjoint file boxes.
- **DYNAMIC AGENT JOIN & IMMEDIATE RESPAWN (NO IDLE WORKERS):** Worker capacity `A` must remain pinned at 3 active agents (System Concurrency Capacity = 6) throughout execution. As soon as ANY single worker completes, the lead immediately spawns a replacement worker for the next pending subtask, or commands the worker to dynamically join an active subtask in collaborative assistance mode (parallel reading, companion linters, memory logging) to keep concurrency at exactly 3 active subagents (100% duty cycle).
- **SOLO EXECUTION IS AN AUTO-REJECT FAILURE:** The lead orchestrator is **STRICTLY FORBIDDEN** from executing planning, spec writing, or code changes by itself without calling the `invoke_subagent` tool.

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

1. **Top-Instruction Priority Verification:** Whatever directives, constraints, checklists, or instructions are given before this section or prompt must be verified as highest priority and non-negotiable.
2. **Extreme "Understood" Turn-1 Showcase (Turn 1 Action):** In your VERY FIRST response turn upon receiving the prompt, output a comprehensive understanding block directly in visible chat before any background tool execution, covering the 5 Mandatory Pillars:
   - **Pillar 1: Strategic Intent & Systems Reasoning**
   - **Pillar 2: Explicit Scope Boundaries (In-Scope vs. Excluded / Deferred)**
   - **Pillar 3: Non-Negotiable Constraints & Negative Guardrails Ingestion**
   - **Pillar 4: Blast Radius & Dependency Topology**
   - **Pillar 5: Multi-Agent Concurrency Topology & Wave Orchestration**
3. **Lossless Verbatim Capture:** Store incoming prompt losslessly under `## User Request (Verbatim)` in canonical spec and parent plan.
4. **Discrete Deliverables Extraction:** Break down user requirements into discrete, actionable items with ordered traceable IDs (`Task-01`, `Task-02`, ...).

### Phase 1B: On-Call Skill Resolution Gate (Step 0.5)

Before proceeding to research or spec generation, evaluate whether the incoming task triggers any **On-Call Skills**:

1. **Time Restriction Check (`gitmap-time-restriction`):**
   - If prompt specifies working hours, schedules, or commit timestamps, activate `.agents/skills/gitmap-time-restriction/skill.md`.
   - Set `isTimeRestricted: true` in task configuration.
   - Enforce the `DD-MMM-YY` UTC standard across all date handling.
2. **AI Automation Bot Check (`gitmap-automation-bot`):**
   - If prompt specifies bot author identity or headless automation commits, activate `.agents/skills/gitmap-automation-bot/skill.md`.
   - Execute pre-flight identity switch via `gitmap change-user "gitmap automation <local@local>"` or `gitmap autobot`.
   - Set `isAutobotActive: true` in task tracking.
3. **Recommit Verification Check:**
   - If task specifies amending, replaying, or recommitting previous changes, verify working tree cleanliness via `gitmap st` and check commit history before modifying files.

### Phase 1C: Planning & Spec Step (Steps 1–100)

1. Spawn 3 `research` subagents to perform discovery across disjoint areas.
2. Spawn 3 `self` authoring subagents to author modular specs in `02-spec/21-app/XX-<slug>/` and plans in `.ai-memory/plans/pending/XX-<slug>.md`.
3. Initialize the SQLite task database at `.ai-memory/temp-agents/<slug>/agent-task.db`.

### Phase 2: Prerequisite Setup & Disjoint Execution (Steps 101–250)

1. Workers execute changes strictly within assigned disjoint file boundaries (`Owned Files`).
2. Adhere strictly to the 12-Pillar Quality Gate:
   - Leaf error sovereignty: leaf functions return `error` or `Result[T]` directly. Zero `os.Exit` or `cliexit.HandleError` calls.
   - Parameter structs for >2-3 parameters.
   - Positive booleans (`is*`, `has*` only; no negative booleans).
   - Function limit: <= 15 lines max, <= 8 lines preferred.
   - Zero nested `if` statements (guard clauses only).
   - Strict relative Git paths only (`02-spec/...`, `.ai-memory/...`, `cmd/...`).

### Phase 3: Verification, Audit & Atomic GitMap Commit (Steps 251–400)

1. Run targeted, file-scoped checks on modified files.
2. If `isTimeRestricted` is active, verify commit timestamp window.
3. If `isAutobotActive` is active, verify author credentials before committing.
4. If recommit was commanded, execute history reconciliation.
5. Move plan from `.ai-memory/plans/pending/XX-<slug>.md` to `.ai-memory/plans/completed/XX-<slug>.md`.
6. Update `.ai-memory/plans/readme.md` and `02-spec/21-app/readme.md`.
7. Execute final atomic commit via GitMap semantic verbs:
   - `gitmap cpf "<module> - <feature summary>"`
   - `gitmap cpb "<module> - <fix summary>"`
   - `gitmap cpc "<module> - <chore summary>"`
   - `gitmap cpr "<module> - <release summary>"`

---

## 1. Core Invariants & Rules Reference

- **R1 Zero Builds, Test Suites & CI/CD Ban:** Routine turns must NEVER run full builds or test suites. Only run targeted checks.
- **R2 Targeted Checks Only:** Run only file-scoped checks on modified files.
- **R3 Evidence or It Did Not Happen:** Every claim must cite a file path, diffstat, or exit code.
- **R4 Never Invent Commands, Flags, or Paths:** Verify commands with harmless calls (`gitmap lf readme.md`).
- **R5 Mandatory Subagents:** Always call `invoke_subagent` (A=3, H=2, capacity=6). Solo execution is forbidden.
- **R6 One Owner Per File:** Disjoint file boundaries for workers.
- **R7 Worker Git Ban:** Workers never run git commands. Only lead executes git/GitMap operations.
- **R8/R9 Atomic Commit via GitMap:** Finish with one atomic semantic GitMap commit. TOTAL BAN on raw git commits, square brackets `[...]`, or colons in commit arguments.
- **R10 Zero Unauthorized Releases:** Never bump versions unless explicitly requested.
- **R11 Leaf Error Sovereignty:** Leaf functions return errors directly; never call exit handlers.
- **R11b Strict Relative Git Paths:** Absolute paths and `file:///` URIs are banned.
- **R12 No Polling:** Stop calling tools after launching subagents.
- **R13 Anti-Looping:** Cap retries at two strikes; report RCA if blocked.
- **R14 Ambiguity Protocol:** Conservative assumptions logged in ledger; prompt user once if blocked.
- **R15 Zero Artifacts Committed:** Never commit generated temporary scripts or caches.
- **R16 Zero Secrets:** Strict secrets gate before committing.
