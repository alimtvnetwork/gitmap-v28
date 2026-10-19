# Database REST CLI Commands and REST API Specification

**Specification Document:** `02-spec/21-app/database-rest-command-implementation/02-cli-and-rest-api-spec.md`  
**Parent Specification Suite:** `02-spec/21-app/database-rest-command-implementation/readme.md`  
**Architecture Specification:** `02-spec/21-app/database-rest-command-implementation/01-architecture-spec.md`  
**Master Plan:** `.ai-memory/plans/database-rest-command-implementation.md`  
**Status:** APPROVED / IN_IMPLEMENTATION  
**Version:** 1.0.0  

---

## 1. Executive Summary & Purpose

This specification defines the exhaustive technical contract for the GitMap SQLite Database REST (`db-rest`) CLI command suite and its underlying high-performance HTTP REST API engine. 

The system enables developers to expose any local SQLite database—or multiple SQLite databases simultaneously—over a standards-compliant RESTful interface with on-the-fly schema introspection, dynamic CRUD generation, column search across tables, view authoring with joins, table alteration/DDL execution, schema migration with Split-DB caching, database cloning, and raw SQL execution.

### Architectural Invariants
1. **Universal Response Envelope**: Every HTTP response MUST return JSON with root keys `{ Status, Attributes, Results }` using PascalCase keys throughout.
2. **Standard Library HTTP Router**: Built on Go standard library `net/http` router with zero third-party web framework dependencies.
3. **CGO-Free Pure Go SQLite**: Built on `modernc.org/sqlite` with `conn.SetMaxOpenConns(1)` and WAL journal mode (`_busy_timeout=5000&_journal_mode=WAL`) to guarantee single-writer safety and eliminate database locks.
4. **Buffer-Before-Next Pattern**: All query results are fully scanned into in-memory data structures before releasing row iterators or executing subsequent queries.
5. **Multi-DB Registry Routing**: Uniformly mounts system aliases (`master`, `pipeline`, `tasks`, `heatmap`), custom registered aliases, and direct SQLite file paths under discrete namespace slugs (`/api/v1/db/{alias}/...`).

---

## 2. CLI Command Syntax & Parsing Contracts

The CLI layer is implemented in package `cli/cmddb` and registered in `cli/cmd`. The commands provide terminal access to start REST servers, bind aliases, delete individual records, and remove databases/aliases.

### 2.1 Command Grammar Summary

| Command Form | Syntax Example | Core Behavior |
|---|---|---|
| **Alias Binding** | `gitmap db-rest "db path" as alias` | Binds an arbitrary SQLite file path to a persistent alias in `DatabaseAlias` table. |
| **Single-DB Server** | `gitmap db-rest aliasName $port,$port2` | Starts REST service mounting `aliasName`, probing ports in order until one binds. |
| **Multi-DB Server (Aliases)** | `gitmap db-rest alias1,alias2 $port,$port2` | Mounts multiple database aliases simultaneously under a unified catalog. |
| **Multi-DB Server (Mixed)** | `gitmap db-rest alias1,dbPath2 $port,$port2` | Mounts registered aliases alongside direct filesystem paths with auto-generated aliases. |
| **Record Deletion** | `gitmap db-record-del <db> <table> <id>` (alias: `dbrd`) | Deletes a record from a specified table by its primary key. |
| **Alias / DB Removal** | `gitmap db-rm <alias> [--purge-file]` (alias: `db-delete`) | Unbinds an alias from registry, optionally deleting the underlying `.db` file. |

---

### 2.2 Alias Binding: `gitmap db-rest "<db-path>" as <alias>`

#### 2.2.1 Syntax & Flag Rules
```bash
gitmap db-rest "<path-to-sqlite-db>" as <alias-name>
```

- **Positional Tokens**:
  1. `<path-to-sqlite-db>`: Quoted or unquoted filesystem path to a `.db`, `.sqlite`, or `.sqlite3` file. Relative paths are cleaned and resolved relative to current working directory.
  2. `as`: Literal case-insensitive keyword token (`as` or `AS`).
  3. `<alias-name>`: Alphanumeric slug representing the alias. May contain hyphens and underscores (`^[a-zA-Z0-9_-]+$`).
