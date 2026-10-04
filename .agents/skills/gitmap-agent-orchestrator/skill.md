---
name: gitmap-agent-orchestrator
description: Autonomously orchestrate, track, and diagnose multi-agent workflows using GitMap's 3-Tier Split-DB architecture, native CLI commands (gitmap agent / gitmap ai-agents), granular telemetry, crash forensics, and real-time Web UI visualizer.
---

# GitMap AI Agent Task Orchestrator & 3-Tier Split-DB Skill

## 1. Overview & Architectural Blueprint

GitMap provides a native, high-performance task orchestration and crash forensics engine for multi-agent LLM systems (`Google Antigravity`, Claude Code, and autonomous pair-programming pipelines). It eliminates concurrency write contention on SQLite by leveraging a **Three-Tier Multi-Agent SQLite Split-DB Hierarchy**.

### 1.1 Three-Tier Multi-Agent Split-DB Architecture

```text
.ai-memory/temp-agents/                       <-- Root Temp Directory (Configurable via gitmap config set agent.temp-dir)
├── ai_agents.db                              <-- TIER 1: Repository Root Master Agent DB (Global registry of all tasks, agents, lifecycles)
└── <sequence>-<task-slug>/                   <-- TIER 2: Run-Scoped Parent Task Directory (e.g. 79-ai-agent-task-orchestrator-and-split)
    ├── agent-task.db (or task.db)            <-- Task Root DB (ParentTask, Subtask manifest, claiming, rollups)
    ├── ledger.md                             <-- Markdown ledger for human/LLM transparency
    └── agents/                               <-- TIER 3: Agent-Scoped Split DBs
        ├── worker-01.db                      <-- Dedicated Agent Telemetry, ActionLog & File-touch DB
        └── worker-02.db                      <-- Dedicated Agent Telemetry, ActionLog & File-touch DB
```

### 1.2 Hierarchy Details & Database Responsibilities

1. **Tier 1: Master Agent Database (`ai_agents.db`)**
   - **Location:** `.ai-memory/temp-agents/ai_agents.db`
   - **Responsibility:** Global ledger tracking all parent tasks, overall status, total budget steps consumed, agent registries, and cross-task crash forensics.
   - **Key Tables:** `ParentTaskRegistry`, `AgentRegistry`, `GlobalLifecycleMetrics`.

2. **Tier 2: Parent Task Root Database (`agent-task.db`)**
   - **Location:** `.ai-memory/temp-agents/<nn>-<slug>/agent-task.db`
   - **Responsibility:** Atomic subtask manifest for a specific run. Coordinates worker claiming (`CLAIMED`), execution state (`IN_PROGRESS`), verification evidence, and completion rollups.
   - **Key Tables:** `ParentTask`, `Subtask`, `SubtaskAuditRollup`.

3. **Tier 3: Worker Split Databases (`agents/<agent-slug>.db`)**
   - **Location:** `.ai-memory/temp-agents/<nn>-<slug>/agents/<agent-slug>.db`
   - **Responsibility:** High-frequency, lock-free action telemetry. Every autonomous worker records every search, file read, file edit, or command execution directly into its private database, eliminating SQLite write lock contention.
   - **Key Tables:** `AgentActionLog`.

---

## 2. Complete CLI Command Reference (`gitmap agent` & `gitmap ai-agents`)

The CLI commands provide native Go execution, bypassing the need for legacy Python runner scripts.

```text
gitmap agent [subcommand] [flags]
gitmap ai-agents [subcommand] [flags]
```

### 2.1 Parent Task Management (`gitmap agent task`)

| Command | Flags | Description |
| :--- | :--- | :--- |
| `gitmap agent task init` | `--name "<name>" [--budget 300] [--dir <path>]` | Initializes parent task, provisions run directory, creates `agent-task.db` with WAL mode, and registers in `ai_agents.db`. |
| `gitmap agent task ls` | `[--limit <n>] [--all] [--status <s>]` | Lists recent parent tasks with status, duration, and completion rates. |
| `gitmap agent task status` | `[task-id]` | Displays comprehensive rollup of active task, subtask counts, budget meter, and in-flight agents. |

### 2.2 Atomic Subtask Management (`gitmap agent subtask`)

