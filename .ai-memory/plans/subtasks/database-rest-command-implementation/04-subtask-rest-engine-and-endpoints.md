# Subtask 04: REST Multi-DB Engine, CRUD, DDL, Views & SQL Querying

**Subtask ID:** `04-subtask-rest-engine-and-endpoints`  
**Parent Task:** `database-rest-command-implementation`  
**Master Plan:** `.ai-memory/plans/database-rest-command-implementation.md`  
**Specification Reference:** `02-spec/21-app/database-rest-command-implementation/02-cli-and-rest-api-spec.md`  
**Assigned Role:** Worker 02  
**Status:** PENDING  
**Concurrency Budget:** A = 3, H = 2, Capacity = 6  

---

## 1. Subtask Overview & Scope

Worker 02 is responsible for architecting and implementing the core HTTP REST API engine, Universal Response Envelope pipeline, dynamic table CRUD controllers, DDL and alteration runners, dynamic view compiler, SQL execution handlers, and Split-DB migration caching subsystem.

### Owned Source Files (Strictly Disjoint):
- `cli/cmddb/rest_server.go`: Server lifecycle, `net/http` route multiplexer, CORS middleware, graceful shutdown.
- `cli/cmddb/rest_envelope.go`: Universal Response Envelope model, PascalCase JSON serializer, success/error writing helpers.
- `cli/cmddb/rest_multi_db_manager.go`: Thread-safe database connection pool manager, dynamic alias resolver, dynamic runtime mounter.
- `cli/cmddb/rest_handlers_catalog.go`: Database catalog endpoint (`GET /api/v1/databases`).
- `cli/cmddb/rest_handlers_schema.go`: Schema, table listing, and cross-table column search endpoints (`/api/v1/db/{alias}/tables`, `/schema`, `/find-columns`).
- `cli/cmddb/rest_handlers_crud.go`: Dynamic table record CRUD handlers (`GET`, `POST`, `PUT`, `DELETE` on `/api/v1/db/{alias}/{table}` and `/{id}`).
- `cli/cmddb/rest_handlers_views.go`: Dynamic view creation, validation, listing, and dropping (`/api/v1/db/{alias}/views`).
- `cli/cmddb/rest_handlers_query.go`: Direct SQL execution (`GET` read-only and `POST` parameterized query / DML on `/api/v1/db/{alias}/query`).
- `cli/cmddb/rest_handlers_ddl.go`: DDL table creation and schema alteration endpoints (`/api/v1/db/{alias}/tables/create`, `/alter`).
- `cli/cmddb/rest_handlers_migrate.go`: Database schema migration runner endpoint (`/api/v1/db/{alias}/migrate`).
- `cli/cmddb/rest_handlers_clone.go`: Database cloning and dynamic mounting endpoint (`/api/v1/db/{alias}/clone`).
- `cli/cmddb/rest_migration_cache.go`: Split-DB cache manager storing migration history and view definitions in `repodb/rest_cache.db`.

---

## 2. Invariants & Coding Standards

1. **Strict Relative Git Paths**: All comments, documentation, and log entries MUST use relative paths (e.g. `02-spec/21-app/...`, `cli/cmddb/...`, `.ai-memory/...`). Absolutely NO absolute filesystem paths or `file:///` URIs.
2. **Universal Response Envelope**: Every REST endpoint MUST return JSON with root keys `{ Status, Attributes, Results }` using PascalCase keys throughout.
3. **Pure Go CGO-Free SQLite**: Operates strictly on `modernc.org/sqlite`. Zero external CGO libraries.
4. **Single-Writer Safety**: Every database connection MUST use `conn.SetMaxOpenConns(1)` and PRAGMAs:
   `PRAGMA journal_mode=WAL; PRAGMA busy_timeout=5000; PRAGMA foreign_keys=ON;`
5. **Buffer-Before-Next Rule**: In all queries, scan all rows into an in-memory `[]map[string]any` buffer before closing rows or executing subsequent operations. Never hold active `*sql.Rows` iterators across operations.
6. **File Size Cap**: Each Go source file MUST be strictly $\le 300$ lines.

---

## 3. Step-by-Step Implementation Instructions

