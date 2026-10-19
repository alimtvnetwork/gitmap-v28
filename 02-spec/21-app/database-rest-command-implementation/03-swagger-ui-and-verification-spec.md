# OpenAPI 3.0 Specification, Interactive Visualizer UI & Verification Spec

> **Document Version:** 1.0.0  
> **Status:** Active / Specification  
> **Task Slug:** `database-rest-command-implementation`  
> **Parent Spec:** `02-spec/21-app/database-rest-command-implementation/readme.md`  
> **Parent Plan:** `.ai-memory/plans/database-rest-command-implementation.md`  
> **Traceability:** Task-02, Subtask-05, Subtask-06  

---

## 1. Executive Summary & Architecture Scope

This specification establishes the technical requirements, architecture, schema synthesis logic, user experience design, and acceptance verification criteria for two foundational subsystems of the GitMap `db-rest` engine:
1. **Dynamic OpenAPI 3.0 Generator & Embedded Swagger UI**: Synthesizes an OpenAPI 3.0.3 specification on-the-fly by inspecting mounted SQLite schemas and delivers an interactive Swagger UI.
2. **Custom Interactive Database Visualizer Web UI**: A dark-theme, responsive developer dashboard for exploring schemas, inspecting tables, editing columns, performing live CRUD operations, running direct SQL, and building dynamic views.
3. **Acceptance Verification & Quality Gate Suite**: An exhaustive, automated verification strategy covering port probing, multi-alias mounting, CRUD operations, DDL alterations, SQL execution, deadlock prevention, and CLI command parity.

### Core Architectural Invariants
- **Dynamic Introspection**: OpenAPI definitions are synthesized dynamically by querying `sqlite_master` and `PRAGMA table_info` for all mounted databases.
- **Zero External CDN Dependencies**: Both the Swagger UI and the custom Visualizer UI are 100% self-contained. All HTML, CSS, JavaScript, icons, and fonts are embedded directly within the Go binary (`cli/cmddb/rest_ui_assets.go`), ensuring full functionality in air-gapped, offline, and enterprise environments.
- **Universal Response Envelope**: All API endpoints return the standardized envelope structure `{ Status, Attributes, Results }` using PascalCase JSON fields.
- **Single-Writer Safety**: SQLite database connections maintain single-writer concurrency (`conn.SetMaxOpenConns(1)`) with all result sets eagerly scanned into memory before releasing connection locks.

---

## 2. Dynamic OpenAPI 3.0 Specification Engine

The dynamic OpenAPI engine generates an OpenAPI 3.0.3 compliant specification JSON document for all mounted SQLite databases upon request at `/swagger/doc.json`.

### 2.1 Introspection Pipeline

```
HTTP GET /swagger/doc.json
  │
  ├── 1. Acquire MultiDbManager Read Lock
  ├── 2. Iterate Mounted Databases (aliases and direct paths)
  ├── 3. For each Database:
  │     ├── Query sqlite_master for Tables and Views
  │     │     `SELECT type, name, sql FROM sqlite_master WHERE type IN ('table', 'view') AND name NOT LIKE 'sqlite_%'`
  │     └── For each Table/View:
  │           `PRAGMA table_info("{tableName}")`
  │           Extract: cid, name, type, notnull, dflt_value, pk
  ├── 4. Map SQLite Data Types to OpenAPI 3.0 Types & Formats
  ├── 5. Synthesize Paths, Parameters, Request Bodies & Component Schemas
  └── 6. Serialize to JSON and Serve with Content-Type: application/json
```

### 2.2 SQLite to OpenAPI 3.0 Type Mapping Matrix

SQLite uses dynamic type affinity. The generator maps declared column types to strict OpenAPI 3.0 data types according to the following normalization table:

