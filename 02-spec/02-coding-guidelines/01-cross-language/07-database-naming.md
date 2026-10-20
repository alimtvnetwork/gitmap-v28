# Cross-Language Database Naming Convention — PascalCase (AI Execution Prompt)

> **/goal** Enforce strict PascalCase naming across all custom database tables, columns, indexes, and associative data keys in PHP, Go, and TypeScript, while strictly preserving snake_case for WordPress core tables.
> **/learn** Understand the cross-language database naming contract, PascalCase table/column/index schema definitions, abbreviation casing rules (`Id`, `Url`, `Md5`), elimination of table prefixes, and zero-redundancy Go struct tagging.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Enforce PascalCase for all custom database tables (`AgentSites`, `Transaction`, `SnapshotProgress`) without table prefixes.
- [ ] `/learn` Keep WordPress core tables and columns strictly in their native `snake_case` (e.g., `wp_posts`, `option_value`).
- [ ] `/goal` Ensure all custom column names and composite index names use PascalCase with single-capitalized abbreviations (`Id`, `Url`, `CreatedAt`, `IdxTransactions_CreatedAt`).
- [ ] `/learn` Verify PHP array record keys and Go struct `db` tags precisely mirror PascalCase schema definitions without redundant JSON tags.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

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
|---|---|---|
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

Custom tables use PascalCase without any `riseup_` or plugin-specific prefix:
- ❌ **Forbidden:** `CREATE TABLE agent_sites (...);`, `CREATE TABLE riseup_transactions (...);`
- ✅ **Required:** `CREATE TABLE AgentSites (...);`, `CREATE TABLE Transaction (...);`

### Rule 2: Column Names — PascalCase

All column names use PascalCase. Abbreviations (`Id`, `Url`, `Md5`) treat only the first letter as capitalized:
- ❌ **Forbidden:** `id`, `plugin_slug`, `created_at`, `agent_site_id`
- ✅ **Required:** `Id`, `PluginSlug`, `CreatedAt`, `AgentSiteId`

### Rule 3: Index Names — PascalCase with `Idx` Prefix

- ❌ **Forbidden:** `CREATE INDEX idx_transactions_created ON Transaction(CreatedAt);`
- ✅ **Required:** `CREATE INDEX IdxTransactions_CreatedAt ON Transaction(CreatedAt);`

### Rule 4: Abbreviation Casing

Abbreviations are NOT fully capitalized — only the first letter is uppercase: `Id` (not `ID`), `Url` (not `URL`), `Md5` (not `MD5`), `Json` (not `JSON`), `Sql` (not `SQL`), `Ip` (not `IP`), `Api` (not `API`).

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

### TableType Enum & Orm Class

The `TableType` enum's backed values match PascalCase table names, and `Orm` column references use PascalCase:

```php
enum TableType: string {
    case Transactions     = 'Transactions';
    case AgentSites       = 'AgentSites';
    case AgentActions     = 'AgentActions';
    case Snapshots        = 'Snapshots';
    case FileCache        = 'FileCache';
}

$orm = new Orm(TableType::Transactions->value);
$orm->findAll(['Status' => 'Pending'], 'CreatedAt DESC', 10);
$orm->insert(['PluginSlug' => $slug, 'CreatedAt' => gmdate('Y-m-d H:i:s')]);
```

### SQL Queries

```php
// ❌ FORBIDDEN: $sql = "SELECT plugin_slug, created_at FROM transactions WHERE status = ?";
// ✅ REQUIRED:  $sql = "SELECT PluginSlug, CreatedAt FROM Transaction WHERE Status = ?";
```

---

## Go Implementation

### Struct Tags & SQL Queries

Go struct `db` tags match PascalCase column names. **Omit redundant `json` tags** — Go marshals PascalCase field names by default:

```go
type Project struct {
    Id          int64  `db:"Id"`
    ProjectId   string `db:"ProjectId"`
    DisplayName string `db:"DisplayName"`
    CreatedAt   string `db:"CreatedAt"`
}

// ❌ FORBIDDEN: const queryGet = `SELECT project_id FROM projects WHERE id = ?`
// ✅ REQUIRED:  const queryGet = `SELECT ProjectId FROM Project WHERE Id = ?`
```

### Migrations

Go migrations use `ALTER TABLE ... RENAME TO` and `RENAME COLUMN` (SQLite ≥ 3.25):