- **Parsing Invariant**:
  When the CLI argument list contains the token `as` (case-insensitive) between two non-flag arguments, `cmddb.RunDbRest` routes directly to alias registration mode rather than starting the HTTP daemon.

#### 2.2.2 Execution Flow
1. Verify target file exists on disk. If file does not exist, return `appfault.NewNotFoundError("database file not found: <path>")`.
2. Call `cmddb.RegisterAlias(aliasName, cleanPath, "User registered via gitmap db-rest")`.
3. Read file metadata (size in bytes, table count via `SELECT count(*) FROM sqlite_master WHERE type='table'`).
4. Output confirmation banner:
   ```text
   [SUCCESS] Registered database alias 'analytics' -> ./data/warehouse.db (2.4 MB, 14 tables)
   You can now start the REST server with: gitmap db-rest analytics 8080
   ```

---

### 2.3 Starting REST Server: `gitmap db-rest <targets> [ports]`

#### 2.3.1 Syntax Variations
```bash
# 1. Single alias with comma-separated port fallback
gitmap db-rest analytics 8080,8081,3000

# 2. Direct path with default ports (8080,8081,8082,3000)
gitmap db-rest ./data/cache.db

# 3. Multiple registered aliases
gitmap db-rest master,pipeline,tasks 8080,8081

# 4. Mixed aliases and raw database paths
gitmap db-rest master,./custom/store.db 9000,9001
```

#### 2.3.2 Target Resolver Rules
The first positional argument `<targets>` is split on commas `,` into individual target identifiers:
- For each token $T$:
  1. If $T$ matches a known system alias (`master`, `pipeline`, `tasks`, `heatmap`), resolve to its standard path in `store` or `pipelinedb`.
  2. If $T$ exists in the `DatabaseAlias` table in `gitmap.db`, resolve to registered `DatabasePath`.
  3. If $T$ is an existing file path, synthesize an alias slug from the base filename (e.g. `./custom/inventory.db` becomes alias `inventory`). If alias collision occurs, append a deterministic hash suffix (`inventory-a1b2`).
  4. If $T$ cannot be found on disk or in the registry, return `appfault.NewNotFoundError("database target not found: " + T)`.

#### 2.3.3 Port Probing & Sequential Fallback Engine
The port argument can be a single port or a comma-separated list of candidate ports (e.g. `8080,8081,3000`). If omitted, defaults to `8080,8081,8082,3000,5000`.

- **Probing Algorithm (`rest_port_picker.go`)**:
  ```go
  func PickAvailablePort(candidates []int) (int, error) {
      for _, port := range candidates {
          ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
          if err == nil {
              _ = ln.Close()
              return port, nil
          }
      }
      return 0, fmt.Errorf("all candidate ports [%s] are in use", formatPorts(candidates))
  }
  ```
- **Fallback Logging**:
  If port `8080` is in use, log to stdout:
  ```text
  [PORT] Port 8080 is already occupied. Falling back to next candidate: 8081...
  [PORT] Successfully bound to port 8081.
  ```

#### 2.3.4 Terminal Startup Banner
Upon binding and launching the HTTP server, the CLI outputs a structured ASCII banner:
```text
================================================================================
                    GITMAP DATABASE REST SERVER (v28)
================================================================================
 Server Listening: http://localhost:8081
 Mounted Databases (3):
   * master    -> store/gitmap.db (1.2 MB)
   * pipeline  -> repodb/pipeline.db (450 KB)
   * analytics -> data/warehouse.db (2.4 MB)

 Endpoints:
   * Catalog API:        http://localhost:8081/api/v1/databases
   * Master Tables:      http://localhost:8081/api/v1/db/master/tables
   * Interactive UI:     http://localhost:8081/ui
   * Swagger / OpenAPI:  http://localhost:8081/swagger
================================================================================
 Press Ctrl+C to stop the server.
```

