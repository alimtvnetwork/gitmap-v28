# Database Schema — UpdateChecker & UpdateStatus

> **Version:** 1.0.0
> **Parent:** [00-overview.md](./00-overview.md)

---

## Purpose

Defines the SQLite tables that persist update-check results. PascalCase
naming throughout — table names, column names, enum values. Default
backend is SQLite via the project ORM; when no database exists the same
data is written to JSON (see [09-json-fallback-store.md](./09-json-fallback-store.md)).

---

## 1. Table: `UpdateChecker`

Singleton-style row table — typically one row per CLI install, updated
in place on every successful check. (Implementations MAY append rows
for audit history; the service always reads `WHERE Id = (SELECT MAX(Id))`.)

| FieldName            | Type     | Null | Default | Notes |
|----------------------|----------|------|---------|-------|
| `UpdateCheckerId`    | INTEGER  | no   | AUTOINCREMENT | Primary key |
| `CurrentVersion`     | TEXT     | no   | —       | Installed CLI version, e.g., `V1.5.0` |
| `LatestVersion`      | TEXT     | yes  | NULL    | Highest version found via discovery |
| `HasUpdate`          | BOOLEAN  | no   | 0       | True if `LatestVersion > CurrentVersion` |
| `UpdateStatusId`     | TINYINT  | no   | 1       | FK → `UpdateStatus.UpdateStatusId` |
| `OwnerKind`          | TEXT     | yes  | NULL    | `User` or `Organization` |
| `Owner`              | TEXT     | yes  | NULL    | GitHub login |
| `CurrentRepo`        | TEXT     | yes  | NULL    | e.g., `repo-v15` |
| `LatestReleaseUrl`   | TEXT     | yes  | NULL    | GitHub release page URL |
| `WindowsInstallUrl`  | TEXT     | yes  | NULL    | Pinned installer URL (Windows) |
| `WindowsInstallCmd`  | TEXT     | yes  | NULL    | One-liner (Windows) |
| `UnixInstallUrl`     | TEXT     | yes  | NULL    | Pinned installer URL (Unix) |
| `UnixInstallCmd`     | TEXT     | yes  | NULL    | One-liner (Unix) |
| `Checksum`           | TEXT     | yes  | NULL    | `Sha256:Hex` of the installer |
| `PublishedAt`        | DATETIME | yes  | NULL    | Release publish time (UTC) |
| `MinSupportedFrom`   | TEXT     | yes  | NULL    | Lowest version that can directly upgrade |
| `NewRepoUrl`         | TEXT     | yes  | NULL    | Migration target repo, when set |
| `Notes`              | TEXT     | yes  | NULL    | Release-notes summary |
| `RawJson`            | TEXT     | yes  | NULL    | Full combined JSON (see `03-combined-json.md`) |
| `LastCheckedAt`      | DATETIME | yes  | NULL    | UTC timestamp of last successful check |
| `NextCheckDueAt`     | DATETIME | yes  | NULL    | `LastCheckedAt + CheckIntervalHours` |
| `CheckIntervalHours` | SMALLINT | no   | 12      | Configured interval |
| `ErrorMessage`       | TEXT     | yes  | NULL    | Last error message |
| `ErrorAt`            | DATETIME | yes  | NULL    | UTC timestamp of last error |
| `Description`        | TEXT     | yes  | NULL    | Schema Rule 10 — entity description |
| `CreatedAt`          | DATETIME | no   | CURRENT_TIMESTAMP | Row creation |
| `UpdatedAt`          | DATETIME | no   | CURRENT_TIMESTAMP | Row update |

> **Schema Rule 10/11/12 compliance:** `UpdateChecker` is a transactional
> table; `Description`, plus `Notes` (already present) and `ErrorMessage`
> satisfy the free-text requirement. All free-text columns are nullable
> with no `DEFAULT`.

---

## 2. Table: `UpdateStatus` (Lookup)

| FieldName        | Type    | Null | Notes |
|------------------|---------|------|-------|
| `UpdateStatusId` | TINYINT | no   | Primary key (smallest int — cardinality ≤ 8) |
| `Name`           | TEXT    | no   | Enum value (PascalCase) |
| `Label`          | TEXT    | no   | Human-readable label |
| `Description`    | TEXT    | yes  | Schema Rule 10 |

### Seed rows

| UpdateStatusId | Name             | Label                       |
|----------------|------------------|-----------------------------|
| 1              | `UpToDate`       | Up to date                  |
| 2              | `UpdateFound`    | Update available            |
| 3              | `UpdateApplied`  | Update applied              |
| 4              | `Failed`         | Last check failed           |
| 5              | `Migrated`       | Project moved to a new repo |

> The enum is mirrored in code as `UpdateStatusEnum` (PascalCase, see
> Enum Standards) with strict `ParseUpdateStatus()` parsing.

---

## 3. Relationships

```
UpdateChecker.UpdateStatusId  ──►  UpdateStatus.UpdateStatusId   (Many-to-One)
```

No other foreign keys. The table is intentionally self-contained so the
update subsystem has zero coupling to user/session tables.

---

## 4. DDL (SQLite)

```sql
CREATE TABLE IF NOT EXISTS UpdateStatus (
  UpdateStatusId TINYINT PRIMARY KEY,
  Name           TEXT    NOT NULL UNIQUE,
  Label          TEXT    NOT NULL,
  Description    TEXT    NULL
);

CREATE TABLE IF NOT EXISTS UpdateChecker (
  UpdateCheckerId    INTEGER PRIMARY KEY AUTOINCREMENT,
  CurrentVersion     TEXT     NOT NULL,
  LatestVersion      TEXT     NULL,
  HasUpdate          BOOLEAN  NOT NULL DEFAULT 0,
  UpdateStatusId     TINYINT  NOT NULL DEFAULT 1
                     REFERENCES UpdateStatus(UpdateStatusId),
  OwnerKind          TEXT     NULL,
  Owner              TEXT     NULL,
  CurrentRepo        TEXT     NULL,
  LatestReleaseUrl   TEXT     NULL,
  WindowsInstallUrl  TEXT     NULL,
  WindowsInstallCmd  TEXT     NULL,
  UnixInstallUrl     TEXT     NULL,
  UnixInstallCmd     TEXT     NULL,
  Checksum           TEXT     NULL,
  PublishedAt        DATETIME NULL,
  MinSupportedFrom   TEXT     NULL,
  NewRepoUrl         TEXT     NULL,
  Notes              TEXT     NULL,
  RawJson            TEXT     NULL,
  LastCheckedAt      DATETIME NULL,
  NextCheckDueAt     DATETIME NULL,
  CheckIntervalHours SMALLINT NOT NULL DEFAULT 12,
  ErrorMessage       TEXT     NULL,
  ErrorAt            DATETIME NULL,
  Description        TEXT     NULL,
  CreatedAt          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UpdatedAt          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS IxUpdateCheckerNextCheckDueAt
  ON UpdateChecker(NextCheckDueAt);
```

---

## 5. Migration Note

The DDL is idempotent. The service runs it on first call to
`UpdateCheckerService.Initialize()`. Lookup rows are upserted by
`UpdateStatusId`.

---

*Database Schema — v1.0.0 — 2026-04-20*
