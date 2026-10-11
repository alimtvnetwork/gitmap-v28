# Subtask 03: Skills & Prompts Update — Agent AI Synchronization & Heatmap Integration

> **Subtask ID:** Task-03  
> **Parent Plan:** `.ai-memory/plans/agent-ai-synchronization-and-task-management.md`  
> **Spec Reference:** `02-spec/21-app/agent-ai-synchronization-and-task-management/02-component-and-cli-spec.md`  
> **Assigned Role:** Worker 02  
> **Status:** PENDING  
> **Owned Files:**  
> - `.agents/skills/gitmap/SKILL.md`  
> - `.cursor/skills/gitmap/skill.md`  
> - `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`  
> - `.cursor/skills/execute-parent-task-with-n-steps-v6/skill.md`  
> - `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md`  

---

## 1. Objectives & Executive Scope

This subtask updates all canonical GitMap skills, Cursor AI skills, and the V6 execution prompt to integrate the **Agent AI Synchronization and Task Management** workflow. Autonomous AI agents and multi-worker swarms must be explicitly instructed through their prompt protocols and skill definitions to:
1. Initialize tasks via `gitmap agent-ai create-task` and consume JSON responses detailing `parentTaskId`, `taskSlug`, and `machineId`.
2. Enqueue subtasks via `gitmap agent-ai subtask create`, propagating `parentTaskId` and `machineId` context.
3. Reserve look-ahead file boxes via `gitmap agent subtask claim-files` to trigger automated collision detection before modifying code.
4. Enforce mandatory pre-touch declarations via `gitmap agent-ai editing` strictly using the canonical `ActionType` enum (`modify`, `analyze`, `read`, `remove`, `refactor`, `create`).
5. Execute ultra-fast heatmap discussions via `gitmap agent-ai heatmap discuss` ($<5\text{ms}$) during planning, review, and peer synchronization.

---

## 2. Target Files & Detailed Modifications

### 2.1 `.agents/skills/gitmap/SKILL.md`

#### A. Command Replacement Matrix Update (Line ~40)
- Update the prohibited legacy pattern row for task coordination:
  - **Prohibited:** Uncoordinated multi-agent file edits, manual lock checks.
  - **GitMap Mandatory Command:** `gitmap agent-ai create-task`, `gitmap agent-ai subtask create`, `gitmap agent subtask claim-files`, `gitmap agent-ai editing`, `gitmap agent-ai heatmap discuss` (fallback: `gitmap task`).
  - **Rationale:** Two-Tier Split-DB architecture (`ai_agents.db` + `agent-task.db`), deterministic `machineId` propagation, sub-millisecond atomic file claiming, crash forensics, pre-touch in-flight logging with `ActionType` enum, and sub-5ms ANSI heatmaps.

#### B. Section 10: AI Agent Coordination & Two-Tier Split-DB Telemetry
- Document `gitmap agent-ai create-task "<name-or-slug>" [--budget <N>] [--json]`:
  - Detail JSON output fields: `isSuccess`, `parentTaskId`, `taskSlug`, `machineId`, `runDirectory`, `dbPath`, `status`.
- Document `gitmap agent-ai subtask create <slug/id> "<title>" [--code <Code>] [--files <paths>] [--role <role>] [--json]`:
  - Detail `parentTaskId` and `machineId` propagation, relative path validation, and JSON output schema.
- Document `gitmap agent subtask claim-files --task-id <id> --subtask <code/id> --files "<paths>" --agent "<role>" [--json]`:
  - Detail look-ahead reservation, collision detection logic, and warning banner.
- Document `gitmap agent-ai editing <subtask> "<file>" --action <action> --lines <x,y> --reasoning "<why>"`:
  - Detail mandatory canonical `ActionType` enum values: `modify`, `analyze`, `read`, `remove`, `refactor`, `create`.
  - Detail $<2\text{ms}$ write latency and crash forensics utility.
- Document `gitmap agent-ai heatmap discuss [<slug>] [--fast] [--json]`:
  - Detail sub-5ms performance guarantee, contention score formula, and prompt discussion workflow.

#### C. Operational Guardrail 13 (Rule R-PRETOUCH)
- Update Rule 13 to explicitly mandate the canonical `ActionType` enum (`modify`, `analyze`, `read`, `remove`, `refactor`, `create`) when declaring in-flight edits.

---

### 2.2 `.cursor/skills/gitmap/skill.md`

- Bring `.cursor/skills/gitmap/skill.md` into 100% parity with `.agents/skills/gitmap/SKILL.md`:
  - Replace Section 10 legacy commands with the full `gitmap agent-ai` command family: `create-task`, `subtask create`, `claim-files`, `editing`, `heatmap discuss`, `heatmap ls/search/grep`, `learn`.
  - Add Rule 13 (`Rule R-PRETOUCH`) to Operational Guardrails.
  - Update the Command Replacement Matrix table to mirror `.agents/skills/gitmap/SKILL.md`.

---

### 2.3 `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md`

#### A. Command Matrix Table (Section 4, ~Line 154)
- Update `Agent Task Engine` row to cite `gitmap agent-ai create-task`, `subtask create`, `claim-files`, `editing`, and `heatmap discuss`.