| SQLite Declared Type | Normalization Rule | OpenAPI Type | OpenAPI Format | Example Default Value |
|---|---|---|---|---|
| `INTEGER`, `INT`, `BIGINT`, `SMALLINT`, `TINYINT` | Contains `INT` | `integer` | `int64` | `0` |
| `REAL`, `FLOAT`, `DOUBLE`, `DECIMAL`, `NUMERIC` | Contains `REAL`, `FLOA`, `DOUB` | `number` | `double` | `0.0` |
| `TEXT`, `VARCHAR`, `NVARCHAR`, `CHAR`, `CLOB` | Contains `TEXT`, `CHAR`, `CLOB` | `string` | *(none)* | `""` |
| `BOOLEAN`, `BOOL` | Contains `BOOL` | `boolean` | *(none)* | `false` |
| `DATETIME`, `TIMESTAMP`, `DATE` | Contains `DATE`, `TIME` | `string` | `date-time` | `2026-01-01T00:00:00Z` |
| `BLOB` | Contains `BLOB` or Empty | `string` | `byte` | `""` |
| Unspecified / Unknown | Fallback | `string` | *(none)* | `""` |

### 2.3 OpenAPI Document Structure

The synthesized OpenAPI document incorporates all system endpoints and dynamic table endpoints:

