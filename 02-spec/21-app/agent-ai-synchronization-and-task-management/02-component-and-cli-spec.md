# Component & CLI Specification: Agent AI Synchronization and Task Management

- **Feature Directory:** `02-spec/21-app/agent-ai-synchronization-and-task-management/`
- **Slug:** `agent-ai-synchronization-and-task-management`
- **Specification Version:** 1.0.0
- **Status:** APPROVED
- **Target Modules:** `cli/cmdagent`, `cli/cmdagentai`, `03-ai-scripts/46-agent-sqlite-task-manager.py`
- **Target Audience:** Autonomous AI Agents, Multi-Agent Workers, Lead Orchestrators, GitMap Engineers
- **Target Release:** Minor Bump

---

## 1. Executive Summary & Architectural Alignment

The **Agent AI Synchronization and Task Management Suite** provides high-speed, concurrency-safe coordination primitives for autonomous multi-agent engineering workflows. Distributed AI agent swarms (such as the V6 orchestrator running $A=2, H=2$ worker pools) require instantaneous task lifecycle tracking, deterministic machine identity propagation, look-ahead file reservation, pre-touch editing declarations, and sub-5ms contention analytics.

### 1.1 Core Problems Addressed
1. **Uncoordinated Multi-Agent Collisions:** Without synchronized file reservations, concurrent agents mutate identical source files simultaneously, triggering OS write-lock contentions, Git merge conflicts, and lost progress.
2. **Missing Machine & Lineage Context:** Subtasks created in headless terminal environments frequently lose parent task lineage and host machine context, causing unrecoverable state loss when diagnosing crashes across heterogeneous nodes.
3. **Unvalidated Action Types:** Ad-hoc free-text action labels prevent automated semantic auditing, crash forensics, and heatmap churn weighting.
4. **Slow Heatmap Inspection:** Real-time multi-agent discussions require immediate file churn and contention telemetry. In-flight discussions cannot tolerate slow multi-second scans.

### 1.2 Two-Tier Split-DB Architecture
The CLI and Python engines operate across a coordinated two-tier SQLite architecture in WAL mode (`PRAGMA journal_mode = WAL; PRAGMA busy_timeout = 5000;`):
- **Tier 1 Master DB (`.ai-memory/temp-agents/ai_agents.db`):** Workspace-wide registry tracking all parent tasks, active file claims, global action indices, and cross-task heuristics.
- **Tier 2 Ephemeral Task DB (`.ai-memory/temp-agents/<slug>/agent-task.db`):** Run-scoped database storing parent task state, subtasks, worker allocations, and fine-grained in-flight action logs.

```
                        ┌────────────────────────────────────────────────────────┐
                        │             Tier 1 Master Registry DB                  │
                        │        (.ai-memory/temp-agents/ai_agents.db)           │
                        └──────────┬─────────────────────────────────┬───────────┘
                                   │                                 │
           ┌───────────────────────▼─────────┐     ┌─────────────────▼─────────┐
           │     ParentTaskRegistry          │     │        FileClaim          │
           │ (Slug, Status, MachineId, Path) │     │ (Path, Subtask, AgentRole)│
           └─────────────────────────────────┘     └───────────────────────────┘
                                   │
              ┌────────────────────┴────────────────────┐
              │                                         │
┌─────────────▼───────────────────────────┐ ┌───────────▼───────────────────────────┐
│     Tier 2 Ephemeral Task DB (Run A)    │ │     Tier 2 Ephemeral Task DB (Run B)    │
│ (.ai-memory/temp-agents/<slug-a>/...db) │ │ (.ai-memory/temp-agents/<slug-b>/...db) │
├─────────────────────────────────────────┤ ├─────────────────────────────────────────┤
│ ParentTask   (MachineId, TotalBudget)   │ │ ParentTask   (MachineId, TotalBudget)   │
│ Subtask      (TaskCode, MachineId, ...) │ │ Subtask      (TaskCode, MachineId, ...) │
│ AgentAction  (ActionType Enum, <2ms)    │ │ AgentAction  (ActionType Enum, <2ms)    │
└─────────────────────────────────────────┘ └─────────────────────────────────────────┘
```

---

## 2. Deterministic Machine Identification Contract (`MachineId`)

To ensure seamless multi-agent telemetry and cluster-wide crash forensics, every parent task and subtask must record a canonical `MachineId`.

