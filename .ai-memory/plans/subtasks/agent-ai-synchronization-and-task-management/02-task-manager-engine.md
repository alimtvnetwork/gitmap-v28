# Subtask 02: Task Manager Engine Updates, ActionType Enum & Heatmap Discuss Engine

> **Subtask ID:** `Subtask-02`  
> **Parent Plan:** `.ai-memory/plans/agent-ai-synchronization-and-task-management.md`  
> **Target Path:** `.ai-memory/plans/subtasks/agent-ai-synchronization-and-task-management/02-task-manager-engine.md`  
> **Status:** `PENDING`  
> **Assigned Agent:** `Engine Worker 02`  
> **Target Script:** `03-ai-scripts/46-agent-sqlite-task-manager.py`  

---

## 1. Objectives & Scope

1. Upgrade `03-ai-scripts/46-agent-sqlite-task-manager.py` to support multi-machine agent coordination by recording deterministic `MachineId` values across task initialization, subtask creation, claiming, and action logging.
2. Standardize JSON responses across all CLI entry points so AI orchestrators receive explicit machine context (`parentTaskId`, `machineId`, `taskSlug`, `runDirectory`, `databasePath`).
3. Enforce the canonical `ActionType` enum (`modify`, `analyze`, `read`, `remove`, `refactor`, `create`) with automated tool alias normalization in `cmd_log_action`, satisfying Rule R-PRETOUCH.
4. Implement the high-speed `<5ms` real-time heatmap and collision discussion engine (`gitmap agent-ai heatmap discuss` / `python 03-ai-scripts/46-agent-sqlite-task-manager.py heatmap discuss`), providing instant file contention metrics and collision warning signals for multi-agent synchronization.

---

## 2. Detailed Technical Deliverables

### 2.1 MachineId Integration in `cmd_init` and Master Ledger Sync
- **Machine Identification:**
  - Invoke `get_machine_id()` to probe host environment (`$GITMAP_MACHINE_ID`, `socket.gethostname()`, `$COMPUTERNAME`, `$HOSTNAME`).
- **Database Stamping:**
  - When creating a new run or resuming an existing run, stamp `MachineId` into `ParentTask`.
  - Register the run in the Tier 1 master ledger (`.ai-memory/temp-agents/ai_agents.db`) under table `AgentRunCatalog` and update `MachineRegistry` heartbeat.
- **Enhanced JSON Output Response:**
  ```json
  {
    "action": "INITIALIZED",
    "parentTaskId": 1,
    "slug": "agent-ai-synchronization-and-task-management",
    "machineId": "devbox-alpha",
    "runDirectory": ".ai-memory/temp-agents/01-agent-ai-synchronization-and-task-management",
    "databasePath": ".ai-memory/temp-agents/01-agent-ai-synchronization-and-task-management/agent-task.db",
    "totalStepsBudget": 300,
    "status": "ACTIVE"
  }
  ```

### 2.2 MachineId Context in `cmd_add_subtasks`
- **Subtask Insertion Updates:**
  - Stamp `MachineId` for each inserted subtask row in `Subtask`.
  - Retain disjoint file ownership validation across subtasks in incoming payload.
- **Structured JSON Output:**
  ```json
  {
    "status": "SUCCESS",
    "parentTaskId": 1,
    "machineId": "devbox-alpha",
    "subtasksCreated": 4,
    "subtasks": [
      {
        "subtaskId": 1,
        "code": "Task-01",
        "title": "Author spec and schema migrations",
        "machineId": "devbox-alpha",
        "assignedAgentRole": "Worker 01",
        "ownedFiles": ["02-spec/21-app/agent-ai-synchronization-and-task-management/01-architecture-spec.md"]
      }
    ]
  }
  ```

### 2.3 Canonical `ActionType` Enum & Alias Normalization in `cmd_log_action`
- **Canonical Action Types:**
  - Mutation: `modify`, `create`, `remove`, `refactor`
  - Inspection: `analyze`, `read`
  - Lifecycle: `START`, `COMPLETE`, `FAIL`
- **Tool Alias Normalization Engine:**
  ```python
  ACTION_TYPE_ALIASES = {
      # In-place mutations
      "replace_file_content": "modify",
      "edit_file": "modify",
      "edit": "modify",
      "patch": "modify",
      "update": "modify",
      # Reading & inspection
      "view_file": "read",
      "read_file": "read",
      "cat": "read",
      "gitmap cat": "read",
      "read_resource": "read",
      # Analysis & linting
      "analyze_file": "analyze",
      "search_web": "analyze",
      "gitmap aum search": "analyze",
      "gitmap find": "analyze",
      "run_linter": "analyze",
      "run_test": "analyze",
      "check": "analyze",
      # Creation & Deletion
      "create_file": "create",
      "touch": "create",
      "new": "create",
      "delete_file": "remove",
      "rm": "remove",
      "unlink": "remove",
      "purge": "remove",
      # Refactoring
      "refactor_module": "refactor",
      "extract_type": "refactor",
  }
  ```
