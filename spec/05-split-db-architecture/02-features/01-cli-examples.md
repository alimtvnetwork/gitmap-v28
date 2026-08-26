# Split DB Architecture: CLI Examples

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**Status:** Active  
**Parent:** [00-overview.md](../00-overview.md)

---

## Overview

Concrete database structure examples for each CLI project using the Split DB pattern. All field names use **PascalCase** (no underscores).

---

## 1. AI Bridge CLI (Primary Example)

### Database Structure

```
data/
├── aibridge.db                                    # ROOT DB
│
└── myproject/                                     # APP FOLDER
    │
    ├── search.db                                  # Search metadata DB
    │
    ├── rag/
    │   ├── cache/
    │   │   └── search/
    │   │       ├── 001-golang-patterns-abc123.db  # Cached search results
    │   │       └── 002-react-hooks-def456.db
    │   │
    │   └── documents/
    │       ├── 001-readme-md.db                   # RAG chunks for readme.md
    │       └── 002-main-go.db                     # RAG chunks for main.go
    │
    ├── ai/
    │   └── chat/
    │       ├── 001-chat-xyz789.db                 # Chat session + messages
    │       └── 002-chat-abc123.db
    │
    └── settings/
        └── config.db                              # App-level settings override
```

### Root DB Schema (`aibridge.db`)

```sql
-- Settings table (seeded from config.seed.json)
CREATE TABLE Setting (
    SettingId INTEGER PRIMARY KEY AUTOINCREMENT,
    Key TEXT UNIQUE NOT NULL,
    Value TEXT NOT NULL,
    ValueType TEXT DEFAULT 'string',
    Source TEXT DEFAULT 'user',
    Description TEXT,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Default settings from seed
INSERT INTO Setting (Key, Value, ValueType, Source, Description) VALUES
('Search.Cache.TtlDays', '5', 'int', 'seed', 'Search cache TTL in days'),
('Search.Cache.MaxEntries', '1000', 'int', 'seed', 'Max cached entries per app'),
('Rag.ChunkSize', '512', 'int', 'seed', 'Default chunk size in tokens'),
('Rag.ChunkOverlap', '50', 'int', 'seed', 'Chunk overlap in tokens');

-- Applications registry
CREATE TABLE Application (
    ApplicationId INTEGER PRIMARY KEY AUTOINCREMENT,
    AppName TEXT UNIQUE NOT NULL,
    DisplayName TEXT NOT NULL,
    Description TEXT,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    LastAccessed DATETIME,
    Status TEXT DEFAULT 'active'
);

-- Sequence counters
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE Counter (
    CounterId INTEGER PRIMARY KEY AUTOINCREMENT,
    ApplicationId INTEGER NOT NULL,
    Category TEXT NOT NULL,
    SubCategory TEXT NOT NULL,
    CurrentCount INTEGER DEFAULT 0,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ApplicationId) REFERENCES Application(ApplicationId),
    UNIQUE(ApplicationId, Category, SubCategory)
);

-- Database registry
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE DbRegistry (
    DbRegistryId INTEGER PRIMARY KEY AUTOINCREMENT,
    ApplicationId INTEGER NOT NULL,
    Category TEXT NOT NULL,
    SubCategory TEXT NOT NULL,
    EntityId TEXT NOT NULL,
    SequenceNum INTEGER NOT NULL,
    Path TEXT NOT NULL,
    DisplayName TEXT,
    SizeBytes INTEGER DEFAULT 0,
    RecordCount INTEGER DEFAULT 0,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    LastAccessed DATETIME,
    ExpiresAt DATETIME,
    Status TEXT DEFAULT 'active',
    FOREIGN KEY (ApplicationId) REFERENCES Application(ApplicationId)
);
```

### Chat Session DB Schema (`data/{app}/ai/chat/001-{id}.db`)

