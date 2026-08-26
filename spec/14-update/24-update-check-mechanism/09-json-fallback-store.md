# JSON Fallback Store

> **Version:** 1.0.0
> **Parent:** [00-overview.md](./00-overview.md)

---

## Purpose

When the host CLI does not have a SQLite database (or any database)
configured, the `UpdateCheckerService` persists results to a single
JSON file with **the same schema** as the `UpdateChecker` table. The
service interface (`UpdateStore`) is identical — callers cannot tell
which backend is in use.

---

## 1. File Path

```
~/.<CliName>/data/UpdateChecker.json
```

Resolved per OS:

| OS | Path |
|----|------|
| Linux / macOS | `$HOME/.<CliName>/data/UpdateChecker.json` |
| Windows | `%USERPROFILE%\.<CliName>\data\UpdateChecker.json` |

The directory is created with mode `0700` if missing. The file is
created with mode `0600`. (CLI tools occasionally include credentials
in the same `.<CliName>/` tree — keep the whole tree user-private.)

---

## 2. Document Schema

A single JSON object whose keys mirror the `UpdateChecker` columns
(see [04-database-schema.md](./04-database-schema.md)) one-to-one,
PascalCase:

```json
{
  "UpdateCheckerId": 1,
  "CurrentVersion": "V1.5.0",
  "LatestVersion": "V1.7.0",
  "HasUpdate": true,
  "UpdateStatusId": 2,
  "UpdateStatusName": "UpdateFound",
  "OwnerKind": "User",
  "Owner": "AlifKarim",
  "CurrentRepo": "repo-v15",
  "LatestReleaseUrl": "https://github.com/AlifKarim/repo-v17/releases/tag/v1.7.0",
  "WindowsInstallUrl": "https://github.com/AlifKarim/repo-v17/releases/download/v1.7.0/Install.ps1",
  "WindowsInstallCmd": "iwr -useb <ScriptUrl> | iex",
  "UnixInstallUrl": "https://github.com/AlifKarim/repo-v17/releases/download/v1.7.0/Install.sh",
  "UnixInstallCmd": "curl -fsSL <ScriptUrl> | bash",
  "Checksum": "Sha256:Abc123...",
  "PublishedAt": "2026-04-15T10:00:00Z",
  "MinSupportedFrom": "V1.0.0",
  "NewRepoUrl": null,
  "Notes": "Bug fixes",
  "RawJson": "<combined JSON document, escaped string>",
  "LastCheckedAt": "2026-04-20T12:00:00Z",
  "NextCheckDueAt": "2026-04-21T00:00:00Z",
  "CheckIntervalHours": 12,
  "ErrorMessage": null,
  "ErrorAt": null,
  "Description": null,
  "CreatedAt": "2026-04-20T12:00:00Z",
  "UpdatedAt": "2026-04-20T12:00:00Z"
}
```

> `UpdateStatusName` is included alongside `UpdateStatusId` so the
> file is self-describing without the lookup table. The DB backend
> joins to `UpdateStatus`; the JSON backend embeds the name directly.

---

## 3. Atomic Write Protocol

The service MUST write atomically to prevent corruption when the user
hits Ctrl-C mid-write or two processes race:

```
1. Serialize payload → bytes
2. Write bytes → "<path>.tmp.<pid>"
3. fsync the tmp file
4. os.Rename(tmp, finalPath)        // POSIX rename = atomic
5. fsync the parent directory       // POSIX only; no-op on Windows
```

On Windows, step 4 uses `MoveFileEx(MOVEFILE_REPLACE_EXISTING |
MOVEFILE_WRITE_THROUGH)`.

If any step fails, the tmp file is deleted and the error is propagated
per [08 §3](./08-error-handling.md#3-subsystem-level-failure).

---

## 4. Concurrency

A POSIX advisory lock (`flock(LOCK_EX)`) on
`<path>.lock` (separate file, never read) serializes writers across
processes. On Windows, `LockFileEx(LOCKFILE_EXCLUSIVE_LOCK)` against
the same `.lock` file. Readers (`GetLastResult`) take a shared lock
(`LOCK_SH`).

If the lock cannot be acquired within 2 seconds, the writer logs a
warning and aborts — the next pre-hook will retry.

---

## 5. Migration to SQLite

If a project later adopts SQLite, the service MUST detect the
JSON file on first `Initialize()` and:

1. Run the DDL.
2. Insert the JSON contents as one row into `UpdateChecker`.
3. Rename the JSON to `UpdateChecker.json.migrated` (kept for one
   release cycle as a rollback safety net).

Migration is one-way. After migration, the JSON path is no longer
read.

---

## 6. Selecting the Backend

`UpdateCheckerService.Initialize()` chooses the backend in this order:

1. If a SQLite database is already configured by the host CLI →
   SQLite.
2. Else if the host CLI declares "no database" in `06-Seedable-Config`
   (`Storage.Backend = "JsonFile"`) → JSON.
3. Else → SQLite (auto-create at `~/.<CliName>/Data.sqlite`).

The choice is logged once per process. There is no automatic
SQLite → JSON downgrade.

---

*JSON Fallback Store — v1.0.0 — 2026-04-20*
