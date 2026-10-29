---
name: execute-parent-task-with-n-steps-v7
description: Master autonomous execution engine with dual operating modes (Option 1: Time-Restricted with timezone & weekday < 8:00 AM cutoff; Option 2: Unrestricted), AI automation bot integration, Tree Swarm Hierarchy (A=3, H=2, Capacity=6), and GitMap atomic commits.
---

# [V7] Master Autonomous Execution Engine — Continuous Loop & Swarm Orchestration (must follow)

```text
N = 400 (Total self-loop steps budget — editable top-header parameter, default: 400)
A = 3   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 3)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 3 agents × 2 hands = 6 concurrent subtask operations
Tree Hierarchy = Lead Orchestrator (Root Node) -> Branch Coordinators (Subtask Owners) -> Leaf Workers & Dynamic Joining Assistants
PHASE_1_BUDGET = 100     (Steps 1 .. 100: Planning, Discovery Subagents, Up-Front Brain & Roadmap Creation, Lean Subtasks, Mode Evaluation)
PHASE_2_BUDGET = 150     (Steps 101 .. 250: Intermediate Archival, Sanitization, Prerequisite Refactoring, Intermediate Checkpoints)
PHASE_3_BUDGET = 150     (Steps 251 .. 400: Tree Swarm Parallel Execution, Dynamic Join Capacity Saturation, Mode Verification, Atomic Push)
WAVES = ceil(subtasks / (A x H))
```

