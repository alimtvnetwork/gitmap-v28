# ResponseKeyType — Case Inventory & Usage Map

> **Enum**: `RiseupAsia\Enums\ResponseKeyType`  
> **File**: `includes/Enums/ResponseKeyType.php`  
> **Go Mirror**: `backend/internal/enums/response_key/variant.go`  
> **TS Mirror**: `src/lib/constants.ts` → `ResponseKeyType`  
> **As of**: v2.1.0 (2026-02-26)  
> **Total cases**: 176  
> **Total usages**: ~3,500+ across 100+ files

---

## Value Casing Convention

All values use **PascalCase** (e.g., `'Success'`, `'SnapshotId'`, `'DeletedByPolicy'`).
Single-word keys are also PascalCase (e.g., `'Success'`, `'Error'`, `'Rows'`, `'Scope'`).

---

## Envelope Keys

Used in nearly every REST response and internal result array.

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `Success` | `Success` | ~56 files | All REST handlers, snapshot providers, helpers, cleaner, orchestrator |
| `Error` | `Error` | ~60 files | All error payloads, log contexts, snapshot CRUD, sync, agents |
| `Message` | `Message` | ~15 files | REST error responses, envelope builders, route traits |
| `Data` | `Data` | ~10 files | SyncManifestTrait, REST response wrappers, status payloads |
| `Code` | `Code` | ~12 files | Error responses with typed codes (SnapshotErrorType, WpErrorCodeType) |
| `Valid` | `Valid` | ~8 files | Import validation, manifest validation, agent validation |
| `Errors` | `Errors` | ~25 files | Batch results, cleaner phases, worker jobs, restore engine |
| `Cached` | `Cached` | ~5 files | FileCache manifest, sync manifest, export handler |
| `Phase` | `Phase` | ~10 files | Snapshot lifecycle logging (initiated, streaming, complete) |
| `Reason` | `Reason` | ~5 files | Retention deletion details, cleaner audit |

---

## Domain Collection Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `Total` | `Total` | ~8 files | Pagination (agents, snapshots, logs), manifest stats |
| `Agents` | `Agents` | ~4 files | AgentCrudReadTrait, agent list responses |
| `Actions` | `Actions` | ~3 files | Action log list responses |
| `Logs` | `Logs` | ~4 files | Log list responses, diagnostics |
| `Snapshots` | `Snapshots` | ~4 files | Snapshot list responses, UpdraftCrudTrait |
| `Sql` | `Sql` | ~3 files | Query builder results, database search |
| `Params` | `Params` | ~3 files | Query builder parameter arrays |
| `Sets` | `Sets` | ~2 files | Batch set operations |
| `Plugins` | `Plugins` | ~5 files | Orchestrator plugin archiving, import execution |
| `Tables` | `Tables` | ~15 files | Snapshot CRUD, restore, export, worker batches |
| `Settings` | `Settings` | ~5 files | SnapshotSettingsHandlerTrait, DetectorSettingsTrait, settings responses |
| `Providers` | `Providers` | ~4 files | DetectorProviderTrait, provider list responses |
| `Dependencies` | `Dependencies` | ~3 files | AnalyzerQueryTrait, table dependency analysis |

---

## File & Size Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `Rows` | `Rows` | ~12 files | Table info, restore results, batch progress, worker exports |
| `Bytes` | `Bytes` | ~4 files | Storage calculations, sync payloads |
| `Size` | `Size` | ~15 files | ZIP sizes, file entries, plugin archiving, export results |
| `FileSize` | `FileSize` | ~8 files | Snapshot records, incremental exports, manifest entries |
| `Path` | `Path` | ~12 files | Log contexts, file manifests, REST responses |
| `Filename` | `Filename` | ~15 files | Snapshot CRUD, export/import, cleaner, manifest |
| `Checksum` | `Checksum` | ~4 files | Incremental exports, file integrity, sync |
| `Duration` | `Duration` | ~12 files | All timed operations (backup, restore, cleanup, sync) |
| `Count` | `Count` | ~5 files | Plugin archiving, orchestrator results |
| `Files` | `Files` | ~8 files | FileCache manifest, sync manifest, plugin list |
| `Directory` | `Directory` | ~6 files | Snapshot creation results, backup responses |
| `Scope` | `Scope` | ~8 files | Snapshot CRUD, export manifest, import records |
| `Exported` | `Exported` | ~4 files | Worker batch progress, export results |
| `Entry` | `Entry` | ~3 files | Plugin archive entries, orchestrator |
| `Computed` | `Computed` | ~3 files | FileCache manifest stats, sync cache stats |
| `Removed` | `Removed` | ~5 files | FileCache pruning, cleaner orphan results |

