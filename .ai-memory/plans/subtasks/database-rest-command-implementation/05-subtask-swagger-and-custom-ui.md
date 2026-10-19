# Subtask 05: Dynamic OpenAPI 3.0 Generator & Custom Database Visualizer UI

> **Document Version:** 1.0.0  
> **Status:** PENDING  
> **Parent Plan:** `.ai-memory/plans/database-rest-command-implementation.md`  
> **Associated Spec:** `02-spec/21-app/database-rest-command-implementation/03-swagger-ui-and-verification-spec.md`  
> **Assigned Role:** Worker 03 / UI & Documentation Specialist  
> **Disjoint Files Owned:**  
> - `cli/cmddb/rest_swagger.go`  
> - `cli/cmddb/rest_ui.go`  
> - `cli/cmddb/rest_ui_assets.go`  

---

## 1. Objective

Implement the dynamic OpenAPI 3.0 documentation generator, the embedded interactive Swagger UI, and the comprehensive custom Database Visualizer Web UI for the GitMap `db-rest` server. All assets and handlers must operate completely offline with zero external CDN dependencies.

---

## 2. Scope & Implementation Boundaries

### 2.1 Owned Implementation Files
1. **`cli/cmddb/rest_swagger.go`**:
   - Dynamic OpenAPI 3.0.3 specification builder.
   - Dynamic introspection of mounted SQLite databases via `PRAGMA table_info` and `sqlite_master`.
   - SQLite affinity to OpenAPI 3.0 type conversion engine.
   - Handlers for `GET /swagger/doc.json`, `GET /swagger`, and `GET /swagger/`.
2. **`cli/cmddb/rest_ui.go`**:
   - HTTP route registration and request dispatch for the custom Database Visualizer UI at `/ui` and `/ui/`.
   - HTML template rendering and dynamic initialization state injection (mounted database catalog, active database alias, server port).
3. **`cli/cmddb/rest_ui_assets.go`**:
   - Embedded standalone static assets: Swagger UI JavaScript/CSS bundle, Visualizer UI CSS styles, and reactive JavaScript application script.
   - Zero external CDN dependencies; all assets bundled directly into Go string constants or embedded byte slices.

---

## 3. Step-by-Step Execution Plan

### Step 5.1: Dynamic OpenAPI 3.0 Specification Builder (`cli/cmddb/rest_swagger.go`)
- **Data Models**: Define Go structs for OpenAPI 3.0 schema: `OpenAPIObject`, `InfoObject`, `ServerObject`, `PathsObject`, `PathItemObject`, `OperationObject`, `ParameterObject`, `RequestBodyObject`, `ResponsesObject`, `ResponseObject`, `ComponentsObject`, `SchemaObject`.
- **Introspection Engine**:
  - Implement `generateOpenAPISpec(dbManager *MultiDbManager) (*OpenAPIObject, error)`.
  - For each mounted database, inspect `sqlite_master` to retrieve all tables and views.
  - Query `PRAGMA table_info({table})` to extract column metadata: name, declared type, primary key flag, nullability, default values.
  - Map SQLite types to OpenAPI types (`INTEGER` -> `integer/int64`, `REAL` -> `number/double`, `TEXT` -> `string`, `BOOLEAN` -> `boolean`, `BLOB` -> `string/byte`).
- **Endpoint Synthesis**:
  - System routes: `/api/v1/databases`.
  - Schema routes: `/api/v1/db/{alias}/schema`, `/api/v1/db/{alias}/tables`, `/api/v1/db/{alias}/find-columns`.
  - CRUD routes per table: `/api/v1/db/{alias}/{table}` (GET, POST) and `/api/v1/db/{alias}/{table}/{id}` (GET, PUT, DELETE).
  - SQL engine routes: `/api/v1/db/{alias}/query` (GET, POST).
  - Dynamic view routes: `/api/v1/db/{alias}/views` (GET, POST, DELETE).
  - DDL & Migration routes: `/api/v1/db/{alias}/alter`, `/api/v1/db/{alias}/tables/create`, `/api/v1/db/{alias}/migrate`, `/api/v1/db/{alias}/clone`.
- **JSON Serialization**: Cache generated JSON in memory with invalidation on DDL or migration events.

### Step 5.2: Embedded Swagger UI Delivery (`cli/cmddb/rest_swagger.go` & `cli/cmddb/rest_ui_assets.go`)
- Bundle Swagger UI 5.x standalone distribution (HTML shell, CSS stylesheet, JavaScript bundle) in `cli/cmddb/rest_ui_assets.go`.
- Configure Swagger UI to point to `/swagger/doc.json` dynamically based on current `window.location.origin`.
- Implement `handleSwaggerUI(w http.ResponseWriter, r *http.Request)`:
  - Serve HTML with `Content-Type: text/html; charset=utf-8`.
  - Handle trailing slash redirects (`/swagger` -> `/swagger/` or serve directly).

