# Database Specification: Multi-AI Heatmap Telemetry & Split-DB Storage Engine

- **Feature Slug:** `heatmap-command-and-database-structure-instructions`
- **Module:** `cli/store`, `cli/heatmap`, `cli/cmd`
- **Specification Version:** 2.0.0
- **Status:** APPROVED
- **Target Audience:** Multi-Agent Workers, Lead Orchestrators, Systems Engineers, Database Architects
- **Target Release:** Minor Bump

---

## 1. Executive Summary & Architectural Motivation

The GitMap AI Agent Telemetry & Heatmap Engine captures granular cognitive traces, pre-touch file editing declarations, and execution hierarchies produced by autonomous AI agents (such as Antigravity, Muse, Cloud Code, Codex, Grok, and Cursor). To ensure telemetry persists across multiple machines and distinct AI agent frameworks without corrupting shared working trees, GitMap implements a **Three-Tier Distributed SQLite Split-DB Architecture**:

1. **Root Heatmap Database (`root-heatmap.db`):**
   Located inside the synchronized repository cache at `.heatmap/01-<ai-name>/root-heatmap.db`. Acts as the central registry for agent identities, repository bindings, parent tasks, and commit provenance.
2. **Repository-Specific Split Database (`heatmap.db`):**
   Located inside the synchronized repository cache at `.heatmap/01-<ai-name>/<repo-slug>/heatmap.db`. Houses high-frequency, append-heavy operational traces, file touch line numbers, code snippets, reasoning state transitions, and execution tree nodes.
3. **Local Machine Tracking Database (`<dataDir>/heatmap.db`):**
   Located on the local developer or worker machine alongside the CLI root database (resolved via `store.EffectiveDataDir()`). Tracks local commit SHAs, transaction timestamps, uncommitted mutation counters, and synchronization states to drive conflict resolution.

All databases enforce **PascalCase** table naming, **positive boolean flags** (`isActive`, `hasCompleted`, `isExpectationMet`, `isDirty`), **SQLite WAL mode** (`PRAGMA journal_mode = WAL`), and structured application errors (`apperror.AppError`).

---

## 2. Distributed Filesystem Topology & Database Paths

### 2.1 Repository Cache (`repo-cache`) Remote Hierarchy

The `repo-cache` (`rc`) companion repository maintains isolated folders per AI sequence name to avoid cross-agent write contention:

```text
<repo-cache>/
└── .heatmap/
    ├── 01-antigravity/
    │   ├── root-heatmap.db                   <-- Tier 1: Agent Sessions, Repos, Parent Tasks, Commits
    │   ├── gitmap/
    │   │   └── heatmap.db                   <-- Tier 2: Granular Actions & Execution Tree for 'gitmap'
    │   └── movie-cli/
    │       └── heatmap.db                   <-- Tier 2: Granular Actions & Execution Tree for 'movie-cli'
    ├── 02-muse/
    │   ├── root-heatmap.db
    │   └── gitmap/
    │       └── heatmap.db
    ├── 03-cloud-code/
    ├── 04-codex/
    ├── 05-grok/
    └── 06-cursor/
```

### 2.2 Local Machine State Tracking Hierarchy

On the local machine, the CLI runtime maintains its local tracking database within the user's configuration and cache store:

```text
<dataDir>/                                   <-- Resolved by store.EffectiveDataDir()
├── gitmap.db                                <-- Global CLI Configuration & State
├── installation.db                          <-- Local Machine Setup & Dependency State
└── heatmap.db                               <-- Tier 3: Local Machine Git Commit & Delta Tracking
```

### 2.3 Strict Relative Path Resolution Matrix

| Logical Database | Physical Relative Path | Scope & Sync Model |
|---|---|---|
| **Root Heatmap DB** | `.heatmap/01-<ai-name>/root-heatmap.db` | Stored in `repo-cache`, committed to Git, pushed upstream |
| **Repo Split DB** | `.heatmap/01-<ai-name>/<repo-slug>/heatmap.db` | Stored in `repo-cache`, committed to Git, pushed upstream |
| **Local Machine Tracking** | `<dataDir>/heatmap.db` | Local machine only, excluded from Git, machine-scoped |
| **Temporary Sync Backup** | `.heatmap/temp-sync-<timestamp>.db` | Ephemeral SQLite export created during conflict resolution |

---

## 3. Root Heatmap Database Schema (`root-heatmap.db`)