### 2.1 Resolution Algorithm
The `MachineId` is evaluated deterministically at process initialization:
1. **Environment Override:** Inspect `GITMAP_MACHINE_ID` environment variable. If non-empty, sanitize to lowercase alphanumeric and hyphen characters.
2. **System Hostname:** Retrieve operating system hostname via standard platform APIs (`os.Hostname()` in Go, `platform.node()` in Python).
3. **Machine Hardware GUID Probe (Fallback):**
   - Windows: Query `HKLM\SOFTWARE\Microsoft\Cryptography\MachineGuid`.
   - Linux: Read `/etc/machine-id` or `/var/lib/dbus/machine-id`.
   - macOS: Query `ioreg -rd1 -c IOPlatformExpertDevice | grep IOPlatformUUID`.
4. **Deterministic Normalization:** If hostname is available, sanitize to lowercase kebab-case (e.g. `workstation-alpha-01`). If hardware GUID is used as fallback, compute SHA-256 and truncate to the first 12 characters prefixed with `node-` (e.g. `node-a1b2c3d4e5f6`).

### 2.2 Propagation Invariant
- Every call to `gitmap agent-ai create-task` writes `MachineId` to the Tier 1 `ParentTaskRegistry` and Tier 2 `ParentTask`.
- Every call to `gitmap agent-ai subtask create` inherits or validates `MachineId` and writes it directly into the Tier 2 `Subtask` record.
- All JSON responses emit `"machineId": "<resolved-id>"` at the top level.

---

## 3. Canonical Action Type Enum (`ActionType`)

Every file operation declared via `gitmap agent-ai editing` must strictly specify a valid `ActionType` enum value. Free-text strings and unapproved action names are rejected with exit code 1.

### 3.1 Allowed ActionType Values

| ActionType Enum | Description | Typical Tooling Triggers | Mutation Type |
|---|---|---|---|
| `modify` | In-place surgical updates, code edits, or line replacements | `replace_file_content` | Write / In-Place |
| `analyze` | Deep AST inspection, linting, complexity checks, or static audits | Custom AST parser, inspection script | Read-Only |
| `read` | Reading source files, specs, plans, or configurations | `view_file`, `gitmap cat` | Read-Only |
| `remove` | Deleting deprecated code, pruning dead functions, or unlinking files | File deletion, refactoring removal | Write / Destructive |
| `refactor` | Structural reorganization, function decomposition, or DRY extraction | Multi-call refactoring sequence | Write / Structural |
| `create` | Authoring brand new files, specs, plans, or tests from scratch | `write_to_file` | Write / Additive |

### 3.2 Enforcement Invariant
The Go CLI and Python engine validate `--action <val>` against the canonical set `{"modify", "analyze", "read", "remove", "refactor", "create"}`. If an unapproved value is provided (e.g. `--action edit` or `--action patch`), the command terminates immediately:
```text
Error: Invalid --action 'patch'. Must be one of: modify, analyze, read, remove, refactor, create.
```

---

## 4. CLI Command Syntax Contracts & Schemas

### 4.1 Parent Task Creation (`gitmap agent-ai create-task`)

Initializes a parent task across Tier 1 and Tier 2 databases with deterministic slug generation, monotonic step budget, and host machine registration.

#### Syntax
```bash
gitmap agent-ai create-task "<title-or-slug>" [--budget <N>] [--dir <relative-dir>] [--json]
```

#### Aliases
`gitmap aai create-task`, `gitmap task init`

#### Arguments & Flags
| Parameter | Short | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `<title-or-slug>` | - | String | Yes | - | Task title (e.g. `"Agent AI Synchronization and Task Management"`) or existing slug |
| `--budget` | `-b` | Integer | No | `300` | Monotonic step budget for the autonomous execution loop |
| `--dir` | `-d` | String | No | auto | Custom relative run directory under `.ai-memory/temp-agents/<slug>` |
| `--json` | `-j` | Boolean | No | `false` | Emits structured JSON payload instead of human ANSI terminal output |

#### Behavioral Logic
1. **Slugification:** Sanitizes `<title-or-slug>` into deterministic lowercase kebab-case (replaces punctuation with hyphens, caps at 6 significant words).
2. **Machine ID Resolution:** Evaluates deterministic `MachineId` for the executing node.
3. **Tier 1 Registration:** Inserts or updates row in `.ai-memory/temp-agents/ai_agents.db` (`ParentTaskRegistry` table) with status `'ACTIVE'`.
4. **Directory Scaffolding:** Creates relative directory `.ai-memory/temp-agents/<slug>/`.
5. **Tier 2 DB Initialization:** Connects to `.ai-memory/temp-agents/<slug>/agent-task.db`, applies `SCHEMA_DDL` (with WAL mode), and inserts record into `ParentTask`.
6. **Output Generation:** If `--json` is supplied, emits JSON envelope. Otherwise, renders ANSI summary table.

