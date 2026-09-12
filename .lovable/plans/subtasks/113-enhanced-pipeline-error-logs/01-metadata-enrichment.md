# Subtask 01: Metadata Enrichment for Pipeline Error Logs

## Objective
Enrich `PipelineErrorLogsPayload` with repository metadata:
- `RepoUrl` (e.g., `https://github.com/alimtvnetwork/gitmap-v28` or remote URL)
- `LastHash` (short commit SHA e.g., `24050fe`)
- `LastReleaseVersion` (e.g., `v6.220.2`)
- `LatestBranch` (active branch or latest run head branch)
- `OpenPRsCount` (number of currently open PRs)

## Implementation Steps
1. Add fields to `PipelineErrorLogsPayload` in `gitmap/cmd/pipeline.go`.
2. Extract metadata during `buildErrorLogsPayload` or `processAndRenderErrorLogs` in `gitmap/cmd/pipeline_logs.go`.
3. Adhere to function size <= 15 lines, affirmative booleans, and blank line before return.
