# Cross-Language Database Naming Convention — PascalCase

> **Version:** 1.0.0  
> **Updated:** 2026-03-09  
> **Applies to:** PHP (SQLite), Go (SQLite), TypeScript (frontend references)

---

## Overview

All **custom** database table names and column names MUST use **PascalCase** across every language in the project. This ensures consistency with the cross-language enum specification, JSON API responses, and TypeScript frontend types.

**WordPress core tables** (e.g., `wp_posts`, `wp_options`, `wp_usermeta`) are **strictly exempt** — they retain their native `snake_case` naming as managed by WordPress itself.

---

## Scope

| Category | Convention | Example |
|----------|-----------|---------|
| Custom SQLite table names | PascalCase | `Transactions`, `AgentSites`, `SnapshotProgress` |
| Custom SQLite column names | PascalCase | `PluginSlug`, `CreatedAt`, `AgentSiteId` |
| Custom SQLite index names | PascalCase with `Idx` prefix | `IdxTransactions_CreatedAt` |
| WordPress core tables | snake_case (exempt) | `wp_posts`, `wp_options` |
| WordPress core columns | snake_case (exempt) | `post_title`, `option_value` |
| Go struct JSON tags | Omit redundant tags | No explicit tag needed — Go marshals PascalCase by default |
| PHP enum-backed table values | PascalCase | `TableType::AgentSites = 'AgentSites'` |

---

## Rules

### Rule 1: Table Names — PascalCase, No Prefix

Custom tables use PascalCase without any `riseup_` or plugin-specific prefix.

```sql
-- ❌ FORBIDDEN
CREATE TABLE agent_sites (...);
CREATE TABLE riseup_transactions (...);
CREATE TABLE snapshot_progress (...);

-- ✅ REQUIRED
CREATE TABLE AgentSites (...);
CREATE TABLE Transaction (...);
CREATE TABLE SnapshotProgress (...);
```

### Rule 2: Column Names — PascalCase

All column names use PascalCase. Abbreviations like `Id`, `Url`, `Md5` are treated as words (only first letter capitalized).

```sql
-- ❌ FORBIDDEN
CREATE TABLE Transaction (
    id INTEGER PRIMARY KEY,
    plugin_slug TEXT,
    created_at TEXT,
    agent_site_id INTEGER
);

-- ✅ REQUIRED
CREATE TABLE Transaction (
    Id INTEGER PRIMARY KEY,
    PluginSlug TEXT,
    CreatedAt TEXT,
    AgentSiteId INTEGER
);
```

### Rule 3: Index Names — PascalCase with `Idx` Prefix

```sql
-- ❌ FORBIDDEN
CREATE INDEX idx_transactions_created ON Transaction(CreatedAt);

-- ✅ REQUIRED
CREATE INDEX IdxTransactions_CreatedAt ON Transaction(CreatedAt);
```

### Rule 4: Abbreviation Casing

Abbreviations are NOT fully capitalized — only the first letter is uppercase:

| ❌ Wrong | ✅ Correct |
|----------|-----------|
| `ID` | `Id` |
| `URL` | `Url` |
| `MD5` | `Md5` |
| `JSON` | `Json` |
| `SQL` | `Sql` |
| `IP` | `Ip` |
| `API` | `Api` |

### Rule 5: WordPress Core — Exempt

Any interaction with WordPress core tables MUST use WordPress's native naming:

```php
// ✅ WordPress core — snake_case required
$wpdb->get_results("SELECT post_title FROM {$wpdb->posts}");
update_option('riseup_schema_version', 13);

// ✅ Custom SQLite — PascalCase required
$this->pdo->query("SELECT PluginSlug FROM Transactions");
```

---

## PHP Implementation

### TableType Enum

The `TableType` enum's backed values must match the PascalCase table names:

```php
enum TableType: string
{
    case Transactions     = 'Transactions';
    case AgentSites       = 'AgentSites';
    case AgentActions     = 'AgentActions';
    case Snapshots        = 'Snapshots';
    case SnapshotProgress = 'SnapshotProgress';
    case SnapshotJobs     = 'SnapshotJobs';
    case SnapshotSettings = 'SnapshotSettings';
    case SnapshotExports  = 'SnapshotExports';
    case FileCache        = 'FileCache';
}
```

### SQL Queries

```php
// ❌ FORBIDDEN
$sql = "SELECT plugin_slug, created_at FROM transactions WHERE status = ?";

// ✅ REQUIRED
$sql = "SELECT PluginSlug, CreatedAt FROM Transaction WHERE Status = ?";
```

### Orm Class

The `Orm` class receives PascalCase table names via `TableType->value` and all column references in conditions, orderBy, and data arrays must use PascalCase keys:

```php
$orm = new Orm(TableType::Transactions->value);
$orm->findAll(['Status' => 'Pending'], 'CreatedAt DESC', 10);
$orm->insert(['PluginSlug' => $slug, 'CreatedAt' => gmdate('Y-m-d H:i:s')]);
```

---

## Go Implementation

### Struct Tags

Go struct `db` tags must match PascalCase column names. **Omit redundant `json` tags** — Go marshals PascalCase field names by default:

