# Phase 1 — v15 Database Rename Plan

**Created:** 2026-04-19 (Malaysia, UTC+8)
**Target version:** v3.1.0
**Spec:** https://github.com/alimtvnetwork/coding-guidelines-v15/blob/main/spec/04-database-conventions/01-naming-conventions.md
**Convention:** PascalCase + **singular** table names + `{TableName}Id` primary keys + FKs match referenced PK name.

---

## 1. Reliability Report (read before executing)

### Blast radius (measured)

| Surface | Count | Notes |
|---|---|---|
| `CREATE TABLE` statements | 26 | Across 18 `constants_*.go` files |
| `CREATE INDEX` statements | 1 | Will become 26+ under v15 |
| SQL string references to renamed objects | 281 | grep'd via plural table names |
| Files referencing renamed objects | 57 | constants/, store/, cmd/, model/ |
| `rows.Scan(...)` call sites | 52 | All in store/ |
| Migration helpers to update | 6 | `migrateLegacyIDs`, `migrateSourceColumn`, etc. |

### Risks

1. **Silent runtime breakage** — Go compiles even when SQL refers to a column that no longer exists. `no such column: Id` only surfaces when the command actually runs.
2. **Data-loss on migration** — SQLite has no native `RENAME COLUMN` for PK. Must `CREATE ... ; INSERT SELECT ; DROP ; ALTER RENAME` per table, with FK enforcement temporarily off. 26 such operations. One typo per table = lost rows.
3. **FK cascade ordering** — `GroupRepo` references both `Group` and `Repo`. Drops must be in reverse-FK order or migration aborts.
4. **Profile multiplicity** — `OpenProfile()` means each user may have N databases. Migration must be idempotent and run on every Open.
5. **Existing `migrateLegacyIDs`** — already does a UUID→INTEGER migration. Adding the v15 rename on top means two migrations run on the same user's DB.

### Mitigations baked into the plan

- **Detect-then-act**: every step checks `tableExists(oldName)` before touching anything; fresh installs skip the migration entirely.
- **Transaction per table**: `BEGIN IMMEDIATE; ...26 statements... ; COMMIT` so a partial failure rolls back.
- **`PRAGMA foreign_keys=OFF`** during migration, `ON` after — required by SQLite docs for table rebuild.
- **Verification step**: after migration, `SELECT COUNT(*)` on every new table must equal the old count. If mismatch → rollback + abort.
- **One table per turn during execution**: I'll do Repos→Repo first, verify build + run a smoke test, then move to the next. No big-bang.

---

## 2. Naming Map (authoritative)

### Tables (plural → singular)

| Old | New |
|---|---|
| `Repos` | `Repo` |
| `Groups` | `Group` |
| `GroupRepos` | `GroupRepo` |
| `Releases` | `Release` |
| `Aliases` | `Alias` |
| `Bookmarks` | `Bookmark` |
| `Amendments` | `Amendment` |
| `CommandHistory` | `CommandHistory` *(already singular)* |
| `CommitTemplates` | `CommitTemplate` |
| `Settings` | `Setting` |
| `SSHKeys` | `SshKey` *(also fixes `SSH`→`Ssh` per v15 abbreviations rule)* |
| `InstalledTools` | `InstalledTool` |
| `TempReleases` | `TempRelease` |
| `ZipGroups` | `ZipGroup` |
| `ZipGroupItems` | `ZipGroupItem` |
| `ProjectTypes` | `ProjectType` |
| `DetectedProjects` | `DetectedProject` |
| `GoProjectMetadata` | `GoProjectMetadata` *(uncountable, stays)* |
| `GoRunnableFiles` | `GoRunnableFile` |
| `CSharpProjectMetadata` | `CsharpProjectMetadata` *(also `CS`→`Cs`)* |
| `CSharpProjectFiles` | `CsharpProjectFile` |
| `CSharpKeyFiles` | `CsharpKeyFile` |
| `TaskType` | `TaskType` *(already singular)* |
| `PendingTask` | `PendingTask` *(already singular)* |
| `CompletedTask` | `CompletedTask` *(already singular)* |
| `RepoVersionHistory` | `RepoVersionHistory` *(uncountable, stays)* |

### Primary keys (`Id` → `{Table}Id`)

