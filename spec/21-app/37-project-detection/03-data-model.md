# Project Type Detection — Data Model

## Table Relationships

```
ProjectTypes (1) ──── (N) DetectedProjects (1) ──── (0..1) GoProjectMetadata
                              │                                │
                              │                      (1) ──── (N) GoRunnableFiles
                              │
                              ├──── (0..1) CSharpProjectMetadata
                              │               │
                              │     (1) ──── (N) CSharpProjectFiles
                              │     (1) ──── (N) CSharpKeyFiles
                              │
Repos (1) ─────────── (N) DetectedProjects
```

---

## ProjectTypes Table

A reference table for all supported project types.

| Column      | Type    | Constraints       | Notes                          |
|-------------|---------|-------------------|--------------------------------|
| Id          | TEXT    | PRIMARY KEY       | UUID                           |
| Key         | TEXT    | NOT NULL, UNIQUE  | `go`, `node`, `react`, `cpp`, `csharp` |
| Name        | TEXT    | NOT NULL          | Display name (e.g., `Go`, `Node.js`) |
| Description | TEXT    | DEFAULT ''        | Human-readable description     |

**Seeding:** This table is seeded during migration with all supported
project types. The `Id` values are stable UUIDs (v4) defined in
constants. Example: `"b3f1a2c4-..."` — never short strings like `"pt-go"`.

---

## DetectedProjects Table

| Column           | Type    | Constraints                                   | Notes                                    |
|------------------|---------|-----------------------------------------------|------------------------------------------|
| Id               | TEXT    | PRIMARY KEY                                   | UUID                                     |
| RepoId           | TEXT    | NOT NULL, FK → Repos(Id) ON DELETE CASCADE    | Link to parent repo                      |
| ProjectTypeId    | TEXT    | NOT NULL, FK → ProjectTypes(Id)               | Link to project type                     |
| ProjectName      | TEXT    | NOT NULL                                      | Parsed from manifest or dir name         |
| AbsolutePath     | TEXT    | NOT NULL                                      | Full filesystem path to project root     |
| RepoPath         | TEXT    | NOT NULL                                      | Absolute path of the Git repo root       |
| RelativePath     | TEXT    | NOT NULL                                      | Path relative to repo root               |
| PrimaryIndicator | TEXT    | NOT NULL                                      | File that triggered detection            |
| DetectedAt       | TEXT    | DEFAULT CURRENT_TIMESTAMP                     |                                          |

**Unique constraint:** `(RepoId, ProjectTypeId, RelativePath)` — one
entry per project type per path per repo.

**Upsert strategy:** On scan, match by the unique constraint. If a
row exists, update `ProjectName`, `AbsolutePath`, `PrimaryIndicator`,
and `DetectedAt`. Otherwise, insert a new row.

### Stale Entry Cleanup

On each scan, after upserting all detected projects for a repo, delete
any `DetectedProjects` rows for that `RepoId` that were **not**
upserted in the current scan. This handles removed projects.

Because `GoProjectMetadata`, `CSharpProjectMetadata`, and their child
tables use `ON DELETE CASCADE` from `DetectedProjects`, deleting a
stale `DetectedProjects` row automatically removes all associated
metadata, runnables, project files, and key files.

### Orphaned Child Record Cleanup

When a parent project **still exists** but child records change
(e.g., a Go `cmd/` subdirectory is deleted, or a C# `.csproj` is
removed), the cleanup works per-table:

- **GoRunnableFiles:** After upserting all runnables for a
  `GoMetadataId`, delete rows for that `GoMetadataId` whose `Id` was
  **not** in the current upsert batch.
- **CSharpProjectFiles:** Same pattern per `CSharpMetadataId`.
- **CSharpKeyFiles:** Same pattern per `CSharpMetadataId`.

---

## SQL Statements

### Create ProjectTypes

```sql
CREATE TABLE IF NOT EXISTS ProjectTypes (
    Id          TEXT PRIMARY KEY,
    Key         TEXT NOT NULL UNIQUE,
    Name        TEXT NOT NULL,
    Description TEXT DEFAULT ''
)
```

### Seed ProjectTypes

```sql
INSERT OR IGNORE INTO ProjectTypes (Id, Key, Name, Description) VALUES
    ('b3f1a2c4-5d6e-4f7a-8b9c-0d1e2f3a4b5c', 'go',     'Go',      'Go modules and packages'),
    ('c4d2b3e5-6f7a-4e8b-9c0d-1e2f3a4b5c6d', 'node',   'Node.js', 'Node.js projects'),
    ('d5e3c4f6-7a8b-4f9c-0d1e-2f3a4b5c6d7e', 'react',  'React',   'React applications'),
    ('e6f4d5a7-8b9c-4a0d-1e2f-3a4b5c6d7e8f', 'cpp',    'C++',     'C and C++ projects'),
    ('f7a5e6b8-9c0d-4b1e-2f3a-4b5c6d7e8f9a', 'csharp', 'C#',      '.NET and C# projects')
```

### Create DetectedProjects

```sql
CREATE TABLE IF NOT EXISTS DetectedProjects (
    Id               TEXT PRIMARY KEY,
    RepoId           TEXT NOT NULL REFERENCES Repos(Id) ON DELETE CASCADE,
    ProjectTypeId    TEXT NOT NULL REFERENCES ProjectTypes(Id),
    ProjectName      TEXT NOT NULL,
    AbsolutePath     TEXT NOT NULL,
    RepoPath         TEXT NOT NULL,
    RelativePath     TEXT NOT NULL,
    PrimaryIndicator TEXT NOT NULL,
    DetectedAt       TEXT DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(RepoId, ProjectTypeId, RelativePath)
)
```

### Upsert DetectedProject

```sql
INSERT INTO DetectedProjects (Id, RepoId, ProjectTypeId, ProjectName,
    AbsolutePath, RepoPath, RelativePath, PrimaryIndicator)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(RepoId, ProjectTypeId, RelativePath) DO UPDATE SET
    ProjectName=excluded.ProjectName,
    AbsolutePath=excluded.AbsolutePath,
    RepoPath=excluded.RepoPath,
    PrimaryIndicator=excluded.PrimaryIndicator,
    DetectedAt=CURRENT_TIMESTAMP
```

### Query by Type Key

```sql
SELECT dp.Id, dp.RepoId, pt.Key AS ProjectType, dp.ProjectName,
    dp.AbsolutePath, dp.RepoPath, dp.RelativePath,
    dp.PrimaryIndicator, dp.DetectedAt,
    r.RepoName
FROM DetectedProjects dp
JOIN ProjectTypes pt ON dp.ProjectTypeId = pt.Id
JOIN Repos r ON dp.RepoId = r.Id
WHERE pt.Key = ?
ORDER BY r.RepoName, dp.RelativePath
```

### Cleanup Stale

```sql
DELETE FROM DetectedProjects
WHERE RepoId = ? AND Id NOT IN (?, ?, ...)
```

### Drop Tables

Drop order must respect foreign key dependencies (children first):

```sql
DROP TABLE IF EXISTS GoRunnableFiles
DROP TABLE IF EXISTS GoProjectMetadata
DROP TABLE IF EXISTS CSharpKeyFiles
DROP TABLE IF EXISTS CSharpProjectFiles
DROP TABLE IF EXISTS CSharpProjectMetadata
DROP TABLE IF EXISTS DetectedProjects
DROP TABLE IF EXISTS ProjectTypes
```