### Step 1: Universal Response Envelope (`cli/cmddb/rest_envelope.go`)
- Define the Go models matching the Universal Envelope specification:
  ```go
  type EnvelopeStatus struct {
      IsSuccess bool   `json:"IsSuccess"`
      IsFailed  bool   `json:"IsFailed"`
      Code      int    `json:"Code"`
      Message   string `json:"Message"`
      Timestamp string `json:"Timestamp"`
  }

  type EnvelopeAttributes struct {
      RequestedAt        string `json:"RequestedAt"`
      RequestDelegatedAt string `json:"RequestDelegatedAt,omitempty"`
      HasAnyErrors       bool   `json:"HasAnyErrors"`
      IsSingle           bool   `json:"IsSingle"`
      IsMultiple         bool   `json:"IsMultiple"`
      IsEmpty            bool   `json:"IsEmpty"`
      TotalRecords       int    `json:"TotalRecords,omitempty"`
      DatabaseAlias      string `json:"DatabaseAlias,omitempty"`
      TableName          string `json:"TableName,omitempty"`
      ExecutionTimeMs    int64  `json:"ExecutionTimeMs,omitempty"`
      PerPage            int    `json:"PerPage,omitempty"`
      CurrentPage        int    `json:"CurrentPage,omitempty"`
      TotalPages         int    `json:"TotalPages,omitempty"`
  }

  type ResponseEnvelope struct {
      Status     EnvelopeStatus     `json:"Status"`
      Attributes EnvelopeAttributes `json:"Attributes"`
      Results    any                `json:"Results"`
  }
  ```
- Implement writing helpers:
  - `WriteSuccess(w http.ResponseWriter, r *http.Request, code int, results any, attrs EnvelopeAttributes)`
  - `WriteError(w http.ResponseWriter, r *http.Request, code int, message string, alias string)`

### Step 2: Multi-DB Connection Manager (`cli/cmddb/rest_multi_db_manager.go`)
- Create `MultiDbManager` to manage active database connections across aliases:
  ```go
  type MultiDbManager struct {
      mu       sync.RWMutex
      targets  map[string]DbMountTarget
      pools    map[string]*sql.DB
  }
  ```
- Implement methods:
  - `NewMultiDbManager(targets []DbMountTarget) (*MultiDbManager, error)`: Opens each database connection with `_busy_timeout=5000&_journal_mode=WAL` and `conn.SetMaxOpenConns(1)`.
  - `GetDB(alias string) (*sql.DB, DbMountTarget, error)`: Looks up database by alias (case-insensitive).
  - `ListTargets() []DbMountTarget`: Returns snapshot of all mounted database targets.
  - `RegisterDynamicTarget(alias string, path string) error`: Safely mounts a newly cloned or attached database at runtime with mutex locking.
  - `CloseAll()`: Closes all database connections gracefully during server shutdown.

### Step 3: HTTP Server Engine & Route Multiplexer (`cli/cmddb/rest_server.go`)
- Create `MultiDbRestServer` using standard library `net/http`:
  - Registers route handlers for all endpoints.
  - Implements route matching for path patterns (`/api/v1/db/{alias}/...`).
  - Attaches standard middleware: CORS headers (`Access-Control-Allow-Origin: *`, allowed headers and methods), request timing tracking (`ExecutionTimeMs`), and panic recovery.
  - Exposes `Start(port int) error` and `Shutdown(ctx context.Context) error`.

### Step 4: Catalog & Schema Introspection Handlers (`rest_handlers_catalog.go` & `rest_handlers_schema.go`)
- `HandleGetDatabases`: Returns `GET /api/v1/databases` listing all mounted aliases with live stats (size, table count, write access).
- `HandleGetTables`: Returns `GET /api/v1/db/{alias}/tables` querying `sqlite_master` for tables and views with row counts.
- `HandleGetSchema`: Returns `GET /api/v1/db/{alias}/schema` executing `PRAGMA table_info` and `PRAGMA index_list` for complete column data types, nullability, defaults, and keys.
- `HandleFindColumns`: Returns `GET /api/v1/db/{alias}/find-columns?name=...` searching for column names across all tables in `{alias}`.

### Step 5: Dynamic Record CRUD Handlers (`cli/cmddb/rest_handlers_crud.go`)
- `HandleListRecords`: `GET /api/v1/db/{alias}/{table}`
  - Parses `limit` (default 50), `offset` (default 0), `sort`, `order`.
  - Parses arbitrary query parameters as `WHERE col = ?` equality filters.
  - Queries rows, scans all columns dynamically into `[]map[string]any` with PascalCase keys, and returns count in `Attributes`.
