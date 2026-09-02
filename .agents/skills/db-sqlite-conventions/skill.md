---
name: db-sqlite-conventions
description: >-
  Autonomously design, migrate, and query SQLite database tables in Gitmap adhering to PascalCase schema,
  SetMaxOpenConns(1), filepath.EvalSymlinks binary anchoring, and zero-swallow error policies.
---

# Gitmap SQLite & Database Conventions Skill

Autonomously implement, migrate, and audit SQLite tables and database operations adhering to `spec/04-database-conventions/`, `spec/05-split-db-architecture/`, and `spec/17-consolidated-guidelines/18-database-conventions.md`.

## Core Checkpoints & Mandatory Invariants

1. **Schema & Identifier Standards:**
   - **Singular Table Names:** Strict PascalCase (e.g. `Repo`, `RepoVersionHistory`, `CloneInteractiveSelection`, `InstalledTools`).
   - **Integer Primary Keys:** Always `{TableName}Id INTEGER PRIMARY KEY AUTOINCREMENT`. No UUID primary keys.
   - **Foreign Keys:** Exact PascalCase matching `{TargetTable}Id`.
   - **PascalCase Columns:** All column names must be PascalCase. No snake_case or camelCase.

2. **Mandatory Context Columns (DB Rules 10-12):**
   - **Entity/Master Data Tables:** Must include `Description TEXT NULL` (Rule 10).
   - **Transactional/Event Tables:** Must include `Notes TEXT NULL` and `Comments TEXT NULL` (Rule 11).
   - Context columns must remain nullable (`NULL`) for optional contextual data (Rule 12).

3. **Boolean Column Policy:**
   - Always use affirmative prefixes: `Is*` or `Has*` (e.g. `IsActive`, `IsArchived`, `HasCustomConfig`).
   - Strictly forbidden prefixes: `can`, `should`, `was`, `will`, `did`, `must`, `not`, `no`.
   - Never persist derived inverse booleans; derive in code queries.

4. **Connection & Runtime Management:**
   - Restricted connection pooling: `db.SetMaxOpenConns(1)` to eliminate SQLite database locks and concurrency races.
   - Binary Execution Anchoring: SQLite database path MUST be anchored to binary execution location via `filepath.EvalSymlinks(os.Executable())`.

5. **Migration & Introspection:**
   - All migrations must be strictly idempotent: `CREATE TABLE IF NOT EXISTS`.
   - When altering or checking table columns, use detect-then-act with `PRAGMA table_info(TableName)` instead of brittle error string matching.

6. **Absolute Release Storage Protection:**
   - NEVER manually create, modify, or delete files within `.gitmap/release/` or `.gitmap/release-assets/`.
