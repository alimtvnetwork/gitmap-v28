# Database & Repo Storage

## Overview

After scanning, gitmap persists all discovered repositories in a local
SQLite database. The database enables slug-based lookup, repo grouping,
batch operations, and release history tracking.

## Naming Convention

All table names and column names use **PascalCase** (e.g. `Repos`, `RepoName`, `AbsolutePath`).

## SQLite Setup

### Package

Use a **CGo-free** SQLite driver. The recommended package is
`modernc.org/sqlite` (pure Go, no C compiler required).

### Database Location

The database is located relative to the **binary's physical installation
directory**, not the current working directory. This ensures all commands
access the same database regardless of where the user invokes gitmap.

| Item | Value |
|------|-------|
| Resolution | `os.Executable()` + `filepath.EvalSymlinks()` |
| Directory | `<binary-dir>/data/` (created automatically) |
| File name | `gitmap.db` |
| Full path | `<binary-dir>/data/gitmap.db` |

**Important**: Scan output files (CSV, JSON, scripts) still write to the
user-specified output directory. Only the database is anchored to the
binary location.

### Auto-Creation

On every `scan` completion, gitmap:

1. Checks if `.gitmap/output/data/gitmap.db` exists.
2. If missing, creates the database and initializes all tables.
3. Upserts all scanned repos into the `Repos` table.

---

## Data Model

### Repos Table

| Column           | Type    | Constraints          | Notes                            |
|------------------|---------|----------------------|----------------------------------|
| Id               | TEXT    | PRIMARY KEY          | UUID from ScanRecord             |
| Slug             | TEXT    | NOT NULL             | Derived from GitHub repo name    |
| RepoName         | TEXT    | NOT NULL             | Display name                     |
| HttpsUrl         | TEXT    | NOT NULL             |                                  |
| SshUrl           | TEXT    | NOT NULL             |                                  |
| Branch           | TEXT    | NOT NULL             |                                  |
| RelativePath     | TEXT    | NOT NULL             |                                  |
| AbsolutePath     | TEXT    | NOT NULL, UNIQUE IDX |                                  |
| CloneInstruction | TEXT    | NOT NULL             |                                  |
| Notes            | TEXT    | DEFAULT ''           |                                  |
| CreatedAt        | TEXT    | DEFAULT CURRENT_TIMESTAMP |                             |
| UpdatedAt        | TEXT    | DEFAULT CURRENT_TIMESTAMP |                             |

**Upsert strategy:** On scan, match by `AbsolutePath`. If a row with
that path exists, update all fields. Otherwise, insert a new row.

### Groups Table

| Column      | Type | Constraints               | Notes                |
|-------------|------|---------------------------|----------------------|
| Id          | TEXT | PRIMARY KEY               | UUID                 |
| Name        | TEXT | NOT NULL, UNIQUE          | Group display name   |
| Description | TEXT | DEFAULT ''                | Optional description |
| Color       | TEXT | DEFAULT ''                | Terminal color       |
| CreatedAt   | TEXT | DEFAULT CURRENT_TIMESTAMP |                      |

### GroupRepos Table (Join)

| Column  | Type | Constraints                              | Notes |
|---------|------|------------------------------------------|-------|
| GroupId | TEXT | NOT NULL, FK → Groups(Id) ON DELETE CASCADE | |
| RepoId  | TEXT | NOT NULL, FK → Repos(Id) ON DELETE CASCADE  | |
| | | PRIMARY KEY (GroupId, RepoId) | |

### Releases Table

| Column       | Type    | Constraints               | Notes                              |
|--------------|---------|---------------------------|------------------------------------|
| Id           | TEXT    | PRIMARY KEY               | UUID                               |
| Version      | TEXT    | NOT NULL                  | Core version string (e.g. 2.14.0)  |
| Tag          | TEXT    | NOT NULL, UNIQUE          | Git tag (e.g. v2.14.0)             |
| Branch       | TEXT    | NOT NULL                  | Release branch name                |
| SourceBranch | TEXT    | NOT NULL                  | Branch release was created from    |
| CommitSha    | TEXT    | NOT NULL                  | Full commit SHA                    |
| Changelog    | TEXT    | DEFAULT ''                | Newline-separated changelog notes  |
| Draft        | INTEGER | DEFAULT 0                 | 1 = draft release                  |
| PreRelease   | INTEGER | DEFAULT 0                 | 1 = pre-release                    |
| IsLatest     | INTEGER | DEFAULT 0                 | 1 = latest stable release          |
| CreatedAt    | TEXT    | DEFAULT CURRENT_TIMESTAMP |                                    |

**Upsert strategy:** On release, match by `Tag`. If a release with that
tag exists, update all fields. Otherwise, insert a new row. When a new
stable release is marked as latest, all other releases have `IsLatest`
cleared to 0 first.

### Amendments Table