Every table's `Id INTEGER PRIMARY KEY AUTOINCREMENT` becomes `{NewTableName}Id`. Examples:

| Table | Old PK | New PK |
|---|---|---|
| `Repo` | `Id` | `RepoId` |
| `Group` | `Id` | `GroupId` |
| `Release` | `Id` | `ReleaseId` |
| `Alias` | `Id` | `AliasId` |
| `Bookmark` | `Id` | `BookmarkId` |
| `Amendment` | `Id` | `AmendmentId` |
| `Setting` | `Id` | `SettingId` |
| `SshKey` | `Id` | `SshKeyId` |
| `TempRelease` | `Id` | `TempReleaseId` |
| `ZipGroup` | `Id` | `ZipGroupId` |
| `ZipGroupItem` | `Id` | `ZipGroupItemId` |
| `ProjectType` | `Id` | `ProjectTypeId` |
| `DetectedProject` | `Id` | `DetectedProjectId` |
| `RepoVersionHistory` | `Id` | `RepoVersionHistoryId` |
| `CommandHistory` | `Id` | `CommandHistoryId` |
| `InstalledTool` | `Id` | `InstalledToolId` |
| `CommitTemplate` | `Id` | `CommitTemplateId` |
| `PendingTask` | `Id` | `PendingTaskId` |
| `CompletedTask` | `Id` | `CompletedTaskId` |
| `TaskType` | `Id` | `TaskTypeId` |
| `GoProjectMetadata` | `Id` | `GoProjectMetadataId` |
| `GoRunnableFile` | `Id` | `GoRunnableFileId` |
| `CsharpProjectMetadata` | `Id` | `CsharpProjectMetadataId` |
| `CsharpProjectFile` | `Id` | `CsharpProjectFileId` |
| `CsharpKeyFile` | `Id` | `CsharpKeyFileId` |

`GroupRepo` keeps its composite `(GroupId, RepoId)` PK — already v15-compliant (FKs match referenced PKs).

### Indexes (rename to `Idx{Table}_{Column}`)

| Old | New |
|---|---|
| `idx_Repos_AbsolutePath` | `IdxRepo_AbsolutePath` |

Plus we'll add explicit indexes per v15 for every FK column and unique business key (e.g., `IdxRelease_Tag`, `IdxAlias_Alias`, `IdxGroup_Name`, `IdxRepoVersionHistory_RepoId`, `IdxBookmark_Slug`, etc.) — measured per-table during execution.

### Boolean columns (audit for `Is`/`Has` prefix)

Existing booleans to verify/rename per v15 Rule 1:

| Table | Old | New (if needed) |
|---|---|---|
| `Release.Draft` | `Draft INTEGER` | `IsDraft INTEGER NOT NULL DEFAULT 0` |
| `Release.PreRelease` | `PreRelease INTEGER` | `IsPreRelease INTEGER NOT NULL DEFAULT 0` |
| `Release.IsLatest` | `IsLatest INTEGER` | already correct |

Other tables: scanned, no other bare-adjective bools found.

### Abbreviation fixes (v15 Rule: `Url`/`Api`/`Id`/`Ssh`, never `URL`/`API`/`ID`/`SSH`)

Already mostly compliant in column names (`HttpsUrl`, `SshUrl`). Table-level fix: `SSHKeys` → `SshKey`, `CSharp*` → `Csharp*`.

### Go struct field-tag impact

Most `model/*.go` structs use Go-idiomatic `ID int64` with no `db:` tag. SQL Scan binds positionally, so renaming the SQL column doesn't require Go field renames — only the SQL strings change. Confirmed by inspecting `scanOneRow` in `store/repo.go`.

---

## 3. Execution Phases (one per "next")

### Phase 1.1 — Repo + GroupRepo + index (small, safe pilot)
- Rewrite `constants_store.go`: `Repos→Repo`, `Id→RepoId`, `GroupRepos→GroupRepo`, index → `IdxRepo_AbsolutePath`.
- Update every SQL string in that file to use `RepoId`.
- Update `store/repo.go` Scan calls (positional, no change needed) + any inline `WHERE Id =` → `WHERE RepoId =`.
- Add `migrateV15Repos()` to `store/migrations.go`: detect old `Repos` table, table-rebuild dance, copy data, verify count, drop old.
- Run `go build ./...` and `go test ./...` to confirm.
- **Stop. Wait for "next".**

