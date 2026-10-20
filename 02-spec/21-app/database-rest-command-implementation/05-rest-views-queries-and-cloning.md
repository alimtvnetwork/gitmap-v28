# 05-rest-views-queries-and-cloning: Views, SQL Queries, DDL, Migration & Cloning

- **Spec ID:** `02-spec/21-app/database-rest-command-implementation/05-rest-views-queries-and-cloning.md`
- **Status:** `APPROVED`
- **Version:** `1.0.0`
- **Subsystem:** Dynamic Views, Direct SQL, DDL Alterations, Split-DB Migrations & Cloning
- **Parent Plan:** `.ai-memory/plans/database-rest-command-implementation.md`
- **Dependencies:** `cli/cmddb`, `cli/store`, `cli/appfault`

---

## 1. Dynamic Views Management

Allows clients to author, list, and destroy complex multi-table SQL views on the fly.

### 1.1 `POST /api/v1/db/{alias}/views` (Create Dynamic View)
- **Request Body**:
  ```json
  {
    "ViewName": "ActiveSubtasksWithParent",
    "Query": "SELECT s.SubtaskId, s.Title, s.Status, t.TaskName, t.Owner FROM Subtask s INNER JOIN Task t ON s.TaskId = t.TaskId WHERE s.Status != 'DONE'"
  }
  ```
- **Validation**:
  1. `ViewName` must be a valid SQLite identifier (`^[a-zA-Z0-9_]+$`).
  2. `Query` must be a read-only `SELECT` statement.
- **Execution**:
  1. Executes `CREATE VIEW IF NOT EXISTS <ViewName> AS <Query>`.
  2. Records view definition in Split-DB cache.
- **Response Status**: `201 Created`

### 1.2 `GET /api/v1/db/{alias}/views` (List Dynamic Views)
- **Response**: List of all views in the database with their underlying SQL definitions.

### 1.3 `DELETE /api/v1/db/{alias}/views/{name}` (Drop View)
- **Execution**: Executes `DROP VIEW IF EXISTS <name>`.
- **Response Status**: `200 OK`

---

## 2. Direct SQL Query Execution Endpoints

### 2.1 `GET /api/v1/db/{alias}/query?sql={sqlQuery}`
Convenience endpoint for read-only queries passed as a URL query parameter.
- **Query Parameters**:
  - `sql` (required, URL-encoded string): Read-only SQL query (e.g. `SELECT count(*), status FROM Task GROUP BY status`).
- **Safety Restriction**:
  Only statements beginning with `SELECT`, `PRAGMA`, or `EXPLAIN` are permitted. Mutation statements return `403 Forbidden`.

### 2.2 `POST /api/v1/db/{alias}/query`
General SQL execution endpoint supporting parameterized queries, joins, and DML statements.
- **Request Body**:
  ```json
  {
    "Sql": "UPDATE Task SET Status = ? WHERE TaskId = ?",
    "Params": ["IN_PROGRESS", "task-99"]
  }
  ```
- **Response**:
  For `SELECT` queries, returns row objects array. For DML (`UPDATE`, `INSERT`, `DELETE`), returns execution metadata `{ RowsAffected, LastInsertId, ExecutionDurationMs }`.

---

## 3. Table Alteration & DDL Endpoints

### 3.1 `POST /api/v1/db/{alias}/tables/create` (Create New Table)
- **Request Body**:
  ```json
  {
    "TableName": "AuditLog",
    "Columns": [
      { "Name": "AuditId", "Type": "INTEGER", "IsPrimaryKey": true, "IsAutoIncrement": true },
      { "Name": "Action", "Type": "TEXT", "IsNullable": false },
      { "Name": "Payload", "Type": "TEXT", "IsNullable": true },
      { "Name": "CreatedAt", "Type": "TEXT", "DefaultValue": "CURRENT_TIMESTAMP" }
    ]
  }
  ```
- **Response Status**: `201 Created`

### 3.2 `POST /api/v1/db/{alias}/alter` (Alter Table Schema)
Supports adding columns, renaming columns, and renaming tables:
- **Request Body**:
  ```json
  {
    "Action": "ADD_COLUMN",
    "TableName": "AuditLog",
    "Column": {
      "Name": "ActorId",
      "Type": "TEXT",
      "IsNullable": true
    }
  }
  ```
- **Supported Actions**:
  - `ADD_COLUMN`: Executes `ALTER TABLE <table> ADD COLUMN <name> <type>`.
  - `RENAME_COLUMN`: Executes `ALTER TABLE <table> RENAME COLUMN <oldName> TO <newName>`.
  - `RENAME_TABLE`: Executes `ALTER TABLE <oldName> RENAME TO <newName>`.

---

## 4. Database Migration Engine & Split-DB Caching

Provides an automated migration runner that stores migration records in a dedicated caching Split-DB:

### `POST /api/v1/db/{alias}/migrate`
- **Request Body**:
  ```json
  {
    "MigrationName": "002_add_index_subtask_status",
    "UpSql": "CREATE INDEX IF NOT EXISTS idx_subtask_status ON Subtask(Status);",
    "DownSql": "DROP INDEX IF EXISTS idx_subtask_status;"
  }
  ```
- **Execution Flow**:
  1. Opens Split-DB cache database (`.gitmap/repodb/dbrest_cache.db`).
  2. Ensures `MigrationLog` table exists.
  3. Checks if `MigrationName` was already applied for this alias.
  4. Executes `UpSql` in transaction and records result in `MigrationLog`.
- **Response Status**: `200 OK`

---

## 5. Database Cloning Endpoint

### `POST /api/v1/db/{alias}/clone`
Creates an exact physical copy of an existing mounted SQLite database, registers it under a new alias, and mounts it into the running REST server dynamically:
- **Request Body**:
  ```json
  {
    "DestinationPath": "./data/clones/tasks_backup.db",
    "NewAlias": "tasks-backup",
    "AutoMount": true
  }
  ```
- **Execution**:
  1. Resolves source database path.
  2. Ensures parent directory exists.
  3. Copies file or uses SQLite `VACUUM INTO ?`.
  4. Registers `NewAlias` in `DatabaseAlias` table.
  5. Mounts alias into active `MultiDbManager`.

---

## 6. Concurrency & Deadlock Prevention Invariants

1. **Single Connection per DB**: Every mounted database pool uses `conn.SetMaxOpenConns(1)`.
2. **WAL Journaling**: Each connection sets `PRAGMA journal_mode=WAL; PRAGMA busy_timeout=5000;`.
3. **Full Row Ingestion**: In all handler functions, rows returned by `conn.Query` are scanned into an in-memory slice and `rows.Close()` is called immediately before returning.
4. **Mutex Guarding**: The `MultiDbManager` maintains a `sync.RWMutex` to guard dynamic registration of new aliases.