```go
type Project struct {
    Id          int64  `db:"Id"`
    ProjectId   string `db:"ProjectId"`
    DisplayName string `db:"DisplayName"`
    CreatedAt   string `db:"CreatedAt"`
}
```

### SQL Queries

```go
// ❌ FORBIDDEN
const queryGetProject = `SELECT project_id, display_name FROM projects WHERE id = ?`

// ✅ REQUIRED
const queryGetProject = `SELECT ProjectId, DisplayName FROM Project WHERE Id = ?`
```

### Migrations

Go migrations use `ALTER TABLE ... RENAME TO` for tables and `ALTER TABLE ... RENAME COLUMN` for columns (SQLite ≥ 3.25):

```go
func migrateV2(db *sql.DB) error {
    statements := []string{
        `ALTER TABLE projects RENAME TO Projects`,
        `ALTER TABLE Projects RENAME COLUMN project_id TO ProjectId`,
        `ALTER TABLE Projects RENAME COLUMN display_name TO DisplayName`,
    }

    for _, stmt := range statements {
        if _, err := db.Exec(stmt); err != nil {
            return apperror.Wrap(err, "E2050", "PascalCase migration failed")
        }
    }

    return nil
}
```

---

## TypeScript Frontend

TypeScript types referencing database columns must use PascalCase to match the API response format:

```typescript
interface Transaction {
    Id: number;
    PluginSlug: string;
    CreatedAt: string;
    Status: string;
    AgentSiteId: number | null;
}
```

---

## Migration Strategy

A 5-phase migration plan converts all existing snake_case tables and columns to PascalCase:

| Phase | Scope | Description |
|-------|-------|-------------|
| Phase 1 | Specs & Standards | This document + memory updates + coding guidelines |
| Phase 2 | Go Backend | SplitDB (3 tables) + E2E Service (4 tables) |
| Phase 3 | PHP Plugin | 12 tables via migration v13 + `TableType` enum update |
| Phase 4 | PHP Root DB | 5 per-snapshot tables + backward compatibility layer |
| Phase 5 | Validation | Full grep sweep + test suite + snapshot round-trip test |

### Safety Rules

1. All table renames use `ALTER TABLE ... RENAME TO ...`
2. All column renames use `ALTER TABLE ... RENAME COLUMN ... TO ...` (SQLite ≥ 3.25)
3. PHP migration v13 runs inside a transaction with rollback on failure
4. Go migrations are version-tracked and idempotent
5. Root DB backward compatibility: new code must detect and handle old snake_case snapshots
6. `TableType` enum values update simultaneously with migration execution
7. Index names must also be updated to PascalCase

---

## Common Mistakes — Database Naming

These are real violations found and fixed. Use as a reference to avoid repeating them.

### Mistake 1: Mixed Casing in PHP Insert Arrays

```php
// ❌ WRONG — mixed casing causes NULL values in the database
$record = array(
    'sequence'       => $seq,          // lowercase — won't match 'Sequence' column
    'filename'       => $name,         // lowercase — won't match 'Filename' column
    'totalRows'      => $rows,         // camelCase — won't match 'TotalRows' column
    'triggerSource'  => $trigger,      // camelCase — won't match 'TriggerSource' column
);

// ✅ CORRECT — every key must be PascalCase
$record = array(
    'Sequence'       => $seq,
    'Filename'       => $name,
    'TotalRows'      => $rows,
    'TriggerSource'  => $trigger,
);
```

**Impact:** SQLite silently ignores unrecognized column names, resulting in NULL values in the database with no error.

### Mistake 2: snake_case in SQL Queries After Migration

```php
// ❌ WRONG — using old column names
$sql = "SELECT plugin_slug, created_at FROM transactions WHERE status = ?";

// ✅ CORRECT — PascalCase post-migration
$sql = "SELECT PluginSlug, CreatedAt FROM Transaction WHERE Status = ?";
```

### Mistake 3: Uppercase Abbreviations in Column Names

```sql
-- ❌ WRONG
CREATE TABLE Project (ID INTEGER, ProjectID TEXT, URL TEXT);

-- ✅ CORRECT
CREATE TABLE Project (Id INTEGER, ProjectId TEXT, Url TEXT);
```

### Mistake 4: Go Struct Tags Not Matching Schema

```go
// ❌ WRONG — tags use Go field casing, not schema casing
type Project struct {
    ID        int64  `db:"id"`
    Name      string `db:"name"`
}

// ✅ CORRECT — db tags match PascalCase schema, no redundant json tags
type Project struct {
    Id        int64  `db:"Id"`
    Name      string `db:"Name"`
}
```

---

## Cross-Reference

- [Cross-Language Code Style](./04-code-style/00-overview.md) — Formatting rules
- [PHP Naming Conventions](../04-php/00-overview.md) — PHP-specific naming
- [Go Coding Standards](../03-golang/04-golang-standards-reference/00-overview.md) — Go-specific naming
- [Master Coding Guidelines](./15-master-coding-guidelines/00-overview.md) — Consolidated cross-language reference
- [Issues & Fixes Log](./01-issues-and-fixes-log.md) — Full historical fixes