#### Exhaustive JSON Output Schema
```json
{
  "isSuccess": true,
  "parentTaskId": 1,
  "taskSlug": "agent-ai-synchronization-and-task-management",
  "taskName": "Agent AI Synchronization and Task Management",
  "machineId": "devbox-worker-01",
  "runDirectory": ".ai-memory/temp-agents/agent-ai-synchronization-and-task-management",
  "dbPath": ".ai-memory/temp-agents/agent-ai-synchronization-and-task-management/agent-task.db",
  "totalStepsBudget": 300,
  "status": "ACTIVE"
}
```

#### Terminal Human Output Format (Default)
```text
[1.1ms] Task Initialized: agent-ai-synchronization-and-task-management
────────────────────────────────────────────────────────────────────────
Parent Task ID : 1
Machine ID     : devbox-worker-01
Status         : ACTIVE
Budget Steps   : 300
Run Directory  : .ai-memory/temp-agents/agent-ai-synchronization-and-task-management
Database Path  : .ai-memory/temp-agents/agent-ai-synchronization-and-task-management/agent-task.db
────────────────────────────────────────────────────────────────────────
```

---

### 4.2 Subtask Creation (`gitmap agent-ai subtask create`)

Enqueues an isolated subtask under the parent task, recording parent lineage, host machine context, initial owned files, and worker role.

#### Syntax
```bash
gitmap agent-ai subtask create <slug-or-id> "<subtask-title>" [--code <Code>] [--files <paths>] [--role <role>] [--json]
```

#### Aliases
`gitmap aai subtask create`, `gitmap task add-subtask`

#### Arguments & Flags
| Parameter | Short | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `<slug-or-id>` | - | String | Yes | - | Parent task ID or slug (resolves via exact match or prefix match) |
| `"<subtask-title>"` | - | String | Yes | - | Concise descriptive title of the subtask deliverable |
| `--code` | `-c` | String | No | auto | Unique subtask code (e.g. `Task-01`). Auto-increments if omitted |
| `--files` | `-f` | String | No | `""` | Comma-separated strictly relative Git file paths owned by this subtask |
| `--role` | `-r` | String | No | `"Worker 01"` | Assigned agent role name |
| `--json` | `-j` | Boolean | No | `false` | Emits structured JSON payload |

#### Behavioral Logic
1. **Parent Resolution:** Resolves target parent task in Tier 1 registry and opens corresponding Tier 2 `agent-task.db`.
2. **Code Assignment:** If `--code` is omitted, counts existing subtasks and formats sequentially (`Task-01`, `Task-02`, etc.).
3. **Relative Path Validation:** Verifies that all comma-separated paths in `--files` are strictly relative Git paths (forward slashes, no drive letters, no `file://`).
4. **Tier 2 Insertion:** Inserts record into `Subtask` table containing `ParentTaskId`, `TaskCode`, `Title`, `AssignedAgentRole`, `OwnedFilesJson`, `Status='PENDING'`, and `MachineId`.
5. **Tier 1 Claim Synchronization:** If files are specified, registers them in Tier 1 `FileClaim` table.
6. **Output Generation:** Returns structured JSON or ANSI confirmation.

#### Exhaustive JSON Output Schema
```json
{
  "isSuccess": true,
  "parentTaskId": 1,
  "machineId": "devbox-worker-01",
  "subtaskId": 3,
  "taskCode": "Task-03",
  "taskSlug": "agent-ai-synchronization-and-task-management",
  "subtaskTitle": "Skills and Prompts Update",
  "assignedAgentRole": "Worker 02",
  "ownedFiles": [
    ".agents/skills/gitmap/SKILL.md",
    ".cursor/skills/gitmap/skill.md",
    ".agents/skills/execute-parent-task-with-n-steps-v6/skill.md",
    ".cursor/skills/execute-parent-task-with-n-steps-v6/skill.md",
    "01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md"
  ],
  "status": "PENDING"
}
```