- **Validation & Enforcement:**
  - If `--action` matches a canonical enum or known alias, normalize it.
  - If `--action == "write_to_file"`:
    - Check if target file exists on disk.
    - If file exists -> normalize to `modify`.
    - If file does not exist -> normalize to `create`.
  - If `--action` is unrecognized:
    - Reject execution with exit code `1` and emit error JSON detailing accepted canonical types and valid aliases.
- **Rule R-PRETOUCH Logging Execution:**
  - Stamping: `MachineId = get_machine_id()`.
  - Target file validation: verify `--file` is a strictly relative Git path.
  - Insert record into `AgentActionLog` with `Status = 'IN_PROGRESS'`.

### 2.4 High-Speed Real-Time Heatmap Engine (`cmd_heatmap_discuss`)
- **Query Architecture:**
  - Query indexed `AgentActionLog` joined with `Subtask` for actions within `--window` (default: 1800s / 30 minutes).
  - Calculate per-file metrics:
    - `touchCount`: total action entries.
    - `uniqueAgents`: count of distinct `AgentRole` values.
    - `mutations`: count of `ActionType IN ('modify', 'create', 'remove', 'refactor')`.
    - `inspections`: count of `ActionType IN ('read', 'analyze')`.
    - `contentionScore`: $(mutations \times 2.0) + (inspections \times 0.5) + (uniqueAgents > 1 ? 5.0 : 0.0)$.
  - Evaluate collision warning:
    `collisionWarning = True` if any file has $contentionScore \ge 10.0$ or ($uniqueAgents \ge 2$ and $mutations > 0$).
- **Dual Output Rendering:**
  1. **Terminal Discussion Mode:** ANSI colored table with severity badges (`COLLISION_RISK` in red, `WARM` in yellow, `COOL` in green) and discussion talking points.
  2. **AI JSON Mode (`--json`):** Clean, structured JSON with machine ID, active agents, hotspots, and suggested arbitration actions.
- **Performance Budget:** Execution time must complete in $<5\text{ ms}$ over typical run database sizes.

### 2.5 CLI Interface Extensions
Register subcommands and flags in `03-ai-scripts/46-agent-sqlite-task-manager.py`:
```bash
# Heatmap discussion command
python 03-ai-scripts/46-agent-sqlite-task-manager.py heatmap discuss [--db <path>] [--window <seconds>] [--json]

# Alias subcommand
python 03-ai-scripts/46-agent-sqlite-task-manager.py heatmap [--db <path>] [--window <seconds>] [--json]
```

---

## 3. Acceptance Criteria & Invariants

1. **Deterministic Machine Context:** `init`, `add-subtasks`, `claim`, and `log-action` consistently write and return the host's `MachineId`.
2. **ActionType Strictness:** Passing arbitrary unknown action strings (e.g. `--action "random_string"`) fails with validation errors. Passing tool aliases (e.g. `--action "replace_file_content"`) succeeds and records `modify`.
3. **Pre-Touch Integrity (Rule R-PRETOUCH):** `log-action` writes `Status = 'IN_PROGRESS'` to `AgentActionLog` before tool dispatch.
4. **Heatmap Execution Speed:** `heatmap discuss` returns valid analysis in $<5\text{ ms}$ for databases with up to 5,000 action records.
5. **Collision Alert Accuracy:** When two distinct agents log mutation actions against the same file within the window, `collisionWarning` evaluates to `true`.
6. **Path Purity:** All logged files and returned JSON objects use strict forward-slash relative Git paths without Windows drive letters or `file:///` prefixes.

---

## 4. Verification Steps & Test Commands

1. **Verify MachineId in Task Initialization:**
   ```bash
   python 03-ai-scripts/46-agent-sqlite-task-manager.py init --name "verify machine id telemetry"
   ```
   *Expected Result:* Output JSON contains `"machineId"`.

2. **Verify Subtask Creation with Machine Context:**
   ```bash
   python 03-ai-scripts/46-agent-sqlite-task-manager.py add-subtasks --db <db> --tasks-json '[{"title":"Test Subtask","owned_files":["02-spec/test.md"]}]'
   ```
   *Expected Result:* Output JSON contains `"parentTaskId"` and `"machineId"`.

3. **Verify ActionType Normalization:**
   ```bash
   python 03-ai-scripts/46-agent-sqlite-task-manager.py log-action --db <db> --subtask-id 1 --agent "Worker 01" --action "replace_file_content" --file "02-spec/test.md"
   ```
   *Expected Result:* Action is accepted; querying SQLite returns `ActionType = 'modify'`.

4. **Verify ActionType Validation Error on Invalid Action:**
   ```bash
   python 03-ai-scripts/46-agent-sqlite-task-manager.py log-action --db <db> --subtask-id 1 --agent "Worker 01" --action "invalid_action_xyz" --file "02-spec/test.md"
   ```
   *Expected Result:* Exits with code `1` and emits error listing valid enums.

5. **Verify Heatmap Discuss JSON & Latency:**
   ```bash
   python 03-ai-scripts/46-agent-sqlite-task-manager.py heatmap discuss --db <db> --json
   ```
   *Expected Result:* Returns JSON payload with `hotspotFiles`, `contentionScore`, and `collisionWarning`.