| Command | Flags | Description |
| :--- | :--- | :--- |
| `gitmap agent subtask add` | `--parent <id> --file <file>` (or `--json '<json>'` / `--code <c> --title <t> --owned '<files>'`) | Enqueues subtasks to the parent task database. |
| `gitmap agent subtask claim` | `--agent "<role>" [--task-id <id>]` | Atomically claims the next available pending subtask for the worker. |
| `gitmap agent subtask start` | `<subtask-id> --agent "<role>"` | Marks claimed subtask as actively `IN_PROGRESS`. |
| `gitmap agent subtask complete` | `<subtask-id> --agent "<role>" --evidence "<ev>"` | Marks subtask as `DONE` with verification evidence (`exit 0`, diffstat). |
| `gitmap agent subtask fail` | `<subtask-id> --agent "<role>" --reason "<reason>"` | Marks subtask as `FAILED` with Root Cause Analysis (RCA) details. |
| `gitmap agent subtask ls` | `[task-id] [--json]` | Lists all decomposed subtasks under the specified parent task. |

### 2.3 Granular Action Logging (`gitmap agent log`)

Autonomous workers MUST log actions immediately BEFORE modifying files or running commands to ensure crash forensics can locate the exact failure point.

```bash
gitmap agent log \
  --agent "Worker 01" \
  --subtask 1 \
  --action "WRITE" \
  --file "cli/cmdagent/agent_cmd.go" \
  --start-line 10 \
  --end-line 25 \
  --details "Refactoring Cobra subcommand routing" \
  --duration-ms 42 \
  --status "SUCCESS"
```

#### ActionType Enum Specification:
- `SEARCH`: Executing codebase searches (`gitmap aum search`, `gitmap find`)
- `READ`: Inspecting file contents (`view_file`, `gitmap cat`)
- `WRITE`: Creating or modifying source code (`write_to_file`, `replace_file_content`)
- `EXEC`: Running shell commands or linters (`run_command`, `gitmap pwsh`)
- `LINT`: Targeted coding guideline checks (`05-guideline-autofixer.py`)
- `CHECK`: Verification checks or preflight validations
- `CLAIM`: Subtask atomic claim event
- `START`: Subtask execution start event
- `COMPLETE`: Subtask verification and completion event
- `FAIL`: Subtask failure event with error message
- `CRASH`: Detected abnormal termination or abort

### 2.4 Crash Forensics & Autopsy (`gitmap agent crashed` & `gitmap agent diagnose`)

When a worker fails to report, crashes, or is terminated prematurely:

```bash
gitmap agent crashed [--task-id <id>]
# Alias:
gitmap agent diagnose [--task-id <id>]
```

#### Forensic Autopsy Report Schema:
```json
{
  "hasCrashesDetected": true,
  "crashedAgents": [
    {
      "subtaskId": 2,
      "taskCode": "Task-02",
      "title": "Implement CLI router",
      "assignedAgentRole": "Worker 01",
      "lastActionType": "WRITE",
      "lastTargetFile": "cli/cmdagent/agent_cmd.go",
      "lastActionDetails": "Refactoring Cobra subcommand routing",
      "lastTimestamp": "2026-10-04T11:22:50Z",
      "diagnosis": "Worker exited while performing WRITE on cli/cmdagent/agent_cmd.go without completing subtask"
    }
  ],
  "counts": {
    "done": 1,
    "pending": 3,
    "inProgress": 1,
    "failed": 0
  },
  "taskDirectory": ".ai-memory/temp-agents/79-ai-agent-task-orchestrator-and-split"
}
```

### 2.5 Real-Time Web UI Visualizer Dashboard (`gitmap agent ui`)

Launches an interactive developer dashboard in the default browser:

```bash
gitmap agent ui [--port 8095] [--browse]
```

- **Collision-Free Port Binding:** Scans from port 8095 upwards until an available local port is acquired.
- **REST Endpoints:**
  - `GET /api/agent/status`: Master orchestrator metrics, active task, and temp directory path.
  - `GET /api/agent/tasks`: Tree hierarchy of parent tasks and subtasks.
  - `GET /api/agent/logs?taskId=<id>&agent=<agent>`: Searchable action timeline with file paths, line ranges, and durations.
  - `GET /api/agent/crashes?taskId=<id>`: Crash forensics modal data.
  - `POST /api/agent/clear`: Trigger task archiving or cleanup directly from browser.