```json
{
  "openapi": "3.0.3",
  "info": {
    "title": "GitMap Database REST API",
    "version": "1.0.0",
    "description": "Dynamic REST API and schema explorer generated for mounted SQLite databases."
  },
  "servers": [
    {
      "url": "/",
      "description": "Local GitMap REST Server"
    }
  ],
  "tags": [
    { "name": "System", "description": "Database catalog and server health" },
    { "name": "Schema", "description": "Table and column discovery endpoints" },
    { "name": "Records", "description": "Generic CRUD operations across tables" },
    { "name": "SQL Engine", "description": "Direct SQL queries and batch execution" },
    { "name": "Views", "description": "Dynamic view creation and query endpoints" },
    { "name": "DDL & Migrations", "description": "Table creation, column alterations, and migrations" }
  ],
  "paths": {
    "/api/v1/databases": {
      "get": {
        "tags": ["System"],
        "summary": "List all mounted database aliases and file paths",
        "responses": {
          "200": {
            "description": "Catalog of mounted databases",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/DatabaseCatalogResponse" }
              }
            }
          }
        }
      }
    },
    "/api/v1/db/{alias}/schema": {
      "get": {
        "tags": ["Schema"],
        "summary": "Inspect full schema for database",
        "parameters": [
          { "name": "alias", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "responses": {
          "200": {
            "description": "Complete database schema tree",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/DatabaseSchemaResponse" }
              }
            }
          }
        }
      }
    },
    "/api/v1/db/{alias}/tables": {
      "get": {
        "tags": ["Schema"],
        "summary": "List all tables and row counts in database",
        "parameters": [
          { "name": "alias", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "responses": {
          "200": {
            "description": "Table list with metadata",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/TableListResponse" }
              }
            }
          }
        }
      }
    },
    "/api/v1/db/{alias}/find-columns": {
      "get": {
        "tags": ["Schema"],
        "summary": "Search across all tables for column name",
        "parameters": [
          { "name": "alias", "in": "path", "required": true, "schema": { "type": "string" } },
          { "name": "name", "in": "query", "required": true, "schema": { "type": "string" } }
        ],
        "responses": {
          "200": {
            "description": "List of tables matching the column name",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/FindColumnsResponse" }
              }
            }
          }
        }
      }
    },
    "/api/v1/db/{alias}/{table}": {
      "get": {
        "tags": ["Records"],
        "summary": "Query paginated records from table",
        "parameters": [
          { "name": "alias", "in": "path", "required": true, "schema": { "type": "string" } },
          { "name": "table", "in": "path", "required": true, "schema": { "type": "string" } },
          { "name": "limit", "in": "query", "required": false, "schema": { "type": "integer", "default": 50 } },
          { "name": "offset", "in": "query", "required": false, "schema": { "type": "integer", "default": 0 } },
          { "name": "sort", "in": "query", "required": false, "schema": { "type": "string" } },
          { "name": "order", "in": "query", "required": false, "schema": { "type": "string", "enum": ["asc", "desc"], "default": "asc" } }
        ],
        "responses": {
          "200": {
            "description": "Paginated records envelope",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/RecordListResponse" }
              }
            }
          }
        }
      },
      "post": {
        "tags": ["Records"],
        "summary": "Insert record into table",
        "parameters": [
          { "name": "alias", "in": "path", "required": true, "schema": { "type": "string" } },
          { "name": "table", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "type": "object", "additionalProperties": true }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Inserted record confirmation",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/RecordMutationResponse" }
              }
            }
          }
        }
      }
    },
    "/api/v1/db/{alias}/{table}/{id}": {
      "get": {
        "tags": ["Records"],
        "summary": "Fetch single record by primary key",
        "parameters": [
          { "name": "alias", "in": "path", "required": true, "schema": { "type": "string" } },
          { "name": "table", "in": "path", "required": true, "schema": { "type": "string" } },
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "responses": {
          "200": { "description": "Single record envelope" },
          "404": { "description": "Record not found" }
        }
      },
      "put": {
        "tags": ["Records"],
        "summary": "Update record by primary key",
        "parameters": [
          { "name": "alias", "in": "path", "required": true, "schema": { "type": "string" } },
          { "name": "table", "in": "path", "required": true, "schema": { "type": "string" } },
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "type": "object", "additionalProperties": true }
            }
          }
        },
        "responses": {
          "200": { "description": "Updated record confirmation" },
          "404": { "description": "Record not found" }
        }
      },
      "delete": {
        "tags": ["Records"],
        "summary": "Delete record by primary key",
        "parameters": [
          { "name": "alias", "in": "path", "required": true, "schema": { "type": "string" } },
          { "name": "table", "in": "path", "required": true, "schema": { "type": "string" } },
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "responses": {
          "200": { "description": "Deleted record confirmation" },
          "404": { "description": "Record not found" }
        }
      }
    },
    "/api/v1/db/{alias}/query": {
      "get": {
        "tags": ["SQL Engine"],
        "summary": "Execute read-only SQL query via query parameter",
        "parameters": [
          { "name": "alias", "in": "path", "required": true, "schema": { "type": "string" } },
          { "name": "q", "in": "query", "required": true, "schema": { "type": "string" } }
        ],
        "responses": {
          "200": { "description": "Query result envelope" }
        }
      },
      "post": {
        "tags": ["SQL Engine"],
        "summary": "Execute SQL statement (read or write) with parameters",
        "parameters": [
          { "name": "alias", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["sql"],
                "properties": {
                  "sql": { "type": "string" },
                  "params": { "type": "array", "items": {} }
                }
              }
            }
          }
        },
        "responses": {
          "200": { "description": "Query execution result envelope" }
        }
      }
    },
    "/api/v1/db/{alias}/views": {
      "get": {
        "tags": ["Views"],
        "summary": "List all dynamic views",
        "parameters": [
          { "name": "alias", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "responses": { "200": { "description": "List of views" } }
      },
      "post": {
        "tags": ["Views"],
        "summary": "Create new dynamic view",
        "parameters": [
          { "name": "alias", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["name", "query"],
                "properties": {
                  "name": { "type": "string" },
                  "query": { "type": "string" }
                }
              }
            }
          }
        },
        "responses": { "201": { "description": "View created confirmation" } }
      },
      "delete": {
        "tags": ["Views"],
        "summary": "Drop dynamic view",
        "parameters": [
          { "name": "alias", "in": "path", "required": true, "schema": { "type": "string" } },
          { "name": "name", "in": "query", "required": true, "schema": { "type": "string" } }
        ],
        "responses": { "200": { "description": "View dropped confirmation" } }
      }
    },
    "/api/v1/db/{alias}/alter": {
      "post": {
        "tags": ["DDL & Migrations"],
        "summary": "Alter table structure (add column, rename column, rename table)",
        "parameters": [
          { "name": "alias", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["table", "operation"],
                "properties": {
                  "table": { "type": "string" },
                  "operation": { "type": "string", "enum": ["add_column", "rename_column", "rename_table"] },
                  "columnName": { "type": "string" },
                  "columnType": { "type": "string" },
                  "defaultValue": { "type": "string" },
                  "isNullable": { "type": "boolean" },
                  "newName": { "type": "string" }
                }
              }
            }
          }
        },
        "responses": { "200": { "description": "DDL operation confirmation" } }
      }
    },
    "/api/v1/db/{alias}/migrate": {
      "post": {
        "tags": ["DDL & Migrations"],
        "summary": "Apply migration script and record in split-db cache",
        "parameters": [
          { "name": "alias", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["version", "script"],
                "properties": {
                  "version": { "type": "string" },
                  "description": { "type": "string" },
                  "script": { "type": "string" }
                }
              }
            }
          }
        },
        "responses": { "200": { "description": "Migration execution confirmation" } }
      }
    },
    "/api/v1/db/{alias}/clone": {
      "post": {
        "tags": ["DDL & Migrations"],
        "summary": "Clone database to new file path and mount under new alias",
        "parameters": [
          { "name": "alias", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["targetPath", "newAlias"],
                "properties": {
                  "targetPath": { "type": "string" },
                  "newAlias": { "type": "string" }
                }
              }
            }
          }
        },
        "responses": { "201": { "description": "Clone completed confirmation" } }
      }
    }
  },
  "components": {
    "schemas": {
      "UniversalEnvelope": {
        "type": "object",
        "required": ["Status", "Attributes", "Results"],
        "properties": {
          "Status": {
            "type": "object",
            "required": ["Code", "Message", "Timestamp"],
            "properties": {
              "Code": { "type": "integer", "example": 200 },
              "Message": { "type": "string", "example": "OK" },
              "Timestamp": { "type": "string", "format": "date-time" }
            }
          },
          "Attributes": {
            "type": "object",
            "additionalProperties": true
          },
          "Results": {
            "type": "array",
            "items": {}
          }
        }
      }
    }
  }
}
```

