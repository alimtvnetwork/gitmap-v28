# Subtask 01: Specification Authoring, Two-Tier Split-DB & MachineId Schema Migration

> **Subtask ID:** `Subtask-01`  
> **Parent Plan:** `.ai-memory/plans/agent-ai-synchronization-and-task-management.md`  
> **Target Path:** `.ai-memory/plans/subtasks/agent-ai-synchronization-and-task-management/01-spec-and-schema.md`  
> **Status:** `PENDING`  
> **Assigned Agent:** `Spec Worker 01`  
> **Owned Files:**  
> - `02-spec/21-app/agent-ai-synchronization-and-task-management/01-architecture-spec.md`  
> - `.ai-memory/plans/subtasks/agent-ai-synchronization-and-task-management/01-spec-and-schema.md`  

---

## 1. Objectives & Scope

1. Author comprehensive, production-grade architecture specification (`02-spec/21-app/agent-ai-synchronization-and-task-management/01-architecture-spec.md`) documenting the Two-Tier SQLite Split-DB topology, `MachineId` tracking, canonical `ActionType` enum, Rule R-PRETOUCH, crash forensics, and real-time heatmap architecture.
2. Formulate SQLite schema definitions and non-destructive migration pathways for `03-ai-scripts/46-agent-sqlite-task-manager.py` to incorporate `MachineId` across `ParentTask`, `Subtask`, and `AgentActionLog`.
3. Specify the master ledger schema for Tier 1 (`.ai-memory/temp-agents/ai_agents.db`) tracking multi-machine runs, active agents, and cross-machine conflict resolution.
4. Define deterministic `MachineId` generation algorithms using hardware host probe and environment fallbacks.
5. Guarantee 100% adherence to repository conventions: PascalCase naming, positive booleans (`IsActive`, `HasCompleted`, `IsBlocked`), WAL mode, and strict relative Git paths.

---

## 2. Detailed Technical Deliverables

### 2.1 Specification Documents
- **Architecture Spec (`02-spec/21-app/agent-ai-synchronization-and-task-management/01-architecture-spec.md`):**
  - Section on Two-Tier Split-DB architecture (`ai_agents.db` + per-task `agent-task.db`).
  - Section on deterministic `MachineId` generation and cross-machine conflict resolution.
  - Section on canonical `ActionType` enum (`modify`, `analyze`, `read`, `remove`, `refactor`, `create`) and alias normalization.
  - Section on Rule R-PRETOUCH pre-touch logging and crash forensic diagnosis.
  - Section on `<5ms` real-time heatmap architecture (`gitmap agent-ai heatmap discuss`).
  - Strict relative Git paths only; zero absolute paths.

### 2.2 Two-Tier SQLite Schema Architecture
- **Tier 1 Master Ledger (`.ai-memory/temp-agents/ai_agents.db`):**
  - Table `MachineRegistry`:
    - `MachineRegistryId INTEGER PRIMARY KEY AUTOINCREMENT`
    - `MachineId TEXT NOT NULL UNIQUE`
    - `Hostname TEXT NOT NULL`
    - `Platform TEXT NOT NULL`
    - `IsOnline INTEGER NOT NULL DEFAULT 1`
    - `LastHeartbeatAt TEXT NOT NULL`
    - `CreatedAt TEXT NOT NULL`
  - Table `AgentRunCatalog`:
    - `RunCatalogId INTEGER PRIMARY KEY AUTOINCREMENT`
    - `TaskSlug TEXT NOT NULL UNIQUE`
    - `TaskName TEXT NOT NULL`
    - `MachineId TEXT NOT NULL`
    - `RunDirectory TEXT NOT NULL`
    - `DatabasePath TEXT NOT NULL`
    - `Status TEXT NOT NULL DEFAULT 'ACTIVE'`
    - `IsActive INTEGER NOT NULL DEFAULT 1`
    - `HasCompleted INTEGER NOT NULL DEFAULT 0`
    - `TotalSubtasks INTEGER NOT NULL DEFAULT 0`
    - `CompletedSubtasks INTEGER NOT NULL DEFAULT 0`
    - `CreatedAt TEXT NOT NULL`
    - `UpdatedAt TEXT NOT NULL`