---

## Pagination Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `Limit` | `Limit` | ~4 files | Paginated list endpoints, query builders |
| `Offset` | `Offset` | ~4 files | Paginated list endpoints, query builders |

---

## Domain Entity Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `Posts` | `Posts` | ~3 files | Post list responses, content sync |
| `Categories` | `Categories` | ~2 files | Category list responses |
| `Category` | `Category` | ~2 files | Single category payloads |
| `Export` | `Export` | ~3 files | Export operation results |
| `Incrementals` | `Incrementals` | ~3 files | Incremental snapshot list results |
| `TotalSize` | `TotalSize` | ~3 files | Storage summary responses, ErrorLogHandlerTrait |
| `Applied` | `Applied` | ~2 files | Incremental apply results |
| `Folder` | `Folder` | ~3 files | Snapshot directory references |

---

## Snapshot-Domain Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `SnapshotId` | `SnapshotId` | ~12 files | All snapshot operations, export, import, restore, audit |
| `Sequence` | `Sequence` | ~6 files | Incremental backups, export manifests, registration |
| `FolderName` | `FolderName` | ~5 files | Incremental backup directories, registration |
| `TablesChanged` | `TablesChanged` | ~4 files | Incremental registration, export results |
| `TotalRows` | `TotalRows` | ~15 files | Snapshot records, worker progress, restore, import |
| `TotalNewRows` | `TotalNewRows` | ~4 files | Incremental registration, export results |
| `ZipPath` | `ZipPath` | ~3 files | ZIP file path references |
| `ZipSize` | `ZipSize` | ~4 files | Backup exec responses, export results |
| `BackupId` | `BackupId` | ~3 files | Pre-restore backup references |
| `ZipFailed` | `ZipFailed` | ~3 files | Snapshot creation error flags |
| `SkipAudit` | `SkipAudit` | ~4 files | Scheduler cron results, no-op cleanup |
| `TablesRestored` | `TablesRestored` | ~4 files | Restore engine results, audit logging |

---

## Cleanup-Pipeline Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `DeletedByPolicy` | `DeletedByPolicy` | ~3 files | SnapshotCleaner::runCleanup, AdminAjaxSnapshotTrait, SchedulerExecutorTrait |
| `DeletedOrphans` | `DeletedOrphans` | ~3 files | SnapshotCleaner::runCleanup, AdminAjaxSnapshotTrait, SchedulerExecutorTrait |
| `DeletedFailed` | `DeletedFailed` | ~3 files | SnapshotCleaner::runCleanup, AdminAjaxSnapshotTrait, SchedulerExecutorTrait |
| `SpaceFreedBytes` | `SpaceFreedBytes` | ~3 files | SnapshotCleaner::runCleanup, AdminAjaxSnapshotTrait, SchedulerExecutorTrait |
| `Retention` | `Retention` | ~3 files | SnapshotCleaner, SnapshotBackupOpsTrait, cleanup settings |
| `Orphans` | `Orphans` | ~3 files | CleanerOrphanTrait, SnapshotCleaner |
| `Stuck` | `Stuck` | ~2 files | SnapshotCleaner stuck-snapshot detection |
| `DryRun` | `DryRun` | ~2 files | SnapshotCleaner dry-run mode flag |
| `BytesFreed` | `BytesFreed` | ~2 files | CleanerStorageTrait, cleanup audit |
| `Deleted` | `Deleted` | ~3 files | PluginLifecycleDeleteTrait, cleanup results |
| `Cleaned` | `Cleaned` | ~2 files | CleanerOrphanTrait orphan cleanup results |

