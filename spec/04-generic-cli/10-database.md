# Local Database

> **Related specs:**
> - [05-configuration.md](05-configuration.md) — config layer that coexists with DB persistence
> - [18-batch-execution.md](18-batch-execution.md) — DB-based repo loading for batch operations
> - [02-project-structure.md](02-project-structure.md) — `store/` package placement for DB code

## Overview

Use a local SQLite database for persistence. This enables slug-based
lookup, grouping, batch operations, and history tracking without
external dependencies.

## SQLite Setup

### Driver

Use a **CGo-free** SQLite driver (e.g., `modernc.org/sqlite` for Go).
No C compiler required.

### Location

| Item | Value |
|------|-------|
| Directory | `toolname-output/data/` (auto-created) |
| File name | `toolname.db` |

### Auto-Creation

On first data-producing command:
1. Check if database exists.
2. If missing, create it and initialize all tables.
3. Upsert data into tables.

## Schema Conventions

| Convention | Detail |
|------------|--------|
| Table names | PascalCase (`Repos`, `Groups`) |
| Column names | PascalCase (`RepoName`, `AbsolutePath`) |
| Primary keys | `Id TEXT PRIMARY KEY` (UUID) |
| Timestamps | `TEXT DEFAULT CURRENT_TIMESTAMP` |
| Booleans | `INTEGER DEFAULT 0` (0/1) |
| String defaults | `DEFAULT ''` (never NULL) |

## Core Tables

### Items Table (e.g., Repos)

```sql
CREATE TABLE IF NOT EXISTS Repos (
    Id            TEXT PRIMARY KEY,
    Slug          TEXT NOT NULL,
    Name          TEXT NOT NULL,
    HttpsUrl      TEXT NOT NULL,
    AbsolutePath  TEXT NOT NULL,
    CreatedAt     TEXT DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt     TEXT DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_repos_path ON Repos(AbsolutePath);
```

### Groups Table

```sql
CREATE TABLE IF NOT EXISTS Groups (
    Id          TEXT PRIMARY KEY,
    Name        TEXT NOT NULL UNIQUE,
    Description TEXT DEFAULT '',
    CreatedAt   TEXT DEFAULT CURRENT_TIMESTAMP
);
```

### Join Table

```sql
CREATE TABLE IF NOT EXISTS GroupItems (
    GroupId TEXT NOT NULL REFERENCES Groups(Id) ON DELETE CASCADE,
    ItemId  TEXT NOT NULL REFERENCES Repos(Id) ON DELETE CASCADE,
    PRIMARY KEY (GroupId, ItemId)
);
```

### History Table

```sql
CREATE TABLE IF NOT EXISTS CommandHistory (
    Id         TEXT PRIMARY KEY,
    Command    TEXT NOT NULL,
    Args       TEXT DEFAULT '',
    StartedAt  TEXT NOT NULL,
    DurationMs INTEGER DEFAULT 0,
    ExitCode   INTEGER DEFAULT 0,
    CreatedAt  TEXT DEFAULT CURRENT_TIMESTAMP
);
```

## Upsert Strategy

Match by unique field (e.g., `AbsolutePath`). If exists, update.
Otherwise, insert.

```go
const UpsertRepo = `
INSERT INTO Repos (Id, Slug, Name, AbsolutePath)
VALUES (?, ?, ?, ?)
ON CONFLICT(AbsolutePath) DO UPDATE SET
    Slug = excluded.Slug,
    Name = excluded.Name,
    UpdatedAt = CURRENT_TIMESTAMP
`
```

## DB-First Lookup with Fallback

Commands that resolve items by slug:
1. Try the database first.
2. Fall back to JSON file if database doesn't exist.

## Store Package Structure

```
store/
├── store.go     DB init, open, close, migration, reset
├── repo.go      Item CRUD (upsert, list, find by slug)
├── group.go     Group CRUD
└── history.go   History insert + query
```

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
