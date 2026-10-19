# Specification: SQLite Database View & Join Engine

- **Spec Path:** `02-spec/21-app/command-and-database-integration-instructions/04-view-and-join-engine.md`
- **Module:** `cli/cmddb`, `cli/store`, `cli/cmd`
- **Status:** APPROVED
- **Version:** 1.0.0
- **Line Count Cap:** <= 300 lines

---

## 1. User Request (Verbatim)

```text
Author the modular specifications in 02-spec/21-app/command-and-database-integration-instructions/:
1. 02-spec/21-app/command-and-database-integration-instructions/03-alias-registry-and-storage.md: SQLite Database Registry architecture. Persistent alias-to-path mapping in gitmap.db (PascalCase table DatabaseAlias or SplitDatabaseRegistry integration), file size detection, search history count, and gitmap db-listing (or dbl) output in both formatted terminal table and --json envelope.
2. 02-spec/21-app/command-and-database-integration-instructions/04-view-and-join-engine.md: Specification for gitmap db-find (wildcard text, text*, *text, regex --grep/-r across columns) and gitmap db-view (-t "t1 as a, t2 as b", --viewname, --find filter condition with automatic join inference, alias support, column projections).
Use strictly relative Git paths and lowercase filenames. All files <= 300 lines.
```

---

## 2. Architectural Motivation & Overview

Navigating SQLite databases across multi-agent environments requires rapid ad-hoc introspection without writing boilerplate SQL scripts. GitMap introduces two complementary database inspection engines:

1. **`gitmap db-find`:** High-speed, multi-table string and regex search engine that scans across text columns to locate records matching keywords, prefixes, suffixes, or regular expressions.
2. **`gitmap db-view`:** Declarative multi-table projection engine that parses table aliases, automatically infers join conditions via foreign keys and `{Entity}Id` heuristics, applies filters, and renders formatted tabular or JSON results.

Both engines leverage the **SQLite Database Registry** (`02-spec/21-app/command-and-database-integration-instructions/03-alias-registry-and-storage.md`) to accept semantic aliases (`--db master`, `--db heatmap`, `--db pipeline`) or direct filepaths.

---

## 3. Specification: `gitmap db-find`

### 3.1 Command Syntax
```bash
gitmap db-find <pattern> [flags]
# Short alias:
gitmap db find <pattern> [flags]
```

### 3.2 Pattern Matching Rules
The search pattern `<pattern>` is evaluated according to the following syntax rules:

| Pattern Syntax | Matching Semantics | Translated SQL / Engine Behavior |
| :--- | :--- | :--- |
| `text` | Substring match | `LIKE '%text%'` (case-insensitive) |
| `text*` | Prefix match | `LIKE 'text%'` |
| `*text` | Suffix match | `LIKE '%text'` |
| `*text*` | Explicit substring match | `LIKE '%text%'` |
| `--grep <regex>` / `-r <regex>` | Regular expression | Evaluated in Go memory or SQLite `REGEXP` UDF |

### 3.3 Flags & Options
- `--db <alias|path>`: Target database (default: `master`).
- `-t, --table <name>`: Limit search to a specific table. If omitted, introspects all non-internal tables via `sqlite_master`.
- `-c, --columns <col1,col2>`: Restrict search to specified columns. If omitted, evaluates all `TEXT`, `VARCHAR`, and `CHAR` columns.
- `-i, --ignore-case`: Force case-insensitive evaluation (default: `true`).
- `-m, --max-matches <n>`: Maximum matched rows per table (default: `20`).
- `--json`: Output typed JSON envelope matching `GitMapEnvelope`.

### 3.4 Processing Pipeline
1. **Schema Discovery:** Resolves target DB connection and executes `PRAGMA table_info(table)` for candidate tables.
2. **Query Generation:** Constructs parameterized queries combining candidate text columns:
   ```sql
   SELECT TableId, Col1, Col2 FROM Table WHERE Col1 LIKE ? OR Col2 LIKE ? LIMIT ?
   ```
3. **Regex In-Memory Filtering:** When `--grep` is provided, fetches candidate string batches and executes compiled Go `regexp.Regexp` to capture exact match offsets.
4. **Search Telemetry:** Increments `SearchHistoryCount` in `DatabaseAlias` upon execution.

### 3.5 Output Presentation

#### Formatted Terminal Table:
```text
[master:Project] 2 matches found
┌───────────┬──────────────────────┬──────────────────────┬─────────────────────────────────────┐
│ ProjectId │ Match Column         │ Value                │ Highlighted Snippet                 │
├───────────┼──────────────────────┼──────────────────────┼─────────────────────────────────────┤
│ 14        │ ProjectName          │ gitmap-v28           │ ...repo: [gitmap-v28] (active)...   │
│ 82        │ Description          │ Core gitmap daemon   │ ...cli daemon for [gitmap] suite... │
└───────────┴──────────────────────┴──────────────────────┴─────────────────────────────────────┘
```

#### JSON Response Envelope:
```json
{
  "apiVersion": "v2",
  "status": "success",
  "data": {
    "database": "master",
    "pattern": "gitmap*",
    "isRegex": false,
    "totalMatches": 2,
    "matches": [
      {
        "table": "Project",
        "primaryKey": "ProjectId",
        "primaryKeyValue": 14,
        "column": "ProjectName",
        "matchValue": "gitmap-v28",
        "matchOffsets": [[0, 6]]
      }
    ]
  },
  "timestamp": "2026-10-19T16:20:00Z"
}
```

---

## 4. Specification: `gitmap db-view`