### 2.6 Lifecycle Cleanup & Storage Governance

Agent temp files are ephemeral and never committed to git:

| Command | Flags | Purpose |
| :--- | :--- | :--- |
| `gitmap agent clear` | `[--task-id <id>] [-y]` | Cleans specific completed task directory or all completed tasks in temp dir. |
| `gitmap agent reset` | `[-y]` | Clears agent split databases and resets tables while preserving run directories. |
| `gitmap agent temp-clear` | `[-y]` | Completely purges `.ai-memory/temp-agents/`, restoring clean state. |

---

## 3. Configurable Temp Directory Hierarchy

GitMap resolves the agent temp directory using the following strict precedence:
1. **CLI Flag:** `--dir <path>`
2. **Environment Variable:** `GITMAP_AGENT_TEMP_DIR`
3. **Database Setting:** `gitmap config set agent.temp-dir <path>`
4. **Default Workspace Path:** `.ai-memory/temp-agents/`

---

## 4. Multi-Agent Orchestration Protocol (Lead & Worker Roles)

### 4.1 Lead Orchestrator Workflow

1. **Step 0 Preflight:**
   - Initialize task:
     ```bash
     gitmap agent task init --name "<task-slug>" --budget 300
     ```
   - Check if resuming from crash:
     ```bash
     gitmap agent crashed
     ```
2. **Phase 1 Planning & Spec:**
   - Decompose subtasks and register in task DB:
     ```bash
     gitmap agent subtask add --parent <id> --file <subtasks-file.json>
     ```
3. **Phase 2 Execution Waves:**
   - Spawn workers (`invoke_subagent`, `TypeName: "self"`).
   - In case of silent worker exit or timeout:
     ```bash
     gitmap agent crashed
     ```
   - Monitor live progress:
     ```bash
     gitmap agent task status
     ```
4. **Phase 3 Consolidation:**
   - Consolidate evidence and push atomically via GitMap:
     ```bash
     gitmap cpf "<module> - <summary>"
     ```
   - Purge ephemeral temp state when desired:
     ```bash
     gitmap agent clear --task-id <id> -y
     ```

### 4.2 Autonomous Worker Subagent Workflow

Every worker spawned via `invoke_subagent` must follow this lifecycle:

1. **Claim Task Atomically:**
   ```bash
   gitmap agent subtask claim --agent "Worker 01"
   ```
2. **Mark Subtask In Progress:**
   ```bash
   gitmap agent subtask start <subtaskId> --agent "Worker 01"
   ```
3. **Log Actions Before File Modification:**
   ```bash
   gitmap agent log --agent "Worker 01" --subtask <subtaskId> --action "WRITE" --file "cli/module.go" --details "Adding handler"
   ```
4. **Run Targeted File Checks (Zero Build/Test Runs):**
   ```bash
   python 03-ai-scripts/05-guideline-autofixer.py <folder> --check-only
   ```
5. **Record Completion with Evidence:**
   ```bash
   gitmap agent subtask complete <subtaskId> --agent "Worker 01" --evidence "PASS exit 0, files: cli/module.go"
   ```
6. **On Failure / Blocker:**
   ```bash
   gitmap agent subtask fail <subtaskId> --agent "Worker 01" --reason "Syntax conflict on line 42"
   ```

---

## 5. Non-Negotiable Operational Rules

1. **Positive Booleans Only:** Use affirmative prefixes `is` and `has` exclusively. Never evaluate explicit `== true`.
2. **Strict Relative Git Paths:** All file paths must be relative to the repository root. Never use absolute paths (`C:\...`) or `file:///` URIs.
3. **Zero Test Runs & Zero Build Runs:** Routine worker execution turns must NEVER run `go test ./...`, `go build`, `npm test`, or `pytest`. Compilation and test suites are strictly verified in CI/CD.
4. **Worker Git Ban:** Subagents must NEVER run any git commands (`git add`, `git commit`, `git push`, `git status`). Staging and committing is exclusively performed by the lead orchestrator via `gitmap cpf` / `gitmap cpb`.