---

## Plugin Lifecycle Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `Activated` | `Activated` | ~3 files | UploadInstallActivateTrait, plugin lifecycle responses |
| `PluginSlug` | `PluginSlug` | ~4 files | UploadInstallActivateTrait, PluginLifecycleDeleteTrait, OrchestratorPluginTrait |
| `IsUpdate` | `IsUpdate` | ~3 files | UploadInstallActivateTrait, upload responses |
| `IsSelfUpdate` | `IsSelfUpdate` | ~2 files | UploadInstallActivateTrait self-update detection |
| `PluginVersion` | `PluginVersion` | ~3 files | UploadInstallActivateTrait, OrchestratorPluginTrait |
| `ActivationError` | `ActivationError` | ~2 files | UploadInstallActivateTrait error capture |
| `Inventory` | `Inventory` | ~2 files | RestoreValidationTrait, import inventory |
| `PluginFile` | `PluginFile` | ~3 files | UploadPipelineTrait, PluginLifecycleHelpersTrait, plugin resolution |

---

## General-Purpose Entity Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `Slug` | `Slug` | ~5 files | OrchestratorPluginTrait, DetectorSettingsTrait, plugin arrays |
| `Title` | `Title` | ~4 files | SchedulerTriggerTrait, snapshot creation, import |
| `Type` | `Type` | ~4 files | Snapshot type indicators, plugin type fields |
| `Action` | `Action` | ~3 files | Transaction logging, action audit |
| `Status` | `Status` | ~8 files | Snapshot records, agent status, transaction status |
| `Percent` | `Percent` | ~3 files | Progress tracking, batch completion |
| `Plugin` | `Plugin` | ~3 files | StatusPayloadTrait, plugin info responses |

---

## Log/Diagnostic Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `ErrorLog` | `ErrorLog` | ~2 files | ErrorLogHandlerTrait log reading |
| `FullLog` | `FullLog` | ~2 files | ErrorLogHandlerTrait full content |
| `StacktraceLog` | `StacktraceLog` | ~2 files | ErrorLogHandlerTrait stacktrace extraction |
| `Exists` | `Exists` | ~2 files | ErrorLogHandlerTrait file existence check |
| `Content` | `Content` | ~3 files | ErrorLogHandlerTrait, log content payloads |
| `Truncated` | `Truncated` | ~2 files | ErrorLogHandlerTrait large-file truncation flag |
| `Lines` | `Lines` | ~2 files | ErrorLogHandlerTrait line-based reading |
| `TotalLines` | `TotalLines` | ~2 files | ErrorLogHandlerTrait line count |

---

## Internal/Domain-Specific Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `Ids` | `Ids` | ~2 files | CleanerOrphanTrait batch ID arrays |
| `TotalSnapshots` | `TotalSnapshots` | ~2 files | CleanerStorageTrait storage summary |
| `TotalSizeBytes` | `TotalSizeBytes` | ~2 files | CleanerStorageTrait storage metrics |
| `TempFile` | `TempFile` | ~2 files | UploadInstallExtractTrait temp path tracking |
| `Stmt` | `Stmt` | ~2 files | IncrementalExportTrait prepared statement key |
| `Columns` | `Columns` | ~3 files | IncrementalExportTrait column metadata |

---

## Temporal Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `CreatedAt` | `CreatedAt` | ~6 files | Snapshot records, audit logs, manifest metadata |
| `UpdatedAt` | `UpdatedAt` | ~4 files | Record update timestamps, sync metadata |
| `Timestamp` | `Timestamp` | ~4 files | StatusOpsTrait opcache result, DateHelper payloads |

---

