# Memory: features/v15-rename-progress
Updated: now (Malaysia, UTC+8)

**Phase 1 of the v15 database naming alignment is COMPLETE as of v3.5.0.** All 22 SQLite tables follow the strict v15 convention. No new commands or features were added — pure naming alignment.

## Spec
PascalCase + **singular** table names + `{TableName}Id` primary keys + FKs match referenced PK name + `IsX` boolean prefix + abbreviations as words. Per https://github.com/alimtvnetwork/coding-guidelines-v15/blob/main/spec/04-database-conventions/01-naming-conventions.md.

## Phase 1 sub-phase status

| Phase | Scope | Status |
|---|---|---|
| 1.1 | `Repos` → `Repo` (RepoId PK), `GroupRepos` → `GroupRepo`, index → `IdxRepo_AbsolutePath` | DONE (v3.1.0) |
| 1.2 | `Groups` → `Group` (GroupId), `Releases` → `Release` (ReleaseId), `Aliases` → `Alias` (AliasId), `Bookmarks` → `Bookmark` (BookmarkId) | DONE (v3.3.0) |
| 1.3 | `Amendments` → `Amendment`, `CommitTemplates` → `CommitTemplate`, `Settings` → `Setting`, `SSHKeys` → `SshKey`, `InstalledTools` → `InstalledTool`, `TempReleases` → `TempRelease` | DONE (v3.3.0) |
| 1.4 | ZipGroup family + Project family (incl. `CSharp` → `Csharp` strict v15) + Task family + History tables | DONE (v3.4.0) |
| 1.5 | Boolean prefix fixes (`Release.Draft` → `IsDraft`, `Release.PreRelease` → `IsPreRelease`) — full IsX consistency across DB, model, ReleaseMeta, Options, with backward-compat JSON aliases | DONE (v3.5.0) |
| 1.6 | Update spec/12, regenerate ERD, bump CHANGELOG entry, update mem://index core | DONE (v3.5.0) |

## Phase 1.6 — final wrap-up (this turn)
- **Regenerated `spec/01-app/gitmap-database-erd.mmd`** to v15 names: every table singular, every PK is `{TableName}Id`, FKs match (`GoRunnableFile.GoProjectMetadataId`, `CsharpProjectFile.CsharpProjectMetadataId`, etc.), `Release` has `IsDraft`/`IsPreRelease`/`IsLatest`, `SshKey`/`Csharp*` reflect abbreviation rule.
- **Updated `spec/12-consolidated-guidelines/11-database.md`** Schema Conventions table — singular names, `{TableName}Id` PKs, IsX boolean prefix row, reserved-word quoting row, abbreviation rules row, link to upstream v15 spec.
- **Added `## v3.5.0` CHANGELOG entry** covering Phases 1.1–1.5: full table rename list, column rename summary, migration safety contract (5 points), Go-side propagation, JSON backward-compat note for `.gitmap/release/*.json`, CLI `--draft` flag retained.
- **Updated `mem://index.md` core**: replaced old "PascalCase, INTEGER PRIMARY KEY AUTOINCREMENT" line with v15 spec one-liner; bumped current version to v3.5.0; added a top-of-list memory pointer to this file.

## Phase 1.2 + 1.3 — what changed in this turn

### New shared infrastructure
- **`gitmap/store/migrate_v15rebuild.go`** — generic `runV15Rebuild(spec)` helper using a `v15RebuildSpec` struct (OldTable, NewTable, NewCreateSQL, OldColumnList, NewColumnList, StartMsg, DoneMsg). Handles PRAGMA foreign_keys toggle, CREATE → INSERT SELECT → row-count parity check → DROP. Idempotent via `tableExists()` detect-then-act.

### Phase 1.2 migrator
- **`gitmap/store/migrate_v15phase2.go`** — `migrateV15Phase2()` runs four `runV15Rebuild` specs in dependency-safe order: Group → Release → Alias → Bookmark. Then calls `rebuildGroupRepoFK()` to rewrite the GroupRepo CREATE so its FK text references the new singular `"Group"(GroupId)` and `Repo(RepoId)` (SQLite stores FK clauses as text in sqlite_master and does NOT auto-update them when parent tables rename).

### Phase 1.3 migrator
- **`gitmap/store/migrate_v15phase3.go`** — `migrateV15Phase3()` runs six `runV15Rebuild` specs: Amendment, CommitTemplate, Setting (Key PK preserved), SshKey (also fixes SSH→Ssh abbreviation), InstalledTool, TempRelease. Setting uses `Key`/`Value` for both old and new column lists since there is no Id column.