---

## 3. Embedded Interactive Swagger UI

The REST server hosts an interactive Swagger UI at `/swagger` and `/swagger/`.

### 3.1 Delivery & Zero CDN Architecture
- All assets (Swagger UI HTML shell, bundled CSS, bundled JavaScript) are pre-compiled and embedded directly into the Go binary in `cli/cmddb/rest_ui_assets.go`.
- No external HTTP requests to `unpkg.com`, `cdnjs.cloudflare.com`, or any remote CDN are permitted.
- The Swagger UI bundle is configured with:
  ```javascript
  window.onload = function() {
    const ui = SwaggerUIBundle({
      url: "/swagger/doc.json",
      dom_id: '#swagger-ui',
      deepLinking: true,
      presets: [
        SwaggerUIBundle.presets.apis,
        SwaggerUIStandalonePreset
      ],
      plugins: [
        SwaggerUIBundle.plugins.DownloadUrl
      ],
      layout: "BaseLayout"
    });
    window.ui = ui;
  };
  ```

### 3.2 HTTP Route Handling
- `GET /swagger`: Serves the embedded HTML page with `Content-Type: text/html; charset=utf-8`.
- `GET /swagger/`: Alias route serving identical content.
- `GET /swagger/doc.json`: Serves dynamic OpenAPI 3.0 specification with `Content-Type: application/json`.
- Missing assets or sub-paths fall back gracefully to the main Swagger document.

---

## 4. Custom Advanced Database Visualizer UI

In addition to Swagger, the REST server serves a dedicated, production-grade Web Visualizer at `/ui` and `/ui/`.

### 4.1 Visual Design Tokens & Palette

The UI adopts the GitMap dark-mode design system with high contrast and dense developer information:

| Token Name | Hex Value | Purpose |
|---|---|---|
| `bg-primary` | `#0b0f19` | Main application background |
| `bg-secondary` | `#111827` | Sidebar and card background |
| `bg-tertiary` | `#1f2937` | Table header, inputs, hover states |
| `border-color` | `#374151` | Subtle divider and modal borders |
| `text-primary` | `#f9fafb` | Primary text and headings |
| `text-secondary` | `#9ca3af` | Secondary labels and helper text |
| `accent-primary` | `#3b82f6` | Active tabs, primary buttons, links |
| `accent-success` | `#10b981` | Success badges, row inserted pills |
| `accent-warning` | `#f59e0b` | DDL modification warning indicators |
| `accent-danger` | `#ef4444` | Row deletion buttons, drop table triggers |

