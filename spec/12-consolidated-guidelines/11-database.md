# 11 — Database Patterns

Schema design, migrations, query optimization, and connection management.

## Schema Conventions (v15)

| Element | Convention | Example |
|---------|-----------|---------|
| Table names | PascalCase, **singular** | `Repo`, `Group`, `Release` |
| Column names | PascalCase | `CreatedAt` |
| Primary key | `{TableName}Id INTEGER PRIMARY KEY AUTOINCREMENT` | `RepoId`, `ReleaseId` |
| Foreign keys | Match referenced PK name | `RepoId`, `GroupId`, `CsharpProjectMetadataId` |
| Booleans | `INTEGER DEFAULT 0`, **`IsX` prefix** | `IsDraft`, `IsPreRelease`, `IsLatest` |
| Strings | `TEXT DEFAULT ''` | |
| Timestamps | `TEXT DEFAULT CURRENT_TIMESTAMP` | `CreatedAt`, `UpdatedAt` |
| Reserved words | Double-quoted in DDL/DML | `"Group"(GroupId)` |
| Abbreviations | Treated as words | `SshKey` (not `SSHKey`), `CsharpProjectMetadata` (not `CSharpProjectMetadata`) |

Every table includes `{TableName}Id` (or a natural key like `Setting.Key`) and `CreatedAt` where applicable.

Source of truth: <https://github.com/alimtvnetwork/coding-guidelines-v15/blob/main/spec/04-database-conventions/01-naming-conventions.md>

## Migrations

Forward-only, append-only. Prefer additive changes (add columns with defaults). For breaking changes: detect legacy schema, provide recovery instructions.

## Query Rules

- Batch over loop (no N+1).
- Upsert with `INSERT ... ON CONFLICT ... DO UPDATE`.
- No `SELECT *` — list columns explicitly.
- Parameterized queries only — no string interpolation.
- Wrap multi-statement writes in transactions.

## Connection Management

CLI tools: single connection, `SetMaxOpenConns(1)`. Enable WAL mode and foreign keys. Always `defer rows.Close()`.

## Constraints

All SQL strings in the `constants` package. SQLite driver: `modernc.org/sqlite` (CGo-free). No UUID primary keys.

---

Source: `spec/05-coding-guidelines/11-database-patterns.md`