## Analysis & Dependency Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `ParentTable` | `ParentTable` | ~2 files | Table dependency analysis, seed ordering |
| `ChildTable` | `ChildTable` | ~2 files | Table dependency analysis |
| `FkColumn` | `FkColumn` | ~2 files | Foreign key column references |
| `RefColumn` | `RefColumn` | ~2 files | Referenced column in FK relationships |
| `SeedOrder` | `SeedOrder` | ~2 files | Table insertion order for seeding |
| `TableCount` | `TableCount` | ~3 files | Snapshot table counts |
| `DepCount` | `DepCount` | ~2 files | Dependency count per table |
| `NewRows` | `NewRows` | ~3 files | Incremental export new row counts |
| `PluginDetails` | `PluginDetails` | ~3 files | Orchestrator plugin metadata |
| `IncludedIds` | `IncludedIds` | ~2 files | Selective export ID lists |
| `IncrementalCount` | `IncrementalCount` | ~3 files | Child incremental count per full snapshot |

---

## Detection & Provider Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `DetectionMethod` | `DetectionMethod` | ~2 files | Snapshot provider auto-detection |
| `SqliteVersion` | `SqliteVersion` | ~2 files | SQLite runtime version reporting |
| `IsCore` | `IsCore` | ~3 files | SnapshotProviderWpReset, WordPress core table flag |

---

## Scheduler Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `ScheduleEnabled` | `ScheduleEnabled` | ~3 files | SchedulerCronTrait, settings responses |
| `NextScheduledSnapshot` | `NextScheduledSnapshot` | ~2 files | Scheduler status responses |
| `NextCleanup` | `NextCleanup` | ~2 files | Cleanup schedule responses |
| `RetentionType` | `RetentionType` | ~3 files | Retention policy settings (days/count/none) |
| `RetentionDays` | `RetentionDays` | ~2 files | Days-based retention value |
| `RetentionCount` | `RetentionCount` | ~2 files | Count-based retention value |
| `SnapshotType` | `SnapshotType` | ~4 files | Full/incremental type indicator |

---

## Error Enrichment Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `ErrorCategory` | `ErrorCategory` | ~2 files | Categorized error reporting |
| `LogHint` | `LogHint` | ~2 files | Contextual hint for log analysis |

---

## Sync Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `FilesUpdated` | `FilesUpdated` | ~3 files | SyncManifestTrait, sync results |
| `FilesDeleted` | `FilesDeleted` | ~3 files | SyncManifestTrait, sync results |
| `FilesIgnored` | `FilesIgnored` | ~2 files | Sync filter results |
| `IgnoredFiles` | `IgnoredFiles` | ~2 files | Sync ignored file list |

---

## Export & Plugin Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `PluginZip` | `PluginZip` | ~2 files | Plugin ZIP export responses |
| `ResolvedUrl` | `ResolvedUrl` | ~2 files | URL resolution for remote resources |
| `TraceLines` | `TraceLines` | ~2 files | Error trace line extraction |

---

## Snapshot Progress & Worker Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `CompletedAt` | `CompletedAt` | ~3 files | Snapshot completion timestamps |
| `ExportedAt` | `ExportedAt` | ~2 files | Export completion timestamps |
| `Format` | `Format` | ~2 files | Snapshot format identifier |
| `FormatVersion` | `FormatVersion` | ~2 files | Snapshot format version metadata |
| `JobId` | `JobId` | ~3 files | Worker job identifiers |
| `TotalTables` | `TotalTables` | ~4 files | Export total table counts |
| `TablesExported` | `TablesExported` | ~3 files | Export progress tracking |
| `PoolSize` | `PoolSize` | ~2 files | Worker pool configuration |
| `TotalBatches` | `TotalBatches` | ~2 files | Batch processing totals |
| `CurrentBatch` | `CurrentBatch` | ~2 files | Batch processing progress |
| `TableProgress` | `TableProgress` | ~3 files | Per-table export progress |
| `IncrementalsApplied` | `IncrementalsApplied` | ~2 files | Incremental restore applied count |
| `SkippedMaster` | `SkippedMaster` | ~2 files | Skipped master snapshot flag |
| `ExportedTables` | `ExportedTables` | ~3 files | List of exported table names |
| `SnapshotDir` | `SnapshotDir` | ~3 files | Snapshot directory path |
| `DirName` | `DirName` | ~2 files | Directory name references |
| `RowCount` | `RowCount` | ~3 files | Per-table row counts |

---