### 4.2 Workspace Layout

```
┌────────────────────────────────────────────────────────────────────────┐
│ GitMap DB Visualizer  [Database: main ▼]  Port: 8080 (Active)  [Swagger]│
├───────────────┬────────────────────────────────────────────────────────┤
│ TABLES (4)    │ [Data Grid]  [Schema Editor]  [SQL Query]  [Views] ... │
│ > users   (12)├────────────────────────────────────────────────────────┤
│ > projects(48)│ Filter: [column] [contains ▼] [value]  [+ Insert Row]  │
│ > commits(120)│ ┌──┬─────────┬──────────────┬───────────────┬────────┐ │
│ > tags     (5)│ │id│name     │email         │created_at     │Actions │ │
│               │ ├──┼─────────┼──────────────┼───────────────┼────────┤ │
│ VIEWS (2)     │ │ 1│Alice    │alice@git.io  │2026-01-01...  │[Edit]  │ │
│ > active_proj │ │ 2│Bob      │bob@git.io    │2026-01-02...  │[Delete]│ │
│               │ └──┴─────────┴──────────────┴───────────────┴────────┘ │
│ [+ New Table] │ Page: [ < Prev ] 1 of 5 [ Next > ]   Limit: [ 50 ▼ ]   │
└───────────────┴────────────────────────────────────────────────────────┘
```

### 4.3 Interactive Capabilities & Tabs

#### Tab 1: Live Interactive Data Grid
- **Dynamic Table Fetching**: Selecting a table in the sidebar triggers `GET /api/v1/db/{alias}/{table}?limit={limit}&offset={offset}`.
- **Column Sorting**: Clicking column headers toggles ascending (`asc`) and descending (`desc`) sorting.
- **Filtering**: A filter bar supports `=`, `!=`, `LIKE`, `>`, and `<` operators.
- **Live Inline Row Editing**: Double-clicking a cell switches the cell to an active input field. Pressing `Enter` or clicking the checkmark sends a `PUT /api/v1/db/{alias}/{table}/{id}` request.
- **Live Inline Row Insertion**: An "+ Insert Row" modal displays input fields matching the table schema with type validation. Submitting sends a `POST /api/v1/db/{alias}/{table}` request.
- **Row Deletion**: Clicking the trash icon opens a confirmation modal before sending a `DELETE /api/v1/db/{alias}/{table}/{id}` request.

#### Tab 2: Schema & Column Editor
- **Table Definition View**: Displays the raw SQLite `CREATE TABLE` statement with syntax highlighting and copy-to-clipboard.
- **Column List**: Displays column properties including Column Name, Declared Type, Primary Key indicator, Not Null constraint, and Default Value.
- **Column Alterations**:
  - "+ Add Column": Modal to specify Column Name, Type dropdown (`INTEGER`, `TEXT`, `REAL`, `BLOB`, `BOOLEAN`), Nullable checkbox, and Default Value. Dispatches `POST /api/v1/db/{alias}/alter`.
  - "Rename Column": Dispatches `POST /api/v1/db/{alias}/alter`.
- **Table Creation Wizard**: Visual builder to define table name and multiple columns before executing `POST /api/v1/db/{alias}/tables/create`.

#### Tab 3: Visual Query & SQL Console
- **Interactive SQL Editor**: Monospace textarea with query templates (`SELECT * FROM ...`, `JOIN`, `AGGREGATE`).
- **Keyboard Shortcut**: `Ctrl+Enter` or `Cmd+Enter` executes the current SQL query immediately.
- **Dual Execution Modes**:
  - `GET /api/v1/db/{alias}/query?q=...` for read-only queries.
  - `POST /api/v1/db/{alias}/query` for mutations with parameter arrays.
- **Execution Telemetry**: Displays query duration in milliseconds, affected row count, or returned row count.
- **Export Facility**: Buttons to export the query results table directly to CSV or formatted JSON.

#### Tab 4: Dynamic View Builder
- **Visual View Creation**: Modal requiring View Name and a `SELECT` SQL statement.
- **Creation Execution**: Dispatches `POST /api/v1/db/{alias}/views`.
- **View Inspection & Deletion**: Lists all existing views with instant query preview and a "Drop View" button calling `DELETE /api/v1/db/{alias}/views?name={name}`.