```go
func migrateV2(db *sql.DB) error {
    stmts := []string{
        `ALTER TABLE projects RENAME TO Projects`,
        `ALTER TABLE Projects RENAME COLUMN project_id TO ProjectId`,
    }
    for _, stmt := range stmts {
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

## Native Database CLI Tooling (GitMap db-*)

To safely inspect, query, and mutate databases adhering to PascalCase schema conventions across polyglot workspaces, developers and AI agents MUST use the compiled native pure-Go `gitmap db-*` command suite instead of ad-hoc Python SQLite one-liners (`python -c "import sqlite3..."`):

| Canonical Command | Short Alias | Primary Function & PascalCase Inspection | Safety Guarantee |
|---|---|---|---|
| `gitmap db-read` | `gitmap dbr` | Query tables, list schema (`ls`), inspect columns (`table <T>`), raw SQL | Strict read-only (`PRAGMA query_only = ON;`) |
| `gitmap db-add` | `gitmap dba` | Insert records via key-value pairs (`--values "PluginSlug=...,Status=..."`) | Atomic transaction (`BEGIN IMMEDIATE`) |
| `gitmap db-edit` | `gitmap dbe` | Update existing records by primary key ID (`--id <id> --set "Status=..."`) | Atomic update (`BEGIN IMMEDIATE`) |
| `gitmap db-find` | `gitmap dbf` | Multi-column string/prefix/suffix and regex search (`--grep`) | Read-only scan across text columns |
| `gitmap db-listing` | `gitmap dbl` | Inventory registered database aliases, file paths, byte sizes, and access counts | In-memory query against `gitmap.db` |
| `gitmap db-view` | `gitmap dbv` | Relational projections with automatic join inference across PascalCase foreign keys | Read-only multi-table view projection |

### Command Syntax & Examples

```bash
# 1. Enumerate tables and inspect row/column counts across Split-DB tiers
gitmap db-read ls
gitmap dbr ls -as install --json

# 2. Inspect table schema and PascalCase column definitions
gitmap db-read table Transaction -c "Id,PluginSlug,Status" -l 10
gitmap dbr table Repo -as master --json

# 3. Read-only SQL query with PascalCase identifiers
gitmap db-read "SELECT Id, PluginSlug, CreatedAt FROM Transaction WHERE Status = 'Active'"

# 4. Search across text columns using wildcards or regex
gitmap db-find "riseup*" -t Transaction -c "PluginSlug"
gitmap dbf --grep "^gitmap-" -as master

# 5. Relational projection with automatic join inference
gitmap db-view -t "Project as p, Task as t" -c "p.ProjectId, p.DisplayName, t.TaskSlug"

# 6. Atomic record insertion and update (single-writer safety)
gitmap db-add Transaction --values "PluginSlug=demo,Status=Active,CreatedAt=2026-03-09"
gitmap db-edit Transaction --id 1 --set "Status=Completed"
```

---

## Migration Strategy

A 5-phase migration plan converts all existing snake_case tables and columns to PascalCase:

| Phase | Scope | Description |
|---|---|---|
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

### Mistake 1: Mixed Casing in PHP Insert Arrays

```php
// ❌ WRONG — mixed casing causes NULL values in the database
$record = ['sequence' => $seq, 'filename' => $name, 'totalRows' => $rows];

// ✅ CORRECT — every key must be PascalCase
$record = ['Sequence' => $seq, 'Filename' => $name, 'TotalRows' => $rows];
```

### Mistake 2: snake_case in SQL Queries After Migration

```php
// ❌ WRONG: $sql = "SELECT plugin_slug, created_at FROM transactions WHERE status = ?";
// ✅ CORRECT: $sql = "SELECT PluginSlug, CreatedAt FROM Transaction WHERE Status = ?";
```

### Mistake 3: Uppercase Abbreviations in Column Names

```sql
-- ❌ WRONG:   CREATE TABLE Project (ID INTEGER, ProjectID TEXT, URL TEXT);
-- ✅ CORRECT: CREATE TABLE Project (Id INTEGER, ProjectId TEXT, Url TEXT);
```

### Mistake 4: Go Struct Tags Not Matching Schema

```go
// ❌ WRONG — tags use Go field casing, not schema casing
type Project struct { ID int64 `db:"id"`; Name string `db:"name"` }

// ✅ CORRECT — db tags match PascalCase schema, no redundant json tags
type Project struct { Id int64 `db:"Id"`; Name string `db:"Name"` }
```

---

## Cross-Reference

- [Cross-Language Code Style](./04-code-style/readme.md) — Formatting rules
- [PHP Naming Conventions](../04-php/readme.md) — PHP-specific naming
- [Go Coding Standards](../03-golang/04-golang-standards-reference/readme.md) — Go-specific naming
- [GitMap DB Suite](../../../.agents/skills/gitmap/SKILL.md) — Native CLI database commands (`db-read`, `db-add`, `db-edit`, `db-find`, `db-listing`, `db-view`, `db-rest`, `db-record-del`, `db-rm`)
- [Master Coding Guidelines](./15-master-coding-guidelines/readme.md) — Consolidated cross-language reference
- [Issues & Fixes Log](./02-issues-and-fixes-log.md) — Full historical fixes

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-NAME-007: PascalCase Database Naming Compliance

**Given** Database schema migrations, ORM entities, SQL queries, and struct definitions across PHP, Go, and TypeScript.
**When** Codebase linters, migration scanners, or CI autofixers validate database definitions and query statements.
**Then** All custom tables, columns, indexes, and data record keys strictly adhere to PascalCase with word-cased abbreviations (`Id`, `Url`), while WordPress core tables remain in exempt snake_case with zero violations.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only
```
**Expected:** exit 0. Zero violations.