- **Tier 2 Task Store Schema Upgrades (`agent-task.db`):**
  - Add `MachineId TEXT NOT NULL` to table `ParentTask`.
  - Add `MachineId TEXT NOT NULL` to table `Subtask`.
  - Add `MachineId TEXT NOT NULL` to table `AgentActionLog`.
  - Add indexes:
    - `idx_subtask_machine ON Subtask(MachineId)`
    - `idx_action_created ON AgentActionLog(CreatedAt)`
    - `idx_action_file ON AgentActionLog(TargetFile)`

### 2.3 Non-Destructive Auto-Migration Logic
In `03-ai-scripts/46-agent-sqlite-task-manager.py`, implement an `ensure_schema_migrations(conn)` routine:
1. Query table column definitions via `PRAGMA table_info(ParentTask);`.
2. If column `MachineId` is absent in `ParentTask`:
   - Execute `ALTER TABLE ParentTask ADD COLUMN MachineId TEXT NOT NULL DEFAULT 'devbox-local';`.
3. If column `MachineId` is absent in `Subtask`:
   - Execute `ALTER TABLE Subtask ADD COLUMN MachineId TEXT NOT NULL DEFAULT 'devbox-local';`.
4. If column `MachineId` is absent in `AgentActionLog`:
   - Execute `ALTER TABLE AgentActionLog ADD COLUMN MachineId TEXT NOT NULL DEFAULT 'devbox-local';`.
5. Execute index creation statements with `IF NOT EXISTS`.

### 2.4 Deterministic Machine ID Detection Helper
Implement `get_machine_id() -> str` in `03-ai-scripts/46-agent-sqlite-task-manager.py`:
```python
def get_machine_id() -> str:
    env_id = os.environ.get("GITMAP_MACHINE_ID", "").strip()
    if env_id:
        return re.sub(r"[^a-z0-9_-]", "-", env_id.lower())[:32]
    try:
        raw = socket.gethostname().strip().lower()
    except Exception:
        raw = os.environ.get("COMPUTERNAME", os.environ.get("HOSTNAME", "devbox-local")).strip().lower()
    cleaned = re.sub(r"[^a-z0-9_-]", "-", raw).strip("-")
    return cleaned[:32] or "devbox-local"
```

---

## 3. Acceptance Criteria & Invariants

1. **Schema Integrity:** Running `python 03-ai-scripts/46-agent-sqlite-task-manager.py schema --ddl` prints valid SQL containing `MachineId` on `ParentTask`, `Subtask`, and `AgentActionLog`.
2. **Backward Compatibility:** Existing legacy SQLite databases without `MachineId` open cleanly and are automatically migrated with zero data loss.
3. **Deterministic Output:** Running `init` on the same machine produces consistent `machineId` values matching `get_machine_id()`.
4. **Positive Booleans:** Schema exclusively employs positive naming conventions (`IsActive`, `HasCompleted`, `IsBlocked`) storing integer `0` or `1`.
5. **Path Hygiene:** `SCHEMA_METADATA` and validation functions strictly enforce relative Git forward-slash paths (`02-spec/...`, `.ai-memory/...`). Zero absolute paths or `file:///` URIs permitted.

---

## 4. Verification Steps & Test Commands

1. **Validate Schema DDL Output:**
   ```bash
   python 03-ai-scripts/46-agent-sqlite-task-manager.py schema --ddl
   ```
   *Expected Output:* SQL script defines `MachineId TEXT NOT NULL` in `ParentTask`, `Subtask`, and `AgentActionLog`.

2. **Validate Schema JSON Metadata:**
   ```bash
   python 03-ai-scripts/46-agent-sqlite-task-manager.py schema --json
   ```
   *Expected Output:* Valid JSON listing `MachineId` in column metadata arrays for each table.

3. **Verify Database Initialization with MachineId:**
   ```bash
   python 03-ai-scripts/46-agent-sqlite-task-manager.py init --name "test machine id schema integration"
   ```
   *Expected Output:* Output JSON contains `"machineId"` field populated with probed host name.

4. **Verify Schema Auto-Migration on Legacy DB Fixture:**
   Execute synthetic script creating legacy DB schema without `MachineId`, connect with updated engine, verify columns are added without error.