---

### 2.4 Record Deletion CLI: `gitmap db-record-del` / `dbrd`

#### 2.4.1 Syntax Rules
```bash
# Standard form
gitmap db-record-del <alias-or-path> <table> <id>

# Flag form
gitmap db-record-del <alias-or-path> -table <table> -id <id> [-y]

# Alias short-code
gitmap dbrd tasks Subtask 42 -y
```

#### 2.4.2 Parameter Specification
- `<alias-or-path>`: Target database alias or file path.
- `<table>`: Target SQLite table name.
- `<id>`: Primary key value (integer or string).
- `-y`, `--yes`, `--force`: Non-interactive bypass flag.

#### 2.4.3 Primary Key Resolution & Execution Logic
1. Resolve database path and open single connection.
2. Query `PRAGMA table_info(<table>)` to discover the primary key column name:
   - Identify column where `pk = 1` (or `pk > 0`).
   - If no primary key is explicitly defined, check for `rowid` / `_rowid_` or return `appfault.NewInvalidInputError("table has no primary key")`.
3. Fetch existing record via `SELECT * FROM <table> WHERE <pkCol> = ?`. If not found, return `appfault.NewNotFoundError("record not found")`.
4. If stdin is a TTY and `-y` flag is absent:
   ```text
   [CONFIRM] Delete record #42 from table 'Subtask' in database 'tasks'? [y/N]: 
   ```
5. Execute `DELETE FROM <table> WHERE <pkCol> = ?`.
6. Output success summary:
   ```text
   [SUCCESS] Deleted 1 record (SubtaskId=42) from table 'Subtask' in database 'tasks'.
   ```

---

### 2.5 Database / Alias Removal CLI: `gitmap db-rm` / `db-delete`

#### 2.5.1 Syntax Rules
```bash
# Remove alias from registry
gitmap db-rm <alias-name> [-y]

# Remove alias AND delete database file from disk
gitmap db-rm <alias-name> --purge-file [-y]

# Long-form command name
gitmap db-delete custom-db --purge-file
```

#### 2.5.2 Safeguards & Protection Rules
- **System Alias Immutability**:
  Aliases `master`, `gitmap`, `pipeline`, `pipelines`, `tasks`, `agents`, `heatmap` are protected system databases. Attempting `gitmap db-rm master` MUST fail with error:
  `appfault.NewForbiddenError("system database alias '<name>' cannot be removed or deleted")`.
- **Physical File Purge Guard**:
  When `--purge-file` is passed, the CLI requires explicit confirmation unless `-y` is supplied:
  ```text
  [WARNING] This will permanently delete the SQLite database file: ./data/test.db
  Are you sure you want to proceed? [y/N]:
  ```
- **Execution Flow**:
  1. Remove row from `DatabaseAlias` in `gitmap.db`.
  2. If `--purge-file` is requested, call `os.Remove(resolvedPath)` along with any corresponding WAL and SHM journal files (`<path>-wal`, `<path>-shm`).
  3. Output status message.

---

## 3. Universal Response Envelope Specification

Every JSON response returned by the REST API MUST adhere to the Universal Response Envelope. No bare arrays, raw primitive values, or irregular JSON payloads are permitted.

