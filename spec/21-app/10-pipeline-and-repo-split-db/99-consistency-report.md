# Consistency Report: Pipeline and Repo Split Databases

## Verification Summary

- **Database Standard Compliance:** 100%
  - Strict PascalCase singular table names (`PipelineRun`, `PipelineErrorLog`, `PipelineSegment`, `RepoScanLog`).
  - Integer primary key autoincrement naming (`PipelineRunId`, `PipelineErrorLogId`, `PipelineSegmentId`, `Id`).
  - Mandatory context columns (`Notes TEXT NULL`, `Comments TEXT NULL`) present on all event tables.
  - Connection pooling limited to `SetMaxOpenConns(1)` to eliminate lock contention.
  - Binary execution anchoring via `filepath.EvalSymlinks(os.Executable())`.

- **Command Architecture Parity:** 100%
  - `gitmap pipeline db <status|clear|reset|optimize|errorlogs|help>`
  - `gitmap repo db <status|log|errorlogs|clear|reset|optimize|help>`
  - `gitmap db <status|optimize|ls|repo-db|sizes|reset>`

- **Path Cleanliness:** 100%
  - No drive-letter URIs or hardcoded absolute filesystem paths.