The Root Heatmap Database manages global metadata, agent sessions, registered repositories, parent task envelopes, and Git commit checkpoints.

### 3.1 Table: `HeatmapAiSession`

Captures the lifecycle, runtime configuration, and operating environment of each active AI agent execution session.

```sql
CREATE TABLE IF NOT EXISTS HeatmapAiSession (
    HeatmapAiSessionId   INTEGER PRIMARY KEY AUTOINCREMENT,
    SessionUuid          TEXT NOT NULL UNIQUE,
    AiName               TEXT NOT NULL,
    SequencePrefix       TEXT NOT NULL,
    MachineHostname      TEXT NOT NULL DEFAULT '',
    OsPlatform           TEXT NOT NULL DEFAULT '',
    AgentRole            TEXT NOT NULL DEFAULT 'ORCHESTRATOR',
    StartedAt            TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CompletedAt          TEXT NULL,
    IsActive             INTEGER NOT NULL DEFAULT 1,
    HasCompleted         INTEGER NOT NULL DEFAULT 0,
    TotalTasksExecuted   INTEGER NOT NULL DEFAULT 0,
    TotalActionsLogged   INTEGER NOT NULL DEFAULT 0,
    Notes                TEXT NULL,
    Comments             TEXT NULL,
    CreatedAt            TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt            TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS IdxHeatmapAiSession_Uuid ON HeatmapAiSession(SessionUuid);
CREATE INDEX IF NOT EXISTS IdxHeatmapAiSession_AiName ON HeatmapAiSession(AiName);
CREATE INDEX IF NOT EXISTS IdxHeatmapAiSession_IsActive ON HeatmapAiSession(IsActive);
```

#### Field Specifications:
- `HeatmapAiSessionId`: Monotonically increasing primary key.
- `SessionUuid`: Unique identifier (UUIDv4) generated at agent invocation.
- `AiName`: Canonical lowercase AI name (`antigravity`, `muse`, `cloud-code`, `codex`, `grok`, `cursor`).
- `SequencePrefix`: Two-digit zero-padded sequence prefix (`01`, `02`, `03`, `04`, `05`, `06`).
- `IsActive`: Positive boolean flag (`1` = session running, `0` = terminated).
- `HasCompleted`: Positive boolean flag (`1` = normal termination, `0` = in-flight or crashed).

---

### 3.2 Table: `HeatmapRepoRegistry`

Tracks all repositories that have telemetry captured under this AI agent sequence.

```sql
CREATE TABLE IF NOT EXISTS HeatmapRepoRegistry (
    HeatmapRepoRegistryId INTEGER PRIMARY KEY AUTOINCREMENT,
    RepoSlug              TEXT NOT NULL UNIQUE,
    RepoUrl               TEXT NOT NULL DEFAULT '',
    DefaultBranch         TEXT NOT NULL DEFAULT 'main',
    FirstSeenAt           TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    LastActiveAt          TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    IsActive              INTEGER NOT NULL DEFAULT 1,
    TotalActionsRecorded  INTEGER NOT NULL DEFAULT 0,
    Notes                 TEXT NULL,
    Comments              TEXT NULL,
    CreatedAt             TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt             TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS IdxHeatmapRepoRegistry_Slug ON HeatmapRepoRegistry(RepoSlug);
CREATE INDEX IF NOT EXISTS IdxHeatmapRepoRegistry_IsActive ON HeatmapRepoRegistry(IsActive);
```

---

### 3.3 Table: `HeatmapParentTask`

Indexes parent tasks across repositories, linking goals, step budgets, and overall completion states.

```sql
CREATE TABLE IF NOT EXISTS HeatmapParentTask (
    HeatmapParentTaskId INTEGER PRIMARY KEY AUTOINCREMENT,
    TaskSlug            TEXT NOT NULL UNIQUE,
    TaskName            TEXT NOT NULL,
    RepoSlug            TEXT NOT NULL,
    Status              TEXT NOT NULL DEFAULT 'IN_PROGRESS',
    StepBudget          INTEGER NOT NULL DEFAULT 300,
    CompletedSteps      INTEGER NOT NULL DEFAULT 0,
    HasCompleted        INTEGER NOT NULL DEFAULT 0,
    IsActive            INTEGER NOT NULL DEFAULT 1,
    StartedAt           TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CompletedAt         TEXT NULL,
    Notes               TEXT NULL,
    Comments            TEXT NULL,
    CreatedAt           TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt           TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (RepoSlug) REFERENCES HeatmapRepoRegistry(RepoSlug) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS IdxHeatmapParentTask_Slug ON HeatmapParentTask(TaskSlug);
CREATE INDEX IF NOT EXISTS IdxHeatmapParentTask_RepoSlug ON HeatmapParentTask(RepoSlug);
CREATE INDEX IF NOT EXISTS IdxHeatmapParentTask_Status ON HeatmapParentTask(Status);
CREATE INDEX IF NOT EXISTS IdxHeatmapParentTask_IsActive ON HeatmapParentTask(IsActive);
```