| Column        | Type    | Constraints               | Notes                              |
|---------------|---------|---------------------------|------------------------------------|
| Id            | TEXT    | PRIMARY KEY               | UUID                               |
| Branch        | TEXT    | NOT NULL                  | Target branch name                 |
| FromCommit    | TEXT    | NOT NULL                  | First commit SHA in range          |
| ToCommit      | TEXT    | NOT NULL                  | Last commit SHA (HEAD at amend)    |
| TotalCommits  | INTEGER | NOT NULL                  | Number of commits rewritten        |
| PreviousName  | TEXT    | DEFAULT ''                | Original author name               |
| PreviousEmail | TEXT    | DEFAULT ''                | Original author email              |
| NewName       | TEXT    | DEFAULT ''                | Replacement author name            |
| NewEmail      | TEXT    | DEFAULT ''                | Replacement author email           |
| Mode          | TEXT    | NOT NULL                  | `all`, `range`, or `head`          |
| ForcePushed   | INTEGER | DEFAULT 0                 | 1 = force-push was executed        |
| CreatedAt     | TEXT    | DEFAULT CURRENT_TIMESTAMP |                                    |

**Insert only:** Each amend operation inserts a new row. No upsert — every
operation is a unique audit record.

### CommandHistory Table

| Column     | Type    | Constraints               | Notes                              |
|------------|---------|---------------------------|------------------------------------|
| Id         | TEXT    | PRIMARY KEY               | Timestamp-based unique ID          |
| Command    | TEXT    | NOT NULL                  | CLI command name                   |
| Alias      | TEXT    | DEFAULT ''                | Alias if used                      |
| Args       | TEXT    | DEFAULT ''                | Positional arguments               |
| Flags      | TEXT    | DEFAULT ''                | Flags passed                       |
| StartedAt  | TEXT    | NOT NULL                  | RFC3339 start timestamp            |
| FinishedAt | TEXT    | DEFAULT ''                | RFC3339 end timestamp              |
| DurationMs | INTEGER | DEFAULT 0                 | Execution time in milliseconds     |
| ExitCode   | INTEGER | DEFAULT 0                 | 0 = success                        |
| Summary    | TEXT    | DEFAULT ''                | Result summary                     |
| RepoCount  | INTEGER | DEFAULT 0                 | Repos affected                     |
| CreatedAt  | TEXT    | DEFAULT CURRENT_TIMESTAMP |                                    |

**Insert + update strategy:** A record is inserted at command start, then
updated with completion details (duration, exit code, summary) at end.

---

## Slug Generation

`mapper.BuildRecords` populates the `Slug` field on every `ScanRecord`
during scan, so it is available in both JSON output and DB upsert.

Extract from HTTPS URL:

```
https://github.com/user/my-api.git  →  my-api
https://github.com/org/my-api.git   →  my-api  (duplicate allowed)
```

**Algorithm:**

1. Parse the HTTPS URL.
2. Take the last path segment.
3. Strip `.git` suffix if present.
4. Lowercase the result.

If the HTTPS URL is empty, fall back to `repoName`.

---

## Package Structure (Database)

### Packages

| Package | Responsibility |
|---------|----------------|
| `store` | SQLite database init, connection, CRUD operations |

### Files

| File | Contents |
|------|----------|
| `store/store.go` | DB init, open, close, migration, reset |
| `store/repo.go` | Repo CRUD (upsert, list, find by slug/path) |
| `store/group.go` | Group CRUD (create, add, remove, list, show, delete) |
| `store/release.go` | Release CRUD (upsert, list, find by tag) |
| `store/amendment.go` | Amendment CRUD (insert, list all, list by branch) |
| `constants/constants_store.go` | DB path, table names, SQL statements, error messages |
| `constants/constants_amend.go` | Amendments table SQL, amend command constants |
| `model/record.go` | ScanRecord, Config, CloneResult, CloneSummary, ScanCache |
| `model/group.go` | Group, GroupRepo |
| `model/release.go` | ReleaseRecord |
| `model/amendment.go` | AmendmentRecord, AuthorInfo, CommitEntry |

---

## DB-First Lookup with JSON Fallback

Commands that resolve repos by slug (`pull`, `exec`, `status`) use a
two-tier lookup strategy:

1. **Try the database first.** Open `.gitmap/output/data/gitmap.db` and
   query the `Repos` table.
2. **Fall back to JSON.** If the database does not exist (no prior scan
   with DB support), load `.gitmap/output/gitmap.json` and match by
   repo name as before.

---

## Error Handling (Database)

| Scenario | Behavior |
|----------|----------|
| DB file cannot be created | Print error, exit 1 |
| Slug not found | `"No repo matches slug: %s"` |
| No database and DB-required flag | `"No database found. Run 'gitmap scan' first."` |
| Duplicate slug without qualifier | Interactive prompt (or error in non-TTY) |

---

## Constraints

- SQLite driver must be **CGo-free** (`modernc.org/sqlite`).
- All table/column names in **PascalCase**.
- All string literals in `constants` package.
- All files under 200 lines.
- All functions 8–15 lines.
- Positive conditions only (no negation).
- Blank line before `return`.