### Pre-Phase-2 column patch
- **`store.go::preV15Phase2EnsureReleaseColumns()`** — runs ALTER on legacy `Releases` to ensure `Source` and `Notes` columns exist before the v15 rebuild SELECTs them by name. Protects very old installs that predate those columns.

### Constants files rewritten (singular + {Table}Id)
- `constants_store.go` — `SQLCreateGroup` (double-quoted reserved word), `SQLCreateRelease` (ReleaseId), `SQLCreateGroupRepo` references `"Group"(GroupId)` and `Repo(RepoId)`. `SQLAddSourceColumn` and `SQLAddNotesColumn` now target singular `Release`. New: `SQLImportInsertGroup`, `SQLDropGroup`, `SQLDropRelease`, `ErrV15Phase2Migration`, `ErrV15Phase3Migration`.
- `constants_alias.go` — `SQLCreateAlias` (AliasId), all DML/JOINs use singular `Alias`. New: `SQLDropAlias`, `LegacyTableAliases`.
- `constants_bookmark.go` — `SQLCreateBookmark` (BookmarkId), new `SQLImportInsertBookmark`, `SQLDropBookmark`.
- `constants_amend.go` — `SQLCreateAmendment` (AmendmentId), all DML uses singular.
- `constants_seo.go` — `SQLCreateCommitTemplate` (CommitTemplateId).
- `constants_settings.go` — `SQLCreateSetting` (Key PK preserved).
- `constants_ssh.go` — `SQLCreateSshKey` (SshKeyId; v15 abbreviation fix). Existing constant names like `SQLInsertSSHKey` retained for callsite stability — they now target the new `SshKey` table.
- `constants_installedtools.go` — `SQLCreateInstalledTool` (InstalledToolId).
- `constants_temprelease.go` — `SQLCreateTempRelease` (TempReleaseId). `SQLMigrateTRCommitSha` still targets legacy `TempReleases` (runs BEFORE the v15 rebuild copies the column).

### Wiring
- **`store.go::Migrate()`** order: `migrateLegacyIDs` → `migrateV15Repo` → `preV15Phase2EnsureReleaseColumns` → `migrateV15Phase2` → `migrateTRCommitSha` → `migrateV15Phase3` → standard CREATE TABLE pass (now uses all v15 singular names) → ALTER pass → seeds.
- **`store.go::Reset()`** — drop list now lists v15 singulars first, then legacy plurals. Each plural drop is a safe no-op when the table is absent. Covers all 10 newly renamed tables.
- **`store.go::migrateNotesColumn`** docstring updated (now says `Release` not `Releases`).
- **`gitmap/store/import.go`** — uses `constants.SQLImportInsertGroup` and `constants.SQLImportInsertBookmark` instead of inline SQL strings.

### Tests touched
- `gitmap/tests/constants_test/seo_constants_test.go` — `TableCommitTemplates` → `TableCommitTemplate`, `SQLCreateCommitTemplates` → `SQLCreateCommitTemplate`, expected column `Id` → `CommitTemplateId`.

### Version
- `gitmap/constants/constants.go` → `v3.3.0`.

## Migration safety contract
1. Detect-then-act on every legacy plural — fresh installs are no-ops.
2. PRAGMA foreign_keys=OFF for the duration of each table rebuild.
3. Row-count parity check between old and new on every rebuild — abort + rollback (via Go-side return) on mismatch.
4. Legacy plural names retained as `LegacyTable*` constants and listed in `Reset()` so cleanup works at any migration state.
5. SQLite-reserved word `Group` is double-quoted in every DDL/DML occurrence.