## Cron & Audit Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `TriggeredBy` | `TriggeredBy` | ~3 files | SchedulerCronTrait, audit logging |
| `AuditData` | `AuditData` | ~2 files | Audit data payloads |
| `LogDataKey` | `LogData` | ~3 files | Log data payloads |

---

## Manifest & Import Metadata Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `OriginalId` | `OriginalId` | ~2 files | Import source record ID |
| `OriginalCreatedAt` | `OriginalCreatedAt` | ~2 files | Import source timestamp |
| `SourceSite` | `SourceSite` | ~2 files | Import source site identifier |
| `OriginalTitle` | `OriginalTitle` | ~2 files | Import source title |
| `OriginalType` | `OriginalType` | ~2 files | Import source snapshot type |
| `WpVersion` | `WpVersion` | ~2 files | WordPress version in manifest |
| `PhpVersion` | `PhpVersion` | ~2 files | PHP version in manifest |
| `MysqlVersion` | `MysqlVersion` | ~2 files | MySQL version in manifest |
| `SiteUrl` | `SiteUrl` | ~3 files | Site URL in manifest |
| `DbPrefix` | `DbPrefix` | ~2 files | Database table prefix |
| `PluginCount` | `PluginCount` | ~2 files | Plugin count in manifest |
| `DurationMs` | `DurationMs` | ~3 files | Duration in milliseconds |
| `TableCounts` | `TableCounts` | ~2 files | Per-table row count map |

---

## Sync Manifest Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `DownloadUrl` | `DownloadUrl` | ~3 files | SyncManifestTrait, file download URLs |
| `FileCount` | `FileCount` | ~2 files | Manifest file count |
| `GeneratedAt` | `GeneratedAt` | ~2 files | Manifest generation timestamp |
| `CacheStats` | `CacheStats` | ~2 files | FileCache statistics |
| `FromCache` | `FromCache` | ~2 files | Cache hit indicator |

---

## Statistics Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `TotalTransactions` | `TotalTransactions` | ~2 files | Log statistics responses |
| `ByAction` | `ByAction` | ~2 files | Stats grouped by action type |
| `ByStatus` | `ByStatus` | ~2 files | Stats grouped by status |
| `Last24h` | `Last24h` | ~2 files | Recent activity window |

---

## Backup Option Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `IncludePlugins` | `IncludePlugins` | ~3 files | Snapshot options, SchedulerTriggerTrait |
| `PluginSelection` | `PluginSelection` | ~2 files | Selective plugin backup list |
| `Compression` | `Compression` | ~2 files | Compression toggle |
| `Async` | `Async` | ~3 files | Async execution flag |
| `Trigger` | `Trigger` | ~2 files | Trigger source identifier |
| `MasterSnapshotId` | `MasterSnapshotId` | ~3 files | SchedulerTriggerTrait, incremental parent ref |
| `MasterDir` | `MasterDir` | ~2 files | Master snapshot directory |
| `Confirm` | `Confirm` | ~2 files | Restore confirmation flag |
| `CreateBackup` | `CreateBackup` | ~2 files | Pre-restore backup creation |
| `RequireBackup` | `RequireBackup` | ~2 files | Backup requirement flag |
| `Mode` | `Mode` | ~5 files | Restore mode, snapshot mode, SnapshotCrudRestoreTrait |

---

## Scheduler Response Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `Frequency` | `Frequency` | ~3 files | SchedulerTriggerTrait, schedule settings |
| `Time` | `Time` | ~2 files | Schedule time setting |
| `Day` | `Day` | ~2 files | Schedule day setting |
| `Scheduled` | `Scheduled` | ~2 files | Scheduling confirmation flag |
| `Trace` | `Trace` | ~2 files | Error trace payloads |
| `Options` | `Options` | ~3 files | Snapshot/restore option payloads |

---