### Step 5.3: Custom Database Visualizer UI Server (`cli/cmddb/rest_ui.go`)
- Implement `handleVisualizerUI(w http.ResponseWriter, r *http.Request)`:
  - Read mounted database list from `MultiDbManager`.
  - Render self-contained HTML single-page application (SPA).
  - Set `Content-Type: text/html; charset=utf-8` and appropriate cache-control headers (`no-cache`).

### Step 5.4: Custom Database Visualizer Frontend Application (`cli/cmddb/rest_ui_assets.go`)
- Implement the complete front-end developer console with dark mode styling:
  1. **Header Navigation**:
     - Database selector dropdown populated from `/api/v1/databases`.
     - Connection status indicator pill and port display.
     - Direct link button to Swagger UI (`/swagger`).
  2. **Sidebar Schema Browser**:
     - Live table list with row-count badges.
     - Dynamic views list with view badges.
     - Fast filter search box for tables and column names (`/api/v1/db/{alias}/find-columns`).
     - "+ New Table" and "+ New View" modal triggers.
  3. **Tab 1: Live Interactive Data Grid**:
     - Paginated table viewer with configurable page sizes (25, 50, 100) and page navigation controls.
     - Column header sort toggle (ASC/DESC).
     - Filter builder with comparison operators (`=`, `!=`, `LIKE`, `>`, `<`).
     - Inline cell editing: double-click cell, validate type, submit `PUT /api/v1/db/{alias}/{table}/{id}`.
     - Inline row insertion modal: schema-driven form fields, submits `POST /api/v1/db/{alias}/{table}`.
     - Row deletion with confirmation dialog: submits `DELETE /api/v1/db/{alias}/{table}/{id}`.
  4. **Tab 2: Schema & Column Editor**:
     - Visual column list with primary key badges, nullability indicators, and default values.
     - Raw DDL viewer with copy-to-clipboard button.
     - "+ Add Column" modal: dispatches `POST /api/v1/db/{alias}/alter`.
     - Visual Table Creation Wizard: dispatches `POST /api/v1/db/{alias}/tables/create`.
  5. **Tab 3: Visual Query & SQL Console**:
     - Monospace SQL text editor with boilerplate templates (`SELECT`, `JOIN`, `UPDATE`).
     - Keyboard shortcut `Ctrl+Enter` to execute.
     - Direct execution calling `POST /api/v1/db/{alias}/query`.
     - Execution metrics: execution time in milliseconds, affected row count.
     - Export query results to CSV and JSON.
  6. **Tab 4: Dynamic View Builder**:
     - Visual form to create views (`name`, `query`), dispatches `POST /api/v1/db/{alias}/views`.
     - View list with quick preview query and "Drop View" action.
  7. **Tab 5: Migration & Clone Manager**:
     - Database clone dialog calling `POST /api/v1/db/{alias}/clone`.
     - Migration runner console calling `POST /api/v1/db/{alias}/migrate`.
     - Split-DB migration history viewer.

---

## 4. Coding Standards & Architectural Invariants

- **File Sizing & Modularity**: Keep files focused; decompose long assets into modular constants in `cli/cmddb/rest_ui_assets.go`.
- **Function Brevity**: Target 8–15 lines per handler and helper function.
- **Max Nesting**: Maintain nesting depth $\le 1$ using early returns and guard clauses.
- **Affirmative Booleans**: Use positive naming conventions (`isOpenApiEnabled`, `hasCustomColumns`, `isTableValid`, `isPrimaryKey`).
- **No External Network Calls**: Zero CDN references (`<script src="https://...">` strictly forbidden).
- **Universal Response Envelope**: All API interactions from the frontend UI adhere to `{ Status, Attributes, Results }`.

---

## 5. Verification & Acceptance Criteria

- **AC-05.1**: `GET /swagger/doc.json` returns HTTP 200 with valid OpenAPI 3.0.3 JSON schema for all mounted databases.
- **AC-05.2**: `GET /swagger` serves embedded interactive Swagger UI with working "Try It Out" capability against local endpoints.
- **AC-05.3**: `GET /ui` serves the custom Database Visualizer dashboard with dark theme styling and zero external CDN script tags.
- **AC-05.4**: Data Grid in Visualizer supports table selection, pagination, sorting, live inline row editing, row insertion, and row deletion.
- **AC-05.5**: SQL Console supports direct query execution with execution duration telemetry and CSV/JSON export.
- **AC-05.6**: Column Editor successfully dispatches `POST /api/v1/db/{alias}/alter` to add columns on live tables.