### 3.1 JSON Schema Definition

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "UniversalResponseEnvelope",
  "type": "object",
  "required": ["Status", "Attributes", "Results"],
  "properties": {
    "Status": {
      "type": "object",
      "required": ["IsSuccess", "IsFailed", "Code", "Message", "Timestamp"],
      "properties": {
        "IsSuccess": { "type": "boolean" },
        "IsFailed": { "type": "boolean" },
        "Code": { "type": "integer" },
        "Message": { "type": "string" },
        "Timestamp": { "type": "string", "format": "date-time" }
      }
    },
    "Attributes": {
      "type": "object",
      "required": ["RequestedAt", "HasAnyErrors", "IsSingle", "IsMultiple", "IsEmpty"],
      "properties": {
        "RequestedAt": { "type": "string" },
        "RequestDelegatedAt": { "type": "string" },
        "HasAnyErrors": { "type": "boolean" },
        "IsSingle": { "type": "boolean" },
        "IsMultiple": { "type": "boolean" },
        "IsEmpty": { "type": "boolean" },
        "TotalRecords": { "type": "integer" },
        "DatabaseAlias": { "type": "string" },
        "TableName": { "type": "string" },
        "ExecutionTimeMs": { "type": "integer" },
        "PerPage": { "type": "integer" },
        "CurrentPage": { "type": "integer" },
        "TotalPages": { "type": "integer" }
      },
      "additionalProperties": true
    },
    "Results": {
      "description": "Array of records, single record object, or execution result object.",
      "type": ["array", "object", "null"]
    }
  }
}
```

### 3.2 Standard Success Envelope Example (List Resource)
```json
{
  "Status": {
    "IsSuccess": true,
    "IsFailed": false,
    "Code": 200,
    "Message": "OK",
    "Timestamp": "2026-10-20T10:15:30Z"
  },
  "Attributes": {
    "RequestedAt": "/api/v1/db/tasks/Subtask?status=IN_PROGRESS",
    "RequestDelegatedAt": "",
    "HasAnyErrors": false,
    "IsSingle": false,
    "IsMultiple": true,
    "IsEmpty": false,
    "TotalRecords": 2,
    "DatabaseAlias": "tasks",
    "TableName": "Subtask",
    "ExecutionTimeMs": 4,
    "PerPage": 20,
    "CurrentPage": 1,
    "TotalPages": 1
  },
  "Results": [
    {
      "SubtaskId": 101,
      "TaskId": "database-rest-command-implementation",
      "Title": "CLI Commands & Port Fallback",
      "Status": "IN_PROGRESS",
      "CreatedAt": "2026-10-20T08:00:00Z"
    },
    {
      "SubtaskId": 102,
      "TaskId": "database-rest-command-implementation",
      "Title": "REST Engine & Endpoints",
      "Status": "IN_PROGRESS",
      "CreatedAt": "2026-10-20T08:05:00Z"
    }
  ]
}
```

### 3.3 Standard Error Envelope Example (404 Not Found)
```json
{
  "Status": {
    "IsSuccess": false,
    "IsFailed": true,
    "Code": 404,
    "Message": "Table 'NonExistent' not found in database 'tasks'",
    "Timestamp": "2026-10-20T10:15:35Z"
  },
  "Attributes": {
    "RequestedAt": "/api/v1/db/tasks/NonExistent",
    "RequestDelegatedAt": "",
    "HasAnyErrors": true,
    "IsSingle": false,
    "IsMultiple": false,
    "IsEmpty": true,
    "TotalRecords": 0,
    "DatabaseAlias": "tasks",
    "ExecutionTimeMs": 1
  },
  "Results": null
}
```

---

## 4. REST API Route Catalog & Operation Contracts

All routes are mounted under the base path `/api/v1/`.

```
REST API Route Hierarchy:
/api/v1/
├── databases                                 [GET]  Catalog of mounted databases
└── db/{alias}/
    ├── tables                                [GET]  Table & view summary
    ├── schema                                [GET]  Comprehensive schema & DDL
    ├── find-columns                          [GET]  Column locator across tables
    ├── query                                 [GET, POST] Direct SQL query execution
    ├── views                                 [GET, POST] Dynamic view management
    │   └── {viewName}                        [DELETE] Drop dynamic view
    ├── tables/create                         [POST] DDL table creation
    ├── alter                                 [POST] DDL alteration (add col, rename)
    ├── migrate                               [POST] Schema migration execution
    ├── clone                                 [POST] Database cloning & registration
    └── {table}                               [GET, POST] Table collection CRUD
        └── {id}                              [GET, PUT, DELETE] Row instance CRUD
