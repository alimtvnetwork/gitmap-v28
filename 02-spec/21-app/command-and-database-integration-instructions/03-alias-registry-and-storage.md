# Specification: SQLite Database Registry Architecture & Alias Storage

- **Spec Path:** `02-spec/21-app/command-and-database-integration-instructions/03-alias-registry-and-storage.md`
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

## 2. Architectural Overview

GitMap utilizes a three-tier Split-DB SQLite storage engine (`gitmap.db`, `installation.db`, `repodb/pipeline.db`), alongside task and heatmap companion stores. To enable seamless AI agent navigation and CLI operations without hardcoding absolute paths, GitMap establishes the **SQLite Database Registry**.

### Core Responsibilities
1. **Persistent Alias Mapping:** Stores semantic nicknames (`master`, `tasks`, `heatmap`, `analytics`, `custom-db`) pointing to physical database files.
2. **Metadata & Health Telemetry:** Tracks physical file existence, live byte size, and cumulative search/query invocations.
3. **Transparent Path Resolution:** Intercepts CLI database flags in `cli/cmddb` to dynamically resolve aliases via `DatabaseAlias` before falling back to filesystem paths.
4. **Dual Output Support:** Renders rich terminal tables via Lipgloss styling and structured JSON payloads for AI consumption.

---

## 3. Database Schema: `DatabaseAlias`

The registry is stored in the master database (`gitmap.db`) using strict singular PascalCase conventions, integer primary keys, and positive boolean prefixes:

```sql
CREATE TABLE IF NOT EXISTS DatabaseAlias (
    DatabaseAliasId     INTEGER PRIMARY KEY AUTOINCREMENT,
    AliasName           TEXT NOT NULL COLLATE NOCASE,
    DatabasePath        TEXT NOT NULL,
    DatabaseType        TEXT NOT NULL DEFAULT 'CUSTOM',
    IsActive            INTEGER NOT NULL DEFAULT 1,
    FileSizeBytes       INTEGER NOT NULL DEFAULT 0,
    SearchHistoryCount  INTEGER NOT NULL DEFAULT 0,
    LastAccessedAt      DATETIME NULL,
    Description         TEXT NULL,
    Notes               TEXT NULL,
    CreatedAt           DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt           DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT UQ_DatabaseAlias_AliasName UNIQUE (AliasName)
);

CREATE INDEX IF NOT EXISTS IX_DatabaseAlias_IsActive ON DatabaseAlias (IsActive);
CREATE INDEX IF NOT EXISTS IX_DatabaseAlias_Type ON DatabaseAlias (DatabaseType);
```

### Go Domain Model (`cli/cmddb/alias_types.go`)

```go
package cmddb

import "time"

type DatabaseAlias struct {
    DatabaseAliasId    int64      `gorm:"column:DatabaseAliasId;primaryKey;autoIncrement" json:"databaseAliasId"`
    AliasName          string     `gorm:"column:AliasName;not null;unique" json:"aliasName"`
    DatabasePath       string     `gorm:"column:DatabasePath;not null" json:"databasePath"`
    DatabaseType       string     `gorm:"column:DatabaseType;not null" json:"databaseType"`
    IsActive           bool       `gorm:"column:IsActive;not null;default:true" json:"isActive"`
    FileSizeBytes      int64      `gorm:"column:FileSizeBytes;not null;default:0" json:"fileSizeBytes"`
    SearchHistoryCount int64      `gorm:"column:SearchHistoryCount;not null;default:0" json:"searchHistoryCount"`
    LastAccessedAt     *time.Time `gorm:"column:LastAccessedAt" json:"lastAccessedAt,omitempty"`
    Description        *string    `gorm:"column:Description" json:"description,omitempty"`
    Notes              *string    `gorm:"column:Notes" json:"notes,omitempty"`
    CreatedAt          time.Time  `gorm:"column:CreatedAt" json:"createdAt"`
    UpdatedAt          time.Time  `gorm:"column:UpdatedAt" json:"updatedAt"`
}

func (DatabaseAlias) TableName() string {
    return "DatabaseAlias"
}
```

---

## 4. Path Resolution & Lifecycle Mechanics

### 4.1 Resolution Hierarchy (`ResolveDbPathWithRegistry`)
When any database command receives `--db <arg>`:
1. **Built-in System Keywords:** Resolves predefined strings (`master`, `gitmap`, `heatmap`, `pipeline`, `tasks`, `agents`).
2. **Registry Lookup:** Queries `DatabaseAlias` for `LOWER(AliasName) = LOWER(?) AND IsActive = 1`.
3. **Filesystem Fallback:** Treats `<arg>` as relative or absolute filepath.
4. **Access Telemetry:** On successful match, updates `LastAccessedAt = CURRENT_TIMESTAMP` and increments `SearchHistoryCount` when invoked via search/find operations.