## Storage Stats Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `TotalSizeFormatted` | `TotalSizeFormatted` | ~2 files | CleanerStorageTrait, formatted size string |
| `OldestTimestamp` | `OldestTimestamp` | ~2 files | CleanerStorageTrait, oldest snapshot date |
| `NewestTimestamp` | `NewestTimestamp` | ~2 files | CleanerStorageTrait, newest snapshot date |
| `DiskFreeBytes` | `DiskFreeBytes` | ~2 files | CleanerStorageTrait, disk space |
| `DiskFreeFormatted` | `DiskFreeFormatted` | ~2 files | CleanerStorageTrait, formatted disk space |
| `SnapshotsCount` | `SnapshotsCount` | ~2 files | CleanerStorageTrait, snapshot count |
| `BytesFormatted` | `BytesFormatted` | ~2 files | Formatted byte strings |

---

## Progress Envelope Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `IsSuccess` | `IsSuccess` | ~2 files | Progress/result envelope wrappers |
| `HasAnyErrors` | `HasAnyErrors` | ~2 files | Error presence flag in batch results |

---

## Cleanup Detail Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `Details` | `Details` | ~3 files | CleanerRetentionTrait, deletion detail arrays |
| `Order` | `Order` | ~2 files | Sort/sequence ordering |

---

## Internal Passing Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `Graph` | `Graph` | ~2 files | AnalyzerQueryTrait, dependency graph |
| `InDegree` | `InDegree` | ~2 files | AnalyzerQueryTrait, topological sort |
| `Manifest` | `Manifest` | ~3 files | Export/import manifest payloads |
| `SqlitePath` | `SqlitePath` | ~3 files | SQLite database file path |
| `RealPath` | `RealPath` | ~2 files | Resolved filesystem path |
| `FilePath` | `FilePath` | ~4 files | NativeSnapshotExecTrait, SnapshotExportHandlerTrait, file references |
| `PkColumn` | `PkColumn` | ~2 files | Primary key column name |
| `TableName` | `TableName` | ~3 files | Table name references in export/import |

---

## Provider & Plugin Info Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `Id` | `Id` | ~5 files | DetectorSettingsTrait, CleanerRetentionTrait, provider identification |
| `Name` | `Name` | ~6 files | SnapshotProviderWpReset, UpdraftCrudTrait, OrchestratorPluginTrait, provider info |
| `Available` | `Available` | ~4 files | DetectorSettingsTrait, DetectorProviderTrait, provider availability |
| `Capabilities` | `Capabilities` | ~3 files | DetectorProviderTrait, provider capability reporting |
| `Version` | `Version` | ~5 files | OrchestratorPluginTrait, plugin version info, provider versioning |
| `Author` | `Author` | ~2 files | Plugin metadata responses |
| `Description` | `Description` | ~2 files | Plugin/provider description |
| `Active` | `Active` | ~3 files | Plugin active status |
| `TotalFiles` | `TotalFiles` | ~2 files | File count in exports |
| `LastSeenId` | `LastSeenId` | ~2 files | Incremental tracking cursor |
| `FileType` | `FileType` | ~2 files | File type classification |
| `Provider` | `Provider` | ~3 files | Provider identifier in responses |
| `Snapshot` | `Snapshot` | ~3 files | Snapshot object in responses |
| `Source` | `Source` | ~2 files | Import source reference |

---

## Capability Sub-Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `FullSite` | `FullSite` | ~2 files | DetectorProviderTrait, capability flag |
| `DatabaseOnly` | `DatabaseOnly` | ~2 files | DetectorProviderTrait, capability flag |
| `Selective` | `Selective` | ~2 files | DetectorProviderTrait, capability flag |
| `Restore` | `Restore` | ~3 files | DetectorProviderTrait, restore capability |
| `Import` | `Import` | ~2 files | DetectorProviderTrait, import capability |

---

## Restore Option Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `Strict` | `Strict` | ~2 files | Strict restore mode flag |
| `ApplyIncrementals` | `ApplyIncrementals` | ~2 files | Incremental application flag |
| `Sqlite` | `Sqlite` | ~2 files | SQLite restore mode |
| `SqliteFile` | `SqliteFile` | ~3 files | SQLite file reference in restore |
| `InternalMode` | `_Mode` | ~2 files | Internal mode discriminator (underscore-prefixed) |

---