```

---

### 4.1 Database Catalog Endpoint

#### `GET /api/v1/databases`
Returns all mounted databases in the active REST server instance with physical file stats and table counts.

- **Request Headers**: `Accept: application/json`
- **Response Status**: `200 OK`
- **Sample Response**:
  ```json
  {
    "Status": {
      "IsSuccess": true,
      "IsFailed": false,
      "Code": 200,
      "Message": "OK",
      "Timestamp": "2026-10-20T10:20:00Z"
    },
    "Attributes": {
      "RequestedAt": "/api/v1/databases",
      "RequestDelegatedAt": "",
      "HasAnyErrors": false,
      "IsSingle": false,
      "IsMultiple": true,
      "IsEmpty": false,
      "TotalRecords": 2
    },
    "Results": [
      {
        "Alias": "master",
        "DatabasePath": "store/gitmap.db",
        "DatabaseType": "system-master",
        "FileSizeBytes": 1245184,
        "TableCount": 8,
        "IsWritable": true
      },
      {
        "Alias": "tasks",
        "DatabasePath": "repodb/tasks.db",
        "DatabaseType": "system-tasks",
        "FileSizeBytes": 327680,
        "TableCount": 4,
        "IsWritable": true
      }
    ]
  }
  ```

---

### 4.2 Database Schema & Table Listing Endpoints

#### 4.2.1 `GET /api/v1/db/{alias}/tables`
Returns a list of all user-defined tables and views inside the target database.

- **URL Parameters**:
  - `alias`: The database alias slug (e.g. `master`, `tasks`).
- **Response Status**: `200 OK` (or `404 Not Found` if alias unknown)
- **Sample Results**:
  ```json
  "Results": [
    {
      "TableName": "Repository",
      "Type": "table",
      "RowCount": 142,
      "ColumnCount": 12,
      "HasPrimaryKey": true
    },
    {
      "TableName": "ActiveReposView",
      "Type": "view",
      "RowCount": 38,
      "ColumnCount": 5,
      "HasPrimaryKey": false
    }
  ]
  ```

#### 4.2.2 `GET /api/v1/db/{alias}/schema`
Returns comprehensive column details, data types, nullability, defaults, primary keys, and indexes for all tables or a specific table.

- **Query Parameters**:
  - `table` (optional): Filter schema introspection to a single table name.
- **Introspection Technique**:
  Scans `sqlite_master` for SQL DDL, runs `PRAGMA table_info(<table>)` for column properties, and `PRAGMA index_list(<table>)` for indexes.
- **Sample Results**:
  ```json
  "Results": [
    {
      "TableName": "Repository",
      "Ddl": "CREATE TABLE Repository (RepositoryId INTEGER PRIMARY KEY AUTOINCREMENT, Name TEXT NOT NULL, ...);",
      "Columns": [
        {
          "ColumnId": 0,
          "Name": "RepositoryId",
          "DataType": "INTEGER",
          "IsNullable": false,
          "DefaultValue": null,
          "IsPrimaryKey": true
        },
        {
          "ColumnId": 1,
          "Name": "Name",
          "DataType": "TEXT",
          "IsNullable": false,
          "DefaultValue": "",
          "IsPrimaryKey": false
        }
      ],
      "Indexes": [
        {
          "Name": "idx_repo_name",
          "IsUnique": true,
          "Columns": ["Name"]
        }
      ]
    }
  ]
  ```

---

### 4.3 Cross-Table Column Search Endpoint

#### `GET /api/v1/db/{alias}/find-columns`
Searches all tables in the specified database for column names matching a query term or pattern.

- **Query Parameters**:
  - `name` (required): Target column name or substring (e.g. `?name=status` or `?name=created`).
  - `exact` (optional, boolean): If `true`, matches exact column name case-insensitively; if `false` (default), matches substring.
- **Sample Results**:
  ```json
  "Results": [
    {
      "TableName": "Task",
      "ColumnName": "Status",
      "DataType": "TEXT",
      "IsPrimaryKey": false,
      "IsNullable": false
    },
    {
      "TableName": "Subtask",
      "ColumnName": "Status",
      "DataType": "TEXT",
      "IsPrimaryKey": false,
      "IsNullable": false
    },
    {
      "TableName": "PipelineRun",
      "ColumnName": "StatusName",
      "DataType": "TEXT",
      "IsPrimaryKey": false,
      "IsNullable": true
    }
  ]
  ```

---

### 4.4 Record CRUD Operations

Provides automated REST resources for any SQLite table without requiring manual route code.

#### 4.4.1 `GET /api/v1/db/{alias}/{table}` (List Records)
- **Query Parameters**:
  - `limit` (default: 50, max: 1000): Number of rows to return.
  - `offset` (default: 0): Row offset for pagination.
  - `sort` (optional): Column name to order by.
  - `order` (optional, `asc` or `desc`, default: `asc`).
  - `{columnName}={val}`: Any additional query parameter is automatically treated as an equality filter (`WHERE columnName = ?`).
- **Response**: List of row objects with PascalCase column names. `Attributes.TotalRecords` reflects total matching rows.

#### 4.4.2 `GET /api/v1/db/{alias}/{table}/{id}` (Get Single Record)
- **URL Parameters**:
  - `id`: Primary key value.
- **Response**: Single record object under `Results`. Returns `404` if not found.

#### 4.4.3 `POST /api/v1/db/{alias}/{table}` (Create Record)
- **Request Body**: JSON object representing the column-value map.
  ```json
  {
    "Title": "Implement port probing engine",
    "TaskId": "database-rest-command-implementation",
    "Status": "PENDING"
  }
  ```
- **Execution**:
  1. Validates table exists and column names match table columns.
  2. Constructs parameterized `INSERT INTO <table> (...) VALUES (...)`.
  3. Executes insert and retrieves last insert row ID via `res.LastInsertId()`.
  4. Fetches and returns the freshly created record.
- **Response Status**: `201 Created`

#### 4.4.4 `PUT /api/v1/db/{alias}/{table}/{id}` (Update Record)
- **Request Body**: JSON object of columns to update.
  ```json
  {
    "Status": "DONE",
    "CompletedAt": "2026-10-20T10:30:00Z"
  }
  ```
- **Execution**:
  1. Locates primary key column via `PRAGMA table_info`.
  2. Builds parameterized `UPDATE <table> SET col1=?, col2=? WHERE <pkCol>=?`.
  3. Executes update inside single connection.
  4. Returns the updated record object.
- **Response Status**: `200 OK`

#### 4.4.5 `DELETE /api/v1/db/{alias}/{table}/{id}` (Delete Record)
- **Execution**:
  Executes `DELETE FROM <table> WHERE <pkCol> = ?`.
- **Response Status**: `200 OK` with execution summary in `Results`:
  ```json
  "Results": {
    "Deleted": true,
    "RowsAffected": 1,
    "PrimaryKey": "SubtaskId",
    "Id": 101
  }
  ```

---

### 4.5 Dynamic Views Management (Authoring, Listing, Deleting)

Allows clients to author, list, and destroy complex multi-table SQL views on the fly.

#### 4.5.1 `POST /api/v1/db/{alias}/views` (Create Dynamic View)
- **Request Body**:
  ```json
  {
    "ViewName": "ActiveSubtasksWithParent",
    "Query": "SELECT s.SubtaskId, s.Title, s.Status, t.TaskName, t.Owner FROM Subtask s INNER JOIN Task t ON s.TaskId = t.TaskId WHERE s.Status != 'DONE'"
  }
  ```
- **Validation**:
  1. `ViewName` must be a valid SQLite identifier (`^[a-zA-Z0-9_]+$`).
  2. `Query` must be a read-only `SELECT` statement. Dangerous statements (`INSERT`, `UPDATE`, `DELETE`, `DROP`, `ALTER`, `ATTACH`) are strictly forbidden and rejected with `400 Bad Request`.
- **Execution**:
  1. Executes `CREATE VIEW IF NOT EXISTS <ViewName> AS <Query>`.
  2. Records view definition into the Split-DB cache table `RestViewMetadata`.
- **Response Status**: `201 Created`

#### 4.5.2 `GET /api/v1/db/{alias}/views` (List Dynamic Views)
- **Response**: List of all views in the database with their underlying SQL definitions.

#### 4.5.3 `DELETE /api/v1/db/{alias}/views/{name}` (Drop View)
- **Execution**: Executes `DROP VIEW IF EXISTS <name>`.
- **Response Status**: `200 OK`

---

### 4.6 Direct SQL Query Execution Endpoints

Enables ad-hoc SQL execution directly from HTTP clients, scripts, or the interactive UI.

#### 4.6.1 `GET /api/v1/db/{alias}/query?sql={sqlQuery}`
Convenience endpoint for read-only queries passed as a URL query parameter.
- **Query Parameters**:
  - `sql` (required, URL-encoded string): Read-only SQL query (e.g. `SELECT count(*), status FROM Task GROUP BY status`).
- **Safety Restriction**:
  Only statements beginning with `SELECT`, `PRAGMA`, or `EXPLAIN` are permitted. Any mutation statement returns `403 Forbidden: GET query endpoint only permits read-only SQL statements. Use POST for mutations.`
- **Sample Results**:
  ```json
  "Results": [
    { "count(*)": 12, "Status": "DONE" },
    { "count(*)": 4, "Status": "IN_PROGRESS" }
  ]
  ```

#### 4.6.2 `POST /api/v1/db/{alias}/query`
General SQL execution endpoint supporting parameterized queries, joins, and DML statements.
- **Request Body**:
  ```json
  {
    "Sql": "UPDATE Task SET Status = ? WHERE TaskId = ?",
    "Params": ["IN_PROGRESS", "task-99"]
  }
  ```
- **Response**:
  For `SELECT` queries, returns row objects array. For DML (`UPDATE`, `INSERT`, `DELETE`), returns execution metadata:
  ```json
  "Results": {
    "RowsAffected": 1,
    "LastInsertId": 0,
    "ExecutionDurationMs": 2
  }
  ```

---

### 4.7 Table Alteration & DDL Endpoints

#### 4.7.1 `POST /api/v1/db/{alias}/tables/create` (Create New Table)
- **Request Body Options**:
  - Structured form:
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
  - Raw DDL form:
    ```json
    {
      "Ddl": "CREATE TABLE AuditLog (AuditId INTEGER PRIMARY KEY AUTOINCREMENT, Action TEXT NOT NULL, Payload TEXT, CreatedAt TEXT DEFAULT CURRENT_TIMESTAMP);"
    }
    ```
- **Response Status**: `201 Created`

#### 4.7.2 `POST /api/v1/db/{alias}/alter` (Alter Table Schema)
Supports adding columns, renaming columns, and renaming tables.
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

### 4.8 Database Migration Engine & Split-DB Caching

Provides an automated database migration runner that stores migration records in a dedicated caching Split-DB to track applied schema versions.

#### `POST /api/v1/db/{alias}/migrate`
- **Request Body**:
  ```json
  {
    "MigrationName": "002_add_index_subtask_status",
    "UpSql": "CREATE INDEX IF NOT EXISTS idx_subtask_status ON Subtask(Status);",
    "DownSql": "DROP INDEX IF EXISTS idx_subtask_status;"
  }
  ```
- **Execution & Split-DB Cache Storage**:
  1. Open Split-DB cache database (`repodb/rest_cache.db`).
  2. Ensure migration log table exists:
     ```sql
     CREATE TABLE IF NOT EXISTS RestMigrationLog (
       MigrationId INTEGER PRIMARY KEY AUTOINCREMENT,
       DatabaseAlias TEXT NOT NULL,
       DatabasePath TEXT NOT NULL,
       MigrationName TEXT NOT NULL,
       UpSql TEXT NOT NULL,
       DownSql TEXT NOT NULL,
       AppliedAt TEXT NOT NULL,
       ExecutionDurationMs INTEGER NOT NULL
     );
     ```
  3. Check if `MigrationName` was already applied for this alias. If so, return `409 Conflict: migration already applied`.
  4. Begin transaction on target database, execute `UpSql`, commit transaction.
  5. Insert record into `RestMigrationLog` in the Split-DB cache.
- **Response Status**: `200 OK`

---

### 4.9 Database Cloning Endpoint

#### `POST /api/v1/db/{alias}/clone`
Creates an exact physical copy of an existing mounted SQLite database, optionally registers it under a new alias, and mounts it into the running REST server dynamically.

- **Request Body**:
  ```json
  {
    "DestinationPath": "./data/clones/tasks_backup.db",
    "NewAlias": "tasks-backup",
    "AutoMount": true
  }
  ```
- **Execution**:
  1. Resolve source database path.
  2. Ensure parent directory of `DestinationPath` exists.
  3. Acquire read lock on source database or use SQLite backup API (`VACUUM INTO ?`).
  4. Register `NewAlias` in `DatabaseAlias` registry table in `gitmap.db`.
  5. If `AutoMount` is `true`, add the new alias connection directly to the active `MultiDbManager` in the running REST server.
- **Sample Results**:
  ```json
  "Results": {
    "SourceAlias": "tasks",
    "NewAlias": "tasks-backup",
    "DestinationPath": "data/clones/tasks_backup.db",
    "FileSizeBytes": 327680,
    "IsMounted": true,
    "ApiBaseUrl": "/api/v1/db/tasks-backup"
  }
  ```

---

## 5. Concurrency & Deadlock Prevention Invariants

To avoid any SQLite deadlocks, file locks, or goroutine starvation:
1. **Single Connection per DB**: Every mounted database pool uses `conn.SetMaxOpenConns(1)`. SQLite operates in serial write mode.
2. **WAL Journaling**: Each connection runs `PRAGMA journal_mode=WAL; PRAGMA busy_timeout=5000; PRAGMA foreign_keys=ON;` immediately upon opening.
3. **Full Row Ingestion**: In all handler functions, rows returned by `conn.Query` or `conn.QueryContext` MUST be scanned and appended into an in-memory slice (`[]map[string]any`) in a loop, and `rows.Close()` MUST be called before any further queries or network responses are executed.
4. **Mutex Guarding on Schema Mutations**: The `MultiDbManager` maintains a `sync.RWMutex` to guard dynamic registration of new aliases (such as after `/clone` operations) without race conditions.

---

## 6. Specification Traceability & Verification Matrix

| Invariant / Feature | CLI Test Target | HTTP API Test Target | Acceptance Criteria |
|---|---|---|---|
| Universal Response Envelope | N/A | All `/api/v1/...` routes | Root keys `Status`, `Attributes`, `Results` present; all keys PascalCase. |
| Port Fallback Engine | `gitmap db-rest tasks 8080,8081` | N/A | Occupy 8080 with dummy listener; server binds to 8081 without crash. |
| Multi-DB Mounting | `gitmap db-rest master,tasks` | `GET /api/v1/databases` | Returns 2 items matching `master` and `tasks`. |
| Cross-Table Column Search | N/A | `GET /api/v1/db/{alias}/find-columns?name=Status` | Returns matching columns across all tables in alias. |
| View Creation & Query | N/A | `POST /api/v1/db/{alias}/views` | View created; queryable via `GET /api/v1/db/{alias}/{viewName}`. |
| Direct SQL Query | N/A | `GET /api/v1/db/{alias}/query?sql=...` | Returns JSON rows matching query; blocks mutations on GET. |
| Record Deletion CLI | `gitmap db-record-del tasks Subtask 101 -y` | `DELETE /api/v1/db/{alias}/{table}/{id}` | Deletes record; returns count = 1. |
| Database Removal CLI | `gitmap db-rm test-db --purge-file -y` | N/A | Removes from `DatabaseAlias` and deletes file from disk. |