#### Terminal Human Output Format (Default)
```text
[0.9ms] Subtask Created: Task-03 (Parent #1 / agent-ai-synchronization-and-task-management)
────────────────────────────────────────────────────────────────────────
Subtask ID    : 3
Task Code     : Task-03
Title         : Skills and Prompts Update
Machine ID    : devbox-worker-01
Assigned Role : Worker 02
Owned Files   : 5 files registered
Status        : PENDING
────────────────────────────────────────────────────────────────────────
```

---

### 4.3 Subtask File Box Claiming (`gitmap agent subtask claim-files`)

Declares look-ahead file box claims for a worker agent before operations commence, executing automated cross-worker collision detection.

#### Syntax
```bash
gitmap agent subtask claim-files --task-id <slug-or-id> --subtask <code-or-id> --files "<paths>" --agent "<role>" [--json]
```

#### Aliases
`gitmap agents claim-files`, `gitmap ai-agents subtask claim-files`

#### Arguments & Flags
| Parameter | Short | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `--task-id` | `-t` | String | Yes | - | Parent task ID or slug |
| `--subtask` | `-s` | String | Yes | - | Subtask ID (integer) or task code (`Task-01`) |
| `--files` | `-f` | String | Yes | - | Comma-separated strictly relative Git file paths |
| `--agent` | `-a` | String | Yes | - | Worker agent role name declaring the claim (e.g. `"Worker 01"`) |
| `--json` | `-j` | Boolean | No | `false` | Emits structured JSON payload |

#### Behavioral Logic
1. **Path Sanitization:** Trims whitespace, replaces backslashes with forward slashes, and validates relative path hygiene.
2. **Collision Evaluation:** Queries Tier 1 `FileClaim` table for any active claim on the same file held by a different subtask or different agent role.
3. **Collision Warning:** If collisions exist, constructs structured collision reports detailing conflicting subtask codes and agent roles.
4. **Persistence:** Inserts or updates claim records in Tier 1 `FileClaim` and serializes `OwnedFilesJson` in Tier 2 `Subtask`.
5. **Output:** Returns collision status and claim registration.

#### Exhaustive JSON Output Schema (No Collision)
```json
{
  "isSuccess": true,
  "parentTaskId": 1,
  "subtaskId": 2,
  "taskCode": "Task-02",
  "agentRole": "Worker 01",
  "claimedFiles": [
    "03-ai-scripts/46-agent-sqlite-task-manager.py",
    "cli/cmdagentai/agent_ai_cmd.go"
  ],
  "hasCollision": false,
  "collisions": [],
  "status": "CLAIMED"
}
```

#### Exhaustive JSON Output Schema (Collision Detected)
```json
{
  "isSuccess": true,
  "parentTaskId": 1,
  "subtaskId": 3,
  "taskCode": "Task-03",
  "agentRole": "Worker 02",
  "claimedFiles": [
    "03-ai-scripts/46-agent-sqlite-task-manager.py"
  ],
  "hasCollision": true,
  "collisions": [
    {
      "filePath": "03-ai-scripts/46-agent-sqlite-task-manager.py",
      "conflictingSubtask": "Task-02",
      "conflictingAgentRole": "Worker 01",
      "conflictSeverity": "CRITICAL"
    }
  ],
  "status": "COLLISION_WARNING"
}
```

#### Terminal Human Output Format with Collision Warning Banner
```text
[1.3ms] File Claim Processed: 1 files claimed
⚠️  FILE CLAIM COLLISION DETECTED:
File                 : 03-ai-scripts/46-agent-sqlite-task-manager.py
Active Claim Holder  : Worker 01 (Task-02)
Attempted Claimant   : Worker 02 (Task-03)
Guidance             : Workers must hold strictly disjoint file boxes. Worker 02 must re-partition!
```

---

### 4.4 Pre-Touch In-Flight Editing Declaration (`gitmap agent-ai editing`)

Mandatory pre-touch command that MUST be executed immediately before invoking any mutating or analyzing tool (`replace_file_content`, `write_to_file`, `view_file`). Commits to SQLite WAL log in $<2\text{ms}$.

#### Syntax
```bash
gitmap agent-ai editing <subtask-code-or-id> "<relative-file>" \
  --action <modify|analyze|read|remove|refactor|create> \
  --lines <start,end> \
  --reasoning "<why>" \
  [--agent "<role>"] \
  [--json]
```

#### Aliases
`gitmap aai editing`, `gitmap aai touch`