#### B. Step 0: Preflight, Platform Handshake & Task DB Initialization (Section 5, ~Line 172)
- Replace legacy single-command initialization with:
  ```bash
  # Primary: High-speed Agent AI task initialization emitting structured JSON
  gitmap agent-ai create-task "<task name>" --budget 300 --json
  # Fallback:
  gitmap task init --name "<task name>" --budget 300
  ```
- Instruct the orchestrator to capture `parentTaskId`, `taskSlug`, `machineId`, `runDirectory`, and `dbPath` from the JSON response and record them in `ledger.md`.

#### C. Phase 1: Planning & Spec Step (Section 6, ~Line 225)
- Replace legacy subtask addition with:
  ```bash
  # Primary: Register subtask with parent task lineage and machine context
  gitmap agent-ai subtask create <taskSlug> "<subtask title>" --code "Task-01" --files "<paths>" --role "Worker 01" --json
  ```
- Add pre-execution heatmap discussion directive:
  ```bash
  # Evaluate active claims and heated files before finalizing execution waves
  gitmap agent-ai heatmap discuss <taskSlug> --fast --json
  ```

#### D. Phase 2: Execution Step (Worker Wave Dispatch, Section 7, ~Lines 250–320)
- Update subagent dispatch prompt templates (`invoke_subagent` payload):
  - Instruct worker subagents to claim file boxes upon wakeup:
    ```bash
    gitmap agent subtask claim-files --task-id "<taskSlug>" --subtask "<taskCode>" --files "<paths>" --agent "<role>" --json
    ```
  - Mandate pre-touch declaration before EVERY file mutation tool call:
    ```bash
    gitmap agent-ai editing <taskCode> "<file>" --action <modify|analyze|read|remove|refactor|create> --lines <x,y> --reasoning "<rationale>"
    ```
  - Note that `--action` MUST be one of `modify`, `analyze`, `read`, `remove`, `refactor`, `create`.

---

### 2.4 `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md` & `.cursor/skills/execute-parent-task-with-n-steps-v6/skill.md`

- Update both skill markdown files to mirror all changes made to `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md`:
  - Preflight Step 0: `gitmap agent-ai create-task` with JSON capture of `machineId`.
  - Phase 1: `gitmap agent-ai subtask create` and `gitmap agent-ai heatmap discuss`.
  - Phase 2: Mandatory `claim-files` and `gitmap agent-ai editing` with `ActionType` enum.
  - Command replacement tables and operational rules.

---

## 3. Acceptance Criteria

| ID | Criterion | Verification Method | Expected Outcome |
|---|---|---|---|
| **AC-03.1** | **GitMap Skills Parity** | Compare Section 10 of `.agents/skills/gitmap/SKILL.md` and `.cursor/skills/gitmap/skill.md` | Both files document `create-task`, `subtask create`, `claim-files`, `editing`, and `heatmap discuss` with identical syntax and options. |
| **AC-03.2** | **V6 Prompt Synchronization** | Review `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md` | Steps 0, Phase 1, and Phase 2 dispatch templates explicitly mandate `agent-ai` commands, `MachineId` capture, and `ActionType` enum. |
| **AC-03.3** | **V6 Skills Parity** | Compare `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md` and `.cursor/skills/.../skill.md` against prompt | 100% content synchronization with zero missing sections or divergent command flags. |
| **AC-03.4** | **ActionType Enum Completeness** | Text search across all 5 files for `--action` | All occurrences document all 6 allowed values: `modify`, `analyze`, `read`, `remove`, `refactor`, `create`. |
| **AC-03.5** | **Strict Relative Path Compliance** | Run `python 03-ai-scripts/50-fastgate.py` | Exit 0 with zero absolute path violations or `file://` URIs in modified files. |

---

## 4. Execution Plan & Action Items

- [ ] **Step 1:** Update `.agents/skills/gitmap/SKILL.md` Command Replacement Matrix, Section 10, and Rule 13.
- [ ] **Step 2:** Update `.cursor/skills/gitmap/skill.md` to achieve 1:1 parity with `.agents/skills/gitmap/SKILL.md`.
- [ ] **Step 3:** Update `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md` (Step 0 preflight, Phase 1 subtask creation, Phase 2 pre-touch dispatch).
- [ ] **Step 4:** Update `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md` with identical updates.
- [ ] **Step 5:** Update `.cursor/skills/execute-parent-task-with-n-steps-v6/skill.md` with identical updates.
- [ ] **Step 6:** Run `python 03-ai-scripts/50-fastgate.py` to verify AST guidelines and relative path hygiene.

---

## 5. Verification Evidence Checklist

- [ ] AST relative path linter passes with zero warnings.
- [ ] All 5 target files contain strictly forward-slash relative Git paths.
- [ ] All action enum references match canonical set `modify, analyze, read, remove, refactor, create`.
- [ ] Task manager commands and flags match specification in `02-spec/21-app/agent-ai-synchronization-and-task-management/02-component-and-cli-spec.md`.