### Phase 1.2 — Group + Release + Alias + Bookmark
- Same procedure for these 4 tables in one batch.
- Build + test.
- **Stop. Wait for "next".**

### Phase 1.3 — Amendment + CommitTemplate + Setting + SshKey + InstalledTool + TempRelease
- 6 tables.
- Build + test.
- **Stop.**

### Phase 1.4 — ZipGroup family + Project family + Task family + History tables
- Remaining 11 tables.
- Build + test.
- **Stop.**

### Phase 1.5 — Boolean prefix fixes (`Draft`→`IsDraft`, `PreRelease`→`IsPreRelease`)
- Smaller cleanup, separate to keep blast radius isolated.

### Phase 1.6 — Spec, ERD, memory, version bump
- Update `spec/12-consolidated-guidelines/11-database.md` to v15 (singular + `{Table}Id`).
- Regenerate both `.mmd` ERDs.
- Update `mem://index.md` core: "Database schema follows v15 PascalCase + singular tables + `{Table}Id` PKs."
- Bump `internal/version` to `v3.1.0`.
- Add CHANGELOG entry.

---

## 4. Migration Code Skeleton (to implement in Phase 1.1)

```go
// migrateV15TableRebuild is the generic 4-step rebuild for renaming a table
// AND its primary key column atomically while preserving data.
func (db *DB) migrateV15TableRebuild(oldTable, newTable, newCreateSQL,
    columnList, oldPkCol, newPkCol string) error {

    if !db.tableExists(oldTable) {
        return nil // fresh install or already migrated
    }
    if db.tableExists(newTable) {
        return nil // migration already ran
    }

    tx, err := db.conn.Begin()
    if err != nil { return err }
    defer tx.Rollback()

    if _, err := tx.Exec("PRAGMA foreign_keys=OFF"); err != nil { return err }

    // 1. Create new table with v15 schema.
    if _, err := tx.Exec(newCreateSQL); err != nil { return err }

    // 2. Copy data, mapping old `Id` → new `{Table}Id`.
    copySQL := fmt.Sprintf(
        "INSERT INTO %s (%s) SELECT %s FROM %s",
        newTable, columnList, columnList, oldTable,
    )
    if _, err := tx.Exec(copySQL); err != nil { return err }

    // 3. Verify counts match.
    var oldN, newN int
    tx.QueryRow("SELECT COUNT(*) FROM " + oldTable).Scan(&oldN)
    tx.QueryRow("SELECT COUNT(*) FROM " + newTable).Scan(&newN)
    if oldN != newN {
        return fmt.Errorf("migration count mismatch on %s: old=%d new=%d",
            oldTable, oldN, newN)
    }

    // 4. Drop old table.
    if _, err := tx.Exec("DROP TABLE " + oldTable); err != nil { return err }

    if _, err := tx.Exec("PRAGMA foreign_keys=ON"); err != nil { return err }

    return tx.Commit()
}
```

---

## 5. What I am NOT doing in Phase 1

- Adding `ScanFolder` table (deferred to Phase 2 per user instruction).
- Adding `VersionProbe` table (deferred to Phase 2).
- Implementing `gitmap find-next` command (deferred to Phase 2).
- Implementing `gitmap pull` parallel runner (deferred to Phase 3).
- Implementing `gitmap cn next all` (deferred to Phase 3).

---

## 6. Open questions to confirm before Phase 1.1 executes

1. The `Releases` table has `Draft` and `PreRelease` columns — these are bare-adjective booleans. v15 Rule 1 mandates `Is`/`Has` prefix. Rename to `IsDraft`/`IsPreRelease` in the same migration, or split into Phase 1.5? **Plan: Phase 1.5 to keep blast radius isolated.**
2. `SSHKeys` table — v15 abbreviation rule says only first letter capitalized: `Ssh`, never `SSH`. Rename to `SshKey`? **Plan: yes, in Phase 1.3.**
3. `CSharp*` tables — same rule, rename to `Csharp*`? **Plan: yes, in Phase 1.4.** But this is a real readability hit; alternative is to declare `CSharp` an established proper noun and exempt. **Need user call.**
```