#### Arguments & Flags
| Parameter | Short | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `<subtask-code-or-id>` | - | String | Yes | - | Subtask code (`Task-01`) or numeric subtask ID |
| `"<relative-file>"` | - | String | Yes | - | Strictly relative Git path of the target file |
| `--action` | `-a` | String | Yes | - | Canonical `ActionType` enum (`modify`, `analyze`, `read`, `remove`, `refactor`, `create`) |
| `--lines` | `-l` | String | Yes | `0,0` | Target line range: `<start>,<end>` (e.g. `120,165` or `0,0` for entire file) |
| `--reasoning` | `-r` | String | Yes | - | Clear technical rationale explaining the impending operation |
| `--agent` | - | String | No | auto | Worker agent role name (auto-resolved from subtask if omitted) |
| `--json` | `-j` | Boolean | No | `false` | Emits structured JSON payload |

#### Behavioral Logic
1. **ActionType Validation:** Verifies `--action` is in `[modify, analyze, read, remove, refactor, create]`. Rejects violations with exit code 1.
2. **Subtask Disambiguation:** Resolves subtask record in Tier 2 DB by `TaskCode` (e.g. `Task-02`) or integer primary key.
3. **Sub-Millisecond Log Insertion:** Inserts record into Tier 2 `AgentActionLog` with status `'IN_PROGRESS'`.
4. **Global Index Mirroring:** Asynchronously updates Tier 1 `GlobalActionIndex` to reflect live churn in global heatmaps.
5. **Output:** Returns action record confirmation in $<2\text{ms}$.

#### Exhaustive JSON Output Schema
```json
{
  "isSuccess": true,
  "actionLogId": 42,
  "subtaskId": 2,
  "taskCode": "Task-02",
  "agentRole": "Worker 01",
  "actionType": "modify",
  "targetFile": "03-ai-scripts/46-agent-sqlite-task-manager.py",
  "lineRange": "120,165",
  "reasoning": "Implement ActionType enum validation and MachineId column in SCHEMA_DDL",
  "status": "IN_PROGRESS",
  "timestamp": "2026-10-11T10:25:00Z"
}
```

#### Terminal Human Output Format (Default)
```text
[0.8ms] In-Flight Action Logged: #42
Subtask    : Task-02 (Worker 01)
Action     : modify
File       : 03-ai-scripts/46-agent-sqlite-task-manager.py [Lines: 120-165]
Reasoning  : Implement ActionType enum validation and MachineId column in SCHEMA_DDL
Status     : IN_PROGRESS
```

---

### 4.5 Heatmap Discussion Summary (`gitmap agent-ai heatmap discuss`)

High-speed diagnostic command designed for autonomous multi-agent planning turns, peer synchronization, and lead orchestrator reviews. Evaluates heated files, active claims, pre-touch rationales, and contention risk scores in $<5\text{ms}$.

#### Syntax
```bash
gitmap agent-ai heatmap discuss [<task-slug-or-id>] [--fast] [--json]
```

#### Aliases
`gitmap aai heatmap discuss`, `gitmap aai hd`

#### Arguments & Flags
| Parameter | Short | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `[<task-slug-or-id>]` | - | String | No | auto | Parent task slug or ID (defaults to active task in `.ai-memory/temp-agents/`) |
| `--fast` | `-f` | Boolean | No | `true` | Skips disk verification probes for raw SQLite index speed ($<3\text{ms}$) |
| `--json` | `-j` | Boolean | No | `false` | Emits structured JSON summary for autonomous AI agent decision loops |

#### Contention Risk Scoring Formula
The contention risk score $S_{\text{contention}}$ evaluates collision probability across concurrent agents:
$$S_{\text{contention}} = (N_{\text{active\_claims}} \times 3) + (N_{\text{in\_flight\_actions}} \times 2) + N_{\text{recent\_touches}}$$

Where:
- $N_{\text{active\_claims}}$: Number of subtasks concurrently claiming the file box.
- $N_{\text{in\_flight\_actions}}$: Number of agents with status `'IN_PROGRESS'` touching the file.
- $N_{\text{recent\_touches}}$: Cumulative touch count within the current task run.

**Contention Level Thresholds:**
- **LOW** ($S \le 3$): Safe for isolated modification; single claimant.
- **MEDIUM** ($4 \le S \le 7$): Active attention; file touched by multiple subtasks sequentially.
- **HIGH** ($8 \le S \le 12$): Significant churn; multiple in-flight declarations.
- **CRITICAL** ($S > 12$ or multiple concurrent claims): Immediate collision risk! Workers must coordinate before writing.

