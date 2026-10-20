# 02-storage-and-split-db-cache: Split-DB Architecture & Storage Tier

- **Spec ID:** `02-spec/21-app/database-rest-command-implementation/02-storage-and-split-db-cache.md`
- **Status:** `APPROVED`
- **Version:** `1.0.0`
- **Subsystem:** Split-DB Migration Cache, Invariant Enforcement, Storage Tier
- **Parent Plan:** `.ai-memory/plans/database-rest-command-implementation.md`
- **Dependencies:** `cli/cmddb`, `cli/store`, `cli/appfault`

---

## 1. Split-DB Migration & Caching Storage

### 1.1 Architecture Rationale
Target databases mounted by `db-rest` may belong to other microservices, external tools, or production assets. Directly injecting metadata or migration tables into target databases risks schema pollution, namespace collisions, or foreign constraint breaks.

GitMap maintains a dedicated Split-DB cache at `.gitmap/repodb/dbrest_cache.db`. All non-destructive auxiliary state—such as migration audit history, DDL alterations, user-defined query views, and database clone lineage—is safely persisted inside this dedicated Split-DB.

### 1.2 Split-DB Schema Definitions
All tables adhere strictly to PascalCase table names, camelCase column names, and positive boolean principles:

```sql
-- Migration execution audit log
CREATE TABLE IF NOT EXISTS MigrationLog (
    MigrationId INTEGER PRIMARY KEY AUTOINCREMENT,
    DatabaseAlias TEXT NOT NULL,
    DatabasePath TEXT NOT NULL,
    VersionTag TEXT NOT NULL,
    Description TEXT,
    SqlStatements TEXT NOT NULL,
    ExecutionDurationMs INTEGER NOT NULL DEFAULT 0,
    IsSuccess INTEGER NOT NULL DEFAULT 1,
    ErrorMessage TEXT,
    AppliedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    AppliedBy TEXT DEFAULT 'rest-api'
);
CREATE INDEX IF NOT EXISTS idx_migration_alias ON MigrationLog (DatabaseAlias);

-- DDL alteration history log
CREATE TABLE IF NOT EXISTS AlterationHistory (
    AlterationId INTEGER PRIMARY KEY AUTOINCREMENT,
    DatabaseAlias TEXT NOT NULL,
    TableName TEXT NOT NULL,
    ActionType TEXT NOT NULL, -- 'ADD_COLUMN', 'ALTER_COLUMN', 'CREATE_TABLE', 'CREATE_INDEX'
    SqlExecuted TEXT NOT NULL,
    ExecutionDurationMs INTEGER NOT NULL DEFAULT 0,
    IsSuccess INTEGER NOT NULL DEFAULT 1,
    ErrorMessage TEXT,
    ExecutedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_alteration_alias_table ON AlterationHistory (DatabaseAlias, TableName);

-- Dynamic saved views created via REST API
CREATE TABLE IF NOT EXISTS ViewDefinition (
    ViewId INTEGER PRIMARY KEY AUTOINCREMENT,
    DatabaseAlias TEXT NOT NULL,
    ViewName TEXT NOT NULL,
    SourceQuery TEXT NOT NULL,
    Description TEXT,
    IsActive INTEGER NOT NULL DEFAULT 1,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(DatabaseAlias, ViewName)
);

-- Registry of on-the-fly database clones
CREATE TABLE IF NOT EXISTS ClonedDatabase (
    CloneId INTEGER PRIMARY KEY AUTOINCREMENT,
    SourceDatabasePath TEXT NOT NULL,
    TargetDatabasePath TEXT NOT NULL,
    TargetAlias TEXT NOT NULL UNIQUE,
    FileSizeBytes INTEGER NOT NULL DEFAULT 0,
    IsActive INTEGER NOT NULL DEFAULT 1,
    ClonedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

---

## 2. Concurrency Governance & SQLite Safety Invariants

### 2.1 Single-Writer Connection Pool Invariant
To prevent `database is locked` panics in Go's `database/sql` driver:
- Every opened SQLite database connection pool sets `conn.SetMaxOpenConns(1)`.
- Concurrent HTTP read/write requests share this serialized pool or utilize WAL mode for non-blocking concurrent reads while writes are serialized.
- Mutating operations (POST, PUT, DELETE, ALTER, MIGRATE, CLONE) acquire a per-database `sync.Mutex` lock in memory before initiating database transactions.

### 2.2 In-Memory Result Buffering Invariant
Nested queries over an open SQLite handle cause immediate lock contention:
- Handlers MUST scan all `*sql.Rows` records into allocated Go slices and immediately close the row iterator (`rows.Close()`) before executing any subsequent queries or formatting responses.
- Handlers MUST NEVER leave an active `*sql.Rows` open while invoking helper functions that open secondary statements.

### 2.3 Mandatory SQLite PRAGMAs
Every database connection opened by `db-rest` executes the standard GitMap SQLite PRAGMA initialization sequence:
```sql
PRAGMA busy_timeout = 5000;
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA foreign_keys = ON;
```

---

## 3. Universal Response Envelope & Error Governance

### 3.1 Universal Response Envelope Structure
All REST responses conform strictly to the standardized GitMap PascalCase envelope:

```json
{
  "Status": "success",
  "Attributes": {
    "Database": "appdb",
    "Table": "Users",
    "RecordCount": 2,
    "Page": 1,
    "PageSize": 50,
    "HasMore": false
  },
  "Results": [
    {
      "UserId": 1,
      "Username": "alice",
      "Email": "alice@example.com",
      "IsActive": true
    }
  ]
}
```

### 3.2 Standard Error Response Structure
When an operation fails, the server responds with a structured error envelope:

```json
{
  "Status": "error",
  "Attributes": {
    "ErrorCode": "RESOURCE_NOT_FOUND",
    "Database": "appdb",
    "ErrorMessage": "table 'InvalidTable' does not exist in database 'appdb'",
    "Timestamp": "2026-10-20T02:00:00Z"
  },
  "Results": []
}
```

### 3.3 Standard Error Codes
| HTTP Status | ErrorCode | Description |
|---|---|---|
| `400 Bad Request` | `VALIDATION_FAILED` | Malformed parameters, invalid SQL syntax, or missing required fields. |
| `404 Not Found` | `DATABASE_NOT_FOUND` | Specified database alias or path is not mounted. |
| `404 Not Found` | `TABLE_NOT_FOUND` | Specified table does not exist in the target database. |
| `404 Not Found` | `RECORD_NOT_FOUND` | Record with specified primary key not found. |
| `409 Conflict` | `ALIAS_CONFLICT` | Attempt to mount or register an alias that already exists. |
| `409 Conflict` | `DATABASE_LOCKED` | SQLite transaction busy timeout exceeded. |
| `500 Internal Error` | `EXECUTION_FAILED` | Unhandled SQL execution or internal filesystem error. |

---

## 4. CLI Companion Commands Architecture

### 4.1 CLI Record Deletion (`db-record-del`)
- **Command Syntax:** `gitmap db-record-del <db> -table <table> -id <id>`
- Resolves `<db>` via `cmddb.ResolveDbPath`.
- Inspects `<table>` schema to determine primary key column name (defaults to `id` or first PK column in `PRAGMA table_info`).
- Executes `DELETE FROM <table> WHERE <pkCol> = ?` using single-writer connection.
- Returns output formatted as JSON Universal Envelope or human-readable summary.

### 4.2 CLI Database Removal (`db-rm`)
- **Command Syntax:** `gitmap db-rm <alias> [--purge-file]`
- Deregisters `<alias>` from `DatabaseAlias` in `gitmap.db`.
- If `--purge-file` flag is provided, additionally removes the physical SQLite file from disk with confirmation guard.
- Cleans up associated Split-DB cache records in `MigrationLog`, `AlterationHistory`, and `ViewDefinition`.

---

## 5. Verification Gates & Acceptance Criteria

```yaml
verificationGates:
  isMultiDbMountingSupported: true
  isPortProbingAndFallbackValidated: true
  isSplitDbMigrationCacheSeparated: true
  isSingleWriterSafetyEnforced: true
  isUniversalResponseEnvelopeCompliant: true
  isZeroExternalWebFrameworkDependencies: true
  isPureGoSqliteDriverEnforced: true
```

- [x] **G1 (Port Probing):** Comma-separated candidate ports probe sequentially and bind to first available port.
- [x] **G2 (Multi-DB Alias Resolution):** Single and multiple aliases/paths are resolved, mounted, and listed in `/api/v1/databases`.
- [x] **G3 (Universal Response Envelope):** All endpoints return PascalCase `{ Status, Attributes, Results }`.
- [x] **G4 (Single-Writer Concurrency):** No `database is locked` panics under concurrent read/write test load.
- [x] **G5 (Split-DB Cache Isolation):** Migrations, views, and alterations record in dedicated cache DB without altering user schema.
- [x] **G6 (CLI Parity):** Companion commands `gitmap db-record-del` and `gitmap db-rm` operate properly from terminal.
