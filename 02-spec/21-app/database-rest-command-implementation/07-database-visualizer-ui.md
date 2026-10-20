# 07-database-visualizer-ui: Embedded Swagger UI & Custom Database Visualizer

- **Spec ID:** `02-spec/21-app/database-rest-command-implementation/07-database-visualizer-ui.md`
- **Status:** `APPROVED`
- **Version:** `1.0.0`
- **Subsystem:** Embedded Swagger UI, Database Visualizer Web UI
- **Parent Plan:** `.ai-memory/plans/database-rest-command-implementation.md`
- **Dependencies:** `cli/cmddb`, `cli/store`, `cli/appfault`

---

## 1. Embedded Interactive Swagger UI

The REST server hosts an interactive Swagger UI at `/swagger` and `/swagger/`.

### 1.1 Delivery & Zero CDN Architecture
- All assets (Swagger UI HTML shell, bundled CSS, bundled JavaScript) are pre-compiled and embedded directly into the Go binary in `cli/cmddb/rest_swagger_asset.go`.
- No external HTTP requests to `unpkg.com`, `cdnjs.cloudflare.com`, or any remote CDN are permitted.

### 1.2 HTTP Route Handling
- `GET /swagger`: Serves embedded HTML page with `Content-Type: text/html; charset=utf-8`.
- `GET /swagger/`: Alias route serving identical content.
- `GET /swagger/doc.json`: Serves dynamic OpenAPI 3.0 specification with `Content-Type: application/json`.

---

## 2. Custom Advanced Database Visualizer UI

In addition to Swagger, the REST server serves a dedicated, production-grade Web Visualizer at `/ui` and `/ui/`.

### 2.1 Visual Design Tokens & Palette

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

### 2.2 Workspace Layout

```text
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

### 2.3 Interactive Capabilities & Tabs

#### Tab 1: Live Interactive Data Grid
- **Dynamic Table Fetching**: Selecting a table in the sidebar triggers `GET /api/v1/db/{alias}/{table}?limit={limit}&offset={offset}`.
- **Column Sorting**: Clicking column headers toggles ascending (`asc`) and descending (`desc`) sorting.
- **Filtering**: A filter bar supports `=`, `!=`, `LIKE`, `>`, and `<` operators.
- **Live Inline Row Editing**: Pressing Edit opens a modal sending `PUT /api/v1/db/{alias}/{table}/{id}`.
- **Live Inline Row Insertion**: An "+ Insert Row" modal sends `POST /api/v1/db/{alias}/{table}`.
- **Row Deletion**: Clicking delete sends `DELETE /api/v1/db/{alias}/{table}/{id}`.

#### Tab 2: Schema & Column Editor
- **Table Definition View**: Displays raw SQLite `CREATE TABLE` statement with syntax styling.
- **Column List**: Displays column properties including Column Name, Declared Type, PK, Not Null, and Default.
- **Column Alterations**:
  - "+ Add Column": Modal to specify Column Name, Type, and Default Value, calling `POST /api/v1/db/{alias}/alter`.
- **Table Creation Wizard**: Visual builder calling `POST /api/v1/db/{alias}/tables/create`.

#### Tab 3: Visual Query & SQL Console
- **Interactive SQL Editor**: Monospace textarea with query execution via `Ctrl+Enter`.
- **Dual Execution Modes**: `GET /api/v1/db/{alias}/query?sql=...` for read-only; `POST /api/v1/db/{alias}/query` for mutations.
- **Execution Telemetry**: Displays query duration in milliseconds and returned/affected rows.
- **Export Facility**: Buttons to export query results directly to CSV or JSON.

#### Tab 4: Dynamic View Builder
- **Visual View Creation**: Modal requiring View Name and `SELECT` query, calling `POST /api/v1/db/{alias}/views`.
- **View Inspection & Deletion**: Lists views with query preview and "Drop View" calling `DELETE /api/v1/db/{alias}/views/{name}`.

#### Tab 5: Database Migration & Clone Manager
- **Database Cloning**: Modal calling `POST /api/v1/db/{alias}/clone`.
- **Migration Runner**: Interface calling `POST /api/v1/db/{alias}/migrate`.
- **Split-DB Cache History**: Displays historical migrations and alteration events.