#### Exhaustive JSON Output Schema
```json
{
  "isSuccess": true,
  "taskSlug": "agent-ai-synchronization-and-task-management",
  "evaluatedAt": "2026-10-11T10:28:15Z",
  "executionTimeMs": 1.45,
  "totalHeatedFiles": 3,
  "totalActiveClaims": 4,
  "highContentionFilesCount": 1,
  "files": [
    {
      "filePath": "03-ai-scripts/46-agent-sqlite-task-manager.py",
      "claimsCount": 2,
      "inFlightTouches": 1,
      "touchCount": 6,
      "contentionScore": 14,
      "contentionLevel": "CRITICAL",
      "activeClaimants": [
        {
          "agentRole": "Worker 01",
          "subtaskCode": "Task-02"
        },
        {
          "agentRole": "Worker 02",
          "subtaskCode": "Task-03"
        }
      ],
      "recentActions": [
        {
          "agentRole": "Worker 01",
          "actionType": "modify",
          "lineRange": "120,165",
          "reasoning": "Implement ActionType enum validation and MachineId column in SCHEMA_DDL",
          "timestamp": "2026-10-11T10:25:00Z"
        }
      ]
    },
    {
      "filePath": "cli/cmdagentai/agent_ai_cmd.go",
      "claimsCount": 1,
      "inFlightTouches": 0,
      "touchCount": 3,
      "contentionScore": 6,
      "contentionLevel": "MEDIUM",
      "activeClaimants": [
        {
          "agentRole": "Worker 01",
          "subtaskCode": "Task-02"
        }
      ],
      "recentActions": [
        {
          "agentRole": "Worker 01",
          "actionType": "refactor",
          "lineRange": "45,90",
          "reasoning": "Route heatmap discuss subcommand with --fast and --json flags",
          "timestamp": "2026-10-11T10:26:30Z"
        }
      ]
    },
    {
      "filePath": "02-spec/21-app/agent-ai-synchronization-and-task-management/02-component-and-cli-spec.md",
      "claimsCount": 1,
      "inFlightTouches": 1,
      "touchCount": 2,
      "contentionScore": 5,
      "contentionLevel": "MEDIUM",
      "activeClaimants": [
        {
          "agentRole": "Spec Worker 02",
          "subtaskCode": "Task-01"
        }
      ],
      "recentActions": [
        {
          "agentRole": "Spec Worker 02",
          "actionType": "create",
          "lineRange": "1,450",
          "reasoning": "Author comprehensive component and CLI specification",
          "timestamp": "2026-10-11T10:27:00Z"
        }
      ]
    }
  ]
}
```

#### Terminal Human Output Format (Default)
```text
[1.4ms] Heatmap Discussion Summary: agent-ai-synchronization-and-task-management
Evaluated: 3 heated files | Active Claims: 4 | Critical Contention: 1

FILE PATH                                    CLAIMS  TOUCHES  SCORE  LEVEL      ACTIVE WORKERS
03-ai-scripts/46-agent-sqlite-task-manager.py     2        6     14  CRITICAL   Worker 01, Worker 02
cli/cmdagentai/agent_ai_cmd.go                    1        3      6  MEDIUM     Worker 01
02-spec/21-app/.../02-component-and-cli-spec.md   1        2      5  MEDIUM     Spec Worker 02

─────────────────────────────────────────────────────────────────────────────────────────────
⚡ DISCUSSION ADVISORY & CONTENTION MITIGATION:
• Critical file detected: 03-ai-scripts/46-agent-sqlite-task-manager.py has overlapping claims!
  - Worker 01 (Task-02) is modifying lines 120-165.
  - Worker 02 (Task-03) must wait for Task-02 completion or operate on disjoint files.
─────────────────────────────────────────────────────────────────────────────────────────────
```

---

## 5. Database Schema & Migration Specification

Both the Go native engine and `03-ai-scripts/46-agent-sqlite-task-manager.py` enforce identical table DDL schemas and indexes.

### 5.1 Updated Table Definitions