```sql
-- Session metadata
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE SessionMeta (
    SessionMetaId INTEGER PRIMARY KEY AUTOINCREMENT,
    SessionId TEXT UNIQUE NOT NULL,
    Title TEXT,
    ModelCategory TEXT NOT NULL,
    ModelUsed TEXT,
    BackendUsed TEXT,
    RagEnabled BOOLEAN DEFAULT FALSE,
    RagSources TEXT,
    MessageCount INTEGER DEFAULT 0,
    TotalTokens INTEGER DEFAULT 0,
    TotalToolCalls INTEGER DEFAULT 0,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    LastMessageAt DATETIME,
    Status TEXT DEFAULT 'active'
);

-- Messages (conversation history)
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE Message (
    MessageId INTEGER PRIMARY KEY AUTOINCREMENT,
    SequenceNum INTEGER NOT NULL,
    Role TEXT NOT NULL,
    Content TEXT NOT NULL,
    Tokens INTEGER,
    Model TEXT,
    RagContext TEXT,
    RagChunkIds TEXT,
    Metadata TEXT,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    CompletedAt DATETIME
);

-- Tool calls
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE ToolCalls (
    ToolCallsId INTEGER PRIMARY KEY AUTOINCREMENT,
    MessageId INTEGER NOT NULL,
    ToolName TEXT NOT NULL,
    Arguments TEXT,
    Result TEXT,
    ResultType TEXT,
    StartedAt DATETIME,
    CompletedAt DATETIME,
    DurationMs INTEGER,
    Status TEXT NOT NULL DEFAULT 'pending',
    ErrorMessage TEXT,
    FOREIGN KEY (MessageId) REFERENCES Message(MessageId)
);
```

### API to DB Mapping

| Endpoint | Action | Target DB |
|----------|--------|-----------|
| `POST /api/v1/chat/sessions` | Create session | Creates `data/{app}/ai/chat/{seq}-{id}.db` |
| `GET /api/v1/chat/sessions?appName=X` | List sessions | Reads from `aibridge.db` → DbRegistry |
| `GET /api/v1/chat/sessions/:id` | Get session | Reads `data/{app}/ai/chat/{seq}-{id}.db` → SessionMeta |
| `POST /api/v1/chat/sessions/:id/messages` | Send message | Writes to `data/{app}/ai/chat/{seq}-{id}.db` → Messages + ToolCalls |
| `GET /api/v1/chat/sessions/:id/messages` | Get messages | Reads `data/{app}/ai/chat/{seq}-{id}.db` → Messages |

---

## 2. GSearch CLI

### Database Structure

```
data/
├── gsearch.db                                     # ROOT DB
│
└── searches/                                      # Search data folder
    │
    ├── search.db                                  # Search history/metadata
    │
    └── cache/
        ├── 001-ai-tools-abc123.db                 # Cached: "AI tools 2026"
        ├── 002-golang-patterns-def456.db          # Cached: "golang patterns"
        └── 003-react-hooks-ghi789.db              # Cached: "react hooks"
```

### Root DB Schema (`gsearch.db`)

```sql
-- Settings
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE Setting (
    SettingId INTEGER PRIMARY KEY AUTOINCREMENT,
    Key TEXT UNIQUE NOT NULL,
    Value TEXT NOT NULL,
    ValueType TEXT DEFAULT 'string',
    Source TEXT DEFAULT 'seed',
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO Setting (Key, Value, ValueType, Source) VALUES
('Cache.TtlDays', '5', 'int', 'seed'),
('Cache.MaxEntries', '500', 'int', 'seed'),
('Search.DefaultEngine', 'google', 'string', 'seed'),
('Search.MaxResults', '10', 'int', 'seed');

-- Counters
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE Counter (
    CounterId INTEGER PRIMARY KEY AUTOINCREMENT,
    Category TEXT NOT NULL,
    CurrentCount INTEGER DEFAULT 0,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Database registry
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE DbRegistry (
    DbRegistryId INTEGER PRIMARY KEY AUTOINCREMENT,
    Category TEXT NOT NULL,
    EntityId TEXT NOT NULL,
    SequenceNum INTEGER NOT NULL,
    Path TEXT NOT NULL,
    SizeBytes INTEGER DEFAULT 0,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    ExpiresAt DATETIME,
    Status TEXT DEFAULT 'active'
);
```

### Search History DB (`searches/search.db`)

```sql
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE SearchLog (
    SearchLogId INTEGER PRIMARY KEY AUTOINCREMENT,
    Query TEXT NOT NULL,
    QueryHash TEXT NOT NULL,
    SearchType TEXT NOT NULL,
    Engine TEXT NOT NULL,
    SearchedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    DurationMs INTEGER,
    ResultCount INTEGER DEFAULT 0,
    CacheHit BOOLEAN DEFAULT FALSE,
    CacheDbPath TEXT,
    CacheExpiresAt DATETIME,
    Status TEXT DEFAULT 'completed'
);

CREATE INDEX IdxSearchLogHash ON SearchLog(QueryHash);
CREATE INDEX IdxSearchLogTime ON SearchLog(SearchedAt DESC);
```