---

### 3.4 Table: `HeatmapCommitLog`

Records each Git commit performed on `repo-cache` for auditability, lineage tracing, and conflict detection.

```sql
CREATE TABLE IF NOT EXISTS HeatmapCommitLog (
    HeatmapCommitLogId          INTEGER PRIMARY KEY AUTOINCREMENT,
    TaskSlug                    TEXT NOT NULL,
    RepoSlug                    TEXT NOT NULL,
    CommitSha                   TEXT NOT NULL,
    CommitMessage               TEXT NOT NULL,
    FilesCommittedCount         INTEGER NOT NULL DEFAULT 0,
    TransactionsCommittedCount  INTEGER NOT NULL DEFAULT 0,
    CommittedAt                 TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PushedSuccessfully          INTEGER NOT NULL DEFAULT 1,
    RemoteOrigin                TEXT NOT NULL DEFAULT '',
    Notes                       TEXT NULL,
    Comments                    TEXT NULL,
    CreatedAt                   TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (TaskSlug) REFERENCES HeatmapParentTask(TaskSlug) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS IdxHeatmapCommitLog_TaskSlug ON HeatmapCommitLog(TaskSlug);
CREATE INDEX IF NOT EXISTS IdxHeatmapCommitLog_RepoSlug ON HeatmapCommitLog(RepoSlug);
CREATE INDEX IF NOT EXISTS IdxHeatmapCommitLog_CommitSha ON HeatmapCommitLog(CommitSha);
CREATE INDEX IF NOT EXISTS IdxHeatmapCommitLog_CommittedAt ON HeatmapCommitLog(CommittedAt);
```

---

## 4. Repository Split Database Schema (`heatmap.db`)

Each target repository receives a dedicated SQLite database located at `.heatmap/01-<ai-name>/<repo-slug>/heatmap.db`. This database stores high-volume action events and hierarchical execution trees.

### 4.1 Table: `HeatmapReasoningState` (Reference Table)

Pre-seeded domain reference table classifying agent cognitive phases.

```sql
CREATE TABLE IF NOT EXISTS HeatmapReasoningState (
    HeatmapReasoningStateId INTEGER PRIMARY KEY AUTOINCREMENT,
    StateCode               TEXT NOT NULL UNIQUE,
    Description             TEXT NOT NULL,
    Notes                   TEXT NULL,
    IsActive                INTEGER NOT NULL DEFAULT 1,
    CreatedAt               TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS IdxHeatmapReasoningState_Code ON HeatmapReasoningState(StateCode);

INSERT OR IGNORE INTO HeatmapReasoningState (StateCode, Description, Notes, IsActive) VALUES
('EXPLORATION', 'Initial exploration, codebase discovery, and file inspection', NULL, 1),
('SPEC_REFINEMENT', 'Spec review, requirement alignment, and task decomposition', NULL, 1),
('CODE_EDITING', 'Direct file modification, creation, and refactoring', NULL, 1),
('VERIFICATION', 'Reviewing diffs, linting checks, and quality validation', NULL, 1),
('REMEDIATION', 'Bug fixing, error resolution, and exception handling', NULL, 1),
('COMPLETED', 'Step or subtask completed successfully', NULL, 1);
```

---

### 4.2 Table: `HeatmapActionLog`

Captures granular pre-touch file declarations, line ranges, code snippets, reasoning statements, and tool execution outcomes.