## OPcache Status Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `OpcacheAvailable` | `OpcacheAvailable` | ~2 files | StatusOpsTrait, opcache detection |
| `OpcacheReset` | `OpcacheReset` | ~2 files | StatusOpsTrait, opcache reset result |
| `FilesInvalidated` | `FilesInvalidated` | ~2 files | StatusOpsTrait, invalidated file count |

---

## Plugin Archive Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `Zip` | `Zip` | ~2 files | OrchestratorPluginTrait, zip filename |
| `ZipFile` | `ZipFile` | ~2 files | OrchestratorPluginTrait, root DB zip path |
| `FileSizeBytes` | `FileSizeBytes` | ~2 files | OrchestratorPluginTrait, root DB file size |
| `ChecksumMd5` | `ChecksumMd5` | ~2 files | OrchestratorPluginTrait, root DB checksum |
| `PluginName` | `PluginName` | ~2 files | OrchestratorPluginTrait, root DB plugin name |

---

## Status Payload Keys

| Case | Value | Usages | Primary Locations |
|------|-------|--------|-------------------|
| `Route` | `Route` | ~2 files | StatusPayloadTrait, registered route list |
| `Methods` | `Methods` | ~2 files | StatusPayloadTrait, HTTP methods per route |
| `Result` | `Result` | ~2 files | AdminAjaxSnapshotTrait, result wrapper key |

---

## Permitted Magic String Exceptions

The following array key patterns intentionally remain as literal strings:

| Pattern | Reason |
|---------|--------|
| `$snapshot['Filename']`, `$snapshot['Filepath']` | Database column reads (PascalCase schema contract) |
| `$error['message']`, `$error['type']` | PHP native `error_get_last()` structure |
| `$upload['error']` | PHP `$_FILES` superglobal structure |
| `$uploadDir['error']` | WordPress `wp_upload_dir()` return structure |
| `$body['scope']`, `$body['tables']` | Incoming REST request body keys |
| `$row['count']`, `$total['count']` | SQL alias reads (`COUNT(*) AS count`) |
| `$dbStats['total_size']`, `$dbStats['oldest']` | SQL alias reads (`SUM(...) AS total_size`) |
| `$result['max_seq']` | SQL alias reads (`MAX(Sequence) AS max_seq`) |
| `$info['Name']`, `$info['Rows']` | MySQL `SHOW TABLE STATUS` column names |
| `$cached['modified_at']`, `$cached['md5_hash']` | SQLite cache table column reads |
| `$where['sql']`, `$filter['params']` | Internal query builder structures |
| `$r['label']`, `$r['file']` | Diagnostic/boot loader internal arrays |
| `$result['failed']`, `$result['loaded']` | Autoloader diagnostics (pre-namespace context) |
| `$result['term_id']` | WordPress `wp_insert_term()` return structure |
| `$result['file']` | SQLite `PRAGMA database_list` return structure |
| `$envelope['Errors']` | PascalCase Go-compatible envelope keys |
| `'methods'`, `'callback'`, `'permission_callback'` | WordPress REST API `register_rest_route()` contract |
| i18n translation domain `'riseup-asia-uploader'` | WordPress gettext extraction requirement |
| Log context keys (`'filepath'`, `'table'`, `'id'`) | Exempt per coding standard — camelCase log keys |

---

## Helper Methods

```php
// Type-safe comparison (preferred over === with ->value)
$key->isEqual(ResponseKeyType::Success);

// Positive boolean logic (P3 — eliminates raw negation)
$key->isOtherThan(ResponseKeyType::Success);

// Multi-match check
$key->isAnyOf(ResponseKeyType::Success, ResponseKeyType::Data);
```

---

## Cross-Language Sync

| Language | Location | Sync Status |
|----------|----------|-------------|
| PHP | `includes/Enums/ResponseKeyType.php` | **Source of truth** (176 cases) |
| Go | `backend/internal/enums/response_key/variant.go` | Pending — must be updated to match 176 PascalCase values |
| TypeScript | `src/lib/constants.ts` → `ResponseKeyType` | Pending — must be updated to match 176 PascalCase values |

All values use **PascalCase** (e.g., `'Success'`, `'SnapshotId'`, `'DeletedByPolicy'`, `'FilePath'`).