### 4.2 File Size Detection & State Refresh
- During query execution or `dbl --refresh`, the runtime executes `os.Stat(resolvedPath)`.
- If file exists, `FileSizeBytes` is updated in memory and committed asynchronously to SQLite.
- If missing, status flags reflect `MISSING` while keeping historical query telemetry intact.

---

## 5. CLI Command Specification: `db-listing` (`dbl`)

### 5.1 Command Invocation & Aliases
- Primary Command: `gitmap db-listing`
- Short Alias: `gitmap dbl`
- Subcommand: `gitmap db list-aliases`

### 5.2 Flags & Options
- `--json`: Output typed JSON envelope matching `GitMapEnvelope` format.
- `-a, --all`: Include inactive and unresolvable databases.
- `-r, --refresh`: Re-scan filesystem to update `FileSizeBytes` and verify file health.
- `-s, --sort <field>`: Sort by `alias`, `size`, `searches`, or `access` (default: `alias`).

### 5.3 Formatted Terminal Output
Formatted using Lipgloss table boundaries and Dracula palette highlights:

```text
┌──────────────┬───────────┬────────────────────────────────────────┬───────────┬──────────┬─────────────────────┬────────┐
│ ALIAS        │ TYPE      │ PATH                                   │ SIZE      │ SEARCHES │ LAST ACCESSED       │ STATUS │
├──────────────┼───────────┼────────────────────────────────────────┼───────────┼──────────┼─────────────────────┼────────┤
│ master       │ SYSTEM    │ ~/.gitmap/gitmap.db                    │ 1.42 MB   │ 128      │ 2026-10-19 15:40:12 │ OK     │
│ heatmap      │ COMPANION │ repo-cache/.heatmap/01-antigravity/... │ 8.12 MB   │ 43       │ 2026-10-19 15:32:00 │ OK     │
│ pipeline     │ REPO      │ cli/repodb/pipeline.db                 │ 420.50 KB │ 912      │ 2026-10-19 16:01:22 │ OK     │
│ tasks        │ LOCAL     │ .ai-memory/temp-agents/agent-task.db   │ 84.00 KB  │ 12       │ 2026-10-19 14:10:05 │ OK     │
│ custom-proj  │ CUSTOM    │ d:/work/my-project/app.db              │ 2.30 MB   │ 5        │ 2026-10-18 09:21:18 │ OK     │
└──────────────┴───────────┴────────────────────────────────────────┴───────────┴──────────┴─────────────────────┴────────┘
Total: 5 databases registered | Active: 5 | Storage: 12.34 MB
```

### 5.4 JSON Response Envelope

```json
{
  "apiVersion": "v2",
  "status": "success",
  "data": {
    "totalCount": 5,
    "activeCount": 5,
    "totalSizeBytes": 12942000,
    "aliases": [
      {
        "aliasName": "master",
        "databaseType": "SYSTEM",
        "databasePath": "C:/Users/Administrator/.gitmap/gitmap.db",
        "isActive": true,
        "fileSizeBytes": 1488972,
        "fileSizeFormatted": "1.42 MB",
        "searchHistoryCount": 128,
        "lastAccessedAt": "2026-10-19T15:40:12Z",
        "status": "OK"
      }
    ]
  },
  "timestamp": "2026-10-19T16:15:00Z"
}
```

---

## 6. Alias Mutation Commands

To manage custom aliases dynamically:
- `gitmap db alias add <name> <path> [--type <type>] [--description <text>]`
- `gitmap db alias remove <name> [--purge-history]`
- `gitmap db alias toggle <name> [--active=true|false]`

---

## 7. Error Handling & Guardrails

- `E1000:INVALID_ARGUMENT`: Empty alias name or invalid characters.
- `E5001:ALIAS_NOT_FOUND`: Specified alias does not exist in registry.
- `E5002:DB_FILE_UNREACHABLE`: Resolved file does not exist on disk when connection requested.
- `E5003:DUPLICATE_ALIAS`: Cannot register existing alias name without `--force`.

---

## 8. Verification & Acceptance Criteria

1. **Schema Migration:** Running `gitmap db migrate` initializes `DatabaseAlias` in `gitmap.db` without regression.
2. **Listing Output:** `gitmap dbl` executes cleanly, returning non-zero row count for default registered databases.
3. **JSON Envelope:** `gitmap dbl --json` produces valid parseable JSON with all required fields.
4. **Resolution:** Setting `--db custom-alias` correctly routes to target database.