## Known limitations / debt for later phases
- The v15 rename touched 10 tables and their entire SQL surface this turn. A `go build` and `go test` on Windows are required to confirm zero compile or runtime regressions; the sandbox has no Go toolchain.
- Phase 1.4 will rebuild ZipGroup, ZipGroupItem, DetectedProject, ProjectType, GoProjectMetadata, GoRunnableFile, CSharpProjectMetadata, CSharpProjectFile, CSharpKeyFile (incl. `CSharp`→`Csharp` strict abbreviation per user's earlier choice), Task tables, and CommandHistory. RepoVersionHistory FK text still references `Repo(RepoId)` correctly (set in Phase 1.1) but its own table will keep that name (uncountable noun).
- Phase 1.5: `Release.Draft` → `IsDraft` and `Release.PreRelease` → `IsPreRelease` (column rename via ALTER … RENAME COLUMN, supported by SQLite ≥ 3.25).
- Phase 1.6: regenerate both ERDs, update spec/12-consolidated-guidelines/11-database.md, CHANGELOG entry, mem://index core update, version bump to v3.4.0 (or higher).

## What's still NOT done in Phase 1

### Phase 1.4 — IN PROGRESS, partial completion as of this turn
**Done this turn:**
- Bulk `sed s/CSharp/Csharp/g` across all 22 Go files that referenced `CSharp` (model, detector, cmd, store, constants). All Go identifiers (`CSharpProjectMetadata` → `CsharpProjectMetadata`, `r.CSharp` → `r.Csharp`, `ProjectKeyCSharp` → `ProjectKeyCsharp`, `CmdCSharpRepos` → `CmdCsharpRepos`, etc.) and SQL string literals inside Go backticks were renamed in one pass. JSON tags (`csharpMetadataId` lowercase) were unaffected.
- Rewrote 5 constants files with v15 singular tables + `{Table}Id` PKs:
  - `constants_project.go` — `TableProjectType`, `TableDetectedProject`, `TableGoRunnableFile`, `TableCsharp*` + `LegacyTable*` for migration detection (incl. pre-Csharp `"CSharp*"` legacy spellings).
  - `constants_project_sql.go` — `SQLCreateProjectType`/`DetectedProject`/`GoProjectMetadata`/`GoRunnableFile`/`CsharpProjectMetadata`/`CsharpProjectFile`/`CsharpKeyFile` all with `{Table}Id` PKs. Added `SQLDropCsharpProjectMetaLegacy` for cleanup of pre-Csharp tables.
  - `constants_zipgroup.go` — `TableZipGroup`/`TableZipGroupItem`, `SQLCreateZipGroup`/`SQLCreateZipGroupItem` with `ZipGroupId` PK; legacy plural drops kept.
  - `constants_history.go` — `CommandHistoryId` PK throughout.
  - `constants_version_history.go` — `RepoVersionHistoryId` PK.
  - `constants_pending_task.go` + `constants_pending_task_sql.go` — `TaskTypeId`/`PendingTaskId`/`CompletedTaskId` PKs throughout.

**NOT done this turn (next turn must finish before bumping version):**
1. **Store-side scan order updates** — every store/*.go file that does `rows.Scan(&r.ID, ...)` for these 14 tables needs to keep working with the new column order. Most are already correct because `{Table}Id` is still the first column, but verify: `store/zipgroup.go`, `store/project.go`, `store/csharpmetadata.go`, `store/gometadata.go`, `store/history.go`, `store/version_history.go`, `store/pendingtask.go`, `store/pendingtaskscan.go`. No code edits expected, just verification.
2. **Constant-name callsite fixes** — the constants file rewrites RENAMED some Go-side identifiers (`SQLCreateZipGroups` → `SQLCreateZipGroup`, `SQLCreateProjectTypes` → `SQLCreateProjectType`, `SQLCreateDetectedProjects` → `SQLCreateDetectedProject`, `SQLCreateGoRunnableFiles` → `SQLCreateGoRunnableFile`, `SQLCreateCsharpProjectFiles` → `SQLCreateCsharpProjectFile`, `SQLCreateCsharpKeyFiles` → `SQLCreateCsharpKeyFile`, `SQLCreateZipGroupItems` → `SQLCreateZipGroupItem`, `SQLDeleteStaleCsharpFiles`/`SQLDeleteStaleCsharpKeyFiles` already match, and `ErrCSharp*` → `ErrCsharp*` from sed). Callsites in `gitmap/store/store.go::Migrate()` and `Reset()` and `gitmap/store/migrateids.go::dropProjectTables()` reference the OLD names and will fail to compile. Must update.
3. **Migrator** — create `gitmap/store/migrate_v15phase4.go` with 14 `runV15Rebuild` specs (incl. CSharp-to-Csharp legacy detection: `OldTable: "CSharpProjectMetadata"` → `NewTable: "CsharpProjectMetadata"`). Wire into `store.go::Migrate()` between Phase 1.3 and the standard CREATE pass.
4. **Version bump** — `constants.go::Version` from `3.3.0` to `3.4.0` AFTER above completes and a clean compile is plausible.
5. **`migrateZipGroupItemPaths()`** in store.go — its constants `SQLMigrateZGI*` still target legacy plural `ZipGroupItems` (correct, these are pre-rename ALTERs that must run BEFORE the v15 rebuild copies the table — same pattern as `preV15Phase2EnsureReleaseColumns`).

### Phase 1.5, 1.6 (unchanged, still TODO)
- Phase 1.5: `Release.Draft` → `IsDraft`, `Release.PreRelease` → `IsPreRelease` (column rename).
- Phase 1.6: regenerate ERDs, update spec/12, CHANGELOG, mem://index core, version bump.

## Deferred to later phases
- ScanFolder table (Phase 2).
- VersionProbe table (Phase 2).
- `gitmap find-next` command (Phase 2).
- `gitmap pull` parallel runner (Phase 3).
- `gitmap cn next all` bulk update (Phase 3).