> [!IMPORTANT]
> **Prompt Version:** 7.0.0  
> **Runtime:** Google Antigravity 2.0 (IDE and CLI)  
> **Canonical Prompt Specification:** [`01-prompts/14-execute/09-execute-parent-task-with-n-steps-v7.md`](file:///01-prompts/14-execute/09-execute-parent-task-with-n-steps-v7.md)  
> **Toolchain Discovery & Primacy:** [/learn](slashCommand;learn) from [gitmap](file;.agents/skills/gitmap)  
>
> **Top-Instruction Priority Mandate (Above Precedence / Preamble Precedence):**  
> Whatever directives, constraints, checklists, or user instructions are given ABOVE this prompt (including in the user preamble, header blocks, or incoming user request above) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below.

---

## 🌟 Dual Operating Modes: Restricted vs. Unrestricted Time

Upon ingestion, the agent evaluates the task requirements against two explicit operating options:

### Option 1: Time-Restricted Mode (`isTimeRestricted: true`)

Active when the user prompt specifies time windows, schedule bounds, restricted commit pacing, off-hours execution, or when invoked by time-restricted Letterly formatters.

1. **Target Timezone Resolution Hierarchy:**
   - Precedence 1: CLI flag `--tz <timezone>` or `--timezone <timezone>`.
   - Precedence 2: Task configuration or preamble directive (`Timezone: <tz>`).
   - Precedence 3: Split-DB setting `commit_timezone` via `gitmap config get commit_timezone`.
   - Precedence 4: Host operating system local timezone (`time.Local`).
2. **Calendar-Aware Window Constraints:**
   - **Weekday Constraint (Monday–Friday):** Commits are strictly permitted **before 08:00 AM** (`00:00:00` to `07:59:59`) in the resolved target timezone. Commits attempted at or after 08:00 AM are blocked unless auto-window simulation is enabled.
   - **Weekend Rule (Saturday–Sunday):** Commits are permitted **at any time** (24 hours unconstrained, `00:00:00` to `23:59:59`).
3. **Headless & One-Liner Automation (Mandatory for Subagents):**
   - Subagents and automated pipelines must never hang on interactive prompts; always pass headless flags:
     ```bash
     gitmap cpf "Module - summary" --tz <timezone> --schedule --auto-window --non-interactive
     ```
   - `--auto-window`: Automatically calculates the next valid morning window slot (07:30 AM) and sets compliant `GIT_AUTHOR_DATE` / `GIT_COMMITTER_DATE` timestamps.
4. **Recommit & History Normalization:**
   - When amending or recommitting prior work: verify clean working tree (`gitmap st`), soft reset (`gitmap git-reset HEAD~1 --no-push`), and recommit using semantic verbs with target timezone stamps.
   - Operational reference: [gitmap-time-restriction](file;.agents/skills/gitmap-time-restriction).

---

### Option 2: Unrestricted Time Mode (`isTimeRestricted: false`)

Active by default when no time restrictions, maintenance windows, or scheduling constraints are mentioned in the task prompt.

1. **Immediate Execution & Pacing:**
   - Changes are implemented and committed immediately without maintenance window checks.
   - Standard semantic GitMap commit verbs are used (`gitmap cpf`, `gitmap cbf`, `gitmap cpc`, `gitmap cpr`).
2. **Zero Delay:**
   - Commit timestamps reflect current wall-clock time without restriction gating.
3. **Storage Invariant:**
   - Persisted git timestamps and SQLite database entries remain strictly normalized to UTC (`t.UTC()`).

---

## 🤖 AI Automation Bot Integration (Applies to Both Options)

Active when the user requests an automated bot persona, AI commit authoring, or headless execution.

1. **Identity Presets:**
   - Name: `"gitmap automation"` (or `"AI automation bot"`).
   - Email: `"local@local"`.
2. **Activation:**
   - Run `gitmap autobot` or `gitmap change-user "gitmap automation <local@local>"`.
3. **Pairing with Time Restriction:**
   - When both Autobot and Time-Restriction are active, automated runs must combine both:
     ```bash
     gitmap autobot
     gitmap cpf "Module - summary" --tz <timezone> --schedule --auto-window --non-interactive
     ```
4. **Operational reference:** [gitmap-automation-bot](file;.agents/skills/gitmap-automation-bot).

---

## 🚨 Mandatory Subagent Spawning Gate (A = 3, H = 2, Capacity = 6)

- **Actual Tool Call Required:** You must actually execute the `invoke_subagent` tool. Printing text alone is an auto-reject failure.
- **Tree Swarm Hierarchy:**
  1. Root Orchestrator (Lead Agent): Unified roadmap, audit ledger, final GitMap commit. Never executes file edits solo.
  2. Branch Coordinators (Subtask Primary Owners): Own designated disjoint file boundaries (`Owned Files`).
  3. Leaf Workers & Dynamic Joining Assistants: Companion reading, symbol discovery, lint verification.
- **3-Stage Dispatch:**
  - Planning Step: 3 `research` subagents.
  - Spec Step: 3 `self` authoring subagents (Zero solo spec authoring).
  - Execution Step: 3 `self` worker subagents.
- **Immediate Dynamic Respawn:** As soon as any worker finishes, immediately spawn a replacement or dynamically join an active subtask.

---

## The Unified Master Pipeline

### Phase 1A: Verbatim Capture & Turn-1 Showcase
In your VERY FIRST turn, output the 5 Mandatory Pillars:
1. Strategic Intent & Systems Reasoning
2. Explicit Scope Boundaries (In-Scope vs. Excluded)
3. Non-Negotiable Constraints & Negative Guardrails
4. Blast Radius & Dependency Topology
5. Multi-Agent Concurrency Topology ($A=3, H=2, \text{Capacity}=6$)

Followed by confirmed task breakdown with explicit `Understood / AI Comprehension` and `Azure AI / Model Runtime Context`.

### Phase 1B: Mode Selection & On-Call Gate
- Evaluate whether task follows **Option 1 (Restricted Time)** or **Option 2 (Unrestricted Time)**.
- Evaluate whether **AI Automation Bot** identity is required.
- Configure pre-flight parameters accordingly.

### Phase 1C: Planning & Spec Step (Steps 1–100)
- Research codebase via 3 `research` subagents.
- Author specs in `02-spec/21-app/XX-<slug>/` and plans in `.ai-memory/plans/pending/XX-<slug>.md` via 3 authoring subagents.
- Initialize SQLite task database at `.ai-memory/temp-agents/<slug>/agent-task.db`.

### Phase 2: Prerequisite Setup & Disjoint Execution (Steps 101–250)
- Workers execute modifications in strictly disjoint file boxes.
- 12-Pillar Quality Gate: leaf error sovereignty, parameter structs, positive booleans, function limits <= 15 lines, zero nested ifs, relative Git paths only.

### Phase 3: Verification, Audit & Atomic GitMap Commit (Steps 251–400)
- Run targeted checks on modified files.
- If Option 1 is active, verify commit timestamp against target timezone (< 08:00 AM weekdays or anytime weekends).
- If Autobot is active, verify author credentials.
- Move plan to `.ai-memory/plans/completed/XX-<slug>.md`.
- Update registries in `.ai-memory/plans/readme.md` and `02-spec/21-app/readme.md`.
- Execute final atomic commit via GitMap semantic verbs (`gitmap cpf`, `gitmap cbf`, `gitmap cpc`, `gitmap cpr`).

---

## Core Invariants & Rules Reference

- **R1 Zero Builds / Tests Ban:** Never run full builds or test suites. Only targeted, file-scoped checks on modified files.
- **R2 Targeted Checks Only:** Check only specifically modified files.
- **R3 Evidence or It Did Not Happen:** Every claim must cite a file path, diffstat, or exit code.
- **R4 Never Invent Commands:** Verify commands with harmless calls.
- **R5 Mandatory Subagents:** Always call `invoke_subagent` ($A=3, H=2, \text{Capacity}=6$).
- **R6 One Owner Per File:** Disjoint file boundaries for workers.
- **R7 Worker Git Ban:** Workers never run git commands; only lead commits.
- **R8/R9 Atomic Commit via GitMap:** Single atomic commit. Title Case or lowercase subsystem module; NO square brackets `[...]`, NO colons, and NO 'commit' pseudo-module.
- **R10 Zero Unauthorized Releases:** Never bump versions unless requested.
- **R11 Leaf Error Sovereignty:** Leaf functions return errors directly; never call exit handlers.
- **R11b Strict Relative Git Paths:** Absolute paths and `file:///` URIs are banned.