```sql
CREATE TABLE IF NOT EXISTS HeatmapActionLog (
    HeatmapActionLogId INTEGER PRIMARY KEY AUTOINCREMENT,
    ParentLogId        INTEGER NULL,
    TaskSlug           TEXT NOT NULL,
    RepoSlug           TEXT NOT NULL DEFAULT '',
    ActionType         TEXT NOT NULL,
    FilePath           TEXT NOT NULL DEFAULT '',
    StartLine          INTEGER NOT NULL DEFAULT 0,
    EndLine            INTEGER NOT NULL DEFAULT 0,
    CodeSnippet        TEXT NULL,
    CommandText        TEXT NULL,
    ReasoningState     TEXT NOT NULL DEFAULT 'EXPLORATION',
    ReasoningNotes     TEXT NULL,
    ExpectedOutcome    TEXT NULL,
    ActualOutcome      TEXT NULL,
    IsExpectationMet   INTEGER NOT NULL DEFAULT 1,
    ExecutionDepth     INTEGER NOT NULL DEFAULT 0,
    ExecutionOrder     INTEGER NOT NULL DEFAULT 1,
    DurationMs         INTEGER NOT NULL DEFAULT 0,
    Notes              TEXT NULL,
    Comments           TEXT NULL,
    CreatedAt          TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ParentLogId) REFERENCES HeatmapActionLog(HeatmapActionLogId) ON DELETE SET NULL,
    FOREIGN KEY (ReasoningState) REFERENCES HeatmapReasoningState(StateCode)
);

CREATE INDEX IF NOT EXISTS IdxHeatmapActionLog_TaskSlug ON HeatmapActionLog(TaskSlug);
CREATE INDEX IF NOT EXISTS IdxHeatmapActionLog_FilePath ON HeatmapActionLog(FilePath);
CREATE INDEX IF NOT EXISTS IdxHeatmapActionLog_ParentId ON HeatmapActionLog(ParentLogId);
CREATE INDEX IF NOT EXISTS IdxHeatmapActionLog_ActionType ON HeatmapActionLog(ActionType);
CREATE INDEX IF NOT EXISTS IdxHeatmapActionLog_ReasoningState ON HeatmapActionLog(ReasoningState);
CREATE INDEX IF NOT EXISTS IdxHeatmapActionLog_CreatedAt ON HeatmapActionLog(CreatedAt);
```

---

### 4.3 Table: `HeatmapExecutionTree`

Maintains the explicit parent-child node hierarchy for structured tree inspection, child traversal, and subtask nesting.

```sql
CREATE TABLE IF NOT EXISTS HeatmapExecutionTree (
    HeatmapExecutionTreeId INTEGER PRIMARY KEY AUTOINCREMENT,
    RootLogId              INTEGER NOT NULL,
    NodeLogId              INTEGER NOT NULL UNIQUE,
    ParentLogId            INTEGER NULL,
    NodeDepth              INTEGER NOT NULL DEFAULT 0,
    SiblingOrder           INTEGER NOT NULL DEFAULT 1,
    PathVector             TEXT NOT NULL DEFAULT '/',
    IsLeaf                 INTEGER NOT NULL DEFAULT 1,
    HasChildren            INTEGER NOT NULL DEFAULT 0,
    CreatedAt              TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt              TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (RootLogId) REFERENCES HeatmapActionLog(HeatmapActionLogId) ON DELETE CASCADE,
    FOREIGN KEY (NodeLogId) REFERENCES HeatmapActionLog(HeatmapActionLogId) ON DELETE CASCADE,
    FOREIGN KEY (ParentLogId) REFERENCES HeatmapActionLog(HeatmapActionLogId) ON DELETE SET NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS IdxHeatmapExecutionTree_Node ON HeatmapExecutionTree(NodeLogId);
CREATE INDEX IF NOT EXISTS IdxHeatmapExecutionTree_Root ON HeatmapExecutionTree(RootLogId);
CREATE INDEX IF NOT EXISTS IdxHeatmapExecutionTree_Parent ON HeatmapExecutionTree(ParentLogId);
CREATE INDEX IF NOT EXISTS IdxHeatmapExecutionTree_Path ON HeatmapExecutionTree(PathVector);
```

---

## 5. Local Machine Tracking Database Schema (`<dataDir>/heatmap.db`)

The local machine tracking database resides close to the CLI root database (`<dataDir>/heatmap.db`). It records local git checkpoints, dirty transaction metrics, and conflict synchronization states across local invocations.

### 5.1 Table: `HeatmapLocalTracking`