#### Tab 5: Database Migration & Clone Manager
- **Database Cloning**: Modal allowing cloning of the current database file to a specified target filesystem path and mounting under a new alias. Dispatches `POST /api/v1/db/{alias}/clone`.
- **Migration Runner**: Interface to execute multi-statement DDL/DML migration scripts. Dispatches `POST /api/v1/db/{alias}/migrate`.
- **Split-DB Cache History**: Displays historical migrations and alteration events recorded in the Split-DB cache database.

---

## 5. Comprehensive Acceptance Verification & Testing Strategy

The test suite ensures 100% verification across all requirements without regressions.

### 5.1 Test Suite Organization

| Test Suite | Implementation File | Verification Focus |
|---|---|---|
| **Port Probing & Fallback** | `cli/cmddb/cmd_rest_test.go` | Sequential port availability checks, occupied port skipping, error handling |
| **Multi-Alias Mounting** | `cli/cmddb/cmd_rest_test.go` | Single alias, comma-separated aliases, mixed alias/path resolution, catalog |
| **Schema & Discovery** | `cli/cmddb/cmd_rest_test.go` | `PRAGMA table_info` extraction, column search across tables (`find-columns`) |
| **Record CRUD Engine** | `cli/cmddb/cmd_rest_test.go` | Pagination, sorting, insertion, primary key update, deletion, 404 responses |
| **DDL & Migrations** | `cli/cmddb/cmd_rest_test.go` | `alter` table add column, create table, migration script caching in split-db |
| **Dynamic Views & SQL** | `cli/cmddb/cmd_rest_test.go` | View creation/query/drop, direct SQL GET/POST, SQL error envelopes |
| **Deadlock & Concurrency** | `cli/cmddb/cmd_rest_test.go` | `SetMaxOpenConns(1)` enforcement, 50 parallel goroutines read/write stress |
| **Swagger & UI Delivery** | `cli/cmddb/cmd_rest_test.go` | `/swagger/doc.json` OpenAPI 3.0 validity, `/swagger` and `/ui` HTTP 200 delivery |
| **CLI Commands Parity** | `cli/cmddb/cmd_record_del_test.go`, `cli/cmddb/cmd_rm_test.go` | `gitmap db-record-del` and `gitmap db-rm` command execution |

### 5.2 Detailed Verification Scenarios

#### Scenario 1: Sequential Port Fallback Engine
1. Occupy port `8080` with a temporary listener.
2. Invoke port picker with candidate ports `8080,8081,8082`.
3. Assert that the engine automatically binds to `8081` without raising an error.
4. Occupy all candidate ports `8080,8081`.
5. Assert that the engine returns a descriptive error indicating that none of the specified ports are available.

#### Scenario 2: Multi-Alias & Path Resolution
1. Create two test SQLite databases: `db1.sqlite` (with table `users`) and `db2.sqlite` (with table `projects`).
2. Mount using command syntax: `main,analytics` pointing to `db1.sqlite` and `db2.sqlite`.
3. Query `GET /api/v1/databases`.
4. Verify response contains both `main` and `analytics` aliases with their corresponding file paths.
5. Mount mixed syntax: `main,./test_extra.db`. Verify `./test_extra.db` is automatically mounted with a derived alias or normalized path key.

#### Scenario 3: End-to-End CRUD Lifecycle
1. Execute `POST /api/v1/db/main/users` with payload `{"name": "Alice", "email": "alice@gitmap.dev", "score": 95.5, "is_active": true}`.
2. Verify HTTP 201 Created and response containing generated Primary Key `id=1`.
3. Execute `GET /api/v1/db/main/users/1`. Verify all fields match types and values.
4. Execute `PUT /api/v1/db/main/users/1` with payload `{"name": "Alice Cooper", "score": 99.0}`.
5. Verify HTTP 200 OK and updated fields.
6. Execute `DELETE /api/v1/db/main/users/1`. Verify HTTP 200 OK.
7. Execute subsequent `GET /api/v1/db/main/users/1`. Verify HTTP 404 Not Found within standard envelope.