- `HandleGetRecord`: `GET /api/v1/db/{alias}/{table}/{id}`
  - Discovers primary key column and executes `SELECT * FROM <table> WHERE <pkCol> = ?`.
  - Returns `404 Not Found` if record does not exist.
- `HandleCreateRecord`: `POST /api/v1/db/{alias}/{table}`
  - Deserializes JSON body into key-value map.
  - Validates column names against table schema.
  - Constructs and executes `INSERT INTO <table> (...) VALUES (...)`.
  - Returns `201 Created` with newly created record including generated PK.
- `HandleUpdateRecord`: `PUT /api/v1/db/{alias}/{table}/{id}`
  - Updates specified columns for record matching primary key.
  - Returns `200 OK` with updated record.
- `HandleDeleteRecord`: `DELETE /api/v1/db/{alias}/{table}/{id}`
  - Deletes row matching primary key and returns execution summary.

### Step 6: Dynamic Views & SQL Handlers (`rest_handlers_views.go` & `rest_handlers_query.go`)
- `HandleCreateView`: `POST /api/v1/db/{alias}/views`
  - Validates view name and verifies SQL statement is a safe read-only `SELECT`.
  - Executes `CREATE VIEW <name> AS <query>`.
  - Caches view metadata in Split-DB cache.
- `HandleListViews`: `GET /api/v1/db/{alias}/views`
  - Queries `sqlite_master WHERE type='view'`.
- `HandleDeleteView`: `DELETE /api/v1/db/{alias}/views/{name}`
  - Executes `DROP VIEW IF EXISTS <name>`.
- `HandleGetQuery`: `GET /api/v1/db/{alias}/query?sql=...`
  - Permits only read-only `SELECT` / `PRAGMA` / `EXPLAIN`.
  - Returns row array.
- `HandlePostQuery`: `POST /api/v1/db/{alias}/query`
  - Supports parameterized queries with `"Sql"` and `"Params"`.
  - Returns rows for queries or rows affected / last insert ID for mutations.

### Step 7: DDL, Migration & Cloning Handlers (`rest_handlers_ddl.go`, `rest_handlers_migrate.go`, `rest_handlers_clone.go`)
- `HandleCreateTable`: `POST /api/v1/db/{alias}/tables/create`
  - Creates table from structured column array or raw DDL string.
- `HandleAlterTable`: `POST /api/v1/db/{alias}/alter`
  - Supports `ADD_COLUMN`, `RENAME_COLUMN`, `RENAME_TABLE`.
- `HandleMigrate`: `POST /api/v1/db/{alias}/migrate`
  - Executes migration in transaction.
  - Records migration entry in `RestMigrationLog` table in Split-DB cache (`repodb/rest_cache.db`).
- `HandleClone`: `POST /api/v1/db/{alias}/clone`
  - Safely copies or vacuums database into new destination path.
  - Automatically registers new alias in `DatabaseAlias` table.
  - Calls `manager.RegisterDynamicTarget` to immediately mount the new alias in the running server.

### Step 8: Split-DB Caching Engine (`cli/cmddb/rest_migration_cache.go`)
- Open or create `repodb/rest_cache.db`.
- Initialize tables `RestMigrationLog` and `RestViewMetadata`.
- Provide atomic methods:
  - `RecordMigration(alias, path, name, upSql, downSql, durationMs) error`
  - `IsMigrationApplied(alias, name) (bool, error)`
  - `RecordView(alias, name, sqlQuery) error`
  - `RemoveViewRecord(alias, name) error`

---

## 4. Verification Checklist & Acceptance Criteria

- [ ] **Universal Response Envelope**: Every endpoint returns JSON with PascalCase root keys `Status`, `Attributes`, `Results`.
- [ ] **Database Catalog**: `GET /api/v1/databases` lists all mounted aliases accurately.
- [ ] **Dynamic CRUD**: Able to Insert, Read, Update, and Delete rows on any table via HTTP.
- [ ] **View Engine**: Creates dynamic view with multi-table join and queries it successfully.
- [ ] **Direct SQL**: Read-only queries execute via GET; parameterized queries and mutations execute via POST.
- [ ] **DDL & Migration**: Adding column and running migration properly updates schema and logs to `repodb/rest_cache.db`.
- [ ] **Dynamic Cloning**: Cloning creates file on disk, registers alias, and mounts it into running server.
- [ ] **No Deadlocks**: `conn.SetMaxOpenConns(1)` with all rows buffered prevents database deadlocks under concurrent requests.
- [ ] **Line Count Cap**: All files are strictly $\le 300$ lines.