```sql
CREATE TABLE IF NOT EXISTS HeatmapLocalTracking (
    HeatmapLocalTrackingId     INTEGER PRIMARY KEY AUTOINCREMENT,
    RepoSlug                   TEXT NOT NULL,
    AiName                     TEXT NOT NULL,
    SequencePrefix             TEXT NOT NULL DEFAULT '01',
    RepoCachePath              TEXT NOT NULL DEFAULT '',
    LastCommittedSha           TEXT NOT NULL DEFAULT '',
    LastCommitTimestamp        TEXT NOT NULL DEFAULT '1970-01-01T00:00:00Z',
    UncommittedChangeCount     INTEGER NOT NULL DEFAULT 0,
    SyncStatus                 TEXT NOT NULL DEFAULT 'IDLE',
    IsDirty                    INTEGER NOT NULL DEFAULT 0,
    LastSyncAttemptAt          TEXT NULL,
    LastSuccessfulSyncAt       TEXT NULL,
    ErrorMessage               TEXT NULL,
    CreatedAt                  TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt                  TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT UqHeatmapLocalTracking_Repo_AI UNIQUE (RepoSlug, AiName)
);

CREATE UNIQUE INDEX IF NOT EXISTS IdxHeatmapLocalTracking_Repo_AI ON HeatmapLocalTracking(RepoSlug, AiName);
CREATE INDEX IF NOT EXISTS IdxHeatmapLocalTracking_SyncStatus ON HeatmapLocalTracking(SyncStatus);
CREATE INDEX IF NOT EXISTS IdxHeatmapLocalTracking_IsDirty ON HeatmapLocalTracking(IsDirty);
```

#### Field Specifications:
- `RepoSlug`: Identifier of the tracked workspace repo (e.g. `gitmap`).
- `AiName`: Canonical lowercase AI agent name (e.g. `antigravity`).
- `SequencePrefix`: Two-digit prefix matching folder (e.g. `01`).
- `LastCommittedSha`: 40-character Git commit hash of the last successful `gitmap cpf` on `repo-cache`.
- `LastCommitTimestamp`: ISO 8601 UTC timestamp of the last successful commit.
- `UncommittedChangeCount`: Number of new action records added locally since `LastCommittedSha`.
- `SyncStatus`: Finite state enum: `IDLE`, `DIRTY`, `SYNCING`, `EXPORTING`, `REBASING`, `MERGING`, `COMMITTED`, `ERROR`.
- `IsDirty`: Positive boolean flag (`1` = uncommitted local changes exist, `0` = synchronized with remote).

---

## 6. SQLite Engine PRAGMAs & High-Concurrency Configuration

To avoid database locking bottlenecks and ensure crash durability across concurrent multi-process writes, all connections to `root-heatmap.db`, `heatmap.db`, and local `heatmap.db` must execute the following PRAGMA sequence immediately upon opening:

```sql
-- Enforce Write-Ahead Logging for non-blocking concurrent reads and single-writer concurrency
PRAGMA journal_mode = WAL;

-- Normal synchronous mode provides durability with minimal disk flush latency
PRAGMA synchronous = NORMAL;

-- 5000ms busy timeout ensures processes wait for lock release instead of failing immediately with SQLITE_BUSY
PRAGMA busy_timeout = 5000;

-- Enforce relational integrity across tasks and logs
PRAGMA foreign_keys = ON;

-- Store temporary tables and indices in memory
PRAGMA temp_store = MEMORY;

-- Optimize query cache capacity
PRAGMA cache_size = -64000;
```

---

## 7. Go Struct Definitions & Domain Models