#### Scenario 4: DDL Alteration & Column Discovery
1. Create table `inventory` with columns `id INTEGER PRIMARY KEY, sku TEXT`.
2. Execute `GET /api/v1/db/main/find-columns?name=sku`. Verify `inventory` is returned.
3. Execute `POST /api/v1/db/main/alter` with payload:
   ```json
   { "table": "inventory", "operation": "add_column", "columnName": "quantity", "columnType": "INTEGER", "defaultValue": "0" }
   ```
4. Verify HTTP 200 OK.
5. Inspect `PRAGMA table_info("inventory")` and confirm column `quantity` exists with type `INTEGER`.
6. Insert record including `quantity` and verify value persistence.

#### Scenario 5: Dynamic Views & SQL Console
1. Execute `POST /api/v1/db/main/views` with payload:
   ```json
   { "name": "active_users_view", "query": "SELECT id, name FROM users WHERE is_active = 1" }
   ```
2. Verify HTTP 201 Created.
3. Query `GET /api/v1/db/main/active_users_view`. Verify view results match expected subset.
4. Execute direct query via `POST /api/v1/db/main/query` with `{"sql": "SELECT COUNT(*) as total FROM users"}`.
5. Verify response contains `total` count attribute.

#### Scenario 6: Deadlock Prevention Under Heavy Concurrency
1. Configure SQLite connection with `conn.SetMaxOpenConns(1)`.
2. Launch 50 concurrent goroutines executing simultaneous `GET` (selects), `POST` (inserts), and `PUT` (updates) operations over a 3-second window.
3. Verify zero instances of `database is locked`, zero deadlocks, and zero goroutine panics.
4. Ensure all connections return cleanly to the pool.

#### Scenario 7: OpenAPI 3.0 Document & UI Delivery
1. Query `GET /swagger/doc.json`.
2. Validate response has `Content-Type: application/json`.
3. Parse JSON and verify OpenAPI 3.0.3 schema compliance (valid `openapi`, `info`, `paths`, and `components` fields).
4. Query `GET /swagger` and `GET /ui`.
5. Verify both respond with HTTP 200 OK, `Content-Type: text/html; charset=utf-8`, and contain embedded CSS/JS without external CDN references.

#### Scenario 8: CLI Command Parity
1. Run `gitmap db-record-del main -table users -id 5`.
2. Verify row with `id=5` is removed from table `users` in database `main`.
3. Run `gitmap db-rm test_alias`.
4. Verify `test_alias` is removed from the registry.

---

## 6. Traceability Matrix & Quality Gates

| Quality Gate | Requirement | Verification Method | Pass Criteria |
|---|---|---|---|
| **G1: Port Probing** | Comma-separated port fallback | Unit Test `TestPortFallback` | Binds to first free port; errors when all busy |
| **G2: Multi-DB Mount** | Multiple aliases and direct paths | Unit Test `TestMultiDbMount` | All DBs mounted and queryable via `/api/v1/databases` |
| **G3: Universal Envelope** | PascalCase `{ Status, Attributes, Results }` | Integration Test `TestEnvelopeFormat` | 100% of REST responses match envelope format |
| **G4: OpenAPI 3.0 & Swagger** | Dynamic spec and embedded UI | Integration Test `TestSwaggerDocAndUI` | `/swagger/doc.json` valid OpenAPI 3.0; `/swagger` returns HTTP 200 |
| **G5: Visualizer UI** | Interactive dark mode dashboard | Integration Test `TestVisualizerUI` | `/ui` returns HTTP 200; zero CDN references |
| **G6: CRUD & DDL Engine** | Insert, read, update, delete, alter | Integration Test `TestCrudAndDdl` | Record mutations and table alterations succeed |
| **G7: Deadlock Prevention** | Single-writer safety under concurrency | Stress Test `TestConcurrentAccess` | 50 concurrent goroutines with 0 lock errors |
| **G8: CLI Parity** | `db-record-del` and `db-rm` commands | CLI Test `TestCliCommands` | Commands execute and modify SQLite state cleanly |