```sql
-- Tier 1 Master Registry (.ai-memory/temp-agents/ai_agents.db)
CREATE TABLE IF NOT EXISTS ParentTaskRegistry (
    RegistryId INTEGER PRIMARY KEY AUTOINCREMENT,
    ParentTaskId INTEGER NOT NULL,
    TaskSlug TEXT NOT NULL UNIQUE,
    TaskName TEXT NOT NULL,
    MachineId TEXT NOT NULL,
    RunDirectory TEXT NOT NULL,
    DatabasePath TEXT NOT NULL,
    Status TEXT NOT NULL DEFAULT 'ACTIVE',
    TotalStepsBudget INTEGER NOT NULL DEFAULT 300,
    CurrentStep INTEGER NOT NULL DEFAULT 1,
    TotalSubtasks INTEGER NOT NULL DEFAULT 0,
    CreatedAt TEXT NOT NULL,
    UpdatedAt TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS FileClaim (
    ClaimId INTEGER PRIMARY KEY AUTOINCREMENT,
    ParentTaskId INTEGER NOT NULL,
    SubtaskCode TEXT NOT NULL,
    AgentRole TEXT NOT NULL,
    FilePath TEXT NOT NULL,
    ClaimStatus TEXT NOT NULL DEFAULT 'ACTIVE',
    CreatedAt TEXT NOT NULL,
    UpdatedAt TEXT NOT NULL
);

-- Tier 2 Ephemeral Database (.ai-memory/temp-agents/<slug>/agent-task.db)
CREATE TABLE IF NOT EXISTS ParentTask (
    ParentTaskId INTEGER PRIMARY KEY AUTOINCREMENT,
    TaskName TEXT NOT NULL,
    TaskSlug TEXT NOT NULL,
    MachineId TEXT NOT NULL,
    RunDirectory TEXT NOT NULL,
    Status TEXT NOT NULL DEFAULT 'ACTIVE',
    IsActive INTEGER NOT NULL DEFAULT 1,
    HasCompleted INTEGER NOT NULL DEFAULT 0,
    TotalStepsBudget INTEGER NOT NULL DEFAULT 300,
    CurrentStep INTEGER NOT NULL DEFAULT 1,
    Notes TEXT NULL,
    CreatedAt TEXT NOT NULL,
    UpdatedAt TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS Subtask (
    SubtaskId INTEGER PRIMARY KEY AUTOINCREMENT,
    ParentTaskId INTEGER NOT NULL,
    TaskCode TEXT NOT NULL,
    Title TEXT NOT NULL,
    MachineId TEXT NOT NULL,
    AssignedAgentRole TEXT NULL,
    OwnedFilesJson TEXT NOT NULL DEFAULT '[]',
    Status TEXT NOT NULL DEFAULT 'PENDING',
    IsBlocked INTEGER NOT NULL DEFAULT 0,
    HasCompleted INTEGER NOT NULL DEFAULT 0,
    Evidence TEXT NULL,
    CreatedAt TEXT NOT NULL,
    UpdatedAt TEXT NOT NULL,
    FOREIGN KEY (ParentTaskId) REFERENCES ParentTask(ParentTaskId)
);

CREATE TABLE IF NOT EXISTS AgentActionLog (
    ActionLogId INTEGER PRIMARY KEY AUTOINCREMENT,
    SubtaskId INTEGER NOT NULL,
    AgentRole TEXT NOT NULL,
    ActionType TEXT NOT NULL CHECK (ActionType IN ('modify', 'analyze', 'read', 'remove', 'refactor', 'create')),
    TargetFile TEXT NULL,
    ActionDetails TEXT NOT NULL,
    Status TEXT NOT NULL DEFAULT 'IN_PROGRESS',
    CreatedAt TEXT NOT NULL,
    FOREIGN KEY (SubtaskId) REFERENCES Subtask(SubtaskId)
);

CREATE INDEX IF NOT EXISTS idx_subtask_status ON Subtask(Status);
CREATE INDEX IF NOT EXISTS idx_subtask_code ON Subtask(TaskCode);
CREATE INDEX IF NOT EXISTS idx_action_subtask ON AgentActionLog(SubtaskId);
CREATE INDEX IF NOT EXISTS idx_action_target_file ON AgentActionLog(TargetFile);
```

---

## 6. Acceptance Criteria (AC-01 through AC-12)