### 4.1 Command Syntax
```bash
gitmap db-view -t "<tables>" [flags]
# Alternative invocation:
gitmap db view -t "<tables>" [flags]
```

### 4.2 Table Alias & Relation Syntax (`-t, --tables`)
Accepts comma-separated tables with optional alias definitions:
- Single table: `-t "Project as p"` or `-t "Project"`
- Multiple tables: `-t "Project as p, Task as t"` or `-t "t1 as a, t2 as b"`

### 4.3 Automatic Join Inference Engine
When multiple tables are passed without explicit `ON` conditions, the engine executes the following deterministic resolution ladder:

1. **Foreign Key Introspection:**
   - Queries `PRAGMA foreign_key_list(Table)`.
   - If Table `B` has foreign key `ProjectId` referencing Table `A(ProjectId)`, joins with:
     ```sql
     FROM Project AS p INNER JOIN Task AS t ON p.ProjectId = t.ProjectId
     ```
2. **Standard Entity PK/FK Matching:**
   - If no explicit SQLite foreign key constraint is declared, inspects column names across both tables.
   - If Table `A` has primary key `{Entity}Id` (e.g. `ProjectId`) and Table `B` possesses an identical column `{Entity}Id`, generates join condition `a.ProjectId = b.ProjectId`.
3. **Ambiguity & Disjoint Guards:**
   - If multiple candidates exist or no shared identifier is discovered, halts query generation and returns error `E5010:JOIN_INFERENCE_FAILED` with candidate columns for manual disambiguation.

### 4.4 Column Projections (`-c, --columns`)
- Syntax: `-c "p.ProjectId, p.ProjectName, t.TaskSlug, t.Status"`
- Defaults: If `-c` is omitted, projects primary keys and leading descriptive text columns from each joined table.
- Expression Support: Supports basic column aliasing (e.g. `p.ProjectName AS Name`).

### 4.5 Filter Condition (`--find`)
- Syntax: `--find "<expression>"`
- Dual-Mode Evaluation:
  - **Structured SQL Filter:** Evaluates SQL WHERE conditions (e.g. `--find "p.IsActive = 1 AND t.Status = 'PENDING'"`).
  - **Fuzzy Token Filter:** If input contains raw text or wildcards without relational operators (`=`, `<`, `>`), expands across projected text columns: `(p.ProjectName LIKE '%query%' OR t.TaskSlug LIKE '%query%')`.

### 4.6 Named Views (`--viewname`)
- When `--viewname <Name>` is specified:
  - Generates deterministic view DDL:
    ```sql
    CREATE VIEW IF NOT EXISTS <Name> AS SELECT ... FROM ...
    ```
  - Registers the view in SQLite catalog or updates existing virtual view schema.

### 4.7 Output Presentation

#### Formatted Terminal Table:
```text
[View: Project ⋈ Task] Filter: Status = 'PENDING' | Limit: 5
┌───────────┬──────────────────────┬──────────────────────┬──────────┐
│ ProjectId │ ProjectName          │ TaskSlug             │ Status   │
├───────────┼──────────────────────┼──────────────────────┼──────────┤
│ 14        │ gitmap               │ db-registry-spec     │ PENDING  │
│ 14        │ gitmap               │ view-engine-spec     │ PENDING  │
│ 22        │ scripts-fixer        │ update-ci-runner     │ PENDING  │
└───────────┴──────────────────────┴──────────────────────┴──────────┘
```

#### JSON Response Envelope:
```json
{
  "apiVersion": "v2",
  "status": "success",
  "data": {
    "viewName": "ActiveTasksView",
    "tables": ["Project as p", "Task as t"],
    "joinCondition": "p.ProjectId = t.ProjectId",
    "columns": ["p.ProjectId", "p.ProjectName", "t.TaskSlug", "t.Status"],
    "rowCount": 3,
    "rows": [
      {
        "ProjectId": 14,
        "ProjectName": "gitmap",
        "TaskSlug": "db-registry-spec",
        "Status": "PENDING"
      }
    ]
  },
  "timestamp": "2026-10-19T16:21:00Z"
}
```

---

## 5. Error Codes & Diagnostics

| Code | Label | Cause & Remediation |
| :--- | :--- | :--- |
| `E1000` | `INVALID_SYNTAX` | Malformed table alias or invalid regex string. |
| `E5010` | `JOIN_INFERENCE_FAILED` | Could not infer relation between tables. Use explicit join condition. |
| `E5011` | `AMBIGUOUS_COLUMN` | Column exists in multiple tables without alias qualifier (e.g. `CreatedAt`). |
| `E5012` | `TABLE_NOT_FOUND` | Specified table does not exist in target database. |
| `E5013` | `VIEW_CREATION_FAILED` | Insufficient permissions or name conflict when creating named view. |

---

## 6. Acceptance Criteria & Verification

1. **Wildcard & Regex Search:** `gitmap db-find "gitmap*"` and `gitmap db-find -r "^v[0-9]+"` return accurate highlighted matches without SQL syntax errors.
2. **Automatic Join Inference:** `gitmap db-view -t "Project as p, Task as t"` automatically connects `p.ProjectId = t.ProjectId` without requiring explicit `--on`.
3. **Projection & Aliasing:** Column projections (`-c "p.ProjectName, t.TaskSlug"`) successfully restrict output columns.
4. **JSON Compliance:** Both commands honor `--json` outputting valid envelopes matching `GitMapEnvelope`.
5. **Database Registry Integration:** Commands accept aliases registered in `DatabaseAlias` via `--db <alias>`.