```go
package store

import "time"

// HeatmapAiSessionRecord maps to table HeatmapAiSession.
type HeatmapAiSessionRecord struct {
	HeatmapAiSessionId int64      `json:"heatmapAiSessionId"`
	SessionUuid        string     `json:"sessionUuid"`
	AiName             string     `json:"aiName"`
	SequencePrefix     string     `json:"sequencePrefix"`
	MachineHostname    string     `json:"machineHostname"`
	OsPlatform         string     `json:"osPlatform"`
	AgentRole          string     `json:"agentRole"`
	StartedAt          string     `json:"startedAt"`
	CompletedAt        *string    `json:"completedAt,omitempty"`
	IsActive           bool       `json:"isActive"`
	HasCompleted       bool       `json:"hasCompleted"`
	TotalTasksExecuted int        `json:"totalTasksExecuted"`
	TotalActionsLogged int        `json:"totalActionsLogged"`
	Notes              string     `json:"notes,omitempty"`
	Comments           string     `json:"comments,omitempty"`
	CreatedAt          string     `json:"createdAt"`
	UpdatedAt          string     `json:"updatedAt"`
}

// HeatmapRepoRegistryRecord maps to table HeatmapRepoRegistry.
type HeatmapRepoRegistryRecord struct {
	HeatmapRepoRegistryId int64  `json:"heatmapRepoRegistryId"`
	RepoSlug              string `json:"repoSlug"`
	RepoUrl               string `json:"repoUrl"`
	DefaultBranch         string `json:"defaultBranch"`
	FirstSeenAt           string `json:"firstSeenAt"`
	LastActiveAt          string `json:"lastActiveAt"`
	IsActive              bool   `json:"isActive"`
	TotalActionsRecorded  int    `json:"totalActionsRecorded"`
	Notes                 string `json:"notes,omitempty"`
	Comments              string `json:"comments,omitempty"`
	CreatedAt             string `json:"createdAt"`
	UpdatedAt             string `json:"updatedAt"`
}

// HeatmapCommitLogRecord maps to table HeatmapCommitLog.
type HeatmapCommitLogRecord struct {
	HeatmapCommitLogId         int64  `json:"heatmapCommitLogId"`
	TaskSlug                   string `json:"taskSlug"`
	RepoSlug                   string `json:"repoSlug"`
	CommitSha                  string `json:"commitSha"`
	CommitMessage              string `json:"commitMessage"`
	FilesCommittedCount        int    `json:"filesCommittedCount"`
	TransactionsCommittedCount int    `json:"transactionsCommittedCount"`
	CommittedAt                string `json:"committedAt"`
	PushedSuccessfully         bool   `json:"pushedSuccessfully"`
	RemoteOrigin               string `json:"remoteOrigin"`
	Notes                      string `json:"notes,omitempty"`
	Comments                   string `json:"comments,omitempty"`
	CreatedAt                  string `json:"createdAt"`
}

// HeatmapLocalTrackingRecord maps to table HeatmapLocalTracking.
type HeatmapLocalTrackingRecord struct {
	HeatmapLocalTrackingId int64   `json:"heatmapLocalTrackingId"`
	RepoSlug               string  `json:"repoSlug"`
	AiName                 string  `json:"aiName"`
	SequencePrefix         string  `json:"sequencePrefix"`
	RepoCachePath          string  `json:"repoCachePath"`
	LastCommittedSha       string  `json:"lastCommittedSha"`
	LastCommitTimestamp    string  `json:"lastCommitTimestamp"`
	UncommittedChangeCount int     `json:"uncommittedChangeCount"`
	SyncStatus             string  `json:"syncStatus"`
	IsDirty                bool    `json:"isDirty"`
	LastSyncAttemptAt      *string `json:"lastSyncAttemptAt,omitempty"`
	LastSuccessfulSyncAt   *string `json:"lastSuccessfulSyncAt,omitempty"`
	ErrorMessage           *string `json:"errorMessage,omitempty"`
	CreatedAt              string  `json:"createdAt"`
	UpdatedAt              string  `json:"updatedAt"`
}
```

---

## 8. Structured Error Codes & Exception Handling

All database operations wrap errors with `apperror.AppError` and deterministic error codes:

| Error Code | Constant | Meaning | Recovery Action |
|---|---|---|---|
| `E_HEATMAP_DB_CONNECT` | `ErrCodeHeatmapDbConnect` | Failed to open SQLite file or invalid path | Verify filesystem permissions and create parent directories |
| `E_HEATMAP_SCHEMA_MIGRATE` | `ErrCodeHeatmapSchemaMigrate` | DDL execution failed during table creation | Inspect database schema integrity or check WAL file |
| `E_HEATMAP_TRACKING_NOT_FOUND`| `ErrCodeHeatmapTrackingNotFound` | Record missing in `HeatmapLocalTracking` | Auto-seed tracking row with zero SHA and epoch timestamp |
| `E_HEATMAP_TX_EXPORT` | `ErrCodeHeatmapTxExport` | Failed to export delta to temporary database | Retry export with exponential backoff before git reset |
| `E_HEATMAP_REPLAY_MERGE` | `ErrCodeHeatmapReplayMerge` | Failed to merge backup records into rebased DB | Preserve `temp-sync-<timestamp>.db` for manual inspection |