### Cache DB Schema (`searches/cache/001-{slug}.db`)

```sql
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE CacheMeta (
    CacheMetaId INTEGER PRIMARY KEY AUTOINCREMENT,
    Query TEXT NOT NULL,
    QueryHash TEXT NOT NULL,
    SearchType TEXT NOT NULL,
    Engine TEXT NOT NULL,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    ExpiresAt DATETIME NOT NULL,
    LastAccessed DATETIME,
    AccessCount INTEGER DEFAULT 1,
    TotalResults INTEGER DEFAULT 0,
    Status TEXT DEFAULT 'active'
);

-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE Result (
    ResultId INTEGER PRIMARY KEY AUTOINCREMENT,
    Rank INTEGER NOT NULL,
    Title TEXT NOT NULL,
    Url TEXT NOT NULL,
    Snippet TEXT,
    Source TEXT,
    Score REAL,
    Metadata TEXT,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

---

## 3. BRun CLI

### Database Structure

```
data/
├── brun.db                                        # ROOT DB
│
└── runs/
    ├── 001-backend-build-abc.db                   # Build run session
    ├── 002-frontend-build-def.db
    └── 003-full-stack-ghi.db
```

### Root DB Schema (`brun.db`)

```sql
-- Settings
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE Setting (
    SettingId INTEGER PRIMARY KEY AUTOINCREMENT,
    Key TEXT UNIQUE NOT NULL,
    Value TEXT NOT NULL,
    ValueType TEXT DEFAULT 'string',
    Source TEXT DEFAULT 'seed',
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO Setting (Key, Value, ValueType, Source) VALUES
('Runs.KeepCount', '100', 'int', 'seed'),
('Runs.VacuumInterval', '24h', 'string', 'seed');

-- Profiles
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE Profile (
    ProfileId INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT UNIQUE NOT NULL,
    Runtime TEXT NOT NULL,
    Command TEXT,
    WorkDir TEXT,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Counters
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE Counter (
    CounterId INTEGER PRIMARY KEY AUTOINCREMENT,
    Category TEXT NOT NULL,
    CurrentCount INTEGER DEFAULT 0,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Database registry
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE DbRegistry (
    DbRegistryId INTEGER PRIMARY KEY AUTOINCREMENT,
    Category TEXT NOT NULL,
    EntityId TEXT NOT NULL,
    SequenceNum INTEGER NOT NULL,
    Path TEXT NOT NULL,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    Status TEXT DEFAULT 'active'
);
```

### Run Session DB (`runs/001-{id}.db`)

```sql
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE BuildRun (
    BuildRunId INTEGER PRIMARY KEY AUTOINCREMENT,
    RunId TEXT UNIQUE NOT NULL,
    ProfileName TEXT,
    Runtime TEXT NOT NULL,
    Command TEXT,
    WorkDir TEXT,
    ExitCode INTEGER NOT NULL DEFAULT 0,
    Success BOOLEAN NOT NULL DEFAULT FALSE,
    Stdout TEXT,
    Stderr TEXT,
    StartTime DATETIME NOT NULL,
    EndTime DATETIME NOT NULL,
    DurationMs INTEGER NOT NULL,
    Port INTEGER DEFAULT 0,
    LogPath TEXT,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE BuildErrors (
    BuildErrorsId INTEGER PRIMARY KEY AUTOINCREMENT,
    BuildRunId INTEGER NOT NULL,
    File TEXT,
    Line INTEGER DEFAULT 0,
    Column INTEGER DEFAULT 0,
    Message TEXT NOT NULL,
    Severity TEXT NOT NULL,
    Code TEXT,
    StackTrace TEXT,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (BuildRunId) REFERENCES BuildRun(BuildRunId)
);

-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE AssetOperations (
    AssetOperationsId INTEGER PRIMARY KEY AUTOINCREMENT,
    BuildRunId INTEGER NOT NULL,
    Source TEXT NOT NULL,
    Destination TEXT NOT NULL,
    Mode TEXT NOT NULL,
    FilesCopied INTEGER NOT NULL DEFAULT 0,
    FilesSkipped INTEGER NOT NULL DEFAULT 0,
    BytesCopied INTEGER NOT NULL DEFAULT 0,
    DurationMs INTEGER NOT NULL DEFAULT 0,
    Success BOOLEAN NOT NULL DEFAULT FALSE,
    ErrorMsg TEXT,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (BuildRunId) REFERENCES BuildRun(BuildRunId)
);
```

---

## 4. Nexus Flow CLI

### Database Structure

```
data/
├── nexusflow.db                                   # ROOT DB
│
└── workflows/
    │
    ├── pipeline-001/
    │   ├── meta.db                                # Pipeline definition
    │   │
    │   ├── executions/
    │   │   ├── 001-exec-abc.db                    # Execution session
    │   │   └── 002-exec-def.db
    │   │
    │   └── checkpoints/
    │       ├── 001-checkpoint-abc.db              # RES checkpoint
    │       └── 002-checkpoint-def.db
    │
    └── pipeline-002/
        └── ...
```

### Root DB Schema (`nexusflow.db`)

```sql
-- Settings
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE Setting (
    SettingId INTEGER PRIMARY KEY AUTOINCREMENT,
    Key TEXT UNIQUE NOT NULL,
    Value TEXT NOT NULL,
    ValueType TEXT DEFAULT 'string',
    Source TEXT DEFAULT 'seed',
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO Setting (Key, Value, ValueType, Source) VALUES
('Execution.Timeout', '30m', 'string', 'seed'),
('Execution.MaxConcurrent', '5', 'int', 'seed'),
('Checkpoint.Enabled', 'true', 'bool', 'seed');

-- Pipelines registry
CREATE TABLE Pipeline (
    PipelineId INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT NOT NULL,
    Description TEXT,
    Version TEXT DEFAULT '1.0.0',
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    Status TEXT DEFAULT 'active'
);

-- Counters
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE Counter (
    CounterId INTEGER PRIMARY KEY AUTOINCREMENT,
    PipelineId INTEGER,
    Category TEXT NOT NULL,
    CurrentCount INTEGER DEFAULT 0,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Database registry
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE DbRegistry (
    DbRegistryId INTEGER PRIMARY KEY AUTOINCREMENT,
    PipelineId INTEGER NOT NULL,
    Category TEXT NOT NULL,
    EntityId TEXT NOT NULL,
    SequenceNum INTEGER NOT NULL,
    Path TEXT NOT NULL,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    Status TEXT DEFAULT 'active',
    FOREIGN KEY (PipelineId) REFERENCES Pipeline(PipelineId)
);
```

### Execution Session DB (`workflows/pipeline-001/executions/001-{id}.db`)

```sql
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE ExecutionMeta (
    ExecutionMetaId INTEGER PRIMARY KEY AUTOINCREMENT,
    ExecutionId TEXT UNIQUE NOT NULL,
    PipelineId INTEGER NOT NULL,
    Input TEXT,
    Output TEXT,
    StartedAt DATETIME,
    CompletedAt DATETIME,
    DurationMs INTEGER,
    Status TEXT NOT NULL DEFAULT 'pending',
    ErrorMessage TEXT
);

-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE BlockExecutions (
    BlockExecutionsId INTEGER PRIMARY KEY AUTOINCREMENT,
    ExecutionId TEXT NOT NULL,
    BlockId TEXT NOT NULL,
    BlockType TEXT NOT NULL,
    Input TEXT,
    Output TEXT,
    StartedAt DATETIME,
    CompletedAt DATETIME,
    DurationMs INTEGER,
    Status TEXT NOT NULL DEFAULT 'pending',
    ErrorMessage TEXT,
    RetryCount INTEGER DEFAULT 0
);

-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE BlockLogs (
    BlockLogsId INTEGER PRIMARY KEY AUTOINCREMENT,
    BlockExecutionsId INTEGER NOT NULL,
    Level TEXT NOT NULL,
    Message TEXT NOT NULL,
    Metadata TEXT,
    Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (BlockExecutionsId) REFERENCES BlockExecutions(BlockExecutionsId)
);
```

---

## JSON Transport Examples

### AI Bridge: Create Session Request

```json
{
  "AppName": "myproject",
  "Title": "Code Review Session",
  "ModelCategory": "coding",
  "RagEnabled": true,
  "RagSources": ["001-readme-md", "002-main-go"]
}
```

### AI Bridge: Send Message Request

```json
{
  "Role": "user",
  "Content": "Explain the main function in main.go",
  "Metadata": {
    "ClientVersion": "1.0.0",
    "Timestamp": "2026-02-02T10:30:00Z"
  }
}
```

### AI Bridge: Message Response

```json
{
  "Id": "msg-xyz789",
  "SessionId": "001-chat-abc123",
  "SequenceNum": 3,
  "Role": "assistant",
  "Content": "The main function initializes the application...",
  "Tokens": 245,
  "Model": "codellama:13b",
  "RagContext": [
    {
      "ChunkId": "chunk-001",
      "Content": "func main() { ... }",
      "Similarity": 0.92
    }
  ],
  "CreatedAt": "2026-02-02T10:30:05Z",
  "CompletedAt": "2026-02-02T10:30:12Z"
}
```

---

## 5. Reset API Tables (All CLIs)

All CLI root databases include the following table for 2-step reset confirmation:

### ResetRequests Table Schema (PascalCase)

```sql
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE ResetRequests (
    ResetRequestsId INTEGER PRIMARY KEY AUTOINCREMENT,                           -- rst_{uuid}
    Scope TEXT NOT NULL,                           -- "all", "app", "cache", etc.
    AppName TEXT,                                  -- For app-scoped resets
    RequestedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    ExpiresAt DATETIME NOT NULL,                   -- RequestedAt + 5 minutes
    AffectedItems TEXT,                            -- JSON: preview of deletion
    ConfirmedAt DATETIME,
    CancelledAt DATETIME,
    CompletedAt DATETIME,
    Status TEXT DEFAULT 'pending',                 -- pending, confirmed, expired, cancelled, completed
    DeletedCount INTEGER,
    FreedBytes INTEGER,
    ErrorMessage TEXT
);

CREATE INDEX IdxResetStatus ON ResetRequests(Status, ExpiresAt);
CREATE INDEX IdxResetApp ON ResetRequests(AppName);
```

### Reset Scopes by CLI

| CLI | Available Scopes |
|-----|------------------|
| AI Bridge | `all`, `app`, `chat`, `rag`, `search`, `seo` |
| GSearch | `all`, `cache`, `history` |
| BRun | `all`, `runs`, `profile:{name}` |
| Nexus Flow | `all`, `executions`, `pipeline:{id}`, `checkpoints` |

### API Endpoints (All CLIs)

```
POST /api/v1/reset/request
  Request:  { "Scope": "<scope>", "AppName": "<app>" }
  Response: { "ResetId": "rst_abc123", "ExpiresAt": "...", "AffectedItems": {...} }

POST /api/v1/reset/confirm
  Request:  { "ResetId": "rst_abc123" }
  Response: { "Status": "completed", "DeletedDatabases": 87, "FreedBytes": 524288000 }

POST /api/v1/reset/cancel
  Request:  { "ResetId": "rst_abc123" }
  Response: { "Status": "cancelled" }
```

### Example: AI Bridge Reset Flow

```json
// Step 1: Request
POST /api/v1/reset/request
{
  "Scope": "app",
  "AppName": "myproject"
}

// Response
{
  "ResetId": "rst_abc123def456",
  "Scope": "app",
  "AppName": "myproject",
  "ExpiresAt": "2026-02-02T10:35:00Z",
  "AffectedItems": {
    "ChatSessions": 12,
    "RagDocuments": 45,
    "SearchCaches": 30,
    "SeoJobs": 5,
    "TotalDatabases": 92,
    "TotalBytes": 524288000
  },
  "Message": "Review affected items and confirm within 5 minutes"
}

// Step 2: Confirm (within 5 minutes)
POST /api/v1/reset/confirm
{
  "ResetId": "rst_abc123def456"
}

// Response
{
  "Status": "completed",
  "Scope": "app",
  "AppName": "myproject",
  "DeletedDatabases": 92,
  "FreedBytes": 524288000,
  "Duration": "2.3s"
}
```

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Split DB Overview | `../00-overview.md` |
| Reset API Standard | `./02-reset-api-standard.md` |
| Database Flow Diagrams | `./03-database-flow-diagrams.md` |
| AI Bridge DB Architecture | `../22-ai-bridge-cli/01-backend/12-database-architecture.md` |
| GSearch DB Schema | `../20-gsearch-cli/01-backend/03-database-schema.md` |
| BRun Data Models | `../21-brun-cli/01-backend/10-data-models.md` |
| Naming Conventions | `../02-coding-guidelines/01-cross-language/07-database-naming.md` |

---

*Concrete examples for every CLI.*
