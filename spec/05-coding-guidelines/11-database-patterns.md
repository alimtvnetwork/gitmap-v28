# Database Patterns

## Overview

Standards for schema design, migrations, query efficiency, and indexing
across all persistent storage layers.

---

## Schema Conventions

### Naming

| Element        | Convention       | Example              |
|----------------|------------------|----------------------|
| Table names    | PascalCase       | `UserRoles`          |
| Column names   | PascalCase       | `CreatedAt`          |
| Index names    | `idx_Table_Col`  | `idx_Repos_AbsPath`  |
| Foreign keys   | `ParentId`       | `GroupId`, `RepoId`  |

### Column Defaults

| Type      | Convention                          |
|-----------|-------------------------------------|
| Primary   | `INTEGER PRIMARY KEY AUTOINCREMENT` |
| Strings   | `TEXT DEFAULT ''`                   |
| Booleans  | `INTEGER DEFAULT 0`                |
| Timestamps| `TEXT DEFAULT CURRENT_TIMESTAMP`    |

### Required Columns

Every table includes `Id` (primary key) and `CreatedAt` (timestamp).
Tables with mutable rows also include `UpdatedAt`.

---

## Migration Strategies

### Forward-Only Migrations

Migrations are append-only. Never modify or delete a previously applied
migration. Each migration is a single SQL statement or a short batch.

### Additive Changes

Prefer additive changes over destructive ones:

| Change           | Approach                                   |
|------------------|--------------------------------------------|
| Add column       | `ALTER TABLE ... ADD COLUMN` with default   |
| Rename column    | Add new column, migrate data, drop old      |
| Remove column    | Leave unused; remove in next major version  |
| Change type      | Add new column, migrate, drop old           |

### Destructive Migrations

When a breaking change is unavoidable:

1. Detect the legacy schema at startup.
2. Provide clear recovery instructions to the user.
3. Offer an automated migration path (e.g. `db-reset --confirm`).

### Migration Ordering

Name migrations with a numeric prefix for deterministic ordering:

```
001_create_repos.sql
002_create_groups.sql
003_add_notes_column.sql
```

---

## Query Optimization

### Batch Over Loop

Never execute queries inside a loop. Batch reads and writes:

```go
// Bad — N+1
for _, id := range ids {
    row := db.QueryRow("SELECT * FROM Repos WHERE Id = ?", id)
}

// Good — single query
rows := db.Query("SELECT * FROM Repos WHERE Id IN (" + placeholders + ")", args...)
```

### Upsert Pattern

Use `INSERT ... ON CONFLICT(...) DO UPDATE SET` for idempotent writes.
Always specify the conflict target column explicitly.

### Select Only Needed Columns

Avoid `SELECT *` in production code. List columns explicitly to prevent
breakage when schema changes and to reduce transfer overhead.

### Parameterized Queries

Always use parameterized queries (`?` placeholders). Never interpolate
user input into SQL strings.

### Transaction Boundaries

Wrap multi-statement writes in a transaction. Keep transactions short —
read data before the transaction, write inside it.

```go
tx, _ := db.Begin()
// writes only
tx.Exec("INSERT INTO ...")
tx.Exec("UPDATE ...")
tx.Commit()
```

---

## Indexing Rules

### When to Index

| Scenario                    | Action                        |
|-----------------------------|-------------------------------|
| Column in `WHERE` clause    | Add index                     |
| Column in `JOIN` condition  | Add index                     |
| Column in `ORDER BY`        | Consider index                |
| Low-cardinality column      | Skip index (minimal benefit)  |
| Write-heavy, rarely queried | Skip index                    |

### Unique Constraints

Use `UNIQUE` indexes for natural keys (e.g. `AbsolutePath`, `Tag`).
These double as both constraint enforcement and query optimization.

### Composite Indexes

Order columns from most selective to least selective. The leftmost
column must appear in the query's `WHERE` clause for the index to apply.

### Index Naming

```
idx_<Table>_<Column>        -- single column
idx_<Table>_<Col1>_<Col2>   -- composite
```

---

## Connection Management

### Single Connection

For CLI tools using SQLite, open one connection at startup and close it
on exit. Enable WAL mode and foreign keys:

```sql
PRAGMA journal_mode = WAL;
PRAGMA foreign_keys = ON;
```

### Connection Pooling

For server applications, configure pool size based on expected
concurrency. Set idle timeout to prevent stale connections.

### Resource Cleanup

Always close rows, statements, and connections. Use `defer` immediately
after acquisition:

```go
rows, err := db.Query("SELECT ...")
if err != nil {
    return err
}
defer rows.Close()
```

---

## Constraints

- All SQL strings live in the `constants` package.
- All table/column names use PascalCase.
- SQLite driver must be CGo-free (`modernc.org/sqlite`).
- No `SELECT *` in production queries.
- No string interpolation in SQL — parameterized only.
- Every migration is forward-only and additive when possible.
