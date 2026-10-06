# 06-database-and-split-db: Split-DB Components, Migration Engine & Agent Task DB Specification

- **Spec ID:** `06-database-and-split-db/02-component-spec.md`
- **Status:** `APPROVED`
- **Version:** `1.0.0`
- **Subsystem:** SQLite Connection Pool, Migration Engine, Split-DB Repositories, AI Agent Tasks
- **Dependencies:** `cli/db`, `cli/repodb`, `cli/pipelinedb`, `cli/cmdagent`
- **Version Baseline:** `v6.498.0`
- **Target Version:** `v6.499.0`

---

## 1. Component Breakdown

The Database and Split-DB cluster encompasses four primary operational components:

```
06-database-and-split-db/
├── 01-architecture-spec.md
└── 02-component-spec.md
```

| Component | Target Source Path | Primary Role |
| :--- | :--- | :--- |
| **Split-DB Pool Manager** | `cli/db/pool.go`, `cli/db/sqlite.go` | Three-tier connection initialization, WAL mode enforcement, single-writer governance. |
| **Schema Migration Engine** | `cli/db/migrations.go` | Deterministic schema evolution, table creation, index management. |
| **Agent Task DB Subsystem** | `cli/cmdagent/task_db.go`, `cli/cmdagent/agent_db.go` | Multi-agent task tracking, worker heartbeat, action logging, crash forensics. |
| **Diagnostics & DB CLI** | `cli/cmddb/db.go`, `cli/cmddb/reset.go` | `gitmap db list`, `gitmap db reset`, diagnostic query inspection. |

---

## 2. Connection Pool & SQLite Lifecycle

### 2.1 Driver & PRAGMA Initialization
Every database connection initialized via `db.OpenSQLite()` executes mandatory PRAGMAs:

```go
package db

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func OpenSQLite(path string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", path+"?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=on")
    if err != nil {
        return nil, err
    }
    // Strict concurrency governance: 1 writer to prevent SQLite locking panics
    db.SetMaxOpenConns(1)
    db.SetMaxIdleConns(1)
    return db, nil
}
```

---

## 3. AI Agent Task Orchestrator Split-DB

### 3.1 Task Hierarchy & Flow
1. **Root Master DB (`ai_agents.db`):** Records parent tasks, agent registration, global status.
2. **Task Root DB (`agent-task.db`):** Tracks subtask allocations, dependency graphs, completion proofs.
3. **Worker Split DB (`worker-XX.db`):** Dedicated worker action logs (file reads, searches, writes, execution ms).

### 3.2 Action Log Table
```sql
CREATE TABLE IF NOT EXISTS AgentActionLog (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    agentId TEXT NOT NULL,
    subtaskId TEXT NOT NULL,
    actionType TEXT NOT NULL, -- 'READ', 'WRITE', 'SEARCH', 'COMMAND'
    targetPath TEXT,
    startLine INTEGER,
    endLine INTEGER,
    durationMs INTEGER,
    status TEXT NOT NULL,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

---

## 4. Database Reset Engine & Foreign Key Protection

When `gitmap db reset` is executed:
- Acquires exclusive write lock.
- Temporarily sets `PRAGMA foreign_keys = OFF;`.
- Truncates child tables (`DetectedProject`, `PullError`, `Bookmark`) before parent tables (`Repo`, `Node`).
- Re-enables `PRAGMA foreign_keys = ON;`.
- Re-runs initial schema migrations cleanly.

---

## 5. Verification & Acceptance Criteria

```yaml
verificationGates:
  isPragmasConfiguredCorrectly: true
  isAgentTaskDbHierarchyOperational: true
  isDbResetSafeAndCascaded: true
  isPositiveBooleansUsed: true
```

- [x] WAL mode and busy timeout active on all three database tiers.
- [x] Agent telemetry records granular action logs with line ranges.
- [x] Database reset executes without foreign key integrity violations.