| ID | Criterion | Verification Command | Expected Outcome |
|---|---|---|---|
| **AC-01** | **Parent Task Creation with Machine Context** | `gitmap agent-ai create-task "Sync Test" --budget 300 --json` | Emits JSON with `isSuccess=true`, deterministic `taskSlug="sync-test"`, non-empty `machineId`, valid relative `runDirectory` and `dbPath`, status `ACTIVE`. Exit 0. |
| **AC-02** | **Subtask Creation with Lineage** | `gitmap agent-ai subtask create sync-test "Worker Task" --code Task-01 --files "cli/a.go" --role "Worker 01" --json` | Emits JSON with `parentTaskId`, `machineId`, `subtaskId`, `taskCode="Task-01"`, `taskSlug`, `ownedFiles=["cli/a.go"]`, status `PENDING`. Exit 0. |
| **AC-03** | **Look-Ahead File Box Claiming** | `gitmap agent subtask claim-files -t sync-test -s Task-01 -f "cli/a.go,cli/b.go" -a "Worker 01" --json` | Inserts into Tier 1 `FileClaim` and Tier 2 `Subtask.OwnedFilesJson`. Emits JSON with `hasCollision=false`. Exit 0. |
| **AC-04** | **Collision Detection & Alerting** | `gitmap agent subtask claim-files -t sync-test -s Task-02 -f "cli/a.go" -a "Worker 02" --json` | Detects overlapping claim on `cli/a.go`. Emits JSON with `hasCollision=true` and conflicting agent details. Renders high-contrast collision warning banner. Exit 0. |
| **AC-05** | **ActionType Enum Enforcement** | `gitmap agent-ai editing Task-01 "cli/a.go" -a "invalid_action" -l 10,20 -r "Testing enum"` | Rejects with error message specifying allowed enum values (`modify`, `analyze`, `read`, `remove`, `refactor`, `create`). Exit code 1. |
| **AC-06** | **Pre-Touch Action Logging Speed** | `gitmap agent-ai editing Task-01 "cli/a.go" -a modify -l 10,20 -r "Refactor parser" --json` | Records action in `AgentActionLog` with status `'IN_PROGRESS'` in $<2\text{ms}$. Emits JSON with `actionLogId`, `actionType="modify"`, `lineRange="10,20"`. Exit 0. |
| **AC-07** | **Heatmap Discussion Speed & Analytics** | `gitmap agent-ai heatmap discuss sync-test --fast --json` | Aggregates heated files, active claims, and contention risk scores in $<5\text{ms}$. Emits JSON containing `files`, `contentionScore`, `contentionLevel`. Exit 0. |
| **AC-08** | **Unified JSON Output Schema Compliance** | Run all commands with `--json` and pipe through schema validator | 100% of outputs conform to documented JSON schemas; valid JSON object envelope with required keys. |
| **AC-09** | **Human ANSI Terminal Formatting** | Run commands without `--json` | Emits clean ANSI tables, colored status badges (`ACTIVE`, `PENDING`, `IN_PROGRESS`), latency indicators, and actionable guidance. |
| **AC-10** | **Strict Relative Git Path Hygiene** | `python 03-ai-scripts/50-fastgate.py` | 0 absolute path violations (`C:\`, `/`), zero `file://` URIs across all database records and emitted payloads. |
| **AC-11** | **Positive Boolean Conventions** | Schema and code inspection | All booleans use positive affirmative prefixes (`isSuccess`, `isActive`, `hasCompleted`, `isBlocked`, `hasCollision`). Zero negative booleans. |
| **AC-12** | **Python Engine & Go CLI Full Parity** | Execute test matrix against both Go binary and `03-ai-scripts/46-agent-sqlite-task-manager.py` | Identical schema DDL, command arguments, JSON response envelopes, and exit codes across both engines. |

---

## 7. Operational Guidelines for Autonomous Agents

1. **Mandatory Turn 1 Task Registration:** When starting a parent task, the lead orchestrator MUST invoke `gitmap agent-ai create-task "<prompt>" --json` in Turn 1 and capture `parentTaskId`, `taskSlug`, and `machineId`.
2. **Subtask Allocation Before Dispatch:** Before spawning subagents via `invoke_subagent`, the lead orchestrator MUST create subtasks via `gitmap agent-ai subtask create` and assign disjoint file boxes.
3. **Mandatory Pre-Touch Rule (Rule R-PRETOUCH):** Subagent workers are strictly forbidden from modifying any file before calling `gitmap agent-ai editing` with a valid `ActionType` enum.
4. **Pre-Flight Heatmap Discussion:** Before resolving contentious files or finalizing multi-file refactors, agents must inspect `gitmap agent-ai heatmap discuss` to verify no peer worker is actively modifying the file box.
5. **Strict Relative Git Paths:** All arguments passed to `--files` and `TargetFile` MUST be relative Git paths with forward slashes. Absolute filesystem paths are prohibited and rejected by AST linters.
